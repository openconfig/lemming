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

package fakedevice

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

	plqpb "github.com/openconfig/gnoi/packet_link_qualification"
	spb "github.com/openconfig/gnoi/system"

	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/gnmi/reconciler"
	"github.com/openconfig/lemming/internal/config"
	configpb "github.com/openconfig/lemming/proto/config"
)

const (
	DefaultNetworkInstance = "DEFAULT"
	StaticRoutingProtocol  = "DEFAULT"
	BGPRoutingProtocol     = "BGP"
)

// Reboot updates the system boot time to the provided Unix time.
func Reboot(ctx context.Context, c *ygnmi.Client, rebootTime int64) error {
	_, err := gnmiclient.Replace(gnmi.AddTimestampMetadata(ctx, rebootTime), c, ocpath.Root().System().BootTime().State(), uint64(rebootTime))
	return err
}

// RebootComponent updates the component's last reboot time and reason.
func RebootComponent(ctx context.Context, c *ygnmi.Client, componentName string, rebootTime int64, cfg *configpb.Config) error {
	log.Infof("Performing component reboot for %s at time %d", componentName, rebootTime)
	timestampedCtx := gnmi.AddTimestampMetadata(ctx, rebootTime)

	// Set component to inactive temporarily
	if _, err := gnmiclient.Replace(timestampedCtx, c, ocpath.Root().Component(componentName).OperStatus().State(), oc.PlatformTypes_COMPONENT_OPER_STATUS_INACTIVE); err != nil {
		return fmt.Errorf("failed to set component %s inactive: %v", componentName, err)
	}

	// Update last reboot time
	if _, err := gnmiclient.Replace(timestampedCtx, c, ocpath.Root().Component(componentName).LastRebootTime().State(), uint64(rebootTime)); err != nil {
		return fmt.Errorf("failed to update component %s reboot time: %v", componentName, err)
	}

	// Update reboot reason
	if _, err := gnmiclient.Replace(timestampedCtx, c, ocpath.Root().Component(componentName).LastRebootReason().State(), oc.PlatformTypes_COMPONENT_REBOOT_REASON_REBOOT_USER_INITIATED); err != nil {
		return fmt.Errorf("failed to update component %s reboot reason: %v", componentName, err)
	}

	// Simulate a configurable reboot period
	time.Sleep(time.Duration(cfg.GetTiming().GetRebootDurationMs()) * time.Millisecond)

	// Now restore the component OperStatus (reboot completed)
	finalState := oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE
	if _, err := gnmiclient.Replace(timestampedCtx, c, ocpath.Root().Component(componentName).OperStatus().State(), finalState); err != nil {
		return fmt.Errorf("failed to restore component %s state after reboot: %v", componentName, err)
	}

	log.Infof("Component %s reboot completed successfully", componentName)
	return nil
}

// SwitchoverSupervisor performs supervisor switchover by swapping the redundant roles and updating related state
func SwitchoverSupervisor(ctx context.Context, c *ygnmi.Client, targetSupervisor string, currentActiveSupervisor string, switchoverTime int64, cfg *configpb.Config) error {
	log.Infof("Performing supervisor switchover from %s to %s at time %d", currentActiveSupervisor, targetSupervisor, switchoverTime)

	timestampedCtx := gnmi.AddTimestampMetadata(ctx, switchoverTime)
	targetPath := ocpath.Root().Component(targetSupervisor)
	currentPath := ocpath.Root().Component(currentActiveSupervisor)

	time.Sleep(time.Duration(cfg.GetTiming().GetSwitchoverDurationMs()) * time.Millisecond)

	batch := &ygnmi.SetBatch{}

	// Swap the redundant roles
	gnmiclient.BatchReplace(batch, targetPath.RedundantRole().State(), oc.PlatformTypes_ComponentRedundantRole_PRIMARY)
	gnmiclient.BatchReplace(batch, currentPath.RedundantRole().State(), oc.PlatformTypes_ComponentRedundantRole_SECONDARY)

	// Update switchover timestamps for both supervisors
	gnmiclient.BatchReplace(batch, targetPath.LastSwitchoverTime().State(), uint64(switchoverTime))
	gnmiclient.BatchReplace(batch, currentPath.LastSwitchoverTime().State(), uint64(switchoverTime))

	// Update switchover reasons for both supervisors
	gnmiclient.BatchReplace(batch, targetPath.LastSwitchoverReason().Trigger().State(),
		oc.PlatformTypes_ComponentRedundantRoleSwitchoverReasonTrigger_USER_INITIATED)
	gnmiclient.BatchReplace(batch, targetPath.LastSwitchoverReason().Details().State(), "user initiated switchover")
	gnmiclient.BatchReplace(batch, currentPath.LastSwitchoverReason().Trigger().State(),
		oc.PlatformTypes_ComponentRedundantRoleSwitchoverReasonTrigger_USER_INITIATED)
	gnmiclient.BatchReplace(batch, currentPath.LastSwitchoverReason().Details().State(), "user initiated switchover")

	if _, err := batch.Set(timestampedCtx, c); err != nil {
		return fmt.Errorf("failed to apply switchover updates: %v", err)
	}

	log.Infof("Successfully completed supervisor switchover from %q to %q", currentActiveSupervisor, targetSupervisor)
	return nil
}

