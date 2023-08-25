

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

#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/hash.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

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

sai_status_t l_create_hash(sai_object_id_t *hash_id, sai_object_id_t switch_id,
                           uint32_t attr_count,
                           const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateHashRequest req;
  lemming::dataplane::sai::CreateHashResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_HASH_ATTR_UDF_GROUP_LIST:
        req.mutable_udf_group_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
      case SAI_HASH_ATTR_FINE_GRAINED_HASH_FIELD_LIST:
        req.mutable_fine_grained_hash_field_list()->Add(
            attr_list[i].value.objlist.list,
            attr_list[i].value.objlist.list + attr_list[i].value.objlist.count);
        break;
    }
  }
  grpc::Status status = hash->CreateHash(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *hash_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_HASH, hash_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_hash(sai_object_id_t hash_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_HASH, hash_id);
}

sai_status_t l_set_hash_attribute(sai_object_id_t hash_id,
                                  const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_HASH, hash_id, attr);
}

sai_status_t l_get_hash_attribute(sai_object_id_t hash_id, uint32_t attr_count,
                                  sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_HASH, hash_id, attr_count,
                                   attr_list);
}

sai_status_t l_create_fine_grained_hash_field(
    sai_object_id_t *fine_grained_hash_field_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateFineGrainedHashFieldRequest req;
  lemming::dataplane::sai::CreateFineGrainedHashFieldResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_NATIVE_HASH_FIELD:
        req.set_native_hash_field(
            static_cast<lemming::dataplane::sai::NativeHashField>(
                attr_list[i].value.s32 + 1));
        break;
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_IPV4_MASK:
        req.set_ipv4_mask(&attr_list[i].value.ip4,
                          sizeof(attr_list[i].value.ip4));
        break;
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_IPV6_MASK:
        req.set_ipv6_mask(attr_list[i].value.ip6,
                          sizeof(attr_list[i].value.ip6));
        break;
      case SAI_FINE_GRAINED_HASH_FIELD_ATTR_SEQUENCE_ID:
        req.set_sequence_id(attr_list[i].value.u32);
        break;
    }
  }
  grpc::Status status = hash->CreateFineGrainedHashField(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  *fine_grained_hash_field_id = resp.oid();

  return translator->create(SAI_OBJECT_TYPE_FINE_GRAINED_HASH_FIELD,
                            fine_grained_hash_field_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_fine_grained_hash_field(
    sai_object_id_t fine_grained_hash_field_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->remove(SAI_OBJECT_TYPE_FINE_GRAINED_HASH_FIELD,
                            fine_grained_hash_field_id);
}

sai_status_t l_set_fine_grained_hash_field_attribute(
    sai_object_id_t fine_grained_hash_field_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->set_attribute(SAI_OBJECT_TYPE_FINE_GRAINED_HASH_FIELD,
                                   fine_grained_hash_field_id, attr);
}

sai_status_t l_get_fine_grained_hash_field_attribute(
    sai_object_id_t fine_grained_hash_field_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return translator->get_attribute(SAI_OBJECT_TYPE_FINE_GRAINED_HASH_FIELD,
                                   fine_grained_hash_field_id, attr_count,
                                   attr_list);
}
