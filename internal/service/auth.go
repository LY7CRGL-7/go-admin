package service

import (
	"admin/internal/conf"
	"admin/internal/data/model"
	"admin/internal/middleware"
	"admin/internal/pkg/logger"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *conf.Config
}

func NewAuthService(db *gorm.DB, cfg *conf.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

// Login 管理员登录
func (s *AuthService) Login(username, password, ip string) (string, *model.Admin, error) {
	// 查询管理员
	var admin model.Admin
	if err := s.db.Where("username = ? AND status = 1", username).Preload("Roles").First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("用户名或密码错误")
		}
		logger.Error("查询管理员失败", "error", err)
		return "", nil, errors.New("系统错误")
	}
	
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}
	
	// 生成 Token
	token, err := middleware.GenerateToken(&s.cfg.JWT, admin.ID, admin.Username)
	if err != nil {
		logger.Error("生成 Token 失败", "error", err)
		return "", nil, errors.New("系统错误")
	}
	
	// 更新最后登录信息
	now := time.Now()
	s.db.Model(&admin).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
	})
	
	return token, &admin, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(adminID uint, oldPassword, newPassword string) error {
	// 查询管理员
	var admin model.Admin
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return errors.New("管理员不存在")
	}
	
	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}
	
	// 验证新密码强度
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}
	
	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("加密密码失败", "error", err)
		return errors.New("系统错误")
	}
	
	// 更新密码
	return s.db.Model(&admin).Update("password", string(hashedPassword)).Error
}

// validatePassword 验证密码强度
func (s *AuthService) validatePassword(password string) error {
	cfg := s.cfg.Security.Password
	
	if len(password) < cfg.MinLength {
		return errors.New("密码长度不能少于8位")
	}
	
	if cfg.RequireUppercase {
		hasUpper := false
		for _, c := range password {
			if c >= 'A' && c <= 'Z' {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			return errors.New("密码必须包含大写字母")
		}
	}
	
	if cfg.RequireLowercase {
		hasLower := false
		for _, c := range password {
			if c >= 'a' && c <= 'z' {
				hasLower = true
				break
			}
		}
		if !hasLower {
			return errors.New("密码必须包含小写字母")
		}
	}
	
	if cfg.RequireNumber {
		hasNumber := false
		for _, c := range password {
			if c >= '0' && c <= '9' {
				hasNumber = true
				break
			}
		}
		if !hasNumber {
			return errors.New("密码必须包含数字")
		}
	}
	
	if cfg.RequireSpecial {
		hasSpecial := false
		for _, c := range password {
			if containsSpecialChar(cfg.SpecialChars, c) {
				hasSpecial = true
				break
			}
		}
		if !hasSpecial {
			return errors.New("密码必须包含特殊字符")
		}
	}
	
	return nil
}

func containsSpecialChar(specialChars string, c rune) bool {
	for _, sc := range specialChars {
		if c == sc {
			return true
		}
	}
	return false
}

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// GetAdminByID 根据 ID 获取管理员信息
func (s *AuthService) GetAdminByID(adminID uint) (*model.Admin, error) {
	var admin model.Admin
	if err := s.db.Preload("Roles").First(&admin, adminID).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// InitAdmin 初始化管理员账号
func (s *AuthService) InitAdmin() error {
	var count int64
	s.db.Model(&model.Admin{}).Count(&count)
	
	if count > 0 {
		return nil // 已有管理员，无需初始化
	}
	
	// 创建初始管理员
	hashedPassword, err := HashPassword(s.cfg.AdminAuth.InitAdmin.Password)
	if err != nil {
		return err
	}
	
	admin := model.Admin{
		Username: s.cfg.AdminAuth.InitAdmin.Username,
		Password: hashedPassword,
		Nickname: s.cfg.AdminAuth.InitAdmin.Nickname,
		Status:   1,
	}
	
	return s.db.Create(&admin).Error
}
