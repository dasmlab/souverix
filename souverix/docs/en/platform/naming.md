# Souverix Naming Conventions

## Component Naming

All Souverix components follow a consistent naming pattern:

- **Full Name**: `Souverix <Component>`
- **Short Name**: `sx-<component>`
- **API Path**: `/api/v1/<component>`

## Component Names

| Full Name | Short Name | API Path | Description |
|-----------|------------|----------|-------------|
| Souverix Coeur | sx-coeur | /api/v1/coeur | IMS Core |
| Souverix Rempart | sx-rempart | /api/v1/rempart | SIG-GW/IBCF |
| Souverix Relais | sx-relais | /api/v1/relais | Media Plane |
| Souverix Autorite | sx-autorite | /api/v1/autorite | PKI/HSM/Vault |
| Souverix Vigie | sx-vigie | /api/v1/vigie | AI Intelligence |
| Souverix Mandat | sx-mandat | /api/v1/mandat | Lawful Intercept |
| Souverix Priorite | sx-priorite | /api/v1/priorite | Emergency Services |
| Souverix Vigile | sx-vigile | /api/v1/vigile | Observability |
| Souverix Federation | sx-federation | /api/v1/federation | Inter-domain |
| Souverix Gouverne | sx-gouverne | /api/v1/gouverne | Policy Control |

## Kubernetes Resources

### Service Names
```yaml
sx-coeur-service
sx-rempart-service
sx-relais-service
# etc.
```

### Deployment Names
```yaml
sx-coeur
sx-rempart
sx-relais
# etc.
```

### ConfigMap Names
```yaml
sx-coeur-config
sx-rempart-config
# etc.
```

## Container Images

```bash
souverix/coeur:latest
souverix/rempart:latest
souverix/relais:latest
# etc.
```

## Environment Variables

```bash
SOUVERIX_COEUR_ENABLED=true
SOUVERIX_REMPART_ENABLED=true
SOUVERIX_AUTORITE_ENABLED=true
# etc.
```

## Naming Principles

1. **ASCII-clean**: No accents in production names
2. **Consistent**: Same pattern across all components
3. **Clear**: Names reflect component function
4. **Scalable**: Pattern works for future components

---

## End of Naming Conventions
