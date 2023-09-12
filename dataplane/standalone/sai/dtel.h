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

#ifndef DATAPLANE_STANDALONE_SAI_DTEL_H_
#define DATAPLANE_STANDALONE_SAI_DTEL_H_

extern "C" {
	#include "inc/sai.h"
	#include "experimental/saiextensions.h"
}

extern const sai_dtel_api_t l_dtel;


sai_status_t l_create_dtel(sai_object_id_t *dtel_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_dtel(sai_object_id_t dtel_id);

sai_status_t l_set_dtel_attribute(sai_object_id_t dtel_id, const sai_attribute_t *attr);

sai_status_t l_get_dtel_attribute(sai_object_id_t dtel_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_dtel_queue_report(sai_object_id_t *dtel_queue_report_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_dtel_queue_report(sai_object_id_t dtel_queue_report_id);

sai_status_t l_set_dtel_queue_report_attribute(sai_object_id_t dtel_queue_report_id, const sai_attribute_t *attr);

sai_status_t l_get_dtel_queue_report_attribute(sai_object_id_t dtel_queue_report_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_dtel_int_session(sai_object_id_t *dtel_int_session_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_dtel_int_session(sai_object_id_t dtel_int_session_id);

sai_status_t l_set_dtel_int_session_attribute(sai_object_id_t dtel_int_session_id, const sai_attribute_t *attr);

sai_status_t l_get_dtel_int_session_attribute(sai_object_id_t dtel_int_session_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_dtel_report_session(sai_object_id_t *dtel_report_session_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_dtel_report_session(sai_object_id_t dtel_report_session_id);

sai_status_t l_set_dtel_report_session_attribute(sai_object_id_t dtel_report_session_id, const sai_attribute_t *attr);

sai_status_t l_get_dtel_report_session_attribute(sai_object_id_t dtel_report_session_id, uint32_t attr_count, sai_attribute_t *attr_list);

sai_status_t l_create_dtel_event(sai_object_id_t *dtel_event_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list);

sai_status_t l_remove_dtel_event(sai_object_id_t dtel_event_id);

sai_status_t l_set_dtel_event_attribute(sai_object_id_t dtel_event_id, const sai_attribute_t *attr);

sai_status_t l_get_dtel_event_attribute(sai_object_id_t dtel_event_id, uint32_t attr_count, sai_attribute_t *attr_list);


#endif  // DATAPLANE_STANDALONE_SAI_DTEL_H_
