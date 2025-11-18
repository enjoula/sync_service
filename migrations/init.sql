/*
 Navicat Premium Dump SQL

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 80405 (8.4.5)
 Source Host           : localhost:5506
 Source Schema         : video_service

 Target Server Type    : MySQL
 Target Server Version : 80405 (8.4.5)
 File Encoding         : 65001

 Date: 19/11/2025 00:40:35
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app_versions
-- ----------------------------
DROP TABLE IF EXISTS `app_versions`;
CREATE TABLE `app_versions` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '版本ID',
  `version_code` bigint NOT NULL COMMENT '版本号(数字)',
  `version_name` varchar(50) NOT NULL COMMENT '版本名称',
  `platform` varchar(20) NOT NULL COMMENT '平台类型(android/ios/web)',
  `download_url` text COMMENT '下载链接',
  `update_content` text COMMENT '更新内容描述',
  `is_force` tinyint(1) DEFAULT '0' COMMENT '是否强制更新(0:否,1:是)',
  `file_size` bigint DEFAULT NULL COMMENT '安装包大小(字节)',
  `is_active` tinyint(1) DEFAULT '1' COMMENT '是否有效(0:无效,1:有效)',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_app_versions_platform_version` (`platform`,`version_code`),
  KEY `idx_app_versions_platform` (`platform`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='应用版本表';

-- ----------------------------
-- Table structure for danmakus
-- ----------------------------
DROP TABLE IF EXISTS `danmakus`;
CREATE TABLE `danmakus` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '弹幕ID',
  `episode_id` bigint NOT NULL COMMENT '所属剧集ID',
  `user_id` bigint DEFAULT NULL COMMENT '发送用户ID',
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
  `channel` varchar(255) DEFAULT NULL COMMENT '频道名称',
  `channel_id` bigint DEFAULT NULL COMMENT '频道ID',
  `video_id` bigint NOT NULL COMMENT '所属视频ID',
  `episode_number` bigint DEFAULT '1' COMMENT '集数编号',
  `name` varchar(255) DEFAULT NULL COMMENT '剧集名称',
  `play_urls` json NOT NULL COMMENT '播放地址列表(JSON格式)',
  `duration_seconds` bigint DEFAULT NULL COMMENT '时长(秒)',
  `subtitle_urls` json DEFAULT NULL COMMENT '字幕地址列表(JSON格式)',
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
  `token` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '登录令牌',
  `device` varchar(100) DEFAULT NULL COMMENT '设备信息',
  `ip_address` varchar(45) DEFAULT NULL COMMENT 'IP地址',
  `expires_at` datetime(3) DEFAULT NULL COMMENT '过期时间',
  `is_active` tinyint(1) DEFAULT '1' COMMENT '是否有效(0:无效,1:有效)',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_tokens_token` (`token`),
  UNIQUE KEY `token` (`token`),
  KEY `idx_user_tokens_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=71 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户登录控制表';

-- ----------------------------
-- Table structure for user_watch_progresses
-- ----------------------------
DROP TABLE IF EXISTS `user_watch_progresses`;
CREATE TABLE `user_watch_progresses` (
  `user_id` bigint NOT NULL COMMENT '用户ID，复合主键',
  `episode_id` bigint NOT NULL COMMENT '剧集ID，复合主键',
  `last_position_ms` bigint DEFAULT '0' COMMENT '最后播放位置(毫秒)',
  `last_played_at` datetime(3) DEFAULT NULL COMMENT '最后播放时间',
  PRIMARY KEY (`user_id`,`episode_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户观看进度表';

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint NOT NULL COMMENT '用户ID，使用算法生成（非自增）',
  `username` varchar(100) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码(加密存储)',
  `nickname` varchar(100) DEFAULT NULL COMMENT '昵称',
  `email` varchar(255) DEFAULT NULL COMMENT '邮箱地址',
  `avatar` text COMMENT '头像URL',
  `acc_web` varchar(255) DEFAULT NULL COMMENT 'Web端访问码',
  `acc_web_create_at` datetime(3) DEFAULT NULL COMMENT 'Web端访问码创建时间',
  `acc_tv` varchar(255) DEFAULT NULL COMMENT 'TV端访问码',
  `acc_tv_create_at` datetime(3) DEFAULT NULL COMMENT 'TV端访问码创建时间',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos` (
  `id` bigint NOT NULL COMMENT '视频ID，使用雪花算法生成（非自增主键）',
  `source_id` int DEFAULT NULL COMMENT '来源站点的视频ID',
  `source` varchar(255) DEFAULT NULL COMMENT '视频来源(如:douban)',
  `title` varchar(255) DEFAULT NULL COMMENT '视频标题',
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '视频类型(电影/电视剧/综艺等)',
  `cover_url` text COMMENT '封面图片URL',
  `description` text COMMENT '视频简介',
  `year` varchar(20) DEFAULT NULL COMMENT '上映年份',
  `rating` varchar(255) DEFAULT NULL COMMENT '评分',
  `country` varchar(50) DEFAULT NULL COMMENT '国家/地区',
  `director` varchar(255) DEFAULT NULL COMMENT '导演',
  `actors` varchar(500) DEFAULT NULL COMMENT '演员列表',
  `tags` varchar(255) DEFAULT NULL COMMENT '标签',
  `status` varchar(255) DEFAULT NULL COMMENT '状态',
  `imdb_id` varchar(20) DEFAULT NULL COMMENT 'IMDB编号',
  `runtime` int DEFAULT NULL COMMENT '时长(分钟)',
  `resolution` varchar(20) DEFAULT NULL COMMENT '分辨率',
  `episode_count` bigint DEFAULT NULL COMMENT '总集数',
  `is_completed` tinyint(1) DEFAULT '0' COMMENT '是否完结(0:未完结,1:已完结)',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='视频表';

SET FOREIGN_KEY_CHECKS = 1;
