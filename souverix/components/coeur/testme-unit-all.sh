#!/usr/bin/env bash
set -euo pipefail

# testme-unit-all.sh - Run unit tests for all Coeur subcomponents in parallel
# Usage: ./testme-unit-all.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

COMPONENTS=("bgcf" "hss" "icscf" "mgcf" "pcscf" "scscf")

echo "🧪 Running unit tests for all Coeur components..."
echo ""

# Create temporary directory for logs
TMPDIR=$(mktemp -d)
trap "rm -rf ${TMPDIR}" EXIT

# Start all tests in parallel
PIDS=()
for comp in "${COMPONENTS[@]}"; do
    if [ -f "${comp}/testme-unit.sh" ]; then
        echo "Starting unit test for ${comp}..."
        (
            "${comp}/testme-unit.sh" > "${TMPDIR}/${comp}.log" 2>&1
            echo $? > "${TMPDIR}/${comp}.exit"
        ) &
        PIDS+=($!)
    else
        echo "⚠️  ${comp}/testme-unit.sh not found, skipping"
    fi
done

# Wait for all tests to complete
echo ""
echo "Waiting for all tests to complete..."
for pid in "${PIDS[@]}"; do
    wait "${pid}"
done

# Show results
echo ""
echo "═══════════════════════════════════════════════════════════"
echo "📊 Unit Test Results Summary"
echo "═══════════════════════════════════════════════════════════"
echo ""

ALL_PASSED=true
PASSED_COUNT=0
FAILED_COUNT=0

for comp in "${COMPONENTS[@]}"; do
    if [ ! -f "${TMPDIR}/${comp}.log" ]; then
        continue
    fi
    
    EXIT_CODE=$(cat "${TMPDIR}/${comp}.exit" 2>/dev/null || echo "1")
    
    if [ "${EXIT_CODE}" == "0" ]; then
        echo "✅ ${comp}: PASSED"
        PASSED_COUNT=$((PASSED_COUNT + 1))
    else
        echo "❌ ${comp}: FAILED"
        FAILED_COUNT=$((FAILED_COUNT + 1))
        ALL_PASSED=false
        echo "   Last 5 lines:"
        tail -n 5 "${TMPDIR}/${comp}.log" | sed 's/^/   /'
    fi
done

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "Total: ${#COMPONENTS[@]} | Passed: ${PASSED_COUNT} | Failed: ${FAILED_COUNT}"
echo "═══════════════════════════════════════════════════════════"

if [ "${ALL_PASSED}" == "true" ]; then
    echo ""
    echo "✅ All unit tests passed!"
    exit 0
else
    echo ""
    echo "❌ Some unit tests failed"
    exit 1
fi
