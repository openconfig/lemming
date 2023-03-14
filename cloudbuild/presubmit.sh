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


set -xeE

printf "\n  apiServerPort: 6443" >> /kne-internal/kind/kind-no-cni.yaml

kne deploy /kne-internal/deploy/kne/kind-bridge.yaml || true
mkdir -p ~/.kube
kind get kubeconfig --internal --name kne > ~/.kube/config
docker network connect kind "$(cat /etc/hostname)"

kne deploy /kne-internal/deploy/kne/kind-bridge.yaml
make itest
