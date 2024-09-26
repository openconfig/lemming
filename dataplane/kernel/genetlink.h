#ifndef DATAPLANE_KERNEL_GENETLINK_H_
#define DATAPLANE_KERNEL_GENETLINK_H_

#include <stdint.h>

int create_port(const char* family, const char* group);
int send_packet(int sock_idx, const void* pkt, uint32_t size, int in_ifindex, int out_ifindex,
                unsigned int context);

#endif