/*
Package lldp is a generated package which contains definitions
of structs which generate gNMI paths for a YANG schema.

This package was generated by ygnmi version: v0.11.1: (ygot: v0.29.21)
using the following YANG input files:
  - public/release/models/acl/openconfig-acl.yang
  - public/release/models/acl/openconfig-packet-match.yang
  - public/release/models/aft/openconfig-aft-network-instance.yang
  - public/release/models/aft/openconfig-aft-summary.yang
  - public/release/models/aft/openconfig-aft.yang
  - public/release/models/bfd/openconfig-bfd.yang
  - public/release/models/bgp/openconfig-bgp-policy.yang
  - public/release/models/bgp/openconfig-bgp-types.yang
  - public/release/models/extensions/openconfig-metadata.yang
  - public/release/models/gnsi/openconfig-gnsi-acctz.yang
  - public/release/models/gnsi/openconfig-gnsi-authz.yang
  - public/release/models/gnsi/openconfig-gnsi-certz.yang
  - public/release/models/gnsi/openconfig-gnsi-credentialz.yang
  - public/release/models/gnsi/openconfig-gnsi-pathz.yang
  - public/release/models/gnsi/openconfig-gnsi.yang
  - public/release/models/gribi/openconfig-gribi.yang
  - public/release/models/interfaces/openconfig-if-aggregate.yang
  - public/release/models/interfaces/openconfig-if-ethernet-ext.yang
  - public/release/models/interfaces/openconfig-if-ethernet.yang
  - public/release/models/interfaces/openconfig-if-ip-ext.yang
  - public/release/models/interfaces/openconfig-if-ip.yang
  - public/release/models/interfaces/openconfig-if-sdn-ext.yang
  - public/release/models/interfaces/openconfig-interfaces.yang
  - public/release/models/isis/openconfig-isis-policy.yang
  - public/release/models/isis/openconfig-isis.yang
  - public/release/models/lacp/openconfig-lacp.yang
  - public/release/models/lldp/openconfig-lldp-types.yang
  - public/release/models/lldp/openconfig-lldp.yang
  - public/release/models/local-routing/openconfig-local-routing.yang
  - public/release/models/mpls/openconfig-mpls-types.yang
  - public/release/models/multicast/openconfig-pim.yang
  - public/release/models/network-instance/openconfig-network-instance.yang
  - public/release/models/openconfig-extensions.yang
  - public/release/models/optical-transport/openconfig-transport-types.yang
  - public/release/models/ospf/openconfig-ospf-policy.yang
  - public/release/models/ospf/openconfig-ospfv2.yang
  - public/release/models/p4rt/openconfig-p4rt.yang
  - public/release/models/platform/openconfig-platform-common.yang
  - public/release/models/platform/openconfig-platform-controller-card.yang
  - public/release/models/platform/openconfig-platform-cpu.yang
  - public/release/models/platform/openconfig-platform-ext.yang
  - public/release/models/platform/openconfig-platform-fabric.yang
  - public/release/models/platform/openconfig-platform-fan.yang
  - public/release/models/platform/openconfig-platform-integrated-circuit.yang
  - public/release/models/platform/openconfig-platform-linecard.yang
  - public/release/models/platform/openconfig-platform-pipeline-counters.yang
  - public/release/models/platform/openconfig-platform-psu.yang
  - public/release/models/platform/openconfig-platform-software.yang
  - public/release/models/platform/openconfig-platform-transceiver.yang
  - public/release/models/platform/openconfig-platform.yang
  - public/release/models/policy-forwarding/openconfig-policy-forwarding.yang
  - public/release/models/policy/openconfig-policy-types.yang
  - public/release/models/qos/openconfig-qos-elements.yang
  - public/release/models/qos/openconfig-qos-interfaces.yang
  - public/release/models/qos/openconfig-qos-types.yang
  - public/release/models/qos/openconfig-qos.yang
  - public/release/models/relay-agent/openconfig-relay-agent.yang
  - public/release/models/rib/openconfig-rib-bgp.yang
  - public/release/models/sampling/openconfig-sampling-sflow.yang
  - public/release/models/segment-routing/openconfig-segment-routing-types.yang
  - public/release/models/system/openconfig-system-bootz.yang
  - public/release/models/system/openconfig-system-controlplane.yang
  - public/release/models/system/openconfig-system-utilization.yang
  - public/release/models/system/openconfig-system.yang
  - public/release/models/types/openconfig-inet-types.yang
  - public/release/models/types/openconfig-types.yang
  - public/release/models/types/openconfig-yang-types.yang
  - public/release/models/vlan/openconfig-vlan.yang
  - public/third_party/ietf/iana-if-type.yang
  - public/third_party/ietf/ietf-inet-types.yang
  - public/third_party/ietf/ietf-interfaces.yang
  - public/third_party/ietf/ietf-yang-types.yang
  - yang/openconfig-bgp-gue.yang

Imported modules were sourced from:
  - public/release/models/...
  - public/third_party/ietf/...
*/
package lldp

