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

#include "dataplane/standalone/translator.h"

#include <glog/logging.h>

#include <string>
#include <unordered_map>

#include "dataplane/standalone/hostif.h"
#include "dataplane/standalone/switch.h"

extern "C" {
#include "inc/sai.h"
#include "meta/saimetadata.h"
}

sai_object_type_t Translator::getObjectType(sai_object_id_t id) {
  return this->attrMgr->get_type(std::to_string(id));
}

// create globally scoped objects.
sai_status_t Translator::create(sai_object_type_t type, sai_object_id_t* id,
                                uint32_t attr_count,
                                const sai_attribute_t* attr_list) {
  *id = this->attrMgr->create(type, "0");
  sai_status_t status = 0;
  switch (type) {
    case SAI_OBJECT_TYPE_SWITCH: {
      auto sw = std::make_shared<Switch>(std::to_string(*id), this->attrMgr,
                                         this->fwd, this->dataplane, this);
      this->switches[std::to_string(*id)] = sw;
      this->apis[std::to_string(*id)] = sw;
      status = this->apis[std::to_string(*id)]->create(attr_count, attr_list);
      break;
    }
    case SAI_OBJECT_TYPE_HOSTIF_TRAP_GROUP:
      this->apis[std::to_string(*id)] = std::make_unique<HostIfTrapGroup>(
          std::to_string(*id), this->attrMgr, this->fwd, this->dataplane);
      break;
    default:
      LOG(INFO) << "Unknown type " << type;
      // TODO(dgrau): handle other types
  }

  if (status != SAI_STATUS_SUCCESS) {
    return status;
  }
  // Save attributes to attr managers.
  for (uint32_t i = 0; i < attr_count; i++) {
    this->attrMgr->set_attribute(std::to_string(*id), attr_list[i]);
  }

  return SAI_STATUS_SUCCESS;
}

// create switch-scoped objects.
sai_status_t Translator::create(sai_object_type_t type, sai_object_id_t* id,
                                sai_object_id_t switch_id, uint32_t attr_count,
                                const sai_attribute_t* attr_list) {
  *id = this->attrMgr->create(
      type, std::to_string(switch_id));  // Allocate new id for object.
  auto status = this->switches[std::to_string(switch_id)]->create_child(
      type, *id, attr_count,
      attr_list);  // Delegate creation to switch instance.
  if (status != SAI_STATUS_SUCCESS) {
    return status;
  }
  return SAI_STATUS_SUCCESS;
}

// create entry type objects.
sai_status_t Translator::create(sai_object_type_t type, common_entry_t entry,
                                uint32_t attr_count,
                                const sai_attribute_t* attr_list) {
  std::string idStr = this->attrMgr->create(type, entry);
  auto switch_id = this->attrMgr->entry_to_switch_id(type, entry);

  auto status = this->switches[switch_id]->create_child(
      type, entry, attr_count,
      attr_list);  // Delegate creation to switch instance.
  if (status != SAI_STATUS_SUCCESS) {
    return status;
  }
  return SAI_STATUS_SUCCESS;
}

sai_status_t Translator::remove(sai_object_type_t type, sai_object_id_t id) {
  return sai_status_t();
}

sai_status_t Translator::remove(sai_object_type_t type, common_entry_t id) {
  return sai_status_t();
}

sai_status_t Translator::set_attribute(sai_object_type_t type,
                                       sai_object_id_t id,
                                       const sai_attribute_t* attr) {
  std::string switch_id = this->attrMgr->get_switch_id(std::to_string(id));

  sai_status_t status;
  if (switch_id != "0") {
    status = this->switches[switch_id]->set_child_attr(type, std::to_string(id),
                                                       attr);
  } else {
    status = this->apis[std::to_string(id)]->set_attribute(attr);
  }
  if (status == SAI_STATUS_SUCCESS) {
    this->attrMgr->set_attribute(std::to_string(id), attr);
  }
  return status;
}

