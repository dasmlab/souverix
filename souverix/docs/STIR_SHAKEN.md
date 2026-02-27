# STIR/SHAKEN Implementation

## Overview

STIR/SHAKEN (Secure Telephone Identity Revisited / Signature-based Handling of Asserted information using toKENs) is the caller identity authentication framework used to combat caller ID spoofing in SIP-based voice networks.

## Our Implementation Approach

### ACME-Based Certificate Management

**Key Innovation**: We use ACME (Automatic Certificate Management Environment) protocol for STIR/SHAKEN certificate management, addressing interoperability issues found in traditional implementations.

#### Why ACME for STIR/SHAKEN?

Traditional STIR/SHAKEN implementations face interoperability challenges:

1. **Proprietary Certificate Distribution**: Many implementations use proprietary mechanisms for certificate distribution, making cross-carrier interoperability difficult.

2. **Manual Certificate Exchange**: Operators must manually exchange certificates, creating operational overhead and potential security gaps.

3. **Inconsistent Trust Models**: Different regions and operators use different certificate authorities and trust models.

4. **Certificate Lifecycle Management**: Manual renewal and rotation of certificates is error-prone.

#### Our Solution

By leveraging ACME (RFC 8555) for STIR/SHAKEN certificates:

- **Standard Protocol**: ACME is a well-established, standardized protocol
- **Automated Lifecycle**: Automatic certificate issuance, renewal, and rotation
- **Interoperability**: Works with any ACME-compatible CA (Let's Encrypt, ZeroSSL, custom)
- **Integration with Zero Trust**: Seamlessly integrates with our Zero Trust Architecture
- **Cloud-Native**: Perfect fit for Kubernetes/OpenShift deployments

## Architecture

### Components

```
┌─────────────────────────────────────────┐
│           SIP Gateway / SBC              │
│                                           │
│  ┌─────────────────────────────────────┐ │
│  │      STIR Signer                    │ │
│  │  - Generates PASSporT tokens       │ │
│  │  - Signs with ECDSA P-256          │ │
│  │  - Adds Identity header            │ │
│  └─────────────────────────────────────┘ │
│                                           │
│  ┌─────────────────────────────────────┐ │
│  │      STIR Verifier                 │ │
│  │  - Verifies Identity header         │ │
│  │  - Fetches certificates via HTTPS  │ │
│  │  - Validates attestation level      │ │
│  └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
           │                    │
           │                    │
           v                    v
┌──────────────────┐  ┌──────────────────┐
│  ACME Certificate│  │  Certificate     │
│  Manager         │  │  Fetcher         │
│                  │  │                  │
│  - ACME client   │  │  - HTTPS fetch   │
│  - Key generation│  │  - PEM parsing    │
│  - Auto-renewal  │  │  - Validation    │
└──────────────────┘  └──────────────────┘
           │
           v
┌──────────────────┐
│  ACME Provider   │
│  (Let's Encrypt, │
│   ZeroSSL, etc.) │
└──────────────────┘
```

## Call Flow with STIR/SHAKEN

### Originating Call (Outbound)

```
1. UE sends INVITE
   ↓
2. P-CSCF receives INVITE
   ↓
3. S-CSCF processes INVITE
   ↓
4. SBC/SIP GW receives INVITE
   ↓
5. STIR Signer:
   - Extracts orig/dest TNs
   - Determines attestation level
   - Generates PASSporT token
   - Signs with ECDSA key
   - Adds Identity header
   ↓
6. INVITE forwarded with Identity header
   ↓
7. External network receives verified call
```

### Terminating Call (Inbound)

```
1. External network sends INVITE with Identity header
   ↓
2. SBC/SIP GW receives INVITE
   ↓
3. STIR Verifier:
   - Extracts Identity header
   - Fetches certificate from x5u URL
   - Verifies ECDSA signature
   - Validates token claims
   - Checks attestation level
   ↓
4. If verified:
   - Adds X-STIR-Verified header
   - Adds X-STIR-Attestation header
   ↓
5. INVITE forwarded to IMS core
   ↓
6. UE receives call with verification status
```

## Attestation Levels

| Level | Meaning | Use Case |
|-------|---------|----------|
| **A** | Full attestation | Provider knows customer & has number control |
| **B** | Partial attestation | Provider knows customer but not full number control |
| **C** | Gateway attestation | Call originated externally (e.g., PSTN gateway) |

### Auto-Detection

Our implementation can automatically determine attestation level based on:
- Subscriber registration status
- Number ownership verification
- Call origin (internal vs. external)

## Configuration

### Environment Variables

```bash
# Enable STIR/SHAKEN
SBC_ENABLE_STIR=true

# Attestation level (A, B, C, or "auto")
SBC_STIR_ATTESTATION=auto

# ACME configuration (for certificate management)
ZTA_ACME_PROVIDER=letsencrypt
ZTA_ACME_EMAIL=admin@ims.local
ZTA_ACME_DOMAIN=ims.local
ZTA_ACME_STAGING=false
```

### Kubernetes ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ims-core-config
data:
  sbc-enable-stir: "true"
  sbc-stir-attestation: "auto"
  zta-acme-provider: "letsencrypt"
  zta-acme-domain: "ims.local"
```

## Certificate Lifecycle

### Initial Certificate Acquisition

1. **ACME Account Creation**: Create account with ACME provider
2. **Domain Validation**: Complete HTTP-01 or DNS-01 challenge
3. **Certificate Issuance**: Download certificate and store
4. **Key Generation**: Generate ECDSA P-256 key pair
5. **Certificate URL**: Publish certificate at `/.well-known/stir/cert.pem`

### Automatic Renewal

- Certificates are automatically renewed before expiration
- Zero-downtime renewal using certificate rotation
- Integration with Kubernetes secrets for secure storage

### Certificate Distribution

Certificates are distributed via:
- HTTPS endpoint: `https://domain/.well-known/stir/cert.pem`
- Standard ACME protocol
- No proprietary mechanisms required

## Interoperability Benefits

### Traditional STIR/SHAKEN Issues

1. **Manual Certificate Exchange**: Operators must manually exchange certificates
2. **Proprietary Protocols**: Different vendors use different certificate distribution methods
3. **Trust Model Complexity**: Multiple CAs and trust chains
4. **Operational Overhead**: Manual certificate lifecycle management

### Our ACME-Based Solution

1. **Automated Exchange**: Certificates automatically available via HTTPS
2. **Standard Protocol**: ACME is IETF standard (RFC 8555)
3. **Unified Trust Model**: Works with any ACME-compatible CA
4. **Zero-Touch Operations**: Fully automated certificate lifecycle

## Security Considerations

### Key Protection

- Private keys stored in Kubernetes secrets
- Optional HSM integration for production
- Key rotation policies

### Certificate Validation

- Certificate chain validation
- Expiration checking
- Revocation list support (future)

### Replay Protection

- Short-lived tokens (5 minutes)
- Call-ID based token IDs
- Timestamp validation

## Integration with Zero Trust

STIR/SHAKEN integrates seamlessly with our Zero Trust Architecture:

- Uses same ACME infrastructure
- Unified certificate management
- Consistent security model
- Single configuration point

## Future Enhancements

1. **Full ACME Client**: Complete RFC 8555 implementation
2. **HSM Integration**: Hardware security module support
3. **Certificate Revocation**: OCSP/CRL support
4. **Fraud Analytics**: Integration with fraud detection
5. **Cross-Border Attestation**: International interoperability

## Standards Compliance

- **RFC 8224**: SIP Identity Header
- **RFC 8225**: PASSporT Token
- **RFC 8588**: Certificate Management
- **RFC 8555**: ACME Protocol
- **ATIS-1000074**: SHAKEN Framework

## Example SIP Message

```
INVITE sip:+15145551234@ims.local SIP/2.0
Via: SIP/2.0/UDP 192.168.1.1:5060
From: <sip:+15145559876@ims.local>;tag=abc123
To: <sip:+15145551234@ims.local>
Call-ID: call-12345@ims.local
CSeq: 1 INVITE
Identity: eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJvcmlnIjp7InRuIjoiMTUxNDU1NTk4NzYifSwiZGVzdCI6eyJ0biI6WyIxNTE0NTU1MTIzNCJdfSwiaWF0IjoxNzA5MDUxMjAwLCJleHAiOjE3MDkwNTQyMDAsImF0dGVzdCI6IkEifQ.signature
Content-Type: application/sdp
Content-Length: 142

v=0
o=alice 2890844526 2890844526 IN IP4 192.168.1.1
s=-
c=IN IP4 192.168.1.1
t=0 0
m=audio 49170 RTP/AVP 0
a=rtpmap:0 PCMU/8000
```

## Troubleshooting

### Certificate Fetch Failures

- Check ACME provider connectivity
- Verify domain DNS resolution
- Check certificate URL accessibility

### Verification Failures

- Verify certificate chain
- Check token expiration
- Validate signature algorithm (must be ES256)

### Attestation Level Issues

- Review subscriber registration status
- Check number ownership verification
- Verify call origin classification
