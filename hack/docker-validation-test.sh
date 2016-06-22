#!/bin/bash

# Copyright 2016 The Kubernetes Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

KUBE_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"
#kube::golang::setup_env

ginkgo=$(kube::util::find-binary "ginkgo")
if [[ -z "${ginkgo}" ]]; then
  echo "You do not appear to have ginkgo built. Try 'hack/build-go.sh github.com/onsi/ginkgo/ginkgo'"
  exit 1
fi

#sudo -v
#"${ginkgo}" "${KUBE_ROOT}/test/docker_validation/conformance"
"${ginkgo}" "${KUBE_ROOT}/test/docker_validation/performance"

# Provided for backwards compatibility
#focus="Conformance"
#sudo -v
#"${ginkgo}" --focus=$focus "${KUBE_ROOT}/test/e2e_node/" -- --alsologtostderr --v 2 --node-name $(hostname) --build-services=true --start-services=true --stop-services=true


exit $?
