#include "proto/forwarding/forwarding_service.pb.h"
#include "proto/forwarding/forwarding_service.grpc.pb.h"
#include "sai.h"

sai_status_t create_switch(_Out_ sai_object_id_t *switch_id, _In_ uint32_t attr_count, _In_ const sai_attribute_t *attr_list);

class Transalator {
    private:
        Transalator() {
        };
    private:
        std::unique_ptr<forwarding::ForwardingService::Stub> stub_;
};