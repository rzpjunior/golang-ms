/*
 Navicat Premium Data Transfer

 Source Server         : Dev Eden ERP-Write
 Source Server Type    : MySQL
 Source Server Version : 80031 (8.0.31-google)
 Source Host           : 10.26.160.2:3306
 Source Schema         : sales

 Target Server Type    : MySQL
 Target Server Version : 80031 (8.0.31-google)
 File Encoding         : 65001

 Date: 05/05/2023 12:28:35
*/
USE sales;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for payment_channel
-- ----------------------------
DROP TABLE IF EXISTS `payment_channel`;
CREATE TABLE `payment_channel`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `payment_method_id` bigint UNSIGNED NULL DEFAULT NULL,
  `code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `value` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `image_url` varchar(300) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `payment_guide_url` varchar(300) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `note` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `status` tinyint(1) NULL DEFAULT 0,
  `publish_iva` tinyint(1) NULL DEFAULT 0,
  `publish_fva` tinyint(1) NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of payment_channel
-- ----------------------------
INSERT INTO `payment_channel` VALUES (1, 2, 'PYC0001', 'BCA', 'BCA Virtual Account', 'https://sgp1.digitaloceanspaces.com/image-prod-eden/image/payment_channel_bca.png', 'https://www.edenfarm.id/payment-instruction-bca-bank', '', 1, 1, 2);
INSERT INTO `payment_channel` VALUES (2, 2, 'PYC0002', 'BRI', 'BRI Virtual Account', 'https://sgp1.digitaloceanspaces.com/image-prod-eden/image/payment_channel_bri.png', 'https://www.edenfarm.id/payment-instruction-bri-bank', '', 1, 1, 2);
INSERT INTO `payment_channel` VALUES (3, 2, 'PYC0003', 'MANDIRI', 'Mandiri Virtual Account', 'https://sgp1.digitaloceanspaces.com/image-prod-eden/image/payment_channel_mandiri.png', 'https://www.edenfarm.id/payment-instruction-mandiri-bank', '', 1, 1, 2);
INSERT INTO `payment_channel` VALUES (4, 2, 'PYC0004', 'BNI', 'BNI Virtual Account', 'https://sgp1.digitaloceanspaces.com/image-prod-eden/image/payment_channel_bni.png', 'https://www.edenfarm.id/payment-instruction-bni-bank', '', 1, 1, 2);
INSERT INTO `payment_channel` VALUES (5, 2, 'PYC0005', 'PERMATA', 'Permata Virtual Account', 'https://sgp1.digitaloceanspaces.com/image-prod-eden/image/payment_channel_permata.png', 'https://www.edenfarm.id/payment-instruction-permata-bank', '', 1, 1, 2);
INSERT INTO `payment_channel` VALUES (6, 2, 'PYC0006', 'BCA_FVA', 'BCA Virtual Account', 'https://sgp1.digitaloceanspaces.com/image-prod-eden/image/payment_channel_bca.png', 'https://www.edenfarm.id/payment-instruction-bca-bank', '', 1, 2, 1);
INSERT INTO `payment_channel` VALUES (7, 2, 'PYC0007', 'PERMATA_FVA', 'Permata Virtual Account', 'https://sgp1.digitaloceanspaces.com/image-prod-eden/image/payment_channel_permata.png', 'https://www.edenfarm.id/payment-instruction-permata-bank', '', 1, 2, 1);
INSERT INTO `payment_channel` VALUES (8, 2, 'PYC0008', 'SAHABAT_SAMPOERNA', 'Sahabat Sampoerna Virtual Account', 'https://sgp1.digitaloceanspaces.com/image-prod-eden/image/payment_channel_sampoerna.png', 'https://www.edenfarm.id/payment-instruction-sampoerna', '', 2, 1, 2);

-- ----------------------------
-- Table structure for payment_group_comb
-- ----------------------------
DROP TABLE IF EXISTS `payment_group_comb`;
CREATE TABLE `payment_group_comb`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `payment_group_sls` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `term_payment_sls` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of payment_group_comb
-- ----------------------------
INSERT INTO `payment_group_comb` VALUES (1, 'Advance', 'PBD');
INSERT INTO `payment_group_comb` VALUES (2, 'on-Delivery', 'COD');
INSERT INTO `payment_group_comb` VALUES (3, 'on-Delivery', 'BNS');
INSERT INTO `payment_group_comb` VALUES (4, 'on-Delivery', 'PWD');
INSERT INTO `payment_group_comb` VALUES (5, 'in-Term', '1 day(s)');
INSERT INTO `payment_group_comb` VALUES (6, 'in-Term', '7 day(s)');
INSERT INTO `payment_group_comb` VALUES (7, 'in-Term', '14 day(s)');
INSERT INTO `payment_group_comb` VALUES (8, 'in-Term', '30 day(s)');
INSERT INTO `payment_group_comb` VALUES (9, 'in-Term', '90 day(s)');
INSERT INTO `payment_group_comb` VALUES (10, 'in-Term', '21 Day(s)');
INSERT INTO `payment_group_comb` VALUES (11, 'in-Term', '3 day(s)');
INSERT INTO `payment_group_comb` VALUES (12, 'in-Term', '45 day(s)');
INSERT INTO `payment_group_comb` VALUES (13, 'in-Term', '60 day(s)');

-- ----------------------------
-- Table structure for payment_method
-- ----------------------------
DROP TABLE IF EXISTS `payment_method`;
CREATE TABLE `payment_method`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `note` varchar(250) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `status` tinyint(1) NULL DEFAULT 0,
  `publish` tinyint(1) NULL DEFAULT 0,
  `maintenance` tinyint(1) NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of payment_method
-- ----------------------------
INSERT INTO `payment_method` VALUES (1, 'PYM0001', 'Cash', '', 1, 1, 1);
INSERT INTO `payment_method` VALUES (2, 'PYM0002', 'Bank Transfer', '', 1, 1, 1);
INSERT INTO `payment_method` VALUES (3, 'PYM0003', 'Giro', '', 1, 1, 1);
INSERT INTO `payment_method` VALUES (4, 'PYM0004', 'Check', '', 1, 1, 1);
INSERT INTO `payment_method` VALUES (5, 'PYM0005', 'Other', '', 1, 1, 1);
INSERT INTO `payment_method` VALUES (6, 'PYM0006', 'Bad debt', '', 1, 1, 1);
INSERT INTO `payment_method` VALUES (7, 'PYM0007', 'Deposit', '', 1, 1, 1);

-- ----------------------------
-- Table structure for sales_order
-- ----------------------------
DROP TABLE IF EXISTS `sales_order`;
CREATE TABLE `sales_order`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `address_id` bigint UNSIGNED NULL DEFAULT NULL,
  `customer_id` bigint UNSIGNED NULL DEFAULT NULL,
  `term_payment_sls_id` bigint UNSIGNED NULL DEFAULT NULL,
  `term_invoice_sls_id` bigint UNSIGNED NULL DEFAULT NULL,
  `salesperson_id` bigint UNSIGNED NULL DEFAULT NULL,
  `sales_group_id` bigint UNSIGNED NULL DEFAULT NULL,
  `sub_district_id` bigint UNSIGNED NULL DEFAULT NULL,
  `site_id` bigint UNSIGNED NULL DEFAULT NULL,
  `wrt_id` bigint UNSIGNED NULL DEFAULT NULL,
  `region_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT 'not FK',
  `voucher_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT 'not FK',
  `price_level_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT 'not FK',
  `payment_group_sls_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT 'not FK',
  `archetype_id` bigint UNSIGNED NULL DEFAULT 0 COMMENT 'not FK',
  `order_type_sls_id` bigint UNSIGNED NULL DEFAULT 0,
  `sales_order_number` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `integration_code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'for integration with talon.one',
  `status` tinyint NULL DEFAULT 0,
  `recognition_date` date NULL DEFAULT NULL,
  `requests_delivery_date` date NULL DEFAULT NULL,
  `billing_address` varchar(350) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `shipping_address` varchar(400) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `shipping_address_note` varchar(250) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `delivery_fee` decimal(10, 2) NULL DEFAULT 0.00,
  `vou_redeem_code` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `vou_disc_amount` decimal(20, 2) NULL DEFAULT 0.00,
  `point_redeem_amount` decimal(20, 2) NULL DEFAULT 0.00,
  `point_redeem_id` bigint UNSIGNED NULL DEFAULT 0,
  `eden_point_campaign_id` bigint UNSIGNED NULL DEFAULT NULL COMMENT 'Not FK',
  `total_sku_disc_amount` decimal(20, 2) NULL DEFAULT 0.00,
  `total_price` decimal(20, 2) NULL DEFAULT 0.00,
  `total_charge` decimal(20, 2) NULL DEFAULT 0.00,
  `total_weight` decimal(10, 2) NULL DEFAULT 0.00,
  `note` varchar(250) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `payment_reminder` tinyint(1) NOT NULL DEFAULT 2,
  `cancel_type` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` bigint UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 380 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sales_order
-- ----------------------------
INSERT INTO `sales_order` VALUES (1, 1, 1, 7, 3, 2, NULL, 4588, 1, 13, 2, NULL, 7, 3, 2, 1, 'SO105C-29071901', '', 1, '2019-07-29', '2019-07-30', 'Jalan Pegangsaan Barat No 75', 'Jalan Pegangsaan Timur No 92', '', 0.00, NULL, 0.00, 0.00, 0, 0, 0.00, 1418424.00, 3465920.00, 0.00, 'Lorem Ipsum', 2, 0, '2019-07-29 19:03:53', 222);
INSERT INTO `sales_order` VALUES (2, 1, 1, 4, 2, 3, NULL, 4568, 1, 13, 2, NULL, 6, 3, 1, 1, 'SO122-29071901', '', 1, '2019-07-29', '2019-07-30', 'Jalan Pegangsaan Barat No 10', 'Jalan Pegangsaan Timur No 9', '', 0.00, NULL, 0.00, 0.00, 0, 0, 0.00, 1375704.00, 3575279.00, 0.00, 'Lorem Ipsum', 2, 0, '2019-07-29 19:23:19', 222);
INSERT INTO `sales_order` VALUES (3, 1, 1, 9, 2, 3, NULL, 1558, 1, 13, 2, NULL, 33, 2, 4, 1, 'SO131-29071901', '', 1, '2019-07-29', '2019-07-30', 'Jalan Pegangsaan Barat No 25', 'Jalan Pegangsaan Timur No 66', '', 0.00, NULL, 0.00, 0.00, 0, 0, 0.00, 2523245.00, 7378638.00, 0.00, '', 2, 0, '2019-07-29 19:40:35', 222);
INSERT INTO `sales_order` VALUES (4, 1, 1, 4, 2, 3, NULL, 1604, 1, 13, 2, NULL, 6, 2, 1, 1, 'SO44-29071901', '', 1, '2019-07-29', '2019-07-30', 'Jalan Pegangsaan Barat No 94', 'Jalan Pegangsaan Timur No 6', '', 0.00, NULL, 0.00, 0.00, 0, 0, 0.00, 8389117.00, 6067353.00, 0.00, 'Lorem Ipsum', 2, 0, '2019-07-29 19:46:37', 222);
INSERT INTO `sales_order` VALUES (5, 1, 1, 7, 3, 2, NULL, 4588, 1, 13, 2, NULL, 7, 3, 2, 1, 'SO105C-29071902', '', 1, '2019-07-29', '2019-07-30', 'Jalan Pegangsaan Barat No 98', 'Jalan Pegangsaan Timur No 28', '', 0.00, NULL, 0.00, 0.00, 0, 0, 0.00, 4275851.00, 8100851.00, 0.00, 'Lorem Ipsum', 2, 0, '2019-07-29 19:47:39', 222);
INSERT INTO `sales_order` VALUES (6, 1, 1, 7, 3, 2, NULL, 4588, 1, 13, 2, NULL, 7, 3, 2, 1, 'SO105C-29071903', '', 1, '2019-07-29', '2019-07-30', 'Jalan Pegangsaan Barat No 9', 'Jalan Pegangsaan Timur No 21', '', 0.00, NULL, 0.00, 0.00, 0, 0, 0.00, 6111905.00, 2202197.00, 0.00, 'Lorem Ipsum', 2, 0, '2019-07-29 19:49:35', 222);
INSERT INTO `sales_order` VALUES (7, 1, 1, 7, 3, 2, NULL, 4588, 1, 13, 2, NULL, 7, 3, 2, 1, 'SO105C-29071904', '', 1, '2019-07-29', '2019-07-30', 'Jalan Pegangsaan Barat No 48', 'Jalan Pegangsaan Timur No 24', '', 0.00, NULL, 0.00, 0.00, 0, 0, 0.00, 7631972.00, 6608432.00, 0.00, 'Lorem Ipsum', 2, 0, '2019-07-29 19:52:19', 222);
INSERT INTO `sales_order` VALUES (8, 1, 1, 7, 3, 2, NULL, 4588, 1, 13, 2, NULL, 7, 3, 2, 1, 'SO105B-29071901', '', 1, '2019-07-29', '2019-07-30', 'Jalan Pegangsaan Barat No 14', 'Jalan Pegangsaan Timur No 56', '', 0.00, NULL, 0.00, 0.00, 0, 0, 0.00, 9724147.00, 6335569.00, 0.00, 'Lorem Ipsum', 2, 0, '2019-07-29 20:07:44', 222);
INSERT INTO `sales_order` VALUES (9, 1, 1, 7, 3, 2, NULL, 4588, 1, 13, 2, NULL, 7, 3, 2, 1, 'SO105B-29071902', '', 1, '2019-07-29', '2019-07-30', 'Jalan Pegangsaan Barat No 26', 'Jalan Pegangsaan Timur No 7', '', 0.00, NULL, 0.00, 0.00, 0, 0, 0.00, 5624819.00, 1752551.00, 0.00, 'Lorem Ipsum', 2, 0, '2019-07-29 20:27:00', 222);
INSERT INTO `sales_order` VALUES (10, 1, 1, 7, 3, 2, NULL, 4588, 1, 13, 2, NULL, 7, 3, 2, 1, 'SO105B-29071903', '', 1, '2019-07-29', '2019-07-30', 'Jalan Pegangsaan Barat No 86', 'Jalan Pegangsaan Timur No 67', '', 0.00, NULL, 0.00, 0.00, 0, 0, 0.00, 8851653.00, 9656050.00, 0.00, 'Lorem Ipsum', 2, 0, '2019-07-29 20:33:30', 222);
INSERT INTO `sales_order` VALUES (11, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, '', '', 1, '2023-03-13', '2023-03-13', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 100000.00, 100000.00, 10.00, '', 0, 0, '2023-03-13 11:17:26', 0);
INSERT INTO `sales_order` VALUES (12, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, '', '', 1, '2023-03-14', '2023-03-14', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-14 04:35:00', 0);
INSERT INTO `sales_order` VALUES (13, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, '', '', 1, '2023-03-14', '2023-03-14', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-14 04:52:15', 0);
INSERT INTO `sales_order` VALUES (14, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, '', '', 1, '2023-03-14', '2023-03-14', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-14 09:52:43', 0);
INSERT INTO `sales_order` VALUES (15, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, '', '', 1, '2023-03-14', '2023-03-14', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-14 10:14:59', 0);
INSERT INTO `sales_order` VALUES (16, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, '', '', 1, '2023-03-14', '2023-03-14', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-14 10:15:22', 0);
INSERT INTO `sales_order` VALUES (17, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, '', '', 1, '2023-03-14', '2023-03-14', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-14 10:18:31', 0);
INSERT INTO `sales_order` VALUES (18, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, '', '', 1, '2023-03-14', '2023-03-14', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-14 10:19:03', 0);
INSERT INTO `sales_order` VALUES (19, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000001', '', 1, '2023-03-14', '2023-10-08', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-14 11:12:53', 0);
INSERT INTO `sales_order` VALUES (20, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000002', '', 1, '2023-03-14', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-14 11:14:38', 0);
INSERT INTO `sales_order` VALUES (21, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000003', '', 1, '2023-03-14', '2023-03-15', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 110000.00, 110000.00, 22.00, '', 0, 0, '2023-03-14 14:11:29', 0);
INSERT INTO `sales_order` VALUES (22, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000004', '', 1, '2023-03-14', '2023-03-15', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 110000.00, 110000.00, 22.00, '', 0, 0, '2023-03-14 14:11:30', 0);
INSERT INTO `sales_order` VALUES (23, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000005', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 04:35:39', 0);
INSERT INTO `sales_order` VALUES (24, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000006', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 04:44:15', 0);
INSERT INTO `sales_order` VALUES (25, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000007', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 04:49:40', 0);
INSERT INTO `sales_order` VALUES (26, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000008', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 04:53:41', 0);
INSERT INTO `sales_order` VALUES (27, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000009', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 04:55:47', 0);
INSERT INTO `sales_order` VALUES (28, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000010', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 04:56:36', 0);
INSERT INTO `sales_order` VALUES (29, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000011', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 04:58:14', 0);
INSERT INTO `sales_order` VALUES (30, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000012', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 04:59:44', 0);
INSERT INTO `sales_order` VALUES (31, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000013', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 05:02:43', 0);
INSERT INTO `sales_order` VALUES (32, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000014', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 06:24:48', 0);
INSERT INTO `sales_order` VALUES (33, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000015', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 06:34:16', 0);
INSERT INTO `sales_order` VALUES (34, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000016', '', 1, '2023-03-15', '2023-10-08', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 06:34:54', 0);
INSERT INTO `sales_order` VALUES (35, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000017', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 06:36:54', 0);
INSERT INTO `sales_order` VALUES (36, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000018', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 07:21:27', 0);
INSERT INTO `sales_order` VALUES (37, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000019', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 08:11:45', 0);
INSERT INTO `sales_order` VALUES (38, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000020', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 08:29:15', 0);
INSERT INTO `sales_order` VALUES (39, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000021', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 08:31:24', 0);
INSERT INTO `sales_order` VALUES (40, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000022', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 08:37:37', 0);
INSERT INTO `sales_order` VALUES (41, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000023', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 08:44:12', 0);
INSERT INTO `sales_order` VALUES (42, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000024', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 08:44:35', 0);
INSERT INTO `sales_order` VALUES (43, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000025', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 09:44:16', 0);
INSERT INTO `sales_order` VALUES (44, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000026', '', 1, '2023-03-15', '2023-10-08', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 09:44:30', 0);
INSERT INTO `sales_order` VALUES (45, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000027', '', 1, '2023-03-15', '2023-10-08', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 09:45:46', 0);
INSERT INTO `sales_order` VALUES (46, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000028', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 15000.00, 3.00, '', 0, 0, '2023-03-15 11:06:17', 0);
INSERT INTO `sales_order` VALUES (47, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000029', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 11:06:37', 0);
INSERT INTO `sales_order` VALUES (48, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000030', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 11:09:13', 0);
INSERT INTO `sales_order` VALUES (49, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000031', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 35000.00, 7.00, '', 0, 0, '2023-03-15 11:12:28', 0);
INSERT INTO `sales_order` VALUES (50, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000032', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 11:12:48', 0);
INSERT INTO `sales_order` VALUES (51, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000033', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 11:25:20', 0);
INSERT INTO `sales_order` VALUES (52, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000034', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 11:28:14', 0);
INSERT INTO `sales_order` VALUES (53, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000035', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 11:31:46', 0);
INSERT INTO `sales_order` VALUES (54, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000036', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0, 0, '2023-03-15 11:32:11', 0);
INSERT INTO `sales_order` VALUES (55, 1, 1, 2, 2, 0, 0, 1, 1, 9, 0, 0, 0, 0, 1, 0, 'SO000037', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 12:33:37', 0);
INSERT INTO `sales_order` VALUES (56, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000038', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 13:03:31', 0);
INSERT INTO `sales_order` VALUES (57, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000039', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 250000.00, 250000.00, 50.00, '', 0, 0, '2023-03-15 13:46:00', 0);
INSERT INTO `sales_order` VALUES (58, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000040', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 250000.00, 250000.00, 50.00, '', 0, 0, '2023-03-15 13:46:05', 0);
INSERT INTO `sales_order` VALUES (59, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000041', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 250000.00, 250000.00, 50.00, '', 0, 0, '2023-03-15 13:46:14', 0);
INSERT INTO `sales_order` VALUES (60, 1, 1, 2, 2, 0, 0, 1, 1, 2, 0, 0, 0, 0, 1, 0, 'SO000042', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 250000.00, 250000.00, 50.00, '', 0, 0, '2023-03-15 13:47:11', 0);
INSERT INTO `sales_order` VALUES (61, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000043', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 250000.00, 250000.00, 50.00, '', 0, 0, '2023-03-15 13:48:17', 0);
INSERT INTO `sales_order` VALUES (62, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000044', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 250000.00, 250000.00, 50.00, '', 0, 0, '2023-03-15 13:48:51', 0);
INSERT INTO `sales_order` VALUES (63, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000045', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 13:54:33', 0);
INSERT INTO `sales_order` VALUES (64, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000046', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 35000.00, 7.00, '', 0, 0, '2023-03-15 13:57:15', 0);
INSERT INTO `sales_order` VALUES (65, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000047', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 35000.00, 7.00, '', 0, 0, '2023-03-15 13:58:01', 0);
INSERT INTO `sales_order` VALUES (66, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000048', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 35000.00, 7.00, '', 0, 0, '2023-03-15 13:58:29', 0);
INSERT INTO `sales_order` VALUES (67, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000049', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 35000.00, 7.00, '', 0, 0, '2023-03-15 13:59:04', 0);
INSERT INTO `sales_order` VALUES (68, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000050', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 14:00:29', 0);
INSERT INTO `sales_order` VALUES (69, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000051', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 35000.00, 7.00, '', 0, 0, '2023-03-15 14:04:11', 0);
INSERT INTO `sales_order` VALUES (70, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000052', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 35000.00, 7.00, '', 0, 0, '2023-03-15 14:06:12', 0);
INSERT INTO `sales_order` VALUES (71, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000053', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 35000.00, 7.00, '', 0, 0, '2023-03-15 14:06:56', 0);
INSERT INTO `sales_order` VALUES (72, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000054', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 35000.00, 7.00, '', 0, 0, '2023-03-15 14:07:30', 0);
INSERT INTO `sales_order` VALUES (73, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000055', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 14:07:56', 0);
INSERT INTO `sales_order` VALUES (74, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000056', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 50000.00, 10.00, '', 0, 0, '2023-03-15 14:08:51', 0);
INSERT INTO `sales_order` VALUES (75, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000057', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 14:10:19', 0);
INSERT INTO `sales_order` VALUES (76, 1, 1, 1, 1, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000058', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 45000.00, 7.00, '', 0, 0, '2023-03-15 14:12:02', 0);
INSERT INTO `sales_order` VALUES (77, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000059', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:15:16', 0);
INSERT INTO `sales_order` VALUES (78, 1, 1, 2, 2, 0, 0, 1, 1, 2, 0, 0, 0, 0, 1, 0, 'SO000060', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:16:34', 0);
INSERT INTO `sales_order` VALUES (79, 1, 1, 1, 1, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000061', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 30000.00, 40000.00, 6.00, '', 0, 0, '2023-03-15 14:16:35', 0);
INSERT INTO `sales_order` VALUES (80, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000062', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 14:17:38', 0);
INSERT INTO `sales_order` VALUES (81, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000063', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 14:18:36', 0);
INSERT INTO `sales_order` VALUES (82, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000064', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:18:43', 0);
INSERT INTO `sales_order` VALUES (83, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000065', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 14:18:56', 0);
INSERT INTO `sales_order` VALUES (84, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000066', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:19:25', 0);
INSERT INTO `sales_order` VALUES (85, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000067', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:19:28', 0);
INSERT INTO `sales_order` VALUES (86, 1, 1, 2, 2, 0, 0, 1, 1, 9, 0, 0, 0, 0, 1, 0, 'SO000068', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:19:36', 0);
INSERT INTO `sales_order` VALUES (87, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000069', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-15 14:19:48', 0);
INSERT INTO `sales_order` VALUES (88, 1, 1, 1, 1, 0, 0, 1, 1, 9, 0, 0, 0, 0, 1, 0, 'SO000070', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:20:59', 0);
INSERT INTO `sales_order` VALUES (89, 1, 1, 1, 1, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000071', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:21:06', 0);
INSERT INTO `sales_order` VALUES (90, 1, 1, 1, 1, 0, 0, 1, 1, 9, 0, 0, 0, 0, 1, 0, 'SO000072', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:22:01', 0);
INSERT INTO `sales_order` VALUES (91, 1, 1, 1, 1, 0, 0, 1, 1, 9, 0, 0, 0, 0, 1, 0, 'SO000073', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:22:08', 0);
INSERT INTO `sales_order` VALUES (92, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000074', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 14:22:28', 0);
INSERT INTO `sales_order` VALUES (93, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000075', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:22:55', 0);
INSERT INTO `sales_order` VALUES (94, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000076', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:24:15', 0);
INSERT INTO `sales_order` VALUES (95, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000077', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 14:26:09', 0);
INSERT INTO `sales_order` VALUES (96, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000078', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 14:30:13', 0);
INSERT INTO `sales_order` VALUES (97, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000079', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:32:33', 0);
INSERT INTO `sales_order` VALUES (98, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000080', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:34:15', 0);
INSERT INTO `sales_order` VALUES (99, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000081', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:35:17', 0);
INSERT INTO `sales_order` VALUES (100, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000082', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 14:35:20', 0);
INSERT INTO `sales_order` VALUES (101, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000083', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:35:37', 0);
INSERT INTO `sales_order` VALUES (102, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000084', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:36:15', 0);
INSERT INTO `sales_order` VALUES (103, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000085', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 14:36:56', 0);
INSERT INTO `sales_order` VALUES (104, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000086', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:37:08', 0);
INSERT INTO `sales_order` VALUES (105, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000087', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 14:38:44', 0);
INSERT INTO `sales_order` VALUES (106, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000088', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:39:21', 0);
INSERT INTO `sales_order` VALUES (107, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000089', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 45000.00, 7.00, '', 0, 0, '2023-03-15 14:40:51', 0);
INSERT INTO `sales_order` VALUES (108, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000090', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 14:40:58', 0);
INSERT INTO `sales_order` VALUES (109, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000091', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:41:06', 0);
INSERT INTO `sales_order` VALUES (110, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000092', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 25000.00, 35000.00, 5.00, '', 0, 0, '2023-03-15 14:41:55', 0);
INSERT INTO `sales_order` VALUES (111, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000093', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 35000.00, 45000.00, 7.00, '', 0, 0, '2023-03-15 14:42:34', 0);
INSERT INTO `sales_order` VALUES (112, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000094', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:43:23', 0);
INSERT INTO `sales_order` VALUES (113, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000095', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:43:50', 0);
INSERT INTO `sales_order` VALUES (114, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000096', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 14:44:01', 0);
INSERT INTO `sales_order` VALUES (115, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000097', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:45:04', 0);
INSERT INTO `sales_order` VALUES (116, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000098', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 14:47:44', 0);
INSERT INTO `sales_order` VALUES (117, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000099', '', 1, '2023-03-15', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-15 14:47:54', 0);
INSERT INTO `sales_order` VALUES (118, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000100', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 15:33:43', 0);
INSERT INTO `sales_order` VALUES (119, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000101', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-15 15:34:11', 0);
INSERT INTO `sales_order` VALUES (120, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000102', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 15:34:39', 0);
INSERT INTO `sales_order` VALUES (121, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000103', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 15:38:17', 0);
INSERT INTO `sales_order` VALUES (122, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000104', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 15:51:21', 0);
INSERT INTO `sales_order` VALUES (123, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000105', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-15 16:18:55', 0);
INSERT INTO `sales_order` VALUES (124, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000106', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 16:25:51', 0);
INSERT INTO `sales_order` VALUES (125, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000107', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 16:26:53', 0);
INSERT INTO `sales_order` VALUES (126, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000108', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 25000.00, 35000.00, 5.00, '', 0, 0, '2023-03-15 16:33:28', 0);
INSERT INTO `sales_order` VALUES (127, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000109', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 25000.00, 35000.00, 5.00, '', 0, 0, '2023-03-15 16:33:42', 0);
INSERT INTO `sales_order` VALUES (128, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000110', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-15 19:20:18', 0);
INSERT INTO `sales_order` VALUES (129, 1, 1, 3, 3, 0, 0, 1, 1, 2, 0, 0, 0, 0, 1, 0, 'SO000111', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 25000.00, 35000.00, 5.00, '', 0, 0, '2023-03-15 19:28:39', 0);
INSERT INTO `sales_order` VALUES (130, 1, 1, 2, 2, 0, 0, 1, 1, 2, 0, 0, 0, 0, 1, 0, 'SO000112', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-15 20:25:27', 0);
INSERT INTO `sales_order` VALUES (131, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000113', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 20:27:01', 0);
INSERT INTO `sales_order` VALUES (132, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000114', '', 1, '2023-03-15', '2023-03-16', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-15 20:28:33', 0);
INSERT INTO `sales_order` VALUES (133, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000115', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 85000.00, 95000.00, 17.00, '', 0, 0, '2023-03-15 20:29:46', 0);
INSERT INTO `sales_order` VALUES (134, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000116', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 22:47:41', 0);
INSERT INTO `sales_order` VALUES (135, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000117', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 22:54:41', 0);
INSERT INTO `sales_order` VALUES (136, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000118', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 22:59:53', 0);
INSERT INTO `sales_order` VALUES (137, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000119', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 23:02:08', 0);
INSERT INTO `sales_order` VALUES (138, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000120', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-15 23:06:42', 0);
INSERT INTO `sales_order` VALUES (139, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000121', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 23:11:45', 0);
INSERT INTO `sales_order` VALUES (140, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000122', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 23:15:16', 0);
INSERT INTO `sales_order` VALUES (141, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000123', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-15 23:18:27', 0);
INSERT INTO `sales_order` VALUES (142, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000124', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 23:23:15', 0);
INSERT INTO `sales_order` VALUES (143, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000125', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 23:27:52', 0);
INSERT INTO `sales_order` VALUES (144, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000126', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 23:29:26', 0);
INSERT INTO `sales_order` VALUES (145, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000127', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-15 23:32:35', 0);
INSERT INTO `sales_order` VALUES (146, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000128', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-15 23:34:07', 0);
INSERT INTO `sales_order` VALUES (147, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000129', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 23:35:51', 0);
INSERT INTO `sales_order` VALUES (148, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000130', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 25000.00, 35000.00, 5.00, '', 0, 0, '2023-03-15 23:40:06', 0);
INSERT INTO `sales_order` VALUES (149, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000131', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-15 23:41:02', 0);
INSERT INTO `sales_order` VALUES (150, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000132', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-15 23:47:15', 0);
INSERT INTO `sales_order` VALUES (151, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000133', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 23:48:18', 0);
INSERT INTO `sales_order` VALUES (152, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000134', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 23:52:43', 0);
INSERT INTO `sales_order` VALUES (153, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000135', '', 1, '2023-03-15', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-15 23:54:22', 0);
INSERT INTO `sales_order` VALUES (154, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000136', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:00:23', 0);
INSERT INTO `sales_order` VALUES (155, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000137', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:02:03', 0);
INSERT INTO `sales_order` VALUES (156, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000138', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:02:57', 0);
INSERT INTO `sales_order` VALUES (157, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000139', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:04:01', 0);
INSERT INTO `sales_order` VALUES (158, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000140', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:05:01', 0);
INSERT INTO `sales_order` VALUES (159, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000141', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:05:30', 0);
INSERT INTO `sales_order` VALUES (160, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000142', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:07:22', 0);
INSERT INTO `sales_order` VALUES (161, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000143', '', 1, '2023-03-16', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-16 00:17:36', 0);
INSERT INTO `sales_order` VALUES (162, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000144', '', 1, '2023-03-16', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-16 00:17:47', 0);
INSERT INTO `sales_order` VALUES (163, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000145', '', 1, '2023-03-16', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-16 00:19:31', 0);
INSERT INTO `sales_order` VALUES (164, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000146', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-16 00:24:21', 0);
INSERT INTO `sales_order` VALUES (165, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000147', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:24:58', 0);
INSERT INTO `sales_order` VALUES (166, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000148', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:29:16', 0);
INSERT INTO `sales_order` VALUES (167, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000149', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:30:03', 0);
INSERT INTO `sales_order` VALUES (168, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000150', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:30:20', 0);
INSERT INTO `sales_order` VALUES (169, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000151', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:30:41', 0);
INSERT INTO `sales_order` VALUES (170, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000152', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-16 00:34:34', 0);
INSERT INTO `sales_order` VALUES (171, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000153', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-16 00:36:22', 0);
INSERT INTO `sales_order` VALUES (172, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000154', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-16 00:38:53', 0);
INSERT INTO `sales_order` VALUES (173, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000155', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 25000.00, 35000.00, 5.00, '', 0, 0, '2023-03-16 00:41:55', 0);
INSERT INTO `sales_order` VALUES (174, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000156', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:43:16', 0);
INSERT INTO `sales_order` VALUES (175, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000157', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:47:12', 0);
INSERT INTO `sales_order` VALUES (176, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000158', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:51:41', 0);
INSERT INTO `sales_order` VALUES (177, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000159', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 00:56:10', 0);
INSERT INTO `sales_order` VALUES (178, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000160', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:00:40', 0);
INSERT INTO `sales_order` VALUES (179, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000161', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:02:46', 0);
INSERT INTO `sales_order` VALUES (180, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000162', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:06:19', 0);
INSERT INTO `sales_order` VALUES (181, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000163', '', 1, '2023-03-16', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-16 01:09:43', 0);
INSERT INTO `sales_order` VALUES (182, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000164', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:10:58', 0);
INSERT INTO `sales_order` VALUES (183, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000165', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:17:28', 0);
INSERT INTO `sales_order` VALUES (184, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000166', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:17:59', 0);
INSERT INTO `sales_order` VALUES (185, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000167', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:20:01', 0);
INSERT INTO `sales_order` VALUES (186, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000168', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:20:49', 0);
INSERT INTO `sales_order` VALUES (187, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000169', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:24:03', 0);
INSERT INTO `sales_order` VALUES (188, 1, 1, 2, 2, 0, 0, 1, 1, 9, 0, 0, 0, 0, 1, 0, 'SO000170', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:26:39', 0);
INSERT INTO `sales_order` VALUES (189, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000171', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:30:25', 0);
INSERT INTO `sales_order` VALUES (190, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000172', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:35:52', 0);
INSERT INTO `sales_order` VALUES (191, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000173', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:39:29', 0);
INSERT INTO `sales_order` VALUES (192, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000174', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:41:05', 0);
INSERT INTO `sales_order` VALUES (193, 1, 1, 2, 2, 0, 0, 1, 1, 9, 0, 0, 0, 0, 1, 0, 'SO000175', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 1.00, '', 0, 0, '2023-03-16 01:44:57', 0);
INSERT INTO `sales_order` VALUES (194, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000176', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-16 01:47:35', 0);
INSERT INTO `sales_order` VALUES (195, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000177', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:48:45', 0);
INSERT INTO `sales_order` VALUES (196, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000178', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-16 01:49:27', 0);
INSERT INTO `sales_order` VALUES (197, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000179', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:53:23', 0);
INSERT INTO `sales_order` VALUES (198, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000180', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:54:54', 0);
INSERT INTO `sales_order` VALUES (199, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000181', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:56:17', 0);
INSERT INTO `sales_order` VALUES (200, 1, 1, 2, 2, 0, 0, 1, 1, 2, 0, 0, 0, 0, 1, 0, 'SO000182', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:57:46', 0);
INSERT INTO `sales_order` VALUES (201, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000183', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 01:59:52', 0);
INSERT INTO `sales_order` VALUES (202, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000184', '', 1, '2023-03-16', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 25000.00, 35000.00, 5.00, '', 0, 0, '2023-03-16 02:34:31', 0);
INSERT INTO `sales_order` VALUES (203, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000185', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 03:10:10', 0);
INSERT INTO `sales_order` VALUES (204, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000186', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 20000.00, 30000.00, 4.00, '', 0, 0, '2023-03-16 03:15:34', 0);
INSERT INTO `sales_order` VALUES (205, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000187', '', 1, '2023-03-16', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 15000.00, 25000.00, 3.00, '', 0, 0, '2023-03-16 03:53:38', 0);
INSERT INTO `sales_order` VALUES (206, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000188', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-16 03:54:37', 0);
INSERT INTO `sales_order` VALUES (207, 1, 1, 3, 3, 0, 0, 1, 1, 2, 0, 0, 0, 0, 1, 0, 'SO000189', '', 1, '2023-03-16', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-16 03:56:47', 0);
INSERT INTO `sales_order` VALUES (208, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000190', '', 1, '2023-03-16', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 30000.00, 40000.00, 6.00, '', 0, 0, '2023-03-16 03:59:30', 0);
INSERT INTO `sales_order` VALUES (209, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000191', '', 1, '2023-03-16', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 45000.00, 55000.00, 9.00, '', 0, 0, '2023-03-16 03:59:44', 0);
INSERT INTO `sales_order` VALUES (210, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000192', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-16 04:00:22', 0);
INSERT INTO `sales_order` VALUES (211, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000193', '', 1, '2023-03-16', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-16 07:03:38', 0);
INSERT INTO `sales_order` VALUES (212, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000194', '', 1, '2023-03-16', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 210000.00, 210000.00, 42.00, '', 0, 0, '2023-03-16 09:34:22', 0);
INSERT INTO `sales_order` VALUES (213, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000195', '', 1, '2023-03-17', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-17 03:40:07', 0);
INSERT INTO `sales_order` VALUES (214, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000196', '', 1, '2023-03-17', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-17 07:07:16', 0);
INSERT INTO `sales_order` VALUES (215, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000197', '', 1, '2023-03-17', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-17 07:21:36', 0);
INSERT INTO `sales_order` VALUES (216, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000198', '', 1, '2023-03-17', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-17 08:07:06', 0);
INSERT INTO `sales_order` VALUES (217, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000199', '', 1, '2023-03-17', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-17 08:20:33', 0);
INSERT INTO `sales_order` VALUES (218, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000200', '', 1, '2023-03-17', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-17 08:22:54', 0);
INSERT INTO `sales_order` VALUES (219, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000201', '', 1, '2023-03-17', '2023-03-18', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-17 08:58:06', 0);
INSERT INTO `sales_order` VALUES (220, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000202', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 02:13:12', 0);
INSERT INTO `sales_order` VALUES (221, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000203', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 02:21:39', 0);
INSERT INTO `sales_order` VALUES (222, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000204', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 02:31:29', 0);
INSERT INTO `sales_order` VALUES (223, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000205', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 02:40:32', 0);
INSERT INTO `sales_order` VALUES (224, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000206', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 02:41:49', 0);
INSERT INTO `sales_order` VALUES (225, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000207', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 02:53:02', 0);
INSERT INTO `sales_order` VALUES (226, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000208', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 02:56:22', 0);
INSERT INTO `sales_order` VALUES (227, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000209', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 02:58:15', 0);
INSERT INTO `sales_order` VALUES (228, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000210', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 02:59:43', 0);
INSERT INTO `sales_order` VALUES (229, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000211', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 03:10:20', 0);
INSERT INTO `sales_order` VALUES (230, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000212', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 100000.00, 110000.00, 20.00, '', 0, 0, '2023-03-20 04:03:16', 0);
INSERT INTO `sales_order` VALUES (231, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000213', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 100000.00, 110000.00, 20.00, '', 0, 0, '2023-03-20 04:05:55', 0);
INSERT INTO `sales_order` VALUES (232, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000214', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 100000.00, 110000.00, 20.00, '', 0, 0, '2023-03-20 04:08:11', 0);
INSERT INTO `sales_order` VALUES (233, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000215', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 04:09:01', 0);
INSERT INTO `sales_order` VALUES (234, 1, 1, 1, 1, 0, 0, 1, 1, 2, 0, 0, 0, 0, 1, 0, 'SO000216', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 04:11:23', 0);
INSERT INTO `sales_order` VALUES (235, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000217', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 04:40:59', 0);
INSERT INTO `sales_order` VALUES (236, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000218', '', 1, '2023-03-20', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 04:43:58', 0);
INSERT INTO `sales_order` VALUES (237, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000219', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 05:01:57', 0);
INSERT INTO `sales_order` VALUES (238, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000220', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 06:00:54', 0);
INSERT INTO `sales_order` VALUES (239, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000221', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 06:30:35', 0);
INSERT INTO `sales_order` VALUES (240, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000222', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 07:25:39', 0);
INSERT INTO `sales_order` VALUES (241, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000223', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 08:09:29', 0);
INSERT INTO `sales_order` VALUES (242, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000224', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 08:24:01', 0);
INSERT INTO `sales_order` VALUES (243, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000225', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 08:42:14', 0);
INSERT INTO `sales_order` VALUES (244, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000226', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 09:54:23', 0);
INSERT INTO `sales_order` VALUES (245, 1, 1, 2, 2, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000227', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 09:57:33', 0);
INSERT INTO `sales_order` VALUES (246, 1, 1, 3, 3, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000228', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 09:59:17', 0);
INSERT INTO `sales_order` VALUES (247, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000229', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 10:18:29', 0);
INSERT INTO `sales_order` VALUES (248, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000230', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 10:18:58', 0);
INSERT INTO `sales_order` VALUES (249, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000231', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 10:29:08', 0);
INSERT INTO `sales_order` VALUES (250, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000232', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 10:30:36', 0);
INSERT INTO `sales_order` VALUES (251, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000233', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 10:32:12', 0);
INSERT INTO `sales_order` VALUES (252, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000234', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 10:34:35', 0);
INSERT INTO `sales_order` VALUES (253, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000235', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 10:36:42', 0);
INSERT INTO `sales_order` VALUES (254, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000236', '', 1, '2023-03-20', '2023-03-21', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 11:01:56', 0);
INSERT INTO `sales_order` VALUES (255, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000237', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 11:02:27', 0);
INSERT INTO `sales_order` VALUES (256, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000238', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 11:12:41', 0);
INSERT INTO `sales_order` VALUES (257, 1, 1, 2, 2, 0, 0, 1, 1, 2, 0, 0, 0, 0, 1, 0, 'SO000239', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-20 12:09:43', 0);
INSERT INTO `sales_order` VALUES (258, 1, 1, 3, 3, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000240', '', 1, '2023-03-20', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-20 12:17:22', 0);
INSERT INTO `sales_order` VALUES (259, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000241', '', 1, '2023-03-21', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 10.00, '', 0, 0, '2023-03-21 01:22:09', 0);
INSERT INTO `sales_order` VALUES (260, 1, 1, 3, 3, 0, 0, 1, 1, 2, 0, 0, 0, 0, 1, 0, 'SO000242', '', 1, '2023-03-21', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 2.00, '', 0, 0, '2023-03-21 04:14:05', 0);
INSERT INTO `sales_order` VALUES (261, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000243', '', 1, '2023-03-21', '2023-03-22', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-21 11:01:21', 0);
INSERT INTO `sales_order` VALUES (262, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000244', '', 1, '2023-03-22', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-22 13:54:41', 0);
INSERT INTO `sales_order` VALUES (263, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000245', '', 1, '2023-03-22', '2023-03-23', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-22 13:59:57', 0);
INSERT INTO `sales_order` VALUES (266, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000248', '', 1, '2023-03-23', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 07:02:17', 0);
INSERT INTO `sales_order` VALUES (267, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000249', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 07:03:44', 0);
INSERT INTO `sales_order` VALUES (268, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000250', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 07:22:20', 0);
INSERT INTO `sales_order` VALUES (269, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000251', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 07:57:44', 0);
INSERT INTO `sales_order` VALUES (270, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000252', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 09:25:17', 0);
INSERT INTO `sales_order` VALUES (271, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000253', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 09:36:14', 0);
INSERT INTO `sales_order` VALUES (272, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000254', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:23:27', 0);
INSERT INTO `sales_order` VALUES (273, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000255', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:24:54', 0);
INSERT INTO `sales_order` VALUES (274, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000256', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:25:00', 0);
INSERT INTO `sales_order` VALUES (275, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000257', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:25:06', 0);
INSERT INTO `sales_order` VALUES (276, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000258', '', 1, '2023-03-23', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:27:15', 0);
INSERT INTO `sales_order` VALUES (277, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000259', '', 1, '2023-03-23', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:28:34', 0);
INSERT INTO `sales_order` VALUES (278, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000260', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:29:08', 0);
INSERT INTO `sales_order` VALUES (279, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000261', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:30:09', 0);
INSERT INTO `sales_order` VALUES (280, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000262', '', 1, '2023-03-23', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:30:18', 0);
INSERT INTO `sales_order` VALUES (281, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000263', '', 1, '2023-03-23', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:31:10', 0);
INSERT INTO `sales_order` VALUES (282, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000264', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:31:26', 0);
INSERT INTO `sales_order` VALUES (283, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000265', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:31:31', 0);
INSERT INTO `sales_order` VALUES (284, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000266', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 25000.00, 35000.00, 0.00, '', 0, 0, '2023-03-23 10:31:54', 0);
INSERT INTO `sales_order` VALUES (285, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000267', '', 1, '2023-03-23', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:32:17', 0);
INSERT INTO `sales_order` VALUES (286, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000268', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:32:30', 0);
INSERT INTO `sales_order` VALUES (287, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000269', '', 1, '2023-03-23', '2023-03-17', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:32:40', 0);
INSERT INTO `sales_order` VALUES (288, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000270', '', 1, '2023-03-23', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:32:52', 0);
INSERT INTO `sales_order` VALUES (289, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000271', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:33:02', 0);
INSERT INTO `sales_order` VALUES (290, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000272', '', 1, '2023-03-23', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:34:27', 0);
INSERT INTO `sales_order` VALUES (291, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000273', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:34:48', 0);
INSERT INTO `sales_order` VALUES (292, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000274', '', 1, '2023-03-23', '2023-03-26', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:34:54', 0);
INSERT INTO `sales_order` VALUES (293, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000275', '', 1, '2023-03-23', '2023-03-26', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:34:58', 0);
INSERT INTO `sales_order` VALUES (294, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000276', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:35:55', 0);
INSERT INTO `sales_order` VALUES (295, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000277', '', 1, '2023-03-23', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:36:15', 0);
INSERT INTO `sales_order` VALUES (296, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000278', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:39:09', 0);
INSERT INTO `sales_order` VALUES (297, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000279', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:40:02', 0);
INSERT INTO `sales_order` VALUES (298, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000280', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:44:15', 0);
INSERT INTO `sales_order` VALUES (299, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000281', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 10:54:19', 0);
INSERT INTO `sales_order` VALUES (300, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000282', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 11:51:56', 0);
INSERT INTO `sales_order` VALUES (301, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000283', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 12:05:15', 0);
INSERT INTO `sales_order` VALUES (302, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000284', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 12:54:18', 0);
INSERT INTO `sales_order` VALUES (303, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000285', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 13:54:09', 0);
INSERT INTO `sales_order` VALUES (304, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000286', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 14:52:36', 0);
INSERT INTO `sales_order` VALUES (305, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000287', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 14:59:35', 0);
INSERT INTO `sales_order` VALUES (306, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000288', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 16:01:07', 0);
INSERT INTO `sales_order` VALUES (307, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000289', '', 1, '2023-03-23', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 16:56:38', 0);
INSERT INTO `sales_order` VALUES (308, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000290', '', 1, '2023-03-23', '2023-03-26', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 17:25:14', 0);
INSERT INTO `sales_order` VALUES (309, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000291', '', 1, '2023-03-23', '2023-03-26', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-23 17:33:06', 0);
INSERT INTO `sales_order` VALUES (310, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000292', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 01:20:39', 0);
INSERT INTO `sales_order` VALUES (311, 1, 1, 2, 2, 0, 0, 1, 1, 10, 0, 0, 0, 0, 1, 0, 'SO000293', '', 1, '2023-03-24', '2023-03-27', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 0.00, '', 0, 0, '2023-03-24 02:08:54', 0);
INSERT INTO `sales_order` VALUES (312, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000294', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 03:22:21', 0);
INSERT INTO `sales_order` VALUES (313, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000295', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 03:32:03', 0);
INSERT INTO `sales_order` VALUES (314, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000296', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 03:45:13', 0);
INSERT INTO `sales_order` VALUES (315, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000297', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 04:36:47', 0);
INSERT INTO `sales_order` VALUES (316, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000298', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 07:03:17', 0);
INSERT INTO `sales_order` VALUES (317, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000299', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 07:35:32', 0);
INSERT INTO `sales_order` VALUES (318, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000300', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 07:50:43', 0);
INSERT INTO `sales_order` VALUES (319, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000301', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 08:05:25', 0);
INSERT INTO `sales_order` VALUES (320, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000302', '', 1, '2023-03-24', '2023-03-24', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 1.00, 0, 0, 0.00, 50000.00, 59999.00, 0.00, '', 0, 0, '2023-03-24 08:15:53', 0);
INSERT INTO `sales_order` VALUES (321, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000303', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 1.00, 0, 0, 0.00, 50000.00, 59999.00, 0.00, '', 0, 0, '2023-03-24 08:19:57', 0);
INSERT INTO `sales_order` VALUES (322, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000304', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 08:20:06', 0);
INSERT INTO `sales_order` VALUES (323, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000305', '', 1, '2023-03-24', '2023-03-26', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 08:21:39', 0);
INSERT INTO `sales_order` VALUES (324, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000306', '', 1, '2023-03-24', '2023-03-26', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 08:31:40', 0);
INSERT INTO `sales_order` VALUES (325, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000307', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 08:32:46', 0);
INSERT INTO `sales_order` VALUES (326, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000308', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 08:34:28', 0);
INSERT INTO `sales_order` VALUES (327, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000309', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 08:39:00', 0);
INSERT INTO `sales_order` VALUES (328, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000310', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 09:01:43', 0);
INSERT INTO `sales_order` VALUES (329, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000311', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 09:36:24', 0);
INSERT INTO `sales_order` VALUES (330, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000312', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 09:36:40', 0);
INSERT INTO `sales_order` VALUES (331, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000313', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 09:37:24', 0);
INSERT INTO `sales_order` VALUES (332, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000314', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 09:59:42', 0);
INSERT INTO `sales_order` VALUES (333, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000315', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 10:00:01', 0);
INSERT INTO `sales_order` VALUES (334, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000316', '', 1, '2023-03-24', '2023-03-25', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-24 10:01:09', 0);
INSERT INTO `sales_order` VALUES (335, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000317', '', 1, '2023-03-27', '2023-03-28', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-27 05:31:42', 0);
INSERT INTO `sales_order` VALUES (336, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000318', '', 1, '2023-03-27', '2023-03-28', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-27 06:38:01', 0);
INSERT INTO `sales_order` VALUES (337, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000319', '', 1, '2023-03-27', '2023-03-28', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-27 06:51:30', 0);
INSERT INTO `sales_order` VALUES (338, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000320', '', 1, '2023-03-29', '2023-03-30', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-29 04:00:16', 0);
INSERT INTO `sales_order` VALUES (339, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000321', '', 1, '2023-03-29', '2023-03-30', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-29 04:49:14', 0);
INSERT INTO `sales_order` VALUES (340, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000322', '', 1, '2023-03-29', '2023-03-30', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-29 05:27:44', 0);
INSERT INTO `sales_order` VALUES (341, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000323', '', 1, '2023-03-29', '2023-03-30', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-29 05:51:58', 0);
INSERT INTO `sales_order` VALUES (342, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000324', '', 1, '2023-03-29', '2023-03-30', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-29 08:30:18', 0);
INSERT INTO `sales_order` VALUES (343, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000325', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 02:40:32', 0);
INSERT INTO `sales_order` VALUES (344, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000326', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 04:49:36', 0);
INSERT INTO `sales_order` VALUES (345, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000327', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 05:10:01', 0);
INSERT INTO `sales_order` VALUES (346, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000328', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 05:12:01', 0);
INSERT INTO `sales_order` VALUES (347, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000329', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 05:13:17', 0);
INSERT INTO `sales_order` VALUES (348, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000330', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 05:13:57', 0);
INSERT INTO `sales_order` VALUES (349, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000331', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 05:14:04', 0);
INSERT INTO `sales_order` VALUES (350, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000332', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 05:14:12', 0);
INSERT INTO `sales_order` VALUES (351, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000333', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 05:14:35', 0);
INSERT INTO `sales_order` VALUES (352, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000334', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 06:11:35', 0);
INSERT INTO `sales_order` VALUES (353, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000335', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 06:44:06', 0);
INSERT INTO `sales_order` VALUES (354, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000336', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 06:52:06', 0);
INSERT INTO `sales_order` VALUES (355, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000337', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 06:53:04', 0);
INSERT INTO `sales_order` VALUES (356, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000338', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 07:11:21', 0);
INSERT INTO `sales_order` VALUES (357, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000339', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 08:59:56', 0);
INSERT INTO `sales_order` VALUES (358, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000340', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 10:19:47', 0);
INSERT INTO `sales_order` VALUES (359, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000341', '', 1, '2023-03-30', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-03-30 10:52:51', 0);
INSERT INTO `sales_order` VALUES (360, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000342', '', 1, '2023-04-03', '2023-04-04', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-03 07:37:56', 0);
INSERT INTO `sales_order` VALUES (361, 1, 1, 0, 0, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000343', '', 1, '2023-04-03', '2023-04-05', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 55000.00, 65000.00, 0.00, '', 0, 0, '2023-04-03 08:23:20', 0);
INSERT INTO `sales_order` VALUES (362, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000344', '', 1, '2023-04-03', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-03 08:29:42', 0);
INSERT INTO `sales_order` VALUES (363, 0, 1, 0, 0, 0, 0, 1, 1, 41, 0, 0, 0, 0, 1, 0, 'SO000345', '', 1, '2023-04-03', '2023-04-04', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 500000.00, 500000.00, 0.00, '', 0, 0, '2023-04-03 09:01:17', 0);
INSERT INTO `sales_order` VALUES (364, 1, 1, 2, 2, 0, 0, 1, 1, 42, 0, 0, 0, 0, 1, 0, 'SO000346', '', 1, '2023-04-03', '2023-03-31', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-03 09:11:24', 0);
INSERT INTO `sales_order` VALUES (365, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000347', '', 1, '2023-04-03', '2023-04-04', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-03 10:14:40', 0);
INSERT INTO `sales_order` VALUES (366, 0, 1, 0, 0, 0, 0, 1, 1, 41, 0, 0, 0, 0, 1, 0, 'SO000348', '', 1, '2023-04-03', '2023-04-05', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-03 15:08:19', 0);
INSERT INTO `sales_order` VALUES (367, 0, 1, 0, 0, 0, 0, 1, 1, 47, 0, 0, 0, 0, 1, 0, 'SO000349', '', 1, '2023-04-04', '2023-04-06', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 5000.00, 15000.00, 0.00, '', 0, 0, '2023-04-04 02:22:10', 0);
INSERT INTO `sales_order` VALUES (368, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000350', '', 1, '2023-04-06', '2023-04-07', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-06 04:42:32', 0);
INSERT INTO `sales_order` VALUES (369, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000351', '', 1, '2023-04-06', '2023-04-07', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-06 05:53:21', 0);
INSERT INTO `sales_order` VALUES (370, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000352', '', 1, '2023-04-06', '2023-04-07', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-06 07:31:29', 0);
INSERT INTO `sales_order` VALUES (371, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000353', '', 1, '2023-04-06', '2023-04-07', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-06 09:05:00', 0);
INSERT INTO `sales_order` VALUES (372, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000354', '', 1, '2023-04-06', '2023-04-07', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-06 10:00:42', 0);
INSERT INTO `sales_order` VALUES (373, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000355', '', 1, '2023-04-10', '2023-04-11', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-10 02:46:22', 0);
INSERT INTO `sales_order` VALUES (374, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000356', '', 1, '2023-04-10', '2023-04-11', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-10 02:58:09', 0);
INSERT INTO `sales_order` VALUES (375, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000357', '', 1, '2023-04-10', '2023-04-11', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-10 04:10:00', 0);
INSERT INTO `sales_order` VALUES (376, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000358', '', 1, '2023-04-10', '2023-04-11', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-10 05:04:57', 0);
INSERT INTO `sales_order` VALUES (377, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000359', '', 1, '2023-04-10', '2023-04-11', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-10 06:48:38', 0);
INSERT INTO `sales_order` VALUES (378, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 'SO000360', '', 1, '2023-04-10', '2023-04-11', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 50000.00, 60000.00, 0.00, '', 0, 0, '2023-04-10 07:46:04', 0);
INSERT INTO `sales_order` VALUES (379, 1, 1, 0, 0, 0, 0, 1, 1, 55, 0, 0, 0, 0, 1, 0, 'SO000361', '', 1, '2023-04-13', '2023-04-14', '', 'Dummy ShippingAddress', '', 10000.00, '', 0.00, 0.00, 0, 0, 0.00, 10000.00, 20000.00, 0.00, '', 0, 0, '2023-04-13 02:55:07', 0);

-- ----------------------------
-- Table structure for sales_order_feedback
-- ----------------------------
DROP TABLE IF EXISTS `sales_order_feedback`;
CREATE TABLE `sales_order_feedback`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `sales_order_id` bigint UNSIGNED NOT NULL,
  `customer_id` bigint UNSIGNED NOT NULL,
  `sales_order_code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `delivery_date` date NULL DEFAULT NULL,
  `rating_score` tinyint UNSIGNED NULL DEFAULT NULL,
  `tags` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `description` varchar(250) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `to_be_contacted` tinyint UNSIGNED NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `ui_sales_order_feedback_idx`(`sales_order_id` ASC) USING BTREE,
  INDEX `sales_order_feedback_2_idx`(`customer_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 188 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sales_order_feedback
-- ----------------------------
INSERT INTO `sales_order_feedback` VALUES (3, 1001, 1, 'DUMMY-SO1001', '2023-01-01', 0, '', '', 0, '2023-03-24 02:29:47');
INSERT INTO `sales_order_feedback` VALUES (101, 1002, 1, 'DUMMY-SO1002', '2023-01-01', 4, '', 'test postman', 1, '2023-03-23 16:01:42');
INSERT INTO `sales_order_feedback` VALUES (102, 1, 1, 'DUMMY-SO1', '2023-01-01', 4, '', 'test postman', 1, '2023-03-24 06:22:16');
INSERT INTO `sales_order_feedback` VALUES (151, 1003, 1, 'DUMMY-SO1003', '2023-01-01', 9, '', 'test postman', 1, '2023-03-24 08:03:46');
INSERT INTO `sales_order_feedback` VALUES (152, 1004, 1, 'DUMMY-SO1004', '2023-01-01', 9, '', 'test postman', 1, '2023-03-24 08:03:52');
INSERT INTO `sales_order_feedback` VALUES (153, 1005, 1, 'DUMMY-SO1005', '2023-01-01', 9, '', 'test postman', 1, '2023-03-24 08:03:58');
INSERT INTO `sales_order_feedback` VALUES (155, 1007, 1, 'DUMMY-SO107', '2023-01-01', 0, '', '', 0, '2023-03-24 09:43:30');
INSERT INTO `sales_order_feedback` VALUES (161, 1006, 1, 'DUMMY-SO106', '2023-01-01', 9, '', 'test postman', 1, '2023-03-24 09:52:55');
INSERT INTO `sales_order_feedback` VALUES (163, 1008, 1, 'DUMMY-SO108', '2023-01-01', 0, '', '', 0, '2023-03-24 09:55:27');
INSERT INTO `sales_order_feedback` VALUES (165, 1009, 1, 'DUMMY-SO109', '2023-01-01', 10, '', 'testing diyana', 2, '2023-03-24 11:47:57');
INSERT INTO `sales_order_feedback` VALUES (167, 1011, 1, 'DUMMY-SO111', '2023-01-01', 10, '', 'testing diyana you want me home a few so I don\'t know if you page and it is not that you care so kan ke sprint so I don\'t know ', 2, '2023-03-24 11:49:23');
INSERT INTO `sales_order_feedback` VALUES (169, 1012, 1, 'DUMMY-SO112', '2023-01-01', 7, '', 'twsting Diana and Roma tomatoes and onions and peppers and onions and peppers ', 2, '2023-03-24 11:50:02');
INSERT INTO `sales_order_feedback` VALUES (171, 1013, 1, 'DUMMY-SO113', '2023-01-01', 8, '', '', 2, '2023-03-24 13:21:07');
INSERT INTO `sales_order_feedback` VALUES (172, 1014, 1, 'DUMMY-SO114', '2023-01-01', 10, '', 'the only one that has to be done by you page aman sih harusnya 2 backlog of the day today I love you too baby girl I love you too baby girl I love you too baby girl I love you too baby girl I love you too baby girl I love you too baby girl I love you', 2, '2023-03-24 13:54:29');
INSERT INTO `sales_order_feedback` VALUES (173, 1015, 1, 'DUMMY-SO115', '2023-01-01', 5, '', 'resting and I will be there in about an hour or so kan jam and yang di video of the kids are in school and then I can go to work and get the rest of the day off and I can get you a new one for the kids to the park and walk around the house and I will', 2, '2023-03-24 13:55:25');
INSERT INTO `sales_order_feedback` VALUES (174, 1016, 1, 'DUMMY-SO116', '2023-01-01', 10, '', 'testing diyna', 2, '2023-03-29 05:44:08');
INSERT INTO `sales_order_feedback` VALUES (175, 1017, 1, 'DUMMY-SO1017', '2023-01-01', 1, '', 'barang jelek', 2, '2023-03-29 07:30:46');
INSERT INTO `sales_order_feedback` VALUES (176, 1018, 1, 'DUMMY-SO1018', '2023-01-01', 9, '', 'Mantap', 2, '2023-03-29 10:22:08');
INSERT INTO `sales_order_feedback` VALUES (177, 1019, 1, 'DUMMY-SO1019', '2023-01-01', 9, '', 'Mantap', 2, '2023-03-29 10:25:54');
INSERT INTO `sales_order_feedback` VALUES (178, 1021, 1, 'DUMMY-SO1021', '2023-01-01', 9, '', 'Sangat bagus', 2, '2023-03-29 10:28:11');
INSERT INTO `sales_order_feedback` VALUES (179, 1022, 1, 'DUMMY-SO1022', '2023-01-01', 8, '', 'Bagus', 2, '2023-03-29 10:29:36');
INSERT INTO `sales_order_feedback` VALUES (180, 1023, 1, 'DUMMY-SO1023', '2023-01-01', 8, '', 'testing', 2, '2023-03-29 10:49:04');
INSERT INTO `sales_order_feedback` VALUES (181, 1024, 1, 'DUMMY-SO1024', '2023-01-01', 10, '', 'Sangat bagus', 2, '2023-03-29 10:51:16');
INSERT INTO `sales_order_feedback` VALUES (182, 5, 1, 'DUMMY-SO1', '2023-01-01', 8, '', 'Mantap', 2, '2023-03-29 10:52:29');
INSERT INTO `sales_order_feedback` VALUES (183, 4, 1, 'DUMMY-SO1', '2023-01-01', 9, '', 'Cocok dong', 2, '2023-03-29 10:52:46');
INSERT INTO `sales_order_feedback` VALUES (184, 10, 1, 'DUMMY-SO1', '2023-01-01', 10, '', 'testing akhir', 2, '2023-03-29 11:02:36');
INSERT INTO `sales_order_feedback` VALUES (185, 2, 1, 'DUMMY-SO1', '2023-01-01', 2, '', 'barang jelek', 2, '2023-03-30 03:28:32');
INSERT INTO `sales_order_feedback` VALUES (186, 3, 1, 'DUMMY-SO1', '2023-01-01', 10, '', 'barang bagus, kurirnya cepat', 2, '2023-03-30 03:29:02');
INSERT INTO `sales_order_feedback` VALUES (187, 9, 1, 'DUMMY-SO1', '2023-01-01', 10, '', 'fteagb', 2, '2023-04-03 08:02:42');

-- ----------------------------
-- Table structure for sales_order_item
-- ----------------------------
DROP TABLE IF EXISTS `sales_order_item`;
CREATE TABLE `sales_order_item`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `sales_order_id` bigint UNSIGNED NULL DEFAULT NULL,
  `item_id` bigint UNSIGNED NULL DEFAULT NULL,
  `item_discount_id` bigint UNSIGNED NULL DEFAULT NULL,
  `order_qty` decimal(10, 2) NULL DEFAULT 0.00,
  `discount_qty` decimal(10, 2) NULL DEFAULT 0.00,
  `forecast_qty` decimal(10, 2) NULL DEFAULT 0.00 COMMENT 'if != with order_qty = order > remaining forecast qty',
  `default_price` decimal(12, 2) NULL DEFAULT 0.00,
  `unit_price` decimal(12, 2) NULL DEFAULT 0.00,
  `unit_price_discount` decimal(10, 2) NULL DEFAULT 0.00,
  `item_disc_amount` decimal(10, 2) NULL DEFAULT 0.00,
  `taxable_item` tinyint(1) NULL DEFAULT 2 COMMENT 'flagging for the item is taxable or not',
  `tax_percentage` decimal(10, 2) NULL DEFAULT 0.00,
  `shadow_price` decimal(12, 2) NULL DEFAULT 0.00,
  `subtotal` decimal(15, 2) NULL DEFAULT 0.00,
  `weight` decimal(10, 2) NULL DEFAULT 0.00,
  `note` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `item_push` tinyint(1) NULL DEFAULT 2 COMMENT 'flag product push in market',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `fk_sales_order_item_1_idx`(`sales_order_id` ASC) USING BTREE,
  CONSTRAINT `fk_sales_order_item_1` FOREIGN KEY (`sales_order_id`) REFERENCES `sales_order` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 786 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of sales_order_item
-- ----------------------------
INSERT INTO `sales_order_item` VALUES (33, 2, 176, NULL, 24.00, 0.00, 0.00, 0.00, 21417.00, 0.00, 0.00, 2, 0.00, 0.00, 559903.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (34, 2, 227, NULL, 84.00, 0.00, 0.00, 0.00, 49802.00, 0.00, 0.00, 2, 0.00, 0.00, 924890.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (35, 2, 134, NULL, 52.00, 0.00, 0.00, 0.00, 83759.00, 0.00, 0.00, 2, 0.00, 0.00, 934740.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (36, 2, 71, NULL, 7.00, 0.00, 0.00, 0.00, 68391.00, 0.00, 0.00, 2, 0.00, 0.00, 889032.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (37, 2, 156, NULL, 76.00, 0.00, 0.00, 0.00, 89678.00, 0.00, 0.00, 2, 0.00, 0.00, 630940.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (38, 2, 180, NULL, 60.00, 0.00, 0.00, 0.00, 42216.00, 0.00, 0.00, 2, 0.00, 0.00, 477604.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (39, 1, 32, NULL, 73.00, 0.00, 0.00, 0.00, 41047.00, 0.00, 0.00, 2, 0.00, 0.00, 485203.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (40, 1, 18, NULL, 86.00, 0.00, 0.00, 0.00, 77590.00, 0.00, 0.00, 2, 0.00, 0.00, 983201.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (41, 1, 22, NULL, 10.00, 0.00, 0.00, 0.00, 63808.00, 0.00, 0.00, 2, 0.00, 0.00, 450396.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (42, 1, 24, NULL, 91.00, 0.00, 0.00, 0.00, 85270.00, 0.00, 0.00, 2, 0.00, 0.00, 292377.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (43, 1, 37, NULL, 28.00, 0.00, 0.00, 0.00, 33928.00, 0.00, 0.00, 2, 0.00, 0.00, 100697.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (44, 1, 44, NULL, 63.00, 0.00, 0.00, 0.00, 12828.00, 0.00, 0.00, 2, 0.00, 0.00, 616354.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (45, 1, 50, NULL, 33.00, 0.00, 0.00, 0.00, 61358.00, 0.00, 0.00, 2, 0.00, 0.00, 769682.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (46, 1, 150, NULL, 73.00, 0.00, 0.00, 0.00, 67307.00, 0.00, 0.00, 2, 0.00, 0.00, 989349.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (47, 1, 262, NULL, 69.00, 0.00, 0.00, 0.00, 51461.00, 0.00, 0.00, 2, 0.00, 0.00, 627698.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (48, 1, 108, NULL, 26.00, 0.00, 0.00, 0.00, 54384.00, 0.00, 0.00, 2, 0.00, 0.00, 160445.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (49, 1, 176, NULL, 24.00, 0.00, 0.00, 0.00, 16539.00, 0.00, 0.00, 2, 0.00, 0.00, 909130.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (50, 1, 76, NULL, 39.00, 0.00, 0.00, 0.00, 18542.00, 0.00, 0.00, 2, 0.00, 0.00, 54318.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (51, 1, 73, NULL, 22.00, 0.00, 0.00, 0.00, 42093.00, 0.00, 0.00, 2, 0.00, 0.00, 534200.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (52, 1, 145, NULL, 93.00, 0.00, 0.00, 0.00, 53840.00, 0.00, 0.00, 2, 0.00, 0.00, 498044.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (53, 1, 143, NULL, 4.00, 0.00, 0.00, 0.00, 41923.00, 0.00, 0.00, 2, 0.00, 0.00, 877622.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (54, 1, 175, NULL, 38.00, 0.00, 0.00, 0.00, 47093.00, 0.00, 0.00, 2, 0.00, 0.00, 883978.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (55, 3, 156, NULL, 76.00, 0.00, 0.00, 0.00, 8699.00, 0.00, 0.00, 2, 0.00, 0.00, 777024.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (56, 3, 71, NULL, 68.00, 0.00, 0.00, 0.00, 1217.00, 0.00, 0.00, 2, 0.00, 0.00, 223186.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (57, 3, 134, NULL, 15.00, 0.00, 0.00, 0.00, 78989.00, 0.00, 0.00, 2, 0.00, 0.00, 774858.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (58, 3, 180, NULL, 66.00, 0.00, 0.00, 0.00, 90293.00, 0.00, 0.00, 2, 0.00, 0.00, 194734.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (59, 3, 143, NULL, 87.00, 0.00, 0.00, 0.00, 13498.00, 0.00, 0.00, 2, 0.00, 0.00, 639097.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (60, 3, 144, NULL, 40.00, 0.00, 0.00, 0.00, 95610.00, 0.00, 0.00, 2, 0.00, 0.00, 601284.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (61, 3, 68, NULL, 35.00, 0.00, 0.00, 0.00, 36559.00, 0.00, 0.00, 2, 0.00, 0.00, 79129.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (62, 3, 35, NULL, 56.00, 0.00, 0.00, 0.00, 94963.00, 0.00, 0.00, 2, 0.00, 0.00, 581793.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (63, 4, 134, NULL, 77.00, 0.00, 0.00, 0.00, 64139.00, 0.00, 0.00, 2, 0.00, 0.00, 661577.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (64, 4, 151, NULL, 18.00, 0.00, 0.00, 0.00, 34806.00, 0.00, 0.00, 2, 0.00, 0.00, 552507.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (65, 4, 38, NULL, 55.00, 0.00, 0.00, 0.00, 80615.00, 0.00, 0.00, 2, 0.00, 0.00, 767804.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (66, 4, 81, NULL, 25.00, 0.00, 0.00, 0.00, 97655.00, 0.00, 0.00, 2, 0.00, 0.00, 171500.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (67, 5, 50, NULL, 56.00, 0.00, 0.00, 0.00, 45432.00, 0.00, 0.00, 2, 0.00, 0.00, 544087.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (70, 6, 46, NULL, 9.00, 0.00, 0.00, 0.00, 33195.00, 0.00, 0.00, 2, 0.00, 0.00, 195935.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (71, 6, 43, NULL, 72.00, 0.00, 0.00, 0.00, 28680.00, 0.00, 0.00, 2, 0.00, 0.00, 337416.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (72, 7, 242, NULL, 36.00, 0.00, 0.00, 0.00, 42814.00, 0.00, 0.00, 2, 0.00, 0.00, 89274.00, 1.00, 'Lorem Ipsum', 2);
INSERT INTO `sales_order_item` VALUES (73, 7, 47, NULL, 64.00, 0.00, 0.00, 0.00, 27030.00, 0.00, 0.00, 2, 0.00, 0.00, 424122.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (74, 7, 37, NULL, 14.00, 0.00, 0.00, 0.00, 5711.00, 0.00, 0.00, 2, 0.00, 0.00, 842791.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (75, 8, 22, NULL, 77.00, 0.00, 0.00, 0.00, 46466.00, 0.00, 0.00, 2, 0.00, 0.00, 931587.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (76, 8, 24, NULL, 42.00, 0.00, 0.00, 0.00, 14197.00, 0.00, 0.00, 2, 0.00, 0.00, 119565.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (77, 8, 18, NULL, 78.00, 0.00, 0.00, 0.00, 30587.00, 0.00, 0.00, 2, 0.00, 0.00, 793063.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (78, 8, 37, NULL, 68.00, 0.00, 0.00, 0.00, 9345.00, 0.00, 0.00, 2, 0.00, 0.00, 596621.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (79, 8, 44, NULL, 5.00, 0.00, 0.00, 0.00, 53964.00, 0.00, 0.00, 2, 0.00, 0.00, 593916.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (80, 8, 46, NULL, 21.00, 0.00, 0.00, 0.00, 40785.00, 0.00, 0.00, 2, 0.00, 0.00, 169718.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (81, 8, 319, NULL, 87.00, 0.00, 0.00, 0.00, 41036.00, 0.00, 0.00, 2, 0.00, 0.00, 56842.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (82, 8, 150, NULL, 73.00, 0.00, 0.00, 0.00, 81822.00, 0.00, 0.00, 2, 0.00, 0.00, 765056.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (83, 8, 80, NULL, 8.00, 0.00, 0.00, 0.00, 85005.00, 0.00, 0.00, 2, 0.00, 0.00, 644755.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (84, 8, 69, NULL, 15.00, 0.00, 0.00, 0.00, 78560.00, 0.00, 0.00, 2, 0.00, 0.00, 918608.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (85, 8, 176, NULL, 54.00, 0.00, 0.00, 0.00, 36783.00, 0.00, 0.00, 2, 0.00, 0.00, 648775.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (86, 8, 76, NULL, 24.00, 0.00, 0.00, 0.00, 47236.00, 0.00, 0.00, 2, 0.00, 0.00, 478054.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (87, 8, 73, NULL, 59.00, 0.00, 0.00, 0.00, 24833.00, 0.00, 0.00, 2, 0.00, 0.00, 433942.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (88, 8, 145, NULL, 23.00, 0.00, 0.00, 0.00, 81456.00, 0.00, 0.00, 2, 0.00, 0.00, 725549.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (89, 8, 164, NULL, 39.00, 0.00, 0.00, 0.00, 31782.00, 0.00, 0.00, 2, 0.00, 0.00, 315922.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (90, 8, 162, NULL, 24.00, 0.00, 0.00, 0.00, 13542.00, 0.00, 0.00, 2, 0.00, 0.00, 392961.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (91, 8, 187, NULL, 3.00, 0.00, 0.00, 0.00, 71365.00, 0.00, 0.00, 2, 0.00, 0.00, 1007040.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (92, 9, 242, NULL, 39.00, 0.00, 0.00, 0.00, 15200.00, 0.00, 0.00, 2, 0.00, 0.00, 846316.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (93, 9, 37, NULL, 90.00, 0.00, 0.00, 0.00, 60904.00, 0.00, 0.00, 2, 0.00, 0.00, 200463.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (94, 9, 297, NULL, 36.00, 0.00, 0.00, 0.00, 57922.00, 0.00, 0.00, 2, 0.00, 0.00, 453367.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (95, 9, 116, NULL, 6.00, 0.00, 0.00, 0.00, 5899.00, 0.00, 0.00, 2, 0.00, 0.00, 655445.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (96, 9, 162, NULL, 22.00, 0.00, 0.00, 0.00, 54731.00, 0.00, 0.00, 2, 0.00, 0.00, 907127.00, 1.00, '', 2);
INSERT INTO `sales_order_item` VALUES (97, 10, 132, NULL, 93.00, 0.00, 0.00, 0.00, 54956.00, 0.00, 0.00, 2, 0.00, 0.00, 559298.00, 1.00, 'Lorem Ipsum', 2);
INSERT INTO `sales_order_item` VALUES (98, 11, 1, NULL, 10.00, 0.00, 0.00, 0.00, 10000.00, 0.00, 0.00, 0, 0.00, 10000.00, 100000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (99, 12, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (100, 13, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (101, 14, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (102, 15, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (103, 16, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (104, 17, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (105, 18, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (106, 19, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (107, 20, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (108, 21, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (109, 21, 1, NULL, 2.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 10000.00, 2.00, '', 0);
INSERT INTO `sales_order_item` VALUES (110, 21, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (111, 22, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (112, 22, 1, NULL, 2.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 10000.00, 2.00, '', 0);
INSERT INTO `sales_order_item` VALUES (113, 22, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (114, 23, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (115, 24, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (116, 25, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (117, 26, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (118, 27, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (119, 28, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (120, 29, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (121, 30, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (122, 31, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (123, 32, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (124, 33, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (125, 34, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (126, 35, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (127, 36, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (128, 37, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (129, 38, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (130, 38, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (131, 39, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (132, 40, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (133, 40, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (134, 41, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (135, 41, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (136, 42, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (137, 42, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (138, 43, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (139, 44, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (140, 45, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (141, 46, 1, NULL, 3.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 15000.00, 3.00, '', 0);
INSERT INTO `sales_order_item` VALUES (142, 47, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (143, 48, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (144, 49, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (145, 50, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (146, 51, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (147, 52, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (148, 53, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (149, 54, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (150, 55, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (151, 56, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (152, 57, 1, NULL, 50.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 250000.00, 50.00, '', 0);
INSERT INTO `sales_order_item` VALUES (153, 58, 1, NULL, 50.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 250000.00, 50.00, '', 0);
INSERT INTO `sales_order_item` VALUES (154, 59, 1, NULL, 50.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 250000.00, 50.00, '', 0);
INSERT INTO `sales_order_item` VALUES (155, 60, 1, NULL, 50.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 250000.00, 50.00, '', 0);
INSERT INTO `sales_order_item` VALUES (156, 61, 1, NULL, 50.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 250000.00, 50.00, '', 0);
INSERT INTO `sales_order_item` VALUES (157, 62, 1, NULL, 50.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 250000.00, 50.00, '', 0);
INSERT INTO `sales_order_item` VALUES (158, 63, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (159, 64, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (160, 65, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (161, 66, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (162, 67, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (163, 68, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (164, 69, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (165, 70, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (166, 71, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (167, 72, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (168, 73, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (169, 74, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (170, 75, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (171, 76, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (172, 77, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (173, 78, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (174, 79, 1, NULL, 6.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 30000.00, 6.00, '', 0);
INSERT INTO `sales_order_item` VALUES (175, 80, 1, NULL, 4.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 20000.00, 4.00, '', 0);
INSERT INTO `sales_order_item` VALUES (176, 81, 1, NULL, 4.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 20000.00, 4.00, '', 0);
INSERT INTO `sales_order_item` VALUES (177, 82, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (178, 83, 1, NULL, 4.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 20000.00, 4.00, '', 0);
INSERT INTO `sales_order_item` VALUES (179, 84, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (180, 85, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (181, 86, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (182, 87, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (183, 87, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (184, 87, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (185, 88, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (186, 89, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (187, 90, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (188, 91, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (189, 92, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (190, 93, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (191, 94, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (192, 95, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (193, 96, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (194, 97, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (195, 98, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (196, 99, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (197, 100, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (198, 101, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (199, 102, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (200, 103, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (201, 104, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (202, 105, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (203, 106, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (204, 107, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (205, 108, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (206, 109, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (207, 110, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, '', 0);
INSERT INTO `sales_order_item` VALUES (208, 111, 1, NULL, 7.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 35000.00, 7.00, '', 0);
INSERT INTO `sales_order_item` VALUES (209, 112, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (210, 113, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (211, 114, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (212, 115, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (213, 116, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (214, 117, 1, NULL, 2.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 10000.00, 2.00, '', 0);
INSERT INTO `sales_order_item` VALUES (215, 118, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (216, 119, 1, NULL, 2.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 10000.00, 2.00, '', 0);
INSERT INTO `sales_order_item` VALUES (217, 120, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (218, 121, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (219, 122, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (220, 123, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (221, 124, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (222, 125, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (223, 126, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, '', 0);
INSERT INTO `sales_order_item` VALUES (224, 127, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, '', 0);
INSERT INTO `sales_order_item` VALUES (225, 128, 1, NULL, 2.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 10000.00, 2.00, '', 0);
INSERT INTO `sales_order_item` VALUES (226, 129, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, 'test', 0);
INSERT INTO `sales_order_item` VALUES (227, 129, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, 'test', 0);
INSERT INTO `sales_order_item` VALUES (228, 129, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, 'test', 0);
INSERT INTO `sales_order_item` VALUES (229, 129, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, 'test', 0);
INSERT INTO `sales_order_item` VALUES (230, 129, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, 'yest', 0);
INSERT INTO `sales_order_item` VALUES (231, 130, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (232, 130, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (233, 130, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (234, 131, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (235, 131, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (236, 131, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (237, 131, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (238, 132, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (239, 133, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (240, 133, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, '', 0);
INSERT INTO `sales_order_item` VALUES (241, 133, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (242, 133, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (243, 134, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (244, 134, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (245, 134, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (246, 134, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (247, 135, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (248, 135, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (249, 135, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (250, 135, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (251, 136, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (252, 136, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (253, 136, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (254, 136, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (255, 137, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (256, 137, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (257, 137, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (258, 137, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (259, 138, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (260, 138, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (261, 138, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (262, 139, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (263, 139, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (264, 139, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (265, 139, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (266, 140, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (267, 140, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (268, 140, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (269, 140, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (270, 141, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (271, 141, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (272, 142, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (273, 142, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (274, 142, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (275, 142, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (276, 143, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (277, 143, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (278, 143, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (279, 143, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (280, 144, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (281, 144, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (282, 144, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (283, 144, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (284, 145, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (285, 145, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (286, 146, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (287, 146, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (288, 147, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (289, 147, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (290, 147, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (291, 147, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (292, 148, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (293, 148, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (294, 148, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (295, 148, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (296, 148, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (297, 149, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (298, 149, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (299, 149, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (300, 150, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (301, 150, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (302, 151, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (303, 151, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (304, 151, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (305, 151, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (306, 152, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (307, 152, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (308, 152, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (309, 152, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (310, 153, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (311, 153, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (312, 153, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (313, 153, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (314, 154, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (315, 154, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (316, 154, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (317, 154, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (318, 155, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (319, 155, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (320, 155, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (321, 155, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (322, 156, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (323, 156, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (324, 156, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (325, 156, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (326, 157, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (327, 157, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (328, 157, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (329, 157, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (330, 158, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (331, 158, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (332, 158, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (333, 158, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (334, 159, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (335, 159, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (336, 159, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (337, 159, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (338, 160, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (339, 160, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (340, 160, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (341, 160, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (342, 161, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (343, 162, 1, NULL, 3.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 15000.00, 3.00, '', 0);
INSERT INTO `sales_order_item` VALUES (344, 163, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (345, 163, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (346, 164, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (347, 164, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (348, 165, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (349, 165, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (350, 165, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (351, 165, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (352, 166, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (353, 166, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (354, 166, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (355, 166, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (356, 167, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (357, 167, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (358, 167, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (359, 167, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (360, 168, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (361, 168, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (362, 168, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (363, 168, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (364, 169, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (365, 169, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (366, 169, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (367, 169, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (368, 170, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (369, 170, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (370, 170, 8, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (371, 171, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (372, 171, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (373, 171, 8, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (374, 172, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (375, 172, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (376, 172, 8, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (377, 173, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (378, 173, 6, NULL, 3.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 15000.00, 3.00, '', 0);
INSERT INTO `sales_order_item` VALUES (379, 173, 8, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (380, 174, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (381, 174, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (382, 174, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (383, 174, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (384, 175, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (385, 175, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (386, 175, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (387, 175, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (388, 176, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (389, 176, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (390, 176, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (391, 176, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (392, 177, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (393, 177, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (394, 177, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (395, 177, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (396, 178, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (397, 178, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (398, 178, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (399, 178, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (400, 179, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (401, 179, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (402, 179, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (403, 179, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (404, 180, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (405, 180, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (406, 180, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (407, 180, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (408, 181, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (409, 181, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (410, 182, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (411, 182, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (412, 182, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (413, 182, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (414, 183, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (415, 183, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (416, 183, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (417, 183, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (418, 184, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (419, 184, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (420, 184, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (421, 184, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (422, 185, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (423, 185, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (424, 185, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (425, 185, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (426, 186, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (427, 186, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (428, 186, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (429, 186, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (430, 187, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (431, 187, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (432, 187, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (433, 187, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (434, 188, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (435, 188, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (436, 188, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (437, 188, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (438, 189, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (439, 189, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (440, 189, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (441, 189, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (442, 190, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (443, 190, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (444, 190, 5, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (445, 190, 6, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (446, 191, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (447, 191, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (448, 191, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (449, 191, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (450, 192, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (451, 192, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (452, 192, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (453, 192, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (454, 193, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (455, 194, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (456, 194, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (457, 194, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (458, 195, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (459, 195, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (460, 195, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (461, 195, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (462, 196, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (463, 196, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (464, 196, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (465, 197, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (466, 197, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (467, 197, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (468, 197, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (469, 198, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (470, 198, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (471, 198, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (472, 198, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (473, 199, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (474, 199, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (475, 199, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (476, 199, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (477, 200, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (478, 200, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (479, 200, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (480, 200, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (481, 201, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (482, 201, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (483, 201, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (484, 201, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (485, 202, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (486, 202, 6, NULL, 2.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 10000.00, 2.00, '', 0);
INSERT INTO `sales_order_item` VALUES (487, 202, 7, NULL, 2.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 10000.00, 2.00, '', 0);
INSERT INTO `sales_order_item` VALUES (488, 203, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (489, 203, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (490, 203, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (491, 203, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (492, 204, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (493, 204, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (494, 204, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (495, 204, 3, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (496, 205, 1, NULL, 3.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 15000.00, 3.00, '', 0);
INSERT INTO `sales_order_item` VALUES (497, 206, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (498, 207, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (499, 207, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (500, 208, 1, NULL, 3.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 15000.00, 3.00, '', 0);
INSERT INTO `sales_order_item` VALUES (501, 208, 2, NULL, 3.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 15000.00, 3.00, '', 0);
INSERT INTO `sales_order_item` VALUES (502, 209, 1, NULL, 3.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 15000.00, 3.00, '', 0);
INSERT INTO `sales_order_item` VALUES (503, 209, 2, NULL, 3.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 15000.00, 3.00, '', 0);
INSERT INTO `sales_order_item` VALUES (504, 209, 3, NULL, 3.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 15000.00, 3.00, '', 0);
INSERT INTO `sales_order_item` VALUES (505, 210, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (506, 210, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (507, 211, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (508, 212, 1, NULL, 42.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 210000.00, 42.00, '', 0);
INSERT INTO `sales_order_item` VALUES (509, 213, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (510, 213, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (511, 214, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (512, 214, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (513, 215, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (514, 215, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (515, 216, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (516, 216, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (517, 217, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (518, 218, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (519, 218, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (520, 219, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (521, 219, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (522, 220, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (523, 220, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (524, 221, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (525, 221, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (526, 222, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (527, 222, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (528, 223, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (529, 223, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (530, 224, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (531, 224, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (532, 225, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (533, 225, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (534, 226, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (535, 226, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (536, 227, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (537, 227, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (538, 228, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (539, 228, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (540, 229, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (541, 229, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (542, 230, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (543, 230, 2, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (544, 231, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (545, 231, 2, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (546, 232, 2, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (547, 232, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (548, 233, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (549, 234, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (550, 235, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (551, 235, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (552, 236, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (553, 237, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (554, 237, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (555, 238, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (556, 238, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (557, 239, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (558, 239, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (559, 240, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (560, 240, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (561, 241, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (562, 241, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (563, 242, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (564, 242, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (565, 243, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (566, 243, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (567, 244, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (568, 244, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (569, 245, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (570, 245, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (571, 246, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (572, 246, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (573, 247, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (574, 247, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (575, 248, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (576, 248, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (577, 249, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (578, 249, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (579, 250, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (580, 250, 4, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (581, 251, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (582, 251, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (583, 252, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (584, 252, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (585, 253, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (586, 253, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (587, 254, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (588, 254, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (589, 255, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (590, 255, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (591, 256, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (592, 256, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (593, 257, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (594, 257, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (595, 258, 1, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 10.00, '', 0);
INSERT INTO `sales_order_item` VALUES (596, 259, 1, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (597, 259, 2, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 5.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (598, 260, 2, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (599, 260, 1, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 1.00, '', 0);
INSERT INTO `sales_order_item` VALUES (600, 261, 182, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (601, 261, 183, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (602, 262, 186, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (603, 263, 182, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (604, 263, 183, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (605, 266, 233, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (606, 267, 233, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (607, 267, 234, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (608, 268, 233, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (609, 268, 234, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (610, 269, 243, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (611, 269, 244, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (612, 270, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (613, 270, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (614, 271, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (615, 271, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (616, 272, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (617, 272, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (618, 273, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (619, 273, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (620, 274, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (621, 274, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (622, 275, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (623, 275, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (624, 276, 304, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (625, 277, 304, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (626, 278, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (627, 278, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (628, 279, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (629, 279, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (630, 280, 305, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (631, 281, 304, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (632, 282, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (633, 282, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (634, 283, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (635, 283, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (636, 284, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (637, 285, 304, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (638, 286, 304, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (639, 287, 304, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (640, 288, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (641, 288, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (642, 289, 305, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (643, 290, 305, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (644, 291, 305, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (645, 292, 305, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (646, 293, 304, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (647, 294, 305, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (648, 295, 305, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (649, 296, 305, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (650, 297, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (651, 297, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (652, 298, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (653, 298, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (654, 299, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (655, 299, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (656, 300, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (657, 300, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (658, 301, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (659, 301, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (660, 302, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (661, 302, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (662, 303, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (663, 303, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (664, 304, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (665, 304, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (666, 305, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (667, 305, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (668, 306, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (669, 306, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (670, 307, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (671, 307, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (672, 308, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (673, 308, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (674, 309, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (675, 309, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (676, 310, 304, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (677, 310, 305, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (678, 311, 304, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (679, 311, 305, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (680, 312, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (681, 312, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (682, 313, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (683, 313, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (684, 314, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (685, 314, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (686, 315, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (687, 315, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (688, 316, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (689, 316, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (690, 317, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (691, 317, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (692, 318, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (693, 318, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (694, 319, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (695, 319, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (696, 320, 318, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (697, 321, 318, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (698, 322, 318, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (699, 323, 318, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (700, 324, 318, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (701, 325, 318, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (702, 326, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (703, 326, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (704, 327, 318, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (705, 328, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (706, 328, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (707, 329, 318, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (708, 330, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (709, 330, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (710, 331, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (711, 331, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (712, 332, 315, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (713, 333, 315, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (714, 334, 315, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (715, 335, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (716, 335, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (717, 336, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (718, 336, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (719, 337, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (720, 337, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (721, 338, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (722, 338, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (723, 339, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (724, 339, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (725, 340, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (726, 340, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (727, 341, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (728, 341, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (729, 342, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (730, 342, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (731, 343, 315, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (732, 344, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (733, 345, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (734, 346, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (735, 347, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (736, 348, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (737, 349, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (738, 350, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (739, 351, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (740, 352, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (741, 352, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (742, 353, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (743, 353, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (744, 354, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (745, 355, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (746, 355, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (747, 356, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (748, 356, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (749, 357, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (750, 357, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (751, 358, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (752, 359, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (753, 360, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (754, 360, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (755, 361, 314, NULL, 11.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 55000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (756, 362, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (757, 363, 314, NULL, 100.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 500000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (758, 364, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (759, 365, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (760, 365, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (761, 366, 314, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (762, 367, 314, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (763, 368, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (764, 368, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (765, 369, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (766, 369, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (767, 370, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (768, 370, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (769, 371, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (770, 371, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (771, 372, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (772, 372, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (773, 373, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (774, 373, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (775, 374, 315, NULL, 10.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 50000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (776, 375, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (777, 375, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (778, 376, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (779, 376, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (780, 377, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (781, 377, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (782, 378, 314, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 1', 0);
INSERT INTO `sales_order_item` VALUES (783, 378, 315, NULL, 5.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 25000.00, 0.00, 'item ke 2', 0);
INSERT INTO `sales_order_item` VALUES (784, 379, 315, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 0.00, '', 0);
INSERT INTO `sales_order_item` VALUES (785, 379, 318, NULL, 1.00, 0.00, 0.00, 0.00, 5000.00, 0.00, 0.00, 0, 0.00, 5000.00, 5000.00, 0.00, '', 0);

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
INSERT INTO `schema_migrations` VALUES (6, 0);

SET FOREIGN_KEY_CHECKS = 1;
