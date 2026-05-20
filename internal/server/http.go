package server

import (
	"time"

	"admin/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer 创建 HTTP 服务器（健康检查 + 可选 REST API）
func NewHTTPServer(c *conf.Server, logger log.Logger) *http.Server {
	var opts []http.ServerOption

	opts = append(opts, http.Middleware(
		recovery.Recovery(),
		logging.Server(logger),
	))

	if c.HTTP.Addr != "" {
		opts = append(opts, http.Address(c.HTTP.Addr))
	}
	if c.HTTP.Timeout > 0 {
		opts = append(opts, http.Timeout(time.Duration(c.HTTP.Timeout)))
	}

	srv := http.NewServer(opts...)

	// 注册健康检查路由
	route := srv.Route("/")
	route.GET("/health", func(ctx http.Context) error {
		return ctx.JSON(200, map[string]string{"status": "ok"})
	})
	route.GET("/ready", func(ctx http.Context) error {
		return ctx.JSON(200, map[string]string{"status": "ready"})
	})

	return srv
}
