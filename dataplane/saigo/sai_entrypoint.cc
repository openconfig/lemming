#include "translator.h"

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

