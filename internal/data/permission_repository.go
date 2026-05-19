package data

import (
	"admin/internal/data/model"
	"context"

	"gorm.io/gorm"
)

// PermissionRepository 权限数据访问接口
type PermissionRepository interface {
	Create(ctx context.Context, permission *model.Permission) error
	GetByID(ctx context.Context, id uint) (*model.Permission, error)
	List(ctx context.Context) ([]*model.Permission, error)
	ListByRoleID(ctx context.Context, roleID uint) ([]*model.Permission, error)
}

type permissionRepo struct {
	db *gorm.DB
}

// NewPermissionRepo 创建权限仓库
func NewPermissionRepo(db *gorm.DB) PermissionRepository {
	return &permissionRepo{db: db}
}

func (r *permissionRepo) Create(ctx context.Context, permission *model.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *permissionRepo) GetByID(ctx context.Context, id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepo) List(ctx context.Context) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.WithContext(ctx).Order("sort ASC").Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepo) ListByRoleID(ctx context.Context, roleID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions rp ON permissions.id = rp.permission_id").
		Where("rp.role_id = ?", roleID).
		Order("permissions.sort ASC").
		Find(&permissions).Error
	return permissions, err
}
