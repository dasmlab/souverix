# OpenShift Container Platform (OCP) Self-Hosted Runner Setup

## Overview

This document describes how to configure GitHub Actions self-hosted runners on OpenShift Container Platform (OCP) for the CI/CD pipeline.

## Prerequisites

- OpenShift Cluster (4.x+)
- Cluster admin or project admin access
- GitHub repository with Actions enabled
- GitHub Personal Access Token (PAT) or GitHub App token

## Runner Deployment Pattern

The CI/CD pipeline requires Docker-in-Docker (dind) for building containers. The runner pod must include a `docker:dind` sidecar container.

### Kubernetes Pattern

```
Runner Pod:
  - actions-runner container (main runner)
  - docker:dind sidecar (Docker daemon)
```

## Deployment Steps

### 1. Create Project/Namespace

```bash
oc new-project github-runners
# Or use existing project
oc project github-runners
```

### 2. Create Service Account

```bash
oc create serviceaccount github-runner-sa
```

### 3. Grant Privileged Access (Required for dind)

**Note**: dind requires privileged containers. In OCP, this requires appropriate SCC (Security Context Constraint).

```bash
# Grant anyuid SCC (or create custom SCC)
oc adm policy add-scc-to-user privileged -z github-runner-sa
```

**Alternative**: Create a custom SCC for more restrictive access:

```yaml
apiVersion: security.openshift.io/v1
kind: SecurityContextConstraints
metadata:
  name: github-runner-scc
allowPrivilegedContainer: true
allowHostDirVolumePlugin: true
allowHostIPC: true
allowHostNetwork: true
allowHostPID: true
allowHostPorts: true
runAsUser:
  type: RunAsAny
seLinuxContext:
  type: RunAsAny
fsGroup:
  type: RunAsAny
users:
  - system:serviceaccount:github-runners:github-runner-sa
```

### 4. Create Runner Deployment

Create a Deployment that includes both the GitHub Actions runner and docker:dind sidecar:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: github-runner
  namespace: github-runners
spec:
  replicas: 1
  selector:
    matchLabels:
      app: github-runner
  template:
    metadata:
      labels:
        app: github-runner
    spec:
      serviceAccountName: github-runner-sa
      containers:
      # GitHub Actions Runner
      - name: runner
        image: myoung34/github-runner:latest
        env:
        - name: RUNNER_NAME
          value: "ocp-runner-1"
        - name: RUNNER_REPOSITORY_URL
          value: "https://github.com/dasmlab/ims"
        - name: RUNNER_TOKEN
          valueFrom:
            secretKeyRef:
              name: github-runner-secret
              key: token
        - name: RUNNER_WORKDIR
          value: "/tmp/_work"
        - name: DOCKER_HOST
          value: "tcp://localhost:2375"
        - name: DOCKER_TLS_CERTDIR
          value: ""
        volumeMounts:
        - name: workdir
          mountPath: /tmp/_work
        - name: dockersock
          mountPath: /var/run/docker.sock
        securityContext:
          privileged: true
          runAsUser: 0
      
      # Docker-in-Docker sidecar
      - name: dind
        image: docker:dind
        securityContext:
          privileged: true
        volumeMounts:
        - name: dockersock
          mountPath: /var/run/docker.sock
        env:
        - name: DOCKER_TLS_CERTDIR
          value: ""
        args:
        - dockerd
        - --host=unix:///var/run/docker.sock
        - --host=tcp://0.0.0.0:2375
      
      volumes:
      - name: workdir
        emptyDir: {}
      - name: dockersock
        emptyDir: {}
```

### 5. Create GitHub Runner Secret

```bash
# Get runner token from GitHub:
# Settings → Actions → Runners → New self-hosted runner

oc create secret generic github-runner-secret \
  --from-literal=token=YOUR_RUNNER_TOKEN
```

### 6. Apply Deployment

```bash
oc apply -f github-runner-deployment.yaml
```

### 7. Verify Runner

```bash
# Check pod status
oc get pods -n github-runners

# Check logs
oc logs -n github-runners -l app=github-runner -c runner
oc logs -n github-runners -l app=github-runner -c dind

# Verify Docker daemon
oc exec -n github-runners -c runner -it <pod-name> -- docker info
```

## Workflow Configuration

The workflow is already configured with:

```yaml
env:
  DOCKER_HOST: tcp://localhost:2375
  DOCKER_TLS_CERTDIR: ""
```

These environment variables tell the workflow to use the dind sidecar's Docker daemon.

## Troubleshooting

### Docker daemon not available

**Symptoms**: `docker info` fails, build steps fail

**Solutions**:
1. Verify dind sidecar is running: `oc logs -c dind <pod-name>`
2. Check DOCKER_HOST environment variable is set
3. Verify privileged mode is enabled
4. Check SCC is applied: `oc get scc github-runner-scc`

### Permission denied errors

**Symptoms**: Cannot access Docker socket

**Solutions**:
1. Verify service account has privileged SCC
2. Check security context allows privileged containers
3. Verify volume mounts are correct

### Runner not connecting to GitHub

**Symptoms**: Runner doesn't appear in GitHub Actions

**Solutions**:
1. Verify RUNNER_TOKEN is correct
2. Check network connectivity from pod
3. Verify RUNNER_REPOSITORY_URL is correct
4. Check runner logs: `oc logs -c runner <pod-name>`

## Alternative: Using Buildah/Podman

For OCP, you might prefer using Buildah/Podman instead of Docker:

```yaml
- name: Build with Buildah
  run: |
    buildah bud -f Dockerfile.pcscf -t $IMAGE_NAME:$TAG .
    buildah push $IMAGE_NAME:$TAG
```

This avoids the need for privileged containers and dind.

## Security Considerations

- **Privileged containers**: Required for dind, increases security risk
- **SCC**: Use custom SCC with minimal required permissions
- **Network policies**: Restrict network access as needed
- **Resource limits**: Set appropriate CPU/memory limits
- **Image scanning**: Scan runner images for vulnerabilities

## References

- [GitHub Actions Self-Hosted Runners](https://docs.github.com/en/actions/hosting-your-own-runners)
- [OpenShift Security Context Constraints](https://docs.openshift.com/container-platform/latest/authentication/managing-security-context-constraints.html)
- [Docker-in-Docker](https://hub.docker.com/_/docker)
