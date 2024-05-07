# PacketIO Handler

The PacketIO handler is a library and systemd service that reads and writes packets from the CPU port of the dataplane and manages the creation of TAP interfaces for SAI hostif.
Packets are transmited and received using a bidirectional gRPC stream.

* [Proto](../proto/packetio/packetio.proto)
* [Client](../standalone/pkthandler/main.go)
* [GENETLINK Ref](https://github.com/sonic-net/sonic-genl-packet)

## HostPortControl RPC

The HostPortControl RPC is a bidirectional gRPC stream. The client initiates the stream, then receives messages from the server. Upon receiving the a message, the client should create or remove the port, then MUST reply with the status.
In the port message, both the `port_id` and `dataplane_port` ID are specified: the client maps the from `dataplane_port` to the ifindex of the port's corresponding netdev (SAI hostif) as required for GENETLINK metadata.

> This is a bidirectional stream because we want the packetIO handler to be the client and the dataplane to be the server (same the proto SAI API). Also, this ensures this client is ready when the stream starts.

The supported port types are netdev (ie TAP interface) and GENETLINK.

## Packet RPC

The Packet RPC is a bidirectional gRPC stream. In this RPC, the send and receives streams are independent. The dataplane sends packets to the client with the ID of the destination host port. The Packet RPC reads from all interfaces created and includes in the ID of the originating host port in the message.

## Client

This repository contains an implementation of the PacketIO handler. All packets are sent and received in a single queue. The library supports reading and writing for netdev ports and writing to GENETLINK ports.
