package envoyserver

import (
	"control-plane/internal/db"
	"control-plane/internal/models"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
)

var (
	map_lds             = map[int64]map[string][]types.Resource{}
	current_version_lds = map[int64]string{}
)

func LoadLDS(envoyClusterId int64) {
	ldsList := []models.Lds{}
	db.Orm.Model(&models.Lds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&ldsList)

	realLdsList := []types.Resource{}
	var version string

	for _, lds := range ldsList {
		realLds := listener.Listener{}
		err := protojson.Unmarshal([]byte(lds.ValueData), &realLds)
		if err == nil {
			realLdsList = append(realLdsList, &realLds)
		} else {
			log.Println(err)
		}
		if version == "" {
			version = lds.Version
		}
	}

	if len(ldsList) > 0 && version != "" {
		if _, ok := map_lds[envoyClusterId]; !ok {
			map_lds[envoyClusterId] = make(map[string][]types.Resource)
		}
		map_lds[envoyClusterId][version] = realLdsList
		current_version_lds[envoyClusterId] = version
	}
}
