#!/bin/bash
# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -xe

cat << EOF > ~/.bazelrc
build --remote_cache https://storage.googleapis.com/lemming-bazel-cache
build --google_default_credentials
EOF

export PATH=${PATH}:/usr/local/go/bin
gopath=$(go env GOPATH)
export PATH=${PATH}:$gopath/bin
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64 && \
sudo install skaffold /usr/local/bin/
curl -Lo bazel https://github.com/bazelbuild/bazelisk/releases/download/v1.16.0/bazelisk-linux-amd64 && \
sudo install bazel /usr/local/bin/

cd /tmp/workspace
kne deploy ~/kne-internal/deploy/kne/kind-bridge.yaml

make load-operator
kubectl set image -n lemming-operator deployment/lemming-controller-manager manager=us-west1-docker.pkg.dev/openconfig-lemming/release/operator:ga
make load
make itest
# Make sure we can populate this file.
stat ~/.ssh/google_compute_known_hosts || true
touch ~/.ssh/google_compute_known_hosts
