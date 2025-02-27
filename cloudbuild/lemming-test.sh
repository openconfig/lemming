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

# shellcheck disable=SC2317
function dumpinfo {
    if [ -d "/tmp/cluster-log" ]; then
        gsutil cp -r -Z /tmp/cluster-log "gs://lemming-test-logs/$BUILD"
    fi
}

cat << EOF > ~/.bazelrc
build --remote_cache https://storage.googleapis.com/lemming-bazel-cache
build --google_default_credentials
EOF

export PATH=${PATH}:/usr/local/go/bin
gopath=$(go env GOPATH)
export PATH=${PATH}:$gopath/bin
curl -Lo bazel https://github.com/bazelbuild/bazelisk/releases/download/v1.16.0/bazelisk-linux-amd64 && \
sudo install bazel /usr/local/bin/
sudo apt-get -y install libpcap-dev libnl-genl-3-dev libnl-3-dev

cd /tmp/workspace
kne deploy ~/kne-internal/deploy/kne/kind-bridge.yaml

make load

set +e
rc=0
trap dumpinfo EXIT
trap 'rc=$?' ERR

make itest

exit "${rc}"