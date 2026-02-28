#!/usr/bin/env bash
set -euo pipefail

# stop-all.sh - Stop all Coeur subcomponent containers
# Stops all running component containers in parallel

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

echo "🛑 Stopping all Coeur subcomponent containers..."
echo ""

# Array to store PIDs
declare -a PIDS=()
declare -a COMPONENTS=()

# Find all subcomponent directories with stopme.sh
for dir in */; do
    comp="${dir%/}"
    stopme="${comp}/stopme.sh"
    
    if [[ -f "${stopme}" && -x "${stopme}" ]]; then
        echo "🛑 Stopping ${comp}..."
        (
            cd "${comp}"
            ./stopme.sh
        ) &
        PIDS+=($!)
        COMPONENTS+=("${comp}")
    fi
done

if [[ ${#PIDS[@]} -eq 0 ]]; then
    echo "ℹ️  No components found with stopme.sh"
    exit 0
fi

# Wait for all stops to complete
for i in "${!PIDS[@]}"; do
    pid="${PIDS[$i]}"
    wait "${pid}" || true  # Don't fail if stop fails (container might not be running)
done

echo ""
echo "✅ All stop operations completed"
