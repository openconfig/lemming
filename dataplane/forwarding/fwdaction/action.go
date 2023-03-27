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

// Package fwdaction contains routines and types to manage forwarding actions.
//
// An Action is an operation that can be performed on a packet.
//
// Various forwarding behaviors within Lucius are implemented by types
// satisfying interface Action. Several actions may contain references
// to other forwarding objects. These actions must implement interface
// fwdobject.Composite.
//
// Actions is a list of actions ordered in their expected order of execution.
// It is always embedded in another object, and hence does not implement
// fwdobject.Object.
package fwdaction

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A State is the processing state of a packet.
type State int

// Enumerate the possible states during packet processing.
const (
	DROP     State = iota // Packet is being dropped.
	CONTINUE              // Packet is being processed further.
	CONSUME               // Packet is being consumed.
	OUTPUT                // Packet should be processed to the output port immediately.
	EVALUATE              // Packet is processed with all actions marked "onEvaluate"
)

// CounterList is a set of counters incremented by various actions.
var CounterList = []fwdpb.CounterId{
	fwdpb.CounterId_COUNTER_ID_ERROR_PACKETS,
	fwdpb.CounterId_COUNTER_ID_ERROR_OCTETS,
	fwdpb.CounterId_COUNTER_ID_DROP_PACKETS,
	fwdpb.CounterId_COUNTER_ID_DROP_OCTETS,
	fwdpb.CounterId_COUNTER_ID_RATELIMIT_PACKETS,
	fwdpb.CounterId_COUNTER_ID_RATELIMIT_OCTETS,
	fwdpb.CounterId_COUNTER_ID_MIRROR_PACKETS,
	fwdpb.CounterId_COUNTER_ID_MIRROR_OCTETS,
	fwdpb.CounterId_COUNTER_ID_MIRROR_ERROR_PACKETS,
	fwdpb.CounterId_COUNTER_ID_MIRROR_ERROR_OCTETS,
	fwdpb.CounterId_COUNTER_ID_ENCAP_ERROR_PACKETS,
	fwdpb.CounterId_COUNTER_ID_ENCAP_ERROR_OCTETS,
	fwdpb.CounterId_COUNTER_ID_DECAP_ERROR_PACKETS,
	fwdpb.CounterId_COUNTER_ID_DECAP_ERROR_OCTETS,
}

// An Action is an operation that can be performed on a packet.
type Action interface {
	// Process processes the packet using the specified counters.
	// Process returns an updated processing state for the packet.
	// It may also return new Actions that are used to process the packet.
	// Process assumes that counters are always valid (not nil).
	Process(packet fwdpacket.Packet, counters fwdobject.Counters) (Actions, State)

	// String formats the state of the action as a string.
	String() string
}

