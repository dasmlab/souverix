#!/usr/bin/env bash
set -euo pipefail

# pushme.sh - Push component container to registry
app=icscf
local_tag="${1:-local}"
repo="ghcr.io/dasmlab"

if command -v podman &> /dev/null && [[ -z "${FORCE_DOCKER:-}" ]]; then
    RUNTIME=podman
else
    RUNTIME=docker
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

LASTBUILD_FILE="${SCRIPT_DIR}/.lastbuild"
LOCAL_IMAGE="${app}:${local_tag}"

if [[ -z "${GITHUB_TOKEN:-}" ]]; then
  echo "Error: GITHUB_TOKEN environment variable is required" >&2
  exit 1
fi

if ! ${RUNTIME} image inspect "${LOCAL_IMAGE}" &> /dev/null; then
  echo "Error: Local image ${LOCAL_IMAGE} not found" >&2
  exit 1
fi

if [[ -f "${LASTBUILD_FILE}" ]]; then
  CURRENT_VERSION="$(cat "${LASTBUILD_FILE}")"
else
  CURRENT_VERSION="0.0.0"
fi

IFS='.' read -r -a VERSION_PARTS <<< "${CURRENT_VERSION}"
MAJOR="${VERSION_PARTS[0]:-0}"
MINOR="${VERSION_PARTS[1]:-0}"
PATCH="${VERSION_PARTS[2]:-0}"

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
VERSION_IMAGE="${repo}/${app}:${NEW_VERSION}"
LATEST_IMAGE="${repo}/${app}:latest"

${RUNTIME} tag "${LOCAL_IMAGE}" "${VERSION_IMAGE}"
${RUNTIME} tag "${LOCAL_IMAGE}" "${LATEST_IMAGE}"

echo "${GITHUB_TOKEN}" | ${RUNTIME} login "${repo}" --username "${repo#ghcr.io/}" --password-stdin

${RUNTIME} push "${VERSION_IMAGE}"
${RUNTIME} push "${LATEST_IMAGE}"

echo "${NEW_VERSION}" > "${LASTBUILD_FILE}"
echo ":white_check_mark: Pushed ${VERSION_IMAGE} and ${LATEST_IMAGE}"
