#!/bin/bash
set -e
cd "$(dirname "$0")"

if [[ ! -f $HOME/go/bin/kne ]]
then
    go install github.com/openconfig/kne/kne_cli
    mv $HOME/go/bin/kne_cli $HOME/go/bin/kne
fi

kne delete ./topo.pb.txt || true
cd ..
DOCKER_BUILDKIT=1 docker build . -t "lemming:latest"
kind load docker-image lemming:latest --name kne
cd -
kne create ./topo.pb.txt 
