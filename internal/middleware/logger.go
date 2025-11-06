// middleware 包提供HTTP请求中间件
// LoggerWithZap 提供HTTP请求日志记录功能
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerWithZap 返回一个HTTP请求日志记录中间件
// 功能：
// 1. 记录每个HTTP请求的详细信息（方法、路径、状态码、响应时间等）
// 2. 记录请求处理耗时
// 3. 记录客户端IP地址
// 4. 记录请求大小和响应大小
func LoggerWithZap(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		
		// 记录请求路径和查询参数
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		
		// 记录请求方法
		method := c.Request.Method
		
		// 记录请求大小
		requestSize := c.Request.ContentLength
		if requestSize < 0 {
			requestSize = 0
		}
		
		// 继续处理请求
		c.Next()
		
		// 计算请求处理耗时
		latency := time.Since(start)
		
		// 记录响应状态码
		statusCode := c.Writer.Status()
		
		// 记录响应大小
		responseSize := c.Writer.Size()
		if responseSize < 0 {
			responseSize = 0
		}
		
		// 获取客户端IP地址
		clientIP := c.ClientIP()
		
		// 获取追踪ID（如果存在）
		traceID, _ := c.Get(TraceIDKey)
		traceIDStr := ""
		if tid, ok := traceID.(string); ok {
			traceIDStr = tid
		}
		
		// 构建日志字段
		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.Int64("request_size", requestSize),
			zap.Int("response_size", responseSize),
			zap.String("user_agent", c.Request.UserAgent()),
		}
		
		// 如果有追踪ID，添加到日志字段中
		if traceIDStr != "" {
			fields = append(fields, zap.String("trace_id", traceIDStr))
		}
		
		// 根据状态码选择日志级别
		if statusCode >= 500 {
			// 服务器错误，记录为Error级别
			logger.Error("HTTP Request", fields...)
		} else if statusCode >= 400 {
			// 客户端错误，记录为Warn级别
			logger.Warn("HTTP Request", fields...)
		} else {
			// 成功请求，记录为Info级别
			logger.Info("HTTP Request", fields...)
		}
	}
}

