package envoyserver

import (
	"control-plane/internal/db"
	"control-plane/internal/models"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
)

var (
	map_rds             = map[int64]map[string][]types.Resource{}
	current_version_rds = map[int64]string{}
)

func LoadRDS(envoyClusterId int64) {
	rdsList := []models.Rds{}
	db.Orm.Model(&models.Rds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&rdsList)

	realRdsList := []types.Resource{}
	var version string

	for _, rds := range rdsList {
		realRds := route.RouteConfiguration{}
		err := protojson.Unmarshal([]byte(rds.ValueData), &realRds)
		if err == nil {
			realRdsList = append(realRdsList, &realRds)
		} else {
			log.Println(err)
		}
		if version == "" {
			version = rds.Version
		}
	}

	if len(rdsList) > 0 && version != "" {
		if _, ok := map_rds[envoyClusterId]; !ok {
			map_rds[envoyClusterId] = make(map[string][]types.Resource)
		}
		map_rds[envoyClusterId][version] = realRdsList
		current_version_rds[envoyClusterId] = version
	}
}
