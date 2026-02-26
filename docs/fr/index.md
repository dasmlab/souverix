# Documentation Souverix
## Doctrine de signalisation souveraine — IMS + SIG-GW

Souverix est une **plateforme de signalisation souveraine**, **nativement IA**, conçue pour des environnements **opérateurs** et **défense**. Développée en Golang moderne, Souverix est une réécriture architecturale propre visant la **performance CPS**, la **maîtrise des frontières d'interconnexion**, et la **souveraineté cryptographique**.

Souverix traite la signalisation comme une infrastructure stratégique :

- **Souveraineté par conception** : racines de confiance, politiques et contrôle d'interconnexion.
- **Résilience par défaut** : modèles active/active, isolation des pannes, continuité d'urgence.
- **Intelligence à la frontière** : détection IA, classification, et application adaptative.
- **Prêt pour la conformité** : STIR/SHAKEN, interception légale, urgence, auditabilité.
- **CNF natif** : conçu pour Kubernetes/OpenShift avec une exploitation prévisible.

---

## Composantes de la Plateforme

### 🧠 Souverix Coeur — Noyau IMS
Noyau de signalisation IMS (pile X-CSCF) infonuagique, responsable du contrôle de session et de l'intégration des politiques.

### 🛡 Souverix Rempart — SIG-GW / IBCF
Passerelle de frontière opérateur/défense : contrôle d'interconnexion, dissimulation de topologie, normalisation SIP, mitigation d'abus, et application STIR/SHAKEN.

### 🎛 Souverix Relais — Plan média
Relais/anchoring média : politiques RTP/SRTP, traversée NAT, QoS, et télémétrie média.

### 🔐 Souverix Autorite — PKI / HSM / Vault
Autorité cryptographique souveraine : gestion de chaîne CA, automatisation des certificats, intégration HSM, application mTLS et rotation des clés.

### 👁 Souverix Vigie — Intelligence IA
Couche IA : détection d'anomalies, signaux de fraude, politiques adaptatives, classification d'attaques, déclencheurs d'auto-rétablissement.

### 🎯 Souverix Mandat — Interception légale
Orchestration d'interception légale : duplication signalisation/média, intégration médiation, suivi de conformité et journaux d'audit.

### 🚨 Souverix Priorite — Urgence & Services Prioritaires
Urgence et priorité nationale : routage PSAP, file prioritaire, contrôles de contournement, continuité sous stress.

### 📊 Souverix Vigile — Observabilité & Audit
Métriques, journaux, traces, télémétrie de conformité et rapports d'audit de niveau réglementaire.

### 🌐 Souverix Federation — Contrôle Inter-domaines
Interopérabilité maîtrisée entre domaines souverains : cartographie de confiance, ententes de peering, politiques multi-locataires.

### ⚙ Souverix Gouverne — Plan de Contrôle
Autorité de configuration et politiques : profils de pairs, bascules d'application, limites de débit, contrôles en temps réel, contournements d'urgence, et gestion des mandats.

---

## Parcours de Lecture Suggéré

1. **Plateforme → Vue d'ensemble / Composants / Nomenclature**
2. **Architecture → Couches / CNF OpenShift**
3. **Conformité → STIR/SHAKEN / Interception légale / Urgence**

---

## Liens Rapides

- [Vue d'ensemble plateforme](plateforme/doctrine.md)
- [Détail des composantes](plateforme/composants.md)
- [Nomenclature & namespaces](plateforme/nomenclature.md)
- [Démarrage](operations/demarrage.md)
- [Hiérarchie d'architecture](architecture/hierarchie.md)

---

## Énoncé Doctrinal

Souverix ne fait pas qu'implanter l'IMS.

Il établit une doctrine moderne de **signalisation souveraine** — où **l'interconnexion**, **la confiance**, **l'intelligence** et **la résilience** sont des préoccupations architecturales de premier ordre.

**Conçu au Canada. Bâti pour la maîtrise souveraine.**

---

## Conformité aux Normes

- **3GPP** : TS 23.228, TS 24.229, TS 29.228, TS 33.107, TS 23.167
- **IETF** : RFC 3261, RFC 8224, RFC 8225, RFC 8588, RFC 8555
- **Réglementaire** : FCC, CRTC, ETSI

---

## Fin de la Documentation