// KillProcess simulates process termination and restart functionality
func KillProcess(ctx context.Context, c *ygnmi.Client, pid uint32, processName string, signal spb.KillProcessRequest_Signal, restart bool, cfg *configpb.Config) error {
	log.Infof("KillProcess called with pid=%d, name=%s, signal=%v, restart=%v", pid, processName, signal, restart)

	processPath := ocpath.Root().System().Process(uint64(pid))
	currentTime := time.Now().UnixNano()
	timestampedCtx := gnmi.AddTimestampMetadata(ctx, currentTime)

	// HUP signal - reload configuration
	if signal == spb.KillProcessRequest_SIGNAL_HUP {
		log.Infof("Reloading process %s (PID: %d) configuration", processName, pid)

		if _, err := gnmiclient.Replace(timestampedCtx, c, processPath.StartTime().State(), uint64(currentTime)); err != nil {
			return fmt.Errorf("failed to update process %s reload time: %v", processName, err)
		}
		log.Infof("Successfully reloaded process %s (PID: %d)", processName, pid)
		return nil
	}

	if _, err := gnmiclient.Delete(timestampedCtx, c, processPath.State()); err != nil {
		return fmt.Errorf("failed to delete process %s (PID: %d): %v", processName, pid, err)
	}

	log.Infof("Process %s (PID: %d) terminated successfully", processName, pid)

	// Restart logic with 2-second delay if restart=true
	if restart {
		go func() {
			time.Sleep(2 * time.Second)

			// PID generation for restarted processes
			newPID, err := generateNewPID(ctx, c, pid)
			if err != nil {
				log.Errorf("Failed to generate new PID: %v", err)
				return
			}
			log.Infof("Restarting process %s with new PID: %d", processName, newPID)

			// Create new process with same name but new PID
			restartTime := time.Now().UnixNano()
			restartCtx := gnmi.AddTimestampMetadata(ctx, restartTime)

			var newProcess *oc.System_Process
			procConfig := config.GetProcessByName(cfg, processName)
			if procConfig != nil {
				newProcess = &oc.System_Process{
					Name:              ygot.String(procConfig.GetName()),
					Pid:               ygot.Uint64(uint64(newPID)),
					StartTime:         ygot.Uint64(uint64(restartTime)),
					CpuUsageUser:      ygot.Uint64(procConfig.GetCpuUsageUser()),
					CpuUsageSystem:    ygot.Uint64(procConfig.GetCpuUsageSystem()),
					CpuUtilization:    ygot.Uint8(uint8(procConfig.GetCpuUtilization())),
					MemoryUsage:       ygot.Uint64(procConfig.GetMemoryUsage()),
					MemoryUtilization: ygot.Uint8(uint8(procConfig.GetMemoryUtilization())),
				}
			} else {
				log.Errorf("Could not find configuration for process %s, skipping restart", processName)
				return
			}
			newProcessPath := ocpath.Root().System().Process(uint64(newPID))
			if _, err := gnmiclient.Replace(restartCtx, c, newProcessPath.State(), newProcess); err != nil {
				log.Errorf("Failed to restart process %s with new PID %d: %v", processName, newPID, err)
				return
			}

			log.Infof("Successfully restarted process %s with new PID: %d", processName, newPID)
		}()
	}
	return nil
}

// PingPacketResult represents the result of a single ping packet
type PingPacketResult struct {
	Sequence int32
	RTT      time.Duration
	Bytes    uint32
	TTL      int32
	Success  bool
}

