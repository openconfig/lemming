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

#include "dataplane/standalone/sai/udf.h"

#include <glog/logging.h>

#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_udf_api_t l_udf = {
    .create_udf = l_create_udf,
    .remove_udf = l_remove_udf,
    .set_udf_attribute = l_set_udf_attribute,
    .get_udf_attribute = l_get_udf_attribute,
    .create_udf_match = l_create_udf_match,
    .remove_udf_match = l_remove_udf_match,
    .set_udf_match_attribute = l_set_udf_match_attribute,
    .get_udf_match_attribute = l_get_udf_match_attribute,
    .create_udf_group = l_create_udf_group,
    .remove_udf_group = l_remove_udf_group,
    .set_udf_group_attribute = l_set_udf_group_attribute,
    .get_udf_group_attribute = l_get_udf_group_attribute,
};

sai_status_t l_create_udf(sai_object_id_t *udf_id, sai_object_id_t switch_id,
                          uint32_t attr_count,
                          const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_UDF, udf_id, switch_id, attr_count,
                            attr_list);
}

sai_status_t l_remove_udf(sai_object_id_t udf_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_UDF, udf_id);
}

sai_status_t l_set_udf_attribute(sai_object_id_t udf_id,
                                 const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_UDF, udf_id, attr);
}

sai_status_t l_get_udf_attribute(sai_object_id_t udf_id, uint32_t attr_count,
                                 sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_UDF, udf_id, attr_count,
                                   attr_list);
}

sai_status_t l_create_udf_match(sai_object_id_t *udf_match_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_UDF_MATCH, udf_match_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_udf_match(sai_object_id_t udf_match_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_UDF_MATCH, udf_match_id);
}

sai_status_t l_set_udf_match_attribute(sai_object_id_t udf_match_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_UDF_MATCH, udf_match_id,
                                   attr);
}

sai_status_t l_get_udf_match_attribute(sai_object_id_t udf_match_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_UDF_MATCH, udf_match_id,
                                   attr_count, attr_list);
}

sai_status_t l_create_udf_group(sai_object_id_t *udf_group_id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->create(SAI_OBJECT_TYPE_UDF_GROUP, udf_group_id, switch_id,
                            attr_count, attr_list);
}

sai_status_t l_remove_udf_group(sai_object_id_t udf_group_id) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->remove(SAI_OBJECT_TYPE_UDF_GROUP, udf_group_id);
}

sai_status_t l_set_udf_group_attribute(sai_object_id_t udf_group_id,
                                       const sai_attribute_t *attr) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->set_attribute(SAI_OBJECT_TYPE_UDF_GROUP, udf_group_id,
                                   attr);
}

sai_status_t l_get_udf_group_attribute(sai_object_id_t udf_group_id,
                                       uint32_t attr_count,
                                       sai_attribute_t *attr_list) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
  return translator->get_attribute(SAI_OBJECT_TYPE_UDF_GROUP, udf_group_id,
                                   attr_count, attr_list);
}
