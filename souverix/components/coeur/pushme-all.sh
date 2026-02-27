#!/usr/bin/env bash
set -euo pipefail

# pushme-all.sh - Push all Coeur subcomponents to registry

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

SUBS=(bgcf hss icscf mgcf pcscf scscf)
TAG="${1:-local}"

if [[ -z "${GITHUB_TOKEN:-}" ]]; then
    echo "Error: GITHUB_TOKEN environment variable is required" >&2
    exit 1
fi

echo "📦 Pushing all Coeur subcomponents..."
echo ""

for subcomp in "${SUBS[@]}"; do
    if [[ -d "$subcomp" ]] && [[ -f "$subcomp/pushme.sh" ]]; then
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "Pushing $subcomp..."
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        (cd "$subcomp" && ./pushme.sh "$TAG")
        echo ""
    else
        echo "⚠️  Skipping $subcomp (no pushme.sh found)"
    fi
done

echo "✅ All subcomponents pushed!"
