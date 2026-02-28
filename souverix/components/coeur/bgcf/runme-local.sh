#!/usr/bin/env bash
set -euo pipefail

# runme-local.sh - Run bgcf component container locally
# Runs the container built by buildme.sh and exposes ports for testing

app="bgcf"
CONTAINER_NAME="${app}-local-instance"
IMAGE_TAG="local"
IMAGE="ghcr.io/dasmlab/${app}:${IMAGE_TAG}"

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

echo "🚀 Starting ${app} container locally..."
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
${RUNTIME} stop "${CONTAINER_NAME}" 2>/dev/null || true
${RUNTIME} rm "${CONTAINER_NAME}" 2>/dev/null || true

# Calculate out-of-band ports (common across all components)
PORT="${PORT:-8084}"                    # Main server
METRICS_PORT="${METRICS_PORT:-9094}"    # Prometheus metrics
DIAG_PORT="${DIAG_PORT:-9084}"              # Diagnostics
TEST_PORT="${TEST_PORT:-9184}"              # Test endpoints

# Run the container in daemon mode
${RUNTIME} run \
    -d \
    --name "${CONTAINER_NAME}" \
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
echo "   Container name: ${CONTAINER_NAME}"
echo ""
echo "📡 Servers (all out of band):"
echo "   Main:      http://localhost:${PORT}"
echo "   Metrics:   http://localhost:${METRICS_PORT}/metrics (Prometheus)"
echo "   Diagnostics: http://localhost:${DIAG_PORT}/diag/health"
echo "   Test:      http://localhost:${TEST_PORT}/test/local"
echo ""
echo "View logs: ${RUNTIME} logs -f ${CONTAINER_NAME}"
echo "Stop container: ${RUNTIME} stop ${CONTAINER_NAME}"
