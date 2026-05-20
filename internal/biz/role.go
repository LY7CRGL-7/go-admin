package biz

import (
	"context"

	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// RoleUsecase 角色业务逻辑
type RoleUsecase struct {
	repo     RoleRepo
	permRepo PermissionRepo
	log      *log.Helper
}

// NewRoleUsecase 创建角色用例
func NewRoleUsecase(repo RoleRepo, permRepo PermissionRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{
		repo:     repo,
		permRepo: permRepo,
		log:      log.NewHelper(log.With(logger, "module", "biz/role")),
	}
}

func (uc *RoleUsecase) Create(ctx context.Context, role *model.Role) error {
	return uc.repo.Create(ctx, role)
}

func (uc *RoleUsecase) Get(ctx context.Context, id uint) (*model.Role, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *RoleUsecase) Update(ctx context.Context, role *model.Role) error {
	return uc.repo.Update(ctx, role)
}

func (uc *RoleUsecase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *RoleUsecase) List(ctx context.Context, tenantID uint, page, pageSize int) ([]*model.Role, int64, error) {
	return uc.repo.List(ctx, tenantID, page, pageSize)
}

func (uc *RoleUsecase) AssignPermissions(ctx context.Context, roleID uint, permIDs []uint) error {
	return uc.repo.AssignPermissions(ctx, roleID, permIDs)
}

// ListPermissions 列出所有权限
func (uc *RoleUsecase) ListPermissions(ctx context.Context) ([]*model.Permission, error) {
	return uc.permRepo.List(ctx)
}

// GetPermissionTree 获取权限树
func (uc *RoleUsecase) GetPermissionTree(ctx context.Context) ([]*model.Permission, error) {
	return uc.permRepo.GetTree(ctx)
}
