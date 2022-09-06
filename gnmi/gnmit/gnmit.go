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
	"encoding/json"
	"fmt"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/cache"
	"github.com/openconfig/lemming/gnmi/internal/oc"
	"github.com/openconfig/lemming/gnmi/subscribe"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	dspb "github.com/openconfig/lemming/proto/datastore"
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

// Queue is an interface that represents a possibly coalescing queue of updates.
type Queue interface {
	Next(ctx context.Context) (interface{}, uint32, error)
	Len() int
	Close()
}

// UpdateFn is a function that takes in a gNMI Notification object and updates
// a gNMI datastore with it.
type UpdateFn func(*gpb.Notification) error

// TaskRoutine is a reactor function that listens for updates from a queue,
// emits updates via an update function. It does this on a target (string
// parameter), and also has a final clean-up function to call when it finishes
// processing.
type TaskRoutine func(func() *oc.Root, Queue, UpdateFn, string, func()) error

// Task defines a particular task that runs on the gNMI datastore.
type Task struct {
	Run    TaskRoutine
	Paths  []ygnmi.PathStruct
	Prefix *gpb.Path
}

// GNMIServer implements the gNMI server interface.
type GNMIServer struct {
	// intendedConfigMu protects the global view of intendedConfig, which
	// can be read by tasks. It protects the storage of the intendedConfig
	// pointer, which gets copied and sent to tasks for reading. Since
	// Set's implementation never alters this view once created, and
	// instead creates a new copy of the root Device struct, tasks can
	// freely read without a race condition.
	intendedConfigMu sync.RWMutex
	// The subscribe Server implements only Subscribe for gNMI.
	*subscribe.Server
	c *Collector
}

// RegisterTask starts up a task on the gNMI datastore.
func (s *GNMIServer) RegisterTask(task Task) error {
	var paths []*gpb.Path
	for _, p := range task.Paths {
		path, _, err := ygnmi.ResolvePath(p)
		if err != nil {
			return fmt.Errorf("gnmit: cannot register task: %v", err)
		}
		paths = append(paths, path)
	}
	queue, remove, err := s.Server.SubscribeLocal(s.c.name, paths, task.Prefix)
	if err != nil {
		return err
	}
	return task.Run(s.getIntendedConfig, queue, s.c.cache.GnmiUpdate, s.c.name, remove)
}

// RegisterTask starts up a task on the gNMI datastore.
func (s *GNMIServer) getIntendedConfig() *oc.Root {
	s.intendedConfigMu.RLock()
	defer s.intendedConfigMu.RUnlock()
	if s.c != nil && s.c.schema != nil {
		return s.c.schema.Root.(*oc.Root)
	}
	return nil
}

// New returns a new collector server implementation that can be registered on
// an existing gRPC server. It takes a string indicating the hostname of the
// target, a boolean indicating whether metadata should be sent, and a slice of
// tasks that are to be launched to run on the server.
func NewServer(ctx context.Context, hostname string, sendMeta bool, tasks []Task) (*Collector, *GNMIServer, error) {
	c := &Collector{
		inCh: make(chan *gpb.SubscribeResponse),
		name: hostname,
	}

	c.cache = cache.New([]string{hostname})
	t := c.cache.GetTarget(hostname)

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

	gnmiserver := &GNMIServer{
		Server: subscribeSrv, // use the 'subscribe' implementation.
		c:      c,
	}

	for _, t := range tasks {
		// TODO(wenbli): We don't current support task re-starts.
		if err := gnmiserver.RegisterTask(t); err != nil {
			return nil, nil, err
		}
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
func New(ctx context.Context, addr, hostname string, sendMeta bool, tasks []Task, opts ...grpc.ServerOption) (*Collector, string, error) {
	c, gnmiserver, err := NewServer(ctx, hostname, sendMeta, tasks)
	if err != nil {
		return nil, "", err
	}

	// Start datastore server.
	if err := os.RemoveAll(DatastoreAddress); err != nil {
		log.Fatal(err)
	}
	srvDS := grpc.NewServer(opts...)
	dspb.RegisterDatastoreServer(srvDS, NewDatastoreServer(gnmiserver.set))
	lisDS, err := net.Listen("unix", DatastoreAddress)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	go func() {
		if err := srvDS.Serve(lisDS); err != nil {
			log.Errorf("Error while serving datastore target: %v", err)
		}
	}()

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
			log.Errorf("Error while serving gnmi target: %v", err)
		}
	}()
	c.stopFn = func() {
		srvDS.GracefulStop()
		srv.GracefulStop()
	}
	return c, lis.Addr().String(), nil
}

