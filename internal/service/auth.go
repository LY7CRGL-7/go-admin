package service

import (
	"context"

	v1 "admin/api/admin/v1"
	"admin/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// AuthService 认证服务
type AuthService struct {
	v1.UnimplementedAuthServiceServer
	uc  *biz.UserUsecase
	log *log.Helper
}

// NewAuthService 创建认证服务
func NewAuthService(uc *biz.UserUsecase, logger log.Logger) *AuthService {
	return &AuthService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/auth")),
	}
}

func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	token, refreshToken, user, err := s.uc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token:        token,
		RefreshToken: refreshToken,
		User:         userToProto(user),
	}, nil
}

func (s *AuthService) Logout(_ context.Context, _ *v1.LogoutRequest) (*v1.LogoutReply, error) {
	return &v1.LogoutReply{Success: true}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, _ *v1.GetProfileRequest) (*v1.UserInfo, error) {
	claims := GetClaimsFromContext(ctx)
	if claims == nil {
		return nil, ErrUnauthorized
	}
	user, err := s.uc.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	return userToProto(user), nil
}

func (s *AuthService) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.ChangePasswordReply, error) {
	claims := GetClaimsFromContext(ctx)
	if claims == nil {
		return nil, ErrUnauthorized
	}
	err := s.uc.ChangePassword(ctx, claims.UserID, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}
	return &v1.ChangePasswordReply{Success: true}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenReply, error) {
	token, refresh, err := s.uc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &v1.RefreshTokenReply{
		Token:        token,
		RefreshToken: refresh,
	}, nil
}
