#!/bin/bash
# Update README.md with current build status badges
# This script should be run as part of CI/CD pipeline

set -e

REPO_OWNER="${GITHUB_REPOSITORY_OWNER:-dasmlab}"
REPO_NAME="${GITHUB_REPOSITORY##*/}"
REPO_NAME="${REPO_NAME:-ims}"

# GitHub Actions badge URLs
BUILD_BADGE="https://github.com/${REPO_OWNER}/${REPO_NAME}/workflows/Build/badge.svg"
TEST_BADGE="https://github.com/${REPO_OWNER}/${REPO_NAME}/workflows/Test/badge.svg"
LINT_BADGE="https://github.com/${REPO_OWNER}/${REPO_NAME}/workflows/Lint/badge.svg"
COVERAGE_BADGE="https://codecov.io/gh/${REPO_OWNER}/${REPO_NAME}/branch/main/graph/badge.svg"

# Workflow URLs
BUILD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/actions/workflows/build.yml"
TEST_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/actions/workflows/test.yml"
LINT_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/actions/workflows/lint.yml"
COVERAGE_URL="https://codecov.io/gh/${REPO_OWNER}/${REPO_NAME}"

# Create status table
STATUS_TABLE="## Build Status

| Component | Build | Tests | Coverage | Lint |
|-----------|-------|-------|----------|------|
| **IMS Core** | [![Build](${BUILD_BADGE})](${BUILD_URL}) | [![Tests](${TEST_BADGE})](${TEST_URL}) | [![Coverage](${COVERAGE_BADGE})](${COVERAGE_URL}) | [![Lint](${LINT_BADGE})](${LINT_URL}) |
| **IBCF/SIG-GW** | [![Build](${BUILD_BADGE})](${BUILD_URL}) | [![Tests](${TEST_BADGE})](${TEST_URL}) | - | - |
| **STIR/SHAKEN** | [![Build](${BUILD_BADGE})](${BUILD_URL}) | [![Tests](${TEST_BADGE})](${TEST_URL}) | - | - |
| **LI/Emergency** | [![Build](${BUILD_BADGE})](${BUILD_URL}) | [![Tests](${TEST_BADGE})](${TEST_URL}) | - | - |

**Latest Build**: [View Details](https://github.com/${REPO_OWNER}/${REPO_NAME}/actions)

---"

# Update README.md
if [ -f README.md ]; then
    # Remove old status section if exists
    sed -i '/^## Build Status$/,/^---$/d' README.md
    
    # Insert new status section after title using a temporary file
    # This avoids sed escaping issues with multi-line variables
    TMP_FILE=$(mktemp)
    # Find the line with "# IMS Core" or "# Souverix" and insert after it
    if grep -q "^# Souverix" README.md; then
        TITLE_PATTERN="^# Souverix"
    elif grep -q "^# IMS Core" README.md; then
        TITLE_PATTERN="^# IMS Core"
    else
        TITLE_PATTERN="^#"
    fi
    
    awk -v table="${STATUS_TABLE}" -v pattern="${TITLE_PATTERN}" '
        BEGIN { inserted = 0 }
        $0 ~ pattern && !inserted {
            print $0
            print ""
            print table
            inserted = 1
            next
        }
        { print }
    ' README.md > "${TMP_FILE}" && mv "${TMP_FILE}" README.md
    
    echo "✅ Updated README.md with build status badges"
else
    echo "❌ README.md not found"
    exit 1
fi
