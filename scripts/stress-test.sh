#!/usr/bin/env bash
# ==============================================================================
# A Bank Mobile Wallet: Enterprise SRE Stress & Concurrency Test Runner
# Evaluates system stability, HTTP/2 multiplexing, ACID balance conservation,
# and idempotency replay deduplication under concurrent worker load.
# ==============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_PATH="${SCRIPT_DIR}/stress-test"

# Compile Go stress testing binary if source is newer or binary does not exist
if [[ ! -f "${BIN_PATH}" || "${SCRIPT_DIR}/stress-test.go" -nt "${BIN_PATH}" ]]; then
  echo "==> Compiling high-performance stress test runner from source..."
  (cd "${REPO_ROOT}" && go build -o "${BIN_PATH}" "${SCRIPT_DIR}/stress-test.go")
fi

# Forward all CLI arguments directly to the compiled Go binary
exec "${BIN_PATH}" "$@"
