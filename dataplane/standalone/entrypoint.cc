// Copyright 2023 Google LLC
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

#include <glog/logging.h>
#include <grpc/grpc.h>
#include <grpcpp/channel.h>
#include <grpcpp/client_context.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>

#include <fstream>

#include "dataplane/standalone/lucius/lucius_clib.h"
#include "dataplane/standalone/proto/acl.grpc.pb.h"
#include "dataplane/standalone/proto/bfd.grpc.pb.h"
#include "dataplane/standalone/proto/bmtor.grpc.pb.h"
#include "dataplane/standalone/proto/bridge.grpc.pb.h"
#include "dataplane/standalone/proto/buffer.grpc.pb.h"
#include "dataplane/standalone/proto/common.grpc.pb.h"
#include "dataplane/standalone/proto/common.pb.h"
#include "dataplane/standalone/proto/counter.grpc.pb.h"
#include "dataplane/standalone/proto/debug_counter.grpc.pb.h"
#include "dataplane/standalone/proto/dtel.grpc.pb.h"
#include "dataplane/standalone/proto/fdb.grpc.pb.h"
#include "dataplane/standalone/proto/generic_programmable.grpc.pb.h"
#include "dataplane/standalone/proto/hash.grpc.pb.h"
#include "dataplane/standalone/proto/hostif.grpc.pb.h"
#include "dataplane/standalone/proto/ipmc.grpc.pb.h"
#include "dataplane/standalone/proto/ipmc_group.grpc.pb.h"
#include "dataplane/standalone/proto/ipsec.grpc.pb.h"
#include "dataplane/standalone/proto/isolation_group.grpc.pb.h"
#include "dataplane/standalone/proto/l2mc.grpc.pb.h"
#include "dataplane/standalone/proto/l2mc_group.grpc.pb.h"
#include "dataplane/standalone/proto/lag.grpc.pb.h"
#include "dataplane/standalone/proto/macsec.grpc.pb.h"
#include "dataplane/standalone/proto/mcast_fdb.grpc.pb.h"
#include "dataplane/standalone/proto/mirror.grpc.pb.h"
#include "dataplane/standalone/proto/mpls.grpc.pb.h"
#include "dataplane/standalone/proto/my_mac.grpc.pb.h"
#include "dataplane/standalone/proto/nat.grpc.pb.h"
#include "dataplane/standalone/proto/neighbor.grpc.pb.h"
#include "dataplane/standalone/proto/next_hop.grpc.pb.h"
#include "dataplane/standalone/proto/next_hop_group.grpc.pb.h"
#include "dataplane/standalone/proto/policer.grpc.pb.h"
#include "dataplane/standalone/proto/port.grpc.pb.h"
#include "dataplane/standalone/proto/qos_map.grpc.pb.h"
#include "dataplane/standalone/proto/queue.grpc.pb.h"
#include "dataplane/standalone/proto/route.grpc.pb.h"
#include "dataplane/standalone/proto/router_interface.grpc.pb.h"
#include "dataplane/standalone/proto/rpf_group.grpc.pb.h"
#include "dataplane/standalone/proto/samplepacket.grpc.pb.h"
#include "dataplane/standalone/proto/scheduler.grpc.pb.h"
#include "dataplane/standalone/proto/scheduler_group.grpc.pb.h"
#include "dataplane/standalone/proto/srv6.grpc.pb.h"
#include "dataplane/standalone/proto/stp.grpc.pb.h"
#include "dataplane/standalone/proto/switch.grpc.pb.h"
#include "dataplane/standalone/proto/system_port.grpc.pb.h"
#include "dataplane/standalone/proto/tam.grpc.pb.h"
#include "dataplane/standalone/proto/tunnel.grpc.pb.h"
#include "dataplane/standalone/proto/udf.grpc.pb.h"
#include "dataplane/standalone/proto/virtual_router.grpc.pb.h"
#include "dataplane/standalone/proto/vlan.grpc.pb.h"
#include "dataplane/standalone/proto/wred.grpc.pb.h"
#include "dataplane/standalone/sai/acl.h"
#include "dataplane/standalone/sai/bfd.h"
#include "dataplane/standalone/sai/bmtor.h"
#include "dataplane/standalone/sai/bridge.h"
#include "dataplane/standalone/sai/buffer.h"
#include "dataplane/standalone/sai/counter.h"
#include "dataplane/standalone/sai/debug_counter.h"
#include "dataplane/standalone/sai/dtel.h"
#include "dataplane/standalone/sai/fdb.h"
#include "dataplane/standalone/sai/generic_programmable.h"
#include "dataplane/standalone/sai/hash.h"
#include "dataplane/standalone/sai/hostif.h"
#include "dataplane/standalone/sai/ipmc.h"
#include "dataplane/standalone/sai/ipmc_group.h"
#include "dataplane/standalone/sai/ipsec.h"
#include "dataplane/standalone/sai/isolation_group.h"
#include "dataplane/standalone/sai/l2mc.h"
#include "dataplane/standalone/sai/l2mc_group.h"
#include "dataplane/standalone/sai/lag.h"
#include "dataplane/standalone/sai/macsec.h"
#include "dataplane/standalone/sai/mcast_fdb.h"
#include "dataplane/standalone/sai/mirror.h"
#include "dataplane/standalone/sai/mpls.h"
#include "dataplane/standalone/sai/my_mac.h"
#include "dataplane/standalone/sai/nat.h"
#include "dataplane/standalone/sai/neighbor.h"
#include "dataplane/standalone/sai/next_hop.h"
#include "dataplane/standalone/sai/next_hop_group.h"
#include "dataplane/standalone/sai/policer.h"
#include "dataplane/standalone/sai/port.h"
#include "dataplane/standalone/sai/qos_map.h"
#include "dataplane/standalone/sai/queue.h"
#include "dataplane/standalone/sai/route.h"
#include "dataplane/standalone/sai/router_interface.h"
#include "dataplane/standalone/sai/rpf_group.h"
#include "dataplane/standalone/sai/samplepacket.h"
#include "dataplane/standalone/sai/scheduler.h"
#include "dataplane/standalone/sai/scheduler_group.h"
#include "dataplane/standalone/sai/srv6.h"
#include "dataplane/standalone/sai/stp.h"
#include "dataplane/standalone/sai/switch.h"
#include "dataplane/standalone/sai/system_port.h"
#include "dataplane/standalone/sai/tam.h"
#include "dataplane/standalone/sai/tunnel.h"
#include "dataplane/standalone/sai/udf.h"
#include "dataplane/standalone/sai/virtual_router.h"
#include "dataplane/standalone/sai/vlan.h"
#include "dataplane/standalone/sai/wred.h"

