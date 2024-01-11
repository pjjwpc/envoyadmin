package config

import (
	v32 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	trace "github.com/envoyproxy/go-control-plane/envoy/config/trace/v3"
	access "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/file/v3"
	rls "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ratelimit/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tcp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	hf "github.com/envoyproxy/go-control-plane/envoy/extensions/http/header_formatters/preserve_case/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	httppo "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"
	pj "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var (
	ClusterParser  pj.UnmarshalOptions
	ClusterFormat  pj.MarshalOptions
	ListenerParser pj.UnmarshalOptions
	ListenerFormat pj.MarshalOptions
)

func init() {
	clusterRegType := protoregistry.Types{}
	httpProtocolOptions := httppo.HttpProtocolOptions{}
	headerformatters := hf.PreserveCaseFormatterConfig{}
	healcheck := v32.HealthCheck{}
	zipkin := trace.ZipkinConfig{}
	otle := trace.OpenTelemetryConfig{}
	tcpproxy := tcp.TcpProxy{}
	clusterRegType.RegisterMessage(httpProtocolOptions.ProtoReflect().Type())
	clusterRegType.RegisterMessage(healcheck.ProtoReflect().Type())
	clusterRegType.RegisterMessage(headerformatters.ProtoReflect().Type())
	ClusterParser = pj.UnmarshalOptions{Resolver: &clusterRegType}
	ClusterFormat = pj.MarshalOptions{Resolver: &clusterRegType}
	//listener反序列化对象注册扩展类型
	listenerRegType := protoregistry.Types{}
	httpConnectionManagerType := hcm.HttpConnectionManager{}
	routerType := router.Router{}
	fileAccesslogType := access.FileAccessLog{}
	httpdown := tls.DownstreamTlsContext{}
	httpup := tls.UpstreamTlsContext{}
	rlstype := rls.RateLimit{}
	listenerRegType.RegisterMessage(httpConnectionManagerType.ProtoReflect().Type())
	listenerRegType.RegisterMessage(routerType.ProtoReflect().Type())
	listenerRegType.RegisterMessage(fileAccesslogType.ProtoReflect().Type())
	listenerRegType.RegisterMessage(httpdown.ProtoReflect().Type())
	listenerRegType.RegisterMessage(httpup.ProtoReflect().Type())
	listenerRegType.RegisterMessage(zipkin.ProtoReflect().Type())
	listenerRegType.RegisterMessage(otle.ProtoReflect().Type())
	listenerRegType.RegisterMessage(tcpproxy.ProtoReflect().Type())
	listenerRegType.RegisterMessage(headerformatters.ProtoReflect().Type())
	listenerRegType.RegisterMessage(rlstype.ProtoReflect().Type())
	ListenerParser = pj.UnmarshalOptions{Resolver: &listenerRegType}
	ListenerFormat = pj.MarshalOptions{Resolver: &listenerRegType}
}
