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
	"context"
	"net"
	"os"
	"sync"

	"github.com/openconfig/lemming/dataplane"
	fgnmi "github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/gnmit"
	"github.com/openconfig/lemming/gnmi/testagentlocal"
	fgnoi "github.com/openconfig/lemming/gnoi"
	fgnsi "github.com/openconfig/lemming/gnsi"
	fgribi "github.com/openconfig/lemming/gribi"
	fp4rt "github.com/openconfig/lemming/p4rt"
	"github.com/openconfig/lemming/sysrib"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/klog/v2"

	log "github.com/golang/glog"
	zpb "github.com/openconfig/lemming/proto/sysrib"
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

// registerTestTask registers a test gothread that reads from the central
// datastore.
//
// Note: This should only be used for testing lemming, since interface paths
// should be owned by the dataplane module.
func registerTestTask(gnmiServer *gnmit.GNMIServer, targetName string) error {
	return gnmiServer.RegisterTask(testagentlocal.InterfaceTask(targetName))
}

// startSysrib starts the sysrib gRPC service at a unix domain socket. This
// should be started prior to routing services to allow them to connect to
// sysrib during their initialization.
func startSysrib(dataplane *sysrib.Dataplane, gnmiServerAddr, target string) {
	if err := os.RemoveAll(sysrib.SockAddr); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("unix", sysrib.SockAddr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	s, err := sysrib.NewServer(dataplane, gnmiServerAddr, target)
	if err != nil {
		log.Fatalf("error while creating sysrib server: %v", err)
	}
	zpb.RegisterSysribServer(grpcServer, s)

	go func() {
		grpcServer.Serve(lis)
	}()
}

// New returns a new initialized device.
func New(lis net.Listener, targetName string, opts ...grpc.ServerOption) (*Device, error) {
	var sysDataplane *sysrib.Dataplane
	var dplane *dataplane.Dataplane
	if viper.GetBool("enable_dataplane") {
		log.Info("enabling dataplane")
		var err error
		dplane, err = dataplane.New()
		if err != nil {
			return nil, err
		}
		hal, err := dplane.HALClient()
		if err != nil {
			return nil, err
		}
		sysDataplane = &sysrib.Dataplane{HALClient: hal}
	}

	log.Infof("starting sysrib: gNMI server(%s, %s)", lis.Addr().String(), targetName)
	startSysrib(sysDataplane, lis.Addr().String(), targetName)

	s := grpc.NewServer(opts...)

	gnmiServer, err := fgnmi.New(s, targetName)
	if err != nil {
		return nil, err
	}
	if err := registerTestTask(gnmiServer.GNMIServer, targetName); err != nil {
		return nil, err
	}

	log.Info("starting gRIBI")
	gribiServer, err := fgribi.New(s)
	if err != nil {
		return nil, err
	}

	d := &Device{
		lis:         lis,
		s:           s,
		gnmiServer:  gnmiServer,
		gnoiServer:  fgnoi.New(s),
		gribiServer: gribiServer,
		gnsiServer:  fgnsi.New(s),
		p4rtServer:  fp4rt.New(s),
	}
	reflection.Register(s)
	d.startServer()
	if dplane != nil {
		if err := dplane.Start(context.Background()); err != nil {
			return nil, err
		}
	}

	log.Info("lemming created")
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
