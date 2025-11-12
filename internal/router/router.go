// router 包提供HTTP路由配置
// 设置Gin路由、中间件和API端点
package router

import (
	"video-service/internal/handler"
	"video-service/internal/middleware"

	"github.com/gin-gonic/gin"
	ginprom "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
)

// SetupRouter 配置并返回Gin路由引擎
// 功能包括：
// 1. 注册全局中间件（Trace、Recovery、Logger）
// 2. 集成Prometheus监控
// 3. 注册公开API端点（健康检查、IP信息查询）
// 4. 注册用户相关API端点（/user/register、/user/login、/user/me）
// 5. 注册监控指标端点（由go-gin-prometheus自动注册）
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
	r.GET("/ping", handler.Ping)         // 健康检查
	r.GET("/ip-info", handler.GetIPInfo) // IP信息查询（用于测试）

	// 创建用户相关路由组
	userGroup := r.Group("/user")
	{
		// 公开接口（无需认证）
		userGroup.POST("/register", handler.Register) // 用户注册
		userGroup.POST("/login", handler.Login)       // 用户登录

		// 需要认证的接口
		// 在该路由组上应用JWT认证中间件
		authGroup := userGroup.Group("")
		authGroup.Use(middleware.JWTAuth())
		authGroup.GET("/me", handler.Me) // 获取当前用户信息
	}

	// 注意：/metrics 路由已由 go-gin-prometheus 中间件自动注册，无需手动注册
	// r.GET("/metrics", gin.WrapH(metrics.Handler()))

	return r
}
