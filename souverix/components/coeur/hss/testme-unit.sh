#!/usr/bin/env bash
set -euo pipefail

# testme-unit.sh - Run unit tests for hss component
# Tests the /diag/unit_test endpoint with call flow simulation

COMPONENT="hss"
DIAG_PORT="${DIAG_PORT:-9086}"
BASE_URL="http://localhost:${DIAG_PORT}"
FLOW_ID="${FLOW_ID:-}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

echo "🧪 Running unit tests for ${COMPONENT} component..."
echo ""

# Check if component is running
if ! curl -s -f "${BASE_URL}/diag/health" > /dev/null 2>&1; then
    echo "❌ Component diagnostics server is not running on port ${DIAG_PORT}"
    echo "   Start it with: ./runme-local.sh"
    exit 1
fi

echo "✅ Component is running"
echo ""

# Build query string
QUERY_PARAMS="base_url=${BASE_URL}"
if [[ -n "${FLOW_ID}" ]]; then
    QUERY_PARAMS="${QUERY_PARAMS}&flow_id=${FLOW_ID}"
fi

# Test /diag/unit_test endpoint
echo "Testing /diag/unit_test endpoint..."
echo "  URL: ${BASE_URL}/diag/unit_test?${QUERY_PARAMS}"
echo ""

RESPONSE=$(curl -s -w "\n%{http_code}" "${BASE_URL}/diag/unit_test?${QUERY_PARAMS}")

HTTP_CODE=$(echo "${RESPONSE}" | tail -n1)
BODY=$(echo "${RESPONSE}" | head -n-1)

if [[ "${HTTP_CODE}" == "200" ]]; then
    echo "✅ Unit test endpoint responded: HTTP ${HTTP_CODE}"
    echo ""
    
    # Parse JSON response
    if command -v jq &> /dev/null; then
        echo "📊 Test Results:"
        echo "${BODY}" | jq -r '
            "Component: " + .component,
            "Flow ID: " + .flow_id,
            "All Passed: " + (.all_passed | tostring),
            "Steps Executed: " + (.steps | length | tostring),
            "",
            "Verification Summary:",
            "  Total: " + (.verification.total | tostring),
            "  Passed: " + (.verification.passed | tostring),
            "  Failed: " + (.verification.failed | tostring),
            ""
        '
        
        # Check if all passed
        ALL_PASSED=$(echo "${BODY}" | jq -r '.all_passed')
        if [[ "${ALL_PASSED}" == "true" ]]; then
            echo "✅ All unit tests passed!"
            exit 0
        else
            echo "❌ Some unit tests failed"
            echo ""
            echo "Failed steps:"
            echo "${BODY}" | jq -r '.steps[] | select(.passed == false) | "  Step " + (.step | tostring) + ": " + .message'
            exit 1
        fi
    else
        # No jq, just show raw response
        echo "Response: ${BODY}"
        echo ""
        echo "⚠️  Install 'jq' for better output formatting"
        
        # Basic check for "all_passed": true
        if echo "${BODY}" | grep -q '"all_passed".*true'; then
            echo "✅ All unit tests passed!"
            exit 0
        else
            echo "❌ Some unit tests failed"
            exit 1
        fi
    fi
else
    echo "❌ Unit test failed: HTTP ${HTTP_CODE}"
    echo "Response: ${BODY}"
    exit 1
fi
