package admin

import (
	"video-service/internal/pkg/errors"
	"video-service/internal/pkg/response"
	"video-service/pkg/infrastructure/scheduler"

	"github.com/gin-gonic/gin"
)

// SchedulerHandler 定时任务开关。状态只保存在进程内存，不写数据库。
type SchedulerHandler struct{}

// NewSchedulerHandler 创建定时任务开关处理器
func NewSchedulerHandler() *SchedulerHandler {
	return &SchedulerHandler{}
}

// Get 返回当前开关状态
func (h *SchedulerHandler) Get(c *gin.Context) {
	response.Success(c, gin.H{
		"enabled":  scheduler.Enabled(),
		"schedule": "每天14:30、20:30",
	})
}

// UpdateRequest 打开或关闭定时任务
type UpdateSchedulerRequest struct {
	Enabled *bool `json:"enabled" binding:"required"`
}

// Update 由超级管理员打开或关闭定时任务
func (h *SchedulerHandler) Update(c *gin.Context) {
	var req UpdateSchedulerRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Enabled == nil {
		response.Error(c, errors.CodeBadRequest, "请求参数无效")
		return
	}

	scheduler.SetEnabled(*req.Enabled)
	response.Success(c, gin.H{"enabled": scheduler.Enabled()})
}
