package service

import (
	"context"

	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// UserService 用户管理服务
type UserService struct {
	v1.UnimplementedUserServiceServer
	uc  *biz.UserUsecase
	log *log.Helper
}

// NewUserService 创建用户服务
func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/user")),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.UserInfo, error) {
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}

	// 设置租户 ID
	if claims := GetClaimsFromContext(ctx); claims != nil {
		user.TenantID = claims.TenantID
	}

	if err := s.uc.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// 分配角色
	if len(req.RoleIds) > 0 {
		var roleIDs []uint
		for _, id := range req.RoleIds {
			roleIDs = append(roleIDs, uint(id))
		}
		if err := s.uc.AssignRoles(ctx, user.ID, roleIDs); err != nil {
			return nil, err
		}
	}

	return userToProto(user), nil
}

func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.UserInfo, error) {
	user, err := s.uc.GetUser(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}
	return userToProto(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UserInfo, error) {
	user, err := s.uc.GetUser(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Status != 0 {
		user.Status = int8(req.Status)
	}

	if err := s.uc.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	// 更新角色
	if len(req.RoleIds) > 0 {
		var roleIDs []uint
		for _, id := range req.RoleIds {
			roleIDs = append(roleIDs, uint(id))
		}
		_ = s.uc.AssignRoles(ctx, user.ID, roleIDs)
	}

	return userToProto(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteReply, error) {
	if err := s.uc.DeleteUser(ctx, uint(req.Id)); err != nil {
		return nil, err
	}
	return &v1.DeleteReply{Success: true}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersReply, error) {
	var tenantID uint
	if claims := GetClaimsFromContext(ctx); claims != nil {
		tenantID = claims.TenantID
	}

	users, total, err := s.uc.ListUsers(ctx, tenantID, int(req.Page), int(req.PageSize), req.Keyword, int(req.Status))
	if err != nil {
		return nil, err
	}

	var items []*v1.UserInfo
	for _, u := range users {
		items = append(items, userToProto(u))
	}

	return &v1.ListUsersReply{Items: items, Total: total}, nil
}

func (s *UserService) AssignRoles(ctx context.Context, req *v1.AssignRolesRequest) (*v1.AssignRolesReply, error) {
	var roleIDs []uint
	for _, id := range req.RoleIds {
		roleIDs = append(roleIDs, uint(id))
	}
	if err := s.uc.AssignRoles(ctx, uint(req.UserId), roleIDs); err != nil {
		return nil, err
	}
	return &v1.AssignRolesReply{Success: true}, nil
}
