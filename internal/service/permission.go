package service

import (
	"context"

	v1 "admin/api/admin/v1"
	"admin/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// PermissionService 权限管理服务
type PermissionService struct {
	v1.UnimplementedPermissionServiceServer
	uc  *biz.RoleUsecase
	log *log.Helper
}

// NewPermissionService 创建权限服务
func NewPermissionService(uc *biz.RoleUsecase, logger log.Logger) *PermissionService {
	return &PermissionService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/permission")),
	}
}

func (s *PermissionService) ListPermissions(ctx context.Context, _ *v1.ListPermissionsRequest) (*v1.ListPermissionsReply, error) {
	perms, err := s.uc.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}

	var items []*v1.PermissionInfo
	for _, p := range perms {
		items = append(items, permToProto(p))
	}
	return &v1.ListPermissionsReply{Items: items}, nil
}

func (s *PermissionService) GetPermissionTree(ctx context.Context, _ *v1.GetPermissionTreeRequest) (*v1.GetPermissionTreeReply, error) {
	perms, err := s.uc.GetPermissionTree(ctx)
	if err != nil {
		return nil, err
	}

	var items []*v1.PermissionInfo
	for _, p := range perms {
		items = append(items, permToProto(p))
	}
	return &v1.GetPermissionTreeReply{Items: items}, nil
}
