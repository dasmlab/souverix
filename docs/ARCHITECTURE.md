# IMS Core Architecture

## Overview

This document describes the architecture of the IMS Core implementation, a cloud-native IP Multimedia Subsystem built with Go.

## Component Architecture

### Core IMS Components

#### P-CSCF (Proxy CSCF)
- **Status**: Scaffolded (ready for implementation)
- **Role**: First contact point for User Equipment (UE)
- **Responsibilities**:
  - SIP message forwarding
  - Security enforcement
  - NAT traversal
  - Compression/decompression

#### I-CSCF (Interrogating CSCF)
- **Status**: Scaffolded (ready for implementation)
- **Role**: Inter-domain routing and HSS query
- **Responsibilities**:
  - HSS query for S-CSCF assignment
  - Inter-domain routing
  - Topology hiding at domain boundary

#### S-CSCF (Serving CSCF)
- **Status**: Scaffolded (ready for implementation)
- **Role**: Core session control and service logic
- **Responsibilities**:
  - Session establishment and management
  - Service profile execution
  - Application Server routing
  - Registration handling

#### HSS/UDM (Home Subscriber Server / Unified Data Management)
- **Status**: ✅ Implemented (in-memory store)
- **Role**: Subscriber database
- **Responsibilities**:
  - Subscriber profile storage
  - Authentication data
  - Service profile management
  - S-CSCF assignment
  - Registration state

#### SBC/IBCF (Session Border Controller / Interconnection Border Control Function)
- **Status**: ✅ Implemented (Priority Component)
- **Role**: SIP Gateway and security boundary
- **Responsibilities**:
  - SIP normalization
  - Topology hiding
  - Security enforcement (DoS protection, rate limiting)
  - Inter-operator SIP peering
  - Enterprise SIP trunking
  - PBX to IMS interworking
  - Fixed Broadband voice to IMS

#### IBCF (Interconnection Border Control Function)
- **Status**: ✅ Implemented (3GPP TS 23.228 Compliant)
- **Role**: Standardized border control between IMS networks
- **Responsibilities**:
  - SIP signaling control (Mw/Mx reference points)
  - Topology hiding (3GPP standardized)
  - Security enforcement (TLS, DoS, policy)
  - Inter-operator peering control
  - STIR/SHAKEN integration
  - Message validation and enforcement

#### MGCF/BGCF
- **Status**: Scaffolded (future implementation)
- **Role**: PSTN interworking
- **Responsibilities**:
  - SIP to ISUP conversion
  - TDM interworking
  - Breakout gateway control

## SIP Gateway Capabilities

### Access SIP Gateway
- Enterprise SIP trunk termination
- PBX to IMS interworking
- Fixed Broadband voice to IMS

### IBCF/SBC Functions
- **Topology Hiding**: Removes internal network topology from SIP messages
- **SIP Normalization**: Standardizes SIP headers and formats
- **Security**:
  - Rate limiting per IP
  - DoS protection
  - TLS/SRTP support (planned)
- **Inter-operator Peering**: Secure SIP peering between carriers

### PSTN Gateway (Future)
- ISUP support
- TDM interworking
- Legacy switching support

## Zero Trust Architecture

### Configuration
Zero Trust Mode is enabled via `ZERO_TRUST_MODE=true` environment variable.

### CA Providers
1. **Internal CA**: Self-signed CA for development/testing
2. **Vault CA**: HashiCorp Vault PKI integration (planned)
3. **ACME CA**: Let's Encrypt and other ACME providers (planned)

### Certificate Management
- Automatic certificate generation
- Certificate rotation
- Mutual TLS support

## Data Flow

### Registration Flow
```
UE → P-CSCF → I-CSCF → HSS → S-CSCF → UE
```

### Call Flow (Intra-domain)
```
UE → P-CSCF → S-CSCF → Application Server → S-CSCF → P-CSCF → UE
```

### Call Flow (Inter-domain via SBC)
```
UE → P-CSCF → S-CSCF → SBC/IBCF → External Network
```

## Technology Stack

### Backend
- **Language**: Go 1.23
- **Web Framework**: Gin
- **Logging**: Logrus
- **Metrics**: Prometheus
- **Tracing**: OpenTelemetry

### Frontend
- **Framework**: Quasar (Vue 3)
- **Build Tool**: Vite
- **Purpose**: Fleet management and monitoring

### Deployment
- **Container Runtime**: Docker/Podman
- **Orchestration**: Kubernetes
- **Automation**: Ansible

## Security Features

1. **Rate Limiting**: Per-IP rate limiting to prevent DoS
2. **Topology Hiding**: Removes internal network information
3. **TLS Support**: SIP over TLS (planned)
4. **SRTP Support**: Secure RTP (planned)
5. **Zero Trust**: Configurable Zero Trust Architecture
6. **STIR/SHAKEN**: Caller identity authentication with ACME-based certificate management
7. **Fraud Detection**: AI agent hooks for fraud analytics

## Scalability

- **Horizontal Scaling**: Stateless components can scale horizontally
- **Session State**: Currently in-memory (can be moved to Redis/etcd)
- **Load Balancing**: Kubernetes service load balancing

## Monitoring and Observability

- **Metrics**: Prometheus metrics exposed on `/metrics`
- **Tracing**: OpenTelemetry distributed tracing
- **Logging**: Structured JSON logging with Logrus
- **Health Checks**: `/health` endpoint for liveness/readiness

## Future Enhancements

1. **Diameter Protocol**: Full Diameter Cx/Dx/Sh interfaces
2. **PSTN Gateway**: Complete ISUP/TDM support
3. **AI Agent Integration**: MCP and extensibility hooks (✅ hooks implemented)
4. **Persistent Storage**: PostgreSQL/Redis for HSS
5. **STIR/SHAKEN**: Call authentication (✅ implemented with ACME)
6. **5G Integration**: NRF, NEF integration
7. **Full ACME Client**: Complete RFC 8555 implementation
8. **HSM Integration**: Hardware security module support
