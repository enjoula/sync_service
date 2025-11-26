// database 包提供MySQL数据库连接和初始化功能
// 使用GORM作为ORM框架
package database

import (
	"fmt"
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
// 4. 执行自动迁移创建表结构
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
		zap.L().Fatal("mysql connect error", zap.Error(err))
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

	// 执行自动迁移，创建或更新表结构
	// 迁移所有定义的模型（以 migrations/init.sql 为准）
	if err := DB.AutoMigrate(
		&model.User{},
		&model.UserToken{},
		&model.Video{},
		&model.Episode{},
		&model.Danmaku{},
		&model.UserFavorite{},
		&model.FilterInfo{},
		&model.AppVersion{},
	); err != nil {
		zap.L().Error("auto migrate failed", zap.Error(err))
	} else {
		zap.L().Info("auto migrate applied")
	}

	// 添加表注释（GORM AutoMigrate 不会自动添加表注释）
	addTableComments()
}

// addTableComments 添加表注释
// GORM 的 AutoMigrate 不会自动添加表注释，需要手动执行 SQL 语句
func addTableComments() {
	tableComments := map[string]string{
		"users":          "用户表",
		"user_tokens":   "用户登录控制表",
		"videos":        "视频表",
		"episodes":      "剧集表",
		"danmakus":      "弹幕表",
		"user_favorites": "用户收藏表",
		"filter_info":   "视频表",
		"app_versions":  "应用版本表",
	}

	for tableName, comment := range tableComments {
		sql := fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", tableName, comment)
		if err := DB.Exec(sql).Error; err != nil {
			zap.L().Warn(fmt.Sprintf("failed to add comment for table %s", tableName), zap.Error(err))
		}
	}
	zap.L().Info("table comments applied")
}
