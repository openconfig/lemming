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

// Package gnmi contains a reference on-device gNMI implementation.
package gnmi

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/cache"
	"github.com/openconfig/gnmi/subscribe"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/reconciler"
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
	OpenConfigOrigin = "openconfig"
)

// Server is a reference gNMI implementation.
type Server struct {
	// The subscribe Server implements only Subscribe for gNMI.
	*subscribe.Server
	c *Collector

	// TODO(wenbli): Implement gnmi.Get and remove this.
	GetResponses []interface{}

	configMu     sync.Mutex
	configSchema *ytypes.Schema

	stateMu     sync.Mutex
	stateSchema *ytypes.Schema

	validators  []func(*oc.Root) error
	reconcilers []reconciler.Reconciler
}

// New creates and registers a reference gNMI server on the given gRPC server.
//
// - targetName is the gNMI target name of the datastore.
func New(srv *grpc.Server, targetName string, recs ...reconciler.Reconciler) (*Server, error) {
	gnmiServer, err := newServer(context.Background(), targetName, true, recs...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gNMI server: %v", err)
	}
	gpb.RegisterGNMIServer(srv, gnmiServer)
	return gnmiServer, nil
}

// newServer returns a new reference gNMI server implementation that can be
// registered on an existing gRPC server.
//
// - configSchema is the specification of the schema if gnmi.Set on config paths is used.
// - stateSchema is the specification of the schema if gnmi.Set on state paths is used.
// - targetName is the name of the target.
func newServer(ctx context.Context, targetName string, enableSet bool, recs ...reconciler.Reconciler) (*Server, error) {
	c := NewCollector(targetName)
	subscribeSrv, err := c.Start(ctx, false)
	if err != nil {
		return nil, err
	}

	gnmiServer := &Server{
		Server:      subscribeSrv, // use the 'subscribe' implementation.
		c:           c,
		reconcilers: recs,
	}

	if !enableSet {
		return gnmiServer, nil
	}

	configSchema, err := oc.Schema()
	if err != nil {
		return nil, fmt.Errorf("cannot create ygot schema object: %v", err)
	}
	// Initialize the cache with the input schema root.
	if configSchema != nil {
		if err := setupSchema(configSchema); err != nil {
			return nil, err
		}
		if err := ygot.PruneConfigFalse(configSchema.RootSchema(), configSchema.Root); err != nil {
			return nil, fmt.Errorf("gnmi: %v", err)
		}
		updateCache(c.cache, configSchema.Root, nil, targetName, OpenConfigOrigin, true)
	}

	stateSchema, err := oc.Schema()
	if err != nil {
		return nil, fmt.Errorf("cannot create ygot schema object: %v", err)
	}
	if stateSchema != nil {
		if err := setupSchema(stateSchema); err != nil {
			return nil, err
		}
		updateCache(c.cache, stateSchema.Root, nil, targetName, OpenConfigOrigin, true)
	}

	for _, rec := range recs {
		if len(rec.ValidationPaths()) > 0 {
			gnmiServer.validators = append(gnmiServer.validators, rec.Validate)
		}
	}

	gnmiServer.configSchema = configSchema
	gnmiServer.stateSchema = stateSchema

	return gnmiServer, nil
}

type populateDefaultser interface {
	PopulateDefaults()
}

// setupSchema takes in a ygot schema object which it assumes to be
// uninitialized. It initializes and validates it, returning any errors
// encountered.
func setupSchema(schema *ytypes.Schema) error {
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

// updateCache updates the cache with the difference between the root ->
// dirtyRoot such that if the root represents the cache, then the dirtyRoot
// will represent the cache afterwards.
//
// If root is nil, then it is assumed the cache is empty, and the entirety of
// the dirtyRoot is put into the cache.
func updateCache(cache *cache.Cache, dirtyRoot, root ygot.GoStruct, target, origin string, preferShadowPath bool) error {
	var nos []*gpb.Notification
	if root == nil {
		var err error
		if nos, err = ygot.TogNMINotifications(dirtyRoot, time.Now().UnixNano(), ygot.GNMINotificationsConfig{
			UsePathElem: true,
		}); err != nil {
			return fmt.Errorf("gnmi: %v", err)
		}
	} else {
		n, err := ygot.Diff(root, dirtyRoot, &ygot.DiffPathOpt{PreferShadowPath: preferShadowPath})
		if err != nil {
			return fmt.Errorf("gnmi: error while creating update notification for Set: %v", err)
		}
		n.Timestamp = time.Now().UnixNano()
		nos = append(nos, n)
	}

	return updateCacheNotifs(cache, nos, target, origin)
}

// updateCacheNotifs updates the target cache with the given notifications.
func updateCacheNotifs(cache *cache.Cache, nos []*gpb.Notification, target, origin string) error {
	cacheTarget := cache.GetTarget(target)
	for _, n := range nos {
		n.Prefix = &gpb.Path{Origin: origin, Target: target}
		if n.Prefix.Origin == "" {
			n.Prefix.Origin = OpenConfigOrigin
		}

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
		if err := cacheTarget.GnmiUpdate(n); err != nil {
			return err
		}
	}
	return nil
}

// set updates the datastore and intended configuration with the SetRequest,
// allowing read-only values to be updated.
//
// update indicates whether to update the cache with the values from the set
// request.
// set returns a gRPC error with the correct code and shouldn't be wrapped again.
func set(schema *ytypes.Schema, cache *cache.Cache, target string, req *gpb.SetRequest, preferShadowPath bool, validators []func(*oc.Root) error) error {
	prevRoot, err := ygot.DeepCopy(schema.Root)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to ygot.DeepCopy the cached root object: %v", err)
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
		return status.Errorf(codes.InvalidArgument, "failed to unmarshal set request %v", err)
	}

	if err := schema.Validate(); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid SetRequest: %v", err)
	}
	for _, validator := range validators {
		if err := validator(schema.Root.(*oc.Root)); err != nil {
			return status.Errorf(codes.InvalidArgument, "validation error: %v", err)
		}
	}

	success = true

	if err := updateCache(cache, schema.Root, prevRoot, target, req.Prefix.Origin, preferShadowPath); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

