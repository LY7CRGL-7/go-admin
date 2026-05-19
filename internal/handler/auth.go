package handler

import (
	"admin/internal/dto"
	"admin/internal/middleware"
	"admin/internal/pkg/logger"
	"admin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 管理员登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 获取 IP
	ip := middleware.GetClientIP(c)

	// 登录
	token, admin, err := h.authService.Login(req.Username, req.Password, ip)
	if err != nil {
		logger.Info("登录失败", "username", req.Username, "ip", ip, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  err.Error(),
		})
		return
	}

	// 构造响应
	adminInfo := &dto.AdminInfo{
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
		adminInfo.LastLoginAt = &timeStr
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"data": dto.LoginResponse{
			Token: token,
			Admin: adminInfo,
		},
	})
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
		})
		return
	}

	adminID, _ := c.Get("admin_id")

	if err := h.authService.ChangePassword(adminID.(uint), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "密码修改成功",
	})
}

// GetProfile 获取当前管理员信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	admin, err := h.authService.GetAdminByID(adminID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取管理员信息失败",
		})
		return
	}

	adminInfo := &dto.AdminInfo{
		ID:       admin.ID,
		Username: admin.Username,
		Nickname: admin.Nickname,
		Email:    admin.Email,
		Phone:    admin.Phone,
		Avatar:   admin.Avatar,
		Status:   admin.Status,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": adminInfo,
	})
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT 是无状态的，客户端删除 Token 即可
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登出成功",
	})
}
