#!/usr/bin/env bash
set -euo pipefail

# buildme.sh - Build IMS Core container image
# Supports both Docker and Podman

app=ims-core
tag="${1:-local}"

# Detect container runtime
if command -v podman &> /dev/null && [[ -z "${FORCE_DOCKER:-}" ]]; then
    RUNTIME=podman
    echo ":whale: Using Podman"
else
    RUNTIME=docker
    echo ":whale: Using Docker"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

echo ":hammer: Building ${app}:${tag}..."

# Build with staged cache mounts for faster rebuilds
${RUNTIME} build \
    --tag "${app}:${tag}" \
    --file Dockerfile \
    --progress=plain \
    .

echo ":white_check_mark: Build complete: ${app}:${tag}"
echo
echo "Run with: ./runme-local.sh"
echo "Or: ${RUNTIME} run -d --name ${app}-local ${app}:${tag}"
