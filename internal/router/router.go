// router 包提供HTTP路由配置
// 设置Gin路由、中间件和API端点
package router

import (
	"video-service/internal/handler"
	"video-service/internal/handler/admin"
	"video-service/internal/middleware"
	"video-service/internal/model"

	"github.com/gin-gonic/gin"
	ginprom "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
)

// SetupRouter 配置并返回Gin路由引擎
// 功能包括：
// 1. 注册全局中间件（Trace、Recovery、Logger）
// 2. 集成Prometheus监控
// 3. 注册公开API端点（健康检查）
// 4. 注册监控指标端点（由go-gin-prometheus自动注册）
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
	// r.GET("/ping", handler.Ping) // 健康检查

	// API路由组
	apiGroup := r.Group("/api")
	{
		// 同步相关接口
		syncGroup := apiGroup.Group("/sync")
		{
			// 豆瓣电影同步接口（手动触发）
			syncGroup.POST("/douban/movies", handler.SyncDoubanMovies)
			// 从豆瓣刷新未下架且未完结视频的集数（不修改 updated_at）
			syncGroup.POST("/douban/episode-counts", handler.SyncRefreshEpisodeCounts)
		}

		// 管理后台接口。除登录外都需要管理员令牌，并按角色权限控制。
		adminGroup := apiGroup.Group("/admin")
		{
			authHandler := admin.NewAuthHandler()
			adminGroup.POST("/login", authHandler.Login)
			// 封面由 <img> 加载，不能带 Authorization；豆瓣直链有防盗链/CORS，需服务端代理
			adminGroup.GET("/cover", admin.Cover)

			authed := adminGroup.Group("")
			authed.Use(middleware.AdminAuth())
			authed.GET("/me", authHandler.Me)
			authed.PUT("/password", authHandler.ChangePassword)

			schedulerHandler := admin.NewSchedulerHandler()
			schedulerGroup := authed.Group("")
			schedulerGroup.Use(middleware.RequireSuperAdmin())
			schedulerGroup.GET("/scheduler", schedulerHandler.Get)
			schedulerGroup.PUT("/scheduler", schedulerHandler.Update)

			videoHandler := admin.NewVideoHandler()
			episodeHandler := admin.NewEpisodeHandler()

			videoRead := authed.Group("")
			videoRead.Use(middleware.RequireAdminPerm(model.PermVideoRead))
			videoRead.GET("/videos", videoHandler.List)
			videoRead.GET("/videos/:id", videoHandler.Get)
			videoRead.GET("/episodes", episodeHandler.ListByVideoID)
			videoRead.GET("/episodes/:id", episodeHandler.Get)

			videoWrite := authed.Group("")
			videoWrite.Use(middleware.RequireAdminPerm(model.PermVideoWrite))
			videoWrite.POST("/videos", videoHandler.Create)
			videoWrite.PUT("/videos/:id", videoHandler.Update)
			videoWrite.DELETE("/videos/:id", videoHandler.Delete)
			videoWrite.POST("/episodes", episodeHandler.Create)
			videoWrite.PUT("/episodes/:id", episodeHandler.Update)
			videoWrite.DELETE("/episodes/:id", episodeHandler.Delete)

			userRead := authed.Group("")
			userRead.Use(middleware.RequireAdminPerm(model.PermUserRead))
			userHandler := admin.NewUserHandler()
			userRead.GET("/users", userHandler.List)
			userRead.GET("/users/watch-records", userHandler.ListWatchRecords)
			userRead.GET("/users/logins", userHandler.ListLogins)

			accountHandler := admin.NewAccountHandler()
			accountManage := authed.Group("")
			accountManage.Use(middleware.RequireAdminPerm(model.PermAdminManage))
			accountManage.GET("/accounts", accountHandler.List)
			accountManage.POST("/accounts", accountHandler.Create)
			accountManage.PUT("/accounts/:id", accountHandler.Update)
		}
	}

	// 静态文件服务 - 先注册，确保静态文件请求不会被NoRoute捕获
	r.Static("/static", "./static")

	// 管理后台只从 /super 提供，未匹配的地址返回 404
	serveAdmin := func(c *gin.Context) {
		c.File("./static/super/index.html")
	}
	r.GET("/super", serveAdmin)
	r.GET("/super/", serveAdmin)
	r.GET("/super/admin.html", serveAdmin)

	// 管理后台静态资源
	r.Static("/super/assets", "./static/super/assets")

	// 禁用根路径页面访问，避免直接打开服务地址返回前端页面
	r.GET("/", func(c *gin.Context) {
		c.Status(404)
	})

	r.NoRoute(func(c *gin.Context) {
		c.Status(404)
	})

	return r
}
