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
	"github.com/openconfig/ygnmi/ygnmi"

	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/lemming/gnmi/oc/ocpath"
)

type interfaceRef struct {
	Name         string
	SubInterface uint32
}

// monitorConnectedIntfs starts a gothread to check for connected prefixes from
// connected interfaces and adds them to the sysrib. It returns an error if
// there is an error before monitoring can begin.
func (s *Server) monitorConnectedIntfs(ctx context.Context, yclient *ygnmi.Client) error {
	b := ygnmi.NewBatch[*oc.Root](ocpath.Root().State())
	b.AddPaths(
		ocpath.Root().InterfaceAny().Enabled().State(),
		ocpath.Root().InterfaceAny().Ifindex().State(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Ipv4().AddressAny().Ip().State(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Ipv4().AddressAny().PrefixLength().State(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Ipv6().AddressAny().Ip().State(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Ipv6().AddressAny().PrefixLength().State(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Ifindex().State(),
		ocpath.Root().InterfaceAny().SubinterfaceAny().Enabled().State(),
		ocpath.Root().NetworkInstanceAny().InterfaceAny().State(),
	)

	prevIntfs := map[connectedRoute]struct{}{}
	log.Infof("starting connected route watcher")
	interfaceWatcher := ygnmi.Watch(
		ctx,
		yclient,
		b.Query(),
		func(intfs *ygnmi.Value[*oc.Root]) error {
			root, ok := intfs.Val()
			if !ok {
				return ygnmi.Continue
			}

			currentIntfs := map[connectedRoute]struct{}{}
			niIntfMap := map[interfaceRef]string{}

			for niName, ni := range root.NetworkInstance {
				for _, intf := range ni.Interface {
					if intf.Interface == nil || intf.Subinterface == nil {
						continue
					}
					niIntfMap[interfaceRef{Name: intf.GetInterface(), SubInterface: intf.GetSubinterface()}] = niName
				}
			}

			for name, intf := range root.Interface {
				log.Infof("got interface update: %v intf.Enabled %v,  intf.Ifindex %v", name, intf.Enabled, intf.Ifindex)
				if intf.Enabled != nil && intf.Ifindex != nil {
					for subintfIdx, subintf := range intf.Subinterface {
						ifindex := subintf.GetIfindex()
						if subintf.Ifindex == nil {
							ifindex = intf.GetIfindex()
						}
						s.setInterface(ctx, name, subintfIdx, int32(ifindex), subintf.GetEnabled() && intf.GetEnabled())
						niName := fakedevice.DefaultNetworkInstance
						if stateNIName := niIntfMap[interfaceRef{Name: intf.GetName(), SubInterface: subintfIdx}]; stateNIName != "" { // Interfaces are in the default network instance implicitly.
							niName = stateNIName
						}

						for _, addr := range subintf.GetOrCreateIpv4().Address {
							if addr.Ip != nil && addr.PrefixLength != nil {
								connected := connectedRoute{
									name:    name,
									ifindex: int32(ifindex),
									prefix:  fmt.Sprintf("%s/%d", addr.GetIp(), addr.GetPrefixLength()),
									niName:  niName,
								}
								if err := s.setConnectedRoute(ctx, connected, false); err != nil {
									log.Warningf("adding connected route failed: %v", err)
								} else {
									currentIntfs[connected] = struct{}{}
								}
							}
						}
						for _, addr := range subintf.GetOrCreateIpv6().Address {
							if addr.Ip != nil && addr.PrefixLength != nil {
								connected := connectedRoute{
									name:    name,
									ifindex: int32(ifindex),
									prefix:  fmt.Sprintf("%s/%d", addr.GetIp(), addr.GetPrefixLength()),
									niName:  niName,
								}
								if err := s.setConnectedRoute(ctx, connected, false); err != nil {
									log.Warningf("adding connected route failed: %v", err)
								} else {
									currentIntfs[connected] = struct{}{}
								}
							}
						}
					}
				}

			}
			for connected := range prevIntfs {
				if _, ok := currentIntfs[connected]; !ok {
					if err := s.setConnectedRoute(ctx, connected, true); err != nil {
						log.Warningf("deleting connected route failed: %v", err)
					}
				}
			}
			prevIntfs = currentIntfs

			return ygnmi.Continue
		},
	)

	go func() {
		if _, err := interfaceWatcher.Await(); err != nil {
			log.Warningf("Sysrib interface watcher has stopped: %v", err)
		}
	}()
	return nil
}
