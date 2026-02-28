#!/usr/bin/env bash
set -euo pipefail

# runme-all-local.sh - Start all Coeur subcomponents locally
# Starts all components in parallel and waits for all to complete

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

echo "🚀 Starting all Coeur subcomponents locally..."
echo ""

# Array to store PIDs
declare -a PIDS=()
declare -a COMPONENTS=()
FAILED=0
FAILED_COMPONENTS=()

# Find all subcomponent directories with runme-local.sh
for dir in */; do
    comp="${dir%/}"
    runme="${comp}/runme-local.sh"
    
    if [[ -f "${runme}" && -x "${runme}" ]]; then
        echo "Starting ${comp}..."
        (
            cd "${comp}"
            ./runme-local.sh
        ) &
        PIDS+=($!)
        COMPONENTS+=("${comp}")
    fi
done

if [[ ${#PIDS[@]} -eq 0 ]]; then
    echo "❌ No components found with runme-local.sh"
    exit 1
fi

echo ""
echo "⏳ Waiting for ${#PIDS[@]} components to start..."
echo ""

# Wait for all starts and collect exit codes
for i in "${!PIDS[@]}"; do
    pid="${PIDS[$i]}"
    comp="${COMPONENTS[$i]}"
    
    if wait "${pid}"; then
        echo "✅ ${comp} started successfully"
    else
        echo "❌ ${comp} failed to start"
        FAILED=1
        FAILED_COMPONENTS+=("${comp}")
    fi
done

echo ""

if [[ ${FAILED} -eq 0 ]]; then
    echo "☑ All subcomponents started!"
    echo ""
    echo "View logs: podman logs -f <component>-local-instance"
    echo "Stop all: ./stop-all.sh"
    exit 0
else
    echo "❌ Failed to start: ${FAILED_COMPONENTS[*]}"
    exit 1
fi
