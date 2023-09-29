

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

#include "dataplane/standalone/sai/srv6.h"

#include <glog/logging.h>

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/srv6.pb.h"
#include "dataplane/standalone/sai/common.h"

const sai_srv6_api_t l_srv6 = {
    .create_srv6_sidlist = l_create_srv6_sidlist,
    .remove_srv6_sidlist = l_remove_srv6_sidlist,
    .set_srv6_sidlist_attribute = l_set_srv6_sidlist_attribute,
    .get_srv6_sidlist_attribute = l_get_srv6_sidlist_attribute,
    .create_srv6_sidlists = l_create_srv6_sidlists,
    .remove_srv6_sidlists = l_remove_srv6_sidlists,
    .create_my_sid_entry = l_create_my_sid_entry,
    .remove_my_sid_entry = l_remove_my_sid_entry,
    .set_my_sid_entry_attribute = l_set_my_sid_entry_attribute,
    .get_my_sid_entry_attribute = l_get_my_sid_entry_attribute,
    .create_my_sid_entries = l_create_my_sid_entries,
    .remove_my_sid_entries = l_remove_my_sid_entries,
    .set_my_sid_entries_attribute = l_set_my_sid_entries_attribute,
    .get_my_sid_entries_attribute = l_get_my_sid_entries_attribute,
};

