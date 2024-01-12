CREATE DATABASE IF NOT EXISTS `envoy_admin` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `envoy_admin`;

CREATE TABLE IF NOT EXISTS `cds` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `envoy_cluster_id` int NOT NULL DEFAULT '1' COMMENT 'envoy集群',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '集群名称',
  `value_data` json NOT NULL,
  `version` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '端点发现类型',
  `health_check` tinyint NOT NULL DEFAULT '0' COMMENT '是否启用健康检查',
  `dns_lookup_family` tinyint NOT NULL DEFAULT '0' COMMENT 'dns解析',
  `lb_policy` tinyint NOT NULL DEFAULT '0' COMMENT '负载均衡策略',
  `enable` tinyint(1) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL,
  `create_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `update_time` datetime NOT NULL,
  `update_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0',
  `err_msg` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `err_code` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=105 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='上游集群配置信息';

INSERT INTO `cds` (`id`, `envoy_cluster_id`, `name`, `value_data`, `version`, `type`, `health_check`, `dns_lookup_family`, `lb_policy`, `enable`, `create_time`, `create_user`, `update_time`, `update_user`, `is_delete`, `err_msg`, `err_code`) VALUES
	(1, 1, 'k8s', '{"name": "k8s", "type": "STRICT_DNS", "connect_timeout": "5s", "load_assignment": {"endpoints": [{"lb_endpoints": [{"endpoint": {"address": {"socketAddress": {"address": "k8s.wangpc", "port_value": 3306}}}}]}], "cluster_name": "k8s"}}', '1', '2', 0, 0, 0, 1, '2024-01-12 18:48:46', 'wangpc', '2024-01-12 18:48:47', 'wangpc', 0, NULL, NULL);



