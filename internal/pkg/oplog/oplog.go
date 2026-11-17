// Package oplog 记录管理后台增删改操作日志，供追溯与查询。
package oplog

import (
	"encoding/json"
	"strings"
	"time"

	"video-service/internal/middleware"
	"video-service/internal/model"
	"video-service/pkg/infrastructure/database"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"

	ObjectVideo     = "video"
	ObjectEpisode   = "episode"
	ObjectAccount   = "account"
	ObjectFeedback  = "feedback"
	ObjectScheduler = "scheduler"
	ObjectPassword  = "password"

	ResultSuccess = "success"
	ResultFailure = "failure"
)

// Entry 一条操作日志
type Entry struct {
	Action     string
	ObjectType string
	ObjectID   string
	ObjectName string
	Before     interface{}
	After      interface{}
	Success    bool
	ErrorMsg   string
}

// Record 异步落库操作日志；失败只打 warn，不影响主流程。
func Record(c *gin.Context, entry Entry) {
	if c == nil {
		return
	}
	adminID, _ := middleware.CurrentAdminID(c)
	adminName, _ := middleware.CurrentAdminUsername(c)
	result := ResultSuccess
	if !entry.Success {
		result = ResultFailure
	}

	logRow := &model.AdminOperationLog{
		AdminID:       adminID,
		AdminUsername: adminName,
		Action:        entry.Action,
		ObjectType:    entry.ObjectType,
		ObjectID:      strings.TrimSpace(entry.ObjectID),
		ObjectName:    strings.TrimSpace(entry.ObjectName),
		BeforeData:    toJSONString(entry.Before),
		AfterData:     toJSONString(entry.After),
		Result:        result,
		ErrorMsg:      strings.TrimSpace(entry.ErrorMsg),
		ClientIP:      c.ClientIP(),
		CreatedAt:     time.Now(),
	}

	go func(row *model.AdminOperationLog) {
		if database.DB == nil {
			return
		}
		if err := database.DB.Create(row).Error; err != nil {
			zap.L().Warn("写入管理员操作日志失败",
				zap.String("action", row.Action),
				zap.String("object_type", row.ObjectType),
				zap.String("object_id", row.ObjectID),
				zap.Error(err),
			)
		}
	}(logRow)
}

func toJSONString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
