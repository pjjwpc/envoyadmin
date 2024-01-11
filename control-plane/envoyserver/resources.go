package envoyserver

import (
	db "control-plane/db"
	"control-plane/models"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
)

var (
	map_cds map[int64]map[string][]types.Resource
	map_eds map[int64]map[string][]types.Resource
)

var (
	l      Logger
	port   uint
	nodeID string
	cachei cache.SnapshotCache
)

func init() {

	map_cds = map[int64]map[string][]types.Resource{}
	map_eds = map[int64]map[string][]types.Resource{}
}
func initCluster(envoyClusterId int64) {
	cdsList := []models.Cds{}
	db.Orm.Model(&models.Cds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&cdsList)
	versionCds := map[string][]types.Resource{}
	realCdsList := []types.Resource{}
	for _, cds := range cdsList {
		realCds := cluster.Cluster{}

		err := ClusterParser.Unmarshal([]byte(cds.ValueData), &realCds)
		if err == nil {
			realCdsList = append(realCdsList, &realCds)
		} else {
			log.Println(err)
		}
	}
	if len(cdsList) > 0 {
		versionCds[cdsList[0].Version] = realCdsList
		map_cds[envoyClusterId] = versionCds
	}
}

func initEndPoint(envoyClusterId int64) {
	edsList := []models.Eds{}
	db.Orm.Model(&models.Eds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&edsList)
	versionEds := map[string][]types.Resource{}
	realEdsList := []types.Resource{}
	for _, eds := range edsList {
		realEds := endpoint.ClusterLoadAssignment{}
		err := protojson.Unmarshal([]byte(eds.ValueData), &realEds)
		if err == nil {
			realEdsList = append(realEdsList, &realEds)
		} else {
			log.Println(err)
		}
	}
	if len(edsList) > 0 {
		versionEds[edsList[0].Version] = realEdsList
		map_eds[envoyClusterId] = versionEds
	}
}
