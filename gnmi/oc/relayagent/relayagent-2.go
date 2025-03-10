/*
Package relayagent is a generated package which contains definitions
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
package relayagent

import (
	oc "github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// RelayAgent_Dhcp_Interface_Counters_DhcpAckSentPath represents the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-ack-sent YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpAckSentPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// RelayAgent_Dhcp_Interface_Counters_DhcpAckSentPathAny represents the wildcard version of the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-ack-sent YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpAckSentPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-ack-sent"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-ack-sent"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpAckSentPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-ack-sent"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpAckSent
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-ack-sent"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-ack-sent"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpAckSentPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-ack-sent"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpAckSent
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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

// RelayAgent_Dhcp_Interface_Counters_DhcpDeclineReceivedPath represents the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-decline-received YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpDeclineReceivedPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// RelayAgent_Dhcp_Interface_Counters_DhcpDeclineReceivedPathAny represents the wildcard version of the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-decline-received YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpDeclineReceivedPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-decline-received"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-decline-received"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpDeclineReceivedPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-decline-received"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpDeclineReceived
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-decline-received"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-decline-received"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpDeclineReceivedPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-decline-received"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpDeclineReceived
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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

// RelayAgent_Dhcp_Interface_Counters_DhcpDiscoverReceivedPath represents the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-discover-received YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpDiscoverReceivedPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// RelayAgent_Dhcp_Interface_Counters_DhcpDiscoverReceivedPathAny represents the wildcard version of the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-discover-received YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpDiscoverReceivedPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-discover-received"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-discover-received"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpDiscoverReceivedPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-discover-received"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpDiscoverReceived
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-discover-received"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-discover-received"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpDiscoverReceivedPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-discover-received"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpDiscoverReceived
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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

// RelayAgent_Dhcp_Interface_Counters_DhcpInformReceivedPath represents the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-inform-received YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpInformReceivedPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// RelayAgent_Dhcp_Interface_Counters_DhcpInformReceivedPathAny represents the wildcard version of the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-inform-received YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpInformReceivedPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-inform-received"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-inform-received"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpInformReceivedPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-inform-received"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpInformReceived
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-inform-received"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-inform-received"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpInformReceivedPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-inform-received"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpInformReceived
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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

// RelayAgent_Dhcp_Interface_Counters_DhcpNackSentPath represents the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-nack-sent YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpNackSentPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// RelayAgent_Dhcp_Interface_Counters_DhcpNackSentPathAny represents the wildcard version of the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-nack-sent YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpNackSentPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-nack-sent"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-nack-sent"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpNackSentPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-nack-sent"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpNackSent
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-nack-sent"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-nack-sent"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpNackSentPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-nack-sent"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpNackSent
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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

// RelayAgent_Dhcp_Interface_Counters_DhcpOfferSentPath represents the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-offer-sent YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpOfferSentPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// RelayAgent_Dhcp_Interface_Counters_DhcpOfferSentPathAny represents the wildcard version of the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-offer-sent YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpOfferSentPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-offer-sent"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-offer-sent"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpOfferSentPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-offer-sent"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpOfferSent
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-offer-sent"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-offer-sent"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpOfferSentPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-offer-sent"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpOfferSent
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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

// RelayAgent_Dhcp_Interface_Counters_DhcpReleaseReceivedPath represents the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-release-received YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpReleaseReceivedPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// RelayAgent_Dhcp_Interface_Counters_DhcpReleaseReceivedPathAny represents the wildcard version of the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-release-received YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpReleaseReceivedPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-release-received"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-release-received"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpReleaseReceivedPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-release-received"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpReleaseReceived
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-release-received"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-release-received"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpReleaseReceivedPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-release-received"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpReleaseReceived
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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

// RelayAgent_Dhcp_Interface_Counters_DhcpRequestReceivedPath represents the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-request-received YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpRequestReceivedPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// RelayAgent_Dhcp_Interface_Counters_DhcpRequestReceivedPathAny represents the wildcard version of the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-request-received YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_DhcpRequestReceivedPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-request-received"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-request-received"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpRequestReceivedPath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-request-received"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpRequestReceived
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "dhcp-request-received"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/dhcp-request-received"
func (n *RelayAgent_Dhcp_Interface_Counters_DhcpRequestReceivedPathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"dhcp-request-received"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).DhcpRequestReceived
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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

// RelayAgent_Dhcp_Interface_Counters_InvalidOpcodePath represents the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/invalid-opcode YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_InvalidOpcodePath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// RelayAgent_Dhcp_Interface_Counters_InvalidOpcodePathAny represents the wildcard version of the /openconfig-relay-agent/relay-agent/dhcp/interfaces/interface/state/counters/invalid-opcode YANG schema element.
type RelayAgent_Dhcp_Interface_Counters_InvalidOpcodePathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "invalid-opcode"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/invalid-opcode"
func (n *RelayAgent_Dhcp_Interface_Counters_InvalidOpcodePath) State() ygnmi.SingletonQuery[uint64] {
	return ygnmi.NewSingletonQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"invalid-opcode"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).InvalidOpcode
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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
//	Defining module:      "openconfig-relay-agent"
//	Instantiating module: "openconfig-relay-agent"
//	Path from parent:     "invalid-opcode"
//	Path from root:       "/relay-agent/dhcp/interfaces/interface/state/counters/invalid-opcode"
func (n *RelayAgent_Dhcp_Interface_Counters_InvalidOpcodePathAny) State() ygnmi.WildcardQuery[uint64] {
	return ygnmi.NewWildcardQuery[uint64](
		"RelayAgent_Dhcp_Interface_Counters",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"invalid-opcode"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint64, bool) {
			ret := gs.(*oc.RelayAgent_Dhcp_Interface_Counters).InvalidOpcode
			if ret == nil {
				var zero uint64
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.RelayAgent_Dhcp_Interface_Counters) },
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
