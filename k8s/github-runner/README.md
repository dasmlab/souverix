# GitHub Actions Runner Deployment for OpenShift

This directory contains Kubernetes manifests for deploying GitHub Actions self-hosted runners on OpenShift Container Platform (OCP) 2.21 SNO (Single Node OpenShift).

## Architecture

- **Platform**: OpenShift Container Platform 2.21 SNO
- **Host**: 28GB RAM / 24 CPU cores
- **Runner Image**: `ghcr.io/dasmlab/ci-cd-github-runner:latest`
- **Build System**: Buildah/Podman (no Docker socket required)
- **Security**: Uses privileged SCC for Buildah operations

## Prerequisites

1. OpenShift cluster access (`oc` CLI configured)
2. Cluster admin permissions (for SCC binding)
3. GitHub organization/repository registration token
4. Image pull secret for `ghcr.io/dasmlab`

## Deployment Steps

### 1. Create the Secret

Copy the template and add your GitHub token:

```bash
cd k8s/github-runner
cp secret.yaml.template secret.yaml
# Edit secret.yaml and replace YOUR_REGISTRATION_TOKEN_HERE with your actual token
```

**Note**: Registration tokens expire after 1 hour. If you need to re-register, get a new token from:
- GitHub → Organization Settings → Actions → Runners → New self-hosted runner

### 2. Create Image Pull Secret

The deployment references `dasmlab-ghcr-pull` secret. Create it if it doesn't exist:

```bash
# Using the script from infra repo (if available)
# Or manually:
oc create secret docker-registry dasmlab-ghcr-pull \
  --docker-server=ghcr.io \
  --docker-username=YOUR_GITHUB_USERNAME \
  --docker-password=YOUR_GITHUB_TOKEN \
  --namespace=github-runner
```

### 3. Deploy All Manifests

Deploy in order:

```bash
oc apply -f namespace.yaml
oc apply -f serviceaccount.yaml
oc apply -f scc-binding.yaml
oc apply -f secret.yaml
oc apply -f pvc.yaml
oc apply -f deployment.yaml
```

Or deploy all at once:

```bash
oc apply -f .
```

### 4. Verify Deployment

Check pod status:

```bash
oc get pods -n github-runner
oc logs -f deployment/github-actions-runner -n github-runner
```

Check runner in GitHub:
- Go to: https://github.com/dasmlab → Settings → Actions → Runners
- You should see `2026-prod-1` runner listed

## Configuration

### Runner Labels

The runner is configured with labels: `self-hosted,linux,x64,ocp-sno`

Use in workflows:
```yaml
jobs:
  build:
    runs-on: self-hosted
```

### Resource Limits

- **Requests**: 2Gi memory, 1000m CPU
- **Limits**: 8Gi memory, 4000m CPU

Adjust in `deployment.yaml` based on your cluster capacity.

### Runner Name

Default: `2026-prod-1` (OCP SNO hostname)

Change in `deployment.yaml`:
```yaml
- name: RUNNER_NAME
  value: "your-custom-name"
```

## Security Context Constraints (SCC)

The deployment uses the `privileged` SCC to allow:
- Running as root (required for Buildah)
- Privileged containers (required for container builds)

The `scc-binding.yaml` creates the necessary RBAC to bind the ServiceAccount to the privileged SCC.

## Troubleshooting

### Pod Not Starting

1. Check SCC binding:
   ```bash
   oc describe rolebinding github-runner-scc-binding -n github-runner
   ```

2. Check pod events:
   ```bash
   oc describe pod -l app=github-actions-runner -n github-runner
   ```

### Runner Not Appearing in GitHub

1. Check logs for token errors:
   ```bash
   oc logs -f deployment/github-actions-runner -n github-runner
   ```

2. Verify token is valid (tokens expire after 1 hour)

3. Check network connectivity:
   ```bash
   oc exec -it deployment/github-actions-runner -n github-runner -- curl -I https://api.github.com
   ```

### Buildah Issues

If container builds fail:
1. Verify privileged SCC is bound
2. Check pod is running as root (UID 0)
3. Verify Buildah is installed in the runner image

## Scaling

To scale runners:

```bash
oc scale deployment github-actions-runner --replicas=3 -n github-runner
```

Each replica will register as a separate runner in GitHub.

## Cleanup

To remove the deployment:

```bash
oc delete -f .
# Then remove runner from GitHub UI
```

## References

- [GitHub Actions Self-Hosted Runners](https://docs.github.com/en/actions/hosting-your-own-runners)
- [OpenShift Security Context Constraints](https://docs.openshift.com/container-platform/latest/authentication/managing-security-context-constraints.html)
- [Buildah Documentation](https://github.com/containers/buildah)
