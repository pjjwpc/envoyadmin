package models

type Cds struct {
	Id              int32  `gorm:"id" json:"id"`
	EnvoyClusterId  int64  `gorm:"envoy_cluster_id" json:"envoyClusterId"`
	Name            string `gorm:"name" json:"name"`
	ValueData       string `gorm:"value_data" json:"valueData"`
	Version         string `gorm:"version" json:"version"`
	Type            string `gorm:"type" json:"type"`
	HealthCheck     bool   `gorm:"health_check" json:"healthCheck"`
	DnsLookupFamily int    `gorm:"dns_lookup_family" json:"dnsLookupFamily"`
	LbPolicy        int    `gorm:"lb_policy" json:"lbPolicy"`
	Enable          bool   `gorm:"enable" json:"enable"`
	ErrMsg          string `gorm:"err_msg" json:"errMsg"`
	ErrCode         int32  `gorm:"err_code" json:"errCode"`
	Audit
	EnvoyCluster EnvoyClusterViewModel `gorm:"foreignKey:id;references:envoy_cluster_id" json:"envoyCluster"`
}

func (Cds) TableName() string {
	return "cds"
}
