# 多阶段构建：构建阶段
# 使用golang:1.21-alpine作为构建镜像
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件（利用Docker缓存层）
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 获取构建参数（目标平台和架构）
ARG TARGETOS
ARG TARGETARCH

# 交叉编译：根据目标平台构建二进制文件
# CGO_ENABLED=0 禁用CGO，生成静态链接的二进制文件，便于跨平台部署
# GOOS和GOARCH根据构建参数动态设置
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/server ./cmd/server

# 多阶段构建：运行阶段
# 使用alpine作为运行时镜像（体积小）
FROM --platform=$TARGETPLATFORM alpine:3.18

# 安装CA证书（用于HTTPS连接）和时区数据
RUN apk add --no-cache ca-certificates tzdata

# 设置时区为东八区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/server .

# 复制配置文件目录
COPY configs ./configs

# 暴露端口
EXPOSE 8080

# 设置入口点
ENTRYPOINT ["./server"]
