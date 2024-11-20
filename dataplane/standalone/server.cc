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

#include <grpc/grpc.h>
#include <grpcpp/server.h>
#include <grpcpp/server_builder.h>
#include <grpcpp/server_context.h>

#include "dataplane/proto/sai/common.grpc.pb.h"
#include "dataplane/proto/sai/common.pb.h"
#include "dataplane/standalone/saiserver/acl.h"
#include "dataplane/standalone/saiserver/bfd.h"
#include "dataplane/standalone/saiserver/bmtor.h"
#include "dataplane/standalone/saiserver/bridge.h"
#include "dataplane/standalone/saiserver/buffer.h"
#include "dataplane/standalone/saiserver/counter.h"
#include "dataplane/standalone/saiserver/debug_counter.h"
#include "dataplane/standalone/saiserver/dtel.h"
#include "dataplane/standalone/saiserver/fdb.h"
#include "dataplane/standalone/saiserver/generic_programmable.h"
#include "dataplane/standalone/saiserver/hash.h"
#include "dataplane/standalone/saiserver/hostif.h"
#include "dataplane/standalone/saiserver/ipmc.h"
#include "dataplane/standalone/saiserver/ipmc_group.h"
#include "dataplane/standalone/saiserver/ipsec.h"
#include "dataplane/standalone/saiserver/isolation_group.h"
#include "dataplane/standalone/saiserver/l2mc.h"
#include "dataplane/standalone/saiserver/l2mc_group.h"
#include "dataplane/standalone/saiserver/lag.h"
#include "dataplane/standalone/saiserver/macsec.h"
#include "dataplane/standalone/saiserver/mcast_fdb.h"
#include "dataplane/standalone/saiserver/mirror.h"
#include "dataplane/standalone/saiserver/mpls.h"
#include "dataplane/standalone/saiserver/my_mac.h"
#include "dataplane/standalone/saiserver/nat.h"
#include "dataplane/standalone/saiserver/neighbor.h"
#include "dataplane/standalone/saiserver/next_hop.h"
#include "dataplane/standalone/saiserver/next_hop_group.h"
#include "dataplane/standalone/saiserver/policer.h"
#include "dataplane/standalone/saiserver/port.h"
#include "dataplane/standalone/saiserver/qos_map.h"
#include "dataplane/standalone/saiserver/queue.h"
#include "dataplane/standalone/saiserver/route.h"
#include "dataplane/standalone/saiserver/router_interface.h"
#include "dataplane/standalone/saiserver/rpf_group.h"
#include "dataplane/standalone/saiserver/samplepacket.h"
#include "dataplane/standalone/saiserver/scheduler.h"
#include "dataplane/standalone/saiserver/scheduler_group.h"
#include "dataplane/standalone/saiserver/srv6.h"
#include "dataplane/standalone/saiserver/stp.h"
#include "dataplane/standalone/saiserver/switch.h"
#include "dataplane/standalone/saiserver/system_port.h"
#include "dataplane/standalone/saiserver/tam.h"
#include "dataplane/standalone/saiserver/tunnel.h"
#include "dataplane/standalone/saiserver/udf.h"
#include "dataplane/standalone/saiserver/virtual_router.h"
#include "dataplane/standalone/saiserver/vlan.h"
#include "dataplane/standalone/saiserver/wred.h"

extern "C" {
#include "inc/sai.h"
}

extern "C" {
#include "experimental/saiextensions.h"
}

