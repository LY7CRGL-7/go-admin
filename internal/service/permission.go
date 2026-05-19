package service

import (
	"admin/internal/data"
	"admin/internal/data/model"
	"context"
)

// PermissionService 权限服务
type PermissionService struct {
	permissionRepo data.PermissionRepository
}

// NewPermissionService 创建权限服务
func NewPermissionService(permissionRepo data.PermissionRepository) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
	}
}

// ListPermissions 列出权限
func (s *PermissionService) ListPermissions(ctx context.Context) ([]*model.Permission, error) {
	return s.permissionRepo.List(ctx)
}

// GetPermission 获取权限
func (s *PermissionService) GetPermission(ctx context.Context, id uint) (*model.Permission, error) {
	return s.permissionRepo.GetByID(ctx, id)
}
