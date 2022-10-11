#!/bin/bash
set -e
cd "$(dirname "$0")"

kne delete "$1" || true
cd ..
DOCKER_BUILDKIT=1 docker build . -f Dockerfile.lemming -t "lemming:latest"
kind load docker-image lemming:latest --name kne
cd -
kne create "$1"
