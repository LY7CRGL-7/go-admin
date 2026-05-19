package grpc

import (
	"context"

	"admin/internal/service"
	pb "admin/proto/admin/v1"

	"go.uber.org/zap"
)

// AuditLogService gRPC 审计日志服务
type AuditLogService struct {
	pb.UnimplementedAuditLogServiceServer
	auditLogService *service.AuditLogService
	logger          *zap.Logger
}

// NewAuditLogService 创建审计日志服务
func NewAuditLogService(auditLogService *service.AuditLogService, logger *zap.Logger) *AuditLogService {
	return &AuditLogService{
		auditLogService: auditLogService,
		logger:          logger,
	}
}

// ListAuditLogs 列出审计日志
func (s *AuditLogService) ListAuditLogs(ctx context.Context, req *pb.AuditLogListRequest) (*pb.AuditLogListResponse, error) {
	s.logger.Info("gRPC ListAuditLogs")

	logs, total, err := s.auditLogService.ListAuditLogs(ctx, int(req.Page), int(req.PageSize), uint(req.AdminId), req.Action)
	if err != nil {
		s.logger.Error("ListAuditLogs failed", zap.Error(err))
		return nil, err
	}

	logList := make([]*pb.AuditLogInfo, len(logs))
	for i, log := range logs {
		logList[i] = &pb.AuditLogInfo{
			Id:            int64(log.ID),
			AdminId:       int64(log.AdminID),
			AdminUsername: log.AdminName,
			Action:        log.Action,
			Resource:      log.Resource,
			RequestMethod: log.Method,
			RequestPath:   log.Path,
			RequestBody:   log.Request,
			ResponseBody:  log.Response,
			Ip:            log.IP,
			UserAgent:     log.UserAgent,
			Duration:      log.Duration,
			StatusCode:    int32(log.Status),
			CreatedAt:     log.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &pb.AuditLogListResponse{
		Logs:  logList,
		Total: total,
	}, nil
}

// GetAuditLog 获取审计日志
func (s *AuditLogService) GetAuditLog(ctx context.Context, req *pb.GetAuditLogRequest) (*pb.AuditLogInfo, error) {
	s.logger.Info("gRPC GetAuditLog", zap.Int64("id", req.Id))

	log, err := s.auditLogService.GetAuditLog(ctx, uint(req.Id))
	if err != nil {
		s.logger.Error("GetAuditLog failed", zap.Error(err))
		return nil, err
	}

	return &pb.AuditLogInfo{
		Id:            int64(log.ID),
		AdminId:       int64(log.AdminID),
		AdminUsername: log.AdminName,
		Action:        log.Action,
		Resource:      log.Resource,
		RequestMethod: log.Method,
		RequestPath:   log.Path,
		RequestBody:   log.Request,
		ResponseBody:  log.Response,
		Ip:            log.IP,
		UserAgent:     log.UserAgent,
		Duration:      log.Duration,
		StatusCode:    int32(log.Status),
		CreatedAt:     log.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
