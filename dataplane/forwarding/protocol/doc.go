// Copyright 2022 Google LLC
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

// Package protocol enables Lucius to uniformly query and mutate network packets
// containing various protocols.
//
// The package registers a parser with the Lucius forwarding infrastructure.
// The parser describes the type of packets that it accepts (network protocols
// and their relative ordering).
//
// The package also provides utilities to implement network protocols.
// Each network protocol is implemented in its own package. They
// implement the interface protocol.Handler and register (optional) parse
// and add functions with package protocol.
//
// The supported protocol are: Ethernet (with VLAN tags), IP4, IP6,
// GRE (with Key and Sequence numbers), ARP, ICMP, TCP, UDP. It also supports
// recursive IP tunnels such as IPv6-over-IPv4-over-GRE-over-IPv6.
package protocol
