# 第一阶段 编译go程序
FROM golang:1.26 AS builder
# 设置容器内工作目录
WORKDIR /app
# 设置国内镜像
ENV GOPROXY=https://goproxy.cn,direct
# COPY [target] [dest]
# 先复制go.mod go.sum等不常变化的文件
COPY go.mod go.sum ./
RUN go mod download
# 复制全部项目代码
COPY . .
# 关闭CGO,不依赖系统C库，后续使用alpine，alpine使用musl libc，Golang基于Debian使用glibc(GNU C Library)，如果程序动态链接glibc在alpine会报错
RUN CGO_ENABLED=0 GOOS=linux go build -o server .
#第二阶段构建 搭建运行环境
#使用一个极小的Linux镜像
FROM alpine:3.22
# 设置工作区
WORKDIR /app
# alpine内运行命令 安装时区数据（程序里面使用了time）
RUN apk add --no-cache tzdata
#从第一阶段编译的内容复制到第二阶段（未复制默认全部丢弃）
COPY --from=builder /app/server .
COPY --from=builder /app/web ./web
COPY --from=builder /app/database ./database
# 声称容器监听8080端口
EXPOSE 8080
# CMD 容器启动时执行的命令
CMD ["./server"]
