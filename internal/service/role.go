package service

import (
	"admin/internal/data"
	"admin/internal/data/model"
	"context"
	"errors"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo       data.RoleRepository
	permissionRepo data.PermissionRepository
}

// NewRoleService 创建角色服务
func NewRoleService(roleRepo data.RoleRepository, permissionRepo data.PermissionRepository) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, role *model.Role) error {
	// 检查角色代码是否已存在
	existing, _ := s.roleRepo.GetByCode(ctx, role.Code)
	if existing != nil {
		return errors.New("角色代码已存在")
	}

	return s.roleRepo.Create(ctx, role)
}

// GetRole 获取角色
func (s *RoleService) GetRole(ctx context.Context, id uint) (*model.Role, error) {
	return s.roleRepo.GetByID(ctx, id)
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, role *model.Role) error {
	_, err := s.roleRepo.GetByID(ctx, role.ID)
	if err != nil {
		return errors.New("角色不存在")
	}

	return s.roleRepo.Update(ctx, role)
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, id uint) error {
	return s.roleRepo.Delete(ctx, id)
}

// ListRoles 列出角色
func (s *RoleService) ListRoles(ctx context.Context, page, pageSize int) ([]*model.Role, int64, error) {
	return s.roleRepo.List(ctx, page, pageSize)
}

// AssignPermissions 分配权限
func (s *RoleService) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	// 验证角色是否存在
	_, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 验证权限是否存在
	for _, permissionID := range permissionIDs {
		_, err := s.permissionRepo.GetByID(ctx, permissionID)
		if err != nil {
			return errors.New("权限不存在")
		}
	}

	return s.roleRepo.AssignPermissions(ctx, roleID, permissionIDs)
}
