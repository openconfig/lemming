#!/bin/bash

git clone https://github.com/openconfig/featureprofiles.git
cd featureprofiles || exit

while read -r test_path; do
    kne_topology=$(metadata_kne_topology "${test_path}")
    echo "$kne_topology"
    echo "$test_path"
    kne create "${kne_topology}"
    go test "./$test_path" -kne-topo "$(pwd)/topologies/kne/openconfig/lemming/topology.textproto " -alsologtostderr
    kne delete "$kne_topology"
done < fp-tests