extern "C" {
#include "experimental/saiextensions.h"
#include "inc/sai.h"
}

std::unique_ptr<lemming::dataplane::sai::Acl::Stub> acl;
std::unique_ptr<lemming::dataplane::sai::Bfd::Stub> bfd;
std::unique_ptr<lemming::dataplane::sai::Buffer::Stub> buffer;
std::unique_ptr<lemming::dataplane::sai::Bmtor::Stub> bmtor;
std::unique_ptr<lemming::dataplane::sai::Bridge::Stub> bridge;
std::unique_ptr<lemming::dataplane::sai::Counter::Stub> counter;
std::unique_ptr<lemming::dataplane::sai::DebugCounter::Stub> debug_counter;
std::unique_ptr<lemming::dataplane::sai::GenericProgrammable::Stub>
    generic_programmable;
std::unique_ptr<lemming::dataplane::sai::Dtel::Stub> dtel;
std::unique_ptr<lemming::dataplane::sai::Fdb::Stub> fdb;
std::unique_ptr<lemming::dataplane::sai::Hash::Stub> hash;
std::unique_ptr<lemming::dataplane::sai::Hostif::Stub> hostif;
std::unique_ptr<lemming::dataplane::sai::IpmcGroup::Stub> ipmc_group;
std::unique_ptr<lemming::dataplane::sai::Ipmc::Stub> ipmc;
std::unique_ptr<lemming::dataplane::sai::Ipsec::Stub> ipsec;
std::unique_ptr<lemming::dataplane::sai::IsolationGroup::Stub> isolation_group;
std::unique_ptr<lemming::dataplane::sai::L2mcGroup::Stub> l2mc_group;
std::unique_ptr<lemming::dataplane::sai::L2mc::Stub> l2mc;
std::unique_ptr<lemming::dataplane::sai::Lag::Stub> lag;
std::unique_ptr<lemming::dataplane::sai::Macsec::Stub> macsec;
std::unique_ptr<lemming::dataplane::sai::Mirror::Stub> mirror;
std::unique_ptr<lemming::dataplane::sai::McastFdb::Stub> mcast_fdb;
std::unique_ptr<lemming::dataplane::sai::Mpls::Stub> mpls;
std::unique_ptr<lemming::dataplane::sai::MyMac::Stub> my_mac;
std::unique_ptr<lemming::dataplane::sai::Nat::Stub> nat;
std::unique_ptr<lemming::dataplane::sai::Neighbor::Stub> neighbor;
std::unique_ptr<lemming::dataplane::sai::NextHopGroup::Stub> next_hop_group;
std::unique_ptr<lemming::dataplane::sai::NextHop::Stub> next_hop;
std::unique_ptr<lemming::dataplane::sai::Policer::Stub> policer;
std::unique_ptr<lemming::dataplane::sai::Port::Stub> port;
std::unique_ptr<lemming::dataplane::sai::QosMap::Stub> qos_map;
std::unique_ptr<lemming::dataplane::sai::Queue::Stub> queue;
std::unique_ptr<lemming::dataplane::sai::Route::Stub> route;
std::unique_ptr<lemming::dataplane::sai::RouterInterface::Stub>
    router_interface;
