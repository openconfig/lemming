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

// testagentlocal is a demo agent that uses ygnmi to receive data from the
// central datastore.
package testagentlocal

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/coalesce"
	"github.com/openconfig/gnmi/ctree"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/gnmit"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/encoding/prototext"
)

func InterfaceTask(target string) gnmit.Task {
	return gnmit.Task{
		Run: interfaceTask,
		Paths: []ygnmi.PathStruct{
			ocpath.Root().InterfaceAny(),
		},
		Prefix: &gpb.Path{
			Origin: "openconfig",
			Target: target,
		},
	}
}

func updateInterface(*oc.Interface) error {
	return nil
}

func deleteInterface(name string) error {
	return nil
}

var (
	enabledPaths, descriptionPaths, namePaths, ipv4AddressPaths, prefixLengthPaths *gpb.Path

	appliedMu sync.Mutex
	// appliedRoot is the SoT for BGP applied configuration. It is maintained locally by the task.
	appliedRoot *oc.Root
)

func initGlobalVars() error {
	appliedRoot = &oc.Root{}

	return initInterfacePaths()
}

func initInterfacePaths() error {
	interfacePath := ocpath.Root().InterfaceAny()
	var err error
	enabledPaths, _, err = ygnmi.ResolvePath(interfacePath.Enabled().Config().PathStruct())
	if err != nil {
		return fmt.Errorf("interfaceTask failed to initialize due to error: %v", err)
	}
	descriptionPaths, _, err = ygnmi.ResolvePath(interfacePath.Description().Config().PathStruct())
	if err != nil {
		return fmt.Errorf("interfaceTask failed to initialize due to error: %v", err)
	}
	namePaths, _, err = ygnmi.ResolvePath(interfacePath.Name().Config().PathStruct())
	if err != nil {
		return fmt.Errorf("interfaceTask failed to initialize due to error: %v", err)
	}
	ipv4AddressPaths, _, err = ygnmi.ResolvePath(interfacePath.SubinterfaceAny().Ipv4().AddressAny().Ip().Config().PathStruct())
	if err != nil {
		return fmt.Errorf("interfaceTask failed to initialize due to error: %v", err)
	}
	prefixLengthPaths, _, err = ygnmi.ResolvePath(interfacePath.SubinterfaceAny().Ipv4().AddressAny().PrefixLength().Config().PathStruct())
	if err != nil {
		return fmt.Errorf("interfaceTask failed to initialize due to error: %v", err)
	}

	return nil
}

func interfaceTask(getIntendedConfig func() *oc.Root, q gnmit.Queue, update gnmit.UpdateFn, target string, remove func()) error {
	if err := initGlobalVars(); err != nil {
		return err
	}

	// updateAppliedConfig computes the diff between a previous applied
	// configuration and the current SoT, and sends the updates to the
	// central DB.
	updateAppliedConfig := func(prevApplied *oc.Root) bool {
		appliedMu.Lock()
		defer appliedMu.Unlock()
		no, err := ygot.Diff(prevApplied, appliedRoot)
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
		} else {
			log.V(1).Info("No applied config updates for interfaceTask")
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

			pendingEvents, err := interfaceNotificationHandler(no)
			if err != nil {
				log.Fatalf("interfaceTask: error at notification handler: %v", err)
			}

			var updateAppliedRoot bool
			for _, triggered := range pendingEvents {
				if triggered {
					updateAppliedRoot = true
				}
			}
			if !updateAppliedRoot {
				continue
			}

			var prevApplied *oc.Root
			appliedMu.Lock()
			prevAppliedGS, err := ygot.DeepCopy(appliedRoot)
			if err != nil {
				log.Fatalf("interfaceTask: Could not copy applied configuration: %v", err)
			}
			prevApplied = prevAppliedGS.(*oc.Root)
			appliedMu.Unlock()

			for reactor, triggered := range pendingEvents {
				if triggered {
					if err := (*reactor)(getIntendedConfig()); err != nil {
						log.Errorf("interfaceTask reactor: %v", err)
					}
					pendingEvents[reactor] = false
				}
			}

			if success := updateAppliedConfig(prevApplied); !success {
				log.Errorf("interfaceTask: updating applied configuration failed")
			}
		}
	}()

	return nil
}

