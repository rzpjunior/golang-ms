/*
 Navicat Premium Data Transfer

 Source Server         : Dev Eden ERP-Write
 Source Server Type    : MySQL
 Source Server Version : 80031 (8.0.31-google)
 Source Host           : 10.26.160.2:3306
 Source Schema         : dynamic

 Target Server Type    : MySQL
 Target Server Version : 80031 (8.0.31-google)
 File Encoding         : 65001

 Date: 05/05/2023 12:28:04
*/
USE dynamic;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for address
-- ----------------------------
DROP TABLE IF EXISTS `address`;
CREATE TABLE `address`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'customer id/number from GP | CSTNMBR',
  `customer_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'address id from GP',
  `archetype_id` bigint NULL DEFAULT NULL COMMENT 'archetype from GP',
  `adm_division_id` bigint NULL DEFAULT NULL COMMENT 'administrative division from GP',
  `site_id` bigint NULL DEFAULT NULL COMMENT 'site id from GP',
  `salesperson_id` bigint NULL DEFAULT NULL,
  `territory_id` bigint NULL DEFAULT NULL COMMENT 'territory id from GP',
  `address_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'address code from GP | ADRSCODE',
  `address_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'address name from GP',
  `contact_person` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'contact person from GP | CNTCPRSN',
  `city` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'city from GP',
  `state` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'state from GP',
  `zip_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'zip_code from GP',
  `country_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'country code from GP',
  `country` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'country from GP',
  `latitude` decimal(17, 14) NULL DEFAULT NULL COMMENT 'latitude from GP',
  `longitude` decimal(17, 14) NULL DEFAULT NULL COMMENT 'longitude from GP',
  `ups_zone` varchar(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'ups zone from GP',
  `shipping_method` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'shipping method from GP | SHIPMTHD',
  `tax_schedule_id` bigint NULL DEFAULT NULL COMMENT 'tax schedule id from GP',
  `print_phone_number` tinyint(1) NULL DEFAULT 0 COMMENT 'phone number to be printed from GP | Print_Phone_NumberGB',
  `phone_1` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'phone number 1 from GP',
  `phone_2` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'phone number 1 from GP',
  `phone_3` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'phone number 1 from GP',
  `fax_number` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'fax number from GP',
  `shipping_address` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'shipping address from GP',
  `bca_va` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'BCA virtual account from GP',
  `other_va` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'other virtual account from GP',
  `note` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT 'note from GP',
  `status` tinyint(1) NULL DEFAULT 0 COMMENT 'status from GP',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'branch' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of address
-- ----------------------------
INSERT INTO `address` VALUES (1, 'ADR001', 'Bajuri', 1, 1, 1, 181, 1, 'Primary', 'Jalan sudirman no 111', 'Justin', 'Jakarta Pusat', 'DKI Jakarta', '10657', '+62', 'Indonesia', 123.12345678000000, 123.12345678000000, 'abc', 'REGULER', 1, 1, '08987233874', '', '', '', 'Jalan Sudirman No 111', '78955558', '12465678', '', 1, '2022-12-20 00:49:30', '2023-02-21 10:20:37');
INSERT INTO `address` VALUES (2, 'ADR002', 'Ucup', 3, 2, 1, 2, 1, 'Primary', 'Jalan Gelora no 44', 'Aidan', 'Jakarta Pusat', 'DKI Jakarta', '10658', '+62', 'Indonesia', 123.12345678000000, 123.12345678000000, 'abc', 'REGULER', 1, 1, '0813243435', '0822343432', '', '', 'Jalan Gelora No 44', '63487394', '3428842', '', 1, '2022-12-20 00:49:30', '2023-01-24 06:52:21');
INSERT INTO `address` VALUES (3, 'ADR003', 'Smith', 2, 1, 1, 3, 1, 'Primary', 'Jalan Benhil 9 no 44', 'Yepi', 'Jakarta Pusat', 'DKI Jakarta', '10578', '+62', 'Indonesia', 123.12345678000000, 123.12345678000000, 'abc', 'REGULER', 1, 1, '08578178332', '', '', '', 'Jalan Satrio No 44', '82798343', '92389324', '', 1, '2022-12-20 00:49:30', '2023-01-24 06:56:44');
INSERT INTO `address` VALUES (4, 'ADR004', 'Siti', 3, 2, 1, 4, 1, 'Primary', 'Jalan Perpustakaan no 33', 'Maryam', 'Jakarta Pusat', 'DKI Jakarta', '12399', '+62', 'Indonesia', 123.12345678000000, 123.12345678000000, 'abc', 'REGULER', 1, 1, '08772389890', '0812324234', '', '', 'Jalan Thamrin No 44', '23748274', '3435453', '', 1, '2022-12-20 00:49:30', '2023-01-24 06:56:44');
INSERT INTO `address` VALUES (5, 'ADR005', 'Bajuri', 4, 1, 1, 5, 1, 'Primary', 'Jalan gondangdia no 111', 'Abdul', 'Jakarta Pusat', 'DKI Jakarta', '10222', '+62', 'Indonesia', 123.12345678000000, 123.12345678000000, 'abc', 'REGULER', 1, 1, '08772323220', '', '', '', 'Jalan Kemuliaan No 111', '82392011', '23234344', '', 1, '2022-12-20 00:49:30', '2023-01-24 06:56:45');
INSERT INTO `address` VALUES (6, 'ADR006', 'Ronaldo', 3, 2, 1, 1, 1, 'Primary', 'Jalan sudirman no 111', 'Justin', 'Jakarta Pusat', 'DKI Jakarta', '10657', '+62', 'Indonesia', 123.12345678000000, 123.12345678000000, 'abc', 'REGULER', 1, 1, '08987233874', '', '', '', 'Jalan Sudirman No 111', '78955558', '12465678', '', 1, '2023-01-24 06:50:54', '2023-01-24 06:56:45');
INSERT INTO `address` VALUES (7, 'ADR007', 'Xavi', 2, 2, 1, 2, 1, 'Primary', 'Jalan Gelora no 44', 'Aidan', 'Jakarta Pusat', 'DKI Jakarta', '10658', '+62', 'Indonesia', 123.12345678000000, 123.12345678000000, 'abc', 'REGULER', 1, 1, '0813243435', '0822343432', '', '', 'Jalan Gelora No 44', '63487394', '3428842', '', 1, '2023-01-24 06:51:02', '2023-01-24 06:53:04');
INSERT INTO `address` VALUES (8, 'ADR008', 'Mbappe', 4, 1, 1, 3, 1, 'Primary', 'Jalan Benhil 9 no 44', 'Yepi', 'Jakarta Pusat', 'DKI Jakarta', '10578', '+62', 'Indonesia', 123.12345678000000, 123.12345678000000, 'abc', 'REGULER', 1, 1, '08578178332', '', '', '', 'Jalan Satrio No 44', '82798343', '92389324', '', 1, '2023-01-24 06:51:05', '2023-01-24 06:56:48');
INSERT INTO `address` VALUES (9, 'ADR009', 'Messi', 2, 1, 1, 4, 1, 'Primary', 'Jalan Perpustakaan no 33', 'Maryam', 'Jakarta Pusat', 'DKI Jakarta', '12399', '+62', 'Indonesia', 123.12345678000000, 123.12345678000000, 'abc', 'REGULER', 1, 1, '08772389890', '0812324234', '', '', 'Jalan Thamrin No 44', '23748274', '3435453', '', 1, '2023-01-24 06:51:07', '2023-01-24 06:56:49');
INSERT INTO `address` VALUES (10, 'ADR010', 'Banzema', 2, 2, 1, 5, 1, 'Primary', 'Jalan gondangdia no 111', 'Abdul', 'Jakarta Pusat', 'DKI Jakarta', '10222', '+62', 'Indonesia', 123.12345678000000, 123.12345678000000, 'abc', 'REGULER', 1, 1, '08772323220', '', '', '', 'Jalan Kemuliaan No 111', '82392011', '23234344', '', 1, '2023-01-24 06:51:12', '2023-01-24 06:56:49');

-- ----------------------------
-- Table structure for adm_division
-- ----------------------------
DROP TABLE IF EXISTS `adm_division`;
CREATE TABLE `adm_division`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'adm.division id from GP',
  `province_id` bigint NULL DEFAULT NULL,
  `city_id` bigint NULL DEFAULT NULL COMMENT 'city from GP | belum tau di gp ada masternya atau tidak',
  `sub_district_id` bigint NULL DEFAULT NULL COMMENT 'sub district from GP | belum tau di gp ada masternya atau tidak',
  `district_id` bigint NULL DEFAULT NULL COMMENT 'district from GP | belum tau di gp ada masternya atau tidak',
  `region_id` bigint NULL DEFAULT NULL COMMENT 'region id from GP',
  `postal_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'postal code from GP | belum tau di gp ada masternya atau tidak',
  `status` tinyint(1) NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `city` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `district` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `region` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `province` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'adm_division' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of adm_division
-- ----------------------------
INSERT INTO `adm_division` VALUES (1, 'ADD0001', 1, 1, 1, 1, 1, '5432101', 1, '2022-12-23 01:36:54', '2023-01-27 10:01:02', '1', '1', '1', '1');
INSERT INTO `adm_division` VALUES (2, 'ADD0001', 2, 2, 2, 2, 2, '5432102', 2, '2022-12-23 01:36:58', '2023-01-27 10:01:02', '1', '1', '1', '1');
INSERT INTO `adm_division` VALUES (3, 'ADD0001', 3, 3, 3, 3, 3, '5432103', 3, '2023-01-24 11:01:14', '2023-01-27 10:01:02', '1', '1', '1', '1');

-- ----------------------------
-- Table structure for archetype
-- ----------------------------
DROP TABLE IF EXISTS `archetype`;
CREATE TABLE `archetype`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'archetype id from GP',
  `business_type_id` bigint NULL DEFAULT NULL COMMENT 'business type id from GP',
  `description` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'description from GP',
  `status` tinyint(1) NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 100 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of archetype
