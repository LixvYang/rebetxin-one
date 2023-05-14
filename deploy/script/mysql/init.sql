CREATE DATABASE IF NOT EXISTS betxin;

CREATE TABLE
    IF NOT EXISTS `topic` (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
        `tid` VARCHAR(36) NOT NULL DEFAULT '',
        `cid` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
        `title` VARCHAR(50) NOT NULL DEFAULT '',
        `intro` VARCHAR(255) NOT NULL DEFAULT '',
        `content` VARCHAR(255) NOT NULL DEFAULT '',
        `yes_ratio` DECIMAL(5, 2) NOT NULL DEFAULT '50.00',
        `no_ratio` DECIMAL(5, 2) NOT NULL DEFAULT '50.00',
        `yes_price` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000',
        `no_price` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000',
        `total_price` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000',
        `collect_count` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
        `read_count` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
        `img_url` VARCHAR(255) NOT NULL DEFAULT '',
        `is_stop` TINYINT UNSIGNED NOT NULL DEFAULT '0',
        `refund_end_time` TIMESTAMP DEFAULT NULL,
        `end_time` TIMESTAMP DEFAULT NULL,
        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` TIMESTAMP DEFAULT NULL,
        PRIMARY KEY (`id`),
        INDEX `idx_cid` (`cid`),
        UNIQUE `idx_tid` (`tid`),
        UNIQUE `title_intro_content_index` (`title`, `intro`, `content`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE
    IF NOT EXISTS `category` (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
        `category_name` VARCHAR(30) NOT NULL DEFAULT '',
        PRIMARY KEY (`id`)
    ) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;

CREATE TABLE
    IF NOT EXISTS `user` (
        `id` int NOT NULL AUTO_INCREMENT,
        `identity_number` varchar(50) NOT NULL,
        `uid` varchar(36) NOT NULL,
        `full_name` varchar(255) DEFAULT NULL,
        `avatar_url` varchar(255) DEFAULT NULL,
        `session_id` varchar(255) DEFAULT NULL,
        `biography` varchar(255) DEFAULT NULL,
        `private_key` varchar(255) DEFAULT NULL,
        `client_id` VARCHAR(255) DEFAULT NULL,
        `contract` VARCHAR(255) DEFAULT NULL,
        `is_mvm_user` TINYINT(1) DEFAULT 0,
        `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_identity_number` (`identity_number`),
        UNIQUE KEY `idx_uid` (`uid`),
        KEY `idx_full_name` (`full_name`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `topicpurchase` (
        `id` bigint(20) NOT NULL AUTO_INCREMENT,
        `uid` varchar(36) NOT NULL DEFAULT '',
        `tid` varchar(255) NOT NULL DEFAULT '',
        `trace_id` varchar(36) NOT NULL DEFAULT '',
        `yes_price` decimal(32, 8) NOT NULL DEFAULT 0,
        `no_price` decimal(32, 8) NOT NULL DEFAULT 0,
        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` TIMESTAMP DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE KEY idx_uid_tid (uid, tid),
        KEY idx_tid (tid),
        KEY idx_yes_price (yes_price),
        KEY idx_no_price (no_price)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE
    `topic_collect` (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
        `uid` VARCHAR(36) NOT NULL DEFAULT '',
        `tid` VARCHAR(36) NOT NULL DEFAULT '',
        `status` TINYINT(1) NOT NULL DEFAULT 0,
        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_uid_tid` (`uid`, `tid`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE
    IF NOT EXISTS `refund` (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
        `uid` VARCHAR(36) NOT NULL DEFAULT '',
        `tid` VARCHAR(36) NOT NULL DEFAULT '',
        `trace_id` VARCHAR(36) NOT NULL DEFAULT '',
        `asset_id` VARCHAR(36) NOT NULL DEFAULT '',
        `yes_price` DECIMAL(32, 8) NOT NULL DEFAULT 0.00000000,
        `no_price` DECIMAL(32, 8) NOT NULL DEFAULT 0.00000000,
        `memo` VARCHAR(255) NOT NULL DEFAULT '',
        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_trace_id` (`trace_id`),
        UNIQUE KEY `idx_uid` (`uid`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
# User snapshort table
CREATE TABLE
    IF NOT EXISTS `snapshot` (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
        `trace_id` VARCHAR(36) NOT NULL DEFAULT '',
        `uid` VARCHAR(36) NOT NULL DEFAULT '',
        `tid` VARCHAR(36) NOT NULL DEFAULT '',
        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        INDEX `idx_trace_id` (`trace_id`),
        INDEX `idx_uid_tid` (`uid`, `tid`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

INSERT INTO category (id, category_name) VALUES (1, 'Buisiness');

INSERT INTO category (id, category_name) VALUES (2, 'Crypto');

INSERT INTO category (id, category_name) VALUES (3, 'Sports');

INSERT INTO category (id, category_name) VALUES (4, 'Politics');

INSERT INTO category (id, category_name) VALUES (5, 'New');

INSERT INTO category (id, category_name) VALUES (6, 'Trending');