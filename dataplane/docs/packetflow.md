# Packet Flow

## Overview

```plaintext
  ┌────────────────────────────────────────────┐
  │                                            │
  │   ┌──────┐                      ┌──────┐   │
  │   │ tap1 │                      │ tap3 │   │
  │   └─┬──▲─┘                      └─▲──┬─┘   │
  │     │  │                          │  │     │
  │     │  │      ┌──────────────┐    │  │     │
  │     │  └──────┤  CPU port    ├────┘  │     │
  │     │         └──────▲───────┘       │     │
  │     │                │               │     │
  │     │ ┌──────────────▼─────────────┐ │     │
  │     │ │      forwarding engine     │ │     │
  │     │ └─▲────────────────────────▲─┘ │     │
  │     │   │                        │   │     │
  │     │   │                        │   │     │
  │   ┌─▼───▼┐                      ┌▼───▼─┐   │
  │   │ eth1 │                      │ eth3 │   │
  └───┴──────┴──────────────────────┴──────┴───┘
```

## Terms

* Port: a l2 "physical" port that usually connected to another device
* Hostif: a "logical" network device (linux TAP) that corresponds a port
* CPU Port: the connection between the physical port and the hostifs.

As the dataplane is a software, this structure only exists to model a real hardware device.

## Ingress

1. Port Read
   1. Reads the packets from the underlying network device
2. Port To Interface
   1. Map from L2 port to L3 interface
   2. TODO: Support VLANs
3. Ingress VRF
   1. Map from L3 interface to VRF
4. Preingress
   1. ACL Stage: Custom match-action rules
   2. Hostif trap actions:
      1. Rules control which packets are sent to the CPU port (eg ARP).
5. Decap
   1. Decap L2 header
6. Ingress
   1. ACL Stage: Custom match-action rules

## Forwarding

1. FIB
   1. Longest prefix match VRF ID + IP DST:
      1. Next hop groups
      2. Next hops
      3. Ports
      4. Interfaces
      5. CPU port
2. CPU Port
   1. Once a packet is sent the CPU port is matched in two a hostif by:
      1. an IP2ME route: the exact IP assigned to a hostif.
      2. the default map: the corresponding hostif for the input port.

## Egress

1. Interface to Port
   1. Maps the egress L3 interface to the L2 port
2. Neighbor
   1. Maps the L3 interface + next hop IP (IP DST for connected routes) to L2 DST mac.
3. Egress
   1. ACL State: Custom match action rules.
4. SRC MAC
   1. Maps the L3 interface to SRC MAC address field
5. Port Write
   1. Write the packet to underlying network device.