package admin

import (
	"net/http"
	"strconv"
	"time"

	"video-service/internal/pkg/response"
	"video-service/internal/repository"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户管理处理器
type UserHandler struct {
	userRepo repository.UserRepository
}

// NewUserHandler 创建用户管理处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{userRepo: repository.NewUserRepository()}
}

// UserListRequest 用户列表查询参数
type UserListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// UserResponse 用户基本信息。不包含密码和四个访问码字段。
type UserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	WatchCount  int64  `json:"watch_count"`
	LastWatchAt string `json:"last_watch_at"`
	LastOpenAt  string `json:"last_open_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// List 获取用户列表
func (h *UserHandler) List(c *gin.Context) {
	var req UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	users, total, err := h.userRepo.ListBasic(req.Keyword, req.PageSize, (req.Page-1)*req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	list := make([]UserResponse, 0, len(users))
	for _, user := range users {
		item := UserResponse{
			ID:         strconv.FormatInt(user.ID, 10),
			Username:   user.Username,
			Nickname:   user.Nickname,
			Email:      user.Email,
			Avatar:     user.Avatar,
			WatchCount: user.WatchCount,
			CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  user.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if user.LastWatchAt != nil {
			item.LastWatchAt = user.LastWatchAt.Format("2006-01-02 15:04:05")
		}
		if user.LastOpenAt != nil {
			item.LastOpenAt = user.LastOpenAt.Format("2006-01-02 15:04:05")
		}
		list = append(list, item)
	}

	response.Success(c, gin.H{
		"total": total,
		"list":  list,
	})
}

// WatchRecordResponse 观看记录。用户 ID 和视频 ID 用字符串，避免前端精度丢失。
type WatchRecordResponse struct {
	ID             int64  `json:"id"`
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	Nickname       string `json:"nickname"`
	VideoID        string `json:"video_id"`
	VideoTitle     string `json:"video_title"`
	EpisodeID      int64  `json:"episode_id"`
	EpisodeName    string `json:"episode_name"`
	EpisodeNumber  int64  `json:"episode_number"`
	PlayCount      int64  `json:"play_count"`
	LastPositionMs int64  `json:"last_position_ms"`
	LastPlayedAt   string `json:"last_played_at"`
	CreatedAt      string `json:"created_at"`
}

// ListWatchRecords 获取用户观看记录
func (h *UserHandler) ListWatchRecords(c *gin.Context) {
	var req UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	rows, total, err := h.userRepo.ListWatchRecords(req.Keyword, req.PageSize, (req.Page-1)*req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取观看记录失败")
		return
	}

	list := make([]WatchRecordResponse, 0, len(rows))
	for _, row := range rows {
		item := WatchRecordResponse{
			ID:             row.ID,
			UserID:         strconv.FormatInt(row.UserID, 10),
			Username:       row.Username,
			Nickname:       row.Nickname,
			VideoID:        strconv.FormatInt(row.VideoID, 10),
			VideoTitle:     row.VideoTitle,
			EpisodeID:      row.EpisodeID,
			EpisodeName:    row.EpisodeName,
			EpisodeNumber:  row.EpisodeNumber,
			PlayCount:      row.PlayCount,
			LastPositionMs: row.LastPositionMs,
			CreatedAt:      formatAdminTime(row.CreatedAt),
		}
		if !row.LastPlayedAt.IsZero() {
			item.LastPlayedAt = row.LastPlayedAt.Format("2006-01-02 15:04:05")
		}
		list = append(list, item)
	}

	response.Success(c, gin.H{
		"total": total,
		"list":  list,
	})
}

// LoginResponse 登录信息。不返回登录令牌。
type LoginResponse struct {
	ID        int64  `json:"id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Device    string `json:"device"`
	IPAddress string `json:"ip_address"`
	ExpiresAt string `json:"expires_at"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// ListLogins 获取用户登录信息
func (h *UserHandler) ListLogins(c *gin.Context) {
	var req UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	rows, total, err := h.userRepo.ListLogins(req.Keyword, req.PageSize, (req.Page-1)*req.PageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取登录信息失败")
		return
	}

	now := time.Now()
	list := make([]LoginResponse, 0, len(rows))
	for _, row := range rows {
		item := LoginResponse{
			ID:        row.ID,
			UserID:    strconv.FormatInt(row.UserID, 10),
			Username:  row.Username,
			Nickname:  row.Nickname,
			Device:    row.Device,
			IPAddress: row.IPAddress,
			Status:    loginStatus(row.IsActive, row.ExpiresAt, now),
			CreatedAt: formatAdminTime(row.CreatedAt),
		}
		if row.ExpiresAt != nil {
			item.ExpiresAt = row.ExpiresAt.Format("2006-01-02 15:04:05")
		}
		list = append(list, item)
	}

	response.Success(c, gin.H{
		"total": total,
		"list":  list,
	})
}

func formatAdminTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func loginStatus(active bool, expiresAt *time.Time, now time.Time) string {
	if !active {
		return "inactive"
	}
	if expiresAt != nil && expiresAt.Before(now) {
		return "expired"
	}
	return "active"
}
