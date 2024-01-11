package models

type EnvoyNode struct {
	Id             int64  `gorm:"id" json:"id"`
	NodeName       string `gorm:"node_name" json:"nodeName"`
	EnvoyClusterId int64  `gorm:"envoy_cluster_id" json:"envoyClusterId"`
}
