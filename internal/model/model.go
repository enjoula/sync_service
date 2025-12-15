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
	ID             int64      `gorm:"primaryKey;comment:用户ID，使用算法生成（非自增）" json:"id"`                         // 用户ID，主键，使用算法生成（非自增）
	Username       string     `gorm:"size:100;uniqueIndex;not null;comment:用户名" json:"username"`             // 用户名，唯一索引，不能为空
	Password       string     `gorm:"column:password;size:255;not null;comment:密码(加密存储)" json:"-"`           // 密码哈希值，不返回给客户端
	Nickname       string     `gorm:"size:100;comment:昵称" json:"nickname"`                                   // 昵称
	Email          string     `gorm:"size:255;comment:邮箱地址" json:"email"`                                    // 邮箱地址
	Avatar         string     `gorm:"type:text;comment:头像URL" json:"avatar"`                                 // 头像URL
	AccWeb         string     `gorm:"column:acc_web;size:255;comment:Web端访问码" json:"acc_web"`                // Web端token
	AccWebCreateAt *time.Time `gorm:"column:acc_web_create_at;comment:Web端访问码创建时间" json:"acc_web_create_at"` // Web token创建时间
	AccTV          string     `gorm:"column:acc_tv;size:255;comment:TV端访问码" json:"acc_tv"`                   // TV端token
	AccTVCreateAt  *time.Time `gorm:"column:acc_tv_create_at;comment:TV端访问码创建时间" json:"acc_tv_create_at"`    // TV token创建时间
	CreatedAt      time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`                         // 创建时间，自动设置
	UpdatedAt      time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`                         // 更新时间，自动更新
}

