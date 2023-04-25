CREATE DATABASE IF NOT EXISTS sbot_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE sbot_db;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

select 'START MIGRATION' AS '';

CREATE TABLE `orders`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `uuid`        CHAR(36),
    `id_external` CHAR(36),
    `number`      CHAR(36),
    `created_at`    TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    `updated_at`    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `state`       CHAR(36),
    `comment`     varchar(2048),
    `total`       int UNSIGNED,
    `phone`       varchar(16),
    PRIMARY KEY (`id`)
);

CREATE TABLE `order_positions`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `id_order`            int UNSIGNED,
    `product_uuid`        CHAR(36),
    `product_name`        varchar(100),
    `measure_name`        varchar(16),
    `price`               int UNSIGNED,
    `price_with_discount` int UNSIGNED,
    `quantity`            int UNSIGNED,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_orders_positions` FOREIGN KEY (`id_order`) REFERENCES `orders` (`id`)
);


CREATE TABLE `products`
(
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `name`         varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `uuid`         CHAR(36),
    `parent_uuid`  CHAR(36),
    `group`        boolean,
    `image`        varchar(100),
    `measure_name` varchar(16),
    `description`  varchar(4096) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `price`        int UNSIGNED,
    `price_with_discount` int UNSIGNED,
     PRIMARY KEY (`id`)
);

CREATE TABLE `chats`
(
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `user_name`  varchar(100),
    `user_phone` varchar(16),
    `chat_state` varchar(100),
    `order`      longtext,
    PRIMARY KEY (`id`)
);


select 'END MIGRATION' AS '';


SET FOREIGN_KEY_CHECKS = 1;
