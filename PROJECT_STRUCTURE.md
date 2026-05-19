# 安全管理端项目结构说明

## 📁 完整目录结构

```
admin/
├── cmd/
│   └── admin/
│       ├── conf/
│       │   ├── config.yaml                      # 开发环境配置
│       │   └── config.production.example.yaml   # 生产环境配置示例
│       └── main.go                              # 主程序入口
├── internal/
│   ├── conf/
│   │   └── config.go                            # 配置结构体定义
│   ├── data/
│   │   ├── model/
│   │   │   └── model.go                        # 数据模型（Admin, Role, Permission等）
│   │   ├── database.go                          # PostgreSQL 数据库连接
│   │   └── redis.go                            # Redis 连接
│   ├── dto/
│   │   └── dto.go                              # 数据传输对象
│   ├── handler/
│   │   ├── auth.go                             # 认证相关处理器（登录、登出、修改密码）
│   │   ├── admin.go                            # 管理员管理处理器
│   │   └── common.go                           # 通用处理器（角色、权限、审计日志）
│   ├── middleware/
│   │   ├── auth.go                             # JWT 认证中间件
│   │   ├── rbac.go                             # RBAC 权限控制中间件
│   │   ├── security.go                         # 安全中间件（限流、登录限制、IP白名单）
│   │   └── audit.go                            # 审计日志中间件
│   ├── pkg/
│   │   └── logger/
│   │       └── logger.go                       # 日志工具（基于 Zap）
│   ├── service/
│   │   └── auth.go                             # 认证业务逻辑
│   └── server/
│       └── http.go                             # HTTP 服务器和路由配置
├── deploy/
│   └── k8s.yaml                                # Kubernetes 部署配置
├── docs/
│   └── init.sql                                # 数据库初始化 SQL
├── Dockerfile                                  # Docker 镜像构建文件
├── Makefile                                    # 构建和管理脚本
├── .gitignore                                  # Git 忽略配置
├── go.mod                                      # Go 模块依赖
├── README.md                                   # 项目说明文档
├── QUICKSTART.md                               # 快速启动指南
└── SECURITY_CHECKLIST.md                       # 安全检查清单
```

## 🎯 核心模块说明

### 1. 认证模块 (Auth)
- **文件**: `internal/handler/auth.go`, `internal/service/auth.go`
- **功能**: 
  - 管理员登录/登出
  - JWT Token 生成和验证
  - 密码修改
  - 密码强度验证
  - 初始化管理员账号

### 2. 中间件模块 (Middleware)
- **JWT 认证** (`middleware/auth.go`): 
  - Token 验证
  - Claims 解析
  - 用户信息注入
  
- **RBAC 权限** (`middleware/rbac.go`):
  - 基于角色的访问控制
  - 接口级权限验证
  - 超级管理员特殊处理
  
- **安全防护** (`middleware/security.go`):
  - 登录失败限制
  - 账号/IP 锁定
  - 多维度限流（全局/用户/IP）
  - IP 白名单
  - 客户端 IP 获取
  
- **审计日志** (`middleware/audit.go`):
  - 请求/响应记录
  - 操作者信息记录
  - 异步日志写入
  - 操作耗时统计

### 3. 数据模型 (Model)
- **Admin**: 管理员信息
- **Role**: 角色信息
- **Permission**: 权限信息
- **AdminRole**: 管理员-角色关联
- **RolePermission**: 角色-权限关联
- **AuditLog**: 审计日志
- **LoginAttempt**: 登录尝试记录

### 4. 业务处理器 (Handler)
- **AuthHandler**: 认证相关接口
- **AdminHandler**: 管理员 CRUD
- **通用 Handler**: 角色、权限、审计日志管理

### 5. 配置管理 (Config)
- **开发配置**: `cmd/admin/conf/config.yaml`
- **生产配置**: `cmd/admin/conf/config.production.example.yaml`
- **配置项**:
  - 服务器配置
  - 数据库配置
  - Redis 配置
  - JWT 配置
  - 安全策略
  - 限流配置
  - 审计配置

## 🔒 安全特性实现

### 1. 密码安全
```go
// bcrypt 加密
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 密码验证
err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
```

### 2. JWT 认证
```go
// Token 生成
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString([]byte(secret))

// Token 验证
token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(secret), nil
})
```

### 3. 限流实现
```go
// Redis 滑动窗口限流
pipe.ZRemRangeByScore(ctx, key, "0", now)
pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: uniqueID})
count := rdb.ZCard(ctx, key).Val()
```

### 4. 审计日志
```go
// 异步写入审计日志
go func() {
    db.Create(&model.AuditLog{
        AdminID:   adminID,
        Action:    action,
        Request:   requestBody,
        Response:  responseBody,
        // ...
    })
}()
```

## 🚀 启动流程

1. **加载配置**: 读取 YAML 配置文件
2. **初始化日志**: 配置 Zap logger
3. **连接数据库**: PostgreSQL 连接 + AutoMigrate
4. **连接 Redis**: 用于缓存和限流
5. **初始化管理员**: 首次启动创建默认管理员
6. **创建路由**: 配置中间件和路由
7. **启动服务器**: 监听端口，处理请求

## 📊 数据流转

```
客户端请求
    ↓
Nginx (反向代理)
    ↓
Middleware 链:
  - CORS
  - Recovery
  - LoginLimit (登录接口)
  - JWTAuth (认证接口)
  - RBAC (权限验证)
  - RateLimiter (限流)
  - AuditLogger (审计)
    ↓
Handler (业务处理)
    ↓
Service (业务逻辑)
    ↓
Database/Redis
    ↓
响应返回
```

## 🔧 扩展指南

### 添加新模块

1. 在 `internal/handler/` 创建新的 handler 文件
2. 在 `internal/service/` 创建对应的 service
3. 在 `internal/server/http.go` 注册路由
4. 在 `internal/data/model/model.go` 添加数据模型（如需要）

### 添加新权限

1. 在数据库 permissions 表插入权限记录
2. 在路由中使用 RBAC 中间件自动验证
3. 或使用 `middleware.RequirePermission` 进行细粒度控制

## 📝 最佳实践

1. **生产环境**: 使用 `config.production.example.yaml` 作为模板
2. **密码管理**: 定期更换管理员密码
3. **日志审计**: 定期审查审计日志
4. **权限分配**: 遵循最小权限原则
5. **备份策略**: 定期备份数据库
6. **监控告警**: 配置系统监控和告警

---

更多详细信息请参考:
- [README.md](README.md) - 完整项目文档
- [QUICKSTART.md](QUICKSTART.md) - 快速启动指南
- [SECURITY_CHECKLIST.md](SECURITY_CHECKLIST.md) - 安全检查清单
