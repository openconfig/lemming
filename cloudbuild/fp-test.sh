#!/bin/bash

set -x

git clone https://github.com/openconfig/featureprofiles.git
cd featureprofiles || exit

rc=0
while read -r test_path; do
    echo "$test_path"
    kne_topology="$(pwd)/topologies/kne/openconfig/lemming/topology.textproto"
    kne create "${kne_topology}"
    if ! go test "./$test_path" -kne-topo "${kne_topology}" -alsologtostderr; then
      rc=$?
      kubectl cluster-info dump --output-directory "/tmp/cluster-log/${test_path/\//-}"  --namespaces openconfig-lemming
    fi
    kne delete "${kne_topology}"
    sleep 10 # Give namespace time to be deleted
done < ../fp-tests

exit "$rc"