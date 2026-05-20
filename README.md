# Kratos Admin Template

基于 [Kratos](https://go-kratos.dev/) 的企业级管理端模板，采用严格四层架构（Server / Service / Biz / Data），纯 Wire 依赖注入，内置 RBAC 权限与多租户支持，支持 MySQL / PostgreSQL 切换。

---

## 目录

- [项目结构](#项目结构)
- [四层架构详解](#四层架构详解)
- [快速开始](#快速开始)
- [配置详解](#配置详解)
- [如何写业务（完整示例）](#如何写业务完整示例)
- [内置 gRPC 服务一览](#内置-grpc-服务一览)
- [Wire 依赖注入](#wire-依赖注入)
- [常用命令](#常用命令)
- [Docker 部署](#docker-部署)
- [作为模板使用](#作为模板使用)

---

## 项目结构

```
.
├── api/admin/v1/                # ① API 层 - Protobuf 定义 + 生成代码
│   ├── admin.proto              #    所有 gRPC 服务和消息的定义
│   ├── admin.pb.go              #    protoc 生成的消息序列化代码
│   └── admin_grpc.pb.go         #    protoc 生成的 gRPC 服务接口
│
├── cmd/admin/                   # ② 应用入口
│   ├── main.go                  #    程序启动、加载配置、Wire 注入
│   ├── wire.go                  #    Wire 注入声明（wireinject 标签）
│   ├── wire_gen.go              #    Wire 自动生成的注入代码
│   └── conf/                    #    配置文件目录
│       ├── config.yaml          #    开发环境配置（可直接运行）
│       ├── config.yaml.example  #    配置模板（占位符）
│       └── config.production.example.yaml  # 生产环境配置示例
│
├── internal/                    # ③ 核心业务代码（不可外部导入）
│   ├── conf/                    #    配置结构体定义
│   │   └── conf.go              #    Bootstrap / Server / Data / Auth / Tenant
│   │
│   ├── server/                  #    Server 层 - 创建服务器实例
│   │   ├── server.go            #    Wire ProviderSet
│   │   ├── grpc.go              #    gRPC 服务器（注册所有服务+中间件）
│   │   └── http.go              #    HTTP 服务器（健康检查 /health /ready）
│   │
│   ├── service/                 #    Service 层 - 实现 proto 定义的 gRPC 接口
│   │   ├── service.go           #    Wire ProviderSet
│   │   ├── auth.go              #    AuthService 实现（登录/登出/刷新Token）
│   │   ├── user.go              #    UserService 实现（用户CRUD+分配角色）
│   │   ├── role.go              #    RoleService 实现（角色CRUD+分配权限）
│   │   ├── permission.go        #    PermissionService 实现（权限列表/树形）
│   │   ├── tenant.go            #    TenantService 实现（租户CRUD）
│   │   ├── audit.go             #    AuditLogService 实现（审计日志查询）
│   │   └── convert.go           #    Model ↔ Proto 类型转换 + Context 工具
│   │
│   ├── biz/                     #    Biz 层 - 核心业务逻辑
│   │   ├── biz.go               #    Wire ProviderSet + Repo 接口定义
│   │   ├── user.go              #    UserUsecase（登录/JWT/RBAC/CRUD）
│   │   ├── role.go              #    RoleUsecase（角色管理）
│   │   └── tenant.go            #    TenantUsecase（租户管理）
│   │
│   ├── data/                    #    Data 层 - 数据库操作
│   │   ├── data.go              #    Wire ProviderSet + DB/Redis 连接初始化
│   │   ├── model/model.go       #    GORM 数据模型（User/Role/Permission/Tenant/AuditLog）
│   │   ├── user.go              #    UserRepo 实现
│   │   ├── role.go              #    RoleRepo 实现
│   │   ├── permission.go        #    PermissionRepo 实现
│   │   ├── tenant.go            #    TenantRepo 实现
│   │   └── audit.go             #    AuditRepo 实现
│   │
│   └── middleware/              #    中间件
│       └── middleware.go        #    JWT 认证 + 白名单 + 多租户上下文
│
├── deploy/                      # 部署配置
│   └── k8s.yaml                 #    Kubernetes Deployment + Service
├── docs/
│   └── init.sql                 #    数据库初始化脚本（可选初始数据）
│
├── Dockerfile                   # Docker 构建（多阶段，固定产物路径）
├── docker-compose.yaml          # PostgreSQL + Redis（可选 MySQL）
├── Makefile                     # 构建/运行/生成/Docker 等命令
├── go.mod / go.sum              # Go 模块依赖
└── .github/workflows/ci.yml    # CI 质量检查
```

---

## 四层架构详解

```
请求 → Server → Service → Biz → Data → 数据库
       (路由)   (转换)    (逻辑)  (存储)
```

### 各层职责

| 层 | 目录 | 做什么 | 不做什么 |
|----|------|--------|----------|
| **Server** | `internal/server/` | 创建 gRPC/HTTP 服务器，挂载中间件，注册路由 | 不写业务逻辑 |
| **Service** | `internal/service/` | 实现 proto 定义的接口，做 Request → Model 转换 | 不写数据库操作 |
| **Biz** | `internal/biz/` | 核心业务逻辑，定义 Repo 接口（依赖倒置） | 不依赖具体数据库实现 |
| **Data** | `internal/data/` | 实现 Repo 接口，操作 DB/Redis | 不写业务判断逻辑 |

### 依赖方向（单向）

```
Server → Service → Biz ← Data
                    ↑
             Biz 定义 Repo 接口
             Data 实现 Repo 接口
```

**关键原则**：Biz 层通过接口依赖 Data 层（依赖倒置），而不是直接 import data 包。

---

## 快速开始

### 环境要求

- **Go 1.25+**
- **Docker & Docker Compose**（用于启动 PostgreSQL 和 Redis）
- `protoc` + `protoc-gen-go`（仅修改 `.proto` 文件时需要）
- `wire`（仅修改依赖注入时需要）

### 5 步启动

```bash
# 1. Clone
git clone https://github.com/LY7CRGL-7/go-admin.git
cd go-admin

# 2. 启动数据库和缓存
docker-compose up -d

# 3. 复制配置文件
cp cmd/admin/conf/config.yaml.example cmd/admin/conf/config.yaml
# 或者直接用项目自带的 config.yaml（已含开发环境默认值）

# 4. 下载依赖
make deps

# 5. 运行
make run
```

启动成功后：
- **HTTP**: http://localhost:8080（健康检查: `GET /health`）
- **gRPC**: localhost:9090

---

## 配置详解

配置文件位于 `cmd/admin/conf/config.yaml`，结构定义在 `internal/conf/conf.go`。

### 完整配置说明

```yaml
# ==================== 服务器 ====================
server:
  http:
    addr: 0.0.0.0:8080    # HTTP 监听地址（健康检查用）
    timeout: 10s           # 请求超时
  grpc:
    addr: 0.0.0.0:9090    # gRPC 监听地址（主要业务协议）
    timeout: 10s

# ==================== 数据源 ====================
data:
  database:
    driver: postgres       # 数据库驱动：postgres 或 mysql
    source: "host=localhost port=5432 user=admin password=admin123 dbname=admin_db sslmode=disable timezone=Asia/Shanghai"
    # MySQL 示例：
    # driver: mysql
    # source: "admin:admin123@tcp(localhost:3306)/admin_db?charset=utf8mb4&parseTime=True&loc=Local"
    max_idle: 10           # 最大空闲连接数
    max_open: 100          # 最大打开连接数
    max_lifetime: 3600s    # 连接最大存活时间

  redis:
    addr: localhost:6379   # Redis 地址
    password: ""           # Redis 密码（空则无密码）
    db: 0                  # Redis 数据库编号
    dial_timeout: 5s
    read_timeout: 3s
    write_timeout: 3s

# ==================== 认证 ====================
auth:
  jwt_secret: your-256-bit-secret-key    # ⚠️ JWT 签名密钥，生产环境必须修改
  token_expire: 24h                       # Access Token 有效期
  refresh_expire: 168h                    # Refresh Token 有效期（7天）
  init_admin:                             # 首次启动自动创建的管理员
    username: admin
    password: Admin@123456                # ⚠️ 生产环境必须修改
    nickname: 系统管理员

# ==================== 多租户 ====================
tenant:
  enabled: false                          # 是否启用多租户
  default_tenant_code: default            # 默认租户编码
```

### 切换数据库（PostgreSQL ↔ MySQL）

**第 1 步**：修改 `config.yaml`

```yaml
# PostgreSQL → MySQL
data:
  database:
    driver: mysql
    source: "admin:admin123@tcp(localhost:3306)/admin_db?charset=utf8mb4&parseTime=True&loc=Local"
```

**第 2 步**：修改 `docker-compose.yaml`

```yaml
# 注释掉 postgres 服务，取消注释 mysql 服务
```

**第 3 步**：重启

```bash
docker-compose down && docker-compose up -d
make run
```

数据库连接代码在 `internal/data/data.go` 的 `newDB()` 函数中，根据 `driver` 字段自动选择驱动。

---

## 如何写业务（完整示例）

以"添加一个商品管理模块"为例，逐步说明在哪写、写什么。

### 第 1 步：定义 API（proto）

编辑 `api/admin/v1/admin.proto`，添加新的 service 和 message：

```protobuf
// ==================== 商品管理 ====================

service ProductService {
  rpc CreateProduct(CreateProductRequest) returns (ProductInfo);
  rpc GetProduct(GetProductRequest) returns (ProductInfo);
  rpc UpdateProduct(UpdateProductRequest) returns (ProductInfo);
  rpc DeleteProduct(DeleteProductRequest) returns (DeleteReply);
  rpc ListProducts(ListProductsRequest) returns (ListProductsReply);
}

message ProductInfo {
  int64 id = 1;
  string name = 2;
  string description = 3;
  int64 price = 4;        // 单位：分
  int32 stock = 5;
  int32 status = 6;       // 1上架 0下架
  int64 tenant_id = 7;
  google.protobuf.Timestamp created_at = 8;
}

message CreateProductRequest {
  string name = 1;
  string description = 2;
  int64 price = 3;
  int32 stock = 4;
}

message GetProductRequest { int64 id = 1; }

message UpdateProductRequest {
  int64 id = 1;
  string name = 2;
  string description = 3;
  int64 price = 4;
  int32 stock = 5;
  int32 status = 6;
}

message DeleteProductRequest { int64 id = 1; }

message ListProductsRequest {
  int32 page = 1;
  int32 page_size = 2;
  string keyword = 3;
}

message ListProductsReply {
  repeated ProductInfo items = 1;
  int64 total = 2;
}
```

然后运行：

```bash
make proto
```

### 第 2 步：添加数据模型

编辑 `internal/data/model/model.go`，添加 Product 结构体：

```go
// Product 商品
type Product struct {
    TenantModel
    Name        string `gorm:"size:200;not null" json:"name"`
    Description string `gorm:"size:1000" json:"description"`
    Price       int64  `gorm:"not null;default:0;comment:单位分" json:"price"`
    Stock       int32  `gorm:"not null;default:0" json:"stock"`
    Status      int8   `gorm:"default:1;comment:1上架0下架" json:"status"`
}

func (Product) TableName() string { return "products" }
```

在 `AutoMigrate` 中添加 `&Product{}`。

### 第 3 步：定义 Repo 接口（Biz 层）

编辑 `internal/biz/biz.go`，添加接口：

```go
// ProductRepo 商品仓储接口
type ProductRepo interface {
    Create(ctx context.Context, product *model.Product) error
    GetByID(ctx context.Context, id uint) (*model.Product, error)
    Update(ctx context.Context, product *model.Product) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, tenantID uint, page, pageSize int, keyword string) ([]*model.Product, int64, error)
}
```

### 第 4 步：编写业务逻辑（Biz 层）

新建 `internal/biz/product.go`：

```go
package biz

import (
    "context"
    "admin/internal/data/model"
    "github.com/go-kratos/kratos/v2/log"
)

type ProductUsecase struct {
    repo ProductRepo
    log  *log.Helper
}

func NewProductUsecase(repo ProductRepo, logger log.Logger) *ProductUsecase {
    return &ProductUsecase{
        repo: repo,
        log:  log.NewHelper(log.With(logger, "module", "biz/product")),
    }
}

func (uc *ProductUsecase) Create(ctx context.Context, p *model.Product) error {
    // 这里写业务校验逻辑，例如：名称不能重复、价格必须大于0
    if p.Price < 0 {
        return errors.New("价格不能为负数")
    }
    return uc.repo.Create(ctx, p)
}

func (uc *ProductUsecase) Get(ctx context.Context, id uint) (*model.Product, error) {
    return uc.repo.GetByID(ctx, id)
}

// ... 其他方法
```

在 `biz.go` 的 ProviderSet 中添加 `NewProductUsecase`。

### 第 5 步：实现数据访问（Data 层）

新建 `internal/data/product.go`：

```go
package data

import (
    "context"
    "admin/internal/biz"
    "admin/internal/data/model"
    "github.com/go-kratos/kratos/v2/log"
)

type productRepo struct {
    data *Data
    log  *log.Helper
}

func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
    return &productRepo{
        data: data,
        log:  log.NewHelper(log.With(logger, "module", "data/product")),
    }
}

func (r *productRepo) Create(ctx context.Context, p *model.Product) error {
    return r.data.DB.WithContext(ctx).Create(p).Error
}

func (r *productRepo) GetByID(ctx context.Context, id uint) (*model.Product, error) {
    var p model.Product
    err := r.data.DB.WithContext(ctx).First(&p, id).Error
    return &p, err
}

// ... 其他方法
```

在 `data.go` 的 ProviderSet 中添加 `NewProductRepo`。

### 第 6 步：实现 gRPC 接口（Service 层）

新建 `internal/service/product.go`：

```go
package service

import (
    "context"
    v1 "admin/api/admin/v1"
    "admin/internal/biz"
    "admin/internal/data/model"
    "github.com/go-kratos/kratos/v2/log"
)

type ProductService struct {
    v1.UnimplementedProductServiceServer
    uc  *biz.ProductUsecase
    log *log.Helper
}

func NewProductService(uc *biz.ProductUsecase, logger log.Logger) *ProductService {
    return &ProductService{
        uc:  uc,
        log: log.NewHelper(log.With(logger, "module", "service/product")),
    }
}

func (s *ProductService) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.ProductInfo, error) {
    p := &model.Product{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
        Status:      1,
    }
    if claims := GetClaimsFromContext(ctx); claims != nil {
        p.TenantID = claims.TenantID
    }
    if err := s.uc.Create(ctx, p); err != nil {
        return nil, err
    }
    return productToProto(p), nil
}

// ... 其他方法
```

在 `service.go` 的 ProviderSet 中添加 `NewProductService`。

在 `convert.go` 中添加转换函数：

```go
func productToProto(p *model.Product) *v1.ProductInfo {
    return &v1.ProductInfo{
        Id: int64(p.ID), Name: p.Name, Description: p.Description,
        Price: p.Price, Stock: p.Stock, Status: int32(p.Status),
        TenantId: int64(p.TenantID), CreatedAt: timestamppb.New(p.CreatedAt),
    }
}
```

### 第 7 步：注册到 gRPC Server

编辑 `internal/server/grpc.go`，在 `NewGRPCServer` 的参数中加入 `productSvc`：

```go
func NewGRPCServer(
    // ... 现有参数
    productSvc *service.ProductService,  // 新增
) *grpc.Server {
    // ...
    v1.RegisterProductServiceServer(srv, productSvc)  // 新增
    return srv
}
```

### 第 8 步：更新 Wire 注入

```bash
make wire
```

Wire 会自动根据依赖关系生成 `wire_gen.go`。

### 第 9 步：运行

```bash
make run
```

### 总结：添加新业务的固定流程

| 步骤 | 文件 | 做什么 |
|------|------|--------|
| 1 | `api/admin/v1/admin.proto` | 定义 gRPC service + message |
| 2 | 执行 `make proto` | 生成 Go 代码 |
| 3 | `internal/data/model/model.go` | 添加 GORM 数据模型 |
| 4 | `internal/biz/biz.go` | 定义 Repo 接口 |
| 5 | `internal/biz/xxx.go` | 编写 Usecase（业务逻辑） |
| 6 | `internal/data/xxx.go` | 实现 Repo（数据库操作） |
| 7 | `internal/service/xxx.go` | 实现 gRPC 接口（协议转换） |
| 8 | `internal/service/convert.go` | 添加 Model → Proto 转换函数 |
| 9 | `internal/server/grpc.go` | 注册新服务到 gRPC Server |
| 10 | 执行 `make wire` | 重新生成依赖注入代码 |

---

## 内置 gRPC 服务一览

定义在 `api/admin/v1/admin.proto`，实现在 `internal/service/` 下。

### AuthService（认证）

| 方法 | 说明 | 需要登录 |
|------|------|----------|
| `Login` | 用户名密码登录，返回 JWT Token | 否 |
| `Logout` | 登出 | 是 |
| `GetProfile` | 获取当前用户信息 | 是 |
| `ChangePassword` | 修改密码 | 是 |
| `RefreshToken` | 刷新 Token | 否 |

### UserService（用户管理）

| 方法 | 说明 |
|------|------|
| `CreateUser` | 创建用户（自动加密密码） |
| `GetUser` | 获取用户详情（含角色信息） |
| `UpdateUser` | 更新用户信息 |
| `DeleteUser` | 删除用户（软删除） |
| `ListUsers` | 分页查询用户（支持关键字搜索） |
| `AssignRoles` | 给用户分配角色 |

### RoleService（角色管理）

| 方法 | 说明 |
|------|------|
| `CreateRole` | 创建角色 |
| `GetRole` | 获取角色详情（含权限列表） |
| `UpdateRole` | 更新角色 |
| `DeleteRole` | 删除角色 |
| `ListRoles` | 分页查询角色 |
| `AssignPermissions` | 给角色分配权限 |

### PermissionService（权限管理）

| 方法 | 说明 |
|------|------|
| `ListPermissions` | 获取权限平铺列表 |
| `GetPermissionTree` | 获取权限树形结构 |

### TenantService（多租户管理）

| 方法 | 说明 |
|------|------|
| `CreateTenant` | 创建租户 |
| `GetTenant` / `UpdateTenant` / `DeleteTenant` | 租户 CRUD |
| `ListTenants` | 分页查询租户 |

### AuditLogService（审计日志）

| 方法 | 说明 |
|------|------|
| `ListAuditLogs` | 查询审计日志（按用户/操作筛选） |

---

## Wire 依赖注入

项目使用 [Google Wire](https://github.com/google/wire) 做编译期依赖注入。

### 核心文件

- `cmd/admin/wire.go` - 注入声明（`wireApp` 函数，`wireinject` 标签）
- `cmd/admin/wire_gen.go` - Wire 自动生成的代码

### 各层 ProviderSet

```go
// internal/data/data.go
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewRoleRepo, NewPermissionRepo, NewTenantRepo, NewAuditRepo)

// internal/biz/biz.go
var ProviderSet = wire.NewSet(NewUserUsecase, NewRoleUsecase, NewTenantUsecase)

// internal/service/service.go
var ProviderSet = wire.NewSet(NewAuthService, NewUserService, NewRoleService, NewPermissionService, NewTenantService, NewAuditService)

// internal/server/server.go
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)
```

### 添加新的 Provider

1. 在对应层的 ProviderSet 中添加构造函数
2. 运行 `make wire`
3. Wire 会自动分析依赖图并生成 `wire_gen.go`

如果没有安装 Wire：

```bash
go install github.com/google/wire/cmd/wire@latest
```

---

## 常用命令

```bash
make build          # 编译二进制到 bin/admin
make run            # 运行（默认读取 cmd/admin/conf/config.yaml）
make test           # 运行测试
make fmt            # 格式化代码
make proto          # 重新生成 proto 代码（修改 .proto 后执行）
make wire           # 重新生成 Wire 代码（修改依赖注入后执行）
make generate       # proto + wire 一起生成
make docker         # 构建 Docker 镜像
make help           # 查看所有命令
```

---

## Docker 部署

### Dockerfile 说明

```dockerfile
# 构建阶段：Go 编译，输出固定路径 /app
FROM golang:1.25-alpine AS builder
RUN go build -o /app ./cmd/admin

# 运行阶段：仅包含二进制 + 配置
FROM alpine:3.20
COPY --from=builder /app .
EXPOSE 8080 9090
ENTRYPOINT ["./admin"]
```

### 构建和运行

```bash
# 构建
docker build -t admin:latest .

# CI 注入版本号
docker build --build-arg VERSION=$(git describe --tags) -t admin:v1.0.0 .

# 运行（挂载配置文件）
docker run -p 8080:8080 -p 9090:9090 \
  -v $(pwd)/cmd/admin/conf/config.yaml:/app/config.yaml \
  admin:latest
```

### Docker Compose（开发环境）

```bash
# 启动 PostgreSQL + Redis
docker-compose up -d

# 停止
docker-compose down
```

如需切换 MySQL，编辑 `docker-compose.yaml`：注释 `postgres` 服务，取消注释 `mysql` 服务。

---

## 作为模板使用

### 第 1 步：Fork / Clone

```bash
git clone https://github.com/LY7CRGL-7/go-admin.git my-project
cd my-project
```

### 第 2 步：修改模块名

编辑 `go.mod`：

```go
module github.com/your-org/my-project  // 替换 admin
```

### 第 3 步：批量替换 import 路径

```bash
# Linux/macOS
find . -name "*.go" -exec sed -i 's|"admin/|"github.com/your-org/my-project/|g' {} +

# Windows PowerShell
Get-ChildItem -Recurse -Filter "*.go" | ForEach-Object {
    (Get-Content $_.FullName -Raw) -replace '"admin/', '"github.com/your-org/my-project/' | Set-Content $_.FullName
}
```

### 第 4 步：修改配置，启动开发

```bash
# 编辑配置
vim cmd/admin/conf/config.yaml

# 启动依赖
docker-compose up -d

# 运行
make deps && make run
```

---

## License

MIT
