package admin

import (
	"strings"
	"time"
	"unicode/utf8"

	"video-service/internal/middleware"
	"video-service/internal/model"
	"video-service/internal/pkg/auth"
	"video-service/internal/pkg/errors"
	"video-service/internal/pkg/response"
	"video-service/internal/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const adminTokenTTL = 12 * time.Hour

// AuthHandler 管理后台登录
type AuthHandler struct {
	adminRepo repository.AdminRepository
}

// NewAuthHandler 创建登录处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{adminRepo: repository.NewAdminRepository()}
}

// LoginRequest 登录参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改当前管理员密码
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// AdminProfile 返回给前端的管理员信息，不包含密码
type AdminProfile struct {
	ID          int64    `json:"id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Role        string   `json:"role"`
	RoleName    string   `json:"role_name"`
	Permissions []string `json:"permissions"`
	Status      int      `json:"status"`
}

func toAdminProfile(admin *model.Admin) AdminProfile {
	return AdminProfile{
		ID:          admin.ID,
		Username:    admin.Username,
		Nickname:    admin.Nickname,
		Role:        admin.Role,
		RoleName:    model.AdminRoleNames[admin.Role],
		Permissions: model.AdminPermissions(admin.Role),
		Status:      admin.Status,
	}
}

// Login 使用管理员表中的用户名和密码登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeBadRequest, "用户名或密码不能为空")
		return
	}
	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		response.Error(c, errors.CodeBadRequest, "用户名或密码不能为空")
		return
	}

	adminUser, err := h.adminRepo.FindByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, errors.CodeUnauthorized, "用户名或密码错误")
			return
		}
		response.Error(c, errors.CodeInternalErr, "登录失败")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(req.Password)) != nil {
		response.Error(c, errors.CodeUnauthorized, "用户名或密码错误")
		return
	}
	if adminUser.Status != 1 {
		response.Error(c, errors.CodeForbidden, "账号已禁用")
		return
	}

	expiresAt := time.Now().Add(adminTokenTTL)
	token, err := auth.GenerateAdminToken(adminUser.ID, adminUser.Username, adminUser.Role, expiresAt)
	if err != nil {
		response.Error(c, errors.CodeInternalErr, "生成登录凭证失败")
		return
	}
	if err := h.adminRepo.UpdateLastLogin(adminUser.ID, time.Now()); err != nil {
		zap.L().Warn("更新管理员登录时间失败", zap.Error(err), zap.Int64("admin_id", adminUser.ID))
	}

	response.Success(c, gin.H{
		"token":      token,
		"expires_at": expiresAt.Format("2006-01-02 15:04:05"),
		"admin":      toAdminProfile(adminUser),
	})
}

// Me 返回当前登录管理员
func (h *AuthHandler) Me(c *gin.Context) {
	adminUser, ok := h.currentAdmin(c)
	if !ok {
		return
	}
	response.Success(c, toAdminProfile(adminUser))
}

// ChangePassword 修改当前管理员密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	adminUser, ok := h.currentAdmin(c)
	if !ok {
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeBadRequest, "请填写原密码和新密码")
		return
	}
	if msg := validateAdminPassword(req.NewPassword); msg != "" {
		response.Error(c, errors.CodeBadRequest, msg)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(req.OldPassword)); err != nil {
		response.Error(c, errors.CodeBadRequest, "原密码错误")
		return
	}
	if req.OldPassword == req.NewPassword {
		response.Error(c, errors.CodeBadRequest, "新密码不能与原密码相同")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, errors.CodeInternalErr, "密码加密失败")
		return
	}
	if err := h.adminRepo.Update(adminUser.ID, map[string]interface{}{"password": string(hashed)}); err != nil {
		response.Error(c, errors.CodeInternalErr, "修改密码失败")
		return
	}
	response.SuccessMsg(c, "密码已更新", nil)
}

func (h *AuthHandler) currentAdmin(c *gin.Context) (*model.Admin, bool) {
	id, ok := middleware.CurrentAdminID(c)
	if !ok {
		response.Error(c, errors.CodeUnauthorized, errors.MsgNotLoggedIn)
		return nil, false
	}
	adminUser, err := h.adminRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, errors.CodeUnauthorized, "账号不存在")
			return nil, false
		}
		response.Error(c, errors.CodeInternalErr, "查询管理员失败")
		return nil, false
	}
	return adminUser, true
}

func validateAdminPassword(password string) string {
	n := utf8.RuneCountInString(password)
	if n < 6 || n > 64 {
		return "密码长度必须在6-64个字符之间"
	}
	return ""
}
