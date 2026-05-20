# ========================================
# 模板项目 Dockerfile
# 使用者可以根据自己的需求修改此文件
# ========================================

# 构建阶段
FROM golang:1.25-alpine AS builder

# 设置构建参数（使用者可以修改）
ARG APP_NAME=admin
ARG APP_VERSION=1.0.0

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN mkdir -p /app/bin && \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${APP_VERSION}" \
    -o /app/bin/${APP_NAME} \
    cmd/admin/main.go

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从 builder 阶段复制二进制文件
COPY --from=builder /app/bin/${APP_NAME} .
COPY --from=builder /app/cmd/admin/conf/config.yaml ./config.yaml

# 暴露端口（根据您的服务修改）
EXPOSE 8080

# 运行
CMD ["./admin", "-config", "config.yaml"]
