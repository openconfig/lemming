
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

#include "dataplane/standalone/sai/segmentroute.h"
#include "dataplane/standalone/log/log.h"

const sai_segmentroute_api_t l_segmentroute = {
	.create_segmentroute_sidlist = l_create_segmentroute_sidlist,
	.remove_segmentroute_sidlist = l_remove_segmentroute_sidlist,
	.set_segmentroute_sidlist_attribute = l_set_segmentroute_sidlist_attribute,
	.get_segmentroute_sidlist_attribute = l_get_segmentroute_sidlist_attribute,
	.create_segmentroute_sidlists = l_create_segmentroute_sidlists,
	.remove_segmentroute_sidlists = l_remove_segmentroute_sidlists,
};


sai_status_t l_create_segmentroute_sidlist(sai_object_id_t *segmentroute_sidlist_id, sai_object_id_t switch_id, uint32_t attr_count, const sai_attribute_t *attr_list) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_remove_segmentroute_sidlist(sai_object_id_t segmentroute_sidlist_id) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_set_segmentroute_sidlist_attribute(sai_object_id_t segmentroute_sidlist_id, const sai_attribute_t *attr) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_get_segmentroute_sidlist_attribute(sai_object_id_t segmentroute_sidlist_id, uint32_t attr_count, sai_attribute_t *attr_list) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_create_segmentroute_sidlists(sai_object_id_t switch_id, uint32_t object_count, const uint32_t *attr_count, const sai_attribute_t **attr_list, sai_bulk_op_error_mode_t mode, sai_object_id_t *object_id, sai_status_t *object_statuses) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


sai_status_t l_remove_segmentroute_sidlists(uint32_t object_count, const sai_object_id_t *object_id, sai_bulk_op_error_mode_t mode, sai_status_t *object_statuses) {
	LUCIUS_LOG_FUNC();
	return SAI_STATUS_NOT_IMPLEMENTED;
}