// PingSimulation simulates ping operation with configurable network conditions
func PingSimulation(ctx context.Context, destination string, count int32, interval time.Duration, wait time.Duration, size uint32, responseChan chan<- *PingPacketResult, cfg *configpb.Config) error {
	log.Infof("Starting ping simulation to %s with count=%d, interval=%v, wait=%v",
		destination, count, interval, wait)

	baseLatencyMs := int64(10)
	ttl := int32(64)

	var jitterMs int64
	var packetLossRate float64

	// Use config values if available
	if cfg != nil && cfg.GetNetworkSimulation() != nil {
		if cfg.GetNetworkSimulation().GetBaseLatencyMs() > 0 {
			baseLatencyMs = cfg.GetNetworkSimulation().GetBaseLatencyMs()
		}
		jitterMs = max(0, cfg.GetNetworkSimulation().GetLatencyJitterMs())
		packetLossRate = max(0.0, min(1.0, float64(cfg.GetNetworkSimulation().GetPacketLossRate())))
		if cfg.GetNetworkSimulation().GetDefaultTtl() > 0 {
			ttl = cfg.GetNetworkSimulation().GetDefaultTtl()
		}
	}

	// Create random generator for this simulation
	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()>>32))) //nolint:gosec // Using math/rand for network simulation

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	packetCount := int32(0)
	maxPackets := count
	if maxPackets == -1 {
		maxPackets = math.MaxInt32
	}

	for packetCount < maxPackets {
		select {
		case <-ctx.Done():
			log.Infof("Ping to %s cancelled", destination)
			return ctx.Err()
		case <-ticker.C:
			packetCount++

			// Simple packet loss simulation
			packetLost := packetLossRate > 0 && rng.Float64() < packetLossRate

			var result *PingPacketResult

			if packetLost {
				// Packet lost - wait for timeout
				waitTimer := time.NewTimer(wait)
				select {
				case <-waitTimer.C:
					waitTimer.Stop()
				case <-ctx.Done():
					waitTimer.Stop()
					return ctx.Err()
				}
				result = &PingPacketResult{
					Sequence: packetCount,
					RTT:      0,
					Bytes:    0,
					TTL:      0,
					Success:  false,
				}
				log.Infof("Ping to %s: seq=%d TIMEOUT (packet lost)", destination, packetCount)
			} else {
				// Simple latency simulation: base + random jitter
				jitterRange := jitterMs * 2 // +/- jitter
				jitterOffset := rng.Int64N(jitterRange+1) - jitterMs
				totalLatencyMs := max(baseLatencyMs+jitterOffset, 1)

				networkLatency := time.Duration(totalLatencyMs) * time.Millisecond

				if networkLatency > wait {
					// Response would arrive after wait timeout
					waitTimer := time.NewTimer(wait)
					select {
					case <-waitTimer.C:
						waitTimer.Stop()
					case <-ctx.Done():
						waitTimer.Stop()
						return ctx.Err()
					}
					result = &PingPacketResult{
						Sequence: packetCount,
						RTT:      0,
						Bytes:    0,
						TTL:      0,
						Success:  false,
					}
					log.Infof("Ping to %s: seq=%d TIMEOUT (latency %v > wait %v)",
						destination, packetCount, networkLatency, wait)
				} else {
					// Successful response within wait time
					latencyTimer := time.NewTimer(networkLatency)
					select {
					case <-latencyTimer.C:
						latencyTimer.Stop()
					case <-ctx.Done():
						latencyTimer.Stop()
						return ctx.Err()
					}
					result = &PingPacketResult{
						Sequence: packetCount,
						RTT:      networkLatency,
						Bytes:    size,
						TTL:      ttl,
						Success:  true,
					}
					log.Infof("Ping to %s: seq=%d time=%v bytes=%d ttl=%d",
						destination, packetCount, networkLatency, size, ttl)
				}
			}

			// Send result with non-blocking send
			select {
			case responseChan <- result:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	log.Infof("Ping simulation to %s completed after %d packets", destination, packetCount)
	return nil
}

// NewBootTimeTask initializes boot-related paths.
func NewBootTimeTask(cfg *configpb.Config) *reconciler.BuiltReconciler {
	chassisName := cfg.GetComponents().GetChassisName()

	rec := reconciler.NewBuilder("boot time").
		WithStart(func(ctx context.Context, c *ygnmi.Client) error {
			now := time.Now().UnixNano()
			if err := Reboot(ctx, c, now); err != nil {
				return err
			}
			if _, err := gnmiclient.Replace(gnmi.AddTimestampMetadata(ctx, now), c, ocpath.Root().Component(chassisName).State(), &oc.Component{
				Name:            ygot.String(chassisName),
				Type:            oc.PlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_CHASSIS,
				OperStatus:      oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE,
				SoftwareVersion: ygot.String("current"),
			}); err != nil {
				return err
			}
			return nil
		}).Build()

	return rec
}

// NewChassisComponentsTask initializes subcomponents for the chassis
func NewChassisComponentsTask(cfg *configpb.Config) *reconciler.BuiltReconciler {
	rec := reconciler.NewBuilder("chassis components").
		WithStart(func(ctx context.Context, c *ygnmi.Client) error {
			now := time.Now().UnixNano()
			timestampedCtx := gnmi.AddTimestampMetadata(ctx, now)
			batch := &ygnmi.SetBatch{}

			chassisName := cfg.GetComponents().GetChassisName()
			supervisor1Name := cfg.GetComponents().GetSupervisor1Name()
			supervisor2Name := cfg.GetComponents().GetSupervisor2Name()

			// Initialize supervisors using explicit names
			var supervisorNames []string
			if supervisor2Name != "" {
				supervisorNames = []string{supervisor1Name, supervisor2Name}
			} else {
				supervisorNames = []string{supervisor1Name}
			}

			for i, componentName := range supervisorNames {
				redundantRole := oc.PlatformTypes_ComponentRedundantRole_PRIMARY
				if i == 1 {
					redundantRole = oc.PlatformTypes_ComponentRedundantRole_SECONDARY
				}
				component := &oc.Component{
					Name:               ygot.String(componentName),
					Type:               oc.PlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_CONTROLLER_CARD,
					OperStatus:         oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE,
					RedundantRole:      redundantRole,
					Parent:             ygot.String(chassisName),
					SoftwareVersion:    ygot.String("current"),
					LastRebootTime:     ygot.Uint64(uint64(now)),
					LastRebootReason:   oc.PlatformTypes_COMPONENT_REBOOT_REASON_UNSET,
					SwitchoverReady:    ygot.Bool(true),
					LastSwitchoverTime: ygot.Uint64(uint64(now)),
					LastSwitchoverReason: &oc.Component_LastSwitchoverReason{
						Trigger: oc.PlatformTypes_ComponentRedundantRoleSwitchoverReasonTrigger_SYSTEM_INITIATED,
						Details: ygot.String("initial system startup"),
					},
				}
				gnmiclient.BatchReplace(batch, ocpath.Root().Component(componentName).State(), component)
				log.Infof("Batching initialization for supervisor component %s", componentName)
			}

			// Initialize line cards
			linecardNames := config.GetAllLinecardNames(cfg)
			for _, componentName := range linecardNames {
				component := &oc.Component{
					Name:             ygot.String(componentName),
					Type:             oc.PlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_LINECARD,
					OperStatus:       oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE,
					Parent:           ygot.String(chassisName),
					SoftwareVersion:  ygot.String("current"),
					LastRebootTime:   ygot.Uint64(uint64(now)),
					LastRebootReason: oc.PlatformTypes_COMPONENT_REBOOT_REASON_UNSET,
				}
				gnmiclient.BatchReplace(batch, ocpath.Root().Component(componentName).State(), component)
				log.Infof("Batching initialization for line card component %s", componentName)
			}

			// Initialize fabric cards
			fabricNames := config.GetAllFabricNames(cfg)
			for _, componentName := range fabricNames {
				component := &oc.Component{
					Name:             ygot.String(componentName),
					Type:             oc.PlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_FABRIC,
					OperStatus:       oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE,
					Parent:           ygot.String(chassisName),
					SoftwareVersion:  ygot.String("current"),
					LastRebootTime:   ygot.Uint64(uint64(now)),
					LastRebootReason: oc.PlatformTypes_COMPONENT_REBOOT_REASON_UNSET,
				}
				gnmiclient.BatchReplace(batch, ocpath.Root().Component(componentName).State(), component)
				log.Infof("Batching initialization for fabric card component %s", componentName)
			}

			if _, err := batch.Set(timestampedCtx, c); err != nil {
				log.Errorf("Error applying batched component initializations: %v", err)
				return err
			}
			log.Infof("Successfully applied batched component initializations.")
			return nil
		}).Build()

	return rec
}

// NewCurrentTimeTask initializes boot-related paths.
func NewCurrentTimeTask() *reconciler.BuiltReconciler {
	rec := reconciler.NewBuilder("current time").
		WithStart(func(ctx context.Context, c *ygnmi.Client) error { // TODO: consider WithPeriodic if this a common pattern.
			tick := time.NewTicker(time.Second)
			periodic := func() error {
				_, err := gnmiclient.Replace(ctx, c, ocpath.Root().System().CurrentDatetime().State(), time.Now().Format(time.RFC3339))
				return err
			}

			if err := periodic(); err != nil {
				return err
			}
			go func() {
				for range tick.C {
					if err := periodic(); err != nil {
						log.Errorf("currentDateTimeTask error: %v", err)
						return
					}
				}
			}()
			return nil
		}).Build()

	return rec
}

// NewSystemBaseTask handles some of the logic for the base systems feature
// profile using ygnmi as the client.
func NewSystemBaseTask() *reconciler.BuiltReconciler {
	b := &ocpath.Batch{}
	b.AddPaths(
		ocpath.Root().System().Hostname().Config().PathStruct(),
		ocpath.Root().System().DomainName().Config().PathStruct(),
		ocpath.Root().System().MotdBanner().Config().PathStruct(),
		ocpath.Root().System().LoginBanner().Config().PathStruct(),
	)

	var hostname, domainName, motdBanner, loginBanner string

	rec := reconciler.NewTypedBuilder[*oc.Root]("system base").
		WithWatch(b.Config(), func(ctx context.Context, c *ygnmi.Client, v *ygnmi.Value[*oc.Root]) error {
			root, ok := v.Val()
			if !ok {
				return ygnmi.Continue
			}
			system := root.GetSystem()
			if system == nil {
				return ygnmi.Continue
			}
			if system.Hostname != nil && system.GetHostname() != hostname {
				if _, err := gnmiclient.Replace(ctx, c, ocpath.Root().System().Hostname().State(), system.GetHostname()); err != nil {
					log.Warningf("unable to update hostname: %v", err)
				} else {
					hostname = system.GetHostname()
					log.Infof("Successfully updated hostname to %q", hostname)
				}
			}
			if system.DomainName != nil && system.GetDomainName() != domainName {
				if _, err := gnmiclient.Replace(ctx, c, ocpath.Root().System().DomainName().State(), system.GetDomainName()); err != nil {
					log.Warningf("unable to update domainName: %v", err)
				} else {
					domainName = system.GetDomainName()
					log.Infof("Successfully updated domainName to %q", domainName)
				}
			}
			if system.MotdBanner != nil && system.GetMotdBanner() != motdBanner {
				if _, err := gnmiclient.Replace(ctx, c, ocpath.Root().System().MotdBanner().State(), system.GetMotdBanner()); err != nil {
					log.Warningf("unable to update motdBanner: %v", err)
				} else {
					motdBanner = system.GetMotdBanner()
					log.Infof("Successfully updated motdBanner to %q", motdBanner)
				}
			}
			if system.LoginBanner != nil && system.GetLoginBanner() != loginBanner {
				if _, err := gnmiclient.Replace(ctx, c, ocpath.Root().System().LoginBanner().State(), system.GetLoginBanner()); err != nil {
					log.Warningf("unable to update loginBanner: %v", err)
				} else {
					loginBanner = system.GetLoginBanner()
					log.Infof("Successfully updated loginBanner to %q", loginBanner)
				}
			}
			return ygnmi.Continue
		}).Build()

	return rec
}

// NewProcessMonitoringTask initializes system processes for monitoring.
func NewProcessMonitoringTask(cfg *configpb.Config) *reconciler.BuiltReconciler {
	rec := reconciler.NewBuilder("process monitoring").
		WithStart(func(ctx context.Context, c *ygnmi.Client) error {
			now := time.Now().UnixNano()
			timestampedCtx := gnmi.AddTimestampMetadata(ctx, now)
			batch := &ygnmi.SetBatch{}

			// Get processes from configuration
			processConfig := cfg.GetProcesses()
			if processConfig == nil {
				log.Warning("No processes configuration found")
			}
			processes := processConfig.GetProcess()
			log.Infof("Initializing %d system processes from configuration", len(processes))

			for _, procConfig := range processes {
				process := &oc.System_Process{
					Name:              ygot.String(procConfig.GetName()),
					Pid:               ygot.Uint64(uint64(procConfig.GetPid())),
					StartTime:         ygot.Uint64(uint64(now)),
					CpuUsageUser:      ygot.Uint64(procConfig.GetCpuUsageUser()),
					CpuUsageSystem:    ygot.Uint64(procConfig.GetCpuUsageSystem()),
					CpuUtilization:    ygot.Uint8(uint8(procConfig.GetCpuUtilization())),
					MemoryUsage:       ygot.Uint64(procConfig.GetMemoryUsage()),
					MemoryUtilization: ygot.Uint8(uint8(procConfig.GetMemoryUtilization())),
				}

				gnmiclient.BatchReplace(batch, ocpath.Root().System().Process(uint64(procConfig.GetPid())).State(), process)
				log.Infof("Batching initialization for process %s (PID: %d)", procConfig.GetName(), procConfig.GetPid())
			}

			if _, err := batch.Set(timestampedCtx, c); err != nil {
				log.Errorf("Error applying batched process initializations: %v", err)
				return err
			}

			log.Infof("Successfully initialized %d system processes", len(processes))
			return nil
		}).Build()

	return rec
}

// generateNewPID generates a new unique PID for restarted processes
func generateNewPID(ctx context.Context, c *ygnmi.Client, excludePID uint32) (uint32, error) {
	processes, err := ygnmi.GetAll(ctx, c, ocpath.Root().System().ProcessAny().State())
	if err != nil {
		return 0, fmt.Errorf("failed to get existing processes: %v", err)
	}

	// Build set of used PIDs
	used := make(map[uint32]bool)
	for _, p := range processes {
		if p.Pid != nil {
			used[uint32(*p.Pid)] = true
		}
	}
	used[excludePID] = true

	// Use a reasonable PID range
	for pid := uint32(1); pid <= 65535; pid++ {
		if !used[pid] {
			return pid, nil
		}
	}
	return 0, fmt.Errorf("no PID available in range 1-65535")
}

// NewInterfaceInitializationTask initializes base network interfaces for link qualification simulation
func NewInterfaceInitializationTask(cfg *configpb.Config) *reconciler.BuiltReconciler {
	rec := reconciler.NewBuilder("interface initialization").
		WithStart(func(ctx context.Context, c *ygnmi.Client) error {
			now := time.Now().UnixNano()
			timestampedCtx := gnmi.AddTimestampMetadata(ctx, now)
			batch := &ygnmi.SetBatch{}

			// Get interfaces from configuration
			interfaceSpecs := cfg.GetInterfaces().GetInterface()
			if len(interfaceSpecs) == 0 {
				log.Warning("No interfaces found in configuration")
				return nil
			}

			log.Infof("Initializing %d network interfaces from configuration", len(interfaceSpecs))

			for _, intfConfig := range interfaceSpecs {
				intf := &oc.Interface{
					Name:        ygot.String(intfConfig.GetName()),
					OperStatus:  oc.Interface_OperStatus_UP,
					Enabled:     ygot.Bool(true),
					Description: ygot.String(intfConfig.GetDescription()),
					Ifindex:     ygot.Uint32(intfConfig.GetIfIndex()),
				}

				gnmiclient.BatchReplace(batch, ocpath.Root().Interface(intfConfig.GetName()).State(), intf)
				log.Infof("Batching initialization for interface %s (ifindex: %d)", intfConfig.GetName(), intfConfig.GetIfIndex())
			}

			if _, err := batch.Set(timestampedCtx, c); err != nil {
				log.Errorf("Error applying batched interface initializations: %v", err)
				return err
			}

			log.Infof("Successfully initialized %d network interfaces", len(interfaceSpecs))
			return nil
		}).Build()

	return rec
}

// LinkQualificationResult represents the complete state of a link qualification operation
type LinkQualificationResult struct {
	mu              sync.Mutex
	State           plqpb.QualificationState
	PacketsSent     uint64
	PacketsReceived uint64
	PacketsDropped  uint64
	PacketsError    uint64
	StartTime       time.Time
	EndTime         time.Time
}

// RunPacketLinkQualification performs a complete packet-based link qualification simulation
func RunPacketLinkQualification(
	ctx context.Context,
	c *ygnmi.Client,
	config *plqpb.QualificationConfiguration,
	updateCallback func(*LinkQualificationResult),
	cfg *configpb.Config,
) error {
	qualID := config.GetId()
	interfaceName := config.GetInterfaceName()
	interfacePath := ocpath.Root().Interface(interfaceName)
	originalInterface, err := ygnmi.Get(ctx, c, interfacePath.State())
	if err != nil {
		return fmt.Errorf("failed to get interface %s: %w", interfaceName, err)
	}
	originalOperStatus := originalInterface.GetOperStatus()

	// Extract configuration parameters
	packetConfig := extractPacketConfiguration(config, cfg)
	timing := extractQualificationTiming(config, cfg)

	log.Infof("Qualification %s: packet_rate=%d PPS, packet_size=%d bytes, duration=%v",
		qualID, packetConfig.PacketRate, packetConfig.PacketSize, timing.TestDuration)

	result := &LinkQualificationResult{
		State:     plqpb.QualificationState_QUALIFICATION_STATE_IDLE,
		StartTime: time.Now(),
	}

	// Send initial state
	if updateCallback != nil {
		updateCallback(result)
	}

	// Always ensure interface is restored on exit
	defer func() {
		if err := restoreInterfaceOperStatus(context.Background(), c, interfaceName, originalOperStatus); err != nil {
			log.Errorf("Failed to restore interface %s status: %v", interfaceName, err)
		}
	}()

	// Execute the qualification state machine
	if err := executeQualificationStateMachine(ctx, c, interfaceName, result, packetConfig, timing, updateCallback, cfg); err != nil {
		result.mu.Lock()
		result.State = plqpb.QualificationState_QUALIFICATION_STATE_ERROR
		result.EndTime = time.Now()
		result.mu.Unlock()

		if updateCallback != nil {
			updateCallback(result)
		}

		return fmt.Errorf("qualification %s failed: %w", qualID, err)
	}

	log.Infof("Packet link qualification %s completed successfully", qualID)
	return nil
}

// PacketConfiguration holds extracted packet generation parameters
type PacketConfiguration struct {
	PacketRate  uint64
	PacketSize  uint64
	IsGenerator bool
	IsReflector bool
}

// extractPacketConfiguration extracts packet generation settings from config
func extractPacketConfiguration(config *plqpb.QualificationConfiguration, cfg *configpb.Config) *PacketConfiguration {
	linkQualConfig := cfg.GetLinkQualification()

	pc := &PacketConfiguration{
		PacketRate: linkQualConfig.GetDefaultPacketRate(),
		PacketSize: uint64(linkQualConfig.GetDefaultPacketSize()),
	}

	// Check for packet generator configuration and override config defaults
	if packetGen := config.GetPacketGenerator(); packetGen != nil {
		pc.IsGenerator = true
		if packetGen.GetPacketRate() > 0 {
			pc.PacketRate = packetGen.GetPacketRate()
		}
		if packetGen.GetPacketSize() > 0 {
			pc.PacketSize = uint64(packetGen.GetPacketSize())
		}
	}

	// Check for reflector configuration
	if config.GetAsicLoopback() != nil || config.GetPmdLoopback() != nil {
		pc.IsReflector = true
	}

	return pc
}

// QualificationTiming holds all timing parameters for a qualification
type QualificationTiming struct {
	SetupDuration    time.Duration
	TestDuration     time.Duration
	TeardownDuration time.Duration
	PreSyncDelay     time.Duration
	PostSyncDelay    time.Duration
}

// extractQualificationTiming extracts all timing parameters from config
func extractQualificationTiming(config *plqpb.QualificationConfiguration, cfg *configpb.Config) *QualificationTiming {
	// Get defaults from configuration
	linkQualConfig := cfg.GetLinkQualification()

	timing := &QualificationTiming{
		SetupDuration:    time.Duration(linkQualConfig.GetMinSetupDurationMs()) * time.Millisecond,
		TestDuration:     time.Duration(linkQualConfig.GetDefaultTestDurationMs()) * time.Millisecond,
		TeardownDuration: time.Duration(linkQualConfig.GetMinTeardownDurationMs()) * time.Millisecond,
	}

	// Extract timing from RPC synchronization
	if rpcTiming := config.GetRpc(); rpcTiming != nil {
		if rpcTiming.GetSetupDuration() != nil {
			timing.SetupDuration = rpcTiming.GetSetupDuration().AsDuration()
		}
		if rpcTiming.GetDuration() != nil {
			timing.TestDuration = rpcTiming.GetDuration().AsDuration()
		}
		if rpcTiming.GetTeardownDuration() != nil {
			timing.TeardownDuration = rpcTiming.GetTeardownDuration().AsDuration()
		}
		if rpcTiming.GetPreSyncDuration() != nil {
			timing.PreSyncDelay = rpcTiming.GetPreSyncDuration().AsDuration()
		}
		if rpcTiming.GetPostSyncDuration() != nil {
			timing.PostSyncDelay = rpcTiming.GetPostSyncDuration().AsDuration()
		}
	}
	return timing
}

// executeQualificationStateMachine runs a complete link qualification test through multiple phases:
// pre-sync delay (optional), SETUP (initialize test), RUNNING (execute packet generation and measurement),
// post-sync delay (optional), and TEARDOWN (cleanup resources).
// It manages state transitions, context cancellation, and result updates throughout the qualification process.
func executeQualificationStateMachine(
	ctx context.Context,
	c *ygnmi.Client,
	interfaceName string,
	result *LinkQualificationResult,
	packetConfig *PacketConfiguration,
	timing *QualificationTiming,
	updateCallback func(*LinkQualificationResult),
	cfg *configpb.Config,
) error {
	// Helper function to check context and send updates
	checkContextAndUpdate := func(state plqpb.QualificationState) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		result.mu.Lock()
		result.State = state
		result.mu.Unlock()
		if updateCallback != nil {
			updateCallback(result)
		}
		return nil
	}

	// Pre-sync delay
	if timing.PreSyncDelay > 0 {
		log.Infof("Waiting pre-sync delay: %v", timing.PreSyncDelay)
		if err := sleep(ctx, timing.PreSyncDelay); err != nil {
			return err
		}
	}

	// SETUP phase
	if err := checkContextAndUpdate(plqpb.QualificationState_QUALIFICATION_STATE_SETUP); err != nil {
		return err
	}
	if err := executeSetupPhase(ctx, c, interfaceName, timing.SetupDuration); err != nil {
		return fmt.Errorf("setup phase failed: %w", err)
	}

	// RUNNING phase
	if err := checkContextAndUpdate(plqpb.QualificationState_QUALIFICATION_STATE_RUNNING); err != nil {
		return err
	}
	if err := executeRunningPhase(ctx, result, packetConfig, timing, updateCallback, cfg); err != nil {
		return fmt.Errorf("running phase failed: %w", err)
	}

	// Post-sync delay
	if timing.PostSyncDelay > 0 {
		log.Infof("Waiting post-sync delay: %v", timing.PostSyncDelay)
		if err := sleep(ctx, timing.PostSyncDelay); err != nil {
			return err
		}
	}

	// TEARDOWN phase
	if err := checkContextAndUpdate(plqpb.QualificationState_QUALIFICATION_STATE_TEARDOWN); err != nil {
		return err
	}
	if err := executeTeardownPhase(ctx, timing.TeardownDuration); err != nil {
		return fmt.Errorf("teardown phase failed: %w", err)
	}

	// COMPLETED phase
	result.mu.Lock()
	result.State = plqpb.QualificationState_QUALIFICATION_STATE_COMPLETED
	result.EndTime = time.Now()
	result.mu.Unlock()
	if updateCallback != nil {
		updateCallback(result)
	}

	return nil
}

// sleep is a context-aware sleep function
func sleep(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}

// executeSetupPhase handles the SETUP state of qualification
func executeSetupPhase(
	ctx context.Context,
	c *ygnmi.Client,
	interfaceName string,
	setupDuration time.Duration,
) error {
	log.Infof("Entering SETUP phase for %v", setupDuration)

	// Set interface to TESTING state
	timestampedCtx := gnmi.AddTimestampMetadata(ctx, time.Now().UnixNano())
	if _, err := gnmiclient.Replace(timestampedCtx, c, ocpath.Root().Interface(interfaceName).OperStatus().State(), oc.Interface_OperStatus_TESTING); err != nil {
		return fmt.Errorf("failed to set interface to TESTING state: %w", err)
	}

	// Wait for setup duration
	if err := sleep(ctx, setupDuration); err != nil {
		return err
	}

	log.Infof("SETUP phase completed")
	return nil
}

// executeRunningPhase handles the RUNNING state with packet simulation
func executeRunningPhase(
	ctx context.Context,
	result *LinkQualificationResult,
	packetConfig *PacketConfiguration,
	timing *QualificationTiming,
	updateCallback func(*LinkQualificationResult),
	cfg *configpb.Config,
) error {
	log.Infof("Entering RUNNING phase for %v", timing.TestDuration)

	// Start packet simulation
	startTime := time.Now()
	sampleInterval := time.Duration(cfg.GetLinkQualification().GetMinSampleIntervalMs()) * time.Millisecond
	updateTicker := time.NewTicker(sampleInterval)
	defer updateTicker.Stop()
	testEndTime := startTime.Add(timing.TestDuration)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-updateTicker.C:
			elapsed := now.Sub(startTime)

			// Update packet statistics based on endpoint type
			updatePacketStatistics(result, packetConfig, elapsed, cfg)

			if updateCallback != nil {
				updateCallback(result)
			}

			// Check if test duration completed
			if now.After(testEndTime) || now.Equal(testEndTime) {
				log.Infof("RUNNING phase completed after %v", elapsed)
				return nil
			}
		}
	}
}

