// handler 包提供HTTP请求处理函数
// auth.go 提供认证相关的HTTP处理器
package handler

import (
	"video-service/internal/pkg/response"
	"video-service/internal/pkg/utils"
	"video-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var userService = service.NewUserService()

// Register 处理用户注册请求
// POST /register
// 请求体: {"username": "string", "password": "string"}
// 请求头: X-Device (可选，如"web"或"tv")
// 功能: 验证用户名和密码，使用bcrypt加密密码后创建用户，并自动生成token实现注册即登录
func Register(c *gin.Context) {
	log := zap.L()

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 获取请求元数据
	device := c.GetHeader("X-Device")
	ipAddress := utils.GetRealIP(c)

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("注册请求绑定失败", zap.Error(err), zap.String("ip", ipAddress))
		response.Error(c, response.CodeBadRequest, "invalid request")
		return
	}

	// 调用服务层注册用户（注册成功后自动生成token）
	token, user, err := userService.Register(req.Username, req.Password, device, ipAddress)
	if err != nil {
		log.Error("注册失败", zap.String("username", req.Username), zap.String("ip", ipAddress), zap.Error(err))
		switch err.Error() {
		case "username/password required":
			response.Error(c, response.CodeBadRequest, err.Error())
		case "用户名重复":
			response.Error(c, response.CodeConflict, err.Error())
		default:
			response.Error(c, response.CodeInternalErr, err.Error())
		}
		return
	}

	log.Info("注册成功", zap.Int64("user_id", user.ID), zap.String("username", req.Username), zap.String("ip", ipAddress))

	// 返回创建成功的用户信息和token（实现注册即登录）
	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"token":    token,
	})
}

// Login 处理用户登录请求
// POST /login
// 请求体: {"username": "string", "password": "string"}
// 请求头: X-Device (可选，如"web"或"tv")
// 功能: 验证用户名密码，生成JWT token，记录登录信息
func Login(c *gin.Context) {
	log := zap.L()

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 获取请求元数据
	device := c.GetHeader("X-Device")
	ipAddress := utils.GetRealIP(c)

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("登录请求绑定失败", zap.Error(err), zap.String("ip", ipAddress))
		response.Error(c, response.CodeBadRequest, "invalid request")
		return
	}

	// 调用服务层登录
	token, user, err := userService.Login(req.Username, req.Password, device, ipAddress)
	if err != nil {
		log.Error("登录失败", zap.String("username", req.Username), zap.String("ip", ipAddress), zap.Error(err))
		if err.Error() == "invalid credentials" {
			response.Error(c, response.CodeUnauthorized, err.Error())
		} else {
			response.Error(c, response.CodeInternalErr, err.Error())
		}
		return
	}

	log.Info("登录成功", zap.Int64("user_id", user.ID), zap.String("username", req.Username), zap.String("ip", ipAddress))

	// 返回用户信息和token给客户端
	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"token":    token,
	})
}
