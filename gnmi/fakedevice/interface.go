package fakedevice

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/coalesce"
	"github.com/openconfig/gnmi/ctree"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/lemming/gnmi/gnmit"
	"github.com/openconfig/lemming/gnmi/internal/config"
	configpath "github.com/openconfig/lemming/gnmi/internal/config/device"
	"github.com/openconfig/lemming/gnmi/internal/telemetry"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/encoding/prototext"
)

// TODO(wenbli): IF we somehow end up using this instead of splitting out the tasks into a separate binary, then we should move this to a different package and add simple testing.
var (
	// Stubs for testing.
	updateInterfaceFn = updateInterface
	deleteInterfaceFn = deleteInterface
)

func updateInterface(*telemetry.Interface) error {
	// TODO: This needs to call into the dataplane to configure the interface.
	return nil
}

func deleteInterface(name string) error {
	// TODO: This needs to call into the dataplane to configure the interface.
	return nil
}

// TODO(wenbli): This file needs to be put in its own package. Common utilities
// need to be factored out to enable this.
var (
	enabledPaths, descriptionPaths, namePaths, ipv4AddressPaths, prefixLengthPaths *gpb.Path
	interfacePendingEvents                                                         map[*func(*config.Device) error]bool
	// appliedRoot is the SoT for BGP applied configuration. It is maintained locally by the task.
	interfaceAppliedRoot *telemetry.Device
	interfaceAppliedMu   sync.Mutex
)

func initInterfaceTaskVars() error {
	interfaceAppliedRoot = &telemetry.Device{}

	interfacePendingEvents = map[*func(*config.Device) error]bool{}
	return initInterfacePaths()
}

func initInterfacePaths() error {
	interfacePath := configpath.DeviceRoot("").InterfaceAny()
	var err []error
	enabledPaths, _, err = ygot.ResolvePath(interfacePath.Enabled())
	if err != nil {
		return fmt.Errorf("interfaceTask failed to initialize due to error: %v", err)
	}
	descriptionPaths, _, err = ygot.ResolvePath(interfacePath.Description())
	if err != nil {
		return fmt.Errorf("interfaceTask failed to initialize due to error: %v", err)
	}
	namePaths, _, err = ygot.ResolvePath(interfacePath.Name())
	if err != nil {
		return fmt.Errorf("interfaceTask failed to initialize due to error: %v", err)
	}
	ipv4AddressPaths, _, err = ygot.ResolvePath(interfacePath.SubinterfaceAny().Ipv4().AddressAny().Ip())
	if err != nil {
		return fmt.Errorf("interfaceTask failed to initialize due to error: %v", err)
	}
	prefixLengthPaths, _, err = ygot.ResolvePath(interfacePath.SubinterfaceAny().Ipv4().AddressAny().PrefixLength())
	if err != nil {
		return fmt.Errorf("interfaceTask failed to initialize due to error: %v", err)
	}

	return nil
}

func interfaceTask(getIntendedConfig func() *config.Device, q gnmit.Queue, update gnmit.UpdateFn, target string, remove func()) error {
	if err := initInterfaceTaskVars(); err != nil {
		return err
	}

	// updateAppliedConfig computes the diff between a previous applied
	// configuration and the current SoT, and sends the updates to the
	// central DB.
	updateAppliedConfig := func(prevApplied *telemetry.Device) bool {
		interfaceAppliedMu.Lock()
		defer interfaceAppliedMu.Unlock()
		no, err := ygot.Diff(prevApplied, interfaceAppliedRoot)
		if err != nil {
			log.Errorf("interfaceTask: error while creating update notification for updating applied configuration: %v", err)
			return false
		}
		if len(no.GetUpdate())+len(no.GetDelete()) > 0 {
			log.V(1).Info("Updating interface applied configuration: ", prototext.Format(no))
			no.Timestamp = time.Now().UnixNano()
			no.Prefix = &gpb.Path{Origin: "openconfig", Target: target}

			if err := update(no); err != nil {
				log.Errorf("interfaceTask: error while writing update to applied configuration: %v", err)
				return false
			}
		}
		return true
	}

	go func() {
		defer remove()
		for {
			item, _, err := q.Next(context.Background())
			if coalesce.IsClosedQueue(err) {
				return
			}
			n, ok := item.(*ctree.Leaf)
			if !ok || n == nil {
				log.Errorf("interfaceTask invalid cache node: %#v", item)
				return
			}
			v := n.Value()
			no, ok := v.(*gpb.Notification)
			if !ok || no == nil {
				log.Errorf("interfaceTask invalid cache node, expected non-nil *gpb.Notification type, got: %#v", v)
				return
			}

			interfaceNotificationHandler(no)

			var updateAppliedRoot bool
			for _, triggered := range interfacePendingEvents {
				if triggered {
					updateAppliedRoot = true
				}
			}

			var prevApplied *telemetry.Device
			if updateAppliedRoot {
				interfaceAppliedMu.Lock()
				prevAppliedGS, err := ygot.DeepCopy(interfaceAppliedRoot)
				if err != nil {
					log.Fatalf("interfaceTask: Could not copy applied configuration: %v", err)
				}
				prevApplied = prevAppliedGS.(*telemetry.Device)
				interfaceAppliedMu.Unlock()
			}

			for reactor, triggered := range interfacePendingEvents {
				if triggered {
					if err := (*reactor)(getIntendedConfig()); err != nil {
						log.Errorf("interfaceTask reactor: %v", err)
					}
					interfacePendingEvents[reactor] = false
				}
			}

			if success := updateAppliedConfig(prevApplied); !success {
				log.Errorf("interfaceTask: updating applied configuration failed")
			}
		}
	}()

	return nil
}

