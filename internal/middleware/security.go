package middleware

import (
	"admin/internal/conf"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

// LoginLimit 登录限制中间件（防止暴力破解）
func LoginLimit(rdb *redis.Client, cfg *conf.LoginConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 IP
		ip := GetClientIP(c)
		username := c.PostForm("username")
		
		ctx := context.Background()
		
		// 检查 IP 是否被锁定
		ipKey := fmt.Sprintf("admin:login:lock:ip:%s", ip)
		if rdb.Exists(ctx, ipKey).Val() > 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "登录尝试过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		
		// 检查用户名是否被锁定
		if username != "" {
			userKey := fmt.Sprintf("admin:login:lock:user:%s", username)
			if rdb.Exists(ctx, userKey).Val() > 0 {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"code": 429,
					"msg":  "该账号登录失败次数过多，已被临时锁定",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

// CheckLoginAttempts 检查并记录登录失败次数
func CheckLoginAttempts(rdb *redis.Client, cfg *conf.LoginConfig, username, ip string, success bool) error {
	ctx := context.Background()
	userKey := fmt.Sprintf("admin:login:attempts:user:%s", username)
	ipKey := fmt.Sprintf("admin:login:attempts:ip:%s", ip)
	
	if success {
		// 登录成功，清除失败记录
		rdb.Del(ctx, userKey, ipKey)
		return nil
	}
	
	// 登录失败，增加计数
	userAttempts := rdb.Incr(ctx, userKey).Val()
	ipAttempts := rdb.Incr(ctx, ipKey).Val()
	
	// 设置过期时间
	if userAttempts == 1 {
		rdb.Expire(ctx, userKey, cfg.AttemptWindow)
	}
	if ipAttempts == 1 {
		rdb.Expire(ctx, ipKey, cfg.AttemptWindow)
	}
	
	// 检查是否超过限制
	if userAttempts >= int64(cfg.MaxAttempts) {
		// 锁定用户
		lockKey := fmt.Sprintf("admin:login:lock:user:%s", username)
		rdb.Set(ctx, lockKey, "1", cfg.LockoutDuration)
	}
	
	if ipAttempts >= int64(cfg.MaxAttempts*2) {
		// 锁定 IP（阈值是用户的2倍）
		lockKey := fmt.Sprintf("admin:login:lock:ip:%s", ip)
		rdb.Set(ctx, lockKey, "1", cfg.LockoutDuration)
	}
	
	return nil
}

// RateLimiter 限流中间件
func RateLimiter(rdb *redis.Client, cfg *conf.RateLimitConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	
	return func(c *gin.Context) {
		ctx := context.Background()
		
		// 全局限流
		globalKey := "admin:ratelimit:global"
		if err := checkRateLimit(ctx, rdb, globalKey, cfg.GlobalQPS); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "系统繁忙，请稍后再试",
			})
			c.Abort()
			return
		}
		
		// IP 限流
		ip := GetClientIP(c)
		ipKey := fmt.Sprintf("admin:ratelimit:ip:%s", ip)
		if err := checkRateLimit(ctx, rdb, ipKey, cfg.PerIPQPS); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		
		// 用户限流（如果已登录）
		adminID, exists := c.Get("admin_id")
		if exists {
			userKey := fmt.Sprintf("admin:ratelimit:user:%v", adminID)
			if err := checkRateLimit(ctx, rdb, userKey, cfg.PerUserQPS); err != nil {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"code": 429,
					"msg":  "您的操作过于频繁，请稍后再试",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

// checkRateLimit 检查限流（滑动窗口）
func checkRateLimit(ctx context.Context, rdb *redis.Client, key string, limit int) error {
	now := fmt.Sprintf("%d", redis.Now().UnixMilli())
	window := 1000 // 1秒窗口
	
	// 使用 Redis Pipeline 实现滑动窗口限流
	pipe := rdb.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", now)
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(redis.Now().UnixMilli()), Member: now})
	pipe.Expire(ctx, key, 2) // 2秒过期
	cmd, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	
	// 获取当前窗口内的请求数
	count := rdb.ZCard(ctx, key).Val()
	if count > int64(limit) {
		return fmt.Errorf("rate limit exceeded")
	}
	
	return nil
}

// IPWhitelist IP 白名单中间件
func IPWhitelist(whitelist []string) gin.HandlerFunc {
	if len(whitelist) == 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	
	// 构建白名单集合
	whitelistMap := make(map[string]bool)
	for _, ip := range whitelist {
		whitelistMap[ip] = true
	}
	
	return func(c *gin.Context) {
		ip := GetClientIP(c)
		
		if !whitelistMap[ip] {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "您的 IP 不在白名单中",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// GetClientIP 获取客户端真实 IP
func GetClientIP(c *gin.Context) string {
	// 尝试从 X-Forwarded-For 获取
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}
	
	// 尝试从 X-Real-IP 获取
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		if net.ParseIP(xRealIP) != nil {
			return xRealIP
		}
	}
	
	// 尝试从 RemoteAddr 获取
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err == nil && net.ParseIP(ip) != nil {
		return ip
	}
	
	return c.ClientIP()
}
