// cache 包提供Redis缓存连接和初始化功能
package cache

import (
	"context"
	"video-service/internal/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Rdb 是全局Redis客户端对象
var Rdb *redis.Client

// InitRedis 初始化Redis缓存连接
// 功能包括：
// 1. 从配置中获取Redis连接信息（地址、密码、数据库编号）
// 2. 创建Redis客户端连接
// 3. 通过Ping命令测试连接是否正常
func InitRedis() {
	// 从配置中获取Redis地址
	addr := config.Cfg.GetString("redis.addr")

	// 如果未配置Redis地址，记录警告并返回
	if addr == "" {
		zap.L().Warn("redis addr empty")
		return
	}

	// 创建Redis客户端，配置连接选项
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,                               // Redis服务器地址
		Password: config.Cfg.GetString("redis.pass"), // Redis密码（如果配置了）
		DB:       config.Cfg.GetInt("redis.db"),      // Redis数据库编号（默认0）
	})

	// 创建上下文用于连接测试
	ctx := context.Background()

	// 通过Ping命令测试Redis连接
	if err := Rdb.Ping(ctx).Err(); err != nil {
		zap.L().Fatal("redis connect error", zap.Error(err))
	}

	zap.L().Info("redis connected")
}
