#!/usr/bin/env bash
set -euo pipefail

# buildme.sh - Build component container image
app=icscf
tag="${1:-local}"

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

VERSION="${VERSION:-dev}"
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

${RUNTIME} build \
    --tag "${app}:${tag}" \
    --file Dockerfile \
    --build-arg VERSION="${VERSION}" \
    --build-arg BUILD_TIME="${BUILD_TIME}" \
    --build-arg GIT_COMMIT="${GIT_COMMIT}" \
    --build-arg package=. \
    --progress=plain \
    .

echo ":white_check_mark: Build complete: ${app}:${tag}"
