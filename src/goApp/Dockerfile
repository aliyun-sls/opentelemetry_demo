# 使用官方的Go镜像作为基础镜像
FROM golang:1.19-alpine

# 设置工作目录
WORKDIR /app

# 将Go模块文件复制到工作目录
COPY go.mod .
COPY go.sum .

# 下载依赖
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 构建应用程序
RUN go build -o main .

# 暴露端口
EXPOSE 8080

# 运行应用程序
CMD ["./main"]