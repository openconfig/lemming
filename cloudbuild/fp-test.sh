#!/bin/bash

git clone https://github.com/openconfig/featureprofiles.git
cd featureprofiles || exit

rc=0
while read -r test_path; do
    kne_topology=$(metadata_kne_topology "${test_path}")
    echo "$kne_topology"
    echo "$test_path"
    kne create "${kne_topology}"
    go test "./$test_path" -kne-topo "$(pwd)/topologies/kne/openconfig/lemming/topology.textproto " -alsologtostderr
    if [[ $? -ne 0 ]]; then
      rc=$?
      kubectl cluster-info dump --output-directory "/tmp/cluster-log/${test_path/\//-}"  --namespaces openconfig-lemming
    fi
    kne delete "$kne_topology"
done < fp-tests

exit "$rc"