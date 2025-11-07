// main 包是应用程序的入口点
package main

import (
	"video-service/internal/router"
	"video-service/pkg/infrastructure/cache"
	"video-service/pkg/infrastructure/config"
	"video-service/pkg/infrastructure/database"
	"video-service/pkg/infrastructure/logger"
	"video-service/pkg/infrastructure/metrics"
	"video-service/pkg/infrastructure/scheduler"

	"go.uber.org/zap"
)

// main 函数是应用程序的启动入口
// 按照以下顺序初始化各个组件：
// 1. 配置管理（支持本地配置文件和Etcd远程配置）
// 2. 日志系统（文件和控制台双输出）
// 3. 数据库连接（MySQL，包含自动迁移）
// 4. 缓存连接（Redis）
// 5. 监控指标（Prometheus）
// 6. 定时任务调度器
// 7. HTTP路由和服务启动
func main() {
	// 初始化配置管理，支持从配置文件和环境变量读取，并可从Etcd获取敏感信息
	config.InitConfig()

	// 初始化日志系统，配置文件和控制台双输出
	logger.InitLogger()
	log := zap.L()

	// 初始化MySQL数据库连接，并执行自动迁移创建表结构
	database.InitMySQL()

	// 初始化Redis缓存连接
	cache.InitRedis()

	// 初始化Prometheus监控指标
	metrics.InitMetrics()

	// 启动定时任务调度器
	scheduler.InitCron()
	// 确保程序退出时停止定时任务
	defer scheduler.Stop()

	// 设置HTTP路由，包括中间件和API端点
	r := router.SetupRouter()

	// 从配置中获取服务器监听地址
	addr := config.Cfg.GetString("server.addr")

	// 记录服务器启动日志
	log.Info("server start", zap.String("addr", addr))

	// 启动HTTP服务器，如果启动失败则记录致命错误并退出
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}
