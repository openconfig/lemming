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

#include "genetlink.h"

#include <linux/netlink.h>
#include <netlink/genl/ctrl.h>
#include <netlink/genl/family.h>
#include <netlink/genl/genl.h>
#include <stdint.h>
#include <stdio.h>

#include "stdlib.h"

enum {
  /* packet metadata */
  GENL_PACKET_ATTR_IIFINDEX,
  GENL_PACKET_ATTR_OIFINDEX,
  GENL_PACKET_ATTR_CONTEXT,
  GENL_PACKET_ATTR_DATA,
};

struct nl_sock** nlsocks = NULL;
int* families = NULL;
int create_idx = 0;
const int max_sockets = 20;

int create_port(const char* family, const char* group) {
  fprintf(stderr, "creating port\n");
  if (nlsocks == NULL) {
    nlsocks = malloc(sizeof(struct nl_sock*) * max_sockets);
    families = malloc(sizeof(int)* max_sockets);
  }
  if (create_idx >= max_sockets) {
    fprintf(stderr,"error: created more ports than max");
    return -1;
  }

  nlsocks[create_idx] = nl_socket_alloc();
  if (nlsocks[create_idx] == NULL) {
    fprintf(stderr,"error: failed to alloc nl socket");
    return -1;
  }
  nl_socket_disable_auto_ack(nlsocks[create_idx]);
  int error = genl_connect(nlsocks[create_idx]);
  if (error < 0) {
    fprintf(stderr,"error: failed to disable auto ack");
    nl_socket_free(nlsocks[create_idx]);
    return error;
  }
  families[create_idx] = genl_ctrl_resolve(nlsocks[create_idx], family);
  if (families[create_idx] < 0) {
    fprintf(stderr,"error: failed to resolve family");
    nl_socket_free(nlsocks[create_idx]);
    return -1;
  }
  int group_id = genl_ctrl_resolve_grp(nlsocks[create_idx], family, group);
  if (group_id < 0) {
    fprintf(stderr,"error: failed to resolve group");
    nl_socket_free(nlsocks[create_idx]);
    return group_id;
  }

  nl_socket_set_peer_groups(nlsocks[create_idx], (1 << (group_id - 1)));
  return create_idx++;
}

int send_packet(int sock_idx, const void* pkt, uint32_t size, int in_ifindex,
                int out_ifindex, unsigned int context) {
  struct nl_msg* msg = nlmsg_alloc();
  genlmsg_put(msg, NL_AUTO_PORT, NL_AUTO_SEQ, families[sock_idx], 0, 0, 0, 1);
  NLA_PUT_S16(msg, GENL_PACKET_ATTR_IIFINDEX, in_ifindex);
  NLA_PUT_S16(msg, GENL_PACKET_ATTR_OIFINDEX, out_ifindex);
  NLA_PUT_U32(msg, GENL_PACKET_ATTR_CONTEXT, context);
  NLA_PUT(msg, GENL_PACKET_ATTR_DATA, size, pkt);
  fprintf(stderr,"sending packet size: %d\n", size);
  if (nl_send(nlsocks[sock_idx], msg) < 0) {
    fprintf(stderr,"failed to send packet\n");
    return -1;
  }
  nlmsg_free(msg);
  return 0;
nla_put_failure:
  nlmsg_free(msg);
  return -1;
}