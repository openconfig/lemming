# The Lemming Dataplane

The dataplane reads and writes packets from ports and the forwards packets between them.

## Overview

The dataplane consists of a several components:

1. forwarding: packet processing actions, tables, and ports.
1. saiserver: implements a "proto-ized" version of the [Switch Abstraction Interface API](https://github.com/opencomputeproject/SAI) and defines the forwarding pipeline
1. apigen: code to generate the proto API for the saiserver and C++ code that calls the API
1. dplanerc: gNMI reconciler that configures the dataplane based on OpenConfig schema
1. standalone: shared library that a implements the C SAI API and container running only the dataplane

There are two ways to run the dataplane. The first as part of Lemming, using the gNMI reconciliation. The second is a SAI implementation (e.g. in SONIC).

## Components

### forwarding

The forwarding directory contains packet processing actions, tables, and ports. These are the building blocks of the dataplane: they can be combined to create the
packet processing pipeline. The forwarding behavior is configured using the [gRPC API](../proto/forwarding/forwarding_service.proto).

Generally, external users should not program the dataplane using this API directly. It is too flexible and low-level. Instead, the saiserver API should be used instead.

Each port has a stack of input and output actions. The forwarding server processes an incoming packet using the input actions of the input port.
If after processing all the actions, if there is an output port, then the server processes the output actions of the output port.

The following sections describe the key parts of the forwarding server. The relevant packages and protobuf files contain more complete documentation.

#### Actions

Actions can modify the contents of the packet, output the packet etc...

Some of the important actions include:

* Update: Set packet headers values
* Transmit: Set the output port
* Lookup: Lookup in a table and process the resulting action. Read more in the [Tables section](#tables).

Implementation: [fwdport](forwarding/fwdaction/actions)

#### Ports

Ports define how to read and write the packets. A port contains a list of input and output actions that define actions to take on a packet after reading and before writing.

Note: If no actions are added to a port: all packets are dropped!

Example ports types:

1. Kernel: uses linux AF_PACKET
2. TAP: uses ioctl and a file descriptor to create/read to TAP interface

Implementation: [fwdport](forwarding/fwdport/ports)

#### Tables

Tables match against packet headers and perform actions. Each table has a default action list, and list of entries. An entry defines a set of field and values to match against
and a list of the corresponding actions.

>Note: After a packet is looked up in a table, the new actions are processed before any pending actions.

Example table types:

1. Prefix: longest matching prefix of header values
2. Exact: exact match

Implementation: [fwdport](forwarding/fwdtable)

### saiserver

The saiserver is a gRPC API generated from the [SAI API](https://github.com/opencomputeproject/SAI). It is the recommended API for programming the dataplane.
It is simpler than the forwarding API, can be call from a variety of languages, and is compatible with software already using the SAI API (eg SONiC).
The saiserver package implements the API by calling the corresponding forwarding API. This package also configures the forwarding pipeline, by setting up the
required tables.

A simple example: the CreateRouteEntry RPC adds a entry to the FIB. This looks something like:

1. Compute forwarding actions
     1. Update the next hop ip metadata field.
     2. If the next hop is port: transmit action.
     3. If next hop is next hop: update packet metadata next hop ID, lookup next-hop table.
     4. If next hop is a group: update packet metadata next hop group, lookup next-hop-group table.
2. Figure out the address family of the prefix.
3. Add forwarding table entry to correct fib: prefix entry with fields vrf and next hop ip and above action.

### dplanerc

The dplanerc implements a gNMI reconciler based on [reconciler package](../gnmi/reconciler/reconciler.go). These reconcilers convert from OpenConfig to
saiserver proto messages to support configuring the dataplane from gNMI.

>Note: The fib uses a custom, simple proto message because OpenConfig does not model configured forwarding tables (only state).

### standalone

The dataplane standalone is meant to be used by a NOS to configure a dataplane using the C SAI API. (eg SONiC).
The current design runs the dataplane separately from the NOS and all communication is over gRPC.

* The sai directory contains generated C++ code that converts the C SAI API to protobuf API and performs the RPCs.
  * Note: the common.cc and common.h contain hand-crafted conversion funcs.
* The packetio package receives packets punted to the CPU port and transmits to the correct interface. It programs the dataplane with the IP assigned to the local interfaces.
* The lucius directory is the Go main that starts the dataplane RPC server.
* The entrypoint.cc is compiled to a shared library that can be loaded in the NOS.

### apigen

These packages generate the protobuf and C++ source and headers based on the [SAI](https://github.com/opencomputeproject/SAI/tree/master/inc) headers.
