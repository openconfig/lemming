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

// Package dataplane is an implementation of the dataplane HAL API.
package dataplane

import (
	"context"
	"fmt"
	"net"

	"github.com/openconfig/lemming/dataplane/forwarding"
	"github.com/openconfig/lemming/dataplane/handlers"
	"github.com/openconfig/lemming/dataplane/internal/engine"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/reconciler"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/grpc/status"

	log "github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	dpb "github.com/openconfig/lemming/proto/dataplane"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Dataplane is an implementation of Dataplane HAL API.
type Dataplane struct {
	dpb.UnimplementedHALServer
	engine      *forwarding.Engine
	srv         *grpc.Server
	lis         net.Listener
	fwd         fwdpb.ServiceClient
	reconcilers []reconciler.Reconciler
}

// New create a new dataplane instance.
func New() (*Dataplane, error) {
	data := &Dataplane{
		engine: forwarding.New("engine"),
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	data.lis = lis
	srv := grpc.NewServer(grpc.Creds(local.NewCredentials()))
	dpb.RegisterHALServer(srv, data)
	fwdpb.RegisterServiceServer(srv, data.engine)
	go srv.Serve(data.lis)

	return data, nil
}

// ID returns the ID of the dataplane reconciler.
func (d *Dataplane) ID() string {
	return "dataplane"
}

// Start starts the HAL gRPC server and packet forwarding engine.
func (d *Dataplane) Start(ctx context.Context, c gpb.GNMIClient, target string) error {
	if d.srv != nil {
		return fmt.Errorf("dataplane already started")
	}

	fc, err := d.FwdClient()
	if err != nil {
		return err
	}
	d.reconcilers = append(d.reconcilers, handlers.NewInterface(fc), handlers.NewRoute(fc))
	d.fwd = fc
	if err := engine.SetupForwardingTables(ctx, fc); err != nil {
		return fmt.Errorf("failed to setup forwarding tables: %v", err)
	}

	for _, rec := range d.reconcilers {
		if err := rec.Start(ctx, c, target); err != nil {
			return fmt.Errorf("failed to stop handler %q: %v", rec.ID(), err)
		}
	}

	return nil
}

// HALClient gets a gRPC client to the dataplane.
func (d *Dataplane) HALClient() (dpb.HALClient, error) {
	conn, err := grpc.Dial(d.lis.Addr().String(), grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}
	return dpb.NewHALClient(conn), nil
}

// FwdClient gets a gRPC client to the packet forwarding engine.
func (d *Dataplane) FwdClient() (fwdpb.ServiceClient, error) {
	conn, err := grpc.Dial(d.lis.Addr().String(), grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}
	return fwdpb.NewServiceClient(conn), nil
}

// Stop gracefully stops the server.
func (d *Dataplane) Stop(ctx context.Context) error {
	for _, rec := range d.reconcilers {
		if err := rec.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop handler %q: %v", rec.ID(), err)
		}
	}
	d.srv.GracefulStop()
	return nil
}

// InsertRoute inserts a route into the dataplane.
func (d *Dataplane) InsertRoute(ctx context.Context, route *dpb.InsertRouteRequest) (*dpb.InsertRouteResponse, error) {
	if len(route.NextHops) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no next hops provided")
	}
	// TODO: support non-default VRF.
	if route.GetVrf() != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "VRF other than DEFAULT (vrfid 0) not supported")
	}

	_, ipNet, err := net.ParseCIDR(route.GetPrefix())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse prefix: %v", err)
	}

	ip := ipNet.IP.To4()
	isIPv4 := true
	if ip == nil {
		ip = ipNet.IP.To16()
		isIPv4 = false
	}

	log.V(1).Infof("inserting route: prefix %s, nexthops %v", route.GetPrefix(), route.GetNextHops())

	if err := engine.AddIPRoute(ctx, d.fwd, isIPv4, ip, ipNet.Mask, route.GetNextHops()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add route: %v", err)
	}

	return &dpb.InsertRouteResponse{}, nil
}

// DeleteRoute deletes a route from the dataplane.
func (d *Dataplane) DeleteRoute(ctx context.Context, route *dpb.DeleteRouteRequest) (*dpb.DeleteRouteResponse, error) {
	// TODO: support non-default VRF.
	if route.GetVrf() != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "VRF other than DEFAULT (vrfid 0) not supported")
	}

	_, ipNet, err := net.ParseCIDR(route.GetPrefix())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse prefix: %v", err)
	}
	isIPv4 := ipNet.IP.To4() != nil

	if err := engine.DeleteIPRoute(ctx, d.fwd, isIPv4, ipNet.IP, ipNet.Mask); err != nil {
		return nil, fmt.Errorf("failed to delete route")
	}

	return &dpb.DeleteRouteResponse{}, nil
}

// Validate is a noop to implement to the reconciler interface.
func (d *Dataplane) Validate(intendedConfig *oc.Root) error {
	return nil
}

// ValidationPaths is a noop to implement to the reconciler interface.
func (d *Dataplane) ValidationPaths() []ygnmi.PathStruct {
	return nil
}
