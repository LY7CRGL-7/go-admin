package grpc

import (
	"context"

	"admin/internal/data/model"
	"admin/internal/service"
	pb "admin/proto/admin/v1"

	"go.uber.org/zap"
)

// RoleService gRPC 角色服务
type RoleService struct {
	pb.UnimplementedRoleServiceServer
	roleService *service.RoleService
	logger      *zap.Logger
}

// NewRoleService 创建角色服务
func NewRoleService(roleService *service.RoleService, logger *zap.Logger) *RoleService {
	return &RoleService{
		roleService: roleService,
		logger:      logger,
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.RoleInfo, error) {
	s.logger.Info("gRPC CreateRole", zap.String("code", req.Code))

	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      1,
	}

	if err := s.roleService.CreateRole(ctx, role); err != nil {
		s.logger.Error("CreateRole failed", zap.Error(err))
		return nil, err
	}

	return &pb.RoleInfo{
		Id:          int64(role.ID),
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
	}, nil
}

// GetRole 获取角色
func (s *RoleService) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.RoleInfo, error) {
	s.logger.Info("gRPC GetRole", zap.Int64("id", req.Id))

	role, err := s.roleService.GetRole(ctx, uint(req.Id))
	if err != nil {
		s.logger.Error("GetRole failed", zap.Error(err))
		return nil, err
	}

	return &pb.RoleInfo{
		Id:          int64(role.ID),
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
	}, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.RoleInfo, error) {
	s.logger.Info("gRPC UpdateRole", zap.Int64("id", req.Id))

	role := &model.Role{
		ID:          uint(req.Id),
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.roleService.UpdateRole(ctx, role); err != nil {
		s.logger.Error("UpdateRole failed", zap.Error(err))
		return nil, err
	}

	return &pb.RoleInfo{
		Id:          int64(role.ID),
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
	}, nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteResponse, error) {
	s.logger.Info("gRPC DeleteRole", zap.Int64("id", req.Id))

	if err := s.roleService.DeleteRole(ctx, uint(req.Id)); err != nil {
		s.logger.Error("DeleteRole failed", zap.Error(err))
		return &pb.DeleteResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteResponse{
		Success: true,
		Message: "删除成功",
	}, nil
}

// ListRoles 列出角色
func (s *RoleService) ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.RoleListResponse, error) {
	s.logger.Info("gRPC ListRoles")

	roles, total, err := s.roleService.ListRoles(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		s.logger.Error("ListRoles failed", zap.Error(err))
		return nil, err
	}

	roleList := make([]*pb.RoleInfo, len(roles))
	for i, role := range roles {
		roleList[i] = &pb.RoleInfo{
			Id:          int64(role.ID),
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
		}
	}

	return &pb.RoleListResponse{
		Roles: roleList,
		Total: total,
	}, nil
}

// AssignPermissions 分配权限
func (s *RoleService) AssignPermissions(ctx context.Context, req *pb.AssignPermissionsRequest) (*pb.AssignPermissionsResponse, error) {
	s.logger.Info("gRPC AssignPermissions", zap.Int64("role_id", req.RoleId))

	permissionIDs := make([]uint, len(req.PermissionIds))
	for i, id := range req.PermissionIds {
		permissionIDs[i] = uint(id)
	}

	if err := s.roleService.AssignPermissions(ctx, uint(req.RoleId), permissionIDs); err != nil {
		s.logger.Error("AssignPermissions failed", zap.Error(err))
		return &pb.AssignPermissionsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.AssignPermissionsResponse{
		Success: true,
		Message: "权限分配成功",
	}, nil
}
