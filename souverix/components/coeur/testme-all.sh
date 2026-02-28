#!/usr/bin/env bash
set -euo pipefail

# testme-all.sh - Test all Coeur subcomponents in parallel
# Runs testme-local.sh for each subcomponent concurrently and waits for all to complete

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

echo "🧪 Testing all Coeur subcomponents in parallel..."
echo ""

# Array to store PIDs
declare -a PIDS=()
declare -a COMPONENTS=()
declare -a LOG_FILES=()

# Find all subcomponent directories with testme-local.sh
for dir in */; do
    comp="${dir%/}"
    testme="${comp}/testme-local.sh"
    
    if [[ -f "${testme}" && -x "${testme}" ]]; then
        echo "🚀 Starting test for ${comp}..."
        
        # Create log file for this test
        LOG_FILE="/tmp/test-${comp}-$$.log"
        LOG_FILES+=("${LOG_FILE}")
        
        (
            cd "${comp}"
            ./testme-local.sh > "${LOG_FILE}" 2>&1
        ) &
        PIDS+=($!)
        COMPONENTS+=("${comp}")
    fi
done

if [[ ${#PIDS[@]} -eq 0 ]]; then
    echo "❌ No components found with testme-local.sh"
    exit 1
fi

echo ""
echo "⏳ Waiting for ${#PIDS[@]} tests to complete..."
echo ""

# Wait for all tests and collect exit codes
FAILED=0
FAILED_COMPONENTS=()

for i in "${!PIDS[@]}"; do
    pid="${PIDS[$i]}"
    comp="${COMPONENTS[$i]}"
    log_file="${LOG_FILES[$i]}"
    
    if wait "${pid}"; then
        echo "✅ ${comp} test passed"
    else
        echo "❌ ${comp} test failed"
        FAILED=1
        FAILED_COMPONENTS+=("${comp}")
        # Show last few lines of failed test
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
    echo "✅ All tests passed!"
    exit 0
else
    echo "❌ Test failed for: ${FAILED_COMPONENTS[*]}"
    echo ""
    echo "💡 To see full test output, run: cd <component> && ./testme-local.sh"
    exit 1
fi
