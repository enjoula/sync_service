// handler 包提供HTTP请求处理器
package handler

import (
	"video-service/internal/pkg/response"
	"video-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SyncDoubanMovies 手动触发豆瓣同步
// @Summary 同步豆瓣数据
// @Description 手动触发豆瓣数据同步任务（包括电影、电视、动漫、综艺、纪录片）
// @Tags 同步
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "同步任务已启动"
// @Failure 500 {object} response.Response "同步失败"
// @Router /api/sync/douban/movies [post]
func SyncDoubanMovies(c *gin.Context) {
	zap.L().Info("手动触发豆瓣同步", zap.String("ip", c.ClientIP()))

	// 创建豆瓣同步服务
	doubanSyncService := service.NewDoubanSyncService()

	// 执行同步（在goroutine中异步执行，避免阻塞请求）
	go func() {
		if err := doubanSyncService.SyncAll(); err != nil {
			zap.L().Error("豆瓣同步失败", zap.Error(err))
		} else {
			zap.L().Info("豆瓣同步成功")
		}
	}()

	// 立即返回响应
	response.SuccessMsg(c, "同步任务已启动，正在后台执行", nil)
}

// SyncRefreshEpisodeCounts 从豆瓣刷新未下架且未完结视频的集数
// @Summary 刷新视频集数
// @Description 查询 status!=0 且 is_completed!=1、类型不是 movie/tvshow 的视频，从豆瓣重新拉取 episode_count；不修改 updated_at
// @Tags 同步
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "刷新任务已启动"
// @Router /api/sync/douban/episode-counts [post]
func SyncRefreshEpisodeCounts(c *gin.Context) {
	zap.L().Info("手动触发豆瓣集数刷新", zap.String("ip", c.ClientIP()))

	doubanSyncService := service.NewDoubanSyncService()
	go func() {
		if err := doubanSyncService.RefreshEpisodeCountsFromDouban(); err != nil {
			zap.L().Error("豆瓣集数刷新失败", zap.Error(err))
		} else {
			zap.L().Info("豆瓣集数刷新成功")
		}
	}()

	response.SuccessMsg(c, "集数刷新任务已启动，正在后台执行", nil)
}
