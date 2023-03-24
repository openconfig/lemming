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

// Package porttestutil contains routines used to create and manage ports for test
// cases.
package porttestutil

import (
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdport"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// DescTestPort builds a port descriptor for a test port using the specified
// unique name. Test ports are created using CPU ports.
func DescTestPort(name string) *fwdpb.PortDesc {
	desc := &fwdpb.PortDesc{
		PortType: fwdpb.PortType_PORT_TYPE_CPU_PORT,
		PortId:   fwdport.MakeID(fwdobject.NewID(name)),
	}
	cpu := &fwdpb.CPUPortDesc{
		QueueId: name,
	}
	desc.Port = &fwdpb.PortDesc_Cpu{
		Cpu: cpu,
	}
	return desc
}

// CreateTestPort creates a test port using the specified name.
func CreateTestPort(t *testing.T, ctx *fwdcontext.Context, name string) fwdport.Port {
	ctx.Lock()
	defer ctx.Unlock()

	desc := DescTestPort(name)
	port, err := fwdport.New(desc, ctx)
	if err != nil {
		t.Fatalf("Port create failed: %v.", err)
	}
	port.State(&fwdpb.PortInfo{Laser: fwdpb.PortLaserState_PORT_LASER_STATE_ENABLED})
	return port
}

// PortMap returns a map of ports that have been selected by the port group.
// Selected ports can be identified by the non-zero tx counters.
func PortMap(ports []fwdport.Port) map[fwdobject.ID]bool {
	pm := make(map[fwdobject.ID]bool)
	for _, p := range ports {
		counters := p.Counters()
		if counter, ok := counters[fwdpb.CounterId_COUNTER_ID_TX_PACKETS]; ok && counter.Value != 0 {
			pm[p.ID()] = true
		}
	}
	return pm
}
