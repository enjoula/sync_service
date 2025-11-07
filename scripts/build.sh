#!/bin/bash

# 编译脚本
# 支持多平台编译

set -e

# 项目信息
PROJECT_NAME="video-service"
BUILD_DIR="bin"
MAIN_FILE="cmd/server/main.go"

# 版本信息
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(go version | awk '{print $3}')

# LDFLAGS
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GoVersion=${GO_VERSION}"

echo "==================================="
echo "编译 ${PROJECT_NAME}"
echo "==================================="
echo "版本: ${VERSION}"
echo "构建时间: ${BUILD_TIME}"
echo "Go版本: ${GO_VERSION}"
echo "==================================="

# 创建输出目录
mkdir -p ${BUILD_DIR}

# 编译当前平台
echo "正在编译当前平台..."
go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/server ${MAIN_FILE}
echo "✅ 编译完成: ${BUILD_DIR}/server"

# 可选：交叉编译其他平台
if [ "$1" == "all" ]; then
    echo ""
    echo "开始交叉编译..."
    
    # Linux AMD64
    echo "编译 Linux AMD64..."
    GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/server-linux-amd64 ${MAIN_FILE}
    echo "✅ ${BUILD_DIR}/server-linux-amd64"
    
    # Linux ARM64
    echo "编译 Linux ARM64..."
    GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/server-linux-arm64 ${MAIN_FILE}
    echo "✅ ${BUILD_DIR}/server-linux-arm64"
    
    # macOS AMD64
    echo "编译 macOS AMD64..."
    GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/server-darwin-amd64 ${MAIN_FILE}
    echo "✅ ${BUILD_DIR}/server-darwin-amd64"
    
    # macOS ARM64 (Apple Silicon)
    echo "编译 macOS ARM64..."
    GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/server-darwin-arm64 ${MAIN_FILE}
    echo "✅ ${BUILD_DIR}/server-darwin-arm64"
    
    echo ""
    echo "✅ 所有平台编译完成！"
fi

echo ""
echo "运行: ./${BUILD_DIR}/server"

