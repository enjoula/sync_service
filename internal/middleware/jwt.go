// middleware 包提供HTTP请求中间件
// JWTAuth 提供JWT认证功能
package middleware

import (
	"strings"
	"video-service/internal/auth"
	"video-service/internal/errors"
	"video-service/internal/repository"
	"video-service/internal/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth 返回一个JWT认证中间件
// 功能：
// 1. 从请求头中提取JWT token
// 2. 检查token在数据库中的is_active状态
// 3. 仅当is_active=1时，验证token的签名和过期时间
// 4. 解析token中的用户信息并存储到context中
// 5. 如果认证失败，返回401错误
func JWTAuth() gin.HandlerFunc {
	// 创建token repository实例
	tokenRepo := repository.NewUserTokenRepository()

	return func(c *gin.Context) {
		// 从Authorization请求头中获取token
		// 格式: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errors.CodeUnauthorized, errors.MsgTokenMissing)
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, errors.CodeUnauthorized, errors.MsgTokenInvalidFormat)
			c.Abort()
			return
		}

		// 提取token字符串
		tokenString := parts[1]

		// 步骤1: 先检查数据库中token的is_active状态
		isActive, err := tokenRepo.IsTokenActive(tokenString)
		if err != nil {
			// 数据库查询错误
			response.Error(c, errors.CodeInternalErr, "token验证失败")
			c.Abort()
			return
		}

		if !isActive {
			// Token不存在或已失效（is_active=0）
			response.Error(c, errors.CodeUnauthorized, "token已过期")
			c.Abort()
			return
		}

		// 步骤2: 只有当is_active=1时，才解析和验证token签名和过期时间
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			response.Error(c, errors.CodeUnauthorized, errors.NewTokenInvalid(err).GetMessage())
			c.Abort()
			return
		}

		// 将用户名存储到context中，供后续处理器使用
		c.Set("user", claims.Username)
		c.Set("user_id", claims.Username) // 可以根据需要添加更多用户信息

		// 认证成功，继续处理请求
		c.Next()
	}
}
