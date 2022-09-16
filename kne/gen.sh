#!/bin/bash
cd "$(dirname "$0")"

DIR=$(pwd)
cat > config.yaml << EOF
topology: $DIR/topo.pb.txt
kubecfg: $HOME/.kube/config
cli: $HOME/go/bin/kne
username: foo
password: fake
EOF

cd -