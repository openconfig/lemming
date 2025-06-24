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
	"github.com/openconfig/ygnmi/ygnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"

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
	pb "github.com/openconfig/gnoi/types"
	wrpb "github.com/openconfig/gnoi/wavelength_router"
)

const (
	supervisor1Name = "Supervisor1"
	supervisor2Name = "Supervisor2"

	// Kill process default
	defaultRestart = true
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
	// switchoverMu protects switchover operations and ensures
	// only one switchover can be in progress at a time
	switchoverMu         sync.Mutex
	hasPendingSwitchover bool
	// processMu protects process operations and ensures
	// only one process operation can be in progress at a time
	processMu sync.Mutex
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

	// Process each subcomponent
	for _, subcompPath := range r.GetSubcomponents() {
		componentName, err := extractComponentNameFromPath(subcompPath)
		if err != nil {
			return nil, err
		}

		// Check if this is an active supervisor
		isActive, err := s.isActiveSupervisor(ctx, componentName)
		if err != nil {
			log.Warningf("Failed to determine supervisor role for %s: %v", componentName, err)
		} else if isActive {
			// reject active control card reboot to enforce standby-only policy
			return nil, status.Errorf(codes.FailedPrecondition, "rebooting active supervisor %s is not allowed, use standby or chassis reboot instead", componentName)
		}

		// Check if the component exists by querying it
		componentPath := ocpath.Root().Component(componentName)
		_, err = ygnmi.Get(ctx, s.c, componentPath.State())
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "component %q not found: %v", componentName, err)
		}

		delay := r.GetDelay()
		if delay == 0 {
			s.componentRebootsMu.Lock()
			if _, exists := s.componentReboots[componentName]; exists {
				s.componentRebootsMu.Unlock()
				return nil, status.Errorf(codes.AlreadyExists, "reboot already pending for component %q", componentName)
			}
			s.componentRebootsMu.Unlock()
			// Immediate reboot
			if err := fakedevice.RebootComponent(context.Background(), s.c, componentName, time.Now().UnixNano()); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to reboot component %q: %v", componentName, err)
			}
			log.Infof("Component %q immediate reboot completed", componentName)
			continue
		}
		// Check if there's already a pending reboot for this component
		s.componentRebootsMu.Lock()
		if _, exists := s.componentReboots[componentName]; exists {
			s.componentRebootsMu.Unlock()
			return nil, status.Errorf(codes.AlreadyExists, "reboot already pending for component %q", componentName)
		}

		// Create a cancellation channel for this component
		cancelCh := make(chan struct{}, 1)
		rebootCtx, cancel := context.WithCancel(context.Background())
		s.componentReboots[componentName] = cancelCh
		s.componentRebootsMu.Unlock()

		// Cleanup function for consistent cleanup
		cleanup := func() {
			cancel()
			s.componentRebootsMu.Lock()
			delete(s.componentReboots, componentName)
			s.componentRebootsMu.Unlock()
		}

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

// SwitchControlProcessor performs supervisor switchover from the current active supervisor to the specified target supervisor
func (s *system) SwitchControlProcessor(ctx context.Context, r *spb.SwitchControlProcessorRequest) (*spb.SwitchControlProcessorResponse, error) {
	log.Infof("Received SwitchControlProcessor request: %v", r)

	if r.GetControlProcessor() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "control_processor path is required")
	}

	targetSupervisor, err := extractComponentNameFromPath(r.GetControlProcessor())
	if err != nil {
		return nil, err
	}

	// Protect against concurrent switchover operations
	s.switchoverMu.Lock()
	defer s.switchoverMu.Unlock()

	if s.hasPendingSwitchover {
		return nil, status.Errorf(codes.FailedPrecondition, "supervisor switchover already in progress")
	}

	// Check if there are any pending reboot operations (system or component level)
	s.rebootMu.Lock()
	systemRebootPending := s.hasPendingReboot
	s.rebootMu.Unlock()

	s.componentRebootsMu.Lock()
	componentRebootsPending := len(s.componentReboots)
	s.componentRebootsMu.Unlock()

	if componentRebootsPending > 0 || systemRebootPending {
		return nil, status.Errorf(codes.FailedPrecondition, "reboot operations pending, cannot perform switchover")
	}

	// Validate supervisor state and get active/standby supervisors
	activeSupervisor, standbySupervisor, err := s.getSupervisorRole(ctx)
	if err != nil {
		return nil, err
	}

	if targetSupervisor != activeSupervisor && targetSupervisor != standbySupervisor {
		return nil, status.Errorf(codes.NotFound, "target supervisor %q does not exist", targetSupervisor)
	}

	// Check if target is already the active supervisor (no-op case)
	if targetSupervisor == activeSupervisor {
		log.Infof("Target supervisor %q is already active, returning current state (no-op)", targetSupervisor)

		componentPath := ocpath.Root().Component(targetSupervisor)
		component, err := ygnmi.Get(ctx, s.c, componentPath.State())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get current active supervisor info: %v", err)
		}

		// Return successful response for no-op case
		return &spb.SwitchControlProcessorResponse{
			ControlProcessor: r.GetControlProcessor(),
			Version:          component.GetSoftwareVersion(),
			Uptime:           0,
		}, nil
	}

	s.hasPendingSwitchover = true

	// Get the target supervisor info for response
	componentPath := ocpath.Root().Component(targetSupervisor)
	component, err := ygnmi.Get(ctx, s.c, componentPath.State())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get target supervisor info: %v", err)
	}

	response := &spb.SwitchControlProcessorResponse{
		ControlProcessor: r.GetControlProcessor(),
		Version:          component.GetSoftwareVersion(),
		Uptime:           0,
	}

	log.Infof("Scheduled supervisor switcover from %s to %s", activeSupervisor, targetSupervisor)
	go func() {
		backgroundctx := context.Background()

		defer func() {
			s.switchoverMu.Lock()
			s.hasPendingSwitchover = false
			s.switchoverMu.Unlock()
		}()

		// Small delay to make sure response is sent
		time.Sleep(100 * time.Millisecond)

		switchoverTime := time.Now().UnixNano()
		err := fakedevice.SwitchoverSupervisor(backgroundctx, s.c, targetSupervisor, activeSupervisor, switchoverTime)
		if err != nil {
			log.Errorf("Background supervisor switchover failed: %v", err)
		}
	}()

	return response, nil
}

