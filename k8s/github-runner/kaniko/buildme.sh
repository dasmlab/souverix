#!/usr/bin/env bash
set -euo pipefail

# buildme.sh - Build Kaniko GitHub Runner container image
# Uses Docker/Podman to build the runner image with Kaniko support

app=dasmlab-ci-cd-kaniko-agent
tag="${1:-latest}"
repo="ghcr.io/dasmlab"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

# Detect container runtime
if command -v podman &> /dev/null && [[ -z "${FORCE_DOCKER:-}" ]]; then
    RUNTIME=podman
    echo ":whale: Using Podman"
else
    RUNTIME=docker
    echo ":whale: Using Docker"
fi

echo ":hammer: Building ${app}:${tag}..."

# Build the runner image with Kaniko
${RUNTIME} build \
    --tag "${app}:${tag}" \
    --tag "${repo}/${app}:${tag}" \
    --file Dockerfile.github-runner \
    --progress=plain \
    .

echo ":white_check_mark: Build complete: ${app}:${tag}"
echo
echo "Push with: ./pushme.sh ${tag}"
