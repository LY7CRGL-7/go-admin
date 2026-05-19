package handler

import (
	"admin/internal/data/model"
	"admin/internal/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListRoles 获取角色列表
func ListRoles(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	var roles []model.Role
	db.Order("sort ASC, created_at DESC").Find(&roles)
	
	roleResponses := make([]dto.RoleResponse, 0, len(roles))
	for _, role := range roles {
		roleResponses = append(roleResponses, dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Status:      role.Status,
			Sort:        role.Sort,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": roleResponses,
	})
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}
	
	// 检查角色代码是否已存在
	var count int64
	db.Model(&model.Role{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "角色代码已存在",
		})
		return
	}
	
	role := model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      1,
	}
	
	if err := db.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建角色失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": gin.H{"id": role.ID},
	})
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的角色 ID",
		})
		return
	}
	
	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
		})
		return
	}
	
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "角色不存在",
		})
		return
	}
	
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	
	if len(updates) > 0 {
		db.Model(&role).Updates(updates)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新成功",
	})
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的角色 ID",
		})
		return
	}
	
	tx := db.Begin()
	// 删除权限关联
	tx.Where("role_id = ?", id).Delete(&model.RolePermission{})
	// 删除管理员关联
	tx.Where("role_id = ?", id).Delete(&model.AdminRole{})
	// 删除角色
	tx.Delete(&model.Role{}, id)
	tx.Commit()
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// AssignPermissions 分配权限
func AssignPermissions(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的角色 ID",
		})
		return
	}
	
	var req dto.AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
		})
		return
	}
	
	tx := db.Begin()
	// 删除旧权限
	tx.Where("role_id = ?", id).Delete(&model.RolePermission{})
	// 添加新权限
	for _, permID := range req.PermissionIDs {
		tx.Create(&model.RolePermission{
			RoleID:       uint(id),
			PermissionID: permID,
		})
	}
	tx.Commit()
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "分配成功",
	})
}

// ListPermissions 获取权限列表
func ListPermissions(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	var permissions []model.Permission
	db.Order("sort ASC, created_at DESC").Find(&permissions)
	
	permResponses := make([]dto.PermissionResponse, 0, len(permissions))
	for _, perm := range permissions {
		permResponses = append(permResponses, dto.PermissionResponse{
			ID:       perm.ID,
			Name:     perm.Name,
			Code:     perm.Code,
			Type:     perm.Type,
			Path:     perm.Path,
			Method:   perm.Method,
			ParentID: perm.ParentID,
			Sort:     perm.Sort,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": permResponses,
	})
}

// ListAuditLogs 获取审计日志列表
func ListAuditLogs(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	var req dto.PaginatedRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
		})
		return
	}
	
	var logs []model.AuditLog
	var total int64
	
	db.Model(&model.AuditLog{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("created_at DESC").
		Find(&logs)
	
	logResponses := make([]dto.AuditLogResponse, 0, len(logs))
	for _, log := range logs {
		logResponses = append(logResponses, dto.AuditLogResponse{
			ID:        log.ID,
			AdminID:   log.AdminID,
			AdminName: log.AdminName,
			Action:    log.Action,
			Resource:  log.Resource,
			Method:    log.Method,
			Path:      log.Path,
			IP:        log.IP,
			UserAgent: log.UserAgent,
			Status:    log.Status,
			Duration:  log.Duration,
			CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"total": total,
			"list":  logResponses,
		},
	})
}
