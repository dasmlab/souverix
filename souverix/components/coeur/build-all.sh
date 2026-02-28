#!/usr/bin/env bash
set -euo pipefail

# build-all.sh - Build all Coeur subcomponents in parallel
# Runs buildme.sh for each subcomponent concurrently and waits for all to complete

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

# Get tag from first argument or default to "latest"
TAG="${1:-latest}"

echo "🔨 Building all Coeur subcomponents in parallel..."
echo "   Tag: ${TAG}"
echo ""

# Array to store PIDs
declare -a PIDS=()
declare -a COMPONENTS=()

# Find all subcomponent directories with buildme.sh
for dir in */; do
    comp="${dir%/}"
    buildme="${comp}/buildme.sh"
    
    if [[ -f "${buildme}" && -x "${buildme}" ]]; then
        echo "🚀 Starting build for ${comp}..."
        (
            cd "${comp}"
            ./buildme.sh "${TAG}"
        ) &
        PIDS+=($!)
        COMPONENTS+=("${comp}")
    fi
done

if [[ ${#PIDS[@]} -eq 0 ]]; then
    echo "❌ No components found with buildme.sh"
    exit 1
fi

echo ""
echo "⏳ Waiting for ${#PIDS[@]} builds to complete..."
echo ""

# Wait for all builds and collect exit codes
FAILED=0
FAILED_COMPONENTS=()

for i in "${!PIDS[@]}"; do
    pid="${PIDS[$i]}"
    comp="${COMPONENTS[$i]}"
    
    if wait "${pid}"; then
        echo "✅ ${comp} build completed successfully"
    else
        echo "❌ ${comp} build failed"
        FAILED=1
        FAILED_COMPONENTS+=("${comp}")
    fi
done

echo ""

if [[ ${FAILED} -eq 0 ]]; then
    echo "✅ All builds completed successfully!"
    exit 0
else
    echo "❌ Build failed for: ${FAILED_COMPONENTS[*]}"
    exit 1
fi