sai_status_t Translator::set_attribute(sai_object_type_t type,
                                       common_entry_t entry,
                                       const sai_attribute_t* attr) {
  std::string idStr = this->attrMgr->serialize_entry(type, entry);
  auto switch_id = this->attrMgr->entry_to_switch_id(type, entry);

  sai_status_t status;
  if (switch_id != "0") {
    status = this->switches[switch_id]->set_child_attr(type, idStr, attr);
  } else {
    status = this->apis[idStr]->set_attribute(attr);
  }
  if (status == SAI_STATUS_SUCCESS) {
    this->attrMgr->set_attribute(idStr, attr);
  }
  return status;
}

sai_status_t Translator::get_attribute(sai_object_type_t type,
                                       sai_object_id_t id, uint32_t attr_count,
                                       sai_attribute_t* attr_list) {
  return this->attrMgr->get_attribute(std::to_string(id), attr_count,
                                      attr_list);
}

sai_status_t Translator::get_attribute(sai_object_type_t type,
                                       common_entry_t id, uint32_t attr_count,
                                       sai_attribute_t* attr_list) {
  std::string idStr = this->attrMgr->serialize_entry(type, id);
  return this->attrMgr->get_attribute(idStr, attr_count, attr_list);
}

sai_status_t Translator::get_stats(sai_object_type_t type, sai_object_id_t id,
                                   uint32_t number_of_counters,
                                   const sai_stat_id_t* counter_ids,
                                   uint64_t* counters) {
  return sai_status_t();
}

sai_status_t Translator::get_stats_ext(sai_object_type_t type,
                                       sai_object_id_t bfd_session_id,
                                       uint32_t number_of_counters,
                                       const sai_stat_id_t* counter_ids,
                                       sai_stats_mode_t mode,
                                       uint64_t* counters) {
  return sai_status_t();
}

sai_status_t Translator::clear_stats(sai_object_type_t type,
                                     sai_object_id_t bfd_session_id,
                                     uint32_t number_of_counters,
                                     const sai_stat_id_t* counter_ids) {
  return sai_status_t();
}

sai_status_t Translator::create_bulk(
    sai_object_type_t type, sai_object_id_t switch_id, uint32_t object_count,
    const uint32_t* attr_count, const sai_attribute_t** attr_list,
    sai_bulk_op_error_mode_t mode, sai_object_id_t* object_id,
    sai_status_t* object_statuses) {
  return sai_status_t();
}

sai_status_t Translator::remove_bulk(sai_object_type_t type,
                                     uint32_t object_count,
                                     const sai_object_id_t* object_id,
                                     sai_bulk_op_error_mode_t mode,
                                     sai_status_t* object_statuses) {
  return sai_status_t();
}

sai_status_t Translator::create_bulk(
    sai_object_type_t type, uint32_t object_count, common_entry_t object_id,
    const uint32_t* attr_count, const sai_attribute_t** attr_list,
    sai_bulk_op_error_mode_t mode, sai_status_t* object_statuses) {
  return sai_status_t();
}

sai_status_t Translator::remove_bulk(sai_object_type_t type,
                                     uint32_t object_count,
                                     common_entry_t object_id,
                                     sai_bulk_op_error_mode_t mode,
                                     sai_status_t* object_statuses) {
  return sai_status_t();
}

sai_status_t Translator::set_attribute_bulk(sai_object_type_t type,
                                            uint32_t object_count,
                                            common_entry_t object_id,
                                            const sai_attribute_t* attr_list,
                                            sai_bulk_op_error_mode_t mode,
                                            sai_status_t* object_statuses) {
  return sai_status_t();
}

sai_status_t Translator::get_attribute_bulk(
    sai_object_type_t type, uint32_t object_count, common_entry_t object_id,
    const uint32_t* attr_count, sai_attribute_t** attr_list,
    sai_bulk_op_error_mode_t mode, sai_status_t* object_statuses) {
  return sai_status_t();
}
