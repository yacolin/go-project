CREATE TABLE `photos` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `album_id` bigint NOT NULL,
  `title` varchar(100) NOT NULL,
  `url` varchar(255) NOT NULL,
  `description` varchar(500) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_album_id` (`album_id`),
  FOREIGN KEY (`album_id`) REFERENCES `albums`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 插入示例数据
INSERT INTO `photos` (`album_id`, `title`, `url`, `description`, `created_at`, `updated_at`)
VALUES
(1, 'Album Cover', 'https://example.com/photos/album1/cover.jpg', 'The main cover photo of the album', NOW(), NOW()),
(1, 'Behind the Scenes', 'https://example.com/photos/album1/bts.jpg', 'Recording session in the studio', NOW(), NOW()),
(2, 'Live Performance', 'https://example.com/photos/album2/live.jpg', 'Live performance at Madison Square Garden', NOW(), NOW()),
(2, 'Band Photo', 'https://example.com/photos/album2/band.jpg', 'Official band photo for the album', NOW(), NOW()),
(3, 'Studio Session', 'https://example.com/photos/album3/studio.jpg', 'Working on the new tracks in the studio', NOW(), NOW()); 