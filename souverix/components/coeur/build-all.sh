#!/usr/bin/env bash
set -euo pipefail

# build-all.sh - Build all Coeur subcomponents

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

SUBS=(bgcf hss icscf mgcf pcscf scscf)
TAG="${1:-local}"

echo "🔨 Building all Coeur subcomponents..."
echo ""

for subcomp in "${SUBS[@]}"; do
    if [[ -d "$subcomp" ]] && [[ -f "$subcomp/buildme.sh" ]]; then
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "Building $subcomp..."
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        (cd "$subcomp" && ./buildme.sh "$TAG")
        echo ""
    else
        echo "⚠️  Skipping $subcomp (no buildme.sh found)"
    fi
done

echo "✅ All subcomponents built!"
