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
	"fmt"
	"net"
	"net/netip"
	"os"
	"sync"

	"github.com/openconfig/lemming/dataplane"
	fgnmi "github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/gnmit"
	fgnoi "github.com/openconfig/lemming/gnoi"
	fgnsi "github.com/openconfig/lemming/gnsi"
	fgribi "github.com/openconfig/lemming/gribi"
	fp4rt "github.com/openconfig/lemming/p4rt"
	"github.com/openconfig/lemming/sysrib"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"
	"google.golang.org/grpc/reflection"
	"k8s.io/klog/v2"

	log "github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
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

// startSysrib starts the sysrib gRPC service at a unix domain socket. This
// should be started prior to routing services to allow them to connect to
// sysrib during their initialization.
func startSysrib(dataplane *sysrib.Dataplane, port int, target string, enableTLS bool) {
	if err := os.RemoveAll(sysrib.SockAddr); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("unix", sysrib.SockAddr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcServer := grpc.NewServer()
	s, err := sysrib.NewServer(dataplane, port, target, enableTLS)
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

	s := grpc.NewServer(opts...)

	gnmiServer, err := fgnmi.New(s, targetName)
	if err != nil {
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
	port, enableTLS := viper.GetInt("port"), viper.GetBool("enable_tls")
	if port == 0 {
		addrport, err := netip.ParseAddrPort(lis.Addr().String())
		if err != nil {
			return nil, err
		}
		port = int(addrport.Port())
	}

	if dplane != nil {
		setConn, err := grpc.Dial(fmt.Sprintf("unix:///%s", gnmit.DatastoreAddress), grpc.WithTransportCredentials(local.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to dial unix socket: %v", err)
		}
		if err := dplane.Start(context.Background(), gnmiclient.NewCacheClient(gnmiServer, gpb.NewGNMIClient(setConn)), targetName); err != nil {
			return nil, err
		}
	}

	log.Infof("starting sysrib")
	startSysrib(sysDataplane, port, targetName, enableTLS)

	fakedevice.StartSystemBaseTask(context.Background(), port, targetName, enableTLS)
	fakedevice.StartBootTimeTask(context.Background(), port, targetName, enableTLS)
	//fakedevice.StartCurrentDateTimeTask(context.Background(), port, targetName, enableTLS)
	fakedevice.StartGoBGPTask(context.Background(), port, targetName, enableTLS)

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
