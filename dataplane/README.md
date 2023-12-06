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

Generally, external users should not program the dataplane using this API directly. It is too flexible and low-table. Instead, the saiserver API should be used instead.

Each port has a stack of input and output actions. The forwarding server processes an incoming packet using the input actions of the output.
If after processing all the actions, there is an output port, then the server processes the output actions of the output port.

The following sections describe the key parts of the forwarding server. The relevant packages and protobuf files contain more complete documentation.

#### Actions

Actions can modify the contents of the packet, output the packet etc...

Some the important actions include:

* Update: Set packet headers values
* Transmit: Set the output port
* Lookup: Lookup in a table

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

The saiserver is a gRPC API generated from the [SAI API](https://github.com/opencomputeproject/SAI). It provides a simpler API for programming the dataplane.
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

These packages generate the protobuf and C++ source and headers based on the SAI headers.

## How to add a new feature

This section describes how to add a new feature to the dataplane using a hypothetical and very contrived example (and just about the worst case).

Let's implement a new feature: based on the `foo` packet header, do a shortest prefix match on dst ip and on match randomize the payload bytes.

### Forwarding support for feature

The first step is to check if the forwarding supports everything we need. The answer is usually yes, but not in this case.

1. There is no shortest prefix forwarding table, so let's add it.
   1. Add a new value to the `TableTable` enum in the [forwarding_table](../proto/forwarding/forwarding_table.proto).
   2. Add a new entry desc proto message for the new table and it to the `oneof entry` in `EntryDesc`.
   3. Implement the table in a new package in the [fwdtable](forwarding/fwdtable/).
   4. Tables must implement an interface and register a builder. (see other packages for examples).
2. There is no randomized payload action, so let's add it.
   1. Follow similar steps as the table.
3. There is no `foo` header parser, so let's add it.
   1. Add a new enum value to the PacketFieldNum [forwarding_common](../proto/forwarding/forwarding_common.proto).
   2. Update the forwarding/protocol package with the logic to parse this new header.

### saiserver implementation

Now that the datataplane supports all the features we need, we have to provide an API to configure it. In the saiserver package, there are structs that
embed the unimplemented server for all the saipb services.

For this example, we are going to assume that this feature is defined in the saipb API as follows.

```go
func (f *foo) CreateFoo(ctx context.Context, req *saipb.CreateFooRequest) (*saipb.CreateFooResponse, error) {
}

func (f *foo) RemoveFoo(ctx context.Context, req *saipb.RemoveFooRequest) (*saipb.RemoveFooResponse, error) {
}
```

Let's assume we want to implement the CreateFoo and RemoveFoo RPCs.

1. First, we need to set up the table and add it as a step in the forwarding pipeline.
   1. We need to create the table. Since there's only one table, creating it in the CreateSwitch RPC is usually the right place.
      1. In the CreateSwitch, create a new shortest prefix table table, which we can call the "foo-table"
   2. Now, we need to determine where in the pipeline should this action take place:
      1. For this example, let's pretend it belongs after updating the DST MAC from the neighbor table.
      2. Most of the pipeline is static and set using input actions on every port. If we look at [ports.go](saiserver/ports.go), we can see all the actions applied on every port.
      3. In this, we are looking up what to put in our "foo-table", so we add a new lookup action to "foo-table" to slices of input actions.
2. Next, we need to implement the CreateFoo and RemoveFoo RPC.
   1. In CreateFoo, we need convert the CreateFooRequest to a `fwdpb.TableEntryAddRequest`
      1. Implementations will vary, there are plenty of examples in the saiserver package.
      2. Make sure that we add randomization action to the entry.
   2. In RemoveFoo, we need to do mostly the same.

>Note: A new feature may just require supporting an additional attribute in already implemented. If the feature is not in SAI API, then it can be configured using the forwarding API directly or by patching the generated protos.

### gNMI reconcilation

Not done yet! Now that the feature is available via the saiserver API, we may also want to configure it using gNMI.

1. In dplanerc, let's add a new [reconciler](../gnmi/reconciler/reconciler.go).
   1. First, we need to determine which OpenConfig are needed to configure this feature.
      1. The [OpenConfig site](https://openconfig.net/projects/models/paths/index.html) is the best reference.
   2. gNMI and Lemming use an eventual consistency model, so in the reconciler we need to subscribe to the config paths for our feature.
      1. Generally, using ygnmi.Watch and a batch subscription is useful.
   3. If the config does not equal the state, we need to reconcile the differences.
      1. This usually looks like, if config != state, call saiserver RPC to fix (Create or Remove)
      2. If the call was successful, then update the state.
      3. Note: Lemming has a special `gnmiclient` client that always ygnmi.Set on state paths.
      4. Note: For performance reasons: it is best to use BatchUpdate (not replace) individual leaves.

Almost there! Now write tests.
