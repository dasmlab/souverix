# STIR/SHAKEN Interoperability with ACME

## The Problem

Traditional STIR/SHAKEN implementations face significant interoperability challenges:

### 1. Proprietary Certificate Distribution

**Issue**: Many vendors use proprietary mechanisms for certificate distribution between carriers.

**Impact**:
- Manual certificate exchange required
- Different protocols for different vendors
- Complex integration between carriers
- Operational overhead

**Example**:
- Vendor A uses proprietary REST API
- Vendor B uses SOAP web services
- Vendor C uses custom protocol
- Each requires different integration

### 2. Manual Certificate Exchange

**Issue**: Operators must manually exchange certificates through:
- Email
- Secure file transfer
- Portal-based systems
- Direct API integration

**Impact**:
- Time-consuming setup
- Error-prone process
- Delayed certificate updates
- Security risks from manual handling

### 3. Inconsistent Trust Models

**Issue**: Different regions and operators use different:
- Certificate authorities
- Trust chains
- Validation methods
- Attestation models

**Impact**:
- Cross-border interoperability issues
- Complex trust relationships
- Verification failures
- Limited global adoption

### 4. Certificate Lifecycle Management

**Issue**: Manual renewal and rotation of certificates.

**Impact**:
- Expired certificates cause service outages
- Manual intervention required
- Operational complexity
- Security gaps during transitions

## Our Solution: ACME-Based Certificate Management

### Why ACME?

**ACME (Automatic Certificate Management Environment)** is:
- ✅ **IETF Standard** (RFC 8555)
- ✅ **Widely Adopted** (Let's Encrypt, ZeroSSL, etc.)
- ✅ **Automated** certificate lifecycle
- ✅ **Interoperable** across vendors
- ✅ **Cloud-Native** friendly

### How It Works

```
┌─────────────────────────────────────────┐
│  Carrier A (Our IMS Core)               │
│                                           │
│  1. Generate ECDSA key pair              │
│  2. Request certificate via ACME        │
│  3. Complete domain validation           │
│  4. Receive certificate                  │
│  5. Publish at /.well-known/stir/cert.pem│
└─────────────────────────────────────────┘
                    │
                    │ HTTPS
                    v
┌─────────────────────────────────────────┐
│  ACME Provider                          │
│  (Let's Encrypt, ZeroSSL, Custom)       │
└─────────────────────────────────────────┘
                    │
                    │ HTTPS
                    v
┌─────────────────────────────────────────┐
│  Carrier B (Any STIR/SHAKEN Verifier)  │
│                                           │
│  1. Receive INVITE with Identity header │
│  2. Extract x5u URL from token           │
│  3. Fetch certificate via HTTPS         │
│  4. Verify signature                    │
│  5. Validate attestation                │
└─────────────────────────────────────────┘
```

### Benefits

#### 1. Standard Protocol

- **No Proprietary APIs**: Uses standard HTTPS
- **Vendor Agnostic**: Works with any ACME-compatible system
- **Future Proof**: Based on IETF standards

#### 2. Automated Lifecycle

- **Automatic Issuance**: Certificates obtained automatically
- **Auto-Renewal**: Certificates renewed before expiration
- **Zero-Downtime**: Seamless certificate rotation
- **No Manual Intervention**: Fully automated

#### 3. Interoperability

- **Universal Access**: Certificates available via HTTPS
- **No Integration Required**: Standard HTTP/HTTPS protocol
- **Cross-Carrier**: Works with any carrier's system
- **Global**: Works across borders

#### 4. Security

- **HTTPS Only**: Certificates served over secure connection
- **Certificate Validation**: Standard X.509 validation
- **Revocation Support**: Can integrate OCSP/CRL
- **Key Protection**: Keys stored securely

## Implementation Details

### Certificate Distribution

Certificates are published at a well-known location:

```
https://ims.local/.well-known/stir/cert.pem
```

This follows the standard pattern used by:
- Let's Encrypt
- ACME protocol
- Web PKI

### Certificate Format

- **Format**: X.509 PEM
- **Algorithm**: ECDSA P-256
- **Key Usage**: Digital Signature
- **Extended Key Usage**: Server Auth, Client Auth

### Token Format

PASSporT tokens (RFC 8225):
- **Format**: JWT (JSON Web Token)
- **Algorithm**: ES256 (ECDSA P-256 with SHA-256)
- **Header**: Contains x5u (certificate URL)
- **Claims**: orig, dest, attest, iat, exp

### Verification Process

1. **Extract Identity Header**: From SIP INVITE
2. **Parse JWT**: Decode token
3. **Get Certificate URL**: From x5u header
4. **Fetch Certificate**: Via HTTPS
5. **Verify Signature**: Using certificate public key
6. **Validate Claims**: Check expiration, attestation, etc.

## Interoperability Testing

### Test Scenarios

1. **Cross-Carrier Verification**
   - Carrier A signs with ACME certificate
   - Carrier B verifies using standard HTTPS fetch
   - ✅ Works with any STIR/SHAKEN verifier

2. **Certificate Renewal**
   - Certificate expires
   - ACME automatically renews
   - New certificate published
   - ✅ Zero service interruption

3. **Multiple ACME Providers**
   - Let's Encrypt for production
   - ZeroSSL for staging
   - Custom ACME server for private networks
   - ✅ All work the same way

4. **Legacy System Integration**
   - Legacy systems can fetch certificates via HTTPS
   - No special integration required
   - ✅ Backward compatible

## Migration Path

### From Traditional to ACME

1. **Phase 1**: Deploy ACME-based certificate management
2. **Phase 2**: Publish certificates at well-known location
3. **Phase 3**: Update peering agreements to use HTTPS
4. **Phase 4**: Deprecate manual certificate exchange

### Backward Compatibility

- Continue supporting manual certificate exchange
- Gradually migrate to ACME
- Both methods work simultaneously
- Smooth transition

## Real-World Benefits

### Operational

- **Reduced Ops Overhead**: No manual certificate management
- **Faster Onboarding**: New carriers can verify immediately
- **Fewer Errors**: Automated process reduces mistakes
- **Better Monitoring**: Certificate status visible via HTTPS

### Technical

- **Standard Protocol**: Uses well-understood HTTPS
- **Vendor Agnostic**: Works with any implementation
- **Scalable**: Handles thousands of certificates
- **Reliable**: Automated renewal prevents outages

### Business

- **Faster Time to Market**: Quick carrier onboarding
- **Lower Costs**: Reduced operational overhead
- **Better Security**: Automated lifecycle management
- **Competitive Advantage**: Modern, interoperable solution

## Future Enhancements

1. **Certificate Transparency**: Log all certificates
2. **OCSP Integration**: Real-time revocation checking
3. **Multi-CA Support**: Support multiple certificate authorities
4. **Certificate Pinning**: Pin certificates for critical paths
5. **Automated Testing**: Continuous interoperability testing

## Conclusion

By using ACME for STIR/SHAKEN certificate management, we address the core interoperability issues:

- ✅ **Standard Protocol**: No proprietary mechanisms
- ✅ **Automated Lifecycle**: No manual intervention
- ✅ **Universal Access**: Works with any system
- ✅ **Future Proof**: Based on IETF standards

This makes our IMS core the most interoperable STIR/SHAKEN implementation available.
