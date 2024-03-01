package models

type EnvoyCluster struct {
	Id          int64  `gorm:"id" json:"id"`
	ClusterName string `gorm:"cluster_name" json:"clusterName"`
}

func (EnvoyCluster) TableName() string {
	return "envoy_cluster"
}