std::unique_ptr<lemming::dataplane::sai::RpfGroup::Stub> rpf_group;
std::unique_ptr<lemming::dataplane::sai::Samplepacket::Stub> samplepacket;
std::unique_ptr<lemming::dataplane::sai::SchedulerGroup::Stub> scheduler_group;
std::unique_ptr<lemming::dataplane::sai::Scheduler::Stub> scheduler;
std::unique_ptr<lemming::dataplane::sai::Srv6::Stub> srv6;
std::unique_ptr<lemming::dataplane::sai::Stp::Stub> stp;
std::unique_ptr<lemming::dataplane::sai::Switch::Stub> switch_;
std::unique_ptr<lemming::dataplane::sai::SystemPort::Stub> system_port;
std::unique_ptr<lemming::dataplane::sai::Tam::Stub> tam;
std::unique_ptr<lemming::dataplane::sai::Tunnel::Stub> tunnel;
std::unique_ptr<lemming::dataplane::sai::Udf::Stub> udf;
std::unique_ptr<lemming::dataplane::sai::VirtualRouter::Stub> virtual_router;
std::unique_ptr<lemming::dataplane::sai::Vlan::Stub> vlan;
std::unique_ptr<lemming::dataplane::sai::Wred::Stub> wred;

std::unique_ptr<lemming::dataplane::sai::Entrypoint::Stub> entry;

