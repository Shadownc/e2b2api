# 第一阶段：构建Go应用
FROM golang:1.19-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（禁用CGO以确保静态链接，为linux/amd64平台构建）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o e2b-gateway main.go

# 第二阶段：创建最小运行镜像
FROM alpine:latest

# 安装必要的CA证书用于HTTPS请求
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/e2b-gateway .

# 复制.env.example文件作为参考配置
COPY .env.example .

# 暴露应用端口
EXPOSE 8080

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1

# 运行应用
ENTRYPOINT ["./e2b-gateway"] 