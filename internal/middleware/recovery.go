// middleware 包提供HTTP请求中间件
// RecoveryWithZap 提供panic恢复和日志记录功能
package middleware

import (
	"net/http"
	"runtime/debug"
	"time"
	"video-service/internal/errors"
	"video-service/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryWithZap 返回一个panic恢复中间件
// 功能：
// 1. 捕获请求处理过程中发生的panic
// 2. 记录详细的错误信息（错误内容、堆栈跟踪、请求路径、处理时间）
// 3. 返回统一的错误响应给客户端
// 4. 防止panic导致整个服务崩溃
func RecoveryWithZap(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 使用defer确保panic被捕获
		defer func() {
			// 捕获panic
			if r := recover(); r != nil {
				// 记录panic详细信息到日志
				logger.Error("panic recovered",
					zap.Any("error", r),                        // panic的错误内容
					zap.ByteString("stack", debug.Stack()),     // 完整的堆栈跟踪
					zap.String("path", c.Request.URL.Path),     // 请求路径
					zap.Duration("elapsed", time.Since(start)), // 请求处理耗时
				)

				// 返回统一的错误响应给客户端
				response.Error(c, errors.CodeInternalErr, errors.NewServerPanic(r).GetMessage())

				// 中止请求处理，返回HTTP 200状态码（业务错误统一返回200，通过code字段区分）
				c.AbortWithStatus(http.StatusOK)
				return
			}
		}()

		// 继续处理请求
		c.Next()
	}
}
