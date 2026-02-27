#!/usr/bin/env bash
set -euo pipefail

# runme-local.sh - Run IMS Core locally in container

app=ims-core
tag="${1:-local}"

# Detect container runtime
if command -v podman &> /dev/null && [[ -z "${FORCE_DOCKER:-}" ]]; then
    RUNTIME=podman
    NETWORK_ARG="--network host"
else
    RUNTIME=docker
    NETWORK_ARG="-p 5060:5060 -p 5061:5061 -p 8080:8080 -p 9443:9443"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

# Stop and remove existing container
${RUNTIME} stop "${app}-local-instance" 2>/dev/null || true
${RUNTIME} rm -f "${app}-local-instance" 2>/dev/null || true

echo ":rocket: Starting ${app} locally..."

${RUNTIME} run -d \
    --name "${app}-local-instance" \
    --restart always \
    ${NETWORK_ARG} \
    -e LOG_LEVEL="${LOG_LEVEL:-info}" \
    -e ZERO_TRUST_MODE="${ZERO_TRUST_MODE:-false}" \
    "${app}:${tag}"

echo ":white_check_mark: Container started: ${app}-local-instance"
echo
echo "View logs: ${RUNTIME} logs -f ${app}-local-instance"
echo "Stop: ${RUNTIME} stop ${app}-local-instance"
