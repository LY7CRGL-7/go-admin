package server

import (
	v1 "admin/api/admin/v1"
	"admin/internal/biz"
	"admin/internal/conf"
	mw "admin/internal/middleware"
	"admin/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer 创建 gRPC 服务器
func NewGRPCServer(
	c *conf.Server,
	logger log.Logger,
	uc *biz.UserUsecase,
	authSvc *service.AuthService,
	userSvc *service.UserService,
	roleSvc *service.RoleService,
	permSvc *service.PermissionService,
	tenantSvc *service.TenantService,
	auditSvc *service.AuditService,
) *grpc.Server {
	var opts []grpc.ServerOption

	// 中间件链
	opts = append(opts, grpc.Middleware(
		recovery.Recovery(),
		logging.Server(logger),
		mw.JWTAuthWithSkip(uc),
		mw.TenantContext(),
	))

	if c.GRPC.Addr != "" {
		opts = append(opts, grpc.Address(c.GRPC.Addr))
	}
	if c.GRPC.Timeout > 0 {
		opts = append(opts, grpc.Timeout(c.GRPC.Timeout))
	}

	srv := grpc.NewServer(opts...)

	// 注册所有 gRPC 服务
	v1.RegisterAuthServiceServer(srv, authSvc)
	v1.RegisterUserServiceServer(srv, userSvc)
	v1.RegisterRoleServiceServer(srv, roleSvc)
	v1.RegisterPermissionServiceServer(srv, permSvc)
	v1.RegisterTenantServiceServer(srv, tenantSvc)
	v1.RegisterAuditLogServiceServer(srv, auditSvc)

	return srv
}
