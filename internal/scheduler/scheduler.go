// scheduler 包提供定时任务调度功能
// 使用cron库实现定时任务
package scheduler

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	// cron调度器实例
	cronScheduler *cron.Cron
)

// InitCron 初始化定时任务调度器
// 功能：创建cron调度器并启动
func InitCron() {
	// 创建cron调度器，使用秒级精度
	cronScheduler = cron.New()

	// 可以在这里添加定时任务
	// 示例：每天凌晨2点执行清理任务
	// _, err := cronScheduler.AddFunc("0 0 2 * * *", func() {
	//     zap.L().Info("执行定时清理任务")
	//     // 执行清理逻辑
	// })
	// if err != nil {
	//     zap.L().Error("添加定时任务失败", zap.Error(err))
	// }

	// 启动调度器
	cronScheduler.Start()
	zap.L().Info("定时任务调度器已启动")
}

// Stop 停止定时任务调度器
// 功能：优雅地停止所有定时任务
func Stop() {
	if cronScheduler != nil {
		// 停止调度器（等待正在执行的任务完成）
		ctx := cronScheduler.Stop()
		<-ctx.Done()
		zap.L().Info("定时任务调度器已停止")
	}
}

// AddFunc 添加定时任务
// 参数：
//
//	spec: cron表达式（支持秒级精度，格式：秒 分 时 日 月 周）
//	cmd: 要执行的函数
//
// 返回：任务ID和错误
func AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	if cronScheduler == nil {
		InitCron()
	}
	return cronScheduler.AddFunc(spec, cmd)
}