std::shared_ptr<Acl> acl;
std::shared_ptr<Bfd> bfd;
std::shared_ptr<Buffer> buffer;
std::shared_ptr<Bmtor> bmtor;
std::shared_ptr<Bridge> bridge;
std::shared_ptr<Counter> counter;
std::shared_ptr<DebugCounter> debug_counter;
std::shared_ptr<GenericProgrammable> generic_programmable;
std::shared_ptr<Dtel> dtel;
std::shared_ptr<Fdb> fdb;
std::shared_ptr<Hash> hash;
std::shared_ptr<Hostif> hostif;
std::shared_ptr<IpmcGroup> ipmc_group;
std::shared_ptr<Ipmc> ipmc;
std::shared_ptr<Ipsec> ipsec;
std::shared_ptr<IsolationGroup> isolation_group;
std::shared_ptr<L2mcGroup> l2mc_group;
std::shared_ptr<L2mc> l2mc;
std::shared_ptr<Lag> lag;
std::shared_ptr<Macsec> macsec;
std::shared_ptr<Mirror> mirror;
std::shared_ptr<McastFdb> mcast_fdb;
std::shared_ptr<Mpls> mpls;
std::shared_ptr<MyMac> my_mac;
std::shared_ptr<Nat> nat;
std::shared_ptr<Neighbor> neighbor;
std::shared_ptr<NextHopGroup> next_hop_group;
std::shared_ptr<NextHop> next_hop;
std::shared_ptr<Policer> policer;
std::shared_ptr<Port> port;
std::shared_ptr<QosMap> qos_map;
std::shared_ptr<Queue> queue;
std::shared_ptr<Route> route;
std::shared_ptr<RouterInterface> router_interface;
std::shared_ptr<RpfGroup> rpf_group;
std::shared_ptr<Samplepacket> samplepacket;
std::shared_ptr<SchedulerGroup> scheduler_group;
std::shared_ptr<Scheduler> scheduler;
std::shared_ptr<Srv6> srv6;
std::shared_ptr<Stp> stp;
std::shared_ptr<Switch> switch_;
std::shared_ptr<SystemPort> system_port;
std::shared_ptr<Tam> tam;
std::shared_ptr<Tunnel> tunnel;
std::shared_ptr<Udf> udf;
std::shared_ptr<VirtualRouter> virtual_router;
std::shared_ptr<Vlan> vlan;
std::shared_ptr<Wred> wred;

const char *table_get_value(_In_ sai_switch_profile_id_t profile_id,
                            _In_ const char *variable) {
  return NULL;
}
int table_get_next_value(_In_ sai_switch_profile_id_t profile_id,
                         _Out_ const char **variable,
                         _Out_ const char **value) {
  return -1;
}
const sai_service_method_table_t table = {table_get_value,
                                          table_get_next_value};

