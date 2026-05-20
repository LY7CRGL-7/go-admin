package data

import (
	"context"

	"admin/internal/biz"
	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type roleRepo struct {
	data *Data
	log  *log.Helper
}

// NewRoleRepo 创建角色仓储
func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/role")),
	}
}

func (r *roleRepo) Create(ctx context.Context, role *model.Role) error {
	return r.data.DB.WithContext(ctx).Create(role).Error
}

func (r *roleRepo) GetByID(ctx context.Context, id uint) (*model.Role, error) {
	var role model.Role
	err := r.data.DB.WithContext(ctx).Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) Update(ctx context.Context, role *model.Role) error {
	return r.data.DB.WithContext(ctx).Save(role).Error
}

func (r *roleRepo) Delete(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除角色-权限关联
		if err := tx.Where("role_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}
		// 删除用户-角色关联
		if err := tx.Where("role_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		// 删除角色
		return tx.Delete(&model.Role{}, id).Error
	})
}

func (r *roleRepo) List(ctx context.Context, tenantID uint, page, pageSize int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	db := r.data.DB.WithContext(ctx).Model(&model.Role{})
	if tenantID > 0 {
		db = db.Where("tenant_id = ?", tenantID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Permissions").Offset(offset).Limit(pageSize).Order("sort ASC, id ASC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *roleRepo) AssignPermissions(ctx context.Context, roleID uint, permIDs []uint) error {
	return r.data.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}
		for _, permID := range permIDs {
			if err := tx.Create(&model.RolePermission{RoleID: roleID, PermissionID: permID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
