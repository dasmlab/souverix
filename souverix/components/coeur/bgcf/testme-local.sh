#!/usr/bin/env bash
set -euo pipefail

# testme-local.sh - Run local tests for COMPONENT component
# Tests the /diag/local_test endpoint

COMPONENT="bgcf"
PORT="${PORT:-8086}"
DIAG_URL="http://localhost:${PORT}/diag/local_test"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

echo "🧪 Testing ${COMPONENT} component locally..."
echo ""

# Check if component is running
if ! curl -s -f "http://localhost:${PORT}/diag/health" > /dev/null 2>&1; then
    echo "❌ Component is not running on port ${PORT}"
    echo "   Start it with: ./runme-local.sh"
    exit 1
fi

echo "✅ Component is running"
echo ""

# Test /diag/local_test endpoint
echo "Testing /diag/local_test endpoint..."
RESPONSE=$(curl -s -w "\n%{http_code}" "${DIAG_URL}")

HTTP_CODE=$(echo "${RESPONSE}" | tail -n1)
BODY=$(echo "${RESPONSE}" | head -n-1)

if [[ "${HTTP_CODE}" == "200" ]]; then
    echo "✅ Test passed: HTTP ${HTTP_CODE}"
    echo "Response: ${BODY}"
    
    # Check if response contains "success"
    if echo "${BODY}" | grep -q '"resp".*"success"'; then
        echo "✅ Response contains 'success'"
        exit 0
    else
        echo "⚠️  Response does not contain 'success'"
        echo "Expected: {\"resp\": \"success\", ...}"
        exit 1
    fi
else
    echo "❌ Test failed: HTTP ${HTTP_CODE}"
    echo "Response: ${BODY}"
    exit 1
fi
