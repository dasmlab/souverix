# Kaniko GitHub Actions Runner

This directory contains the setup for a GitHub Actions self-hosted runner with Kaniko support for building containers in Kubernetes/OpenShift environments.

## Overview

Kaniko is a daemonless container image builder that executes builds inside containers. It's ideal for Kubernetes environments where you don't want to run Docker-in-Docker or require privileged containers.

## Files

- **Dockerfile.github-runner**: Custom runner image based on `myoung34/github-runner` with Kaniko executor
- **deployment.yaml**: Kubernetes deployment for the runner (2026-prod-3)
- **secret.yaml**: Secret template for GitHub PAT (Personal Access Token)
- **buildme.sh**: Script to build the runner image locally
- **pushme.sh**: Script to push the runner image with SemVer versioning

## Setup

### 1. Create the Secret

The secret requires a GitHub Personal Access Token (PAT) with the following scopes:
- `repo` (full control)
- `workflow` (update GitHub Actions workflows)
- `write:packages` (push to GitHub Container Registry)

**Important**: The PAT is stored at `/home/dasm/gh_token` on the build machine. Update `secret.yaml` with the actual token before applying:

```bash
# Read token from file
TOKEN=$(cat /home/dasm/gh_token)

# Update secret.yaml (manually or via sed)
sed -i "s/REPLACE_WITH_PAT_FROM_\/home\/dasm\/gh_token/$TOKEN/" secret.yaml

# Apply secret
kubectl apply -f secret.yaml
```

### 2. Deploy the Runner

```bash
kubectl apply -f deployment.yaml
```

The runner will:
- Register as `2026-prod-3`
- Use labels: `self-hosted,linux,x64,ocp-sno,kaniko`
- Run as root (required for Kaniko)
- Use the image: `ghcr.io/dasmlab/dasmlab-ci-cd-kaniko-agent:latest`

### 3. Verify Runner Registration

Check GitHub Actions → Settings → Runners to see the runner registered.

## Building and Pushing the Runner Image

### Build Locally

```bash
cd k8s/github-runner/kaniko
export GITHUB_TOKEN=$(cat /home/dasm/gh_token)
./buildme.sh latest
```

### Push to Registry

```bash
export GITHUB_TOKEN=$(cat /home/dasm/gh_token)
./pushme.sh latest
```

The script will:
- Bump the SemVer version (patch by default)
- Tag the image with the new version and `latest`
- Push to `ghcr.io/dasmlab/dasmlab-ci-cd-kaniko-agent`

## Workflow Usage

In your GitHub Actions workflows, use the `kaniko` label to run jobs on this runner:

```yaml
jobs:
  build:
    runs-on: [self-hosted, kaniko]
    steps:
      - name: Build with Kaniko
        run: |
          /kaniko/executor \
            --context . \
            --dockerfile Dockerfile \
            --destination ghcr.io/user/image:tag
```

## Kaniko Configuration

Kaniko uses Docker config for authentication. The workflow sets up `/kaniko/.docker/config.json` with credentials from `GITHUB_TOKEN`.

## Differences from Docker/Buildah

- **No daemon**: Kaniko doesn't require a Docker daemon
- **No privileged mode**: Kaniko runs without privileged containers
- **Direct push**: Kaniko builds and pushes directly to registries
- **Kubernetes-native**: Designed for Kubernetes environments

## Troubleshooting

### Runner not registering

- Check the secret is correctly applied: `kubectl get secret github-runner-secrets -n github-runner`
- Verify the PAT has correct scopes
- Check runner pod logs: `kubectl logs -n github-runner -l app=github-actions-runner-kaniko`

### Kaniko build fails

- Verify `/kaniko/executor` exists: `kubectl exec -it <pod> -- ls -la /kaniko/executor`
- Check Docker config: `kubectl exec -it <pod> -- cat /kaniko/.docker/config.json`
- Ensure registry credentials are correct
