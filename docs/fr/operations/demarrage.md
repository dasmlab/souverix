# Démarrage avec Souverix

## Démarrage Rapide

### Prérequis

- Go 1.23.0 ou plus récent
- Docker ou Podman
- Cluster Kubernetes/OpenShift (pour le déploiement)

### Construction

```bash
./buildme.sh
```

### Exécution Locale

```bash
./runme-local.sh
```

### Push vers le Registre

```bash
export GITHUB_TOKEN=votre_token
./pushme.sh
```

## Configuration

### Configuration de Base

Définir les variables d'environnement :

```bash
export SOUVERIX_COEUR_ENABLED=true
export SOUVERIX_REMPART_ENABLED=true
export SOUVERIX_AUTORITE_ENABLED=true
```

### Mode Zero Trust

Activer l'architecture Zero Trust :

```bash
export ZERO_TRUST_MODE=true
export AUTORITE_VAULT_ENABLED=true
```

### STIR/SHAKEN

Activer STIR/SHAKEN :

```bash
export REMPART_STIR_ENABLED=true
export REMPART_STIR_ATTESTATION=auto
```

## Déploiement

### Kubernetes

```bash
kubectl apply -f k8s/
```

### OpenShift

```bash
oc apply -f k8s/
```

## Prochaines Étapes

- [Guide de Configuration](configuration.md)
- [Aperçu de l'Architecture](../architecture/hierarchie.md)
- [Documentation des Composants](../plateforme/composants.md)

---

## Fin du Guide de Démarrage