// TODO(dgrau): implement this without using gRPC.
sai_status_t sai_api_initialize(
    _In_ uint64_t flags, _In_ const sai_service_method_table_t *services) {
  FLAGS_log_dir = "/var/log/syncd";
  google::InitGoogleLogging("lucius");
  google::InstallFailureSignalHandler();

  LOG(WARNING) << "iniitializing";
  startAsync(50000);

  auto chan = grpc::CreateChannel("localhost:50000",
                                  grpc::InsecureChannelCredentials());

  acl = std::make_unique<lemming::dataplane::sai::Acl::Stub>(chan);
  bfd = std::make_unique<lemming::dataplane::sai::Bfd::Stub>(chan);
  buffer = std::make_unique<lemming::dataplane::sai::Buffer::Stub>(chan);
  bridge = std::make_unique<lemming::dataplane::sai::Bridge::Stub>(chan);
  counter = std::make_unique<lemming::dataplane::sai::Counter::Stub>(chan);
  debug_counter =
      std::make_unique<lemming::dataplane::sai::DebugCounter::Stub>(chan);
  dtel = std::make_unique<lemming::dataplane::sai::Dtel::Stub>(chan);
  generic_programmable =
      std::make_unique<lemming::dataplane::sai::GenericProgrammable::Stub>(
          chan);
  fdb = std::make_unique<lemming::dataplane::sai::Fdb::Stub>(chan);
  hash = std::make_unique<lemming::dataplane::sai::Hash::Stub>(chan);
  hostif = std::make_unique<lemming::dataplane::sai::Hostif::Stub>(chan);
  ipmc_group = std::make_unique<lemming::dataplane::sai::IpmcGroup::Stub>(chan);
  ipmc = std::make_unique<lemming::dataplane::sai::Ipmc::Stub>(chan);
  ipsec = std::make_unique<lemming::dataplane::sai::Ipsec::Stub>(chan);
  isolation_group =
      std::make_unique<lemming::dataplane::sai::IsolationGroup::Stub>(chan);
  l2mc_group = std::make_unique<lemming::dataplane::sai::L2mcGroup::Stub>(chan);
  l2mc = std::make_unique<lemming::dataplane::sai::L2mc::Stub>(chan);
  lag = std::make_unique<lemming::dataplane::sai::Lag::Stub>(chan);
  macsec = std::make_unique<lemming::dataplane::sai::Macsec::Stub>(chan);
  mirror = std::make_unique<lemming::dataplane::sai::Mirror::Stub>(chan);
  mcast_fdb = std::make_unique<lemming::dataplane::sai::McastFdb::Stub>(chan);
  mpls = std::make_unique<lemming::dataplane::sai::Mpls::Stub>(chan);
  my_mac = std::make_unique<lemming::dataplane::sai::MyMac::Stub>(chan);
  nat = std::make_unique<lemming::dataplane::sai::Nat::Stub>(chan);
  neighbor = std::make_unique<lemming::dataplane::sai::Neighbor::Stub>(chan);
  next_hop_group =
      std::make_unique<lemming::dataplane::sai::NextHopGroup::Stub>(chan);
  next_hop = std::make_unique<lemming::dataplane::sai::NextHop::Stub>(chan);
  policer = std::make_unique<lemming::dataplane::sai::Policer::Stub>(chan);
  port = std::make_unique<lemming::dataplane::sai::Port::Stub>(chan);
  qos_map = std::make_unique<lemming::dataplane::sai::QosMap::Stub>(chan);
  queue = std::make_unique<lemming::dataplane::sai::Queue::Stub>(chan);
  route = std::make_unique<lemming::dataplane::sai::Route::Stub>(chan);
  router_interface =
      std::make_unique<lemming::dataplane::sai::RouterInterface::Stub>(chan);
  rpf_group = std::make_unique<lemming::dataplane::sai::RpfGroup::Stub>(chan);
  samplepacket =
      std::make_unique<lemming::dataplane::sai::Samplepacket::Stub>(chan);
  scheduler_group =
      std::make_unique<lemming::dataplane::sai::SchedulerGroup::Stub>(chan);
  scheduler = std::make_unique<lemming::dataplane::sai::Scheduler::Stub>(chan);
  srv6 = std::make_unique<lemming::dataplane::sai::Srv6::Stub>(chan);
  stp = std::make_unique<lemming::dataplane::sai::Stp::Stub>(chan);
  switch_ = std::make_unique<lemming::dataplane::sai::Switch::Stub>(chan);
  system_port =
      std::make_unique<lemming::dataplane::sai::SystemPort::Stub>(chan);
  tam = std::make_unique<lemming::dataplane::sai::Tam::Stub>(chan);
  tunnel = std::make_unique<lemming::dataplane::sai::Tunnel::Stub>(chan);
  udf = std::make_unique<lemming::dataplane::sai::Udf::Stub>(chan);
  virtual_router =
      std::make_unique<lemming::dataplane::sai::VirtualRouter::Stub>(chan);
  vlan = std::make_unique<lemming::dataplane::sai::Vlan::Stub>(chan);
  wred = std::make_unique<lemming::dataplane::sai::Wred::Stub>(chan);
  entry = std::make_unique<lemming::dataplane::sai::Entrypoint::Stub>(chan);

  return SAI_STATUS_SUCCESS;
}

