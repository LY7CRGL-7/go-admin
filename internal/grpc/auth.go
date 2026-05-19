package grpc

import (
	"context"

	"admin/internal/pkg/logger"
	"admin/internal/service"
	pb "admin/proto/admin/v1"

	"go.uber.org/zap"
)

// AuthService gRPC 认证服务
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
	logger      *zap.Logger
}

// NewAuthService 创建认证服务
func NewAuthService(authService *service.AuthService, logger *zap.Logger) *AuthService {
	return &AuthService{
		authService: authService,
		logger:      logger,
	}
}

// Login 登录
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.logger.Info("gRPC Login", zap.String("username", req.Username))

	token, admin, err := s.authService.Login(ctx, req.Username, req.Password, "grpc")
	if err != nil {
		s.logger.Error("Login failed", zap.Error(err))
		return nil, err
	}

	return &pb.LoginResponse{
		Token: token,
		Admin: &pb.AdminInfo{
			Id:       int64(admin.ID),
			Username: admin.Username,
			Nickname: admin.Nickname,
			Email:    admin.Email,
			Phone:    admin.Phone,
			Status:   int32(admin.Status),
		},
	}, nil
}

// GetProfile 获取管理员信息
func (s *AuthService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.AdminInfo, error) {
	s.logger.Info("gRPC GetProfile", zap.Int64("admin_id", req.AdminId))

	admin, err := s.authService.GetAdminByID(ctx, uint(req.AdminId))
	if err != nil {
		s.logger.Error("GetProfile failed", zap.Error(err))
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

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	s.logger.Info("gRPC ChangePassword")

	err := s.authService.ChangePassword(ctx, 1, req.OldPassword, req.NewPassword)
	if err != nil {
		s.logger.Error("ChangePassword failed", zap.Error(err))
		return &pb.ChangePasswordResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ChangePasswordResponse{
		Success: true,
		Message: "密码修改成功",
	}, nil
}

// Logout 登出
func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	s.logger.Info("gRPC Logout", zap.Int64("admin_id", req.AdminId))

	// TODO: 实现 Token 黑名单机制

	return &pb.LogoutResponse{
		Success: true,
		Message: "登出成功",
	}, nil
}

// InitAdmin 初始化管理员
func (s *AuthService) InitAdmin(ctx context.Context) error {
	if err := s.authService.InitAdmin(ctx); err != nil {
		logger.Error("初始化管理员失败", "error", err)
		return err
	}
	logger.Info("初始化管理员成功")
	return nil
}
