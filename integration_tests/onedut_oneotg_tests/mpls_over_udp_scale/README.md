# TE-18.3 MPLS in UDP Encapsulation with QoS Scheduler Scale Test

Building on TE-18.1 and TE-18.2, add scaling parameters

## Topology

* 32 ports as the 'input port set'
* 4 ports as "uplink facing"
* VLAN configurations
  * input vlans are distributed evenly across the 'input port set'

## Test setup

TODO: Complete test environment setup steps

inner_ipv6_dst_A = "2001:aa:bb::1/128"
inner_ipv6_dst_B = "2001:aa:bb::2/128"
inner_ipv6_default = "::/0"

ipv4_inner_dst_A = "10.5.1.1/32"
ipv4_inner_dst_B = "10.5.1.2/32"
ipv4_inner_default = "0.0.0.0/0"

outer_ipv6_src =      "2001:f:a:1::0"
outer_ipv6_dst_A =    "2001:f:c:e::1"
outer_ipv6_dst_B =    "2001:f:c:e::2"
outer_ipv6_dst_def =  "2001:1:1:1::0"
outer_dst_udp_port =  "5555"
outer_dscp =          "26"
outer_ip-ttl =        "64"

## Procedure

### TE-18.3.1 Scale

#### Scale targets

* Flow scale
  * 20,000 IPv4/IPv6 destinations
  * 1,000 vlans
  * Inner IP address space should be reused for each network-instance.
  * gRIBI client update rate `flow_r` = 1 update per second
  * Each gRIBI update include ip entries in batches of `flow_q` = 200
  * DUT packet forwarding updated within 1 second after adding entries

* Scheduler (policer) scale
  * 1,000 policer rates
  * 20,000 policer-policies / token buckets instantiations
  * Update policer-policies at 1 per `sched_r` = 60 seconds
  * Update policer-policies in a batch of `sched_q` = 1,000
  * Policer-policies changes should take effect within `sched_r` / 2 time

#### Scale profile A - 1 Network instance, 20k gRIBI NHG, 20k gRIBI exact match prefixes, 1 NH per NHG with same MPLS label
* 1 network-instance
* Create 20K gRIBI NHG with 1 NH per NHG with same MPLS label.
* Create 1 unique gRIBI prefixes which point to 1 NHG and repeat this 20k times, resulting in 20K total prefixes, pointing to 20k NHG.
* 1 MPLS label across all 20k tuples.

#### Scale profile B - 1 Network instance, 20k gRIBI NHG, 20k gRIBI exact match prefixes, 1 NH per NHG with different MPLS label
* 1024 network-instance.
* 1 VLAN per VRF.
* Create 20K gRIBI NHG with 1 NH per NHG with different MPLS label.
* Create 20 unique gRIBI prefixes which point to 1 NHG and repeat this 1k times, resulting in 20K total prefixes, pointing to 20k NHG.
* 20 gRIBI prefixes with one MPLS label.
* Create 1 unique gRIBI prefixes which point to 1 NHG and repeat this 20k times, resulting in 20K total prefixes, pointing to 20k NHG.
* 1 MPLS label across all 20k tuples.

#### Scale profile C - 1 Network instance, 20k gRIBI NHG, 20k gRIBI exact match prefixes, 8 NH per NHG with same MPLS label
* 1 network-instance
* Create 20K gRIBI NHG with 8 NH per NHG with same MPLS label.
* Create 1 unique gRIBI prefixes which point to 1 NHG and repeat this 20k times, resulting in 20K total prefixes, pointing to 20k NHG.
* Each NH has a different prefix.
* 1 MPLS label across all 20k tuples.

#### Scale profile D - 1 Network instance, 20k gRIBI NHG, 20k gRIBI exact match prefixes, 1 NH per NHG with same MPLS label to test gRIBI server scaling
* 1 network-instance
* Create 20K gRIBI NHG with 1 NH per NHG with same MPLS label.
* Create 1 unique gRIBI prefixes which point to 1 NHG and repeat this 20k times, resulting in 20K total prefixes, pointing to 20k NHG.
* 1 MPLS label across all 20k tuples.
* At traffic start time, gRIBI client to send 50 aft entries and their
    related NHG and NH at rate 50 per second.

