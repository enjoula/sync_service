package repository

import (
	"strings"
	"time"

	"video-service/internal/model"
	"video-service/pkg/infrastructure/database"

	"gorm.io/gorm"
)

// UserWatchRecordItem 管理端观看记录，附带用户和视频标题。
type UserWatchRecordItem struct {
	ID             int64
	UserID         int64
	Username       string
	Nickname       string
	VideoID        int64
	VideoTitle     string
	EpisodeID      int64
	EpisodeName    string
	EpisodeNumber  int64
	PlayCount      int64
	LastPositionMs int64
	LastPlayedAt   time.Time
	CreatedAt      time.Time
}

// UserLoginItem 管理端登录记录。不包含登录令牌。
type UserLoginItem struct {
	ID        int64
	UserID    int64
	Username  string
	Nickname  string
	Device    string
	IPAddress string
	ExpiresAt *time.Time
	IsActive  bool
	CreatedAt time.Time
}

// UserRepository 用户仓库
type UserRepository interface {
	// ListBasic 分页查询用户基本信息，不读取密码和四个访问码字段
	ListBasic(keyword string, limit, offset int) ([]model.User, int64, error)
	// ListWatchRecords 分页查询观看记录，按最后观看时间倒序
	ListWatchRecords(keyword string, limit, offset int) ([]UserWatchRecordItem, int64, error)
	// ListLogins 分页查询登录记录，按登录时间倒序，不读取令牌
	ListLogins(keyword string, limit, offset int) ([]UserLoginItem, int64, error)
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

// ListWatchRecords 分页查询观看记录。
func (r *userRepository) ListWatchRecords(keyword string, limit, offset int) ([]UserWatchRecordItem, int64, error) {
	query := database.DB.Table("user_watch_records AS r").
		Joins("LEFT JOIN users AS u ON u.id = r.user_id").
		Joins("LEFT JOIN videos AS v ON v.id = r.video_id")
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("u.username LIKE ? OR u.nickname LIKE ? OR v.title LIKE ?", like, like, like)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []UserWatchRecordItem
	err := query.Session(&gorm.Session{}).
		Select(`r.id, r.user_id, u.username, u.nickname, r.video_id, v.title AS video_title,
			r.episode_id, r.episode_name, r.episode_number, r.play_count, r.last_position_ms,
			r.last_played_at, r.created_at`).
		Order("r.last_played_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListLogins 分页查询登录记录，不读取 token 字段。
func (r *userRepository) ListLogins(keyword string, limit, offset int) ([]UserLoginItem, int64, error) {
	query := database.DB.Table("user_tokens AS t").
		Joins("LEFT JOIN users AS u ON u.id = t.user_id")
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("u.username LIKE ? OR u.nickname LIKE ? OR t.ip_address LIKE ? OR t.device LIKE ?", like, like, like, like)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []UserLoginItem
	err := query.Session(&gorm.Session{}).
		Select("t.id, t.user_id, u.username, u.nickname, t.device, t.ip_address, t.expires_at, t.is_active, t.created_at").
		Order("t.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
