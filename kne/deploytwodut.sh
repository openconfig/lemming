#!/bin/bash
set -e
cd "$(dirname "$0")"

kne delete ./twodut_topo.pb.txt || true
cd ..
DOCKER_BUILDKIT=1 docker build . -t "lemming:latest"
kind load docker-image lemming:latest --name kne
cd -
kne create ./twodut_topo.pb.txt 
