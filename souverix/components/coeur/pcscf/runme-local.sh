#!/usr/bin/env bash
set -euo pipefail

# runme-local.sh - Run pcscf component container locally
# Runs the container built by buildme.sh and exposes ports for testing

COMPONENT="pcscf"
PORT="${PORT:-8081}"
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

# Calculate out-of-band ports (common across all components)
METRICS_PORT="${METRICS_PORT:-9091}"    # Prometheus metrics
DIAG_PORT="${DIAG_PORT:-9081}"              # Diagnostics
TEST_PORT="${TEST_PORT:-9181}"              # Test endpoints

# Run the container in daemon mode
${RUNTIME} run \
    -d \
    --name "${COMPONENT}-local" \
    --rm \
    -p "${PORT}:${PORT}" \
    -p "${METRICS_PORT}:${METRICS_PORT}" \
    -p "${DIAG_PORT}:${DIAG_PORT}" \
    -p "${TEST_PORT}:${TEST_PORT}" \
    -e "PORT=${PORT}" \
    -e "METRICS_PORT=${METRICS_PORT}" \
    -e "DIAG_PORT=${DIAG_PORT}" \
    -e "TEST_PORT=${TEST_PORT}" \
    "${IMAGE}"

echo ""
echo "✅ Container started in daemon mode"
echo "   Container name: ${COMPONENT}-local"
echo ""
echo "📡 Servers (all out of band):"
echo "   Main:      http://localhost:${PORT}"
echo "   Metrics:   http://localhost:${METRICS_PORT}/metrics (Prometheus)"
echo "   Diagnostics: http://localhost:${DIAG_PORT}/diag/health"
echo "   Test:      http://localhost:${TEST_PORT}/test/local"
echo ""
echo "View logs: ${RUNTIME} logs -f ${COMPONENT}-local"
echo "Stop container: ${RUNTIME} stop ${COMPONENT}-local"