func interfaceNotificationHandler(no *gpb.Notification) (map[*func(*oc.Root) error]bool, error) {
	pendingEvents := map[*func(*oc.Root) error]bool{}

	for _, u := range no.Update {
		interfacePathHandler(pendingEvents, u.Path)
	}
	for _, u := range no.Delete {
		log.V(2).Infof("Received delete path: %s", prototext.Format(u))
		switch {
		case len(u.Elem) > 0:
		case len(u.Element) > 0: //nolint:staticcheck //lint:ignore SA1019 gnmi cache currently doesn't support PathElem for deletions.
			// Since gNMI still sends delete paths using the deprecated Element field, we need to translate it into path-elems first.
			// We also need to strip the first element for origin.
			//nolint:staticcheck //lint:ignore SA1019 gnmi cache currently doesn't support PathElem for deletions.
			elems, err := fakedevice.PathTranslator.PathElem(u.Element[1:])
			if err != nil {
				return nil, fmt.Errorf("interfaceTask: failed to translate delete path: %s", prototext.Format(u))
			}
			u.Elem = elems
		default:
			return nil, fmt.Errorf("Unhandled: delete at root: %s", prototext.Format(u))
		}
		interfacePathHandler(pendingEvents, u)
	}
	return pendingEvents, nil
}

// matchingPath returns true iff the path matches the given matcher path in
// length and in values; wildcards are allowed in the matcher path.
func matchingPath(path, matcher *gpb.Path) bool {
	return len(path.Elem) == len(matcher.Elem) && util.PathMatchesQuery(path, matcher)
}

// interfacePathHandler sets the pending events that should be triggered based on the input path.
func interfacePathHandler(pendingEvents map[*func(*oc.Root) error]bool, path *gpb.Path) {
	switch {
	case matchingPath(path, descriptionPaths):
		log.V(2).Infof("interfaceTask: Received update path: %s", prototext.Format(path))
		pendingEvents[&intfDescriptionReactor] = true
	case matchingPath(path, enabledPaths), matchingPath(path, namePaths), matchingPath(path, ipv4AddressPaths), matchingPath(path, prefixLengthPaths):
		log.V(2).Infof("interfaceTask: Received name or address/prefix path: %s", prototext.Format(path))
		pendingEvents[&interfaceReactor] = true
	default:
		log.V(2).Infof("interfaceTask: update path received isn't matched by any handlers: %s", prototext.Format(path))
	}
}

var (
	intfDescriptionReactor = func(intendedRoot *oc.Root) error {
		for intfName, intf := range intendedRoot.Interface {
			log.V(1).Infof("interfaceTask: adding new interface %q", intfName)
			curIntf, ok := appliedRoot.Interface[intfName]
			if !ok {
				var err error
				if curIntf, err = appliedRoot.NewInterface(intfName); err != nil {
					return fmt.Errorf("interfaceTask: %v", err)
				}
			}
			curIntf.Description = intf.Description
		}

		for intfName := range appliedRoot.Interface {
			log.V(1).Infof("interfaceTask: deleting interface %q", intfName)
			if _, ok := intendedRoot.Interface[intfName]; !ok {
				delete(appliedRoot.Interface, intfName)
			}
		}

		return nil
	}

	interfaceReactor = func(intendedRoot *oc.Root) error {
		for intfName, intf := range intendedRoot.Interface {
			curIntf, ok := appliedRoot.Interface[intfName]
			if !ok {
				var err error
				if curIntf, err = appliedRoot.NewInterface(intfName); err != nil {
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
			if err := updateInterface(curIntf); err != nil {
				return err
			}
		}

		for intfName := range appliedRoot.Interface {
			if _, ok := intendedRoot.Interface[intfName]; !ok {
				if err := deleteInterface(intfName); err != nil {
					return err
				}
			}
		}

		return nil
	}
)
