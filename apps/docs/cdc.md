# CDC

Change Data Capture (CDC) is included in the **Community** edition of DCraft Fusion. You do not need Enterprise to run the Fusion CDC engine.

## How CDC fits

1. Operators register sources and pipelines through the Fusion control plane.
2. The **Fusion CDC engine** captures inserts, updates, and deletes from supported databases.
3. Events stream to warehouses, analytics systems, or other consumers with checkpoints and operational controls.
4. Run state, metadata, and audit remain visible in the control plane.

CDC is one muscle under the Fusion control plane — not a separate product license for Community users.

## Engine documentation

Authoritative CDC docs live with the engine:

**[engines/fusion-cdc-engine/docs](https://github.com/DCraft-Labs/dcraft-fusion/tree/main/engines/fusion-cdc-engine/docs)**

Start at the index:

- [00_INDEX.md](https://github.com/DCraft-Labs/dcraft-fusion/blob/main/engines/fusion-cdc-engine/docs/00_INDEX.md) — documentation map
- [01_ARCHITECTURE.md](https://github.com/DCraft-Labs/dcraft-fusion/blob/main/engines/fusion-cdc-engine/docs/01_ARCHITECTURE.md) — CDC architecture
- [02_DATA_FLOW.md](https://github.com/DCraft-Labs/dcraft-fusion/blob/main/engines/fusion-cdc-engine/docs/02_DATA_FLOW.md) — pipelines and event flow
- [04_OPERATIONS_RUNBOOK.md](https://github.com/DCraft-Labs/dcraft-fusion/blob/main/engines/fusion-cdc-engine/docs/04_OPERATIONS_RUNBOOK.md) — operations

Helm chart source for CDC: `engines/fusion-cdc-engine/helm/fusion-cdc` (also published under `oci://ghcr.io/dcraft-labs/charts/...`).

## Open core note

Community includes CDC integration. Enterprise gates identity, isolation, and commercial offerings — not CDC itself. See [Open core](/open-core).
