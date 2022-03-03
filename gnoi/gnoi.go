package gnoi

import (
	bpb "github.com/openconfig/gnoi/bgp"
	cmpb "github.com/openconfig/gnoi/cert"
	dpb "github.com/openconfig/gnoi/diag"
	frpb "github.com/openconfig/gnoi/factory_reset"
	fpb "github.com/openconfig/gnoi/file"
	hpb "github.com/openconfig/gnoi/healthz"
	ipb "github.com/openconfig/gnoi/interface"
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
	dpb.UnimplementedDiagServer
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

type interfac struct {
	ipb.UnimplementedInterfaceServer
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
	interfaceServer        *interfac
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
		interfaceServer:        &interfac{},
		layer2Server:           &layer2{},
		mplsServer:             &mpls{},
		osServer:               &os{},
		otdrServer:             &otdr{},
		systemServer:           &system{},
		wavelengthRouterServer: &wavelengthRouter{},
	}
	bpb.RegisterBGPServer(s, srv.bgpServer)
	cmpb.RegisterCertificateManagementServer(s, srv.certServer)
	dpb.RegisterDiagServer(s, srv.diagServer)
	fpb.RegisterFileServer(s, srv.fileServer)
	frpb.RegisterFactoryResetServer(s, srv.resetServer)
	hpb.RegisterHealthzServer(s, srv.healthzServer)
	ipb.RegisterInterfaceServer(s, srv.interfaceServer)
	lpb.RegisterLayer2Server(s, srv.layer2Server)
	mpb.RegisterMPLSServer(s, srv.mplsServer)
	ospb.RegisterOSServer(s, srv.osServer)
	otpb.RegisterOTDRServer(s, srv.otdrServer)
	spb.RegisterSystemServer(s, srv.systemServer)
	wrpb.RegisterWavelengthRouterServer(s, srv.wavelengthRouterServer)
	return srv
}
