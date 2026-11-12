package admin

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"video-service/internal/middleware"
	"video-service/internal/model"
	"video-service/internal/pkg/errors"
	"video-service/internal/pkg/response"
	"video-service/internal/repository"
	"video-service/pkg/infrastructure/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const feedbackRemarkMaxRunes = 1000

// FeedbackHandler 用户反馈管理
type FeedbackHandler struct {
	feedbackRepo repository.FeedbackRepository
	adminRepo    repository.AdminRepository
}

// NewFeedbackHandler 创建反馈管理处理器
func NewFeedbackHandler() *FeedbackHandler {
	return &FeedbackHandler{
		feedbackRepo: repository.NewFeedbackRepository(),
		adminRepo:    repository.NewAdminRepository(),
	}
}

// FeedbackListRequest 列表查询参数
type FeedbackListRequest struct {
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
	Source      string `form:"source"`
	Status      string `form:"status"`
	ProblemType string `form:"problem_type"`
	VideoID     string `form:"video_id"`
	UserID      string `form:"user_id"`
	Keyword     string `form:"keyword"`
}

// UpdateFeedbackRequest 更新处理状态
type UpdateFeedbackRequest struct {
	Status      string  `json:"status" binding:"required"`
	AdminRemark *string `json:"admin_remark"`
}

