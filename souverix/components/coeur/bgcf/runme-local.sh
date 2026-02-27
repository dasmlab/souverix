#!/usr/bin/env bash
set -euo pipefail

# runme-local.sh - Run bgcf component container locally
# Runs the container built by buildme.sh and exposes ports for testing

COMPONENT="bgcf"
PORT="${PORT:-8084}"
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

# Calculate metrics port (main port + 1000)
METRICS_PORT=${METRICS_PORT:-9094}

# Run the container in daemon mode
${RUNTIME} run \
    -d \
    --name "${COMPONENT}-local" \
    --rm \
    -p "${PORT}:${PORT}" \
    -p "${METRICS_PORT}:${METRICS_PORT}" \
    -e "PORT=${PORT}" \
    -e "METRICS_PORT=${METRICS_PORT}" \
    "${IMAGE}"

echo ""
echo "✅ Container started in daemon mode"
echo "   Container name: ${COMPONENT}-local"
echo "   Main port: ${PORT}"
echo "   Metrics port: ${METRICS_PORT}"
echo ""
echo "View logs: ${RUNTIME} logs -f ${COMPONENT}-local"
echo "Stop container: ${RUNTIME} stop ${COMPONENT}-local"
