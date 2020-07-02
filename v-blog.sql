/*
 Navicat MySQL Data Transfer

 Source Server         : 47.97.196.203
 Source Server Type    : MySQL
 Source Server Version : 50721
 Source Host           : 47.97.196.203:3306
 Source Schema         : v-blog

 Target Server Type    : MySQL
 Target Server Version : 50721
 File Encoding         : 65001

 Date: 12/04/2020 19:05:19
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for article_tags
-- ----------------------------
DROP TABLE IF EXISTS `article_tags`;
CREATE TABLE `article_tags`  (
  `article_id` int(11) NOT NULL COMMENT '文章id',
  `tag_id` int(11) NOT NULL COMMENT '标签id',
  INDEX `index_article_id`(`article_id`) USING BTREE,
  INDEX `index_tag_id`(`tag_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of article_tags
-- ----------------------------
INSERT INTO `article_tags` VALUES (8, 2);
INSERT INTO `article_tags` VALUES (8, 3);
INSERT INTO `article_tags` VALUES (8, 4);
INSERT INTO `article_tags` VALUES (9, 5);
INSERT INTO `article_tags` VALUES (9, 6);
INSERT INTO `article_tags` VALUES (9, 7);

-- ----------------------------
-- Table structure for articles
-- ----------------------------
DROP TABLE IF EXISTS `articles`;
CREATE TABLE `articles`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文章标题',
  `head_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文章首图',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文章内容',
  `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文章简介',
  `category_id` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章分类id',
  `views` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章浏览量',
  `comment_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章评论数量',
  `is_draft` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否草稿 0-否；1-是',
  `published_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '发布时间',
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `deleted_at` datetime(0) NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章表\n' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of articles
-- ----------------------------
INSERT INTO `articles` VALUES (1, 'falksdfjl', 'klfajsldkf', 'kaskf', 'fasjdfkl', 1, 1, 1, 0, '2020-02-05 15:44:58', '2020-02-05 15:45:02', '2020-02-05 15:45:02', NULL);
INSERT INTO `articles` VALUES (2, '测试啊', 'fsdfasdf', 'fsadfasfasdf', 'fasfasdf', 1, 0, 0, 0, '2020-03-01 18:10:10', '2020-02-10 23:31:48', '2020-02-10 23:31:48', NULL);
INSERT INTO `articles` VALUES (3, '测试啊', 'fsdfasdf', 'fsadfasfasdf', 'fasfasdf', 1, 0, 0, 0, '2020-02-10 23:36:08', '2020-02-10 23:36:08', '2020-02-10 23:36:08', NULL);
INSERT INTO `articles` VALUES (4, '测试啊', 'fsdfasdf', 'fsadfasfasdf', 'fasfasdf', 2, 0, 0, 0, '2020-02-12 20:26:16', '2020-02-12 20:26:17', '2020-02-12 20:26:17', NULL);
INSERT INTO `articles` VALUES (5, '测试啊', 'fsdfasdf', 'fsadfasfasdf', 'fasfasdf', 2, 0, 0, 0, '2020-02-12 20:26:54', '2020-02-12 20:26:54', '2020-02-12 20:26:54', NULL);
INSERT INTO `articles` VALUES (6, '测试啊', 'fsdfasdf', 'fsadfasfasdf', 'fasfasdf', 2, 0, 0, 0, '2020-02-12 20:29:45', '2020-02-12 20:29:45', '2020-02-12 20:29:45', NULL);
INSERT INTO `articles` VALUES (7, '11xiugai', 'fasdf', 'fasf', 'fasdfsdfaf', 2, 0, 0, 0, '2020-02-20 08:00:00', '2020-02-12 20:32:35', '2020-02-13 17:40:14', NULL);
INSERT INTO `articles` VALUES (8, '测试啊', 'fsdfasdf', 'fsadfasfasdf', 'fasfasdf', 2, 0, 0, 0, '2020-02-12 20:36:25', '2020-02-12 20:36:25', '2020-02-12 20:36:25', NULL);
INSERT INTO `articles` VALUES (9, '看发家克里斯放假埃里克设计费蓝', 'fsdfasdf', 'fsadfasfasdf', 'fasfasdf', 2, 0, 0, 0, '2020-02-13 17:35:15', '2020-02-13 17:35:15', '2020-02-13 17:35:15', NULL);
INSERT INTO `articles` VALUES (10, '10xiugai', 'fasdf', 'fasf', 'fasdfsdfaf', 2, 0, 0, 0, '2020-02-20 08:00:00', '2020-02-13 17:37:26', '2020-02-13 17:43:29', NULL);
INSERT INTO `articles` VALUES (11, '11xiugai', 'fasdf', 'fasf', 'fasdfsdfaf', 2, 0, 0, 0, '2020-02-20 08:00:00', '2020-02-13 17:37:34', '2020-02-13 17:40:38', NULL);

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '分类名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '分类描述',
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `deleted_at` datetime(0) NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 23 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章分类表\n' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of categories
-- ----------------------------
INSERT INTO `categories` VALUES (1, 'php', '世界上最好的语言', '2020-02-05 21:43:51', '2020-02-11 22:49:48', '2020-02-11 22:49:48');
INSERT INTO `categories` VALUES (2, '测试修改1', '房间爱丽丝卡积分了卡', '2020-02-11 21:41:27', '2020-03-08 20:25:47', NULL);
INSERT INTO `categories` VALUES (3, 'js', '房间爱丽丝卡积分了卡', '2020-02-11 22:28:28', '2020-02-12 18:59:43', NULL);
INSERT INTO `categories` VALUES (4, 'java', '后端语言', '2020-03-06 23:02:41', '2020-03-06 23:02:41', NULL);
INSERT INTO `categories` VALUES (5, 'php', '后端语言', '2020-03-06 23:03:02', '2020-03-06 23:03:02', NULL);
INSERT INTO `categories` VALUES (6, '测试1', '后端语言', '2020-03-06 23:03:08', '2020-03-06 23:03:08', NULL);
INSERT INTO `categories` VALUES (7, '测试2', '后端语言', '2020-03-06 23:03:16', '2020-03-06 23:03:16', NULL);
INSERT INTO `categories` VALUES (8, '测试3', '后端语言', '2020-03-06 23:03:18', '2020-03-06 23:03:18', NULL);
INSERT INTO `categories` VALUES (9, '测试4', '后端语言', '2020-03-06 23:03:22', '2020-03-06 23:03:22', NULL);
INSERT INTO `categories` VALUES (10, '测试5', '后端语言', '2020-03-06 23:03:25', '2020-03-06 23:03:25', NULL);
INSERT INTO `categories` VALUES (11, '测试6', '后端语言', '2020-03-06 23:03:28', '2020-03-06 23:03:28', NULL);
INSERT INTO `categories` VALUES (12, '测试7', '后端语言', '2020-03-06 23:03:31', '2020-03-06 23:03:31', NULL);
INSERT INTO `categories` VALUES (13, '测试8', '后端语言', '2020-03-06 23:03:34', '2020-03-06 23:03:34', NULL);
INSERT INTO `categories` VALUES (14, '测试9', '后端语言', '2020-03-06 23:03:54', '2020-03-06 23:03:54', NULL);
INSERT INTO `categories` VALUES (15, '测试10', '后端语言', '2020-03-06 23:03:56', '2020-03-06 23:03:56', NULL);
INSERT INTO `categories` VALUES (16, '测试11', '后端语言', '2020-03-06 23:04:01', '2020-03-06 23:04:01', NULL);
INSERT INTO `categories` VALUES (17, '测试12', '后端语言', '2020-03-06 23:04:04', '2020-03-06 23:04:04', NULL);
INSERT INTO `categories` VALUES (18, '测试13', '后端语言', '2020-03-06 23:04:20', '2020-03-06 23:04:20', NULL);
INSERT INTO `categories` VALUES (19, '测试14', '后端语言', '2020-03-06 23:04:23', '2020-03-06 23:04:23', NULL);
INSERT INTO `categories` VALUES (20, 'asfasdfa', 'fasdfasfasdf', '2020-03-08 20:27:15', '2020-03-08 20:27:15', NULL);
INSERT INTO `categories` VALUES (21, '飞洒地方', '发射点发顺丰', '2020-03-08 20:27:25', '2020-03-08 23:04:13', '2020-03-08 23:04:14');
INSERT INTO `categories` VALUES (22, '鬼地方鬼地方', '公司豆腐干山豆根副食店广泛的', '2020-03-08 22:43:18', '2020-03-08 23:03:48', '2020-03-08 23:03:49');

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `article_id` int(11) UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章id',
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '评论者昵称',
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '评论者邮箱',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论内容',
  `state` tinyint(4) NOT NULL COMMENT '状态 0-待审核；1-审核通过；2-驳回',
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `deleted_at` datetime(0) NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章评论表\n' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for tags
-- ----------------------------
DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标签名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标签描述',
  `article_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '标签下的文章数量',
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `deleted_at` datetime(0) NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章标签表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of tags
-- ----------------------------
INSERT INTO `tags` VALUES (1, 'php', '世界上最好的语言', 0, '2020-02-05 21:44:07', '2020-02-12 20:04:37', '2020-02-12 20:04:37');
INSERT INTO `tags` VALUES (2, '后端1', '大概念', 0, '2020-02-05 21:44:17', '2020-03-08 23:28:14', NULL);
INSERT INTO `tags` VALUES (3, 'ps', 'falsjdflka', 0, '2020-02-12 19:55:25', '2020-02-12 20:36:25', NULL);
INSERT INTO `tags` VALUES (4, 'java', '第十六届法兰克世纪东方', 0, '2020-02-12 19:56:28', '2020-02-13 13:56:02', NULL);
INSERT INTO `tags` VALUES (5, 'lu', '第十六届法兰克世纪东方', 0, '2020-02-12 19:58:02', '2020-02-13 17:39:42', NULL);
INSERT INTO `tags` VALUES (6, 'js', '', 0, '2020-02-12 19:59:58', '2020-02-13 17:39:39', NULL);
INSERT INTO `tags` VALUES (7, 'haha', '第十六届法兰克世纪东方', 0, '2020-02-13 13:05:02', '2020-02-13 17:39:40', NULL);
INSERT INTO `tags` VALUES (8, 'fasdf', 'fasdfafds', 0, '2020-03-08 23:25:12', '2020-03-08 23:25:12', NULL);
INSERT INTO `tags` VALUES (9, '123', 'fasdfasdfasdefa', 0, '2020-03-08 23:26:40', '2020-03-08 23:26:40', NULL);
INSERT INTO `tags` VALUES (10, 'fsdfsd', 'fdffdf', 0, '2020-03-08 23:27:04', '2020-03-08 23:28:00', '2020-03-08 23:28:00');
INSERT INTO `tags` VALUES (11, 'fdsf', 'afsdfasdf', 0, '2020-03-08 23:27:39', '2020-03-08 23:27:48', '2020-03-08 23:27:49');
INSERT INTO `tags` VALUES (12, 'jjjj', '卢卡斯JFK了', 0, '2020-03-25 23:12:16', '2020-03-25 23:12:16', NULL);

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户昵称',
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户邮箱',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录密码',
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `deleted_at` datetime(0) NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表\n' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, '防火身份卡', 'abc@local.com', '1E3c+INIqKtnyTznK3Q1QMuv38sCTdbQakx3HS1wr64=', '2020-02-07 00:53:32', '2020-02-08 22:40:37', NULL);

SET FOREIGN_KEY_CHECKS = 1;
