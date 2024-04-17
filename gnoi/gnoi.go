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

package gnoi

import (
	"context"
	"sync"
	"time"

	log "github.com/golang/glog"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	bpb "github.com/openconfig/gnoi/bgp"
	cmpb "github.com/openconfig/gnoi/cert"
	diagpb "github.com/openconfig/gnoi/diag"
	frpb "github.com/openconfig/gnoi/factory_reset"
	fpb "github.com/openconfig/gnoi/file"
	hpb "github.com/openconfig/gnoi/healthz"
	lpb "github.com/openconfig/gnoi/layer2"
	mpb "github.com/openconfig/gnoi/mpls"
	ospb "github.com/openconfig/gnoi/os"
	otpb "github.com/openconfig/gnoi/otdr"
	spb "github.com/openconfig/gnoi/system"
	wrpb "github.com/openconfig/gnoi/wavelength_router"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type bgp struct {
	bpb.UnimplementedBGPServer
}

type cert struct {
	cmpb.UnimplementedCertificateManagementServer
}

type diag struct {
	diagpb.UnimplementedDiagServer
}

type factoryReset struct {
	frpb.UnimplementedFactoryResetServer
}

type file struct {
	fpb.UnimplementedFileServer
}

type healthz struct {
	hpb.UnimplementedHealthzServer
}

type layer2 struct {
	lpb.UnimplementedLayer2Server
}

type mpls struct {
	mpb.UnimplementedMPLSServer
}

type os struct {
	ospb.UnimplementedOSServer
}

type otdr struct {
	otpb.UnimplementedOTDRServer
}

type system struct {
	spb.UnimplementedSystemServer

	c *ygnmi.Client

	// rebootMu has the following roles:
	// * ensures that writes to hasPendingReboot are free from race
	//   conditions
	// * ensures consistency between reboot operations and the current
	//   state of hasPendingReboot (i.e. prevent TOCCTOU race conditions).
	rebootMu         sync.Mutex
	hasPendingReboot bool
	// These channels ensure that cancellation is a blocking operation to
	// avoid future reboots from conflicting with cancelled pending
	// reboots.
	cancelReboot       chan struct{}
	cancelRebootFinish chan struct{}
}

func newSystem(c *ygnmi.Client) *system {
	return &system{
		c:                  c,
		cancelReboot:       make(chan struct{}, 1),
		cancelRebootFinish: make(chan struct{}),
	}
}

func (*system) Time(context.Context, *spb.TimeRequest) (*spb.TimeResponse, error) {
	return &spb.TimeResponse{Time: uint64(time.Now().UnixNano())}, nil
}

func (s *system) Reboot(ctx context.Context, r *spb.RebootRequest) (*spb.RebootResponse, error) {
	if r.Method == spb.RebootMethod_POWERUP {
		return &spb.RebootResponse{}, nil
	}

	s.rebootMu.Lock()
	defer s.rebootMu.Unlock()
	if s.hasPendingReboot {
		return nil, status.Errorf(codes.AlreadyExists, "reboot already pending")
	}

	delay := r.GetDelay()
	if delay == 0 {
		now := time.Now().UnixNano()
		if err := fakedevice.Reboot(ctx, s.c, now); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &spb.RebootResponse{}, nil
	}

	s.hasPendingReboot = true
	go func() { // wait the delay time for reboot
		select {
		case <-s.cancelReboot:
			log.Infof("delayed reboot cancelled")
			s.cancelRebootFinish <- struct{}{}
		case <-time.After(time.Duration(delay) * time.Nanosecond):
			now := time.Now().UnixNano()
			if err := fakedevice.Reboot(ctx, s.c, now); err != nil {
				log.Errorf("delayed reboot failed: %v", err)
			}
			s.rebootMu.Lock()
			defer s.rebootMu.Unlock()
			s.hasPendingReboot = false
		}
	}()

	return &spb.RebootResponse{}, nil
}

func (s *system) CancelReboot(context.Context, *spb.CancelRebootRequest) (*spb.CancelRebootResponse, error) {
	s.rebootMu.Lock()
	hasPendingReboot := s.hasPendingReboot
	s.rebootMu.Unlock()
	if !hasPendingReboot {
		return &spb.CancelRebootResponse{}, nil
	}

	s.cancelReboot <- struct{}{} // signal cancellation
	for {
		select {
		case <-s.cancelRebootFinish:
			s.rebootMu.Lock()
			defer s.rebootMu.Unlock()
			s.hasPendingReboot = false
			return &spb.CancelRebootResponse{}, nil
		case <-time.After(time.Second): // It's possible for reboot to happen after cancellation signal -- use polling to check that.
			s.rebootMu.Lock()
			if !s.hasPendingReboot {
				s.rebootMu.Unlock()
				<-s.cancelReboot // clean-up cancellation signal that's not needed since reboot actually happened.
				return &spb.CancelRebootResponse{}, nil
			}
			s.rebootMu.Unlock()
		}
	}
}

type wavelengthRouter struct {
	wrpb.UnimplementedWavelengthRouterServer
}

type Server struct {
	s                      *grpc.Server
	bgpServer              *bgp
	certServer             *cert
	diagServer             *diag
	fileServer             *file
	resetServer            *factoryReset
	healthzServer          *healthz
	layer2Server           *layer2
	mplsServer             *mpls
	osServer               *os
	otdrServer             *otdr
	systemServer           *system
	wavelengthRouterServer *wavelengthRouter
}

func New(s *grpc.Server, gClient gpb.GNMIClient, target string) (*Server, error) {
	yclient, err := ygnmi.NewClient(gClient, ygnmi.WithTarget(target), ygnmi.WithRequestLogLevel(2))
	if err != nil {
		return nil, err
	}

	srv := &Server{
		s:                      s,
		bgpServer:              &bgp{},
		certServer:             &cert{},
		diagServer:             &diag{},
		fileServer:             &file{},
		resetServer:            &factoryReset{},
		healthzServer:          &healthz{},
		layer2Server:           &layer2{},
		mplsServer:             &mpls{},
		osServer:               &os{},
		otdrServer:             &otdr{},
		systemServer:           newSystem(yclient),
		wavelengthRouterServer: &wavelengthRouter{},
	}
	bpb.RegisterBGPServer(s, srv.bgpServer)
	cmpb.RegisterCertificateManagementServer(s, srv.certServer)
	diagpb.RegisterDiagServer(s, srv.diagServer)
	fpb.RegisterFileServer(s, srv.fileServer)
	frpb.RegisterFactoryResetServer(s, srv.resetServer)
	hpb.RegisterHealthzServer(s, srv.healthzServer)
	lpb.RegisterLayer2Server(s, srv.layer2Server)
	mpb.RegisterMPLSServer(s, srv.mplsServer)
	ospb.RegisterOSServer(s, srv.osServer)
	otpb.RegisterOTDRServer(s, srv.otdrServer)
	spb.RegisterSystemServer(s, srv.systemServer)
	wrpb.RegisterWavelengthRouterServer(s, srv.wavelengthRouterServer)
	return srv, nil
}
