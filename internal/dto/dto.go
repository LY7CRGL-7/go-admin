package dto

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	Admin *AdminInfo `json:"admin"`
}

// AdminInfo 管理员信息
type AdminInfo struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Avatar      string `json:"avatar"`
	Status      int8   `json:"status"`
	LastLoginAt *string `json:"last_login_at"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// CreateAdminRequest 创建管理员请求
type CreateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RoleIDs  []uint `json:"role_ids"`
}

// UpdateAdminRequest 更新管理员请求
type UpdateAdminRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   *int8  `json:"status"`
	RoleIDs  []uint `json:"role_ids"`
}

// AdminListResponse 管理员列表响应
type AdminListResponse struct {
	Total int64        `json:"total"`
	List  []AdminInfo `json:"list"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      *int8  `json:"status"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
	Sort        int    `json:"sort"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Type     string `json:"type"`
	Path     string `json:"path"`
	Method   string `json:"method"`
	ParentID uint   `json:"parent_id"`
	Sort     int    `json:"sort"`
}

// AssignPermissionsRequest 分配权限请求
type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

// AuditLogResponse 审计日志响应
type AuditLogResponse struct {
	ID         uint   `json:"id"`
	AdminID    uint   `json:"admin_id"`
	AdminName  string `json:"admin_name"`
	Action     string `json:"action"`
	Resource   string `json:"resource"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	Status     int    `json:"status"`
	Duration   int64  `json:"duration"`
	CreatedAt  string `json:"created_at"`
}

// PaginatedRequest 分页请求
type PaginatedRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}
