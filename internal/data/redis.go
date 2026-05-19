package data

import (
	"admin/internal/conf"
	"admin/internal/pkg/logger"
	"context"

	"github.com/redis/go-redis/v9"
)

// NewRedis 创建 Redis 连接
func NewRedis(cfg *conf.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})
	
	// 测试连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("Redis 连接失败", "error", err)
		return nil, err
	}
	
	logger.Info("Redis 连接成功")
	return rdb, nil
}
