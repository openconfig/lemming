# TE-18.3 MPLS in UDP Encapsulation Scale Test

Building on TE-18.1 and TE-18.2, this test focuses on scaling gRIBI-programmed MPLS-over-UDP tunnels and associated forwarding entries, parameterized by key scaling dimensions.

## Topology

- 32 ports as the 'input port set' (Ingress)
- 4 ports as "uplink facing" (Egress)
- Network Instances (VRFs) will be mapped from ingress ports/subinterfaces as needed by scale profiles.

## Test setup

TODO: Complete test environment setup steps

inner_ipv6_dst_A = "2001:aa:bb::1/128" inner_ipv6_dst_B = "2001:aa:bb::2/128" inner_ipv6_default = "::/0"

ipv4_inner_dst_A = "10.5.1.1/32" ipv4_inner_dst_B = "10.5.1.2/32" ipv4_inner_default = "0.0.0.0/0"

outer_ipv6_src = "2001:f:a:1::0" outer_ipv6_dst_A = "2001:f:c:e::1" outer_ipv6_dst_B = "2001:f:c:e::2" outer_ipv6_dst_def = "2001:1:1:1::0" outer_dst_udp_port = "5555" outer_dscp = "26" outer_ip-ttl = "64"

## Procedure

### TE-18.3.1 Scaling Dimensions and Targets

This test evaluates scaling across the following dimensions using gRIBI. The test profiles below represent different parameter combinations of these dimensions.

