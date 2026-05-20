FROM golang:1.25-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=1.0.0" -o /app/bin/admin cmd/admin/main.go

# 最终镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从 builder 阶段复制二进制文件
COPY --from=builder /app/bin/admin .
COPY --from=builder /app/cmd/admin/conf/config.yaml ./config.yaml

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./admin", "-config", "config.yaml"]
