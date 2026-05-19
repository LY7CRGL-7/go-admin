package middleware

import (
	"admin/internal/conf"
	"admin/internal/data/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RBAC 基于角色的权限控制中间件
func RBAC(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID, exists := c.Get("admin_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "未获取到管理员信息",
			})
			c.Abort()
			return
		}

		// 获取请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 检查是否是超级管理员（拥有所有权限）
		var admin model.Admin
		if err := db.Preload("Roles").First(&admin, adminID).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "管理员不存在",
			})
			c.Abort()
			return
		}

		// 检查管理员状态
		if admin.Status != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "管理员账号已被禁用",
			})
			c.Abort()
			return
		}

		// 检查是否为超级管理员
		for _, role := range admin.Roles {
			if role.Code == "super_admin" && role.Status == 1 {
				c.Next()
				return
			}
		}

		// 查询权限
		hasPermission := CheckPermission(db, adminID, path, method)
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "没有权限访问该资源",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckPermission 检查权限
func CheckPermission(db *gorm.DB, adminID interface{}, path, method string) bool {
	var count int64
	
	query := `
		SELECT COUNT(*)
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN admin_roles ar ON rp.role_id = ar.role_id
		WHERE ar.admin_id = ? 
			AND p.status = 1
			AND p.path = ?
			AND p.method = ?
	`
	
	if err := db.Raw(query, adminID, path, method).Count(&count).Error; err != nil {
		return false
	}
	
	return count > 0
}

// RequirePermission 权限检查装饰器（用于 Handler 内部细粒度权限控制）
func RequirePermission(db *gorm.DB, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID, exists := c.Get("admin_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "未获取到管理员信息",
			})
			c.Abort()
			return
		}

		var count int64
		query := `
			SELECT COUNT(*)
			FROM permissions p
			JOIN role_permissions rp ON p.id = rp.permission_id
			JOIN admin_roles ar ON rp.role_id = ar.role_id
			WHERE ar.admin_id = ? 
				AND p.status = 1
				AND p.code = ?
		`
		
		if err := db.Raw(query, adminID, permissionCode).Count(&count).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 500,
				"msg":  "权限检查失败",
			})
			c.Abort()
			return
		}

		if count == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "没有执行该操作的权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
