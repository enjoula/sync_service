#!/bin/bash

# 豆瓣电影同步功能测试脚本

# 服务地址
SERVER_URL="http://localhost:8080"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== 豆瓣电影同步功能测试 ===${NC}\n"

# 测试1: 健康检查
echo -e "${YELLOW}[1] 测试服务健康状态...${NC}"
response=$(curl -s "${SERVER_URL}/ping")
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 服务正常运行${NC}"
    echo "响应: $response"
else
    echo -e "${RED}✗ 服务未启动或无法访问${NC}"
    echo "请先启动服务: make run 或 ./bin/server"
    exit 1
fi
echo ""

# 测试2: 手动触发同步
echo -e "${YELLOW}[2] 手动触发豆瓣电影同步...${NC}"
response=$(curl -s -X POST "${SERVER_URL}/api/sync/douban/movies" \
    -H "Content-Type: application/json")

if [ $? -eq 0 ]; then
    # 检查响应中的code字段
    code=$(echo "$response" | grep -o '"code":[0-9]*' | grep -o '[0-9]*')
    if [ "$code" = "0" ]; then
        echo -e "${GREEN}✓ 同步任务已成功启动${NC}"
        echo "响应: $response"
    else
        echo -e "${RED}✗ 同步任务启动失败${NC}"
        echo "响应: $response"
    fi
else
    echo -e "${RED}✗ 请求失败${NC}"
    exit 1
fi
echo ""

# 测试3: 查看日志
echo -e "${YELLOW}[3] 查看同步日志（最近20行）...${NC}"
if [ -f "logs/app.log" ]; then
    echo -e "${GREEN}最新日志:${NC}"
    tail -20 logs/app.log | grep --color=always -E "豆瓣|电影|同步|$"
else
    echo -e "${YELLOW}⚠ 日志文件不存在${NC}"
fi
echo ""

# 测试4: 监控日志
echo -e "${YELLOW}[4] 实时监控同步日志...${NC}"
echo -e "${YELLOW}提示: 按 Ctrl+C 退出监控${NC}\n"

if [ -f "logs/app.log" ]; then
    tail -f logs/app.log | grep --color=always -E "豆瓣|电影|同步|错误|失败|成功"
else
    echo -e "${RED}✗ 日志文件不存在${NC}"
fi

