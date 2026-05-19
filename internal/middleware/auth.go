package middleware

import (
	"net/http"
	"strings"
	"time"

	"admin/internal/conf"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT Claims
type Claims struct {
	AdminID  uint   `json:"admin_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTAuth JWT 认证中间件
func JWTAuth(cfg *conf.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// Bearer Token 格式验证
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "认证令牌格式错误",
			})
			c.Abort()
			return
		}

		// 解析 Token
		token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "认证令牌无效或已过期",
			})
			c.Abort()
			return
		}

		// 验证 Claims
		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "认证令牌无效",
			})
			c.Abort()
			return
		}

		// 验证 Issuer
		if claims.Issuer != cfg.Issuer {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "认证令牌来源不正确",
			})
			c.Abort()
			return
		}

		// 将用户信息存入 Context
		c.Set("admin_id", claims.AdminID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// GenerateToken 生成 JWT Token
func GenerateToken(cfg *conf.JWTConfig, adminID uint, username string) (string, error) {
	claims := Claims{
		AdminID:  adminID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}
