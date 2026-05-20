package service

import (
	"context"

	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// TenantService 租户管理服务
type TenantService struct {
	v1.UnimplementedTenantServiceServer
	uc  *biz.TenantUsecase
	log *log.Helper
}

// NewTenantService 创建租户服务
func NewTenantService(uc *biz.TenantUsecase, logger log.Logger) *TenantService {
	return &TenantService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/tenant")),
	}
}

func (s *TenantService) CreateTenant(ctx context.Context, req *v1.CreateTenantRequest) (*v1.TenantInfo, error) {
	tenant := &model.Tenant{
		Name:         req.Name,
		Code:         req.Code,
		Domain:       req.Domain,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		MaxUsers:     int(req.MaxUsers),
		Status:       1,
	}
	if err := s.uc.Create(ctx, tenant); err != nil {
		return nil, err
	}
	return tenantToProto(tenant), nil
}

func (s *TenantService) GetTenant(ctx context.Context, req *v1.GetTenantRequest) (*v1.TenantInfo, error) {
	tenant, err := s.uc.Get(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return tenantToProto(tenant), nil
}

func (s *TenantService) UpdateTenant(ctx context.Context, req *v1.UpdateTenantRequest) (*v1.TenantInfo, error) {
	tenant, err := s.uc.Get(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Domain != "" {
		tenant.Domain = req.Domain
	}
	if req.ContactName != "" {
		tenant.ContactName = req.ContactName
	}
	if req.ContactPhone != "" {
		tenant.ContactPhone = req.ContactPhone
	}
	if req.Status != 0 {
		tenant.Status = int8(req.Status)
	}
	if req.MaxUsers != 0 {
		tenant.MaxUsers = int(req.MaxUsers)
	}

	if err := s.uc.Update(ctx, tenant); err != nil {
		return nil, err
	}
	return tenantToProto(tenant), nil
}

func (s *TenantService) DeleteTenant(ctx context.Context, req *v1.DeleteTenantRequest) (*v1.DeleteReply, error) {
	if err := s.uc.Delete(ctx, uint(req.Id)); err != nil {
		return nil, err
	}
	return &v1.DeleteReply{Success: true}, nil
}

func (s *TenantService) ListTenants(ctx context.Context, req *v1.ListTenantsRequest) (*v1.ListTenantsReply, error) {
	tenants, total, err := s.uc.List(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var items []*v1.TenantInfo
	for _, t := range tenants {
		items = append(items, tenantToProto(t))
	}
	return &v1.ListTenantsReply{Items: items, Total: total}, nil
}
