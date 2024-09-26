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

	"github.com/openconfig/lemming/dataplane/kernel"
)

// GenetlinkPort is connect to a netlink socket that be written to.
type GenetlinkPort struct {
	socketIndex int
}

// NewGenetlinkPort creates netlink socket for the given family and multicast group.
func NewGenetlinkPort(family, group string) (*GenetlinkPort, error) {
	log.Errorf("creating genl port: %s %s", family, group)

	cFamily := C.CString(family)
	defer C.free(unsafe.Pointer(cFamily))
	cGroup := C.CString(group)
	defer C.free(unsafe.Pointer(cGroup))

	idx := C.create_port(cFamily, cGroup)
	if idx < 0 {
		return nil, fmt.Errorf("failed to create port: %d", idx)
	}

	return &GenetlinkPort{
		socketIndex: int(idx),
	}, nil
}

// Writes writes a layer2 frame to the port.
func (p GenetlinkPort) Write(frame []byte, md *kernel.PacketMetadata) (int, error) {
	log.Errorf("writing genl packet: %x", frame)

	packet := C.CBytes(frame)
	defer C.free(unsafe.Pointer(packet))

	res := C.send_packet(C.int(p.socketIndex), packet, C.uint(uint32(len(frame))), C.int(md.SrcIfIndex), C.int(md.DstIfIndex), C.uint(md.Context))
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
	return nil
}
