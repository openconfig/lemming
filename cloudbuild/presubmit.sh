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

export PATH=${PATH}:/usr/local/go/bin
gopath=$(go env GOPATH)
export PATH=${PATH}:$gopath/bin

cd /tmp/workspace
kne deploy ~/kne-internal/deploy/kne/kind-bridge.yaml
make deploy itest
make clean
make deploy2 itest2
make clean
sleep 5
make deploy3 itest3
make clean
sleep 5
make deploy3 itest4
