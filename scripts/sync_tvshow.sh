#!/bin/bash

# 豆瓣热门剧集同步脚本

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}    豆瓣热门剧集同步脚本${NC}"
echo -e "${BLUE}================================================${NC}\n"

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 切换到项目根目录
cd "$PROJECT_ROOT" || exit 1

echo -e "${YELLOW}[1] 检查配置文件...${NC}"
if [ ! -f "configs/config.yaml" ]; then
    echo -e "${RED}✗ 配置文件不存在: configs/config.yaml${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 配置文件存在${NC}"

# 检查是否设置了MYSQL_DSN环境变量
if [ -z "$MYSQL_DSN" ]; then
    # 默认使用本地5506端口
    export MYSQL_DSN="root:123456@tcp(localhost:5506)/video_service?charset=utf8mb4&parseTime=True&loc=Local"
    echo -e "${YELLOW}提示: 使用默认数据库配置 (localhost:5506)${NC}"
    echo -e "${YELLOW}如需自定义，可设置环境变量 MYSQL_DSN${NC}\n"
else
    echo -e "${GREEN}✓ 检测到环境变量 MYSQL_DSN${NC}\n"
fi

echo -e "${YELLOW}[2] 编译脚本...${NC}"
go build -o /tmp/sync_tvshow scripts/sync_tvshow.go
if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 编译失败${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 编译成功${NC}\n"

echo -e "${YELLOW}[3] 执行同步任务...${NC}"
echo -e "${BLUE}------------------------------------------------${NC}"
/tmp/sync_tvshow
EXIT_CODE=$?
echo -e "${BLUE}------------------------------------------------${NC}\n"

# 清理临时文件
rm -f /tmp/sync_tvshow

if [ $EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✓ 同步任务完成${NC}"
else
    echo -e "${RED}✗ 同步任务失败 (退出码: $EXIT_CODE)${NC}"
    exit $EXIT_CODE
fi

echo -e "\n${BLUE}================================================${NC}"
echo -e "${GREEN}  同步完成！${NC}"
echo -e "${BLUE}================================================${NC}"

