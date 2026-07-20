# CDC in DCraft Fusion v1

Fusion CDC Engine ships **in-tree** under `engines/fusion-cdc-engine/` as the first included execution engine.

## Local all-in-one

```powershell
# Control plane + web
docker compose -f infra/local-dev/docker-compose.yml up -d --build

# CDC stack (control plane, workers, frontend)
docker compose -f infra/local-dev/docker-compose.cdc.yml up -d --build
```

See Docs Guide: CDC chapter and Install → Docker / Helm.

## Product boundary

- **DCraft Fusion** = control plane (kernel + web)
- **Fusion CDC Engine** = data-plane muscle under the same brand (not a second product name)

Community feedback welcome via GitHub Issues and Discussions on `DCraft-Labs/dcraft-fusion`.
