/*
Package sampling is a generated package which contains definitions
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
package sampling

import (
	oc "github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// Sampling_Sflow_Interface_PollingIntervalPath represents the /openconfig-sampling/sampling/sflow/interfaces/interface/state/polling-interval YANG schema element.
type Sampling_Sflow_Interface_PollingIntervalPath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Sampling_Sflow_Interface_PollingIntervalPathAny represents the wildcard version of the /openconfig-sampling/sampling/sflow/interfaces/interface/state/polling-interval YANG schema element.
type Sampling_Sflow_Interface_PollingIntervalPathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling-sflow"
//	Path from parent:     "state/polling-interval"
//	Path from root:       "/sampling/sflow/interfaces/interface/state/polling-interval"
func (n *Sampling_Sflow_Interface_PollingIntervalPath) State() ygnmi.SingletonQuery[uint16] {
	return ygnmi.NewSingletonQuery[uint16](
		"Sampling_Sflow_Interface",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"state", "polling-interval"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint16, bool) {
			ret := gs.(*oc.Sampling_Sflow_Interface).PollingInterval
			if ret == nil {
				var zero uint16
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Sampling_Sflow_Interface) },
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
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling-sflow"
//	Path from parent:     "state/polling-interval"
//	Path from root:       "/sampling/sflow/interfaces/interface/state/polling-interval"
func (n *Sampling_Sflow_Interface_PollingIntervalPathAny) State() ygnmi.WildcardQuery[uint16] {
	return ygnmi.NewWildcardQuery[uint16](
		"Sampling_Sflow_Interface",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"state", "polling-interval"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint16, bool) {
			ret := gs.(*oc.Sampling_Sflow_Interface).PollingInterval
			if ret == nil {
				var zero uint16
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Sampling_Sflow_Interface) },
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

// Config returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling-sflow"
//	Path from parent:     "config/polling-interval"
//	Path from root:       "/sampling/sflow/interfaces/interface/config/polling-interval"
func (n *Sampling_Sflow_Interface_PollingIntervalPath) Config() ygnmi.ConfigQuery[uint16] {
	return ygnmi.NewConfigQuery[uint16](
		"Sampling_Sflow_Interface",
		false,
		true,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"config", "polling-interval"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint16, bool) {
			ret := gs.(*oc.Sampling_Sflow_Interface).PollingInterval
			if ret == nil {
				var zero uint16
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Sampling_Sflow_Interface) },
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

// Config returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling-sflow"
//	Path from parent:     "config/polling-interval"
//	Path from root:       "/sampling/sflow/interfaces/interface/config/polling-interval"
func (n *Sampling_Sflow_Interface_PollingIntervalPathAny) Config() ygnmi.WildcardQuery[uint16] {
	return ygnmi.NewWildcardQuery[uint16](
		"Sampling_Sflow_Interface",
		false,
		true,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"config", "polling-interval"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (uint16, bool) {
			ret := gs.(*oc.Sampling_Sflow_Interface).PollingInterval
			if ret == nil {
				var zero uint16
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Sampling_Sflow_Interface) },
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

// Sampling_Sflow_InterfacePath represents the /openconfig-sampling/sampling/sflow/interfaces/interface YANG schema element.
type Sampling_Sflow_InterfacePath struct {
	*ygnmi.NodePath
}

// Sampling_Sflow_InterfacePathAny represents the wildcard version of the /openconfig-sampling/sampling/sflow/interfaces/interface YANG schema element.
type Sampling_Sflow_InterfacePathAny struct {
	*ygnmi.NodePath
}

// Sampling_Sflow_InterfacePathMap represents the /openconfig-sampling/sampling/sflow/interfaces/interface YANG schema element.
type Sampling_Sflow_InterfacePathMap struct {
	*ygnmi.NodePath
}

// Sampling_Sflow_InterfacePathMapAny represents the wildcard version of the /openconfig-sampling/sampling/sflow/interfaces/interface YANG schema element.
type Sampling_Sflow_InterfacePathMapAny struct {
	*ygnmi.NodePath
}

// EgressSamplingRate (leaf): Sets the egress packet sampling rate.  The rate is expressed
// as an integer N, where the intended sampling rate is 1/N
// packets.  An implementation may implement the sampling rate as
// a statistical average, rather than a strict periodic sampling.
//
// The allowable sampling rate range is generally a property of
// the system, e.g., determined by the capability of the
// hardware.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "*/egress-sampling-rate"
//	Path from root:       "/sampling/sflow/interfaces/interface/*/egress-sampling-rate"
func (n *Sampling_Sflow_InterfacePath) EgressSamplingRate() *Sampling_Sflow_Interface_EgressSamplingRatePath {
	ps := &Sampling_Sflow_Interface_EgressSamplingRatePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "egress-sampling-rate"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// EgressSamplingRate (leaf): Sets the egress packet sampling rate.  The rate is expressed
// as an integer N, where the intended sampling rate is 1/N
// packets.  An implementation may implement the sampling rate as
// a statistical average, rather than a strict periodic sampling.
//
// The allowable sampling rate range is generally a property of
// the system, e.g., determined by the capability of the
// hardware.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "*/egress-sampling-rate"
//	Path from root:       "/sampling/sflow/interfaces/interface/*/egress-sampling-rate"
func (n *Sampling_Sflow_InterfacePathAny) EgressSamplingRate() *Sampling_Sflow_Interface_EgressSamplingRatePathAny {
	ps := &Sampling_Sflow_Interface_EgressSamplingRatePathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "egress-sampling-rate"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Enabled (leaf): Enables or disables sFlow on the interface.  If sFlow is
// globally disabled, this leaf is ignored.  If sFlow
// is globally enabled, this leaf may be used to disable it
// for a specific interface.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "*/enabled"
//	Path from root:       "/sampling/sflow/interfaces/interface/*/enabled"
func (n *Sampling_Sflow_InterfacePath) Enabled() *Sampling_Sflow_Interface_EnabledPath {
	ps := &Sampling_Sflow_Interface_EnabledPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "enabled"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Enabled (leaf): Enables or disables sFlow on the interface.  If sFlow is
// globally disabled, this leaf is ignored.  If sFlow
// is globally enabled, this leaf may be used to disable it
// for a specific interface.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "*/enabled"
//	Path from root:       "/sampling/sflow/interfaces/interface/*/enabled"
func (n *Sampling_Sflow_InterfacePathAny) Enabled() *Sampling_Sflow_Interface_EnabledPathAny {
	ps := &Sampling_Sflow_Interface_EnabledPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "enabled"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// IngressSamplingRate (leaf): Sets the ingress packet sampling rate.  The rate is expressed
// as an integer N, where the intended sampling rate is 1/N
// packets.  An implementation may implement the sampling rate as
// a statistical average, rather than a strict periodic sampling.
//
// The allowable sampling rate range is generally a property of
// the system, e.g., determined by the capability of the
// hardware.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "*/ingress-sampling-rate"
//	Path from root:       "/sampling/sflow/interfaces/interface/*/ingress-sampling-rate"
func (n *Sampling_Sflow_InterfacePath) IngressSamplingRate() *Sampling_Sflow_Interface_IngressSamplingRatePath {
	ps := &Sampling_Sflow_Interface_IngressSamplingRatePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "ingress-sampling-rate"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// IngressSamplingRate (leaf): Sets the ingress packet sampling rate.  The rate is expressed
// as an integer N, where the intended sampling rate is 1/N
// packets.  An implementation may implement the sampling rate as
// a statistical average, rather than a strict periodic sampling.
//
// The allowable sampling rate range is generally a property of
// the system, e.g., determined by the capability of the
// hardware.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "*/ingress-sampling-rate"
//	Path from root:       "/sampling/sflow/interfaces/interface/*/ingress-sampling-rate"
func (n *Sampling_Sflow_InterfacePathAny) IngressSamplingRate() *Sampling_Sflow_Interface_IngressSamplingRatePathAny {
	ps := &Sampling_Sflow_Interface_IngressSamplingRatePathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "ingress-sampling-rate"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Name (leaf): Reference to the interface for sFlow configuration and
// state.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "*/name"
//	Path from root:       "/sampling/sflow/interfaces/interface/*/name"
func (n *Sampling_Sflow_InterfacePath) Name() *Sampling_Sflow_Interface_NamePath {
	ps := &Sampling_Sflow_Interface_NamePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "name"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Name (leaf): Reference to the interface for sFlow configuration and
// state.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "*/name"
//	Path from root:       "/sampling/sflow/interfaces/interface/*/name"
func (n *Sampling_Sflow_InterfacePathAny) Name() *Sampling_Sflow_Interface_NamePathAny {
	ps := &Sampling_Sflow_Interface_NamePathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "name"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// PacketsSampled (leaf): Total number of packets sampled from the interface.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "state/packets-sampled"
//	Path from root:       "/sampling/sflow/interfaces/interface/state/packets-sampled"
func (n *Sampling_Sflow_InterfacePath) PacketsSampled() *Sampling_Sflow_Interface_PacketsSampledPath {
	ps := &Sampling_Sflow_Interface_PacketsSampledPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"state", "packets-sampled"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// PacketsSampled (leaf): Total number of packets sampled from the interface.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "state/packets-sampled"
//	Path from root:       "/sampling/sflow/interfaces/interface/state/packets-sampled"
func (n *Sampling_Sflow_InterfacePathAny) PacketsSampled() *Sampling_Sflow_Interface_PacketsSampledPathAny {
	ps := &Sampling_Sflow_Interface_PacketsSampledPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"state", "packets-sampled"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// PollingInterval (leaf): Sets the traffic sampling polling interval.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "*/polling-interval"
//	Path from root:       "/sampling/sflow/interfaces/interface/*/polling-interval"
func (n *Sampling_Sflow_InterfacePath) PollingInterval() *Sampling_Sflow_Interface_PollingIntervalPath {
	ps := &Sampling_Sflow_Interface_PollingIntervalPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "polling-interval"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// PollingInterval (leaf): Sets the traffic sampling polling interval.
//
//	Defining module:      "openconfig-sampling-sflow"
//	Instantiating module: "openconfig-sampling"
//	Path from parent:     "*/polling-interval"
//	Path from root:       "/sampling/sflow/interfaces/interface/*/polling-interval"
func (n *Sampling_Sflow_InterfacePathAny) PollingInterval() *Sampling_Sflow_Interface_PollingIntervalPathAny {
	ps := &Sampling_Sflow_Interface_PollingIntervalPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "polling-interval"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// State returns a Query that can be used in gNMI operations.
func (n *Sampling_Sflow_InterfacePath) State() ygnmi.SingletonQuery[*oc.Sampling_Sflow_Interface] {
	return ygnmi.NewSingletonQuery[*oc.Sampling_Sflow_Interface](
		"Sampling_Sflow_Interface",
		true,
		false,
		false,
		false,
		true,
		false,
		n,
		nil,
		nil,
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
func (n *Sampling_Sflow_InterfacePathAny) State() ygnmi.WildcardQuery[*oc.Sampling_Sflow_Interface] {
	return ygnmi.NewWildcardQuery[*oc.Sampling_Sflow_Interface](
		"Sampling_Sflow_Interface",
		true,
		false,
		false,
		false,
		true,
		false,
		n,
		nil,
		nil,
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

// Config returns a Query that can be used in gNMI operations.
func (n *Sampling_Sflow_InterfacePath) Config() ygnmi.ConfigQuery[*oc.Sampling_Sflow_Interface] {
	return ygnmi.NewConfigQuery[*oc.Sampling_Sflow_Interface](
		"Sampling_Sflow_Interface",
		false,
		true,
		false,
		false,
		true,
		false,
		n,
		nil,
		nil,
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

// Config returns a Query that can be used in gNMI operations.
func (n *Sampling_Sflow_InterfacePathAny) Config() ygnmi.WildcardQuery[*oc.Sampling_Sflow_Interface] {
	return ygnmi.NewWildcardQuery[*oc.Sampling_Sflow_Interface](
		"Sampling_Sflow_Interface",
		false,
		true,
		false,
		false,
		true,
		false,
		n,
		nil,
		nil,
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
func (n *Sampling_Sflow_InterfacePathMap) State() ygnmi.SingletonQuery[map[string]*oc.Sampling_Sflow_Interface] {
	return ygnmi.NewSingletonQuery[map[string]*oc.Sampling_Sflow_Interface](
		"Sampling_Sflow",
		true,
		false,
		false,
		false,
		true,
		true,
		n,
		func(gs ygot.ValidatedGoStruct) (map[string]*oc.Sampling_Sflow_Interface, bool) {
			ret := gs.(*oc.Sampling_Sflow).Interface
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.Sampling_Sflow) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		&ygnmi.CompressionInfo{
			PreRelPath:  []string{"openconfig-sampling-sflow:interfaces"},
			PostRelPath: []string{"openconfig-sampling-sflow:interface"},
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *Sampling_Sflow_InterfacePathMapAny) State() ygnmi.WildcardQuery[map[string]*oc.Sampling_Sflow_Interface] {
	return ygnmi.NewWildcardQuery[map[string]*oc.Sampling_Sflow_Interface](
		"Sampling_Sflow",
		true,
		false,
		false,
		false,
		true,
		true,
		n,
		func(gs ygot.ValidatedGoStruct) (map[string]*oc.Sampling_Sflow_Interface, bool) {
			ret := gs.(*oc.Sampling_Sflow).Interface
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.Sampling_Sflow) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		&ygnmi.CompressionInfo{
			PreRelPath:  []string{"openconfig-sampling-sflow:interfaces"},
			PostRelPath: []string{"openconfig-sampling-sflow:interface"},
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Sampling_Sflow_InterfacePathMap) Config() ygnmi.ConfigQuery[map[string]*oc.Sampling_Sflow_Interface] {
	return ygnmi.NewConfigQuery[map[string]*oc.Sampling_Sflow_Interface](
		"Sampling_Sflow",
		false,
		true,
		false,
		false,
		true,
		true,
		n,
		func(gs ygot.ValidatedGoStruct) (map[string]*oc.Sampling_Sflow_Interface, bool) {
			ret := gs.(*oc.Sampling_Sflow).Interface
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.Sampling_Sflow) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		&ygnmi.CompressionInfo{
			PreRelPath:  []string{"openconfig-sampling-sflow:interfaces"},
			PostRelPath: []string{"openconfig-sampling-sflow:interface"},
		},
	)
}

// Config returns a Query that can be used in gNMI operations.
func (n *Sampling_Sflow_InterfacePathMapAny) Config() ygnmi.WildcardQuery[map[string]*oc.Sampling_Sflow_Interface] {
	return ygnmi.NewWildcardQuery[map[string]*oc.Sampling_Sflow_Interface](
		"Sampling_Sflow",
		false,
		true,
		false,
		false,
		true,
		true,
		n,
		func(gs ygot.ValidatedGoStruct) (map[string]*oc.Sampling_Sflow_Interface, bool) {
			ret := gs.(*oc.Sampling_Sflow).Interface
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.Sampling_Sflow) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		&ygnmi.CompressionInfo{
			PreRelPath:  []string{"openconfig-sampling-sflow:interfaces"},
			PostRelPath: []string{"openconfig-sampling-sflow:interface"},
		},
	)
}
