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

#include "genetlink.h"                    // NOLINT(build/include_subdir)

#include <linux/netlink.h>
#include <netlink/genl/ctrl.h>
#include <netlink/genl/family.h>
#include <netlink/genl/genl.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

#define PACKET_BUFFER_SIZE 16384          // 16KB buffer
#define NL_SOCKET_BUFFER_SIZE 2097152     // 2MB socket buffer

enum {
  /* packet metadata */
  GENL_PACKET_ATTR_IIFINDEX,
  GENL_PACKET_ATTR_OIFINDEX,
  GENL_PACKET_ATTR_CONTEXT,
  GENL_PACKET_ATTR_DATA,
};

struct nl_sock* create_port(const char* family, const char* group) {
  fprintf(stderr, "creating port\n");

  struct nl_sock* sock = nl_socket_alloc();
  if (sock == NULL) {
    fprintf(stderr, "error: failed to alloc nl socket");
    return NULL;
  }
  nl_socket_disable_auto_ack(sock);

  int error = genl_connect(sock);
  if (error < 0) {
    fprintf(stderr, "error: failed to connect to genetlink: err %d\n", error);
    nl_socket_free(sock);
    return NULL;
  }

  int err = nl_socket_set_buffer_size(sock, NL_SOCKET_BUFFER_SIZE,
                                      NL_SOCKET_BUFFER_SIZE);

  // If increased size cannot be set, log error without crashing pkthandler.
  if (err < 0) {
    fprintf(stderr, "error: failed to set buffer size: err %d\n", err);
  }
  int group_id = genl_ctrl_resolve_grp(sock, family, group);
  if (group_id < 0) {
    fprintf(stderr, "error: failed to resolve group: err %d", group_id);
    nl_socket_free(sock);
    return NULL;
  }
  nl_socket_set_peer_groups(sock, (1 << (group_id - 1)));
  return sock;
}

void delete_port(void* sock) { nl_socket_free(sock); }

int send_packet(void* sock, int family, const void* pkt, uint32_t size,
                int in_ifindex, int out_ifindex, unsigned int context) {
  struct nl_msg* msg = nlmsg_alloc_size(PACKET_BUFFER_SIZE);
  if (msg == NULL) {
    fprintf(stderr, "failed to allocate packet\n");
    return -1;
  }
  genlmsg_put(msg, NL_AUTO_PORT, NL_AUTO_SEQ, family, 0, 0, 0, 1);
  NLA_PUT_S16(msg, GENL_PACKET_ATTR_IIFINDEX, in_ifindex);
  NLA_PUT_S16(msg, GENL_PACKET_ATTR_OIFINDEX, out_ifindex);
  NLA_PUT_U32(msg, GENL_PACKET_ATTR_CONTEXT, context);
  NLA_PUT(msg, GENL_PACKET_ATTR_DATA, size, pkt);
  fprintf(stderr, "sending packet size: %d\n", size);
  int err = nl_send(sock, msg);
  if (err < 0) {
    fprintf(stderr, "failed to send packet: %d\n", err);
    return -1;
  }
  nlmsg_free(msg);
  return 0;
nla_put_failure:
  fprintf(stderr, "nla_put_failure: packet exceeds nlmsg allcation size\n");
  nlmsg_free(msg);
  return -1;
}
