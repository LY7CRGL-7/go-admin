package model

import (
	"time"
)

// Admin 管理员模型
type Admin struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Username      string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password      string    `gorm:"size:255;not null" json:"-"` // 不返回给前端
	Nickname      string    `gorm:"size:100" json:"nickname"`
	Email         string    `gorm:"size:100;uniqueIndex" json:"email"`
	Phone         string    `gorm:"size:20" json:"phone"`
	Avatar        string    `gorm:"size:255" json:"avatar"`
	Status        int8      `gorm:"default:1;comment:1启用 0禁用" json:"status"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	LastLoginIP   string    `gorm:"size:50" json:"last_login_ip"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// 关联
	Roles []Role `gorm:"many2many:admin_roles;" json:"roles,omitempty"`
}

// Role 角色模型
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description string    `gorm:"size:255" json:"description"`
	Status      int8      `gorm:"default:1;comment:1启用 0禁用" json:"status"`
	Sort        int       `gorm:"default:0" json:"sort"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// 关联
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission 权限模型
type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Code      string    `gorm:"uniqueIndex;size:100;not null" json:"code"`
	Type      string    `gorm:"size:20;comment:menu-菜单 button-按钮 api-API" json:"type"`
	Path      string    `gorm:"size:255" json:"path"`
	Method    string    `gorm:"size:10" json:"method"`
	ParentID  uint      `gorm:"default:0" json:"parent_id"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AdminRole 管理员角色关联表
type AdminRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AdminID   uint      `gorm:"uniqueIndex:idx_admin_role" json:"admin_id"`
	RoleID    uint      `gorm:"uniqueIndex:idx_admin_role" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"uniqueIndex:idx_role_permission" json:"role_id"`
	PermissionID uint      `gorm:"uniqueIndex:idx_role_permission" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuditLog 审计日志模型
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AdminID    uint      `gorm:"index" json:"admin_id"`
	AdminName  string    `gorm:"size:50" json:"admin_name"`
	Action     string    `gorm:"size:100;not null;index" json:"action"`
	Resource   string    `gorm:"size:100;index" json:"resource"`
	ResourceID string    `gorm:"size:100" json:"resource_id"`
	Method     string    `gorm:"size:10" json:"method"`
	Path       string    `gorm:"size:255" json:"path"`
	IP         string    `gorm:"size:50" json:"ip"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	Request    string    `gorm:"type:text" json:"request"`
	Response   string    `gorm:"type:text" json:"response"`
	Status     int       `json:"status"`
	Error      string    `gorm:"type:text" json:"error"`
	Duration   int64     `json:"duration"` // 毫秒
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;not null;index" json:"username"`
	IP        string    `gorm:"size:50" json:"ip"`
	Success   bool      `json:"success"`
	FailReason string   `gorm:"size:255" json:"fail_reason"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admins"
}

func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (AdminRole) TableName() string {
	return "admin_roles"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

func (LoginAttempt) TableName() string {
	return "login_attempts"
}
