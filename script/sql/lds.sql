
CREATE TABLE IF NOT EXISTS `lds` (

  `id` int NOT NULL AUTO_INCREMENT,
  `envoy_cluster_id` int NOT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `value_data` json NOT NULL,
  `version` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `protocol` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `listener_protocol` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `port` int NOT NULL DEFAULT '0',
  `rds` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0',
  `enable` tinyint NOT NULL DEFAULT '0',
  `is_delete` tinyint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL,
  `create_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `update_time` datetime DEFAULT NULL,
  `update_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `err_msg` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `err_code` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `lds` (`id`, `envoy_cluster_id`, `name`, `value_data`, `version`, `protocol`, `listener_protocol`, `port`, `rds`, `enable`, `is_delete`, `create_time`, `create_user`, `update_time`, `update_user`, `err_msg`, `err_code`) VALUES
	(1, 1, 'k8s_proxy', '{"name": "k8s-proxy", "address": {"socket_address": {"address": "0.0.0.0", "port_value": 80}}, "filter_chains": [{"filters": [{"name": "envoy.filters.network.http_connection_manager", "typed_config": {"@type": "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager", "stat_prefix": "ingress_http", "http_filters": [{"name": "envoy.filters.http.router", "typed_config": {"@type": "type.googleapis.com/envoy.extensions.filters.http.router.v3.Router"}}], "route_config": {"name": "local_route", "virtual_hosts": [{"name": "harbor", "routes": [{"match": {"prefix": "/"}, "route": {"cluster": "k8s", "timeout": "1200s"}}], "domains": ["harbor.wangpc", "harbor.wangpc:*"]}, {"name": "kuboard", "routes": [{"match": {"prefix": "/"}, "route": {"cluster": "k8s"}}], "domains": ["kuboard.wangpc", "kuboard.wangpc:*"]}, {"name": "local_service", "routes": [{"match": {"prefix": "/", "headers": [{"name": ":method", "exact_match": "HEAD"}]}, "direct_response": {"body": {"inline_string": "heihei"}, "status": 200}}], "domains": ["*"]}]}}}]}]}', '1', 'tcp', 'tcp', 80, '0', 1, 0, '2024-01-18 12:04:10', 'wangpc', '2024-01-18 12:04:15', 'wangpc', NULL, NULL);

