package service

import (
	"context"

	v1 "admin/api/admin/v1"
	"admin/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// AuditService 审计日志服务
type AuditService struct {
	v1.UnimplementedAuditLogServiceServer
	repo biz.AuditRepo
	log  *log.Helper
}

// NewAuditService 创建审计日志服务
func NewAuditService(repo biz.AuditRepo, logger log.Logger) *AuditService {
	return &AuditService{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "service/audit")),
	}
}

func (s *AuditService) ListAuditLogs(ctx context.Context, req *v1.ListAuditLogsRequest) (*v1.ListAuditLogsReply, error) {
	var tenantID uint
	if claims := GetClaimsFromContext(ctx); claims != nil {
		tenantID = claims.TenantID
	}

	logs, total, err := s.repo.List(ctx, tenantID, int(req.Page), int(req.PageSize), uint(req.UserId), req.Action)
	if err != nil {
		return nil, err
	}

	var items []*v1.AuditLogInfo
	for _, l := range logs {
		items = append(items, auditToProto(l))
	}
	return &v1.ListAuditLogsReply{Items: items, Total: total}, nil
}
