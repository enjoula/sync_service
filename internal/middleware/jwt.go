// middleware 包提供HTTP请求中间件
// JWTAuth 提供JWT认证功能
package middleware

import (
	"strings"
	"video-service/internal/auth"
	"video-service/internal/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth 返回一个JWT认证中间件
// 功能：
// 1. 从请求头中提取JWT token
// 2. 验证token的有效性（签名、过期时间等）
// 3. 解析token中的用户信息并存储到context中
// 4. 如果认证失败，返回401错误
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization请求头中获取token
		// 格式: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, response.CodeUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, response.CodeUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		// 提取token字符串
		tokenString := parts[1]

		// 解析和验证token
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			response.Error(c, response.CodeUnauthorized, "invalid token: "+err.Error())
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

