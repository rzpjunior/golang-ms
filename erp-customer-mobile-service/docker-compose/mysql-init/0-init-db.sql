/* note for db
note for DB:
- Account
- Audit = mongo 
- Configuration
- Bridge = dynamic
- Catalog = inventory 
- Campaign
- Promotion
- Sales
- Crm
- Settlement
- Storage
- Notification
*/

-- create databases
CREATE DATABASE IF NOT EXISTS `account`;
CREATE DATABASE IF NOT EXISTS `configuration`;
-- db dynamic - bridge
CREATE DATABASE IF NOT EXISTS `dynamic`;
-- db catalog - inventory
CREATE DATABASE IF NOT EXISTS `inventory`;
CREATE DATABASE IF NOT EXISTS `campaign`;
CREATE DATABASE IF NOT EXISTS `promotion`;
CREATE DATABASE IF NOT EXISTS `sales`;
CREATE DATABASE IF NOT EXISTS `crm`;
CREATE DATABASE IF NOT EXISTS `settlement`;
CREATE DATABASE IF NOT EXISTS `notification`;
CREATE DATABASE IF NOT EXISTS `mobile_customer`;

-- SET GITD
SET @@GLOBAL.gtid_purged = '';

-- create root user and grant rights
/*CREATE USER 'root'@'%' IDENTIFIED BY 'secret';
GRANT ALL ON *.* TO 'root'@'%';
CREATE USER 'edenfarm'@'%' IDENTIFIED BY 'secret'; */
-- GRANT ALL PRIVILEGES ON *.* TO 'edenfarm'@'%';