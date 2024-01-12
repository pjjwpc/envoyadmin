package query

import (
	"manage-plane/models"
)

type CdsQuery struct {
	models.PageInfo
	Name           string `json:"name" form:"name"`
	EnvoyClusterId int64  `json:"envoyClusterId" form:"envoyClusterId"`
}
