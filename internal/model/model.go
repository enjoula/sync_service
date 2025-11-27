// model 包定义数据库模型结构
// 使用GORM标签定义数据库字段约束和JSON序列化规则
package model

import (
	"time"

	"gorm.io/datatypes"
)

// User 用户模型
// 存储用户基本信息、认证信息和多设备token
type User struct {
	ID             int64      `gorm:"primaryKey;comment:用户ID，使用算法生成（非自增）" json:"id"`
	Username       string     `gorm:"size:100;uniqueIndex;not null;comment:用户名" json:"username"`
	Password       string     `gorm:"column:password;size:255;not null;comment:密码(加密存储)" json:"-"`
	Nickname       string     `gorm:"size:100;comment:昵称" json:"nickname"`
	Email          string     `gorm:"size:255;comment:邮箱地址" json:"email"`
	Avatar         string     `gorm:"type:text;comment:头像URL" json:"avatar"`
	AccWeb         string     `gorm:"column:acc_web;size:255;comment:Web端访问码" json:"acc_web"`
	AccWebCreateAt *time.Time `gorm:"column:acc_web_create_at;comment:Web端访问码创建时间" json:"acc_web_create_at"`
	AccTV          string     `gorm:"column:acc_tv;size:255;comment:TV端访问码" json:"acc_tv"`
	AccTVCreateAt  *time.Time `gorm:"column:acc_tv_create_at;comment:TV端访问码创建时间" json:"acc_tv_create_at"`
	CreatedAt      time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// UserToken 用户Token记录模型
// 记录用户每次登录生成的token信息，支持多设备登录
type UserToken struct {
	ID        int64      `gorm:"primaryKey;autoIncrement;comment:令牌ID" json:"id"`
	UserID    int64      `gorm:"column:user_id;index;not null;comment:用户ID" json:"user_id"`
	Token     string     `gorm:"size:512;not null;uniqueIndex;comment:登录令牌" json:"token"`
	Device    string     `gorm:"size:100;comment:设备信息" json:"device"`
	IPAddress string     `gorm:"column:ip_address;size:45;comment:IP地址" json:"ip_address"`
	ExpiresAt *time.Time `gorm:"column:expires_at;comment:过期时间" json:"expires_at"`
	IsActive  bool       `gorm:"column:is_active;default:true;comment:是否有效(0:无效,1:有效)" json:"is_active"`
	CreatedAt time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
}

// Video 视频模型
// 存储视频/电影/电视剧的基本信息
type Video struct {
	ID            int64          `gorm:"primaryKey;comment:视频ID，使用雪花算法生成（非自增主键）" json:"id"`
	SourceID      int64          `gorm:"column:source_id;comment:来源站点的视频ID" json:"source_id"`
	Source        string         `gorm:"size:255;comment:视频来源(如:douban、xiaoya)" json:"source"`
	Title         string         `gorm:"size:255;comment:视频标题" json:"title"`
	Type          string         `gorm:"size:32;comment:视频类型(movie/tv/tvshow等)" json:"type"`
	CoverURL      string         `gorm:"column:cover_url;type:text;comment:封面图片地址" json:"cover_url"`
	Description   string         `gorm:"type:text;comment:视频简介" json:"description"`
	ReleaseDate   *time.Time     `gorm:"column:release_date;type:date;comment:上映日期（用于排序和范围查询）" json:"release_date"`
	Score         *float64       `gorm:"column:score;type:decimal(3,1);comment:评分（数值类型，用于排序和范围查询）" json:"score"`
	CountryJSON   datatypes.JSON `gorm:"column:country_json;type:json;comment:国家/地区（JSON数组，支持多值筛选）" json:"country_json"`
	DirectorJSON  datatypes.JSON `gorm:"column:director_json;type:json;comment:导演（JSON数组，支持多值筛选）" json:"director_json"`
	ActorsJSON    datatypes.JSON `gorm:"column:actors_json;type:json;comment:演员列表（JSON数组，支持多值筛选）" json:"actors_json"`
	TagsJSON      datatypes.JSON `gorm:"column:tags_json;type:json;comment:标签（JSON数组，支持多值筛选）" json:"tags_json"`
	Status        string         `gorm:"size:255;comment:状态(用于列表是否返回，0:不 1:返回)" json:"status"`
	IMDbID        string         `gorm:"column:imdb_id;size:20;comment:IMDB 主键" json:"imdb_id"`
	Runtime       *int64         `gorm:"column:runtime;comment:时长" json:"runtime"`
	Resolution    string         `gorm:"size:20;comment:清晰度" json:"resolution"`
	EpisodeCount  *int64         `gorm:"column:episode_count;comment:集数" json:"episode_count"`
	IsCompleted   bool           `gorm:"column:is_completed;default:0;comment:是否完结(0:未完结,1:已完结)" json:"is_completed"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// Episode 剧集/集数模型
// 存储视频的每一集信息，一个Video可以有多个Episode
type Episode struct {
	ID              int64          `gorm:"primaryKey;autoIncrement;comment:剧集ID" json:"id"`
	Channel         string         `gorm:"size:255;comment:频道名称" json:"channel"`
	ChannelID       *int64         `gorm:"column:channel_id;comment:频道ID" json:"channel_id"`
	VideoID         int64          `gorm:"column:video_id;index;not null;comment:所属视频ID" json:"video_id"`
	EpisodeNumber   *int64         `gorm:"column:episode_number;default:1;comment:集数编号" json:"episode_number"`
	Name            string         `gorm:"size:255;comment:剧集名称" json:"name"`
	PlayURLs        string         `gorm:"column:play_urls;size:255;not null;comment:播放地址" json:"play_urls"`
	DurationSeconds *int64         `gorm:"column:duration_seconds;comment:时长(秒)" json:"duration_seconds"`
	SubtitleURLs    datatypes.JSON `gorm:"column:subtitle_urls;type:json;comment:字幕地址列表(JSON格式)" json:"subtitle_urls"`
	CreatedAt       time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// Danmaku 弹幕模型
// 存储视频播放时的弹幕信息
type Danmaku struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:弹幕ID" json:"id"`
	EpisodeID int64     `gorm:"column:episode_id;index;not null;comment:所属剧集ID" json:"episode_id"`
	UserID    *int64    `gorm:"column:user_id;index;comment:发送用户ID" json:"user_id"`
	Content   string    `gorm:"size:255;not null;comment:弹幕内容" json:"content"`
	TimeMs    int64     `gorm:"column:time_ms;not null;comment:弹幕出现时间(毫秒)" json:"time_ms"`
	Color     string    `gorm:"size:20;default:'#FFFFFF';comment:弹幕颜色" json:"color"`
	FontSize  int64     `gorm:"column:font_size;default:16;comment:字体大小" json:"font_size"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
}

// UserFavorite 用户收藏模型
// 记录用户收藏的视频
type UserFavorite struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null;comment:用户ID" json:"user_id"`
	VideoID   int64     `gorm:"column:video_id;not null;comment:视频ID" json:"video_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:收藏时间" json:"created_at"`
}

// FilterInfo 筛选信息模型
// 存储视频筛选条件信息
// 注意：SQL 文件中 filter_info 表的字段没有注释
type FilterInfo struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string `gorm:"size:255" json:"name"`
	Type    string `gorm:"size:255" json:"type"`
	Country string `gorm:"size:255" json:"country"`
	Year    string `gorm:"size:255" json:"year"`
	Tags    string `gorm:"size:255" json:"tags"`
}

// TableName 指定表名
func (FilterInfo) TableName() string {
	return "filter_info"
}

// AppVersion 应用版本模型
// 存储应用的版本信息，用于版本检测和更新管理
type AppVersion struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;comment:版本ID" json:"id"`
	VersionCode   int64     `gorm:"column:version_code;not null;comment:版本号(数字)" json:"version_code"`
	VersionName   string    `gorm:"column:version_name;size:50;not null;comment:版本名称" json:"version_name"`
	Platform      string    `gorm:"size:20;not null;index:idx_app_versions_platform;comment:平台类型(android/ios/web)" json:"platform"`
	DownloadURL   string    `gorm:"column:download_url;type:text;comment:下载链接" json:"download_url"`
	UpdateContent string    `gorm:"column:update_content;type:text;comment:更新内容描述" json:"update_content"`
	IsForce       bool      `gorm:"column:is_force;default:0;comment:是否强制更新(0:否,1:是)" json:"is_force"`
	FileSize      int64     `gorm:"column:file_size;comment:安装包大小(字节)" json:"file_size"`
	IsActive      bool      `gorm:"column:is_active;default:1;comment:是否有效(0:无效,1:有效)" json:"is_active"`
	CreatedAt     time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}
