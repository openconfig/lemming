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

#ifndef DATAPLANE_STANDALONE_TRANSLATOR_H_
#define DATAPLANE_STANDALONE_TRANSLATOR_H_

#include <grpcpp/channel.h>
#include <grpcpp/security/credentials.h>

#include <memory>
#include <string>
#include <unordered_map>
#include <utility>

#include "dataplane/standalone/common.h"
#include "dataplane/standalone/sai/entry.h"
#include "dataplane/standalone/switch.h"
#include "proto/dataplane/dataplane.grpc.pb.h"
#include "proto/dataplane/dataplane.pb.h"
#include "proto/forwarding/forwarding_service.grpc.pb.h"
#include "proto/forwarding/forwarding_service.pb.h"

extern "C" {
#include "inc/sai.h"
}

class Translator {
 public:
  explicit Translator(std::shared_ptr<grpc::Channel> chan) {
    attrMgr = std::make_shared<AttributeManager>();
    fwd = std::shared_ptr<forwarding::Forwarding::Stub>(
        forwarding::Forwarding::NewStub(chan));
    dataplane = std::shared_ptr<lemming::dataplane::Dataplane::Stub>(
        lemming::dataplane::Dataplane::NewStub(chan));
  }
  sai_object_type_t getObjectType(sai_object_id_t id);

  sai_status_t create(sai_object_type_t type, sai_object_id_t *id,
                      uint32_t attr_count, const sai_attribute_t *attr_list);

  sai_status_t create(sai_object_type_t type, common_entry_t id,
                      uint32_t attr_count, const sai_attribute_t *attr_list);

  sai_status_t create(sai_object_type_t type, sai_object_id_t *id,
                      sai_object_id_t switch_id, uint32_t attr_count,
                      const sai_attribute_t *attr_list);

  sai_status_t remove(sai_object_type_t type, sai_object_id_t id);
  sai_status_t remove(sai_object_type_t type, common_entry_t id);

  sai_status_t set_attribute(sai_object_type_t type, sai_object_id_t id,
                             const sai_attribute_t *attr);
  sai_status_t get_attribute(sai_object_type_t type, sai_object_id_t id,
                             uint32_t attr_count, sai_attribute_t *attr_list);

  sai_status_t set_attribute(sai_object_type_t type, common_entry_t id,
                             const sai_attribute_t *attr);
  sai_status_t get_attribute(sai_object_type_t type, common_entry_t id,
                             uint32_t attr_count, sai_attribute_t *attr_list);

  sai_status_t get_stats(sai_object_type_t type, sai_object_id_t id,
                         uint32_t number_of_counters,
                         const sai_stat_id_t *counter_ids, uint64_t *counters);
  sai_status_t get_stats_ext(sai_object_type_t type,
                             sai_object_id_t bfd_session_id,
                             uint32_t number_of_counters,
                             const sai_stat_id_t *counter_ids,
                             sai_stats_mode_t mode, uint64_t *counters);
  sai_status_t clear_stats(sai_object_type_t type,
                           sai_object_id_t bfd_session_id,
                           uint32_t number_of_counters,
                           const sai_stat_id_t *counter_ids);

  sai_status_t create_bulk(sai_object_type_t type, sai_object_id_t switch_id,
                           uint32_t object_count, const uint32_t *attr_count,
                           const sai_attribute_t **attr_list,
                           sai_bulk_op_error_mode_t mode,
                           sai_object_id_t *object_id,
                           sai_status_t *object_statuses);
  sai_status_t remove_bulk(sai_object_type_t type, uint32_t object_count,
                           const sai_object_id_t *object_id,
                           sai_bulk_op_error_mode_t mode,
                           sai_status_t *object_statuses);

  sai_status_t create_bulk(sai_object_type_t type, uint32_t object_count,
                           common_entry_t object_id, const uint32_t *attr_count,
                           const sai_attribute_t **attr_list,
                           sai_bulk_op_error_mode_t mode,
                           sai_status_t *object_statuses);

  sai_status_t remove_bulk(sai_object_type_t type, uint32_t object_count,
                           common_entry_t object_id,
                           sai_bulk_op_error_mode_t mode,
                           sai_status_t *object_statuses);

  sai_status_t set_attribute_bulk(sai_object_type_t type, uint32_t object_count,
                                  common_entry_t object_id,
                                  const sai_attribute_t *attr_list,
                                  sai_bulk_op_error_mode_t mode,
                                  sai_status_t *object_statuses);

  sai_status_t get_attribute_bulk(sai_object_type_t type, uint32_t object_count,
                                  common_entry_t object_id,
                                  const uint32_t *attr_count,
                                  sai_attribute_t **attr_list,
                                  sai_bulk_op_error_mode_t mode,
                                  sai_status_t *object_statuses);

  std::shared_ptr<AttributeManager> attrMgr;

 private:
  std::shared_ptr<forwarding::Forwarding::Stub> fwd;
  std::shared_ptr<lemming::dataplane::Dataplane::Stub> dataplane;
  std::unordered_map<std::string, std::shared_ptr<Switch>> switches;
  std::unordered_map<std::string, std::shared_ptr<APIBase>>
      apis;  // TODO(dgrau): Confirm that switch is the only global API and
             // remove.
};

#endif  // DATAPLANE_STANDALONE_TRANSLATOR_H_
