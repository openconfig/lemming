#!/bin/bash

set -xe

export PATH=${PATH}:/usr/local/go/bin
gopath=$(go env GOPATH)
export PATH=${PATH}:$gopath/bin

cd /tmp/workspace
kne deploy ~/kne-internal/deploy/kne/kind-bridge.yaml
make deploy itest
make clean
make deploy2 itest2