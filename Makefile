# ========================================
# 模板项目 Makefile
# 使用者可以根据自己的需求修改此文件
# ========================================

# 项目配置（使用者需要修改这些变量）
BINARY_NAME ?= your-project-name  # 修改为您的项目名称
MAIN_PATH ?= cmd/admin/main.go    # 主程序入口
CONFIG_PATH ?= cmd/admin/conf/config.yaml  # 配置文件路径
VERSION ?= 1.0.0

.PHONY: build run clean test docker proto wire deps fmt help

# 构建
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build completed: bin/$(BINARY_NAME)"

# 运行
run:
	@echo "🚀 Starting $(BINARY_NAME)..."
	go run $(MAIN_PATH) -config $(CONFIG_PATH)

# 清理
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	@echo "✅ Clean completed!"

# 测试
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# 格式化代码
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# 下载依赖
deps:
	@echo "📦 Downloading dependencies..."
	go mod download

# 生成 proto 代码
proto:
	@echo "📝 Generating proto code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/admin/v1/admin.proto
	@echo "✅ Proto generation completed!"

# 生成 wire 代码
wire:
	@echo "🔌 Generating wire code..."
	cd cmd/admin && wire
	@echo "✅ Wire generation completed!"

# 代码检查（可选）
lint:
	@echo "🔍 Running linter..."
	golangci-lint run

# Docker 构建（使用者需要修改镜像名）
docker:
	@echo "🐳 Building Docker image..."
	docker build -t YOUR_USERNAME/$(BINARY_NAME):latest .
	@echo "✅ Docker image built: YOUR_USERNAME/$(BINARY_NAME):latest"

# 帮助
help:
	@echo ""
	@echo "📋 Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make fmt      - Format code"
	@echo "  make deps     - Download dependencies"
	@echo "  make proto    - Generate proto code"
	@echo "  make wire     - Generate wire code"
	@echo "  make lint     - Run linter (requires golangci-lint)"
	@echo "  make docker   - Build Docker image"
	@echo "  make help     - Show this help"
	@echo ""
	@echo "💡 Quick Start:"
	@echo "  1. make deps          # Install dependencies"
	@echo "  2. make proto         # Generate proto code"
	@echo "  3. cp cmd/admin/conf/config.yaml.example cmd/admin/conf/config.yaml"
	@echo "  4. make run           # Run the application"
	@echo ""
