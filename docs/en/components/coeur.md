# Souverix Coeur
IMS Core - Signaling Brain

## Overview

Souverix Coeur is the IMS Core signaling brain, containing all X-CSCF components and core session control logic.

**Coeur** = heart/core (French)

## Role

Coeur is the sovereign signaling intelligence core that orchestrates all IMS session control and service logic.

## Components

### P-CSCF (Proxy CSCF)
- First contact point for User Equipment (UE)
- SIP proxy functionality
- Security enforcement
- Emergency detection

### I-CSCF (Interrogating CSCF)
- Inter-domain routing
- HSS query and selection
- Load balancing
- Topology hiding

### S-CSCF (Serving CSCF)
- Core session control
- Service logic execution
- Subscriber management
- Routing decisions

### BGCF (Breakout Gateway Control Function)
- PSTN breakout routing
- Gateway selection
- Interconnect decisions

### MGCF (Media Gateway Control Function)
- SIP to ISUP conversion
- Media gateway control
- PSTN interworking

## Integration

Coeur integrates with:
- **Souverix Rempart** - Border control
- **Souverix Autorite** - Certificate management
- **Souverix Vigie** - AI intelligence
- **Souverix Gouverne** - Policy control

## Standards Compliance

- **3GPP TS 23.228** - IMS Architecture
- **3GPP TS 24.229** - IP multimedia call control
- **3GPP TS 29.228** - Cx and Dx interfaces

---

## End of Coeur Documentation
