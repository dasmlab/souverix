#!/usr/bin/env bash
set -euo pipefail

# runme-local.sh - Run icscf component container locally
# Runs the container built by buildme.sh and exposes ports for testing

COMPONENT="icscf"
PORT="${PORT:-8082}"
IMAGE_TAG="${IMAGE_TAG:-local}"
IMAGE="ghcr.io/dasmlab/${COMPONENT}:${IMAGE_TAG}"

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

echo "🚀 Starting ${COMPONENT} container locally..."
echo ""

# Check if image exists
if ! ${RUNTIME} image inspect "${IMAGE}" &>/dev/null; then
    echo "❌ Container image ${IMAGE} not found"
    echo "   Build it first with: ./buildme.sh ${IMAGE_TAG}"
    exit 1
fi

echo "✅ Found container image: ${IMAGE}"
echo ""

# Stop any existing container
${RUNTIME} stop "${COMPONENT}-local" 2>/dev/null || true
${RUNTIME} rm "${COMPONENT}-local" 2>/dev/null || true

echo "📦 Starting container..."
echo "   Image: ${IMAGE}"
echo "   Port: ${PORT}"
echo "   Diagnostic endpoint: http://localhost:${PORT}/diag/health"
echo "   Local test endpoint: http://localhost:${PORT}/diag/local_test"
echo "   Press Ctrl+C to stop"
echo ""

# Run the container
${RUNTIME} run \
    --name "${COMPONENT}-local" \
    --rm \
    -p "${PORT}:${PORT}" \
    -e "PORT=${PORT}" \
    "${IMAGE}"
