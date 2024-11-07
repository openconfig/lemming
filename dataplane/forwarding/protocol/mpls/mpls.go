// Copyright 2024 Google LLC
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

package mpls

import (
	"github.com/openconfig/lemming/dataplane/forwarding/protocol"
	"github.com/openconfig/lemming/dataplane/forwarding/util/frame"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

type mpls struct {
	protocol.Handler
	desc *protocol.Desc
}

func parseMPLS(f *frame.Frame, desc *protocol.Desc) (protocol.Handler, fwdpb.PacketHeaderId, error) {
	h := &mpls{
		desc: desc,
	}

	for {
		label, err := f.Peek(0, 3)
		if err != nil {
			return nil, 0, err
		}
		tc := label.BitField(20, 3)
		bot := label.BitField(23, 1)
		if bot.Value() == 1 {
			break
		}
	}

	return h, fwdpb.PacketHeaderId_PACKET_HEADER_ID_OPAQUE, nil
}

func add(id fwdpb.PacketHeaderId, _ *protocol.Desc) (protocol.Handler, error) {
	return nil, nil
}

func init() {
	protocol.Register(fwdpb.PacketHeaderId_PACKET_HEADER_ID_MPLS, parseMPLS, add)
}
