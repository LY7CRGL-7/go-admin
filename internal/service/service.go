package service

import "github.com/google/wire"

// ProviderSet service 层依赖注入集合
var ProviderSet = wire.NewSet(
	NewAuthService,
	NewUserService,
	NewRoleService,
	NewPermissionService,
	NewTenantService,
	NewAuditService,
)
