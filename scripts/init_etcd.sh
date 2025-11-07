#!/bin/bash

# 初始化Etcd配置脚本
# 用于在Etcd中设置JWT密钥和数据库连接信息

echo "==================================="
echo "初始化Etcd配置"
echo "==================================="

# 检查Etcd容器是否运行
if ! docker ps | grep -q "Etcd"; then
    echo "错误: Etcd容器未运行"
    echo "请先运行: make docker-up"
    exit 1
fi

# 设置Etcd密钥和值
ETCD_KEY="/video-service/secret"
ETCD_VALUE='{
  "jwt_key": "super-secret-key-change-me-in-production",
  "mysql_dsn": "root:123456@tcp(mysql:3306)/video_service?charset=utf8mb4&parseTime=True&loc=Local"
}'

echo "正在设置Etcd配置..."
echo "Key: $ETCD_KEY"

# 写入Etcd
docker exec Etcd etcdctl put "$ETCD_KEY" "$ETCD_VALUE"

if [ $? -eq 0 ]; then
    echo "✅ Etcd配置设置成功！"
    echo ""
    echo "提示:"
    echo "1. 请在生产环境中修改jwt_key为安全的密钥"
    echo "2. 如需修改配置，可直接编辑此脚本后重新运行"
    echo ""
    echo "验证配置:"
    echo "docker exec Etcd etcdctl get /video-service/secret"
else
    echo "❌ Etcd配置设置失败"
    exit 1
fi

