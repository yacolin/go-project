-- 评论表，photo与comment一对多
CREATE TABLE `comments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `photo_id` bigint NOT NULL,
  `content` varchar(500) NOT NULL,
  `author` varchar(100) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_photo_id` (`photo_id`),
  FOREIGN KEY (`photo_id`) REFERENCES `photos`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- mock数据
INSERT INTO `comments` (`photo_id`, `content`, `author`, `created_at`, `updated_at`) VALUES
(1, '很棒的照片！', 'Alice', NOW(), NOW()),
(1, '色彩很美，喜欢！', 'Bob', NOW(), NOW()),
(2, '这张照片拍得真好。', 'Charlie', NOW(), NOW()),
(2, '请问是用什么相机拍的？', 'Diana', NOW(), NOW()),
(3, '有故事感的一张图。', 'Eve', NOW(), NOW()),
(3, '光影处理很棒。', 'Frank', NOW(), NOW()),
(4, '很有氛围，点赞！', 'Grace', NOW(), NOW()),
(4, '请继续分享更多作品。', 'Heidi', NOW(), NOW()),
(5, '喜欢这个角度。', 'Ivan', NOW(), NOW()),
(5, '色调很舒服。', 'Judy', NOW(), NOW());
