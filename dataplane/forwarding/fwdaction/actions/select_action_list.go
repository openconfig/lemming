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

package actions

import (
	"fmt"
	"hash/crc32"
	"math"
	"math/rand"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	"github.com/openconfig/lemming/dataplane/forwarding/util/hash/crc16"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A selectActionList is an action that selects an action list from a set of
// such lists, using the specified algorithm and packet fields.
type selectActionList struct {
	fwdobject.Base
	fields    []fwdpacket.FieldID // packet fields used to create a packet hash
	set       []fwdaction.Actions // set of action lists
	weights   []int               // weights associates the set of action lists
	weightSum uint64
	hashFn    func(key []byte, max uint64) int                 // function used to hash a set of bytes
	hash      fwdpb.SelectActionListActionDesc_SelectAlgorithm // hash algorithm used to select the action list
}

// String returns the action as a formatted string.
func (s *selectActionList) String() string {
	desc := fmt.Sprintf("Type=%v;<Fields=%v>;<Hash=%v>;%v;", fwdpb.ActionType_ACTION_TYPE_SELECT_ACTION_LIST, s.fields, s.hash, s.BaseInfo())
	for _, a := range s.set {
		desc += fmt.Sprintf("<%v>;", a.String())
	}
	return desc
}

// Cleanup releases references held by the action lists .
func (s *selectActionList) Cleanup() {
	for _, a := range s.set {
		a.Cleanup()
	}
	s.set = nil
}

// Process processes the packet by selecting an action from a list of actions.
func (s *selectActionList) Process(packet fwdpacket.Packet, counters fwdobject.Counters) (fwdaction.Actions, fwdaction.State) {
	if s.hashFn == nil {
		counters.Increment(fwdpb.CounterId_COUNTER_ID_ERROR_PACKETS, 1)
		counters.Increment(fwdpb.CounterId_COUNTER_ID_ERROR_OCTETS, 1)
		return nil, fwdaction.DROP
	}

	var key []byte
	for _, id := range s.fields {
		if f, err := packet.Field(id); err == nil {
			key = append(key, f...)
		}
	}

	h := s.hashFn(key, s.weightSum)
	var index int
	for i, w := range s.weights {
		index = i
		if w < h {
			break
		}
	}
	a := s.set[index]
	packet.Logf(fwdpacket.LogDebugMessage, "hash selected %v", a)
	return a, fwdaction.CONTINUE
}

// hashCRC32 computes the CRC32 checksum of the key.
func hashCRC32(key []byte, max uint64) int {
	rand := crc32.ChecksumIEEE(key)
	return int(uint64(rand) * max / math.MaxUint32)
}

// hashCRC16 computes the CRC16 checksum of the key.
func hashCRC16(key []byte, max uint64) int {
	rand := crc16.ChecksumANSI(key)
	return int(uint64(rand) * max / math.MaxUint16)
}

// random selects a random index.
func random(_ []byte, max uint64) int {
	//nolint:gosec
	return int(rand.Intn(int(max)))
}

// A selectActionListBuilder builds selectActionList actions.
type selectActionListBuilder struct{}

// init registers a builder for the selectActionList action type.
func init() {
	fwdaction.Register(fwdpb.ActionType_ACTION_TYPE_SELECT_ACTION_LIST, &selectActionListBuilder{})
}

// Build creates a new selectActionList action.
func (*selectActionListBuilder) Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (fwdaction.Action, error) {
	sal, ok := desc.Action.(*fwdpb.ActionDesc_Select)
	if !ok {
		return nil, fmt.Errorf("actions: Build for selectActionList action failed, missing desc")
	}

	s := &selectActionList{
		hash: sal.Select.GetSelectAlgorithm(),
	}

	// Setup the fields for the packet hash.
	s.fields = make([]fwdpacket.FieldID, 0, len(sal.Select.GetFieldIds()))
	for _, field := range sal.Select.GetFieldIds() {
		s.fields = append(s.fields, fwdpacket.NewFieldID(field))
	}

	// Setup the packet hash function.
	switch s.hash {
	case fwdpb.SelectActionListActionDesc_SELECT_ALGORITHM_CRC32:
		s.hashFn = hashCRC32
	case fwdpb.SelectActionListActionDesc_SELECT_ALGORITHM_CRC16:
		s.hashFn = hashCRC16
	case fwdpb.SelectActionListActionDesc_SELECT_ALGORITHM_RANDOM:
		s.hashFn = random
	default:
		return nil, fmt.Errorf("actions: Unable to find select function %v", s.hash)
	}

	allZeros := true
	for _, l := range sal.Select.GetActionLists() {
		a, err := fwdaction.NewActions(l.GetActions(), ctx)
		if err != nil {
			return nil, fmt.Errorf("actions: Unable to create actions %v, err %v", l, err)
		}
		s.set = append(s.set, a)
		if l.GetWeight() != 0 {
			allZeros = false
		}
		s.weights = append(s.weights, int(l.GetWeight()+s.weightSum))
		s.weightSum += l.GetWeight()
	}
	if allZeros {
		s.weightSum = uint64(len(s.weights))
		for i := range s.weights {
			s.weights[i] = i + 1
		}
	}
	return s, nil
}
