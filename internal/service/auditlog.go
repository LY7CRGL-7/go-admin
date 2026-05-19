package service

import (
	"admin/internal/data"
	"admin/internal/data/model"
	"context"
)

// AuditLogService 审计日志服务
type AuditLogService struct {
	auditLogRepo data.AuditLogRepository
}

// NewAuditLogService 创建审计日志服务
func NewAuditLogService(auditLogRepo data.AuditLogRepository) *AuditLogService {
	return &AuditLogService{
		auditLogRepo: auditLogRepo,
	}
}

// CreateAuditLog 创建审计日志
func (s *AuditLogService) CreateAuditLog(ctx context.Context, log *model.AuditLog) error {
	return s.auditLogRepo.Create(ctx, log)
}

// GetAuditLog 获取审计日志
func (s *AuditLogService) GetAuditLog(ctx context.Context, id uint) (*model.AuditLog, error) {
	return s.auditLogRepo.GetByID(ctx, id)
}

// ListAuditLogs 列出审计日志
func (s *AuditLogService) ListAuditLogs(ctx context.Context, page, pageSize int, adminID uint, action string) ([]*model.AuditLog, int64, error) {
	return s.auditLogRepo.List(ctx, page, pageSize, adminID, action)
}
