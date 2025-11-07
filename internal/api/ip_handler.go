// api 包提供HTTP请求处理函数
// ip_handler.go 提供IP信息查询的HTTP处理器
package api

import (
	"video-service/internal/response"
	"video-service/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetIPInfo 获取客户端IP信息
// GET /ip-info
// 用于测试和调试真实IP获取功能
// 返回详细的IP信息，包括从各个头部获取的IP地址
func GetIPInfo(c *gin.Context) {
	log := zap.L()

	// 获取基本IP信息
	ipInfo := utils.GetIPInfo(c)

	// 获取所有请求头（用于调试）
	allHeaders := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			allHeaders[key] = values[0]
		}
	}

	// 详细日志记录所有IP相关信息
	log.Info("IP调试信息",
		zap.String("real_ip", ipInfo["real_ip"]),
		zap.String("remote_addr", c.Request.RemoteAddr),
		zap.String("x_real_ip", c.GetHeader("X-Real-IP")),
		zap.String("x_forwarded_for", c.GetHeader("X-Forwarded-For")),
		zap.String("cf_connecting_ip", c.GetHeader("CF-Connecting-IP")),
		zap.String("true_client_ip", c.GetHeader("True-Client-IP")),
		zap.String("x_forwarded_proto", c.GetHeader("X-Forwarded-Proto")),
		zap.String("x_forwarded_host", c.GetHeader("X-Forwarded-Host")),
		zap.Any("all_headers", allHeaders),
	)

	// 返回详细信息
	response.Success(c, gin.H{
		"ip_info":     ipInfo,
		"all_headers": allHeaders,
		"remote_addr": c.Request.RemoteAddr,
		"client_ip":   c.ClientIP(),
	})
}
