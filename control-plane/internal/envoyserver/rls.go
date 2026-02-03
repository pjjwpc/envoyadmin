package envoyserver

import (
	"control-plane/internal/db"
	"control-plane/internal/models"
	rateLimit "github.com/envoyproxy/go-control-plane/ratelimit/config/ratelimit/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
)

var (
	map_rls             = map[int64]map[string][]types.Resource{}
	current_version_rls = map[int64]string{}
)

func LoadRLS(envoyClusterId int64) {
	rlsList := []models.Rls{}
	db.Orm.Model(&models.Rls{}).Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).Find(&rlsList)

	realRlsList := []types.Resource{}
	var version string

	for _, rls := range rlsList {
		realRls := rateLimit.RateLimitConfig{}
		err := protojson.Unmarshal([]byte(rls.ValueData), &realRls)
		if err == nil {
			realRlsList = append(realRlsList, &realRls)
		} else {
			log.Println(err)
		}
		if version == "" {
			version = rls.Version
		}
	}

	if len(rlsList) > 0 && version != "" {
		if _, ok := map_rls[envoyClusterId]; !ok {
			map_rls[envoyClusterId] = make(map[string][]types.Resource)
		}
		map_rls[envoyClusterId][version] = realRlsList
		current_version_rls[envoyClusterId] = version
	}
}
