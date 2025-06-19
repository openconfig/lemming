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
	"math/rand/v2"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

	spb "github.com/openconfig/gnoi/system"

	"github.com/openconfig/lemming/gnmi"
	"github.com/openconfig/lemming/gnmi/gnmiclient"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/lemming/gnmi/reconciler"
)

const (
	DefaultNetworkInstance = "DEFAULT"
	StaticRoutingProtocol  = "DEFAULT"
	BGPRoutingProtocol     = "BGP"

	// Component names
	chassisComponentName     = "chassis"
	linecardComponentName    = "Linecard"
	fabricComponentName      = "Fabric"
	controlcardComponentName = "Supervisor"

	// TODO: Make lemming chassis configurable
	// Number of each component type
	numLineCard       = 8
	numFabricCard     = 6
	numSupervisorCard = 2

	// Simulation duration
	switchoverDuration = 2 * time.Second
	rebootDuration     = 2 * time.Second

	// Process configuration constants
	basePID    = 1000
	ceilingPID = 1100

	// Default mock process configurations
	defaultCpuUsageUser      = 1000000  // 1ms in nanoseconds
	defaultCpuUsageSystem    = 500000   // 0.5ms in nanoseconds
	defaultCpuUtilization    = 1        // 1% CPU
	defaultMemoryUsage       = 10485760 // 10MB
	defaultMemoryUtilization = 2        // 2% memory

	// Ping simulation configuration - final latency = baseLatency ± random jitter
	defaultBaseLatency    = 50 * time.Millisecond // Base latency for ping responses
	defaultLatencyJitter  = 20 * time.Millisecond // ±20ms jitter variation
	defaultPacketLossRate = 0.0                   // 0% packet loss by default
	defaultTTL            = 64                    // Default TTL value
)

var (
	// Thread-safe random source for ping simulation
	randMu  sync.Mutex
	randSrc = rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()>>32))
	randGen = rand.New(randSrc) //nolint:gosec // Using math/rand for network simulation
)

// TODO: Add per-destination ping configuration:
// - Host unreachable
// - IPv4/IPv6 differences
// - Network congestion

// PingPacketResult represents the result of a single ping packet
type PingPacketResult struct {
	Sequence int32
	RTT      time.Duration
	Bytes    uint32
	TTL      int32
	Success  bool
}

// Reboot updates the system boot time to the provided Unix time.
func Reboot(ctx context.Context, c *ygnmi.Client, rebootTime int64) error {
	_, err := gnmiclient.Replace(gnmi.AddTimestampMetadata(ctx, rebootTime), c, ocpath.Root().System().BootTime().State(), uint64(rebootTime))
	return err
}

// RebootComponent updates the component's last reboot time and reason.
func RebootComponent(ctx context.Context, c *ygnmi.Client, componentName string, rebootTime int64) error {
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

	// Simulate a brief reboot period
	time.Sleep(rebootDuration)

	// Now restore the component OperStatus (reboot completed)
	finalState := oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE
	if _, err := gnmiclient.Replace(timestampedCtx, c, ocpath.Root().Component(componentName).OperStatus().State(), finalState); err != nil {
		return fmt.Errorf("failed to restore component %s state after reboot: %v", componentName, err)
	}

	log.Infof("Component %s reboot completed successfully", componentName)
	return nil
}

