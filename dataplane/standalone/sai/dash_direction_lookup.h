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

#ifndef DATAPLANE_STANDALONE_SAI_DASH_DIRECTION_LOOKUP_H_
#define DATAPLANE_STANDALONE_SAI_DASH_DIRECTION_LOOKUP_H_

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

extern const sai_dash_direction_lookup_api_t l_dash_direction_lookup;

#endif  // DATAPLANE_STANDALONE_SAI_DASH_DIRECTION_LOOKUP_H_
