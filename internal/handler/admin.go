package handler

import (
	"admin/internal/conf"
	"admin/internal/data/model"
	"admin/internal/dto"
	"admin/internal/pkg/logger"
	"admin/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	db  *gorm.DB
	cfg *conf.Config
}

func NewAdminHandler(db *gorm.DB, cfg *conf.Config) *AdminHandler {
	return &AdminHandler{
		db:  db,
		cfg: cfg,
	}
}

// ListAdmins 获取管理员列表
func (h *AdminHandler) ListAdmins(c *gin.Context) {
	var req dto.PaginatedRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
		})
		return
	}
	
	var admins []model.Admin
	var total int64
	
	h.db.Model(&model.Admin{}).Count(&total)
	h.db.Preload("Roles").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("created_at DESC").
		Find(&admins)
	
	adminInfos := make([]dto.AdminInfo, 0, len(admins))
	for _, admin := range admins {
		info := dto.AdminInfo{
			ID:       admin.ID,
			Username: admin.Username,
			Nickname: admin.Nickname,
			Email:    admin.Email,
			Phone:    admin.Phone,
			Avatar:   admin.Avatar,
			Status:   admin.Status,
		}
		if admin.LastLoginAt != nil {
			timeStr := admin.LastLoginAt.Format("2006-01-02 15:04:05")
			info.LastLoginAt = &timeStr
		}
		adminInfos = append(adminInfos, info)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": dto.AdminListResponse{
			Total: total,
			List:  adminInfos,
		},
	})
}

// CreateAdmin 创建管理员
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req dto.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}
	
	// 检查用户名是否已存在
	var count int64
	h.db.Model(&model.Admin{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}
	
	// 加密密码
	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		logger.Error("加密密码失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统错误",
		})
		return
	}
	
	// 创建管理员
	admin := model.Admin{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}
	
	tx := h.db.Begin()
	if err := tx.Create(&admin).Error; err != nil {
		tx.Rollback()
		logger.Error("创建管理员失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建管理员失败",
		})
		return
	}
	
	// 分配角色
	if len(req.RoleIDs) > 0 {
		for _, roleID := range req.RoleIDs {
			tx.Create(&model.AdminRole{
				AdminID: admin.ID,
				RoleID:  roleID,
			})
		}
	}
	
	tx.Commit()
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": gin.H{"id": admin.ID},
	})
}

// UpdateAdmin 更新管理员
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的管理员 ID",
		})
		return
	}
	
	var req dto.UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
		})
		return
	}
	
	var admin model.Admin
	if err := h.db.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}
	
	// 更新字段
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	
	tx := h.db.Begin()
	if len(updates) > 0 {
		tx.Model(&admin).Updates(updates)
	}
	
	// 更新角色
	if len(req.RoleIDs) > 0 {
		// 删除旧角色
		tx.Where("admin_id = ?", admin.ID).Delete(&model.AdminRole{})
		// 添加新角色
		for _, roleID := range req.RoleIDs {
			tx.Create(&model.AdminRole{
				AdminID: admin.ID,
				RoleID:  roleID,
			})
		}
	}
	
	tx.Commit()
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新成功",
	})
}

// DeleteAdmin 删除管理员
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的管理员 ID",
		})
		return
	}
	
	// 不能删除自己
	adminID, _ := c.Get("admin_id")
	if uint(id) == adminID.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "不能删除当前登录的管理员账号",
		})
		return
	}
	
	tx := h.db.Begin()
	// 删除角色关联
	tx.Where("admin_id = ?", id).Delete(&model.AdminRole{})
	// 删除管理员
	tx.Delete(&model.Admin{}, id)
	tx.Commit()
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}
