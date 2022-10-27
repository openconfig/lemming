#!/bin/bash

set -xe

cd /tmp/workspace
kne deploy ~/kne-internal/deploy/kne/kind-bridge.yaml
make deploy itest
make clean
make deploy2 itest2