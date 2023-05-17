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

#include "dataplane/standalone/common.h"

#include <glog/logging.h>

extern "C" {
#include "common.h"
#include "inc/sai.h"
#include "meta/saimetadata.h"
}

// create a object with the given type and allocate a new id.
sai_object_id_t AttributeManager::create(sai_object_type_t type,
                                         std::string switch_id) {
  sai_object_id_t id = this->objects.size();
  this->objects[std::to_string(id)] = {
      .type = type,
      .switch_id = switch_id,
      .attributes = std::unordered_map<sai_attr_id_t, sai_attribute_value_t>(),
  };
  return id;
}

// create a object with the given type with a serialized version of the object
// used as the id.
std::string AttributeManager::create(sai_object_type_t type,
                                     common_entry_t entry) {
  std::string id = this->serialize_entry(type, entry);
  auto switch_id = this->entry_to_switch_id(type, entry);

  this->objects[id] = {
      .type = type,
      .switch_id = switch_id,
      .attributes = std::unordered_map<sai_attr_id_t, sai_attribute_value_t>(),
  };
  return id;
}

sai_object_type_t AttributeManager::get_type(std::string id) {
  auto iter = this->objects.find(id);
  if (iter == this->objects.end()) {
    return SAI_OBJECT_TYPE_NULL;
  }
  return iter->second.type;
}

std::string AttributeManager::get_switch_id(std::string id) {
  auto iter = this->objects.find(id);
  if (iter == this->objects.end()) {
    return SAI_NULL_OBJECT_ID;
  }
  return iter->second.switch_id;
}

void AttributeManager::set_attribute(std::string id,
                                     const sai_attribute_t* attr) {
  this->objects[id].attributes[attr->id] = attr->value;
}

void AttributeManager::set_attribute(std::string id, sai_attribute_t attr) {
  this->objects[id].attributes[attr.id] = attr.value;
}

sai_status_t AttributeManager::get_attribute(std::string id,
                                             uint32_t attr_count,
                                             sai_attribute_t* attr_list) {
  for (uint32_t i = 0; i < attr_count; i++) {
    LOG(INFO) << "Get Attr object id " << id << " attr id " << attr_list[i].id;
    auto iter = this->objects[id].attributes.find(attr_list[i].id);
    if (iter == this->objects[id].attributes.end()) {
      return SAI_STATUS_ITEM_NOT_FOUND;
    }
    attr_list[i].value = iter->second;
  }
  return SAI_STATUS_SUCCESS;
}

std::string AttributeManager::serialize_entry(sai_object_type_t type,
                                              common_entry_t id) {
  char serialize_buf[0x4000];
  switch (type) {
    case SAI_OBJECT_TYPE_FDB_ENTRY:
      sai_serialize_fdb_entry(serialize_buf, id.fdb_entry);
      break;
    case SAI_OBJECT_TYPE_INSEG_ENTRY:
      sai_serialize_inseg_entry(serialize_buf, id.inseg_entry);
      break;
    case SAI_OBJECT_TYPE_IPMC_ENTRY:
      sai_serialize_ipmc_entry(serialize_buf, id.ipmc_entry);
      break;
    case SAI_OBJECT_TYPE_L2MC_ENTRY:
      sai_serialize_l2mc_entry(serialize_buf, id.l2mc_entry);
      break;
    case SAI_OBJECT_TYPE_MCAST_FDB_ENTRY:
      sai_serialize_mcast_fdb_entry(serialize_buf, id.mcast_fdb_entry);
      break;
    case SAI_OBJECT_TYPE_NEIGHBOR_ENTRY:
      sai_serialize_neighbor_entry(serialize_buf, id.neighbor_entry);
      break;
    case SAI_OBJECT_TYPE_ROUTE_ENTRY:
      sai_serialize_route_entry(serialize_buf, id.route_entry);
      break;
    case SAI_OBJECT_TYPE_NAT_ENTRY:
      sai_serialize_nat_entry(serialize_buf, id.nat_entry);
      break;
    case SAI_OBJECT_TYPE_MY_SID_ENTRY:
      sai_serialize_my_sid_entry(serialize_buf, id.my_sid_entry);
      break;
    default:
      throw "Invalid type";
  }

  return std::string(serialize_buf);
}

std::string AttributeManager::entry_to_switch_id(sai_object_type_t type,
                                                 common_entry_t id) {
  sai_object_id_t swID;
  switch (type) {
    case SAI_OBJECT_TYPE_FDB_ENTRY:
      swID = id.fdb_entry->switch_id;
      break;
    case SAI_OBJECT_TYPE_INSEG_ENTRY:
      swID = id.inseg_entry->switch_id;
      break;
    case SAI_OBJECT_TYPE_IPMC_ENTRY:
      swID = id.ipmc_entry->switch_id;
      break;
    case SAI_OBJECT_TYPE_L2MC_ENTRY:
      swID = id.l2mc_entry->switch_id;
      break;
    case SAI_OBJECT_TYPE_MCAST_FDB_ENTRY:
      swID = id.mcast_fdb_entry->switch_id;
      break;
    case SAI_OBJECT_TYPE_NEIGHBOR_ENTRY:
      swID = id.neighbor_entry->switch_id;
      break;
    case SAI_OBJECT_TYPE_ROUTE_ENTRY:
      swID = id.route_entry->switch_id;
      break;
    case SAI_OBJECT_TYPE_NAT_ENTRY:
      swID = id.nat_entry->switch_id;
      break;
    case SAI_OBJECT_TYPE_MY_SID_ENTRY:
      swID = id.my_sid_entry->switch_id;
      break;
    default:
      throw "Invalid type";
  }
  return std::to_string(swID);
}

sai_status_t APIBase::create(common_entry_t id, uint32_t attr_count,
                             const sai_attribute_t* attr_list) {
  for (uint32_t i = 0; i < attr_count; i++) {
    this->attrMgr->set_attribute(this->id, attr_list[i]);
  }
  return SAI_STATUS_SUCCESS;
}

sai_status_t APIBase::create(uint32_t attr_count,
                             const sai_attribute_t* attr_list) {
  for (uint32_t i = 0; i < attr_count; i++) {
    this->attrMgr->set_attribute(this->id, attr_list[i]);
  }
  return SAI_STATUS_SUCCESS;
}

sai_status_t APIBase::set_attribute(const sai_attribute_t* attr) {
  this->attrMgr->set_attribute(this->id, attr);
  return SAI_STATUS_SUCCESS;
}
