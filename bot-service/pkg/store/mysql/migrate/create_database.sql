CREATE TABLE `orders`
(
    `id`          INT AUTO_INCREMENT,
    `uuid`        CHAR(36),
    `id_external` CHAR(36),
    `number`      CHAR(36),
    created_at    TIMESTAMP DEFAULT UTC_TIMESTAMP(),
    updated_at    TIMESTAMP DEFAULT UTC_TIMESTAMP() ON UPDATE UTC_TIMESTAMP(),
    `state`       CHAR(36),
    `comment`     varchar(2048),
    `total`       int,
    `phone`       varchar(16),
    PRIMARY KEY (`id`)
);

CREATE TABLE `order_positions`
(
    `id`                  bigint unsigned,
    `position`            int,
    `id_order`            varchar(100),
    `product_uuid`        CHAR(36),
    `product_name`        varchar(100),
    `measure_name`        varchar(16),
    `price`               int,
    `price_with_discount` int,
    `quantity`            int,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_orders_positions` FOREIGN KEY (`id`) REFERENCES `orders` (`id`)
);


CREATE TABLE `products`
(
    `name`         varchar(100),
    `store_id`     longtext,
    `uuid`         CHAR(36),
    `parent_uuid`  CHAR(36),
    `group`        boolean,
    `image`        varchar(100),
    `measure_name` varchar(16),
    `description`  varchar(2048),
    `price`        DECIMAL(10, 2)
);


CREATE TABLE `chats`
(
    `id`         bigint AUTO_INCREMENT,
    `user_name`  varchar(100),
    `user_phone` varchar(16),
    `chat_state` varchar(16),
    `order`      longtext,
    PRIMARY KEY (`id`)
)