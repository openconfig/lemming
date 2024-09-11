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

#include "dataplane/standalone/sai/bmtor.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/bmtor.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_bmtor_api_t l_bmtor = {
    .create_table_bitmap_classification_entry =
        l_create_table_bitmap_classification_entry,
    .remove_table_bitmap_classification_entry =
        l_remove_table_bitmap_classification_entry,
    .set_table_bitmap_classification_entry_attribute =
        l_set_table_bitmap_classification_entry_attribute,
    .get_table_bitmap_classification_entry_attribute =
        l_get_table_bitmap_classification_entry_attribute,
    .get_table_bitmap_classification_entry_stats =
        l_get_table_bitmap_classification_entry_stats,
    .get_table_bitmap_classification_entry_stats_ext =
        l_get_table_bitmap_classification_entry_stats_ext,
    .clear_table_bitmap_classification_entry_stats =
        l_clear_table_bitmap_classification_entry_stats,
    .create_table_bitmap_router_entry = l_create_table_bitmap_router_entry,
    .remove_table_bitmap_router_entry = l_remove_table_bitmap_router_entry,
    .set_table_bitmap_router_entry_attribute =
        l_set_table_bitmap_router_entry_attribute,
    .get_table_bitmap_router_entry_attribute =
        l_get_table_bitmap_router_entry_attribute,
    .get_table_bitmap_router_entry_stats =
        l_get_table_bitmap_router_entry_stats,
    .get_table_bitmap_router_entry_stats_ext =
        l_get_table_bitmap_router_entry_stats_ext,
    .clear_table_bitmap_router_entry_stats =
        l_clear_table_bitmap_router_entry_stats,
    .create_table_meta_tunnel_entry = l_create_table_meta_tunnel_entry,
    .remove_table_meta_tunnel_entry = l_remove_table_meta_tunnel_entry,
    .set_table_meta_tunnel_entry_attribute =
        l_set_table_meta_tunnel_entry_attribute,
    .get_table_meta_tunnel_entry_attribute =
        l_get_table_meta_tunnel_entry_attribute,
    .get_table_meta_tunnel_entry_stats = l_get_table_meta_tunnel_entry_stats,
    .get_table_meta_tunnel_entry_stats_ext =
        l_get_table_meta_tunnel_entry_stats_ext,
    .clear_table_meta_tunnel_entry_stats =
        l_clear_table_meta_tunnel_entry_stats,
};

