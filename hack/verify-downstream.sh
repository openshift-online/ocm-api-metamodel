#!/usr/bin/env bash
#
# Copyright (c) 2026 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -o errexit
set -o nounset
set -o pipefail

TARGET="${1:-}"
if [[ "${TARGET}" != "model" && "${TARGET}" != "sdk" ]]; then
	echo "usage: $0 <model|sdk>" >&2
	exit 1
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
METAMODEL="${ROOT}/metamodel"
WORK_DIR="$(mktemp -d -t verify-downstream-XXXX)"
trap 'rm -rf "${WORK_DIR}"' EXIT

if [[ ! -x "${METAMODEL}" ]]; then
	echo "metamodel binary not found at ${METAMODEL}; run 'make binary' first" >&2
	exit 1
fi

install_goimports() {
	local bin_path="${1}"
	GOBIN="${bin_path}" go install golang.org/x/tools/cmd/goimports@v0.4.0
	export PATH="${bin_path}:${PATH}"
}

require_ref() {
	local ref_name="${1}"
	if [[ -z "${!ref_name:-}" ]]; then
		echo "error: ${ref_name} must be set to a tag or branch ref" >&2
		exit 1
	fi
}

clone_downstream() {
	local url="${1}"
	local repo_dir="${2}"
	local ref_name="${3}"

	require_ref "${ref_name}"
	git clone --depth=1 --branch "${!ref_name}" "${url}" "${repo_dir}"
}

verify_model() {
	local repo_dir="${WORK_DIR}/ocm-api-model"
	local bin_path="${WORK_DIR}/bin"

	clone_downstream \
		"https://github.com/openshift-online/ocm-api-model.git" \
		"${repo_dir}" \
		OCM_API_MODEL_REF
	install_goimports "${bin_path}"
	cp "${METAMODEL}" "${repo_dir}/metamodel_generator/metamodel"
	chmod +x "${repo_dir}/metamodel_generator/metamodel"

	(
		cd "${repo_dir}"
		make check
		make clientapi
		cd clientapi
		go build ./...
	)
}

verify_sdk() {
	local repo_dir="${WORK_DIR}/ocm-sdk-go"
	local bin_path="${WORK_DIR}/bin"

	clone_downstream \
		"https://github.com/openshift-online/ocm-sdk-go.git" \
		"${repo_dir}" \
		OCM_SDK_GO_REF
	install_goimports "${bin_path}"
	cp "${METAMODEL}" "${repo_dir}/metamodel_generator/metamodel"
	chmod +x "${repo_dir}/metamodel_generator/metamodel"

	(
		cd "${repo_dir}"
		make model goimports-install
		hack/generate-client.sh "${repo_dir}/metamodel_generator/metamodel" .
		go build ./...
	)
}

case "${TARGET}" in
model)
	verify_model
	;;
sdk)
	verify_sdk
	;;
esac

echo "downstream ${TARGET} verification succeeded"
