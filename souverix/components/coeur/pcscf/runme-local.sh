#!/usr/bin/env bash
set -euo pipefail

# runme-local.sh - Run pcscf component locally

COMPONENT="pcscf"
PORT="${PORT:-8081}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

echo "🚀 Starting ${COMPONENT} component locally..."
echo ""

# Build if binary doesn't exist
if [[ ! -f "app/${COMPONENT}" ]] && [[ ! -f "app/main" ]] && [[ ! -f "${COMPONENT}-local" ]]; then
    echo "📦 Building ${COMPONENT}..."
    cd app
    go build -o "../${COMPONENT}-local" ./main.go || {
        echo "❌ Build failed"
        exit 1
    }
    cd ..
fi

# Find binary
BINARY=""
if [[ -f "${COMPONENT}-local" ]]; then
    BINARY="${COMPONENT}-local"
elif [[ -f "app/${COMPONENT}" ]]; then
    BINARY="app/${COMPONENT}"
elif [[ -f "app/main" ]]; then
    BINARY="app/main"
else
    echo "❌ Binary not found. Run ./buildme.sh first or build manually"
    exit 1
fi

echo "✅ Starting ${COMPONENT} on port ${PORT}..."
echo "   Diagnostic endpoint: http://localhost:${PORT}/diag/health"
echo "   Local test endpoint: http://localhost:${PORT}/diag/local_test"
echo "   Press Ctrl+C to stop"
echo ""

# Run the component
PORT="${PORT}" exec "${BINARY}"