lemming::dataplane::sai::CreateTableBitmapClassificationEntryRequest
convert_create_table_bitmap_classification_entry(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateTableBitmapClassificationEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ACTION:
        msg.set_action(
            convert_sai_table_bitmap_classification_entry_action_t_to_proto(
                attr_list[i].value.s32));
        break;
      case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ROUTER_INTERFACE_KEY:
        msg.set_router_interface_key(attr_list[i].value.oid);
        break;
      case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IS_DEFAULT:
        msg.set_is_default(attr_list[i].value.booldata);
        break;
      case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IN_RIF_METADATA:
        msg.set_in_rif_metadata(attr_list[i].value.u32);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateTableBitmapRouterEntryRequest
convert_create_table_bitmap_router_entry(sai_object_id_t switch_id,
                                         uint32_t attr_count,
                                         const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateTableBitmapRouterEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_ACTION:
        msg.set_action(convert_sai_table_bitmap_router_entry_action_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_PRIORITY:
        msg.set_priority(attr_list[i].value.u32);
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_KEY:
        msg.set_in_rif_metadata_key(attr_list[i].value.u32);
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_MASK:
        msg.set_in_rif_metadata_mask(attr_list[i].value.u32);
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_TUNNEL_INDEX:
        msg.set_tunnel_index(attr_list[i].value.u16);
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_NEXT_HOP:
        msg.set_next_hop(attr_list[i].value.oid);
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_ROUTER_INTERFACE:
        msg.set_router_interface(attr_list[i].value.oid);
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_TRAP_ID:
        msg.set_trap_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateTableMetaTunnelEntryRequest
convert_create_table_meta_tunnel_entry(sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateTableMetaTunnelEntryRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_ACTION:
        msg.set_action(convert_sai_table_meta_tunnel_entry_action_t_to_proto(
            attr_list[i].value.s32));
        break;
      case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_METADATA_KEY:
        msg.set_metadata_key(attr_list[i].value.u16);
        break;
      case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_IS_DEFAULT:
        msg.set_is_default(attr_list[i].value.booldata);
        break;
      case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_TUNNEL_ID:
        msg.set_tunnel_id(attr_list[i].value.oid);
        break;
      case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_UNDERLAY_DIP:
        msg.set_underlay_dip(
            convert_from_ip_address(attr_list[i].value.ipaddr));
        break;
    }
  }
  return msg;
}

sai_status_t l_create_table_bitmap_classification_entry(
    sai_object_id_t *table_bitmap_classification_entry_id,
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTableBitmapClassificationEntryRequest req =
      convert_create_table_bitmap_classification_entry(switch_id, attr_count,
                                                       attr_list);
  lemming::dataplane::sai::CreateTableBitmapClassificationEntryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status =
      bmtor->CreateTableBitmapClassificationEntry(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (table_bitmap_classification_entry_id) {
    *table_bitmap_classification_entry_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_table_bitmap_classification_entry(
    sai_object_id_t table_bitmap_classification_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveTableBitmapClassificationEntryRequest req;
  lemming::dataplane::sai::RemoveTableBitmapClassificationEntryResponse resp;
  grpc::ClientContext context;
  req.set_oid(table_bitmap_classification_entry_id);

  grpc::Status status =
      bmtor->RemoveTableBitmapClassificationEntry(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_table_bitmap_classification_entry_attribute(
    sai_object_id_t table_bitmap_classification_entry_id,
    const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_table_bitmap_classification_entry_attribute(
    sai_object_id_t table_bitmap_classification_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetTableBitmapClassificationEntryAttributeRequest
      req;
  lemming::dataplane::sai::GetTableBitmapClassificationEntryAttributeResponse
      resp;
  grpc::ClientContext context;

  req.set_oid(table_bitmap_classification_entry_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_table_bitmap_classification_entry_attr_t_to_proto(
            attr_list[i].id));
  }
  grpc::Status status =
      bmtor->GetTableBitmapClassificationEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ACTION:
        attr_list[i].value.s32 =
            convert_sai_table_bitmap_classification_entry_action_t_to_sai(
                resp.attr().action());
        break;
      case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_ROUTER_INTERFACE_KEY:
        attr_list[i].value.oid = resp.attr().router_interface_key();
        break;
      case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IS_DEFAULT:
        attr_list[i].value.booldata = resp.attr().is_default();
        break;
      case SAI_TABLE_BITMAP_CLASSIFICATION_ENTRY_ATTR_IN_RIF_METADATA:
        attr_list[i].value.u32 = resp.attr().in_rif_metadata();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_table_bitmap_classification_entry_stats(
    sai_object_id_t table_bitmap_classification_entry_id,
    uint32_t number_of_counters, const sai_stat_id_t *counter_ids,
    uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetTableBitmapClassificationEntryStatsRequest req;
  lemming::dataplane::sai::GetTableBitmapClassificationEntryStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(table_bitmap_classification_entry_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(
        convert_sai_table_bitmap_classification_entry_stat_t_to_proto(
            counter_ids[i]));
  }
  grpc::Status status =
      bmtor->GetTableBitmapClassificationEntryStats(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_table_bitmap_classification_entry_stats_ext(
    sai_object_id_t table_bitmap_classification_entry_id,
    uint32_t number_of_counters, const sai_stat_id_t *counter_ids,
    sai_stats_mode_t mode, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_table_bitmap_classification_entry_stats(
    sai_object_id_t table_bitmap_classification_entry_id,
    uint32_t number_of_counters, const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_table_bitmap_router_entry(
    sai_object_id_t *table_bitmap_router_entry_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTableBitmapRouterEntryRequest req =
      convert_create_table_bitmap_router_entry(switch_id, attr_count,
                                               attr_list);
  lemming::dataplane::sai::CreateTableBitmapRouterEntryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status =
      bmtor->CreateTableBitmapRouterEntry(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (table_bitmap_router_entry_id) {
    *table_bitmap_router_entry_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_table_bitmap_router_entry(
    sai_object_id_t table_bitmap_router_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveTableBitmapRouterEntryRequest req;
  lemming::dataplane::sai::RemoveTableBitmapRouterEntryResponse resp;
  grpc::ClientContext context;
  req.set_oid(table_bitmap_router_entry_id);

  grpc::Status status =
      bmtor->RemoveTableBitmapRouterEntry(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_table_bitmap_router_entry_attribute(
    sai_object_id_t table_bitmap_router_entry_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_table_bitmap_router_entry_attribute(
    sai_object_id_t table_bitmap_router_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetTableBitmapRouterEntryAttributeRequest req;
  lemming::dataplane::sai::GetTableBitmapRouterEntryAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(table_bitmap_router_entry_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_table_bitmap_router_entry_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status =
      bmtor->GetTableBitmapRouterEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_ACTION:
        attr_list[i].value.s32 =
            convert_sai_table_bitmap_router_entry_action_t_to_sai(
                resp.attr().action());
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_PRIORITY:
        attr_list[i].value.u32 = resp.attr().priority();
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_KEY:
        attr_list[i].value.u32 = resp.attr().in_rif_metadata_key();
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_IN_RIF_METADATA_MASK:
        attr_list[i].value.u32 = resp.attr().in_rif_metadata_mask();
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_TUNNEL_INDEX:
        attr_list[i].value.u16 = resp.attr().tunnel_index();
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_NEXT_HOP:
        attr_list[i].value.oid = resp.attr().next_hop();
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_ROUTER_INTERFACE:
        attr_list[i].value.oid = resp.attr().router_interface();
        break;
      case SAI_TABLE_BITMAP_ROUTER_ENTRY_ATTR_TRAP_ID:
        attr_list[i].value.oid = resp.attr().trap_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_table_bitmap_router_entry_stats(
    sai_object_id_t table_bitmap_router_entry_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetTableBitmapRouterEntryStatsRequest req;
  lemming::dataplane::sai::GetTableBitmapRouterEntryStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(table_bitmap_router_entry_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(
        convert_sai_table_bitmap_router_entry_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status =
      bmtor->GetTableBitmapRouterEntryStats(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_table_bitmap_router_entry_stats_ext(
    sai_object_id_t table_bitmap_router_entry_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids, sai_stats_mode_t mode,
    uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_table_bitmap_router_entry_stats(
    sai_object_id_t table_bitmap_router_entry_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_table_meta_tunnel_entry(
    sai_object_id_t *table_meta_tunnel_entry_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateTableMetaTunnelEntryRequest req =
      convert_create_table_meta_tunnel_entry(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateTableMetaTunnelEntryResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = bmtor->CreateTableMetaTunnelEntry(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  if (table_meta_tunnel_entry_id) {
    *table_meta_tunnel_entry_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_table_meta_tunnel_entry(
    sai_object_id_t table_meta_tunnel_entry_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveTableMetaTunnelEntryRequest req;
  lemming::dataplane::sai::RemoveTableMetaTunnelEntryResponse resp;
  grpc::ClientContext context;
  req.set_oid(table_meta_tunnel_entry_id);

  grpc::Status status = bmtor->RemoveTableMetaTunnelEntry(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_table_meta_tunnel_entry_attribute(
    sai_object_id_t table_meta_tunnel_entry_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_table_meta_tunnel_entry_attribute(
    sai_object_id_t table_meta_tunnel_entry_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetTableMetaTunnelEntryAttributeRequest req;
  lemming::dataplane::sai::GetTableMetaTunnelEntryAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(table_meta_tunnel_entry_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_table_meta_tunnel_entry_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status =
      bmtor->GetTableMetaTunnelEntryAttribute(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_ACTION:
        attr_list[i].value.s32 =
            convert_sai_table_meta_tunnel_entry_action_t_to_sai(
                resp.attr().action());
        break;
      case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_METADATA_KEY:
        attr_list[i].value.u16 = resp.attr().metadata_key();
        break;
      case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_IS_DEFAULT:
        attr_list[i].value.booldata = resp.attr().is_default();
        break;
      case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_TUNNEL_ID:
        attr_list[i].value.oid = resp.attr().tunnel_id();
        break;
      case SAI_TABLE_META_TUNNEL_ENTRY_ATTR_UNDERLAY_DIP:
        attr_list[i].value.ipaddr =
            convert_to_ip_address(resp.attr().underlay_dip());
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_table_meta_tunnel_entry_stats(
    sai_object_id_t table_meta_tunnel_entry_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids, uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetTableMetaTunnelEntryStatsRequest req;
  lemming::dataplane::sai::GetTableMetaTunnelEntryStatsResponse resp;
  grpc::ClientContext context;
  req.set_oid(table_meta_tunnel_entry_id);

  for (uint32_t i = 0; i < number_of_counters; i++) {
    req.add_counter_ids(
        convert_sai_table_meta_tunnel_entry_stat_t_to_proto(counter_ids[i]));
  }
  grpc::Status status =
      bmtor->GetTableMetaTunnelEntryStats(&context, req, &resp);
  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Trace ID " << it->second << " " << status.error_message();
    } else {
      LOG(ERROR) << status.error_message();
    }
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0;
       i < number_of_counters && i < uint32_t(resp.values_size()); i++) {
    counters[i] = resp.values(i);
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_table_meta_tunnel_entry_stats_ext(
    sai_object_id_t table_meta_tunnel_entry_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids, sai_stats_mode_t mode,
    uint64_t *counters) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_clear_table_meta_tunnel_entry_stats(
    sai_object_id_t table_meta_tunnel_entry_id, uint32_t number_of_counters,
    const sai_stat_id_t *counter_ids) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}
