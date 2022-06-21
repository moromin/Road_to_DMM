CREATE TABLE `account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL UNIQUE,
  `password_hash` varchar(255) NOT NULL,
  `display_name` varchar(255),
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `followers_count` bigint(20) NOT NULL DEFAULT 0,
  `following_count` bigint(20) NOT NULL DEFAULT 0,
  `note` text,
  `avatar` text,
  `header` text,
  PRIMARY KEY (`id`)
);

CREATE TABLE `status` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `account_id` bigint(20) NOT NULL,
  `content` text NOT NULL,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_account_id` (`account_id`),
  CONSTRAINT `fk_status_account_id` FOREIGN KEY (`account_id`) REFERENCES  `account` (`id`)
);

-- TODO: automate migration
CREATE TABLE `follow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `follower_id` bigint(20) NOT NULL,
  `followee_id` bigint(20) NOT NULL,
  -- INDEX `idx_follower_id` (`follower_id`),
  -- INDEX `idx_followee_id` (`followee_id`),
  CONSTRAINT `fk_follower_account_id` FOREIGN KEY (`follower_id`) REFERENCES  `account` (`id`),
  CONSTRAINT `fk_followee_account_id` FOREIGN KEY (`followee_id`) REFERENCES  `account` (`id`),
  PRIMARY KEY (`id`),
  UNIQUE follow_combination (follower_id, followee_id)
);

CREATE TABLE `attachment` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `type` text NOT NULL,
  `url` text NOT NULL,
  `description` varchar(420),
  PRIMARY KEY (`id`)
);
