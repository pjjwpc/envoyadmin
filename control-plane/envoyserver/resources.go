package envoyserver

import (
	"context"
	db "control-plane/db"
	"control-plane/models"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"os"
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
func InitCache() cache.SnapshotCache {
	l = Logger{}
	l.Debug = true
	cachei = cache.NewSnapshotCache(false, cache.IDHash{}, l)
	initData() //加载数据到内存
	envoynodeList := []models.EnvoyNode{}
	db.Orm.Model(&models.EnvoyNode{}).
		Where("1=1 and enable = 1 and is_delete = 0").
		Find(&envoynodeList)
	for _, node := range envoynodeList {
		snap, _ := setSnap(node.EnvoyClusterId)
		if err := cachei.SetSnapshot(context.Background(), node.NodeName, snap); err != nil {
			l.Errorf("snapshot error %q for %+v", err, snap)
			os.Exit(1)
		}
	}

	return cachei
}

func setSnap(envoyClusterId int64) (*cache.Snapshot, error) {
	cc := cache.Snapshot{}

	for version, cdsList := range map_cds[envoyClusterId] {
		cc.Resources[types.Cluster] = cache.NewResources(version, cdsList)
	}
	for version, edsList := range map_eds[envoyClusterId] {
		cc.Resources[types.Endpoint] = cache.NewResources(version, edsList)
	}

	return &cc, nil
}
func initData() {
	envoyClusterList := []models.EnvoyCluster{}
	db.Orm.Model(&models.EnvoyCluster{}).
		Where("is_delete=0 and enable=1").
		Find(&envoyClusterList)
	//遍历集群查找cds、eds
	for _, envoyCluster := range envoyClusterList {
		initCluster(envoyCluster.Id)
		initEndPoint(envoyCluster.Id)
	}
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
