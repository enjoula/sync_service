// handler 包提供HTTP请求处理函数
// user.go 提供用户相关的HTTP处理器
package handler

import (
	"video-service/internal/pkg/errors"
	"video-service/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// Me 获取当前登录用户信息
// GET /user/me
// 需要JWT认证中间件，用户信息从context中获取
func Me(c *gin.Context) {
	username, err := userService.GetCurrentUser(c)
	if err != nil {
		// 检查是否是业务错误类型
		if bizErr, ok := err.(*errors.BusinessError); ok {
			response.Error(c, bizErr.GetCode(), bizErr.GetMessage())
		} else {
			response.Error(c, errors.CodeUnauthorized, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"user": username})
}
