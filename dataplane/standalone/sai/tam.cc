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

#include "dataplane/standalone/sai/tam.h"
#include <glog/logging.h>
#include "dataplane/standalone/sai/common.h"
#include "dataplane/standalone/sai/entry.h"

const sai_tam_api_t l_tam = {
	.create_tam = l_create_tam,
	.remove_tam = l_remove_tam,
	.set_tam_attribute = l_set_tam_attribute,
	.get_tam_attribute = l_get_tam_attribute,
	.create_tam_math_func = l_create_tam_math_func,
	.remove_tam_math_func = l_remove_tam_math_func,
	.set_tam_math_func_attribute = l_set_tam_math_func_attribute,
	.get_tam_math_func_attribute = l_get_tam_math_func_attribute,
	.create_tam_report = l_create_tam_report,
	.remove_tam_report = l_remove_tam_report,
	.set_tam_report_attribute = l_set_tam_report_attribute,
	.get_tam_report_attribute = l_get_tam_report_attribute,
	.create_tam_event_threshold = l_create_tam_event_threshold,
	.remove_tam_event_threshold = l_remove_tam_event_threshold,
	.set_tam_event_threshold_attribute = l_set_tam_event_threshold_attribute,
	.get_tam_event_threshold_attribute = l_get_tam_event_threshold_attribute,
	.create_tam_int = l_create_tam_int,
	.remove_tam_int = l_remove_tam_int,
	.set_tam_int_attribute = l_set_tam_int_attribute,
	.get_tam_int_attribute = l_get_tam_int_attribute,
	.create_tam_tel_type = l_create_tam_tel_type,
	.remove_tam_tel_type = l_remove_tam_tel_type,
	.set_tam_tel_type_attribute = l_set_tam_tel_type_attribute,
	.get_tam_tel_type_attribute = l_get_tam_tel_type_attribute,
	.create_tam_transport = l_create_tam_transport,
	.remove_tam_transport = l_remove_tam_transport,
	.set_tam_transport_attribute = l_set_tam_transport_attribute,
	.get_tam_transport_attribute = l_get_tam_transport_attribute,
	.create_tam_telemetry = l_create_tam_telemetry,
	.remove_tam_telemetry = l_remove_tam_telemetry,
	.set_tam_telemetry_attribute = l_set_tam_telemetry_attribute,
	.get_tam_telemetry_attribute = l_get_tam_telemetry_attribute,
	.create_tam_collector = l_create_tam_collector,
	.remove_tam_collector = l_remove_tam_collector,
	.set_tam_collector_attribute = l_set_tam_collector_attribute,
	.get_tam_collector_attribute = l_get_tam_collector_attribute,
	.create_tam_event_action = l_create_tam_event_action,
	.remove_tam_event_action = l_remove_tam_event_action,
	.set_tam_event_action_attribute = l_set_tam_event_action_attribute,
	.get_tam_event_action_attribute = l_get_tam_event_action_attribute,
	.create_tam_event = l_create_tam_event,
	.remove_tam_event = l_remove_tam_event,
	.set_tam_event_attribute = l_set_tam_event_attribute,
	.get_tam_event_attribute = l_get_tam_event_attribute,
};


sai_status_t l_create_tam(sai_object_id_t *tam_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM, tam_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam(sai_object_id_t tam_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM, tam_id);
}

sai_status_t l_set_tam_attribute(sai_object_id_t tam_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM, tam_id, attr);
}

sai_status_t l_get_tam_attribute(sai_object_id_t tam_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM, tam_id, attr_count, attr_list);
}

sai_status_t l_create_tam_math_func(sai_object_id_t *tam_math_func_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM_MATH_FUNC, tam_math_func_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_math_func(sai_object_id_t tam_math_func_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM_MATH_FUNC, tam_math_func_id);
}

sai_status_t l_set_tam_math_func_attribute(sai_object_id_t tam_math_func_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM_MATH_FUNC, tam_math_func_id, attr);
}

sai_status_t l_get_tam_math_func_attribute(sai_object_id_t tam_math_func_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM_MATH_FUNC, tam_math_func_id, attr_count, attr_list);
}

sai_status_t l_create_tam_report(sai_object_id_t *tam_report_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM_REPORT, tam_report_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_report(sai_object_id_t tam_report_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM_REPORT, tam_report_id);
}

sai_status_t l_set_tam_report_attribute(sai_object_id_t tam_report_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM_REPORT, tam_report_id, attr);
}

sai_status_t l_get_tam_report_attribute(sai_object_id_t tam_report_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM_REPORT, tam_report_id, attr_count, attr_list);
}

