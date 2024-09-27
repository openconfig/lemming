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
int create_idx = 0;
const int max_sockets = 16;

int create_port(const char* family, const char* group) {
  fprintf(stderr, "creating port\n");
  if (nlsocks == NULL) {
    nlsocks = malloc(sizeof(struct nl_sock*) * max_sockets);
  }
  if (create_idx >= max_sockets) {
    return -1;
  }

  nlsocks[create_idx] = nl_socket_alloc();
  if (nlsocks[create_idx] == NULL) {
    return -1;
  }
  nl_socket_disable_auto_ack(nlsocks[create_idx]);
  int error = genl_connect(nlsocks[create_idx]);
  if (error < 0) {
    nl_socket_free(nlsocks[create_idx]);
    return -1;
  }
  family_id = genl_ctrl_resolve(nlsocks[create_idx], family);
  if (family_id < 0) {
    nl_socket_free(nlsocks[create_idx]);
    return -1;
  }
  int group_id = genl_ctrl_resolve_grp(nlsocks[create_idx], family, group);
  if (group_id < 0) {
    nl_socket_free(nlsocks[create_idx]);
    return -1;
  }

  nl_socket_set_peer_groups(nlsocks[create_idx], (1 << (group_id - 1)));

  return create_idx++;
}

int send_packet(int sock_idx, const void* pkt, uint32_t size, int in_ifindex,
                int out_ifindex, unsigned int context) {
  fprintf(stderr,"sending packet\n");
  fprintf(stderr,"populating packet to index %d\n", sock_idx);
  struct nl_msg* msg = nlmsg_alloc();
  genlmsg_put(msg, NL_AUTO_PORT, NL_AUTO_SEQ, family_id, 0, 0, 0, 1);
  fprintf(stderr, "putting src if index nl msg: %d\n", in_ifindex);
  NLA_PUT_S16(msg, GENL_PACKET_ATTR_IIFINDEX, in_ifindex);
  fprintf(stderr, "putting dst if index nl msg: %d\n", out_ifindex);
  NLA_PUT_S16(msg, GENL_PACKET_ATTR_OIFINDEX, out_ifindex);
  fprintf(stderr, "putting context nl msg: %d\n", context);
  NLA_PUT_U32(msg, GENL_PACKET_ATTR_CONTEXT, context);
  fprintf(stderr, "putting data nl msg, size: %d\n", size);
  NLA_PUT(msg, GENL_PACKET_ATTR_DATA, size, pkt);
  fprintf(stderr,"sending to index %d\n", sock_idx);
  if (nl_send(nlsocks[sock_idx], msg) < 0) {
    fprintf(stderr,"failed to send message\n");
    return -1;
  }
  fprintf(stderr,"sent packet\n");
  nlmsg_free(msg);
  return 0;
nla_put_failure:
  nlmsg_free(msg);
  return -1;
}