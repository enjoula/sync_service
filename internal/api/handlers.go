// api 包提供HTTP请求处理函数
package api

import (
	"time"
	"video-service/internal/config"
	"video-service/internal/db"
	"video-service/internal/models"
	"video-service/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// defaultKey 是JWT签名的默认密钥，生产环境应通过配置或Etcd获取
var defaultKey = []byte("change-me-default")

// Claims 定义JWT token的载荷结构
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GetJWTKey 获取JWT签名密钥
// 优先从Etcd配置中获取，如果不存在则使用默认密钥
func GetJWTKey() []byte {
	if config.Secrets != nil && config.Secrets.JWTKey != "" {
		return []byte(config.Secrets.JWTKey)
	}
	return defaultKey
}

// Register 处理用户注册请求
// POST /register
// 请求体: {"username": "string", "password": "string"}
// 功能: 验证用户名和密码，使用bcrypt加密密码后创建用户
func Register(c *gin.Context) {
	var req struct {
		Username, Password string
	}

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, "invalid request")
		return
	}

	// 验证用户名和密码不能为空
	if req.Username == "" || req.Password == "" {
		response.Error(c, response.CodeBadRequest, "username/password required")
		return
	}

	// 使用bcrypt加密密码，默认成本为10
	pw, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// 创建用户记录
	u := models.User{
		Username:     req.Username,
		PasswordHash: string(pw),
	}

	if err := db.DB.Create(&u).Error; err != nil {
		response.Error(c, response.CodeInternalErr, "create user failed")
		return
	}

	// 返回创建成功的用户信息（不包含密码）
	response.Success(c, gin.H{"id": u.ID, "username": u.Username})
}

// Login 处理用户登录请求
// POST /login
// 请求体: {"username": "string", "password": "string"}
// 请求头: X-Device (可选，如"web"或"tv")
// 功能: 验证用户名密码，生成JWT token，记录登录信息
func Login(c *gin.Context) {
	var req struct {
		Username, Password string
	}

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, "invalid request")
		return
	}

	// 根据用户名查询用户
	var u models.User
	if err := db.DB.Where("username = ?", req.Username).First(&u).Error; err != nil {
		response.Error(c, response.CodeUnauthorized, "invalid credentials")
		return
	}

	// 验证密码是否匹配
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, response.CodeUnauthorized, "invalid credentials")
		return
	}

	// 设置token过期时间为24小时后
	expiration := time.Now().Add(24 * time.Hour)

	// 创建JWT claims，包含用户名和过期时间
	claims := Claims{
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	// 使用HS256算法签名生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, _ := token.SignedString(GetJWTKey())

	// 记录用户token信息到数据库，包含设备类型和IP地址
	ut := models.UserToken{
		UserID:    u.ID,
		Token:     ts,
		Device:    c.GetHeader("X-Device"),
		IPAddress: c.ClientIP(),
		ExpiresAt: &expiration,
		IsActive:  true,
	}
	_ = db.DB.Create(&ut)

	// 如果是web设备，更新用户的web_token字段
	if c.GetHeader("X-Device") == "web" {
		now := time.Now()
		db.DB.Model(&u).Updates(map[string]interface{}{
			"web_token":            ts,
			"web_token_created_at": now,
		})
	}

	// 返回token给客户端
	response.Success(c, gin.H{"token": ts})
}

// Me 获取当前登录用户信息
// GET /user/me
// 需要JWT认证中间件，用户信息从context中获取
func Me(c *gin.Context) {
	// 从context中获取用户信息（由JWT中间件设置）
	user := c.MustGet("user").(string)
	response.Success(c, gin.H{"user": user})
}

// Ping 健康检查接口
// GET /ping
// 用于检查服务是否正常运行
func Ping(c *gin.Context) {
	response.SuccessMsg(c, "pong", gin.H{"time": "ok"})
}
