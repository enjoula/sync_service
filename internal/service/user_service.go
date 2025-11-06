// service 包提供业务逻辑层（Business Logic Layer）
// 封装业务逻辑，协调repository和外部服务
package service

import (
	"errors"
	"time"
	"video-service/internal/auth"
	"video-service/internal/idgen"
	"video-service/internal/models"
	"video-service/internal/nickname"
	"video-service/internal/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	Register(username, password string, device, ipAddress string) (string, *models.User, error)
	Login(username, password string, device, ipAddress string) (string, error)
	GetCurrentUser(c *gin.Context) (string, error)
}

// userService 用户服务实现
type userService struct {
	userRepo      repository.UserRepository
	userTokenRepo repository.UserTokenRepository
}

// NewUserService 创建用户服务
func NewUserService() UserService {
	return &userService{
		userRepo:      repository.NewUserRepository(),
		userTokenRepo: repository.NewUserTokenRepository(),
	}
}

// Register 用户注册
// 注册成功后自动生成token并写入user_tokens表，实现注册即登录
func (s *userService) Register(username, password string, device, ipAddress string) (string, *models.User, error) {
	// 验证输入
	if username == "" || password == "" {
		return "", nil, errors.New("username/password required")
	}

	// 检查用户名是否已存在
	_, err := s.userRepo.FindByUsername(username)
	if err == nil {
		return "", nil, errors.New("用户名重复")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, errors.New("查询用户失败")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, errors.New("密码加密失败")
	}

	// 生成用户ID
	userID := idgen.GenerateUserID()

	// 生成随机昵称
	randomNickname := nickname.GenerateRandomNickname()

	// 创建用户
	user := &models.User{
		ID:       userID,
		Username: username,
		Password: string(hashedPassword),
		Nickname: randomNickname,
	}

	if err := s.userRepo.Create(user); err != nil {
		if repository.IsDuplicateError(err) {
			return "", nil, errors.New("用户名重复")
		}
		return "", nil, errors.New("创建用户失败")
	}

	// 注册成功后自动生成token，实现注册即登录
	expiration := time.Now().Add(24 * time.Hour)
	token, err := auth.GenerateToken(username, expiration)
	if err != nil {
		return "", nil, errors.New("生成token失败")
	}

	// 保存token记录到user_tokens表
	userToken := &models.UserToken{
		UserID:    user.ID,
		Token:     token,
		Device:    device,
		IPAddress: ipAddress,
		ExpiresAt: &expiration,
		IsActive:  true,
	}
	if err := s.userTokenRepo.Create(userToken); err != nil {
		return "", nil, errors.New("保存token失败: " + err.Error())
	}

	// 管理活跃token数量，只保留最新的3个token为活跃状态
	if err := s.userTokenRepo.ManageActiveTokens(user.ID, 3); err != nil {
		// 记录错误但不影响注册流程
		// log.Warn("管理活跃token失败", zap.Error(err))
	}

	// 如果是web设备，更新用户的acc_web字段
	if device == "web" {
		now := time.Now()
		_ = s.userRepo.UpdateWebToken(user.ID, token, now)
	}

	return token, user, nil
}

// Login 用户登录
func (s *userService) Login(username, password string, device, ipAddress string) (string, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// 生成token
	expiration := time.Now().Add(24 * time.Hour)
	token, err := auth.GenerateToken(username, expiration)
	if err != nil {
		return "", errors.New("生成token失败")
	}

	// 保存token记录
	userToken := &models.UserToken{
		UserID:    user.ID,
		Token:     token,
		Device:    device,
		IPAddress: ipAddress,
		ExpiresAt: &expiration,
		IsActive:  true,
	}
	_ = s.userTokenRepo.Create(userToken)

	// 管理活跃token数量，只保留最新的3个token为活跃状态
	if err := s.userTokenRepo.ManageActiveTokens(user.ID, 3); err != nil {
		// 记录错误但不影响登录流程
		// log.Warn("管理活跃token失败", zap.Error(err))
	}

	return token, nil
}

// GetCurrentUser 获取当前登录用户
func (s *userService) GetCurrentUser(c *gin.Context) (string, error) {
	user, exists := c.Get("user")
	if !exists {
		return "", errors.New("未登录")
	}
	username, ok := user.(string)
	if !ok {
		return "", errors.New("用户信息格式错误")
	}
	return username, nil
}
