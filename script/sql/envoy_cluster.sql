CREATE TABLE IF NOT EXISTS `envoy_cluster` (

  `id` int NOT NULL AUTO_INCREMENT,
  `cluster_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0',
  `display_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0',
  `enable` tinyint NOT NULL DEFAULT '0',
  `is_delete` tinyint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL,
  `create_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0',
  `update_time` datetime DEFAULT NULL,
  `update_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


INSERT INTO `envoy_cluster` (`id`, `cluster_name`, `display_name`, `enable`, `is_delete`, `create_time`, `create_user`, `update_time`, `update_user`) VALUES
	(1, 'dev', '开发', 1, 0, '2024-01-11 21:09:54', 'wangpc', '2024-01-11 21:10:00', 'wangpc'),
	(2, 'test', '测试', 1, 0, '2024-01-11 21:10:29', 'wangpc', '2024-01-11 21:10:33', 'wangpc');

CREATE TABLE IF NOT EXISTS `envoy_node` (
  `id` int NOT NULL AUTO_INCREMENT,
  `envoy_cluster_id` int NOT NULL,
  `node_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `socket_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `enable` tinyint NOT NULL DEFAULT '0',
  `describe` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,

  `create_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `update_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `is_delete` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `envoy_node` (`id`, `envoy_cluster_id`, `node_name`, `socket_address`, `enable`, `describe`, `create_time`, `update_time`, `create_user`, `update_user`, `is_delete`) VALUES
	(1, 1, 'dev-1', '192.168.0.199', 1, '开发环境节点1', '2024-01-11 21:13:49', '2024-01-11 21:13:51', 'wangpc', 'wangpc', 0);
