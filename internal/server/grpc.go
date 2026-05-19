package server

import (
	"context"
	"fmt"
	"net"

	"admin/internal/conf"
	pb "admin/proto/admin/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GRPCServer gRPC 服务器
type GRPCServer struct {
	server   *grpc.Server
	config   *conf.Config
	logger   *zap.Logger
	listener net.Listener
}

// NewGRPCServer 创建 gRPC 服务器
func NewGRPCServer(cfg *conf.Config, logger *zap.Logger) (*GRPCServer, error) {
	if !cfg.GRPC.Enabled {
		logger.Info("gRPC server is disabled")
		return nil, nil
	}

	// 创建监听器
	addr := fmt.Sprintf(":%d", cfg.GRPC.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor(logger)),
	)

	s := &GRPCServer{
		server:   grpcServer,
		config:   cfg,
		logger:   logger,
		listener: listener,
	}

	return s, nil
}

// RegisterServices 注册 gRPC 服务
func (s *GRPCServer) RegisterServices(
	authService pb.AuthServiceServer,
	adminService pb.AdminServiceServer,
	roleService pb.RoleServiceServer,
	permissionService pb.PermissionServiceServer,
	auditLogService pb.AuditLogServiceServer,
) {
	pb.RegisterAuthServiceServer(s.server, authService)
	pb.RegisterAdminServiceServer(s.server, adminService)
	pb.RegisterRoleServiceServer(s.server, roleService)
	pb.RegisterPermissionServiceServer(s.server, permissionService)
	pb.RegisterAuditLogServiceServer(s.server, auditLogService)

	s.logger.Info("gRPC services registered")
}

// Start 启动 gRPC 服务器
func (s *GRPCServer) Start() error {
	if s == nil {
		return nil
	}

	s.logger.Info("Starting gRPC server", zap.Int("port", s.config.GRPC.Port))

	if err := s.server.Serve(s.listener); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

// Stop 停止 gRPC 服务器
func (s *GRPCServer) Stop() {
	if s == nil {
		return
	}

	s.logger.Info("Stopping gRPC server")
	s.server.GracefulStop()
}

// loggerInterceptor gRPC 日志拦截器
func loggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Info("gRPC request",
			zap.String("method", info.FullMethod),
		)

		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error("gRPC error",
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
		}

		return resp, err
	}
}
