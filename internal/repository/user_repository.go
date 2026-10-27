package repository

import (
	"strings"

	"video-service/internal/model"
	"video-service/pkg/infrastructure/database"

	"gorm.io/gorm"
)

// UserRepository 用户仓库
type UserRepository interface {
	// ListBasic 分页查询用户基本信息，不读取密码和四个访问码字段
	ListBasic(keyword string, limit, offset int) ([]model.User, int64, error)
}

type userRepository struct{}

// NewUserRepository 创建用户仓库
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// ListBasic 分页查询用户基本信息。
func (r *userRepository) ListBasic(keyword string, limit, offset int) ([]model.User, int64, error) {
	query := database.DB.Model(&model.User{}).
		Select("id", "username", "nickname", "email", "avatar", "watch_count", "last_watch_at", "last_open_at", "created_at", "updated_at")

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", like, like, like)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []model.User
	err := query.Session(&gorm.Session{}).Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
