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

#ifndef DATAPLANE_STANDALONE_SAI_ENTRY_H_
#define DATAPLANE_STANDALONE_SAI_ENTRY_H_

extern "C" {
#include "inc/sai.h"
}

typedef union _common_entry_t {
  const sai_fdb_entry_t* fdb_entry;
  const sai_inseg_entry_t* inseg_entry;
  const sai_ipmc_entry_t* ipmc_entry;
  const sai_l2mc_entry_t* l2mc_entry;
  const sai_mcast_fdb_entry_t* mcast_fdb_entry;
  const sai_neighbor_entry_t* neighbor_entry;
  const sai_route_entry_t* route_entry;
  const sai_nat_entry_t* nat_entry;
  const sai_my_sid_entry_t* my_sid_entry;
} common_entry_t;

#endif  // DATAPLANE_STANDALONE_SAI_ENTRY_H_
