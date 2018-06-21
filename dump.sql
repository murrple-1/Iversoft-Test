DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `user_addresses`;
DROP TABLE IF EXISTS `user_roles`;

CREATE TABLE `user_roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `label` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO `user_roles` (`id`, `label`)
VALUES
	(1, 'Admin'),
	(2, 'Publisher'),
	(3, 'Public User');

CREATE TABLE `user_addresses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `address` varchar(255) DEFAULT NULL,
  `province` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  `postal_code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO `user_addresses` (`id`, `address`, `province`, `city`, `country`, `postal_code`)
VALUES
	(1, '123 fake street', 'Ontario', 'Ottawa', 'Canada', '123 w4t'),
	(2, '123 queen street', 'Quebec','Gatineau','Canada', '123 tdf'),
	(3, '123 major road', 'Ontariofdgdgdfg', 'Ottawa','Canada', '145 w4t'),
	(4, '123 blue street', 'Ontario', 'Ottawa', 'Canada', '145 lpo');

CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `user_roles_id` int NOT NULL,
  `user_addresses_id` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE (`username`),
  UNIQUE (`email`),
  FOREIGN KEY (`user_roles_id`) REFERENCES `user_roles` (`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_addresses_id`) REFERENCES `user_addresses` (`id`) ON DELETE CASCADE
);

INSERT INTO `users` (`id`, `user_roles_id`, `user_addresses_id`, `username`, `email`, `created_at`)
VALUES
	(1, 1, 1, 'I_Admin', 'admin@test.com', '2017-05-20 12:42:53'),
	(2, 2, 2, 'I_Publish', 'publisher@test.com', '2017-05-20 13:05:53'),
	(3, 3, 3, 'I_User', 'user@test.com', '2017-05-21 13:05:53'),
	(4, 3, 4, 'I_User_Too', 'user2@test.com', '2017-05-22 14:05:53');