// Set is a prototype for a gNMI Set operation.
func (s *Server) Set(ctx context.Context, req *gpb.SetRequest) (*gpb.SetResponse, error) {
	// Use ConfigMode by default so that external users don't need to set metadata.
	gnmiMode := ConfigMode
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		switch {
		case slices.Contains(md.Get(GNMIModeMetadataKey), string(ConfigMode)):
			gnmiMode = ConfigMode
		case slices.Contains(md.Get(GNMIModeMetadataKey), string(StateMode)):
			gnmiMode = StateMode
		}
	}

	switch gnmiMode {
	case ConfigMode:
		s.configMu.Lock()
		defer s.configMu.Unlock()

		log.V(1).Infof("config datastore service received SetRequest: %v", req)
		if s.configSchema == nil {
			return s.UnimplementedGNMIServer.Set(ctx, req)
		}

		// TODO(wenbli): Reject paths that try to modify read-only values.
		// TODO(wenbli): Question: what to do if there are operational-state values in a container that is specified to be replaced or deleted?
		err := set(s.configSchema, s.c.cache, s.c.name, req, true, s.validators)

		// TODO(wenbli): Currently the SetResponse is not filled.
		return &gpb.SetResponse{
			Timestamp: time.Now().UnixNano(),
		}, err
	case StateMode:
		s.stateMu.Lock()
		defer s.stateMu.Unlock()

		log.V(1).Infof("operational state datastore service received SetRequest: %v", req)
		if s.stateSchema == nil {
			return s.UnimplementedGNMIServer.Set(ctx, req)
		}
		// TODO(wenbli): Reject values that modify config values. We only allow modifying state in this mode.
		if err := set(s.stateSchema, s.c.cache, s.c.name, req, false, nil); err != nil {
			return &gpb.SetResponse{}, err
		}

		// This mode is intended to be used internally, and the SetResponse doesn't matter.
		return &gpb.SetResponse{}, nil
	default:
		return nil, status.Errorf(codes.Internal, "incoming gNMI SetRequest must specify a valid gnmi.Mode via context metadata: %v", md)
	}
}

func (s *Server) Capabilities(ctx context.Context, req *gpb.CapabilityRequest) (*gpb.CapabilityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Reference Implementation Unimplemented")
}

func (s *Server) Get(ctx context.Context, req *gpb.GetRequest) (*gpb.GetResponse, error) {
	if len(s.GetResponses) == 0 {
		return nil, status.Errorf(codes.Unimplemented, "Reference Implementation Unimplemented")
	}
	resp := s.GetResponses[0]
	s.GetResponses = s.GetResponses[1:]
	switch v := resp.(type) {
	case error:
		return nil, v
	case *gpb.GetResponse:
		return v, nil
	default:
		return nil, status.Errorf(codes.DataLoss, "Unknown message type: %T", resp)
	}
}

// LocalClient returns a gNMI client for the server.
func (s *Server) LocalClient() gpb.GNMIClient {
	return newLocalClient(s)
}

// StartReconcilers starts all the reconcilers.
func (s *Server) StartReconcilers(ctx context.Context) error {
	c := s.LocalClient()
	for _, rec := range s.reconcilers {
		if err := rec.Start(ctx, c, s.c.name); err != nil {
			return err
		}
	}
	return nil
}

// StopReconcilers stops all the reconcilers.
func (s *Server) StopReconcilers(ctx context.Context) error {
	for _, rec := range s.reconcilers {
		if err := rec.Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}
