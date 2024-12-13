// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package saiutil

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/openconfig/ondatra"

	"github.com/openconfig/lemming/internal/binding"
)

// ConfigOp is interface for applying a dataplane configuration option.
type ConfigOp interface {
	Config(t testing.TB, s *Suite, dut *ondatra.DUTDevice)
	UnConfig(t testing.TB, s *Suite, dut *ondatra.DUTDevice)
}

// NewSuite returns a empty test suite.
func NewSuite() *Suite {
	return &Suite{
		interfaceMap:     map[InterfaceRef]uint64{},
		oc2SAINextHop:    map[uint64]uint64{},
		oc2SAINextHopGrp: map[uint64]uint64{},
		oc2SAIVRF:        map[string]uint64{},
	}
}

// Suite contains a baseconfig and sets of tests dataplane packet tests.
type Suite struct {
	BaseConfig []ConfigOp
	Case       []*Case

	interfaceMap     map[InterfaceRef]uint64
	oc2SAINextHop    map[uint64]uint64
	oc2SAINextHopGrp map[uint64]uint64
	oc2SAIVRF        map[string]uint64
}

// Packet is a list of layers and a corresponding port.
type Packet struct {
	Layers []gopacket.SerializableLayer
	Port   string
}

// A Case is single test case containing the config and the input and expect output.
type Case struct {
	In, Out *Packet
	Config  []ConfigOp
}

// Run runs all the cases and reports any diffs.
func (cases *Suite) Run(t testing.TB, pm *binding.PortMgr) {
	dut := ondatra.DUT(t, "dut")
	for _, cfg := range cases.BaseConfig {
		cfg.Config(t, cases, dut)
	}
	for i, cs := range cases.Case {
		t.Log("Running case ", i)
		for _, cfg := range cs.Config {
			cfg.Config(t, cases, dut)
		}
		buf := gopacket.NewSerializeBuffer()
		if err := gopacket.SerializeLayers(buf, gopacket.SerializeOptions{FixLengths: true}, cs.In.Layers...); err != nil {
			t.Fatalf("failed to serialize headers: %v", err)
		}
		p1 := pm.GetPort(dut.Port(t, cs.In.Port))
		p1.RXQueue.Write(buf.Bytes())
		p2 := pm.GetPort(dut.Port(t, cs.Out.Port))

		select {
		case packet := (<-p2.TXQueue.Receive()):
			p := gopacket.NewPacket(packet.([]byte), layers.LayerTypeEthernet, gopacket.Default)
			t.Logf("Got packet:\n%s", p.Dump())

			got := []gopacket.SerializableLayer{}
			for _, l := range p.Layers() {
				got = append(got, l.(gopacket.SerializableLayer))
			}
			// Skip the payload when comparing layers.
			if d := cmp.Diff(got[0:len(got)-1], cs.Out.Layers[0:len(cs.Out.Layers)-1], cmpopts.IgnoreUnexported(layers.IPv6{}, layers.IPv4{}, layers.UDP{}, layers.MPLS{}),
				cmpopts.IgnoreFields(layers.IPv4{}, "BaseLayer"), cmpopts.IgnoreFields(layers.UDP{}, "BaseLayer"), cmpopts.IgnoreFields(layers.Ethernet{}, "BaseLayer"), cmpopts.IgnoreFields(layers.IPv6{}, "BaseLayer")); d != "" {
				t.Error(d)
			}
		case <-time.After(10 * time.Millisecond):
		}

		for i := len(cs.Config) - 1; i >= 0; i-- {
			cs.Config[i].UnConfig(t, cases, dut)
		}
	}
}
