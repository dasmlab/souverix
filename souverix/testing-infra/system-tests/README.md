# System Tests

System-level integration tests for Souverix components.

## Philosophy

See [Testing Philosophy](../../docs/TESTING_PHILOSOPHY.md) for details.

**Key Points:**
- Platform and other components are always assumed pristine
- Only one component is "new" (under test) at a time
- Failed components remain on system (still considered pristine by others)
- Component workflows are infrastructure-agnostic

## Test System Configuration

Each test system is configured via environment variables:

```bash
TEST-SYSTEM-NAME=souverix-test-001
TEST-SYSTEM-TYPE=ocp-sno
TEST-SYSTEM-VERSION=4.15
```

These variables are:
- Set in GitHub Actions repository variables
- Passed to component workflows
- Used by deployment scripts
- Identified in test results

## Test Execution Flow

1. **Component Build & Unit Test** (in component workflow)
2. **Container Publish** (to registry)
3. **System Test Deployment** (deploy to TEST-SYSTEM-NAME)
4. **System Test Execution** (run component-specific tests)
5. **Test Results** (report success/failure)

## Component Test Structure

Each component has its own test suite in `tests/<component>/`:

```
system-tests/
├── tests/
│   ├── coeur/
│   │   ├── pcscf/
│   │   ├── icscf/
│   │   ├── scscf/
│   │   ├── bgcf/
│   │   ├── mgcf/
│   │   └── hss/
│   ├── rempart/
│   └── ...
└── deploy/
    └── deploy-component.sh
```

## Deployment Scripts

Deployment scripts in `deploy/` handle:
- Reading TEST-SYSTEM-* variables
- Applying Kubernetes/OpenShift manifests
- Waiting for deployment readiness
- Health check validation

Components only need to provide:
- Deployment manifests (k8s/*.yaml)
- Deployment credentials (from secrets)
- Environment variables (TEST-SYSTEM-*)

## Parallel Execution

Multiple test systems can run in parallel:
- Each has unique TEST-SYSTEM-NAME
- Each maintains its own state
- Components can be tested simultaneously on different systems
- No cross-system interference

## Related Documentation

- [Testing Philosophy](../../docs/TESTING_PHILOSOPHY.md)
- [Testing Infrastructure](../README.md)
