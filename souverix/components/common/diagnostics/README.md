# Diagnostics Framework

The diagnostics framework provides a common library for all Souverix components to implement standardized diagnostic endpoints and call flow testing.

## Architecture Principles

### 1. Framework, Not Server
- **Diag common is NEVER a server** - it's a library and framework
- Components expose `/diag` endpoints on their own servers (R3 - diagnostics server)
- Common library provides the structure, templates, and utilities

### 2. Server Abstraction
- Components may have different numbers of servers (R1, R2, R3, R4, etc.)
- Diagnostics framework abstracts server health checking
- Each component registers its servers with the framework
- Framework provides standardized health checks for all registered servers

### 3. Call Flow Templates
- Diag common maintains the "mother list" of call flows at ETSI/3GPP standard level
- Map table: Component → Call flows → Sequence steps that component participates in
- Components only know their 1-hop neighbors in call flows
- Framework provides call flow context to components

### 4. Component-Specific Diagnostics
- Each component has diagnostics relative to their specific function
- Local map table of their functions and view
- Messages are "hardwired" for diagnostic purposes
- When uncertain, diag common provides framework and call path info

## Endpoints

### `/diag/health`
Basic health check - verifies component is running

### `/diag/status`
Component status, version, and build information

### `/diag/local_test`
Local testing endpoint - returns success for manual testing

### `/diag/unit_test`
Unit testing endpoint - executes call flow simulation for component's portion

## Unit Test Framework

### Purpose
- Simulate calling operations and responses for selected call paths
- Verify component's data changes, state updates, storage operations
- **NOT testing the actual call** - testing component's internal operations

### How It Works

1. **Call Path Selection**: Unit test selects a call flow from diag common's master list
2. **Faux Call Generation**: Diag common generates "faux" calling parts and responses
3. **Component Execution**: Component executes its portion of the call flow
4. **State Verification**: Verify data elements, state changes, storage operations

### Call Flow Simulation

```
caller --> component (under test) --> something else
caller <--- component (under test) <--- something else
```

- The "something else" callback must be faked out enough to make the caller work
- Always uses curl (never cross-cutting)
- Fakes out actual service endpoints
- Key: Verify component's internal operations, not the actual call

### Runtime Integration

- Components can run unit tests during normal runtime via Go subroutines
- Backed by tickets/queues (abstracted in diag library)
- Components call their own `/diag/unit_test` endpoint
- Execute call flow portion relative to them
- All call path info comes from diag common library (not stored in containers)

## Future: System-Level Diagnostics

- Platform connectivity (message buses, anycasts, SSEs)
- Network-level diagnostics (Mx interface, NODE level)
- Abstracted functions in the network
- System test layer integration
