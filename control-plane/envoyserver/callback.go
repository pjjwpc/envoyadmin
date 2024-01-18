package envoyserver

import (
	"context"
	"log"
	"sync"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

type Callbacks struct {
	Signal         chan struct{}
	Debug          bool
	Fetches        int
	Requests       int
	DeltaRequests  int
	DeltaResponses int
	mu             sync.Mutex
}

func (cb *Callbacks) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
}
func (cb *Callbacks) OnStreamOpen(c context.Context, id int64, typ string) error {
	log.Println("StreamOpen")
	return nil
}
func (cb *Callbacks) OnStreamClosed(id int64, node *corev3.Node) {

}
func (cb *Callbacks) OnDeltaStreamOpen(_ context.Context, id int64, typ string) error {

	return nil
}
func (cb *Callbacks) OnDeltaStreamClosed(id int64, _ *corev3.Node) {

}
func (cb *Callbacks) OnStreamRequest(id int64, req *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Requests++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	if req.ErrorDetail != nil {
		log.Println(req.ErrorDetail.Message)
	}
	return nil
}

func (cb *Callbacks) OnStreamResponse(
	c context.Context,
	id int64,
	req *discovery.DiscoveryRequest,
	resp *discovery.DiscoveryResponse) {
}

func (cb *Callbacks) OnStreamDeltaResponse(id int64,
	req *discovery.DeltaDiscoveryRequest,
	resp *discovery.DeltaDiscoveryResponse) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.DeltaResponses++
}

func (cb *Callbacks) OnStreamDeltaRequest(id int64, req *discovery.DeltaDiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.DeltaRequests++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	if req.ErrorDetail != nil {
		log.Println(req.ErrorDetail.Message)
	} else {
	}
	if req.ResponseNonce != "" {

	}

	return nil
}

func (cb *Callbacks) OnFetchRequest(_ context.Context, req *discovery.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.Fetches++
	if cb.Signal != nil {
		close(cb.Signal)
		cb.Signal = nil
	}
	if req.ErrorDetail != nil {
		log.Println(req.ErrorDetail.Message)
	}
	if cb.Debug {
	}
	return nil
}
func (cb *Callbacks) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {

}
