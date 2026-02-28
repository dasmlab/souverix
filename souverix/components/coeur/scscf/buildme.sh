#!/usr/bin/env bash
set -euo pipefail

# buildme.sh - Build COMPONENT component container image
# Supports both Docker and Podman

app=scscf
tag="${1:-local}"
repo="ghcr.io/dasmlab"

# Detect container runtime
if command -v podman &> /dev/null && [[ -z "${FORCE_DOCKER:-}" ]]; then
    RUNTIME=podman
    echo "🐳 Using Podman"
else
    RUNTIME=docker
    echo "🐳 Using Docker"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

# Build from souverix root so common module is available
SOUVERIX_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"
cd "${SOUVERIX_ROOT}"

echo "🔨 Building ${app}:${tag}..."
echo "   Build context: ${SOUVERIX_ROOT}"
echo "   Dockerfile: components/coeur/${app}/Dockerfile"
echo "   Common library: github.com/dasmlab/souverix/common (imported at compile time)"

VERSION="${VERSION:-dev}"
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build with staged cache mounts for faster rebuilds
${RUNTIME} build \
    --tag "${repo}/${app}:${tag}" \
    --file "components/coeur/${app}/Dockerfile" \
    --build-arg VERSION="${VERSION}" \
    --build-arg BUILD_TIME="${BUILD_TIME}" \
    --build-arg GIT_COMMIT="${GIT_COMMIT}" \
    --build-arg package=. \
    --progress=plain \
    .

echo "✅ Build complete: ${repo}/${app}:${tag}"
echo ""
echo "Push with: ./pushme.sh ${tag}"
