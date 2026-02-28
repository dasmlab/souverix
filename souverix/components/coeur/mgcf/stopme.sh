#!/usr/bin/env bash
set -euo pipefail

# stopme.sh - Stop mgcf component container

app="mgcf"
CONTAINER_NAME="${app}-local-instance"

# Detect container runtime
if command -v podman &> /dev/null && [[ -z "${FORCE_DOCKER:-}" ]]; then
    RUNTIME=podman
else
    RUNTIME=docker
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

echo "🛑 Stopping ${app} container..."

if ${RUNTIME} ps -a --format "{{.Names}}" | grep -q "^$"; then
    ${RUNTIME} stop "${CONTAINER_NAME}" 2>/dev/null || true
    echo "✅ Container ${CONTAINER_NAME} stopped"
else
    echo "ℹ️  Container ${CONTAINER_NAME} not running"
fi
