# 🚀 Go 微服务项目模板

这是一个企业级 Go 微服务项目模板，基于 IM 项目架构标准构建。

## ✨ 特性

- ✅ **Wire 依赖注入** - Google 官方依赖注入框架
- ✅ **Protocol Buffers** - 语言无关的接口定义
- ✅ **Kafka 消息队列** - 异步处理能力
- ✅ **MinIO 对象存储** - 文件存储支持
- ✅ **Prometheus 监控** - 完整的可观测性
- ✅ **Docker Compose** - 一键启动基础设施
- ✅ **GitHub Actions** - CI/CD 自动化

## 📋 技术栈

- **Go 1.23+**
- **Gin** - HTTP 框架
- **GORM** - ORM
- **PostgreSQL** - 数据库
- **Redis** - 缓存
- **Wire** - 依赖注入
- **Kafka** - 消息队列
- **MinIO** - 对象存储
- **Prometheus** - 监控

## 🚀 快速开始

### 1. 使用模板

```bash
# 克隆模板
git clone https://github.com/LY7CRGL-7/admin.git my-project
cd my-project

# 删除 .git 并重新初始化
rm -rf .git
git init
```

### 2. 启动基础设施

```bash
docker-compose up -d
```

### 3. 安装依赖

```bash
make deps
```

### 4. 配置项目

```bash
cp cmd/admin/conf/config.yaml.example cmd/admin/conf/config.yaml
# 编辑配置文件
```

### 5. 运行

```bash
make run
```

## 📁 项目结构

```
admin/
├── proto/              # Protocol Buffers 定义
├── cmd/                # 应用入口
├── internal/           # 内部代码
│   ├── kafka/         # Kafka 封装
│   ├── storage/       # MinIO 封装
│   ├── middleware/    # 中间件
│   ├── handler/       # HTTP 处理器
│   ├── service/       # 业务逻辑
│   └── data/          # 数据层
├── deploy/            # 部署配置
└── docker-compose.yaml
```

## 🛠️ 开发工具

```bash
# 生成 Proto 代码
make proto

# 生成 Wire 代码
make wire

# 运行测试
make test

# 构建
make build

# Docker 构建
make docker
```

## 📊 监控

- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **MinIO**: http://localhost:9001 (minioadmin/minioadmin)

## 🔒 安全提示

1. 修改所有默认密码
2. 不要提交 `.env` 文件
3. 使用强 JWT 密钥
4. 启用 HTTPS

## 📝 自定义

1. 修改 `go.mod` 中的模块名
2. 更新 `README.md`
3. 修改配置文件
4. 实现你的业务逻辑

## 📄 许可证

MIT License

---

**作者**: LY7CRGL-7  
**版本**: v1.0.0
