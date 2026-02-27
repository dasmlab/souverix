# Souverix Repository Skeleton

---

## ENGLISH

# Repository Structure

```
souverix/
├── cmd/                    # Executable services
│   ├── coeur/             # Souverix Coeur (IMS Core)
│   ├── rempart/           # Souverix Rempart (SIG-GW/IBCF)
│   ├── relais/            # Souverix Relais (Media Plane)
│   ├── autorite/          # Souverix Autorite (PKI/HSM/Vault)
│   ├── vigie/             # Souverix Vigie (AI Intelligence)
│   ├── mandat/            # Souverix Mandat (Lawful Intercept)
│   ├── priorite/          # Souverix Priorite (Emergency Services)
│   ├── vigile/            # Souverix Vigile (Observability)
│   ├── federation/        # Souverix Federation (Inter-domain)
│   └── gouverne/          # Souverix Gouverne (Policy Control)
│
├── internal/              # Internal packages
│   ├── signaling/         # SIP/Diameter signaling
│   ├── media/             # RTP/SRTP media handling
│   ├── security/          # Security and cryptography
│   ├── ai/                # AI integration (MCP, hooks)
│   ├── policy/            # Policy engine
│   └── compliance/        # Compliance and audit
│
├── api/                   # API definitions
│   ├── v1/                # REST API v1
│   └── grpc/              # gRPC service definitions
│
├── proto/                 # Protocol buffers
│   ├── signaling/         # SIP/Diameter protos
│   ├── media/             # Media control protos
│   └── control/           # Control plane protos
│
├── deployments/           # Deployment manifests
│   ├── kubernetes/        # Kubernetes manifests
│   └── openshift/         # OpenShift CNF manifests
│
├── test/                  # Test suites
│   ├── e2e/               # End-to-end tests
│   ├── performance/       # Performance/load tests
│   ├── chaos/             # Chaos engineering tests
│   └── compliance/        # Compliance validation tests
│
├── docs/                  # Documentation
│   ├── SOUVERIX_PLATFORM.md
│   ├── MANIFESTO.md
│   ├── DOCTRINE.md
│   └── ...
│
├── scripts/               # Build and deployment scripts
│   ├── buildme.sh
│   ├── pushme.sh
│   └── runme-local.sh
│
├── buildme.sh            # Build script
├── pushme.sh             # Push script with SemVer
├── runme-local.sh        # Local run script
├── Dockerfile            # Multi-stage Dockerfile
├── go.mod                # Go module definition
└── README.md             # Main README
```

---

## FRANÇAIS (FR-CA)

# Structure du Dépôt

```
souverix/
├── cmd/                   # Services exécutables
│   ├── coeur/            # Souverix Coeur (Cœur IMS)
│   ├── rempart/          # Souverix Rempart (SIG-GW/IBCF)
│   ├── relais/           # Souverix Relais (Plan média)
│   ├── autorite/         # Souverix Autorite (PKI/HSM/Vault)
│   ├── vigie/            # Souverix Vigie (Intelligence IA)
│   ├── mandat/           # Souverix Mandat (Interception légale)
│   ├── priorite/         # Souverix Priorite (Services d'urgence)
│   ├── vigile/           # Souverix Vigile (Observabilité)
│   ├── federation/       # Souverix Federation (Inter-domaine)
│   └── gouverne/         # Souverix Gouverne (Contrôle politique)
│
├── internal/             # Packages internes
│   ├── signaling/        # Signalisation SIP/Diameter
│   ├── media/            # Gestion média RTP/SRTP
│   ├── security/         # Sécurité et cryptographie
│   ├── ai/               # Intégration IA (MCP, hooks)
│   ├── policy/           # Moteur de politique
│   └── compliance/       # Conformité et audit
│
├── api/                  # Définitions d'API
│   ├── v1/               # API REST v1
│   └── grpc/             # Définitions de services gRPC
│
├── proto/                # Protocoles buffers
│   ├── signaling/        # Protos SIP/Diameter
│   ├── media/            # Protos contrôle média
│   └── control/          # Protos plan de contrôle
│
├── deployments/          # Manifests de déploiement
│   ├── kubernetes/       # Manifests Kubernetes
│   └── openshift/        # Manifests OpenShift CNF
│
├── test/                 # Suites de tests
│   ├── e2e/              # Tests bout en bout
│   ├── performance/      # Tests de performance/charge
│   ├── chaos/            # Tests d'ingénierie du chaos
│   └── compliance/       # Tests de validation de conformité
│
├── docs/                 # Documentation
│   ├── SOUVERIX_PLATFORM.md
│   ├── MANIFESTO.md
│   ├── DOCTRINE.md
│   └── ...
│
├── scripts/              # Scripts de construction et déploiement
│   ├── buildme.sh
│   ├── pushme.sh
│   └── runme-local.sh
│
├── buildme.sh           # Script de construction
├── pushme.sh            # Script de push avec SemVer
├── runme-local.sh       # Script d'exécution locale
├── Dockerfile           # Dockerfile multi-étapes
├── go.mod               # Définition du module Go
└── README.md            # README principal
```

---

## End of Repository Structure
