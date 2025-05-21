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
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"

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

	// Number of each component type
	numLineCard       = 10
	numFabricCard     = 6
	numSupervisorCard = 2
)

// Reboot updates the system boot time to the provided Unix time.
func Reboot(ctx context.Context, c *ygnmi.Client, rebootTime int64) error {
	_, err := gnmiclient.Replace(gnmi.AddTimestampMetadata(ctx, rebootTime), c, ocpath.Root().System().BootTime().State(), uint64(rebootTime))
	return err
}

// RebootComponent updates the component's last reboot time and reason.
func RebootComponent(ctx context.Context, c *ygnmi.Client, componentName string, rebootTime int64) error {
	log.Infof("Performing component reboot for %s at time %d", componentName, rebootTime)

	// Get current component state to ensure it exists
	_, err := ygnmi.Get(ctx, c, ocpath.Root().Component(componentName).State())
	if err != nil {
		return fmt.Errorf("failed to get component %s state: %v", componentName, err)
	}

	// Set component to inactive temporarily
	if _, err := gnmiclient.Replace(ctx, c, ocpath.Root().Component(componentName).OperStatus().State(), oc.PlatformTypes_COMPONENT_OPER_STATUS_INACTIVE); err != nil {
		return fmt.Errorf("failed to set component %s inactive: %v", componentName, err)
	}

	// Update last reboot time
	if _, err := gnmiclient.Replace(ctx, c, ocpath.Root().Component(componentName).LastRebootTime().State(), uint64(rebootTime)); err != nil {
		return fmt.Errorf("failed to update component %s reboot time: %v", componentName, err)
	}

	// Update reboot reason
	if _, err := gnmiclient.Replace(ctx, c, ocpath.Root().Component(componentName).LastRebootReason().State(), oc.PlatformTypes_COMPONENT_REBOOT_REASON_REBOOT_USER_INITIATED); err != nil {
		return fmt.Errorf("failed to update component %s reboot reason: %v", componentName, err)
	}

	// Simulate a brief reboot period
	time.Sleep(2 * time.Minute)

	// Now restore the component OperStatus (reboot completed)
	finalState := oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE
	if _, err := gnmiclient.Replace(ctx, c, ocpath.Root().Component(componentName).OperStatus().State(), finalState); err != nil {
		return fmt.Errorf("failed to restore component %s state after reboot: %v", componentName, err)
	}

	log.Infof("Component %s reboot completed successfully", componentName)
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
			for i := 0; i < numSupervisorCard; i++ {
				componentName := fmt.Sprintf("%s%d", controlcardComponentName, i)
				redundantRole := oc.PlatformTypes_ComponentRedundantRole_PRIMARY
				if i == 1 {
					redundantRole = oc.PlatformTypes_ComponentRedundantRole_SECONDARY
				}
				component := &oc.Component{
					Name:             ygot.String(componentName),
					Type:             oc.PlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_CONTROLLER_CARD,
					OperStatus:       oc.PlatformTypes_COMPONENT_OPER_STATUS_ACTIVE,
					RedundantRole:    redundantRole,
					Parent:           ygot.String(chassisComponentName),
					SoftwareVersion:  ygot.String("current"),
					LastRebootTime:   ygot.Uint64(uint64(now)),
					LastRebootReason: oc.PlatformTypes_COMPONENT_REBOOT_REASON_UNSET,
				}
				gnmiclient.BatchUpdate(batch, ocpath.Root().Component(componentName).State(), component)
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
				gnmiclient.BatchUpdate(batch, ocpath.Root().Component(componentName).State(), component)
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