lemming::dataplane::sai::CreateSrv6SidlistRequest convert_create_srv6_sidlist(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateSrv6SidlistRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SRV6_SIDLIST_ATTR_TYPE:
        msg.set_type(static_cast<lemming::dataplane::sai::Srv6SidlistType>(
            attr_list[i].value.s32 + 1));
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateMySidEntryRequest convert_create_my_sid_entry(
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateMySidEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR:
        msg.set_endpoint_behavior(
            static_cast<lemming::dataplane::sai::MySidEntryEndpointBehavior>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR_FLAVOR:
        msg.set_endpoint_behavior_flavor(
            static_cast<
                lemming::dataplane::sai::MySidEntryEndpointBehaviorFlavor>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_MY_SID_ENTRY_ATTR_PACKET_ACTION:
        msg.set_packet_action(
            static_cast<lemming::dataplane::sai::PacketAction>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_MY_SID_ENTRY_ATTR_TRAP_PRIORITY:
        msg.set_trap_priority(attr_list[i].value.u8);
        break;
      case SAI_MY_SID_ENTRY_ATTR_NEXT_HOP_ID:
        msg.set_next_hop_id(attr_list[i].value.oid);
        break;
      case SAI_MY_SID_ENTRY_ATTR_TUNNEL_ID:
        msg.set_tunnel_id(attr_list[i].value.oid);
        break;
      case SAI_MY_SID_ENTRY_ATTR_VRF:
        msg.set_vrf(attr_list[i].value.oid);
        break;
      case SAI_MY_SID_ENTRY_ATTR_COUNTER_ID:
        msg.set_counter_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_srv6_sidlist(sai_object_id_t *srv6_sidlist_id,
                                   sai_object_id_t switch_id,
                                   uint32_t attr_count,
                                   const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateSrv6SidlistRequest req =
      convert_create_srv6_sidlist(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateSrv6SidlistResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = srv6->CreateSrv6Sidlist(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *srv6_sidlist_id = resp.oid();

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_srv6_sidlist(sai_object_id_t srv6_sidlist_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveSrv6SidlistRequest req;
  lemming::dataplane::sai::RemoveSrv6SidlistResponse resp;
  grpc::ClientContext context;
  req.set_oid(srv6_sidlist_id);

  grpc::Status status = srv6->RemoveSrv6Sidlist(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_srv6_sidlist_attribute(sai_object_id_t srv6_sidlist_id,
                                          const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_srv6_sidlist_attribute(sai_object_id_t srv6_sidlist_id,
                                          uint32_t attr_count,
                                          sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetSrv6SidlistAttributeRequest req;
  lemming::dataplane::sai::GetSrv6SidlistAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(srv6_sidlist_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::Srv6SidlistAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = srv6->GetSrv6SidlistAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_SRV6_SIDLIST_ATTR_TYPE:
        attr_list[i].value.s32 = static_cast<int>(resp.attr().type() - 1);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_srv6_sidlists(sai_object_id_t switch_id,
                                    uint32_t object_count,
                                    const uint32_t *attr_count,
                                    const sai_attribute_t **attr_list,
                                    sai_bulk_op_error_mode_t mode,
                                    sai_object_id_t *object_id,
                                    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateSrv6SidlistsRequest req;
  lemming::dataplane::sai::CreateSrv6SidlistsResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    auto r =
        convert_create_srv6_sidlist(switch_id, attr_count[i], attr_list[i]);
    *req.add_reqs() = r;
  }

  grpc::Status status = srv6->CreateSrv6Sidlists(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < object_count; i++) {
    switch_id = object_id[i] = resp.resps(i).oid();
    object_statuses[i] = SAI_STATUS_SUCCESS;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_srv6_sidlists(uint32_t object_count,
                                    const sai_object_id_t *object_id,
                                    sai_bulk_op_error_mode_t mode,
                                    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_create_my_sid_entry(const sai_my_sid_entry_t *my_sid_entry,
                                   uint32_t attr_count,
                                   const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateMySidEntryRequest req =
      convert_create_my_sid_entry(attr_count, attr_list);
  lemming::dataplane::sai::CreateMySidEntryResponse resp;
  grpc::ClientContext context;

  grpc::Status status = srv6->CreateMySidEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_my_sid_entry(const sai_my_sid_entry_t *my_sid_entry) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveMySidEntryRequest req;
  lemming::dataplane::sai::RemoveMySidEntryResponse resp;
  grpc::ClientContext context;

  grpc::Status status = srv6->RemoveMySidEntry(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_my_sid_entry_attribute(
    const sai_my_sid_entry_t *my_sid_entry, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetMySidEntryAttributeRequest req;
  lemming::dataplane::sai::SetMySidEntryAttributeResponse resp;
  grpc::ClientContext context;

  switch (attr->id) {
    case SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR:
      req.set_endpoint_behavior(
          static_cast<lemming::dataplane::sai::MySidEntryEndpointBehavior>(
              attr->value.s32 + 1));
      break;
    case SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR_FLAVOR:
      req.set_endpoint_behavior_flavor(
          static_cast<
              lemming::dataplane::sai::MySidEntryEndpointBehaviorFlavor>(
              attr->value.s32 + 1));
      break;
    case SAI_MY_SID_ENTRY_ATTR_PACKET_ACTION:
      req.set_packet_action(static_cast<lemming::dataplane::sai::PacketAction>(
          attr->value.s32 + 1));
      break;
    case SAI_MY_SID_ENTRY_ATTR_TRAP_PRIORITY:
      req.set_trap_priority(attr->value.u8);
      break;
    case SAI_MY_SID_ENTRY_ATTR_NEXT_HOP_ID:
      req.set_next_hop_id(attr->value.oid);
      break;
    case SAI_MY_SID_ENTRY_ATTR_TUNNEL_ID:
      req.set_tunnel_id(attr->value.oid);
      break;
    case SAI_MY_SID_ENTRY_ATTR_VRF:
      req.set_vrf(attr->value.oid);
      break;
    case SAI_MY_SID_ENTRY_ATTR_COUNTER_ID:
      req.set_counter_id(attr->value.oid);
      break;
  }

  grpc::Status status = srv6->SetMySidEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_my_sid_entry_attribute(
    const sai_my_sid_entry_t *my_sid_entry, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetMySidEntryAttributeRequest req;
  lemming::dataplane::sai::GetMySidEntryAttributeResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(static_cast<lemming::dataplane::sai::MySidEntryAttr>(
        attr_list[i].id + 1));
  }
  grpc::Status status = srv6->GetMySidEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().endpoint_behavior() - 1);
        break;
      case SAI_MY_SID_ENTRY_ATTR_ENDPOINT_BEHAVIOR_FLAVOR:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().endpoint_behavior_flavor() - 1);
        break;
      case SAI_MY_SID_ENTRY_ATTR_PACKET_ACTION:
        attr_list[i].value.s32 =
            static_cast<int>(resp.attr().packet_action() - 1);
        break;
      case SAI_MY_SID_ENTRY_ATTR_TRAP_PRIORITY:
        attr_list[i].value.u8 = resp.attr().trap_priority();
        break;
      case SAI_MY_SID_ENTRY_ATTR_NEXT_HOP_ID:
        attr_list[i].value.oid = resp.attr().next_hop_id();
        break;
      case SAI_MY_SID_ENTRY_ATTR_TUNNEL_ID:
        attr_list[i].value.oid = resp.attr().tunnel_id();
        break;
      case SAI_MY_SID_ENTRY_ATTR_VRF:
        attr_list[i].value.oid = resp.attr().vrf();
        break;
      case SAI_MY_SID_ENTRY_ATTR_COUNTER_ID:
        attr_list[i].value.oid = resp.attr().counter_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_my_sid_entries(uint32_t object_count,
                                     const sai_my_sid_entry_t *my_sid_entry,
                                     const uint32_t *attr_count,
                                     const sai_attribute_t **attr_list,
                                     sai_bulk_op_error_mode_t mode,
                                     sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateMySidEntriesRequest req;
  lemming::dataplane::sai::CreateMySidEntriesResponse resp;
  grpc::ClientContext context;

  for (uint32_t i = 0; i < object_count; i++) {
    auto r = convert_create_my_sid_entry(attr_count[i], attr_list[i]);
    *req.add_reqs() = r;
  }

  grpc::Status status = srv6->CreateMySidEntries(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < object_count; i++) {
    object_statuses[i] = SAI_STATUS_SUCCESS;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_my_sid_entries(uint32_t object_count,
                                     const sai_my_sid_entry_t *my_sid_entry,
                                     sai_bulk_op_error_mode_t mode,
                                     sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_set_my_sid_entries_attribute(
    uint32_t object_count, const sai_my_sid_entry_t *my_sid_entry,
    const sai_attribute_t *attr_list, sai_bulk_op_error_mode_t mode,
    sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t l_get_my_sid_entries_attribute(
    uint32_t object_count, const sai_my_sid_entry_t *my_sid_entry,
    const uint32_t *attr_count, sai_attribute_t **attr_list,
    sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return SAI_STATUS_NOT_IMPLEMENTED;
}
