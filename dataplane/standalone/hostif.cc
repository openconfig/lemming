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

#include "dataplane/standalone/hostif.h"

#include <glog/logging.h>

#include <vector>

#include "dataplane/standalone/translator.h"

sai_status_t HostIf::create(_In_ uint32_t attr_count,
                            _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);
  sai_object_id_t external_port = SAI_NULL_OBJECT_ID;
  std::string name = "";
  for (auto attr : attrs) {
    LOG(INFO) << "create hostif attr id " << attr.id;
    switch (attr.id) {
      case SAI_HOSTIF_ATTR_TYPE:
        if (attr.value.s32 ==
            SAI_HOSTIF_TYPE_GENETLINK) {  // TODO(dgrau): figure what this is.
          return SAI_STATUS_SUCCESS;
        } else if (attr.value.s32 != SAI_HOSTIF_TYPE_NETDEV) {
          return SAI_STATUS_NOT_SUPPORTED;
        } else {
          break;
        }
      case SAI_HOSTIF_ATTR_OBJ_ID:
        external_port = attr.value.oid;
        break;
      case SAI_HOSTIF_ATTR_NAME:
        name = std::string(attr.value.chardata);
        break;
    }
  }
  if (name == "" || external_port == SAI_NULL_OBJECT_ID) {
    return SAI_STATUS_MANDATORY_ATTRIBUTE_MISSING;
  }

  grpc::ClientContext context;
  lemming::dataplane::CreatePortRequest req;
  lemming::dataplane::CreatePortResponse resp;

  req.set_id(this->id);
  req.set_type(forwarding::PORT_TYPE_TAP);
  req.set_kernel_dev(name);

  // Link the ports, if the the corresponding port is present.
  sai_attribute_t portSt = {.id = SAI_PORT_ATTR_OPER_STATUS};
  this->attrMgr->get_attribute(std::to_string(external_port), 1, &portSt);
  if (portSt.value.s32 != SAI_PORT_OPER_STATUS_NOT_PRESENT) {
    req.set_external_port(std::to_string(external_port));
  }

  auto status = this->dataplane->CreatePort(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << "Failed to create port: " << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  APIBase::create(attrs.size(), attrs.data());
  return SAI_STATUS_SUCCESS;
}

sai_status_t HostIf::set_attribute(_In_ const sai_attribute_t* attr) {
  LOG(INFO) << "set hostif attr id " << attr->id;
  switch (attr->id) {
    case SAI_HOSTIF_ATTR_OPER_STATUS:
      grpc::ClientContext context;
      forwarding::PortStateRequest req;
      req.mutable_context_id()->set_id("lucius");
      req.mutable_port_id()->mutable_object_id()->set_id(this->id);

      if (attr->value.booldata) {
        LOG(INFO) << "Setting to hostif up, id " << this->id;
        req.mutable_operation()->set_admin_status(
            forwarding::PORT_STATE_ENABLED_UP);
      } else {
        LOG(INFO) << "Setting to hostif down, id " << this->id;
        req.mutable_operation()->set_admin_status(

            forwarding::PORT_STATE_DISABLED_DOWN);
      }

      forwarding::PortStateReply resp;
      auto status = this->fwd->PortState(&context, req, &resp);
      if (!status.ok()) {
        LOG(ERROR) << "Failed to hostif state: " << status.error_message();
        return SAI_STATUS_FAILURE;
      }
      break;
  }
  return SAI_STATUS_SUCCESS;
}

sai_status_t HostIfTableEntry::create(_In_ uint32_t attr_count,
                                      _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);
  APIBase::create(attrs.size(), attrs.data());
  return SAI_STATUS_SUCCESS;
}

sai_status_t HostIfTableEntry::set_attribute(_In_ const sai_attribute_t* attr) {
  return SAI_STATUS_SUCCESS;
}

sai_status_t HostIfTrap::create(_In_ uint32_t attr_count,
                                _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);
  APIBase::create(attrs.size(), attrs.data());
  return SAI_STATUS_SUCCESS;
}

sai_status_t HostIfTrap::set_attribute(_In_ const sai_attribute_t* attr) {
  return SAI_STATUS_SUCCESS;
}

sai_status_t HostIfTrapGroup::create(_In_ uint32_t attr_count,
                                     _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);
  APIBase::create(attrs.size(), attrs.data());
  return SAI_STATUS_SUCCESS;
}

sai_status_t HostIfTrapGroup::set_attribute(_In_ const sai_attribute_t* attr) {
  return SAI_STATUS_SUCCESS;
}
