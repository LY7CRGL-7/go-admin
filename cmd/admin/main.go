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

// go build -ldflags "-X main.Version=x.y.z"
var (
	Version = "1.0.0"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "cmd/admin/conf/config.yaml", "配置文件路径")
	flag.Parse()

	fmt.Printf("Admin Management System v%s\n", Version)

	// 加载配置
	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log.FilePath, cfg.Log.Level, cfg.Log.MaxSize, cfg.Log.MaxBackups, cfg.Log.MaxAge); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	logger.Info("启动管理端服务", "version", Version)

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

	// 启动 HTTP 服务器
	logger.Info("HTTP 服务器启动中", "port", cfg.Server.Port)
	if err := server.StartServer(cfg, db, rdb); err != nil {
		logger.Error("启动 HTTP 服务器失败", "error", err)
		os.Exit(1)
	}
}

// loadConfig 加载配置文件
func loadConfig(path string) (*conf.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg conf.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}
