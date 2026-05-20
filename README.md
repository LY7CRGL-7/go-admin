# 🚀 企业级 Go 微服务模板

一个开箱即用的 Go 微服务项目模板，基于纯 gRPC 架构，包含完整的基础设施、安全机制和 CI/CD 配置。

> 📌 **这是一个可复用的项目模板**，你可以基于此快速构建自己的微服务项目。

---

## 📋 目录

- [核心特性](#核心特性)
- [技术栈](#技术栈)
- [快速开始（5分钟）](#快速开始5分钟)
- [使用模板构建项目](#使用模板构建项目)
  - [1. Fork/Clone 模板](#1-forkclone-模板)
  - [2. 重命名项目](#2-重命名项目)
  - [3. 配置项目](#3-配置项目)
  - [4. 启动和运行](#4-启动和运行)
- [项目结构](#项目结构)
- [配置文件说明](#配置文件说明)
- [开发指南](#开发指南)
  - [添加新的 gRPC 服务](#添加新的-grpc-服务)
  - [常用命令](#常用命令)
  - [开发工具安装](#开发工具安装)
- [Docker 使用](#docker-使用)
- [CI/CD 配置](#cicd-配置)
- [监控与可观测性](#监控与可观测性)
- [安全最佳实践](#安全最佳实践)
- [常见问题](#常见问题)
- [扩展阅读](#扩展阅读)

---

## ✨ 核心特性

### 安全与权限
- ✅ **JWT 认证** - 基于 Token 的无状态认证
- ✅ **RBAC 权限控制** - 基于角色的访问控制
- ✅ **密码安全** - bcrypt 加密 + 强密码策略
- ✅ **防暴力破解** - 登录失败限制和账号锁定
- ✅ **限流保护** - 全局/用户/IP 多维度限流
- ✅ **审计日志** - 完整的操作日志记录
- ✅ **IP 白名单** - 可选的 IP 访问控制

### 微服务架构
- ✅ **Wire 依赖注入** - Google 官方依赖注入框架
- ✅ **Protocol Buffers** - 语言无关的接口定义
- ✅ **gRPC 服务** - 高性能 RPC 框架
- ✅ **Kafka 消息队列** - 异步处理能力
- ✅ **MinIO 对象存储** - 文件存储支持
- ✅ **Prometheus 监控** - 完整的可观测性
- ✅ **Docker Compose** - 一键启动基础设施
- ✅ **GitHub Actions** - CI/CD 自动化

---

## 🛠 技术栈

| 分类 | 技术 | 版本 |
|------|------|------|
| **语言** | Go | 1.25.0+ |
| **RPC** | gRPC | v1.81.1 |
| **依赖注入** | Google Wire | v0.6.0 |
| **接口定义** | Protocol Buffers | v3 |
| **数据库** | PostgreSQL | 13+ |
| **缓存** | Redis | 6.0+ |
| **消息队列** | Kafka | 2.8+ |
| **对象存储** | MinIO | Latest |
| **监控** | Prometheus + Grafana | Latest |
| **日志** | Zap + Lumberjack | Latest |

---

## 🚀 快速开始（5分钟）

### 环境要求

- **Go 1.25.0+**（支持 toolchain 自动升级）
- Docker & Docker Compose
- protoc 编译器（**仅修改 .proto 文件时需要**，初次使用无需安装）

> 💡 **Go Toolchain 自动升级**：  
> 项目使用 `go 1.25.0` 和 `toolchain go1.25.0` 声明，CI/CD 配置 `GOTOOLCHAIN=auto`，支持自动下载更高版本工具链。

> 📦 **proto 生成文件已内置**：  
> 仓库中已包含 `*.pb.go` 文件，clone 后**无需运行** `make proto` 即可直接编译。  
> 只有修改了 `.proto` 文件时，才需要重新执行 `make proto`。

### 一键启动

```bash
# 1. Clone 项目
git clone https://github.com/LY7CRGL-7/go-admin.git
cd go-admin

# 2. 启动基础设施
docker-compose up -d

# 3. 下载依赖
make deps

# 4. 配置项目
cp cmd/admin/conf/config.yaml.example cmd/admin/conf/config.yaml

# 5. 运行
make run
```

就这么简单！🎉

---

## 📝 使用模板构建项目

### 1. Fork/Clone 模板

**推荐方式**：使用 GitHub 的 "Use this template" 按钮  
访问：https://github.com/LY7CRGL-7/go-admin

**或者手动 Clone**：
```bash
git clone https://github.com/YOUR_USERNAME/your-project.git
cd your-project
```

### 2. 重命名项目

#### 2.1 修改模块名称

编辑 `go.mod`：

```go
// 修改前
module admin

// 修改后
module github.com/YOUR_USERNAME/your-project
```

#### 2.2 批量替换 import 路径

**Linux/macOS:**
```bash
find . -name "*.go" -type f -exec sed -i 's|"admin/|"github.com/YOUR_USERNAME/your-project/|g' {} +
```

**Windows PowerShell:**
```powershell
Get-ChildItem -Recurse -Filter "*.go" | ForEach-Object {
    (Get-Content $_.FullName -Raw) -replace '"admin/', '"github.com/YOUR_USERNAME/your-project/' | Set-Content $_.FullName
}
```

#### 2.3 更新 Makefile

```makefile
# 修改这些变量
BINARY_NAME ?= your-project-name
MAIN_PATH ?= cmd/your-project/main.go
```

### 3. 配置项目

#### 3.1 必须修改的配置文件

| 文件 | 修改内容 | 说明 |
|------|---------|------|
| `go.mod` | `module` 名称 | 改为您自己的模块路径 |
| `Makefile` | `BINARY_NAME` | 改为您项目的二进制名称 |
| `docker-compose.yaml` | 所有 `YOUR_PROJECT` | 替换为您的项目名 |
| `docker-compose.yaml` | 所有密码 | ⚠️ 使用强密码 |
| `config.yaml` | 数据库、Redis等 | 根据实际环境配置 |

#### 3.2 复制并编辑配置

```bash
cp cmd/admin/conf/config.yaml.example cmd/admin/conf/config.yaml
```

编辑 `config.yaml`：

```yaml
database:
  host: localhost
  port: 5432
  user: your_user          # ⚠️ 修改
  password: your_password  # ⚠️ 修改为强密码
  name: your_database      # ⚠️ 修改

redis:
  addr: localhost:6379
  password: your_redis_password  # ⚠️ 修改

jwt:
  secret: your-jwt-secret-key-change-this  # ⚠️ 修改为强密钥
```

### 4. 启动和运行

```bash
# 启动基础设施（PostgreSQL, Redis等）
docker-compose up -d

# 下载依赖
make deps

# 运行项目（proto 文件已内置，无需先生成）
make run
```

验证服务：
```bash
# 使用 grpcurl 测试
grpcurl -plaintext localhost:9090 list

# 或使用 grpcui（可视化）
grpcui -plaintext localhost:9090
```

---

## 📂 项目结构

```
your-project/
├── cmd/your-project/
│   ├── conf/
│   │   ├── config.yaml          # 配置文件（不提交到 Git）
│   │   └── config.yaml.example  # 配置模板（提交到 Git）
│   ├── main.go                  # 主程序入口
│   └── wire.go                  # Wire 依赖注入配置
├── internal/
│   ├── conf/                    # 配置结构体
│   ├── data/                    # 数据层（数据库、Redis）
│   │   ├── model/              # 数据模型
│   │   ├── database.go
│   │   └── redis.go
│   ├── dto/                     # 数据传输对象
│   ├── grpc/                    # gRPC 服务实现
│   ├── handler/                 # HTTP 处理器（如有）
│   ├── middleware/              # 中间件
│   │   ├── auth.go             # JWT 认证
│   │   ├── rbac.go             # RBAC 权限
│   │   └── security.go         # 安全中间件
│   ├── service/                 # 业务逻辑层
│   └── server/                  # 服务器
├── proto/your-service/v1/
│   └── service.proto            # Protocol Buffers 定义
├── deploy/
│   ├── k8s.yaml                # Kubernetes 部署配置
│   └── prometheus.yml          # Prometheus 配置
├── .github/workflows/
│   └── ci.yml                  # GitHub Actions CI/CD
├── docker-compose.yaml          # Docker Compose 配置
├── Dockerfile                   # Docker 镜像构建
├── Makefile                     # 构建脚本
├── go.mod
└── README.md
```

---

## 📋 配置文件说明

### Makefile

```makefile
# 项目配置（使用者需要修改这些变量）
BINARY_NAME ?= your-project-name  # 修改为您的项目名称
MAIN_PATH ?= cmd/admin/main.go    # 主程序入口
CONFIG_PATH ?= cmd/admin/conf/config.yaml
VERSION ?= 1.0.0
```

**常用命令**：
```bash
make help     # 查看所有可用命令
make build    # 构建项目
make run      # 运行项目
make clean    # 清理构建产物
make test     # 运行测试
make proto    # 生成 proto 代码
make wire     # 生成 wire 依赖注入代码
make docker   # 构建 Docker 镜像
```

### Dockerfile

使用构建参数自定义：

```bash
docker build \
  --build-arg APP_NAME=your-project \
  --build-arg APP_VERSION=1.0.0 \
  -t YOUR_USERNAME/your-project:latest .
```

### docker-compose.yaml

包含的可选服务：

| 服务 | 必需性 | 端口 | 说明 |
|------|--------|------|------|
| PostgreSQL | ⭐ 必需 | 5432 | 主数据库 |
| Redis | 推荐 | 6379 | 缓存、限流 |
| Kafka | 可选 | 9092 | 消息队列 |
| MinIO | 可选 | 9000/9001 | 对象存储 |
| Prometheus | 可选 | 9090 | 监控指标 |
| Grafana | 可选 | 3000 | 可视化仪表板 |

**只启动必需服务**：
```bash
docker-compose up -d postgres redis
```

---

## 💻 开发指南

### 添加新的 gRPC 服务

#### 1. 定义 Proto 文件

```protobuf
// proto/your-service/v1/service.proto
syntax = "proto3";

package yourservice.v1;

option go_package = "github.com/YOUR_USERNAME/your-project/proto/your-service/v1;yourservicev1";

service YourService {
  rpc CreateItem(CreateItemRequest) returns (ItemResponse);
}

message CreateItemRequest {
  string name = 1;
}

message ItemResponse {
  int64 id = 1;
  string name = 2;
}
```

#### 2. 生成代码

```bash
make proto

# 或手动执行
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/your-service/v1/service.proto
```

#### 3. 实现服务

```go
// internal/grpc/your-service.go
package grpc

import (
    "context"
    pb "github.com/YOUR_USERNAME/your-project/proto/your-service/v1"
)

type YourService struct {
    pb.UnimplementedYourServiceServer
}

func (s *YourService) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.ItemResponse, error) {
    // 实现业务逻辑
    return &pb.ItemResponse{
        Id:   1,
        Name: req.Name,
    }, nil
}
```

#### 4. 注册服务

在 `internal/server/grpc.go` 中注册您的服务。

### 开发工具安装

```bash
# 1. protoc 编译器
# macOS
brew install protobuf

# Linux
wget https://github.com/protocolbuffers/protobuf/releases/download/v29.3/protoc-29.3-linux-x86_64.zip
sudo unzip protoc.zip -d /usr/local/

# 2. Go protobuf 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 3. Wire 依赖注入
go install github.com/google/wire/cmd/wire@latest

# 4. gRPC 调试工具
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
```

---

## 🐳 Docker 使用

### 构建镜像

```bash
# 方式1：使用 Makefile
make docker

# 方式2：直接使用 docker 命令
docker build -t YOUR_USERNAME/your-project:latest .

# 方式3：使用构建参数
docker build \
  --build-arg APP_NAME=your-project \
  --build-arg APP_VERSION=1.0.0 \
  -t YOUR_USERNAME/your-project:1.0.0 .
```

### 运行容器

```bash
# 确保基础设施已启动
docker-compose up -d postgres redis

# 运行您的服务
docker run -d \
  --name your-project \
  --network your-project-network \
  -p 9090:9090 \
  -v $(pwd)/config.yaml:/root/config.yaml \
  YOUR_USERNAME/your-project:latest
```

### 多阶段构建说明

Dockerfile 使用多阶段构建：
- **Builder 阶段**：编译 Go 代码
- **Runner 阶段**：只包含二进制文件和运行时依赖

最终镜像约 20MB，保持最小化。

---

## 🔄 CI/CD 配置

### GitHub Actions 质量检查

模板包含的 CI 用于**验证模板质量**，不是必须的构建流程：

- ✅ 编译验证（直接编译，无需安装 protoc）
- ✅ Docker 构建验证
- 📝 代码测试（可选，失败不阻断）
- 📝 格式检查（可选，失败不阻断）

### 自定义 CI

编辑 `.github/workflows/ci.yml` 添加您的检查：

```yaml
jobs:
  quality-check:
    steps:
      # ... 现有步骤 ...
      
      # 添加您的自定义检查
      - name: Custom Check
        run: |
          echo "Running custom checks..."
```

### 添加 Docker Hub 推送

1. 在 GitHub 仓库 → Settings → Secrets 添加：
   - `DOCKERHUB_USERNAME`
   - `DOCKERHUB_TOKEN`

2. 修改 CI 中的 Docker job：

```yaml
- name: Login to Docker Hub
  uses: docker/login-action@v3
  with:
    username: ${{ secrets.DOCKERHUB_USERNAME }}
    password: ${{ secrets.DOCKERHUB_TOKEN }}

- name: Build and push
  uses: docker/build-push-action@v5
  with:
    context: .
    push: true
    tags: |
      YOUR_USERNAME/your-project:latest
      YOUR_USERNAME/your-project:${{ github.sha }}
```

---

## 📊 监控与可观测性

### Prometheus 指标

模板已集成 Prometheus，暴露以下指标：

- HTTP/gRPC 请求计数
- 请求延迟
- 错误率
- Goroutine 数量

### 访问 Grafana

1. 访问 `http://localhost:3000`
2. 使用配置的用户名密码登录
3. 添加 Prometheus 数据源：`http://prometheus:9090`
4. 导入 Go 应用监控仪表板

---

## 🔒 安全最佳实践

### 1. 密码管理

- ⚠️ **永远不要**将密码提交到 Git
- 使用 `.gitignore` 排除 `config.yaml`
- 只提交 `config.yaml.example` 模板
- 使用环境变量或密钥管理服务

### 2. JWT 密钥

```yaml
# config.yaml
jwt:
  secret: ${JWT_SECRET}  # 从环境变量读取
```

```bash
# 生成强密钥
openssl rand -base64 32
```

### 3. 生产环境配置

```bash
# 使用环境变量覆盖
export DB_PASSWORD=production_password
export JWT_SECRET=production_jwt_secret

# 或使用 .env 文件（不提交到 Git）
cat .env
DB_PASSWORD=production_password
JWT_SECRET=production_jwt_secret
```

---

## ❓ 常见问题

### Q: 如何只启动必需的服务？

```bash
# 只启动 PostgreSQL 和 Redis
docker-compose up -d postgres redis

# 或注释掉 docker-compose.yaml 中不需要的服务
```

### Q: 需要安装 protoc 吗？

**初次使用不需要**。仓库已内置生成好的 `*.pb.go` 文件，clone 后直接 `make run` 即可。

只有**修改了 `.proto` 文件**后才需要重新生成：

```bash
# 安装 protoc（如需修改 proto 文件）
# macOS
brew install protobuf

# Linux
wget https://github.com/protocolbuffers/protobuf/releases/download/v29.3/protoc-29.3-linux-x86_64.zip
sudo unzip protoc.zip -d /usr/local/

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 修改 .proto 后重新生成，并提交 .pb.go
make proto
git add proto/
git commit -m "feat: 更新 proto 定义"
```

### Q: 数据库连接失败？

```bash
# 检查 PostgreSQL 是否运行
docker-compose ps postgres

# 查看日志
docker-compose logs postgres

# 测试连接
psql -h localhost -U your_user -d your_database
```

### Q: 如何升级到新的 Go 版本？

```bash
# 1. 更新 go.mod
go mod edit -go=1.26

# 2. 更新 Dockerfile
# FROM golang:1.26-alpine AS builder

# 3. 更新 .github/workflows/ci.yml
# GO_VERSION: '1.26'
```

---

## 📚 扩展阅读

- [Go 官方文档](https://golang.org/doc/)
- [gRPC 官方文档](https://grpc.io/docs/)
- [Protocol Buffers 文档](https://developers.google.com/protocol-buffers)
- [Google Wire 文档](https://github.com/google/wire)
- [Docker 最佳实践](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)

---

## 🤝 贡献

如果您在使用模板过程中发现问题或有改进建议，欢迎提交 Issue 或 Pull Request。

---

## 📄 许可证

MIT License

---

**作者**: LY7CRGL-7  
**版本**: v1.0.0  
**仓库**: https://github.com/LY7CRGL-7/go-admin

**祝您开发愉快！** 🎉
