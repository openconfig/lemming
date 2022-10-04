// Copyright 2021 Google LLC
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

// Package gnmit is a single-target gNMI collector implementation that can be
// used as an on-device/fake device implementation. It supports the Subscribe RPC
// using the libraries from openconfig/gnmi.
package gnmit

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/cache"
	"github.com/openconfig/gnmi/subscribe"
	"github.com/openconfig/lemming/gnmi/gnmistore"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	OpenconfigOrigin = "openconfig"
)

var (
	// metadataUpdatePeriod is the period of time after which the metadata for the collector
	// is updated to the client.
	metadataUpdatePeriod = 30 * time.Second
	// sizeUpdatePeriod is the period of time after which the storage size information for
	// the collector is updated to the client.
	sizeUpdatePeriod = 30 * time.Second
)

// periodic runs the function fn every period.
func periodic(period time.Duration, fn func()) {
	if period == 0 {
		return
	}
	t := time.NewTicker(period)
	defer t.Stop()
	for range t.C {
		fn()
	}
}

// GNMIServer implements the gNMI server interface.
type GNMIServer struct {
	// The subscribe Server implements only Subscribe for gNMI.
	*subscribe.Server
	c *Collector

	stateMu     sync.Mutex
	stateSchema *ytypes.Schema
}

