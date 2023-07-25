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

#include "dataplane/standalone/next_hop_group.h"

#include <glog/logging.h>

#include <bitset>
#include <string>
#include <vector>

#include "dataplane/standalone/translator.h"

sai_status_t NextHopGroup::create(_In_ uint32_t attr_count,
                                  _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);

  sai_next_hop_group_type_t type;
  for (auto attr : attrs) {
    switch (attr.id) {
      case SAI_NEXT_HOP_GROUP_ATTR_TYPE:
        type = static_cast<sai_next_hop_group_type_t>(attr.value.s32);
        break;
    }
  }
  if (type != SAI_NEXT_HOP_GROUP_TYPE_DYNAMIC_UNORDERED_ECMP) {
    return SAI_STATUS_NOT_SUPPORTED;
  }
  grpc::ClientContext context;
  lemming::dataplane::AddNextHopGroupRequest req;
  lemming::dataplane::AddNextHopGroupResponse resp;
  req.set_id(std::stoul(this->id));
  auto status = this->dataplane->AddNextHopGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << "Failed to create route: " << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  APIBase::create(attrs.size(), attrs.data());
  return SAI_STATUS_SUCCESS;
}

sai_status_t NextHopGroup::set_attribute(_In_ const sai_attribute_t* attr) {
  return SAI_STATUS_SUCCESS;
}

sai_status_t NextHopGroupMember::create(_In_ uint32_t attr_count,
                                        _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);

  sai_object_id_t group_id, hop_id;
  for (auto attr : attrs) {
    switch (attr.id) {
      case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_GROUP_ID:
        group_id = attr.value.oid;
        break;
      case SAI_NEXT_HOP_GROUP_MEMBER_ATTR_NEXT_HOP_ID:
        hop_id = attr.value.oid;
        break;
    }
  }
  grpc::ClientContext context;
  lemming::dataplane::AddNextHopGroupRequest req;
  lemming::dataplane::AddNextHopGroupResponse resp;
  req.set_id(group_id);
  req.mutable_list()->add_hops(hop_id);
  req.mutable_list()->add_weights(1);
  req.set_mode(lemming::dataplane::GROUP_UPDATE_MODE_APPEND);
  auto status = this->dataplane->AddNextHopGroup(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << "Failed to create route: " << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  APIBase::create(attrs.size(), attrs.data());
  return SAI_STATUS_SUCCESS;
}

sai_status_t NextHopGroupMember::set_attribute(
    _In_ const sai_attribute_t* attr) {
  return SAI_STATUS_SUCCESS;
}
