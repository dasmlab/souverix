# Souverix Common Library

Common packages and frameworks for all Souverix components.

## Module

```
github.com/dasmlab/souverix/common
```

## Versioning Contract

**N+1 and N-1 Compatibility:**
- The common library maintains backward and forward compatibility
- Components can use N+1 (newer) or N-1 (older) versions
- Breaking changes are avoided; new features are additive
- Local extensions in components are supported without breaking others

## Packages

### Diagnostics Framework
- **Package**: `github.com/dasmlab/souverix/common/diagnostics`
- **Purpose**: Common diagnostic endpoints for all components
- **Endpoints**: `/diag/health`, `/diag/status`, `/diag/local_test`, `/diag/unit_test`

### Testing Framework
- **Package**: `github.com/dasmlab/souverix/common/testing`
- **Purpose**: Common testing utilities and frameworks
- **Features**: Test helpers, mocks, fixtures

### Authentication
- **Package**: `github.com/dasmlab/souverix/common/auth`
- **Purpose**: OAuth and authentication utilities
- **Note**: No certificate management (ZTA - Zero Trust Architecture)

### SIP
- **Package**: `github.com/dasmlab/souverix/common/sip`
- **Purpose**: SIP message handling and parsing

### HSS
- **Package**: `github.com/dasmlab/souverix/common/hss`
- **Purpose**: HSS client and Cx interface simulation

## Usage

```go
import (
    "github.com/dasmlab/souverix/common/diagnostics"
    "github.com/dasmlab/souverix/common/testing"
)
```

## Publishing

The common library is published as a versioned Go module:
- Version tags: `v1.0.0`, `v1.1.0`, etc.
- Components reference via `go.mod`
- Local development uses `replace` directive

## Extension Support

Components can extend common functionality locally:
- Create component-specific extensions
- Won't break other components
- Data payloads remain compatible
