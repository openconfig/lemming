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

#ifndef DATAPLANE_STANDALONE_SAI_TAM_H_
#define DATAPLANE_STANDALONE_SAI_TAM_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/tam.grpc.pb.h"
#include "dataplane/proto/sai/tam.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Tam final : public lemming::dataplane::sai::Tam::Service {
 public:
  grpc::Status CreateTam(grpc::ServerContext* context,
                         const lemming::dataplane::sai::CreateTamRequest* req,
                         lemming::dataplane::sai::CreateTamResponse* resp);

  grpc::Status RemoveTam(grpc::ServerContext* context,
                         const lemming::dataplane::sai::RemoveTamRequest* req,
                         lemming::dataplane::sai::RemoveTamResponse* resp);

  grpc::Status SetTamAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamAttributeRequest* req,
      lemming::dataplane::sai::SetTamAttributeResponse* resp);

  grpc::Status GetTamAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamAttributeRequest* req,
      lemming::dataplane::sai::GetTamAttributeResponse* resp);

  grpc::Status CreateTamMathFunc(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTamMathFuncRequest* req,
      lemming::dataplane::sai::CreateTamMathFuncResponse* resp);

  grpc::Status RemoveTamMathFunc(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTamMathFuncRequest* req,
      lemming::dataplane::sai::RemoveTamMathFuncResponse* resp);

  grpc::Status SetTamMathFuncAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamMathFuncAttributeRequest* req,
      lemming::dataplane::sai::SetTamMathFuncAttributeResponse* resp);

  grpc::Status GetTamMathFuncAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamMathFuncAttributeRequest* req,
      lemming::dataplane::sai::GetTamMathFuncAttributeResponse* resp);

  grpc::Status CreateTamReport(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTamReportRequest* req,
      lemming::dataplane::sai::CreateTamReportResponse* resp);

  grpc::Status RemoveTamReport(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTamReportRequest* req,
      lemming::dataplane::sai::RemoveTamReportResponse* resp);

  grpc::Status SetTamReportAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamReportAttributeRequest* req,
      lemming::dataplane::sai::SetTamReportAttributeResponse* resp);

  grpc::Status GetTamReportAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamReportAttributeRequest* req,
      lemming::dataplane::sai::GetTamReportAttributeResponse* resp);

  grpc::Status CreateTamEventThreshold(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTamEventThresholdRequest* req,
      lemming::dataplane::sai::CreateTamEventThresholdResponse* resp);

  grpc::Status RemoveTamEventThreshold(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTamEventThresholdRequest* req,
      lemming::dataplane::sai::RemoveTamEventThresholdResponse* resp);

  grpc::Status SetTamEventThresholdAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamEventThresholdAttributeRequest* req,
      lemming::dataplane::sai::SetTamEventThresholdAttributeResponse* resp);

  grpc::Status GetTamEventThresholdAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamEventThresholdAttributeRequest* req,
      lemming::dataplane::sai::GetTamEventThresholdAttributeResponse* resp);

  grpc::Status CreateTamInt(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTamIntRequest* req,
      lemming::dataplane::sai::CreateTamIntResponse* resp);

  grpc::Status RemoveTamInt(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTamIntRequest* req,
      lemming::dataplane::sai::RemoveTamIntResponse* resp);

  grpc::Status SetTamIntAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamIntAttributeRequest* req,
      lemming::dataplane::sai::SetTamIntAttributeResponse* resp);

  grpc::Status GetTamIntAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamIntAttributeRequest* req,
      lemming::dataplane::sai::GetTamIntAttributeResponse* resp);

  grpc::Status CreateTamTelType(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTamTelTypeRequest* req,
      lemming::dataplane::sai::CreateTamTelTypeResponse* resp);

  grpc::Status RemoveTamTelType(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTamTelTypeRequest* req,
      lemming::dataplane::sai::RemoveTamTelTypeResponse* resp);

  grpc::Status SetTamTelTypeAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamTelTypeAttributeRequest* req,
      lemming::dataplane::sai::SetTamTelTypeAttributeResponse* resp);

  grpc::Status GetTamTelTypeAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamTelTypeAttributeRequest* req,
      lemming::dataplane::sai::GetTamTelTypeAttributeResponse* resp);

  grpc::Status CreateTamTransport(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTamTransportRequest* req,
      lemming::dataplane::sai::CreateTamTransportResponse* resp);

  grpc::Status RemoveTamTransport(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTamTransportRequest* req,
      lemming::dataplane::sai::RemoveTamTransportResponse* resp);

  grpc::Status SetTamTransportAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamTransportAttributeRequest* req,
      lemming::dataplane::sai::SetTamTransportAttributeResponse* resp);

  grpc::Status GetTamTransportAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamTransportAttributeRequest* req,
      lemming::dataplane::sai::GetTamTransportAttributeResponse* resp);

  grpc::Status CreateTamTelemetry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTamTelemetryRequest* req,
      lemming::dataplane::sai::CreateTamTelemetryResponse* resp);

  grpc::Status RemoveTamTelemetry(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTamTelemetryRequest* req,
      lemming::dataplane::sai::RemoveTamTelemetryResponse* resp);

  grpc::Status SetTamTelemetryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamTelemetryAttributeRequest* req,
      lemming::dataplane::sai::SetTamTelemetryAttributeResponse* resp);

  grpc::Status GetTamTelemetryAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamTelemetryAttributeRequest* req,
      lemming::dataplane::sai::GetTamTelemetryAttributeResponse* resp);

  grpc::Status CreateTamCollector(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTamCollectorRequest* req,
      lemming::dataplane::sai::CreateTamCollectorResponse* resp);

  grpc::Status RemoveTamCollector(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTamCollectorRequest* req,
      lemming::dataplane::sai::RemoveTamCollectorResponse* resp);

  grpc::Status SetTamCollectorAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamCollectorAttributeRequest* req,
      lemming::dataplane::sai::SetTamCollectorAttributeResponse* resp);

  grpc::Status GetTamCollectorAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamCollectorAttributeRequest* req,
      lemming::dataplane::sai::GetTamCollectorAttributeResponse* resp);

  grpc::Status CreateTamEventAction(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTamEventActionRequest* req,
      lemming::dataplane::sai::CreateTamEventActionResponse* resp);

  grpc::Status RemoveTamEventAction(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTamEventActionRequest* req,
      lemming::dataplane::sai::RemoveTamEventActionResponse* resp);

  grpc::Status SetTamEventActionAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamEventActionAttributeRequest* req,
      lemming::dataplane::sai::SetTamEventActionAttributeResponse* resp);

  grpc::Status GetTamEventActionAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamEventActionAttributeRequest* req,
      lemming::dataplane::sai::GetTamEventActionAttributeResponse* resp);

  grpc::Status CreateTamEvent(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateTamEventRequest* req,
      lemming::dataplane::sai::CreateTamEventResponse* resp);

  grpc::Status RemoveTamEvent(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveTamEventRequest* req,
      lemming::dataplane::sai::RemoveTamEventResponse* resp);

  grpc::Status SetTamEventAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetTamEventAttributeRequest* req,
      lemming::dataplane::sai::SetTamEventAttributeResponse* resp);

  grpc::Status GetTamEventAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetTamEventAttributeRequest* req,
      lemming::dataplane::sai::GetTamEventAttributeResponse* resp);

  sai_tam_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_TAM_H_
