// middleware 包提供HTTP请求中间件
// Trace 提供请求追踪ID生成功能
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceIDKey 是context中存储追踪ID的键名
const TraceIDKey = "trace_id"

// Trace 返回一个请求追踪中间件
// 功能：为每个HTTP请求生成唯一的追踪ID，并存储到context中
// 追踪ID可用于分布式系统追踪和日志关联
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头中获取追踪ID（如果上游服务已设置）
		traceID := c.GetHeader("X-Trace-ID")

		// 如果请求头中没有追踪ID，则生成一个新的UUID
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 将追踪ID存储到context中，供后续中间件和处理器使用
		c.Set(TraceIDKey, traceID)

		// 将追踪ID添加到响应头中，方便客户端追踪
		c.Header("X-Trace-ID", traceID)

		// 继续处理请求
		c.Next()
	}
}
