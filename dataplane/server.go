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
	"google.golang.org/grpc/credentials/local"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// Dataplane is an implementation of Dataplane HAL API.
type Dataplane struct {
	fwdSrv      *forwarding.Engine
	e           *engine.Engine
	srv         *grpc.Server
	lis         net.Listener
	fwd         fwdpb.ServiceClient
	reconcilers []reconciler.Reconciler
}

// New create a new dataplane instance.
func New() (*Dataplane, error) {
	data := &Dataplane{
		fwdSrv: forwarding.New("engine"),
	}

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	data.lis = lis
	srv := grpc.NewServer(grpc.Creds(local.NewCredentials()))
	fwdpb.RegisterServiceServer(srv, data.fwdSrv)
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
	d.fwd = fc

	d.e, err = engine.New(ctx, "default", fc)
	if err != nil {
		return err
	}

	d.reconcilers = append(d.reconcilers, handlers.NewInterface(d.e), handlers.NewRoute(d.e))

	for _, rec := range d.reconcilers {
		if err := rec.Start(ctx, c, target); err != nil {
			return fmt.Errorf("failed to stop handler %q: %v", rec.ID(), err)
		}
	}

	return nil
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

// Validate is a noop to implement to the reconciler interface.
func (d *Dataplane) Validate(intendedConfig *oc.Root) error {
	return nil
}

// ValidationPaths is a noop to implement to the reconciler interface.
func (d *Dataplane) ValidationPaths() []ygnmi.PathStruct {
	return nil
}
