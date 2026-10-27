package admin

import (
	"net/http"
	"strconv"

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
