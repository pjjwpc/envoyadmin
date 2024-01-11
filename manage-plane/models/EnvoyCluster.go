package models

type EnvoyClusterModel struct {
	Id          int64  `gorm:"id" json:"id"`
	Version     string `gorm:"version" json:"version"`
	ClusterName string `gorm:"cluster_name" json:"clusterName"`
	DisplayName string `gorm:"display_name" json:"displayName"`
	Enable      bool   `gorm:"enable" json:"enable"`
	Audit
}
type EnvoyClusterViewModel struct {
	Id          int64  `gorm:"id" json:"id"`
	ClusterName string `gorm:"cluster_name" json:"clusterName"`
	Version     string `gorm:"version" json:"version"`
}

func (EnvoyClusterViewModel) TableName() string {
	return "envoy_cluster"
}

func (EnvoyClusterModel) TableName() string {
	return "envoy_cluster"
}
