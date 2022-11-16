/*
 Navicat Premium Data Transfer

 Source Server         : sanheRec
 Source Server Type    : MySQL
 Source Server Version : 50737
 Source Host           : 39.99.155.73:3306
 Source Schema         : sanheRec

 Target Server Type    : MySQL
 Target Server Version : 50737
 File Encoding         : 65001

 Date: 16/11/2022 20:34:55
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for articles
-- ----------------------------
DROP TABLE IF EXISTS `articles`;
CREATE TABLE `articles`  (
  `art_id` int(11) NOT NULL AUTO_INCREMENT,
  `career_job_id` int(11) NULL DEFAULT NULL,
  `job_label` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `boss_id` int(11) NULL DEFAULT NULL,
  `company_id` int(11) NULL DEFAULT NULL,
  `title` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `create_time` datetime(0) NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  `collect` int(11) NULL DEFAULT NULL,
  `view` int(11) NULL DEFAULT NULL,
  `delivery_num` int(11) NULL DEFAULT NULL,
  `weight` double NULL DEFAULT NULL,
  `status` int(11) NULL DEFAULT NULL,
  `region` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `salary_min` int(11) NULL DEFAULT NULL,
  `salary_max` int(11) NULL DEFAULT NULL,
  `show` int(11) NULL DEFAULT NULL,
  `recommend` int(11) NULL DEFAULT NULL,
  `art_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`art_id`) USING BTREE,
  INDEX `recommend`(`recommend`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 302 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for collections
-- ----------------------------
DROP TABLE IF EXISTS `collections`;
CREATE TABLE `collections`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `art_id` int(11) NULL DEFAULT NULL,
  `classification` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 232 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for companies
-- ----------------------------
DROP TABLE IF EXISTS `companies`;
CREATE TABLE `companies`  (
  `com_id` int(11) NOT NULL AUTO_INCREMENT,
  `pic_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `company_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `scale_tag` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `person_scale` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  `update_person` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `com_level` int(11) NULL DEFAULT NULL,
  `vip` int(11) NULL DEFAULT NULL,
  `phone` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `scale_level` int(11) NULL DEFAULT NULL,
  `com_status` int(11) NULL DEFAULT NULL,
  `applicant` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`com_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 31 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for connections
-- ----------------------------
DROP TABLE IF EXISTS `connections`;
CREATE TABLE `connections`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `des_person` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `connect` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `upload_time` datetime(0) NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 19 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for dailysavers
-- ----------------------------
DROP TABLE IF EXISTS `dailysavers`;
CREATE TABLE `dailysavers`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `art_id` int(11) NULL DEFAULT NULL,
  `day_view` int(11) NULL DEFAULT NULL,
  `day_delivery` int(11) NULL DEFAULT NULL,
  `date` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1303 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for deliveries
-- ----------------------------
DROP TABLE IF EXISTS `deliveries`;
CREATE TABLE `deliveries`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `boss_id` int(11) NULL DEFAULT NULL,
  `art_id` int(11) NULL DEFAULT NULL,
  `from_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `qualification` int(11) NULL DEFAULT NULL,
  `read` int(11) NULL DEFAULT NULL,
  `delivery_time` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 157 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for descriptions
-- ----------------------------
DROP TABLE IF EXISTS `descriptions`;
CREATE TABLE `descriptions`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `module` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  `upload_time` datetime(0) NULL DEFAULT NULL,
  `uploader` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for docks
-- ----------------------------
DROP TABLE IF EXISTS `docks`;
CREATE TABLE `docks`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `com_id` int(11) NULL DEFAULT NULL,
  `boss_id` int(11) NULL DEFAULT NULL,
  `user_id` int(11) NULL DEFAULT NULL,
  `art_id` int(11) NULL DEFAULT NULL,
  `dock_time` datetime(0) NULL DEFAULT NULL,
  `boss_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `pub_title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 455 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for educations
-- ----------------------------
DROP TABLE IF EXISTS `educations`;
CREATE TABLE `educations`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `school` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `major` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `start_time` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `end_time` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `degree` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `degree_level` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for invitations
-- ----------------------------
DROP TABLE IF EXISTS `invitations`;
CREATE TABLE `invitations`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `art_id` int(11) NULL DEFAULT NULL,
  `boss_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `des_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `invite_time` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 100 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for labels
-- ----------------------------
DROP TABLE IF EXISTS `labels`;
CREATE TABLE `labels`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `level` int(11) NULL DEFAULT NULL,
  `parent_id` int(11) NULL DEFAULT NULL,
  `label` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `recommend` int(11) NULL DEFAULT NULL,
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 314 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for mass_sends
-- ----------------------------
DROP TABLE IF EXISTS `mass_sends`;
CREATE TABLE `mass_sends`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `publish_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `publish_time` datetime(0) NULL DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `des_role` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1421 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for msgobjs
-- ----------------------------
DROP TABLE IF EXISTS `msgobjs`;
CREATE TABLE `msgobjs`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `from_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `to_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `msg_time` int(11) NULL DEFAULT NULL,
  `read` int(11) NULL DEFAULT NULL,
  `last_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `message_type` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 525 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for notices
-- ----------------------------
DROP TABLE IF EXISTS `notices`;
CREATE TABLE `notices`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `title` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `upload_time` datetime(0) NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  `uploader` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `recommend` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for open_pubids
-- ----------------------------
DROP TABLE IF EXISTS `open_pubids`;
CREATE TABLE `open_pubids`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `openid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `pubid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for propagandas
-- ----------------------------
DROP TABLE IF EXISTS `propagandas`;
CREATE TABLE `propagandas`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uploader` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `upload_time` datetime(0) NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  `type` int(11) NULL DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `title` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `recommend` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 49 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for trainers
-- ----------------------------
DROP TABLE IF EXISTS `trainers`;
CREATE TABLE `trainers`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `start_time` int(11) NULL DEFAULT NULL,
  `end_time` int(11) NULL DEFAULT NULL,
  `read` int(11) NULL DEFAULT NULL,
  `message_type` int(11) NULL DEFAULT NULL,
  `from_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `to_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 120884 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for upgrades
-- ----------------------------
DROP TABLE IF EXISTS `upgrades`;
CREATE TABLE `upgrades`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `company_id` int(11) NULL DEFAULT NULL,
  `target_level` int(11) NULL DEFAULT NULL,
  `qualification` int(11) NULL DEFAULT NULL,
  `from_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `apply_time` datetime(0) NULL DEFAULT NULL,
  `company_exist` int(11) NULL DEFAULT NULL,
  `show` int(11) NULL DEFAULT NULL,
  `time_id` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 103 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `company_id` int(11) NULL DEFAULT NULL,
  `head_pic` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `role` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `identy_pin` int(11) NULL DEFAULT NULL,
  `user_level` int(11) NULL DEFAULT NULL,
  `gender` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `resume` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `age` int(11) NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `phone` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `personal_skill` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `project_experience` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `personal_resume` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `intended_position` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `president` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `make_time` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`) USING BTREE,
  INDEX `username`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 72 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for vip_shows
-- ----------------------------
DROP TABLE IF EXISTS `vip_shows`;
CREATE TABLE `vip_shows`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `create_time` datetime(0) NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  `clicks` int(11) NULL DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `publisher` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `recommend` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 36 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for vouchers
-- ----------------------------
DROP TABLE IF EXISTS `vouchers`;
CREATE TABLE `vouchers`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `voucher_pic` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `company_id` int(11) NULL DEFAULT NULL,
  `time_id` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 93 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
