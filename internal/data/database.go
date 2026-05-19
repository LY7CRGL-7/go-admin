package data

import (
	"admin/internal/conf"
	"admin/internal/data/model"
	"admin/internal/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "gorm.io/gorm/logger"
)

// NewDatabase 创建数据库连接
func NewDatabase(cfg *conf.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 自动迁移数据表
	if err := autoMigrate(db); err != nil {
		logger.Error("数据库迁移失败", "error", err)
		return nil, err
	}

	logger.Info("数据库连接成功")
	return db, nil
}

// autoMigrate 自动迁移数据表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Admin{},
		&model.Role{},
		&model.Permission{},
		&model.AdminRole{},
		&model.RolePermission{},
		&model.AuditLog{},
		&model.LoginAttempt{},
	)
}
