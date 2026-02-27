#!/usr/bin/env bash
set -euo pipefail

# buildme.sh - Build Souverix Priorite container image
# Supports both Docker and Podman

app=priorite
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

# Get version info
VERSION="${VERSION:-dev}"
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build with staged cache mounts for faster rebuilds
${RUNTIME} build \
    --tag "${app}:${tag}" \
    --file Dockerfile.priorite \
    --build-arg COMPONENT=priorite \
    --build-arg VERSION="${VERSION}" \
    --build-arg BUILD_TIME="${BUILD_TIME}" \
    --build-arg GIT_COMMIT="${GIT_COMMIT}" \
    --progress=plain \
    .

echo ":white_check_mark: Build complete: ${app}:${tag}"
echo
echo "Run with: ./runme-local.sh"
echo "Or: ${RUNTIME} run -d --name ${app}-local ${app}:${tag}"