// A builder is an entity that can build an Action of the specified type.
type builder interface {
	// Build builds an action.
	Build(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (Action, error)
}

// builders is a map of builders for various types of actions.
var builders = make(map[fwdpb.ActionType]builder)

// Register registers a builder for the specified action type. Note that
// builders are expected to be registered serially during initialization.
func Register(atype fwdpb.ActionType, builder builder) {
	builders[atype] = builder
}

// New creates a new action based on the descriptor.
func New(desc *fwdpb.ActionDesc, ctx *fwdcontext.Context) (Action, error) {
	if desc == nil {
		return nil, errors.New("fwdaction: new failed, no descriptor for action")
	}
	if builder, ok := builders[desc.GetActionType()]; ok {
		return builder.Build(desc, ctx)
	}
	return nil, fmt.Errorf("fwdaction: new failed, no builder for action %s", desc)
}

// An ActionAttr contains atttributes of a given action.
type ActionAttr struct {
	onEvaluate bool // true if the action excutes only during evaluation
	action     Action
	hash       uint32 // A hash used to compare equality of actions
}

// Action returns the action associated with an attribute set.
func (a ActionAttr) Action() Action {
	return a.action
}

// String formats an ActionAttr
func (a ActionAttr) String() string {
	if a.action != nil {
		return fmt.Sprintf("<%v<onEvaluate=%v><Hash=%v>;>", a.action.String(), a.onEvaluate, a.hash)
	}
	return "Continue"
}

// NewActionAttr creates a type that describes attributes associated with the
// specified action.
func NewActionAttr(action Action, onEvaluate bool) *ActionAttr {
	a := &ActionAttr{
		action:     action,
		onEvaluate: onEvaluate,
		hash:       0, // seed value for the hash
	}

	h := fnv.New32a()
	h.Write([]byte(a.String()))
	a.hash = h.Sum32()
	return a
}

// Actions is a list of actions and their corresponding attributes.
type Actions []*ActionAttr

// NewActions builds a list of actions from a slice of descriptors.
func NewActions(descs []*fwdpb.ActionDesc, ctx *fwdcontext.Context) (Actions, error) {
	var actions Actions
	var err error
	defer func() {
		if err != nil {
			actions.Cleanup()
		}
	}()

	for _, desc := range descs {
		var action Action
		if action, err = New(desc, ctx); err != nil {
			return actions, err
		}
		if action != nil {
			actions = append(actions, NewActionAttr(action, desc.GetOnEvaluate()))
		}
	}
	return actions, nil
}

// Cleanup releases all references held by the actions.
func (actions Actions) Cleanup() {
	for _, a := range actions {
		if a.action == nil {
			continue
		}
		if composite, ok := a.action.(fwdobject.Composite); ok {
			composite.Cleanup()
		}
	}
}

// String returns the actions as a formatted string.
func (actions Actions) String() string {
	if len(actions) == 0 {
		return "Continue"
	}
	s := make([]string, len(actions))
	for id, a := range actions {
		s[id] = a.String()
	}
	return strings.Join(s, ";")
}

// IsEqual returns true if a two series of actions are equal.
func (actions Actions) IsEqual(arg Actions) bool {
	if len(actions) != len(arg) {
		return false
	}
	for pos, action := range actions {
		if action.hash != arg[pos].hash {
			return false
		}
	}
	return true
}

// maxActions is the maxumum numbers of actions to perform on a packet.
const maxActions = 100

// evaluatePacket processes a network packet with a series of actions to
// determine the state of the packet. It returns the next state and count of
// executed actions.
func evaluatePacket(packet fwdpacket.Packet, actions Actions, counters fwdobject.Counters) (State, int, error) {
	state := CONTINUE

	curr := 0 // index of the current executing action within actions
	exec := 0 // count of total number of executed actions
	for curr < len(actions) {
		a := actions[curr]
		curr++

		var next Actions
		packet.Logf(fwdpacket.LogDebugMessage, "curr action %v", a)
		next, state = a.action.Process(packet, counters)
		packet.Logf(fwdpacket.LogDebugMessage, "result state %v action %v", state, next)
		exec++

		// Any intervening "EVALUATE" is treated as a "CONTINUE".
		if state == EVALUATE {
			state = CONTINUE
		}

		switch {
		case exec == maxActions:
			return DROP, exec, errors.New("fwdaction: Actions %v exceed max number of actions")

		case state != CONTINUE:
			// Processing should terminate.
			return state, exec, nil

		case len(next) != 0:
			if curr == len(actions) {
				actions = next
			} else {
				pending := actions[curr:]
				actions = append(next, pending...)
			}
			curr = 0
		}
		// Continue next action by default.
	}
	return state, exec, nil
}

// ProcessPacket processes a network packet with a series of actions and
// determines the processing state of the packet.
func ProcessPacket(packet fwdpacket.Packet, actions Actions, counters fwdobject.Counters) (State, error) {
	if packet == nil {
		return CONTINUE, nil
	}
	state := CONTINUE

	var evaluate Actions // actions to be executed on an evaluate
	curr := 0            // index of the current executing action within actions
	exec := 0            // count of total number of executed actions
	for curr < len(actions) {
		a := actions[curr]
		curr++

		// If the action is tagged as onEvaluate, append it to the evaluate list.
		if a.onEvaluate {
			packet.Logf(fwdpacket.LogDebugMessage, "defer action %v", a)
			evaluate = append(evaluate, a)
			continue
		}

		var next Actions
		packet.Logf(fwdpacket.LogDebugMessage, "curr action %v", a)
		next, state = a.action.Process(packet, counters)
		packet.Logf(fwdpacket.LogDebugMessage, "result state %v action %v", state, next)
		exec++

		// If the state is Evaluate, then evaluate the packet by running the
		// accumulated evaluate list. These actions are then used to update the
		// number of executed action and the next state of processing.
		if state == EVALUATE {
			var count int
			var err error
			if state, count, err = evaluatePacket(packet, evaluate, counters); err != nil {
				return DROP, err
			}
			evaluate = nil
			exec += count
		}

		switch {
		case exec == maxActions:
			return DROP, errors.New("fwdaction: Actions %v exceed max number of actions")

		case state != CONTINUE:
			// Processing should terminate.
			return state, nil

		case len(next) != 0:
			if curr == len(actions) {
				actions = next
			} else {
				pending := actions[curr:]
				actions = append(next, pending...)
			}
			curr = 0
		}
		// Continue next action by default.
	}
	return state, nil
}
