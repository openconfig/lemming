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
int family_id;
int create_index = 0;
const int max_sockets = 16;

int create_port(const char* family, const char* group) {
  if (nlsocks == NULL) {
    nlsocks = malloc(sizeof(struct nl_sock*) * max_sockets);
  }
  if (create_index >= max_sockets) {
    return -1;
  }

  struct nl_sock* nlsock = nlsocks[create_index];

  nlsock = nl_socket_alloc();
  if (nlsock == NULL) {
    return -1;
  }
  nl_socket_disable_auto_ack(nlsock);
  int error = genl_connect(nlsock);
  if (error < 0) {
    nl_socket_free(nlsock);
    return -1;
  }
  family_id = genl_ctrl_resolve(nlsock, family);
  if (family_id < 0) {
    nl_socket_free(nlsock);
    return -1;
  }
  int group_id = genl_ctrl_resolve_grp(nlsock, family, group);
  if (group_id < 0) {
    nl_socket_free(nlsock);
    return -1;
  }

  nl_socket_set_peer_groups(nlsock, (1 << (group_id - 1)));

  return create_index++;
}

int send_packet(int sock_idx, const void* pkt, uint32_t size, int in_ifindex,
                int out_ifindex, unsigned int context) {
  printf("creating nl msg sock idx: %d", sock_idx);
  struct nl_msg* msg = nlmsg_alloc();
  genlmsg_put(msg, NL_AUTO_PORT, NL_AUTO_SEQ, family_id, 0, 0, 0, 1);
  printf("putting src if index nl msg: %d", in_ifindex);
  NLA_PUT_S16(msg, GENL_PACKET_ATTR_IIFINDEX, in_ifindex);
  printf("putting dst if index nl msg: %d", out_ifindex);
  NLA_PUT_S16(msg, GENL_PACKET_ATTR_OIFINDEX, out_ifindex);
  printf("putting context nl msg: %d", context);
  NLA_PUT_U32(msg, GENL_PACKET_ATTR_CONTEXT, context);
  printf("putting data nl msg, size: %d", size);
  // NLA_PUT(msg, GENL_PACKET_ATTR_DATA, size, pkt);
  printf("sending to index %d", sock_idx);
  if (nl_send(nlsocks[sock_idx], msg) < 0) {
    printf("failed to send message");
    return -1;
  }
  printf("sent message");
  nlmsg_free(msg);
  return 0;
nla_put_failure:
  nlmsg_free(msg);
  return -1;
}