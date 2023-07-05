# Build阶段
FROM golang:1.20 AS build

# 设置工作目录
WORKDIR /app

# 将go mod和src文件复制到工作目录
COPY go.mod .
COPY src ./src

# 下载所有依赖包
RUN go mod download

# 编译程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/main.go

# Final阶段
FROM alpine:latest AS final

# 在/alpine中创建/app目录
WORKDIR /app

# 从build阶段复制编译好的二进制文件到当前目录
COPY --from=build /app/main .

# 暴露端口
# EXPOSE 8080

# 运行程序
CMD ["./main"]