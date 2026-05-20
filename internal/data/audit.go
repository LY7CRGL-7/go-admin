package data

import (
	"context"

	"admin/internal/biz"
	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type auditRepo struct {
	data *Data
	log  *log.Helper
}

// NewAuditRepo 创建审计日志仓储
func NewAuditRepo(data *Data, logger log.Logger) biz.AuditRepo {
	return &auditRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/audit")),
	}
}

func (r *auditRepo) Create(ctx context.Context, auditLog *model.AuditLog) error {
	return r.data.DB.WithContext(ctx).Create(auditLog).Error
}

func (r *auditRepo) List(ctx context.Context, tenantID uint, page, pageSize int, userID uint, action string) ([]*model.AuditLog, int64, error) {
	var logs []*model.AuditLog
	var total int64

	db := r.data.DB.WithContext(ctx).Model(&model.AuditLog{})

	if tenantID > 0 {
		db = db.Where("tenant_id = ?", tenantID)
	}
	if userID > 0 {
		db = db.Where("user_id = ?", userID)
	}
	if action != "" {
		db = db.Where("action = ?", action)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
