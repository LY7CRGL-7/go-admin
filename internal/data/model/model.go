package model

import (
	"time"

	"gorm.io/gorm"
)

// ==================== 基础模型 ====================

// BaseModel 基础字段
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TenantModel 含租户字段的基础模型
type TenantModel struct {
	BaseModel
	TenantID uint `gorm:"index;not null;default:0;comment:租户ID" json:"tenant_id"`
}

// ==================== 租户 ====================

// Tenant 租户
type Tenant struct {
	BaseModel
	Name         string     `gorm:"size:100;not null" json:"name"`
	Code         string     `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Domain       string     `gorm:"size:255" json:"domain"`
	ContactName  string     `gorm:"size:100" json:"contact_name"`
	ContactPhone string     `gorm:"size:20" json:"contact_phone"`
	Status       int8       `gorm:"default:1;comment:1启用 0禁用" json:"status"`
	MaxUsers     int        `gorm:"default:0;comment:0不限制" json:"max_users"`
	ExpireAt     *time.Time `json:"expire_at"`
}

func (Tenant) TableName() string { return "tenants" }

// ==================== 用户 ====================

// User 用户（管理员）
type User struct {
	TenantModel
	Username    string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password    string     `gorm:"size:255;not null" json:"-"`
	Nickname    string     `gorm:"size:100" json:"nickname"`
	Email       string     `gorm:"size:100;index" json:"email"`
	Phone       string     `gorm:"size:20" json:"phone"`
	Avatar      string     `gorm:"size:255" json:"avatar"`
	Status      int8       `gorm:"default:1;comment:1启用 0禁用 -1锁定" json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `gorm:"size:50" json:"last_login_ip"`
	Roles       []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (User) TableName() string { return "users" }

// ==================== RBAC ====================

// Role 角色
type Role struct {
	TenantModel
	Name        string       `gorm:"size:50;not null" json:"name"`
	Code        string       `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Description string       `gorm:"size:255" json:"description"`
	Status      int8         `gorm:"default:1" json:"status"`
	Sort        int          `gorm:"default:0" json:"sort"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

func (Role) TableName() string { return "roles" }

// Permission 权限
type Permission struct {
	BaseModel
	Name     string `gorm:"size:100;not null" json:"name"`
	Code     string `gorm:"size:100;uniqueIndex;not null" json:"code"`
	Type     string `gorm:"size:20;comment:menu/button/api" json:"type"`
	Path     string `gorm:"size:255" json:"path"`
	Method   string `gorm:"size:10" json:"method"`
	ParentID uint   `gorm:"default:0;index" json:"parent_id"`
	Sort     int    `gorm:"default:0" json:"sort"`
	Status   int8   `gorm:"default:1" json:"status"`
}

func (Permission) TableName() string { return "permissions" }

// UserRole 用户-角色关联
type UserRole struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"uniqueIndex:idx_user_role"`
	RoleID    uint      `gorm:"uniqueIndex:idx_user_role"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserRole) TableName() string { return "user_roles" }

// RolePermission 角色-权限关联
type RolePermission struct {
	ID           uint      `gorm:"primaryKey"`
	RoleID       uint      `gorm:"uniqueIndex:idx_role_perm"`
	PermissionID uint      `gorm:"uniqueIndex:idx_role_perm"`
	CreatedAt    time.Time `json:"created_at"`
}

func (RolePermission) TableName() string { return "role_permissions" }

// ==================== 审计日志 ====================

// AuditLog 审计日志
type AuditLog struct {
	BaseModel
	TenantID   uint   `gorm:"index" json:"tenant_id"`
	UserID     uint   `gorm:"index" json:"user_id"`
	Username   string `gorm:"size:50" json:"username"`
	Action     string `gorm:"size:100;not null;index" json:"action"`
	Resource   string `gorm:"size:100;index" json:"resource"`
	Method     string `gorm:"size:10" json:"method"`
	Path       string `gorm:"size:255" json:"path"`
	IP         string `gorm:"size:50" json:"ip"`
	UserAgent  string `gorm:"size:500" json:"user_agent"`
	StatusCode int    `json:"status_code"`
	Duration   int64  `json:"duration"` // 毫秒
}

func (AuditLog) TableName() string { return "audit_logs" }

// ==================== 登录记录 ====================

// LoginAttempt 登录尝试
type LoginAttempt struct {
	ID         uint      `gorm:"primaryKey"`
	Username   string    `gorm:"size:50;not null;index" json:"username"`
	IP         string    `gorm:"size:50" json:"ip"`
	Success    bool      `json:"success"`
	FailReason string    `gorm:"size:255" json:"fail_reason"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

func (LoginAttempt) TableName() string { return "login_attempts" }

// AutoMigrate 自动迁移所有表
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Tenant{},
		&User{},
		&Role{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
		&AuditLog{},
		&LoginAttempt{},
	)
}
