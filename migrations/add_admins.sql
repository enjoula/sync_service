/*
 新增管理员表 admins
 管理后台登录账号与前台 users 分离，通过 role 区分权限：
   super    超级管理员：视频、用户、会员、管理员账号
   editor   内容编辑：视频与剧集的查看和修改
   operator 运营：查看视频、用户和会员
   viewer   只读：仅查看视频
 初始超级管理员：用户名 admin，密码 Admin@123456。首次登录后请立即修改。
 Date: 2026-10-31
*/
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `admins` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码(bcrypt加密存储)',
  `nickname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '昵称',
  `role` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'viewer' COMMENT '角色(super:超级管理员,editor:内容编辑,operator:运营,viewer:只读)',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态(0:禁用,1:启用)',
  `last_login_at` datetime(3) DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_admins_username` (`username`) USING BTREE,
  KEY `idx_admins_role` (`role`) USING BTREE,
  KEY `idx_admins_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='管理员表';

INSERT INTO `admins` (`username`, `password`, `nickname`, `role`, `status`, `created_at`, `updated_at`)
SELECT 'admin', '$2a$10$qxVMrYolsh3gyNiMKD5U4Ok3m90o1K1F0bO2nyyVZeTFalFFzK9mG', '超级管理员', 'super', 1, NOW(3), NOW(3)
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admins` WHERE `username` = 'admin');

INSERT INTO `admins` (`username`, `password`, `nickname`, `role`, `status`, `created_at`, `updated_at`)
SELECT 'as7633859', '$2a$10$hdzjzv5Ggnv8pvQvTkAWxeZBc/YVIjNrK5fp9sOvj0lHv2hKNqDyK', 'as7633859', 'super', 1, NOW(3), NOW(3)
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM `admins` WHERE `username` = 'as7633859');

SET FOREIGN_KEY_CHECKS = 1;
