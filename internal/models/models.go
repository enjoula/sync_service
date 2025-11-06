// models 包定义数据库模型结构
// 使用GORM标签定义数据库字段约束和JSON序列化规则
package models

import (
	"time"

	"gorm.io/datatypes"
)

// User 用户模型
// 存储用户基本信息、认证信息和多设备token
type User struct {
	ID             int64      `gorm:"primaryKey" json:"id"`                              // 用户ID，主键，使用算法生成（非自增）
	Username       string     `gorm:"size:100;uniqueIndex;not null" json:"username"`     // 用户名，唯一索引，不能为空
	Password       string     `gorm:"column:password;size:255;not null" json:"-"`        // 密码哈希值，不返回给客户端
	Nickname       string     `gorm:"size:100" json:"nickname"`                          // 昵称
	Email          string     `gorm:"size:255" json:"email"`                             // 邮箱地址
	Avatar         string     `gorm:"type:text" json:"avatar"`                           // 头像URL
	AccWeb         string     `gorm:"column:acc_web;size:255" json:"acc_web"`            // Web端token
	AccWebCreateAt *time.Time `gorm:"column:acc_web_create_at" json:"acc_web_create_at"` // Web token创建时间
	AccTV          string     `gorm:"column:acc_tv;size:255" json:"acc_tv"`              // TV端token
	AccTVCreateAt  *time.Time `gorm:"column:acc_tv_create_at" json:"acc_tv_create_at"`   // TV token创建时间
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`                  // 创建时间，自动设置
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`                  // 更新时间，自动更新
}

// UserToken 用户Token记录模型
// 记录用户每次登录生成的token信息，支持多设备登录
type UserToken struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`         // Token记录ID，主键，自增
	UserID    int64      `gorm:"index;not null" json:"user_id"`              // 用户ID，索引，不能为空
	Token     string     `gorm:"size:255;not null;uniqueIndex" json:"token"` // JWT token字符串，唯一索引
	Device    string     `gorm:"size:100" json:"device"`                     // 设备类型（如web、tv、mobile等）
	IPAddress string     `gorm:"size:45" json:"ip_address"`                  // 登录IP地址（IPv6最长45字符）
	ExpiresAt *time.Time `json:"expires_at"`                                 // Token过期时间
	IsActive  bool       `gorm:"default:true" json:"is_active"`              // 是否激活，默认true
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`           // 创建时间，自动设置
}

// Video 视频模型
// 存储视频/电影/电视剧的基本信息
type Video struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"` // 视频ID，主键，自增
	Title       string    `gorm:"size:255;not null" json:"title"`     // 视频标题，不能为空
	Type        string    `gorm:"size:32;not null" json:"type"`       // 视频类型（如movie、tv、documentary等）
	Description string    `gorm:"type:text" json:"description"`       // 视频描述，文本类型
	Year        int       `json:"year"`                               // 发布年份
	Country     string    `gorm:"size:50" json:"country"`             // 国家/地区
	Director    string    `gorm:"size:255" json:"director"`           // 导演
	Actors      string    `gorm:"size:500" json:"actors"`             // 演员列表（逗号分隔）
	CoverURL    string    `gorm:"type:text" json:"cover_url"`         // 封面图片URL
	Tags        string    `gorm:"size:255" json:"tags"`               // 标签（逗号分隔）
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`   // 创建时间，自动设置
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`   // 更新时间，自动更新
}

// Episode 剧集/集数模型
// 存储视频的每一集信息，一个Video可以有多个Episode
type Episode struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`  // 剧集ID，主键，自增
	VideoID         int64          `gorm:"index;not null" json:"video_id"`      // 所属视频ID，索引，不能为空
	EpisodeNumber   int            `gorm:"default:1" json:"episode_number"`     // 集数，默认第1集
	Name            string         `gorm:"size:255" json:"name"`                // 集名称
	PlayURLs        datatypes.JSON `gorm:"type:json;not null" json:"play_urls"` // 播放地址列表（JSON数组），支持多源
	DurationSeconds int            `json:"duration_seconds"`                    // 时长（秒）
	SubtitleURLs    datatypes.JSON `gorm:"type:json" json:"subtitle_urls"`      // 字幕文件URL列表（JSON数组），支持多语言
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`    // 创建时间，自动设置
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`    // 更新时间，自动更新
}

// Danmaku 弹幕模型
// 存储视频播放时的弹幕信息
type Danmaku struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`     // 弹幕ID，主键，自增
	EpisodeID int64     `gorm:"index;not null" json:"episode_id"`       // 所属剧集ID，索引，不能为空
	UserID    *int64    `gorm:"index" json:"user_id"`                   // 发送用户ID（可为空，支持匿名弹幕），索引
	Content   string    `gorm:"size:255;not null" json:"content"`       // 弹幕内容，不能为空
	TimeMs    int       `gorm:"not null" json:"time_ms"`                // 弹幕出现时间点（毫秒），不能为空
	Color     string    `gorm:"size:20;default:'#FFFFFF'" json:"color"` // 弹幕颜色，默认白色
	FontSize  int       `gorm:"default:16" json:"font_size"`            // 字体大小，默认16
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`       // 创建时间，自动设置
}

// UserFavorite 用户收藏模型
// 记录用户收藏的视频，使用复合主键（UserID + VideoID）
type UserFavorite struct {
	UserID    int64     `gorm:"primaryKey" json:"user_id"`        // 用户ID，复合主键
	VideoID   int64     `gorm:"primaryKey" json:"video_id"`       // 视频ID，复合主键
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 收藏时间，自动设置
}

// UserWatchProgress 用户观看进度模型
// 记录用户观看每个剧集的进度信息，使用复合主键（UserID + EpisodeID）
type UserWatchProgress struct {
	UserID         int64     `gorm:"primaryKey" json:"user_id"`            // 用户ID，复合主键
	EpisodeID      int64     `gorm:"primaryKey" json:"episode_id"`         // 剧集ID，复合主键
	LastPositionMs int64     `gorm:"default:0" json:"last_position_ms"`    // 最后观看位置（毫秒），默认0
	LastPlayedAt   time.Time `gorm:"autoCreateTime" json:"last_played_at"` // 最后播放时间，自动设置
}
