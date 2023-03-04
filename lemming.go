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

type gRPCService struct {
	s       *grpc.Server
	lis     net.Listener
	stopped chan struct{}
}

// Device is the reference device implementation.
type Device struct {
	gnmignoignsiService *gRPCService
	gribiService        *gRPCService
	p4rtService         *gRPCService
	stop                func()

	gnmiServer  *fgnmi.Server
	gnoiServer  *fgnoi.Server
	gribiServer *fgribi.Server
	gnsiServer  *fgnsi.Server
	p4rtServer  *fp4rt.Server
	// Stores the errors if the server fails will be returned on call to stop.
	errsMu sync.Mutex
	errs   []error

	stopped chan struct{}
}

// DevOpt is an interface that is implemented by options that can be handed to New()
// for the device.
type DevOpt interface {
	isDevOpt()
}

// gRIBIAddr is the internal implementation that specifies the port that gRIBI should
// listen on.
type gRIBIAddr struct {
	host string
	port int
}

// isDevOpt implements the DevOpt interface.
func (*gRIBIAddr) isDevOpt() {}

// GRIBIAddr is a device option that specifies that the port that should be listened on
// is i.
func GRIBIAddr(host string, i int) *gRIBIAddr {
	return &gRIBIAddr{host: host, port: i}
}

// gNMIAddress is the internal implementation that specifies the port that gNMI should
// listen on.
type gNMIAddr struct {
	host string
	port int
}

// isDevOpt implements the DevOpt interface.
func (*gNMIAddr) isDevOpt() {}

// GNMIAddr specifies the host and port that the gNMI server should listen on.
func GNMIAddr(host string, i int) *gNMIAddr {
	return &gNMIAddr{host: host, port: i}
}

// p4RTAddress is the internal implementation that specifies the port that p4RT should
// listen on.
type p4RTAddr struct {
	host string
	port int
}

// isDevOpt implements the DevOpt interface.
func (*p4RTAddr) isDevOpt() {}

