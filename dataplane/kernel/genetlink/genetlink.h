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

#ifndef DATAPLANE_KERNEL_GENETLINK_GENETLINK_H_
#define DATAPLANE_KERNEL_GENETLINK_GENETLINK_H_

#include <stdint.h>

// create_port create genetlink socket.
struct nl_sock* create_port(const char* family, const char* group);

void delete_port(void* sock);

// send_packet sends a packet with given metadata to specified port.
int send_packet(void* sock, int family, const void* pkt, uint32_t size,
                int in_ifindex, int out_ifindex, unsigned int context);

#endif  // DATAPLANE_KERNEL_GENETLINK_GENETLINK_H_
