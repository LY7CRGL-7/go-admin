# 安全管理端后端系统

一个基于 Go 和 Gin 框架构建的企业级安全管理端后端系统，具备完善的安全机制、权限控制和审计功能。

## 📋 目录

- [核心特性](#核心特性)
- [安全特性](#安全特性)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [快速开始](#快速开始)
- [API 文档](#api-文档)
- [安全最佳实践](#安全最佳实践)
- [部署](#部署)

## 🎯 核心特性

- ✅ **JWT 认证** - 基于 Token 的无状态认证
- ✅ **RBAC 权限控制** - 基于角色的访问控制
- ✅ **密码安全** - bcrypt 加密 + 强密码策略
- ✅ **防暴力破解** - 登录失败限制和账号锁定
- ✅ **限流保护** - 全局/用户/IP 多维度限流
- ✅ **审计日志** - 完整的操作日志记录
- ✅ **IP 白名单** - 可选的 IP 访问控制
- ✅ **CORS 支持** - 安全的跨域配置
- ✅ **数据验证** - 严格的输入参数验证

## 🔒 安全特性

### 1. 认证安全
- JWT Token 认证，支持过期时间控制
- Bearer Token 格式验证
- Token 签发者验证
- 密码 bcrypt 加密存储

### 2. 授权安全
- 基于角色的访问控制（RBAC）
- 超级管理员特殊权限
- 接口级权限验证
- 细粒度权限控制

### 3. 密码安全
- 最小长度要求（默认 8 位）
- 必须包含大写字母
- 必须包含小写字母
- 必须包含数字
- 必须包含特殊字符
- 密码强度可配置

### 4. 防攻击
- 登录失败次数限制
- 账号锁定机制
- IP 锁定机制
- 滑动窗口限流
- SQL 注入防护（GORM ORM）
- XSS 防护

### 5. 审计追踪
- 所有操作记录审计日志
- 记录请求和响应
- 记录操作者信息
- 记录 IP 和 User-Agent
- 记录操作耗时

## 🛠 技术栈

### 核心框架
- **Go 1.23+** - 编程语言
- **Gin** - HTTP Web 框架
- **Wire** - 依赖注入框架
- **Protocol Buffers** - 接口定义语言

### 数据存储
- **PostgreSQL** - 关系型数据库
- **Redis** - 缓存和限流
- **MinIO** - 对象存储

### 消息队列
- **Kafka** - 消息队列（审计日志等异步处理）

### 监控与日志
- **Prometheus** - 指标采集
- **Zap** - 结构化日志
- **Lumberjack** - 日志轮转

### 安全
- **JWT** - Token 认证
- **bcrypt** - 密码加密

## 📂 项目结构

```
admin/
├── cmd/admin/
│   ├── conf/
│   │   └── config.yaml          # 配置文件
│   └── main.go                  # 主程序入口
├── internal/
│   ├── conf/                    # 配置结构体
│   │   └── config.go
│   ├── data/                    # 数据层
│   │   ├── model/
│   │   │   └── model.go        # 数据模型
│   │   ├── database.go          # 数据库连接
│   │   └── redis.go            # Redis 连接
│   ├── dto/                     # 数据传输对象
│   │   └── dto.go
│   ├── handler/                 # HTTP 处理器
│   │   ├── auth.go             # 认证相关
│   │   ├── admin.go            # 管理员管理
│   │   └── common.go           # 通用功能
│   ├── middleware/              # 中间件
│   │   ├── auth.go             # JWT 认证
│   │   ├── rbac.go             # RBAC 权限
│   │   ├── security.go         # 安全中间件
│   │   └── audit.go            # 审计日志
│   ├── pkg/
│   │   └── logger/             # 日志工具
│   │       └── logger.go
│   ├── service/                 # 业务逻辑层
│   │   └── auth.go
│   └── server/                  # 服务器
│       └── http.go
├── deploy/
│   └── k8s.yaml                # Kubernetes 部署配置
├── Dockerfile                   # Docker 镜像构建
├── Makefile                     # 构建脚本
├── go.mod
└── README.md
```

## 🚀 快速开始

### 环境要求

- Go 1.23+
- PostgreSQL 13+
- Redis 6.0+
- Kafka 2.8+ （可选）
- MinIO （可选）
- protoc 编译器 （可选）

### 1. 克隆项目

```bash
cd D:\golang\code
# 或从 IM 项目复制
cp -r D:\golang\company\im\app\admin D:\golang\code\admin
cd admin
```

### 2. 安装依赖

```bash
# 安装 Go 依赖
make deps

# 安装 protoc 编译器（如果需要生成 proto）
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 安装 Wire 依赖注入工具
go install github.com/google/wire/cmd/wire@latest

# 生成 proto 代码（可选）
make proto

# 生成 wire 代码（可选）
make wire
```

### 3. 配置数据库

创建 PostgreSQL 数据库：

```sql
CREATE DATABASE admin_db;
```

### 4. 修改配置

编辑 `cmd/admin/conf/config.yaml`：

```yaml
database:
  dsn: host=localhost port=5432 user=admin password=admin123 dbname=admin_db sslmode=disable timezone=Asia/Shanghai

redis:
  addr: localhost:6379
  password: "your-redis-password"

jwt:
  secret: your-256-bit-secret-key-change-in-production
```

### 5. 运行

```bash
# 开发模式
make run

# 或
go run cmd/admin/main.go -config cmd/admin/conf/config.yaml
```

### 6. 构建

```bash
make build
```

## 📡 API 文档

### 基础信息

- 基础路径: `/api/v1`
- 认证方式: Bearer Token
- 数据格式: JSON

### 认证接口

#### 1. 管理员登录

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "Admin@123456"
}
```

响应：

```json
{
  "code": 0,
  "msg": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "admin": {
      "id": 1,
      "username": "admin",
      "nickname": "系统管理员",
      "email": "",
      "phone": "",
      "status": 1
    }
  }
}
```

#### 2. 获取当前管理员信息

```http
GET /api/v1/auth/profile
Authorization: Bearer {token}
```

#### 3. 修改密码

```http
POST /api/v1/auth/change-password
Authorization: Bearer {token}
Content-Type: application/json

{
  "old_password": "OldPass@123",
  "new_password": "NewPass@456"
}
```

#### 4. 登出

```http
POST /api/v1/auth/logout
Authorization: Bearer {token}
```

### 管理员管理

#### 1. 获取管理员列表

```http
GET /api/v1/admins?page=1&page_size=20
Authorization: Bearer {token}
```

#### 2. 创建管理员

```http
POST /api/v1/admins
Authorization: Bearer {token}
Content-Type: application/json

{
  "username": "newadmin",
  "password": "Secure@123",
  "nickname": "新管理员",
  "email": "admin@example.com",
  "role_ids": [1, 2]
}
```

#### 3. 更新管理员

```http
PUT /api/v1/admins/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "nickname": "更新昵称",
  "status": 1
}
```

#### 4. 删除管理员

```http
DELETE /api/v1/admins/1
Authorization: Bearer {token}
```

### 角色管理

#### 1. 获取角色列表

```http
GET /api/v1/roles
Authorization: Bearer {token}
```

#### 2. 创建角色

```http
POST /api/v1/roles
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "内容管理员",
  "code": "content_admin",
  "description": "负责内容管理"
}
```

#### 3. 分配权限

```http
POST /api/v1/roles/1/permissions
Authorization: Bearer {token}
Content-Type: application/json

{
  "permission_ids": [1, 2, 3, 4]
}
```

### 权限管理

#### 获取权限列表

```http
GET /api/v1/permissions
Authorization: Bearer {token}
```

### 审计日志

#### 获取审计日志

```http
GET /api/v1/audit-logs?page=1&page_size=20
Authorization: Bearer {token}
```

## 🔐 安全最佳实践

### 1. 生产环境配置

```yaml
server:
  mode: release  # 必须设置为 release

jwt:
  secret: # 使用强随机密钥，至少 32 位
  expire: 2h  # 缩短 Token 过期时间

security:
  password:
    min_length: 12  # 增加最小密码长度
  
  login:
    max_attempts: 3  # 减少允许尝试次数
    lockout_duration: 1h  # 增加锁定时间
  
  ip_whitelist:  # 启用 IP 白名单
    - 10.0.0.1
    - 10.0.0.2
```

### 2. 数据库安全

- 使用强密码
- 限制数据库访问 IP
- 启用 SSL 连接
- 定期备份数据

### 3. Redis 安全

- 设置访问密码
- 限制访问 IP
- 禁用危险命令（FLUSHALL, CONFIG 等）

### 4. 网络安全

- 使用 HTTPS/TLS
- 配置防火墙规则
- 使用反向代理（Nginx）
- 启用 IP 白名单

### 5. 日志安全

- 定期审查审计日志
- 保护日志文件访问权限
- 日志文件加密存储
- 设置合理的日志保留时间

### 6. 密码管理

- 定期更换管理员密码
- 不使用默认密码
- 启用双因素认证（可选）
- 密码历史检查

## 🚢 部署

### Docker 部署

```bash
# 构建镜像
make docker

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v /path/to/config.yaml:/root/config.yaml \
  -v /path/to/logs:/opt/apps/admin/logs \
  --name admin-service \
  admin:latest
```

### Kubernetes 部署

```bash
kubectl apply -f deploy/k8s.yaml
```

### Nginx 反向代理配置

```nginx
server {
    listen 443 ssl;
    server_name admin.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 📝 初始化管理员

系统首次启动时会自动创建初始管理员账号：

- 用户名: `admin`
- 密码: `Admin@123456`（请在配置文件中修改）

**⚠️ 重要：首次登录后请立即修改密码！**

## 🔧 开发

```bash
# 格式化代码
make fmt

# 运行测试
make test

# 清理构建
make clean
```

## 📄 许可证

本项目为公司内部使用，未经授权不得使用。

## 📧 联系方式

如有问题，请联系开发团队。

---

## 📊 监控与可观测性

### Prometheus 指标

系统暴露以下 Prometheus 指标：

- `admin_http_requests_total` - HTTP 请求总数
- `admin_http_request_duration_seconds` - HTTP 请求持续时间
- `admin_active_connections` - 活跃连接数

访问指标：`http://localhost:8080/metrics`

### Grafana 仪表板

1. 访问 Grafana：`http://localhost:3000`（用户名/密码：admin/admin）
2. 添加 Prometheus 数据源：`http://prometheus:9090`
3. 导入仪表板或创建自定义图表

---

## 🔧 开发工具

### 代码生成

```bash
# 生成 Proto 代码
make proto

# 生成 Wire 依赖注入代码
make wire

# 格式化代码
make fmt

# 代码检查
make lint
```

### 快速参考

- **Proto 定义**：`proto/admin/v1/admin.proto`
- **Wire 配置**：`cmd/admin/wire.go`
- **Docker Compose**：`docker-compose.yaml`
- **Prometheus 配置**：`deploy/prometheus.yml`
