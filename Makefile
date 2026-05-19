.PHONY: build run clean test docker proto wire

# 构建
build:
	@echo "Building admin service..."
	go build -ldflags "-X main.Version=1.0.0" -o bin/admin cmd/admin/main.go
	@echo "Build completed!"

# 运行
run:
	@echo "Starting admin service..."
	go run cmd/admin/main.go -config cmd/admin/conf/config.yaml

# 清理
clean:
	@echo "Cleaning..."
	rm -rf bin/
	@echo "Clean completed!"

# 测试
test:
	@echo "Running tests..."
	go test -v ./...

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 下载依赖
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# 生成 proto 代码
proto:
	@echo "Generating proto code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/admin/v1/admin.proto
	@echo "Proto generation completed!"

# 生成 wire 代码
wire:
	@echo "Generating wire code..."
	cd cmd/admin && wire
	@echo "Wire generation completed!"

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run

# Docker 构建
docker:
	@echo "Building Docker image..."
	docker build -t admin:latest .

# 帮助
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make fmt      - Format code"
	@echo "  make deps     - Download dependencies"
	@echo "  make proto    - Generate proto code"
	@echo "  make wire     - Generate wire code"
	@echo "  make lint     - Run linter"
	@echo "  make docker   - Build Docker image"
