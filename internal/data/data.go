package data

import (
	"context"
	"fmt"
	"time"

	"admin/internal/conf"
	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet data 层依赖注入集合
var ProviderSet = wire.NewSet(
	NewData,
	NewUserRepo,
	NewRoleRepo,
	NewPermissionRepo,
	NewTenantRepo,
	NewAuditRepo,
)

// Data 数据层封装，持有 DB 和 Redis 连接
type Data struct {
	DB    *gorm.DB
	Redis *redis.Client
	log   *log.Helper
}

// NewData 创建数据层实例
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "data"))

	// 初始化数据库
	db, err := newDB(c.Database)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移
	if err := model.AutoMigrate(db); err != nil {
		return nil, nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	// 初始化 Redis
	rdb := newRedis(c.Redis)

	cleanup := func() {
		helper.Info("closing data resources")
		if err := rdb.Close(); err != nil {
			helper.Errorf("failed to close redis: %v", err)
		}
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			if err := sqlDB.Close(); err != nil {
				helper.Errorf("failed to close db: %v", err)
			}
		}
	}

	return &Data{DB: db, Redis: rdb, log: helper}, cleanup, nil
}

// newDB 根据 driver 创建数据库连接（支持 mysql / postgres）
func newDB(c conf.Database) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch c.Driver {
	case "mysql":
		dialector = mysql.Open(c.Source)
	case "postgres":
		dialector = postgres.Open(c.Source)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s (use mysql or postgres)", c.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if c.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(c.MaxIdle)
	}
	if c.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(c.MaxOpen)
	}
	if c.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime))
	} else {
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db, nil
}

// newRedis 创建 Redis 连接
func newRedis(c conf.Redis) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		DialTimeout:  time.Duration(c.DialTimeout),
		ReadTimeout:  time.Duration(c.ReadTimeout),
		WriteTimeout: time.Duration(c.WriteTimeout),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		// Redis 连接失败不阻止启动，仅记录警告
		fmt.Printf("⚠️  Redis connection failed: %v (service will start without cache)\n", err)
	}

	return rdb
}
