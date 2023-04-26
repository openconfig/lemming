#include "translator.h"

using namespace Translator;

sai_status_t Translator::create_switch(_Out_ sai_object_id_t *switch_id, _In_ uint32_t attr_count, _In_ const sai_attribute_t *attr_list)
{
    grpc::ClientContext context;
    forwarding::ContextCreateRequest req;
    forwarding::ContextCreateReply reply;

    *switch_id = 1;
    forwarding::ContextId contextID;
    contextID.set_id("1");
    req.set_allocated_context_id(&contextID);

    for (uint32_t i = 0; i < attr_count; i++)
    {
    }

    auto status = client->ContextCreate(&context, req, &reply);
    if (!status.ok())
    {
        std::cout << status.error_code() << ": " << status.error_message()
                  << std::endl;
        return SAI_STATUS_FAILURE;
    }
    return SAI_STATUS_SUCCESS;
}