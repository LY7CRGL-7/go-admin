package service

import (
	"context"

	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// RoleService 角色管理服务
type RoleService struct {
	v1.UnimplementedRoleServiceServer
	uc  *biz.RoleUsecase
	log *log.Helper
}

// NewRoleService 创建角色服务
func NewRoleService(uc *biz.RoleUsecase, logger log.Logger) *RoleService {
	return &RoleService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/role")),
	}
}

func (s *RoleService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.RoleInfo, error) {
	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      1,
	}
	if claims := GetClaimsFromContext(ctx); claims != nil {
		role.TenantID = claims.TenantID
	}

	if err := s.uc.Create(ctx, role); err != nil {
		return nil, err
	}
	return roleToProto(role), nil
}

func (s *RoleService) GetRole(ctx context.Context, req *v1.GetRoleRequest) (*v1.RoleInfo, error) {
	role, err := s.uc.Get(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return roleToProto(role), nil
}

func (s *RoleService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.RoleInfo, error) {
	role, err := s.uc.Get(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != 0 {
		role.Status = int8(req.Status)
	}
	role.Sort = int(req.Sort)

	if err := s.uc.Update(ctx, role); err != nil {
		return nil, err
	}
	return roleToProto(role), nil
}

func (s *RoleService) DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*v1.DeleteReply, error) {
	if err := s.uc.Delete(ctx, uint(req.Id)); err != nil {
		return nil, err
	}
	return &v1.DeleteReply{Success: true}, nil
}

func (s *RoleService) ListRoles(ctx context.Context, req *v1.ListRolesRequest) (*v1.ListRolesReply, error) {
	var tenantID uint
	if claims := GetClaimsFromContext(ctx); claims != nil {
		tenantID = claims.TenantID
	}

	roles, total, err := s.uc.List(ctx, tenantID, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var items []*v1.RoleInfo
	for _, r := range roles {
		items = append(items, roleToProto(r))
	}
	return &v1.ListRolesReply{Items: items, Total: total}, nil
}

func (s *RoleService) AssignPermissions(ctx context.Context, req *v1.AssignPermissionsRequest) (*v1.AssignPermissionsReply, error) {
	var permIDs []uint
	for _, id := range req.PermissionIds {
		permIDs = append(permIDs, uint(id))
	}
	if err := s.uc.AssignPermissions(ctx, uint(req.RoleId), permIDs); err != nil {
		return nil, err
	}
	return &v1.AssignPermissionsReply{Success: true}, nil
}
