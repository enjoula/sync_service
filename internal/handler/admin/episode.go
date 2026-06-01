// admin 包提供管理后台的处理器
// 处理剧集相关的HTTP请求
package admin

import (
	"net/http"
	"strconv"

	"video-service/internal/model"
	"video-service/internal/pkg/response"
	"video-service/internal/repository"

	"github.com/gin-gonic/gin"
)

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

// ListByVideoID 根据视频ID获取剧集列表
func (h *EpisodeHandler) ListByVideoID(c *gin.Context) {
	videoIDStr := c.Query("video_id")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的视频ID")
		return
	}

	// 直接查询剧集列表，不检查视频是否存在

	// 获取剧集列表
	episodes, err := h.episodeRepo.FindByVideoID(videoID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取剧集列表失败")
		return
	}

	// 构建响应
	var list []EpisodeResponse
	for _, episode := range episodes {
		resp := EpisodeResponse{
			ID:              strconv.FormatInt(episode.ID, 10),
			Channel:         episode.Channel,
			ChannelID:       strconv.FormatInt(episode.ChannelID, 10),
			VideoID:         strconv.FormatInt(episode.VideoID, 10),
			EpisodeNumber:   episode.EpisodeNumber,
			Name:            episode.Name,
			PlayURLs:        episode.PlayURLs,
			DurationSeconds: episode.DurationSeconds,
			CreatedAt:       episode.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       episode.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		list = append(list, resp)
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

	resp := EpisodeResponse{
		ID:              strconv.FormatInt(episode.ID, 10),
		Channel:         episode.Channel,
		ChannelID:       strconv.FormatInt(episode.ChannelID, 10),
		VideoID:         strconv.FormatInt(episode.VideoID, 10),
		EpisodeNumber:   episode.EpisodeNumber,
		Name:            episode.Name,
		PlayURLs:        episode.PlayURLs,
		DurationSeconds: episode.DurationSeconds,
		CreatedAt:       episode.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       episode.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	response.Success(c, resp)
}

// CreateRequest 创建剧集请求
type CreateEpisodeRequest struct {
	VideoID         int64  `json:"video_id" binding:"required"`
	Channel         string `json:"channel"`
	ChannelID       int64  `json:"channel_id"`
	EpisodeNumber   int64  `json:"episode_number" binding:"required"`
	Name            string `json:"name"`
	PlayURLs        string `json:"play_urls" binding:"required"`
	DurationSeconds int64  `json:"duration_seconds"`
	SubtitleURLs    string `json:"subtitle_urls"`
}

// Create 创建剧集
func (h *EpisodeHandler) Create(c *gin.Context) {
	var req CreateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效")
		return
	}

	videoID := req.VideoID

	// 检查视频是否存在
	_, err := h.videoRepo.FindByID(videoID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "视频不存在")
		return
	}

	// 构建剧集模型
	episode := &model.Episode{
		Channel:         req.Channel,
		ChannelID:       req.ChannelID,
		VideoID:         videoID,
		EpisodeNumber:   req.EpisodeNumber,
		Name:            req.Name,
		PlayURLs:        req.PlayURLs,
		DurationSeconds: req.DurationSeconds,
	}

	// 保存剧集
	if err := h.episodeRepo.Create(episode); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建剧集失败")
		return
	}

	response.Success(c, gin.H{"id": strconv.FormatInt(episode.ID, 10)})
}

// UpdateRequest 更新剧集请求
type UpdateEpisodeRequest struct {
	Channel         string `json:"channel"`
	ChannelID       int64  `json:"channel_id"`
	EpisodeNumber   int64  `json:"episode_number" binding:"required"`
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

	// 检查剧集是否存在
	existingEpisode, err := h.episodeRepo.FindByID(episodeID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "剧集不存在")
		return
	}

	var req UpdateEpisodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效")
		return
	}

	// 更新字段
	existingEpisode.Channel = req.Channel
	existingEpisode.ChannelID = req.ChannelID
	existingEpisode.EpisodeNumber = req.EpisodeNumber
	existingEpisode.Name = req.Name
	existingEpisode.PlayURLs = req.PlayURLs
	existingEpisode.DurationSeconds = req.DurationSeconds

	// 保存更新
	if err := h.episodeRepo.Update(existingEpisode); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新剧集失败")
		return
	}

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

	// 检查剧集是否存在
	_, err = h.episodeRepo.FindByID(episodeID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "剧集不存在")
		return
	}

	// 删除剧集
	if err := h.episodeRepo.DeleteByID(episodeID); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除剧集失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
