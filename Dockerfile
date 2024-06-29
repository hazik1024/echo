# 使用官方的Golang运行时作为父镜像
FROM golang:latest AS build

# 设置工作目录为/app
WORKDIR /app

# 将当前目录内容复制到容器的/app内
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o echoserver .

# 使用Alpine Linux作为生产环境的镜像，减小镜像体积
FROM alpine:latest AS production

RUN mkdir /server

# 将上一步构建的应用复制到当前镜像的/server/目录下
COPY --from=build /app/echoserver /server/
# 设置执行权限
RUN chmod +x /server/echoserver

RUN ls -l /server

# 设置容器启动时运行的应用
CMD ["/server/echoserver"]
