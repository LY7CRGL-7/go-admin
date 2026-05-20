package conf

import (
	"encoding/json"
	"time"
)

// Duration 支持从字符串（如 "10s"、"24h"）解析的 time.Duration 包装
// Kratos config.Scan 底层用 JSON 反序列化，原生 time.Duration 不支持字符串
type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case float64:
		*d = Duration(time.Duration(int64(val)))
	case string:
		dur, err := time.ParseDuration(val)
		if err != nil {
			return err
		}
		*d = Duration(dur)
	default:
		*d = 0
	}
	return nil
}

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
	Network string   `json:"network" yaml:"network"`
	Addr    string   `json:"addr" yaml:"addr"`
	Timeout Duration `json:"timeout" yaml:"timeout"`
}

// Data 数据源配置
type Data struct {
	Database Database `json:"database" yaml:"database"`
	Redis    Redis    `json:"redis" yaml:"redis"`
}

// Database 数据库配置
type Database struct {
	Driver      string   `json:"driver" yaml:"driver"`             // mysql | postgres
	Source      string   `json:"source" yaml:"source"`             // DSN
	MaxIdle     int      `json:"max_idle" yaml:"max_idle"`         // 最大空闲连接
	MaxOpen     int      `json:"max_open" yaml:"max_open"`         // 最大打开连接
	MaxLifetime Duration `json:"max_lifetime" yaml:"max_lifetime"` // 连接最大存活时间
}

// Redis 缓存配置
type Redis struct {
	Addr         string   `json:"addr" yaml:"addr"`
	Password     string   `json:"password" yaml:"password"`
	DB           int      `json:"db" yaml:"db"`
	DialTimeout  Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout Duration `json:"write_timeout" yaml:"write_timeout"`
}

// Auth 认证配置
type Auth struct {
	JWTSecret     string    `json:"jwt_secret" yaml:"jwt_secret"`
	TokenExpire   Duration  `json:"token_expire" yaml:"token_expire"`
	RefreshExpire Duration  `json:"refresh_expire" yaml:"refresh_expire"`
	InitAdmin     InitAdmin `json:"init_admin" yaml:"init_admin"`
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