-- ----------------------------
INSERT INTO `archetype` VALUES (1, 'Cafe', 2, 'Cafe', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (2, 'Restaurant', 2, 'Restoran', 23, '2022-12-20 00:50:32', '2023-02-14 10:31:14');
INSERT INTO `archetype` VALUES (3, 'Hotel', 2, 'Hotel', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (4, 'Catering', 2, 'Catering (Tempo)', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (5, 'Other (Horeca)', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (6, 'Warung Nasi', 3, 'Warung Nasi/Rumah Makan', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (7, 'Retail', 3, 'Toko sembako/Kelontong', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (8, 'Street Food', 3, 'Jajanan Makanan', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (9, 'Minuman', 3, 'Jajanan Minuman', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (10, 'Martabak', 3, 'Martabak', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (11, 'Warung Kopi', 3, 'Warung Kopi', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (12, 'Warung Sayur', 3, 'Warung Sayur & Buah', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (13, 'Other (Traditional Culinary B2B)', 3, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (14, 'Other (Partnership)', 4, 'Perusahaan Start-up', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (15, 'Other (Retail)', 5, '', 2, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (16, 'Reseller', 6, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (17, 'None', 7, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (18, 'Supplier (Traditional Culinary B2B)', 3, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (19, 'Supplier (Horeca)', 2, 'Supplier B2B', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (20, 'E-Commerce (Horeca)', 2, '', 2, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (21, 'Reseller', 3, '', 2, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (22, 'All Archetype', 1, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (23, 'all', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (24, 'all', 3, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (25, 'all', 4, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (26, 'all', 5, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (27, 'all', 6, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (28, 'Mitra Eden', 8, 'Pelapak Pasar', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (29, 'Retail Eden', 8, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (30, 'Modern Trade (Horeca)', 2, 'Supermarket / Minimarket', 2, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (31, 'all', 8, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (32, 'Personal', 9, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (33, 'Industry', 2, '', 2, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (34, 'Membership', 6, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (35, 'Rekan Eden', 6, 'Rekan Eden', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (36, 'Catering LM', 3, 'Catering (COD)', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (37, 'Pedagang Buah', 8, 'Pedagang Buah', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (38, 'E-Commerce', 4, 'E-commerce', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (39, 'Duta Institusi', 6, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (40, 'Membership Institusi', 6, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (41, 'Agen Telur', 3, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (42, 'Internal', 9, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (43, 'Supplier Funding', 11, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (44, 'Central Market Trade', 11, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (45, 'Importir', 11, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (46, 'Catering Enterprise', 10, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (47, 'Food Industry', 10, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (48, 'Exportir', 10, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (49, 'Corporate Enterprise', 10, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (50, 'Modern Trade', 10, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (51, 'Pesantren', 3, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (52, 'Agen Bawang', 3, '', 2, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (53, 'Agen COG/Veg Fresh', 3, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (54, 'Agen dalam pasar', 8, '', 2, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (55, 'Zero Waste', 12, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (56, 'Third Party', 12, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (57, 'Distributor Bawang dan Cabe', 13, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (58, 'Distributor Dry Goods', 13, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (59, 'Distributor Kentang', 13, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (60, 'Distributor Telur', 13, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (61, 'Lapak Dalam Pasar', 8, 'Lapak Dalam Pasar', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (62, 'Lapak Luar Pasar', 8, 'Lapak Luar Pasar', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (63, 'Other (Wet Market)', 8, 'Other (Wet Market)', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (64, 'Toko Kelontong', 8, 'Toko Kelontong', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (65, 'Toko Kelontong Dalam Pasar', 8, 'Toko Kelontong Dalam Pasar', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (66, 'Bakery & Snack', 3, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (67, 'Direct ECF', 14, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (68, 'EDN Mobile', 14, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (69, 'ED Mobile', 15, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (70, 'Pasar Tier 2', 15, '', 2, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (71, 'Pasar Tier 3', 15, '', 2, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (72, 'Catering (Corporate)', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (73, 'Catering (Personal)', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (74, 'Hotel 4 Star up', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (75, 'Hotel 3 star below', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (76, 'International Chain', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (77, 'Medium Low Chain', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (78, 'Medium High Chain', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (79, 'Indonesian/Chinese Restaurant', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (80, 'Western Restaurant', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (81, 'Asian Restaurant', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (82, 'Coffee Shop', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (83, 'Bakery', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (84, 'Medium High MT', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (85, 'Medium Low MT', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (86, 'Minimarket', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (87, 'EDN Jatiuwung', 15, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (88, 'ECF-Partnership', 14, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (89, 'Retort', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (90, 'EDN Cibitung', 15, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (91, 'EDN Cikopo', 15, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (92, 'EDN TU Bogor', 15, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (93, 'EDN Johar', 15, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (94, 'EDN Pabean', 15, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (95, 'EDN Krian', 15, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (96, 'Coffee Shop - F&C', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (97, 'West Restaurant-F&C', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (98, 'Asian Restaurant-F&C', 2, '', 1, '2022-12-20 00:50:32', NULL);
INSERT INTO `archetype` VALUES (99, 'Indo/Chi resto F&C', 2, '', 1, '2022-12-20 00:50:32', NULL);

-- ----------------------------
-- Table structure for business_type
-- ----------------------------
DROP TABLE IF EXISTS `business_type`;
CREATE TABLE `business_type`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'business type id from GP',
  `description` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'description from GP',
  `group_type` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'group type from GP',
  `abbreviation` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `status` tinyint(1) NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of business_type
-- ----------------------------
INSERT INTO `business_type` VALUES (1, 'All', 'All Business Type', '1', 'ALL', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (2, 'Horeca', 'Horeca', '1', 'MH', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (3, 'Traditional Culinary', 'Traditional Culinary', '1', 'LM', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (4, 'Partnership', 'Partnership', '1', 'EN', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (5, 'Retail', 'Retail', '1', 'RE', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (6, 'Reseller', 'Reseller', '1', 'RS', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (7, 'None', 'None', '1', 'NO', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (8, 'Wet Market', 'Wet Market', '1', 'ED', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (9, 'Personal', 'Personal', '1', 'PE', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (10, 'Modern Enterprise', 'Modern Enterprise', '1', 'ME', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (11, 'Traditional Enterprise', 'Traditional Enterprise', '1', 'TE', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (12, 'Eden Waste Center', 'Eden Waste Center', '1', 'EWC', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (13, 'Traditional Distributor', 'Traditional Distributor', '1', 'TD', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (14, 'Direct ECF', 'Direct ECF', '1', 'ECF', 1, '2022-12-20 00:52:21', NULL);
INSERT INTO `business_type` VALUES (15, 'EDN', 'EDN', '1', 'EDM', 1, '2022-12-20 00:52:21', NULL);

-- ----------------------------
-- Table structure for class
-- ----------------------------
DROP TABLE IF EXISTS `class`;
CREATE TABLE `class`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `description` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `status` tinyint(1) NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of class
-- ----------------------------
INSERT INTO `class` VALUES (1, 'CLS001', 'Class 1', 1, '2022-12-23 01:46:04', NULL);

-- ----------------------------
-- Table structure for customer
-- ----------------------------
DROP TABLE IF EXISTS `customer`;
CREATE TABLE `customer`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'Customer ID From GP',
  `name` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `short_name` varchar(15) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'Belum fix akan digunakan atau tidak',
  `business_type_id` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `archetype_id` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `business_type_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `payment_term_id` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `id_card_number` varchar(16) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `taxpayer_number` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `referrer_code` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT 'Belum tahu nama field di GPnya apa',
  `referral_code` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT 'Belum tahu nama field di GPnya apa',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = 'customer' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of customer
-- ----------------------------

-- ----------------------------
-- Table structure for item
-- ----------------------------
DROP TABLE IF EXISTS `item`;
CREATE TABLE `item`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'item ID from GP',
  `uom_id` bigint NULL DEFAULT NULL COMMENT 'uom ID from GP',
  `class_id` bigint NULL DEFAULT NULL COMMENT 'class ID from GP',
  `item_category_id` bigint NULL DEFAULT NULL COMMENT 'product tags from GP',
  `description` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'item name from GP',
  `unit_weight_conversion` decimal(5, 2) NULL DEFAULT 0.00 COMMENT 'conversion from GP',
  `order_min_qty` decimal(5, 2) NULL DEFAULT 0.00,
  `order_max_qty` decimal(5, 2) NULL DEFAULT 0.00,
  `item_type` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'item type from GP',
  `packability` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'packability flag from GP',
  `capitalize` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'capital flag from GP',
  `note` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'item description in catalog',
  `exclude_archetype` varchar(250) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL COMMENT 'array of archetype from GP',
  `max_day_delivery_date` tinyint(1) NULL DEFAULT 0 COMMENT 'maximum delivery date option from GP',
  `fragile_goods` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'fragile flag from GP',
  `taxable` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'flagging for the product is taxable or not',
  `order_channel_restriction` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `status` tinyint(1) NULL DEFAULT 0 COMMENT 'status from GP',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 45 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = 'product' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of item
-- ----------------------------
INSERT INTO `item` VALUES (1, 'GP0001', 1, 1, 1, 'Product GP0001', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (2, 'GP0002', 1, 1, 1, 'Product GP0002', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (3, 'GP0003', 1, 1, 1, 'Product GP0003', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (4, 'GP0004', 1, 1, 1, 'Product GP0004', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (5, 'GP0005', 1, 1, 1, 'Product GP0005', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (6, 'GP0006', 1, 1, 1, 'Product GP0006', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (7, 'GP0007', 1, 1, 1, 'Product GP0007', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (8, 'GP0008', 1, 1, 1, 'Product GP0008', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (9, 'GP0009', 1, 1, 1, 'Product GP0009', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (10, 'GP0010', 1, 1, 1, 'Product GP0010', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (11, 'GP0011', 1, 1, 1, 'Product GP0011', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (12, 'GP0012', 1, 1, 2, 'Product GP0012', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (13, 'GP0013', 1, 1, 2, 'Product GP0013', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (14, 'GP0014', 1, 1, 2, 'Product GP0014', 1.00, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (15, 'GP0015', 1, 1, 2, 'Product GP0015', 1.00, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (16, 'GP0016', 1, 1, 2, 'Product GP0016', 1.00, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (17, 'GP0017', 1, 1, 2, 'Product GP0017', 1.00, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (18, 'GP0018', 1, 1, 2, 'Product GP0018', 1.00, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (19, 'GP0019', 1, 1, 2, 'Product GP0019', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (20, 'GP0020', 1, 1, 2, 'Product GP0020', 1.00, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (21, 'GP0021', 1, 1, 2, 'Product GP0021', 1.00, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (22, 'GP0022', 1, 1, 1, 'Product GP0022', 1.00, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (23, 'GP0023', 1, 1, 1, 'Product GP0023', 1.00, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (24, 'GP0024', 1, 1, 1, 'Product GP0024', 1.00, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (25, 'GP0025', 1, 1, 1, 'Product GP0025', 1.00, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (26, 'GP0026', 1, 1, 1, 'Product GP0026', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (27, 'GP0027', 1, 1, 1, 'Product GP0027', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (28, 'GP0028', 1, 1, 1, 'Product GP0028', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (29, 'GP0029', 1, 1, 1, 'Product GP0029', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (30, 'GP0030', 1, 1, 3, 'Product GP0030', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (31, 'GP0031', 1, 1, 3, 'Product GP0031', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (32, 'GP0032', 1, 1, 3, 'Product GP0032', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (33, 'GP0033', 1, 1, 3, 'Product GP0033', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (34, 'GP0034', 1, 1, 3, 'Product GP0034', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (35, 'GP0035', 1, 1, 3, 'Product GP0035', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (36, 'GP0036', 1, 1, 3, 'Product GP0036', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (37, 'GP0037', 1, 1, 3, 'Product GP0037', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (38, 'GP0038', 1, 1, 1, 'Product GP0038', 0.50, 1.00, 0.00, 'Sales Inventory', 'non packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (39, 'GP0039', 1, 1, 1, 'Product GP0039', 0.50, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (40, 'GP0040', 1, 1, 1, 'Product GP0040', 0.50, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (41, 'GP0041', 1, 1, 1, 'Product GP0041', 0.50, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (42, 'GP0042', 1, 1, 1, 'Product GP0042', 0.50, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (43, 'GP0043', 1, 1, 1, 'Product GP0043', 0.50, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', NULL);
INSERT INTO `item` VALUES (44, 'GP0044', 2, 1, 1, 'SKU Pack', 0.50, 1.00, 0.00, 'Sales Inventory', 'packable', 'non capital', 'It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.', NULL, 3, 'non fragile', 'taxable', 'dashboard', 1, '2022-12-20 00:57:46', '2023-02-14 06:41:18');

-- ----------------------------
-- Table structure for item_category
-- ----------------------------
DROP TABLE IF EXISTS `item_category`;
CREATE TABLE `item_category`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `region_id` bigint NULL DEFAULT NULL COMMENT 'region name from GP',
  `name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `status` tinyint(1) NULL DEFAULT 0 COMMENT 'status from GP',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = 'product_tag' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of item_category
-- ----------------------------
INSERT INTO `item_category` VALUES (1, NULL, 1, 'Sayuran Segar', 1, '2022-12-20 00:57:46', '2022-12-20 01:03:41');
INSERT INTO `item_category` VALUES (2, NULL, 1, 'Bumbu Dapur', 1, '2022-12-20 00:57:46', '2022-12-20 01:03:41');
INSERT INTO `item_category` VALUES (3, NULL, 1, 'Cabai & Bawang', 1, '2022-12-20 00:57:46', '2022-12-20 01:03:41');
INSERT INTO `item_category` VALUES (4, NULL, 1, 'Telur, Daging & Ikan', 1, '2022-12-20 00:57:46', '2022-12-20 01:03:41');
INSERT INTO `item_category` VALUES (5, NULL, 1, 'Buah Segar', 1, '2022-12-20 00:57:46', '2022-12-20 01:03:41');
INSERT INTO `item_category` VALUES (6, NULL, 1, 'Minuman Sehat', 1, '2022-12-20 00:57:46', '2022-12-20 01:03:41');
INSERT INTO `item_category` VALUES (7, NULL, 1, 'Bahan Pangan Olahan', 1, '2022-12-20 00:57:46', '2022-12-20 01:03:41');
INSERT INTO `item_category` VALUES (8, NULL, 1, 'Sayuran Hidro', 1, '2022-12-20 00:57:46', '2022-12-20 01:03:41');

-- ----------------------------
-- Table structure for region
-- ----------------------------
DROP TABLE IF EXISTS `region`;
CREATE TABLE `region`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'region_id from GP',
  `description` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '' COMMENT 'description from GP',
  `status` tinyint(1) NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of region
-- ----------------------------
INSERT INTO `region` VALUES (1, 'all', 'Dummy All Area', 1, '2022-12-20 00:57:28', '2023-02-14 07:26:14');
INSERT INTO `region` VALUES (2, 'jakarta', 'Dummy Jakarta', 1, '2022-12-20 00:57:28', '2023-02-14 07:25:49');
INSERT INTO `region` VALUES (3, 'bandung', 'Dummy Bandung', 1, '2022-12-20 00:57:28', '2023-02-14 07:25:49');
INSERT INTO `region` VALUES (4, 'semarang', 'Dummy Semarang', 1, '2022-12-20 00:57:28', '2023-02-14 07:25:49');
INSERT INTO `region` VALUES (5, 'surabaya', 'Dummy Surabaya', 1, '2022-12-20 00:57:28', '2023-02-14 07:25:50');
INSERT INTO `region` VALUES (6, 'medan', 'Dummy Medan', 1, '2022-12-20 00:57:28', '2023-02-14 07:25:50');

-- ----------------------------
-- Table structure for sales_order
-- ----------------------------
DROP TABLE IF EXISTS `sales_order`;
CREATE TABLE `sales_order`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `doc_number` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `address_id` bigint NULL DEFAULT NULL,
  `customer_id` bigint UNSIGNED NULL DEFAULT NULL,
  `salesperson_id` bigint NULL DEFAULT NULL,
  `application` tinyint(1) NULL DEFAULT 0,
  `status` tinyint NULL DEFAULT 0,
  `order_date` date NULL DEFAULT NULL,
  `total` decimal(20, 2) NULL DEFAULT 0.00,
  `created_date` date NULL DEFAULT NULL,
  `modified_date` date NULL DEFAULT NULL,
  `finished_date` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 227 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sales_order
-- ----------------------------
INSERT INTO `sales_order` VALUES (1, 'SO-000001', 'DOC00001', 1, 1, 1, 1, 1, '2023-02-10', 350000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (2, 'SO-000002', 'DOC00002', 3, 3, 1, 1, 1, '2023-02-10', 12000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (3, 'SO-000003', 'DOC00003', 2, 2, 2, 1, 1, '2023-02-10', 8000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (4, 'SO-000004', 'DOC00004', 4, 4, 1, 1, 1, '2023-02-10', 16000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (5, 'SO-000005', 'DOC00005', 3, 3, 2, 1, 1, '2023-02-10', 12000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (6, 'SO-000006', 'DOC00006', 1, 1, 1, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (7, 'SO-000007', 'DOC00007', 1, 1, 2, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (8, 'SO-000008', 'DOC00008', 3, 3, 1, 1, 1, '2023-02-10', 12000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (9, 'SO-000009', 'DOC00009', 1, 1, 1, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (10, 'SO-000010', 'DOC00010', 2, 2, 1, 1, 1, '2023-02-10', 8000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (11, 'SO-000011', 'DOC00011', 1, 1, 2, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (12, 'SO-000012', 'DOC00012', 2, 2, 1, 1, 1, '2023-02-10', 8000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (13, 'SO-000013', 'DOC00013', 2, 2, 1, 1, 1, '2023-02-10', 8000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (14, 'SO-000014', 'DOC00014', 1, 1, 1, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (15, 'SO-000015', 'DOC00015', 3, 3, 1, 1, 1, '2023-02-10', 12000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (16, 'SO-000016', 'DOC00016', 2, 2, 2, 1, 1, '2023-02-10', 8000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (17, 'SO-000017', 'DOC00017', 4, 4, 1, 1, 1, '2023-02-10', 16000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (18, 'SO-000018', 'DOC00018', 3, 3, 1, 1, 1, '2023-02-10', 12000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (19, 'SO-000019', 'DOC00019', 1, 1, 2, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (20, 'SO-000020', 'DOC00020', 1, 1, 1, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (21, 'SO-000021', 'DOC00021', 3, 3, 2, 1, 1, '2023-02-10', 12000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (22, 'SO-000022', 'DOC00022', 1, 1, 1, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (23, 'SO-000023', 'DOC00023', 1, 1, 2, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (24, 'SO-000024', 'DOC00024', 3, 3, 1, 1, 1, '2023-02-10', 12000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (25, 'SO-000025', 'DOC00025', 2, 2, 1, 1, 1, '2023-02-10', 8000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (26, 'SO-000026', 'DOC00026', 4, 4, 2, 1, 1, '2023-02-10', 16000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (27, 'SO-000027', 'DOC00027', 3, 3, 1, 1, 1, '2023-02-10', 12000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (28, 'SO-000028', 'DOC00028', 1, 1, 2, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (29, 'SO-000029', 'DOC00029', 1, 1, 1, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (30, 'SO-000030', 'DOC00030', 3, 3, 2, 1, 1, '2023-02-10', 12000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (31, 'SO-000031', 'DOC00031', 1, 1, 1, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (32, 'SO-000032', 'DOC00032', 2, 2, 1, 1, 1, '2023-02-10', 8000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (33, 'SO-000033', 'DOC00033', 1, 1, 1, 1, 1, '2023-02-10', 4000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (34, 'SO-000034', 'DOC00034', 2, 2, 2, 1, 1, '2023-02-10', 8000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (35, 'SO-000035', 'DOC00035', 4, 1, 3, 1, 1, '2023-02-10', 220000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (36, 'SO-000036', 'DOC00036', 10, 1, 1, 1, 1, '2023-02-10', 130000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (37, 'SO-000037', 'DOC00037', 5, 1, 2, 1, 1, '2023-02-10', 380000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (38, 'SO-000038', 'DOC00038', 3, 1, 4, 1, 1, '2023-02-10', 30000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (39, 'SO-000039', 'DOC00039', 1, 1, 1, 1, 1, '2023-02-10', 730000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (40, 'SO-000040', 'DOC00040', 3, 1, 5, 1, 1, '2023-02-10', 620000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (41, 'SO-000041', 'DOC00041', 5, 1, 2, 1, 1, '2023-02-10', 1600000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (42, 'SO-000042', 'DOC00042', 9, 1, 4, 1, 1, '2023-02-10', 360000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (43, 'SO-000043', 'DOC00043', 7, 1, 2, 1, 1, '2023-02-10', 880000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (44, 'SO-000044', 'DOC00044', 7, 1, 5, 1, 1, '2023-02-10', 980000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (45, 'SO-000045', 'DOC00045', 9, 1, 2, 1, 1, '2023-02-10', 1240000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (46, 'SO-000046', 'DOC00046', 5, 1, 4, 1, 1, '2023-02-10', 1120000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (47, 'SO-000047', 'DOC00047', 5, 1, 3, 1, 1, '2023-02-10', 2280000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (48, 'SO-000048', 'DOC00048', 8, 1, 1, 1, 1, '2023-02-10', 560000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (49, 'SO-000049', 'DOC00049', 8, 1, 3, 1, 1, '2023-02-10', 2430000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (50, 'SO-000050', 'DOC00050', 10, 1, 4, 1, 1, '2023-02-10', 2670000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (51, 'SO-000051', 'DOC00051', 7, 1, 2, 1, 1, '2023-02-10', 4350000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (52, 'SO-000052', 'DOC00052', 1, 1, 3, 1, 1, '2023-02-10', 2320000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (53, 'SO-000053', 'DOC00053', 8, 1, 4, 1, 1, '2023-02-10', 1120000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (54, 'SO-000054', 'DOC00054', 5, 1, 2, 1, 1, '2023-02-10', 3300000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (55, 'SO-000055', 'DOC00055', 7, 1, 5, 1, 1, '2023-02-10', 1440000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (56, 'SO-000056', 'DOC00056', 9, 1, 5, 1, 1, '2023-02-10', 2400000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (57, 'SO-000057', 'DOC00057', 9, 1, 5, 1, 1, '2023-02-10', 2850000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (58, 'SO-000058', 'DOC00058', 8, 1, 5, 1, 1, '2023-02-10', 4850000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (59, 'SO-000059', 'DOC00059', 10, 1, 3, 1, 1, '2023-02-10', 1380000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (60, 'SO-000060', 'DOC00060', 9, 1, 4, 1, 1, '2023-02-10', 3750000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (61, 'SO-000061', 'DOC00061', 8, 1, 4, 1, 1, '2023-02-10', 510000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (62, 'SO-000062', 'DOC00062', 1, 1, 2, 1, 1, '2023-02-10', 1360000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (63, 'SO-000063', 'DOC00063', 10, 1, 3, 1, 1, '2023-02-10', 2220000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (64, 'SO-000064', 'DOC00064', 7, 1, 5, 1, 1, '2023-02-10', 840000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (65, 'SO-000065', 'DOC00065', 10, 1, 3, 1, 1, '2023-02-10', 110000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (66, 'SO-000066', 'DOC00066', 3, 1, 4, 1, 1, '2023-02-10', 690000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (67, 'SO-000067', 'DOC00067', 3, 1, 1, 1, 1, '2023-02-10', 800000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (68, 'SO-000068', 'DOC00068', 9, 1, 5, 1, 1, '2023-02-10', 470000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (69, 'SO-000069', 'DOC00069', 2, 1, 3, 1, 1, '2023-02-10', 2580000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (70, 'SO-000070', 'DOC00070', 10, 1, 3, 1, 1, '2023-02-10', 1700000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (71, 'SO-000071', 'DOC00071', 4, 1, 4, 1, 1, '2023-02-10', 1040000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (72, 'SO-000072', 'DOC00072', 6, 1, 4, 1, 1, '2023-02-10', 820000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (73, 'SO-000073', 'DOC00073', 7, 1, 2, 1, 1, '2023-02-10', 4900000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (74, 'SO-000074', 'DOC00074', 6, 1, 2, 1, 1, '2023-02-10', 1820000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (75, 'SO-000075', 'DOC00075', 2, 1, 5, 1, 1, '2023-02-10', 10000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (76, 'SO-000076', 'DOC00076', 4, 1, 2, 1, 1, '2023-02-10', 5000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (77, 'SO-000077', 'DOC00077', 2, 1, 4, 1, 1, '2023-02-10', 500000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (78, 'SO-000078', 'DOC00078', 6, 1, 5, 1, 1, '2023-02-10', 930000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (79, 'SO-000079', 'DOC00079', 7, 1, 5, 1, 1, '2023-02-10', 3500000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (80, 'SO-000080', 'DOC00080', 6, 1, 5, 1, 1, '2023-02-10', 80000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (81, 'SO-000081', 'DOC00081', 2, 1, 5, 1, 1, '2023-02-10', 3850000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (82, 'SO-000082', 'DOC00082', 2, 1, 1, 1, 1, '2023-02-10', 50000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (83, 'SO-000083', 'DOC00083', 9, 1, 5, 1, 1, '2023-02-10', 920000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (84, 'SO-000084', 'DOC00084', 1, 1, 3, 1, 1, '2023-02-10', 800000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (85, 'SO-000085', 'DOC00085', 3, 1, 5, 1, 1, '2023-02-10', 1940000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (86, 'SO-000086', 'DOC00086', 1, 1, 3, 1, 1, '2023-02-10', 730000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (87, 'SO-000087', 'DOC00087', 6, 1, 2, 1, 1, '2023-02-10', 1080000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (88, 'SO-000088', 'DOC00088', 7, 1, 1, 1, 1, '2023-02-10', 100000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (89, 'SO-000089', 'DOC00089', 6, 1, 2, 1, 1, '2023-02-10', 2010000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (90, 'SO-000090', 'DOC00090', 2, 1, 1, 1, 1, '2023-02-10', 840000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (91, 'SO-000091', 'DOC00091', 10, 1, 5, 1, 1, '2023-02-10', 2000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (92, 'SO-000092', 'DOC00092', 3, 1, 3, 1, 1, '2023-02-10', 3880000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (93, 'SO-000093', 'DOC00093', 3, 1, 2, 1, 1, '2023-02-10', 2430000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (94, 'SO-000094', 'DOC00094', 1, 1, 4, 1, 1, '2023-02-10', 1020000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (95, 'SO-000095', 'DOC00095', 1, 1, 1, 1, 1, '2023-02-10', 680000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (96, 'SO-000096', 'DOC00096', 7, 1, 4, 1, 1, '2023-02-10', 670000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (97, 'SO-000097', 'DOC00097', 4, 1, 3, 1, 1, '2023-02-10', 3300000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (98, 'SO-000098', 'DOC00098', 4, 1, 5, 1, 1, '2023-02-10', 720000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (99, 'SO-000099', 'DOC00099', 6, 1, 2, 1, 1, '2023-02-10', 2760000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (100, 'SO-000100', 'DOC00100', 1, 1, 3, 1, 1, '2023-02-10', 730000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (101, 'SO-000101', 'DOC00101', 4, 1, 4, 1, 1, '2023-02-10', 1200000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (102, 'SO-000102', 'DOC00102', 8, 1, 4, 1, 1, '2023-02-10', 2730000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (103, 'SO-000103', 'DOC00103', 9, 1, 5, 1, 1, '2023-02-10', 4450000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (104, 'SO-000104', 'DOC00104', 4, 1, 4, 1, 1, '2023-02-10', 220000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (105, 'SO-000105', 'DOC00105', 3, 1, 1, 1, 1, '2023-02-10', 3520000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (106, 'SO-000106', 'DOC00106', 6, 1, 1, 1, 1, '2023-02-10', 960000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (107, 'SO-000107', 'DOC00107', 2, 1, 2, 1, 1, '2023-02-10', 1950000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (108, 'SO-000108', 'DOC00108', 3, 1, 5, 1, 1, '2023-02-10', 720000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (109, 'SO-000109', 'DOC00109', 8, 1, 2, 1, 1, '2023-02-10', 3250000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (110, 'SO-000110', 'DOC00110', 7, 1, 5, 1, 1, '2023-02-10', 850000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (111, 'SO-000111', 'DOC00111', 8, 1, 2, 1, 1, '2023-02-10', 1120000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (112, 'SO-000112', 'DOC00112', 1, 1, 2, 1, 1, '2023-02-10', 640000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (113, 'SO-000113', 'DOC00113', 2, 1, 3, 1, 1, '2023-02-10', 520000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (114, 'SO-000114', 'DOC00114', 2, 1, 2, 1, 1, '2023-02-10', 4250000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (115, 'SO-000115', 'DOC00115', 1, 1, 2, 1, 1, '2023-02-10', 640000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (116, 'SO-000116', 'DOC00116', 10, 1, 2, 1, 1, '2023-02-10', 1750000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (117, 'SO-000117', 'DOC00117', 7, 1, 4, 1, 1, '2023-02-10', 1200000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (118, 'SO-000118', 'DOC00118', 8, 1, 3, 1, 1, '2023-02-10', 1230000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (119, 'SO-000119', 'DOC00119', 7, 1, 2, 1, 1, '2023-02-10', 550000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (120, 'SO-000120', 'DOC00120', 3, 1, 4, 1, 1, '2023-02-10', 1080000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (121, 'SO-000121', 'DOC00121', 6, 1, 2, 1, 1, '2023-02-10', 1000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (122, 'SO-000122', 'DOC00122', 7, 1, 5, 1, 1, '2023-02-10', 2200000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (123, 'SO-000123', 'DOC00123', 5, 1, 4, 1, 1, '2023-02-10', 650000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (124, 'SO-000124', 'DOC00124', 8, 1, 1, 1, 1, '2023-02-10', 340000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (125, 'SO-000125', 'DOC00125', 2, 1, 5, 1, 1, '2023-02-10', 720000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (126, 'SO-000126', 'DOC00126', 6, 1, 3, 1, 1, '2023-02-10', 390000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (127, 'SO-000127', 'DOC00127', 2, 1, 3, 1, 1, '2023-02-10', 1200000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (128, 'SO-000128', 'DOC00128', 5, 1, 2, 1, 1, '2023-02-10', 4800000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (129, 'SO-000129', 'DOC00129', 4, 1, 5, 1, 1, '2023-02-10', 1240000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (130, 'SO-000130', 'DOC00130', 9, 1, 2, 1, 1, '2023-02-10', 410000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (131, 'SO-000131', 'DOC00131', 9, 1, 1, 1, 1, '2023-02-10', 640000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (132, 'SO-000132', 'DOC00132', 2, 1, 4, 1, 1, '2023-02-10', 1700000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (133, 'SO-000133', 'DOC00133', 1, 1, 5, 1, 1, '2023-02-10', 4150000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (134, 'SO-000134', 'DOC00134', 1, 1, 4, 1, 1, '2023-02-10', 870000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (135, 'SO-000135', 'DOC00135', 3, 1, 5, 1, 1, '2023-02-10', 1280000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (136, 'SO-000136', 'DOC00136', 7, 1, 3, 1, 1, '2023-02-10', 3950000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (137, 'SO-000137', 'DOC00137', 1, 1, 4, 1, 1, '2023-02-10', 2000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (138, 'SO-000138', 'DOC00138', 1, 1, 4, 1, 1, '2023-02-10', 280000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (139, 'SO-000139', 'DOC00139', 10, 1, 2, 1, 1, '2023-02-10', 3800000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (140, 'SO-000140', 'DOC00140', 9, 1, 5, 1, 1, '2023-02-10', 2800000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (141, 'SO-000141', 'DOC00141', 10, 1, 5, 1, 1, '2023-02-10', 1060000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (142, 'SO-000142', 'DOC00142', 10, 1, 5, 1, 1, '2023-02-10', 4850000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (143, 'SO-000143', 'DOC00143', 6, 1, 1, 1, 1, '2023-02-10', 1710000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (144, 'SO-000144', 'DOC00144', 3, 1, 5, 1, 1, '2023-02-10', 240000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (145, 'SO-000145', 'DOC00145', 9, 1, 2, 1, 1, '2023-02-10', 1080000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (146, 'SO-000146', 'DOC00146', 5, 1, 2, 1, 1, '2023-02-10', 440000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (147, 'SO-000147', 'DOC00147', 6, 1, 1, 1, 1, '2023-02-10', 380000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (148, 'SO-000148', 'DOC00148', 5, 1, 5, 1, 1, '2023-02-10', 480000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (149, 'SO-000149', 'DOC00149', 7, 1, 1, 1, 1, '2023-02-10', 1580000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (150, 'SO-000150', 'DOC00150', 10, 1, 5, 1, 1, '2023-02-10', 2650000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (151, 'SO-000151', 'DOC00151', 1, 1, 2, 1, 1, '2023-02-10', 900000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (152, 'SO-000152', 'DOC00152', 2, 1, 3, 1, 1, '2023-02-10', 520000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (153, 'SO-000153', 'DOC00153', 6, 1, 5, 1, 1, '2023-02-10', 2760000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (154, 'SO-000154', 'DOC00154', 4, 1, 5, 1, 1, '2023-02-10', 480000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (155, 'SO-000155', 'DOC00155', 8, 1, 4, 1, 1, '2023-02-10', 10000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (156, 'SO-000156', 'DOC00156', 10, 1, 4, 1, 1, '2023-02-10', 1590000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (157, 'SO-000157', 'DOC00157', 10, 1, 4, 1, 1, '2023-02-10', 3120000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (158, 'SO-000158', 'DOC00158', 7, 1, 5, 1, 1, '2023-02-10', 1900000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (159, 'SO-000159', 'DOC00159', 5, 1, 5, 1, 1, '2023-02-10', 1160000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (160, 'SO-000160', 'DOC00160', 9, 1, 4, 1, 1, '2023-02-10', 2910000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (161, 'SO-000161', 'DOC00161', 4, 1, 2, 1, 1, '2023-02-10', 50000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (162, 'SO-000162', 'DOC00162', 6, 1, 2, 1, 1, '2023-02-10', 330000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (163, 'SO-000163', 'DOC00163', 9, 1, 5, 1, 1, '2023-02-10', 1860000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (164, 'SO-000164', 'DOC00164', 3, 1, 5, 1, 1, '2023-02-10', 640000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (165, 'SO-000165', 'DOC00165', 10, 1, 3, 1, 1, '2023-02-10', 1400000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (166, 'SO-000166', 'DOC00166', 7, 1, 5, 1, 1, '2023-02-10', 610000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (167, 'SO-000167', 'DOC00167', 9, 1, 2, 1, 1, '2023-02-10', 480000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (168, 'SO-000168', 'DOC00168', 6, 1, 5, 1, 1, '2023-02-10', 1260000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (169, 'SO-000169', 'DOC00169', 5, 1, 2, 1, 1, '2023-02-10', 470000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (170, 'SO-000170', 'DOC00170', 4, 1, 1, 1, 1, '2023-02-10', 830000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (171, 'SO-000171', 'DOC00171', 8, 1, 2, 1, 1, '2023-02-10', 840000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (172, 'SO-000172', 'DOC00172', 1, 1, 3, 1, 1, '2023-02-10', 2580000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (173, 'SO-000173', 'DOC00173', 4, 1, 2, 1, 1, '2023-02-10', 450000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (174, 'SO-000174', 'DOC00174', 2, 1, 2, 1, 1, '2023-02-10', 540000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (175, 'SO-000175', 'DOC00175', 2, 1, 2, 1, 1, '2023-02-10', 250000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (176, 'SO-000176', 'DOC00176', 5, 1, 3, 1, 1, '2023-02-10', 540000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (177, 'SO-000177', 'DOC00177', 2, 1, 1, 1, 1, '2023-02-10', 300000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (178, 'SO-000178', 'DOC00178', 2, 1, 3, 1, 1, '2023-02-10', 2050000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (179, 'SO-000179', 'DOC00179', 9, 1, 4, 1, 1, '2023-02-10', 3120000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (180, 'SO-000180', 'DOC00180', 3, 1, 3, 1, 1, '2023-02-10', 140000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (181, 'SO-000181', 'DOC00181', 8, 1, 1, 1, 1, '2023-02-10', 4100000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (182, 'SO-000182', 'DOC00182', 2, 1, 5, 1, 1, '2023-02-10', 3000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (183, 'SO-000183', 'DOC00183', 6, 1, 2, 1, 1, '2023-02-10', 2910000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (184, 'SO-000184', 'DOC00184', 6, 1, 5, 1, 1, '2023-02-10', 1940000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (185, 'SO-000185', 'DOC00185', 6, 1, 2, 1, 1, '2023-02-10', 120000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (186, 'SO-000186', 'DOC00186', 10, 1, 4, 1, 1, '2023-02-10', 1860000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (187, 'SO-000187', 'DOC00187', 6, 1, 5, 1, 1, '2023-02-10', 2000000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (188, 'SO-000188', 'DOC00188', 4, 1, 1, 1, 1, '2023-02-10', 720000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (189, 'SO-000189', 'DOC00189', 10, 1, 2, 1, 1, '2023-02-10', 2430000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (190, 'SO-000190', 'DOC00190', 8, 1, 1, 1, 1, '2023-02-10', 450000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (191, 'SO-000191', 'DOC00191', 8, 1, 3, 1, 1, '2023-02-10', 1960000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (192, 'SO-000192', 'DOC00192', 9, 1, 3, 1, 1, '2023-02-10', 620000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (193, 'SO-000193', 'DOC00193', 3, 1, 2, 1, 1, '2023-02-10', 3960000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (194, 'SO-000194', 'DOC00194', 1, 1, 5, 1, 1, '2023-02-10', 1040000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (195, 'SO-000195', 'DOC00195', 6, 1, 1, 1, 1, '2023-02-10', 220000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (196, 'SO-000196', 'DOC00196', 8, 1, 5, 1, 1, '2023-02-10', 690000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (197, 'SO-000197', 'DOC00197', 3, 1, 2, 1, 1, '2023-02-10', 630000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (198, 'SO-000198', 'DOC00198', 77, 1, 1, 1, 1, '2023-02-10', 10000.00, '2023-01-01', NULL, NULL, NULL, NULL);
INSERT INTO `sales_order` VALUES (199, 'SO-000199', 'DOC00199', 5, 1, 1, 1, 1, '2023-02-10', 700000.00, '2023-01-01', '2023-01-25', '2023-01-25 00:00:00', '2023-01-25 00:00:00', '2023-01-25 00:00:00');
INSERT INTO `sales_order` VALUES (200, 'SO-000200', 'DOC00200', 4, 1, 1, 1, 1, '2023-02-10', 100000.00, '2023-01-01', '2023-01-25', '2023-01-25 00:00:00', '2023-01-25 00:00:00', '2023-01-25 00:00:00');
INSERT INTO `sales_order` VALUES (225, 'SO-000201', 'DOC00201', 2, 1, 1, 1, 1, '2023-02-10', 200000.00, '2023-01-01', '2023-01-25', '2023-01-25 00:00:00', '2023-01-25 00:00:00', '2023-01-25 00:00:00');
INSERT INTO `sales_order` VALUES (226, 'SO-000202', 'DOC00202', 1, 1, 1, 1, 1, '2023-02-10', 200000.00, '2023-01-01', '2023-01-26', '2023-01-26 00:00:00', '2023-01-26 00:00:00', '2023-01-26 00:00:00');

-- ----------------------------
-- Table structure for salesperson
-- ----------------------------
DROP TABLE IF EXISTS `salesperson`;
CREATE TABLE `salesperson`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `firstname` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `middlename` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `lastname` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `status` tinyint(1) NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 182 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of salesperson
-- ----------------------------
INSERT INTO `salesperson` VALUES (1, 'EPI000000001', 'Andi', 'Robby', 'Yono', 1, '2022-12-20 00:59:04', '2023-01-20 23:25:29');
INSERT INTO `salesperson` VALUES (2, 'EPI000000002', 'Adan', 'Firman', 'Bono', 1, '2022-12-20 00:59:04', '2023-01-20 23:25:29');
INSERT INTO `salesperson` VALUES (3, 'EPI000000003', 'Dirga', 'Tony', 'Messi', 1, '2022-12-20 00:59:04', '2023-01-20 23:25:29');
INSERT INTO `salesperson` VALUES (4, 'EPI000000004', 'Martin', 'Louis', 'Rooney', 1, '2022-12-20 00:59:04', '2023-01-20 23:25:29');
INSERT INTO `salesperson` VALUES (5, 'EPI000000005', 'Yono', 'Bintang', 'Ahmad', 1, '2022-12-20 00:59:04', '2023-01-20 23:25:29');
INSERT INTO `salesperson` VALUES (6, 'EPI000000006', 'Bono', 'Ardi', 'Charlie', 1, '2022-12-20 00:59:04', '2023-01-27 07:31:57');
INSERT INTO `salesperson` VALUES (7, 'EPI000000007', 'Messi', 'Gigih', 'Ulil', 1, '2022-12-20 00:59:04', '2023-01-27 07:31:58');
INSERT INTO `salesperson` VALUES (8, 'EPI000000008', 'Rooney', 'Ibrahim', 'Khalid', 1, '2022-12-20 00:59:04', '2023-01-27 07:31:59');
INSERT INTO `salesperson` VALUES (9, 'EPI000000009', 'Ahmad', 'Arief', 'Desta', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:22');
INSERT INTO `salesperson` VALUES (10, 'EPI000000010', 'Charlie', 'Ario', 'Vincent', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:24');
INSERT INTO `salesperson` VALUES (11, 'EPI000000011', 'Ulil', 'Adrian', 'Wahid', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:11');
INSERT INTO `salesperson` VALUES (12, 'EPI000000012', 'Khalid', 'Tegar', 'Tri', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:35');
INSERT INTO `salesperson` VALUES (13, 'EPI000000013', 'Desta', 'Panca', 'Prabowo', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:39');
INSERT INTO `salesperson` VALUES (14, 'EPI000000014', 'Vincent', 'Imron', 'Dede', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:40');
INSERT INTO `salesperson` VALUES (15, 'EPI000000015', 'Wahid', 'Qomar', 'Kamil', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:41');
INSERT INTO `salesperson` VALUES (16, 'EPI000000016', 'Tri', 'Malik', 'Husni', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:42');
INSERT INTO `salesperson` VALUES (17, 'EPI000000017', 'Prabowo', 'Dennis', 'Junior', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:44');
INSERT INTO `salesperson` VALUES (18, 'EPI000000018', 'Dede', 'Johan', 'Joy', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:45');
INSERT INTO `salesperson` VALUES (19, 'EPI000000019', 'Kamil', 'Tyo', 'Faisal', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:46');
INSERT INTO `salesperson` VALUES (20, 'EPI000000020', 'Husni', 'Andra', 'Nino', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:48');
INSERT INTO `salesperson` VALUES (21, 'EPI000000021', 'Junior', 'Kevin', 'Paul', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:49');
INSERT INTO `salesperson` VALUES (22, 'EPI000000022', 'Joy', 'Rizky', 'Ibrahim', 1, '2022-12-20 00:59:04', '2023-01-27 07:32:50');
INSERT INTO `salesperson` VALUES (181, 'EPI00000002', 'Sales', 'BLABLA', 'User', 1, '2023-02-21 09:31:39', '2023-02-22 10:13:14');

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
INSERT INTO `schema_migrations` VALUES (4, 0);

-- ----------------------------
-- Table structure for site
-- ----------------------------
DROP TABLE IF EXISTS `site`;
CREATE TABLE `site`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `description` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `status` tinyint(1) NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci COMMENT = 'warehouse' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of site
-- ----------------------------
INSERT INTO `site` VALUES (1, 'ST0001', 'Site 1', 1, '2023-02-06 09:57:50', NULL);
INSERT INTO `site` VALUES (2, 'ST0002', 'Site 2', 1, '2023-02-13 11:28:26', '2023-02-13 11:28:57');

-- ----------------------------
-- Table structure for sub_district
-- ----------------------------
DROP TABLE IF EXISTS `sub_district`;
CREATE TABLE `sub_district`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `description` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `status` tinyint(1) NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 45 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sub_district
-- ----------------------------
INSERT INTO `sub_district` VALUES (1, 'Kl.BndngnHlr', 'Kel. Bendungan Hilir', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (2, 'Kl.Glr', 'Kel. Gelora', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (3, 'Kl.KmpngBl', 'Kel. Kampung Bali', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (4, 'Kl.KrtTngsn', 'Kel. Karet Tengsin', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (5, 'Kl.KbnKcng', 'Kel. Kebon Kacang', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (6, 'Kl.KbnMlt', 'Kel. Kebon Melati', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (7, 'Kl.Ptmbrn', 'Kel. Petamburan', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (8, 'Kl.Bngr', 'Kel. Bungur', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (9, 'Kl.Knr', 'Kel. Kenari', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (10, 'Kl.Krmt', 'Kel. Kramat', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (11, 'Kl.Kwtng', 'Kel. Kwitang', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (12, 'Kl.Psbn', 'Kel. Paseban', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (13, 'Kl.Snn', 'Kel. Senen', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (14, 'Kl.GnngShrtr', 'Kel. Gunung Sahari Utara', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (15, 'Kl.Krngnyr', 'Kel. Karang Anyar', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (16, 'Kl.Krtn', 'Kel. Kartini', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (17, 'Kl.MnggDSltn', 'Kel. Mangga Dua Selatan', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:40');
INSERT INTO `sub_district` VALUES (18, 'Kl.PsrBr', 'Kel. Pasar Baru', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (19, 'Kl.Ckn', 'Kel. Cikini', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (20, 'Kl.Gndngd', 'Kel. Gondangdia', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (21, 'Kl.KbnSrh', 'Kel. Kebon Sirih', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (22, 'Kl.Mntng', 'Kel. Menteng', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (23, 'Kl.Pgngsn', 'Kel. Pegangsaan', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (24, 'Kl.CmpkBr', 'Kel. Cempaka Baru', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (25, 'Kl.GnngShrSltn', 'Kel. Gunung Sahari Selatan', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (26, 'Kl.HrpnMly', 'Kel. Harapan Mulya', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (27, 'Kl.KbnKsng', 'Kel. Kebon Kosong', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (28, 'Kl.Kmyrn', 'Kel. Kemayoran', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (29, 'Kl.Srdng', 'Kel. Serdang', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (30, 'Kl.SmrBt', 'Kel. Sumur Batu', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (31, 'Kl.tnPnjng', 'Kel. Utan Panjang', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (32, 'Kl.Glr', 'Kel. Galur', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (33, 'Kl.JhrBr', 'Kel. Johar Baru', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (34, 'Kl.KmpngRw', 'Kel. Kampung Rawa', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (35, 'Kl.TnhTngg', 'Kel. Tanah Tinggi', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (36, 'Kl.Cdng', 'Kel. Cideng', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (37, 'Kl.DrPl', 'Kel. Duri Pulo', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (38, 'Kl.Gmbr', 'Kel. Gambir', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (39, 'Kl.KbnKlp', 'Kel. Kebon Kelapa', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (40, 'Kl.PtjSltn', 'Kel. Petojo Selatan', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:41');
INSERT INTO `sub_district` VALUES (41, 'Kl.Ptjtr', 'Kel. Petojo Utara', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:42');
INSERT INTO `sub_district` VALUES (42, 'Kl.CmpkPthBrt', 'Kel. Cempaka Putih Barat', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:42');
INSERT INTO `sub_district` VALUES (43, 'Kl.CmpkPthTmr', 'Kel. Cempaka Putih Timur', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:42');
INSERT INTO `sub_district` VALUES (44, 'Kl.Rwsr', 'Kel. Rawasari', 1, '2022-12-20 01:00:22', '2022-12-21 05:18:42');

-- ----------------------------
-- Table structure for territory
-- ----------------------------
DROP TABLE IF EXISTS `territory`;
CREATE TABLE `territory`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'territory id from GP',
  `description` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `region_id` bigint NULL DEFAULT NULL,
  `salesperson_id` bigint NULL DEFAULT NULL,
  `business_type_id` bigint NULL DEFAULT NULL,
  `sub_district_id` bigint NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 181 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'sales_group' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of territory
-- ----------------------------
INSERT INTO `territory` VALUES (1, 'T000001', 'Territory 1', 2, 1, 2, 1, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (2, 'T000002', 'Territory 2', 2, 2, 2, 2, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (3, 'T000003', 'Territory 3', 2, 3, 2, 3, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (4, 'T000004', 'Territory 4', 2, 4, 2, 4, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (5, 'T000005', 'Territory 5', 2, 5, 2, 5, '2022-12-20 01:02:05', '2023-01-24 10:50:16');
INSERT INTO `territory` VALUES (6, 'T000006', 'Territory 6', 2, 6, 2, 6, '2022-12-20 01:02:05', '2022-12-27 02:37:07');
INSERT INTO `territory` VALUES (7, 'T000007', 'Territory 7', 2, 7, 2, 7, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (8, 'T000008', 'Territory 8', 2, 8, 2, 8, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (9, 'T000009', 'Territory 9', 2, 9, 2, 9, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (10, 'T000010', 'Territory 10', 2, 10, 2, 10, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (11, 'T000011', 'Territory 11', 2, 11, 2, 11, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (12, 'T000012', 'Territory 12', 2, 12, 2, 12, '2022-12-20 01:02:05', '2023-01-24 10:50:16');
INSERT INTO `territory` VALUES (13, 'T000015', 'Territory 13', 2, 13, 2, 13, '2022-12-20 01:02:05', '2023-01-24 10:50:16');
INSERT INTO `territory` VALUES (14, 'T000016', 'Territory 14', 2, 14, 2, 14, '2022-12-20 01:02:05', '2023-01-24 10:50:17');
INSERT INTO `territory` VALUES (15, 'T000017', 'Territory 15', 2, 15, 2, 15, '2022-12-20 01:02:05', '2023-01-24 10:50:17');
INSERT INTO `territory` VALUES (16, 'T000018', 'Territory 16', 2, 16, 2, 16, '2022-12-20 01:02:05', '2023-01-24 10:50:17');
INSERT INTO `territory` VALUES (17, 'T000019', 'Territory 17', 2, 17, 2, 17, '2022-12-20 01:02:05', '2023-01-24 10:50:17');
INSERT INTO `territory` VALUES (18, 'T000020', 'Territory 18', 2, 18, 2, 18, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (19, 'T000021', 'Territory 19', 2, 19, 2, 19, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (20, 'T000022', 'Territory 20', 2, 20, 2, 20, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (21, 'T000023', 'Territory 21', 2, 21, 2, 21, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (22, 'T000024', 'Territory 22', 2, 22, 2, 22, '2022-12-20 01:02:05', '2023-01-24 10:50:17');
INSERT INTO `territory` VALUES (23, 'T000025', 'Territory 23', 2, 23, 2, 23, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (24, 'T000026', 'Territory 24', 2, 24, 2, 24, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (25, 'T000027', 'Territory 25', 2, 25, 2, 25, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (26, 'T000028', 'Territory 26', 2, 26, 2, 26, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (27, 'T000029', 'Territory 27', 2, 27, 2, 27, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (28, 'T000030', 'Territory 28', 2, 28, 2, 28, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (29, 'T000031', 'Territory 29', 2, 29, 2, 29, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (30, 'T000032', 'Territory 30', 2, 30, 2, 30, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (31, 'T000033', 'Territory 31', 2, 31, 2, 31, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (32, 'T000034', 'Territory 32', 2, 32, 2, 32, '2022-12-20 01:02:05', '2023-01-24 10:50:17');
INSERT INTO `territory` VALUES (33, 'T000035', 'Territory 33', 2, 33, 2, 33, '2022-12-20 01:02:05', '2023-01-24 10:50:17');
INSERT INTO `territory` VALUES (34, 'T000036', 'Territory 34', 2, 34, 2, 34, '2022-12-20 01:02:05', '2023-01-24 10:50:17');
INSERT INTO `territory` VALUES (35, 'T000037', 'Territory 35', 2, 35, 2, 35, '2022-12-20 01:02:05', '2023-01-24 10:50:17');
INSERT INTO `territory` VALUES (36, 'T000038', 'Territory 36', 2, 36, 2, 36, '2022-12-20 01:02:05', '2023-01-24 10:50:17');
INSERT INTO `territory` VALUES (37, 'T000039', 'Territory 37', 2, 37, 2, 37, '2022-12-20 01:02:05', '2023-01-24 10:50:18');
INSERT INTO `territory` VALUES (38, 'T000040', 'Territory 38', 2, 38, 2, 38, '2022-12-20 01:02:05', '2023-01-24 10:50:18');
INSERT INTO `territory` VALUES (39, 'T000041', 'Territory 39', 2, 39, 2, 39, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (40, 'T000042', 'Territory 40', 2, 40, 2, 40, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (41, 'T000043', 'Territory 41', 2, 41, 2, 41, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (42, 'T000044', 'Territory 42', 2, 42, 2, 42, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (43, 'T000045', 'Territory 43', 2, 43, 2, 43, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (44, 'T000046', 'Territory 44', 2, 44, 2, 44, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (45, 'T000001', 'Territory 45', 2, 45, 3, 1, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (46, 'T000002', 'Territory 46', 2, 46, 3, 2, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (47, 'T000003', 'Territory 47', 2, 47, 3, 3, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (48, 'T000004', 'Territory 48', 2, 48, 3, 4, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (49, 'T000005', 'Territory 49', 2, 49, 3, 5, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (50, 'T000006', 'Territory 50', 2, 50, 3, 6, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (51, 'T000007', 'Territory 51', 2, 51, 3, 7, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (52, 'T000008', 'Territory 52', 2, 52, 3, 8, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (53, 'T000009', 'Territory 53', 2, 53, 3, 9, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (54, 'T000010', 'Territory 54', 2, 54, 3, 10, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (55, 'T000011', 'Territory 55', 2, 55, 3, 11, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (56, 'T000012', 'Territory 56', 2, 56, 3, 12, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (57, 'T000015', 'Territory 57', 2, 57, 3, 13, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (58, 'T000016', 'Territory 58', 2, 58, 3, 14, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (59, 'T000017', 'Territory 59', 2, 59, 3, 15, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (60, 'T000018', 'Territory 60', 2, 60, 3, 16, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (61, 'T000019', 'Territory 61', 2, 61, 3, 17, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (62, 'T000020', 'Territory 62', 2, 62, 3, 18, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (63, 'T000021', 'Territory 63', 2, 63, 3, 19, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (64, 'T000022', 'Territory 64', 2, 64, 3, 20, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (65, 'T000023', 'Territory 65', 2, 65, 3, 21, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (66, 'T000024', 'Territory 66', 2, 66, 3, 22, '2022-12-20 01:02:05', '2023-01-24 10:50:19');
INSERT INTO `territory` VALUES (67, 'T000025', 'Territory 67', 2, 67, 3, 23, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (68, 'T000026', 'Territory 68', 2, 68, 3, 24, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (69, 'T000027', 'Territory 69', 2, 69, 3, 25, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (70, 'T000028', 'Territory 70', 2, 70, 3, 26, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (71, 'T000029', 'Territory 71', 2, 71, 3, 27, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (72, 'T000030', 'Territory 72', 2, 72, 3, 28, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (73, 'T000031', 'Territory 73', 2, 73, 3, 29, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (74, 'T000032', 'Territory 74', 2, 74, 3, 30, '2022-12-20 01:02:05', '2023-01-24 10:50:19');
INSERT INTO `territory` VALUES (75, 'T000033', 'Territory 75', 2, 75, 3, 31, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (76, 'T000034', 'Territory 76', 2, 76, 3, 32, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (77, 'T000035', 'Territory 77', 2, 77, 3, 33, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (78, 'T000036', 'Territory 78', 2, 78, 3, 34, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (79, 'T000037', 'Territory 79', 2, 79, 3, 35, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (80, 'T000038', 'Territory 80', 2, 80, 3, 36, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (81, 'T000039', 'Territory 81', 2, 81, 3, 37, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (82, 'T000040', 'Territory 82', 2, 82, 3, 38, '2022-12-20 01:02:05', '2023-01-24 10:50:19');
INSERT INTO `territory` VALUES (83, 'T000041', 'Territory 83', 2, 83, 3, 39, '2022-12-20 01:02:05', '2023-01-24 10:50:19');
INSERT INTO `territory` VALUES (84, 'T000042', 'Territory 84', 2, 84, 3, 40, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (85, 'T000043', 'Territory 85', 2, 85, 3, 41, '2022-12-20 01:02:05', '2023-01-24 10:50:20');
INSERT INTO `territory` VALUES (86, 'T000044', 'Territory 86', 2, 86, 3, 42, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (87, 'T000045', 'Territory 87', 2, 87, 3, 43, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (88, 'T000046', 'Territory 88', 2, 88, 3, 44, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (89, 'T000001', 'Territory 89', 2, 89, 4, 1, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (90, 'T000002', 'Territory 90', 2, 90, 4, 2, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (91, 'T000003', 'Territory 91', 2, 91, 4, 3, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (92, 'T000004', 'Territory 92', 2, 92, 4, 4, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (93, 'T000005', 'Territory 93', 2, 93, 4, 5, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (94, 'T000006', 'Territory 94', 2, 94, 4, 6, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (95, 'T000007', 'Territory 95', 2, 95, 4, 7, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (96, 'T000008', 'Territory 96', 2, 96, 4, 8, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (97, 'T000009', 'Territory 97', 2, 97, 4, 9, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (98, 'T000010', 'Territory 98', 2, 98, 4, 10, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (99, 'T000011', 'Territory 99', 2, 99, 4, 11, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (100, 'T000012', 'Territory 100', 2, 100, 4, 12, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (101, 'T000015', 'Territory 101', 2, 101, 4, 13, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (102, 'T000016', 'Territory 102', 2, 102, 4, 14, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (103, 'T000017', 'Territory 103', 2, 103, 4, 15, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (104, 'T000018', 'Territory 104', 2, 104, 4, 16, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (105, 'T000019', 'Territory 105', 2, 105, 4, 17, '2022-12-20 01:02:05', '2023-01-24 10:50:20');
INSERT INTO `territory` VALUES (106, 'T000020', 'Territory 106', 2, 106, 4, 18, '2022-12-20 01:02:05', '2023-01-24 10:50:20');
INSERT INTO `territory` VALUES (107, 'T000021', 'Territory 107', 2, 107, 4, 19, '2022-12-20 01:02:05', '2023-01-24 10:50:20');
INSERT INTO `territory` VALUES (108, 'T000022', 'Territory 108', 2, 108, 4, 20, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (109, 'T000023', 'Territory 109', 2, 109, 4, 21, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (110, 'T000024', 'Territory 110', 2, 110, 4, 22, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (111, 'T000025', 'Territory 111', 2, 111, 4, 23, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (112, 'T000026', 'Territory 112', 2, 112, 4, 24, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (113, 'T000027', 'Territory 113', 2, 113, 4, 25, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (114, 'T000028', 'Territory 114', 2, 114, 4, 26, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (115, 'T000029', 'Territory 115', 2, 115, 4, 27, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (116, 'T000030', 'Territory 116', 2, 116, 4, 28, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (117, 'T000031', 'Territory 117', 2, 117, 4, 29, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (118, 'T000032', 'Territory 118', 2, 118, 4, 30, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (119, 'T000033', 'Territory 119', 2, 119, 4, 31, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (120, 'T000034', 'Territory 120', 2, 120, 4, 32, '2022-12-20 01:02:05', '2023-01-24 10:50:21');
INSERT INTO `territory` VALUES (121, 'T000035', 'Territory 121', 2, 121, 4, 33, '2022-12-20 01:02:05', '2023-01-24 10:50:21');
INSERT INTO `territory` VALUES (122, 'T000036', 'Territory 122', 2, 122, 4, 34, '2022-12-20 01:02:05', '2023-01-24 10:50:21');
INSERT INTO `territory` VALUES (123, 'T000037', 'Territory 123', 2, 123, 4, 35, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (124, 'T000038', 'Territory 124', 2, 124, 4, 36, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (125, 'T000039', 'Territory 125', 2, 125, 4, 37, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (126, 'T000040', 'Territory 126', 2, 126, 4, 38, '2022-12-20 01:02:05', '2023-01-24 10:50:21');
INSERT INTO `territory` VALUES (127, 'T000041', 'Territory 127', 2, 127, 4, 39, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (128, 'T000042', 'Territory 128', 2, 128, 4, 40, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (129, 'T000043', 'Territory 129', 2, 129, 4, 41, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (130, 'T000044', 'Territory 130', 2, 130, 4, 42, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (131, 'T000045', 'Territory 131', 2, 131, 4, 43, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (132, 'T000046', 'Territory 132', 2, 132, 4, 44, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (133, 'T000001', 'Territory 133', 2, 133, 5, 1, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (134, 'T000002', 'Territory 134', 2, 134, 5, 2, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (135, 'T000003', 'Territory 135', 2, 135, 5, 3, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (136, 'T000004', 'Territory 136', 2, 136, 5, 4, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (137, 'T000005', 'Territory 137', 2, 137, 5, 5, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (138, 'T000006', 'Territory 138', 2, 138, 5, 6, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (139, 'T000007', 'Territory 139', 2, 139, 5, 7, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (140, 'T000008', 'Territory 140', 2, 140, 5, 8, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (141, 'T000009', 'Territory 141', 2, 141, 5, 9, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (142, 'T000010', 'Territory 142', 2, 142, 5, 10, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (143, 'T000011', 'Territory 143', 2, 143, 5, 11, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (144, 'T000012', 'Territory 144', 2, 144, 5, 12, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (145, 'T000015', 'Territory 145', 2, 145, 5, 13, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (146, 'T000016', 'Territory 146', 2, 146, 5, 14, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (147, 'T000017', 'Territory 147', 2, 147, 5, 15, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (148, 'T000018', 'Territory 148', 2, 148, 5, 16, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (149, 'T000019', 'Territory 149', 2, 149, 5, 17, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (150, 'T000020', 'Territory 150', 2, 150, 5, 18, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (151, 'T000021', 'Territory 151', 2, 151, 5, 19, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (152, 'T000022', 'Territory 152', 2, 152, 5, 20, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (153, 'T000023', 'Territory 153', 2, 153, 5, 21, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (154, 'T000024', 'Territory 154', 2, 154, 5, 22, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (155, 'T000025', 'Territory 155', 2, 155, 5, 23, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (156, 'T000026', 'Territory 156', 2, 156, 5, 24, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (157, 'T000027', 'Territory 157', 2, 157, 5, 25, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (158, 'T000028', 'Territory 158', 2, 158, 5, 26, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (159, 'T000029', 'Territory 159', 2, 159, 5, 27, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (160, 'T000030', 'Territory 160', 2, 160, 5, 28, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (161, 'T000031', 'Territory 161', 2, 161, 5, 29, '2022-12-20 01:02:05', '2023-01-24 10:50:23');
INSERT INTO `territory` VALUES (162, 'T000032', 'Territory 162', 2, 162, 5, 30, '2022-12-20 01:02:05', '2023-01-24 10:50:23');
INSERT INTO `territory` VALUES (163, 'T000033', 'Territory 163', 2, 163, 5, 31, '2022-12-20 01:02:05', '2023-01-24 10:50:23');
INSERT INTO `territory` VALUES (164, 'T000034', 'Territory 164', 2, 164, 5, 32, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (165, 'T000035', 'Territory 165', 2, 165, 5, 33, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (166, 'T000036', 'Territory 166', 2, 166, 5, 34, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (167, 'T000037', 'Territory 167', 2, 167, 5, 35, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (168, 'T000038', 'Territory 168', 2, 168, 5, 36, '2022-12-20 01:02:05', '2023-01-24 10:50:23');
INSERT INTO `territory` VALUES (169, 'T000039', 'Territory 169', 2, 169, 5, 37, '2022-12-20 01:02:05', '2023-01-24 10:50:23');
INSERT INTO `territory` VALUES (170, 'T000040', 'Territory 170', 2, 170, 5, 38, '2022-12-20 01:02:05', '2023-01-24 10:50:23');
INSERT INTO `territory` VALUES (171, 'T000041', 'Territory 171', 2, 171, 5, 39, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (172, 'T000042', 'Territory 172', 2, 172, 5, 40, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (173, 'T000043', 'Territory 173', 2, 173, 5, 41, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (174, 'T000044', 'Territory 174', 2, 174, 5, 42, '2022-12-20 01:02:05', '2023-01-24 10:50:23');
INSERT INTO `territory` VALUES (175, 'T000045', 'Territory 175', 2, 175, 5, 43, '2022-12-20 01:02:05', '2023-01-24 10:50:23');
INSERT INTO `territory` VALUES (176, 'T000046', 'Territory 176', 2, 176, 5, 44, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (177, 'T000001', 'Territory 177', 2, 177, 6, 1, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (178, 'T000002', 'Territory 178', 2, 178, 6, 2, '2022-12-20 01:02:05', '2023-01-06 03:44:20');
INSERT INTO `territory` VALUES (179, 'T000003', 'Territory 179', 2, 179, 6, 3, '2022-12-20 01:02:05', '2023-01-24 10:50:23');
INSERT INTO `territory` VALUES (180, 'T000004', 'Territory 180', 2, 180, 6, 4, '2022-12-20 01:02:05', '2023-01-24 10:50:23');

-- ----------------------------
-- Table structure for uom
-- ----------------------------
DROP TABLE IF EXISTS `uom`;
CREATE TABLE `uom`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `description` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT NULL,
  `status` tinyint(1) NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of uom
-- ----------------------------
INSERT INTO `uom` VALUES (1, 'UOM001', 'Uom 1', 1, '2022-12-23 01:53:04', NULL);
INSERT INTO `uom` VALUES (2, 'UOM002', 'Uom 2', 1, '2022-12-23 01:53:12', NULL);

-- ----------------------------
-- Table structure for wrt
-- ----------------------------
DROP TABLE IF EXISTS `wrt`;
CREATE TABLE `wrt`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `region_id` bigint UNSIGNED NULL DEFAULT NULL,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '',
  `start_time` time NULL DEFAULT NULL,
  `end_time` time NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `code`(`code` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of wrt
-- ----------------------------
INSERT INTO `wrt` VALUES (1, 2, 'WRT0001', '04:00:00', '05:00:00');
INSERT INTO `wrt` VALUES (2, 2, 'WRT0002', '05:00:00', '07:00:00');
INSERT INTO `wrt` VALUES (3, 2, 'WRT0003', '07:00:00', '09:00:00');
INSERT INTO `wrt` VALUES (4, 2, 'WRT0004', '09:00:00', '11:00:00');
INSERT INTO `wrt` VALUES (5, 2, 'WRT0005', '11:00:00', '13:00:00');
INSERT INTO `wrt` VALUES (6, 2, 'WRT0006', '13:00:00', '15:00:00');
INSERT INTO `wrt` VALUES (7, 2, 'WRT0007', '08:00:00', '17:00:00');
INSERT INTO `wrt` VALUES (8, 2, 'WRT0008', '10:00:00', '15:00:00');
INSERT INTO `wrt` VALUES (9, 3, 'WRT0009', '04:00:00', '05:00:00');
INSERT INTO `wrt` VALUES (10, 3, 'WRT0010', '05:00:00', '07:00:00');
INSERT INTO `wrt` VALUES (11, 3, 'WRT0011', '07:00:00', '09:00:00');
INSERT INTO `wrt` VALUES (12, 3, 'WRT0012', '09:00:00', '11:00:00');
INSERT INTO `wrt` VALUES (13, 3, 'WRT0013', '11:00:00', '13:00:00');
INSERT INTO `wrt` VALUES (14, 3, 'WRT0014', '13:00:00', '15:00:00');
INSERT INTO `wrt` VALUES (15, 3, 'WRT0015', '08:00:00', '17:00:00');
INSERT INTO `wrt` VALUES (16, 3, 'WRT0016', '10:00:00', '15:00:00');
INSERT INTO `wrt` VALUES (17, 4, 'WRT0017', '04:00:00', '05:00:00');
INSERT INTO `wrt` VALUES (18, 4, 'WRT0018', '05:00:00', '07:00:00');
INSERT INTO `wrt` VALUES (19, 4, 'WRT0019', '07:00:00', '09:00:00');
INSERT INTO `wrt` VALUES (20, 4, 'WRT0020', '09:00:00', '11:00:00');
INSERT INTO `wrt` VALUES (21, 4, 'WRT0021', '11:00:00', '13:00:00');
INSERT INTO `wrt` VALUES (22, 4, 'WRT0022', '13:00:00', '15:00:00');
INSERT INTO `wrt` VALUES (23, 4, 'WRT0023', '08:00:00', '17:00:00');
INSERT INTO `wrt` VALUES (24, 4, 'WRT0024', '10:00:00', '15:00:00');
INSERT INTO `wrt` VALUES (25, 5, 'WRT0025', '04:00:00', '05:00:00');
INSERT INTO `wrt` VALUES (26, 5, 'WRT0026', '05:00:00', '07:00:00');
INSERT INTO `wrt` VALUES (27, 5, 'WRT0027', '07:00:00', '09:00:00');
INSERT INTO `wrt` VALUES (28, 5, 'WRT0028', '09:00:00', '11:00:00');
INSERT INTO `wrt` VALUES (29, 5, 'WRT0029', '11:00:00', '13:00:00');
INSERT INTO `wrt` VALUES (30, 5, 'WRT0030', '13:00:00', '15:00:00');
INSERT INTO `wrt` VALUES (31, 5, 'WRT0031', '08:00:00', '17:00:00');
INSERT INTO `wrt` VALUES (32, 5, 'WRT0032', '10:00:00', '15:00:00');
INSERT INTO `wrt` VALUES (33, 6, 'WRT0033', '04:00:00', '05:00:00');
INSERT INTO `wrt` VALUES (34, 6, 'WRT0034', '05:00:00', '07:00:00');
INSERT INTO `wrt` VALUES (35, 6, 'WRT0035', '07:00:00', '09:00:00');
INSERT INTO `wrt` VALUES (36, 6, 'WRT0036', '09:00:00', '11:00:00');
INSERT INTO `wrt` VALUES (37, 6, 'WRT0037', '11:00:00', '13:00:00');
INSERT INTO `wrt` VALUES (38, 6, 'WRT0038', '13:00:00', '15:00:00');
INSERT INTO `wrt` VALUES (39, 6, 'WRT0039', '08:00:00', '17:00:00');
INSERT INTO `wrt` VALUES (40, 6, 'WRT0040', '10:00:00', '15:00:00');

SET FOREIGN_KEY_CHECKS = 1;
