// service 包提供业务逻辑层（Business Logic Layer）
// 封装业务逻辑，协调repository和外部服务
package service

import (
	"errors"
	"regexp"
	"time"
	"video-service/internal/auth"
	bizerrors "video-service/internal/errors"
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
	Login(username, password string, device, ipAddress string) (string, *models.User, error)
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
		return "", nil, bizerrors.ErrUsernamePasswordEmpty
	}

	// 验证用户名长度（4-15个字符）
	if len(username) < 4 || len(username) > 15 {
		return "", nil, bizerrors.ErrUsernameLengthInvalid
	}

	// 验证用户名格式（只允许字母和数字）
	matched, err := regexp.MatchString("^[a-zA-Z0-9]+$", username)
	if err != nil {
		return "", nil, bizerrors.ErrUsernameFormatInvalid
	}
	if !matched {
		return "", nil, bizerrors.ErrUsernameInvalidChars
	}

	// 检查用户名是否已存在
	_, err = s.userRepo.FindByUsername(username)
	if err == nil {
		return "", nil, bizerrors.ErrUsernameDuplicate
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, bizerrors.ErrUserQueryFailed
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, bizerrors.ErrPasswordEncryptFailed
	}

	// 生成用户ID
	userID := idgen.GenerateUserID()

	// 生成随机卡通人物昵称
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
			return "", nil, bizerrors.ErrUsernameDuplicate
		}
		return "", nil, bizerrors.ErrUserCreateFailed
	}

	// 注册成功后自动生成token，实现注册即登录
	// token过期时间设置为1年
	expiration := time.Now().Add(365 * 24 * time.Hour)
	token, err := auth.GenerateToken(username, expiration)
	if err != nil {
		return "", nil, bizerrors.ErrTokenGenerateFailed
	}

	// 检查该用户的token数量，如果大于等于3则将最早创建时间的那条数据的is_active置为0
	tokenCount, err := s.userTokenRepo.CountByUserID(user.ID)
	if err != nil {
		return "", nil, bizerrors.ErrTokenQueryFailed
	}
	if tokenCount >= 3 {
		if err := s.userTokenRepo.DeactivateOldestToken(user.ID); err != nil {
			return "", nil, bizerrors.ErrTokenDeactivateFailed
		}
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
		// 检查是否是token重复错误
		if repository.IsDuplicateError(err) {
			return "", nil, bizerrors.ErrTokenDuplicate
		}
		return "", nil, bizerrors.ErrTokenSaveFailed
	}

	// 注册时不更新users表中的web_token和tv_token_created_at字段
	// token信息仅保存在user_tokens表中

	return token, user, nil
}

// Login 用户登录
func (s *userService) Login(username, password string, device, ipAddress string) (string, *models.User, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", nil, bizerrors.ErrUsernamePasswordError
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, bizerrors.ErrUsernamePasswordError
	}

	// 生成token，过期时间设置为1年
	expiration := time.Now().Add(365 * 24 * time.Hour)
	token, err := auth.GenerateToken(username, expiration)
	if err != nil {
		return "", nil, bizerrors.ErrTokenGenerateFailed
	}

	// 检查该用户的token数量，如果大于3则将最新创建时间的那条数据的is_active置为0
	tokenCount, err := s.userTokenRepo.CountByUserID(user.ID)
	if err != nil {
		return "", nil, bizerrors.ErrTokenQueryFailed
	}
	if tokenCount >= 3 {
		if err := s.userTokenRepo.DeactivateOldestToken(user.ID); err != nil {
			return "", nil, bizerrors.ErrTokenDeactivateFailed
		}
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
	if err := s.userTokenRepo.Create(userToken); err != nil {
		// 检查是否是token重复错误
		if repository.IsDuplicateError(err) {
			return "", nil, bizerrors.ErrTokenDuplicate
		}
		return "", nil, bizerrors.ErrTokenSaveFailed
	}

	// 登录时不更新users表中的web_token和tv_token_created_at字段
	// token信息仅保存在user_tokens表中

	return token, user, nil
}

// GetCurrentUser 获取当前登录用户
func (s *userService) GetCurrentUser(c *gin.Context) (string, error) {
	user, exists := c.Get("user")
	if !exists {
		return "", bizerrors.ErrNotLoggedIn
	}
	username, ok := user.(string)
	if !ok {
		return "", bizerrors.ErrUserInfoFormatError
	}
	return username, nil
}
