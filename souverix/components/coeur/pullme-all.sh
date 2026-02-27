#!/usr/bin/env bash
set -euo pipefail

# pullme-all.sh - Pull all Coeur subcomponents from registry

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

SUBS=(bgcf hss icscf mgcf pcscf scscf)
TAG="${1:-latest}"
repo="ghcr.io/dasmlab"

if command -v podman &> /dev/null && [[ -z "${FORCE_DOCKER:-}" ]]; then
    RUNTIME=podman
else
    RUNTIME=docker
fi

echo "⬇️  Pulling all Coeur subcomponents..."
echo ""

for subcomp in "${SUBS[@]}"; do
    if [[ -d "$subcomp" ]]; then
        IMAGE="${repo}/${subcomp}:${TAG}"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "Pulling $subcomp..."
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        ${RUNTIME} pull "${IMAGE}" || echo "⚠️  Failed to pull ${IMAGE}"
        ${RUNTIME} tag "${IMAGE}" "${subcomp}:local" 2>/dev/null || true
        echo ""
    fi
done

echo "✅ All subcomponents pulled!"
