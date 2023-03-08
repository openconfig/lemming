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

package gnmi

import (
	"context"
	"io"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// Mode indicates the mode in which the gNMI service operates.
type Mode string

const (
	// GNMIModeMetadataKey is the context metadata key used to specify the
	// mode in which the gNMI server should operate.
	GNMIModeMetadataKey = "gnmi-mode"
	// ConfigMode indicates that the gNMI service will allow updates to
	// intended configuration, but not operational state values.
	ConfigMode Mode = "config"
	// StateMode indicates that the gNMI service will allow updates to
	// operational state, but not intended configuration values.
	StateMode Mode = "state"

	// TimestampMetadataKey is the context metadata key used to specify a
	// custom timestamp for the values in the SetRequest instead of using
	// the time at which the SetRequest is received by the server.
	TimestampMetadataKey = "gnmi-timestamp"
)

// appendToIncomingContext returns a new context with the provided kv merged
// with any existing metadata in the context. Please refer to the documentation
// of Pairs for a description of kv.
func appendToIncomingContext(ctx context.Context, kv ...string) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	return metadata.NewIncomingContext(ctx, metadata.Join(md, metadata.Pairs(kv...)))
}

// AddTimestampMetadata adds a gNMI timestamp metadata to the context.
//
// - ctx is the context to be used for accessing lemming's internal datastore.
// - timestamp is the number of nanoseconds since Epoch.
//
// NOTE: The output of this function should only be used to call into the
// internal lemming gNMI server. This is because it adds an incoming rather
// than an outgoing context metadata to skip regular protobuf handling.
func AddTimestampMetadata(ctx context.Context, timestamp int64) context.Context {
	return appendToIncomingContext(ctx, TimestampMetadataKey, strconv.FormatInt(timestamp, 10))
}

// new creates a state-based gNMI client for the gNMI cache.
// The client calls the server gRPC implementation with a custom streaming gRPC implementation
// in order to bypass the regular gRPC wire marshalling/unmarshalling handling.
func newLocalClient(srv gpb.GNMIServer) gpb.GNMIClient {
	return &localClient{
		gnmiMode: StateMode,
		srv:      srv,
	}
}

// localClient is a gNMI client talks directly to a server, without sending messages over the wire.
type localClient struct {
	gpb.GNMIClient
	gnmiMode Mode
	srv      gpb.GNMIServer
}

// Set uses the datastore client for Set, instead of the public cache endpoint.
func (c *localClient) Set(ctx context.Context, in *gpb.SetRequest, _ ...grpc.CallOption) (*gpb.SetResponse, error) {
	return c.srv.Set(appendToIncomingContext(ctx, GNMIModeMetadataKey, string(c.gnmiMode)), in)
}

// Subscribe implements gNMI Subscribe, by calling a gNMI server directly.
func (c *localClient) Subscribe(ctx context.Context, _ ...grpc.CallOption) (gpb.GNMI_SubscribeClient, error) {
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
		err := c.srv.Subscribe(sub)
		if err != nil {
			errCh <- err
		}
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