#### Scale profile E - many vlans 1 NH per NHG with same MPLS label

* 20 ip destinations * 1,000 vlans = 20,000 'flows'
* 1000 network-instances have each have 1 static route pointing a static  NHG with 8 NH. The encap-headers for the NH should indicate MPLS in GRE encapsulation as per the canonical OC defined in test TE-18.1
* Configure static LSP's for the MPLS label to VLAN mapping.
* 1 default route in each VRF pointing to policy forwarding NHG.
* Create 20K gRIBI NHG with 1 NH per NHG with same MPLS label.
* Create 1 unique gRIBI prefixes which point to 1 NHG and repeat this 20k times, resulting in 20K total prefixes, pointing to 20k NHG.
* Each ingress vlan has 20 policer-policies = 10,000 'token buckets'
* Each policer is assigned rate limits matching one of 800 different possible limits between 1Gbps to 400Gbps in 0.5Gbps increments

#### Scale profile F - many vlans 1 NH per NHG with different MPLS label

* 20 ip destinations * 1,000 vlans = 20,000 'flows'
* 1000 network-instances have each have 1 static route pointing a static  NHG with 8 NH. The encap-headers for the NH should indicate MPLS in GRE encapsulation as per the canonical OC defined in test TE-18.1
* Configure static LSP's for the MPLS label to VLAN mapping.
* 1 default route in each VRF pointing to policy forwarding NHG.
* Create 20K gRIBI NHG with 1 NH per NHG with different MPLS label.
* Create 1 unique gRIBI prefixes which point to 1 NHG and repeat this 20k times, resulting in 20K total prefixes, pointing to 20k NHG.
* Each ingress vlan has 20 policer-policies = 10,000 'token buckets'
* Each policer is assigned rate limits matching one of 800 different possible limits between 1Gbps to 400Gbps in 0.5Gbps increments

#### Scale profile G - many destinations, few vlans

* 200 ip destinations * 100 vlans = 20,000 'flows'
* Each ingress vlan has 4 policer-policies = 4,000 'token buckets'
* The 200 ip destinations are split evenly between the 4 policers
* Each policer is assigned rate limits matching one of 800 different possible limits between 1Gbps to 400Gbps in 0.5Gbps increments

#### Procedure - 1 VRF flow Scale
* For scale profile A, create the following subsets TE-18.3.1.n
* Program all 20,000 flows
* Validate RIB_AND_FIB_ACK with FIB_PROGRAMMED is received from DUT for each one of them.

#### Procedure - 1024 VRF flow Scale
* For scale profile B, create the following subsets TE-18.3.1.n
* Program all 20,000 flows
* Validate RIB_AND_FIB_ACK with FIB_PROGRAMMED is received from DUT for each one of them.

#### Procedure - 1024 VRF flow Scale with 8 NHs per NHG
* For scale profile C, create the following subsets TE-18.3.1.n
* Program all 20,000 flows
* Validate config is applied at DUT for each one of them.

#### Procedure - gRIBI server scaling
* For scale profile D, create the following subsets TE-18.3.1.n
* At traffic start time, gRIBI client to send 50 aft entries and their
  related NHG and NH at rate 50 per second.
* Total operations will be 60k with 1k per second.
* Send 20 operations in a gRIBI call and 50 gRIBI calls per second.

#### Procedure - gRIBI server scaling from traffic perspective
* For scale profile D, create the following subsets TE-18.3.1.n
* At traffic start time, gRIBI client to send 50 aft entries and their
  related NHG and NH at rate 50 per second.
* Total operations will be 60k with 1k per second.
* Send 20 operations in a gRIBI call and 50 gRIBI calls per second.
* Check the traffic received has correct encaps.

#### Procedure - Flow Scale

