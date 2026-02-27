# Getting Started with IMS Core

## Quick Start

### Prerequisites

- Go 1.23 or later
- Docker or Podman
- Kubernetes cluster (for K8s deployment)
- Ansible (for Ansible deployment)

### Build

```bash
./buildme.sh
```

This will build the IMS Core container image with tag `ims-core:local`.

### Run Locally

```bash
./runme-local.sh
```

Or manually:

```bash
docker run -d \
  --name ims-core-local \
  -p 5060:5060/udp \
  -p 5060:5060/tcp \
  -p 8080:8080 \
  -p 9443:9443 \
  -e LOG_LEVEL=debug \
  ims-core:local
```

### Configuration

Key environment variables:

- `LOG_LEVEL`: Log level (debug, info, warn, error)
- `IMS_DOMAIN`: IMS domain name (default: ims.local)
- `ENABLE_SBC`: Enable SBC/IBCF (default: true)
- `ENABLE_HSS`: Enable HSS (default: true)
- `ZERO_TRUST_MODE`: Enable Zero Trust Mode (default: false)
- `SBC_TOPOLOGY_HIDING`: Enable topology hiding (default: true)
- `SBC_DOS_PROTECTION`: Enable DoS protection (default: true)

### Zero Trust Mode

To enable Zero Trust Mode:

```bash
export ZERO_TRUST_MODE=true
export ZTA_CA_PROVIDER=internal  # or "vault", "acme"
./runme-local.sh
```

### API Endpoints

- Health: `http://localhost:8080/health`
- Metrics: `http://localhost:9443/metrics`
- Subscribers: `http://localhost:8080/api/v1/subscribers`
- Registrations: `http://localhost:8080/api/v1/registrations/:impi`

### SIP Endpoints

- UDP: `0.0.0.0:5060`
- TCP: `0.0.0.0:5060`
- TLS: `0.0.0.0:5061` (when TLS enabled)

## Kubernetes Deployment

### Using Manifests

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/configmap.yaml
```

### Using Ansible

```bash
ansible-playbook -i ansible/inventory.yml ansible/deploy.yml
```

## Testing

### Test SIP Message

Use `sipgrep` or `sipsak` to send test SIP messages:

```bash
# Register
sipsak -U -s sip:alice@ims.local -a alice -w secret123 sip:ims.local:5060

# INVITE
sipsak -U -s sip:alice@ims.local -a alice -w secret123 -M -B "sip:bob@ims.local" sip:ims.local:5060
```

### Test API

```bash
# Health check
curl http://localhost:8080/health

# List subscribers
curl http://localhost:8080/api/v1/subscribers

# Get subscriber
curl http://localhost:8080/api/v1/subscribers/alice@ims.local
```

## Next Steps

1. Review [Architecture Hierarchy](../architecture/hierarchy.md) for detailed architecture
2. Configure Zero Trust Mode for production
3. Set up persistent storage for HSS
4. Configure TLS/SRTP for secure SIP
5. Integrate with external systems
