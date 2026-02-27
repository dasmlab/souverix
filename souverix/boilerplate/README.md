# Souverix Component Boilerplate

This directory contains templates and scripts for scaffolding new Souverix components.

## Structure

```
boilerplate/
├── templates/          # Template files for component scaffolding
│   ├── Dockerfile.COMPONENT
│   ├── cmd/main.go
│   ├── buildme.sh
│   ├── pushme.sh
│   ├── pullme.sh
│   ├── runme-local.sh
│   └── k8s/
│       ├── deployment.yaml
│       └── configmap.yaml
└── scripts/
    └── scaffold-component.sh   # Main scaffolding script
```

## Usage

To scaffold a new component:

```bash
./boilerplate/scripts/scaffold-component.sh <component> <component-name> <internal-path> [ports]
```

### Parameters

- `component`: Short component name (e.g., `coeur`, `rempart`)
- `component-name`: Full display name (e.g., `Souverix Coeur`)
- `internal-path`: Path within `internal/` (e.g., `internal/coeur`)
- `ports`: Space-separated port mappings (e.g., `"5060:5060 8081:8081"`)

### Examples

```bash
# Scaffold Souverix Coeur (IMS Core)
./boilerplate/scripts/scaffold-component.sh \
  coeur \
  "Souverix Coeur" \
  "internal/coeur" \
  "5060:5060 5061:5061 8081:8081"

# Scaffold Souverix Rempart (Border Control)
./boilerplate/scripts/scaffold-component.sh \
  rempart \
  "Souverix Rempart" \
  "internal/rempart" \
  "5060:5060 5061:5061 8081:8081"

# Scaffold a simple component with just health port
./boilerplate/scripts/scaffold-component.sh \
  mandat \
  "Souverix Mandat" \
  "internal/mandat" \
  "8081:8081"
```

## Generated Files

Each component scaffold generates:

1. **Dockerfile.`<component>`** - Multi-stage build with golang:1.23-alpine
2. **internal/`<path>`/cmd/main.go** - Component entrypoint with standard logging
3. **buildme-`<component>`.sh** - Build script (Docker/Podman)
4. **pushme-`<component>`.sh** - Push to registry with SemVer
5. **pullme-`<component>`.sh** - Pull from registry
6. **runme-local-`<component>`.sh** - Run locally in container
7. **k8s/`<component>`/deployment.yaml** - Kubernetes deployment
8. **k8s/`<component>`/configmap.yaml** - Kubernetes ConfigMap

## Standard Component Structure

All components follow this pattern:

### main.go

- First line logs: `"Souverix - <Component Name> - Version: <version> Build: <git-commit>"`
- Uses `gouverneConfig.Load()` for configuration
- Implements graceful shutdown with context timeout
- Uses logrus for structured logging

### Dockerfile

- Multi-stage build: `golang:1.23-alpine` → `alpine:latest`
- Builds with version, build time, and git commit
- Exposes component-specific ports + health port (8081)
- Includes health check endpoint

### Scripts

- **buildme.sh**: Detects Docker/Podman, builds with version info
- **pushme.sh**: Pushes to `ghcr.io/dasmlab/<component>` with SemVer
- **pullme.sh**: Pulls from registry and tags as `local`
- **runme-local.sh**: Runs container locally with proper port mapping

### Kubernetes

- Namespace: `souverix-<component>`
- Labels: `app: souverix-<component>`, `component: <component>`, `platform: souverix`
- Health checks on port 8081
- ConfigMap for component-specific settings

## Customization

After scaffolding, customize:

1. **Component logic**: Implement actual functionality in `internal/<component>/`
2. **Configuration**: Add component-specific config in `ConfigFromGouverne()` or `New()`
3. **Ports**: Update Dockerfile and k8s YAMLs if ports change
4. **Dependencies**: Add to `go.mod` as needed

## CI/CD Integration

All scaffolded components automatically integrate with:

- **GitHub Actions workflows**: `.github/workflows/<component>-full-pipeline.yml`
- **Kaniko builds**: Uses self-hosted runner with `kaniko` label
- **Path triggers**: Builds on changes to `souverix/internal/<component>/**`

## Notes

- Components start as stubs and need implementation
- All components use the same logging and configuration patterns
- Health endpoint is always on port 8081
- Build scripts support both Docker and Podman
