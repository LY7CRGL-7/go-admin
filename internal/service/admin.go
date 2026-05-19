package service

import (
	"admin/internal/data"
	"admin/internal/data/model"
	"admin/internal/pkg/logger"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AdminService 管理员服务
type AdminService struct {
	adminRepo data.AdminRepository
	roleRepo  data.RoleRepository
}

// NewAdminService 创建管理员服务
func NewAdminService(adminRepo data.AdminRepository, roleRepo data.RoleRepository) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
		roleRepo:  roleRepo,
	}
}

// CreateAdmin 创建管理员
func (s *AdminService) CreateAdmin(ctx context.Context, admin *model.Admin, roleIDs []uint) error {
	// 检查用户名是否已存在
	existing, _ := s.adminRepo.GetByUsername(ctx, admin.Username)
	if existing != nil {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("加密密码失败", "error", err)
		return errors.New("系统错误")
	}
	admin.Password = string(hashedPassword)

	// 创建管理员
	if err := s.adminRepo.Create(ctx, admin); err != nil {
		return err
	}

	// 分配角色
	if len(roleIDs) > 0 {
		return s.adminRepo.AssignRoles(ctx, admin.ID, roleIDs)
	}

	return nil
}

// GetAdmin 获取管理员
func (s *AdminService) GetAdmin(ctx context.Context, id uint) (*model.Admin, error) {
	return s.adminRepo.GetByID(ctx, id)
}

// UpdateAdmin 更新管理员
func (s *AdminService) UpdateAdmin(ctx context.Context, admin *model.Admin) error {
	existing, err := s.adminRepo.GetByID(ctx, admin.ID)
	if err != nil {
		return errors.New("管理员不存在")
	}

	// 保持密码不变
	admin.Password = existing.Password
	return s.adminRepo.Update(ctx, admin)
}

// DeleteAdmin 删除管理员
func (s *AdminService) DeleteAdmin(ctx context.Context, id uint) error {
	return s.adminRepo.Delete(ctx, id)
}

// ListAdmins 列出管理员
func (s *AdminService) ListAdmins(ctx context.Context, page, pageSize int) ([]*model.Admin, int64, error) {
	return s.adminRepo.List(ctx, page, pageSize)
}

// AssignRoles 分配角色
func (s *AdminService) AssignRoles(ctx context.Context, adminID uint, roleIDs []uint) error {
	// 验证管理员是否存在
	_, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		return errors.New("管理员不存在")
	}

	// 验证角色是否存在
	for _, roleID := range roleIDs {
		_, err := s.roleRepo.GetByID(ctx, roleID)
		if err != nil {
			return errors.New("角色不存在")
		}
	}

	return s.adminRepo.AssignRoles(ctx, adminID, roleIDs)
}
