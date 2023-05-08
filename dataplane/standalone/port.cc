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

#include "dataplane/standalone/log/log.h"
#include "dataplane/standalone/port.h"
#include "dataplane/standalone/translator.h"

sai_status_t Port::create_port(_Out_ sai_object_id_t *port_id,
                                      _In_ sai_object_id_t switch_id,
                                      _In_ uint32_t attr_count,
                                      _In_ const sai_attribute_t *attr_list) {
  *port_id = this->translator->createObject(SAI_OBJECT_TYPE_PORT);
  for (uint32_t i = 0; i < attr_count; i++) {
    LOG(attr_list[i].id);
    this->translator->setAttribute(*port_id, sai_attribute_t{
                                                   .id = attr_list[i].id,
                                                   .value = attr_list[i].value,
                                               });
    switch (attr_list[i].id) {
      // TODO(dgrau): handle this attributes specially.
    }
  }

  LOG("created switch");
  return SAI_STATUS_SUCCESS;
}

sai_status_t Port::set_port_attribute(_In_ sai_object_id_t switch_id,
                                        _In_ const sai_attribute_t *attr) {
  this->translator->setAttribute(switch_id, *attr);
  return SAI_STATUS_SUCCESS;
}

sai_status_t Port::get_port_attribute(_In_ sai_object_id_t switch_id,
                                        _In_ uint32_t attr_count,
                                        _Inout_ sai_attribute_t *attr_list) {
  for (uint32_t i = 0; i < attr_count; i++) {
    LOG(attr_list[i].id);
    if (auto ret = this->translator->getAttribute(switch_id, &attr_list[i]);
        ret != SAI_STATUS_SUCCESS) {
      return ret;
    }
  }
  return SAI_STATUS_SUCCESS;
}
