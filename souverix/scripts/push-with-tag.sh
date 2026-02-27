#!/usr/bin/env bash
set -euo pipefail

# push-with-tag.sh - Commit, tag, and push with SemVer
# Usage: ./scripts/push-with-tag.sh [commit message] [--push-tag]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}/.."

COMMIT_MSG="${1:-chore: update}"
PUSH_TAG="${2:-}"

# Stage all changes
echo "📦 Staging changes..."
git add -A

# Check if there are changes to commit
if git diff --staged --quiet; then
    echo "ℹ️  No changes to commit"
else
    # Commit changes
    echo "💾 Committing changes..."
    git commit -m "${COMMIT_MSG}"
fi

# Create and push tag
echo "🏷️  Creating version tag..."
export PUSH_TAG="yes"
"${SCRIPT_DIR}/tag-version.sh" --push

# Push commits
echo "📤 Pushing commits..."
git push origin "$(git branch --show-current)"

echo ""
echo "✅ All done! Commit and tag pushed."
