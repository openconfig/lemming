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
git clone https://github.com/openconfig/gnsi.git

EXCLUDE_MODULES=ietf-interfaces,openconfig-bfd,openconfig-messages

YANG_FILES=(
  gnsi/yang/gnsi-telemetry.yang
  public/release/models/aft/openconfig-aft.yang
  public/release/models/bgp/openconfig-bgp-policy.yang
  public/release/models/bgp/openconfig-bgp-types.yang
  public/release/models/interfaces/openconfig-if-aggregate.yang
  public/release/models/interfaces/openconfig-if-ethernet.yang
  public/release/models/interfaces/openconfig-if-ip-ext.yang
  public/release/models/interfaces/openconfig-if-ip.yang
  public/release/models/interfaces/openconfig-interfaces.yang
  public/release/models/lacp/openconfig-lacp.yang
  public/release/models/lldp/openconfig-lldp-types.yang
  public/release/models/lldp/openconfig-lldp.yang
  public/release/models/local-routing/openconfig-local-routing.yang
  public/release/models/mpls/openconfig-mpls-types.yang
  public/release/models/network-instance/openconfig-network-instance.yang
  public/release/models/openconfig-extensions.yang
  public/release/models/optical-transport/openconfig-transport-types.yang
  public/release/models/platform/openconfig-platform-cpu.yang
  public/release/models/platform/openconfig-platform-integrated-circuit.yang
  public/release/models/platform/openconfig-platform-software.yang
  public/release/models/platform/openconfig-platform-transceiver.yang
  public/release/models/platform/openconfig-platform.yang
  public/release/models/rib/openconfig-rib-bgp.yang
  public/release/models/system/openconfig-system.yang
  public/release/models/types/openconfig-inet-types.yang
  public/release/models/types/openconfig-types.yang
  public/release/models/types/openconfig-yang-types.yang
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
  --output_dir=oc \
  --paths=public/release/models/...,public/third_party/ietf/...,gnsi/... \
  "${YANG_FILES[@]}"

find oc -name "*.go" -exec goimports -w {} +
find oc -name "*.go" -exec gofmt -w -s {} +
rm -rf public gnsi
