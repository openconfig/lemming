#!/bin/bash

set -xe

export PATH=${PATH}:/usr/local/go/bin
GOPATH=$(go env GOPATH)
export PATH=${PATH}:$GOPATH/bin

cd /tmp/workspace
# Make sure kne is the version we expect
go install github.com/openconfig/kne/kne_cli
mv "$HOME/go/bin/kne_cli" "$HOME/go/bin/kne"

kne deploy ~/kne-internal/deploy/kne/kind-bridge.yaml
make deploy itest
make clean
make deploy2 itest2