// SwitchoverSupervisor performs supervisor switchover by swapping the redundant roles and updating related state
func SwitchoverSupervisor(ctx context.Context, c *ygnmi.Client, targetSupervisor string, currentActiveSupervisor string, switchoverTime int64) error {
	log.Infof("Performing supervisor switchover from %s to %s at time %d", currentActiveSupervisor, targetSupervisor, switchoverTime)

	timestampedCtx := gnmi.AddTimestampMetadata(ctx, switchoverTime)
	targetPath := ocpath.Root().Component(targetSupervisor)
	currentPath := ocpath.Root().Component(currentActiveSupervisor)

	time.Sleep(switchoverDuration)

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
func KillProcess(ctx context.Context, c *ygnmi.Client, pid uint32, processName string, signal spb.KillProcessRequest_Signal, restart bool) error {
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
			newPID, err := generateNewPID(ctx, c)
			if err != nil {
				log.Errorf("Failed to generate new PID: %v", err)
				return
			}
			log.Infof("Restarting process %s with new PID: %d", processName, newPID)

			// Create new process with same name but new PID
			restartTime := time.Now().UnixNano()
			restartCtx := gnmi.AddTimestampMetadata(ctx, restartTime)

			newProcess := &oc.System_Process{
				Name:      ygot.String(processName),
				Pid:       ygot.Uint64(uint64(newPID)),
				StartTime: ygot.Uint64(uint64(restartTime)),
				// Simulate realistic resource usage
				CpuUsageUser:      ygot.Uint64(defaultCpuUsageUser),
				CpuUsageSystem:    ygot.Uint64(defaultCpuUsageSystem),
				CpuUtilization:    ygot.Uint8(defaultCpuUtilization),
				MemoryUsage:       ygot.Uint64(defaultMemoryUsage),
				MemoryUtilization: ygot.Uint8(defaultMemoryUtilization),
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

// PingSimulation simulates ping operation with configurable network conditions
func PingSimulation(ctx context.Context, destination string, count int32, interval time.Duration, wait time.Duration, size uint32, responseChan chan<- *PingPacketResult) error {
	log.Infof("Starting ping simulation to %s with count=%d, interval=%v, wait=%v",
		destination, count, interval, wait)

	baseLatency := defaultBaseLatency
	latencyJitter := defaultLatencyJitter
	packetLossRate := float32(defaultPacketLossRate)
	ttl := int32(defaultTTL)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	packetCount := int32(0)
	maxPackets := count
	if maxPackets == -1 {
		maxPackets = 1<<31 - 1
	}

	for packetCount < maxPackets {
		select {
		case <-ctx.Done():
			log.Infof("Ping to %s cancelled", destination)
			return ctx.Err()
		case <-ticker.C:
			packetCount++

			// Simulate network conditions for this packet
			networkLatency := simulateLatency(baseLatency, latencyJitter)
			packetLost := packetLossRate > 0 && shouldDropPacket(packetLossRate)

			var result *PingPacketResult

			switch {
			case packetLost:
				// Packet lost - wait for timeout
				waitTimer := time.NewTimer(wait)
				select {
				case <-waitTimer.C:
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

			case networkLatency > wait:
				// Response would arrive after wait timeout
				waitTimer := time.NewTimer(wait)
				select {
				case <-waitTimer.C:
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

			default:
				// Successful response within wait time
				latencyTimer := time.NewTimer(networkLatency)
				select {
				case <-latencyTimer.C:
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
func NewBootTimeTask() *reconciler.BuiltReconciler {
	rec := reconciler.NewBuilder("boot time").
		WithStart(func(ctx context.Context, c *ygnmi.Client) error {
			now := time.Now().UnixNano()
			if err := Reboot(ctx, c, now); err != nil {
				return err
			}
			if _, err := gnmiclient.Replace(gnmi.AddTimestampMetadata(ctx, now), c, ocpath.Root().Component(chassisComponentName).State(), &oc.Component{
				Name:            ygot.String(chassisComponentName),
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
func NewChassisComponentsTask() *reconciler.BuiltReconciler {
	rec := reconciler.NewBuilder("chassis components").
		WithStart(func(ctx context.Context, c *ygnmi.Client) error {
			now := time.Now().UnixNano()
			timestampedCtx := gnmi.AddTimestampMetadata(ctx, now)
			batch := &ygnmi.SetBatch{}

			// Initialize supervisors
			for i := 1; i <= numSupervisorCard; i++ {
				componentName := fmt.Sprintf("%s%d", controlcardComponentName, i)
				redundantRole := oc.PlatformTypes_ComponentRedundantRole_PRIMARY
				if i == 2 {
					redundantRole = oc.PlatformTypes_ComponentRedundantRole_SECONDARY
				}
				component := &oc.Component{
					Name:               ygot.String(componentName),
					Type:               oc.PlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_CONTROLLER_CARD,
					OperStatus:         oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE,
					RedundantRole:      redundantRole,
					Parent:             ygot.String(chassisComponentName),
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
			for i := 0; i < numLineCard; i++ {
				componentName := fmt.Sprintf("%s%d", linecardComponentName, i)
				component := &oc.Component{
					Name:             ygot.String(componentName),
					Type:             oc.PlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_LINECARD,
					OperStatus:       oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE,
					Parent:           ygot.String(chassisComponentName),
					SoftwareVersion:  ygot.String("current"),
					LastRebootTime:   ygot.Uint64(uint64(now)),
					LastRebootReason: oc.PlatformTypes_COMPONENT_REBOOT_REASON_UNSET,
				}
				gnmiclient.BatchReplace(batch, ocpath.Root().Component(componentName).State(), component)
				log.Infof("Batching initialization for line card component %s", componentName)
			}

			// Initialize fabric cards
			for i := 0; i < numFabricCard; i++ {
				componentName := fmt.Sprintf("%s%d", fabricComponentName, i)
				component := &oc.Component{
					Name:             ygot.String(componentName),
					Type:             oc.PlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_FABRIC,
					OperStatus:       oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE,
					Parent:           ygot.String(chassisComponentName),
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

// NewProcessMonitoringTask initializes mock system processes for monitoring.
func NewProcessMonitoringTask() *reconciler.BuiltReconciler {
	rec := reconciler.NewBuilder("process monitoring").
		WithStart(func(ctx context.Context, c *ygnmi.Client) error {
			now := time.Now().UnixNano()
			timestampedCtx := gnmi.AddTimestampMetadata(ctx, now)
			batch := &ygnmi.SetBatch{}

			// Mock daemon processes with their PIDs
			processes := map[string]uint64{
				"bgpd":        basePID + 1,
				"ospfd":       basePID + 2,
				"gnmi-server": basePID + 3,
				"sysrib":      basePID + 4,
			}

			log.Infof("Initializing %d mock system processes", len(processes))

			for processName, pid := range processes {
				process := &oc.System_Process{
					Name:      ygot.String(processName),
					Pid:       ygot.Uint64(pid),
					StartTime: ygot.Uint64(uint64(now)),
					// Simulate realistic resource usage
					CpuUsageUser:      ygot.Uint64(defaultCpuUsageUser),
					CpuUsageSystem:    ygot.Uint64(defaultCpuUsageSystem),
					CpuUtilization:    ygot.Uint8(defaultCpuUtilization),
					MemoryUsage:       ygot.Uint64(defaultMemoryUsage),
					MemoryUtilization: ygot.Uint8(defaultMemoryUtilization),
				}

				gnmiclient.BatchReplace(batch, ocpath.Root().System().Process(pid).State(), process)
				log.Infof("Batching initialization for process %s (PID: %d)", processName, pid)
			}

			if _, err := batch.Set(timestampedCtx, c); err != nil {
				log.Errorf("Error applying batched process initializations: %v", err)
				return err
			}

			log.Infof("Successfully initialized %d mock system processes", len(processes))
			return nil
		}).Build()

	return rec
}

// generateNewPID generates a new unique PID for restarted processes
func generateNewPID(ctx context.Context, c *ygnmi.Client) (uint32, error) {
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
	for pid := uint32(1005); pid <= ceilingPID; pid++ {
		if !used[pid] {
			return pid, nil
		}
	}
	return 0, fmt.Errorf("no PID available")
}

// getRandomInt64 returns a random int64 in the range [min, max)
func getRandomInt64(min, max int64) int64 {
	randMu.Lock()
	defer randMu.Unlock()
	if min >= max {
		return min
	}
	return randGen.Int64N(max-min) + min
}

// shouldDropPacket determines if a packet should be dropped based on loss rate
func shouldDropPacket(lossRate float32) bool {
	if lossRate <= 0 {
		return false
	}
	if lossRate >= 1.0 {
		return true
	}
	lossPercentage := int64(lossRate * 100)
	randomPercent := getRandomInt64(0, 100)
	return randomPercent < lossPercentage
}

// simulateLatency generates realistic network latency with jitter
func simulateLatency(baseLatency, jitter time.Duration) time.Duration {
	if jitter == 0 {
		return baseLatency
	}
	// Generate random jitter between -jitter and +jitter
	jitterNs := jitter.Nanoseconds()
	randomJitter := getRandomInt64(-jitterNs, jitterNs+1)
	finalLatency := baseLatency + time.Duration(randomJitter)

	// Ensure latency is never negative
	if finalLatency <= 0 {
		finalLatency = time.Millisecond
	}

	return finalLatency
}
