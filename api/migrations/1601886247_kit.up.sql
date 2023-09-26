-- MySQL dump 10.13  Distrib 5.7.13, for linux-glibc2.5 (x86_64)
--
-- Host: localhost    Database: eden_v2
-- ------------------------------------------------------
-- Server version	5.7.32-0ubuntu0.16.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Temporary view structure for view `adm_division`
--

DROP TABLE IF EXISTS `adm_division`;
/*!50001 DROP VIEW IF EXISTS `adm_division`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `adm_division` AS SELECT
 1 AS `area_id`,
 1 AS `area_name`,
 1 AS `sub_district_id`,
 1 AS `sub_district_name`,
 1 AS `postal_code`,
 1 AS `district_id`,
 1 AS `district_name`,
 1 AS `city_id`,
 1 AS `city_name`,
 1 AS `province_id`,
 1 AS `province_name`,
 1 AS `country_id`,
 1 AS `country_name`,
 1 AS `concate_address`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `archetype`
--

DROP TABLE IF EXISTS `archetype`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `archetype` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `business_type_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `alias_name_idn` varchar(100) DEFAULT '',
  `abbreviation` varchar(3) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `aux_data` tinyint(1) DEFAULT '2' COMMENT 'auxiliary data',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_archetype_1_idx` (`business_type_id`),
  CONSTRAINT `fk_archetype_1` FOREIGN KEY (`business_type_id`) REFERENCES `business_type` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `area`
--

DROP TABLE IF EXISTS `area`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `area` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `aux_data` tinyint(1) DEFAULT '2' COMMENT 'auxiliary data',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `area_policy`
--

DROP TABLE IF EXISTS `area_policy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `area_policy` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `area_id` bigint(20) unsigned DEFAULT NULL,
  `min_order` decimal(12,2) DEFAULT '0.00',
  `delivery_fee` decimal(12,2) DEFAULT '0.00',
  `order_time_limit` varchar(5) DEFAULT '',
  `default_price_set` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_delivery_config_1_idx` (`area_id`),
  KEY `fk_area_policy_2` (`default_price_set`),
  CONSTRAINT `fk_area_policy_1` FOREIGN KEY (`area_id`) REFERENCES `area` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_area_policy_2` FOREIGN KEY (`default_price_set`) REFERENCES `price_set` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_delivery_config_1` FOREIGN KEY (`area_id`) REFERENCES `area` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `audit_log`
--

DROP TABLE IF EXISTS `audit_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `audit_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `staff_id` bigint(20) unsigned DEFAULT NULL,
  `merchant_id` bigint(20) unsigned DEFAULT NULL,
  `ref_id` bigint(20) unsigned DEFAULT NULL,
  `type` varchar(30) DEFAULT '',
  `function` varchar(30) DEFAULT '',
  `timestamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_audit_log_1_idx` (`staff_id`),
  KEY `fk_audit_log_2_idx` (`merchant_id`),
  CONSTRAINT `fk_audit_log_1` FOREIGN KEY (`staff_id`) REFERENCES `staff` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_audit_log_2` FOREIGN KEY (`merchant_id`) REFERENCES `merchant` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `branch`
--

DROP TABLE IF EXISTS `branch`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `branch` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint(20) unsigned DEFAULT NULL,
  `area_id` bigint(20) unsigned DEFAULT NULL,
  `archetype_id` bigint(20) unsigned DEFAULT NULL,
  `price_set_id` bigint(20) unsigned DEFAULT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `salesperson_id` bigint(20) unsigned DEFAULT NULL,
  `sub_district_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `pic_name` varchar(100) DEFAULT '',
  `phone_number` varchar(15) DEFAULT '',
  `alt_phone_number` varchar(100) DEFAULT '',
  `address_name` varchar(100) DEFAULT '',
  `shipping_address` varchar(350) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `main_branch` tinyint(1) DEFAULT '0',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_branch_1_idx` (`merchant_id`),
  KEY `fk_branch_2_idx` (`area_id`),
  KEY `fk_branch_3_idx` (`archetype_id`),
  KEY `fk_branch_4_idx` (`price_set_id`),
  KEY `fk_branch_5_idx` (`warehouse_id`),
  KEY `fk_branch_6_idx` (`salesperson_id`),
  KEY `fk_branch_7_idx` (`sub_district_id`),
  CONSTRAINT `fk_branch_1` FOREIGN KEY (`merchant_id`) REFERENCES `merchant` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_branch_2` FOREIGN KEY (`area_id`) REFERENCES `area` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_branch_3` FOREIGN KEY (`archetype_id`) REFERENCES `archetype` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_branch_4` FOREIGN KEY (`price_set_id`) REFERENCES `price_set` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_branch_5` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_branch_6` FOREIGN KEY (`salesperson_id`) REFERENCES `staff` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_branch_7` FOREIGN KEY (`sub_district_id`) REFERENCES `sub_district` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `business_type`
--

DROP TABLE IF EXISTS `business_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `business_type` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `aux_data` tinyint(1) DEFAULT '2' COMMENT 'auxiliary data',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `category`
--

DROP TABLE IF EXISTS `category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `category` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='product category';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `city`
--

DROP TABLE IF EXISTS `city`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `city` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `province_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `value` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_city_1_idx` (`province_id`),
  CONSTRAINT `fk_city_1` FOREIGN KEY (`province_id`) REFERENCES `province` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=516 DEFAULT CHARSET=utf8 COMMENT='kabupaten/kota';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `code_generator`
--

DROP TABLE IF EXISTS `code_generator`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `code_generator` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(100) NOT NULL,
  `code_name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ui_code_generator_1_idx` (`code`,`code_name`),
  KEY `xi_code_generator_1_idx` (`code`),
  KEY `xi_code_generator_2_idx` (`code_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `config_app`
--

DROP TABLE IF EXISTS `config_app`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `config_app` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `application` tinyint(2) DEFAULT '0',
  `field` varchar(50) DEFAULT '',
  `attribute` varchar(30) DEFAULT '',
  `value` varchar(50) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8 COMMENT='application configuration';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country`
--

DROP TABLE IF EXISTS `country`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `value` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='provinsi';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `delivery_order`
--

DROP TABLE IF EXISTS `delivery_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `delivery_order` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sales_order_id` bigint(20) unsigned DEFAULT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `wrt_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `doc_status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `shipping_address` varchar(350) DEFAULT '',
  `receipt_note` varchar(250) DEFAULT '',
  `total_weight` decimal(5,2) DEFAULT '0.00',
  `delta_print` int(3) DEFAULT '0',
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_delivery_order_1_idx` (`sales_order_id`),
  KEY `fk_delivery_order_2_idx` (`warehouse_id`),
  KEY `fk_delivery_order_3_idx` (`wrt_id`),
  CONSTRAINT `fk_delivery_order_1` FOREIGN KEY (`sales_order_id`) REFERENCES `sales_order` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_delivery_order_2` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_delivery_order_3` FOREIGN KEY (`wrt_id`) REFERENCES `wrt` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `delivery_order_item`
--

DROP TABLE IF EXISTS `delivery_order_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `delivery_order_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `delivery_order_id` bigint(20) unsigned DEFAULT NULL,
  `sales_order_item_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `deliver_qty` decimal(10,2) DEFAULT '0.00',
  `receive_qty` decimal(10,2) DEFAULT '0.00',
  `receipt_item_note` varchar(100) DEFAULT '',
  `order_item_note` varchar(100) DEFAULT '',
  `weight` decimal(5,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_delivery_order_item_1_idx` (`delivery_order_id`),
  KEY `fk_delivery_order_item_2_idx` (`sales_order_item_id`),
  KEY `fk_delivery_order_item_3_idx` (`product_id`),
  CONSTRAINT `fk_delivery_order_item_1` FOREIGN KEY (`delivery_order_id`) REFERENCES `delivery_order` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_delivery_order_item_2` FOREIGN KEY (`sales_order_item_id`) REFERENCES `sales_order_item` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_delivery_order_item_3` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `delivery_return`
--

DROP TABLE IF EXISTS `delivery_return`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `delivery_return` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `delivery_order_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `doc_status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_delivery_return_1_idx` (`warehouse_id`),
  KEY `fk_delivery_return_2_idx` (`delivery_order_id`),
  CONSTRAINT `fk_delivery_return_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_delivery_return_2` FOREIGN KEY (`delivery_order_id`) REFERENCES `delivery_order` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `delivery_return_item`
--

DROP TABLE IF EXISTS `delivery_return_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `delivery_return_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `delivery_return_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `return_good_qty` decimal(10,2) DEFAULT '0.00',
  `return_waste_qty` decimal(10,2) DEFAULT '0.00',
  `unit_cost` decimal(12,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_delivery_return_item_1_idx` (`delivery_return_id`),
  KEY `fk_delivery_return_item_2_idx` (`product_id`),
  CONSTRAINT `fk_delivery_return_item_1` FOREIGN KEY (`delivery_return_id`) REFERENCES `delivery_return` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_delivery_return_item_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `district`
--

DROP TABLE IF EXISTS `district`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `district` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `city_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `value` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_district_1_idx` (`city_id`),
  CONSTRAINT `fk_district_1` FOREIGN KEY (`city_id`) REFERENCES `city` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=7069 DEFAULT CHARSET=utf8 COMMENT='kecamatan';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `division`
--

DROP TABLE IF EXISTS `division`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `division` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `glossary`
--

DROP TABLE IF EXISTS `glossary`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `glossary` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `table` varchar(30) DEFAULT '',
  `attribute` varchar(30) DEFAULT '',
  `value_int` tinyint(3) DEFAULT '0',
  `value_name` varchar(50) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=94 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goods_receipt`
--

DROP TABLE IF EXISTS `goods_receipt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goods_receipt` (
  `id` bigint(20) unsigned NOT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `purchase_order_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `ata_date` date DEFAULT NULL COMMENT 'actual time arrival',
  `ata_time` varchar(5) DEFAULT '',
  `total_weight` decimal(5,2) DEFAULT '0.00',
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_goods_receipt_1_idx` (`warehouse_id`),
  KEY `fk_goods_receipt_2_idx` (`purchase_order_id`),
  CONSTRAINT `fk_goods_receipt_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_goods_receipt_2` FOREIGN KEY (`purchase_order_id`) REFERENCES `purchase_order` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goods_receipt_item`
--

DROP TABLE IF EXISTS `goods_receipt_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goods_receipt_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `goods_receipt_id` bigint(20) unsigned DEFAULT NULL,
  `purchase_order_item_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `deliver_qty` decimal(10,2) DEFAULT '0.00',
  `reject_qty` decimal(10,2) DEFAULT '0.00',
  `receive_qty` decimal(10,2) DEFAULT '0.00',
  `weight` decimal(5,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_goods_receipt_item_1_idx` (`goods_receipt_id`),
  KEY `fk_goods_receipt_item_2_idx` (`purchase_order_item_id`),
  KEY `fk_goods_receipt_item_3_idx` (`product_id`),
  CONSTRAINT `fk_goods_receipt_item_1` FOREIGN KEY (`goods_receipt_id`) REFERENCES `goods_receipt` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_goods_receipt_item_2` FOREIGN KEY (`purchase_order_item_id`) REFERENCES `purchase_order_item` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_goods_receipt_item_3` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goods_transfer`
--

DROP TABLE IF EXISTS `goods_transfer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goods_transfer` (
  `id` bigint(20) unsigned NOT NULL,
  `origin_id` bigint(20) unsigned DEFAULT NULL,
  `destination_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `doc_status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `eta_date` date DEFAULT NULL,
  `eta_time` varchar(5) DEFAULT '',
  `total_weight` decimal(5,2) DEFAULT '0.00',
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_goods_transfer_1_idx` (`origin_id`),
  KEY `fk_goods_transfer_2_idx` (`destination_id`),
  CONSTRAINT `fk_goods_transfer_1` FOREIGN KEY (`origin_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_goods_transfer_2` FOREIGN KEY (`destination_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `goods_transfer_item`
--

DROP TABLE IF EXISTS `goods_transfer_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `goods_transfer_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `goods_transfer_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `deliver_qty` decimal(10,2) DEFAULT '0.00',
  `receive_qty` decimal(10,2) DEFAULT '0.00',
  `receive_note` varchar(100) DEFAULT '',
  `weight` decimal(5,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_goods_transfer_item_1_idx` (`goods_transfer_id`),
  KEY `fk_goods_transfer_item_2_idx` (`product_id`),
  CONSTRAINT `fk_goods_transfer_item_1` FOREIGN KEY (`goods_transfer_id`) REFERENCES `goods_transfer` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_goods_transfer_item_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `merchant`
--

DROP TABLE IF EXISTS `merchant`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `merchant` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_merchant_id` bigint(20) unsigned DEFAULT NULL,
  `term_invoice_sls_id` bigint(20) unsigned DEFAULT NULL,
  `term_payment_sls_id` bigint(20) unsigned DEFAULT NULL,
  `payment_method_id` bigint(20) unsigned DEFAULT NULL,
  `business_type_id` bigint(20) unsigned DEFAULT NULL,
  `finance_area_id` bigint(20) unsigned DEFAULT NULL,
  `tag_customer` varchar(100) DEFAULT '',
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `pic_name` varchar(100) DEFAULT '',
  `phone_number` varchar(15) DEFAULT '',
  `alt_phone_number` varchar(100) DEFAULT '',
  `email` varchar(100) DEFAULT '',
  `billing_address` varchar(350) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_merchant_1_idx` (`user_merchant_id`),
  KEY `fk_merchant_3_idx` (`term_payment_sls_id`),
  KEY `fk_merchant_4_idx` (`payment_method_id`),
  KEY `fk_merchant_5_idx` (`business_type_id`),
  KEY `fk_merchant_2_idx` (`term_invoice_sls_id`),
  KEY `fk_merchant_6_idx` (`finance_area_id`),
  CONSTRAINT `fk_merchant_1` FOREIGN KEY (`user_merchant_id`) REFERENCES `user_merchant` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_merchant_2` FOREIGN KEY (`term_invoice_sls_id`) REFERENCES `term_invoice_sls` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_merchant_3` FOREIGN KEY (`term_payment_sls_id`) REFERENCES `term_payment_sls` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_merchant_4` FOREIGN KEY (`payment_method_id`) REFERENCES `payment_method` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_merchant_5` FOREIGN KEY (`business_type_id`) REFERENCES `business_type` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_merchant_6` FOREIGN KEY (`finance_area_id`) REFERENCES `area` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `notification`
--

DROP TABLE IF EXISTS `notification`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `notification` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `type` tinyint(1) DEFAULT '0',
  `title` varchar(50) DEFAULT '',
  `description` varchar(200) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `notification_log`
--

DROP TABLE IF EXISTS `notification_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `notification_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint(20) unsigned DEFAULT NULL,
  `ref_id` varchar(45) DEFAULT '',
  `type` tinyint(1) DEFAULT '0',
  `message` varchar(300) DEFAULT '',
  `read` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_notification_log_1_idx` (`merchant_id`),
  CONSTRAINT `fk_notification_log_1` FOREIGN KEY (`merchant_id`) REFERENCES `user_merchant` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `payment_method`
--

DROP TABLE IF EXISTS `payment_method`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_method` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `permission`
--

DROP TABLE IF EXISTS `permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `permission` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `value` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ui_permission_1_idx` (`value`),
  UNIQUE KEY `ui_permission_2_idx` (`code`),
  UNIQUE KEY `ui_permission_3_idx` (`name`),
  KEY `fk_permission_1_idx` (`parent_id`),
  CONSTRAINT `fk_permission_1` FOREIGN KEY (`parent_id`) REFERENCES `permission` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=329 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `price`
--

DROP TABLE IF EXISTS `price`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `price` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `price_set_id` bigint(20) unsigned DEFAULT NULL,
  `unit_price` decimal(12,2) DEFAULT '0.00',
  `shadow_price` decimal(12,2) DEFAULT '0.00',
  `shadow_price_pct` int(3) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_product_price_2_idx` (`price_set_id`),
  KEY `fk_product_price_1_idx` (`product_id`),
  CONSTRAINT `fk_product_price_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_product_price_2` FOREIGN KEY (`price_set_id`) REFERENCES `price_set` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `price_log`
--

DROP TABLE IF EXISTS `price_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `price_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `price_id` bigint(20) unsigned DEFAULT NULL,
  `unit_price` decimal(12,2) DEFAULT '0.00',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_price_log_1_idx` (`price_id`),
  KEY `fk_price_log_2_idx` (`created_by`),
  CONSTRAINT `fk_price_log_1` FOREIGN KEY (`price_id`) REFERENCES `price` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `price_set`
--

DROP TABLE IF EXISTS `price_set`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `price_set` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `product`
--

DROP TABLE IF EXISTS `product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `uom_id` bigint(20) unsigned DEFAULT NULL,
  `category_id` bigint(20) unsigned DEFAULT NULL,
  `warehouse_sto` varchar(100) DEFAULT '',
  `warehouse_pur` varchar(100) DEFAULT '',
  `warehouse_sal` varchar(100) DEFAULT '',
  `tag_product` varchar(100) DEFAULT '',
  `code` varchar(50) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `unit_weight` decimal(5,2) DEFAULT '0.00' COMMENT 'conversion to kg',
  `storability` tinyint(1) DEFAULT '0',
  `purchasability` tinyint(1) DEFAULT '0',
  `salability` tinyint(1) DEFAULT '0',
  `description` varchar(500) DEFAULT NULL COMMENT 'external use',
  `note` varchar(250) DEFAULT NULL COMMENT 'internal use',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_item_1_idx` (`uom_id`),
  KEY `fk_item_2_idx` (`category_id`),
  CONSTRAINT `fk_item_1` FOREIGN KEY (`uom_id`) REFERENCES `uom` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_2` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `product_image`
--

DROP TABLE IF EXISTS `product_image`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_image` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `image_url` varchar(300) DEFAULT '',
  `main_image` tinyint(1) DEFAULT '2',
  PRIMARY KEY (`id`),
  KEY `fk_product_image_1_idx` (`product_id`),
  CONSTRAINT `fk_product_image_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `prospect_customer`
--

DROP TABLE IF EXISTS `prospect_customer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `prospect_customer` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sub_district_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `business_type_name` varchar(50) DEFAULT '',
  `pic_name` varchar(100) DEFAULT '',
  `phone_number` varchar(15) DEFAULT '',
  `street_address` varchar(350) DEFAULT '',
  `time_consent` tinyint(1) DEFAULT '0',
  `reg_status` tinyint(1) DEFAULT '0' COMMENT 'registration status',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='prospective customer';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `prospect_supplier`
--

DROP TABLE IF EXISTS `prospect_supplier`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `prospect_supplier` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sub_district_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `phone_number` varchar(15) DEFAULT '',
  `alt_phone_number` varchar(50) DEFAULT '',
  `street_address` varchar(350) DEFAULT '',
  `pic_name` varchar(100) DEFAULT '',
  `pic_phone_number` varchar(15) DEFAULT '',
  `commodity` varchar(250) DEFAULT '',
  `time_consent` tinyint(1) DEFAULT '0',
  `reg_status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='prospective supplier';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `province`
--

DROP TABLE IF EXISTS `province`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `province` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `country_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `value` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_province_1_idx` (`country_id`),
  CONSTRAINT `fk_province_1` FOREIGN KEY (`country_id`) REFERENCES `country` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8 COMMENT='provinsi';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `purchase_invoice`
--

DROP TABLE IF EXISTS `purchase_invoice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `purchase_invoice` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `purchase_order_id` bigint(20) unsigned DEFAULT NULL,
  `term_payment_pur_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `due_date` date DEFAULT NULL,
  `tax_pct` decimal(5,2) DEFAULT '0.00',
  `tax_amount` decimal(15,2) DEFAULT '0.00',
  `delivery_fee` decimal(10,2) DEFAULT '0.00',
  `adjustment` tinyint(1) DEFAULT '0',
  `adj_amount` decimal(20,2) DEFAULT '0.00',
  `total_price` decimal(20,2) DEFAULT '0.00',
  `total_charge` decimal(20,2) DEFAULT '0.00',
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_purchase_invoice_1_idx` (`purchase_order_id`),
  KEY `fk_purchase_invoice_2_idx` (`term_payment_pur_id`),
  CONSTRAINT `fk_purchase_invoice_1` FOREIGN KEY (`purchase_order_id`) REFERENCES `purchase_order` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_invoice_2` FOREIGN KEY (`term_payment_pur_id`) REFERENCES `term_payment_pur` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `purchase_invoice_item`
--

DROP TABLE IF EXISTS `purchase_invoice_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `purchase_invoice_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `purchase_invoice_id` bigint(20) unsigned DEFAULT NULL,
  `purchase_order_item_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `invoice_qty` decimal(10,2) DEFAULT '0.00',
  `unit_price` decimal(12,2) DEFAULT '0.00',
  `subtotal` decimal(15,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_purchase_invoice_item_1_idx` (`purchase_invoice_id`),
  KEY `fk_purchase_invoice_item_2_idx` (`purchase_order_item_id`),
  KEY `fk_purchase_invoice_item_3_idx` (`product_id`),
  CONSTRAINT `fk_purchase_invoice_item_1` FOREIGN KEY (`purchase_invoice_id`) REFERENCES `purchase_invoice` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_invoice_item_2` FOREIGN KEY (`purchase_order_item_id`) REFERENCES `purchase_order_item` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_invoice_item_3` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `purchase_order`
--

DROP TABLE IF EXISTS `purchase_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `purchase_order` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `supplier_id` bigint(20) unsigned DEFAULT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `term_payment_pur_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT NULL,
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `eta_date` date DEFAULT NULL,
  `eta_time` varchar(5) DEFAULT NULL,
  `tax_pct` int(3) unsigned DEFAULT '0',
  `delivery_fee` decimal(10,2) DEFAULT '0.00',
  `total_price` decimal(20,2) DEFAULT '0.00',
  `total_charge` decimal(20,2) DEFAULT '0.00',
  `total_invoice` decimal(20,2) DEFAULT '0.00',
  `total_weight` decimal(5,2) DEFAULT '0.00',
  `note` varchar(250) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_purchase_order_1_idx` (`supplier_id`),
  KEY `fk_purchase_order_2_idx` (`warehouse_id`),
  KEY `fk_purchase_order_3_idx` (`term_payment_pur_id`),
  CONSTRAINT `fk_purchase_order_1` FOREIGN KEY (`supplier_id`) REFERENCES `supplier` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_order_2` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_order_3` FOREIGN KEY (`term_payment_pur_id`) REFERENCES `term_payment_pur` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `purchase_order_item`
--

DROP TABLE IF EXISTS `purchase_order_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `purchase_order_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `purchase_order_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `order_qty` decimal(10,2) DEFAULT '0.00',
  `receive_qty` decimal(10,2) DEFAULT '0.00',
  `unit_price` decimal(12,2) DEFAULT '0.00',
  `subtotal` decimal(15,2) DEFAULT '0.00',
  `weight` decimal(5,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_purchase_order_item_1_idx` (`purchase_order_id`),
  KEY `fk_purchase_order_item_2_idx` (`product_id`),
  CONSTRAINT `fk_purchase_order_item_1` FOREIGN KEY (`purchase_order_id`) REFERENCES `purchase_order` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_order_item_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `purchase_payment`
--

DROP TABLE IF EXISTS `purchase_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `purchase_payment` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `purchase_invoice_id` bigint(20) unsigned DEFAULT NULL,
  `payment_method_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT NULL,
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `amount` decimal(20,2) DEFAULT '0.00',
  `paid_off` tinyint(1) DEFAULT '0',
  `note` varchar(250) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_purchase_payment_1_idx` (`purchase_invoice_id`),
  KEY `fk_purchase_payment_2_idx` (`payment_method_id`),
  CONSTRAINT `fk_purchase_payment_1` FOREIGN KEY (`purchase_invoice_id`) REFERENCES `purchase_invoice` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_payment_2` FOREIGN KEY (`payment_method_id`) REFERENCES `payment_method` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `role`
--

DROP TABLE IF EXISTS `role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `division_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_role_1_idx` (`division_id`),
  CONSTRAINT `fk_role_1` FOREIGN KEY (`division_id`) REFERENCES `division` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `role_permission`
--

DROP TABLE IF EXISTS `role_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role_permission` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` bigint(20) unsigned DEFAULT NULL,
  `permission_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_role_permission_1_idx` (`role_id`),
  KEY `fk_role_permission_2_idx` (`permission_id`),
  CONSTRAINT `fk_role_permission_1` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_role_permission_2` FOREIGN KEY (`permission_id`) REFERENCES `permission` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sales_inv_recap`
--

DROP TABLE IF EXISTS `sales_inv_recap`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sales_inv_recap` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint(20) unsigned DEFAULT NULL,
  `issuer_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `due_date` date DEFAULT NULL,
  `billing_address` varchar(350) DEFAULT '',
  `grand_total` decimal(20,2) DEFAULT '0.00',
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_sales_inv_recap_1_idx` (`merchant_id`),
  KEY `fk_sales_inv_recap_2_idx` (`issuer_id`),
  CONSTRAINT `fk_sales_inv_recap_1` FOREIGN KEY (`merchant_id`) REFERENCES `merchant` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_inv_recap_2` FOREIGN KEY (`issuer_id`) REFERENCES `staff` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='sales invoice recap';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sales_inv_recap_item`
--

DROP TABLE IF EXISTS `sales_inv_recap_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sales_inv_recap_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sales_invoice_id` bigint(20) unsigned DEFAULT NULL,
  `delivery_order_id` bigint(20) unsigned DEFAULT NULL,
  `branch_id` bigint(20) unsigned DEFAULT NULL,
  `subtotal` decimal(20,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_sales_inv_recap_item_1_idx` (`sales_invoice_id`),
  KEY `fk_sales_inv_recap_item_2_idx` (`delivery_order_id`),
  KEY `fk_sales_inv_recap_item_3_idx` (`branch_id`),
  CONSTRAINT `fk_sales_inv_recap_item_1` FOREIGN KEY (`sales_invoice_id`) REFERENCES `sales_invoice` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_inv_recap_item_2` FOREIGN KEY (`delivery_order_id`) REFERENCES `delivery_order` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_inv_recap_item_3` FOREIGN KEY (`branch_id`) REFERENCES `branch` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sales_invoice`
--

DROP TABLE IF EXISTS `sales_invoice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sales_invoice` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sales_order_id` bigint(20) unsigned DEFAULT NULL,
  `term_payment_sls_id` bigint(20) unsigned DEFAULT NULL,
  `voucher_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `billing_address` varchar(350) DEFAULT '',
  `delivery_fee` decimal(10,2) DEFAULT '0.00',
  `vou_redeem_code` varchar(20) DEFAULT '',
  `vou_disc_amount` decimal(20,2) DEFAULT '0.00',
  `adjustment` tinyint(1) DEFAULT '0',
  `adj_amount` decimal(20,2) DEFAULT '0.00',
  `total_price` decimal(20,2) DEFAULT '0.00',
  `total_charge` decimal(20,2) DEFAULT '0.00',
  `delta_print` int(3) DEFAULT '0',
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_sales_invoice_1_idx` (`sales_order_id`),
  KEY `fk_sales_invoice_2_idx` (`term_payment_sls_id`),
  CONSTRAINT `fk_sales_invoice_1` FOREIGN KEY (`sales_order_id`) REFERENCES `sales_order` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_invoice_2` FOREIGN KEY (`term_payment_sls_id`) REFERENCES `term_payment_sls` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sales_invoice_item`
--

DROP TABLE IF EXISTS `sales_invoice_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sales_invoice_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sales_invoice_id` bigint(20) unsigned DEFAULT NULL,
  `sales_order_item_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `invoice_qty` decimal(10,2) DEFAULT '0.00',
  `unit_price` decimal(12,2) DEFAULT '0.00',
  `subtotal` decimal(15,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_sales_invoice_item_1_idx` (`sales_invoice_id`),
  KEY `fk_sales_invoice_item_2_idx` (`sales_order_item_id`),
  KEY `fk_sales_invoice_item_3_idx` (`product_id`),
  CONSTRAINT `fk_sales_invoice_item_1` FOREIGN KEY (`sales_invoice_id`) REFERENCES `sales_invoice` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_invoice_item_2` FOREIGN KEY (`sales_order_item_id`) REFERENCES `sales_order_item` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_invoice_item_3` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sales_order`
--

DROP TABLE IF EXISTS `sales_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sales_order` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `branch_id` bigint(20) unsigned DEFAULT NULL,
  `term_payment_sls_id` bigint(20) unsigned DEFAULT NULL,
  `term_invoice_sls_id` bigint(20) unsigned DEFAULT NULL,
  `salesperson_id` bigint(20) unsigned DEFAULT NULL,
  `sub_district_id` bigint(20) unsigned DEFAULT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `wrt_id` bigint(20) unsigned DEFAULT NULL,
  `area_id` bigint(20) unsigned DEFAULT NULL COMMENT 'not FK',
  `voucher_id` bigint(20) unsigned DEFAULT NULL COMMENT 'not FK',
  `price_set_id` bigint(20) unsigned DEFAULT NULL COMMENT 'not FK',
  `archetype_id` bigint(20) DEFAULT '0' COMMENT 'not FK',
  `business_type_id` bigint(20) DEFAULT '0' COMMENT 'not FK',
  `order_channel` tinyint(1) DEFAULT '0',
  `code` varchar(50) DEFAULT '',
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `delivery_date` date DEFAULT NULL,
  `shipping_address` varchar(350) DEFAULT '',
  `delivery_fee` decimal(10,2) DEFAULT '0.00',
  `vou_redeem_code` varchar(20) DEFAULT '',
  `vou_disc_amount` decimal(20,2) DEFAULT '0.00',
  `total_price` decimal(20,2) DEFAULT '0.00',
  `total_charge` decimal(20,2) DEFAULT '0.00',
  `total_weight` decimal(5,2) DEFAULT '0.00',
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_sales_order_1_idx` (`branch_id`),
  KEY `fk_sales_order_2_idx` (`term_payment_sls_id`),
  KEY `fk_sales_order_3_idx` (`term_invoice_sls_id`),
  KEY `fk_sales_order_4_idx` (`salesperson_id`),
  KEY `fk_sales_order_5_idx` (`sub_district_id`),
  KEY `fk_sales_order_6_idx` (`warehouse_id`),
  KEY `fk_sales_order_7_idx` (`wrt_id`),
  CONSTRAINT `fk_sales_order_1` FOREIGN KEY (`branch_id`) REFERENCES `branch` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_2` FOREIGN KEY (`term_payment_sls_id`) REFERENCES `term_payment_sls` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_3` FOREIGN KEY (`term_invoice_sls_id`) REFERENCES `term_invoice_sls` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_4` FOREIGN KEY (`salesperson_id`) REFERENCES `staff` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_5` FOREIGN KEY (`sub_district_id`) REFERENCES `sub_district` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_6` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_7` FOREIGN KEY (`wrt_id`) REFERENCES `wrt` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sales_order_item`
--

DROP TABLE IF EXISTS `sales_order_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sales_order_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sales_order_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `order_qty` decimal(10,2) DEFAULT '0.00',
  `unit_price` decimal(12,2) DEFAULT '0.00',
  `shadow_price` decimal(12,2) DEFAULT '0.00',
  `subtotal` decimal(15,2) DEFAULT '0.00',
  `weight` decimal(5,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_sales_order_item_1_idx` (`sales_order_id`),
  KEY `fk_sales_order_item_2_idx` (`product_id`),
  CONSTRAINT `fk_sales_order_item_1` FOREIGN KEY (`sales_order_id`) REFERENCES `sales_order` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_order_item_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sales_payment`
--

DROP TABLE IF EXISTS `sales_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sales_payment` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sales_invoice_id` bigint(20) unsigned DEFAULT NULL,
  `payment_method_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `amount` decimal(20,2) DEFAULT '0.00',
  `paid_off` tinyint(1) DEFAULT '0',
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_sales_payment_1_idx` (`sales_invoice_id`),
  KEY `fk_sales_payment_2_idx` (`payment_method_id`),
  CONSTRAINT `fk_sales_payment_1` FOREIGN KEY (`sales_invoice_id`) REFERENCES `sales_invoice` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sales_payment_2` FOREIGN KEY (`payment_method_id`) REFERENCES `payment_method` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `staff`
--

DROP TABLE IF EXISTS `staff`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `staff` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `area_id` bigint(20) unsigned DEFAULT NULL,
  `parent_id` bigint(20) unsigned DEFAULT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `display_name` varchar(100) DEFAULT '',
  `employee_code` varchar(50) DEFAULT '',
  `role_group` tinyint(1) DEFAULT '0',
  `phone_number` varchar(15) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_staff_1_idx` (`role_id`),
  KEY `fk_staff_2_idx` (`user_id`),
  KEY `fk_staff_3_idx` (`area_id`),
  KEY `fk_staff_4_idx` (`parent_id`),
  KEY `fk_staff_5_idx` (`warehouse_id`),
  CONSTRAINT `fk_staff_1` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_staff_2` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_staff_3` FOREIGN KEY (`area_id`) REFERENCES `area` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_staff_4` FOREIGN KEY (`parent_id`) REFERENCES `staff` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_staff_5` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `stock`
--

DROP TABLE IF EXISTS `stock`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `stock` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `available_stock` decimal(10,2) DEFAULT '0.00',
  `waste_stock` decimal(10,2) DEFAULT '0.00',
  `safety_stock` decimal(10,2) DEFAULT '0.00',
  `commited_in_stock` decimal(10,2) DEFAULT '0.00',
  `commited_out_stock` decimal(10,2) DEFAULT '0.00',
  `salable` tinyint(1) DEFAULT '0',
  `purchasable` tinyint(1) DEFAULT '0',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_stock_1_idx` (`product_id`),
  KEY `fk_stock_2_idx` (`warehouse_id`),
  CONSTRAINT `fk_stock_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_2` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `stock_log`
--

DROP TABLE IF EXISTS `stock_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `stock_log` (
  `id` bigint(20) unsigned NOT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `ref_id` bigint(20) unsigned DEFAULT NULL,
  `ref_type` tinyint(1) DEFAULT '0',
  `type` tinyint(1) DEFAULT '0',
  `initial_stock` decimal(10,2) DEFAULT '0.00',
  `quantity` decimal(10,2) DEFAULT '0.00',
  `final_stock` decimal(10,2) DEFAULT '0.00',
  `unit_cost` decimal(12,2) DEFAULT '0.00',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_stock_log_1_idx` (`warehouse_id`),
  KEY `fk_stock_log_2_idx` (`product_id`),
  CONSTRAINT `fk_stock_log_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_log_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `stock_opname`
--

DROP TABLE IF EXISTS `stock_opname`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `stock_opname` (
  `id` bigint(20) unsigned NOT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `category_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_stock_opname_1_idx` (`warehouse_id`),
  KEY `fk_stock_opname_2_idx` (`category_id`),
  CONSTRAINT `fk_stock_opname_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_opname_2` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `stock_opname_item`
--

DROP TABLE IF EXISTS `stock_opname_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `stock_opname_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `stock_opname_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `initial_stock` decimal(20,2) DEFAULT '0.00',
  `adjust_qty` decimal(20,2) DEFAULT '0.00' COMMENT 'adjustment quantity',
  `final_stock` decimal(20,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_stock_opname_item_1_idx` (`stock_opname_id`),
  KEY `fk_stock_opname_item_2_idx` (`product_id`),
  CONSTRAINT `fk_stock_opname_item_1` FOREIGN KEY (`stock_opname_id`) REFERENCES `stock_opname` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_opname_item_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sub_district`
--

DROP TABLE IF EXISTS `sub_district`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sub_district` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `district_id` bigint(20) unsigned DEFAULT NULL,
  `area_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `value` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `postal_code` varchar(10) DEFAULT '',
  `concat_address` varchar(250) DEFAULT '' COMMENT 'concatenated address',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_sub_district_1_idx` (`district_id`),
  KEY `fk_sub_district_2_idx` (`area_id`),
  CONSTRAINT `fk_sub_district_1` FOREIGN KEY (`district_id`) REFERENCES `district` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_sub_district_2` FOREIGN KEY (`area_id`) REFERENCES `area` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=82103 DEFAULT CHARSET=utf8 COMMENT='kelurahan';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `supplier`
--

DROP TABLE IF EXISTS `supplier`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `supplier` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `supplier_type_id` bigint(20) unsigned DEFAULT NULL,
  `term_payment_pur_id` bigint(20) unsigned DEFAULT NULL,
  `payment_method_id` bigint(20) unsigned DEFAULT NULL,
  `sub_district_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `email` varchar(100) DEFAULT '',
  `phone_number` varchar(15) DEFAULT '',
  `alt_phone_number` varchar(50) DEFAULT '',
  `pic_name` varchar(100) DEFAULT '',
  `address` varchar(350) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_supplier_1_idx` (`supplier_type_id`),
  KEY `fk_supplier_3_idx` (`payment_method_id`),
  KEY `fk_supplier_4_idx` (`sub_district_id`),
  KEY `fk_supplier_2_idx` (`term_payment_pur_id`),
  CONSTRAINT `fk_supplier_1` FOREIGN KEY (`supplier_type_id`) REFERENCES `supplier_type` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_supplier_2` FOREIGN KEY (`term_payment_pur_id`) REFERENCES `term_payment_pur` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_supplier_3` FOREIGN KEY (`payment_method_id`) REFERENCES `payment_method` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_supplier_4` FOREIGN KEY (`sub_district_id`) REFERENCES `sub_district` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `supplier_type`
--

DROP TABLE IF EXISTS `supplier_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `supplier_type` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `abbreviation` varchar(3) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tag_customer`
--

DROP TABLE IF EXISTS `tag_customer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tag_customer` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tag_product`
--

DROP TABLE IF EXISTS `tag_product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tag_product` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `image_url` varchar(300) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `term_condition`
--

DROP TABLE IF EXISTS `term_condition`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `term_condition` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `application` tinyint(1) DEFAULT '0',
  `title` varchar(50) DEFAULT '',
  `attribute` varchar(50) DEFAULT '',
  `value` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `term_invoice_sls`
--

DROP TABLE IF EXISTS `term_invoice_sls`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `term_invoice_sls` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `term_payment_pur`
--

DROP TABLE IF EXISTS `term_payment_pur`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `term_payment_pur` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `days_value` int(2) DEFAULT '0',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `term_payment_sls`
--

DROP TABLE IF EXISTS `term_payment_sls`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `term_payment_sls` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `days_value` int(2) DEFAULT '0',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `uom`
--

DROP TABLE IF EXISTS `uom`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `uom` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `decimal_enabled` tinyint(1) DEFAULT '0',
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COMMENT='unit of measurement';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `email` varchar(100) DEFAULT '',
  `password` varchar(250) DEFAULT '',
  `last_login_at` timestamp NULL DEFAULT NULL,
  `note` varchar(250) DEFAULT '',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='internal user';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_merchant`
--

DROP TABLE IF EXISTS `user_merchant`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_merchant` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(50) DEFAULT '',
  `uid` varchar(100) DEFAULT '',
  `password` varchar(250) DEFAULT '',
  `firebase_id` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `verification` tinyint(1) DEFAULT '0',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_permission`
--

DROP TABLE IF EXISTS `user_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_permission` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `permission_id` bigint(20) unsigned DEFAULT NULL,
  `permission_value` varchar(50) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_user_permission_1_idx` (`user_id`),
  CONSTRAINT `fk_user_permission_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `voucher`
--

DROP TABLE IF EXISTS `voucher`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `voucher` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `area_id` bigint(20) unsigned DEFAULT NULL,
  `archetype_id` bigint(20) unsigned DEFAULT NULL,
  `tag_customer` varchar(100) DEFAULT '',
  `code` varchar(50) DEFAULT '',
  `redeem_code` varchar(20) DEFAULT '',
  `type` tinyint(1) DEFAULT '0',
  `name` varchar(100) DEFAULT '',
  `start_timestamp` varchar(45) DEFAULT '',
  `end_timestamp` varchar(45) DEFAULT '',
  `overall_quota` int(10) DEFAULT '0',
  `user_quota` int(10) DEFAULT '0',
  `rem_overall_quota` int(10) DEFAULT '0' COMMENT 'remaining overall quota',
  `min_order` decimal(20,2) DEFAULT '0.00' COMMENT 'minimum order',
  `disc_amount` decimal(20,2) DEFAULT '0.00',
  `note` varchar(250) DEFAULT '',
  `void_reason` tinyint(1) DEFAULT '0',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_promo_1_idx` (`area_id`),
  KEY `fk_promo_2_idx` (`archetype_id`),
  CONSTRAINT `fk_promo_1` FOREIGN KEY (`area_id`) REFERENCES `area` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_promo_2` FOREIGN KEY (`archetype_id`) REFERENCES `archetype` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `voucher_log`
--

DROP TABLE IF EXISTS `voucher_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `voucher_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `voucher_id` bigint(20) unsigned DEFAULT NULL,
  `merchant_id` bigint(20) unsigned DEFAULT NULL,
  `branch_id` bigint(20) unsigned DEFAULT NULL,
  `sales_order_id` bigint(20) unsigned DEFAULT NULL,
  `tag_customer` varchar(100) DEFAULT '',
  `vou_disc_amount` decimal(20,2) DEFAULT '0.00',
  `timestamp` timestamp NULL DEFAULT NULL,
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_voucher_log_1_idx` (`voucher_id`),
  KEY `fk_voucher_log_2_idx` (`merchant_id`),
  KEY `fk_voucher_log_3_idx` (`branch_id`),
  KEY `fk_voucher_log_4_idx` (`sales_order_id`),
  CONSTRAINT `fk_voucher_log_1` FOREIGN KEY (`voucher_id`) REFERENCES `voucher` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_voucher_log_2` FOREIGN KEY (`merchant_id`) REFERENCES `merchant` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_voucher_log_3` FOREIGN KEY (`branch_id`) REFERENCES `branch` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_voucher_log_4` FOREIGN KEY (`sales_order_id`) REFERENCES `sales_order` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `warehouse`
--

DROP TABLE IF EXISTS `warehouse`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `warehouse` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `area_id` bigint(20) unsigned DEFAULT NULL,
  `sub_district_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `pic_name` varchar(100) DEFAULT '',
  `phone_number` varchar(15) DEFAULT '',
  `alt_phone_number` varchar(50) DEFAULT '',
  `street_address` varchar(350) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `main_warehouse` tinyint(1) DEFAULT '0',
  `aux_data` tinyint(1) DEFAULT '2' COMMENT 'auxiliary data',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_warehouse_1_idx` (`area_id`),
  KEY `fk_warehouse_2_idx` (`sub_district_id`),
  CONSTRAINT `fk_warehouse_1` FOREIGN KEY (`area_id`) REFERENCES `area` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_warehouse_2` FOREIGN KEY (`sub_district_id`) REFERENCES `sub_district` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `waste_disposal`
--

DROP TABLE IF EXISTS `waste_disposal`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `waste_disposal` (
  `id` bigint(20) unsigned NOT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_waste_disposal_1_idx` (`warehouse_id`),
  CONSTRAINT `fk_waste_disposal_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `waste_disposal_item`
--

DROP TABLE IF EXISTS `waste_disposal_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `waste_disposal_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `waste_disposal_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `dispose_qty` decimal(10,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_waste_disposal_item_1_idx` (`waste_disposal_id`),
  KEY `fk_waste_disposal_item_2_idx` (`product_id`),
  CONSTRAINT `fk_waste_disposal_item_1` FOREIGN KEY (`waste_disposal_id`) REFERENCES `waste_disposal` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_waste_disposal_item_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `waste_entry`
--

DROP TABLE IF EXISTS `waste_entry`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `waste_entry` (
  `id` bigint(20) unsigned NOT NULL,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `status` tinyint(2) DEFAULT '0',
  `recognition_date` date DEFAULT NULL,
  `note` varchar(250) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_waste_entry_1_idx` (`warehouse_id`),
  CONSTRAINT `fk_waste_entry_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `waste_entry_item`
--

DROP TABLE IF EXISTS `waste_entry_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `waste_entry_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `waste_entry_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `waste_qty` decimal(10,2) DEFAULT '0.00',
  `note` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_waste_entry_item_2_idx` (`product_id`),
  KEY `fk_waste_entry_item_1_idx` (`waste_entry_id`),
  CONSTRAINT `fk_waste_entry_item_1` FOREIGN KEY (`waste_entry_id`) REFERENCES `waste_entry` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_waste_entry_item_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `waste_log`
--

DROP TABLE IF EXISTS `waste_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `waste_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `warehouse_id` bigint(20) unsigned DEFAULT NULL,
  `waste_id` bigint(20) unsigned DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `type` tinyint(1) DEFAULT '0',
  `initial_stock` decimal(10,2) DEFAULT '0.00',
  `quantity` decimal(10,2) DEFAULT '0.00',
  `final_stock` decimal(10,2) DEFAULT '0.00',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_waste_log_1_idx` (`warehouse_id`),
  KEY `fk_waste_log_2_idx` (`waste_id`),
  KEY `fk_waste_log_3_idx` (`product_id`),
  CONSTRAINT `fk_waste_log_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouse` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_waste_log_2` FOREIGN KEY (`waste_id`) REFERENCES `waste_entry` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_waste_log_3` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `wrt`
--

DROP TABLE IF EXISTS `wrt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `wrt` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `area_id` bigint(20) unsigned DEFAULT NULL,
  `code` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  `note` varchar(250) DEFAULT '',
  `aux_data` tinyint(1) DEFAULT '2' COMMENT 'auxiliary data',
  `status` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `fk_wrt_1_idx` (`area_id`),
  CONSTRAINT `fk_wrt_1` FOREIGN KEY (`area_id`) REFERENCES `area` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='window receiving time';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Final view structure for view `adm_division`
--

/*!50001 DROP VIEW IF EXISTS `adm_division`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8 */;
/*!50001 SET character_set_results     = utf8 */;
/*!50001 SET collation_connection      = utf8_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `adm_division` AS select `a`.`id` AS `area_id`,`a`.`name` AS `area_name`,`sd`.`id` AS `sub_district_id`,`sd`.`name` AS `sub_district_name`,`sd`.`postal_code` AS `postal_code`,`d`.`id` AS `district_id`,`d`.`name` AS `district_name`,`c`.`id` AS `city_id`,`c`.`name` AS `city_name`,`p`.`id` AS `province_id`,`p`.`name` AS `province_name`,`co`.`id` AS `country_id`,`co`.`name` AS `country_name`,`sd`.`concat_address` AS `concate_address` from (((((`sub_district` `sd` join `area` `a` on((`a`.`id` = `sd`.`area_id`))) join `district` `d` on((`d`.`id` = `sd`.`district_id`))) join `city` `c` on((`c`.`id` = `d`.`city_id`))) join `province` `p` on((`p`.`id` = `c`.`province_id`))) join `country` `co` on((`co`.`id` = `p`.`country_id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-11-25 11:30:40
