# Vault Integration for PKI and Certificate Management

## Overview

Integration with HashiCorp Vault for certificate management in OpenShift Container Platform (OCP), providing secure PKI operations for IMS Core and SIP Gateway functions.

## Architecture

```
┌─────────────────────────────────────────┐
│      OpenShift Cluster                  │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │  Vault (PKI Secret Engine)        │ │
│  │  - Root CA                         │ │
│  │  - Intermediate CAs               │ │
│  │  - Certificate Issuance           │ │
│  │  - Certificate Rotation           │ │
│  └───────────────────────────────────┘ │
│           │                             │
│           v                             │
│  ┌───────────────────────────────────┐ │
│  │  IMS Core Pods                    │ │
│  │  - SBC/IBCF                       │ │
│  │  - STIR/SHAKEN                    │ │
│  │  - TLS Termination                │ │
│  └───────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

## PKI Structure

### Recommended Hierarchy

```
Root CA (Vault)
  │
  ├── Intermediate CA - Factory (Internal)
  │     ├── Component Certificates
  │     └── Service Certificates
  │
  ├── Intermediate CA - Border (Interconnect)
  │     ├── STIR/SHAKEN Certificates
  │     ├── mTLS Peer Certificates
  │     └── TLS Termination Certificates
  │
  └── Intermediate CA - Edge (External)
        ├── Public-facing Certificates
        └── ACME Integration
```

### Certificate Types

1. **Factory (Internal)**
   - Component-to-component communication
   - Service mesh certificates
   - Internal API TLS

2. **Border (Interconnect)**
   - STIR/SHAKEN signing certificates
   - mTLS for carrier peering
   - IBCF border certificates

3. **Edge (External)**
   - Public-facing endpoints
   - ACME-managed certificates
   - User-facing services

## Vault Configuration

### PKI Secret Engine Setup

```hcl
# Root CA
path "pki/root/generate/internal" {
  capabilities = ["update"]
}

path "pki/root/sign/intermediate" {
  capabilities = ["create", "update"]
}

# Intermediate CA - Factory
path "pki-factory/issue/component" {
  capabilities = ["create", "update"]
}

# Intermediate CA - Border
path "pki-border/issue/stir-shaken" {
  capabilities = ["create", "update"]
}

path "pki-border/issue/mtls-peer" {
  capabilities = ["create", "update"]
}

# Certificate Rotation
path "pki-*/issue/*" {
  capabilities = ["create", "update"]
}
```

### Certificate Roles

```hcl
# STIR/SHAKEN Certificate Role
{
  "allowed_domains": ["ims.local", "*.ims.local"],
  "allow_subdomains": true,
  "max_ttl": "720h",  # 30 days
  "key_type": "ec",
  "key_bits": 256,
  "signature_bits": 256
}

# mTLS Peer Certificate Role
{
  "allowed_domains": ["peer1.com", "peer2.com"],
  "allow_subdomains": false,
  "max_ttl": "168h",  # 7 days
  "key_type": "ec",
  "key_bits": 256
}
```

## Integration Points

### 1. Certificate Issuance

```go
// Vault client for certificate issuance
type VaultPKIClient struct {
    client *api.Client
    pkiPath string
}

func (v *VaultPKIClient) IssueCertificate(role string, commonName string) (*Certificate, error) {
    secret, err := v.client.Logical().Write(
        fmt.Sprintf("%s/issue/%s", v.pkiPath, role),
        map[string]interface{}{
            "common_name": commonName,
            "ttl": "720h",
        },
    )
    // Parse and return certificate
}
```

### 2. Certificate Rotation

```go
// Automatic certificate rotation
func (v *VaultPKIClient) RotateCertificate(certID string) error {
    // Issue new certificate
    newCert, err := v.IssueCertificate(role, commonName)
    
    // Update in Kubernetes secret
    // Graceful reload in application
    // Revoke old certificate after grace period
}
```

### 3. AI-Powered Rapid Rotation

Future enhancement: AI engine for certificate rotation based on:
- Threat intelligence
- Anomaly detection
- Risk scoring
- Automated rotation triggers

## OpenShift Integration

### Service Account and RBAC

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ims-vault-auth
  namespace: ims-core
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: vault-auth
rules:
- apiGroups: [""]
  resources: ["serviceaccounts/token"]
  verbs: ["create"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ims-vault-auth
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: vault-auth
subjects:
- kind: ServiceAccount
  name: ims-vault-auth
  namespace: ims-core
```

### Vault Auth Method (Kubernetes)

```hcl
# Vault Kubernetes auth configuration
path "auth/kubernetes/role/ims-core" {
  capabilities = ["read"]
  bound_service_account_names = ["ims-vault-auth"]
  bound_service_account_namespaces = ["ims-core"]
  policies = ["ims-pki-policy"]
  ttl = "1h"
}
```

### Kubernetes Secrets Integration

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: stir-shaken-cert
  namespace: ims-core
type: Opaque
data:
  certificate: <from-vault>
  private_key: <from-vault>
  ca_chain: <from-vault>
```

## Certificate Lifecycle

### 1. Initial Issuance

- Component starts
- Authenticates to Vault via Kubernetes auth
- Requests certificate from appropriate PKI path
- Stores in Kubernetes secret
- Application loads certificate

### 2. Rotation

- Monitor certificate expiration (e.g., 80% of TTL)
- Issue new certificate from Vault
- Update Kubernetes secret
- Graceful reload in application
- Revoke old certificate after grace period

### 3. Revocation

- Certificate compromise detected
- Immediate revocation via Vault
- Issue emergency replacement
- Update all components

## Security Considerations

1. **Key Storage**: Private keys never leave Vault
2. **Access Control**: Role-based access to PKI paths
3. **Audit Logging**: All certificate operations logged
4. **Rotation Policy**: Automated rotation before expiration
5. **Revocation**: OCSP/CRL support for revoked certificates

## Monitoring

- Certificate expiration alerts
- Rotation success/failure metrics
- Vault API latency
- Certificate issuance rate
- Revocation events

## Future: AI-Powered Rotation

- Threat intelligence integration
- Anomaly detection triggers
- Risk-based rotation frequency
- Automated response to threats
- Predictive certificate management
