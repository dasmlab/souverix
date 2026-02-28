#!/usr/bin/env bash
set -euo pipefail

# build-all.sh - Build all Coeur subcomponents in parallel
# Runs buildme.sh for each subcomponent concurrently and waits for all to complete

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

# Get tag from first argument or default to "latest"
TAG="${1:-local}"

echo "🔨 Building all Coeur subcomponents in parallel..."
echo "   Tag: ${TAG}"
echo ""

# Array to store PIDs
declare -a PIDS=()
declare -a COMPONENTS=()
declare -a LOG_FILES=()

# Find all subcomponent directories with buildme.sh
for dir in */; do
    comp="${dir%/}"
    buildme="${comp}/buildme.sh"
    
    if [[ -f "${buildme}" && -x "${buildme}" ]]; then
        echo "🚀 Starting build for ${comp}..."
        
        # Create log file for this build
        LOG_FILE="/tmp/build-${comp}-$$.log"
        LOG_FILES+=("${LOG_FILE}")
        
        (
            cd "${comp}"
            ./buildme.sh "${TAG}" > "${LOG_FILE}" 2>&1
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
    log_file="${LOG_FILES[$i]}"
    
    if wait "${pid}"; then
        echo "✅ ${comp} build completed successfully"
    else
        echo "❌ ${comp} build failed"
        FAILED=1
        FAILED_COMPONENTS+=("${comp}")
        # Show last few lines of failed build
        echo "   Last output:"
        tail -5 "${log_file}" | sed 's/^/   /' || true
    fi
done

# Clean up log files
for log_file in "${LOG_FILES[@]}"; do
    rm -f "${log_file}" 2>/dev/null || true
done

echo ""

if [[ ${FAILED} -eq 0 ]]; then
    echo "✅ All builds completed successfully!"
    exit 0
else
    echo "❌ Build failed for: ${FAILED_COMPONENTS[*]}"
    echo ""
    echo "💡 To see full build output, run: cd <component> && ./buildme.sh ${TAG}"
    exit 1
fi
