// Copyright 2024 Google LLC
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

#include "dataplane/standalone/saiserver/hostif.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/hostif.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Hostif::CreateHostif(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHostifRequest* req,
    lemming::dataplane::sai::CreateHostifResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::RemoveHostif(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHostifRequest* req,
    lemming::dataplane::sai::RemoveHostifResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_hostif(req.get_oid());

  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return grpc::Status::INTERNAL;
  }

  return grpc::Status::OK;
}

grpc::Status Hostif::SetHostifAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetHostifAttributeRequest* req,
    lemming::dataplane::sai::SetHostifAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::GetHostifAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHostifAttributeRequest* req,
    lemming::dataplane::sai::GetHostifAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::CreateHostifTableEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHostifTableEntryRequest* req,
    lemming::dataplane::sai::CreateHostifTableEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::RemoveHostifTableEntry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHostifTableEntryRequest* req,
    lemming::dataplane::sai::RemoveHostifTableEntryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_hostif_table_entry(req.get_oid());

  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return grpc::Status::INTERNAL;
  }

  return grpc::Status::OK;
}

grpc::Status Hostif::GetHostifTableEntryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHostifTableEntryAttributeRequest* req,
    lemming::dataplane::sai::GetHostifTableEntryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::CreateHostifTrapGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHostifTrapGroupRequest* req,
    lemming::dataplane::sai::CreateHostifTrapGroupResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::RemoveHostifTrapGroup(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHostifTrapGroupRequest* req,
    lemming::dataplane::sai::RemoveHostifTrapGroupResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_hostif_trap_group(req.get_oid());

  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return grpc::Status::INTERNAL;
  }

  return grpc::Status::OK;
}

grpc::Status Hostif::SetHostifTrapGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetHostifTrapGroupAttributeRequest* req,
    lemming::dataplane::sai::SetHostifTrapGroupAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::GetHostifTrapGroupAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHostifTrapGroupAttributeRequest* req,
    lemming::dataplane::sai::GetHostifTrapGroupAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::CreateHostifTrap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHostifTrapRequest* req,
    lemming::dataplane::sai::CreateHostifTrapResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::RemoveHostifTrap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHostifTrapRequest* req,
    lemming::dataplane::sai::RemoveHostifTrapResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_hostif_trap(req.get_oid());

  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return grpc::Status::INTERNAL;
  }

  return grpc::Status::OK;
}

grpc::Status Hostif::SetHostifTrapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetHostifTrapAttributeRequest* req,
    lemming::dataplane::sai::SetHostifTrapAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::GetHostifTrapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHostifTrapAttributeRequest* req,
    lemming::dataplane::sai::GetHostifTrapAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::CreateHostifUserDefinedTrap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateHostifUserDefinedTrapRequest* req,
    lemming::dataplane::sai::CreateHostifUserDefinedTrapResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::RemoveHostifUserDefinedTrap(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveHostifUserDefinedTrapRequest* req,
    lemming::dataplane::sai::RemoveHostifUserDefinedTrapResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_hostif_user_defined_trap(req.get_oid());

  if (!status.ok()) {
    auto it = context.GetServerTrailingMetadata().find("traceparent");
    if (it != context.GetServerTrailingMetadata().end()) {
      LOG(ERROR) << "Lucius RPC error: Trace ID " << it->second
                 << " msg: " << status.error_message();
    } else {
      LOG(ERROR) << "Lucius RPC error: " << status.error_message();
    }
    return grpc::Status::INTERNAL;
  }

  return grpc::Status::OK;
}

grpc::Status Hostif::SetHostifUserDefinedTrapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetHostifUserDefinedTrapAttributeRequest*
        req,
    lemming::dataplane::sai::SetHostifUserDefinedTrapAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Hostif::GetHostifUserDefinedTrapAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetHostifUserDefinedTrapAttributeRequest*
        req,
    lemming::dataplane::sai::GetHostifUserDefinedTrapAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
