CREATE TABLE `songs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `album_id` bigint NOT NULL,
  `title` varchar(255) NOT NULL,
  `duration` int DEFAULT NULL,  -- 单位：秒
  `track_number` int DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`album_id`) REFERENCES `albums`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


INSERT INTO `songs` (`album_id`, `title`, `duration`, `track_number`, `created_at`, `updated_at`)
VALUES
(3, 'the lakes', 240, 1, NOW(), NOW()),
(3, 'you all over me', 215, 2, NOW(), NOW()),
(3, 'it’s time to go', 250, 3, NOW(), NOW()),
(3, 'right where you left me', 260, 4, NOW(), NOW()),
(3, 'mad woman (demo)', 230, 5, NOW(), NOW());