// UserToken 用户Token记录模型
// 记录用户每次登录生成的token信息，支持多设备登录
type UserToken struct {
	ID        int64      `gorm:"primaryKey;autoIncrement;comment:令牌ID" json:"id"`                        // Token记录ID，主键，自增
	UserID    int64      `gorm:"column:user_id;index;not null;comment:用户ID" json:"user_id"`              // 用户ID，索引，不能为空
	Token     string     `gorm:"size:512;not null;uniqueIndex;comment:登录令牌" json:"token"`                // JWT token字符串，唯一索引（必须512，JWT token超过255字符）
	Device    string     `gorm:"size:100;comment:设备信息" json:"device"`                                    // 设备类型（如web、tv、mobile等）
	IPAddress string     `gorm:"column:ip_address;size:45;comment:IP地址" json:"ip_address"`               // 登录IP地址（IPv6最长45字符）
	ExpiresAt *time.Time `gorm:"column:expires_at;comment:过期时间" json:"expires_at"`                       // Token过期时间
	IsActive  bool       `gorm:"column:is_active;default:true;comment:是否有效(0:无效,1:有效)" json:"is_active"` // 是否激活，默认true
	CreatedAt time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`                          // 创建时间，自动设置
}

// Video 视频模型
// 存储视频/电影/电视剧的基本信息
type Video struct {
	ID           int64          `gorm:"primaryKey;comment:视频ID，使用雪花算法生成（非自增主键）" json:"id"`                              // 视频ID，主键，使用雪花算法生成（非自增）
	SourceID     *int64         `gorm:"column:source_id;comment:来源站点的视频ID" json:"source_id"`                            // 视频来源三方ID
	Source       string         `gorm:"size:255;comment:视频来源(如:douban、xiaoya)" json:"source"`                           // 视频来源三方
	Title        string         `gorm:"size:255;comment:视频标题" json:"title"`                                             // 视频标题
	Type         string         `gorm:"size:32;comment:视频类型(movie/tv/tvshow等)" json:"type"`                             // 视频类型（movie/tv/tvshow/anime/variety/doc）
	CoverURL     string         `gorm:"column:cover_url;type:text;comment:封面图片地址" json:"cover_url"`                     // 封面图片URL
	Description  string         `gorm:"type:text;comment:视频简介" json:"description"`                                      // 视频描述，文本类型
	ReleaseDate  *time.Time     `gorm:"column:release_date;type:date;comment:上映日期（用于排序和范围查询）" json:"release_date"`      // 上映日期（用于排序和范围查询）
	Score        *float64       `gorm:"column:score;type:decimal(3,1);comment:评分（数值类型，用于排序和范围查询）" json:"score"`         // 评分（数值类型，用于排序和范围查询）
	CountryJSON  datatypes.JSON `gorm:"column:country_json;type:json;comment:国家/地区（JSON数组，支持多值筛选）" json:"country_json"` // 国家/地区（JSON数组，支持多值筛选）
	DirectorJSON datatypes.JSON `gorm:"column:director_json;type:json;comment:导演（JSON数组，支持多值筛选）" json:"director_json"`  // 导演（JSON数组，支持多值筛选）
	ActorsJSON   datatypes.JSON `gorm:"column:actors_json;type:json;comment:演员列表（JSON数组，支持多值筛选）" json:"actors_json"`    // 演员列表（JSON数组，支持多值筛选）
	TagsJSON     datatypes.JSON `gorm:"column:tags_json;type:json;comment:标签（JSON数组，支持多值筛选）" json:"tags_json"`          // 标签（JSON数组，支持多值筛选）
	Status       string         `gorm:"size:255;comment:状态(用于列表是否返回，0:不 1:返回)" json:"status"`                           // 状态（用于列表是否返回，0:不 1:返回）
	IMDbID       string         `gorm:"column:imdb_id;size:20;comment:IMDB 主键" json:"imdb_id"`                          // IMDb ID
	Runtime      *int64         `gorm:"comment:时长" json:"runtime"`                                                      // 时长（分钟）
	Resolution   string         `gorm:"size:20;comment:清晰度" json:"resolution"`                                          // 清晰度（如：1080p、720p、4K）
	EpisodeCount int64          `gorm:"column:episode_count;comment:集数" json:"episode_count"`                           // 集数
	IsCompleted  bool           `gorm:"column:is_completed;default:0;comment:是否完结(0:未完结,1:已完结)" json:"is_completed"`    // 是否完结（1:是 0:否）
	IsUpdate     bool           `gorm:"column:is_update;default:0;comment:是否有更新(0:无更新,1:有更新)" json:"is_update"`         // 是否有更新（1:有更新 0:无更新）
	CreatedAt    time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`                                  // 创建时间，自动设置
	UpdatedAt    time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`                                  // 更新时间，自动更新
}

// Episode 剧集/集数模型
// 存储视频的每一集信息，一个Video可以有多个Episode
type Episode struct {
	ID              int64          `gorm:"primaryKey;autoIncrement;comment:剧集ID" json:"id"`                            // 剧集ID，主键，自增
	Channel         string         `gorm:"size:255;comment:频道名称" json:"channel"`                                       // 同步渠道
	ChannelID       int64          `gorm:"column:channel_id;comment:频道ID" json:"channel_id"`                           // 渠道视频ID
	VideoID         int64          `gorm:"column:video_id;index;not null;comment:所属视频ID" json:"video_id"`              // 所属视频ID，索引，不能为空
	EpisodeNumber   int64          `gorm:"column:episode_number;default:1;comment:集数编号" json:"episode_number"`         // 集数，默认第1集
	Name            string         `gorm:"size:255;comment:剧集名称" json:"name"`                                          // 集名称
	PlayURLs        string         `gorm:"column:play_urls;size:255;not null;comment:播放地址" json:"play_urls"`           // 播放地址
	DurationSeconds int64          `gorm:"column:duration_seconds;comment:时长(秒)" json:"duration_seconds"`              // 时长（秒）
	SubtitleURLs    datatypes.JSON `gorm:"column:subtitle_urls;type:json;comment:字幕地址列表(JSON格式)" json:"subtitle_urls"` // 字幕文件URL列表（JSON数组），支持多语言
	CreatedAt       time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`                              // 创建时间，自动设置
	UpdatedAt       time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`                              // 更新时间，自动更新
}

// Danmaku 弹幕模型
// 存储视频播放时的弹幕信息
type Danmaku struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:弹幕ID" json:"id"`                   // 弹幕ID，主键，自增
	EpisodeID int64     `gorm:"column:episode_id;index;not null;comment:所属剧集ID" json:"episode_id"` // 所属剧集ID，索引，不能为空
	UserID    *int64    `gorm:"column:user_id;index;comment:发送用户ID" json:"user_id"`                // 发送用户ID（可为空，支持匿名弹幕），索引
	Content   string    `gorm:"size:255;not null;comment:弹幕内容" json:"content"`                     // 弹幕内容，不能为空
	TimeMs    int64     `gorm:"column:time_ms;not null;comment:弹幕出现时间(毫秒)" json:"time_ms"`         // 弹幕出现时间点（毫秒），不能为空
	Color     string    `gorm:"size:20;default:'#FFFFFF';comment:弹幕颜色" json:"color"`               // 弹幕颜色，默认白色
	FontSize  int64     `gorm:"column:font_size;default:16;comment:字体大小" json:"font_size"`         // 字体大小，默认16
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`                     // 创建时间，自动设置
}

// UserFavorite 用户收藏模型
// 记录用户收藏的视频
type UserFavorite struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`                    // 收藏ID，主键，自增
	UserID    int64     `gorm:"column:user_id;not null;comment:用户ID" json:"user_id"`   // 用户ID
	VideoID   int64     `gorm:"column:video_id;not null;comment:视频ID" json:"video_id"` // 视频ID
	CreatedAt time.Time `gorm:"autoCreateTime;comment:收藏时间" json:"created_at"`         // 收藏时间，自动设置
}

