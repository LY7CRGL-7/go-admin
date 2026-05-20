package middleware

import (
	"context"
	"strings"

	"admin/internal/biz"
	"admin/internal/service"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// JWTAuth JWT 认证中间件
// 从请求头提取 Bearer Token，验证后注入用户信息到 Context
func JWTAuth(uc *biz.UserUsecase) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 从 transport 获取请求头
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.Unauthorized("UNAUTHORIZED", "no transport context")
			}

			// 提取 Authorization 头
			authHeader := tr.RequestHeader().Get("Authorization")
			if authHeader == "" {
				return nil, errors.Unauthorized("UNAUTHORIZED", "missing authorization header")
			}

			// 解析 Bearer Token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return nil, errors.Unauthorized("UNAUTHORIZED", "invalid authorization format")
			}

			tokenStr := parts[1]

			// 验证 Token
			claims, err := uc.ValidateToken(tokenStr)
			if err != nil {
				return nil, errors.Unauthorized("UNAUTHORIZED", "invalid token: "+err.Error())
			}

			// 注入用户信息到 Context
			ctx = service.SetClaimsToContext(ctx, claims)

			return handler(ctx, req)
		}
	}
}

// SkipAuth 白名单路由（不需要认证的接口）
var SkipAuthRoutes = map[string]bool{
	"/admin.v1.AuthService/Login":        true,
	"/admin.v1.AuthService/RefreshToken": true,
}

// JWTAuthWithSkip 带白名单的 JWT 认证中间件
func JWTAuthWithSkip(uc *biz.UserUsecase) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			// 检查白名单
			operation := tr.Operation()
			if SkipAuthRoutes[operation] {
				return handler(ctx, req)
			}

			// 执行认证
			return JWTAuth(uc)(handler)(ctx, req)
		}
	}
}

// TenantContext 多租户上下文中间件
// 从认证信息中提取租户 ID，注入到 Context
func TenantContext() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 租户 ID 已在 JWT 认证中通过 claims 注入
			// 此中间件可扩展为从请求头 X-Tenant-ID 获取
			return handler(ctx, req)
		}
	}
}
