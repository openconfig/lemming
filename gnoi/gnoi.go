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
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openconfig/lemming/gnmi/oc"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
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

	// componentRebootsMu protects the componentReboots map
	// Map to track pending component reboots by component name
	componentRebootsMu sync.Mutex
	componentReboots   map[string]chan struct{}
}

func newSystem(c *ygnmi.Client) *system {
	return &system{
		c:                  c,
		cancelReboot:       make(chan struct{}, 1),
		cancelRebootFinish: make(chan struct{}),
		componentReboots:   make(map[string]chan struct{}),
	}
}

func (*system) Time(context.Context, *spb.TimeRequest) (*spb.TimeResponse, error) {
	return &spb.TimeResponse{Time: uint64(time.Now().UnixNano())}, nil
}

func (s *system) Reboot(ctx context.Context, r *spb.RebootRequest) (*spb.RebootResponse, error) {
	log.Infof("Received reboot request: %v", r)
	if r.Method == spb.RebootMethod_POWERUP {
		return &spb.RebootResponse{}, nil
	}

	// If subcomponents are specified, handle component-specific reboot
	if len(r.GetSubcomponents()) > 0 {
		return s.handleComponentReboot(ctx, r)
	}

	// Otherwise handle system-wide reboot
	if err := s.handleSystemReboot(ctx, r); err != nil {
		return nil, err
	}

	log.Infof("successful reboot with delay %v, type %v, and force %v", r.GetDelay(), r.GetMethod(), r.GetForce())
	return &spb.RebootResponse{}, nil
}

// handleComponentReboot processes a reboot request for specific components
func (s *system) handleComponentReboot(ctx context.Context, r *spb.RebootRequest) (*spb.RebootResponse, error) {
	// Check if there's a system-wide reboot pending, which would block all component reboots
	s.rebootMu.Lock()
	systemRebootPending := s.hasPendingReboot
	s.rebootMu.Unlock()

	if systemRebootPending {
		return nil, status.Errorf(codes.FailedPrecondition, "system-wide reboot already pending, cannot reboot components")
	}

	// Check for any pending component reboots
	s.componentRebootsMu.Lock()
	hasPendingComponentReboots := len(s.componentReboots) > 0
	s.componentRebootsMu.Unlock()

	// Process each subcomponent
	for _, subcompPath := range r.GetSubcomponents() {
		var componentName string
		if elem := subcompPath.GetElem()[0]; elem.GetName() == "component" {
			componentName = elem.GetKey()["name"]
		} else {
			componentName = elem.GetName()
		}

		// Check if this is an active supervisor
		isActive, err := s.IsActiveSupervisor(ctx, componentName)
		if err != nil {
			log.Warningf("Failed to determine supervisor role for %s: %v", componentName, err)
		} else if isActive {
			// Reject active supervisor reboot if there are pending reboots
			if hasPendingComponentReboots {
				return nil, status.Errorf(codes.FailedPrecondition,
					"cannot reboot active supervisor %s while component reboots are pending", componentName)
			}
			// If this is the first reboot request and it's for an active supervisor,
			// reject it to enforce standby-only policy
			return nil, status.Errorf(codes.FailedPrecondition, "rebooting active supervisor %s is not allowed, use standby or chassis reboot instead", componentName)
		}

		// Check if the component exists by querying it
		componentPath := ocpath.Root().Component(componentName)
		_, err = ygnmi.Get(ctx, s.c, componentPath.State())
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "component %q not found: %v", componentName, err)
		}

		// Check if there's already a pending reboot for this component
		s.componentRebootsMu.Lock()
		if _, exists := s.componentReboots[componentName]; exists {
			s.componentRebootsMu.Unlock()
			return nil, status.Errorf(codes.AlreadyExists, "reboot already pending for component %q", componentName)
		}

		// Create a cancellation channel for this component
		cancelCh := make(chan struct{}, 1)
		rebootCtx, cancel := context.WithCancel(ctx)
		s.componentReboots[componentName] = cancelCh
		s.componentRebootsMu.Unlock()

		// Cleanup function for consistent cleanup
		cleanup := func() {
			cancel()
			s.componentRebootsMu.Lock()
			delete(s.componentReboots, componentName)
			s.componentRebootsMu.Unlock()
		}

		delay := r.GetDelay()
		if delay == 0 {
			// Immediate reboot
			if err := fakedevice.RebootComponent(ctx, s.c, componentName, time.Now().UnixNano()); err != nil {
				cleanup()
				return nil, status.Errorf(codes.Internal, "failed to reboot component %q: %v", componentName, err)
			}

			// Remove from pending list after reboot is complete
			cleanup()
			log.Infof("Component %q immediate reboot completed", componentName)
		} else {
			// Handle delayed reboot
			go func(compName string) {
				defer cleanup()
				select {
				case <-cancelCh:
					log.Infof("delayed component reboot for %q cancelled", compName)
				case <-rebootCtx.Done():
					log.Infof("delayed component reboot for %q cancelled due to context", compName)
				case <-time.After(time.Duration(delay) * time.Nanosecond):
					now := time.Now().UnixNano()
					if err := fakedevice.RebootComponent(rebootCtx, s.c, compName, now); err != nil {
						log.Errorf("delayed component reboot for %q failed: %v", compName, err)
						return
					}
					log.Infof("Component %q delayed reboot completed", compName)
				}
			}(componentName)

			log.Infof("scheduled component reboot for %q with delay %v", componentName, r.GetDelay())
		}
	}
	return &spb.RebootResponse{}, nil
}

// handleSystemReboot processes a reboot request for chassis
func (s *system) handleSystemReboot(ctx context.Context, r *spb.RebootRequest) error {
	s.rebootMu.Lock()
	defer s.rebootMu.Unlock()
	if s.hasPendingReboot {
		return status.Errorf(codes.AlreadyExists, "reboot already pending")
	}

	delay := r.GetDelay()
	if delay == 0 {
		now := time.Now().UnixNano()
		if err := fakedevice.Reboot(ctx, s.c, now); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		return nil
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
	return nil
}

func (s *system) CancelReboot(ctx context.Context, c *spb.CancelRebootRequest) (*spb.CancelRebootResponse, error) {
	log.Infof("Received cancel reboot request %v", c)

	// Check if there are any component reboots to cancel
	s.componentRebootsMu.Lock()
	componentRebootsPending := len(s.componentReboots)
	for component, cancelCh := range s.componentReboots {
		select {
		case cancelCh <- struct{}{}:
			log.Infof("Sent cancellation signal for component %q reboot", component)
		default:
			// Channel already has a message or is closed
			log.Infof("Component %q reboot already completed or in progress, couldn't cancel", component)
		}
		delete(s.componentReboots, component)
	}
	s.componentRebootsMu.Unlock()

	// Check for system-wide reboot to cancel
	s.rebootMu.Lock()
	hasPendingReboot := s.hasPendingReboot
	s.rebootMu.Unlock()

	if !hasPendingReboot && componentRebootsPending == 0 {
		// No reboots of any kind to cancel
		return &spb.CancelRebootResponse{}, nil
	}

	if hasPendingReboot {
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

	return &spb.CancelRebootResponse{}, nil
}

// IsActiveSupervisor checks for the redundant role of a supervisor
func (s *system) IsActiveSupervisor(ctx context.Context, componentName string) (bool, error) {
	componentPath := ocpath.Root().Component(componentName)
	roleVal, err := ygnmi.Get(ctx, s.c, componentPath.RedundantRole().State())
	if err != nil {
		return false, err
	}
	return roleVal == oc.PlatformTypes_ComponentRedundantRole_PRIMARY, nil
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
