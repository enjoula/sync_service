// repository 包提供数据访问层（Data Access Layer）
// 封装数据库操作，提供统一的数据访问接口
package repository

import (
	"errors"
	"strings"
	"video-service/internal/db"
	"video-service/internal/models"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	FindByID(id int64) (*models.User, error)
	UpdateWebToken(userID int64, token string, createdAt interface{}) error
}

// userRepository 用户数据访问实现
type userRepository struct{}

// NewUserRepository 创建用户数据访问对象
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// Create 创建用户
func (r *userRepository) Create(user *models.User) error {
	return db.DB.Create(user).Error
}

// FindByUsername 根据用户名查找用户
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	err := db.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateWebToken 更新用户的Web token
func (r *userRepository) UpdateWebToken(userID int64, token string, createdAt interface{}) error {
	return db.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"acc_web":           token,
		"acc_web_create_at": createdAt,
	}).Error
}

// UserTokenRepository 用户Token数据访问接口
type UserTokenRepository interface {
	Create(token *models.UserToken) error
	CountByUserID(userID int64) (int64, error)
	ManageActiveTokens(userID int64, keepCount int) error
	IsTokenActive(token string) (bool, error)
}

// userTokenRepository 用户Token数据访问实现
type userTokenRepository struct{}

// NewUserTokenRepository 创建用户Token数据访问对象
func NewUserTokenRepository() UserTokenRepository {
	return &userTokenRepository{}
}

// Create 创建用户Token记录
func (r *userTokenRepository) Create(token *models.UserToken) error {
	return db.DB.Create(token).Error
}

// CountByUserID 统计指定用户的token数量
func (r *userTokenRepository) CountByUserID(userID int64) (int64, error) {
	var count int64
	err := db.DB.Model(&models.UserToken{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// ManageActiveTokens 管理用户的活跃token
// 保持最新的keepCount个token为活跃状态，其余设为非活跃
func (r *userTokenRepository) ManageActiveTokens(userID int64, keepCount int) error {
	// 首先将该用户所有token设为非活跃
	if err := db.DB.Model(&models.UserToken{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error; err != nil {
		return err
	}

	// 查询最新的keepCount个token的ID
	var tokenIDs []int64
	if err := db.DB.Model(&models.UserToken{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(keepCount).
		Pluck("id", &tokenIDs).Error; err != nil {
		return err
	}

	// 如果有token需要设为活跃
	if len(tokenIDs) > 0 {
		if err := db.DB.Model(&models.UserToken{}).
			Where("id IN ?", tokenIDs).
			Update("is_active", true).Error; err != nil {
			return err
		}
	}

	return nil
}

// IsTokenActive 检查token是否在数据库中且处于活跃状态
func (r *userTokenRepository) IsTokenActive(token string) (bool, error) {
	var userToken models.UserToken
	err := db.DB.Where("token = ? AND is_active = ?", token, true).First(&userToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Token不存在或未激活
			return false, nil
		}
		// 数据库查询错误
		return false, err
	}
	// Token存在且活跃
	return true, nil
}

// IsDuplicateError 检查是否是重复键错误
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	errStrLower := strings.ToLower(errStr)
	return errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(errStrLower, "duplicate entry") ||
		strings.Contains(errStrLower, "1062")
}
