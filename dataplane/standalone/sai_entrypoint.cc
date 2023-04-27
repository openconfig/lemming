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

#ifdef __cplusplus
extern "C" {
#endif

// TODO: implement this without using gRPC.
sai_status_t sai_api_initialize(_In_ uint64_t flags, _In_ const sai_service_method_table_t *services)
{
    initialize(GoInt(50000));
    auto chan = grpc::CreateChannel("localhost:50000",grpc::InsecureChannelCredentials());
    Translator::client = forwarding::Forwarding::NewStub(chan);
    return SAI_STATUS_SUCCESS;
}

sai_status_t sai_api_query(_In_ sai_api_t api, _Out_ void **api_method_table)
{
    switch (api)
    {
    case SAI_API_SWITCH:
    {
        sai_switch_api_t *swapi = (sai_switch_api_t *)malloc(sizeof(sai_switch_api_t));
        swapi->create_switch = Translator::create_switch;
        *api_method_table = swapi;
        break;
    }
    default:
        return SAI_STATUS_FAILURE;
    }
    return SAI_STATUS_SUCCESS;
}

sai_status_t sai_log_set(_In_ sai_api_t api, _In_ sai_log_level_t log_level)
{
    return SAI_STATUS_SUCCESS;
}

#ifdef __cplusplus
}
#endif