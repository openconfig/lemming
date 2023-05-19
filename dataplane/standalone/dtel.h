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

#ifndef DATAPLANE_STANDALONE_DTEL_H_
#define DATAPLANE_STANDALONE_DTEL_H_

#include <memory>
#include <string>
#include <unordered_map>

#include "dataplane/standalone/common.h"
#include "proto/dataplane/dataplane.grpc.pb.h"
#include "proto/dataplane/dataplane.pb.h"
#include "proto/forwarding/forwarding_service.grpc.pb.h"
#include "proto/forwarding/forwarding_service.pb.h"

extern "C" {
#include "inc/sai.h"
}

class DTEL : public APIBase {
 public:
  DTEL(std::string id, std::shared_ptr<AttributeManager> mgr,
       std::shared_ptr<forwarding::Forwarding::Stub> fwd,
       std::shared_ptr<lemming::dataplane::Dataplane::Stub> dplane)
      : APIBase(id, mgr, fwd, dplane) {}
  ~DTEL() = default;
  sai_status_t create(_In_ uint32_t attr_count,
                      _In_ const sai_attribute_t* attr_list);
  sai_status_t set_attribute(_In_ const sai_attribute_t* attr);
};

#endif  // DATAPLANE_STANDALONE_DTEL_H_