sai_status_t l_create_tam_event_threshold(sai_object_id_t *tam_event_threshold_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM_EVENT_THRESHOLD, tam_event_threshold_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_event_threshold(sai_object_id_t tam_event_threshold_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM_EVENT_THRESHOLD, tam_event_threshold_id);
}

sai_status_t l_set_tam_event_threshold_attribute(sai_object_id_t tam_event_threshold_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM_EVENT_THRESHOLD, tam_event_threshold_id, attr);
}

sai_status_t l_get_tam_event_threshold_attribute(sai_object_id_t tam_event_threshold_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM_EVENT_THRESHOLD, tam_event_threshold_id, attr_count, attr_list);
}

sai_status_t l_create_tam_int(sai_object_id_t *tam_int_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM_INT, tam_int_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_int(sai_object_id_t tam_int_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM_INT, tam_int_id);
}

sai_status_t l_set_tam_int_attribute(sai_object_id_t tam_int_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM_INT, tam_int_id, attr);
}

sai_status_t l_get_tam_int_attribute(sai_object_id_t tam_int_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM_INT, tam_int_id, attr_count, attr_list);
}

sai_status_t l_create_tam_tel_type(sai_object_id_t *tam_tel_type_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM_TEL_TYPE, tam_tel_type_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_tel_type(sai_object_id_t tam_tel_type_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM_TEL_TYPE, tam_tel_type_id);
}

sai_status_t l_set_tam_tel_type_attribute(sai_object_id_t tam_tel_type_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM_TEL_TYPE, tam_tel_type_id, attr);
}

sai_status_t l_get_tam_tel_type_attribute(sai_object_id_t tam_tel_type_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM_TEL_TYPE, tam_tel_type_id, attr_count, attr_list);
}

sai_status_t l_create_tam_transport(sai_object_id_t *tam_transport_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM_TRANSPORT, tam_transport_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_transport(sai_object_id_t tam_transport_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM_TRANSPORT, tam_transport_id);
}

sai_status_t l_set_tam_transport_attribute(sai_object_id_t tam_transport_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM_TRANSPORT, tam_transport_id, attr);
}

sai_status_t l_get_tam_transport_attribute(sai_object_id_t tam_transport_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM_TRANSPORT, tam_transport_id, attr_count, attr_list);
}

sai_status_t l_create_tam_telemetry(sai_object_id_t *tam_telemetry_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM_TELEMETRY, tam_telemetry_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_telemetry(sai_object_id_t tam_telemetry_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM_TELEMETRY, tam_telemetry_id);
}

sai_status_t l_set_tam_telemetry_attribute(sai_object_id_t tam_telemetry_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM_TELEMETRY, tam_telemetry_id, attr);
}

sai_status_t l_get_tam_telemetry_attribute(sai_object_id_t tam_telemetry_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM_TELEMETRY, tam_telemetry_id, attr_count, attr_list);
}

sai_status_t l_create_tam_collector(sai_object_id_t *tam_collector_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM_COLLECTOR, tam_collector_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_collector(sai_object_id_t tam_collector_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM_COLLECTOR, tam_collector_id);
}

sai_status_t l_set_tam_collector_attribute(sai_object_id_t tam_collector_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM_COLLECTOR, tam_collector_id, attr);
}

sai_status_t l_get_tam_collector_attribute(sai_object_id_t tam_collector_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM_COLLECTOR, tam_collector_id, attr_count, attr_list);
}

sai_status_t l_create_tam_event_action(sai_object_id_t *tam_event_action_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM_EVENT_ACTION, tam_event_action_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_event_action(sai_object_id_t tam_event_action_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM_EVENT_ACTION, tam_event_action_id);
}

sai_status_t l_set_tam_event_action_attribute(sai_object_id_t tam_event_action_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM_EVENT_ACTION, tam_event_action_id, attr);
}

sai_status_t l_get_tam_event_action_attribute(sai_object_id_t tam_event_action_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM_EVENT_ACTION, tam_event_action_id, attr_count, attr_list);
}

sai_status_t l_create_tam_event(sai_object_id_t *tam_event_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->create(SAI_OBJECT_TYPE_TAM_EVENT, tam_event_id, switch_id, attr_count, attr_list);
}

sai_status_t l_remove_tam_event(sai_object_id_t tam_event_id) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->remove(SAI_OBJECT_TYPE_TAM_EVENT, tam_event_id);
}

sai_status_t l_set_tam_event_attribute(sai_object_id_t tam_event_id, const sai_attribute_t *attr) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->set_attribute(SAI_OBJECT_TYPE_TAM_EVENT, tam_event_id, attr);
}

sai_status_t l_get_tam_event_attribute(sai_object_id_t tam_event_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;
	return translator->get_attribute(SAI_OBJECT_TYPE_TAM_EVENT, tam_event_id, attr_count, attr_list);
}

