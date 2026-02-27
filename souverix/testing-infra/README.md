# Testing Infrastructure

This directory contains infrastructure provisioning, cluster management, and system test orchestration for Souverix.

## Philosophy

See [Testing Philosophy](../docs/TESTING_PHILOSOPHY.md) for core principles:
- Pristine platform assumption
- One component under test at a time
- Component workflow independence
- Parallel system test runs

## Directory Structure

```
testing-infra/
├── README.md              # This file
├── call-flow-test/        # Basic call flow validation tests
├── system-tests/          # System-level integration tests
│   ├── README.md         # System test documentation
│   ├── deploy/           # Deployment scripts
│   └── tests/            # Test suites per component
└── provisioning/         # Infrastructure provisioning (future)
    └── clusters/         # Cluster auto-provisioning
```

## Test System Variables

Each test system is identified by environment variables:

- **TEST-SYSTEM-NAME**: Unique identifier (e.g., `souverix-test-001`)
- **TEST-SYSTEM-TYPE**: System type (e.g., `ocp-sno`, `k8s-cluster`)
- **TEST-SYSTEM-VERSION**: Platform version (e.g., `4.15`, `1.28`)

These are set in GitHub Actions repository variables and passed to workflows.

## System Test Execution

System tests assume:
- Platform is pristine and operational
- All other components are at their last known good state
- Only the component under test is "new"
- Infrastructure is already provisioned and available

## Future: Auto-Provisioning

Auto-provisioning of clusters and infrastructure will be added here:
- Cluster creation scripts
- Environment setup automation
- Resource allocation and cleanup
- Multi-cluster management

## Related Documentation

- [Testing Philosophy](../docs/TESTING_PHILOSOPHY.md)
- [System Tests](./system-tests/README.md)