// FeedbackResponse 反馈列表/详情响应
type FeedbackResponse struct {
	ID             int64   `json:"id"`
	UserID         string  `json:"user_id"`
	Username       string  `json:"username,omitempty"`
	Nickname       string  `json:"nickname,omitempty"`
	Source         string  `json:"source"`
	SourceLabel    string  `json:"source_label"`
	VideoID        *string `json:"video_id,omitempty"`
	EpisodeID      *int64  `json:"episode_id,omitempty"`
	EpisodeNumber  *int64  `json:"episode_number,omitempty"`
	VideoTitle     string  `json:"video_title"`
	SearchKeyword  string  `json:"search_keyword,omitempty"`
	ProblemType    string  `json:"problem_type"`
	ProblemLabel   string  `json:"problem_label"`
	Content        string  `json:"content"`
	Contact        string  `json:"contact,omitempty"`
	Status         string  `json:"status"`
	StatusLabel    string  `json:"status_label"`
	AdminRemark    string  `json:"admin_remark,omitempty"`
	HandledBy      *int64  `json:"handled_by,omitempty"`
	HandledByName  string  `json:"handled_by_name,omitempty"`
	HandledAt      string  `json:"handled_at,omitempty"`
	ClientDevice   string  `json:"client_device,omitempty"`
	ClientIP       string  `json:"client_ip,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// List 分页查询用户反馈
func (h *FeedbackHandler) List(c *gin.Context) {
	var req FeedbackListRequest
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

	source := strings.TrimSpace(req.Source)
	if source != "" && !model.IsValidFeedbackSource(source) {
		response.Error(c, errors.CodeBadRequest, "来源方式无效")
		return
	}
	status := strings.TrimSpace(req.Status)
	if status != "" && !model.IsValidFeedbackStatus(status) {
		response.Error(c, errors.CodeBadRequest, "处理状态无效")
		return
	}
	problemType := strings.TrimSpace(req.ProblemType)
	if problemType != "" && !model.IsValidFeedbackType(problemType) {
		response.Error(c, errors.CodeBadRequest, "问题类型无效")
		return
	}

	filter := repository.FeedbackListFilter{
		Source:      source,
		Status:      status,
		ProblemType: problemType,
		Keyword:     strings.TrimSpace(req.Keyword),
		Page:        req.Page,
		PageSize:    req.PageSize,
	}

	if videoIDStr := strings.TrimSpace(req.VideoID); videoIDStr != "" {
		id, err := strconv.ParseInt(videoIDStr, 10, 64)
		if err != nil || id <= 0 {
			response.Error(c, errors.CodeBadRequest, "视频ID格式错误")
			return
		}
		filter.VideoID = &id
	}
	if userIDStr := strings.TrimSpace(req.UserID); userIDStr != "" {
		id, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil || id <= 0 {
			response.Error(c, errors.CodeBadRequest, "用户ID格式错误")
			return
		}
		filter.UserID = &id
	}

	list, total, err := h.feedbackRepo.List(filter)
	if err != nil {
		response.Error(c, errors.CodeInternalErr, "查询反馈列表失败")
		return
	}

	userNames := h.loadUserNames(collectFeedbackUserIDs(list))
	adminNames := h.loadAdminNames(collectHandledByIDs(list))

	items := make([]FeedbackResponse, 0, len(list))
	for i := range list {
		items = append(items, toFeedbackResponse(&list[i], userNames, adminNames))
	}

	response.Success(c, gin.H{
		"list":      items,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// Get 获取反馈详情
func (h *FeedbackHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, errors.CodeBadRequest, "反馈ID格式错误")
		return
	}

	feedback, err := h.feedbackRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, errors.CodeBadRequest, "反馈不存在")
			return
		}
		response.Error(c, errors.CodeInternalErr, "查询反馈失败")
		return
	}

	userNames := h.loadUserNames([]int64{feedback.UserID})
	adminNames := h.loadAdminNames(collectHandledByIDs([]model.VideoFeedback{*feedback}))
	response.Success(c, toFeedbackResponse(feedback, userNames, adminNames))
}

// Update 更新反馈处理状态与备注
func (h *FeedbackHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, errors.CodeBadRequest, "反馈ID格式错误")
		return
	}

	var req UpdateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeBadRequest, "请求参数错误")
		return
	}

	status := strings.TrimSpace(req.Status)
	if !model.IsValidFeedbackStatus(status) {
		response.Error(c, errors.CodeBadRequest, "处理状态无效")
		return
	}

	feedback, err := h.feedbackRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, errors.CodeBadRequest, "反馈不存在")
			return
		}
		response.Error(c, errors.CodeInternalErr, "查询反馈失败")
		return
	}

	feedback.Status = status
	if req.AdminRemark != nil {
		remark := strings.TrimSpace(*req.AdminRemark)
		if utf8.RuneCountInString(remark) > feedbackRemarkMaxRunes {
			response.Error(c, errors.CodeBadRequest, "管理员备注过长")
			return
		}
		feedback.AdminRemark = remark
	}

	if adminID, ok := middleware.CurrentAdminID(c); ok {
		feedback.HandledBy = &adminID
	}
	now := time.Now()
	feedback.HandledAt = &now

	if err := h.feedbackRepo.Update(feedback); err != nil {
		response.Error(c, errors.CodeInternalErr, "更新反馈失败")
		return
	}

	userNames := h.loadUserNames([]int64{feedback.UserID})
	adminNames := h.loadAdminNames(collectHandledByIDs([]model.VideoFeedback{*feedback}))
	response.Success(c, toFeedbackResponse(feedback, userNames, adminNames))
}

func collectFeedbackUserIDs(list []model.VideoFeedback) []int64 {
	seen := make(map[int64]struct{}, len(list))
	ids := make([]int64, 0, len(list))
	for _, item := range list {
		if _, ok := seen[item.UserID]; ok || item.UserID <= 0 {
			continue
		}
		seen[item.UserID] = struct{}{}
		ids = append(ids, item.UserID)
	}
	return ids
}

func collectHandledByIDs(list []model.VideoFeedback) []int64 {
	seen := make(map[int64]struct{}, len(list))
	ids := make([]int64, 0, len(list))
	for _, item := range list {
		if item.HandledBy == nil || *item.HandledBy <= 0 {
			continue
		}
		id := *item.HandledBy
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

func (h *FeedbackHandler) loadUserNames(ids []int64) map[int64][2]string {
	out := make(map[int64][2]string, len(ids))
	if len(ids) == 0 || database.DB == nil {
		return out
	}
	type row struct {
		ID       int64
		Username string
		Nickname string
	}
	var rows []row
	if err := database.DB.Model(&model.User{}).
		Select("id, username, nickname").
		Where("id IN ?", ids).
		Find(&rows).Error; err != nil {
		return out
	}
	for _, item := range rows {
		out[item.ID] = [2]string{item.Username, item.Nickname}
	}
	return out
}

func (h *FeedbackHandler) loadAdminNames(ids []int64) map[int64]string {
	out := make(map[int64]string, len(ids))
	for _, id := range ids {
		admin, err := h.adminRepo.FindByID(id)
		if err != nil {
			continue
		}
		name := strings.TrimSpace(admin.Nickname)
		if name == "" {
			name = admin.Username
		}
		out[id] = name
	}
	return out
}

func toFeedbackResponse(feedback *model.VideoFeedback, userNames map[int64][2]string, adminNames map[int64]string) FeedbackResponse {
	resp := FeedbackResponse{
		ID:            feedback.ID,
		UserID:        strconv.FormatInt(feedback.UserID, 10),
		Source:        feedback.Source,
		SourceLabel:   model.ValidFeedbackSources[feedback.Source],
		EpisodeID:     feedback.EpisodeID,
		EpisodeNumber: feedback.EpisodeNumber,
		VideoTitle:    feedback.VideoTitle,
		SearchKeyword: feedback.SearchKeyword,
		ProblemType:   feedback.ProblemType,
		ProblemLabel:  model.ValidFeedbackTypes[feedback.ProblemType],
		Content:       feedback.Content,
		Contact:       feedback.Contact,
		Status:        feedback.Status,
		StatusLabel:   model.ValidFeedbackStatuses[feedback.Status],
		AdminRemark:   feedback.AdminRemark,
		HandledBy:     feedback.HandledBy,
		ClientDevice:  feedback.ClientDevice,
		ClientIP:      feedback.ClientIP,
		CreatedAt:     feedback.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     feedback.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if feedback.VideoID != nil {
		s := strconv.FormatInt(*feedback.VideoID, 10)
		resp.VideoID = &s
	}
	if names, ok := userNames[feedback.UserID]; ok {
		resp.Username = names[0]
		resp.Nickname = names[1]
	}
	if feedback.HandledBy != nil {
		resp.HandledByName = adminNames[*feedback.HandledBy]
	}
	if feedback.HandledAt != nil {
		resp.HandledAt = feedback.HandledAt.Format("2006-01-02 15:04:05")
	}
	return resp
}
