// Copyright 2023 Google LLC
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

#ifndef DATAPLANE_STANDALONE_SAI_COMMON_H_
#define DATAPLANE_STANDALONE_SAI_COMMON_H_

#include <algorithm>
#include <memory>
#include <string>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/translator.h"

extern "C" {
#include "inc/sai.h"
}

extern std::shared_ptr<Translator> translator;

std::string convert_from_ip_addr(sai_ip_addr_family_t addr_family,
                                 const sai_ip_addr_t &addr);
std::string convert_from_ip_address(const sai_ip_address_t &val);
lemming::dataplane::sai::RouteEntry convert_from_route_entry(
    const sai_route_entry_t &entry);
lemming::dataplane::sai::IpPrefix convert_from_ip_prefix(
    const sai_ip_prefix_t &ip_prefix);

sai_ip_addr_t convert_to_ip_addr(std::string val);
sai_ip_address_t convert_to_ip_address(std::string str);
sai_route_entry_t convert_to_route_entry(
    const lemming::dataplane::sai::RouteEntry &entry);
sai_ip_prefix_t convert_to_ip_prefix(
    const lemming::dataplane::sai::IpPrefix &ip_prefix);

// copy_list copies a scalar proto list to an attribute.
// Note: It is expected that the attribute list contains preallocated memory.
template <typename T, typename S>
void copy_list(S *dst, const google::protobuf::RepeatedField<T> &src,
               int attr_len) {
  // It's not safe to just memcpy this because in some cases to proto types are
  // larger than the corresponding sai types.
  for (int i = 0; i < std::min(attr_len, src.size()); i++) {
    dst[i] = src[i];
  }
}

#endif  // DATAPLANE_STANDALONE_SAI_COMMON_H_
