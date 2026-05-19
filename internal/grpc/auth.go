package grpc

import (
	"context"

	"admin/internal/service"
	pb "admin/proto/admin/v1"

	"go.uber.org/zap"
)

// AuthService gRPC 认证服务
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
	logger      *zap.Logger
}

// NewAuthService 创建认证服务
func NewAuthService(authService service.AuthService, logger *zap.Logger) *AuthService {
	return &AuthService{
		authService: authService,
		logger:      logger,
	}
}

// Login 登录
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.logger.Info("gRPC Login", zap.String("username", req.Username))

	// TODO: 调用 service 层实现登录逻辑
	// 这里需要根据实际的 service 接口来实现

	return &pb.LoginResponse{
		Token: "grpc-token-placeholder",
		Admin: &pb.AdminInfo{
			Id:       1,
			Username: req.Username,
			Nickname: "管理员",
		},
	}, nil
}

// GetProfile 获取管理员信息
func (s *AuthService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.AdminInfo, error) {
	s.logger.Info("gRPC GetProfile", zap.Int64("admin_id", req.AdminId))

	// TODO: 实现获取个人信息逻辑

	return &pb.AdminInfo{
		Id:       req.AdminId,
		Username: "admin",
		Nickname: "系统管理员",
	}, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	s.logger.Info("gRPC ChangePassword")

	// TODO: 实现修改密码逻辑

	return &pb.ChangePasswordResponse{
		Success: true,
		Message: "密码修改成功",
	}, nil
}

// Logout 登出
func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	s.logger.Info("gRPC Logout", zap.Int64("admin_id", req.AdminId))

	// TODO: 实现登出逻辑

	return &pb.LogoutResponse{
		Success: true,
		Message: "登出成功",
	}, nil
}