sai_status_t sai_api_uninitialize(void) { return SAI_STATUS_SUCCESS; }

sai_status_t sai_api_query(_In_ sai_api_t api, _Out_ void **api_method_table) {
  switch (api) {
    case SAI_API_SWITCH: {
      *api_method_table = const_cast<sai_switch_api_t *>(&l_switch);
      break;
    }
    case SAI_API_PORT: {
      *api_method_table = const_cast<sai_port_api_t *>(&l_port);
      break;
    }
    case SAI_API_FDB: {
      *api_method_table = const_cast<sai_fdb_api_t *>(&l_fdb);
      break;
    }
    case SAI_API_VLAN: {
      *api_method_table = const_cast<sai_vlan_api_t *>(&l_vlan);
      break;
    }
    case SAI_API_VIRTUAL_ROUTER: {
      *api_method_table =
          const_cast<sai_virtual_router_api_t *>(&l_virtual_router);
      break;
    }
    case SAI_API_ROUTE: {
      *api_method_table = const_cast<sai_route_api_t *>(&l_route);
      break;
    }
    case SAI_API_NEXT_HOP: {
      *api_method_table = const_cast<sai_next_hop_api_t *>(&l_next_hop);
      break;
    }
    case SAI_API_NEXT_HOP_GROUP: {
      *api_method_table =
          const_cast<sai_next_hop_group_api_t *>(&l_next_hop_group);
      break;
    }
    case SAI_API_ROUTER_INTERFACE: {
      *api_method_table =
          const_cast<sai_router_interface_api_t *>(&l_router_interface);
      break;
    }
    case SAI_API_NEIGHBOR: {
      *api_method_table = const_cast<sai_neighbor_api_t *>(&l_neighbor);
      break;
    }
    case SAI_API_ACL: {
      *api_method_table = const_cast<sai_acl_api_t *>(&l_acl);
      break;
    }
    case SAI_API_HOSTIF: {
      *api_method_table = const_cast<sai_hostif_api_t *>(&l_hostif);
      break;
    }
    case SAI_API_MIRROR: {
      *api_method_table = const_cast<sai_mirror_api_t *>(&l_mirror);
      break;
    }
    case SAI_API_SAMPLEPACKET: {
      *api_method_table = const_cast<sai_samplepacket_api_t *>(&l_samplepacket);
      break;
    }
    case SAI_API_STP: {
      *api_method_table = const_cast<sai_stp_api_t *>(&l_stp);
      break;
    }
    case SAI_API_LAG: {
      *api_method_table = const_cast<sai_lag_api_t *>(&l_lag);
      break;
    }
    case SAI_API_POLICER: {
      *api_method_table = const_cast<sai_policer_api_t *>(&l_policer);
      break;
    }
    case SAI_API_WRED: {
      *api_method_table = const_cast<sai_wred_api_t *>(&l_wred);
      break;
    }
    case SAI_API_QOS_MAP: {
      *api_method_table = const_cast<sai_qos_map_api_t *>(&l_qos_map);
      break;
    }
    case SAI_API_QUEUE: {
      *api_method_table = const_cast<sai_queue_api_t *>(&l_queue);
      break;
    }
    case SAI_API_SCHEDULER: {
      *api_method_table = const_cast<sai_scheduler_api_t *>(&l_scheduler);
      break;
    }
    case SAI_API_SCHEDULER_GROUP: {
      *api_method_table =
          const_cast<sai_scheduler_group_api_t *>(&l_scheduler_group);
      break;
    }
    case SAI_API_BUFFER: {
      *api_method_table = const_cast<sai_buffer_api_t *>(&l_buffer);
      break;
    }
    case SAI_API_HASH: {
      *api_method_table = const_cast<sai_hash_api_t *>(&l_hash);
      break;
    }
    case SAI_API_UDF: {
      *api_method_table = const_cast<sai_udf_api_t *>(&l_udf);
      break;
    }
    case SAI_API_TUNNEL: {
      *api_method_table = const_cast<sai_tunnel_api_t *>(&l_tunnel);
      break;
    }
    case SAI_API_L2MC: {
      *api_method_table = const_cast<sai_l2mc_api_t *>(&l_l2mc);
      break;
    }
    case SAI_API_IPMC: {
      *api_method_table = const_cast<sai_ipmc_api_t *>(&l_ipmc);
      break;
    }
    case SAI_API_RPF_GROUP: {
      *api_method_table = const_cast<sai_rpf_group_api_t *>(&l_rpf_group);
      break;
    }
    case SAI_API_L2MC_GROUP: {
      *api_method_table = const_cast<sai_l2mc_group_api_t *>(&l_l2mc_group);
      break;
    }
    case SAI_API_IPMC_GROUP: {
      *api_method_table = const_cast<sai_ipmc_group_api_t *>(&l_ipmc_group);
      break;
    }
    case SAI_API_MCAST_FDB: {
      *api_method_table = const_cast<sai_mcast_fdb_api_t *>(&l_mcast_fdb);
      break;
    }
    case SAI_API_BRIDGE: {
      *api_method_table = const_cast<sai_bridge_api_t *>(&l_bridge);
      break;
    }
    case SAI_API_TAM: {
      *api_method_table = const_cast<sai_tam_api_t *>(&l_tam);
      break;
    }
    case SAI_API_SRV6: {
      *api_method_table = const_cast<sai_srv6_api_t *>(&l_srv6);
      break;
    }
    case SAI_API_MPLS: {
      *api_method_table = const_cast<sai_mpls_api_t *>(&l_mpls);
      break;
    }
    case SAI_API_DTEL: {
      *api_method_table = const_cast<sai_dtel_api_t *>(&l_dtel);
      break;
    }
    case SAI_API_BFD: {
      *api_method_table = const_cast<sai_bfd_api_t *>(&l_bfd);
      break;
    }
    case SAI_API_ISOLATION_GROUP: {
      *api_method_table =
          const_cast<sai_isolation_group_api_t *>(&l_isolation_group);
      break;
    }
    case SAI_API_NAT: {
      *api_method_table = const_cast<sai_nat_api_t *>(&l_nat);
      break;
    }
    case SAI_API_COUNTER: {
      *api_method_table = const_cast<sai_counter_api_t *>(&l_counter);
      break;
    }
    case SAI_API_DEBUG_COUNTER: {
      *api_method_table =
          const_cast<sai_debug_counter_api_t *>(&l_debug_counter);
      break;
    }
    case SAI_API_MACSEC: {
      *api_method_table = const_cast<sai_macsec_api_t *>(&l_macsec);
      break;
    }
    case SAI_API_SYSTEM_PORT: {
      *api_method_table = const_cast<sai_system_port_api_t *>(&l_system_port);
      break;
    }
    case SAI_API_MY_MAC: {
      *api_method_table = const_cast<sai_my_mac_api_t *>(&l_my_mac);
      break;
    }
    case SAI_API_IPSEC: {
      *api_method_table = const_cast<sai_ipsec_api_t *>(&l_ipsec);
      break;
    }
    case SAI_API_GENERIC_PROGRAMMABLE: {
      *api_method_table =
          const_cast<sai_generic_programmable_api_t *>(&l_generic_programmable);
      break;
    }
    case SAI_API_BMTOR: {
      *api_method_table = const_cast<sai_bmtor_api_t *>(&l_bmtor);
      break;
    }
    default:
      LOG(WARNING) << "unknown API type " << api;
      return SAI_STATUS_NOT_IMPLEMENTED;
  }
  return SAI_STATUS_SUCCESS;
}

