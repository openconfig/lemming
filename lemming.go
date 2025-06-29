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

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"k8s.io/klog/v2"

	gribis "github.com/openconfig/gribigo/server"

	"github.com/openconfig/lemming/bgp"
	"github.com/openconfig/lemming/dataplane"
	"github.com/openconfig/lemming/dataplane/dplaneopts"
	"github.com/openconfig/lemming/fault"
	fgnmi "github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/reconciler"
	fgnoi "github.com/openconfig/lemming/gnoi"
	fgnsi "github.com/openconfig/lemming/gnsi"
	fgribi "github.com/openconfig/lemming/gribi"
	"github.com/openconfig/lemming/internal/config"
	fp4rt "github.com/openconfig/lemming/p4rt"
	configpb "github.com/openconfig/lemming/proto/config"
	faultpb "github.com/openconfig/lemming/proto/fault"
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
	faultService        *gRPCService
	stop                func()

	gnmiServer   *fgnmi.Server
	gnoiServer   *fgnoi.Server
	gribiServer  *fgribi.Server
	gnsiServer   *fgnsi.Server
	p4rtServer   *fp4rt.Server
	dplaneServer *dataplane.Dataplane

	// Configuration
	config *configpb.LemmingConfig

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
	sysribAddr     string
	faultAddr      string
	bgpPort        uint16
	dataplane      bool
	faultInject    bool
	dataplaneOpts  []dplaneopts.Option
	gribiOpts      []gribis.ServerOpt
	configFile     string
}

// resolveOpts applies all the options and returns a struct containing the result.
func resolveOpts(opts []Option) *opt {
	o := &opt{
		sysribAddr: "/tmp/sysrib.api",
	}
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

func WithDataplane(enable bool) Option {
	return func(o *opt) {
		o.dataplane = enable
	}
}

func WithDataplaneOpts(opts ...dplaneopts.Option) Option {
	return func(o *opt) {
		o.dataplaneOpts = opts
	}
}

// WithSysribAddr specifies a unix domain socket path for sysrib.
// Default: "/tmp/sysrib.api"
func WithSysribAddr(sysribAddr string) Option {
	return func(o *opt) {
		o.sysribAddr = sysribAddr
	}
}

// WithGRIBIOpts specifies the set of gRIBI options that should be passed to
// the gRIBI server that lemming runs.
func WithGRIBIOpts(opts ...gribis.ServerOpt) Option {
	return func(o *opt) {
		o.gribiOpts = opts
	}
}

// WithFaultInjection enables the fault injection service.
func WithFaultInjection(enable bool) Option {
	return func(o *opt) {
		o.faultInject = enable
	}
}

// WithFaultAddr sets the address of the fault service.
func WithFaultAddr(addr string) Option {
	return func(o *opt) {
		o.faultAddr = addr
	}
}

// WithConfigFile specifies a configuration file path for lemming device settings.
// The file has to be in protobuf text (.textproto/.pb.txt) format.
func WithConfigFile(configFile string) Option {
	return func(o *opt) {
		o.configFile = configFile
	}
}

// New returns a new initialized device.
func New(targetName, zapiURL string, opts ...Option) (*Device, error) {
	var dplane *dataplane.Dataplane
	var recs []reconciler.Reconciler

	resolvedOpts := resolveOpts(opts)

	// Load configuration in startup
	lemmingConfig, err := config.Load(resolvedOpts.configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	if resolvedOpts.dataplane {
		if runtime.GOOS != "linux" {
			return nil, fmt.Errorf("dataplane only supported on linux, GOOS is %s", runtime.GOOS)
		}

		log.Info("enabling dataplane")
		var err error
		dplane, err = dataplane.New(context.Background(), resolvedOpts.dataplaneOpts...)
		if err != nil {
			return nil, err
		}
		recs = append(recs, dplane)
	}

	log.Info("starting gNMI")

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
	streamInt := []grpc.StreamServerInterceptor{fgnmi.NewSubscribeTargetUpdateInterceptor(targetName)}
	unaryInt := []grpc.UnaryServerInterceptor{}

	creds := resolvedOpts.tlsCredentials
	if creds != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	var faultService *gRPCService

	if resolvedOpts.faultInject {
		faultInt := fault.NewInterceptor()
		l, err := net.Listen("tcp", resolvedOpts.faultAddr)
		if err != nil {
			return nil, err
		}
		srv := grpc.NewServer(grpc.Creds(creds))
		faultpb.RegisterFaultInjectServer(srv, faultInt)

		streamInt = append(streamInt, faultInt.Stream)
		unaryInt = append(unaryInt, faultInt.Unary)
		faultService = &gRPCService{
			lis:     l,
			s:       srv,
			stopped: make(chan struct{}),
		}
	}

	grpcOpts = append(grpcOpts, grpc.ChainStreamInterceptor(streamInt...), grpc.ChainUnaryInterceptor(unaryInt...))

	s := grpc.NewServer(grpcOpts...)

	recs = append(recs,
		fakedevice.NewSystemBaseTask(),
		fakedevice.NewBootTimeTask(lemmingConfig),
		fakedevice.NewCurrentTimeTask(),
		fakedevice.NewChassisComponentsTask(lemmingConfig),
		fakedevice.NewProcessMonitoringTask(lemmingConfig),
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
	if err := sysribServer.Start(context.Background(), cacheClient, targetName, zapiURL, resolvedOpts.sysribAddr); err != nil {
		return nil, fmt.Errorf("sysribServer failed to start: %v", err)
	}

	log.Info("starting gRIBI")
	// TODO(wenbli): Use gRIBIs once we change lemming's KNE config to use different ports.
	// gRIBIs := grpc.NewServer()
	gribiServer, err := fgribi.New(s, cacheClient, targetName, root, fmt.Sprintf("unix:%s", resolvedOpts.sysribAddr), resolvedOpts.gribiOpts...)
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

	gnoiServer, err := fgnoi.New(s, cacheClient, targetName, lemmingConfig)
	if err != nil {
		return nil, err
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
		faultService: faultService,
		gnmiServer:   gnmiServer,
		gnoiServer:   gnoiServer,
		gribiServer:  gribiServer,
		gnsiServer:   gnsiServer,
		p4rtServer:   fp4rt.New(P4RTs),
		dplaneServer: dplane,
		config:       lemmingConfig,
	}
	reflection.Register(s)
	d.startServer()

	if err := gnmiServer.StartReconcilers(context.Background()); err != nil {
		return nil, err
	}

	m := otel.GetMeterProvider().Meter("openconfig/lemming")
	c, err := m.Int64Counter("lemming-instance")
	if err != nil {
		return nil, err
	}
	c.Add(context.Background(), 1)

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

// Dataplane returns the dataplane server implementation.
func (d *Device) Dataplane() *dataplane.Dataplane {
	return d.dplaneServer
}

// Config returns the lemming configuration.
func (d *Device) Config() *configpb.LemmingConfig {
	return d.config
}

func (d *Device) startServer() {
	d.stopped = make(chan struct{})
	services := map[string]*gRPCService{
		"gNMI/gNOI/gNSI": d.gnmignoignsiService,
		"gRIBI":          d.gribiService,
		"P4RT":           d.p4rtService,
		"Fault":          d.faultService,
	}
	for svcName, svc := range services {
		if svc == nil {
			continue
		}
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
			if svc != nil {
				svc.s.Stop()
			}
		}
		for _, svc := range services {
			if svc != nil {
				<-svc.stopped
			}
		}
	}
}