- **Network Instances (VRFs):** Number of separate routing instances.
- **Next Hop Groups (NHGs):** Total number of NHGs programmed. Target: **Up to 20,000** (profile-dependent).
- **Next Hops (NHs):** Total number of NHs programmed (NHs per NHG * Total NHGs).
- **NHs per NHG:** Number of NH entries within each NHG (e.g., 1 or 8).
- **Prefixes:** Total number of unique IPv4/IPv6 exact-match forwarding entries (routes). Target: **20,000**.
- **(Unique Destination IP + MPLS) Tuples:** The combination of the inner destination IP and the MPLS label used in the NH encapsulation. Target: **Up to 20,000 unique tuples**.
- **MPLS Labels:** Number and uniqueness of MPLS labels used in NH encapsulation. **Constraint:** The number of unique MPLS labels must equal the number of VRFs (#MPLS Labels == #VRFs).
- **gRIBI Operations Rate (QPS):** Rate of gRIBI Modify requests or operations per second.
- **gRIBI Batch Size:** Number of AFT entries (or operations) per ModifyRequest.
- **Convergence:** DUT packet forwarding updated within **1 second** after receiving FIB_PROGRAMMED acknowledgement for added entries (baseline).
- **IP Address Reuse:** Inner IP destination prefixes should be reused across different Network Instances where applicable.

### TE-18.3.2 Scale Profiles

#### Profile 1 (Single VRF)

- **Goal:** Baseline single VRF scale (Exact Label Match scenario).
- **Network Instances (VRFs):** 1 (DEFAULT).
- **Total NHGs:** 20,000.
- **NHs per NHG:** 1.
- **MPLS Labels:** 1 (consistent with #VRFs = 1). Same label used for all NHs.
- **Total Prefixes:** 20,000 (e.g., 10k IPv4, 10k IPv6).
- **Unique (Dest IP + MPLS) Tuples:** 20,000 (different destination IPs, same MPLS label).
- **Prefix Mapping:** 1 unique prefix -> 1 unique NHG (1:1).
- **gRIBI Rate/Batch:** Baseline (e.g., 1 ModifyRequest/sec, 200 entries/request) - QPS not the primary focus here.

#### Profile 2 (Multi-VRF)

- **Goal:** Scale across multiple VRFs with unique labels per VRF.
- **Network Instances (VRFs):** 1024.
- **Total NHGs:** 20,000 (distributed across VRFs, ~19-20 NHGs/VRF).
- **NHs per NHG:** 1.
- **MPLS Labels:** 1024 unique labels (1 label assigned per VRF, consistent with #VRFs = 1024).
- **Total Prefixes:** 20,000 (distributed across VRFs, ~19-20 prefixes/VRF).
- **Unique (Dest IP + MPLS) Tuples:** 20,000 (e.g., 20 unique destination IPs reused per MPLS label/VRF).
- **Prefix Mapping:** Prefixes within a VRF map to NHGs using that VRF's unique MPLS label.
- **Inner IP Reuse:** Required.
- **gRIBI Rate/Batch:** Baseline - QPS not the primary focus here.

#### Profile 3 (Multi-VRF)

- **Goal:** Similar to Profile 2, but test potentially skewed distribution of prefixes/routes per VRF/label.
- **Network Instances (VRFs):** 1024.
- **Total NHGs:** 20,000.
- **NHs per NHG:** 1.
- **MPLS Labels:** 1024 unique labels (1 per VRF).
- **Total Prefixes:** 20,000.
- **Unique (Dest IP + MPLS) Tuples:** 20,000.
- **Prefix Mapping:** Similar to Profile 2, but the distribution of the 20k prefixes across the 1024 VRFs/labels might be intentionally uneven (e.g., some VRFs have many more prefixes than others). _Exact skew pattern TBD._
- **Inner IP Reuse:** Required.
- **gRIBI Rate/Batch:** Baseline - QPS not the primary focus here.

#### Profile 4 (Single VRF)

- **Goal:** Test ECMP scale within a single VRF (Deridex/Zakim use-case).
- **Network Instances (VRFs):** 1 (DEFAULT).
- **Total NHGs:** 2,500.
- **NHs per NHG:** 8 (each NH having a different destination IP).
- **Total NHs:** 20,000 (2500 NHGs * 8 NH/NHG).
- **MPLS Labels:** 1 (consistent with #VRFs = 1). Same label used for all NHs.
- **Total Prefixes:** 20,000 (e.g., 10k IPv4, 10k IPv6).
- **Unique (Dest IP + MPLS) Tuples:** 20,000 (different destination IPs across all NHs, same MPLS label).
- **Prefix Mapping:** 8 unique prefixes -> 1 unique NHG (8:1 mapping, repeated 2500 times).
- **gRIBI Rate/Batch:** Baseline - QPS not the primary focus here.

#### Profile 5 (Single VRF)

- **Goal:** Test gRIBI control plane QPS scaling and impact on dataplane. Uses Profile 1 as the base state.
- **Network Instances (VRFs):** 1 (DEFAULT).
- **Total NHGs:** 20,000.
- **NHs per NHG:** 1.
- **MPLS Labels:** 1.
- **Total Prefixes:** 20,000.
- **Unique (Dest IP + MPLS) Tuples:** 20,000.
- **Prefix Mapping:** 1:1.
- **gRIBI Operations:** Program/Modify the full 20k entries (Prefix + NHG + NH = ~3 ops/entry = ~60k ops total).

  - Target Rate: **6,000 operations/second** (aiming to update the full table in ~10 seconds).
  - Target Batch Size: **60 operations/ModifyRequest**.
  - Target Request Rate: **100 ModifyRequests/second**.
  - Operation Mix: Test with **50% ADD, 50% DELETE** operations during high-rate phase.

- **Dataplane Validation:** Ensure live traffic forwarding remains stable and correct during high-rate gRIBI operations.

### TE-18.3.3 Validation Procedures

#### Procedure - Single VRF Validation (Profiles 1, 4)

- Program all gRIBI entries (NHs, NHGs, Prefixes) according to the profile using baseline rate/batch.
- Validate `RIB_ACK` / `FIB_PROGRAMMED` status is received from DUT for all entries.
- Verify AFT state on DUT for a sample of entries (NH, NHG, Prefix -> NHG mapping).
- Send traffic matching programmed prefixes from appropriate ingress ports.
- Verify traffic is received on egress ports with correct MPLS-over-UDP encapsulation (correct outer IPs, UDP port, MPLS label).
- Measure packet loss (target: <= 1% steady state).
- Delete all gRIBI entries.
- Verify AFT state shows entries removed.
- Verify traffic loss is 100%.

#### Procedure - Multi-VRF Validation (Profiles 2, 3)

- Program all gRIBI entries across all specified VRFs according to the profile using baseline rate/batch.
- Validate `RIB_ACK` / `FIB_PROGRAMMED` status for all entries.
- Verify AFT state on DUT for a sample of entries within different VRFs.
- Send traffic matching programmed prefixes, ensuring traffic is directed to the correct VRF (e.g., via appropriate ingress interface mapping).
- Verify traffic is received with correct MPLS-over-UDP encapsulation, including the VRF-specific MPLS label.
- Measure packet loss (target: <= 1% steady state).
- Delete all gRIBI entries.
- Verify AFT state shows entries removed across VRFs.
- Verify traffic loss is 100%.

#### Procedure - ECMP Validation (Profile 4)

- Perform Single VRF Validation steps.
- Additionally, verify that traffic sent towards prefixes mapped to the ECMP NHG is distributed across the multiple NHs within that NHG (requires ATE support for flow analysis or DUT counter validation for NH packet/octet counters).

#### Procedure - gRIBI Rate Validation (Profile 5)

- Establish the baseline state (e.g., program 20k entries as per Profile 1).
- Start traffic flows matching the programmed entries. Verify baseline forwarding and low loss.
- Initiate high-rate gRIBI Modify operations (e.g., 100 ModifyRequests/sec, 60 ops/request, 50% ADD/50% DELETE mix targeting existing/new entries).
- Monitor gRIBI operation results (ACKs) for success/failure and latency.
- Continuously monitor traffic forwarding during the high-rate gRIBI phase.

  - Verify traffic uses correct encapsulation based on the programmed state.
  - Measure packet loss (target: minimal loss, allowing for brief transient loss during updates, but stable low loss overall).

- Validate `RIB_ACK` / `FIB_PROGRAMMED` status is received promptly for updates.

- Verify AFT state on DUT reflects the changes made during the high-rate phase.

- Stop high-rate programming and measure steady-state loss again.

#### Investigation - VRF Impact on QPS

- As an extension, investigate if the number of VRFs impacts gRIBI QPS or dataplane stability during high-rate updates. This could involve running a variation of Profile 5 using the multi-VRF setup from Profile 2 or 3 as the baseline state.

### TE-18.3.4 OpenConfig Path and RPC Coverage

```yaml
paths:
  # AFTs Next-Hop state (Verification)
  /network-instances/network-instance/afts/next-hops/next-hop/state/counters/packets-forwarded:
  /network-instances/network-instance/afts/next-hops/next-hop/state/counters/octets-forwarded:
  /network-instances/network-instance/afts/next-hops/next-hop/state/ip-address: # NH IP

  # AFTs Next-Hop-Group state (Verification)
  /network-instances/network-instance/afts/next-hop-groups/next-hop-group/state/next-hop: # Verify NHs in NHG

  # AFTs Prefix Entry state (Verification)
  /network-instances/network-instance/afts/ipv4-unicast/ipv4-entry/state/next-hop-group:
  /network-instances/network-instance/afts/ipv6-unicast/ipv6-entry/state/next-hop-group:

rpcs:
  gnmi:
    # Primarily used for verification (Subscribe/Get)
    gNMI.Subscribe:
      on_change: true
    gNMI.Get:
  gribi:
    # Used for programming all AFT entries
    gRIBI.Modify:
    gRIBI.Flush:
```

## Required DUT platform

- FFF
