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

package fwdport

import (
	"fmt"
	"testing"

	"github.com/openconfig/lemming/dataplane/forwarding/fwdaction"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdcontext"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdobject"
	"github.com/openconfig/lemming/dataplane/forwarding/infra/fwdpacket"
	fwdpb "github.com/openconfig/lemming/proto/forwarding"
)

// A testPort is an Port that tracks its state.
type testPort struct {
	fwdobject.Base
	id        int
	cleaned   bool
	allocated bool
}

// Name returns the name of the port.
func (testPort) Name() string { return "" }

// String returns the state of the port as a formatted string.
func (port *testPort) String() string {
	return fmt.Sprintf("ID=%v, cleaned=%v, allocated=%v, ", port.id, port.cleaned, port.allocated)
}

// Type returns the type.
func (testPort) Type() fwdpb.PortType { return fwdpb.PortType_PORT_TYPE_UNSPECIFIED }

// Cleanup releases all held references (satisfies interface Composite).
func (port *testPort) Cleanup() {
	port.cleaned = true
}

// Update updates the port.
func (testPort) Update(*fwdpb.PortUpdateDesc) error {
	return nil
}

// Write writes a packet out.
func (testPort) Write(fwdpacket.Packet) (fwdaction.State, error) {
	return fwdaction.CONSUME, nil
}

// Actions returns the port actions of the specified type.
func (testPort) Actions(fwdpb.PortAction) fwdaction.Actions {
	return nil
}

// State manages the state of the port.
func (testPort) State(*fwdpb.PortInfo) (*fwdpb.PortStateReply, error) {
	return &fwdpb.PortStateReply{}, nil
}

// testBuilder builds test ports using a prebuilt port.
type testBuilder struct {
	port *testPort
}

// unregister unregisters a builder for the specified port type.
func unregister(portType fwdpb.PortType) {
	delete(builders, portType)
}

// Build uses the prebuilt port.
func (builder *testBuilder) Build(*fwdpb.PortDesc, *fwdcontext.Context) (Port, error) {
	builder.port.allocated = true
	return builder.port, nil
}

// newTestBuilder creates a new test builder and registers it.
func newTestBuilder(id int, portType fwdpb.PortType) *testBuilder {
	builder := &testBuilder{
		port: &testPort{
			id:        id,
			allocated: false,
			cleaned:   false,
		},
	}
	Register(portType, builder)
	return builder
}

// TestPort tests various port operations.
func TestPort(t *testing.T) {
	portType := fwdpb.PortType_PORT_TYPE_CPU_PORT
	unregister(portType)

	ctx := fwdcontext.New("test", "fwd")

	pd := &fwdpb.PortDesc{
		PortType: portType,
		PortId:   MakeID(fwdobject.NewID("TestPort")),
	}

	// Create a port, no builder registered.
	port, err := New(pd, ctx)
	if err != nil {
		t.Logf("Got expected error %v.", err)
	} else {
		t.Errorf("Unexpected object %v created.", port)
	}

	// Create a port, builder registered.
	builder := newTestBuilder(10, portType)
	port, err = New(pd, ctx)
	if err == nil {
		t.Logf("Created object %v.", port)
	} else {
		t.Errorf("Port creation failed, err %v.", err)
	}
	if !builder.port.allocated {
		t.Errorf("Port not allocated as expected, port %v.", builder.port)
	}
	id := string(port.ID())

	// Find the port, invalid object id.
	invalid := fwdobject.NewID(id + "1")
	if _, err := Find(MakeID(invalid), ctx); err != nil {
		t.Logf("Got expected error %v.", err)
	} else {
		t.Errorf("Unexpected object %v found.", port)
	}

	// Find the port, valid object id.
	if _, err := Find(GetID(port), ctx); err != nil {
		t.Errorf("Unable to find port, error %v.", err)
	}
}