// UserWatchProgress 用户观看进度模型
// 记录用户观看每个剧集的进度信息，使用复合主键（UserID + EpisodeID）
type UserWatchProgress struct {
	UserID         int64     `gorm:"column:user_id;primaryKey" json:"user_id"`                  // 用户ID，复合主键
	EpisodeID      int64     `gorm:"column:episode_id;primaryKey" json:"episode_id"`            // 剧集ID，复合主键
	LastPositionMs int64     `gorm:"column:last_position_ms;default:0" json:"last_position_ms"` // 最后观看位置（毫秒），默认0
	LastPlayedAt   time.Time `gorm:"column:last_played_at" json:"last_played_at"`               // 最后播放时间
}

// AppVersion 应用版本模型
// 存储应用的版本信息和公告信息，通过Type字段区分（version:版本更新, announcement:公告）
type AppVersion struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;comment:版本ID" json:"id"`                                                // 版本ID，主键，自增
	Type          string    `gorm:"size:20;not null;default:'version';index;comment:类型(version:版本更新,announcement:公告)" json:"type"`  // 类型（version:版本更新, announcement:公告），索引
	VersionCode   int64     `gorm:"column:version_code;comment:版本号(数字)" json:"version_code"`                                        // 版本号（数字），公告类型可为空
	VersionName   string    `gorm:"column:version_name;size:50;comment:版本名称" json:"version_name"`                                   // 版本名称（如1.0.0），公告类型可为空
	Platform      string    `gorm:"size:20;not null;index:idx_app_versions_platform;comment:平台类型(android/ios/web)" json:"platform"` // 平台（android/ios/windows/macos/linux）
	Title         string    `gorm:"size:255;comment:公告标题" json:"title"`                                                             // 公告标题，版本更新类型可为空
	Content       string    `gorm:"type:text;comment:公告内容" json:"content"`                                                          // 公告内容，版本更新类型可为空
	DownloadURL   string    `gorm:"column:download_url;type:text;comment:下载链接" json:"download_url"`                                 // 下载地址
	UpdateContent string    `gorm:"column:update_content;type:text;comment:更新内容描述" json:"update_content"`                           // 更新内容
	IsForce       bool      `gorm:"column:is_force;default:0;comment:是否强制更新(0:否,1:是)" json:"is_force"`                              // 是否强制更新（1:是 0:否）
	FileSize      int64     `gorm:"column:file_size;comment:安装包大小(字节)" json:"file_size"`                                            // 文件大小（字节）
	SortOrder     int       `gorm:"column:sort_order;default:0;comment:排序顺序" json:"sort_order"`                                     // 排序顺序，数字越大越靠前
	IsActive      bool      `gorm:"column:is_active;default:1;index;comment:是否有效(0:无效,1:有效)" json:"is_active"`                      // 是否启用（1:是 0:否），索引
	CreatedAt     time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`                                                  // 创建时间，自动设置
	UpdatedAt     time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`                                                  // 更新时间，自动更新
}

// FilterInfo 筛选信息模型
// 存储视频筛选条件（国家、年份、标签等）
type FilterInfo struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"` // 筛选信息ID，主键，自增
	Name    string `gorm:"size:255" json:"name"`               // 名称
	Type    string `gorm:"size:255" json:"type"`               // 视频类型（movie/tv/anime/variety/doc）
	Country string `gorm:"size:255" json:"country"`            // 国家/地区
	Year    string `gorm:"size:255" json:"year"`               // 年份
	Tags    string `gorm:"size:255" json:"tags"`               // 标签
}

// TableName 指定表名
func (FilterInfo) TableName() string {
	return "filter_info"
}

