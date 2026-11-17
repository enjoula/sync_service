// admin 包提供管理后台的处理器
// 处理剧集相关的HTTP请求
package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"video-service/internal/model"
	"video-service/internal/pkg/oplog"
	"video-service/internal/pkg/response"
	"video-service/internal/pkg/utils"
	"video-service/internal/repository"

	"github.com/gin-gonic/gin"
)

const defaultEpisodeChannel = "电影天堂资源"

// EpisodeHandler 剧集处理器
type EpisodeHandler struct {
	episodeRepo repository.EpisodeRepository
	videoRepo   repository.VideoRepository
}

// NewEpisodeHandler 创建剧集处理器实例
func NewEpisodeHandler() *EpisodeHandler {
	return &EpisodeHandler{
		episodeRepo: repository.NewEpisodeRepository(),
		videoRepo:   repository.NewVideoRepository(),
	}
}

// EpisodeResponse 剧集响应结构
type EpisodeResponse struct {
	ID              string      `json:"id"`
	Channel         string      `json:"channel"`
	ChannelID       string      `json:"channel_id"`
	VideoID         string      `json:"video_id"`
	EpisodeNumber   int64       `json:"episode_number"`
	Name            string      `json:"name"`
	PlayURLs        string      `json:"play_urls"`
	DurationSeconds int64       `json:"duration_seconds"`
	SubtitleURLs    interface{} `json:"subtitle_urls"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
}

func toEpisodeResponse(episode *model.Episode) EpisodeResponse {
	return EpisodeResponse{
		ID:              strconv.FormatInt(episode.ID, 10),
		Channel:         episode.Channel,
		ChannelID:       strconv.FormatInt(episode.ChannelID, 10),
		VideoID:         strconv.FormatInt(episode.VideoID, 10),
		EpisodeNumber:   episode.EpisodeNumber,
		Name:            episode.Name,
		PlayURLs:        utils.MaskPlayURL(episode.PlayURLs),
		DurationSeconds: episode.DurationSeconds,
		CreatedAt:       episode.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       episode.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func episodeSnapshot(episode *model.Episode) gin.H {
	if episode == nil {
		return nil
	}
	return gin.H{
		"id":               episode.ID,
		"video_id":         episode.VideoID,
		"channel":          episode.Channel,
		"channel_id":       episode.ChannelID,
		"episode_number":   episode.EpisodeNumber,
		"name":             episode.Name,
		"play_urls":        utils.MaskPlayURL(episode.PlayURLs),
		"duration_seconds": episode.DurationSeconds,
	}
}

// ListByVideoID 根据视频ID获取剧集列表
func (h *EpisodeHandler) ListByVideoID(c *gin.Context) {
	videoIDStr := c.Query("video_id")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的视频ID")
		return
	}

	episodes, err := h.episodeRepo.FindByVideoID(videoID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取剧集列表失败")
		return
	}

	list := make([]EpisodeResponse, 0, len(episodes))
	for _, episode := range episodes {
		list = append(list, toEpisodeResponse(episode))
	}

	response.Success(c, list)
}

// Get 获取单个剧集详情
func (h *EpisodeHandler) Get(c *gin.Context) {
	episodeIDStr := c.Param("id")
	episodeID, err := strconv.ParseInt(episodeIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的剧集ID")
		return
	}

	episode, err := h.episodeRepo.FindByID(episodeID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "剧集不存在")
		return
	}

	response.Success(c, toEpisodeResponse(episode))
}

// CreateRequest 创建剧集请求。video_id 兼容数字与字符串，避免前端雪花 ID 精度/类型问题。
// episode_number 允许 0（电影），故不用 binding:"required"。
type CreateEpisodeRequest struct {
	VideoID         interface{} `json:"video_id" binding:"required"`
	Channel         string      `json:"channel"`
	ChannelID       int64       `json:"channel_id"`
	EpisodeNumber   *int64      `json:"episode_number"`
	Name            string      `json:"name"`
	PlayURLs        string      `json:"play_urls" binding:"required"`
	DurationSeconds int64       `json:"duration_seconds"`
	SubtitleURLs    string      `json:"subtitle_urls"`
}

func parseFlexibleInt64(v interface{}) (int64, error) {
	switch n := v.(type) {
	case float64:
		return int64(n), nil
	case int64:
		return n, nil
	case int:
		return int64(n), nil
	case json.Number:
		return n.Int64()
	case string:
		return strconv.ParseInt(n, 10, 64)
	default:
		return 0, strconv.ErrSyntax
	}
}

// Create 创建剧集
func (h *EpisodeHandler) Create(c *gin.Context) {
	var req CreateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效")
		return
	}
	if req.EpisodeNumber == nil {
		response.Error(c, http.StatusBadRequest, "请填写集数")
		return
	}
	if *req.EpisodeNumber < 0 {
		response.Error(c, http.StatusBadRequest, "集数无效")
		return
	}

	videoID, err := parseFlexibleInt64(req.VideoID)
	if err != nil || videoID <= 0 {
		response.Error(c, http.StatusBadRequest, "无效的视频ID")
		return
	}

	video, err := h.videoRepo.FindByID(videoID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "视频不存在")
		return
	}

	channel := strings.TrimSpace(req.Channel)
	if channel == "" {
		channel = defaultEpisodeChannel
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = video.Title
	}

	episode := &model.Episode{
		Channel:         channel,
		ChannelID:       req.ChannelID,
		VideoID:         videoID,
		EpisodeNumber:   *req.EpisodeNumber,
		Name:            name,
		PlayURLs:        strings.TrimSpace(req.PlayURLs),
		DurationSeconds: req.DurationSeconds,
	}

	if err := h.episodeRepo.Create(episode); err != nil {
		oplog.Record(c, oplog.Entry{
			Action:     oplog.ActionCreate,
			ObjectType: oplog.ObjectEpisode,
			ObjectName: name,
			After:      episodeSnapshot(episode),
			Success:    false,
			ErrorMsg:   "创建剧集失败",
		})
		response.Error(c, http.StatusInternalServerError, "创建剧集失败")
		return
	}

	oplog.Record(c, oplog.Entry{
		Action:     oplog.ActionCreate,
		ObjectType: oplog.ObjectEpisode,
		ObjectID:   strconv.FormatInt(episode.ID, 10),
		ObjectName: episode.Name,
		After:      episodeSnapshot(episode),
		Success:    true,
	})
	response.Success(c, gin.H{"id": strconv.FormatInt(episode.ID, 10)})
}

// UpdateRequest 更新剧集请求
type UpdateEpisodeRequest struct {
	Channel         string `json:"channel"`
	ChannelID       int64  `json:"channel_id"`
	EpisodeNumber   *int64 `json:"episode_number"`
	Name            string `json:"name"`
	PlayURLs        string `json:"play_urls" binding:"required"`
	DurationSeconds int64  `json:"duration_seconds"`
	SubtitleURLs    string `json:"subtitle_urls"`
}

// Update 更新剧集
func (h *EpisodeHandler) Update(c *gin.Context) {
	episodeIDStr := c.Param("id")
	episodeID, err := strconv.ParseInt(episodeIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的剧集ID")
		return
	}

	existingEpisode, err := h.episodeRepo.FindByID(episodeID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "剧集不存在")
		return
	}
	before := episodeSnapshot(existingEpisode)

	var req UpdateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效")
		return
	}
	if req.EpisodeNumber == nil {
		response.Error(c, http.StatusBadRequest, "请填写集数")
		return
	}
	if *req.EpisodeNumber < 0 {
		response.Error(c, http.StatusBadRequest, "集数无效")
		return
	}

	channel := strings.TrimSpace(req.Channel)
	if channel == "" {
		channel = existingEpisode.Channel
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = existingEpisode.Name
	}

	playURLs := strings.TrimSpace(req.PlayURLs)
	// 接口返回已脱敏；编辑回写脱敏值时保留原地址，避免把 **** 写进库
	if utils.IsMaskedPlayURL(playURLs) || playURLs == utils.MaskPlayURL(existingEpisode.PlayURLs) {
		playURLs = existingEpisode.PlayURLs
	}

	existingEpisode.Channel = channel
	existingEpisode.ChannelID = req.ChannelID
	existingEpisode.EpisodeNumber = *req.EpisodeNumber
	existingEpisode.Name = name
	existingEpisode.PlayURLs = playURLs
	existingEpisode.DurationSeconds = req.DurationSeconds

	if err := h.episodeRepo.Update(existingEpisode); err != nil {
		oplog.Record(c, oplog.Entry{
			Action:     oplog.ActionUpdate,
			ObjectType: oplog.ObjectEpisode,
			ObjectID:   strconv.FormatInt(existingEpisode.ID, 10),
			ObjectName: existingEpisode.Name,
			Before:     before,
			After:      episodeSnapshot(existingEpisode),
			Success:    false,
			ErrorMsg:   "更新剧集失败",
		})
		response.Error(c, http.StatusInternalServerError, "更新剧集失败")
		return
	}

	oplog.Record(c, oplog.Entry{
		Action:     oplog.ActionUpdate,
		ObjectType: oplog.ObjectEpisode,
		ObjectID:   strconv.FormatInt(existingEpisode.ID, 10),
		ObjectName: existingEpisode.Name,
		Before:     before,
		After:      episodeSnapshot(existingEpisode),
		Success:    true,
	})
	response.Success(c, gin.H{"id": strconv.FormatInt(existingEpisode.ID, 10)})
}

// Delete 删除剧集
func (h *EpisodeHandler) Delete(c *gin.Context) {
	episodeIDStr := c.Param("id")
	episodeID, err := strconv.ParseInt(episodeIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的剧集ID")
		return
	}

	existing, err := h.episodeRepo.FindByID(episodeID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "剧集不存在")
		return
	}
	before := episodeSnapshot(existing)

	if err := h.episodeRepo.DeleteByID(episodeID); err != nil {
		oplog.Record(c, oplog.Entry{
			Action:     oplog.ActionDelete,
			ObjectType: oplog.ObjectEpisode,
			ObjectID:   strconv.FormatInt(episodeID, 10),
			ObjectName: existing.Name,
			Before:     before,
			Success:    false,
			ErrorMsg:   "删除剧集失败",
		})
		response.Error(c, http.StatusInternalServerError, "删除剧集失败")
		return
	}

	oplog.Record(c, oplog.Entry{
		Action:     oplog.ActionDelete,
		ObjectType: oplog.ObjectEpisode,
		ObjectID:   strconv.FormatInt(episodeID, 10),
		ObjectName: existing.Name,
		Before:     before,
		Success:    true,
	})
	response.Success(c, gin.H{"message": "删除成功"})
}
