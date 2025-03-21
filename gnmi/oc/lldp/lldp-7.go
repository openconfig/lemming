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
	"reflect"

	oc "github.com/openconfig/lemming/gnmi/oc"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/ygot/ytypes"
)

// Lldp_Interface_Neighbor_Tlv_TypePath represents the /openconfig-lldp/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/state/type YANG schema element.
type Lldp_Interface_Neighbor_Tlv_TypePath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Lldp_Interface_Neighbor_Tlv_TypePathAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/state/type YANG schema element.
type Lldp_Interface_Neighbor_Tlv_TypePathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "state/type"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/state/type"
func (n *Lldp_Interface_Neighbor_Tlv_TypePath) State() ygnmi.SingletonQuery[int32] {
	return ygnmi.NewSingletonQuery[int32](
		"Lldp_Interface_Neighbor_Tlv",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"state", "type"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int32, bool) {
			ret := gs.(*oc.Lldp_Interface_Neighbor_Tlv).Type
			if ret == nil {
				var zero int32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Neighbor_Tlv) },
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
//	Path from parent:     "state/type"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/state/type"
func (n *Lldp_Interface_Neighbor_Tlv_TypePathAny) State() ygnmi.WildcardQuery[int32] {
	return ygnmi.NewWildcardQuery[int32](
		"Lldp_Interface_Neighbor_Tlv",
		true,
		false,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"state", "type"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int32, bool) {
			ret := gs.(*oc.Lldp_Interface_Neighbor_Tlv).Type
			if ret == nil {
				var zero int32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Neighbor_Tlv) },
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
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "type"
//	Path from root:       ""
func (n *Lldp_Interface_Neighbor_Tlv_TypePath) Config() ygnmi.ConfigQuery[int32] {
	return ygnmi.NewConfigQuery[int32](
		"Lldp_Interface_Neighbor_Tlv",
		false,
		true,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"type"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int32, bool) {
			ret := gs.(*oc.Lldp_Interface_Neighbor_Tlv).Type
			if ret == nil {
				var zero int32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Neighbor_Tlv) },
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
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "type"
//	Path from root:       ""
func (n *Lldp_Interface_Neighbor_Tlv_TypePathAny) Config() ygnmi.WildcardQuery[int32] {
	return ygnmi.NewWildcardQuery[int32](
		"Lldp_Interface_Neighbor_Tlv",
		false,
		true,
		true,
		true,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"type"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (int32, bool) {
			ret := gs.(*oc.Lldp_Interface_Neighbor_Tlv).Type
			if ret == nil {
				var zero int32
				return zero, false
			}
			return *ret, true
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Neighbor_Tlv) },
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

// Lldp_Interface_Neighbor_Tlv_ValuePath represents the /openconfig-lldp/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/state/value YANG schema element.
type Lldp_Interface_Neighbor_Tlv_ValuePath struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// Lldp_Interface_Neighbor_Tlv_ValuePathAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/state/value YANG schema element.
type Lldp_Interface_Neighbor_Tlv_ValuePathAny struct {
	*ygnmi.NodePath
	parent ygnmi.PathStruct
}

// State returns a Query that can be used in gNMI operations.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "state/value"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/state/value"
func (n *Lldp_Interface_Neighbor_Tlv_ValuePath) State() ygnmi.SingletonQuery[oc.Binary] {
	return ygnmi.NewSingletonQuery[oc.Binary](
		"Lldp_Interface_Neighbor_Tlv",
		true,
		false,
		true,
		false,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"state", "value"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (oc.Binary, bool) {
			ret := gs.(*oc.Lldp_Interface_Neighbor_Tlv).Value
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Neighbor_Tlv) },
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
//	Path from parent:     "state/value"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/state/value"
func (n *Lldp_Interface_Neighbor_Tlv_ValuePathAny) State() ygnmi.WildcardQuery[oc.Binary] {
	return ygnmi.NewWildcardQuery[oc.Binary](
		"Lldp_Interface_Neighbor_Tlv",
		true,
		false,
		true,
		false,
		true,
		false,
		ygnmi.NewNodePath(
			[]string{"state", "value"},
			nil,
			n.parent,
		),
		func(gs ygot.ValidatedGoStruct) (oc.Binary, bool) {
			ret := gs.(*oc.Lldp_Interface_Neighbor_Tlv).Value
			return ret, !reflect.ValueOf(ret).IsZero()
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Neighbor_Tlv) },
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

// Lldp_Interface_Neighbor_TlvPath represents the /openconfig-lldp/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv YANG schema element.
type Lldp_Interface_Neighbor_TlvPath struct {
	*ygnmi.NodePath
}

// Lldp_Interface_Neighbor_TlvPathAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv YANG schema element.
type Lldp_Interface_Neighbor_TlvPathAny struct {
	*ygnmi.NodePath
}

// Lldp_Interface_Neighbor_TlvPathMap represents the /openconfig-lldp/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv YANG schema element.
type Lldp_Interface_Neighbor_TlvPathMap struct {
	*ygnmi.NodePath
}

// Lldp_Interface_Neighbor_TlvPathMapAny represents the wildcard version of the /openconfig-lldp/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv YANG schema element.
type Lldp_Interface_Neighbor_TlvPathMapAny struct {
	*ygnmi.NodePath
}

// Oui (leaf): The organizationally unique identifier field shall contain
// the organization's OUI as defined in Clause 9 of IEEE Std
// 802. The high-order octet is 0 and the low-order 3 octets
// are the SMI Network Management Private Enterprise Code of
// the Vendor in network byte order, as defined in the
// 'Assigned Numbers' RFC [RFC3232].
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "*/oui"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/*/oui"
func (n *Lldp_Interface_Neighbor_TlvPath) Oui() *Lldp_Interface_Neighbor_Tlv_OuiPath {
	ps := &Lldp_Interface_Neighbor_Tlv_OuiPath{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "oui"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Oui (leaf): The organizationally unique identifier field shall contain
// the organization's OUI as defined in Clause 9 of IEEE Std
// 802. The high-order octet is 0 and the low-order 3 octets
// are the SMI Network Management Private Enterprise Code of
// the Vendor in network byte order, as defined in the
// 'Assigned Numbers' RFC [RFC3232].
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "*/oui"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/*/oui"
func (n *Lldp_Interface_Neighbor_TlvPathAny) Oui() *Lldp_Interface_Neighbor_Tlv_OuiPathAny {
	ps := &Lldp_Interface_Neighbor_Tlv_OuiPathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "oui"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// OuiSubtype (leaf): The organizationally defined subtype field shall contain a
// unique subtype value assigned by the defining organization.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "*/oui-subtype"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/*/oui-subtype"
func (n *Lldp_Interface_Neighbor_TlvPath) OuiSubtype() *Lldp_Interface_Neighbor_Tlv_OuiSubtypePath {
	ps := &Lldp_Interface_Neighbor_Tlv_OuiSubtypePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "oui-subtype"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// OuiSubtype (leaf): The organizationally defined subtype field shall contain a
// unique subtype value assigned by the defining organization.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "*/oui-subtype"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/*/oui-subtype"
func (n *Lldp_Interface_Neighbor_TlvPathAny) OuiSubtype() *Lldp_Interface_Neighbor_Tlv_OuiSubtypePathAny {
	ps := &Lldp_Interface_Neighbor_Tlv_OuiSubtypePathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "oui-subtype"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Type (leaf): The integer value identifying the type of information
// contained in the value field.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "*/type"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/*/type"
func (n *Lldp_Interface_Neighbor_TlvPath) Type() *Lldp_Interface_Neighbor_Tlv_TypePath {
	ps := &Lldp_Interface_Neighbor_Tlv_TypePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "type"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Type (leaf): The integer value identifying the type of information
// contained in the value field.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "*/type"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/*/type"
func (n *Lldp_Interface_Neighbor_TlvPathAny) Type() *Lldp_Interface_Neighbor_Tlv_TypePathAny {
	ps := &Lldp_Interface_Neighbor_Tlv_TypePathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"*", "type"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Value (leaf): A variable-length octet-string containing the
// instance-specific information for this TLV.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "state/value"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/state/value"
func (n *Lldp_Interface_Neighbor_TlvPath) Value() *Lldp_Interface_Neighbor_Tlv_ValuePath {
	ps := &Lldp_Interface_Neighbor_Tlv_ValuePath{
		NodePath: ygnmi.NewNodePath(
			[]string{"state", "value"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// Value (leaf): A variable-length octet-string containing the
// instance-specific information for this TLV.
//
//	Defining module:      "openconfig-lldp"
//	Instantiating module: "openconfig-lldp"
//	Path from parent:     "state/value"
//	Path from root:       "/lldp/interfaces/interface/neighbors/neighbor/custom-tlvs/tlv/state/value"
func (n *Lldp_Interface_Neighbor_TlvPathAny) Value() *Lldp_Interface_Neighbor_Tlv_ValuePathAny {
	ps := &Lldp_Interface_Neighbor_Tlv_ValuePathAny{
		NodePath: ygnmi.NewNodePath(
			[]string{"state", "value"},
			map[string]interface{}{},
			n,
		),
		parent: n,
	}
	return ps
}

// State returns a Query that can be used in gNMI operations.
func (n *Lldp_Interface_Neighbor_TlvPath) State() ygnmi.SingletonQuery[*oc.Lldp_Interface_Neighbor_Tlv] {
	return ygnmi.NewSingletonQuery[*oc.Lldp_Interface_Neighbor_Tlv](
		"Lldp_Interface_Neighbor_Tlv",
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
func (n *Lldp_Interface_Neighbor_TlvPathAny) State() ygnmi.WildcardQuery[*oc.Lldp_Interface_Neighbor_Tlv] {
	return ygnmi.NewWildcardQuery[*oc.Lldp_Interface_Neighbor_Tlv](
		"Lldp_Interface_Neighbor_Tlv",
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
func (n *Lldp_Interface_Neighbor_TlvPathMap) State() ygnmi.SingletonQuery[map[oc.Lldp_Interface_Neighbor_Tlv_Key]*oc.Lldp_Interface_Neighbor_Tlv] {
	return ygnmi.NewSingletonQuery[map[oc.Lldp_Interface_Neighbor_Tlv_Key]*oc.Lldp_Interface_Neighbor_Tlv](
		"Lldp_Interface_Neighbor",
		true,
		false,
		false,
		false,
		true,
		true,
		n,
		func(gs ygot.ValidatedGoStruct) (map[oc.Lldp_Interface_Neighbor_Tlv_Key]*oc.Lldp_Interface_Neighbor_Tlv, bool) {
			ret := gs.(*oc.Lldp_Interface_Neighbor).Tlv
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Neighbor) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		&ygnmi.CompressionInfo{
			PreRelPath:  []string{"openconfig-lldp:custom-tlvs"},
			PostRelPath: []string{"openconfig-lldp:tlv"},
		},
	)
}

// State returns a Query that can be used in gNMI operations.
func (n *Lldp_Interface_Neighbor_TlvPathMapAny) State() ygnmi.WildcardQuery[map[oc.Lldp_Interface_Neighbor_Tlv_Key]*oc.Lldp_Interface_Neighbor_Tlv] {
	return ygnmi.NewWildcardQuery[map[oc.Lldp_Interface_Neighbor_Tlv_Key]*oc.Lldp_Interface_Neighbor_Tlv](
		"Lldp_Interface_Neighbor",
		true,
		false,
		false,
		false,
		true,
		true,
		n,
		func(gs ygot.ValidatedGoStruct) (map[oc.Lldp_Interface_Neighbor_Tlv_Key]*oc.Lldp_Interface_Neighbor_Tlv, bool) {
			ret := gs.(*oc.Lldp_Interface_Neighbor).Tlv
			return ret, ret != nil
		},
		func() ygot.ValidatedGoStruct { return new(oc.Lldp_Interface_Neighbor) },
		func() *ytypes.Schema {
			return &ytypes.Schema{
				Root:       &oc.Root{},
				SchemaTree: oc.SchemaTree,
				Unmarshal:  oc.Unmarshal,
			}
		},
		nil,
		&ygnmi.CompressionInfo{
			PreRelPath:  []string{"openconfig-lldp:custom-tlvs"},
			PostRelPath: []string{"openconfig-lldp:tlv"},
		},
	)
}
