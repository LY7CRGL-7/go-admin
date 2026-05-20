package biz

import (
	"context"

	"admin/internal/data/model"

	"github.com/google/wire"
)

// ProviderSet biz 层依赖注入集合
var ProviderSet = wire.NewSet(
	NewUserUsecase,
	NewRoleUsecase,
	NewTenantUsecase,
)

// ==================== Repo 接口定义 ====================
// 由 data 层实现，biz 层通过接口依赖（依赖倒置）

// UserRepo 用户仓储接口
type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, tenantID uint, page, pageSize int, keyword string, status int) ([]*model.User, int64, error)
	AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error
}

// RoleRepo 角色仓储接口
type RoleRepo interface {
	Create(ctx context.Context, role *model.Role) error
	GetByID(ctx context.Context, id uint) (*model.Role, error)
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, tenantID uint, page, pageSize int) ([]*model.Role, int64, error)
	AssignPermissions(ctx context.Context, roleID uint, permIDs []uint) error
}

// PermissionRepo 权限仓储接口
type PermissionRepo interface {
	List(ctx context.Context) ([]*model.Permission, error)
	GetTree(ctx context.Context) ([]*model.Permission, error)
	GetByRoleIDs(ctx context.Context, roleIDs []uint) ([]*model.Permission, error)
}

// TenantRepo 租户仓储接口
type TenantRepo interface {
	Create(ctx context.Context, tenant *model.Tenant) error
	GetByID(ctx context.Context, id uint) (*model.Tenant, error)
	GetByCode(ctx context.Context, code string) (*model.Tenant, error)
	Update(ctx context.Context, tenant *model.Tenant) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, pageSize int) ([]*model.Tenant, int64, error)
}

// AuditRepo 审计日志仓储接口
type AuditRepo interface {
	Create(ctx context.Context, log *model.AuditLog) error
	List(ctx context.Context, tenantID uint, page, pageSize int, userID uint, action string) ([]*model.AuditLog, int64, error)
}
