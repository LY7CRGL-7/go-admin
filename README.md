# Kratos Admin Template

基于 [Kratos](https://go-kratos.dev/) 的企业级管理端模板，采用严格四层架构（Server / Service / Biz / Data），内置 RBAC 权限与多租户支持，适用于 SaaS 后台系统。

## 特性

- **Kratos 框架** - gRPC + HTTP 双协议，Protobuf 定义 API
- **Wire 依赖注入** - 编译期安全，零反射
- **严格四层架构** - Server → Service → Biz → Data，职责清晰
- **RBAC 权限模型** - 用户 / 角色 / 权限，多对多关联
- **多租户支持** - 数据隔离，按需开启
- **双数据库支持** - MySQL / PostgreSQL 一键切换
- **JWT 认证** - Access Token + Refresh Token，中间件自动鉴权
- **审计日志** - 操作记录，可追溯
- **Docker 支持** - Dockerfile + Compose + K8s 部署清单

## 技术栈

| 类别 | 技术 |
|------|------|
| 框架 | Go 1.25 + Kratos v2 |
| 协议 | gRPC + HTTP (Protobuf) |
| 依赖注入 | Google Wire |
| ORM | GORM (MySQL / PostgreSQL) |
| 缓存 | Redis |
| 认证 | JWT (golang-jwt/v5) |
| 部署 | Docker / Kubernetes |

## 项目结构

```
.
├── api/admin/v1/            # Protobuf API 定义 + 生成代码
│   ├── admin.proto
│   ├── admin.pb.go
│   └── admin_grpc.pb.go
├── cmd/admin/               # 应用入口
│   ├── main.go              # 启动 + 配置加载
│   ├── wire.go              # Wire 注入声明
│   ├── wire_gen.go          # Wire 生成代码
│   └── conf/                # 配置文件
├── internal/
│   ├── conf/                # 配置结构体
│   ├── server/              # Server 层 - gRPC/HTTP 服务器
│   ├── service/             # Service 层 - gRPC 接口实现
│   ├── biz/                 # Biz 层 - 业务逻辑 + Repo 接口
│   ├── data/                # Data 层 - 数据库实现
│   │   └── model/           # GORM 数据模型
│   └── middleware/          # 中间件（JWT / 多租户）
├── deploy/                  # 部署配置（K8s / Prometheus）
├── docs/                    # 数据库初始化脚本
├── Dockerfile
├── docker-compose.yaml
└── Makefile
```

### 四层架构说明

```
请求 → Server（路由/中间件）→ Service（协议转换）→ Biz（业务逻辑）→ Data（数据访问）
                                                      ↑ 定义 Repo 接口
                                                                        ↑ 实现 Repo 接口
```

| 层 | 职责 | 依赖方向 |
|----|------|----------|
| Server | 创建 gRPC/HTTP 服务器，挂载中间件 | → Service |
| Service | 实现 proto 定义的 gRPC 接口，做协议转换 | → Biz |
| Biz | 核心业务逻辑，定义 Repo 接口（依赖倒置） | → Repo Interface |
| Data | 实现 Repo 接口，操作数据库/缓存 | 实现 Biz 的接口 |

## 快速开始

### 环境要求

- Go 1.25+
- Docker & Docker Compose
- protoc + protoc-gen-go（仅修改 `.proto` 时需要）

### 启动

```bash
# 1. Clone
git clone https://github.com/LY7CRGL-7/go-admin.git
cd go-admin

# 2. 启动 PostgreSQL + Redis
docker-compose up -d

# 3. 复制配置
cp cmd/admin/conf/config.yaml.example cmd/admin/conf/config.yaml

# 4. 下载依赖
make deps

# 5. 运行
make run
```

服务启动后：
- HTTP: `http://localhost:8080`（健康检查: `/health`）
- gRPC: `localhost:9090`

## 配置说明

配置文件：`cmd/admin/conf/config.yaml`

```yaml
server:
  http:
    addr: 0.0.0.0:8080
    timeout: 10s
  grpc:
    addr: 0.0.0.0:9090
    timeout: 10s

data:
  database:
    driver: postgres  # 或 mysql
    source: "host=localhost port=5432 user=admin password=admin123 dbname=admin_db sslmode=disable"
    max_idle: 10
    max_open: 100
    max_lifetime: 3600s
  redis:
    addr: localhost:6379

auth:
  jwt_secret: your-secret-key
  token_expire: 24h
  refresh_expire: 168h
  init_admin:
    username: admin
    password: Admin@123456

tenant:
  enabled: false
  default_tenant_code: default
```

### 切换数据库

修改 `config.yaml` 中 `data.database` 即可：

```yaml
# PostgreSQL
driver: postgres
source: "host=localhost port=5432 user=admin password=admin123 dbname=admin_db sslmode=disable"

# MySQL
driver: mysql
source: "admin:admin123@tcp(localhost:3306)/admin_db?charset=utf8mb4&parseTime=True&loc=Local"
```

同时修改 `docker-compose.yaml` 中启用对应的数据库服务。

## gRPC 服务

模板内置以下 gRPC 服务（定义在 `api/admin/v1/admin.proto`）：

| 服务 | 方法 | 说明 |
|------|------|------|
| AuthService | Login, Logout, GetProfile, ChangePassword, RefreshToken | 认证 |
| UserService | CreateUser, GetUser, UpdateUser, DeleteUser, ListUsers, AssignRoles | 用户管理 |
| RoleService | CreateRole, GetRole, UpdateRole, DeleteRole, ListRoles, AssignPermissions | 角色管理 |
| PermissionService | ListPermissions, GetPermissionTree | 权限查询 |
| TenantService | CreateTenant, GetTenant, UpdateTenant, DeleteTenant, ListTenants | 多租户 |
| AuditLogService | ListAuditLogs | 审计日志 |

## 常用命令

```bash
make build          # 编译
make run            # 运行
make test           # 测试
make proto          # 重新生成 proto 代码
make wire           # 重新生成 wire 代码
make generate       # proto + wire 一键生成
make docker         # 构建 Docker 镜像
make help           # 查看所有命令
```

## 开发指南

### 添加新的 gRPC 服务

1. 在 `api/admin/v1/admin.proto` 中定义新服务和消息
2. 运行 `make proto` 生成代码
3. 在 `internal/biz/` 中定义 Repo 接口和 Usecase
4. 在 `internal/data/` 中实现 Repo 接口
5. 在 `internal/service/` 中实现 gRPC 服务
6. 在 `internal/server/grpc.go` 中注册服务
7. 更新 Wire Provider Set，运行 `make wire`

### 开发工具安装

```bash
# protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Wire
go install github.com/google/wire/cmd/wire@latest
```

## Docker 部署

```bash
# 构建镜像
docker build -t admin:latest .

# 或通过 CI 注入版本
docker build --build-arg VERSION=$(git describe --tags) -t admin:v1.0.0 .

# 运行
docker run -p 8080:8080 -p 9090:9090 -v ./config.yaml:/app/config.yaml admin:latest
```

## 作为模板使用

1. Fork 或使用 GitHub 的 "Use this template"
2. 修改 `go.mod` 中的 module 名称
3. 批量替换 import 路径：
   ```bash
   # Linux/macOS
   find . -name "*.go" -exec sed -i 's|"admin/|"your-module/|g' {} +
   
   # Windows PowerShell
   Get-ChildItem -Recurse -Filter "*.go" | ForEach-Object {
       (Get-Content $_.FullName -Raw) -replace '"admin/', '"your-module/' | Set-Content $_.FullName
   }
   ```
4. 修改配置文件中的数据库连接信息
5. 运行 `make deps && make run`

## License

MIT
