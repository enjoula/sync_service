.PHONY: help build run test clean docker-build docker-up docker-down init lint

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make build        - 编译应用程序"
	@echo "  make run          - 运行应用程序"
	@echo "  make test         - 运行测试"
	@echo "  make clean        - 清理编译文件"
	@echo "  make docker-build - 构建Docker镜像"
	@echo "  make docker-up    - 启动Docker容器"
	@echo "  make docker-down  - 停止Docker容器"
	@echo "  make init         - 初始化Etcd配置"
	@echo "  make lint         - 运行代码检查"

# 编译应用程序
build:
	@echo "正在编译应用程序..."
	@go build -o bin/server ./cmd/server/main.go
	@echo "编译完成！二进制文件：bin/server"

# 运行应用程序
run:
	@echo "正在启动应用程序..."
	@go run ./cmd/server/main.go

# 运行测试
test:
	@echo "正在运行测试..."
	@go test -v ./...

# 清理编译文件
clean:
	@echo "正在清理..."
	@rm -rf bin/
	@rm -rf logs/*.log
	@echo "清理完成！"

# 构建Docker镜像
docker-build:
	@echo "正在构建Docker镜像..."
	@docker build -t video-service:latest -f deployments/docker/Dockerfile .
	@echo "Docker镜像构建完成！"

# 启动Docker容器
docker-up:
	@echo "正在启动Docker容器..."
	@cd deployments/docker && docker-compose up -d
	@echo "Docker容器已启动！"
	@echo "运行 'make init' 初始化Etcd配置"

# 停止Docker容器
docker-down:
	@echo "正在停止Docker容器..."
	@cd deployments/docker && docker-compose down
	@echo "Docker容器已停止！"

# 初始化Etcd配置
init:
	@echo "正在初始化Etcd配置..."
	@bash scripts/init_etcd.sh
	@echo "Etcd配置初始化完成！"

# 代码检查
lint:
	@echo "正在运行代码检查..."
	@go vet ./...
	@echo "代码检查完成！"

# 格式化代码
fmt:
	@echo "正在格式化代码..."
	@go fmt ./...
	@echo "代码格式化完成！"

# 更新依赖
deps:
	@echo "正在更新依赖..."
	@go mod tidy
	@go mod download
	@echo "依赖更新完成！"