class Entry final : public lemming::dataplane::sai::Entrypoint::Service {
  grpc::Status ObjectTypeQuery(
      ::grpc::ServerContext *context,
      const ::lemming::dataplane::sai::ObjectTypeQueryRequest *request,
      ::lemming::dataplane::sai::ObjectTypeQueryResponse *response) {
    return grpc::Status::OK;
  }
  grpc::Status Initialize(
      ::grpc::ServerContext *context,
      const ::lemming::dataplane::sai::InitializeRequest *request,
      ::lemming::dataplane::sai::InitializeResponse *response) {
    sai_status_t st = sai_api_initialize(0, &table);
    if (st != SAI_STATUS_SUCCESS) {
      return grpc::Status(grpc::StatusCode::INTERNAL,
                          "failed to initialize api");
    }

    sai_api_query(SAI_API_ACL, (void **)acl->api);
    sai_api_query(SAI_API_BFD, (void **)bfd->api);
    sai_api_query(SAI_API_BUFFER, (void **)buffer->api);
    sai_api_query((sai_api_t)SAI_API_BMTOR, (void **)bmtor->api);
    sai_api_query(SAI_API_BRIDGE, (void **)bridge->api);
    sai_api_query(SAI_API_COUNTER, (void **)counter->api);
    sai_api_query(SAI_API_DEBUG_COUNTER, (void **)debug_counter->api);
    sai_api_query(SAI_API_GENERIC_PROGRAMMABLE,
                  (void **)generic_programmable->api);
    sai_api_query(SAI_API_DTEL, (void **)dtel->api);
    sai_api_query(SAI_API_FDB, (void **)fdb->api);
    sai_api_query(SAI_API_HASH, (void **)hash->api);
    sai_api_query(SAI_API_HOSTIF, (void **)hostif->api);
    sai_api_query(SAI_API_IPMC_GROUP, (void **)ipmc_group->api);
    sai_api_query(SAI_API_IPMC, (void **)ipmc->api);
    sai_api_query(SAI_API_IPSEC, (void **)ipsec->api);
    sai_api_query(SAI_API_ISOLATION_GROUP, (void **)isolation_group->api);
    sai_api_query(SAI_API_L2MC_GROUP, (void **)l2mc_group->api);
    sai_api_query(SAI_API_L2MC, (void **)l2mc->api);
    sai_api_query(SAI_API_LAG, (void **)lag->api);
    sai_api_query(SAI_API_MACSEC, (void **)macsec->api);
    sai_api_query(SAI_API_MIRROR, (void **)mirror->api);
    sai_api_query(SAI_API_MCAST_FDB, (void **)mcast_fdb->api);
    sai_api_query(SAI_API_MPLS, (void **)mpls->api);
    sai_api_query(SAI_API_MY_MAC, (void **)my_mac->api);
    sai_api_query(SAI_API_NAT, (void **)nat->api);
    sai_api_query(SAI_API_NEIGHBOR, (void **)neighbor->api);
    sai_api_query(SAI_API_NEXT_HOP_GROUP, (void **)next_hop_group->api);
    sai_api_query(SAI_API_NEXT_HOP, (void **)next_hop->api);
    sai_api_query(SAI_API_POLICER, (void **)policer->api);
    sai_api_query(SAI_API_PORT, (void **)port->api);
    sai_api_query(SAI_API_QOS_MAP, (void **)qos_map->api);
    sai_api_query(SAI_API_QUEUE, (void **)queue->api);
    sai_api_query(SAI_API_ROUTE, (void **)route->api);
    sai_api_query(SAI_API_ROUTER_INTERFACE, (void **)router_interface->api);
    sai_api_query(SAI_API_RPF_GROUP, (void **)rpf_group->api);
    sai_api_query(SAI_API_SAMPLEPACKET, (void **)samplepacket->api);
    sai_api_query(SAI_API_SCHEDULER_GROUP, (void **)scheduler_group->api);
    sai_api_query(SAI_API_SCHEDULER, (void **)scheduler->api);
    sai_api_query(SAI_API_SRV6, (void **)srv6->api);
    sai_api_query(SAI_API_STP, (void **)stp->api);
    sai_api_query(SAI_API_SWITCH, (void **)switch_->api);
    sai_api_query(SAI_API_SYSTEM_PORT, (void **)system_port->api);
    sai_api_query(SAI_API_TAM, (void **)tam->api);
    sai_api_query(SAI_API_TUNNEL, (void **)tunnel->api);
    sai_api_query(SAI_API_UDF, (void **)udf->api);
    sai_api_query(SAI_API_VIRTUAL_ROUTER, (void **)virtual_router->api);
    sai_api_query(SAI_API_VLAN, (void **)vlan->api);
    sai_api_query(SAI_API_WRED, (void **)wred->api);

    return grpc::Status::OK;
  }
  grpc::Status Uninitialize(
      ::grpc::ServerContext *context,
      const ::lemming::dataplane::sai::UninitializeRequest *request,
      ::lemming::dataplane::sai::UninitializeResponse *response) {
    return grpc::Status::OK;
  }
};