// executeTeardownPhase handles the TEARDOWN state
func executeTeardownPhase(
	ctx context.Context,
	teardownDuration time.Duration,
) error {
	log.Infof("Entering TEARDOWN phase for %v", teardownDuration)

	// Wait for teardown duration
	if err := sleep(ctx, teardownDuration); err != nil {
		return err
	}

	log.Infof("TEARDOWN phase completed")
	return nil
}

// updatePacketStatistics calculates realistic packet statistics for both generators and reflectors
func updatePacketStatistics(
	result *LinkQualificationResult,
	packetConfig *PacketConfiguration,
	elapsed time.Duration,
	cfg *configpb.Config,
) {
	elapsedSeconds := elapsed.Seconds()
	if elapsedSeconds <= 0 {
		return
	}
	rate := packetConfig.PacketRate
	expectedPackets := uint64(elapsedSeconds * float64(rate))

	// Apply network simulation parameters
	var lossRate float32
	var errorRate float32
	if cfg != nil && cfg.GetNetworkSimulation() != nil {
		lossRate = cfg.GetNetworkSimulation().GetPacketLossRate()
		errorRate = cfg.GetNetworkSimulation().GetPacketErrorRate()
	}

	dropped := uint64(float64(expectedPackets) * float64(lossRate))
	errored := uint64(float64(expectedPackets) * float64(errorRate))
	successful := expectedPackets - dropped

	// Update statistics atomically
	result.mu.Lock()
	if packetConfig.IsGenerator {
		// Generator: sends packets, receives responses
		result.PacketsSent = expectedPackets
		result.PacketsReceived = successful
	} else {
		// Reflector: receives packets, reflects them back
		result.PacketsReceived = successful
		result.PacketsSent = successful
	}
	result.PacketsDropped = dropped
	result.PacketsError = errored
	result.mu.Unlock()
}

// restoreInterfaceOperStatus restores interface to original operational state
func restoreInterfaceOperStatus(ctx context.Context, c *ygnmi.Client, interfaceName string, originalStatus oc.E_Interface_OperStatus) error {
	timestampedCtx := gnmi.AddTimestampMetadata(ctx, time.Now().UnixNano())
	if _, err := gnmiclient.Replace(timestampedCtx, c, ocpath.Root().Interface(interfaceName).OperStatus().State(), originalStatus); err != nil {
		return fmt.Errorf("failed to restore interface operational status to %v: %v", originalStatus, err)
	}
	log.Infof("Restored interface operational status to %v", originalStatus)
	return nil
}
