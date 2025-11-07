// router 包提供HTTP路由配置
// 设置Gin路由、中间件和API端点
package router

import (
	"video-service/internal/api"
	"video-service/internal/middleware"

	"github.com/gin-gonic/gin"
	ginprom "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
)

// SetupRouter 配置并返回Gin路由引擎
// 功能包括：
// 1. 注册全局中间件（Trace、Recovery、Logger）
// 2. 集成Prometheus监控
// 3. 注册公开API端点（注册、登录、健康检查）
// 4. 注册需要认证的API端点（用户信息）
// 5. 注册监控指标端点
func SetupRouter() *gin.Engine {
	// 创建新的Gin引擎（不使用默认中间件）
	r := gin.New()
	log := zap.L()

	// 配置信任的代理IP地址，以便正确获取客户端真实IP
	// 信任本地和私有网络IP段（适用于Nginx、Docker等反向代理场景）
	if err := r.SetTrustedProxies([]string{
		"127.0.0.1",      // 本地回环
		"10.0.0.0/8",     // 私有网络 A类
		"172.16.0.0/12",  // 私有网络 B类
		"192.168.0.0/16", // 私有网络 C类
	}); err != nil {
		log.Warn("设置信任代理失败", zap.Error(err))
	}

	// 注册全局中间件（按顺序执行）
	// Trace: 生成请求追踪ID
	r.Use(middleware.Trace())
	// RecoveryWithZap: 捕获panic并记录日志
	r.Use(middleware.RecoveryWithZap(log))
	// LoggerWithZap: 记录HTTP请求日志
	r.Use(middleware.LoggerWithZap(log))

	// 集成Prometheus监控中间件，服务名称为"video_service"
	p := ginprom.NewPrometheus("video_service")
	p.Use(r)

	// 注册公开API端点（无需认证）
	r.POST("/register", api.Register) // 用户注册
	r.POST("/login", api.Login)       // 用户登录
	r.GET("/ping", api.Ping)          // 健康检查
	r.GET("/ip-info", api.GetIPInfo)  // IP信息查询（用于测试）

	// 创建需要认证的路由组
	auth := r.Group("/user")
	// 在该路由组上应用JWT认证中间件
	auth.Use(middleware.JWTAuth())
	// 注册需要认证的API端点
	auth.GET("/me", api.Me) // 获取当前用户信息

	// 注意：/metrics 路由已由 go-gin-prometheus 中间件自动注册，无需手动注册
	// r.GET("/metrics", gin.WrapH(metrics.Handler()))

	return r
}
