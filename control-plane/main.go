package main

import (
	"context"
	cpc "control-plane/config"
	"control-plane/db"
	ecpl "control-plane/envoyserver"
	"encoding/json"
	"log"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/protobuf/types/known/durationpb"
)

func makeEndpoint(clusterName string) *endpoint.ClusterLoadAssignment {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,

									Address: "k8s.wangpc",
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: 3306,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}
func main() {
	realCds := cluster.Cluster{
		Name:                 "k8s",
		ConnectTimeout:       durationpb.New(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeEndpoint("k8s"),
		DnsLookupFamily:      cluster.Cluster_V4_ONLY,
	}
	cds, _ := json.Marshal(&realCds)
	eds, _ := json.Marshal(realCds.LoadAssignment)
	log.Println(string(cds))
	log.Println(string(eds))

	db.InitDb()
	cache := ecpl.InitCache()
	// Run the xDS server
	ctx := context.Background()
	cb := &ecpl.Callbacks{Debug: false}

	srv := server.NewServer(ctx, cache, cb)

	ecpl.RunServer(ctx, srv, cpc.Config.Port)
}
