#!/usr/bin/env bash
set -euo pipefail

# tag-version.sh - Create and push SemVer git tag for current commit
# Format: 0.0.X-alpha

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}/.."

VERSION_FILE="${SCRIPT_DIR}/../.version"
MAJOR=0
MINOR=0
PATCH=1
PRERELEASE="alpha"

# Get latest tag matching pattern 0.0.X-alpha
LATEST_TAG=$(git tag --sort=-version:refname | grep -E "^0\.0\.[0-9]+-alpha$" | head -1 || echo "")

if [[ -n "${LATEST_TAG}" ]]; then
    # Extract patch version from latest tag (e.g., 0.0.5-alpha -> 5)
    PATCH=$(echo "${LATEST_TAG}" | sed -E 's/^0\.0\.([0-9]+)-alpha$/\1/')
    PATCH=$((PATCH + 1))
fi

# If .version file exists, use it (allows manual override)
if [[ -f "${VERSION_FILE}" ]]; then
    VERSION=$(cat "${VERSION_FILE}")
    # Parse version if it's in format 0.0.X-alpha
    if [[ "${VERSION}" =~ ^0\.0\.([0-9]+)-alpha$ ]]; then
        PATCH="${BASH_REMATCH[1]}"
        # Increment if not explicitly set to keep this version
        if [[ -z "${KEEP_VERSION:-}" ]]; then
            PATCH=$((PATCH + 1))
        fi
    fi
fi

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}-${PRERELEASE}"
NEW_TAG="v${NEW_VERSION}"

# Check if tag already exists
if git rev-parse "${NEW_TAG}" >/dev/null 2>&1; then
    echo "⚠️  Tag ${NEW_TAG} already exists, skipping..."
    exit 0
fi

# Get current commit SHA
CURRENT_SHA=$(git rev-parse HEAD)
SHORT_SHA=$(git rev-parse --short HEAD)

# Create annotated tag
echo "🏷️  Creating tag: ${NEW_TAG}"
echo "   Commit: ${SHORT_SHA}"
echo "   Version: ${NEW_VERSION}"

git tag -a "${NEW_TAG}" -m "Release ${NEW_VERSION}

Commit: ${SHORT_SHA}
Date: $(date -u +%Y-%m-%dT%H:%M:%SZ)"

# Save version to file
echo "${NEW_VERSION}" > "${VERSION_FILE}"

# Push tag if requested or if PUSH_TAG is set
if [[ "${1:-}" == "--push" ]] || [[ -n "${PUSH_TAG:-}" ]]; then
    echo "📤 Pushing tag ${NEW_TAG} to origin..."
    git push origin "${NEW_TAG}"
    echo "✅ Tag ${NEW_TAG} pushed successfully"
else
    echo "💡 To push the tag, run: git push origin ${NEW_TAG}"
    echo "   Or use: ./scripts/tag-version.sh --push"
fi

echo ""
echo "✅ Tag created: ${NEW_TAG}"
echo "📝 Version saved to: ${VERSION_FILE}"