func interfaceNotificationHandler(no *gpb.Notification) error {
	for _, u := range no.Update {
		interfacePathHandler(u.Path)
	}
	for _, u := range no.Delete {
		log.V(1).Infof("Received delete path: %s", prototext.Format(u))
		switch {
		case len(u.Elem) > 0:
		case len(u.Element) > 0: //nolint:staticcheck //lint:ignore SA1019 gnmi cache currently doesn't support PathElem for deletions.
			// Since gNMI still sends delete paths using the deprecated Element field, we need to translate it into path-elems first.
			// We also need to strip the first element for origin.
			//nolint:staticcheck //lint:ignore SA1019 gnmi cache currently doesn't support PathElem for deletions.
			elems, err := pathTranslator.PathElem(u.Element[1:])
			if err != nil {
				return fmt.Errorf("interfaceTask: failed to translate delete path: %s", prototext.Format(u))
			}
			u.Elem = elems
		default:
			return fmt.Errorf("Unhandled: delete at root: %s", prototext.Format(u))
		}
		interfacePathHandler(u)
	}
	return nil
}

func interfacePathHandler(path *gpb.Path) {
	switch {
	case matchingPath(path, descriptionPaths):
		log.V(1).Infof("Received update path: %s", prototext.Format(path))
		interfacePendingEvents[&intfDescriptionReactor] = true
	case matchingPath(path, enabledPaths), matchingPath(path, namePaths), matchingPath(path, ipv4AddressPaths), matchingPath(path, prefixLengthPaths):
		interfacePendingEvents[&interfaceReactor] = true
	default:
		log.V(1).Infof("interfaceTask: update path received isn't matched by any handlers: %s", prototext.Format(path))
	}
}

var (
	intfDescriptionReactor = func(intendedRoot *config.Device) error {
		for intfName, intf := range intendedRoot.Interface {
			curIntf, ok := interfaceAppliedRoot.Interface[intfName]
			if !ok {
				var err error
				if curIntf, err = interfaceAppliedRoot.NewInterface(intfName); err != nil {
					return fmt.Errorf("interfaceTask: %v", err)
				}
			}
			curIntf.Description = intf.Description
		}

		for intfName := range interfaceAppliedRoot.Interface {
			if _, ok := intendedRoot.Interface[intfName]; !ok {
				delete(interfaceAppliedRoot.Interface, intfName)
			}
		}

		return nil
	}

	interfaceReactor = func(intendedRoot *config.Device) error {
		for intfName, intf := range intendedRoot.Interface {
			curIntf, ok := interfaceAppliedRoot.Interface[intfName]
			if !ok {
				var err error
				if curIntf, err = interfaceAppliedRoot.NewInterface(intfName); err != nil {
					return fmt.Errorf("interfaceTask: %v", err)
				}
			}
			curIntf.Name = intf.Name
			curIntf.Enabled = intf.Enabled
			// TODO(wenbli): Handle more than subinterface 0th index.
			if subintf, ok := intf.Subinterface[0]; ok {
				for addrKey, addr := range subintf.GetIpv4().Address {
					curSubintf := curIntf.GetOrCreateSubinterface(0)
					curAddr := curSubintf.GetOrCreateIpv4().GetOrCreateAddress(addrKey)
					curAddr.Ip = addr.Ip
					curAddr.PrefixLength = addr.PrefixLength
				}
			}
			updateInterface(curIntf)
		}

		for intfName := range interfaceAppliedRoot.Interface {
			if _, ok := intendedRoot.Interface[intfName]; !ok {
				deleteInterface(intfName)
			}
		}

		return nil
	}
)
