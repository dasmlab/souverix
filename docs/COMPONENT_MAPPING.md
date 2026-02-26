# Souverix Component Mapping
Legacy to Souverix Architecture Mapping

This document maps the original component names to the new Souverix architecture.

---

## Component Name Mapping

| Legacy Name | Souverix Name | Description |
|-------------|---------------|-------------|
| IMS Core | **Souverix Coeur** | IMS Core (X-CSCF stack) |
| SIG-GW / IBCF / SBC | **Souverix Rempart** | Border control and SIP gateway |
| Media Relay | **Souverix Relais** | Media plane and RTP anchoring |
| PKI / CA / Vault | **Souverix Autorite** | Cryptographic authority |
| AI Intelligence | **Souverix Vigie** | AI-driven intelligence engine |
| Lawful Intercept | **Souverix Mandat** | LI orchestration |
| Emergency Services | **Souverix Priorite** | Emergency and priority services |
| Observability | **Souverix Vigile** | Monitoring and audit |
| Federation | **Souverix Federation** | Inter-domain control |
| Policy Engine | **Souverix Gouverne** | Policy and control plane |

---

## Code Structure Mapping

### Souverix Coeur (IMS Core)
```
cmd/ims/              → Souverix Coeur main
internal/sip/         → SIP handling (used by all)
internal/store/hss.go → HSS/UDM
pkg/ims/              → IMS core types
```

### Souverix Rempart (SIG-GW / IBCF)
```
internal/sbc/         → Souverix Rempart core
internal/ibcf/        → IBCF functionality
internal/sbc/ratelimiter.go → DoS protection
internal/sbc/stir.go  → STIR/SHAKEN integration
internal/sbc/emergency.go → Emergency handling
```

### Souverix Relais (Media Plane)
```
internal/media/       → Media relay (future)
internal/rtp/         → RTP handling (future)
```

### Souverix Autorite (PKI / HSM / Vault)
```
internal/zta/         → Zero Trust / PKI
internal/stir/acme_cert.go → Certificate management
docs/VAULT_INTEGRATION.md → Vault integration
```

### Souverix Vigie (AI Intelligence)
```
internal/ai/          → AI hooks and integration
internal/ai/hooks.go   → MCP integration
```

### Souverix Mandat (Lawful Intercept)
```
internal/li/          → Lawful Intercept
internal/li/intercept.go → Intercept controller
```

### Souverix Priorite (Emergency Services)
```
internal/emergency/   → Emergency services
internal/emergency/emergency.go → Emergency routing
```

### Souverix Vigile (Observability)
```
internal/metrics/     → Prometheus metrics
internal/logutil/      → Logging
internal/diagnostics/  → Diagnostic APIs
```

### Souverix Gouverne (Policy & Control)
```
internal/config/      → Configuration management
internal/ibcf/policy.go → Policy engine
```

### Souverix Federation (Inter-domain)
```
internal/federation/  → Federation control (future)
```

---

## Configuration Mapping

### Environment Variables

| Legacy | Souverix | Component |
|--------|----------|-----------|
| `SBC_*` | `REMPART_*` | Souverix Rempart |
| `STIR_*` | `REMPART_STIR_*` | STIR in Rempart |
| `LI_*` | `MANDAT_*` | Souverix Mandat |
| `EMERGENCY_*` | `PRIORITE_*` | Souverix Priorite |
| `ZERO_TRUST_*` | `AUTORITE_*` | Souverix Autorite |
| `AI_*` | `VIGIE_*` | Souverix Vigie |

---

## Documentation Mapping

| Legacy Document | Souverix Document |
|-----------------|-------------------|
| `ARCHITECTURE.md` | `SOUVERIX_PLATFORM.md` |
| `SIP_GATEWAY.md` | `REMPART.md` (future) |
| `STIR_SHAKEN.md` | `REMPART_STIR.md` (future) |
| `LI_EMERGENCY_*.md` | `MANDAT_PRIORITE.md` (future) |
| `VAULT_INTEGRATION.md` | `AUTORITE.md` (future) |

---

## Deployment Mapping

### Kubernetes Resources

| Legacy | Souverix |
|--------|----------|
| `ims-core` | `souverix-coeur` |
| `sbc` | `souverix-rempart` |
| `media-relay` | `souverix-relais` |
| `vault` | `souverix-autorite` |
| `ai-engine` | `souverix-vigie` |
| `li-controller` | `souverix-mandat` |
| `emergency-router` | `souverix-priorite` |
| `monitoring` | `souverix-vigile` |

---

## API Endpoints Mapping

| Legacy Path | Souverix Path | Component |
|-------------|---------------|-----------|
| `/health` | `/health` | All components |
| `/diagnostics/*` | `/vigile/*` | Souverix Vigile |
| `/metrics` | `/vigile/metrics` | Souverix Vigile |
| `/api/v1/stir/*` | `/rempart/stir/*` | Souverix Rempart |
| `/api/v1/li/*` | `/mandat/*` | Souverix Mandat |
| `/api/v1/emergency/*` | `/priorite/*` | Souverix Priorite |

---

## Container Images

| Legacy Image | Souverix Image |
|--------------|----------------|
| `ims-core:latest` | `souverix-coeur:latest` |
| `sbc:latest` | `souverix-rempart:latest` |
| `media-relay:latest` | `souverix-relais:latest` |
| `vault-integration:latest` | `souverix-autorite:latest` |
| `ai-engine:latest` | `souverix-vigie:latest` |
| `li-controller:latest` | `souverix-mandat:latest` |
| `emergency-router:latest` | `souverix-priorite:latest` |
| `monitoring:latest` | `souverix-vigile:latest` |

---

## Migration Notes

1. **Gradual Migration**: Components can be migrated incrementally
2. **Backward Compatibility**: Legacy names supported during transition
3. **Configuration**: Environment variables can use either naming
4. **Documentation**: Both naming conventions documented
5. **API**: Legacy endpoints remain for compatibility

---

## End of Component Mapping