sai_status_t sai_log_set(_In_ sai_api_t api, _In_ sai_log_level_t log_level) {
  return SAI_STATUS_SUCCESS;
}

sai_object_type_t sai_object_type_query(_In_ sai_object_id_t object_id) {
  lemming::dataplane::sai::ObjectTypeQueryRequest req;
  lemming::dataplane::sai::ObjectTypeQueryResponse resp;
  grpc::ClientContext context;
  LOG(WARNING) << "type query " << object_id;
  req.set_object(object_id);
  grpc::Status status = entry->ObjectTypeQuery(&context, req, &resp);
  if (!status.ok()) {
    LOG(ERROR) << status.error_message();
    return SAI_OBJECT_TYPE_NULL;
  }
  LOG(WARNING) << "type query done " << resp.type();
  return static_cast<sai_object_type_t>(resp.type() - 1);
}

sai_status_t sai_query_attribute_capability(
    _In_ sai_object_id_t switch_id, _In_ sai_object_type_t object_type,
    _In_ sai_attr_id_t attr_id, _Out_ sai_attr_capability_t *attr_capability) {
  *attr_capability = {true, true, true};
  return SAI_STATUS_SUCCESS;
}

sai_status_t sai_query_attribute_enum_values_capability(
    _In_ sai_object_id_t switch_id, _In_ sai_object_type_t object_type,
    _In_ sai_attr_id_t attr_id,
    _Inout_ sai_s32_list_t *enum_values_capability) {

  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t sai_object_type_get_availability(
    _In_ sai_object_id_t switch_id, _In_ sai_object_type_t object_type,
    _In_ uint32_t attr_count, _In_ const sai_attribute_t *attr_list,
    _Out_ uint64_t *count) {
  *count = 1024;
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t sai_query_stats_capability(
    _In_ sai_object_id_t switch_id, _In_ sai_object_type_t object_type,
    _Inout_ sai_stat_capability_list_t *stats_capability) {
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t sai_get_maximum_attribute_count(_In_ sai_object_id_t switch_id,
                                             _In_ sai_object_type_t object_type,
                                             _Out_ uint32_t *count) {
  return SAI_STATUS_NOT_IMPLEMENTED;
}
sai_status_t sai_get_object_count(_In_ sai_object_id_t switch_id,
                                  _In_ sai_object_type_t object_type,
                                  _Out_ uint32_t *count) {
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t sai_get_object_key(_In_ sai_object_id_t switch_id,
                                _In_ sai_object_type_t object_type,
                                _Inout_ uint32_t *object_count,
                                _Inout_ sai_object_key_t *object_list) {
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t sai_bulk_get_attribute(_In_ sai_object_id_t switch_id,
                                    _In_ sai_object_type_t object_type,
                                    _In_ uint32_t object_count,
                                    _In_ const sai_object_key_t *object_key,
                                    _Inout_ uint32_t *attr_count,
                                    _Inout_ sai_attribute_t **attr_list,
                                    _Inout_ sai_status_t *object_statuses) {
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t sai_bulk_object_get_stats(
    _In_ sai_object_id_t switch_id, _In_ sai_object_type_t object_type,
    _In_ uint32_t object_count, _In_ const sai_object_key_t *object_key,
    _In_ uint32_t number_of_counters, _In_ const sai_stat_id_t *counter_ids,
    _In_ sai_stats_mode_t mode, _Inout_ sai_status_t *object_statuses,
    _Out_ uint64_t *counters) {
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t sai_bulk_object_clear_stats(
    _In_ sai_object_id_t switch_id, _In_ sai_object_type_t object_type,
    _In_ uint32_t object_count, _In_ const sai_object_key_t *object_key,
    _In_ uint32_t number_of_counters, _In_ const sai_stat_id_t *counter_ids,
    _In_ sai_stats_mode_t mode, _Inout_ sai_status_t *object_statuses) {
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t sai_query_object_stage(_In_ sai_object_id_t switch_id,
                                    _In_ sai_object_type_t object_type,
                                    _In_ uint32_t attr_count,
                                    _In_ const sai_attribute_t *attr_list,
                                    _Out_ sai_object_stage_t *stage) {
  return SAI_STATUS_NOT_IMPLEMENTED;
}

sai_status_t sai_query_stats_capability(
    _In_ sai_object_id_t switch_id, _In_ sai_object_type_t object_type,
    _Inout_ sai_stat_capability_list_t *stats_capability) {
  return SAI_STATUS_SUCCESS;
}

int main() {}
