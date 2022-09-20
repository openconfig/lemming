#!/bin/bash
set -e
cd "$(dirname "$0")"

kne delete ./topo.pb.txt || true
cd ..
DOCKER_BUILDKIT=1 docker build . -t "lemming:latest"
kind load docker-image lemming:latest --name kne
cd -
kne create ./topo.pb.txt 
