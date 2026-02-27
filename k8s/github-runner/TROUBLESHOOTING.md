# Troubleshooting GitHub Runner Registration

## Common Issues

### 404 Not Found Error

If you see `404 Not Found` when registering the runner:

1. **Check Token Format**
   - The token must be a **registration token** (not a PAT)
   - Registration tokens expire after **1 hour**
   - Get a fresh token from: Repository → Settings → Actions → Runners → New self-hosted runner

2. **Verify REPO_URL Format**
   - For repo-level runners: `dasmlab/souverix` (owner/repo format)
   - NOT: `https://github.com/dasmlab/souverix`
   - Note: The local disk path is `ims`, but the GitHub repository is `souverix`

3. **Check Token Permissions**
   - The registration token must be generated from the **same repository** you're trying to register
   - If the repo is private, ensure the token has access

4. **Verify Repository Exists**
   - Confirm the repository `dasmlab/souverix` exists and is accessible
   - Check you have admin access to the repository

### Invalid Token Configuration

If you see "Invalid configuration provided for token":

1. **Environment Variable Name**
   - For repo-level runners: Use `RUNNER_TOKEN` (not `GITHUB_TOKEN`)
   - For org-level runners: Use `GITHUB_TOKEN`

2. **Secret Key Name**
   - Ensure the secret key is named `github_token` in your Kubernetes secret
   - Verify the secret exists: `oc get secret github-runner-secrets -n github-runner`

### Getting a Fresh Registration Token

1. Go to: https://github.com/dasmlab/souverix/settings/actions/runners/new
2. Copy the registration token (starts with something like `AXXXXXXXXXXXXXXXXXXXXX`)
3. Update the secret:
   ```bash
   oc create secret generic github-runner-secrets \
     --from-literal=github_token="YOUR_NEW_TOKEN_HERE" \
     --namespace=github-runner \
     --dry-run=client -o yaml | oc apply -f -
   ```
4. Restart the pod:
   ```bash
   oc delete pod -l app=github-runner -n github-runner
   ```

### Checking Runner Logs

```bash
# Get pod name
oc get pods -n github-runner

# View logs
oc logs -f <pod-name> -n github-runner -c runner
```

### Verifying Network Connectivity

```bash
# Test GitHub API access from the pod
oc exec -it <pod-name> -n github-runner -c runner -- \
  curl -I https://api.github.com
```

## Environment Variables Reference

For `myoung34/github-runner` repo-level runners:

- `REPO_URL`: `dasmlab/souverix` (owner/repo format)
- `RUNNER_TOKEN`: Registration token from GitHub
- `RUNNER_NAME`: Unique name for this runner instance
- `RUNNER_WORKDIR`: Working directory (default: `/tmp/_work`)
- `DOCKER_HOST`: `tcp://localhost:2375` (for dind sidecar)
- `DOCKER_TLS_CERTDIR`: `""` (empty string to disable TLS)
