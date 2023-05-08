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

#include <grpc/grpc.h>
#include <grpcpp/channel.h>
#include <grpcpp/client_context.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>

#include <fstream>

#include "dataplane/standalone/log/log.h"
#include "dataplane/standalone/lucius/lucius_clib.h"
#include "dataplane/standalone/sai/acl.h"
#include "dataplane/standalone/sai/bfd.h"
#include "dataplane/standalone/sai/bridge.h"
#include "dataplane/standalone/sai/buffer.h"
#include "dataplane/standalone/sai/counter.h"
#include "dataplane/standalone/sai/debug_counter.h"
#include "dataplane/standalone/sai/dtel.h"
#include "dataplane/standalone/sai/fdb.h"
#include "dataplane/standalone/sai/hash.h"
#include "dataplane/standalone/sai/hostif.h"
#include "dataplane/standalone/sai/ipmc.h"
#include "dataplane/standalone/sai/ipmc_group.h"
#include "dataplane/standalone/sai/isolation_group.h"
#include "dataplane/standalone/sai/l2mc.h"
#include "dataplane/standalone/sai/l2mc_group.h"
#include "dataplane/standalone/sai/lag.h"
#include "dataplane/standalone/sai/macsec.h"
#include "dataplane/standalone/sai/mcast_fdb.h"
#include "dataplane/standalone/sai/mirror.h"
#include "dataplane/standalone/sai/mpls.h"
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
#include "dataplane/standalone/sai/segmentroute.h"
#include "dataplane/standalone/sai/stp.h"
#include "dataplane/standalone/sai/switch.h"
#include "dataplane/standalone/sai/system_port.h"
#include "dataplane/standalone/sai/tam.h"
#include "dataplane/standalone/sai/tunnel.h"
#include "dataplane/standalone/sai/udf.h"
#include "dataplane/standalone/sai/virtual_router.h"
#include "dataplane/standalone/sai/vlan.h"
#include "dataplane/standalone/sai/wred.h"
#include "dataplane/standalone/translator.h"

extern "C" {
#include "inc/sai.h"
}

std::shared_ptr<Translator> translator;

// TODO(dgrau): implement this without using gRPC.
sai_status_t sai_api_initialize(
    _In_ uint64_t flags, _In_ const sai_service_method_table_t *services) {
  LUCIUS_LOG_FUNC();
  initialize(GoInt(50000));

  auto chan = grpc::CreateChannel("localhost:50000",
                                  grpc::InsecureChannelCredentials());
  translator = std::make_shared<Translator>(chan);
  return SAI_STATUS_SUCCESS;
}

sai_status_t sai_api_query(_In_ sai_api_t api, _Out_ void **api_method_table) {
  LUCIUS_LOG_FUNC();
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
    case SAI_API_SEGMENTROUTE: {
      *api_method_table = const_cast<sai_segmentroute_api_t *>(&l_segmentroute);
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
    default:
      return SAI_STATUS_SUCCESS;
  }
  return SAI_STATUS_SUCCESS;
}

sai_status_t sai_log_set(_In_ sai_api_t api, _In_ sai_log_level_t log_level) {
  return SAI_STATUS_SUCCESS;
}

sai_object_type_t sai_object_type_query(_In_ sai_object_id_t object_id) {
  LUCIUS_LOG_FUNC();
  return translator->getObjectType(object_id);
}

int main() {}
