#!/bin/bash
set -e
cd "$(dirname "$0")"

if ! which kne
then
    go install github.com/openconfig/kne/kne_cli
    mv "$HOME/go/bin/kne_cli" "$HOME/go/bin/kne"
fi


DIR=$(pwd)
cat > config.yaml << EOF
topology: $DIR/$1
kubecfg: $HOME/.kube/config
cli: $GOPATH/bin/kne
username: foo
password: fake
EOF

cd -
