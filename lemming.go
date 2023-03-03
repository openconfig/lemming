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
	"sync"

	"github.com/openconfig/lemming/bgp"
	"github.com/openconfig/lemming/dataplane"
	fgnmi "github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/reconciler"
	fgnoi "github.com/openconfig/lemming/gnoi"
	fgnsi "github.com/openconfig/lemming/gnsi"
	fgribi "github.com/openconfig/lemming/gribi"
	fp4rt "github.com/openconfig/lemming/p4rt"
	"github.com/openconfig/lemming/sysrib"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"k8s.io/klog/v2"

	log "github.com/golang/glog"
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

// DevOpt is an interface that is implemented by options that can be handed to New()
// for the device.
type DevOpt interface {
	isDevOpt()
}

// deviceConfig is a wrapper for an input OpenConfig RFC7951-marshalled JSON
// configuration for the device.
type deviceConfig struct {
	// json is the contents of the JSON document (prior to unmarshal).
	json []byte
}

// isDevOpt marks deviceConfig as a device option.
func (*deviceConfig) isDevOpt() {}

// DeviceConfig sets the startup config of the device to c.
// Today we do not allow the configuration to be changed in flight, but this
// can be implemented in the future.
//
// This DeviceOption is intended for standalone device testing.
func DeviceConfig(c []byte) *deviceConfig {
	return &deviceConfig{json: c}
}

// tlsCreds returns TLS credentials that can be used for a device.
type tlsCreds struct {
	c credentials.TransportCredentials
}

// TLSCredsFromFile loads the credentials from the specified cert and key file
// and returns them such that they can be used for the gNMI and gRIBI servers.
func TLSCredsFromFile(certFile, keyFile string) (*tlsCreds, error) {
	t, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tlsCreds{c: t}, nil
}

// TLSCreds returns a wrapper of TransportCredentials into a DevOpt.
func TLSCreds(c credentials.TransportCredentials) *tlsCreds {
	return &tlsCreds{c: c}
}

// IsDevOpt implements the DevOpt interface for tlsCreds.
func (*tlsCreds) isDevOpt() {}

// New returns a new initialized device.
func New(lis net.Listener, targetName, zapiURL string, opts ...DevOpt) (*Device, error) {
	var dplane *dataplane.Dataplane
	var recs []reconciler.Reconciler

	if viper.GetBool("enable_dataplane") {
		log.Info("enabling dataplane")
		var err error
		dplane, err = dataplane.New()
		if err != nil {
			return nil, err
		}
		recs = append(recs, dplane)
	}

	jcfg := optDeviceCfg(opts)
	root := &oc.Root{}
	switch jcfg {
	case nil:
		root.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
	default:
		if err := oc.Unmarshal(jcfg, root); err != nil {
			return nil, fmt.Errorf("cannot unmarshal JSON configuration, %v", err)
		}
	}

	var grpcOpts []grpc.ServerOption
	creds := optTLSCreds(opts)
	if creds != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(creds.c))
	}

	s := grpc.NewServer(grpcOpts...)

	recs = append(recs,
		fakedevice.NewSystemBaseTask(),
		fakedevice.NewBootTimeTask(),
		bgp.NewGoBGPTaskDecl(zapiURL),
	)

	gnsiServer := fgnsi.New(s)

	gnmiServer, err := fgnmi.New(s, targetName, gnsiServer.GetPathZ(), recs...)
	if err != nil {
		return nil, err
	}

	log.Info("starting gRIBI")
	gribiServer, err := fgribi.New(s, root)
	if err != nil {
		return nil, err
	}

	d := &Device{
		lis:         lis,
		s:           s,
		gnmiServer:  gnmiServer,
		gnoiServer:  fgnoi.New(s),
		gribiServer: gribiServer,
		gnsiServer:  gnsiServer,
		p4rtServer:  fp4rt.New(s),
	}
	reflection.Register(s)
	d.startServer()

	cacheClient := gnmiServer.LocalClient()

	log.Infof("starting sysrib")
	sysribServer, err := sysrib.New(root)
	if err != nil {
		return nil, err
	}
	if err := sysribServer.Start(cacheClient, targetName, zapiURL); err != nil {
		return nil, fmt.Errorf("sysribServer failed to start: %v", err)
	}

	if err := gnmiServer.StartReconcilers(context.Background()); err != nil {
		return nil, err
	}

	log.Info("lemming created")
	return d, nil
}

// optDeviceCfg finds the first occurrence of the DeviceConfig option in opts.
func optDeviceCfg(opts []DevOpt) []byte {
	for _, o := range opts {
		if v, ok := o.(*deviceConfig); ok {
			return v.json
		}
	}
	return nil
}

// optTLSCreds finds the first occurrence of the tlsCreds option in opts.
func optTLSCreds(opts []DevOpt) *tlsCreds {
	for _, o := range opts {
		if v, ok := o.(*tlsCreds); ok {
			return v
		}
	}
	return nil
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
	if err := d.gnmiServer.StopReconcilers(context.Background()); err != nil {
		return err
	}

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
