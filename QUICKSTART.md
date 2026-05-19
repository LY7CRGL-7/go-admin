# 快速启动指南

## 🚀 5 分钟快速开始

### 前置要求

- Docker & Docker Compose
- Go 1.23+
- Make (可选)

### 步骤 1: 启动基础设施

```bash
# 启动所有基础设施（PostgreSQL, Redis, Kafka, MinIO, Prometheus, Grafana）
docker-compose up -d

# 查看服务状态
docker-compose ps
```

### 步骤 2: 安装依赖

```bash
# 下载 Go 依赖
make deps

# 或者手动执行
go mod download
go mod tidy
```

### 步骤 3: 配置数据库

基础设施启动后，数据库会自动初始化。如果需要手动初始化：

```bash
# 连接 PostgreSQL
docker exec -it admin-postgres psql -U admin -d admin_db

# 查看表
\dt
```

### 步骤 4: 运行应用

```bash
# 开发模式运行
make run

# 或者
go run cmd/admin/main.go -config cmd/admin/conf/config.yaml
```

### 步骤 5: 验证服务

```bash
# 健康检查
curl http://localhost:8080/health

# 查看 Prometheus 指标
curl http://localhost:8080/metrics

# 登录测试
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin@123456"}'
```

## 📊 访问监控面板

### Grafana
- URL: http://localhost:3000
- 用户名: `admin`
- 密码: `admin`

### Prometheus
- URL: http://localhost:9090

### MinIO Console
- URL: http://localhost:9001
- 用户名: `minioadmin`
- 密码: `minioadmin`

## 🔧 常用命令

```bash
# 停止所有服务
docker-compose down

# 查看日志
docker-compose logs -f

# 重启某个服务
docker-compose restart kafka

# 清理数据（谨慎使用）
docker-compose down -v

# 生成 Proto 代码
make proto

# 生成 Wire 代码
make wire

# 构建应用
make build

# 运行测试
make test
```

## ⚙️ 配置说明

### 启用 Kafka

编辑 `cmd/admin/conf/config.yaml`:

```yaml
kafka:
  enabled: true  # 改为 true
  brokers:
    - localhost:9092
  audit_log_topic: audit-logs
```

### 启用 MinIO

编辑 `cmd/admin/conf/config.yaml`:

```yaml
minio:
  enabled: true  # 改为 true
  endpoint: localhost:9000
  access_key_id: minioadmin
  secret_access_key: minioadmin
  use_ssl: false
  bucket_name: admin-files
```

## 🐛 故障排查

### 服务启动失败

```bash
# 查看日志
docker-compose logs postgres
docker-compose logs kafka

# 检查端口占用
netstat -tulpn | grep -E '5432|6379|9092|9000'
```

### 数据库连接失败

```bash
# 测试数据库连接
docker exec -it admin-postgres pg_isready -U admin

# 查看数据库状态
docker exec -it admin-postgres psql -U admin -d admin_db -c "SELECT version();"
```

### Kafka 连接失败

```bash
# 查看 Kafka 日志
docker-compose logs kafka

# 测试 Kafka 连接
docker exec -it admin-kafka kafka-topics --bootstrap-server localhost:9092 --list
```

## 📝 下一步

1. 查看 [README.md](README.md) 了解完整功能
2. 查看 [REFACTORING.md](REFACTORING.md) 了解重构详情
3. 查看 API 文档开始开发
4. 配置 Grafana 监控面板

## 💡 提示

- 开发时不需要启动所有服务，可以只启动需要的组件
- Kafka 和 MinIO 默认禁用，按需启用
- 完整启动所有服务需要约 4GB 内存
- 建议使用 SSD 以获得更好的数据库性能

---

有问题？查看 [REFACTORING.md](REFACTORING.md) 或联系开发团队。
