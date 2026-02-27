#!/usr/bin/env bash
set -euo pipefail

# buildme.sh - Build BGCF component container image
# Supports both Docker and Podman

app=bgcf
tag="${1:-latest}"
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

echo "🔨 Building ${app}:${tag}..."

# Build with staged cache mounts for faster rebuilds
${RUNTIME} build \
    --tag "${repo}/${app}:${tag}" \
    --file Dockerfile \
    --progress=plain \
    .

echo "✅ Build complete: ${repo}/${app}:${tag}"
echo ""
echo "Push with: ./pushme.sh ${tag}"
