# Lucius Dataplane Tests

## Packet Traces

### Setup

1.  Checkout lemming repository `git clone https://github.com/openconfig/lemming.git`
1.  Install lemctl `go install cmd/lemctl`

The lemctl CLI contains useful tools for debugging Lemming and Lucius.

To enable packet traces, the lucius containers must be logging at level 3 or
higher. Add the "-v=3" arg in the topology file for all lucius containers.

### Collect and Analyze

The verbose logging generates lots of logs, it is best to copy the log directory
from the pod to inspect it. This command copies the /tmp folder to local folder
tmp: `kubectl exec -n twodut-alpine "alpine-dut" -c dataplane -- tar cf -
"/tmp/" | tar xf -`

The packet traces are logged in the INFO file.

The findtrace command is a search tool (like grep), it searches the logs for
matches and prints the full trace that matches the search string.

Example: `cat tmp/lucius.INFO | lemctl findtrace -e "test packet #1:"`

This searches the file for the "test packet #1:" as a hex-encode string. It will
output the full trace for the first test packet.

### Interpreting Traces

```
I1119 18:05:15.964276       1 packet.go:207] "level"=3 "msg"="process current" "action"="<Type=ACTION_TYPE_LOOKUP;<Table=input-iface><onEvaluate=false><Hash=773172257>;>"
I1119 18:05:15.964279       1 packet.go:207] "level"=3 "msg"="exact table entry matched" "table"="input-iface" "entry"="<key=0000000000000035>;<actions=<Type=ACTION_TYPE_UPDATE;Field={false 0 PACKET_FIELD_NUM_INPUT_IFACE PACKET_HEADER_GROUP_UNSPECIFIED 0 0};Op=UPDATE_TYPE_SET;ByteArg=0000000000000867;FieldArg={false 0 PACKET_FIELD_NUM_UNSPECIFIED PACKET_HEADER_GROUP_UNSPECIFIED 0 0};BitCount=0;BitOffset=0<onEvaluate=false><Hash=3333249902>;>;<Type=ACTION_TYPE_FLOW_COUNTER;<FlowCounter=2151-in-counter><onEvaluate=false><Hash=3798394713>;>>;<timeout=0001-01-01 00:00:00 +0000 UTC>;"
I1119 18:05:15.964281       1 packet.go:207] "level"=3 "msg"="process result" "state"=1 "action"="<Type=ACTION_TYPE_UPDATE;Field={false 0 PACKET_FIELD_NUM_INPUT_IFACE PACKET_HEADER_GROUP_UNSPECIFIED 0 0};Op=UPDATE_TYPE_SET;ByteArg=0000000000000867;FieldArg={false 0 PACKET_FIELD_NUM_UNSPECIFIED PACKET_HEADER_GROUP_UNSPECIFIED 0 0};BitCount=0;BitOffset=0<onEvaluate=false><Hash=3333249902>;>;<Type=ACTION_TYPE_FLOW_COUNTER;<FlowCounter=2151-in-counter><onEvaluate=false><Hash=3798394713>;>"
I1119 18:05:15.964284       1 packet.go:207] "level"=3 "msg"="process current" "action"="<Type=ACTION_TYPE_UPDATE;Field={false 0 PACKET_FIELD_NUM_INPUT_IFACE PACKET_HEADER_GROUP_UNSPECIFIED 0 0};Op=UPDATE_TYPE_SET;ByteArg=0000000000000867;FieldArg={false 0 PACKET_FIELD_NUM_UNSPECIFIED PACKET_HEADER_GROUP_UNSPECIFIED 0 0};BitCount=0;BitOffset=0<onEvaluate=false><Hash=3333249902>;>"
I1119 18:05:15.964287       1 packet.go:207] "level"=3 "msg"="process result" "state"=1 "action"="Continue"
```

These traces provide detailed information about how a specific packet flows
through the Lucius forwarding pipeline.

#### Finding Trace for specific test case

The Lucius log file containing traces can be very large. In a dataplane test,
every packet has a unique payload string and the log records the full packet
frame as hex. Use the findtrace command on the unique payload string to find the
beginning of the relevant packet trace. Each packet trace starts with an input
message, indicating that a packet has been received and specifying the Ethernet
device it came from on the DUT (Device Under Test) switch.

#### Lucius Trace Format

The trace logs every step inside the forwarding pipeline. These are all the
operations Lucius takes on the packet. There are a few important actions logged:
lookup: This means Lucius is looking up the packet in one of its tables. The
trace will indicate which table is being searched. update: This means a field in
the packet, sometimes a metadata field (data associated with the packet but not
a real header), is being changed.

After a lookup, the trace indicates if an entry was matched (hit) or if no entry
was matched (miss). A hit means an entry in the table matches the packet's
header. The trace will show the entry that matched and the corresponding actions
defined for that entry. A miss means no entry in the table matched. Depending on
the table, this might lead to a default action, like dropping the packet.

Tracing the packet involves following the sequence of lookup and update actions
as listed in the log. Actions within a table might trigger additional lookups in
other tables, chaining the logic. Note: The trace uses Lucius OIDs, which are
not the same as those seen in gPINS logs and sairedic.rec. There is a table in
Redis that maps these, (RIDTOVID and VIDTORID). but entries might be deleted
after tests, making correlation difficult. It is often easier to rely on fields
within the packet itself like MAC addresses or IP addresses to correlate entries
between different logs or between the trace and the SAI redis data.

#### Key Tables in the Lucius Pipeline

1.  The pipeline involves multiple tables applied to packets received on a port:
2.  tun-term: IP and IP tunnel terminations (e.g., decapsulation).
3.  input-iface: Maps the input port to the L3 interface.
4.  ingress-vrf: Maps the VRF for the interface.
5.  preingress-table: Contains entries like pre-ingress ACL rules. This is where
    pre-ingress ACL rules usually are.
6.  my-mac-table: Determines if the packet's destination MAC is the router's
    MAC.
7.  fib-selector: Decides whether to look up in the IPv4 or IPv6 FIB based on
    the IP version.
8.  fib (fib-v4/fib-v6): The main routing table. It's a single table for all
    VRFs, organized by concatenating the VRF ID and the destination for prefix
    matching.
9.  ingress-table: Runs ingress ACL rules.
10. output-interface: Maps output interface to output port
11. egress-action-table: Egress ACL
12. neighbor: Used to find the next hop MAC address. Its key is typically the
    output interface ID and the next top IP address.
13. output-table: Makes the final decision on forwarding, dropping, or copying
    to the CPU.

#### Debugging Workflow Using Traces

When a test fails, look at the test result (e.g., in the sponge log) to see why
it failed (e.g., mismatched header fields, wrong packet count). For example, if
the test result shows a wrong source MAC address, this immediately tells you the
output interface is wrong, because the output interface sets the source MAC. Use
the unique packet payload to find the specific packet trace in the Lucius log.

**Other sources of errors**

Looking at the Lucius source code for the relevant API implementation (e.g.,
create neighbor) or table definitions to understand the exact logic and which
fields are supported. Lucius only supports a subset of fields for some API; if a
required field for matching the expected entry is not supported or used, it
could lead to matching the wrong entry. The Lucius logs might show details about
the entry being matched.
