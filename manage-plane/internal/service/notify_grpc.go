package service

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"sync"

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

type NotifyServiceServer interface {
	WatchChanges(*WatchRequest, NotifyService_WatchChangesServer) error
}

type NotifyService_WatchChangesServer interface {
	Send(*ChangeEvent) error
	grpc.ServerStream
}

type notifyService struct {
	mu          sync.Mutex
	subscribers map[*subscriber]struct{}
}

type subscriber struct {
	ch chan *ChangeEvent
}

func newNotifyService() *notifyService {
	return &notifyService{
		subscribers: make(map[*subscriber]struct{}),
	}
}

var defaultNotifyService = newNotifyService()

func (s *notifyService) addSubscriber(sub *subscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribers[sub] = struct{}{}
}

func (s *notifyService) removeSubscriber(sub *subscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subscribers, sub)
	close(sub.ch)
}

func (s *notifyService) broadcast(ev *ChangeEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for sub := range s.subscribers {
		select {
		case sub.ch <- ev:
		default:
		}
	}
}

func (s *notifyService) WatchChanges(req *WatchRequest, stream NotifyService_WatchChangesServer) error {
	sub := &subscriber{ch: make(chan *ChangeEvent, 16)}
	s.addSubscriber(sub)
	defer s.removeSubscriber(sub)
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case ev := <-sub.ch:
			if err := stream.Send(ev); err != nil {
				return err
			}
		}
	}
}

func PublishChange(ev *ChangeEvent) {
	defaultNotifyService.broadcast(ev)
}

func RegisterNotifyServiceServer(s *grpc.Server, srv NotifyServiceServer) {
	s.RegisterService(&_NotifyService_serviceDesc, srv)
}

var _NotifyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "notify.NotifyService",
	HandlerType: (*NotifyServiceServer)(nil),
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchChanges",
			Handler:       _NotifyService_WatchChanges_Handler,
			ServerStreams: true,
		},
	},
	Methods: []grpc.MethodDesc{},
}

func _NotifyService_WatchChanges_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NotifyServiceServer).WatchChanges(m, &notifyServiceWatchChangesServer{stream})
}

type notifyServiceWatchChangesServer struct {
	grpc.ServerStream
}

func (x *notifyServiceWatchChangesServer) Send(m *ChangeEvent) error {
	return x.ServerStream.SendMsg(m)
}

func StartNotifyServer() {
	addr := os.Getenv("MANAGE_PLANE_NOTIFY_ADDR")
	if addr == "" {
		addr = ":8091"
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("notify grpc listen error", err)
		return
	}
	server := grpc.NewServer()
	RegisterNotifyServiceServer(server, defaultNotifyService)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Println("notify grpc serve error", err)
		}
	}()
}
