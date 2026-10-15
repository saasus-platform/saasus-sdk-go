#!/bin/bash

# Generate every module's client from the OpenAPI specs in SPEC_DIR.
#
# SPEC_DIR must contain one <module>.yml per module below.
#
# The generator version is pinned here so that regenerating always produces the
# same code for the same spec: a version bump would otherwise show up as a diff
# that looks like a spec change.

set -eu

SPEC_DIR="${SPEC_DIR:?SPEC_DIR is required}"
OAPI_CODEGEN_VERSION="v1.12.4"

MODULES="apilog auth awsmarketplace billing communication integration pricing"

GO_BIN_PATH="$(go env GOBIN)"
if [ -z "${GO_BIN_PATH}" ]; then
    GO_BIN_PATH="$(go env GOPATH)/bin"
fi
export PATH="${PATH}:${GO_BIN_PATH}"

go install "github.com/deepmap/oapi-codegen/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION}"

for module in ${MODULES}; do
    sh generate.sh "${module}" "${SPEC_DIR}/${module}.yml"
done
