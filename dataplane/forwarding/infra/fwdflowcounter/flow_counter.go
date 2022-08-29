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

// Package fwdflowcounter implements the functionality of Flow Counters.
// Flow Counters are objects that are managed by the client. They track the
// number of packets and octets that match each row in the flow table.
package fwdflowcounter

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

var counterList = []fwdpb.CounterId{
	fwdpb.CounterId_FLOW_COUNTER_OCTETS,
	fwdpb.CounterId_FLOW_COUNTER_PACKETS,
}

// FlowCounter implements the per flow counter.
type FlowCounter struct {
	fwdobject.Base
}

// String method formats the FlowCounter.
func (fc *FlowCounter) String() string {
	return fc.BaseInfo()
}

// Acquire acquires a reference to a FlowCounter.
func Acquire(ctx *fwdcontext.Context, id *fwdpb.FlowCounterId) (*FlowCounter, error) {
	if id == nil {
		return nil, errors.New("fwdflowcounter: Acquire failed, no flow counter specified")
	}
	obj, err := ctx.Objects.Acquire(id.GetObjectId())
	if err != nil {
		return nil, err
	}
	if fc, ok := obj.(*FlowCounter); ok {
		return fc, nil
	}

	// Release the object if there was an error
	_ = obj.Release(false /*forceCleanup*/)
	return nil, fmt.Errorf("fwdflowcounter: Acquire failed, %v is not a FlowCounter", id)
}

// Release releases a reference to a FlowCounter.
func Release(fc *FlowCounter) error {
	if fc == nil {
		return errors.New("fwdflowcounter: Release failed, no flow counter specified")
	}
	return fc.Release(false /*forceCleanup*/)
}

// Process is used to process a packet and increment the octet and packet counts.
func (fc *FlowCounter) Process(octetCount, packetCount uint32) error {
	fc.Increment(fwdpb.CounterId_FLOW_COUNTER_OCTETS, octetCount)
	fc.Increment(fwdpb.CounterId_FLOW_COUNTER_PACKETS, packetCount)
	return nil
}

// Query is used to read out the packet and octet counts stored in flow counters.
func (fc *FlowCounter) Query() (*fwdpb.FlowCounter, error) {
	ctrs := fc.Counters()
	octets, berr := ctrs[fwdpb.CounterId_FLOW_COUNTER_OCTETS]
	packets, perr := ctrs[fwdpb.CounterId_FLOW_COUNTER_PACKETS]

	if berr == false || perr == false {
		return nil, errors.New("fwdflowcounter: FlowCounterQuery failed, counter not available")
	}

	retfc := &fwdpb.FlowCounter{
		Id: &fwdpb.FlowCounterId{
			ObjectId: &fwdpb.ObjectId{
				Id: proto.String(string(fc.ID())),
			},
		},
		Octets:  &octets.Value,
		Packets: &packets.Value,
	}
	return retfc, nil
}

// New is used to create a per flow counter, which will store byte and packet counts.
func New(ctx *fwdcontext.Context, req *fwdpb.FlowCounterCreateRequest) (*FlowCounter, error) {
	if ctx == nil {
		return nil, errors.New("fwdflowcounter: New failed, no context specified")
	}
	if req == nil {
		return nil, errors.New("fwdflowcounter: New failed, no flow counter specified")
	}

	fc := &FlowCounter{}
	if err := fc.InitCounters("", "", counterList...); err != nil {
		return nil, fmt.Errorf("Failed to init FlowCounter %v: %v", fc.ID(), err)
	}

	if err := ctx.Objects.Insert(fc, req.GetId().GetObjectId()); err != nil {
		return nil, err
	}
	return fc, nil
}
