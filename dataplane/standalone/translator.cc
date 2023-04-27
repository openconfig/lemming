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

#include "translator.h"

using namespace Translator;

sai_status_t Translator::create_switch(_Out_ sai_object_id_t *switch_id, _In_ uint32_t attr_count, _In_ const sai_attribute_t *attr_list)
{
    return SAI_STATUS_SUCCESS;
}