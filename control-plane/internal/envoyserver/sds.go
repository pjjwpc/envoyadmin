package envoyserver

import (
	"control-plane/internal/db"
	"control-plane/internal/models"
	secret "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
)

var (
	map_sds             = map[int64]map[string][]types.Resource{}
	current_version_sds = map[int64]string{}
)

func LoadSDS(envoyClusterId int64) {
	sdsList := []models.Sds{}
	db.Orm.Model(&models.Sds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&sdsList)

	realSdsList := []types.Resource{}
	var version string

	for _, sds := range sdsList {
		realSds := secret.Secret{}
		err := protojson.Unmarshal([]byte(sds.ValueData), &realSds)
		if err == nil {
			realSdsList = append(realSdsList, &realSds)
		} else {
			log.Println(err)
		}
		if version == "" {
			version = sds.Version
		}
	}

	if len(sdsList) > 0 && version != "" {
		if _, ok := map_sds[envoyClusterId]; !ok {
			map_sds[envoyClusterId] = make(map[string][]types.Resource)
		}
		map_sds[envoyClusterId][version] = realSdsList
		current_version_sds[envoyClusterId] = version
	}
}
