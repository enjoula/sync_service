// scheduler 包提供定时任务调度功能
// 使用cron库实现定时任务
package scheduler

import (
	"sync"
	"sync/atomic"

	"video-service/internal/service"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	// cron调度器实例
	cronScheduler *cron.Cron
	// enabled 定时任务开关，只存在内存中，默认打开。进程重启后回到打开。
	enabled atomic.Bool
	mu      sync.Mutex
)

func init() {
	enabled.Store(true)
}

// InitCron 初始化定时任务调度器
func InitCron() {
	cronScheduler = cron.New(cron.WithSeconds())

	// 豆瓣同步：每天 14:30、20:30
	doubanSyncService := service.NewDoubanSyncService()
	_, err := cronScheduler.AddFunc("0 30 14,20 * * *", func() {
		if !Enabled() {
			zap.L().Info("定时任务已关闭，跳过本次豆瓣同步")
			return
		}
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
		zap.L().Info("豆瓣同步定时任务已添加", zap.String("schedule", "每天14:30、20:30执行"))
	}

	cronScheduler.Start()
	zap.L().Info("定时任务调度器已启动")
}

// Enabled 返回定时任务是否打开。状态只在内存中。
func Enabled() bool {
	return enabled.Load()
}

// SetEnabled 打开或关闭定时任务，不写入数据库。
// 关闭后不再触发后续执行；正在跑的那一次会继续跑完。
func SetEnabled(on bool) {
	mu.Lock()
	defer mu.Unlock()

	enabled.Store(on)
	if cronScheduler == nil {
		return
	}
	if on {
		cronScheduler.Start()
		zap.L().Info("定时任务已打开")
		return
	}
	// 不等待正在执行的任务结束，避免管理接口被长时间堵住
	cronScheduler.Stop()
	zap.L().Info("定时任务已关闭")
}

// Stop 停止定时任务调度器
func Stop() {
	if cronScheduler != nil {
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
