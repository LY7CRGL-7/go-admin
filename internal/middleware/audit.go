package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"admin/internal/data/model"
	"admin/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditLogger 审计日志中间件
func AuditLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()
		
		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// 重新设置 Body
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}
		
		// 处理请求
		c.Next()
		
		// 记录审计日志
		duration := time.Since(startTime).Milliseconds()
		
		// 获取管理员信息
		adminID := uint(0)
		adminName := ""
		if id, exists := c.Get("admin_id"); exists {
			if uid, ok := id.(uint); ok {
				adminID = uid
			}
		}
		if name, exists := c.Get("username"); exists {
			if uname, ok := name.(string); ok {
				adminName = uname
			}
		}
		
		// 获取响应体
		var responseBody string
		if writer, ok := c.Writer.(*responseWriter); ok {
			responseBody = string(writer.body.Bytes())
		}
		
		// 创建审计日志
		auditLog := model.AuditLog{
			AdminID:   adminID,
			AdminName: adminName,
			Action:    getAction(c.Request.URL.Path, c.Request.Method),
			Resource:  getResource(c.Request.URL.Path),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			IP:        GetClientIP(c),
			UserAgent: c.Request.UserAgent(),
			Request:   truncateString(requestBody, 2000),
			Response:  truncateString(responseBody, 2000),
			Status:    c.Writer.Status(),
			Duration:  duration,
			CreatedAt: startTime,
		}
		
		// 异步保存审计日志
		go func() {
			if err := db.Create(&auditLog).Error; err != nil {
				logger.Error("保存审计日志失败", "error", err)
			}
		}()
	}
}

// responseWriter 自定义 ResponseWriter 用于捕获响应体
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// getAction 从路径和操作中提取动作
func getAction(path, method string) string {
	actions := map[string]string{
		"POST":   "创建",
		"PUT":    "更新",
		"PATCH":  "更新",
		"DELETE": "删除",
		"GET":    "查询",
	}
	
	if action, ok := actions[method]; ok {
		return action
	}
	return method
}

// getResource 从路径中提取资源
func getResource(path string) string {
	// 例如: /api/v1/admins/123 -> admins
	// 简化实现，实际可以更复杂
	parts := splitPath(path)
	if len(parts) >= 3 {
		return parts[2] // /api/v1/{resource}/...
	}
	return path
}

func splitPath(path string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if start < i {
				parts = append(parts, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		parts = append(parts, path[start:])
	}
	return parts
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
