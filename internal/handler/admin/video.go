// admin 包提供管理后台的处理器
// 处理视频相关的HTTP请求
package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"video-service/internal/model"
	"video-service/internal/pkg/response"
	"video-service/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// VideoHandler 视频处理器
type VideoHandler struct {
	videoRepo   repository.VideoRepository
	episodeRepo repository.EpisodeRepository
}

// NewVideoHandler 创建视频处理器实例
func NewVideoHandler() *VideoHandler {
	return &VideoHandler{
		videoRepo:   repository.NewVideoRepository(),
		episodeRepo: repository.NewEpisodeRepository(),
	}
}

// applyJSONStringField 将管理端传入的 JSON 字符串写入 datatypes.JSON（空串置空，非法 JSON 则写 []）
func applyJSONStringField(dst *datatypes.JSON, raw string) {
	if raw == "" {
		*dst = nil
		return
	}
	if !json.Valid([]byte(raw)) {
		*dst = datatypes.JSON([]byte("[]"))
		return
	}
	*dst = datatypes.JSON([]byte(raw))
}

// VideoListRequest 视频列表请求参数
type VideoListRequest struct {
	Page        int    `form:"page" binding:"min=1"`
	PageSize    int    `form:"page_size" binding:"min=1,max=100"`
	Title       string `form:"title"`
	Type        string `form:"type"`
	Status      string `form:"status"`
	IsCompleted string `form:"is_completed"`
	IsUpdate    string `form:"is_update"`
}

