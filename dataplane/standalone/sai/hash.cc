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

#include "dataplane/standalone/sai/hash.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/hash.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_hash_api_t l_hash = {
    .create_hash = l_create_hash,
    .remove_hash = l_remove_hash,
    .set_hash_attribute = l_set_hash_attribute,
    .get_hash_attribute = l_get_hash_attribute,
    .create_fine_grained_hash_field = l_create_fine_grained_hash_field,
    .remove_fine_grained_hash_field = l_remove_fine_grained_hash_field,
    .set_fine_grained_hash_field_attribute =
        l_set_fine_grained_hash_field_attribute,
    .get_fine_grained_hash_field_attribute =
        l_get_fine_grained_hash_field_attribute,
};

lemming::dataplane::sai::CreateHashRequest convert_create_hash(
    sai_object_id_t switch_id, uint32_t attr_count,
    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateHashRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HASH_ATTR_UDF_GROUP_LIST:
        msg.mutable_udf_group_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_HASH_ATTR_FINE_GRAINED_HASH_FIELD_LIST:
        msg.mutable_fine_grained_hash_field_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
    }
  }
  return msg;
}

lemming::dataplane::sai::CreateFineGrainedHashFieldRequest
convert_create_fine_grained_hash_field(sai_object_id_t switch_id,
                                       uint32_t attr_count,
                                       const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateFineGrainedHashFieldRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_NATIVE_HASH_FIELD:
        msg.set_native_hash_field(
            convert_sai_native_hash_field_t_to_proto(attr_list[i].value.s32));
        break;
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_IPV4_MASK:
        msg.set_ipv4_mask(&attr_list[i].value.ip4,
                          sizeof(attr_list[i].value.ip4));
        break;
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_IPV6_MASK:
        msg.set_ipv6_mask(attr_list[i].value.ip6,
                          sizeof(attr_list[i].value.ip6));
        break;
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_SEQUENCE_ID:
        msg.set_sequence_id(attr_list[i].value.u32);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_hash(sai_object_id_t *hash_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHashRequest req =
      convert_create_hash(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateHashResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = hash->CreateHash(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  if (hash_id) {
    *hash_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_hash(sai_object_id_t hash_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveHashRequest req;
  lemming::dataplane::sai::RemoveHashResponse resp;
  grpc::ClientContext context;
  req.set_oid(hash_id);

  grpc::Status status = hash->RemoveHash(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_hash_attribute(sai_object_id_t hash_id,
                                  const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetHashAttributeRequest req;
  lemming::dataplane::sai::SetHashAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(hash_id);

  switch (attr->id) {
    case SAI_HASH_ATTR_UDF_GROUP_LIST:
      req.mutable_udf_group_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
    case SAI_HASH_ATTR_FINE_GRAINED_HASH_FIELD_LIST:
      req.mutable_fine_grained_hash_field_list()->Add(
          attr->value.objlist.list,
          attr->value.objlist.list + attr->value.objlist.count);
      break;
  }

  grpc::Status status = hash->SetHashAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_hash_attribute(sai_object_id_t hash_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetHashAttributeRequest req;
  lemming::dataplane::sai::GetHashAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(hash_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(convert_sai_hash_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = hash->GetHashAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HASH_ATTR_UDF_GROUP_LIST:
        copy_list(attr_list[i].value.objlist.list, resp.attr().udf_group_list(),
                  &attr_list[i].value.objlist.count);
        break;
      case SAI_HASH_ATTR_FINE_GRAINED_HASH_FIELD_LIST:
        copy_list(attr_list[i].value.objlist.list,
                  resp.attr().fine_grained_hash_field_list(),
                  &attr_list[i].value.objlist.count);
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_create_fine_grained_hash_field(
    sai_object_id_t *fine_grained_hash_field_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateFineGrainedHashFieldRequest req =
      convert_create_fine_grained_hash_field(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateFineGrainedHashFieldResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status = hash->CreateFineGrainedHashField(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  if (fine_grained_hash_field_id) {
    *fine_grained_hash_field_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_fine_grained_hash_field(
    sai_object_id_t fine_grained_hash_field_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveFineGrainedHashFieldRequest req;
  lemming::dataplane::sai::RemoveFineGrainedHashFieldResponse resp;
  grpc::ClientContext context;
  req.set_oid(fine_grained_hash_field_id);

  grpc::Status status = hash->RemoveFineGrainedHashField(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_fine_grained_hash_field_attribute(
    sai_object_id_t fine_grained_hash_field_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_fine_grained_hash_field_attribute(
    sai_object_id_t fine_grained_hash_field_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetFineGrainedHashFieldAttributeRequest req;
  lemming::dataplane::sai::GetFineGrainedHashFieldAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(fine_grained_hash_field_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_fine_grained_hash_field_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status =
      hash->GetFineGrainedHashFieldAttribute(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_NATIVE_HASH_FIELD:
        attr_list[i].value.s32 = convert_sai_native_hash_field_t_to_sai(
            resp.attr().native_hash_field());
        break;
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_IPV4_MASK:
        memcpy(&attr_list[i].value.ip4, resp.attr().ipv4_mask().data(),
               sizeof(sai_ip4_t));
        break;
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_IPV6_MASK:
        memcpy(attr_list[i].value.ip6, resp.attr().ipv6_mask().data(),
               sizeof(sai_ip6_t));
        break;
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_SEQUENCE_ID:
        attr_list[i].value.u32 = resp.attr().sequence_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