int main() {
  std::string server_address("0.0.0.0:50000");

  acl = std::make_shared<Acl>();
  bfd = std::make_shared<Bfd>();
  buffer = std::make_shared<Buffer>();
  bmtor = std::make_shared<Bmtor>();
  bridge = std::make_shared<Bridge>();
  counter = std::make_shared<Counter>();
  debug_counter = std::make_shared<DebugCounter>();
  generic_programmable = std::make_shared<GenericProgrammable>();
  dtel = std::make_shared<Dtel>();
  fdb = std::make_shared<Fdb>();
  hash = std::make_shared<Hash>();
  hostif = std::make_shared<Hostif>();
  ipmc_group = std::make_shared<IpmcGroup>();
  ipmc = std::make_shared<Ipmc>();
  ipsec = std::make_shared<Ipsec>();
  isolation_group = std::make_shared<IsolationGroup>();
  l2mc_group = std::make_shared<L2mcGroup>();
  l2mc = std::make_shared<L2mc>();
  lag = std::make_shared<Lag>();
  macsec = std::make_shared<Macsec>();
  mirror = std::make_shared<Mirror>();
  mcast_fdb = std::make_shared<McastFdb>();
  mpls = std::make_shared<Mpls>();
  my_mac = std::make_shared<MyMac>();
  nat = std::make_shared<Nat>();
  neighbor = std::make_shared<Neighbor>();
  next_hop_group = std::make_shared<NextHopGroup>();
  next_hop = std::make_shared<NextHop>();
  policer = std::make_shared<Policer>();
  port = std::make_shared<Port>();
  qos_map = std::make_shared<QosMap>();
  queue = std::make_shared<Queue>();
  route = std::make_shared<Route>();
  router_interface = std::make_shared<RouterInterface>();
  rpf_group = std::make_shared<RpfGroup>();
  samplepacket = std::make_shared<Samplepacket>();
  scheduler_group = std::make_shared<SchedulerGroup>();
  scheduler = std::make_shared<Scheduler>();
  srv6 = std::make_shared<Srv6>();
  stp = std::make_shared<Stp>();
  switch_ = std::make_shared<Switch>();
  system_port = std::make_shared<SystemPort>();
  tam = std::make_shared<Tam>();
  tunnel = std::make_shared<Tunnel>();
  udf = std::make_shared<Udf>();
  virtual_router = std::make_shared<VirtualRouter>();
  vlan = std::make_shared<Vlan>();
  wred = std::make_shared<Wred>();

  grpc::ServerBuilder builder;

  Entry entry;

  builder.RegisterService(&entry);
  builder.RegisterService(acl.get());
  builder.RegisterService(bfd.get());
  builder.RegisterService(buffer.get());
  builder.RegisterService(bmtor.get());
  builder.RegisterService(bridge.get());
  builder.RegisterService(counter.get());
  builder.RegisterService(debug_counter.get());
  builder.RegisterService(generic_programmable.get());
  builder.RegisterService(dtel.get());
  builder.RegisterService(fdb.get());
  builder.RegisterService(hash.get());
  builder.RegisterService(hostif.get());
  builder.RegisterService(ipmc_group.get());
  builder.RegisterService(ipmc.get());
  builder.RegisterService(ipsec.get());
  builder.RegisterService(isolation_group.get());
  builder.RegisterService(l2mc_group.get());
  builder.RegisterService(l2mc.get());
  builder.RegisterService(lag.get());
  builder.RegisterService(macsec.get());
  builder.RegisterService(mirror.get());
  builder.RegisterService(mcast_fdb.get());
  builder.RegisterService(mpls.get());
  builder.RegisterService(my_mac.get());
  builder.RegisterService(nat.get());
  builder.RegisterService(neighbor.get());
  builder.RegisterService(next_hop_group.get());
  builder.RegisterService(next_hop.get());
  builder.RegisterService(policer.get());
  builder.RegisterService(port.get());
  builder.RegisterService(qos_map.get());
  builder.RegisterService(queue.get());
  builder.RegisterService(route.get());
  builder.RegisterService(router_interface.get());
  builder.RegisterService(rpf_group.get());
  builder.RegisterService(samplepacket.get());
  builder.RegisterService(scheduler_group.get());
  builder.RegisterService(scheduler.get());
  builder.RegisterService(srv6.get());
  builder.RegisterService(stp.get());
  builder.RegisterService(switch_.get());
  builder.RegisterService(system_port.get());
  builder.RegisterService(tam.get());
  builder.RegisterService(tunnel.get());
  builder.RegisterService(udf.get());
  builder.RegisterService(virtual_router.get());
  builder.RegisterService(vlan.get());
  builder.RegisterService(wred.get());

  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());

  std::shared_ptr<grpc::Server> server(builder.BuildAndStart());

  server->Wait();
  return 0;
}