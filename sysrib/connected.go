// Copyright 2023 Google LLC
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

package sysrib

import (
	"context"
	"fmt"

	log "github.com/golang/glog"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
	"github.com/openconfig/ygnmi/ygnmi"
)

// monitorConnectedIntfs starts a gothread to check for connected prefixes from
// connected interfaces and adds them to the sysrib. It returns an error if
// there is an error before monitoring can begin.
//
// - gnmiServerAddr is the address of the central gNMI datastore.
// - target is the name of the gNMI target.
func (s *Server) monitorConnectedIntfs(ctx context.Context, yclient *ygnmi.Client) error {
	b := &ocpath.Batch{}
	b.AddPaths(
		ocpath.Root().InterfaceAny().Enabled().State().PathStruct(),
		ocpath.Root().InterfaceAny().Ifindex().State().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv4().AddressAny().Ip().State().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv4().AddressAny().PrefixLength().State().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv6().AddressAny().Ip().State().PathStruct(),
		ocpath.Root().InterfaceAny().Subinterface(0).Ipv6().AddressAny().PrefixLength().State().PathStruct(),
	)

	interfaceWatcher := ygnmi.Watch(
		ctx,
		yclient,
		b.State(),
		func(root *ygnmi.Value[*oc.Root]) error {
			rootVal, ok := root.Val()
			if !ok {
				return ygnmi.Continue
			}
			for name, intf := range rootVal.Interface {
				if intf.Enabled != nil {
					if intf.Ifindex != nil {
						ifindex := intf.GetIfindex()
						s.setInterface(ctx, name, int32(ifindex), intf.GetEnabled())
						// TODO(wenbli): Support other VRFs.
						if subintf := intf.GetSubinterface(0); subintf != nil {
							for _, addr := range subintf.GetOrCreateIpv4().Address {
								if addr.Ip != nil && addr.PrefixLength != nil {
									if err := s.addInterfacePrefix(ctx, name, int32(ifindex), fmt.Sprintf("%s/%d", addr.GetIp(), addr.GetPrefixLength()), fakedevice.DefaultNetworkInstance); err != nil {
										log.Warningf("adding interface prefix failed: %v", err)
									}
								}
							}
							for _, addr := range subintf.GetOrCreateIpv6().Address {
								if addr.Ip != nil && addr.PrefixLength != nil {
									if err := s.addInterfacePrefix(ctx, name, int32(ifindex), fmt.Sprintf("%s/%d", addr.GetIp(), addr.GetPrefixLength()), fakedevice.DefaultNetworkInstance); err != nil {
										log.Warningf("adding interface prefix failed: %v", err)
									}
								}
							}
						}
					}
				}
			}
			return ygnmi.Continue
		},
	)

	// TODO(wenbli): Support interface removal.
	go func() {
		if _, err := interfaceWatcher.Await(); err != nil {
			log.Warningf("Sysrib interface watcher has stopped: %v", err)
		}
	}()
	return nil
}
