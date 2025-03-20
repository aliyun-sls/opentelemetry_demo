# 使用 Go 官方基础镜像
FROM golang:1.22

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod .
COPY go.sum .

# 复制 common 目录
COPY common/ ./common/

# 下载依赖
RUN go mod download

# 复制所有文件到工作目录
COPY . .

# 构建应用
RUN go build -o main .

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./main"]