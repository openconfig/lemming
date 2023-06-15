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

#include "dataplane/standalone/port.h"

#include <glog/logging.h>

#include <fstream>
#include <string>
#include <vector>

#include "absl/strings/str_split.h"
#include "dataplane/standalone/translator.h"

sai_status_t Port::create(_In_ uint32_t attr_count,
                          _In_ const sai_attribute_t* attr_list) {
  std::vector<sai_attribute_t> attrs(attr_list, attr_list + attr_count);
  std::vector<int> lanes;
  std::string name;
  sai_port_type_t type;
  for (auto attr : attrs) {
    switch (attr.id) {
      case SAI_PORT_ATTR_TYPE:
        type = static_cast<sai_port_type_t>(attr.value.s32);
        break;
      case SAI_PORT_ATTR_HW_LANE_LIST:
        lanes = std::vector<int>(
            attr.value.u32list.list,
            attr.value.u32list.list + attr.value.u32list.count);

        for (auto port : Port::laneMap) {
          if (port.second == lanes) {
            name = port.first;
            break;
          }
        }
        break;
    }
  }
  grpc::ClientContext context;
  lemming::dataplane::CreatePortRequest req;
  lemming::dataplane::CreatePortResponse resp;
  req.set_id(this->id);
  if (type == SAI_PORT_TYPE_CPU) {
    req.set_type(forwarding::PORT_TYPE_CPU_PORT);
    auto status = this->dataplane->CreatePort(&context, req, &resp);
    if (!status.ok()) {
      LOG(ERROR) << "Failed to create port: " << status.error_message();
      return SAI_STATUS_FAILURE;
    }
    attrs.push_back({
        .id = SAI_PORT_ATTR_OPER_STATUS,
        .value = {.s32 = SAI_PORT_OPER_STATUS_UP},
    });
  } else {
    if (name != "") {
      req.set_type(forwarding::PORT_TYPE_KERNEL);
      req.set_kernel_dev(name);
      auto status = this->dataplane->CreatePort(&context, req, &resp);
      if (!status.ok()) {
        LOG(ERROR) << "Failed to create port: " << status.error_message();
        return SAI_STATUS_FAILURE;
      }
      attrs.push_back({
          .id = SAI_PORT_ATTR_OPER_STATUS,
          .value = {.s32 = SAI_PORT_OPER_STATUS_UP},
      });
      LOG(INFO) << "Created port with id " << this->id;
    } else {  // TODO(dgrau): Figure out what to do for this ports.
      attrs.push_back({
          .id = SAI_PORT_ATTR_OPER_STATUS,
          .value = {.s32 = SAI_PORT_OPER_STATUS_NOT_PRESENT},
      });
      LOG(WARNING) << "Skipped port for SAI interface without kernel device"
                   << std::to_string(lanes[0]);
    }
  }

  attrs.push_back({
      .id = SAI_PORT_ATTR_NUMBER_OF_INGRESS_PRIORITY_GROUPS,
      .value = {.u32 = 0},
  });
  attrs.push_back({
      .id = SAI_PORT_ATTR_QOS_NUMBER_OF_QUEUES,
      .value = {.u32 = 0},
  });
  attrs.push_back({
      .id = SAI_PORT_ATTR_QOS_MAXIMUM_HEADROOM_SIZE,
      .value = {.u32 = 0},
  });
  attrs.push_back({
      .id = SAI_PORT_ATTR_ADMIN_STATE,
      .value = {.booldata = true},
  });
  attrs.push_back({
      .id = SAI_PORT_ATTR_AUTO_NEG_MODE,
      .value = {.booldata = false},
  });
  attrs.push_back({
      .id = SAI_PORT_ATTR_MTU,
      .value = {.u32 = 1514},
  });
  attrs.push_back({
      .id = SAI_PORT_ATTR_SUPPORTED_SPEED,
      .value = {.u32list = {.count = 0}},
  });
  attrs.push_back({
      .id = SAI_PORT_ATTR_SUPPORTED_FEC_MODE,
      .value = {.s32list = {.count = 0}},
  });
  APIBase::create(attrs.size(), attrs.data());
  return SAI_STATUS_SUCCESS;
}

sai_status_t Port::set_attribute(_In_ const sai_attribute_t* attr) {
  return SAI_STATUS_SUCCESS;
}

std::unordered_map<std::string, std::vector<int>> Port::parseLaneMap() {
  std::ifstream file("/usr/share/sonic/hwsku/lanemap.ini");
  std::string line;
  std::unordered_map<std::string, std::vector<int>> res;
  while (std::getline(file, line)) {
    std::vector<std::string> elems = absl::StrSplit(line, ":");
    std::vector<std::string> lanes = absl::StrSplit(elems[1], ",");
    for (auto lane : lanes) {
      res[elems[0]].push_back(std::stoi(lane));
    }
  }
  return res;
}

std::unordered_map<std::string, std::vector<int>> Port::laneMap =
    Port::parseLaneMap();
