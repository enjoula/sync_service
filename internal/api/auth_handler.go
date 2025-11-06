// api 包提供HTTP请求处理函数
// auth_handler.go 提供认证相关的HTTP处理器
package api

import (
	"video-service/internal/errors"
	"video-service/internal/response"
	"video-service/internal/service"

	"github.com/gin-gonic/gin"
)

var userService = service.NewUserService()

// Register 处理用户注册请求
// POST /register
// 请求体: {"username": "string", "password": "string"}
// 请求头: X-Device (可选，如"web"或"tv")
// 功能: 验证用户名和密码，使用bcrypt加密密码后创建用户，并自动生成token实现注册即登录
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeBadRequest, errors.MsgBadRequest)
		return
	}

	// 调用服务层注册用户（注册成功后自动生成token）
	device := c.GetHeader("X-Device")
	ipAddress := c.ClientIP()
	token, user, err := userService.Register(req.Username, req.Password, device, ipAddress)
	if err != nil {
		// 检查是否是业务错误类型
		if bizErr, ok := err.(*errors.BusinessError); ok {
			response.Error(c, bizErr.GetCode(), bizErr.GetMessage())
		} else {
			response.Error(c, errors.CodeInternalErr, err.Error())
		}
		return
	}

	// 返回创建成功的用户信息和token（实现注册即登录）
	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"token":    token,
	})
}

// Login 处理用户登录请求
// POST /login
// 请求体: {"username": "string", "password": "string", "device": "string"}
// 功能: 验证用户名密码，生成JWT token（过期时间一年），记录登录信息
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Device   string `json:"device"`
	}

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeBadRequest, errors.MsgBadRequest)
		return
	}

	// 验证username和password参数不能为空
	if req.Username == "" || req.Password == "" {
		response.Error(c, errors.CodeBadRequest, errors.MsgUsernamePasswordEmpty)
		return
	}

	// 调用服务层登录
	ipAddress := c.ClientIP()
	token, user, err := userService.Login(req.Username, req.Password, req.Device, ipAddress)
	if err != nil {
		// 检查是否是业务错误类型
		if bizErr, ok := err.(*errors.BusinessError); ok {
			response.Error(c, bizErr.GetCode(), bizErr.GetMessage())
		} else {
			response.Error(c, errors.CodeInternalErr, err.Error())
		}
		return
	}

	// 返回用户信息和token给客户端
	response.Success(c, gin.H{
		"username": user.Username,
		"avatar":   user.Avatar,
		"nickname": user.Nickname,
		"token":    token,
	})
}
