package data

import (
	"admin/internal/data/model"
	"context"

	"gorm.io/gorm"
)

// AuditLogRepository 审计日志数据访问接口
type AuditLogRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	GetByID(ctx context.Context, id uint) (*model.AuditLog, error)
	List(ctx context.Context, page, pageSize int, adminID uint, action string) ([]*model.AuditLog, int64, error)
}

type auditLogRepo struct {
	db *gorm.DB
}

// NewAuditLogRepo 创建审计日志仓库
func NewAuditLogRepo(db *gorm.DB) AuditLogRepository {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditLogRepo) GetByID(ctx context.Context, id uint) (*model.AuditLog, error) {
	var log model.AuditLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *auditLogRepo) List(ctx context.Context, page, pageSize int, adminID uint, action string) ([]*model.AuditLog, int64, error) {
	var logs []*model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	if adminID > 0 {
		query = query.Where("admin_id = ?", adminID)
	}
	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}
