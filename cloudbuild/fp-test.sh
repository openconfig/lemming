#!/bin/bash


git clone https://github.com/openconfig/featureprofiles.git
cd featureprofiles

function metadata_kne_topology() {
  local metadata_test_path
  metadata_test_path="${1}"
  local topo_prefix
  topo_prefix=$(echo "${platform}" | tr "_" "/")
  declare -A kne_topology_file
  kne_topology_file["TESTBED_DUT"]="${topo_prefix}/dut.textproto"
  kne_topology_file["TESTBED_DUT_DUT_4LINKS"]="${topo_prefix}/dutdut.textproto"
  kne_topology_file["TESTBED_DUT_ATE_2LINKS"]="${topo_prefix}/dutate.textproto"
  kne_topology_file["TESTBED_DUT_ATE_4LINKS"]="${topo_prefix}/dutate.textproto"
  kne_topology_file["TESTBED_DUT_ATE_9LINKS_LAG"]="${topo_prefix}/dutate_lag.textproto"
  for p in "${!kne_topology_file[@]}"; do
    if grep -q "testbed.*${p}$" "${metadata_test_path}"/metadata.textproto; then
      echo "${kne_topology_file[${p}]}"
      return
    fi
  done
  echo "UNKNOWN"
}

platform="topologies/kne/openconfig/lemming"

for test_path in $(cat fp-tests); do
    kne_topology=$(metadata_kne_topology "${test_path}")
    echo "$kne_topology"
    echo "$test_paths"
    kne create "${kne_topology}"
    go test "./$test_path" -kne-topo "$(pwd)/${kne_topology}" -alsologtostderr
    kne delete $kne_topology
done
