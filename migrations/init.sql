/*
 Navicat Premium Dump SQL

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 80407 (8.4.7)
 Source Host           : localhost:5506
 Source Schema         : video_service

 Target Server Type    : MySQL
 Target Server Version : 80407 (8.4.7)
 File Encoding         : 65001

 Date: 11/11/2025 14:34:46
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for danmakus
-- ----------------------------
DROP TABLE IF EXISTS `danmakus`;
CREATE TABLE `danmakus` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '弹幕ID',
  `episode_id` bigint NOT NULL COMMENT '剧集ID',
  `user_id` bigint DEFAULT NULL COMMENT '用户ID',
  `content` varchar(255) NOT NULL COMMENT '弹幕内容',
  `time_ms` bigint NOT NULL COMMENT '弹幕出现时间(毫秒)',
  `color` varchar(20) DEFAULT '#FFFFFF' COMMENT '弹幕颜色',
  `font_size` bigint DEFAULT '16' COMMENT '字体大小',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_danmakus_episode_id` (`episode_id`),
  KEY `idx_danmakus_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='弹幕表';

-- ----------------------------
-- Table structure for episodes
-- ----------------------------
DROP TABLE IF EXISTS `episodes`;
CREATE TABLE `episodes` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '剧集ID',
  `channel` varchar(255) DEFAULT NULL COMMENT '同步渠道',
  `channel_id` int DEFAULT NULL COMMENT '渠道视频ID',
  `video_id` bigint NOT NULL COMMENT '视频ID',
  `episode_number` bigint DEFAULT '1' COMMENT '第几集',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '剧集名称',
  `play_urls` json NOT NULL COMMENT '播放地址(JSON数组)',
  `duration_seconds` bigint DEFAULT NULL COMMENT '时长(秒)',
  `subtitle_urls` json DEFAULT NULL COMMENT '字幕地址(JSON数组)',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_episodes_video_id` (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='剧集表';

-- ----------------------------
-- Table structure for user_favorites
-- ----------------------------
DROP TABLE IF EXISTS `user_favorites`;
CREATE TABLE `user_favorites` (
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `video_id` bigint NOT NULL COMMENT '视频ID',
  `created_at` datetime(3) DEFAULT NULL COMMENT '收藏时间',
  PRIMARY KEY (`user_id`,`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户收藏表';

-- ----------------------------
-- Table structure for user_tokens
-- ----------------------------
DROP TABLE IF EXISTS `user_tokens`;
CREATE TABLE `user_tokens` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '令牌ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `token` varchar(255) NOT NULL COMMENT '令牌字符串',
  `device` varchar(100) DEFAULT NULL COMMENT '设备信息',
  `ip_address` varchar(45) DEFAULT NULL COMMENT 'IP地址',
  `expires_at` datetime(3) DEFAULT NULL COMMENT '过期时间',
  `is_active` tinyint(1) DEFAULT '1' COMMENT '是否激活(1:是 0:否)',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_tokens_token` (`token`),
  KEY `idx_user_tokens_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=71 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户登录控制表';

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(100) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码(加密)',
  `nickname` varchar(100) DEFAULT NULL COMMENT '昵称',
  `email` varchar(255) DEFAULT NULL COMMENT '邮箱',
  `avatar` text COMMENT '头像URL',
  `acc_web` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'Web授权码',
  `acc_web_create_at` datetime(3) DEFAULT NULL COMMENT 'Web授权码创建时间',
  `acc_tv` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'TV授权码',
  `acc_tv_create_at` datetime(3) DEFAULT NULL COMMENT 'TV授权码创建时间',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=246149288072368129 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '视频ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '视频标题',
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '类型(电影、电视剧、综艺、动漫、纪录片)',
  `cover_url` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '封面图片URL',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '视频描述',
  `year` bigint DEFAULT NULL COMMENT '上映年份',
  `rating` varchar(255) DEFAULT NULL COMMENT '评分',
  `country` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '制作国家/地区',
  `director` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '导演',
  `actors` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '主演(多个用逗号分隔)',
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '标签(多个用逗号分隔)',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频表';

-- ----------------------------
-- Table structure for app_versions
-- ----------------------------
DROP TABLE IF EXISTS `app_versions`;
CREATE TABLE `app_versions` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '版本ID',
  `version_code` bigint NOT NULL COMMENT '版本号(数字)',
  `version_name` varchar(50) NOT NULL COMMENT '版本名称(如1.0.0)',
  `platform` varchar(20) NOT NULL COMMENT '平台(android/ios/windows/macos/linux)',
  `download_url` text COMMENT '下载地址',
  `update_content` text COMMENT '更新内容',
  `is_force` tinyint(1) DEFAULT '0' COMMENT '是否强制更新(1:是 0:否)',
  `file_size` bigint DEFAULT NULL COMMENT '文件大小(字节)',
  `is_active` tinyint(1) DEFAULT '1' COMMENT '是否启用(1:是 0:否)',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_app_versions_platform_version` (`platform`,`version_code`),
  KEY `idx_app_versions_platform` (`platform`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='应用版本表';

SET FOREIGN_KEY_CHECKS = 1;
