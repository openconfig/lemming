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

#include "dataplane/standalone/saiserver/tam.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/tam.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Tam::CreateTam(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamRequest* req,
    lemming::dataplane::sai::CreateTamResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTam(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamRequest* req,
    lemming::dataplane::sai::RemoveTamResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam(req.get_oid());

  auto status = api->remove_tam(entry);
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

grpc::Status Tam::SetTamAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamAttributeRequest* req,
    lemming::dataplane::sai::SetTamAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamAttributeRequest* req,
    lemming::dataplane::sai::GetTamAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::CreateTamMathFunc(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamMathFuncRequest* req,
    lemming::dataplane::sai::CreateTamMathFuncResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTamMathFunc(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamMathFuncRequest* req,
    lemming::dataplane::sai::RemoveTamMathFuncResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam_math_func(req.get_oid());

  auto status = api->remove_tam_math_func(entry);
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

grpc::Status Tam::SetTamMathFuncAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamMathFuncAttributeRequest* req,
    lemming::dataplane::sai::SetTamMathFuncAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamMathFuncAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamMathFuncAttributeRequest* req,
    lemming::dataplane::sai::GetTamMathFuncAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::CreateTamReport(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamReportRequest* req,
    lemming::dataplane::sai::CreateTamReportResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTamReport(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamReportRequest* req,
    lemming::dataplane::sai::RemoveTamReportResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam_report(req.get_oid());

  auto status = api->remove_tam_report(entry);
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

grpc::Status Tam::SetTamReportAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamReportAttributeRequest* req,
    lemming::dataplane::sai::SetTamReportAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamReportAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamReportAttributeRequest* req,
    lemming::dataplane::sai::GetTamReportAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::CreateTamEventThreshold(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamEventThresholdRequest* req,
    lemming::dataplane::sai::CreateTamEventThresholdResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTamEventThreshold(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamEventThresholdRequest* req,
    lemming::dataplane::sai::RemoveTamEventThresholdResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam_event_threshold(req.get_oid());

  auto status = api->remove_tam_event_threshold(entry);
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

grpc::Status Tam::SetTamEventThresholdAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamEventThresholdAttributeRequest* req,
    lemming::dataplane::sai::SetTamEventThresholdAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamEventThresholdAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamEventThresholdAttributeRequest* req,
    lemming::dataplane::sai::GetTamEventThresholdAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::CreateTamInt(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamIntRequest* req,
    lemming::dataplane::sai::CreateTamIntResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTamInt(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamIntRequest* req,
    lemming::dataplane::sai::RemoveTamIntResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam_int(req.get_oid());

  auto status = api->remove_tam_int(entry);
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

grpc::Status Tam::SetTamIntAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamIntAttributeRequest* req,
    lemming::dataplane::sai::SetTamIntAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamIntAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamIntAttributeRequest* req,
    lemming::dataplane::sai::GetTamIntAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::CreateTamTelType(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamTelTypeRequest* req,
    lemming::dataplane::sai::CreateTamTelTypeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTamTelType(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamTelTypeRequest* req,
    lemming::dataplane::sai::RemoveTamTelTypeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam_tel_type(req.get_oid());

  auto status = api->remove_tam_tel_type(entry);
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

grpc::Status Tam::SetTamTelTypeAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamTelTypeAttributeRequest* req,
    lemming::dataplane::sai::SetTamTelTypeAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamTelTypeAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamTelTypeAttributeRequest* req,
    lemming::dataplane::sai::GetTamTelTypeAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::CreateTamTransport(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamTransportRequest* req,
    lemming::dataplane::sai::CreateTamTransportResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTamTransport(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamTransportRequest* req,
    lemming::dataplane::sai::RemoveTamTransportResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam_transport(req.get_oid());

  auto status = api->remove_tam_transport(entry);
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

grpc::Status Tam::SetTamTransportAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamTransportAttributeRequest* req,
    lemming::dataplane::sai::SetTamTransportAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamTransportAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamTransportAttributeRequest* req,
    lemming::dataplane::sai::GetTamTransportAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::CreateTamTelemetry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamTelemetryRequest* req,
    lemming::dataplane::sai::CreateTamTelemetryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTamTelemetry(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamTelemetryRequest* req,
    lemming::dataplane::sai::RemoveTamTelemetryResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam_telemetry(req.get_oid());

  auto status = api->remove_tam_telemetry(entry);
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

grpc::Status Tam::SetTamTelemetryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamTelemetryAttributeRequest* req,
    lemming::dataplane::sai::SetTamTelemetryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamTelemetryAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamTelemetryAttributeRequest* req,
    lemming::dataplane::sai::GetTamTelemetryAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::CreateTamCollector(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamCollectorRequest* req,
    lemming::dataplane::sai::CreateTamCollectorResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTamCollector(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamCollectorRequest* req,
    lemming::dataplane::sai::RemoveTamCollectorResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam_collector(req.get_oid());

  auto status = api->remove_tam_collector(entry);
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

grpc::Status Tam::SetTamCollectorAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamCollectorAttributeRequest* req,
    lemming::dataplane::sai::SetTamCollectorAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamCollectorAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamCollectorAttributeRequest* req,
    lemming::dataplane::sai::GetTamCollectorAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::CreateTamEventAction(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamEventActionRequest* req,
    lemming::dataplane::sai::CreateTamEventActionResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTamEventAction(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamEventActionRequest* req,
    lemming::dataplane::sai::RemoveTamEventActionResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam_event_action(req.get_oid());

  auto status = api->remove_tam_event_action(entry);
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

grpc::Status Tam::SetTamEventActionAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamEventActionAttributeRequest* req,
    lemming::dataplane::sai::SetTamEventActionAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamEventActionAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamEventActionAttributeRequest* req,
    lemming::dataplane::sai::GetTamEventActionAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::CreateTamEvent(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateTamEventRequest* req,
    lemming::dataplane::sai::CreateTamEventResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::RemoveTamEvent(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveTamEventRequest* req,
    lemming::dataplane::sai::RemoveTamEventResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  grpc::ClientContext context;
  auto status = api->remove_tam_event(req.get_oid());

  auto status = api->remove_tam_event(entry);
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

grpc::Status Tam::SetTamEventAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetTamEventAttributeRequest* req,
    lemming::dataplane::sai::SetTamEventAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}

grpc::Status Tam::GetTamEventAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetTamEventAttributeRequest* req,
    lemming::dataplane::sai::GetTamEventAttributeResponse* resp) {
  LOG(INFO) << "Func: " << __PRETTY_FUNCTION__;

  return grpc::Status::OK;
}
