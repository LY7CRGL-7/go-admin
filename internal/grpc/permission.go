package grpc

import (
	"context"

	"admin/internal/service"
	pb "admin/proto/admin/v1"

	"go.uber.org/zap"
)

// PermissionService gRPC 权限服务
type PermissionService struct {
	pb.UnimplementedPermissionServiceServer
	permissionService *service.PermissionService
	logger            *zap.Logger
}

// NewPermissionService 创建权限服务
func NewPermissionService(permissionService *service.PermissionService, logger *zap.Logger) *PermissionService {
	return &PermissionService{
		permissionService: permissionService,
		logger:            logger,
	}
}

// ListPermissions 列出权限
func (s *PermissionService) ListPermissions(ctx context.Context, req *pb.ListPermissionsRequest) (*pb.PermissionListResponse, error) {
	s.logger.Info("gRPC ListPermissions")

	permissions, err := s.permissionService.ListPermissions(ctx)
	if err != nil {
		s.logger.Error("ListPermissions failed", zap.Error(err))
		return nil, err
	}

	permissionList := make([]*pb.PermissionInfo, len(permissions))
	for i, perm := range permissions {
		permissionList[i] = &pb.PermissionInfo{
			Id:       int64(perm.ID),
			Name:     perm.Name,
			Code:     perm.Code,
			Resource: perm.Path,
			Action:   perm.Method,
			ParentId: int64(perm.ParentID),
		}
	}

	return &pb.PermissionListResponse{
		Permissions: permissionList,
	}, nil
}
