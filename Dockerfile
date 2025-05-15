FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod ./
COPY go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN go build -o toolbox-api ./cmd/main.go

# 使用更小的基础镜像
FROM alpine:latest

WORKDIR /app

# 安装依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建toolbox_data目录
RUN mkdir -p /app/toolbox_data

# 从构建阶段复制二进制文件
COPY --from=builder /app/toolbox-api .

# 暴露端口
EXPOSE 4005

# 运行服务
CMD ["./toolbox-api"] 