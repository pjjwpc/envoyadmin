package envoyserver

import (
	"control-plane/internal/db"
	"control-plane/internal/models"
	"log"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

var (
	map_cds             = map[int64]map[string][]types.Resource{}
	current_version_cds = map[int64]string{}
)

func LoadCDS(envoyClusterId int64) {
	cdsList := []models.Cds{}
	db.Orm.Model(&models.Cds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&cdsList)

	realCdsList := []types.Resource{}
	var version string

	for _, cds := range cdsList {
		realCds := cluster.Cluster{}
		err := ClusterParser.Unmarshal([]byte(cds.ValueData), &realCds)
		if err == nil {
			realCdsList = append(realCdsList, &realCds)
		} else {
			log.Println(err)
		}
		// 假设同一批加载的数据版本一致，取第一个非空的版本号
		if version == "" {
			version = cds.Version
		}
	}

	if len(cdsList) > 0 && version != "" {
		// 初始化 map
		if _, ok := map_cds[envoyClusterId]; !ok {
			map_cds[envoyClusterId] = make(map[string][]types.Resource)
		}
		// 保存新版本（保留旧版本）
		map_cds[envoyClusterId][version] = realCdsList
		// 更新当前版本指针
		current_version_cds[envoyClusterId] = version
	}
}
