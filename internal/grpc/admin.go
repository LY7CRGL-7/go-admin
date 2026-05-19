package grpc

import (
	"context"

	"admin/internal/data/model"
	"admin/internal/service"
	pb "admin/proto/admin/v1"

	"go.uber.org/zap"
)

// AdminService gRPC 管理员服务
type AdminService struct {
	pb.UnimplementedAdminServiceServer
	adminService *service.AdminService
	logger       *zap.Logger
}

// NewAdminService 创建管理员服务
func NewAdminService(adminService *service.AdminService, logger *zap.Logger) *AdminService {
	return &AdminService{
		adminService: adminService,
		logger:       logger,
	}
}

// CreateAdmin 创建管理员
func (s *AdminService) CreateAdmin(ctx context.Context, req *pb.CreateAdminRequest) (*pb.AdminInfo, error) {
	s.logger.Info("gRPC CreateAdmin", zap.String("username", req.Username))

	admin := &model.Admin{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   int8(req.Status),
	}

	roleIDs := make([]uint, len(req.RoleIds))
	for i, id := range req.RoleIds {
		roleIDs[i] = uint(id)
	}

	if err := s.adminService.CreateAdmin(ctx, admin, roleIDs); err != nil {
		s.logger.Error("CreateAdmin failed", zap.Error(err))
		return nil, err
	}

	return &pb.AdminInfo{
		Id:       int64(admin.ID),
		Username: admin.Username,
		Nickname: admin.Nickname,
		Email:    admin.Email,
		Phone:    admin.Phone,
		Status:   int32(admin.Status),
	}, nil
}

// GetAdmin 获取管理员
func (s *AdminService) GetAdmin(ctx context.Context, req *pb.GetAdminRequest) (*pb.AdminInfo, error) {
	s.logger.Info("gRPC GetAdmin", zap.Int64("id", req.Id))

	admin, err := s.adminService.GetAdmin(ctx, uint(req.Id))
	if err != nil {
		s.logger.Error("GetAdmin failed", zap.Error(err))
		return nil, err
	}

	return &pb.AdminInfo{
		Id:       int64(admin.ID),
		Username: admin.Username,
		Nickname: admin.Nickname,
		Email:    admin.Email,
		Phone:    admin.Phone,
		Status:   int32(admin.Status),
	}, nil
}

// UpdateAdmin 更新管理员
func (s *AdminService) UpdateAdmin(ctx context.Context, req *pb.UpdateAdminRequest) (*pb.AdminInfo, error) {
	s.logger.Info("gRPC UpdateAdmin", zap.Int64("id", req.Id))

	admin := &model.Admin{
		ID:       uint(req.Id),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   int8(req.Status),
	}

	if err := s.adminService.UpdateAdmin(ctx, admin); err != nil {
		s.logger.Error("UpdateAdmin failed", zap.Error(err))
		return nil, err
	}

	return &pb.AdminInfo{
		Id:       int64(admin.ID),
		Username: admin.Username,
		Nickname: admin.Nickname,
		Email:    admin.Email,
		Phone:    admin.Phone,
		Status:   int32(admin.Status),
	}, nil
}

// DeleteAdmin 删除管理员
func (s *AdminService) DeleteAdmin(ctx context.Context, req *pb.DeleteAdminRequest) (*pb.Response, error) {
	s.logger.Info("gRPC DeleteAdmin", zap.Int64("id", req.Id))

	if err := s.adminService.DeleteAdmin(ctx, uint(req.Id)); err != nil {
		s.logger.Error("DeleteAdmin failed", zap.Error(err))
		return &pb.Response{
			Code: 500,
			Msg:  err.Error(),
		}, nil
	}

	return &pb.Response{
		Code: 0,
		Msg:  "删除成功",
	}, nil
}

// ListAdmins 列出管理员
func (s *AdminService) ListAdmins(ctx context.Context, req *pb.ListAdminsRequest) (*pb.AdminListResponse, error) {
	s.logger.Info("gRPC ListAdmins")

	admins, total, err := s.adminService.ListAdmins(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		s.logger.Error("ListAdmins failed", zap.Error(err))
		return nil, err
	}

	adminList := make([]*pb.AdminInfo, len(admins))
	for i, admin := range admins {
		adminList[i] = &pb.AdminInfo{
			Id:       int64(admin.ID),
			Username: admin.Username,
			Nickname: admin.Nickname,
			Email:    admin.Email,
			Phone:    admin.Phone,
			Status:   int32(admin.Status),
		}
	}

	return &pb.AdminListResponse{
		Admins: adminList,
		Total:  total,
	}, nil
}
