package main

import (
	"admin/internal/conf"
	"admin/internal/data"
	"admin/internal/grpc"
	"admin/internal/pkg/logger"
	grpcserver "admin/internal/server"
	"admin/internal/service"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 读取配置文件
	cfg := &conf.Config{}
	configData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(configData, cfg); err != nil {
		fmt.Printf("解析配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log.FilePath, cfg.Log.Level, cfg.Log.MaxSize, cfg.Log.MaxBackups, cfg.Log.MaxAge); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	logger.Info("启动管理端服务", "version", "1.0.0")

	// 初始化数据库
	db, err := data.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Error("初始化数据库失败", "error", err)
		os.Exit(1)
	}

	// 启动 gRPC 服务器
	if cfg.GRPC.Enabled {
		grpcLogger, _ := zap.NewProduction()
		defer grpcLogger.Sync()

		grpcServer, err := grpcserver.NewGRPCServer(cfg, grpcLogger)
		if err != nil {
			logger.Error("初始化 gRPC 服务器失败", "error", err)
			os.Exit(1)
		}

		// 创建 Service 层
		adminRepo := data.NewAdminRepo(db)
		roleRepo := data.NewRoleRepo(db)
		permissionRepo := data.NewPermissionRepo(db)
		auditLogRepo := data.NewAuditLogRepo(db)

		authService := grpc.NewAuthService(
			service.NewAuthService(adminRepo, cfg),
			grpcLogger,
		)
		adminService := grpc.NewAdminService(
			service.NewAdminService(adminRepo, roleRepo),
			grpcLogger,
		)
		roleService := grpc.NewRoleService(
			service.NewRoleService(roleRepo, permissionRepo),
			grpcLogger,
		)
		permissionService := grpc.NewPermissionService(
			service.NewPermissionService(permissionRepo),
			grpcLogger,
		)
		auditLogService := grpc.NewAuditLogService(
			service.NewAuditLogService(auditLogRepo),
			grpcLogger,
		)

		// 注册 gRPC 服务
		grpcServer.RegisterServices(
			authService,
			adminService,
			roleService,
			permissionService,
			auditLogService,
		)

		// 在后台启动 gRPC 服务器
		go func() {
			if err := grpcServer.Start(); err != nil {
				logger.Error("gRPC 服务器启动失败", "error", err)
			}
		}()

		logger.Info("gRPC 服务器已启动", "port", cfg.GRPC.Port)
	} else {
		logger.Warn("gRPC 服务器未启用")
		os.Exit(1)
	}

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")
	logger.Info("服务器已退出")
}
