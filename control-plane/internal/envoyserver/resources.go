package envoyserver

import (
	"context"
	db "control-plane/internal/db"
	"control-plane/internal/models"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	secret "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	rateLimit "github.com/envoyproxy/go-control-plane/ratelimit/config/ratelimit/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"os"
)

var (
	map_cds  map[int64]map[string][]types.Resource
	map_eds  map[int64]map[string][]types.Resource
	map_lds  map[int64]map[string][]types.Resource
	map_rds  map[int64]map[string][]types.Resource
	map_vhds map[int64]map[string][]types.Resource
	map_sds  map[int64]map[string][]types.Resource
	map_rls  map[int64]map[string][]types.Resource
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
	map_lds = map[int64]map[string][]types.Resource{}
	map_rds = map[int64]map[string][]types.Resource{}
	map_vhds = map[int64]map[string][]types.Resource{}
	map_sds = map[int64]map[string][]types.Resource{}
	map_rls = map[int64]map[string][]types.Resource{}
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
	for version, ldsList := range map_lds[envoyClusterId] {
		cc.Resources[types.Listener] = cache.NewResources(version, ldsList)
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
		initLds(envoyCluster.Id)
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

func initLds(envoyClusterId int64) {
	ldsList := []models.Lds{}
	db.Orm.Model(&models.Lds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&ldsList)
	versionLds := map[string][]types.Resource{}
	realLdsList := []types.Resource{}
	for _, lds := range ldsList {
		realLds := listener.Listener{}
		err := protojson.Unmarshal([]byte(lds.ValueData), &realLds)
		if err == nil {
			realLdsList = append(realLdsList, &realLds)
		} else {
			log.Println(err)
		}
	}
	if len(ldsList) > 0 {
		versionLds[ldsList[0].Version] = realLdsList
		map_lds[envoyClusterId] = versionLds
	}
}

func initRds(envoyClusterId int64) {
	rdsList := []models.Rds{}
	db.Orm.Model(&models.Rds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&rdsList)
	versionRds := map[string][]types.Resource{}
	realLdsList := []types.Resource{}
	for _, rds := range rdsList {
		realRds := route.RouteConfiguration{}
		err := protojson.Unmarshal([]byte(rds.ValueData), &realRds)
		if err == nil {
			realLdsList = append(realLdsList, &realRds)
		} else {
			log.Println(err)
		}
	}
	if len(rdsList) > 0 {
		versionRds[rdsList[0].Version] = realLdsList
		map_rds[envoyClusterId] = versionRds
	}
}
func initVhds(envoyClusterId int64) {
	vhdsList := []models.Vhds{}
	db.Orm.Model(&models.Vhds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&vhdsList)
	versionVhds := map[string][]types.Resource{}
	realVhdsList := []types.Resource{}
	for _, vhds := range vhdsList {
		realVhds := route.VirtualHost{}
		err := protojson.Unmarshal([]byte(vhds.ValueData), &realVhds)
		if err == nil {
			realVhdsList = append(realVhdsList, &realVhds)
		} else {
			log.Println(err)
		}
	}
	if len(vhdsList) > 0 {
		versionVhds[vhdsList[0].Version] = realVhdsList
		map_vhds[envoyClusterId] = versionVhds
	}
}
func initSds(envoyClusterId int64) {
	sdsList := []models.Sds{}
	db.Orm.Model(&models.Sds{}).
		Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).
		Find(&sdsList)
	versionSds := map[string][]types.Resource{}
	realSdsList := []types.Resource{}
	for _, sds := range sdsList {
		realSds := secret.Secret{}
		err := protojson.Unmarshal([]byte(sds.ValueData), &realSds)
		if err == nil {
			realSdsList = append(realSdsList, &realSds)
		} else {
			log.Println(err)
		}
	}
	if len(sdsList) > 0 {
		versionSds[sdsList[0].Version] = realSdsList
		map_sds[envoyClusterId] = versionSds
	}
}
func initRls(envoyClusterId int64) {
	rlsList := []models.Rls{}
	db.Orm.Model(&models.Rls{}).Where("is_delete=0 and enable=1 and envoy_cluster_id=?", envoyClusterId).Find(&rlsList)
	versionRls := map[string][]types.Resource{}
	realRlsList := []types.Resource{}
	for _, rls := range rlsList {
		realRls := rateLimit.RateLimitConfig{}
		err := protojson.Unmarshal([]byte(rls.ValueData), &realRls)
		if err == nil {
			realRlsList = append(realRlsList, &realRls)
		} else {
			log.Println(err)
		}
	}
	if len(rlsList) > 0 {
		versionRls[rlsList[0].Version] = realRlsList
		map_rls[envoyClusterId] = versionRls
	}

}
