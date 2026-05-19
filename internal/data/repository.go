package data

import (
	"admin/internal/data/model"
	"context"
	"time"

	"gorm.io/gorm"
)

// AdminRepository 管理员数据访问接口
type AdminRepository interface {
	Create(ctx context.Context, admin *model.Admin) error
	GetByID(ctx context.Context, id uint) (*model.Admin, error)
	GetByUsername(ctx context.Context, username string) (*model.Admin, error)
	Update(ctx context.Context, admin *model.Admin) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, pageSize int) ([]*model.Admin, int64, error)
	UpdateLastLogin(ctx context.Context, id uint, ip string) error
	AssignRoles(ctx context.Context, adminID uint, roleIDs []uint) error
}

type adminRepo struct {
	db *gorm.DB
}

// NewAdminRepo 创建管理员仓库
func NewAdminRepo(db *gorm.DB) AdminRepository {
	return &adminRepo{db: db}
}

func (r *adminRepo) Create(ctx context.Context, admin *model.Admin) error {
	return r.db.WithContext(ctx).Create(admin).Error
}

func (r *adminRepo) GetByID(ctx context.Context, id uint) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.WithContext(ctx).Preload("Roles").First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepo) GetByUsername(ctx context.Context, username string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.WithContext(ctx).Preload("Roles.Permissions").Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepo) Update(ctx context.Context, admin *model.Admin) error {
	return r.db.WithContext(ctx).Save(admin).Error
}

func (r *adminRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Admin{}, id).Error
}

func (r *adminRepo) List(ctx context.Context, page, pageSize int) ([]*model.Admin, int64, error) {
	var admins []*model.Admin
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Admin{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Roles").Offset(offset).Limit(pageSize).Find(&admins).Error
	return admins, total, err
}

func (r *adminRepo) UpdateLastLogin(ctx context.Context, id uint, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.Admin{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
	}).Error
}

func (r *adminRepo) AssignRoles(ctx context.Context, adminID uint, roleIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧的角色关联
		if err := tx.Where("admin_id = ?", adminID).Delete(&model.AdminRole{}).Error; err != nil {
			return err
		}
		// 创建新的角色关联
		for _, roleID := range roleIDs {
			if err := tx.Create(&model.AdminRole{
				AdminID: adminID,
				RoleID:  roleID,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
