package data

import (
	"context"

	"admin/internal/biz"
	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type permissionRepo struct {
	data *Data
	log  *log.Helper
}

// NewPermissionRepo 创建权限仓储
func NewPermissionRepo(data *Data, logger log.Logger) biz.PermissionRepo {
	return &permissionRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/permission")),
	}
}

func (r *permissionRepo) List(ctx context.Context) ([]*model.Permission, error) {
	var perms []*model.Permission
	err := r.data.DB.WithContext(ctx).Order("sort ASC, id ASC").Find(&perms).Error
	return perms, err
}

func (r *permissionRepo) GetTree(ctx context.Context) ([]*model.Permission, error) {
	var perms []*model.Permission
	err := r.data.DB.WithContext(ctx).Order("sort ASC, id ASC").Find(&perms).Error
	if err != nil {
		return nil, err
	}
	return buildPermissionTree(perms, 0), nil
}

func (r *permissionRepo) GetByRoleIDs(ctx context.Context, roleIDs []uint) ([]*model.Permission, error) {
	var perms []*model.Permission
	err := r.data.DB.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Distinct().
		Find(&perms).Error
	return perms, err
}

// buildPermissionTree 递归构建权限树
func buildPermissionTree(perms []*model.Permission, parentID uint) []*model.Permission {
	var tree []*model.Permission
	for _, p := range perms {
		if p.ParentID == parentID {
			tree = append(tree, p)
		}
	}
	return tree
}