// NewServer returns a new collector server implementation that can be registered on
// an existing gRPC server.
//
// - schema is the specification of the schema if gnmi.Set is used.
// - hostname is the name of the target.
// - sendMeta indicates whether metadata should be sent
func NewServer(ctx context.Context, schema *ytypes.Schema, hostname string, sendMeta bool) (*Collector, *GNMIServer, error) {
	c := &Collector{
		inCh:   make(chan *gpb.SubscribeResponse),
		name:   hostname,
		schema: schema,
	}

	log.V(1).Infof("Starting cache target: %v", hostname)
	c.cache = cache.New([]string{hostname})
	t := c.cache.GetTarget(hostname)

	// Initialize the cache with the input schema root.
	if schema != nil {
		notifs, err := ygot.TogNMINotifications(schema.Root, time.Now().UnixNano(), ygot.GNMINotificationsConfig{
			UsePathElem: true,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("gnmit: %v", err)
		}
		for _, notif := range notifs {
			if notif.Prefix == nil {
				notif.Prefix = &gpb.Path{}
			}
			notif.Prefix.Origin = OpenconfigOrigin
			notif.Prefix.Target = hostname
			c.cache.GnmiUpdate(notif)
		}
	}

	if sendMeta {
		go periodic(metadataUpdatePeriod, c.cache.UpdateMetadata)
		go periodic(sizeUpdatePeriod, c.cache.UpdateSize)
	}
	t.Connect()

	// start our single collector from the input channel.
	go func() {
		for {
			select {
			case msg := <-c.inCh:
				if err := c.handleUpdate(msg); err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	subscribeSrv, err := subscribe.NewServer(c.cache)
	if err != nil {
		return nil, nil, fmt.Errorf("could not instantiate gNMI server: %v", err)
	}

	stateSchema, err := oc.Schema()
	if err != nil {
		return nil, nil, err
	}
	if err := SetupSchema(stateSchema); err != nil {
		return nil, nil, err
	}

	gnmiserver := &GNMIServer{
		Server:      subscribeSrv, // use the 'subscribe' implementation.
		c:           c,
		stateSchema: stateSchema,
	}

	c.cache.SetClient(subscribeSrv.Update)

	return c, gnmiserver, nil
}

// New returns a new collector that listens on the specified addr (in the form host:port),
// supporting a single downstream target named hostname. sendMeta controls whether the
// metadata *other* than meta/sync and meta/connected is sent by the collector.
//
// New returns the new collector, the address it is listening on in the form hostname:port
// or any errors encounted whilst setting it up.
func New(ctx context.Context, schema *ytypes.Schema, addr, hostname string, sendMeta bool, opts ...grpc.ServerOption) (*Collector, string, error) {
	c, gnmiserver, err := NewServer(ctx, schema, hostname, sendMeta)
	if err != nil {
		return nil, "", err
	}

	// Start gNMI server.
	srv := grpc.NewServer(opts...)
	gpb.RegisterGNMIServer(srv, gnmiserver)
	// Forward streaming updates to clients.
	// Register listening port and start serving.
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to listen: %v", err)
	}

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Error while serving gnmi target: %v", err)
		}
	}()
	c.stopFn = func() {
		srv.GracefulStop()
	}
	return c, lis.Addr().String(), nil
}

type populateDefaultser interface {
	PopulateDefaults()
}

// SetupSchema takes in a ygot schema object which it assumes to be
// uninitialized. It initializes and validates it, returning any errors
// encountered.
func SetupSchema(schema *ytypes.Schema) error {
	if !schema.IsValid() {
		return fmt.Errorf("cannot obtain valid schema for GoStructs: %v", schema)
	}

	// Initialize the root with default values.
	schema.Root.(populateDefaultser).PopulateDefaults()
	if err := schema.Validate(); err != nil {
		return fmt.Errorf("default root of input schema fails validation: %v", err)
	}

	return nil
}

// NewSettable returns a new collector that listens on the specified addr (in the form host:port),
// supporting a single downstream target named hostname. sendMeta controls whether the
// metadata *other* than meta/sync and meta/connected is sent by the collector.
//
// New returns the new collector, the address it is listening on in the form hostname:port
// or any errors encounted whilst setting it up.
//
// NewSettable is different from New in that the returned collector is
// schema-aware and supports gNMI Set. Currently it is not possible to change
// the schema of a Collector after it is created.
func NewSettable(ctx context.Context, addr string, hostname string, sendMeta bool, schema *ytypes.Schema, opts ...grpc.ServerOption) (*Collector, string, error) {
	if err := SetupSchema(schema); err != nil {
		return nil, "", err
	}
	collector, addr, err := New(ctx, schema, addr, hostname, sendMeta, opts...)
	if err != nil {
		return nil, "", err
	}
	return collector, addr, nil
}

// Stop halts the running collector.
func (c *Collector) Stop() {
	c.stopFn()
}

// handleUpdate handles an input gNMI SubscribeResponse that is received by
// the target.
func (c *Collector) handleUpdate(resp *gpb.SubscribeResponse) error {
	t := c.cache.GetTarget(c.name)
	switch v := resp.Response.(type) {
	case *gpb.SubscribeResponse_Update:
		return t.GnmiUpdate(v.Update)
	case *gpb.SubscribeResponse_SyncResponse:
		t.Sync()
	case *gpb.SubscribeResponse_Error:
		return fmt.Errorf("error in response: %s", v)
	default:
		return fmt.Errorf("unknown response %T: %s", v, v)
	}
	return nil
}

// Collector is a basic gNMI target that supports only the Subscribe
// RPC, and acts as a cache for exactly one target.
type Collector struct {
	cache *cache.Cache

	smu    sync.Mutex
	schema *ytypes.Schema

	// name is the hostname of the client.
	name string
	// inCh is a channel use to write new SubscribeResponses to the client.
	inCh chan *gpb.SubscribeResponse
	// stopFn is the function used to stop the server.
	stopFn func()
}

// TargetUpdate provides an input gNMI SubscribeResponse to update the
// cache and clients with.
func (c *Collector) TargetUpdate(m *gpb.SubscribeResponse) {
	c.inCh <- m
}

func updateCache(cache *cache.Cache, root ygot.GoStruct, target string, dirtyRoot ygot.GoStruct, origin string, preferShadowPath bool) error {
	n, err := ygot.Diff(root, dirtyRoot, &ygot.DiffPathOpt{PreferShadowPath: preferShadowPath})
	if err != nil {
		return fmt.Errorf("gnmit: error while creating update notification for Set: %v", err)
	}
	n.Timestamp = time.Now().UnixNano()
	n.Prefix = &gpb.Path{Origin: origin, Target: target}
	if n.Prefix.Origin == "" {
		n.Prefix.Origin = OpenconfigOrigin
	}

	// Update cache
	t := cache.GetTarget(target)
	var pathsForDelete []string
	for _, path := range n.Delete {
		p, err := ygot.PathToString(path)
		if err != nil {
			return fmt.Errorf("cannot convert deleted path to string: %v", err)
		}
		pathsForDelete = append(pathsForDelete, p)
	}
	log.V(1).Infof("datastore: updating the following values: %+v", n.Update)
	log.V(1).Infof("datastore: deleting the following paths: %+v", pathsForDelete)
	if err := t.GnmiUpdate(n); err != nil {
		return err
	}
	return nil
}

// set updates the datastore and intended configuration with the SetRequest,
// allowing read-only values to be updated.
//
// update indicates whether to update the cache with the values from the set
// request.
func set(schema *ytypes.Schema, cache *cache.Cache, target string, req *gpb.SetRequest, preferShadowPath bool) error {
	prevRoot, err := ygot.DeepCopy(schema.Root)
	if err != nil {
		return fmt.Errorf("gnmit: failed to ygot.DeepCopy the cached root object: %v", err)
	}

	success := false

	// Rollback function
	defer func() {
		if !success {
			log.V(1).Infof("Rolling back set request: %v", req)
			schema.Root = prevRoot
		}
	}()

	var unmarshalOpts []ytypes.UnmarshalOpt
	if preferShadowPath {
		unmarshalOpts = append(unmarshalOpts, &ytypes.PreferShadowPath{})
	}
	if err := ytypes.UnmarshalSetRequest(schema, req, unmarshalOpts...); err != nil {
		return fmt.Errorf("gnmit: %v", err)
	}

	if err := schema.Validate(); err != nil {
		return fmt.Errorf("gnmit: invalid SetRequest: %v", err)
	}

	success = true

	if err := updateCache(cache, prevRoot, target, schema.Root, req.Prefix.Origin, preferShadowPath); err != nil {
		return err
	}
	return nil
}

// Set is a prototype for a gNMI Set operation.
// TODO(wenbli): Add unit test.
func (s *GNMIServer) Set(ctx context.Context, req *gpb.SetRequest) (*gpb.SetResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "incoming gNMI SetRequest must specify GNMIMode via the user agent string.")
	}

	// Use ConfigMode by default so that external users don't need to set metadata.
	gnmiMode := gnmistore.ConfigMode
	switch {
	case slices.Contains(md.Get(gnmistore.GNMIModeMetadataKey), string(gnmistore.ConfigMode)):
		gnmiMode = gnmistore.ConfigMode
	case slices.Contains(md.Get(gnmistore.GNMIModeMetadataKey), string(gnmistore.StateMode)):
		gnmiMode = gnmistore.StateMode
	}

	switch gnmiMode {
	case gnmistore.ConfigMode:
		s.c.smu.Lock()
		defer s.c.smu.Unlock()

		log.V(1).Infof("config datastore service received SetRequest: %v", req)
		if s.c.schema == nil {
			return s.UnimplementedGNMIServer.Set(ctx, req)
		}

		// TODO(wenbli): Reject paths that try to modify read-only values.
		// TODO(wenbli): Question: what to do if there are operational-state values in a container that is specified to be replaced or deleted?
		err := set(s.c.schema, s.c.cache, s.c.name, req, true)

		// TODO(wenbli): Currently the SetResponse is not filled.
		return &gpb.SetResponse{
			Timestamp: time.Now().UnixNano(),
		}, err
	case gnmistore.StateMode:
		s.stateMu.Lock()
		defer s.stateMu.Unlock()

		log.V(1).Infof("operational state datastore service received SetRequest: %v", req)
		if s.stateSchema == nil {
			return s.UnimplementedGNMIServer.Set(ctx, req)
		}
		// TODO(wenbli): Reject values that modify config values. We only allow modifying state in this mode.
		if err := set(s.stateSchema, s.c.cache, s.c.name, req, false); err != nil {
			return &gpb.SetResponse{}, status.Errorf(codes.Aborted, "%v", err)
		}

		// This mode is intended to be used internally, and the SetResponse doesn't matter.
		return &gpb.SetResponse{}, nil
	default:
		return nil, status.Errorf(codes.Internal, "incoming gNMI SetRequest must specify a valid GNMIMode via context metadata: %v", md)
	}
}
