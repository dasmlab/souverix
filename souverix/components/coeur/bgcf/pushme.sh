#!/usr/bin/env bash
set -euo pipefail

# pushme.sh - Push COMPONENT component container to registry with SemVer versioning
# Tags from :local to :latest and creates SemVer tag, then pushes both

app=bgcf
local_tag="local"
repo="ghcr.io/dasmlab"

# Detect container runtime
if command -v podman &> /dev/null && [[ -z "${FORCE_DOCKER:-}" ]]; then
    RUNTIME=podman
else
    RUNTIME=docker
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

LOCAL_IMAGE="${repo}/${app}:${local_tag}"
LASTBUILD_FILE="${SCRIPT_DIR}/.lastbuild"

if [[ -z "${GITHUB_TOKEN:-}" ]]; then
  echo "Error: GITHUB_TOKEN environment variable is required" >&2
  echo "Create a token at: https://github.com/settings/tokens" >&2
  echo "Required scope: write:packages" >&2
  exit 1
fi

# Check if local image exists
if ! ${RUNTIME} image inspect "${LOCAL_IMAGE}" &>/dev/null; then
  echo "Error: Local image ${LOCAL_IMAGE} not found" >&2
  echo "Run ./buildme.sh ${local_tag} first" >&2
  exit 1
fi

# Read or initialize version
if [[ -f "${LASTBUILD_FILE}" ]]; then
  CURRENT_VERSION="$(cat "${LASTBUILD_FILE}")"
else
  CURRENT_VERSION="0.0.0"
fi

# Parse version components
IFS='.' read -r -a VERSION_PARTS <<< "${CURRENT_VERSION}"
MAJOR="${VERSION_PARTS[0]:-0}"
MINOR="${VERSION_PARTS[1]:-0}"
PATCH="${VERSION_PARTS[2]:-0}"

# Bump patch version by default (can be overridden with env var)
if [[ -n "${BUMP_MAJOR:-}" ]]; then
  MAJOR=$((MAJOR + 1))
  MINOR=0
  PATCH=0
elif [[ -n "${BUMP_MINOR:-}" ]]; then
  MINOR=$((MINOR + 1))
  PATCH=0
else
  PATCH=$((PATCH + 1))
fi

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"

echo ":package: Pushing image:"
echo " Current version: ${CURRENT_VERSION}"
echo " New version: ${NEW_VERSION}"
echo " Source: ${LOCAL_IMAGE}"
echo ""

# Tag images: :latest and SemVer
LATEST_IMAGE="${repo}/${app}:latest"
VERSION_IMAGE="${repo}/${app}:${NEW_VERSION}"

echo ":label: Tagging images..."
${RUNTIME} tag "${LOCAL_IMAGE}" "${LATEST_IMAGE}"
${RUNTIME} tag "${LOCAL_IMAGE}" "${VERSION_IMAGE}"

echo ":key: Logging in to ${repo}..."
echo "${GITHUB_TOKEN}" | ${RUNTIME} login "${repo}" --username "${repo#ghcr.io/}" --password-stdin

echo ":arrow_up: Pushing ${LATEST_IMAGE}..."
${RUNTIME} push "${LATEST_IMAGE}"

echo ":arrow_up: Pushing ${VERSION_IMAGE}..."
${RUNTIME} push "${VERSION_IMAGE}"

# Save new version to .lastbuild
echo "${NEW_VERSION}" > "${LASTBUILD_FILE}"

echo ""
echo ":white_check_mark: Successfully pushed:"
echo "  ${LATEST_IMAGE}"
echo "  ${VERSION_IMAGE}"
echo ""
echo ":bookmark: Version ${NEW_VERSION} saved to .lastbuild"
