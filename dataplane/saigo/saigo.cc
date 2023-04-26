#include "saigo.h"

sai_status_t sai_api_initialize(_In_ uint64_t flags, _In_ const sai_service_method_table_t *services)
{
    return SAI_STATUS_SUCCESS;
}

sai_status_t sai_api_query(_In_ sai_api_t api, _Out_ void **api_method_table)
{
    switch (api)
    {
    case SAI_API_SWITCH:
    {
        sai_switch_api_t *swapi = (sai_switch_api_t *)malloc(sizeof(sai_switch_api_t));
        swapi->create_switch = create_switch;
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

sai_status_t create_switch(_Out_ sai_object_id_t *switch_id, _In_ uint32_t attr_count, _In_ const sai_attribute_t *attr_list)
{
    for (int i = 0; i < attr_count; i++)
    {
        switch (attr_list[i].id)
        { // handle union with different ids
        default:
            printf("%d\n", attr_list[i].value);
        }
    }
    // call some go func
    *switch_id = 1;
    return SAI_STATUS_SUCCESS;
    return SAI_STATUS_SUCCESS;
}