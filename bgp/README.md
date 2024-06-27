# BGP Reconciler

## OpenConfig BGP Paths Supported by Lemming (as of 2024-06-27)

### Session Establishment

```text
/network-instances/network-instance/protocols/protocol/bgp/global/config/as
/network-instances/network-instance/protocols/protocol/bgp/global/config/router-id
/network-instances/network-instance/protocols/protocol/bgp/neighbors/neighbor/config/peer-as
/network-instances/network-instance/protocols/protocol/bgp/neighbors/neighbor/config/neighbor-address
/network-instances/network-instance/protocols/protocol/bgp/neighbors/neighbor/config/neighbor-port
/network-instances/network-instance/protocols/protocol/bgp/neighbors/neighbor/transport/config/local-address
```

### BGP Policy

```text
/routing-policy/defined-sets/prefix-sets/prefix-set/config/name
/routing-policy/defined-sets/prefix-sets/prefix-set/config/mode
/routing-policy/defined-sets/prefix-sets/prefix-set/state/name
/routing-policy/defined-sets/prefix-sets/prefix-set/state/mode
/routing-policy/defined-sets/prefix-sets/prefix-set/prefixes/prefix/ip-prefix
/routing-policy/defined-sets/prefix-sets/prefix-set/prefixes/prefix/masklength-range
/routing-policy/defined-sets/prefix-sets/prefix-set/prefixes/prefix/config/ip-prefix
/routing-policy/defined-sets/prefix-sets/prefix-set/prefixes/prefix/config/masklength-range
/routing-policy/defined-sets/prefix-sets/prefix-set/prefixes/prefix/state/ip-prefix
/routing-policy/defined-sets/prefix-sets/prefix-set/prefixes/prefix/state/masklength-range
/routing-policy/defined-sets/bgp-defined-sets/community-sets/community-set/config/community-set-name
/routing-policy/defined-sets/bgp-defined-sets/community-sets/community-set/config/community-member
/routing-policy/defined-sets/bgp-defined-sets/community-sets/community-set/state/community-set-name
/routing-policy/defined-sets/bgp-defined-sets/community-sets/community-set/state/community-member
/routing-policy/defined-sets/bgp-defined-sets/as-path-sets/as-path-set/config/as-path-set-name
/routing-policy/defined-sets/bgp-defined-sets/as-path-sets/as-path-set/config/as-path-set-member
/routing-policy/defined-sets/bgp-defined-sets/as-path-sets/as-path-set/state/as-path-set-name
/routing-policy/defined-sets/bgp-defined-sets/as-path-sets/as-path-set/state/as-path-set-member

/routing-policy/policy-definitions/policy-definition/statements/statement/config/name
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/match-prefix-set/config/prefix-set
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/match-prefix-set/config/match-set-options
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/config/med-eq
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/config/origin-eq
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/config/local-pref-eq
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/config/route-type
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/config/community-set
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/config/next-hop-in
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/community-count/config/operator
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/community-count/config/value
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/match-as-path-set/config/as-path-set
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/match-as-path-set/config/match-set-options
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/match-community-set/config/community-set
/routing-policy/policy-definitions/policy-definition/statements/statement/conditions/bgp-conditions/match-community-set/config/match-set-options
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/config/set-route-origin
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/config/set-local-pref
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/config/set-next-hop
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/config/set-med
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/set-as-path-prepend/config/repeat-n
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/set-as-path-prepend/config/asn
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/set-community/config/method
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/set-community/config/options
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/set-community/inline/config/communities
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/set-community/reference/config/community-set-ref (deprecated)
/routing-policy/policy-definitions/policy-definition/statements/statement/actions/bgp-actions/set-community/reference/config/community-set-refs
```

### BGP RIB

```text
/network-instances/network-instance/protocols/protocol/bgp/rib/attr-sets/attr-set/index
/network-instances/network-instance/protocols/protocol/bgp/rib/attr-sets/attr-set/state/index
/network-instances/network-instance/protocols/protocol/bgp/rib/attr-sets/attr-set/state/origin
/network-instances/network-instance/protocols/protocol/bgp/rib/attr-sets/attr-set/state/med
/network-instances/network-instance/protocols/protocol/bgp/rib/attr-sets/attr-set/state/local-pref
/network-instances/network-instance/protocols/protocol/bgp/rib/attr-sets/attr-set/as-path/as-segment/index
/network-instances/network-instance/protocols/protocol/bgp/rib/attr-sets/attr-set/as-path/as-segment/state/index
/network-instances/network-instance/protocols/protocol/bgp/rib/attr-sets/attr-set/as-path/as-segment/state/type
/network-instances/network-instance/protocols/protocol/bgp/rib/attr-sets/attr-set/as-path/as-segment/state/member
/network-instances/network-instance/protocols/protocol/bgp/rib/communities/community/index
/network-instances/network-instance/protocols/protocol/bgp/rib/communities/community/state/index
/network-instances/network-instance/protocols/protocol/bgp/rib/communities/community/state/community
```

## Design

For design see [Lemming and GoBGP](docs/gobgp_integration.md).
