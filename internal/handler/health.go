// handler 包提供HTTP请求处理函数
// health.go 提供健康检查相关的HTTP处理器
package handler

import (
	"video-service/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// Ping 健康检查接口
// GET /ping
// 用于检查服务是否正常运行
func Ping(c *gin.Context) {
	response.SuccessMsg(c, "pong", gin.H{"time": "ok"})
}
