package biz

import (
	"context"

	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// TenantUsecase 租户业务逻辑
type TenantUsecase struct {
	repo TenantRepo
	log  *log.Helper
}

// NewTenantUsecase 创建租户用例
func NewTenantUsecase(repo TenantRepo, logger log.Logger) *TenantUsecase {
	return &TenantUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "biz/tenant")),
	}
}

func (uc *TenantUsecase) Create(ctx context.Context, tenant *model.Tenant) error {
	return uc.repo.Create(ctx, tenant)
}

func (uc *TenantUsecase) Get(ctx context.Context, id uint) (*model.Tenant, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *TenantUsecase) GetByCode(ctx context.Context, code string) (*model.Tenant, error) {
	return uc.repo.GetByCode(ctx, code)
}

func (uc *TenantUsecase) Update(ctx context.Context, tenant *model.Tenant) error {
	return uc.repo.Update(ctx, tenant)
}

func (uc *TenantUsecase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *TenantUsecase) List(ctx context.Context, page, pageSize int) ([]*model.Tenant, int64, error) {
	return uc.repo.List(ctx, page, pageSize)
}
