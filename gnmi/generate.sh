#!/bin/bash
#
# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script is used to generate the Ondatra Telemetry and Config Go APIs.

set -e

git clone https://github.com/openconfig/public.git
cd public && git checkout v5.0.0 && cd ..

EXCLUDE_MODULES=ietf-interfaces,openconfig-bfd,openconfig-messages

YANG_FILES=(
  public/release/models/acl/openconfig-acl.yang
  public/release/models/acl/openconfig-packet-match.yang
  public/release/models/aft/openconfig-aft-network-instance.yang
  public/release/models/aft/openconfig-aft-summary.yang
  public/release/models/aft/openconfig-aft.yang
  public/release/models/bfd/openconfig-bfd.yang
  public/release/models/bgp/openconfig-bgp-policy.yang
  public/release/models/bgp/openconfig-bgp-types.yang
  public/release/models/extensions/openconfig-metadata.yang
  public/release/models/gnsi/openconfig-gnsi-acctz.yang
  public/release/models/gnsi/openconfig-gnsi-authz.yang
  public/release/models/gnsi/openconfig-gnsi-certz.yang
  public/release/models/gnsi/openconfig-gnsi-credentialz.yang
  public/release/models/gnsi/openconfig-gnsi-pathz.yang
  public/release/models/gnsi/openconfig-gnsi.yang
  public/release/models/gribi/openconfig-gribi.yang
  public/release/models/interfaces/openconfig-if-aggregate.yang
  public/release/models/interfaces/openconfig-if-ethernet-ext.yang
  public/release/models/interfaces/openconfig-if-ethernet.yang
  public/release/models/interfaces/openconfig-if-ip-ext.yang
  public/release/models/interfaces/openconfig-if-ip.yang
  public/release/models/interfaces/openconfig-if-sdn-ext.yang
  public/release/models/interfaces/openconfig-interfaces.yang
  public/release/models/isis/openconfig-isis-policy.yang
  public/release/models/isis/openconfig-isis.yang
  public/release/models/lacp/openconfig-lacp.yang
  public/release/models/lldp/openconfig-lldp-types.yang
  public/release/models/lldp/openconfig-lldp.yang
  public/release/models/local-routing/openconfig-local-routing.yang
  public/release/models/mpls/openconfig-mpls-types.yang
  public/release/models/multicast/openconfig-pim.yang
  public/release/models/network-instance/openconfig-network-instance.yang
  public/release/models/openconfig-extensions.yang
  public/release/models/optical-transport/openconfig-transport-types.yang
  public/release/models/ospf/openconfig-ospf-policy.yang
  public/release/models/ospf/openconfig-ospfv2.yang
  public/release/models/p4rt/openconfig-p4rt.yang
  public/release/models/platform/openconfig-platform-common.yang
  public/release/models/platform/openconfig-platform-controller-card.yang
  public/release/models/platform/openconfig-platform-cpu.yang
  public/release/models/platform/openconfig-platform-ext.yang
  public/release/models/platform/openconfig-platform-fabric.yang
  public/release/models/platform/openconfig-platform-fan.yang
  public/release/models/platform/openconfig-platform-integrated-circuit.yang
  public/release/models/platform/openconfig-platform-linecard.yang
  public/release/models/platform/openconfig-platform-pipeline-counters.yang
  public/release/models/platform/openconfig-platform-psu.yang
  public/release/models/platform/openconfig-platform-software.yang
  public/release/models/platform/openconfig-platform-transceiver.yang
  public/release/models/platform/openconfig-platform.yang
  public/release/models/policy-forwarding/openconfig-policy-forwarding.yang
  public/release/models/policy/openconfig-policy-types.yang
  public/release/models/qos/openconfig-qos-elements.yang
  public/release/models/qos/openconfig-qos-interfaces.yang
  public/release/models/qos/openconfig-qos-types.yang
  public/release/models/qos/openconfig-qos.yang
  public/release/models/relay-agent/openconfig-relay-agent.yang
  public/release/models/rib/openconfig-rib-bgp.yang
  public/release/models/sampling/openconfig-sampling-sflow.yang
  public/release/models/segment-routing/openconfig-segment-routing-types.yang
  public/release/models/system/openconfig-system-bootz.yang
  public/release/models/system/openconfig-system-controlplane.yang
  public/release/models/system/openconfig-system-utilization.yang
  public/release/models/system/openconfig-system.yang
  public/release/models/types/openconfig-inet-types.yang
  public/release/models/types/openconfig-types.yang
  public/release/models/types/openconfig-yang-types.yang
  public/release/models/vlan/openconfig-vlan.yang
  public/third_party/ietf/iana-if-type.yang
  public/third_party/ietf/ietf-inet-types.yang
  public/third_party/ietf/ietf-interfaces.yang
  public/third_party/ietf/ietf-yang-types.yang
  yang/openconfig-bgp-gue.yang
)

rm -r oc || true
mkdir oc

go run github.com/openconfig/ygnmi/app/ygnmi generator \
  --trim_module_prefix=openconfig \
  --exclude_modules="${EXCLUDE_MODULES}" \
  --base_package_path=github.com/openconfig/lemming/gnmi/oc \
  --split_package_paths="/network-instances/network-instance/protocols/protocol/isis=netinstisis,/network-instances/network-instance/protocols/protocol/bgp=netinstbgp" \
  --structs_split_files_count=8 \
  --pathstructs_split_files_count=8 \
  --output_dir=oc \
  --paths=public/release/models/...,public/third_party/ietf/... \
  "${YANG_FILES[@]}"

find oc -name "*.go" -exec goimports -w {} +
find oc -name "*.go" -exec gofmt -w -s {} +
rm -rf public