// MembershipPlan 会员套餐模型
// 存储会员套餐配置信息（如VIP、SVIP等）
type MembershipPlan struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;comment:套餐ID" json:"id"`                                  // 套餐ID，主键，自增
	Name          string    `gorm:"size:100;not null;comment:套餐名称" json:"name"`                                       // 套餐名称（如：VIP、SVIP）
	Code          string    `gorm:"size:50;uniqueIndex;not null;comment:套餐代码" json:"code"`                            // 套餐代码（如：vip、svip），唯一索引
	Description   string    `gorm:"type:text;comment:套餐描述" json:"description"`                                        // 套餐描述
	Price         float64   `gorm:"type:decimal(10,2);not null;comment:价格（元）" json:"price"`                           // 价格（元）
	DurationDays  int       `gorm:"column:duration_days;not null;comment:有效期（天数）" json:"duration_days"`               // 有效期（天数）
	MaxResolution string    `gorm:"column:max_resolution;size:20;default:'720p';comment:最大清晰度" json:"max_resolution"` // 最大清晰度（720p/1080p/4K）
	CanSkipAd     bool      `gorm:"column:can_skip_ad;default:0;comment:是否可跳过广告" json:"can_skip_ad"`                  // 是否可跳过广告
	CanDownload   bool      `gorm:"column:can_download;default:0;comment:是否可下载" json:"can_download"`                  // 是否可下载
	MaxDevices    int       `gorm:"column:max_devices;default:1;comment:最大同时登录设备数" json:"max_devices"`                // 最大同时登录设备数
	IsActive      bool      `gorm:"column:is_active;default:1;comment:是否启用" json:"is_active"`                         // 是否启用
	SortOrder     int       `gorm:"column:sort_order;default:0;comment:排序顺序" json:"sort_order"`                       // 排序顺序
	CreatedAt     time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`                                    // 创建时间
	UpdatedAt     time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`                                    // 更新时间
}

// UserMembership 用户会员记录模型
// 记录用户的会员状态和有效期
type UserMembership struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:记录ID" json:"id"`                 // 记录ID，主键，自增
	UserID    int64     `gorm:"column:user_id;index;not null;comment:用户ID" json:"user_id"`       // 用户ID，索引
	PlanID    int64     `gorm:"column:plan_id;index;not null;comment:套餐ID" json:"plan_id"`       // 套餐ID，索引
	PlanCode  string    `gorm:"column:plan_code;size:50;not null;comment:套餐代码" json:"plan_code"` // 套餐代码（冗余字段，便于查询）
	StartTime time.Time `gorm:"column:start_time;not null;comment:开始时间" json:"start_time"`       // 开始时间
	EndTime   time.Time `gorm:"column:end_time;not null;index;comment:结束时间" json:"end_time"`     // 结束时间，索引
	IsActive  bool      `gorm:"column:is_active;default:1;index;comment:是否有效" json:"is_active"`  // 是否有效，索引
	AutoRenew bool      `gorm:"column:auto_renew;default:0;comment:是否自动续费" json:"auto_renew"`    // 是否自动续费
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`                   // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`                   // 更新时间
}

// MembershipOrder 会员订单模型
// 记录用户购买会员的订单信息
type MembershipOrder struct {
	ID            int64      `gorm:"primaryKey;autoIncrement;comment:订单ID" json:"id"`                                    // 订单ID，主键，自增
	OrderNo       string     `gorm:"column:order_no;size:64;uniqueIndex;not null;comment:订单号" json:"order_no"`           // 订单号，唯一索引
	UserID        int64      `gorm:"column:user_id;index;not null;comment:用户ID" json:"user_id"`                          // 用户ID，索引
	PlanID        int64      `gorm:"column:plan_id;not null;comment:套餐ID" json:"plan_id"`                                // 套餐ID
	PlanCode      string     `gorm:"column:plan_code;size:50;not null;comment:套餐代码" json:"plan_code"`                    // 套餐代码
	Amount        float64    `gorm:"type:decimal(10,2);not null;comment:订单金额（元）" json:"amount"`                          // 订单金额（元）
	PaymentMethod string     `gorm:"column:payment_method;size:50;comment:支付方式" json:"payment_method"`                   // 支付方式（alipay/wechat/other）
	PaymentStatus string     `gorm:"column:payment_status;size:20;default:'pending';comment:支付状态" json:"payment_status"` // 支付状态（pending/paid/failed/refunded）
	PaidAt        *time.Time `gorm:"column:paid_at;comment:支付时间" json:"paid_at"`                                         // 支付时间
	ExpiresAt     time.Time  `gorm:"column:expires_at;not null;comment:会员到期时间" json:"expires_at"`                        // 会员到期时间
	Remark        string     `gorm:"type:text;comment:备注" json:"remark"`                                                 // 备注
	CreatedAt     time.Time  `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`                                      // 创建时间
	UpdatedAt     time.Time  `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`                                      // 更新时间
}
