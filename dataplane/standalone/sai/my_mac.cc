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

#include "dataplane/standalone/sai/my_mac.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/my_mac.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_my_mac_api_t l_my_mac = {
    .create_my_mac = l_create_my_mac,
    .remove_my_mac = l_remove_my_mac,
    .set_my_mac_attribute = l_set_my_mac_attribute,
    .get_my_mac_attribute = l_get_my_mac_attribute,
};

lemming::dataplane::sai::CreateMyMacRequest convert_create_my_mac(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateMyMacRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MY_MAC_ATTR_PRIORITY:
        msg.set_priority(attr_list[i].value.u32);
        break;
      case SAI_MY_MAC_ATTR_PORT_ID:
        msg.set_port_id(attr_list[i].value.oid);
        break;
      case SAI_MY_MAC_ATTR_VLAN_ID:
        msg.set_vlan_id(attr_list[i].value.u16);
        break;
      case SAI_MY_MAC_ATTR_MAC_ADDRESS:
        msg.set_mac_address(attr_list[i].value.mac,
                            sizeof(attr_list[i].value.mac));
        break;
      case SAI_MY_MAC_ATTR_MAC_ADDRESS_MASK:
        msg.set_mac_address_mask(attr_list[i].value.mac,
                                 sizeof(attr_list[i].value.mac));
        break;
    }
  }
  return msg;
}

sai_status_t l_create_my_mac(sai_object_id_t *my_mac_id,
                             sai_object_id_t switch_id, uint32_t attr_count,
                             const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateMyMacRequest req =
      convert_create_my_mac(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateMyMacResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = my_mac->CreateMyMac(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  if (my_mac_id) {
    *my_mac_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_my_mac(sai_object_id_t my_mac_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveMyMacRequest req;
  lemming::dataplane::sai::RemoveMyMacResponse resp;
  grpc::ClientContext context;
  req.set_oid(my_mac_id);

  grpc::Status status = my_mac->RemoveMyMac(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_my_mac_attribute(sai_object_id_t my_mac_id,
                                    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetMyMacAttributeRequest req;
  lemming::dataplane::sai::SetMyMacAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(my_mac_id);

  switch (attr->id) {
    case SAI_MY_MAC_ATTR_PRIORITY:
      req.set_priority(attr->value.u32);
      break;
  }

  grpc::Status status = my_mac->SetMyMacAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_my_mac_attribute(sai_object_id_t my_mac_id,
                                    uint32_t attr_count,
                                    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMyMacAttributeRequest req;
  lemming::dataplane::sai::GetMyMacAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(my_mac_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_my_mac_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = my_mac->GetMyMacAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MY_MAC_ATTR_PRIORITY:
        attr_list[i].value.u32 = resp.attr().priority();
        break;
      case SAI_MY_MAC_ATTR_PORT_ID:
        attr_list[i].value.oid = resp.attr().port_id();
        break;
      case SAI_MY_MAC_ATTR_VLAN_ID:
        attr_list[i].value.u16 = resp.attr().vlan_id();
        break;
      case SAI_MY_MAC_ATTR_MAC_ADDRESS:
        memcpy(attr_list[i].value.mac, resp.attr().mac_address().data(),
               sizeof(sai_mac_t));
        break;
      case SAI_MY_MAC_ATTR_MAC_ADDRESS_MASK:
        memcpy(attr_list[i].value.mac, resp.attr().mac_address_mask().data(),
               sizeof(sai_mac_t));
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
