package server

import (
	"admin/internal/conf"
	"admin/internal/handler"
	"admin/internal/middleware"
	"admin/internal/service"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewRouter 创建路由
func NewRouter(cfg *conf.Config, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)
	
	r := gin.New()
	
	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(CORS(cfg.Security.CORS))
	
	// 将 db 添加到 context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
		})
	})
	
	// 创建服务和处理器
	authService := service.NewAuthService(db, cfg)
	authHandler := handler.NewAuthHandler(authService)
	adminHandler := handler.NewAdminHandler(db, cfg)
	
	// 初始化管理员账号
	if err := authService.InitAdmin(); err != nil {
		fmt.Printf("初始化管理员账号失败: %v\n", err)
	}
	
	// 公开路由（无需认证）
	public := r.Group("/api/v1")
	public.Use(middleware.LoginLimit(rdb, &cfg.Security.Login))
	{
		public.POST("/auth/login", authHandler.Login)
	}
	
	// 需要认证的路由
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuth(&cfg.JWT))
	auth.Use(middleware.RateLimiter(rdb, &cfg.RateLimit))
	auth.Use(middleware.AuditLogger(db))
	{
		auth.POST("/auth/logout", authHandler.Logout)
		auth.POST("/auth/change-password", authHandler.ChangePassword)
		auth.GET("/auth/profile", authHandler.GetProfile)
		
		// 管理员管理
		admins := auth.Group("/admins")
		{
			admins.GET("", adminHandler.ListAdmins)
			admins.POST("", adminHandler.CreateAdmin)
			admins.PUT("/:id", adminHandler.UpdateAdmin)
			admins.DELETE("/:id", adminHandler.DeleteAdmin)
		}
		
		// 角色管理
		roles := auth.Group("/roles")
		{
			roles.GET("", handler.ListRoles)
			roles.POST("", handler.CreateRole)
			roles.PUT("/:id", handler.UpdateRole)
			roles.DELETE("/:id", handler.DeleteRole)
			roles.POST("/:id/permissions", handler.AssignPermissions)
		}
		
		// 权限管理
		permissions := auth.Group("/permissions")
		{
			permissions.GET("", handler.ListPermissions)
		}
		
		// 审计日志
		auditLogs := auth.Group("/audit-logs")
		{
			auditLogs.GET("", handler.ListAuditLogs)
		}
	}
	
	return r
}

// CORS 跨域中间件
func CORS(cfg conf.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 检查是否在允许的源列表中
		allowed := false
		for _, allowedOrigin := range cfg.AllowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,Accept,Origin,X-Requested-With")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.MaxAge))
		}
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// StartServer 启动服务器
func StartServer(cfg *conf.Config, db *gorm.DB, rdb *redis.Client) error {
	router := NewRouter(cfg, db, rdb)
	
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	
	s := &HTTPServer{
		addr: addr,
		router: router,
		readTimeout: cfg.Server.ReadTimeout,
		writeTimeout: cfg.Server.WriteTimeout,
	}
	
	return s.Start()
}

// HTTPServer HTTP 服务器
type HTTPServer struct {
	addr         string
	router       *gin.Engine
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// Start 启动 HTTP 服务器
func (s *HTTPServer) Start() error {
	return s.router.Run(s.addr)
}
