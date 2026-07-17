#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

TMP_DIR="$(mktemp -d)"
cleanup() {
	rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

SOURCE_REPO="https://github.com/magodo/azure-rest-api-cov-terraform-reports"
SOURCE_DIR="${TMP_DIR}/azure-rest-api-cov-terraform-reports"

git clone --depth 1 "${SOURCE_REPO}" "${SOURCE_DIR}"

SOURCE_SHA="$(git -C "${SOURCE_DIR}" rev-parse HEAD)"

cp "${SOURCE_DIR}/tf.json" "${REPO_ROOT}/azurerm/coverage/tf.json"

echo "Updated ./azurerm/coverage/tf.json from magodo/azure-rest-api-cov-terraform-reports @ ${SOURCE_SHA}"
echo ""
echo "Run:"
echo "git add ./azurerm/coverage/tf.json && git commit -m \"./azurerm/coverage/tf.json updated to magodo/azure-rest-api-cov-terraform-reports commit ${SOURCE_SHA}\""
