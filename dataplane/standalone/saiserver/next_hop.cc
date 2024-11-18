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

#include "dataplane/standalone/saiserver/next_hop.h"

#include <glog/logging.h>

#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/proto/sai/next_hop.pb.h"
#include "dataplane/standalone/saiserver/common.h"
#include "dataplane/standalone/saiserver/enum.h"

grpc::Status NextHop::CreateNextHop(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNextHopRequest* req,
    lemming::dataplane::sai::CreateNextHopResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status NextHop::RemoveNextHop(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNextHopRequest* req,
    lemming::dataplane::sai::RemoveNextHopResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status NextHop::SetNextHopAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::SetNextHopAttributeRequest* req,
    lemming::dataplane::sai::SetNextHopAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status NextHop::GetNextHopAttribute(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::GetNextHopAttributeRequest* req,
    lemming::dataplane::sai::GetNextHopAttributeResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status NextHop::CreateNextHops(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::CreateNextHopsRequest* req,
    lemming::dataplane::sai::CreateNextHopsResponse* resp) {
  return grpc::Status::OK;
}

grpc::Status NextHop::RemoveNextHops(
    grpc::ServerContext* context,
    const lemming::dataplane::sai::RemoveNextHopsRequest* req,
    lemming::dataplane::sai::RemoveNextHopsResponse* resp) {
  return grpc::Status::OK;
}
