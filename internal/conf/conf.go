package conf

import "time"

// Bootstrap 应用总配置
type Bootstrap struct {
	Server Server `json:"server" yaml:"server"`
	Data   Data   `json:"data" yaml:"data"`
	Auth   Auth   `json:"auth" yaml:"auth"`
	Tenant Tenant `json:"tenant" yaml:"tenant"`
}

// Server 服务器配置
type Server struct {
	HTTP ServerItem `json:"http" yaml:"http"`
	GRPC ServerItem `json:"grpc" yaml:"grpc"`
}

// ServerItem 单个服务器配置
type ServerItem struct {
	Network string        `json:"network" yaml:"network"`
	Addr    string        `json:"addr" yaml:"addr"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

// Data 数据源配置
type Data struct {
	Database Database `json:"database" yaml:"database"`
	Redis    Redis    `json:"redis" yaml:"redis"`
}

// Database 数据库配置
type Database struct {
	Driver      string        `json:"driver" yaml:"driver"`           // mysql | postgres
	Source      string        `json:"source" yaml:"source"`           // DSN
	MaxIdle     int           `json:"max_idle" yaml:"max_idle"`       // 最大空闲连接
	MaxOpen     int           `json:"max_open" yaml:"max_open"`       // 最大打开连接
	MaxLifetime time.Duration `json:"max_lifetime" yaml:"max_lifetime"` // 连接最大存活时间
}

// Redis 缓存配置
type Redis struct {
	Addr         string        `json:"addr" yaml:"addr"`
	Password     string        `json:"password" yaml:"password"`
	DB           int           `json:"db" yaml:"db"`
	DialTimeout  time.Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
}

// Auth 认证配置
type Auth struct {
	JWTSecret      string        `json:"jwt_secret" yaml:"jwt_secret"`
	TokenExpire    time.Duration `json:"token_expire" yaml:"token_expire"`
	RefreshExpire  time.Duration `json:"refresh_expire" yaml:"refresh_expire"`
	InitAdmin      InitAdmin     `json:"init_admin" yaml:"init_admin"`
}

// InitAdmin 初始管理员配置
type InitAdmin struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Nickname string `json:"nickname" yaml:"nickname"`
}

// Tenant 多租户配置
type Tenant struct {
	Enabled           bool   `json:"enabled" yaml:"enabled"`
	DefaultTenantCode string `json:"default_tenant_code" yaml:"default_tenant_code"`
}
