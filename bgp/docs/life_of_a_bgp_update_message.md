# GoBGP: Life of a BGP UPDATE Message for Policy Evaluation

This document captures some important moments during the life of BGP Update
message in GoBGP where relevant to policy evaluations.

This investigation was conducted during May-August, 2023

## BGP Message Policy Evaluation Call Chain

The hierarchy below indicate call chain locations.

**In bold** are where NLRIs are added to tables.

1.  [server.go](server.go): addNeighbor
    1.  newPeer
    2.  s.addIncoming(peer.fsm.incomingCh)
        1.  This hooks the UPDATE messages as received by the FSM into the BGP
            server via the channel, which BgpServer.Serve() reads.
    3.  Serve(): handlefsmMsg(e)
2.  [server.go](server.go): handlefsmMsg()
    1.  pathList, eor, notification := peer.handleUpdate(e)
        1.  [peer.go](peer.go): handleUpdate()
            1.  **peer.adjRibIn.Update(e.PathList)**
    2.  s.propagateUpdate(peer, pathList)
        1.  if p := s.policy.ApplyPolicy(tableId, table.POLICY_DIRECTION_IMPORT,
            path, policyOptions); p != nil {
            1.  ***IMPORT policies are applied***
        2.  if dsts := **rib.Update(path)**; len(dsts) > 0 {
            1.  This populates the Global RIB. The Global RIB is the same as the
                local RIB in GoBGP. The only difference is that when querying
                for local RIB routes, the routes are first filtered by
                `path.GetSource().Address.String() == id` in
                table/[destination.go](destination.go)’s rsFilter().
                1.  This comes from [server.go’s](server.go’s) ListPath() ->
                    getRib() -> table.Select() -> destination.Select() ->
                    GetKnownPathList() -> rsFilter().
            2.  However, this is actually wrong because getRib() always uses
                `id := table.GLOBAL_RIB_NAME`, which means that it’s always
                querying from the global RIB.
        3.  propagateUpdateToNeighbors()
            1.  Here remember that dsts is what’s going to be sent since
                rib.Update(path) computes the BGP UPDATE messages that need to
                be sent.
            2.  s.processOutgoingPaths(targetPeer, bestList, oldList);
                1.  calls s.filterPath() to ***apply EXPORT policies*** to get
                    the final outgoing messages.
            3.  sendfsmOutgoingMsg() places the BGP UPDATE messages in the
                infinite channel of the FSM back for sending out.

## RIB Information

GoBGP only stores the per-neighbor AdjRIB table and the Global RIB.

*   [server.go](server.go): ListPath()
    *   getAdjRib()
        *   Here you will notice that adj-rib-in-post comes from applying the
            IMPORT policies to the stored AdjRib.
        *   adj-rib-out-pre comes from filtering the peer’s “local RIB” via
            s.getBestFromLocal() -> s.getPossibleBest(peer, family)
            *   However, the “local RIB” is actually the global RIB because in
                [server.go’s](server.go’s) addNeighbor(), we have peer :=
                newPeer(&s.bgpConfig.Global, c, rib, s.policy, s.logger), and
                when the server is not a route server, we have rib :=
                s.globalRib.
        *   adj-rib-out-post comes from applying EXPORT policies after getting
            essentially adj-rib-out-pre using steps similar to the above.

## Miscellaneous

*   table/[policy.go](policy.go): ApplyPolicy returns true if the route was
    accepted, and nil if the route was rejected.
*   Searching for “RouteServerClient” will reveal that local policies and local
    RIBs will be used if this is true for a neighbour, and the global RIB and
    its own policies will be used if this is false.
*   table/[destination.go](destination.go): contains BGP best path selection
    logic.
