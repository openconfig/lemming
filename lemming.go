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
	"runtime"
	"sync"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"k8s.io/klog/v2"

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

// Option are device startup options for lemming.
type Option func(*opt)

type opt struct {
	// deviceConfigJSON is the contents of the JSON document (prior to unmarshal).
	deviceConfigJSON []byte
	// tlsCredentials contains TLS credentials that can be used for a device.
	tlsCredentials credentials.TransportCredentials
	gribiAddr      string
	gnmiAddr       string
	p4rtAddr       string
	bgpPort        uint16
}

// resolveOpts applies all the options and returns a struct containing the result.
func resolveOpts(opts []Option) *opt {
	o := &opt{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithInitialConfig sets the startup config of the device to c.
// Today we do not allow the configuration to be changed in flight, but this
// can be implemented in the future.
//
// This option is intended for standalone device testing.
func WithInitialConfig(c []byte) Option {
	return func(o *opt) {
		o.deviceConfigJSON = c
	}
}

// WithGRIBIAddr is a device option that specifies that the gRIBI address.
func WithGRIBIAddr(addr string) Option {
	return func(o *opt) {
		o.gribiAddr = addr
	}
}

// WithGNMIAddr is a device option that specifies that the gRIBI address.
func WithGNMIAddr(addr string) Option {
	return func(o *opt) {
		o.gnmiAddr = addr
	}
}

// WithP4RTAddr is a device option that specifies that the P4RT address.
func WithP4RTAddr(addr string) Option {
	return func(o *opt) {
		o.p4rtAddr = addr
	}
}

// WithBGPPort is a device option that specifies that the BGP port.
func WithBGPPort(port uint16) Option {
	return func(o *opt) {
		o.bgpPort = port
	}
}

// WithTLSCredsFromFile loads the credentials from the specified cert and key file
// and returns them such that they can be used for the gNMI and gRIBI servers.
func WithTLSCredsFromFile(certFile, keyFile string) (Option, error) {
	t, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return func(o *opt) {
		o.tlsCredentials = t
	}, nil
}

// WithTransportCreds returns a wrapper of TransportCredentials into a DevOpt.
func WithTransportCreds(c credentials.TransportCredentials) Option {
	return func(o *opt) {
		o.tlsCredentials = c
	}
}

// New returns a new initialized device.
func New(targetName, zapiURL string, opts ...Option) (*Device, error) {
	var dplane *dataplane.Dataplane
	var recs []reconciler.Reconciler

	if viper.GetBool("enable_dataplane") {
		if runtime.GOOS != "linux" {
			return nil, fmt.Errorf("dataplane only supported on linux, GOOS is %s", runtime.GOOS)
		}

		log.Info("enabling dataplane")
		var err error
		dplane, err = dataplane.New(context.Background())
		if err != nil {
			return nil, err
		}
		recs = append(recs, dplane)
	}

	log.Info("starting gNMI")

	resolvedOpts := resolveOpts(opts)

	root := &oc.Root{}
	if jcfg := resolvedOpts.deviceConfigJSON; jcfg != nil {
		if err := oc.Unmarshal(jcfg, root); err != nil {
			return nil, fmt.Errorf("cannot unmarshal JSON configuration, %v", err)
		}
	} else {
		// The initial config may specify a differently-named network
		// instance, so only add when it is not present.
		root.GetOrCreateNetworkInstance(fakedevice.DefaultNetworkInstance).Type = oc.NetworkInstanceTypes_NETWORK_INSTANCE_TYPE_DEFAULT_INSTANCE
	}

	var grpcOpts []grpc.ServerOption
	creds := resolvedOpts.tlsCredentials
	if creds != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}
	grpcOpts = append(grpcOpts, grpc.StreamInterceptor(fgnmi.NewStreamInterceptorFn(targetName)))

	s := grpc.NewServer(grpcOpts...)

	recs = append(recs,
		fakedevice.NewSystemBaseTask(),
		fakedevice.NewBootTimeTask(),
		bgp.NewGoBGPTask(targetName, zapiURL, resolvedOpts.bgpPort),
	)

	log.Info("starting gNSI")
	gnsiServer := fgnsi.New(s)

	gnmiServer, err := fgnmi.New(s, targetName, gnsiServer.GetPathZ(), recs...)
	if err != nil {
		return nil, err
	}

	cacheClient := gnmiServer.LocalClient()

	log.Infof("starting sysrib")
	sysribServer, err := sysrib.New(root)
	if err != nil {
		return nil, err
	}
	if err := sysribServer.Start(cacheClient, targetName, zapiURL); err != nil {
		return nil, fmt.Errorf("sysribServer failed to start: %v", err)
	}

	log.Info("starting gRIBI")
	// TODO(wenbli): Use gRIBIs once we change lemming's KNE config to use different ports.
	// gRIBIs := grpc.NewServer()
	gribiServer, err := fgribi.New(s, cacheClient, targetName, root)
	if err != nil {
		return nil, err
	}

	log.Info("starting P4RT (there is nothing here yet)")
	P4RTs := grpc.NewServer()

	log.Info("Create listeners")
	lgnmi, err := net.Listen("tcp", resolvedOpts.gnmiAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC server for gNMI/gNOI/gNSI, %v", err)
	}

	lgribi, err := net.Listen("tcp", resolvedOpts.gribiAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC server for gRIBI, %v", err)
	}

	lp4rt, err := net.Listen("tcp", resolvedOpts.p4rtAddr)
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

	if err := gnmiServer.StartReconcilers(context.Background()); err != nil {
		return nil, err
	}

	log.Info("lemming created")
	return d, nil
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

	if len(d.errs) == 0 {
		return nil
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
				d.errs = append(d.errs, err)
				d.errsMu.Unlock()
			}
			d.errsMu.Lock()
			klog.Infof("%s server stopped: %v", svcName, d.errs)
			d.errsMu.Unlock()
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
