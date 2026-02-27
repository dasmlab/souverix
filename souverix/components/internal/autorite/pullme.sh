#!/usr/bin/env bash
set -euo pipefail

# pullme.sh - Pull Souverix Autorite container from registry

app=autorite
repo="ghcr.io/dasmlab"
tag="${1:-latest}"

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

IMAGE="${repo}/${app}:${tag}"

echo ":arrow_down: Pulling ${IMAGE}..."

${RUNTIME} pull "${IMAGE}"

# Tag as local for convenience
${RUNTIME} tag "${IMAGE}" "${app}:local"

echo ":white_check_mark: Successfully pulled and tagged as ${app}:local"
echo
echo "Run with: ./runme-local.sh"
