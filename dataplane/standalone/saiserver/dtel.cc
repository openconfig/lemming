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

#include "dataplane/standalone/saiserver/dtel.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/dtel.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status Dtel::CreateDtel(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateDtelRequest* req,
    lemming::dataplane::sai::CreateDtelResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::RemoveDtel(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveDtelRequest* req,
    lemming::dataplane::sai::RemoveDtelResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::SetDtelAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetDtelAttributeRequest* req,
    lemming::dataplane::sai::SetDtelAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::GetDtelAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetDtelAttributeRequest* req,
    lemming::dataplane::sai::GetDtelAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::CreateDtelQueueReport(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateDtelQueueReportRequest* req,
    lemming::dataplane::sai::CreateDtelQueueReportResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::RemoveDtelQueueReport(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveDtelQueueReportRequest* req,
    lemming::dataplane::sai::RemoveDtelQueueReportResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::SetDtelQueueReportAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetDtelQueueReportAttributeRequest* req,
    lemming::dataplane::sai::SetDtelQueueReportAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::GetDtelQueueReportAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetDtelQueueReportAttributeRequest* req,
    lemming::dataplane::sai::GetDtelQueueReportAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::CreateDtelIntSession(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateDtelIntSessionRequest* req,
    lemming::dataplane::sai::CreateDtelIntSessionResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::RemoveDtelIntSession(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveDtelIntSessionRequest* req,
    lemming::dataplane::sai::RemoveDtelIntSessionResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::SetDtelIntSessionAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetDtelIntSessionAttributeRequest* req,
    lemming::dataplane::sai::SetDtelIntSessionAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::GetDtelIntSessionAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetDtelIntSessionAttributeRequest* req,
    lemming::dataplane::sai::GetDtelIntSessionAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::CreateDtelReportSession(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateDtelReportSessionRequest* req,
    lemming::dataplane::sai::CreateDtelReportSessionResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::RemoveDtelReportSession(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveDtelReportSessionRequest* req,
    lemming::dataplane::sai::RemoveDtelReportSessionResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::SetDtelReportSessionAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetDtelReportSessionAttributeRequest* req,
    lemming::dataplane::sai::SetDtelReportSessionAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::GetDtelReportSessionAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetDtelReportSessionAttributeRequest* req,
    lemming::dataplane::sai::GetDtelReportSessionAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::CreateDtelEvent(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateDtelEventRequest* req,
    lemming::dataplane::sai::CreateDtelEventResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::RemoveDtelEvent(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveDtelEventRequest* req,
    lemming::dataplane::sai::RemoveDtelEventResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::SetDtelEventAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetDtelEventAttributeRequest* req,
    lemming::dataplane::sai::SetDtelEventAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status Dtel::GetDtelEventAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetDtelEventAttributeRequest* req,
    lemming::dataplane::sai::GetDtelEventAttributeResponse* resp) {
  return grpc::Status::OK;
}
