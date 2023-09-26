/*
 Navicat Premium Data Transfer

 Source Server         : Dev Eden ERP-Write
 Source Server Type    : MySQL
 Source Server Version : 80031 (8.0.31-google)
 Source Host           : 10.26.160.2:3306
 Source Schema         : settlement

 Target Server Type    : MySQL
 Target Server Version : 80031 (8.0.31-google)
 File Encoding         : 65001

 Date: 05/05/2023 12:28:44
*/
USE settlement;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sales_invoice_external
-- ----------------------------
DROP TABLE IF EXISTS `sales_invoice_external`;
CREATE TABLE `sales_invoice_external`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `sales_order_id` bigint UNSIGNED NULL DEFAULT NULL,
  `xendit_invoice_id` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `created_at` timestamp NULL DEFAULT NULL,
  `cancelled_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 233 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sales_invoice_external
-- ----------------------------
INSERT INTO `sales_invoice_external` VALUES (1, 56, '6411c223bda105790b234580', '2023-03-15 13:03:32', NULL);
INSERT INTO `sales_invoice_external` VALUES (2, 63, '6411ce19bda1053311234deb', '2023-03-15 13:54:34', NULL);
INSERT INTO `sales_invoice_external` VALUES (3, 68, '6411cf7d1d644638d61f07f4', '2023-03-15 14:00:30', NULL);
INSERT INTO `sales_invoice_external` VALUES (4, 73, '6411d13b1d6446f6d31f08ee', '2023-03-15 14:07:57', NULL);
INSERT INTO `sales_invoice_external` VALUES (5, 74, '6411d1734dce5e63be902c74', '2023-03-15 14:08:52', NULL);
INSERT INTO `sales_invoice_external` VALUES (6, 75, '6411d1cb4dce5e5833902cb9', '2023-03-15 14:10:21', NULL);
INSERT INTO `sales_invoice_external` VALUES (7, 76, '6411d232bda1056d732350d0', '2023-03-15 14:12:04', NULL);
INSERT INTO `sales_invoice_external` VALUES (8, 79, '6411d3434dce5e93e3902d9e', '2023-03-15 14:16:37', NULL);
INSERT INTO `sales_invoice_external` VALUES (9, 83, '6411d3d05558cd0a2362d037', '2023-03-15 14:18:57', NULL);
INSERT INTO `sales_invoice_external` VALUES (10, 87, '6411d4034dce5eab7f902e0f', '2023-03-15 14:19:49', NULL);
INSERT INTO `sales_invoice_external` VALUES (11, 92, '6411d4a44dce5eab1b902e75', '2023-03-15 14:22:29', NULL);
INSERT INTO `sales_invoice_external` VALUES (12, 95, '6411d5804dce5e6342902f14', '2023-03-15 14:26:10', NULL);
INSERT INTO `sales_invoice_external` VALUES (13, 96, '6411d6755558cd7cc662d1e8', '2023-03-15 14:30:14', NULL);
INSERT INTO `sales_invoice_external` VALUES (14, 100, '6411d7a8bda105257523545c', '2023-03-15 14:35:21', NULL);
INSERT INTO `sales_invoice_external` VALUES (15, 103, '6411d8085558cd3d9762d2f6', '2023-03-15 14:36:57', NULL);
INSERT INTO `sales_invoice_external` VALUES (16, 105, '6411d8741d6446f3a21f0da7', '2023-03-15 14:38:45', NULL);
INSERT INTO `sales_invoice_external` VALUES (17, 108, '6411d8fa4dce5e4687903171', '2023-03-15 14:40:59', NULL);
INSERT INTO `sales_invoice_external` VALUES (18, 110, '6411d933bda1052c3e23557e', '2023-03-15 14:41:57', NULL);
INSERT INTO `sales_invoice_external` VALUES (19, 114, '6411d9b0bda1051e7e2355d8', '2023-03-15 14:44:02', NULL);
INSERT INTO `sales_invoice_external` VALUES (20, 117, '6411da9a4dce5ecf53903266', '2023-03-15 14:47:55', NULL);
INSERT INTO `sales_invoice_external` VALUES (21, 119, '6411e573bda1058d8a235dce', '2023-03-15 15:34:12', NULL);
INSERT INTO `sales_invoice_external` VALUES (22, 123, '6411eff04dce5e4d45903fe0', '2023-03-15 16:18:57', NULL);
INSERT INTO `sales_invoice_external` VALUES (23, 124, '6411f18f4dce5e336b90411c', '2023-03-15 16:25:52', NULL);
INSERT INTO `sales_invoice_external` VALUES (24, 125, '6411f1cd5558cd377862e409', '2023-03-15 16:26:54', NULL);
INSERT INTO `sales_invoice_external` VALUES (25, 127, '6411f3661d6446d2c11f1ed8', '2023-03-15 16:33:43', NULL);
INSERT INTO `sales_invoice_external` VALUES (26, 128, '64121a721d6446218a1f3f72', '2023-03-15 19:20:19', NULL);
INSERT INTO `sales_invoice_external` VALUES (27, 129, '64121c671d6446044e1f40b8', '2023-03-15 19:28:40', NULL);
INSERT INTO `sales_invoice_external` VALUES (28, 130, '641229b7bda1057994239829', '2023-03-15 20:25:28', NULL);
INSERT INTO `sales_invoice_external` VALUES (29, 131, '64122a151d64467e551f4d55', '2023-03-15 20:27:02', NULL);
INSERT INTO `sales_invoice_external` VALUES (30, 132, '64122a711d6446111a1f4d8d', '2023-03-15 20:28:34', NULL);
INSERT INTO `sales_invoice_external` VALUES (31, 133, '64122aba5558cd15a16314b6', '2023-03-15 20:29:47', NULL);
INSERT INTO `sales_invoice_external` VALUES (32, 134, '64124b0d4dce5e8169908827', '2023-03-15 22:47:42', NULL);
INSERT INTO `sales_invoice_external` VALUES (33, 135, '64124cb2bda105919223b355', '2023-03-15 22:54:43', NULL);
INSERT INTO `sales_invoice_external` VALUES (34, 136, '64124de9bda1056fe523b47c', '2023-03-15 22:59:54', NULL);
INSERT INTO `sales_invoice_external` VALUES (35, 137, '64124e704dce5e11d2908b12', '2023-03-15 23:02:09', NULL);
INSERT INTO `sales_invoice_external` VALUES (36, 138, '64124f82bda105fd9f23b5c6', '2023-03-15 23:06:43', NULL);
INSERT INTO `sales_invoice_external` VALUES (37, 139, '641250b1bda1057ff023b6fd', '2023-03-15 23:11:46', NULL);
INSERT INTO `sales_invoice_external` VALUES (38, 140, '641251844dce5e2d4b908dcc', '2023-03-15 23:15:17', NULL);
INSERT INTO `sales_invoice_external` VALUES (39, 141, '64125243bda105307323b897', '2023-03-15 23:18:28', NULL);
INSERT INTO `sales_invoice_external` VALUES (40, 142, '641253635558cd7cf463346f', '2023-03-15 23:23:16', NULL);
INSERT INTO `sales_invoice_external` VALUES (41, 143, '64125478bda105c1cf23bab1', '2023-03-15 23:27:53', NULL);
INSERT INTO `sales_invoice_external` VALUES (42, 144, '641254d65558cd3b396335de', '2023-03-15 23:29:27', NULL);
INSERT INTO `sales_invoice_external` VALUES (43, 145, '641255935558cd7dca6336a2', '2023-03-15 23:32:36', NULL);
INSERT INTO `sales_invoice_external` VALUES (44, 146, '641255ef5558cd75bf633708', '2023-03-15 23:34:08', NULL);
INSERT INTO `sales_invoice_external` VALUES (45, 147, '641256575558cd5b6a633789', '2023-03-15 23:35:52', NULL);
INSERT INTO `sales_invoice_external` VALUES (46, 148, '641257564dce5e587e909381', '2023-03-15 23:40:07', NULL);
INSERT INTO `sales_invoice_external` VALUES (47, 149, '6412578d4dce5e4a2e9093c2', '2023-03-15 23:41:02', NULL);
INSERT INTO `sales_invoice_external` VALUES (48, 150, '641259035558cdb584633a1c', '2023-03-15 23:47:16', NULL);
INSERT INTO `sales_invoice_external` VALUES (49, 151, '641259424dce5e5e6c90953c', '2023-03-15 23:48:19', NULL);
INSERT INTO `sales_invoice_external` VALUES (50, 152, '64125a4b5558cd7c4f633b05', '2023-03-15 23:52:44', NULL);
INSERT INTO `sales_invoice_external` VALUES (51, 153, '64125aae1d644617431f73d6', '2023-03-15 23:54:23', NULL);
INSERT INTO `sales_invoice_external` VALUES (52, 154, '64125c175558cd3978633ca7', '2023-03-16 00:00:24', NULL);
INSERT INTO `sales_invoice_external` VALUES (53, 155, '64125c7b4dce5e190b909820', '2023-03-16 00:02:04', NULL);
INSERT INTO `sales_invoice_external` VALUES (54, 156, '64125cb15558cd1e0a633d43', '2023-03-16 00:02:58', NULL);
INSERT INTO `sales_invoice_external` VALUES (55, 157, '64125cf14dce5e530a909866', '2023-03-16 00:04:02', NULL);
INSERT INTO `sales_invoice_external` VALUES (56, 158, '64125d2cbda105269023c2e2', '2023-03-16 00:05:02', NULL);
INSERT INTO `sales_invoice_external` VALUES (57, 159, '64125d4a5558cd27c1633ddf', '2023-03-16 00:05:31', NULL);
INSERT INTO `sales_invoice_external` VALUES (58, 160, '64125dba4dce5e5e25909932', '2023-03-16 00:07:23', NULL);
INSERT INTO `sales_invoice_external` VALUES (59, 161, '64126020bda105c8f523c556', '2023-03-16 00:17:37', NULL);
INSERT INTO `sales_invoice_external` VALUES (60, 162, '6412602abda105038523c563', '2023-03-16 00:17:48', NULL);
INSERT INTO `sales_invoice_external` VALUES (61, 163, '64126093bda1056f8123c5b8', '2023-03-16 00:19:32', NULL);
INSERT INTO `sales_invoice_external` VALUES (62, 164, '641261b54dce5ea665909c3d', '2023-03-16 00:24:22', NULL);
INSERT INTO `sales_invoice_external` VALUES (63, 165, '641261dabda105835123c6b1', '2023-03-16 00:24:59', NULL);
INSERT INTO `sales_invoice_external` VALUES (64, 166, '641262db5558cd0482634243', '2023-03-16 00:29:16', NULL);
INSERT INTO `sales_invoice_external` VALUES (65, 167, '6412630b1d6446847b1f7ad0', '2023-03-16 00:30:04', NULL);
INSERT INTO `sales_invoice_external` VALUES (66, 168, '6412631c4dce5e86e2909d51', '2023-03-16 00:30:21', NULL);
INSERT INTO `sales_invoice_external` VALUES (67, 169, '641263315558cda5ac634296', '2023-03-16 00:30:42', NULL);
INSERT INTO `sales_invoice_external` VALUES (68, 170, '6412641a5558cd8e0c63436c', '2023-03-16 00:34:35', NULL);
INSERT INTO `sales_invoice_external` VALUES (69, 171, '64126486bda10566c723c911', '2023-03-16 00:36:23', NULL);
INSERT INTO `sales_invoice_external` VALUES (70, 172, '6412651d5558cd07f7634450', '2023-03-16 00:38:54', NULL);
INSERT INTO `sales_invoice_external` VALUES (71, 173, '641265d35558cd9ec6634509', '2023-03-16 00:41:56', NULL);
INSERT INTO `sales_invoice_external` VALUES (72, 174, '641266234dce5e2fff90a017', '2023-03-16 00:43:16', NULL);
INSERT INTO `sales_invoice_external` VALUES (73, 175, '641267101d64466e3d1f7e22', '2023-03-16 00:47:13', NULL);
INSERT INTO `sales_invoice_external` VALUES (74, 176, '6412681d5558cd32756346e7', '2023-03-16 00:51:42', NULL);
INSERT INTO `sales_invoice_external` VALUES (75, 177, '6412692a1d644633c21f8031', '2023-03-16 00:56:11', NULL);
INSERT INTO `sales_invoice_external` VALUES (76, 178, '64126a37bda105161723cdcb', '2023-03-16 01:00:41', NULL);
INSERT INTO `sales_invoice_external` VALUES (77, 179, '64126ab55558cd5f0b634911', '2023-03-16 01:02:46', NULL);
INSERT INTO `sales_invoice_external` VALUES (78, 180, '64126b8b1d644645d51f8251', '2023-03-16 01:06:21', NULL);
INSERT INTO `sales_invoice_external` VALUES (79, 181, '64126c571d64460c8b1f8305', '2023-03-16 01:09:44', NULL);
INSERT INTO `sales_invoice_external` VALUES (80, 182, '64126ca2bda105299023d00b', '2023-03-16 01:10:59', NULL);
INSERT INTO `sales_invoice_external` VALUES (81, 183, '64126e285558cd5e98634be6', '2023-03-16 01:17:29', NULL);
INSERT INTO `sales_invoice_external` VALUES (82, 184, '64126e461d644613f31f8490', '2023-03-16 01:18:00', NULL);
INSERT INTO `sales_invoice_external` VALUES (83, 185, '64126ec1bda1052c8523d1d7', '2023-03-16 01:20:02', NULL);
INSERT INTO `sales_invoice_external` VALUES (84, 186, '64126ef15558cd4610634ca3', '2023-03-16 01:20:50', NULL);
INSERT INTO `sales_invoice_external` VALUES (85, 187, '64126fb21d644643301f85e0', '2023-03-16 01:24:04', NULL);
INSERT INTO `sales_invoice_external` VALUES (86, 188, '6412704e1d644620c41f865e', '2023-03-16 01:26:39', NULL);
INSERT INTO `sales_invoice_external` VALUES (87, 189, '641271301d644602741f871d', '2023-03-16 01:30:26', NULL);
INSERT INTO `sales_invoice_external` VALUES (88, 190, '64127278bda105abcc23d512', '2023-03-16 01:35:53', NULL);
INSERT INTO `sales_invoice_external` VALUES (89, 191, '64127351bda105205623d609', '2023-03-16 01:39:30', NULL);
INSERT INTO `sales_invoice_external` VALUES (90, 192, '641273b01d644628f31f8988', '2023-03-16 01:41:06', NULL);
INSERT INTO `sales_invoice_external` VALUES (91, 193, '64127499bda105234d23d712', '2023-03-16 01:44:58', NULL);
INSERT INTO `sales_invoice_external` VALUES (92, 194, '641275371d64465ddf1f8a99', '2023-03-16 01:47:36', NULL);
INSERT INTO `sales_invoice_external` VALUES (93, 195, '6412757d4dce5e21ed90ad17', '2023-03-16 01:48:46', NULL);
INSERT INTO `sales_invoice_external` VALUES (94, 196, '641275a6bda105127223d7e4', '2023-03-16 01:49:27', NULL);
INSERT INTO `sales_invoice_external` VALUES (95, 197, '641276934dce5e457090ae40', '2023-03-16 01:53:24', NULL);
INSERT INTO `sales_invoice_external` VALUES (96, 198, '641276ed4dce5eba0690ae9a', '2023-03-16 01:54:54', NULL);
INSERT INTO `sales_invoice_external` VALUES (97, 199, '64127741bda105bfb223d995', '2023-03-16 01:56:18', NULL);
INSERT INTO `sales_invoice_external` VALUES (98, 200, '6412779a5558cdfd3f635484', '2023-03-16 01:57:47', NULL);
INSERT INTO `sales_invoice_external` VALUES (99, 201, '641278181d6446deb91f8d87', '2023-03-16 01:59:53', NULL);
INSERT INTO `sales_invoice_external` VALUES (100, 202, '64128037bda105662b23e889', '2023-03-16 02:34:32', NULL);
INSERT INTO `sales_invoice_external` VALUES (101, 203, '641288921d6446319a1fa8db', '2023-03-16 03:10:11', NULL);
INSERT INTO `sales_invoice_external` VALUES (102, 204, '641289d61d644623071fa9f1', '2023-03-16 03:15:35', NULL);
INSERT INTO `sales_invoice_external` VALUES (103, 205, '641292c35558cd716b637e01', '2023-03-16 03:53:40', NULL);
INSERT INTO `sales_invoice_external` VALUES (104, 206, '641292fdbda10537c4240407', '2023-03-16 03:54:38', NULL);
INSERT INTO `sales_invoice_external` VALUES (105, 207, '6412937e5558cddfd5637ebd', '2023-03-16 03:56:48', NULL);
INSERT INTO `sales_invoice_external` VALUES (106, 208, '641294225558cda678637f7e', '2023-03-16 03:59:31', NULL);
INSERT INTO `sales_invoice_external` VALUES (107, 209, '6412942f5558cd4106637f89', '2023-03-16 03:59:44', NULL);
INSERT INTO `sales_invoice_external` VALUES (108, 210, '641294565558cdf266637fb0', '2023-03-16 04:00:23', NULL);
INSERT INTO `sales_invoice_external` VALUES (109, 211, '6412bf4a5558cd0a2b63a8a4', '2023-03-16 07:03:39', NULL);
INSERT INTO `sales_invoice_external` VALUES (110, 212, '6412e29e95caa54fa963b038', '2023-03-16 09:34:23', NULL);
INSERT INTO `sales_invoice_external` VALUES (111, 213, '6413e11797daddf0fa77463e', '2023-03-17 03:40:09', NULL);
INSERT INTO `sales_invoice_external` VALUES (112, 220, '6417c138a4e305c21597b60b', '2023-03-20 02:13:13', NULL);
INSERT INTO `sales_invoice_external` VALUES (113, 221, '6417c33395caa5a34f685c95', '2023-03-20 02:21:40', NULL);
INSERT INTO `sales_invoice_external` VALUES (114, 222, '6417c581b872d67bf75acba0', '2023-03-20 02:31:30', NULL);
INSERT INTO `sales_invoice_external` VALUES (115, 223, '6417c7a0a4e3055e2b97bfda', '2023-03-20 02:40:33', NULL);
INSERT INTO `sales_invoice_external` VALUES (116, 224, '6417c7eda4e305450297c02f', '2023-03-20 02:41:50', NULL);
INSERT INTO `sales_invoice_external` VALUES (117, 225, '6417ca8ea4e305fd7797c2ae', '2023-03-20 02:53:03', NULL);
INSERT INTO `sales_invoice_external` VALUES (118, 226, '6417cb55a4e3057b7a97c361', '2023-03-20 02:56:23', NULL);
INSERT INTO `sales_invoice_external` VALUES (119, 227, '6417cbc697dadd41b17ad685', '2023-03-20 02:58:15', NULL);
INSERT INTO `sales_invoice_external` VALUES (120, 228, '6417cc1eb872d62a8f5ad1b4', '2023-03-20 02:59:44', NULL);
INSERT INTO `sales_invoice_external` VALUES (121, 229, '6417ce9cb872d6294d5ada27', '2023-03-20 03:10:21', NULL);
INSERT INTO `sales_invoice_external` VALUES (122, 230, '6417db04a4e305832197dc5f', '2023-03-20 04:03:17', NULL);
INSERT INTO `sales_invoice_external` VALUES (123, 231, '6417dba3a4e305c67097dd79', '2023-03-20 04:05:56', NULL);
INSERT INTO `sales_invoice_external` VALUES (124, 232, '6417dc2ab872d669875ae941', '2023-03-20 04:08:12', NULL);
INSERT INTO `sales_invoice_external` VALUES (125, 233, '6417dc5c95caa5227068819c', '2023-03-20 04:09:01', NULL);
INSERT INTO `sales_invoice_external` VALUES (126, 234, '6417dceb97dadd01957af093', '2023-03-20 04:11:24', NULL);
INSERT INTO `sales_invoice_external` VALUES (127, 235, '6417e3db97daddd85a7af6a8', '2023-03-20 04:41:00', NULL);
INSERT INTO `sales_invoice_external` VALUES (128, 236, '6417e48e97dadd68b37af751', '2023-03-20 04:43:59', NULL);
INSERT INTO `sales_invoice_external` VALUES (129, 237, '6417e8c597dadd5d127afb88', '2023-03-20 05:01:58', NULL);
INSERT INTO `sales_invoice_external` VALUES (130, 238, '6417f69697dadd02c27b084f', '2023-03-20 06:00:55', NULL);
INSERT INTO `sales_invoice_external` VALUES (131, 239, '6417fd8ba4e305cd6797fead', '2023-03-20 06:30:36', NULL);
INSERT INTO `sales_invoice_external` VALUES (132, 240, '64180a739eb9952d68d8a92d', '2023-03-20 07:25:40', NULL);
INSERT INTO `sales_invoice_external` VALUES (133, 241, '641814b80286c2009521eb42', '2023-03-20 08:09:30', NULL);
INSERT INTO `sales_invoice_external` VALUES (134, 242, '641818209eb995b026d8be70', '2023-03-20 08:24:01', NULL);
INSERT INTO `sales_invoice_external` VALUES (135, 243, '64181c660286c24da521f7d6', '2023-03-20 08:42:15', NULL);
INSERT INTO `sales_invoice_external` VALUES (136, 244, '64182d4f1a27d8e0987cface', '2023-03-20 09:54:25', NULL);
INSERT INTO `sales_invoice_external` VALUES (137, 245, '64182e0d0286c260b5220b12', '2023-03-20 09:57:34', NULL);
INSERT INTO `sales_invoice_external` VALUES (138, 246, '64182e769eb9955637d8d592', '2023-03-20 09:59:19', NULL);
INSERT INTO `sales_invoice_external` VALUES (139, 247, '641832f53ffa7a49f79e73d2', '2023-03-20 10:18:30', NULL);
INSERT INTO `sales_invoice_external` VALUES (140, 248, '641833121a27d8d4a97d001f', '2023-03-20 10:18:59', NULL);
INSERT INTO `sales_invoice_external` VALUES (141, 249, '641835749eb995e9abd8dc48', '2023-03-20 10:29:09', NULL);
INSERT INTO `sales_invoice_external` VALUES (142, 250, '641835cc3ffa7a1fe89e76a0', '2023-03-20 10:30:37', NULL);
INSERT INTO `sales_invoice_external` VALUES (143, 251, '6418362c0286c27d552213a2', '2023-03-20 10:32:13', NULL);
INSERT INTO `sales_invoice_external` VALUES (144, 252, '641836bb0286c2ea49221423', '2023-03-20 10:34:36', NULL);
INSERT INTO `sales_invoice_external` VALUES (145, 253, '641837390286c26a1422149d', '2023-03-20 10:36:43', NULL);
INSERT INTO `sales_invoice_external` VALUES (146, 255, '64183d433ffa7ae5399e7e14', '2023-03-20 11:02:28', NULL);
INSERT INTO `sales_invoice_external` VALUES (147, 256, '64183fa99eb9957e56d8e543', '2023-03-20 11:12:42', NULL);
INSERT INTO `sales_invoice_external` VALUES (148, 257, '64184d082bb1c47c3b2b1adb', '2023-03-20 12:09:45', NULL);
INSERT INTO `sales_invoice_external` VALUES (149, 258, '64184ed22bb1c454ba2b1c39', '2023-03-20 12:17:23', NULL);
INSERT INTO `sales_invoice_external` VALUES (150, 259, '641906c257b68e54e347cf5e', '2023-03-21 01:22:11', NULL);
INSERT INTO `sales_invoice_external` VALUES (151, 260, '64192f0ce17c6341f4f1a542', '2023-03-21 04:14:05', NULL);
INSERT INTO `sales_invoice_external` VALUES (152, 262, '641b08a257b68e109d4a59e4', '2023-03-22 13:54:43', NULL);
INSERT INTO `sales_invoice_external` VALUES (153, 266, '641bf97a2bb1c40cd92f2edc', '2023-03-23 07:02:20', NULL);
INSERT INTO `sales_invoice_external` VALUES (154, 267, '641bf9d02bb1c46a902f2f44', '2023-03-23 07:03:45', NULL);
INSERT INTO `sales_invoice_external` VALUES (155, 268, '641bfe2c2bb1c4a7b42f33dd', '2023-03-23 07:22:21', NULL);
INSERT INTO `sales_invoice_external` VALUES (156, 269, '641c06782bb1c471382f3c08', '2023-03-23 07:57:45', NULL);
INSERT INTO `sales_invoice_external` VALUES (157, 276, '641c29832bb1c49bc52f5df5', '2023-03-23 10:27:16', NULL);
INSERT INTO `sales_invoice_external` VALUES (158, 277, '641c29d157b68e3f684b5bf6', '2023-03-23 10:28:35', NULL);
INSERT INTO `sales_invoice_external` VALUES (159, 280, '641c2a3a2bb1c45e542f5e6e', '2023-03-23 10:30:19', NULL);
INSERT INTO `sales_invoice_external` VALUES (160, 281, '641c2a6d55c3159dee7d1bd6', '2023-03-23 10:31:10', NULL);
INSERT INTO `sales_invoice_external` VALUES (161, 285, '641c2ab12bb1c433212f5ed2', '2023-03-23 10:32:18', NULL);
INSERT INTO `sales_invoice_external` VALUES (162, 287, '641c2ac7e17c635b2af51c36', '2023-03-23 10:32:41', NULL);
INSERT INTO `sales_invoice_external` VALUES (163, 288, '641c2ad32bb1c4e1dc2f5eef', '2023-03-23 10:32:52', NULL);
INSERT INTO `sales_invoice_external` VALUES (164, 290, '641c2b322bb1c421822f5f53', '2023-03-23 10:34:27', NULL);
INSERT INTO `sales_invoice_external` VALUES (165, 292, '641c2b4ee17c63f6c4f51cba', '2023-03-23 10:34:55', NULL);
INSERT INTO `sales_invoice_external` VALUES (166, 293, '641c2b512bb1c421392f5f6b', '2023-03-23 10:34:59', NULL);
INSERT INTO `sales_invoice_external` VALUES (167, 295, '641c2b9f55c31583d27d1cff', '2023-03-23 10:36:16', NULL);
INSERT INTO `sales_invoice_external` VALUES (168, 296, '641c2c4d2bb1c4f0322f6052', '2023-03-23 10:39:10', NULL);
INSERT INTO `sales_invoice_external` VALUES (169, 297, '641c2c8155c315d0577d1dab', '2023-03-23 10:40:03', NULL);
INSERT INTO `sales_invoice_external` VALUES (170, 298, '641c2d7e2bb1c40fc72f6192', '2023-03-23 10:44:16', NULL);
INSERT INTO `sales_invoice_external` VALUES (171, 299, '641c2fdb57b68e6a644b60df', '2023-03-23 10:54:20', NULL);
INSERT INTO `sales_invoice_external` VALUES (172, 300, '641c3d5c57b68ec7484b6b45', '2023-03-23 11:51:57', NULL);
INSERT INTO `sales_invoice_external` VALUES (173, 301, '641c407b57b68e7da54b6d37', '2023-03-23 12:05:16', NULL);
INSERT INTO `sales_invoice_external` VALUES (174, 302, '641c4bfb57b68efaad4b7633', '2023-03-23 12:54:20', NULL);
INSERT INTO `sales_invoice_external` VALUES (175, 310, '641cfae72bb1c451fc3016fa', '2023-03-24 01:20:41', NULL);
INSERT INTO `sales_invoice_external` VALUES (176, 311, '641d0636e17c634422f5e323', '2023-03-24 02:08:55', NULL);
INSERT INTO `sales_invoice_external` VALUES (177, 312, '641d176d57b68ee8514c2d4b', '2023-03-24 03:22:22', NULL);
INSERT INTO `sales_invoice_external` VALUES (178, 313, '641d19b357b68e44ff4c2fdf', '2023-03-24 03:32:04', NULL);
INSERT INTO `sales_invoice_external` VALUES (179, 314, '641d1cc957b68e17334c3285', '2023-03-24 03:45:14', NULL);
INSERT INTO `sales_invoice_external` VALUES (180, 315, '641d28df2bb1c4c2b8304637', '2023-03-24 04:36:48', NULL);
INSERT INTO `sales_invoice_external` VALUES (181, 316, '641d4b3555c315b6537e19a7', '2023-03-24 07:03:18', NULL);
INSERT INTO `sales_invoice_external` VALUES (182, 317, '641d52c457b68e7d454c65cf', '2023-03-24 07:35:33', NULL);
INSERT INTO `sales_invoice_external` VALUES (183, 318, '641d56522bb1c4f9cd307348', '2023-03-24 07:50:44', NULL);
INSERT INTO `sales_invoice_external` VALUES (184, 319, '641d59c42bb1c41e5d307696', '2023-03-24 08:05:25', NULL);
INSERT INTO `sales_invoice_external` VALUES (185, 320, '641d5c39e17c635864f63af9', '2023-03-24 08:15:54', NULL);
INSERT INTO `sales_invoice_external` VALUES (186, 321, '641d5d2d57b68e64644c7056', '2023-03-24 08:19:58', NULL);
INSERT INTO `sales_invoice_external` VALUES (187, 322, '641d5d362bb1c47603307a19', '2023-03-24 08:20:07', NULL);
INSERT INTO `sales_invoice_external` VALUES (188, 323, '641d5d92e17c635959f63c85', '2023-03-24 08:21:40', NULL);
INSERT INTO `sales_invoice_external` VALUES (189, 324, '641d5fece17c63cec1f63f07', '2023-03-24 08:31:41', NULL);
INSERT INTO `sales_invoice_external` VALUES (190, 325, '641d602d57b68e67bb4c737a', '2023-03-24 08:32:46', NULL);
INSERT INTO `sales_invoice_external` VALUES (191, 326, '641d609355c31531bb7e2eb2', '2023-03-24 08:34:29', NULL);
INSERT INTO `sales_invoice_external` VALUES (192, 327, '641d61a457b68ec89d4c7498', '2023-03-24 08:39:01', NULL);
INSERT INTO `sales_invoice_external` VALUES (193, 328, '641d66f7e17c632b01f64666', '2023-03-24 09:01:44', NULL);
INSERT INTO `sales_invoice_external` VALUES (194, 329, '641d6f18e17c633d3df64e37', '2023-03-24 09:36:25', NULL);
INSERT INTO `sales_invoice_external` VALUES (195, 330, '641d6f282bb1c4b2e1308b5d', '2023-03-24 09:36:41', NULL);
INSERT INTO `sales_invoice_external` VALUES (196, 331, '641d6f5457b68e7d614c815b', '2023-03-24 09:37:25', NULL);
INSERT INTO `sales_invoice_external` VALUES (197, 333, '641d74a155c315931e7e410d', '2023-03-24 10:00:02', NULL);
INSERT INTO `sales_invoice_external` VALUES (198, 334, '641d74e52bb1c410963090da', '2023-03-24 10:01:10', NULL);
INSERT INTO `sales_invoice_external` VALUES (199, 335, '64212a3eb326425b5cb70ad7', '2023-03-27 05:31:43', NULL);
INSERT INTO `sales_invoice_external` VALUES (200, 336, '642139c9014272f75c3d04ce', '2023-03-27 06:38:02', NULL);
INSERT INTO `sales_invoice_external` VALUES (201, 337, '64213cf25967e2a632f23277', '2023-03-27 06:51:31', NULL);
INSERT INTO `sales_invoice_external` VALUES (202, 338, '6423b7d07c4e478ddd018da0', '2023-03-29 04:00:17', NULL);
INSERT INTO `sales_invoice_external` VALUES (203, 339, '6423c34aefa902759f500f40', '2023-03-29 04:49:15', NULL);
INSERT INTO `sales_invoice_external` VALUES (204, 340, '6423cc50b5cb303c74b3537e', '2023-03-29 05:27:45', NULL);
INSERT INTO `sales_invoice_external` VALUES (205, 341, '6423d1feb5cb3017e8b35bd7', '2023-03-29 05:51:59', NULL);
INSERT INTO `sales_invoice_external` VALUES (206, 342, '6423f71aefa9029078505fd6', '2023-03-29 08:30:19', NULL);
INSERT INTO `sales_invoice_external` VALUES (207, 343, '6424f6a1b5cb3071a5b507b7', '2023-03-30 02:40:34', NULL);
INSERT INTO `sales_invoice_external` VALUES (208, 344, '642514e0efa902b1b251ff29', '2023-03-30 04:49:37', NULL);
INSERT INTO `sales_invoice_external` VALUES (209, 345, '642519a9b5cb302a04b52ff5', '2023-03-30 05:10:02', NULL);
INSERT INTO `sales_invoice_external` VALUES (210, 346, '64251a211b451015927323dd', '2023-03-30 05:12:02', NULL);
INSERT INTO `sales_invoice_external` VALUES (211, 347, '64251a6db5cb300fccb530b2', '2023-03-30 05:13:18', NULL);
INSERT INTO `sales_invoice_external` VALUES (212, 348, '64251a941b4510df92732441', '2023-03-30 05:13:57', NULL);
INSERT INTO `sales_invoice_external` VALUES (213, 349, '64251a9c7c4e47df29038416', '2023-03-30 05:14:05', NULL);
INSERT INTO `sales_invoice_external` VALUES (214, 350, '64251aa4efa9020f98520418', '2023-03-30 05:14:13', NULL);
INSERT INTO `sales_invoice_external` VALUES (215, 351, '64251abbefa902eaba520428', '2023-03-30 05:14:36', NULL);
INSERT INTO `sales_invoice_external` VALUES (216, 352, '64252817b5cb30a490b53cf5', '2023-03-30 06:11:37', NULL);
INSERT INTO `sales_invoice_external` VALUES (217, 353, '64252fb6b5cb3039a6b544de', '2023-03-30 06:44:07', NULL);
INSERT INTO `sales_invoice_external` VALUES (218, 354, '64253195b5cb302cf7b54713', '2023-03-30 06:52:07', NULL);
INSERT INTO `sales_invoice_external` VALUES (219, 355, '642531d07c4e4718d6039a8f', '2023-03-30 06:53:05', NULL);
INSERT INTO `sales_invoice_external` VALUES (220, 356, '642536197c4e470a0e039e2c', '2023-03-30 07:11:22', NULL);
INSERT INTO `sales_invoice_external` VALUES (221, 357, '64254f8c1b4510f6987359d8', '2023-03-30 08:59:57', NULL);
INSERT INTO `sales_invoice_external` VALUES (222, 358, '64256243efa902d453524a2f', '2023-03-30 10:19:48', NULL);
INSERT INTO `sales_invoice_external` VALUES (223, 361, '642a8cf7a0f58330e02133a8', '2023-04-03 08:23:21', NULL);
INSERT INTO `sales_invoice_external` VALUES (224, 362, '642a8e76f750fd3570d63cc5', '2023-04-03 08:29:43', NULL);
INSERT INTO `sales_invoice_external` VALUES (225, 364, '642a983ca0f583e231214017', '2023-04-03 09:11:26', NULL);
INSERT INTO `sales_invoice_external` VALUES (226, 366, '642aebe3a565a7785008f3c2', '2023-04-03 15:08:20', NULL);
INSERT INTO `sales_invoice_external` VALUES (227, 367, '642b89d2a565a7f1930978b1', '2023-04-04 02:22:11', NULL);
INSERT INTO `sales_invoice_external` VALUES (228, 373, '6433787e7dfc9c29a36e0ae4', '2023-04-10 02:46:23', NULL);
INSERT INTO `sales_invoice_external` VALUES (229, 374, '64337b417dfc9c988c6e0d3b', '2023-04-10 02:58:10', NULL);
INSERT INTO `sales_invoice_external` VALUES (230, 375, '64338c18ffd2b96618b47e76', '2023-04-10 04:10:01', NULL);
INSERT INTO `sales_invoice_external` VALUES (231, 376, '643398f902f2fccaaff7fa03', '2023-04-10 05:04:59', NULL);
INSERT INTO `sales_invoice_external` VALUES (232, 379, '64376f0b92279a03d53207a0', '2023-04-13 02:55:09', NULL);

-- ----------------------------
-- Table structure for schema_migrations
-- ----------------------------
DROP TABLE IF EXISTS `schema_migrations`;
CREATE TABLE `schema_migrations`  (
  `version` bigint NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of schema_migrations
-- ----------------------------
INSERT INTO `schema_migrations` VALUES (2, 0);

-- ----------------------------
-- Table structure for txn_xendit_fva
-- ----------------------------
DROP TABLE IF EXISTS `txn_xendit_fva`;
CREATE TABLE `txn_xendit_fva`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `external_id` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `paid_amount` decimal(20, 2) NULL DEFAULT 0.00,
  `payment_channel` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `account_number` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `transaction_time` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of txn_xendit_fva
-- ----------------------------
INSERT INTO `txn_xendit_fva` VALUES (2, 'PERMATA_FVA-ECXHMB3318', 123321.00, 'PERMATA', '82149999446180', '2023-05-03 08:11:36', '2023-05-03 08:11:40');

-- ----------------------------
-- Table structure for txn_xendit_iva
-- ----------------------------
DROP TABLE IF EXISTS `txn_xendit_iva`;
CREATE TABLE `txn_xendit_iva`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `xendit_invoice_id` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `external_id` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `paid_amount` decimal(20, 2) NULL DEFAULT 0.00,
  `created_so_status` tinyint(1) NULL DEFAULT 0 COMMENT '1_created 2_not_yet_created',
  `status` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `payment_method` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `payment_channel` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `account_number` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `transaction_time` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of txn_xendit_iva
-- ----------------------------
INSERT INTO `txn_xendit_iva` VALUES (1, '579c8d61f23fa4ca35e52da4', 'invoice_123124123', 50000.00, 2, 'PAID', 'BANK_TRANSFER', 'PERMATA', '888888888888', '2016-10-12 08:15:03', '2023-05-04 10:35:01');

SET FOREIGN_KEY_CHECKS = 1;
