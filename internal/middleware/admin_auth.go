package middleware

import (
	"strings"

	"video-service/internal/model"
	"video-service/internal/pkg/auth"
	"video-service/internal/pkg/errors"
	"video-service/internal/pkg/response"
	"video-service/internal/repository"

	"github.com/gin-gonic/gin"
)

const (
	ctxAdminID          = "admin_id"
	ctxAdminRole        = "admin_role"
	ctxAdminUsername    = "admin_username"
	ctxAdminPermissions = "admin_permissions"
)

// AdminAuth 校验管理后台登录令牌，并确认账号仍然启用。
// 每次请求都以数据库中的角色为准，角色变更和禁用会立即生效。
func AdminAuth() gin.HandlerFunc {
	repo := repository.NewAdminRepository()
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errors.CodeUnauthorized, errors.MsgNotLoggedIn)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, errors.CodeUnauthorized, errors.MsgTokenInvalidFormat)
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil || claims.TokenType != auth.TokenTypeAdmin || claims.AdminID == 0 {
			response.Error(c, errors.CodeUnauthorized, "登录已失效，请重新登录")
			c.Abort()
			return
		}

		admin, err := repo.FindByID(claims.AdminID)
		if err != nil || admin.Status != 1 {
			response.Error(c, errors.CodeUnauthorized, "账号已禁用或不存在")
			c.Abort()
			return
		}

		c.Set(ctxAdminID, admin.ID)
		c.Set(ctxAdminRole, admin.Role)
		c.Set(ctxAdminUsername, admin.Username)
		c.Set(ctxAdminPermissions, model.AdminPermissions(admin.Role))
		c.Next()
	}
}

// RequireSuperAdmin 只允许超级管理员继续。角色以本次请求从数据库读到的为准。
func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ctxAdminRole)
		if role != model.AdminRoleSuper {
			response.Error(c, errors.CodeForbidden, "只有超级管理员可以操作")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAdminPerm 要求当前管理员拥有指定权限
func RequireAdminPerm(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, _ := c.Get(ctxAdminPermissions)
		list, _ := perms.([]string)
		for _, item := range list {
			if item == perm {
				c.Next()
				return
			}
		}
		response.Error(c, errors.CodeForbidden, "没有权限执行该操作")
		c.Abort()
	}
}

// CurrentAdminID 读取当前登录管理员 ID
func CurrentAdminID(c *gin.Context) (int64, bool) {
	value, ok := c.Get(ctxAdminID)
	if !ok {
		return 0, false
	}
	id, ok := value.(int64)
	return id, ok
}
