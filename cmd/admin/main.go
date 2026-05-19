package main

import (
	"admin/internal/conf"
	"admin/internal/data"
	grpcservice "admin/internal/grpc"
	"admin/internal/pkg/logger"
	"admin/internal/server"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 读取配置文件
	cfg := &conf.Config{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		fmt.Printf("解析配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	logger, err := logger.NewLogger(&cfg.Logger)
	if err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

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
		grpcLogger, _ := zap.NewProduction()
		grpcServer, err = server.NewGRPCServer(cfg, grpcLogger)
		if err != nil {
			logger.Error("初始化 gRPC 服务器失败", "error", err)
			os.Exit(1)
		}

		// 注册 gRPC 服务
		authGRPCService := grpcservice.NewAuthService(nil, grpcLogger)
		grpcServer.RegisterServices(
			authGRPCService,
			nil, // adminService
			nil, // roleService
			nil, // permissionService
			nil, // auditLogService
		)

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
