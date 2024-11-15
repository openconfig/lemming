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

#ifndef DATAPLANE_STANDALONE_SAI_DTEL_H_
#define DATAPLANE_STANDALONE_SAI_DTEL_H_

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/dtel.grpc.pb.h"
#include "dataplane/proto/sai/dtel.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

class Dtel final : public lemming::dataplane::sai::Dtel::Service {
 public:
  grpc::Status CreateDtel(grpc::ServerContext* context,
                          const lemming::dataplane::sai::CreateDtelRequest* req,
                          lemming::dataplane::sai::CreateDtelResponse* resp);

  grpc::Status RemoveDtel(grpc::ServerContext* context,
                          const lemming::dataplane::sai::RemoveDtelRequest* req,
                          lemming::dataplane::sai::RemoveDtelResponse* resp);

  grpc::Status SetDtelAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetDtelAttributeRequest* req,
      lemming::dataplane::sai::SetDtelAttributeResponse* resp);

  grpc::Status GetDtelAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetDtelAttributeRequest* req,
      lemming::dataplane::sai::GetDtelAttributeResponse* resp);

  grpc::Status CreateDtelQueueReport(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateDtelQueueReportRequest* req,
      lemming::dataplane::sai::CreateDtelQueueReportResponse* resp);

  grpc::Status RemoveDtelQueueReport(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveDtelQueueReportRequest* req,
      lemming::dataplane::sai::RemoveDtelQueueReportResponse* resp);

  grpc::Status SetDtelQueueReportAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetDtelQueueReportAttributeRequest* req,
      lemming::dataplane::sai::SetDtelQueueReportAttributeResponse* resp);

  grpc::Status GetDtelQueueReportAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetDtelQueueReportAttributeRequest* req,
      lemming::dataplane::sai::GetDtelQueueReportAttributeResponse* resp);

  grpc::Status CreateDtelIntSession(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateDtelIntSessionRequest* req,
      lemming::dataplane::sai::CreateDtelIntSessionResponse* resp);

  grpc::Status RemoveDtelIntSession(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveDtelIntSessionRequest* req,
      lemming::dataplane::sai::RemoveDtelIntSessionResponse* resp);

  grpc::Status SetDtelIntSessionAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetDtelIntSessionAttributeRequest* req,
      lemming::dataplane::sai::SetDtelIntSessionAttributeResponse* resp);

  grpc::Status GetDtelIntSessionAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetDtelIntSessionAttributeRequest* req,
      lemming::dataplane::sai::GetDtelIntSessionAttributeResponse* resp);

  grpc::Status CreateDtelReportSession(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateDtelReportSessionRequest* req,
      lemming::dataplane::sai::CreateDtelReportSessionResponse* resp);

  grpc::Status RemoveDtelReportSession(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveDtelReportSessionRequest* req,
      lemming::dataplane::sai::RemoveDtelReportSessionResponse* resp);

  grpc::Status SetDtelReportSessionAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetDtelReportSessionAttributeRequest* req,
      lemming::dataplane::sai::SetDtelReportSessionAttributeResponse* resp);

  grpc::Status GetDtelReportSessionAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetDtelReportSessionAttributeRequest* req,
      lemming::dataplane::sai::GetDtelReportSessionAttributeResponse* resp);

  grpc::Status CreateDtelEvent(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::CreateDtelEventRequest* req,
      lemming::dataplane::sai::CreateDtelEventResponse* resp);

  grpc::Status RemoveDtelEvent(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::RemoveDtelEventRequest* req,
      lemming::dataplane::sai::RemoveDtelEventResponse* resp);

  grpc::Status SetDtelEventAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::SetDtelEventAttributeRequest* req,
      lemming::dataplane::sai::SetDtelEventAttributeResponse* resp);

  grpc::Status GetDtelEventAttribute(
      grpc::ServerContext* context,
      const lemming::dataplane::sai::GetDtelEventAttributeRequest* req,
      lemming::dataplane::sai::GetDtelEventAttributeResponse* resp);

  sai_dtel_api_t* api;
};

#endif  // DATAPLANE_STANDALONE_SAI_DTEL_H_
