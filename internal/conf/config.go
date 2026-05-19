package conf

import "time"

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	GRPC       GRPCConfig       `yaml:"grpc"`
	Log        LogConfig        `yaml:"log"`
	JWT        JWTConfig        `yaml:"jwt"`
	AdminAuth  AdminAuthConfig  `yaml:"admin_auth"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	Security   SecurityConfig   `yaml:"security"`
	RateLimit  RateLimitConfig  `yaml:"rate_limit"`
	Audit      AuditConfig      `yaml:"audit"`
	Kafka      KafkaConfig      `yaml:"kafka"`
	MinIO      MinIOConfig      `yaml:"minio"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
}

type ServerConfig struct {
	Port         int           `yaml:"port"`
	Mode         string        `yaml:"mode"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type GRPCConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

type JWTConfig struct {
	Secret string        `yaml:"secret"`
	Expire time.Duration `yaml:"expire"`
	Issuer string        `yaml:"issuer"`
}

type AdminAuthConfig struct {
	InitAdmin InitAdminConfig `yaml:"init_admin"`
}

type InitAdminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Nickname string `yaml:"nickname"`
}

type DatabaseConfig struct {
	Driver          string        `yaml:"driver"`
	DSN             string        `yaml:"dsn"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

type SecurityConfig struct {
	Password    PasswordConfig    `yaml:"password"`
	Login       LoginConfig       `yaml:"login"`
	IPWhitelist []string          `yaml:"ip_whitelist"`
	CORS        CORSConfig        `yaml:"cors"`
}

type PasswordConfig struct {
	MinLength        int    `yaml:"min_length"`
	RequireUppercase bool   `yaml:"require_uppercase"`
	RequireLowercase bool   `yaml:"require_lowercase"`
	RequireNumber    bool   `yaml:"require_number"`
	RequireSpecial   bool   `yaml:"require_special"`
	SpecialChars     string `yaml:"special_chars"`
}

type LoginConfig struct {
	MaxAttempts     int           `yaml:"max_attempts"`
	LockoutDuration time.Duration `yaml:"lockout_duration"`
	AttemptWindow   time.Duration `yaml:"attempt_window"`
}

type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type RateLimitConfig struct {
	Enabled     bool `yaml:"enabled"`
	GlobalQPS   int  `yaml:"global_qps"`
	PerUserQPS  int  `yaml:"per_user_qps"`
	PerIPQPS    int  `yaml:"per_ip_qps"`
}

type AuditConfig struct {
	Enabled    bool   `yaml:"enabled"`
	LogFile    string `yaml:"log_file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

type KafkaConfig struct {
	Enabled       bool     `yaml:"enabled"`
	Brokers       []string `yaml:"brokers"`
	AuditLogTopic string   `yaml:"audit_log_topic"`
}

type MinIOConfig struct {
	Enabled         bool   `yaml:"enabled"`
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	UseSSL          bool   `yaml:"use_ssl"`
	BucketName      string `yaml:"bucket_name"`
}

type PrometheusConfig struct {
	Enabled     bool   `yaml:"enabled"`
	MetricsPath string `yaml:"metrics_path"`
}
