// database 包提供MySQL数据库连接和初始化功能
// 使用GORM作为ORM框架
package database

import (
	"time"
	"video-service/internal/model"
	"video-service/pkg/infrastructure/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 是全局数据库连接对象
var DB *gorm.DB

// InitMySQL 初始化MySQL数据库连接
// 功能包括：
// 1. 从配置中获取数据库连接字符串（优先从Etcd获取）
// 2. 建立数据库连接
// 3. 配置连接池参数
// 注意：表结构需要通过 migrations/init.sql 手动初始化
func InitMySQL() {
	// 从配置中获取MySQL连接字符串
	dsn := config.Cfg.GetString("mysql.dsn")

	// 如果配置文件中没有，尝试从Etcd获取
	if dsn == "" {
		if config.Secrets != nil && config.Secrets.MySQLDsn != "" {
			dsn = config.Secrets.MySQLDsn
		}
	}

	// 如果仍然没有配置，记录警告并返回
	if dsn == "" {
		zap.L().Warn("mysql dsn empty")
		return
	}

	// 使用GORM打开MySQL连接
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("mysql connect error", zap.Error(err))
		return
	}

	// 获取底层sql.DB对象以配置连接池
	sqlDB, _ := dbConn.DB()
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(50)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置连接的最大生命周期（1小时）
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 保存全局数据库连接
	DB = dbConn

	zap.L().Info("mysql connected")

	// 注意：不再执行自动迁移
	// 表结构需要通过执行 migrations/init.sql 文件手动初始化
	// 这样可以更好地控制数据库结构变更，避免意外的表结构修改
	// 管理员表已存在时，补齐内置超级管理员（已有同名账号则跳过）
	ensureDefaultAdmin()
}

// 与 migrations/add_admins.sql 中的密码哈希保持一致
const (
	defaultAdminPasswordHash = "$2a$10$qxVMrYolsh3gyNiMKD5U4Ok3m90o1K1F0bO2nyyVZeTFalFFzK9mG" // admin / Admin@123456
	namedAdminPasswordHash   = "$2a$10$hdzjzv5Ggnv8pvQvTkAWxeZBc/YVIjNrK5fp9sOvj0lHv2hKNqDyK" // as7633859
)

func ensureDefaultAdmin() {
	if DB == nil {
		return
	}
	if err := DB.Model(&model.Admin{}).Limit(1).Count(new(int64)).Error; err != nil {
		zap.L().Warn("管理员表不可用，请执行 migrations/add_admins.sql", zap.Error(err))
		return
	}
	ensureAdmin("admin", defaultAdminPasswordHash, "超级管理员")
	ensureAdmin("as7633859", namedAdminPasswordHash, "as7633859")
}

func ensureAdmin(username, passwordHash, nickname string) {
	var count int64
	if err := DB.Model(&model.Admin{}).Where("username = ?", username).Count(&count).Error; err != nil {
		zap.L().Warn("查询管理员失败", zap.String("username", username), zap.Error(err))
		return
	}
	if count > 0 {
		return
	}
	adminUser := &model.Admin{
		Username: username,
		Password: passwordHash,
		Nickname: nickname,
		Role:     model.AdminRoleSuper,
		Status:   1,
	}
	if err := DB.Create(adminUser).Error; err != nil {
		zap.L().Warn("创建管理员失败", zap.String("username", username), zap.Error(err))
		return
	}
	zap.L().Info("已创建超级管理员", zap.String("username", username))
}

// addTableComments 添加表注释
// GORM 的 AutoMigrate 不会自动添加表注释，需要手动执行 SQL 语句

// func addTableComments() {
// 	tableComments := map[string]string{
// 		"users":             "用户表",
// 		"user_tokens":       "用户登录控制表",
// 		"videos":            "视频表",
// 		"episodes":          "剧集表",
// 		"danmakus":          "弹幕表",
// 		"user_favorites":    "用户收藏表",
// 		"filter_info":       "视频表",
// 		"app_versions":      "应用版本表",
// 		"membership_plans":  "会员套餐表",
// 		"user_memberships":  "用户会员记录表",
// 		"membership_orders": "会员订单表",
// 	}

// 	for tableName, comment := range tableComments {
// 		sql := fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", tableName, comment)
// 		if err := DB.Exec(sql).Error; err != nil {
// 			zap.L().Warn(fmt.Sprintf("failed to add comment for table %s", tableName), zap.Error(err))
// 		}
// 	}
// 	zap.L().Info("table comments applied")
