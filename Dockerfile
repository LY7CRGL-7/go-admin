# ========================================
# Kratos Admin Template - Dockerfile
# 工业标准：固定产物路径，CI 注入版本
# ========================================

# ---------- Build ----------
FROM golang:1.25-alpine AS builder

ARG VERSION=dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X main.Version=${VERSION}" \
    -o /app ./cmd/admin

# ---------- Runtime ----------
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app .
COPY cmd/admin/conf/config.yaml.example config.yaml

EXPOSE 8080 9090

ENTRYPOINT ["./admin"]
CMD ["-conf", "config.yaml"]
