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
	"time"

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
	"google.golang.org/grpc"
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
}

func (*system) Time(context.Context, *spb.TimeRequest) (*spb.TimeResponse, error) {
	return &spb.TimeResponse{Time: uint64(time.Now().UnixNano())}, nil
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

func New(s *grpc.Server) *Server {
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
		systemServer:           &system{},
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
	return srv
}
