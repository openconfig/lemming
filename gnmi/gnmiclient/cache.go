// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gnmiclient

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// NewCacheClient creates a gNMI client for the gNMI cache.
// The client calls the server gRPC implementation with a custom streaming gRPC implementation
// in order to bypass the regular gRPC wire marshalling/unmarshalling handling.
// TODO: refactor gNMI so that a single server can be used here.
func NewCacheClient(srv gpb.GNMIServer, setClient gpb.GNMIClient) gpb.GNMIClient {
	return &cacheClient{
		srv:       srv,
		setClient: setClient,
	}
}

// cacheClient is a gNMI client talks directly to a server, without sending messages over the wire.
type cacheClient struct {
	gpb.GNMIClient
	srv       gpb.GNMIServer
	setClient gpb.GNMIClient
}

// Set uses the datastore client for Set, instead of the public cache endpoint.
func (cc *cacheClient) Set(ctx context.Context, in *gpb.SetRequest, opts ...grpc.CallOption) (*gpb.SetResponse, error) {
	return cc.setClient.Set(ctx, in, opts...)
}

// Subscribe implements gNMI Subscribe, by calling a gNMI server directly.
func (cc *cacheClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (gpb.GNMI_SubscribeClient, error) {
	errCh := make(chan error)
	respCh := make(chan *gpb.SubscribeResponse, 10)
	reqCh := make(chan *gpb.SubscribeRequest)

	sub := &subServer{
		respCh: respCh,
		reqCh:  reqCh,
		ctx:    peer.NewContext(ctx, &peer.Peer{}), // Add empty Peer, since the cache expects to be set.
	}
	client := &subClient{
		errCh:  errCh,
		respCh: respCh,
		reqCh:  reqCh,
	}

	go func() {
		err := cc.srv.Subscribe(sub)
		errCh <- err
	}()
	return client, nil
}

// subClient is an implementation of GNMI_SubscribeClient that use channels to pass messages.
type subClient struct {
	gpb.GNMI_SubscribeClient
	errCh  chan error
	respCh chan *gpb.SubscribeResponse
	reqCh  chan *gpb.SubscribeRequest
}

func (sc *subClient) CloseSend() error {
	close(sc.reqCh)
	return nil
}

func (sc *subClient) Send(req *gpb.SubscribeRequest) error {
	sc.reqCh <- req
	return nil
}

func (sc *subClient) Recv() (*gpb.SubscribeResponse, error) {
	for {
		select {
		case resp := <-sc.respCh:
			return resp, nil
		case err := <-sc.errCh:
			return nil, err
		}
	}
}

// subServer is an implementation of GNMI_SubscribeServer that use channels to pass messages.
type subServer struct {
	gpb.GNMI_SubscribeServer
	respCh chan *gpb.SubscribeResponse
	reqCh  chan *gpb.SubscribeRequest
	ctx    context.Context
}

func (ss *subServer) Context() context.Context {
	return ss.ctx
}

func (ss *subServer) Send(resp *gpb.SubscribeResponse) error {
	ss.respCh <- resp
	return nil
}

func (ss *subServer) Recv() (*gpb.SubscribeRequest, error) {
	req, ok := <-ss.reqCh
	if !ok {
		return nil, io.EOF
	}
	return req, nil
}
