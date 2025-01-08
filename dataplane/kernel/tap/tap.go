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

package tap

import (
	"os"

	"github.com/vishvananda/netlink"

	"github.com/openconfig/lemming/dataplane/kernel"
	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
	"github.com/openconfig/lemming/dataplane/standalone/pkthandler/pktiohandler"
)

// New creates a new tap interface.
func New(msg *pktiopb.HostPortControlMessage) (pktiohandler.PortIO, error) {
	tap := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name: msg.GetNetdev().GetName(),
		},
		Mode:       netlink.TUNTAP_MODE_TAP,
		Flags:      netlink.TUNTAP_MULTI_QUEUE_DEFAULTS,
		NonPersist: true,
		Queues:     1,
	}
	if err := netlink.LinkAdd(tap); err != nil {
		return nil, err
	}
	return &TapInterface{
		name:    msg.GetNetdev().GetName(),
		File:    tap.Fds[0],
		ifIndex: tap.Index,
	}, nil
}

// TapInterface is Linux tap interface.
type TapInterface struct {
	name string
	*os.File
	ifIndex int
}

// Delete deletes the tap interface.
func (t *TapInterface) Delete() error {
	l, err := netlink.LinkByName(t.name)
	if err != nil {
		return err
	}
	return netlink.LinkDel(l)
}

func (t *TapInterface) Write(frame []byte, _ *kernel.PacketMetadata) (int, error) {
	return t.File.Write(frame)
}

func (t *TapInterface) IfIndex() int {
	return t.ifIndex
}

func (t *TapInterface) SetAdminState(up bool) error {
	l, err := netlink.LinkByName(t.name)
	if err != nil {
		return err
	}
	if up {
		return netlink.LinkSetUp(l)
	}
	return netlink.LinkSetDown(l)
}

func init() {
	pktiohandler.Register(pktiopb.PortType_PORT_TYPE_NETDEV, New)
}
