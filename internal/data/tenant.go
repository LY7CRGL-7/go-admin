package data

import (
	"context"

	"admin/internal/biz"
	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type tenantRepo struct {
	data *Data
	log  *log.Helper
}

// NewTenantRepo 创建租户仓储
func NewTenantRepo(data *Data, logger log.Logger) biz.TenantRepo {
	return &tenantRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/tenant")),
	}
}

func (r *tenantRepo) Create(ctx context.Context, tenant *model.Tenant) error {
	return r.data.DB.WithContext(ctx).Create(tenant).Error
}

func (r *tenantRepo) GetByID(ctx context.Context, id uint) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.data.DB.WithContext(ctx).First(&tenant, id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepo) GetByCode(ctx context.Context, code string) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.data.DB.WithContext(ctx).Where("code = ?", code).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepo) Update(ctx context.Context, tenant *model.Tenant) error {
	return r.data.DB.WithContext(ctx).Save(tenant).Error
}

func (r *tenantRepo) Delete(ctx context.Context, id uint) error {
	return r.data.DB.WithContext(ctx).Delete(&model.Tenant{}, id).Error
}

func (r *tenantRepo) List(ctx context.Context, page, pageSize int) ([]*model.Tenant, int64, error) {
	var tenants []*model.Tenant
	var total int64

	db := r.data.DB.WithContext(ctx).Model(&model.Tenant{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tenants).Error; err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}
