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

package genetlink

// #cgo LDFLAGS: -lnl-3 -lnl-genl-3
// #cgo CFLAGS: -I/usr/include/libnl3
// #include "genetlink.h"
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"io"
	"unsafe"

	log "github.com/golang/glog"

	"github.com/mdlayher/genetlink"

	"github.com/openconfig/lemming/dataplane/kernel"
	pktiopb "github.com/openconfig/lemming/dataplane/proto/packetio"
	"github.com/openconfig/lemming/dataplane/standalone/pkthandler/pktiohandler"
)

// GenetlinkPort is connect to a netlink socket that be written to.
type GenetlinkPort struct {
	sock     unsafe.Pointer
	familyID int
}

// NewGenetlinkPort creates netlink socket for the given family and multicast group.
func New(msg *pktiopb.HostPortControlMessage) (pktiohandler.PortIO, error) {
	log.Errorf("creating genl port: %s %s", msg.GetGenetlink().GetFamily(), msg.GetGenetlink().GetGroup())

	cFamily := C.CString(msg.GetGenetlink().GetFamily())
	defer C.free(unsafe.Pointer(cFamily))
	cGroup := C.CString(msg.GetGenetlink().GetGroup())
	defer C.free(unsafe.Pointer(cGroup))

	conn, err := genetlink.Dial(nil)
	if err != nil {
		return nil, err
	}
	fam, err := conn.GetFamily(msg.GetGenetlink().GetFamily())
	if err != nil {
		return nil, err
	}
	familyID := -1
	for _, grp := range fam.Groups {
		if grp.Name == msg.GetGenetlink().GetGroup() {
			familyID = int(grp.ID)
		}
	}
	if familyID == -1 {
		return nil, fmt.Errorf("failed to find multicast group")
	}

	sockAddr := C.create_port(cFamily, cGroup)
	if sockAddr == nil {
		return nil, fmt.Errorf("failed to create port")
	}

	log.Errorf("creating genl port")
	return &GenetlinkPort{
		sock:     unsafe.Pointer(sockAddr),
		familyID: familyID,
	}, nil
}

// Writes writes a layer2 frame to the port.
func (p GenetlinkPort) Write(frame []byte, md *kernel.PacketMetadata) (int, error) {
	log.Errorf("writing genl packet: %x", frame)

	packet := C.CBytes(frame)
	defer C.free(packet)

	res := C.send_packet(p.sock, C.int(p.familyID), packet, C.uint(uint32(len(frame))), C.int(md.SrcIfIndex), C.int(md.DstIfIndex), C.uint(md.Context))
	if res < 0 {
		return 0, fmt.Errorf("failed to write packet")
	}

	return len(frame), nil
}

// Read is not implemented.
func (p GenetlinkPort) Read([]byte) (int, error) {
	return 0, io.EOF
}

// Delete closes the netlink connection.
func (p GenetlinkPort) Delete() error {
	C.delete_port(p.sock)
	return nil
}

func (t GenetlinkPort) SetAdminState(up bool) error {
	return nil
}

func init() {
	pktiohandler.Register(pktiopb.PortType_PORT_TYPE_GENETLINK, New)
}
