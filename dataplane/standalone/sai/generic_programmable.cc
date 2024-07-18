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

#include "dataplane/standalone/sai/generic_programmable.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/generic_programmable.pb.h"
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/enum.h"

const sai_generic_programmable_api_t l_generic_programmable = {
    .create_generic_programmable = l_create_generic_programmable,
    .remove_generic_programmable = l_remove_generic_programmable,
    .set_generic_programmable_attribute = l_set_generic_programmable_attribute,
    .get_generic_programmable_attribute = l_get_generic_programmable_attribute,
};

lemming::dataplane::sai::CreateGenericProgrammableRequest
convert_create_generic_programmable(sai_object_id_t switch_id,
                                    uint32_t attr_count,
                                    const sai_attribute_t *attr_list) {
  lemming::dataplane::sai::CreateGenericProgrammableRequest msg;

  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_GENERIC_PROGRAMMABLE_ATTR_OBJECT_NAME:
        msg.mutable_object_name()->Add(
            attr_list[i].value.s8list.list,
            attr_list[i].value.s8list.list + attr_list[i].value.s8list.count);
        break;
      case SAI_GENERIC_PROGRAMMABLE_ATTR_COUNTER_ID:
        msg.set_counter_id(attr_list[i].value.oid);
        break;
    }
  }
  return msg;
}

sai_status_t l_create_generic_programmable(
    sai_object_id_t *generic_programmable_id, sai_object_id_t switch_id,
    uint32_t attr_count, const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::CreateGenericProgrammableRequest req =
      convert_create_generic_programmable(switch_id, attr_count, attr_list);
  lemming::dataplane::sai::CreateGenericProgrammableResponse resp;
  grpc::ClientContext context;
  req.set_switch_(switch_id);

  grpc::Status status =
      generic_programmable->CreateGenericProgrammable(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  if (generic_programmable_id) {
    *generic_programmable_id = resp.oid();
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_remove_generic_programmable(
    sai_object_id_t generic_programmable_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::RemoveGenericProgrammableRequest req;
  lemming::dataplane::sai::RemoveGenericProgrammableResponse resp;
  grpc::ClientContext context;
  req.set_oid(generic_programmable_id);

  grpc::Status status =
      generic_programmable->RemoveGenericProgrammable(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_set_generic_programmable_attribute(
    sai_object_id_t generic_programmable_id, const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::SetGenericProgrammableAttributeRequest req;
  lemming::dataplane::sai::SetGenericProgrammableAttributeResponse resp;
  grpc::ClientContext context;
  req.set_oid(generic_programmable_id);

  switch (attr->id) {
    case SAI_GENERIC_PROGRAMMABLE_ATTR_COUNTER_ID:
      req.set_counter_id(attr->value.oid);
      break;
  }

  grpc::Status status = generic_programmable->SetGenericProgrammableAttribute(
      &context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }

  return SAI_STATUS_SUCCESS;
}

sai_status_t l_get_generic_programmable_attribute(
    sai_object_id_t generic_programmable_id, uint32_t attr_count,
    sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  lemming::dataplane::sai::GetGenericProgrammableAttributeRequest req;
  lemming::dataplane::sai::GetGenericProgrammableAttributeResponse resp;
  grpc::ClientContext context;

  req.set_oid(generic_programmable_id);

  for (uint32_t i = 0; i < attr_count; i++) {
    req.add_attr_type(
        convert_sai_generic_programmable_attr_t_to_proto(attr_list[i].id));
  }
  grpc::Status status = generic_programmable->GetGenericProgrammableAttribute(
      &context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_STATUS_FAILURE;
  }
  for (uint32_t i = 0; i < attr_count; i++) {
    switch (attr_list[i].id) {
      case SAI_GENERIC_PROGRAMMABLE_ATTR_OBJECT_NAME:
        copy_list(attr_list[i].value.s8list.list, resp.attr().object_name(),
                  &attr_list[i].value.s8list.count);
        break;
      case SAI_GENERIC_PROGRAMMABLE_ATTR_COUNTER_ID:
        attr_list[i].value.oid = resp.attr().counter_id();
        break;
    }
  }

  return SAI_STATUS_SUCCESS;
}
