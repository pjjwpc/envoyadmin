package envoyserver

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	// ecds "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	// listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	// routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	// runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	// secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	// rls "github.com/envoyproxy/go-control-plane/ratelimit/service/ratelimit/v3"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
)

var (
	grpcServer *grpc.Server
)

func registerServer(grpcServer *grpc.Server, server server.Server) {
	// register services
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	// routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	// listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	// secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, server)
	// runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)
	// routeservice.RegisterVirtualHostDiscoveryServiceServer(grpcServer, server)
	// rls.RegisterRateLimitConfigDiscoveryServiceServer(grpcServer, server)
	// ecds.RegisterExtensionConfigDiscoveryServiceServer(grpcServer, server)
}

// RunServer starts an xDS server at the given port.
func RunServer(_ context.Context, srv server.Server, port uint) error {
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	)
	grpcServer = grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", ":8060")
	if err != nil {
		log.Fatal(err)
	}

	registerServer(grpcServer, srv)

	log.Printf("management server listening on %d\n", 8060)
	return grpcServer.Serve(lis)
}
