package biz

import (
	"context"
	"errors"
	"time"

	"admin/internal/conf"
	"admin/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecase 用户业务逻辑
type UserUsecase struct {
	repo     UserRepo
	roleRepo RoleRepo
	permRepo PermissionRepo
	conf     *conf.Auth
	log      *log.Helper
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepo, roleRepo RoleRepo, permRepo PermissionRepo, c *conf.Auth, logger log.Logger) *UserUsecase {
	uc := &UserUsecase{
		repo:     repo,
		roleRepo: roleRepo,
		permRepo: permRepo,
		conf:     c,
		log:      log.NewHelper(log.With(logger, "module", "biz/user")),
	}

	// 自动初始化管理员账号
	if c.InitAdmin.Username != "" {
		if err := uc.InitAdmin(context.Background(), c.InitAdmin.Username, c.InitAdmin.Password, c.InitAdmin.Nickname); err != nil {
			uc.log.Errorf("failed to init admin: %v", err)
		}
	}

	return uc
}

// ==================== 认证 ====================

// TokenClaims JWT 自定义 Claims
type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	TenantID uint   `json:"tenant_id"`
	jwt.RegisteredClaims
}

// Login 用户登录
func (uc *UserUsecase) Login(ctx context.Context, username, password string) (token, refreshToken string, user *model.User, err error) {
	user, err = uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return "", "", nil, errors.New("账号已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", nil, errors.New("用户名或密码错误")
	}

	// 生成 Token
	token, err = uc.generateToken(user, time.Duration(uc.conf.TokenExpire))
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, err = uc.generateToken(user, time.Duration(uc.conf.RefreshExpire))
	if err != nil {
		return "", "", nil, err
	}

	// 更新登录时间
	now := time.Now()
	user.LastLoginAt = &now
	_ = uc.repo.Update(ctx, user)

	return token, refreshToken, user, nil
}

// ValidateToken 验证 Token
func (uc *UserUsecase) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(uc.conf.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的 Token")
	}

	return claims, nil
}

// RefreshToken 刷新 Token
func (uc *UserUsecase) RefreshToken(ctx context.Context, refreshTokenStr string) (string, string, error) {
	claims, err := uc.ValidateToken(refreshTokenStr)
	if err != nil {
		return "", "", errors.New("无效的刷新令牌")
	}

	user, err := uc.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", "", err
	}

	token, err := uc.generateToken(user, time.Duration(uc.conf.TokenExpire))
	if err != nil {
		return "", "", err
	}

	newRefresh, err := uc.generateToken(user, time.Duration(uc.conf.RefreshExpire))
	if err != nil {
		return "", "", err
	}

	return token, newRefresh, nil
}

// ChangePassword 修改密码
func (uc *UserUsecase) ChangePassword(ctx context.Context, userID uint, oldPwd, newPwd string) error {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd)); err != nil {
		return errors.New("原密码错误")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	return uc.repo.Update(ctx, user)
}

// generateToken 生成 JWT Token
func (uc *UserUsecase) generateToken(user *model.User, expire time.Duration) (string, error) {
	if expire == 0 {
		expire = 24 * time.Hour
	}
	claims := &TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		TenantID: user.TenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.conf.JWTSecret))
}

// ==================== CRUD ====================

// CreateUser 创建用户
func (uc *UserUsecase) CreateUser(ctx context.Context, user *model.User) error {
	return uc.repo.Create(ctx, user)
}

// GetUser 获取用户
func (uc *UserUsecase) GetUser(ctx context.Context, id uint) (*model.User, error) {
	return uc.repo.GetByID(ctx, id)
}

// UpdateUser 更新用户
func (uc *UserUsecase) UpdateUser(ctx context.Context, user *model.User) error {
	return uc.repo.Update(ctx, user)
}

// DeleteUser 删除用户
func (uc *UserUsecase) DeleteUser(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

// ListUsers 列表查询用户
func (uc *UserUsecase) ListUsers(ctx context.Context, tenantID uint, page, pageSize int, keyword string, status int) ([]*model.User, int64, error) {
	return uc.repo.List(ctx, tenantID, page, pageSize, keyword, status)
}

// AssignRoles 分配角色
func (uc *UserUsecase) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	return uc.repo.AssignRoles(ctx, userID, roleIDs)
}

// ==================== RBAC 鉴权 ====================

// GetUserPermissions 获取用户的所有权限编码
func (uc *UserUsecase) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var roleIDs []uint
	for _, role := range user.Roles {
		roleIDs = append(roleIDs, role.ID)
	}
	if len(roleIDs) == 0 {
		return nil, nil
	}

	perms, err := uc.permRepo.GetByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	var codes []string
	for _, p := range perms {
		codes = append(codes, p.Code)
	}
	return codes, nil
}

// HasPermission 检查用户是否拥有指定权限
func (uc *UserUsecase) HasPermission(ctx context.Context, userID uint, permCode string) (bool, error) {
	codes, err := uc.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, c := range codes {
		if c == permCode || c == "*" {
			return true, nil
		}
	}
	return false, nil
}

// InitAdmin 初始化管理员账号（首次启动时）
func (uc *UserUsecase) InitAdmin(ctx context.Context, username, password, nickname string) error {
	_, err := uc.repo.GetByUsername(ctx, username)
	if err == nil {
		return nil // 已存在，跳过
	}

	admin := &model.User{
		Username: username,
		Password: password,
		Nickname: nickname,
		Status:   1,
	}
	return uc.repo.Create(ctx, admin)
}
