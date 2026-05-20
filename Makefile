# ========================================
# Kratos 企业级管理端模板 - Makefile
# ========================================

BINARY_NAME = admin
MAIN_PATH   = ./cmd/admin
CONF_PATH   = ./cmd/admin/conf/config.yaml
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

PROTO_DIR   = api/admin/v1
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

.PHONY: all build run clean test fmt deps proto wire generate lint docker help

# 默认目标
all: generate build

# ==================== 构建 ====================

build:
	@echo "🔨 Building $(BINARY_NAME) $(VERSION)..."
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build completed: bin/$(BINARY_NAME)"

run:
	@echo "🚀 Starting $(BINARY_NAME)..."
	go run $(MAIN_PATH) -conf $(CONF_PATH)

clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	@echo "✅ Clean completed"

# ==================== 代码质量 ====================

test:
	@echo "🧪 Running tests..."
	go test -v -cover ./...

fmt:
	@echo "🎨 Formatting..."
	gofmt -w .

lint:
	@echo "🔍 Linting..."
	golangci-lint run ./...

# ==================== 代码生成 ====================

proto:
	@echo "📝 Generating proto code..."
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=paths=source_relative:$(PROTO_DIR) \
		--go-grpc_out=paths=source_relative:$(PROTO_DIR) \
		$(PROTO_FILES)
	@echo "✅ Proto generation completed"

wire:
	@echo "🔌 Generating wire code..."
	cd cmd/admin && wire
	@echo "✅ Wire generation completed"

# proto + wire 一键生成
generate: proto wire

deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy

# ==================== Docker ====================

docker:
	@echo "🐳 Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .
	@echo "✅ Docker image built: $(BINARY_NAME):$(VERSION)"

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

# ==================== 帮助 ====================

help:
	@echo ""
	@echo "📋 Kratos Admin Template - Available Commands"
	@echo "============================================="
	@echo ""
	@echo "  Build & Run:"
	@echo "    make build          Build binary"
	@echo "    make run            Run application"
	@echo "    make clean          Clean build artifacts"
	@echo ""
	@echo "  Code Quality:"
	@echo "    make test           Run tests"
	@echo "    make fmt            Format code"
	@echo "    make lint           Run linter"
	@echo ""
	@echo "  Code Generation:"
	@echo "    make proto          Generate proto code"
	@echo "    make wire           Generate wire code"
	@echo "    make generate       Generate all (proto + wire)"
	@echo ""
	@echo "  Docker:"
	@echo "    make docker              Build Docker image"
	@echo "    make docker-compose-up   Start all services"
	@echo "    make docker-compose-down Stop all services"
	@echo ""
	@echo "  Quick Start:"
	@echo "    1. make deps"
	@echo "    2. cp cmd/admin/conf/config.yaml.example cmd/admin/conf/config.yaml"
	@echo "    3. docker-compose up -d  (start PostgreSQL & Redis)"
	@echo "    4. make run"
	@echo ""
