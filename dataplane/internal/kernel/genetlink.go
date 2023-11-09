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

package kernel

import (
	"fmt"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
)

// GenetlinkPort is connect to a netlink socket that be written to.
type GenetlinkPort struct {
	conn     *genetlink.Conn
	familyID uint16
}

// NewGenetlinkPort creates netlink socket for the given family and multicast group.
func NewGenetlinkPort(family, group string) (*GenetlinkPort, error) {
	conn, err := genetlink.Dial(nil)
	if err != nil {
		return nil, err
	}
	fam, err := conn.GetFamily(family)
	if err != nil {
		return nil, fmt.Errorf("could not find %v family", family)
	}
	grpID := -1
	for _, grp := range fam.Groups {
		if grp.Name == group {
			grpID = int(grp.ID)
			break
		}
	}
	if grpID == -1 {
		return nil, fmt.Errorf("could not find multicast group in the %v family", family)
	}
	if err := conn.JoinGroup(uint32(grpID)); err != nil {
		return nil, err
	}
	return &GenetlinkPort{
		conn: conn,
	}, nil
}

// Writes writes a layer2 frame to the port.
func (p GenetlinkPort) Write(frame []byte) error {
	// TODO: Support src, dst, and context
	data, err := (&NLPacket{
		payload:      frame,
		srcIfIndex:   1,
		dstIfIndex:   1,
		contextValue: 1,
	}).Encode()
	if err != nil {
		return err
	}

	_, err = p.conn.Send(genetlink.Message{
		Header: genetlink.Header{
			Command: 1,
			Version: 1,
		},
		Data: data,
	}, p.familyID, 0)
	return err
}

// NLPacket contains a packet data.
type NLPacket struct {
	srcIfIndex   int16
	dstIfIndex   int16
	contextValue uint32
	payload      []byte
}

// Constants sourced from https://github.com/sonic-net/sonic-pins/blob/main/p4rt_app/sonic/receive_genetlink.cc#L32
const (
	AttrDstIfIndex uint16 = iota
	AttrSrcIfIndex
	AttrContextValue
	AttrPayload
)

// Encode encodes the packet into a netlink-compatible byte slice.
func (nl *NLPacket) Encode() ([]byte, error) {
	enc := netlink.NewAttributeEncoder()
	enc.Int16(AttrSrcIfIndex, nl.srcIfIndex)
	enc.Int16(AttrSrcIfIndex, nl.dstIfIndex)
	enc.Uint32(AttrContextValue, nl.contextValue)
	enc.Bytes(AttrPayload, nl.payload)
	return enc.Encode()
}