import (
	oc "github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// Lldp_Interface_Counters_FrameDiscardPath represents the /openconfig-lldp/lldp/interfaces/interface/state/counters/frame-discard YANG schema element.
type Lldp_Interface_Counters_FrameDiscardPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Lldp_Interface_Counters_FrameDiscardPathAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/state/counters/frame-discard YANG schema element.
type Lldp_Interface_Counters_FrameDiscardPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "frame-discard"
//	Path from root:       "/lldp/interfaces/interface/state/counters/frame-discard"
func (n *Lldp_Interface_Counters_FrameDiscardPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"frame-discard"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).FrameDiscard
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "frame-discard"
//	Path from root:       "/lldp/interfaces/interface/state/counters/frame-discard"
func (n *Lldp_Interface_Counters_FrameDiscardPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"frame-discard"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).FrameDiscard
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// Lldp_Interface_Counters_FrameErrorInPath represents the /openconfig-lldp/lldp/interfaces/interface/state/counters/frame-error-in YANG schema element.
type Lldp_Interface_Counters_FrameErrorInPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Lldp_Interface_Counters_FrameErrorInPathAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/state/counters/frame-error-in YANG schema element.
type Lldp_Interface_Counters_FrameErrorInPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "frame-error-in"
//	Path from root:       "/lldp/interfaces/interface/state/counters/frame-error-in"
func (n *Lldp_Interface_Counters_FrameErrorInPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"frame-error-in"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).FrameErrorIn
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "frame-error-in"
//	Path from root:       "/lldp/interfaces/interface/state/counters/frame-error-in"
func (n *Lldp_Interface_Counters_FrameErrorInPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"frame-error-in"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).FrameErrorIn
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// Lldp_Interface_Counters_FrameErrorOutPath represents the /openconfig-lldp/lldp/interfaces/interface/state/counters/frame-error-out YANG schema element.
type Lldp_Interface_Counters_FrameErrorOutPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Lldp_Interface_Counters_FrameErrorOutPathAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/state/counters/frame-error-out YANG schema element.
type Lldp_Interface_Counters_FrameErrorOutPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "frame-error-out"
//	Path from root:       "/lldp/interfaces/interface/state/counters/frame-error-out"
func (n *Lldp_Interface_Counters_FrameErrorOutPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"frame-error-out"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).FrameErrorOut
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "frame-error-out"
//	Path from root:       "/lldp/interfaces/interface/state/counters/frame-error-out"
func (n *Lldp_Interface_Counters_FrameErrorOutPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"frame-error-out"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).FrameErrorOut
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// Lldp_Interface_Counters_FrameInPath represents the /openconfig-lldp/lldp/interfaces/interface/state/counters/frame-in YANG schema element.
type Lldp_Interface_Counters_FrameInPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Lldp_Interface_Counters_FrameInPathAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/state/counters/frame-in YANG schema element.
type Lldp_Interface_Counters_FrameInPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "frame-in"
//	Path from root:       "/lldp/interfaces/interface/state/counters/frame-in"
func (n *Lldp_Interface_Counters_FrameInPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"frame-in"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).FrameIn
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "frame-in"
//	Path from root:       "/lldp/interfaces/interface/state/counters/frame-in"
func (n *Lldp_Interface_Counters_FrameInPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"frame-in"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).FrameIn
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// Lldp_Interface_Counters_FrameOutPath represents the /openconfig-lldp/lldp/interfaces/interface/state/counters/frame-out YANG schema element.
type Lldp_Interface_Counters_FrameOutPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Lldp_Interface_Counters_FrameOutPathAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/state/counters/frame-out YANG schema element.
type Lldp_Interface_Counters_FrameOutPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "frame-out"
//	Path from root:       "/lldp/interfaces/interface/state/counters/frame-out"
func (n *Lldp_Interface_Counters_FrameOutPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"frame-out"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).FrameOut
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "frame-out"
//	Path from root:       "/lldp/interfaces/interface/state/counters/frame-out"
func (n *Lldp_Interface_Counters_FrameOutPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"frame-out"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).FrameOut
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// Lldp_Interface_Counters_LastClearPath represents the /openconfig-lldp/lldp/interfaces/interface/state/counters/last-clear YANG schema element.
type Lldp_Interface_Counters_LastClearPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Lldp_Interface_Counters_LastClearPathAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/state/counters/last-clear YANG schema element.
type Lldp_Interface_Counters_LastClearPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "last-clear"
//	Path from root:       "/lldp/interfaces/interface/state/counters/last-clear"
func (n *Lldp_Interface_Counters_LastClearPath) State() ygnmi.SingletonQuery[string] {
	return ygnmi.NewSingletonQuery[string](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"last-clear"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).LastClear
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "last-clear"
//	Path from root:       "/lldp/interfaces/interface/state/counters/last-clear"
func (n *Lldp_Interface_Counters_LastClearPathAny) State() ygnmi.WildcardQuery[string] {
	return ygnmi.NewWildcardQuery[string](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"last-clear"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (string, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).LastClear
			if ret == nil {
				var zero string
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// Lldp_Interface_Counters_TlvDiscardPath represents the /openconfig-lldp/lldp/interfaces/interface/state/counters/tlv-discard YANG schema element.
type Lldp_Interface_Counters_TlvDiscardPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Lldp_Interface_Counters_TlvDiscardPathAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/state/counters/tlv-discard YANG schema element.
type Lldp_Interface_Counters_TlvDiscardPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "tlv-discard"
//	Path from root:       "/lldp/interfaces/interface/state/counters/tlv-discard"
func (n *Lldp_Interface_Counters_TlvDiscardPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"tlv-discard"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).TlvDiscard
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "tlv-discard"
//	Path from root:       "/lldp/interfaces/interface/state/counters/tlv-discard"
func (n *Lldp_Interface_Counters_TlvDiscardPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"Lldp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"tlv-discard"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.Lldp_Interface_Counters).TlvDiscard
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Counters) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		nil,
	)
}
