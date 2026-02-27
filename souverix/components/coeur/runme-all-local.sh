#!/usr/bin/env bash
set -euo pipefail

# runme-all-local.sh - Run all Coeur subcomponents locally

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

SUBS=(bgcf hss icscf mgcf pcscf scscf)
TAG="${1:-local}"

if command -v podman &> /dev/null && [[ -z "${FORCE_DOCKER:-}" ]]; then
    RUNTIME=podman
    NETWORK_ARG="--network host"
else
    RUNTIME=docker
    NETWORK_ARG=""
fi

echo "🚀 Starting all Coeur subcomponents locally..."
echo ""

for subcomp in "${SUBS[@]}"; do
    if [[ -d "$subcomp" ]]; then
        CONTAINER_NAME="coeur-${subcomp}-local"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "Starting $subcomp..."
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        
        ${RUNTIME} stop "${CONTAINER_NAME}" 2>/dev/null || true
        ${RUNTIME} rm -f "${CONTAINER_NAME}" 2>/dev/null || true
        
        ${RUNTIME} run -d \
            --name "${CONTAINER_NAME}" \
            --restart always \
            ${NETWORK_ARG} \
            -e LOG_LEVEL="${LOG_LEVEL:-info}" \
            "${subcomp}:${TAG}" || echo "⚠️  Failed to start ${subcomp}"
        
        echo ""
    fi
done

echo "✅ All subcomponents started!"
echo ""
echo "View logs: ${RUNTIME} logs -f coeur-<subcomponent>-local"
echo "Stop all: ${RUNTIME} stop coeur-*-local"
