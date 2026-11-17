-- 用户反馈表
-- C 端用户反馈视频/剧集问题、缺少视频，以及搜索无结果记录，供管理后台查看并修正
-- source 字段区分：feedback=主动反馈，search=搜索无结果
-- 启动时 AutoMigrate 也会创建/更新该表；本脚本可用于手动建表或存量环境

CREATE TABLE IF NOT EXISTS `feedbacks` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '反馈ID',
  `user_id` bigint NOT NULL COMMENT '反馈/搜索用户ID',
  `source` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'feedback' COMMENT '来源方式(feedback:主动反馈,search:搜索无结果)',
  `video_id` bigint DEFAULT NULL COMMENT '视频ID(缺少视频/搜索无结果时可空)',
  `episode_id` bigint DEFAULT NULL COMMENT '剧集ID(可选)',
  `episode_number` bigint DEFAULT NULL COMMENT '集数编号(可选)',
  `video_title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '视频标题快照或缺失/搜索视频名称',
  `search_keyword` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '搜索关键词(source=search时填写)',
  `problem_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '问题类型(play_error/wrong_content/missing_episode/missing_video/subtitle_error/audio_error/other)',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '问题描述或搜索说明',
  `contact` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '联系方式(可选)',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'pending' COMMENT '处理状态(pending/processing/resolved/rejected/closed)',
  `admin_remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '管理员备注',
  `handled_by` bigint DEFAULT NULL COMMENT '处理人管理员ID',
  `handled_at` datetime(3) DEFAULT NULL COMMENT '处理时间',
  `client_device` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '客户端设备信息',
  `client_ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '客户端IP',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间(反馈/搜索时间)',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_feedbacks_user_created` (`user_id`,`created_at`) USING BTREE,
  KEY `idx_feedbacks_video_id` (`video_id`) USING BTREE,
  KEY `idx_feedbacks_status_created` (`status`,`created_at`) USING BTREE,
  KEY `idx_feedbacks_type_status` (`problem_type`,`status`) USING BTREE,
  KEY `idx_feedbacks_source_created` (`source`,`created_at`) USING BTREE,
  KEY `idx_feedbacks_search_keyword` (`search_keyword`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户反馈表';

-- 若此前已创建旧表 video_feedbacks，可执行：
-- RENAME TABLE `video_feedbacks` TO `feedbacks`;
