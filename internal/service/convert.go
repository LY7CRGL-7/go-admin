package service

import (
	"context"
	"errors"

	"admin/internal/biz"
	"admin/internal/data/model"

	v1 "admin/api/admin/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Context Key 定义
type contextKey string

const claimsKey contextKey = "claims"

// ErrUnauthorized 未认证错误
var ErrUnauthorized = errors.New("未认证，请先登录")

// SetClaimsToContext 设置认证信息到 context
func SetClaimsToContext(ctx context.Context, claims *biz.TokenClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// GetClaimsFromContext 从 context 获取认证信息
func GetClaimsFromContext(ctx context.Context) *biz.TokenClaims {
	val := ctx.Value(claimsKey)
	if val == nil {
		return nil
	}
	claims, ok := val.(*biz.TokenClaims)
	if !ok {
		return nil
	}
	return claims
}

// ==================== Proto 类型转换 ====================

func userToProto(u *model.User) *v1.UserInfo {
	if u == nil {
		return nil
	}
	info := &v1.UserInfo{
		Id:        int64(u.ID),
		Username:  u.Username,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Status:    int32(u.Status),
		TenantId:  int64(u.TenantID),
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
	for _, role := range u.Roles {
		info.Roles = append(info.Roles, role.Code)
	}
	return info
}

func roleToProto(r *model.Role) *v1.RoleInfo {
	if r == nil {
		return nil
	}
	info := &v1.RoleInfo{
		Id:          int64(r.ID),
		Name:        r.Name,
		Code:        r.Code,
		Description: r.Description,
		Status:      int32(r.Status),
		Sort:        int32(r.Sort),
		TenantId:    int64(r.TenantID),
		CreatedAt:   timestamppb.New(r.CreatedAt),
	}
	for _, p := range r.Permissions {
		info.Permissions = append(info.Permissions, permToProto(&p))
	}
	return info
}

func permToProto(p *model.Permission) *v1.PermissionInfo {
	if p == nil {
		return nil
	}
	return &v1.PermissionInfo{
		Id:       int64(p.ID),
		Name:     p.Name,
		Code:     p.Code,
		Type:     p.Type,
		Path:     p.Path,
		Method:   p.Method,
		ParentId: int64(p.ParentID),
		Sort:     int32(p.Sort),
		Status:   int32(p.Status),
	}
}

func tenantToProto(t *model.Tenant) *v1.TenantInfo {
	if t == nil {
		return nil
	}
	info := &v1.TenantInfo{
		Id:           int64(t.ID),
		Name:         t.Name,
		Code:         t.Code,
		Domain:       t.Domain,
		ContactName:  t.ContactName,
		ContactPhone: t.ContactPhone,
		Status:       int32(t.Status),
		MaxUsers:     int32(t.MaxUsers),
		CreatedAt:    timestamppb.New(t.CreatedAt),
	}
	if t.ExpireAt != nil {
		info.ExpireAt = timestamppb.New(*t.ExpireAt)
	}
	return info
}

func auditToProto(a *model.AuditLog) *v1.AuditLogInfo {
	if a == nil {
		return nil
	}
	return &v1.AuditLogInfo{
		Id:         int64(a.ID),
		UserId:     int64(a.UserID),
		Username:   a.Username,
		Action:     a.Action,
		Resource:   a.Resource,
		Method:     a.Method,
		Path:       a.Path,
		Ip:         a.IP,
		Status:     int32(a.StatusCode),
		DurationMs: a.Duration,
		CreatedAt:  timestamppb.New(a.CreatedAt),
	}
}