// VideoResponse 视频响应结构
type VideoResponse struct {
	ID           string      `json:"id"`
	SourceID     *int64      `json:"source_id"`
	Source       string      `json:"source"`
	Title        string      `json:"title"`
	Type         string      `json:"type"`
	CoverURL     string      `json:"cover_url"`
	Description  string      `json:"description"`
	ReleaseDate  *string     `json:"release_date"`
	Score        *float64    `json:"score"`
	CountryJSON  interface{} `json:"country_json"`
	DirectorJSON interface{} `json:"director_json"`
	ActorsJSON   interface{} `json:"actors_json"`
	TagsJSON     interface{} `json:"tags_json"`
	Status       string      `json:"status"`
	IMDbID       string      `json:"imdb_id"`
	Runtime      *int64      `json:"runtime"`
	EpisodeCount int64       `json:"episode_count"`
	IsCompleted  bool        `json:"is_completed"`
	IsUpdate     bool        `json:"is_update"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
}

// VideoListResponse 视频列表响应
type VideoListResponse struct {
	Total int64           `json:"total"`
	List  []VideoResponse `json:"list"`
}

// List 获取视频列表
func (h *VideoHandler) List(c *gin.Context) {
	var req VideoListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效")
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20 // 默认每页20条
	}

	// 计算偏移量
	offset := (req.Page - 1) * req.PageSize

	// 处理筛选参数
	var isCompleted *bool
	if req.IsCompleted != "" {
		completed := req.IsCompleted == "1"
		isCompleted = &completed
	}

	var isUpdate *bool
	if req.IsUpdate != "" {
		update := req.IsUpdate == "1"
		isUpdate = &update
	}

	// 获取视频列表（带分页、排序和筛选）
	videos, total, err := h.videoRepo.FindVideosWithPagination(req.PageSize, offset, req.Title, req.Type, req.Status, isCompleted, isUpdate)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取视频列表失败")
		return
	}

	// 构建响应
	var list []VideoResponse
	for _, video := range videos {
		resp := VideoResponse{
			ID:           strconv.FormatInt(video.ID, 10),
			SourceID:     video.SourceID,
			Source:       video.Source,
			Title:        video.Title,
			Type:         video.Type,
			CoverURL:     video.CoverURL,
			Description:  video.Description,
			Status:       video.Status,
			IMDbID:       video.IMDbID,
			Runtime:      video.Runtime,
			EpisodeCount: video.EpisodeCount,
			IsCompleted:  video.IsCompleted,
			IsUpdate:     video.IsUpdate,
			CreatedAt:    video.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    video.UpdatedAt.Format("2006-01-02 15:04:05"),
			Score:        video.Score,
			CountryJSON:  video.CountryJSON,
			DirectorJSON: video.DirectorJSON,
			ActorsJSON:   video.ActorsJSON,
			TagsJSON:     video.TagsJSON,
		}

		// 处理日期
		if video.ReleaseDate != nil {
			dateStr := video.ReleaseDate.Format("2006-01-02")
			resp.ReleaseDate = &dateStr
		}

		list = append(list, resp)
	}

	response.Success(c, VideoListResponse{
		Total: total,
		List:  list,
	})
}

// Get 获取单个视频详情
func (h *VideoHandler) Get(c *gin.Context) {
	videoIDStr := c.Param("id")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的视频ID")
		return
	}

	video, err := h.videoRepo.FindByID(videoID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "视频不存在")
		return
	}

	resp := VideoResponse{
		ID:           strconv.FormatInt(video.ID, 10),
		SourceID:     video.SourceID,
		Source:       video.Source,
		Title:        video.Title,
		Type:         video.Type,
		CoverURL:     video.CoverURL,
		Description:  video.Description,
		Status:       video.Status,
		IMDbID:       video.IMDbID,
		Runtime:      video.Runtime,
		EpisodeCount: video.EpisodeCount,
		IsCompleted:  video.IsCompleted,
		IsUpdate:     video.IsUpdate,
		CreatedAt:    video.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    video.UpdatedAt.Format("2006-01-02 15:04:05"),
		Score:        video.Score,
		CountryJSON:  video.CountryJSON,
		DirectorJSON: video.DirectorJSON,
		ActorsJSON:   video.ActorsJSON,
		TagsJSON:     video.TagsJSON,
	}

	// 处理日期
	if video.ReleaseDate != nil {
		dateStr := video.ReleaseDate.Format("2006-01-02")
		resp.ReleaseDate = &dateStr
	}

	response.Success(c, resp)
}

// CreateRequest 创建视频请求
type CreateRequest struct {
	SourceID     *int64   `json:"source_id"`
	Source       string   `json:"source"`
	Title        string   `json:"title" binding:"required"`
	Type         string   `json:"type" binding:"required"`
	CoverURL     string   `json:"cover_url"`
	Description  string   `json:"description"`
	ReleaseDate  *string  `json:"release_date"`
	Score        *float64 `json:"score"`
	CountryJSON  string   `json:"country_json"`
	DirectorJSON string   `json:"director_json"`
	ActorsJSON   string   `json:"actors_json"`
	TagsJSON     string   `json:"tags_json"`
	Status       string   `json:"status"`
	IMDbID       string   `json:"imdb_id"`
	Runtime      *int64   `json:"runtime"`
	EpisodeCount int64    `json:"episode_count"`
	IsCompleted  bool     `json:"is_completed"`
	IsUpdate     bool     `json:"is_update"`
}

// Create 创建视频
func (h *VideoHandler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效")
		return
	}

	// 构建视频模型
	video := &model.Video{
		SourceID:     req.SourceID,
		Source:       req.Source,
		Title:        req.Title,
		Type:         req.Type,
		CoverURL:     req.CoverURL,
		Description:  req.Description,
		Status:       req.Status,
		IMDbID:       req.IMDbID,
		Runtime:      req.Runtime,
		EpisodeCount: req.EpisodeCount,
		IsCompleted:  req.IsCompleted,
		IsUpdate:     req.IsUpdate,
	}

	// 处理日期
	if req.ReleaseDate != nil {
		releaseDate, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err == nil {
			video.ReleaseDate = &releaseDate
		}
	}

	// 处理JSON字段
	video.Score = req.Score
	applyJSONStringField(&video.CountryJSON, req.CountryJSON)
	applyJSONStringField(&video.DirectorJSON, req.DirectorJSON)
	applyJSONStringField(&video.ActorsJSON, req.ActorsJSON)
	applyJSONStringField(&video.TagsJSON, req.TagsJSON)

	// 保存视频
	if err := h.videoRepo.Create(video); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建视频失败")
		return
	}

	response.Success(c, gin.H{"id": strconv.FormatInt(video.ID, 10)})
}

// UpdateRequest 更新视频请求
type UpdateRequest struct {
	SourceID     *int64   `json:"source_id"`
	Source       string   `json:"source"`
	Title        string   `json:"title" binding:"required"`
	Type         string   `json:"type" binding:"required"`
	CoverURL     string   `json:"cover_url"`
	Description  string   `json:"description"`
	ReleaseDate  *string  `json:"release_date"`
	Score        *float64 `json:"score"`
	CountryJSON  string   `json:"country_json"`
	DirectorJSON string   `json:"director_json"`
	ActorsJSON   string   `json:"actors_json"`
	TagsJSON     string   `json:"tags_json"`
	Status       string   `json:"status"`
	IMDbID       string   `json:"imdb_id"`
	Runtime      *int64   `json:"runtime"`
	EpisodeCount int64    `json:"episode_count"`
	IsCompleted  bool     `json:"is_completed"`
	IsUpdate     bool     `json:"is_update"`
}

// Update 更新视频
func (h *VideoHandler) Update(c *gin.Context) {
	videoIDStr := c.Param("id")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的视频ID")
		return
	}

	// 检查视频是否存在
	existingVideo, err := h.videoRepo.FindByID(videoID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "视频不存在")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效")
		return
	}

	// 更新字段
	existingVideo.SourceID = req.SourceID
	existingVideo.Source = req.Source
	existingVideo.Title = req.Title
	existingVideo.Type = req.Type
	existingVideo.CoverURL = req.CoverURL
	existingVideo.Description = req.Description
	existingVideo.Status = req.Status
	existingVideo.IMDbID = req.IMDbID
	existingVideo.Runtime = req.Runtime
	existingVideo.EpisodeCount = req.EpisodeCount
	existingVideo.IsCompleted = req.IsCompleted
	existingVideo.IsUpdate = req.IsUpdate

	// 处理日期
	if req.ReleaseDate != nil {
		releaseDate, err := time.Parse("2006-01-02", *req.ReleaseDate)
		if err == nil {
			existingVideo.ReleaseDate = &releaseDate
		}
	}

	// 处理JSON字段
	existingVideo.Score = req.Score
	applyJSONStringField(&existingVideo.CountryJSON, req.CountryJSON)
	applyJSONStringField(&existingVideo.DirectorJSON, req.DirectorJSON)
	applyJSONStringField(&existingVideo.ActorsJSON, req.ActorsJSON)
	applyJSONStringField(&existingVideo.TagsJSON, req.TagsJSON)

	// 保存更新
	if err := h.videoRepo.Update(existingVideo); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新视频失败")
		return
	}

	response.Success(c, gin.H{"id": strconv.FormatInt(existingVideo.ID, 10)})
}

// Delete 删除视频
func (h *VideoHandler) Delete(c *gin.Context) {
	videoIDStr := c.Param("id")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的视频ID")
		return
	}

	// 检查视频是否存在
	_, err = h.videoRepo.FindByID(videoID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "视频不存在")
		return
	}

	// 删除视频
	if err := h.videoRepo.DeleteByID(videoID); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除视频失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
