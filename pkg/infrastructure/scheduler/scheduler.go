// scheduler 包提供定时任务调度功能
// 使用cron库实现定时任务
package scheduler

import (
	"video-service/internal/service"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	// cron调度器实例
	cronScheduler *cron.Cron
)

// InitCron 初始化定时任务调度器
func InitCron() {
	cronScheduler = cron.New(cron.WithSeconds())

	// 添加豆瓣同步任务：每8小时执行一次
	// Cron表达式: 0 0 6,14,22 * * * (每8小时执行一次)
	doubanSyncService := service.NewDoubanSyncService()
	_, err := cronScheduler.AddFunc("0 0 6,14,22 * * *", func() {
		zap.L().Info("开始执行豆瓣同步任务")
		if err := doubanSyncService.SyncAll(); err != nil {
			zap.L().Error("豆瓣同步任务执行失败", zap.Error(err))
		} else {
			zap.L().Info("豆瓣同步任务执行成功")
		}
	})
	if err != nil {
		zap.L().Error("添加豆瓣同步定时任务失败", zap.Error(err))
	} else {
		zap.L().Info("豆瓣同步定时任务已添加", zap.String("schedule", "6,14,22,每8小时执行一次"))
	}

	// 启动调度器
	cronScheduler.Start()
	zap.L().Info("定时任务调度器已启动")
}

// Stop 停止定时任务调度器
func Stop() {
	if cronScheduler != nil {
		// 停止调度器（等待正在执行的任务完成）
		ctx := cronScheduler.Stop()
		<-ctx.Done()
		zap.L().Info("定时任务调度器已停止")
	}
}

// AddFunc 添加定时任务
//
//	spec: cron表达式（支持秒级精度，格式：秒 分 时 日 月 周）
//
// 返回：任务ID和错误
func AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	if cronScheduler == nil {
		InitCron()
	}
	return cronScheduler.AddFunc(spec, cmd)
}
