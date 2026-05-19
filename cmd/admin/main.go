package main

import (
	"admin/internal/conf"
	"admin/internal/data"
	"admin/internal/pkg/logger"
	"admin/internal/server"
	"flag"
	"fmt"
	"os"

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

	// 初始化 Redis
	rdb, err := data.NewRedis(&cfg.Redis)
	if err != nil {
		logger.Error("初始化 Redis 失败", "error", err)
		os.Exit(1)
	}

	// 启动 gRPC 服务器
	var grpcServer *server.GRPCServer
	if cfg.GRPC.Enabled {
		grpcLogger := logger.GetLogger()
		grpcServer, err = server.NewGRPCServer(cfg, grpcLogger)
		if err != nil {
			logger.Error("初始化 gRPC 服务器失败", "error", err)
			os.Exit(1)
		}

		// TODO: 注册 gRPC 服务（需要先实现 service 层）
		// authService := grpcservice.NewAuthService(authService, grpcLogger)
		// grpcServer.RegisterServices(
		// 	authService,
		// 	nil, // adminService
		// 	nil, // roleService
		// 	nil, // permissionService
		// 	nil, // auditLogService
		// )

		// 在后台启动 gRPC 服务器
		go func() {
			if err := grpcServer.Start(); err != nil {
				logger.Error("gRPC 服务器启动失败", "error", err)
			}
		}()

		logger.Info("gRPC 服务器启动中", "port", cfg.GRPC.Port)
	}

	// 启动 HTTP 服务器
	logger.Info("HTTP 服务器启动中", "port", cfg.Server.Port)
	if err := server.StartServer(cfg, db, rdb); err != nil {
		logger.Error("启动 HTTP 服务器失败", "error", err)
		os.Exit(1)
	}
}
