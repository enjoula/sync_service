package admin

import (
	"strconv"
	"strings"

	"video-service/internal/model"
	"video-service/internal/pkg/errors"
	"video-service/internal/pkg/response"
	"video-service/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OpLogHandler 管理员操作日志查询
type OpLogHandler struct {
	opLogRepo repository.OpLogRepository
}

// NewOpLogHandler 创建操作日志处理器
func NewOpLogHandler() *OpLogHandler {
	return &OpLogHandler{opLogRepo: repository.NewOpLogRepository()}
}

// OpLogListRequest 列表查询参数
type OpLogListRequest struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	AdminID    string `form:"admin_id"`
	Action     string `form:"action"`
	ObjectType string `form:"object_type"`
	Result     string `form:"result"`
	Keyword    string `form:"keyword"`
}

// OpLogResponse 日志响应
type OpLogResponse struct {
	ID            int64  `json:"id"`
	AdminID       int64  `json:"admin_id"`
	AdminUsername string `json:"admin_username"`
	Action        string `json:"action"`
	ActionLabel   string `json:"action_label"`
	ObjectType    string `json:"object_type"`
	ObjectLabel   string `json:"object_label"`
	ObjectID      string `json:"object_id"`
	ObjectName    string `json:"object_name"`
	BeforeData    string `json:"before_data"`
	AfterData     string `json:"after_data"`
	Result        string `json:"result"`
	ResultLabel   string `json:"result_label"`
	ErrorMsg      string `json:"error_msg,omitempty"`
	ClientIP      string `json:"client_ip,omitempty"`
	CreatedAt     string `json:"created_at"`
}

var actionLabels = map[string]string{
	"create": "新增",
	"update": "编辑",
	"delete": "删除",
}

var objectLabels = map[string]string{
	"video":     "视频",
	"episode":   "剧集",
	"account":   "管理员账号",
	"feedback":  "用户反馈",
	"scheduler": "定时任务",
	"password":  "登录密码",
}

var resultLabels = map[string]string{
	"success": "成功",
	"failure": "失败",
}

// List 分页查询操作日志
func (h *OpLogHandler) List(c *gin.Context) {
	var req OpLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errors.CodeBadRequest, "请求参数无效")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	filter := repository.OpLogListFilter{
		Action:     strings.TrimSpace(req.Action),
		ObjectType: strings.TrimSpace(req.ObjectType),
		Result:     strings.TrimSpace(req.Result),
		Keyword:    strings.TrimSpace(req.Keyword),
		Page:       req.Page,
		PageSize:   req.PageSize,
	}
	if adminIDStr := strings.TrimSpace(req.AdminID); adminIDStr != "" {
		id, err := strconv.ParseInt(adminIDStr, 10, 64)
		if err != nil || id <= 0 {
			response.Error(c, errors.CodeBadRequest, "管理员ID格式错误")
			return
		}
		filter.AdminID = &id
	}

	list, total, err := h.opLogRepo.List(filter)
	if err != nil {
		response.Error(c, errors.CodeInternalErr, "查询操作日志失败")
		return
	}

	items := make([]OpLogResponse, 0, len(list))
	for i := range list {
		items = append(items, toOpLogResponse(&list[i]))
	}
	response.Success(c, gin.H{
		"list":      items,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// Get 获取操作日志详情
func (h *OpLogHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, errors.CodeBadRequest, "日志ID格式错误")
		return
	}
	row, err := h.opLogRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, errors.CodeBadRequest, "日志不存在")
			return
		}
		response.Error(c, errors.CodeInternalErr, "查询操作日志失败")
		return
	}
	response.Success(c, toOpLogResponse(row))
}

func toOpLogResponse(row *model.AdminOperationLog) OpLogResponse {
	return OpLogResponse{
		ID:            row.ID,
		AdminID:       row.AdminID,
		AdminUsername: row.AdminUsername,
		Action:        row.Action,
		ActionLabel:   actionLabels[row.Action],
		ObjectType:    row.ObjectType,
		ObjectLabel:   objectLabels[row.ObjectType],
		ObjectID:      row.ObjectID,
		ObjectName:    row.ObjectName,
		BeforeData:    row.BeforeData,
		AfterData:     row.AfterData,
		Result:        row.Result,
		ResultLabel:   resultLabels[row.Result],
		ErrorMsg:      row.ErrorMsg,
		ClientIP:      row.ClientIP,
		CreatedAt:     row.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
