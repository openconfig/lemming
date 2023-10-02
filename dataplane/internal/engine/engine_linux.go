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

//go:build linux

package engine

import (
	"context"

	"github.com/vishvananda/netlink"

	log "github.com/golang/glog"
)

func (e *Engine) handleIPUpdates(ctx context.Context) {
	updCh := make(chan netlink.AddrUpdate)
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case upd := <-updCh:
				l, err := netlink.LinkByIndex(upd.LinkIndex)
				if err != nil {
					log.Warningf("failed to get link: %v", err)
					continue
				}
				e.ipToDevNameMu.Lock()
				if upd.NewAddr {
					log.Infof("added new ip %s to device %s", upd.LinkAddress.IP.String(), l.Attrs().Name)
					e.ipToDevName[upd.LinkAddress.IP.String()] = l.Attrs().Name
				} else {
					delete(e.ipToDevName, upd.LinkAddress.IP.String())
				}
				e.ipToDevNameMu.Unlock()
			case <-ctx.Done():
				close(doneCh)
				return
			}
		}
	}()

	netlink.AddrSubscribe(updCh, doneCh)
}