// KillProcess simulates process termination and restart functionality
func (s *system) KillProcess(ctx context.Context, r *spb.KillProcessRequest) (*spb.KillProcessResponse, error) {
	log.Infof("Received kill process request: %v", r)

	if r.GetPid() == 0 && r.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "either pid or name must be specified")
	}

	signal := r.GetSignal()
	if signal == spb.KillProcessRequest_SIGNAL_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "signal must be specified")
	}

	targetPID, processName, err := s.resolvePIDAndName(ctx, r.GetPid(), r.GetName())
	if err != nil {
		return nil, err
	}

	// HUP is for reload, restart should be false by default
	restart := defaultRestart
	if signal == spb.KillProcessRequest_SIGNAL_HUP {
		restart = false
	}

	// Protect against concurrent process operations
	s.processMu.Lock()
	defer s.processMu.Unlock()

	if err := fakedevice.KillProcess(context.Background(), s.c, targetPID, processName, signal, restart); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to kill process: %v", err)
	}

	return &spb.KillProcessResponse{}, nil
}

// resolvePIDAndName resolves either PID or name to a validated PID and process name that exists in the system
func (s *system) resolvePIDAndName(ctx context.Context, pid uint32, name string) (uint32, string, error) {
	var targetPID uint32
	var targetProcessName string

	if pid != 0 {
		// PID provided - get process info to extract name
		process, err := ygnmi.Get(ctx, s.c, ocpath.Root().System().Process(uint64(pid)).State())
		if err != nil {
			return 0, "", status.Errorf(codes.NotFound, "PID %d not found in process path: %v", pid, err)
		}
		targetPID = pid
		targetProcessName = process.GetName()
	} else {
		// Name provided - look up PID and process name from process monitoring system
		processes, err := ygnmi.GetAll(ctx, s.c, ocpath.Root().System().ProcessAny().State())
		if err != nil {
			return 0, "", status.Errorf(codes.Internal, "failed to query processes: %v", err)
		}
		var foundPID uint64
		var found bool
		for _, process := range processes {
			if process.Name != nil && *process.Name == name {
				if process.Pid != nil {
					foundPID = *process.Pid
					found = true
					break
				}
			}
		}
		if !found {
			return 0, "", status.Errorf(codes.NotFound, "process %q not found", name)
		}
		targetPID = uint32(foundPID)
		targetProcessName = name
	}

	return targetPID, targetProcessName, nil
}

// extractComponentNameFromPath extracts the component name from the gNMI path
func extractComponentNameFromPath(path *pb.Path) (string, error) {
	elems := path.GetElem()
	// Handle Arista format
	if len(elems) == 1 {
		componentName := elems[0].GetName()
		if componentName == "" {
			return "", status.Errorf(codes.InvalidArgument, "Invalid component path, element name is empty, got: %v", path)
		}
		return componentName, nil
	}
	if len(elems) == 2 &&
		elems[0].GetName() == "components" &&
		elems[1].GetName() == "component" &&
		elems[1].GetKey()["name"] != "" {
		return elems[1].GetKey()["name"], nil
	}
	return "", status.Errorf(codes.InvalidArgument,
		"invalid component path, expected either single element or OpenConfig format (/componets/component[name=...]), got: %v", path)
}

// isActiveSupervisor checks for the redundant role of a supervisor
func (s *system) isActiveSupervisor(ctx context.Context, componentName string) (bool, error) {
	componentPath := ocpath.Root().Component(componentName)
	roleVal, err := ygnmi.Get(ctx, s.c, componentPath.RedundantRole().State())
	if err != nil {
		return false, err
	}
	return roleVal == oc.PlatformTypes_ComponentRedundantRole_PRIMARY, nil
}

// getSupervisorRole validates and returns the active and standby supervisors
func (s *system) getSupervisorRole(ctx context.Context) (activeSupervisor, standbySupervisor string, err error) {
	// Check if Supervisor1 is active
	supervisor1Active, err := s.isActiveSupervisor(ctx, supervisor1Name)
	if err != nil {
		return "", "", status.Errorf(codes.Internal, "failed to check supervisor %q state: %v", supervisor1Name, err)
	}

	// Check if Supervisor2 is active
	supervisor2Active, err := s.isActiveSupervisor(ctx, supervisor2Name)
	if err != nil {
		return "", "", status.Errorf(codes.Internal, "failed to check supervisor %q state: %v", supervisor2Name, err)
	}

	// Determine active and standby based on the results
	switch {
	case supervisor1Active && !supervisor2Active:
		return supervisor1Name, supervisor2Name, nil
	case supervisor2Active && !supervisor1Active:
		return supervisor2Name, supervisor1Name, nil
	case supervisor1Active && supervisor2Active:
		return "", "", status.Errorf(codes.FailedPrecondition, "both supervisors are active")
	default:
		return "", "", status.Errorf(codes.FailedPrecondition, "no active supervisor found")
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
