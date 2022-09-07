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

// Package lemming provides reference device to be used with ondatra.
package lemming

import (
	"net"
	"sync"

	fgnmi "github.com/openconfig/lemming/gnmi"
	fgnoi "github.com/openconfig/lemming/gnoi"
	fgnsi "github.com/openconfig/lemming/gnsi"
	fgribi "github.com/openconfig/lemming/gribi"
	fp4rt "github.com/openconfig/lemming/p4rt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/klog/v2"
)

// Device is the reference device implementation.
type Device struct {
	s           *grpc.Server
	lis         net.Listener
	stop        func()
	gnmiServer  *fgnmi.Server
	gnoiServer  *fgnoi.Server
	gribiServer *fgribi.Server
	gnsiServer  *fgnsi.Server
	p4rtServer  *fp4rt.Server
	// Stores the error if the server fails will be returned on call to stop.
	mu      sync.Mutex
	err     error
	stopped chan struct{}
}

// New returns a new initialized device.
func New(lis net.Listener, opts ...grpc.ServerOption) (*Device, error) {
	s := grpc.NewServer(opts...)
	d := &Device{
		lis:         lis,
		s:           s,
		gnmiServer:  fgnmi.New(s),
		gnoiServer:  fgnoi.New(s),
		gribiServer: fgribi.New(s),
		gnsiServer:  fgnsi.New(s),
		p4rtServer:  fp4rt.New(s),
	}
	reflection.Register(s)
	d.startServer()
	return d, nil
}

// Addr returns the currently configured ip:port for the listening services.
func (d *Device) Addr() string {
	return d.lis.Addr().String()
}

// Stop stops the listening services.
// If error is not nil, it will contain why the server failed.
func (d *Device) Stop() error {
	klog.Info("Stopping server")
	select {
	case <-d.stopped:
		klog.Info("Server already stopped: ", d.err)
	default:
		d.stop()
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.err
}

// GNMI returns the gNMI server implementation.
func (d *Device) GNMI() *fgnmi.Server {
	return d.gnmiServer
}

// GNSI returns the gNSI server implementation.
func (d *Device) GNSI() *fgnsi.Server {
	return d.gnsiServer
}

func (d *Device) startServer() {
	d.stopped = make(chan struct{})
	go func() {
		err := d.s.Serve(d.lis)
		d.mu.Lock()
		defer d.mu.Unlock()
		d.err = err
		klog.Infof("Server stopped: %v", err)
		close(d.stopped)
	}()
	d.stop = func() {
		d.s.Stop()
		<-d.stopped
	}
}
