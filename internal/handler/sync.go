// handler 包提供HTTP请求处理器
package handler

import (
	"video-service/internal/pkg/response"
	"video-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SyncDoubanMovies 手动触发豆瓣电影同步
// @Summary 同步豆瓣电影数据
// @Description 手动触发豆瓣电影数据同步任务
// @Tags 同步
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "同步任务已启动"
// @Failure 500 {object} response.Response "同步失败"
// @Router /api/sync/douban/movies [post]
func SyncDoubanMovies(c *gin.Context) {
	zap.L().Info("手动触发豆瓣电影同步", zap.String("ip", c.ClientIP()))

	// 创建豆瓣同步服务
	doubanSyncService := service.NewDoubanSyncService()

	// 执行同步（在goroutine中异步执行，避免阻塞请求）
	go func() {
		if err := doubanSyncService.SyncMovies(); err != nil {
			zap.L().Error("豆瓣电影同步失败", zap.Error(err))
		} else {
			zap.L().Info("豆瓣电影同步成功")
		}
	}()

	// 立即返回响应
	response.SuccessMsg(c, "同步任务已启动，正在后台执行", nil)
}
