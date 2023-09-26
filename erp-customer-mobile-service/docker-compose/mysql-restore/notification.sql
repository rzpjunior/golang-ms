/*
 Navicat Premium Data Transfer

 Source Server         : Local
 Source Server Type    : MySQL
 Source Server Version : 80031 (8.0.31)
 Source Host           : localhost:3306
 Source Schema         : notification

 Target Server Type    : MySQL
 Target Server Version : 80031 (8.0.31)
 File Encoding         : 65001

 Date: 05/05/2023 12:28:57
*/
USE notification;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for notification
-- ----------------------------
DROP TABLE IF EXISTS `notification`;
CREATE TABLE `notification`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `type` tinyint(1) NULL DEFAULT 0,
  `title` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `message` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT '',
  `status` tinyint(1) NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 34 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of notification
-- ----------------------------
INSERT INTO `notification` VALUES (1, 'NOT0001', 1, 'Notifikasi terbaru dari Eden Farm', 'Pesanan #sales_order_code# sudah kami terima.', 1);
INSERT INTO `notification` VALUES (2, 'NOT0002', 1, 'Notifikasi terbaru dari Eden Farm', 'Pembayaran #sales_order_code# sudah dikonfirmasi. Pesanan akan kami proses.', 1);
INSERT INTO `notification` VALUES (3, 'NOT0003', 1, 'Notifikasi terbaru dari Eden Farm', 'Pesanan #sales_order_code# sedang dikirimkan.', 1);
INSERT INTO `notification` VALUES (4, 'NOT0004', 1, 'Notifikasi terbaru dari Eden Farm', 'Pesanan #sales_order_code# sudah dibayar.', 1);
INSERT INTO `notification` VALUES (5, 'NOT0005', 1, 'Notifikasi terbaru dari Eden Farm', 'Pesanan #sales_order_code# dibatalkan.', 1);
INSERT INTO `notification` VALUES (6, 'NOT0006', 1, 'Notifikasi terbaru dari Eden Farm', 'Pesanan #sales_order_code# sudah selesai.', 1);
INSERT INTO `notification` VALUES (7, 'NOT0007', 1, 'Silakan lakukan pembayaran', 'Mohon melakukan pembayaran sebelum #current_date# #time_limit#. Bila tidak pesanan #sales_order_code# akan dibatalkan secara otomatis.', 1);
INSERT INTO `notification` VALUES (8, 'NOT0008', 1, 'Pesanan Telah Dibuat', 'Pesanan #sales_order_code# sudah kami terima, silahkan selesaikan pembayaranmu.', 1);
INSERT INTO `notification` VALUES (9, 'NOT0009', 3, 'Pesanan Telah Sampai', 'Pesananmu berhasil diantar. Ceritakan pengalaman belanjamu disini', 1);
INSERT INTO `notification` VALUES (10, 'NOT0010', 4, 'Notifikasi terbaru dari Eden Farm', 'Pesanan Sales Order anda telah ditolak oleh Supervisor anda', 1);
INSERT INTO `notification` VALUES (11, 'NOT0011', 4, 'Notifikasi terbaru dari Eden Farm', 'Pesanan Sales Order anda telah diterima oleh Supervisor anda', 1);
INSERT INTO `notification` VALUES (12, 'NOT0012', 5, 'Yay, Kamu dapat EdenPoint!', 'EdenPoint dari transaksimu sebelumnya telah masuk. Pakai EdenPointmu dan dapatkan potongan belanja langsung', 1);
INSERT INTO `notification` VALUES (13, 'NOT0013', 1, 'Update Order', 'Order #customer# - #sales_order_code#.', 1);
INSERT INTO `notification` VALUES (14, 'NOT0014', 1, 'Cancel Order', 'Order #customer# - #sales_order_code#.', 1);
INSERT INTO `notification` VALUES (15, 'NOT0015', 4, 'Notifikasi terbaru dari Eden Farm', 'Sales Order anda telah ditolak oleh Checker', 1);
INSERT INTO `notification` VALUES (16, 'NOT0016', 4, 'Notifikasi terbaru dari Eden Farm', 'Sales Order anda telah diterima oleh Checker', 1);
INSERT INTO `notification` VALUES (17, 'NOT0017', 2, 'Registered Prospect Customer', '#Name# registered successfully.', 1);
INSERT INTO `notification` VALUES (18, 'NOT0018', 2, 'Declined Prospect Customer', '#Name# has been declined, reason : #type#.', 1);
INSERT INTO `notification` VALUES (19, 'NOT0019', 6, '#supplier_name#', 'Hai #name#, ada tugas baru ini buat kamu. #smile#', 1);
INSERT INTO `notification` VALUES (20, 'NOT0020', 6, 'Ada Purchase Plan baru!', 'Purchase Plan dengan kode #purchase_plan_code# telah ditambahkan, segera tugaskan!', 1);
INSERT INTO `notification` VALUES (21, 'NOT0021', 6, 'Kamu mendapatkan tugas baru!', 'Tugas dengan kode #purchase_plan_code# telah ditambahkan.', 1);
INSERT INTO `notification` VALUES (22, 'NOT0022', 6, 'Penugasan dibatalkan!', 'Tugas dengan kode #purchase_plan_code# telah dibatalkan oleh #purchasing_manager_name# .', 1);
INSERT INTO `notification` VALUES (23, 'NOT0023', 6, 'Tugaskan ulang!', 'Tugaskan ulang untuk Purchase Plan dengan kode #purchase_plan_code# yang telah dibatalkan dari #field_purchaser_name#', 1);
INSERT INTO `notification` VALUES (24, 'NOT0024', 6, 'Purchase Plan telah dibatalkan!', 'Tugas dengan kode #purchase_plan_code# telah dibatalkan.', 1);
INSERT INTO `notification` VALUES (25, 'NOT0025', 6, 'Ada Purchase Plan baru!', 'Purchase Plan baru dengan kode #purchase_plan_code# telah ditugaskan ke #field_purchaser_name#, segera cek!', 1);
INSERT INTO `notification` VALUES (26, 'NOT0026', 4, 'Notifikasi terbaru dari Eden Farm', 'Anda mendapatkan tugas baru dari lead picker', 1);
INSERT INTO `notification` VALUES (27, 'NOT0027', 4, 'Notifikasi terbaru dari Eden Farm', 'Tugas anda telah dibatalkan oleh lead picker', 1);
INSERT INTO `notification` VALUES (28, 'NOT0028', 7, 'Bonus EdenPoint diterima', 'Selamat kamu mendapatkan Bonus EdenPoint dari program Eden Rewards, yuk belanja sekarang!', 1);
INSERT INTO `notification` VALUES (29, 'NOT0029', 8, 'Kejutan Karung Bawang terbuka', 'Ada hadiah spesial di karung kamu, cek di sini', 1);
INSERT INTO `notification` VALUES (30, 'NOT0030', 9, 'Selamat, kamu jadi Juragan', 'Sekarang kamu dapat CS prioritas dan voucher lebih banyak! Yuk, semangat tingkatkan transaksi Eden Rewards dan raih level tertinggi yaitu Konglomerat. Cek di sini ya', 1);
INSERT INTO `notification` VALUES (31, 'NOT0031', 10, 'Selamat, kamu jadi Konglomerat', 'Kamu sudah berada di level tertinggi Eden Rewards! Terus tingkatkan transaksi untuk mendapat bonus EdenPoint s.d. 175.000. Cek semua keuntunganmu di sini', 1);
INSERT INTO `notification` VALUES (32, 'NOT0032', 1, 'Notifikasi terbaru dari Eden Farm', 'Pesanan #sales_order_code# siap untuk kamu ambil.', 1);
INSERT INTO `notification` VALUES (33, 'NOT0033', 3, 'Pesanan Sudah Diambil', 'Pesananmu berhasil kamu ambil. Ceritakan pengalaman belanjamu disini.', 1);

SET FOREIGN_KEY_CHECKS = 1;
