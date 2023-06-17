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

#include "dataplane/standalone/router_interface.h"

#include <glog/logging.h>

#include <vector>

#include "dataplane/standalone/translator.h"

sai_status_t RouterInterface::create(_In_ uint32_t attr_count,
                                     _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);
  sai_router_interface_type_t type;
  for (auto attr : attrs) {
    switch (attr.id) {
      // TODO(dgrau): Handle vr and mtu.
      case SAI_ROUTER_INTERFACE_ATTR_VIRTUAL_ROUTER_ID:
        break;
      case SAI_ROUTER_INTERFACE_ATTR_MTU:
        break;
      case SAI_ROUTER_INTERFACE_ATTR_TYPE:
        type = static_cast<sai_router_interface_type_t>(attr.value.s32);
        break;
    }
  }
  if (type == SAI_ROUTER_INTERFACE_TYPE_LOOPBACK) {
    grpc::ClientContext context;
    lemming::dataplane::CreatePortRequest req;
    lemming::dataplane::CreatePortResponse resp;
    req.set_id(this->id);
    req.set_type(forwarding::PORT_TYPE_KERNEL);
    req.set_kernel_dev("lo");
    auto status = this->dataplane->CreatePort(&context, req, &resp);
    if (!status.ok()) {
      LOG(ERROR) << "Failed to create port: " << status.error_message();
      return SAI_STATUS_FAILURE;
    }
  }
  return SAI_STATUS_SUCCESS;
}

sai_status_t RouterInterface::set_attribute(_In_ const sai_attribute_t* attr) {
  return SAI_STATUS_SUCCESS;
}
