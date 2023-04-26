#ifndef TRANSLATOR_H
#define TRANSLATOR_H

#include "proto/forwarding/forwarding_service.pb.h"
#include "proto/forwarding/forwarding_service.grpc.pb.h"
#include "sai.h"
#include <grpc/grpc.h>
#include <grpcpp/channel.h>
#include <grpcpp/client_context.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>
#include "dataplane/saigo/standalone/standalone.h"


namespace Translator {
    sai_status_t create_switch(_Out_ sai_object_id_t *switch_id, _In_ uint32_t attr_count, _In_ const sai_attribute_t *attr_list);
    extern std::unique_ptr<forwarding::Forwarding::Stub> client;        
};

#endif