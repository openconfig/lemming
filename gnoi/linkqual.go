// Copyright 2025 Google LLC
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
	plqpb "github.com/openconfig/gnoi/packet_link_qualification"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	configpb "github.com/openconfig/lemming/proto/config"
	"github.com/openconfig/ygnmi/ygnmi"
	rpc "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type linkQualification struct {
	plqpb.UnimplementedLinkQualificationServer

	c      *ygnmi.Client
	config *configpb.Config

	// Protect the qualification state tracking and historical results
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
