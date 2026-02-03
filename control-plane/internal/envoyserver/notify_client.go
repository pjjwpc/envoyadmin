package envoyserver

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

type JSONCodec struct{}

func (JSONCodec) Name() string {
	return "json"
}

func (JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func init() {
	encoding.RegisterCodec(JSONCodec{})
}

type WatchRequest struct {
	NodeId string `json:"nodeId,omitempty"`
}

type ChangeEvent struct {
	EventType      string `json:"eventType"`
	EnvoyClusterId int64  `json:"envoyClusterId"`
	ResourceType   string `json:"resourceType"`
	Version        string `json:"version"`
}

type NotifyServiceClient interface {
	WatchChanges(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (NotifyService_WatchChangesClient, error)
}

type notifyServiceClient struct {
	cc *grpc.ClientConn
}

type NotifyService_WatchChangesClient interface {
	Recv() (*ChangeEvent, error)
	grpc.ClientStream
}

type notifyServiceWatchChangesClient struct {
	grpc.ClientStream
}

func (x *notifyServiceWatchChangesClient) Recv() (*ChangeEvent, error) {
	m := new(ChangeEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func NewNotifyServiceClient(cc *grpc.ClientConn) NotifyServiceClient {
	return &notifyServiceClient{cc}
}

var _NotifyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "notify.NotifyService",
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchChanges",
			ServerStreams: true,
		},
	},
	Methods: []grpc.MethodDesc{},
}

func (c *notifyServiceClient) WatchChanges(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (NotifyService_WatchChangesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_NotifyService_serviceDesc.Streams[0], "/notify.NotifyService/WatchChanges", opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := stream.CloseSend(); err != nil {
		return nil, err
	}
	return &notifyServiceWatchChangesClient{stream}, nil
}

func handleChangeEvent(ev *ChangeEvent) {
	switch ev.ResourceType {
	case "CDS":
		LoadCDS(ev.EnvoyClusterId)
	case "EDS":
		LoadEDS(ev.EnvoyClusterId)
	case "LDS":
		LoadLDS(ev.EnvoyClusterId)
	case "RDS":
		LoadRDS(ev.EnvoyClusterId)
	case "VHDS":
		LoadVHDS(ev.EnvoyClusterId)
	case "SDS":
		LoadSDS(ev.EnvoyClusterId)
	case "RLS":
		LoadRLS(ev.EnvoyClusterId)
	default:
	}
	refreshClusterSnapshots(ev.EnvoyClusterId)
	saveClusterBackup(ev.EnvoyClusterId)
}

func StartNotifyStream(ctx context.Context) {
	addr := os.Getenv("MANAGE_PLANE_NOTIFY_ADDR")
	if addr == "" {
		addr = "manage-plane:8091"
	}
	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(JSONCodec{}.Name())),
	)
	if err != nil {
		log.Println("notify grpc dial error", err)
		return
	}
	client := NewNotifyServiceClient(conn)
	go func() {
		defer conn.Close()
		for {
			stream, err := client.WatchChanges(ctx, &WatchRequest{})
			if err != nil {
				log.Println("notify watch error", err)
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					continue
				}
			}
			for {
				ev, err := stream.Recv()
				if err != nil {
					log.Println("notify stream recv error", err)
					break
				}
				handleChangeEvent(ev)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
	}()
}
