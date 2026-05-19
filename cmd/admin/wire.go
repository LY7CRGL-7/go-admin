//go:build wireinject
// +build wireinject

package main

import (
	"admin/internal/conf"
	"admin/internal/data"
	"admin/internal/handler"
	"admin/internal/server"
	"admin/internal/service"

	"github.com/google/wire"
)

// InitializeApp 初始化应用
func InitializeApp(cfg *conf.Config) (*server.HTTPServer, error) {
	wire.Build(
		// 数据层
		data.NewDatabase,
		data.NewRedis,
		data.NewAdminRepo,
		data.NewRoleRepo,
		data.NewPermissionRepo,
		data.NewAuditLogRepo,

		// 服务层
		service.NewAuthService,
		service.NewAdminService,
		service.NewRoleService,
		service.NewPermissionService,
		service.NewAuditLogService,

		// 处理器层
		handler.NewAuthHandler,
		handler.NewAdminHandler,
		handler.NewRoleHandler,
		handler.NewPermissionHandler,
		handler.NewAuditLogHandler,

		// 服务器
		server.NewHTTPServer,
	)

	return &server.HTTPServer{}, nil
}