// P4RTAddr specifies the host and port that the p4RT server should listen on.
func P4RTAddr(host string, i int) *p4RTAddr {
	return &p4RTAddr{host: host, port: i}
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
func New(targetName, zapiURL string, opts ...DevOpt) (*Device, error) {
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

	log.Info("starting gNMI")
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

	log.Info("starting gNSI")
	gnsiServer := fgnsi.New(s)

	gnmiServer, err := fgnmi.New(s, targetName, gnsiServer.GetPathZ(), recs...)
	if err != nil {
		return nil, err
	}

	log.Info("starting gRIBI")
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

	// TODO(wenbli): Use gRIBIs once we change lemming's KNE config to use different ports.
	//gRIBIs := grpc.NewServer()
	gribiServer, err := fgribi.New(s, root)
	if err != nil {
		return nil, err
	}

	log.Info("starting P4RT (there is nothing here yet)")
	P4RTs := grpc.NewServer()

	log.Info("Create listeners")
	gr := optGRIBIAddr(opts)
	gn := optGNMIAddr(opts)
	p4 := optP4RTAddr(opts)

	lgnmi, err := net.Listen("tcp", fmt.Sprintf("%s:%d", gn.host, gn.port))
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC server for gNMI/gNOI/gNSI, %v", err)
	}

	lgribi, err := net.Listen("tcp", fmt.Sprintf("%s:%d", gr.host, gr.port))
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC server for gRIBI, %v", err)
	}

	lp4rt, err := net.Listen("tcp", fmt.Sprintf("%s:%d", p4.host, p4.port))
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC server for P4RT, %v", err)
	}

	d := &Device{
		gnmignoignsiService: &gRPCService{
			s:       s,
			lis:     lgnmi,
			stopped: make(chan struct{}),
		},
		gribiService: &gRPCService{
			// TODO(wenbli): Change s to gRIBIs once we change lemming's KNE config to use different ports.
			s:       s,
			lis:     lgribi,
			stopped: make(chan struct{}),
		},
		p4rtService: &gRPCService{
			s:       P4RTs,
			lis:     lp4rt,
			stopped: make(chan struct{}),
		},
		gnmiServer:  gnmiServer,
		gnoiServer:  fgnoi.New(s),
		gribiServer: gribiServer,
		gnsiServer:  gnsiServer,
		p4rtServer:  fp4rt.New(P4RTs),
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

// optGRIBIAddr finds the first occurrence of the GRIBIAddr option in opts.
// If no GRIBIAddr option is found, the default of localhost:0 is returned.
func optGRIBIAddr(opts []DevOpt) *gRIBIAddr {
	for _, o := range opts {
		if v, ok := o.(*gRIBIAddr); ok {
			return v
		}
	}
	return &gRIBIAddr{host: "localhost", port: 0}
}

// optGNMIAddr finds the first occurrence of the GNMIAddr option in opts.
// If no GNMIAddr option is found, the default of localhost:0 is returned.
func optGNMIAddr(opts []DevOpt) *gNMIAddr {
	for _, o := range opts {
		if v, ok := o.(*gNMIAddr); ok {
			return v
		}
	}
	return &gNMIAddr{host: "localhost", port: 0}
}

// optP4RTAddr finds the first occurrence of the P4RTAddr option in opts.
// If no P4RTAddr option is found, the default of localhost:0 is returned.
func optP4RTAddr(opts []DevOpt) *p4RTAddr {
	for _, o := range opts {
		if v, ok := o.(*p4RTAddr); ok {
			return v
		}
	}
	return &p4RTAddr{host: "localhost", port: 0}
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

// GRIBIAddr returns the address that the gRIBI server is listening on.
func (d *Device) GRIBIAddr() string {
	return d.gribiService.lis.Addr().String()
}

// GNMIAddr returns the address that the gNMI/gNOI/gNSI server is listening on.
func (d *Device) GNMIAddr() string {
	return d.gnmignoignsiService.lis.Addr().String()
}

// P4RTAddr returns the address that the P4RT server is listening on.
func (d *Device) P4RTAddr() string {
	return d.p4rtService.lis.Addr().String()
}

// GRIBIListener returns the listener that the gRIBI server is listening on.
func (d *Device) GRIBIListener() net.Listener {
	return d.gribiService.lis
}

// GNMIListener returns the listener that the gNMI/gNOI/gNSI server is listening on.
func (d *Device) GNMIListener() net.Listener {
	return d.gnmignoignsiService.lis
}

// P4RTListener returns the listener that the P4RT server is listening on.
func (d *Device) P4RTListener() net.Listener {
	return d.p4rtService.lis
}

// Stop stops the listening services.
// If error is not nil, it will contain why the server failed.
func (d *Device) Stop() error {
	klog.Info("Stopping server")
	select {
	case <-d.stopped:
		klog.Infof("Server already stopped: %v", d.errs)
	default:
		d.stop()
	}
	d.errsMu.Lock()
	defer d.errsMu.Unlock()
	if err := d.gnmiServer.StopReconcilers(context.Background()); err != nil {
		d.errs = append(d.errs, err)
	}

	return fmt.Errorf("%v", d.errs)
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
	services := map[string]*gRPCService{
		"gNMI/gNOI/gNSI": d.gnmignoignsiService,
		"gRIBI":          d.gribiService,
		"P4RT":           d.p4rtService,
	}
	for svcName, svc := range services {
		// Capture loop variables by value instead of reference.
		svcName, svc := svcName, svc
		go func() {
			if err := svc.s.Serve(svc.lis); err != nil {
				d.errsMu.Lock()
				defer d.errsMu.Unlock()
				d.errs = append(d.errs, err)
			}
			klog.Infof("%s server stopped: %v", svcName, d.errs)
			close(svc.stopped)
		}()
	}

	d.stop = func() {
		for _, svc := range services {
			svc.s.Stop()
		}
		for _, svc := range services {
			<-svc.stopped
		}
	}
}