* For each scale profile, create the following subsets TE-18.3.1.n
  * Configure ATE flows to send 100 pps per flow and wait for ARP
  * Send traffic for q flows (destination IP prefixes) for 2 seconds
  * At traffic start time, gRIBI client to send `flow_q` aft entries and their
    related NHG and NH at rate `flow_r`
  * Validate RIB_AND_FIB_ACK with FIB_PROGRAMMED is received from DUT within
    1 second
  * Measure packet loss.  Target packet loss <= 50%.
  * Repeat adding 200 flows until 20,000 flows have been added
  * Once reaching 20,000 flows, perform 1 iteration of modifying the first
    `flow_q` flows to use different NH,NHG

#### Procedure - Policer + Flow Scale

* For each scale profile, create the following subsets TE-18.3.1.n
  * Program all 20,000 flows
  * Every `sched_r` interval use gnmi.Set to replace `sched_q` scheduler policies
  * Verify packet loss changes for all flows within `sched_r` / 2 time


#### Procedure - VRF Scale

* For each scale profile, create the following subsets TE-18.3.1.n
  * Generate traffic stream for policy forwarding NHG and gRIBI NHG
  * Observe that MPLS over GRE and MPLS over UDP encapsulation are working properly.

#### OpenConfig Path and RPC Coverage

```yaml
paths:
  # qos scheduler config
  /qos/scheduler-policies/scheduler-policy/config/name:
  /qos/scheduler-policies/scheduler-policy/schedulers/scheduler/config/type:
  /qos/scheduler-policies/scheduler-policy/schedulers/scheduler/one-rate-two-color/config/cir:
  /qos/scheduler-policies/scheduler-policy/schedulers/scheduler/one-rate-two-color/config/bc:
  /qos/scheduler-policies/scheduler-policy/schedulers/scheduler/one-rate-two-color/config/queuing-behavior:
  /qos/scheduler-policies/scheduler-policy/schedulers/scheduler/one-rate-two-color/exceed-action/config/drop:

  # qos classifier config
  /qos/classifiers/classifier/config/name:
  /qos/classifiers/classifier/terms/term/config/id:
  #/qos/classifiers/classifier/terms/term/conditions/next-hop-group/config/name: # TODO: new OC leaf to be added

  # qos input-policies config - TODO: a new OC subtree (/qos/input-policies)
  # /qos/input-policies/input-policy/config/name:
  # /qos/input-policies/input-policy/config/classifier:
  # /qos/input-policies/input-policy/config/scheduler-policy:

  # qos interface config
  #/qos/interfaces/interface/subinterface/input/config/policies:   # TODO:  new OC leaf-list (/qos/interfaces/interface/input/config/policies)

  # qos interface scheduler counters
  /qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/conforming-pkts:
  /qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/conforming-octets:
  /qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/exceeding-pkts:
  /qos/interfaces/interface/input/scheduler-policy/schedulers/scheduler/state/exceeding-octets:

  # afts next-hop counters
  /network-instances/network-instance/afts/next-hops/next-hop/state/counters/packets-forwarded:
  /network-instances/network-instance/afts/next-hops/next-hop/state/counters/octets-forwarded:

  # afts state paths set via gRIBI
  # TODO: https://github.com/openconfig/public/pull/1153
  #/network-instances/network-instance/afts/next-hops/next-hop/mpls-in-udp/state/src-ip:
  #/network-instances/network-instance/afts/next-hops/next-hop/mpls-in-udp/state/dst-ip:
  #/network-instances/network-instance/afts/next-hops/next-hop/mpls-in-udp/state/ip-ttl:
  #/network-instances/network-instance/afts/next-hops/next-hop/mpls-in-udp/state/dst-udp-port:
  #/network-instances/network-instance/afts/next-hops/next-hop/mpls-in-udp/state/dscp:

rpcs:
  gnmi:
    gNMI.Set:
      union_replace: true
      replace: true
    gNMI.Subscribe:
      on_change: true
  gribi:
    gRIBI.Modify:
    gRIBI.Flush:
```

## Required DUT platform

* FFF
