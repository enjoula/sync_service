package admin

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"video-service/internal/middleware"
	"video-service/internal/model"
	"video-service/internal/pkg/errors"
	"video-service/internal/pkg/response"
	"video-service/internal/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var adminUsernamePattern = regexp.MustCompile(`^[A-Za-z0-9_]{4,32}$`)

// AccountHandler 管理员账号管理，仅超级管理员可调用
type AccountHandler struct {
	adminRepo repository.AdminRepository
}

// NewAccountHandler 创建管理员账号处理器
func NewAccountHandler() *AccountHandler {
	return &AccountHandler{adminRepo: repository.NewAdminRepository()}
}

// AccountListRequest 管理员列表查询参数
type AccountListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
}

// CreateAccountRequest 新建管理员
type CreateAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Role     string `json:"role" binding:"required"`
}

// UpdateAccountRequest 更新管理员。指针字段用于区分“未传”和“传了空值”。
type UpdateAccountRequest struct {
	Nickname *string `json:"nickname"`
	Role     *string `json:"role"`
	Status   *int    `json:"status"`
	Password *string `json:"password"`
}

// AccountResponse 管理员列表项
type AccountResponse struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Role        string `json:"role"`
	RoleName    string `json:"role_name"`
	Status      int    `json:"status"`
	LastLoginAt string `json:"last_login_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// List 分页查询管理员
func (h *AccountHandler) List(c *gin.Context) {
	var req AccountListRequest
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

	admins, total, err := h.adminRepo.List(req.Keyword, req.PageSize, (req.Page-1)*req.PageSize)
	if err != nil {
		response.Error(c, errors.CodeInternalErr, "获取管理员列表失败")
		return
	}

	list := make([]AccountResponse, 0, len(admins))
	for i := range admins {
		list = append(list, toAccountResponse(&admins[i]))
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// Create 新建管理员
func (h *AccountHandler) Create(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeBadRequest, "请填写用户名、密码和角色")
		return
	}
	username := strings.TrimSpace(req.Username)
	if !adminUsernamePattern.MatchString(username) {
		response.Error(c, errors.CodeBadRequest, "用户名只能包含4-32位字母、数字或下划线")
		return
	}
	if msg := validateAdminPassword(req.Password); msg != "" {
		response.Error(c, errors.CodeBadRequest, msg)
		return
	}
	if !model.IsValidAdminRole(req.Role) {
		response.Error(c, errors.CodeBadRequest, "角色无效")
		return
	}
	nickname := strings.TrimSpace(req.Nickname)
	if utf8.RuneCountInString(nickname) > 100 {
		response.Error(c, errors.CodeBadRequest, "昵称不能超过100个字符")
		return
	}

	if _, err := h.adminRepo.FindByUsername(username); err == nil {
		response.Error(c, errors.CodeConflict, "用户名已存在")
		return
	} else if err != gorm.ErrRecordNotFound {
		response.Error(c, errors.CodeInternalErr, "查询管理员失败")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, errors.CodeInternalErr, "密码加密失败")
		return
	}
	adminUser := &model.Admin{
		Username: username,
		Password: string(hashed),
		Nickname: nickname,
		Role:     req.Role,
		Status:   1,
	}
	if err := h.adminRepo.Create(adminUser); err != nil {
		response.Error(c, errors.CodeInternalErr, "创建管理员失败")
		return
	}
	response.Success(c, toAccountResponse(adminUser))
}

// Update 更新昵称、角色、状态或重置密码
func (h *AccountHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	current, err := h.adminRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, errors.CodeBadRequest, "管理员不存在")
			return
		}
		response.Error(c, errors.CodeInternalErr, "查询管理员失败")
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeBadRequest, "请求参数无效")
		return
	}

	selfID, _ := middleware.CurrentAdminID(c)
	nextRole := current.Role
	nextStatus := current.Status
	fields := map[string]interface{}{}

	if req.Nickname != nil {
		nickname := strings.TrimSpace(*req.Nickname)
		if utf8.RuneCountInString(nickname) > 100 {
			response.Error(c, errors.CodeBadRequest, "昵称不能超过100个字符")
			return
		}
		fields["nickname"] = nickname
	}
	if req.Role != nil {
		if !model.IsValidAdminRole(*req.Role) {
			response.Error(c, errors.CodeBadRequest, "角色无效")
			return
		}
		nextRole = *req.Role
		fields["role"] = nextRole
	}
	if req.Status != nil {
		if *req.Status != 0 && *req.Status != 1 {
			response.Error(c, errors.CodeBadRequest, "状态无效")
			return
		}
		nextStatus = *req.Status
		fields["status"] = nextStatus
	}
	if req.Password != nil && *req.Password != "" {
		if msg := validateAdminPassword(*req.Password); msg != "" {
			response.Error(c, errors.CodeBadRequest, msg)
			return
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Error(c, errors.CodeInternalErr, "密码加密失败")
			return
		}
		fields["password"] = string(hashed)
	}

	if id == selfID && (nextRole != current.Role || nextStatus != current.Status) {
		response.Error(c, errors.CodeBadRequest, "不能修改自己的角色或状态")
		return
	}
	if err := h.ensureEnabledSuper(current, nextRole, nextStatus); err != nil {
		if err == errLastSuperAdmin {
			response.Error(c, errors.CodeBadRequest, "至少保留一名启用的超级管理员")
			return
		}
		response.Error(c, errors.CodeInternalErr, "更新管理员失败")
		return
	}
	if len(fields) == 0 {
		response.Success(c, toAccountResponse(current))
		return
	}
	if err := h.adminRepo.Update(id, fields); err != nil {
		response.Error(c, errors.CodeInternalErr, "更新管理员失败")
		return
	}

	updated, err := h.adminRepo.FindByID(id)
	if err != nil {
		response.Error(c, errors.CodeInternalErr, "查询管理员失败")
		return
	}
	response.Success(c, toAccountResponse(updated))
}

func (h *AccountHandler) ensureEnabledSuper(current *model.Admin, nextRole string, nextStatus int) error {
	if current.Role != model.AdminRoleSuper || current.Status != 1 {
		return nil
	}
	if nextRole == model.AdminRoleSuper && nextStatus == 1 {
		return nil
	}
	count, err := h.adminRepo.CountEnabledByRole(model.AdminRoleSuper, current.ID)
	if err != nil {
		return err
	}
	if count == 0 {
		return errLastSuperAdmin
	}
	return nil
}

func toAccountResponse(adminUser *model.Admin) AccountResponse {
	item := AccountResponse{
		ID:        adminUser.ID,
		Username:  adminUser.Username,
		Nickname:  adminUser.Nickname,
		Role:      adminUser.Role,
		RoleName:  model.AdminRoleNames[adminUser.Role],
		Status:    adminUser.Status,
		CreatedAt: adminUser.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: adminUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if adminUser.LastLoginAt != nil {
		item.LastLoginAt = adminUser.LastLoginAt.Format("2006-01-02 15:04:05")
	}
	return item
}

func parseIDParam(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, errors.CodeBadRequest, "管理员ID无效")
		return 0, false
	}
	return id, true
}

type sentinelError string

func (e sentinelError) Error() string { return string(e) }

const errLastSuperAdmin sentinelError = "last super admin"
