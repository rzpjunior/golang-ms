SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `eden_v2`.`voucher`
CHANGE COLUMN `start_timestamp` `start_timestamp` TIMESTAMP NULL DEFAULT NULL ,
CHANGE COLUMN `end_timestamp` `end_timestamp` TIMESTAMP NULL DEFAULT NULL ;
