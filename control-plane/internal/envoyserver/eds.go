package envoyserver

import (
	"control-plane/internal/db"
	"control-plane/internal/models"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
)

var (
	map_eds             = map[int64]map[string][]types.Resource{}
	current_version_eds = map[int64]string{}
)

func LoadEDS(envoyClusterId int64) {
	edsList := []models.Eds{}
	db.Orm.Model(&models.Eds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&edsList)

	realEdsList := []types.Resource{}
	var version string

	for _, eds := range edsList {
		realEds := endpoint.ClusterLoadAssignment{}
		err := protojson.Unmarshal([]byte(eds.ValueData), &realEds)
		if err == nil {
			realEdsList = append(realEdsList, &realEds)
		} else {
			log.Println(err)
		}
		if version == "" {
			version = eds.Version
		}
	}

	if len(edsList) > 0 && version != "" {
		if _, ok := map_eds[envoyClusterId]; !ok {
			map_eds[envoyClusterId] = make(map[string][]types.Resource)
		}
		map_eds[envoyClusterId][version] = realEdsList
		current_version_eds[envoyClusterId] = version
	}
}
