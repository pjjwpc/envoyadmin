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
	map_vhds             = map[int64]map[string][]types.Resource{}
	current_version_vhds = map[int64]string{}
)

func LoadVHDS(envoyClusterId int64) {
	vhdsList := []models.Vhds{}
	db.Orm.Model(&models.Vhds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&vhdsList)

	realVhdsList := []types.Resource{}
	var version string

	for _, vhds := range vhdsList {
		realVhds := route.VirtualHost{}
		err := protojson.Unmarshal([]byte(vhds.ValueData), &realVhds)
		if err == nil {
			realVhdsList = append(realVhdsList, &realVhds)
		} else {
			log.Println(err)
		}
		if version == "" {
			version = vhds.Version
		}
	}

	if len(vhdsList) > 0 && version != "" {
		if _, ok := map_vhds[envoyClusterId]; !ok {
			map_vhds[envoyClusterId] = make(map[string][]types.Resource)
		}
		map_vhds[envoyClusterId][version] = realVhdsList
		current_version_vhds[envoyClusterId] = version
	}
}
