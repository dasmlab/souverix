# Souverix Testing Philosophy

## Core Principles

### 1. Pristine Platform Assumption

**The underlying platform and all other components are always assumed to be "pristine" (working correctly).**

- When testing Component A, we assume Components B, C, D... are all functioning correctly
- The platform (Kubernetes/OpenShift) is assumed to be stable and operational
- Infrastructure dependencies (databases, message queues, etc.) are assumed to be available and working
- This allows us to isolate failures to the component under test

### 2. One Component Under Test

**On a given system test run, there can only be one component that holds the "new" title (and others must wait in sequence).**

- System tests run sequentially per component to avoid conflicts
- Only one component can be in "testing" state at a time on a given test system
- Other components remain at their last known good state
- This ensures clean test isolation and predictable results

### 3. Failure Handling

**Even if a new component fails its test, it remains on the system (and is pristine from the point of view of others - meaning we always assume everyone else is doing their job).**

- Failed components are not automatically rolled back
- Other components continue to assume the failed component is working
- This allows for debugging and investigation without disrupting other tests
- Manual intervention may be required for cleanup

### 4. Component Workflow Independence

**Component workflows should not know/nor care about the infrastructure other than credentials or variables to call whatever deploy (oc apply, etc.).**

- Component workflows only need:
  - Deployment credentials (kubeconfig, tokens)
  - Environment variables (TEST-SYSTEM-NAME, TEST-SYSTEM-TYPE, etc.)
  - Deployment commands (oc apply, kubectl apply, etc.)
- Infrastructure provisioning is handled separately in `testing-infra/`
- Components are infrastructure-agnostic

### 5. Parallel System Test Runs

**You can have parallel runs of components on parallel system tests, each in their own loops.**

- Multiple test systems can run simultaneously
- Each test system has its own namespace/environment
- Components can be tested in parallel on different systems
- Each system maintains its own state and isolation

## Test System Naming Convention

Each test system is identified by:

- **TEST-SYSTEM-NAME**: Unique identifier for the test system (e.g., `souverix-test-001`)
- **TEST-SYSTEM-TYPE**: Type of system (e.g., `ocp-sno`, `k8s-cluster`, `local`)
- **TEST-SYSTEM-VERSION**: Version of the system/platform (e.g., `4.15`, `1.28`)

These variables are set in GitHub Actions workflows and passed to deployment scripts.

## Test Execution Flow

```
1. Unit Tests (local/CI)
   ↓
2. Build & Publish Container
   ↓
3. System Test Deployment
   - Deploy to TEST-SYSTEM-NAME
   - Use TEST-SYSTEM-TYPE for deployment method
   - Apply component with TEST-SYSTEM-VERSION
   ↓
4. System Test Execution
   - Run component-specific system tests
   - Validate integration with pristine components
   ↓
5. Test Results & Reporting
```

## Component Workflow Integration

After **UNIT test** and **publish**, the workflow calls deploy on a common target:

```yaml
- name: Deploy to System Test
  env:
    TEST-SYSTEM-NAME: ${{ vars.TEST_SYSTEM_NAME }}
    TEST-SYSTEM-TYPE: ${{ vars.TEST_SYSTEM_TYPE }}
    TEST-SYSTEM-VERSION: ${{ vars.TEST_SYSTEM_VERSION }}
  run: |
    oc apply -f k8s/deployment.yaml
```

## Infrastructure Separation

- **Component Workflows**: Handle component build, test, and deployment
- **testing-infra/**: Handles infrastructure provisioning, cluster setup, auto-provisioning
- **Clear Separation**: Components don't manage infrastructure, infrastructure doesn't manage component logic

## Related Documentation

- [Testing Infrastructure](../testing-infra/README.md)
- [Component Workflows](../.github/workflows/)
- [Deployment Guide](../docs/DEPLOYMENT.md)
