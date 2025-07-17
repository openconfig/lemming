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
	"fmt"
	"maps"
	"math"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/ygnmi/ygnmi"
	rpc "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	configpb "github.com/openconfig/lemming/proto/config"

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
	plqpb "github.com/openconfig/gnoi/packet_link_qualification"
	spb "github.com/openconfig/gnoi/system"
	pb "github.com/openconfig/gnoi/types"
	wrpb "github.com/openconfig/gnoi/wavelength_router"
)

const (
	// Kill process default
	defaultRestart = true

	// Ping simulation default values
	defaultPingCount    = 5
	defaultPingInterval = 1000000000
	defaultPingWait     = 2000000000
	defaultPingSize     = 56
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

	c      *ygnmi.Client
	config *configpb.Config

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

func newSystem(c *ygnmi.Client, config *configpb.Config) *system {
	return &system{
		c:                  c,
		config:             config,
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
			if err := fakedevice.RebootComponent(context.Background(), s.c, componentName, time.Now().UnixNano(), s.config); err != nil {
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
				if err := fakedevice.RebootComponent(rebootCtx, s.c, compName, now, s.config); err != nil {
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
		err := fakedevice.SwitchoverSupervisor(backgroundctx, s.c, targetSupervisor, activeSupervisor, switchoverTime, s.config)
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

	if err := fakedevice.KillProcess(context.Background(), s.c, targetPID, processName, signal, restart, s.config); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to kill process: %v", err)
	}

	return &spb.KillProcessResponse{}, nil
}

// Ping simulates ICMP ping operations with configurable network conditions
func (s *system) Ping(r *spb.PingRequest, stream spb.System_PingServer) error {
	log.Infof("Received ping request: %v", r)

	if r.GetDestination() == "" {
		return status.Errorf(codes.InvalidArgument, "destination address is required")
	}

	ctx := stream.Context()
	destination := r.GetDestination()

	count := r.GetCount()
	if count == 0 {
		count = defaultPingCount
	} else if count < -1 {
		return status.Errorf(codes.InvalidArgument, "count must be >= -1, got %d", count)
	}

	interval := r.GetInterval()
	switch {
	case interval == 0:
		interval = defaultPingInterval
	case interval == -1:
		// Flood ping - 1ms minimum interval for safety
		interval = 1000000
	case interval < -1:
		return status.Errorf(codes.InvalidArgument, "interval must be >= -1, got %d", interval)
	}

	wait := r.GetWait()
	if wait == 0 {
		wait = defaultPingWait
	} else if wait < 0 {
		return status.Errorf(codes.InvalidArgument, "wait must be >= 0, got %d", wait)
	}

	size := r.GetSize()
	if size == 0 {
		size = defaultPingSize
	} else if size < 8 || size > 65507 {
		return status.Errorf(codes.InvalidArgument, "packet size must be between 8 and 65507 bytes, got %d", size)
	}

	// TODO: Add support for do_not_fragment, do_not_resolve, l3protocol, network_instance parameters

	// Calculate appropriate buffer size based on interval and count
	bufferSize := 100
	if count > 0 && count < 1000 {
		bufferSize = int(count) + 10
	} else if interval < 10000000 {
		bufferSize = 1000 // Larger buffer for flood ping
	}

	responseChan := make(chan *fakedevice.PingPacketResult, bufferSize)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		if err := fakedevice.PingSimulation(ctx, destination, count, time.Duration(interval), time.Duration(wait), uint32(size), responseChan, s.config); err != nil {
			log.Errorf("Ping simulation error: %v", err)
			select {
			case errorChan <- err:
			default:
			}
		}
		close(errorChan)
	}()

	startTime := time.Now()
	var totalSent, totalReceived int32
	var minTime, maxTime, totalTime time.Duration
	var rtts []time.Duration
	firstPacket := true
	const maxRTTSamples = 10000

	// Process results and stream responses
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errorChan:
			if err != nil {
				return status.Errorf(codes.Internal, "ping simulation failed: %v", err)
			}
		case result, ok := <-responseChan:
			if !ok {
				// Channel closed - check for any remaining errors
				select {
				case err := <-errorChan:
					if err != nil {
						return status.Errorf(codes.Internal, "ping simulation failed: %v", err)
					}
				default:
				}

				// Send summary and finish
				summary := &spb.PingResponse{
					Source:   destination,
					Time:     time.Since(startTime).Nanoseconds(),
					Sent:     totalSent,
					Received: totalReceived,
				}
				if totalReceived > 0 {
					summary.MinTime = minTime.Nanoseconds()
					summary.AvgTime = (totalTime / time.Duration(totalReceived)).Nanoseconds()
					summary.MaxTime = maxTime.Nanoseconds()

					// Calculate standard deviation from collected RTTs
					if len(rtts) < 2 {
						summary.StdDev = 0
					} else {
						var sumSquaredDiff float64
						avgTime := float64(summary.AvgTime)
						for _, rtt := range rtts {
							diff := float64(rtt.Nanoseconds()) - avgTime
							sumSquaredDiff += diff * diff
						}
						variance := sumSquaredDiff / float64(len(rtts)-1)
						summary.StdDev = int64(math.Round(math.Sqrt(variance)))
					}
				}
				return stream.Send(summary)
			}

			// Update stats
			totalSent++
			if result.Success {
				totalReceived++
				totalTime += result.RTT
				if firstPacket {
					minTime = result.RTT
					maxTime = result.RTT
					firstPacket = false
				} else {
					if result.RTT < minTime {
						minTime = result.RTT
					}
					if result.RTT > maxTime {
						maxTime = result.RTT
					}
				}
				// Only collect RTT samples up to the limit
				if len(rtts) < maxRTTSamples {
					rtts = append(rtts, result.RTT)
				}
			}

			// Send individual response
			response := &spb.PingResponse{
				Source:   destination,
				Time:     result.RTT.Nanoseconds(),
				Bytes:    int32(result.Bytes),
				Sequence: result.Sequence,
				Ttl:      result.TTL,
			}

			if err := stream.Send(response); err != nil {
				log.Errorf("Failed to send ping response: %v", err)
				return err
			}
		}
	}
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
	supervisor1Name := s.config.GetComponents().GetSupervisor1Name()
	supervisor2Name := s.config.GetComponents().GetSupervisor2Name()

	// Check if supervisor1 is active
	supervisor1Active, err := s.isActiveSupervisor(ctx, supervisor1Name)
	if err != nil {
		return "", "", status.Errorf(codes.Internal, "failed to check supervisor %q state: %v", supervisor1Name, err)
	}

	// Check if supervisor2 is active
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

type linkQualification struct {
	plqpb.UnimplementedLinkQualificationServer

	c      *ygnmi.Client
	config *configpb.Config

	mu sync.RWMutex
	// Qualification state tracking
	// qual_id -> state
	qualifications map[string]*QualificationState
	// Historical results tracking per interface
	// interface -> historical results
	historicalResults map[string][]*plqpb.QualificationResult
	maxHistorical     uint64
}

// QualificationState represents the state of a single qualification
type QualificationState struct {
	ID            string
	InterfaceName string
	State         plqpb.QualificationState
	StartTime     time.Time
	EndTime       time.Time

	// Packet statistics - use atomic operations for thread safety
	packetsSent     atomic.Uint64
	packetsReceived atomic.Uint64
	packetsDropped  atomic.Uint64
	packetsError    atomic.Uint64

	// Configuration
	IsGenerator bool
	IsReflector bool
	Config      *plqpb.QualificationConfiguration

	// Control channels for cancellation
	cancelCh chan struct{}
	done     bool
	// Protect individual state updates
	mu sync.Mutex
}

func newLinkQualification(c *ygnmi.Client, config *configpb.Config) *linkQualification {
	return &linkQualification{
		c:                 c,
		config:            config,
		qualifications:    make(map[string]*QualificationState),
		historicalResults: make(map[string][]*plqpb.QualificationResult),
		maxHistorical:     uint64(config.GetLinkQualification().GetMaxHistoricalResults()),
	}
}

// Capabilities returns the capabilities of the LinkQualification service
func (lq *linkQualification) Capabilities(ctx context.Context, req *plqpb.CapabilitiesRequest) (*plqpb.CapabilitiesResponse, error) {
	log.Infof("Received LinkQualification Capabilities request")

	return &plqpb.CapabilitiesResponse{
		Time:      timestamppb.Now(),
		NtpSynced: true,
		Generator: &plqpb.GeneratorCapabilities{
			PacketGenerator: &plqpb.PacketGeneratorCapabilities{
				MaxBps:              lq.config.GetLinkQualification().GetMaxBps(),
				MaxPps:              lq.config.GetLinkQualification().GetMaxPps(),
				MinMtu:              lq.config.GetLinkQualification().GetMinMtu(),
				MaxMtu:              lq.config.GetLinkQualification().GetMaxMtu(),
				MinSetupDuration:    durationpb.New(time.Duration(lq.config.GetLinkQualification().GetMinSetupDurationMs()) * time.Millisecond),
				MinTeardownDuration: durationpb.New(time.Duration(lq.config.GetLinkQualification().GetMinTeardownDurationMs()) * time.Millisecond),
				MinSampleInterval:   durationpb.New(time.Duration(lq.config.GetLinkQualification().GetMinSampleIntervalMs()) * time.Millisecond),
			},
			// PacketInjector intentionally omitted - unimplemented in simulation
		},
		Reflector: &plqpb.ReflectorCapabilities{
			AsicLoopback: &plqpb.AsicLoopbackCapabilities{
				MinSetupDuration:    durationpb.New(time.Duration(lq.config.GetLinkQualification().GetMinSetupDurationMs()) * time.Millisecond),
				MinTeardownDuration: durationpb.New(time.Duration(lq.config.GetLinkQualification().GetMinTeardownDurationMs()) * time.Millisecond),
				Fields:              []plqpb.HeaderMatchField{plqpb.HeaderMatchField_HEADER_MATCH_FIELD_L2},
			},
			PmdLoopback: &plqpb.PmdLoopbackCapabilities{
				MinSetupDuration:    durationpb.New(time.Duration(lq.config.GetLinkQualification().GetMinSetupDurationMs()) * time.Millisecond),
				MinTeardownDuration: durationpb.New(time.Duration(lq.config.GetLinkQualification().GetMinTeardownDurationMs()) * time.Millisecond),
			},
		},
		MaxHistoricalResultsPerInterface: lq.maxHistorical,
	}, nil
}

// Create starts link qualification on specified interfaces with multi-port support
func (lq *linkQualification) Create(ctx context.Context, req *plqpb.CreateRequest) (*plqpb.CreateResponse, error) {
	log.Infof("Received LinkQualification Create request with %d interfaces", len(req.GetInterfaces()))

	if len(req.GetInterfaces()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no interfaces specified")
	}

	lq.mu.Lock()
	defer lq.mu.Unlock()

	// Batch validation - validate ALL interfaces before starting ANY
	validationErrors := make(map[string]*rpc.Status)
	qualificationStates := make([]*QualificationState, 0, len(req.GetInterfaces()))

	// Track IDs and interfaces within this request to detect duplicates
	requestIDs := make(map[string]bool)
	requestInterfaces := make(map[string]bool)

	for _, config := range req.GetInterfaces() {
		if config.GetId() == "" {
			return nil, status.Errorf(codes.InvalidArgument, "qualification id is required")
		}
		if config.GetInterfaceName() == "" {
			return nil, status.Errorf(codes.InvalidArgument, "interface name is required")
		}

		id := config.GetId()
		interfaceName := config.GetInterfaceName()

		// Check for duplicates within this request first
		if requestIDs[id] {
			validationErrors[id] = &rpc.Status{
				Code:    int32(codes.AlreadyExists),
				Message: fmt.Sprintf("duplicate qualification id %s in request", id),
			}
			continue
		}
		if requestInterfaces[interfaceName] {
			validationErrors[id] = &rpc.Status{
				Code:    int32(codes.AlreadyExists),
				Message: fmt.Sprintf("duplicate interface %s in request", interfaceName),
			}
			continue
		}

		requestIDs[id] = true
		requestInterfaces[interfaceName] = true

		// Validate against existing state and configuration
		if err := lq.validateQualificationConfig(config); err != nil {
			if grpcErr, ok := status.FromError(err); ok {
				validationErrors[id] = &rpc.Status{
					Code:    int32(grpcErr.Code()),
					Message: grpcErr.Message(),
				}
			} else {
				validationErrors[id] = &rpc.Status{
					Code:    int32(codes.Internal),
					Message: err.Error(),
				}
			}
		} else {
			// Create qualification state for valid configs
			qualState := lq.createQualificationState(config)
			qualificationStates = append(qualificationStates, qualState)
		}
	}

	response := &plqpb.CreateResponse{
		Status: make(map[string]*rpc.Status),
	}
	maps.Copy(response.Status, validationErrors)

	// Proceed with creating qualifications if we have valid ones
	if len(qualificationStates) > 0 {
		for _, qualState := range qualificationStates {
			lq.qualifications[qualState.ID] = qualState
			if _, hasError := validationErrors[qualState.ID]; !hasError {
				response.Status[qualState.ID] = &rpc.Status{
					Code:    int32(codes.OK),
					Message: "qualification created successfully",
				}
			}
		}

		// Start individual qualifications
		for _, qualState := range qualificationStates {
			go lq.executeQualification(context.Background(), qualState)
		}
	}

	log.Infof("Created %d valid qualifications, %d validation errors",
		len(qualificationStates), len(validationErrors))

	return response, nil
}

// Get returns the status for the provided qualification ids
func (lq *linkQualification) Get(ctx context.Context, req *plqpb.GetRequest) (*plqpb.GetResponse, error) {
	log.Infof("Received LinkQualification Get request: %v", req.GetIds())

	if len(req.GetIds()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no qualification ids specified")
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	response := &plqpb.GetResponse{
		Results: make(map[string]*plqpb.QualificationResult),
	}

	lq.mu.RLock()
	defer lq.mu.RUnlock()

	for _, id := range req.GetIds() {
		qualState, exists := lq.qualifications[id]
		if !exists {
			// Return NOT_FOUND result for missing qualifications
			response.Results[id] = &plqpb.QualificationResult{
				Id: id,
				Status: &rpc.Status{
					Code:    int32(codes.NotFound),
					Message: fmt.Sprintf("qualification %s not found", id),
				},
			}
			continue
		}

		// Get a thread-safe snapshot
		snapshot := qualState.getSnapshot()

		// Create a new result from the snapshot
		result := lq.buildQualificationResult(snapshot, true)

		// Set status field for ERROR states
		if snapshot.State == plqpb.QualificationState_QUALIFICATION_STATE_ERROR {
			result.Status = &rpc.Status{
				Code:    int32(codes.Internal),
				Message: "qualification encountered an error",
			}
		}

		response.Results[id] = result
	}

	return response, nil
}

// Delete removes the qualification results for the provided ids
func (lq *linkQualification) Delete(ctx context.Context, req *plqpb.DeleteRequest) (*plqpb.DeleteResponse, error) {
	log.Infof("Received LinkQualification Delete request: %v", req.GetIds())

	if len(req.GetIds()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no qualification ids specified")
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	response := &plqpb.DeleteResponse{
		Results: make(map[string]*rpc.Status),
	}

	lq.mu.Lock()
	defer lq.mu.Unlock()

	for _, id := range req.GetIds() {
		qualState, exists := lq.qualifications[id]
		if !exists {
			response.Results[id] = &rpc.Status{
				Code:    int32(codes.NotFound),
				Message: fmt.Sprintf("qualification %s not found", id),
			}
			continue
		}

		// Cancel the qualification if it is not completed.
		qualState.mu.Lock()
		var cancellationFailed bool
		if !qualState.done {
			select {
			case qualState.cancelCh <- struct{}{}:
				log.Infof("Sent cancellation signal for qualification %s", id)
				// Update state to reflect cancellation.
				qualState.done = true
				qualState.State = plqpb.QualificationState_QUALIFICATION_STATE_ERROR
				qualState.EndTime = time.Now()
			default:
				// Channel is full or closed - cancellation failed.
				cancellationFailed = true
				log.Warningf("Failed to send cancellation signal for qualification %s - channel full or closed", id)
			}
		}
		// Create a snapshot of the final state after attempting cancellation.
		snapshot := qualState.getSnapshotLocked()
		qualState.mu.Unlock()

		if cancellationFailed {
			response.Results[id] = &rpc.Status{
				Code:    int32(codes.FailedPrecondition),
				Message: fmt.Sprintf("qualification %s cannot be stopped", id),
			}
			continue
		}

		// Store completed qualification in historical results if it completed successfully
		if snapshot.State == plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED {
			interfaceName := snapshot.InterfaceName

			// Create qualification result for historical storage
			result := lq.buildQualificationResult(snapshot, false)

			history := lq.historicalResults[interfaceName]
			history = append(history, result)

			// Keep only the most recent results up to maxHistorical
			if uint64(len(history)) > lq.maxHistorical {
				history = history[uint64(len(history))-lq.maxHistorical:]
			}
			lq.historicalResults[interfaceName] = history
		}

		delete(lq.qualifications, id)

		response.Results[id] = &rpc.Status{
			Code:    int32(codes.OK),
			Message: "qualification deleted successfully",
		}

		log.Infof("Deleted qualification %s for interface %s", id, snapshot.InterfaceName)
	}

	return response, nil
}

// List qualifications currently on the target
func (lq *linkQualification) List(ctx context.Context, req *plqpb.ListRequest) (*plqpb.ListResponse, error) {
	log.Infof("Received LinkQualification List request")

	lq.mu.RLock()
	defer lq.mu.RUnlock()

	var results []*plqpb.ListResult

	for _, qualState := range lq.qualifications {
		qualState.mu.Lock()
		result := &plqpb.ListResult{
			Id:            qualState.ID,
			State:         qualState.State,
			InterfaceName: qualState.InterfaceName,
		}
		qualState.mu.Unlock()
		results = append(results, result)
	}

	log.Infof("List results: %d total qualifications", len(results))

	return &plqpb.ListResponse{
		Results: results,
	}, nil
}

// validateQualificationConfig validates a single qualification configuration
func (lq *linkQualification) validateQualificationConfig(config *plqpb.QualificationConfiguration) error {
	id := config.GetId()

	// Check for duplicate ID across all link qualification operations,
	// in order to prevent two qualification tests from being created at the same time.
	_, exists := lq.qualifications[id]
	if exists {
		return status.Errorf(codes.AlreadyExists, "qualification id already exists")
	}

	// Validate interface exists in the system
	interfaceName := config.GetInterfaceName()
	interfacePath := ocpath.Root().Interface(interfaceName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := ygnmi.Get(ctx, lq.c, interfacePath.State())
	if err != nil {
		return status.Errorf(codes.NotFound, "interface %s not found", interfaceName)
	}

	if config.GetEndpointType() == nil {
		return status.Errorf(codes.InvalidArgument, "endpoint type is required")
	}
	if config.GetTiming() == nil {
		return status.Errorf(codes.InvalidArgument, "timing configuration is required")
	}

	// Validate timing configuration content
	timing := config.GetTiming()
	switch t := timing.(type) {
	case *plqpb.QualificationConfiguration_Rpc:
		rpcTiming := t.Rpc
		if rpcTiming.GetDuration() == nil {
			return status.Errorf(codes.InvalidArgument, "test duration is required")
		}
		duration := rpcTiming.GetDuration().AsDuration()
		if duration <= 0 {
			return status.Errorf(codes.InvalidArgument, "test duration must be positive")
		}
	case *plqpb.QualificationConfiguration_Ntp:
		// NTP timing is not implemented in simulation
		return status.Errorf(codes.Unimplemented, "NTP timing is not implemented in simulation")
	default:
		return status.Errorf(codes.InvalidArgument, "unknown timing configuration type")
	}

	// Validate endpoint type configuration
	switch et := config.GetEndpointType().(type) {
	case *plqpb.QualificationConfiguration_PacketGenerator:
		pg := et.PacketGenerator
		if pg.GetPacketRate() == 0 {
			return status.Errorf(codes.InvalidArgument, "packet rate must be greater than 0")
		}
	case *plqpb.QualificationConfiguration_PacketInjector:
		// PacketInjector is not implemented in simulation
		return status.Errorf(codes.Unimplemented, "PacketInjector endpoint type is not implemented in simulation")
	case *plqpb.QualificationConfiguration_AsicLoopback:
		// ASIC loopback is valid with any non-nil configuration
	case *plqpb.QualificationConfiguration_PmdLoopback:
		// PMD loopback is valid with any non-nil configuration
	default:
		return status.Errorf(codes.InvalidArgument, "unknown endpoint type")
	}

	return nil
}

// createQualificationState creates a qualification state from config
func (lq *linkQualification) createQualificationState(config *plqpb.QualificationConfiguration) *QualificationState {
	isGenerator := config.GetPacketGenerator() != nil
	isReflector := config.GetAsicLoopback() != nil || config.GetPmdLoopback() != nil

	state := &QualificationState{
		ID:            config.GetId(),
		InterfaceName: config.GetInterfaceName(),
		State:         plqpb.QualificationState_QUALIFICATION_STATE_IDLE,
		Config:        config,
		StartTime:     time.Now(),
		IsGenerator:   isGenerator,
		IsReflector:   isReflector,
		cancelCh:      make(chan struct{}, 1),
		done:          false,
	}

	log.Infof("Created qualification state: id=%s, interface=%s, generator=%v, reflector=%v",
		state.ID, state.InterfaceName, state.IsGenerator, state.IsReflector)

	return state
}

// executeQualification runs a single qualification by calling fakedevice simulation
func (lq *linkQualification) executeQualification(ctx context.Context, qual *QualificationState) {
	// Create cancellable context for the simulation
	qualCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Handle cancellation immediately
	go func() {
		select {
		case <-qual.cancelCh:
			log.Infof("Qualification %s cancelled", qual.ID)
			cancel()
		case <-qualCtx.Done():
		}
	}()

	// Update callback function to replace channel communication
	updateCallback := func(result *fakedevice.LinkQualificationResult) {
		qual.mu.Lock()
		defer qual.mu.Unlock()

		if qual.done {
			return
		}

		qual.State = result.State
		qual.packetsSent.Store(result.PacketsSent)
		qual.packetsReceived.Store(result.PacketsReceived)
		qual.packetsDropped.Store(result.PacketsDropped)
		qual.packetsError.Store(result.PacketsError)
		qual.StartTime = result.StartTime
		if !result.EndTime.IsZero() {
			qual.EndTime = result.EndTime
		}

		// Mark as done if in terminal state
		if result.State == plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED ||
			result.State == plqpb.QualificationState_QUALIFICATION_STATE_ERROR {
			qual.done = true
		}

		log.Infof("Qualification %s transitioned to state %v", qual.ID, result.State)
	}

	// Run the simulation with the callback.
	log.Infof("Starting RunPacketLinkQualification for %s", qual.ID)
	if err := fakedevice.RunPacketLinkQualification(qualCtx, lq.c, qual.Config, updateCallback, lq.config); err != nil {
		log.Errorf("Link qualification simulation failed for %s: %v", qual.ID, err)
		// Mark qualification as failed
		qual.mu.Lock()
		qual.State = plqpb.QualificationState_QUALIFICATION_STATE_ERROR
		qual.done = true
		qual.EndTime = time.Now()
		qual.mu.Unlock()
	} else {
		log.Infof("RunPacketLinkQualification completed successfully for %s", qual.ID)
	}
}

// getSnapshot creates a thread-safe snapshot of the qualification state
func (qs *QualificationState) getSnapshot() QualificationSnapshot {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	return qs.getSnapshotLocked()
}

// getSnapshotLocked creates a snapshot of the qualification state.
// It assumes the caller holds the lock on qs.mu.
func (qs *QualificationState) getSnapshotLocked() QualificationSnapshot {
	isComplete := qs.done
	return QualificationSnapshot{
		ID:              qs.ID,
		InterfaceName:   qs.InterfaceName,
		State:           qs.State,
		StartTime:       qs.StartTime,
		EndTime:         qs.EndTime,
		PacketsSent:     qs.packetsSent.Load(),
		PacketsReceived: qs.packetsReceived.Load(),
		PacketsDropped:  qs.packetsDropped.Load(),
		PacketsError:    qs.packetsError.Load(),
		Config:          qs.Config,
		IsComplete:      isComplete,
	}
}

// QualificationSnapshot is a snapshot of qualification state
type QualificationSnapshot struct {
	ID              string
	InterfaceName   string
	State           plqpb.QualificationState
	StartTime       time.Time
	EndTime         time.Time
	PacketsSent     uint64
	PacketsReceived uint64
	PacketsDropped  uint64
	PacketsError    uint64
	Config          *plqpb.QualificationConfiguration
	IsComplete      bool
}

// buildQualificationResult creates a QualificationResult from a snapshot with rate calculations
func (lq *linkQualification) buildQualificationResult(snapshot QualificationSnapshot, useCurrentTimeIfOngoing bool) *plqpb.QualificationResult {
	result := &plqpb.QualificationResult{
		Id:              snapshot.ID,
		InterfaceName:   snapshot.InterfaceName,
		State:           snapshot.State,
		StartTime:       timestamppb.New(snapshot.StartTime),
		PacketsSent:     snapshot.PacketsSent,
		PacketsReceived: snapshot.PacketsReceived,
		PacketsDropped:  snapshot.PacketsDropped,
		PacketsError:    snapshot.PacketsError,
	}

	// Set end time based on completion status
	if snapshot.IsComplete {
		result.EndTime = timestamppb.New(snapshot.EndTime)
	} else if useCurrentTimeIfOngoing {
		result.EndTime = timestamppb.Now()
	}

	// Calculate rates if we have packet data (sent or received)
	if snapshot.PacketsSent > 0 || snapshot.PacketsReceived > 0 {
		var duration time.Duration
		if snapshot.IsComplete {
			duration = snapshot.EndTime.Sub(snapshot.StartTime)
		} else if useCurrentTimeIfOngoing {
			duration = time.Since(snapshot.StartTime)
		}

		if duration > 0 {
			durationSeconds := duration.Seconds()
			if durationSeconds > 0 {
				// Get packet size from configuration
				packetSize := uint64(8184)
				if packetGen := snapshot.Config.GetPacketGenerator(); packetGen != nil && packetGen.GetPacketSize() > 0 {
					packetSize = uint64(packetGen.GetPacketSize())
				}
				// Calculate expected rate.
				if packetGen := snapshot.Config.GetPacketGenerator(); packetGen != nil {
					configuredRate := packetGen.GetPacketRate()
					if configuredRate > 0 && packetSize > 0 && configuredRate <= math.MaxUint64/packetSize {
						result.ExpectedRateBytesPerSecond = configuredRate * packetSize
					}
				} else {
					// For non-generator endpoints, calculate based on sent packets and duration
					totalBytes := snapshot.PacketsSent * packetSize
					if totalBytes > 0 {
						result.ExpectedRateBytesPerSecond = uint64(float64(totalBytes) / durationSeconds)
					}
				}

				// Calculate actual rate.
				actualBytes := snapshot.PacketsReceived * packetSize
				if actualBytes > 0 {
					result.QualificationRateBytesPerSecond = uint64(float64(actualBytes) / durationSeconds)
				}

				// Log rate calculations for ongoing qualifications
				if useCurrentTimeIfOngoing {
					var lossPercent float64
					if snapshot.PacketsSent > 0 {
						lossPercent = float64(snapshot.PacketsDropped) / float64(snapshot.PacketsSent) * 100
					}
					log.Infof("Calculated rates for qualification %s: packet_size=%d, expected_rate=%d Bps, actual_rate=%d Bps, loss=%.4f%%",
						snapshot.ID, packetSize,
						result.ExpectedRateBytesPerSecond, result.QualificationRateBytesPerSecond,
						lossPercent)
				}
			}
		}
	}

	return result
}

type wavelengthRouter struct {
	wrpb.UnimplementedWavelengthRouterServer
}

type Server struct {
	s                       *grpc.Server
	bgpServer               *bgp
	certServer              *cert
	diagServer              *diag
	fileServer              *file
	resetServer             *factoryReset
	healthzServer           *healthz
	layer2Server            *layer2
	linkQualificationServer *linkQualification
	mplsServer              *mpls
	osServer                *os
	otdrServer              *otdr
	systemServer            *system
	wavelengthRouterServer  *wavelengthRouter
}

func New(s *grpc.Server, gClient gpb.GNMIClient, target string, config *configpb.Config) (*Server, error) {
	yclient, err := ygnmi.NewClient(gClient, ygnmi.WithTarget(target), ygnmi.WithRequestLogLevel(2))
	if err != nil {
		return nil, err
	}

	srv := &Server{
		s:                       s,
		bgpServer:               &bgp{},
		certServer:              &cert{},
		diagServer:              &diag{},
		fileServer:              &file{},
		resetServer:             &factoryReset{},
		healthzServer:           &healthz{},
		layer2Server:            &layer2{},
		mplsServer:              &mpls{},
		osServer:                &os{},
		otdrServer:              &otdr{},
		linkQualificationServer: newLinkQualification(yclient, config),
		systemServer:            newSystem(yclient, config),
		wavelengthRouterServer:  &wavelengthRouter{},
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
	plqpb.RegisterLinkQualificationServer(s, srv.linkQualificationServer)
	spb.RegisterSystemServer(s, srv.systemServer)
	wrpb.RegisterWavelengthRouterServer(s, srv.wavelengthRouterServer)
	return srv, nil
}
