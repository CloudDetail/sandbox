CREATE DATABASE IF NOT EXISTS `sandbox` DEFAULT CHARACTER SET utf8;

-- 切换到数据库
USE `sandbox`;

-- 创建 users 表
CREATE TABLE IF NOT EXISTS `users` (
  `id` VARCHAR(36) NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `email` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 插入数据
INSERT INTO `users` (`id`, `name`, `email`) VALUES
('049c5806-0a6e-4011-bb58-10e35b3dff34','DB_Saved_User 8','db_saved_user8@apo.com'),
('3aad51f0-f050-4f3a-a3e7-dc63977d47f2','DB_Saved_User 7','db_saved_user7@apo.com'),
('7a4c6533-dab1-4f02-9063-68c17a016114','DB_Saved_User 10','db_saved_user10@apo.com'),
('7c9de1c3-e856-407f-820e-a10fd4f7c52f','DB_Saved_User 3','db_saved_user3@apo.com'),
('8c2c3e50-4c56-4002-9a63-7271ff3f053a','DB_Saved_User 5','db_saved_user5@apo.com'),
('95edabc6-b5ee-4898-b138-369c6f9c8b1a','DB_Saved_User 2','db_saved_user2@apo.com'),
('d990dd72-c662-42d6-96d5-e44d55dec8e8','DB_Saved_User 6','db_saved_user6@apo.com'),
('daadd4cd-dc62-497f-8912-18d586b158b9','DB_Saved_User 4','db_saved_user4@apo.com'),
('effa3f90-fcae-45c4-a167-30f330a2c41c','DB_Saved_User 1','db_saved_user1@apo.com'),
('fba6b264-2d4b-4c6a-b426-1dad543d31c6','DB_Saved_User 9','db_saved_user9@apo.com');