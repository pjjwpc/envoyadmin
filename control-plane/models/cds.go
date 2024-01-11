package models

type Cds struct {
	Id             int32  `gorm:"id" json:"id"`
	EnvoyClusterId int64  `gorm:"envoy_cluster_id" json:"envoyClusterId"`
	Name           string `gorm:"name" json:"name"`
	ValueData      string `gorm:"value_data" json:"valueData"`
	Version        string `gorm:"version" json:"version"`
}

func (Cds) TableName() string {
	return "cds"
}