type populateDefaultser interface {
	PopulateDefaults()
}

// New returns a new collector that listens on the specified addr (in the form host:port),
// supporting a single downstream target named hostname. sendMeta controls whether the
// metadata *other* than meta/sync and meta/connected is sent by the collector.
//
// New returns the new collector, the address it is listening on in the form hostname:port
// or any errors encounted whilst setting it up.
//
// NewSettable is different from New in that the returned collector is
// schema-aware and supports gNMI Set. Currently it is not possible to change
// the schema of a Collector after it is created.
func NewSettable(ctx context.Context, addr string, hostname string, sendMeta bool, schema *ytypes.Schema, tasks []Task, opts ...grpc.ServerOption) (*Collector, string, error) {
	if !schema.IsValid() {
		return nil, "", fmt.Errorf("cannot obtain valid schema for GoStructs: %v", schema)
	}

	vr, ok := schema.Root.(ygot.ValidatedGoStruct)
	if !ok {
		return nil, "", fmt.Errorf("invalid schema root, %v", schema.Root)
	}

	// Initialize the root with default values.
	schema.Root.(populateDefaultser).PopulateDefaults()
	if err := vr.Validate(); err != nil {
		return nil, "", fmt.Errorf("default root of input schema fails validation: %v", err)
	}

	// FIXME(wenbli): initialize the collector with default values.
	collector, addr, err := New(ctx, addr, hostname, sendMeta, tasks, opts...)
	if err != nil {
		return nil, "", err
	}
	collector.schema = schema
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

func (c *Collector) Schema() *ytypes.Schema {
	c.smu.Lock()
	defer c.smu.Unlock()
	return c.schema
}

// TargetUpdate provides an input gNMI SubscribeResponse to update the
// cache and clients with.
func (c *Collector) TargetUpdate(m *gpb.SubscribeResponse) {
	c.inCh <- m
}

// setNode is a function that is able to unmarshal either a JSON-encoded value or a gNMI-encoded value.
func setNode(schema *ytypes.Schema, goStruct ygot.GoStruct, update *gpb.Update) error {
	nodeName := reflect.TypeOf(goStruct).Elem().Name()
	nodeI, targetSchema, err := ytypes.GetOrCreateNode(schema.SchemaTree[nodeName], goStruct, update.Path, &ytypes.PreferShadowPath{})
	if err != nil {
		return fmt.Errorf("gnmit: failed to GetOrCreate a node: %v", err)
	}

	// TODO(wenbli): Populate default values using PopulateDefaults.
	jsonUpdate := update.Val.GetJsonIetfVal()
	if jsonUpdate == nil {
		if err := ytypes.SetNode(schema.SchemaTree[nodeName], goStruct, update.Path, update.Val, &ytypes.PreferShadowPath{}); err != nil {
			return fmt.Errorf("gnmit: SetNode failed on leaf node: %v", err)
		}
		return nil
	}

	// TODO(wenbli): Use SetNode natively instead of this. SetNode now is supposed to support setting JSON.
	node, ok := nodeI.(ygot.GoStruct)
	path := proto.Clone(update.Path).(*gpb.Path)
	for i := len(path.Elem) - 1; i >= 0 && !ok; i-- {
		path.Elem = path.Elem[:i]
		nodeI, _, err := ytypes.GetOrCreateNode(schema.SchemaTree[nodeName], goStruct, path, &ytypes.PreferShadowPath{})
		if err != nil {
			continue
		}
		node, ok = nodeI.(ygot.GoStruct)
	}
	if ok {
		var jsonTree interface{}
		if err := json.Unmarshal(jsonUpdate, &jsonTree); err != nil {
			return err
		}
		return ytypes.Unmarshal(targetSchema, node, jsonTree, &ytypes.PreferShadowPath{})
	}

	return fmt.Errorf("gnmit: cannot find GoStruct parent into which to ummarshal update message: %s", prototext.Format(update))
}

func (s *GNMIServer) getOrCreateNode(path *gpb.Path) (ygot.ValidatedGoStruct, ygot.GoStruct, string, error) {
	// Create a copy so that we can rollback the transaction when validation fails.
	dirtyRootG, err := ygot.DeepCopy(s.c.Schema().Root)
	if err != nil {
		return nil, nil, "", fmt.Errorf("gnmit: failed to ygot.DeepCopy the cached root object: %v", err)
	}
	dirtyRoot, ok := dirtyRootG.(ygot.ValidatedGoStruct)
	if !ok {
		return nil, nil, "", fmt.Errorf("gnmit: cannot convert root object to ValidatedGoStruct")
	}
	// Operate at the prefix level.
	nodeI, _, err := ytypes.GetOrCreateNode(s.c.Schema().RootSchema(), dirtyRoot, path, &ytypes.PreferShadowPath{})
	if err != nil {
		return nil, nil, "", fmt.Errorf("gnmit: failed to GetOrCreate the prefix node: %v", err)
	}
	node, ok := nodeI.(ygot.GoStruct)
	if !ok {
		return nil, nil, "", fmt.Errorf("gnmit: prefix path points to a non-GoStruct, this is not allowed: %T, %v", nodeI, nodeI)
	}
	nodeName := reflect.TypeOf(nodeI).Elem().Name()
	return dirtyRoot, node, nodeName, nil
}

func (s *GNMIServer) update(dirtyRoot ygot.ValidatedGoStruct, origin string) error {
	if err := dirtyRoot.Validate(); err != nil {
		return fmt.Errorf("gnmit: invalid SetRequest: %v", err)
	}

	n, err := ygot.Diff(s.c.Schema().Root, dirtyRoot, &ygot.DiffPathOpt{PreferShadowPath: true})
	if err != nil {
		return fmt.Errorf("gnmit: error while creating update notification for Set: %v", err)
	}
	n.Timestamp = time.Now().UnixNano()
	n.Prefix = &gpb.Path{Origin: origin, Target: s.c.name}
	if n.Prefix.Origin == "" {
		n.Prefix.Origin = "openconfig"
	}

	// Update cache
	t := s.c.cache.GetTarget(s.c.name)
	var pathsForDelete []string
	for _, path := range n.Delete {
		p, err := ygot.PathToString(path)
		if err != nil {
			return fmt.Errorf("cannot convert deleted path to string: %v", err)
		}
		pathsForDelete = append(pathsForDelete, p)
	}
	log.V(1).Infof("gnmi.Set: deleting the following paths: %+v", pathsForDelete)
	if err := t.GnmiUpdate(n); err != nil {
		return err
	}
	s.c.Schema().Root = dirtyRoot
	return nil
}

// set updates the datastore and intended configuration with the values from the input notification.
func (s *GNMIServer) set(notif *gpb.Notification) error {
	s.intendedConfigMu.Lock()
	defer s.intendedConfigMu.Unlock()
	dirtyRoot, node, nodeName, err := s.getOrCreateNode(notif.Prefix)
	if err != nil {
		return err
	}

	// Process deletes, then replace, then updates.
	for _, path := range notif.Delete {
		if err := ytypes.DeleteNode(s.c.Schema().SchemaTree[nodeName], node, path, &ytypes.PreferShadowPath{}); err != nil {
			return fmt.Errorf("gnmit: DeleteNode error: %v", err)
		}
	}
	for _, update := range notif.Update {
		if err := setNode(s.c.schema, node, update); err != nil {
			return err
		}
	}

	return s.update(dirtyRoot, notif.Prefix.Origin)
}

// Set is a prototype for a gNMI Set operation.
func (s *GNMIServer) Set(ctx context.Context, req *gpb.SetRequest) (*gpb.SetResponse, error) {
	s.intendedConfigMu.Lock()
	defer s.intendedConfigMu.Unlock()
	if s.c.schema == nil {
		return s.UnimplementedGNMIServer.Set(ctx, req)
	}
	dirtyRoot, node, nodeName, err := s.getOrCreateNode(req.Prefix)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	// TODO(wenbli): Reject paths that try to modify read-only values.
	// TODO(wenbli): Question: what to do if there are operational-state values in a container that is specified to be replaced or deleted?

	// Process deletes, then replace, then updates.
	for _, path := range req.Delete {
		if err := ytypes.DeleteNode(s.c.Schema().SchemaTree[nodeName], node, path, &ytypes.PreferShadowPath{}); err != nil {
			return nil, fmt.Errorf("gnmit: DeleteNode error: %v", err)
		}
	}
	for _, update := range req.Replace {
		if err := ytypes.DeleteNode(s.c.Schema().SchemaTree[nodeName], node, update.Path, &ytypes.PreferShadowPath{}); err != nil {
			return nil, fmt.Errorf("gnmit: DeleteNode error: %v", err)
		}
		if err := setNode(s.c.schema, node, update); err != nil {
			return nil, err
		}
	}
	for _, update := range req.Update {
		if err := setNode(s.c.schema, node, update); err != nil {
			return nil, err
		}
	}

	err = s.update(dirtyRoot, req.Prefix.Origin)
	// TODO(wenbli): Currently the SetResponse is not filled.
	return &gpb.SetResponse{
		Timestamp: time.Now().UnixNano(),
	}, err
}
