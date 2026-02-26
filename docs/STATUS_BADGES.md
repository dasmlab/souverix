# Build Status Badges

## Overview

The repository uses GitHub Actions workflows to automatically update build status badges in the README.md file. These badges show the current status of builds, tests, coverage, and linting.

## Badge Types

### Build Badge
Shows the status of the build workflow:
- ✅ Green: Build successful
- ❌ Red: Build failed
- ⏳ Yellow: Build in progress

### Test Badge
Shows the status of the test workflow:
- ✅ Green: All tests passing
- ❌ Red: Tests failing
- ⏳ Yellow: Tests running

### Coverage Badge
Shows code coverage percentage (via Codecov):
- Shows current coverage percentage
- Links to detailed coverage report

### Lint Badge
Shows the status of the lint workflow:
- ✅ Green: Linting passed
- ❌ Red: Linting failed
- ⏳ Yellow: Linting in progress

## Workflows

### Build Workflow
- **File**: `.github/workflows/build.yml`
- **Triggers**: Push to main/develop, PRs, manual
- **Actions**: Builds Go binary and Docker image

### Test Workflow
- **File**: `.github/workflows/test.yml`
- **Triggers**: Push to main/develop, PRs, manual
- **Actions**: Runs unit tests, generates coverage

### Lint Workflow
- **File**: `.github/workflows/lint.yml`
- **Triggers**: Push to main/develop, PRs, manual
- **Actions**: Runs golangci-lint and format checks

### Update README Workflow
- **File**: `.github/workflows/update-readme.yml`
- **Triggers**: After Build/Test/Lint complete, hourly schedule
- **Actions**: Updates README.md with current badge status

## Manual Update

To manually update the README badges:

```bash
./scripts/update-status.sh
```

This script:
1. Generates badge URLs based on repository name
2. Updates the status table in README.md
3. Commits changes (if run in CI/CD)

## Customization

### Repository Name
Set environment variables:
```bash
export GITHUB_REPOSITORY_OWNER=dasmlab
export GITHUB_REPOSITORY=ims
```

### Badge URLs
Badge URLs follow this pattern:
```
https://github.com/{owner}/{repo}/workflows/{workflow}/badge.svg
```

### Coverage Badge
Requires Codecov integration:
1. Sign up at https://codecov.io
2. Add repository
3. Get token and add as `CODECOV_TOKEN` secret
4. Badge URL: `https://codecov.io/gh/{owner}/{repo}/branch/{branch}/graph/badge.svg`

## Status Table Format

The status table in README.md shows:

| Component | Build | Tests | Coverage | Lint |
|-----------|-------|-------|----------|------|
| IMS Core | Badge | Badge | Badge | Badge |
| IBCF/SIG-GW | Badge | Badge | - | - |
| STIR/SHAKEN | Badge | Badge | - | - |
| LI/Emergency | Badge | Badge | - | - |

Each badge links to the corresponding workflow run.

## Troubleshooting

### Badges Not Updating
1. Check workflow runs are completing
2. Verify `update-readme.yml` workflow is enabled
3. Check repository permissions for GitHub Actions

### Badge Shows "Unknown"
- Workflow may not have run yet
- Check workflow file exists and is valid
- Verify workflow is enabled in repository settings

### Coverage Badge Missing
- Codecov integration may not be set up
- Check `CODECOV_TOKEN` secret is configured
- Verify Codecov is tracking the repository
