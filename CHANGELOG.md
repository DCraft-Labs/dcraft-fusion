# Changelog

All notable changes to DCraft Fusion (public repo) are documented here.
This project follows [Keep a Changelog](https://keepachangelog.com/) and
uses [Semantic Versioning](https://semver.org/).

## [1.3.6] — 2026-07-24

### Fix-forward: `dcraft-fusion` chart version + default image tags
Phase B had already published OCI `fusion-cdc:1.3.6` (committer `drainBatch: 300`,
2000m/2048Mi) and retagged local-dev / `deploy.ps1` to 1.3.6, but
`infra/helm/dcraft-fusion/Chart.yaml` was still `1.2.35` and default
`values.yaml` image tags were still `1.2.30`. The publish workflow overrode
the chart *version* to 1.3.6 at package time, so GHCR looked correct while
a plain `helm install … --version 1.3.6` would still pull stale kernel/web
images.

- Bumped `dcraft-fusion` `Chart.yaml` to `version` / `appVersion` `1.3.6`.
- Bumped default + minimal example image tags to `1.3.6`.
- Re-published both OCI charts under `v1.3.6` (same tag move pattern as
  the v1.3.5 provenance fix-forward).

## [1.2.33] — 2026-07-24

### Coordinated release with `fusion-cdc-engine` v1.2.33 — parallel-load contention + duplicate-row fixes
This release ships the chart + image-tag bumps that point the public Fusion
chart at the v1.2.33 CDC engine images, which fix three P0 bugs found in the
v1.2.32 parallel-load (K=6) test:

- **Bug #20** — unbounded final partition premature DONE due to an adaptive
  chunk-size race in `loader.py` (the unbounded branch compared `row_count`
  against the LIVE `cur_chunk_size`, which the adaptive sizer grows between
  the fetch and the check). Fixed by capturing `requested_size` before the
  fetch and comparing against that.
- **Bug #21** — Iceberg commit contention dead-lettered 4 of 6 partitions in
  ~6 minutes (optimistic-concurrency losers exhausted the 10-retry budget).
  Fixed with jittered backoff, a higher `INITIAL_LOAD_MAX_RETRIES=30` budget
  for initial-load tasks, and a per-table Redis commit mutex that serializes
  COMMITS while keeping FETCHES parallel.
- **Bug #22** — commit-batching + retry-after-conflict produced ~21%
  duplicate rows (up to 48% in the partitions that retried the most). The
  checkpoint (`last_pk`) advanced on every buffered chunk, not when the
  batch's commit succeeded; a retry resumed past durable rows and re-appended
  them. Fixed by (1) advancing the checkpoint only after commit success,
  (2) dedup-on-PK (delete-then-append) before each batch append, and
  (3) defaulting `INITIAL_LOAD_COMMIT_BATCH=1` (one commit per chunk).

See `fusion-cdc-engine` v1.2.33 CHANGELOG for full details. This repo only
bumps versions: `infra/helm/dcraft-fusion/Chart.yaml`,
`infra/local-dev/k8s/values-cdc-local.yaml`,
`infra/local-dev/k8s/values-fusion-local.yaml`, and
`infra/local-dev/k8s/deploy.ps1`.

## [1.2.30] — 2026-07-23

### Coordinated release with `fusion-cdc-engine` v1.2.30 — parallel-load correctness fix
This release ships the chart + image-tag bumps that point the public Fusion
chart at the v1.2.30 CDC engine images, which fix five P0 defects in the
multi-pod parallel initial load (confirmed live against a 118M-row MySQL
table — 6 pods, 6 disjoint PK ranges, only 380k rows / 0.32% before
plateauing with 1 of 6 checkpoints and a fake `progress_pct:100`):

- **Defect A — premature DONE on a short chunk** (`loader.py`): the
  single-partition-era `if row_count < chunk_size: break` heuristic was
  replaced with an explicit `last_pk >= pk_end` upper-bound check. A short
  chunk near a bounded partition's boundary is expected and does NOT end
  the partition; the loop continues until the boundary is reached (or the
  range is genuinely exhausted, for the unbounded last partition).
- **Defect B — missing checkpoints** (`loader.py`): the convert+write block
  is wrapped in try/except so an exception (e.g. Iceberg "snapshot id
  changed" conflict from a duplicate-dequeue sibling pod) persists a
  `state="failed"` checkpoint for the `chunk_seq` before re-raising to the
  worker retry/dead-letter path. Every exit path (normal, premature-DONE
  fix, error, exception, fetch-thread failure) now reports a checkpoint.
- **Defect C — fake `rows_estimated`** (`partitioning.py`, `connections.py`,
  `internal.py`): the per-partition row estimate is now density-based
  (`table_rows * span / total_span`, from the instant
  `information_schema.tables.table_rows` / `pg_class.reltuples` count) and
  stamped at ENQUEUE time in the task payload. The worker stamps it on the
  FIRST checkpoint for the partition and the control-plane never overwrites
  a non-null `rows_estimated`, so `progress_pct = rows_written /
  rows_estimated * 100` reflects real progress instead of always reading
  100%.
- **Defect D — duplicate dequeue** (`worker.py`): replaced the non-atomic
  `LRANGE` + `LREM` dequeue with `BLMOVE` (atomic) from the main queue to a
  per-worker in-flight list. Two pods can never dequeue the same `task_id`
  concurrently. The task is removed from in-flight only on ack (success or
  dead-letter); during retry/backoff it stays in in-flight (no sibling pod
  can grab it), then is atomically moved back to the main queue.
- **Defect E — ProxySQL pooling auth** (`loader.py`): the pooled source
  connection now uses the EXACT same param extraction as the per-chunk path
  (`database_name`/`database`, `username`/`user`, `connect_timeout=10`,
  MySQL `autocommit=True` + `DictCursor`), so ProxySQL no longer rejects
  the pooled connection with "Access denied" and the worker keeps the
  pooling win.

### Tests
Five regression tests added in
`transform-worker/tests/test_parallel_load_correctness.py` (re-exported
from `test_v130_correctness.py`):
1. `test_partition_loop_continues_past_short_chunk` — Defect A.
2. `test_all_partitions_get_checkpoint` — Defect B (K=4).
3. `test_rows_estimated_from_partitioning` — Defect C.
4. `test_no_duplicate_dequeue` — Defect D (two workers, atomic BLMOVE).
5. `test_premature_done_fix_regression` — Defect A regression (25M-key
   range, chunk_size 10k, must NOT mark DONE at 50k rows).

### Version
- `dcraft-fusion` chart: `1.2.30` (`infra/helm/dcraft-fusion/Chart.yaml`).
- `fusion-cdc` chart: `1.2.30` (`infra/helm/fusion-cdc/Chart.yaml`).
- Image tags bumped to `"1.2.30"` in both charts' `values.yaml` and
  `examples/values-minimal.yaml`, the local-dev overrides
  (`values-cdc-local.yaml`, `values-fusion-local.yaml`), and
  `infra/local-dev/k8s/deploy.ps1` (`--version 1.2.30`).

## [1.2.29] — 2026-07-23

### Coordinated release with `fusion-cdc-engine` v1.2.29
This release ships the chart + image-tag bumps that point the public Fusion
chart at the v1.2.29 CDC engine images, which contain the performance,
observability, and reliability work:

- **DuckDB native scanner for bulk initial load** (Task 1) — direct Arrow
  reads from MySQL/Postgres, gated by `INITIAL_LOAD_BULK_MODE`.
- **Per-chunk Prometheus metrics** (Task 2) — `initial_load_*` counters /
  histograms / gauges on `TRANSFORM_WORKER_PROMETHEUS_PORT`.
- **CDC streaming idempotency** (Task 4) — `cdc_applied_events` ledger gives
  exactly-once apply, gated by `CDC_IDEMPOTENCY_ENABLED`.
- **Source-DB connection pooling per partition** (Task 5).
- **Backpressure handling** (Task 6) — bounded prefetch queue + depth gauge.
- **Real-time UI progress + ETA** (Task 3) — new
  `GET /connections/{id}/initial-load/progress` endpoint + frontend polling.

### Version
- `dcraft-fusion` chart: `1.2.29` (`infra/helm/dcraft-fusion/Chart.yaml`).
- `fusion-cdc` chart: `1.2.29` (`infra/helm/fusion-cdc/Chart.yaml`).
- Image tags bumped to `"1.2.29"` in both charts' `values.yaml` and
  `examples/values-minimal.yaml`, the local-dev overrides
  (`values-cdc-local.yaml`, `values-fusion-local.yaml`), and
  `infra/local-dev/k8s/deploy.ps1` (`--version 1.2.29`).

## [1.2.28] — 2026-07-23

### Coordinated release with `fusion-cdc-engine` v1.2.28 (CI fix)
v1.2.27's tag pointed to a commit that failed a contract test (the test did a
static source check for `def _trigger_dag_or_worker(` but v1.2.27 made the
function `async def`). v1.2.28 re-tags on top of the CI fix. No production code
changed on the public repo — this is a chart/image-only bump to `1.2.28` so
the public charts pull the fixed images from GHCR.

- **Chart version + image tags bumped to `1.2.28`** across `fusion-cdc` and
  `dcraft-fusion` charts, `values-*-local.yaml`, `deploy.ps1`, and
  `examples/values-minimal.yaml`.

## [1.2.27] — 2026-07-23

### Coordinated release with `fusion-cdc-engine` v1.2.27 (P0 partitioning fix)
This is a chart/image-only release on the public repo — the partitioning fix
itself lives in the private `fusion-cdc-engine` repo. The public repo bumps
chart versions + image tags so the public `fusion-cdc` and `dcraft-fusion`
charts pull the fixed `1.2.27` images from GHCR.

- **Chart version + image tags bumped to `1.2.27`** across `fusion-cdc` and
  `dcraft-fusion` charts, `values-*-local.yaml`, `deploy.ps1`, and
  `examples/values-minimal.yaml`. See the private repo's v1.2.27 changelog
  for the full P0 fix details (non-blocking partitioning: threadpool offload
  + `information_schema` count + timeout/KILL fallback + async 202 endpoint).

## [1.2.26] — 2026-07-23

### Infrastructure
- **Configurable KEDA max scale for initial sync (Task 3):** documented the
  `transformWorker.keda.maxReplicaCount` value (already wired through the
  ScaledObject template) and added a clarifying comment explaining its
  relationship to the per-connection `resource_limits.parallelism` (K).
  Local-dev `values-cdc-local.yaml` now ships a `keda` block with
  `maxReplicaCount: 4` for documentation/local-dev.
- **Chart podSecurityContext defaults completed (Task 8):** pinned
  `runAsUser`/`fsGroup` for `controlPlane` (999), `cdcWorkers` (2000),
  `sparkConsumer` (1000), and `frontend` (101) — `transformWorker` already
  had them (v1.2.18). Without these, strict admission controllers can reject
  pods with `CreateContainerConfigError` because Kubernetes cannot verify the
  image's non-numeric `USER` directive.

### Version
- Chart version + all image tags bumped to `1.2.26` across `fusion-cdc` and
  `dcraft-fusion` charts, `values-*-local.yaml`, `deploy.ps1`, and
  `examples/values-minimal.yaml`.

## [1.2.25] — 2026-07-23

### Infrastructure
- **Remove Kafka (Task 1):** Kafka was unused dead infrastructure in the
  `fusion-cdc` Helm chart — the CDC engine uses Redis Streams (XADD/XREADGROUP)
  and KEDA scales on Redis list depth, not Kafka. Removed:
  - `templates/kafka.yaml` (deleted)
  - `kafka` keyword from `Chart.yaml` (replaced with `redis-streams`)
  - `kafka:` block from `values.yaml` and `values.schema.json`
  - `fusion-cdc.kafkaBootstrapServers` helper from `templates/_helpers.tpl`
  - `KAFKA_BOOTSTRAP_SERVERS` env injection from `templates/cdc-workers.yaml`,
    `templates/spark-consumer.yaml`, `templates/control-plane.yaml`
  - `kafka:` override block from `infra/local-dev/k8s/values-cdc-local.yaml`
  This frees ~480Mi of cluster RAM and removes a dead dependency.

### Version
- Chart version + all image tags bumped to `1.2.25` across `fusion-cdc` and
  `dcraft-fusion` charts, `values-*-local.yaml`, `deploy.ps1`, and
  `examples/values-minimal.yaml`.

## [1.2.24] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.24 — CI fix
for the v1.2.23 `test` job (missing `redis` dep for the transform-worker
unit tests, which import `loader.py` and `loader.py` does `import redis`
at module level). The public repo has no source change in this release —
only chart/image tag bumps to keep the public charts in sync with the
rebuilt GHCR images.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.24`
  / `appVersion: "1.2.24"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.23` → `1.2.24` in the values files
  (`infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, `*/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`) and
  `--version 1.2.24` in `infra/local-dev/k8s/deploy.ps1`.

## [1.2.23] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.23 — CI fix
for the v1.2.22 `test` job (missing `requests` / `pymysql` / `psycopg2-binary`
deps for the new transform-worker unit tests). The public repo has no source
change in this release — only chart/image tag bumps to keep the public charts
in sync with the rebuilt GHCR images.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.23`
  / `appVersion: "1.2.23"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.22` → `1.2.23` in the values files
  (`infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, `*/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`) and
  `--version 1.2.23` in `infra/local-dev/k8s/deploy.ps1`.

## [1.2.22] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.22 — critical
fix release for the transform-worker (Iceberg destination path). The
public repo has no source change in this release — only chart/image tag
bumps to keep the public charts in sync with the rebuilt GHCR images.

The v1.2.22 fix set (in the private repo) covers:
- **Bug A**: all-NULL columns → `pa.null()` → PyIceberg rejects — fixed by
  fetching the source schema once from `information_schema` and passing
  it through to `pa.Table.from_pylist(rows, schema=...)`.
- **Bug B**: DuckDB `$1` binding fails for `list[dict]` — fixed by
  converting rows to a PyArrow Table and registering it as a view.
- 3 additional step-handler bugs (date_op on VARCHAR, json_flatten_child
  `from_json('[]')`, mask hash `sha256(BLOB)`, udf `duckdb.create_function`
  vs `conn.create_function`).
- Compute efficiency: source schema fetched ONCE per stream (not per
  chunk), READ ONLY + autocommit on source fetches so the source DB is
  not locked across destination writes.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.22`
  / `appVersion: "1.2.22"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.21` → `1.2.22` in the values files
  (`infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, `*/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`) and
  `--version 1.2.22` in `infra/local-dev/k8s/deploy.ps1`.

## [1.2.21] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.21 — CI fix
for the v1.2.20 `test` job (`cryptography` import in the postgres
initial-load unit tests). The public repo has no source change in this
release — only chart/image tag bumps to keep the public charts in sync
with the rebuilt GHCR images.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.21`
  / `appVersion: "1.2.21"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.20` → `1.2.21` in the values files
  (`infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, `*/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`) and
  `--version 1.2.21` in `infra/local-dev/k8s/deploy.ps1`.

## [1.2.20] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.20 —
**bulletproof connection lifecycle for every source × destination combo**.
The public repo has no source change in this release — only chart/image
tag bumps to keep the public charts in sync with the rebuilt GHCR images.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.20`
  / `appVersion: "1.2.20"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.19` → `1.2.20` in the values files
  (`infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, `*/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`) and
  `--version 1.2.20` in `infra/local-dev/k8s/deploy.ps1`.

### Notes
- No source-code change in the public repo. The bulletproof routing
  (`_dest_needs_transform_worker`), the new Postgres source initial-load
  path, and the integration/contract tests live in the private
  `fusion-cdc-engine` repo. See the private `CHANGELOG.md` v1.2.20 entry
  for the full Fix A/B/C/D details and the 6-combination audit matrix.

## [1.2.19] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.19. The
private v1.2.18 wrongly deleted `cdc-workers/cdc_consumer.py` (the CDC
consumer deployed via `kubernetes/base/cdc-consumer.yaml`); v1.2.19
restores it and reverts the `snapshot_mode` default back to `inline`.
The public repo has no source change in this release — only chart/image
tag bumps to keep the public charts in sync with the rebuilt GHCR images.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.19`
  / `appVersion: "1.2.19"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.18` → `1.2.19` in the values files
  (`infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, `*/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`) and
  `--version 1.2.19` in `infra/local-dev/k8s/deploy.ps1` /
  `infra/helm/README.md`.

### Notes
- No source-code change in the public repo. The `cdc_consumer.py`
  restoration and `snapshot_mode` default revert live in the private
  `fusion-cdc-engine` repo. See the private `CHANGELOG.md` v1.2.19 entry
  for the full P0 regression fix details.

## [1.2.18] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.18. Fixes the
chart + UX issues found in the user investigation of v1.2.16 (v1.2.17 fixed
the `fetchall()` OOM regression in the transform-worker; v1.2.18 fixes the
chart bugs that prevented the transform-worker from starting at all, plus
the missing retry/snapshot-mode UX).

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.18`
  / `appVersion: "1.2.18"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.17` → `1.2.18` in the values files
  (`infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, `*/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`) and
  `--version 1.2.18` in `infra/local-dev/k8s/deploy.ps1` /
  `infra/helm/README.md`.

### Added
- **Default `podSecurityContext` for the transform-worker** in
  `infra/helm/fusion-cdc/values.yaml` — pins `runAsUser: 2001` /
  `fsGroup: 2001` to match the image's non-numeric `transform` user
  (see `docker/Dockerfile.transform-worker`). Without these, Kubernetes
  rejects the pod with `CreateContainerConfigError`. The values are
  overridable via `.Values.transformWorker.podSecurityContext`.
- **LimitRange `dcraft-local-limits`** in
  `infra/local-dev/k8s/00-infra.yaml` with `max.memory: 2Gi` (was 1 Gi
  when applied out-of-band — the 1 Gi ceiling rejected the
  transform-worker pod, which needs ~1.5 Gi for baseline Python + pyarrow
  + polars + duckdb). The LimitRange is now committed to the repo so the
  constraint is explicit and adequate.

### Removed
- `infra/local-dev/k8s/patch-cdc-worker-metadata-dsn.json` — obsolete
  after the transform-worker `worker.py` was renamed to read `DATABASE_URL`
  (the env var the chart already injects) instead of `METADATA_DB_DSN`.

### Fixed (in fusion-cdc-engine v1.2.18, referenced here for coordination)
- Removed orphaned `cdc-workers/cdc_consumer.py` (dead code — the chart's
  `Dockerfile.cdc-worker` runs `python -m cdc_worker.worker`, which never
  imported it).
- Changed the default `snapshot_mode` from `inline` (no-op, cdc_consumer.py
  removed) to `transform_worker` (canonical) in
  `control-plane/app/api/connections.py`.
- Added `POST /api/v1/connections/{id}/retry-initial-load` endpoint so
  users can re-enqueue the initial-load snapshot without deleting +
  recreating the connection.
- Renamed `METADATA_DB_DSN` → `DATABASE_URL` in
  `transform-worker/worker.py` to match the env var the chart injects.
- Added a `snapshot_mode` select field to the Iceberg destination form
  (`frontend/src/components/iceberg/IcebergDestinationForm.tsx`) and a
  "Retry Initial Load" button to the connection detail page
  (`frontend/src/pages/connections/ConnectionDetailPage.tsx`).

## [1.2.17] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.17. This repo
holds the public Helm charts and local-dev values that reference the CDC
images, so the chart + image tags move `1.2.16` → `1.2.17` to match the
transform-worker initial-load chunking fix shipped in `fusion-cdc-engine`
v1.2.17 (PK-bounded chunked fetch + checkpoint resume, fixing the OOM
regression on 2 GB+ source tables).

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.17`
  / `appVersion: "1.2.17"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.16` → `1.2.17` in the values files
  (`infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, `*/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`) and
  `--version 1.2.17` in `infra/local-dev/k8s/deploy.ps1` /
  `infra/helm/README.md`.

### Fixed (in fusion-cdc-engine v1.2.17, referenced here for coordination)
- Transform-worker `InitialLoadTask` no longer materializes the entire
  source table into memory (Postgres/MySQL `fetchall()`, MongoDB
  `list(find())`). It now loops over PK-bounded chunks of 10000 rows and
  writes a checkpoint after each chunk, so a 2 GB initial load completes
  without OOM and resumes from the last PK after a worker restart.

## [1.2.16] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.16. This repo
holds the public Helm charts and local-dev values that reference the CDC
images, so the chart + image tags move `1.2.15` → `1.2.16` to match the three
v1.2.14 gap-closing fixes shipped in `fusion-cdc-engine` v1.2.16 (initial_load
producer + MySQL/Mongo destination connector defs + InitialLoadTask direct
source fetch + /internal/load-checkpoints endpoint).

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.16`
  / `appVersion: "1.2.16"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.15` → `1.2.16` in the values files
  (`infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, `*/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`) and `--version 1.2.16`
  in `infra/local-dev/k8s/deploy.ps1` / `infra/helm/README.md`.

### Fixed (in fusion-cdc-engine v1.2.16, referenced here for coordination)
- Wired an initial_load task producer in `connections.py` that enqueues
  `initial_load` tasks to `fusion:transforms:high` when a destination's
  `connection_config.snapshot_mode` is `transform_worker` (default `inline`
  preserves the existing `cdc_consumer.py` snapshot path).
- Seeded `MySQL Destination` and `MongoDB Destination` connector definitions
  in `seed-admin.sql` so users can create MySQL/MongoDB destinations (the
  v1.2.14 DSN builders now have matching connector defs).
- Rewrote `InitialLoadTask._fetch_rows` to connect to the source DB directly
  (psycopg2/pymysql/pymongo) instead of the non-existent
  `/internal/data-proxy/fetch` endpoint; added `POST /internal/load-checkpoints`
  for chunk progress upserts.

## [1.2.15] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.15. This repo
holds the public Helm charts and local-dev values that reference the CDC
images, so the chart + image tags move `1.2.14` → `1.2.15` to match the
Iceberg write-path fixes shipped in `fusion-cdc-engine` v1.2.15
(`TableNotFound` import drop + `s3fs` dep + catalog-type-aware warehouse
hint).

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.15`
  / `appVersion: "1.2.15"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.14` → `1.2.15` in the values files
  (`infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, `*/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`) and `--version 1.2.15`
  in `infra/local-dev/k8s/deploy.ps1` / `infra/helm/README.md`.

### Fixed (in fusion-cdc-engine v1.2.15, referenced here for coordination)
- Dropped `TableNotFound` from `iceberg_tester.py` imports (does not exist
  in pyiceberg 0.7.1; only `NoSuchTableError` is used).
- Added `s3fs==2024.6.1` to `control-plane/requirements.txt` and
  `transform-worker/requirements.txt` so PyIceberg's fsspec-based S3 FileIO
  can write to Nessie/REST + MinIO/S3 (previously
  `ModuleNotFoundError: No module named 's3fs'` on `table.append()`).
- Catalog-type-aware warehouse hint in the Create Destination form
  (`frontend/src/lib/iceberg-config.ts` +
  `IcebergDestinationForm.tsx`): Nessie/REST show the warehouse NAME,
  Hive/Glue/SQL/DynamoDB show the S3 path.

## [1.2.14] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.14. Closes the
two real gaps the v1.2.13 audit flagged: the initial-load (snapshot) path no
longer calls a non-existent control-plane endpoint to fetch the destination
DSN (it derives the DSN from the task payload, mirroring the v1.2.13 CDC
fix), and MySQL/MongoDB DSN builders have been added so non-Postgres
destinations can be routed when their connector definitions are added.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.14`
  / `appVersion: "1.2.14"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.13` → `1.2.14` in the values files and
  `--version 1.2.14` in `deploy.ps1` / `infra/helm/README.md`.

### Added
- `.tmp/v114-verify/verify_v114.py` (in the private repo) — live E2E
  verification script the operator runs after deploying v1.2.14. Extends
  the v1.2.13 script with DSN-builder unit assertions (Postgres / MySQL /
  MongoDB / unknown) and a Postgres initial-load (snapshot) E2E that
  creates a fresh connection and polls until the snapshot completes.

## [1.2.13] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.13. Closes the
three honest gaps the v1.2.12 worker audit flagged: Postgres-bound CDC now
derives `dest_dsn` from the destination block, the Create Destination wizard
actually calls the validate-write endpoint before allowing finish, and a
live E2E verification script is included for post-deploy validation.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.13`
  / `appVersion: "1.2.13"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.12` → `1.2.13` in the values files and
  `--version 1.2.13` in `deploy.ps1` / `infra/helm/README.md`.

### Added
- `infra/local-dev/v113-verify/verify_v113.py` (in the private repo) — live
  E2E verification script the operator runs after deploying v1.2.13 to
  confirm Iceberg test-conn, validate-write, and Postgres → Postgres CDC.

## [1.2.12] — 2026-07-23

Coordinated release with the private `fusion-cdc-engine` v1.2.12. This repo
holds the public Helm charts and local-dev values that reference the CDC
images, so the chart + image tags move `1.2.11` → `1.2.12` to match the
CDC fixes shipped in `fusion-cdc-engine` v1.2.12.

### Fixed (in fusion-cdc-engine v1.2.12, referenced here for coordination)
- Real Iceberg Test Connection (PyIceberg catalog load + namespace list +
  S3 HeadBucket) in the control plane.
- Real Iceberg write-permission check (create/insert/delete/drop a test
  namespace + table).
- CDC → transform-worker bridge: cdc-worker now LPUSHes `cdc_transform`
  tasks to `fusion:transforms:normal` so the transform-worker BRPOPs them
  and writes to Postgres / Iceberg (previously CDC events went to a Redis
  Stream the transform-worker never read).
- Seeded Iceberg `auth_mode: "static"` now resolves to the `s3_*` creds.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.12`
  / `appVersion: "1.2.12"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.11` → `1.2.12` in the values files and
  `--version 1.2.12` in `deploy.ps1` / `infra/helm/README.md`.

## [1.2.11] — 2026-07-22

Follow-up to v1.2.10. The v1.2.10 tag shipped corrupted Helm `values.yaml`
files (a PowerShell re-encode pass introduced UTF-8 mojibake on the em-dash
and box-drawing comment characters), which made `helm lint` fail in the
`Publish Helm charts` workflow. v1.2.11 re-applies the v1.2.10 changes with
UTF-8-safe file writes and re-tags all images/charts to `1.2.11`.

### Fixed
- Re-encoded `infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, the two `examples/values-minimal.yaml`
  overlays, `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`,
  `infra/local-dev/k8s/deploy.ps1`, and `infra/helm/README.md` as UTF-8 (no
  BOM) so `helm lint` passes. No semantic changes vs v1.2.10 — only the
  image/chart version tags move from `1.2.10` → `1.2.11`.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.11`
  / `appVersion: "1.2.11"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.10` → `1.2.11` in the values files and
  `--version 1.2.11` in `deploy.ps1` / `infra/helm/README.md`.

## [1.2.10] — 2026-07-22

Coordinated release with the private `fusion-cdc-engine` v1.2.10. v1.2.9
verified the UI was live and solid, but the E2E CDC audit found the runtime
non-functional. v1.2.10 ships the local-dev infrastructure pods and the
re-tagged Helm charts that point at the v1.2.10 CDC images.

### Added
- **Local-dev infra pods (BLOCKER 7)** — `infra/local-dev/k8s/00-infra.yaml`
  now deploys `mysql-source` (MySQL 8.0 with binlog + self-seeding init SQL
  for customers/products/orders), `mongo-source` (MongoDB 7, no auth),
  `minio` (S3-compatible API :9000 + console :9001 + PVC + bucket-init Job
  that creates `iceberg-warehouse`), and `nessie` (Iceberg REST catalog
  :19120 + mgmt :19121). All with startup/readiness probes and small
  (128Mi–256Mi) resource requests. LOCAL-DEV ONLY — production must use
  managed equivalents (RDS, DocumentDB, S3, Nessie/Polaris, etc.).

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.10`
  / `appVersion: "1.2.10"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.9` → `1.2.10` in
  `infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`, the `examples/values-minimal.yaml`
  overlays, and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.10` in `infra/local-dev/k8s/deploy.ps1` and the
  `infra/helm/README.md` install snippets.

### Coordinated
- The CDC runtime fixes (worker assignment, MongoDB connections, Iceberg
  seed, pg discovery, UI form fields, port defaults) ship in the private
  `fusion-cdc-engine` repo's v1.2.10 release. This public repo re-tags the
  Fusion web + control-plane-kernel images to `1.2.10` and adds the
  local-dev infra pods so the seeded Mongo/MySQL/MinIO/Nessie hostnames
  resolve.

## [1.2.9] — 2026-07-22

UI polish pass on top of the verified-stable v1.2.8 backend. The v1.2.8 UX
audit (`E:\Dcraft\.tmp\v128-ux-audit\REPORT_v128_ux_audit.md`) found that
the Fusion workspace nav active state was invisible (no CSS for
`.nav-item.active`), the Audit Center table had no timestamp column, and
Recent Runs status had no semantic color. This release fixes those plus
the emilkowalski button `:active` scale finding. Coordinated with the
private `fusion-cdc-engine` v1.2.9 release which fixes the superadmin
role mapping (BLOCKER), the connectors "Used by" count, the localhost
GraphQL link, and adds keyboard-accessible dropdowns + button active scale
on the CDC frontend.

### Changed (Fusion public repo)
- Added `.nav-item` / `.nav-item.active` / `.nav-item:hover` styles in
  `apps/web/src/styles.css` — bolder font weight, brand-tinted background,
  and a 2px left accent border so the active workspace page is clearly
  distinguishable.
- Added `:active { transform: scale(0.97) }` to `.btn`, `.primary-button`,
  `.secondary-button`, `.primary-link`, `.secondary-link` in
  `apps/web/src/styles.css` (emilkowalski polish — buttons no longer feel
  dead on click).
- Added a "Timestamp" column to the Audit Center table
  (`apps/web/src/ui/App.tsx` `AuditCenter`) showing relative time
  ("2 min ago") with the absolute timestamp on hover via the `title`
  attribute. `AuditEvent.occurredAt` added to
  `apps/web/src/ui/workspace.ts` (optional) and populated in
  `apps/web/src/ui/demo-workspace.ts`. `DataTable` cell type relaxed from
  `string` to `ReactNode` so cells can carry hover titles.
- Added semantic color badges to the Recent Runs status column
  (`apps/web/src/ui/App.tsx` `RunCenter`): `succeeded` → green
  (`--success`/`--success-text`), `failed`/`canceled` → red
  (`--error-bg`/`--error-text`), `running` → brand blue/teal
  (`--brand-soft`/`--brand-strong`), `queued`/`unknown` → muted gray.
  New `.status[data-tone="success|failed|running|pending"]` rules in
  `apps/web/src/styles.css` use existing design tokens (no hardcoded hex).
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.9`
  / `appVersion: "1.2.9"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.8` → `1.2.9` in
  `infra/helm/*/values.yaml`, `infra/helm/*/examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.9` in `infra/local-dev/k8s/deploy.ps1` and the
  helm install examples in `infra/helm/README.md`.

### Notes
- The CDC-side fixes (superadmin role mapping, connector usage_count,
  localhost GraphQL link, button active scale, keyboard-accessible
  dropdowns) ship in the private `fusion-cdc-engine` repo's v1.2.9
  release.
- No backend changes in this public repo — the v1.2.8 control-plane
  kernel is verified live and stable; v1.2.9 is UI + chart re-tag only.

## [1.2.8] — 2026-07-22

Coordinated release with the private `fusion-cdc-engine` v1.2.8. The
v1.2.7 live stability audit found that `/api/v1/alerts/suppressions`
still returned HTTP 500 (the v1.2.5 `f7a8b9c0d1e2` migration only added
`rule_ids`/`connection_ids`; the `is_recurring`, `recurrence_pattern`,
and `updated_by` columns the `AlertSuppression` model declares were
never created) and `/api/v1/data-quality/templates` returned HTTP 501
from a stub handler. The CDC-side fixes ship in the private repo's
v1.2.8 release (new migration `a8b9c0d1e2f3` + templates stub replaced
with an empty 200 response). This public release re-tags all images
and charts to `1.2.8` so the public `fusion-cdc` chart deploys the
fixed CDC images.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.8`
  / `appVersion: "1.2.8"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.7` → `1.2.8` in
  `infra/helm/*/values.yaml`, `infra/helm/*/examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.8` in `infra/local-dev/k8s/deploy.ps1` and the
  helm install examples in `infra/helm/README.md`.

### Notes
- The CDC-side code fixes (alert_suppressions migration
  `a8b9c0d1e2f3`, data-quality/templates 501 → 200, FastAPI app version
  bump) ship in the private `fusion-cdc-engine` repo's v1.2.8 release.
- pg-source discovery returning 0 tables is operational and already
  documented in `docs/POST_DEPLOY_CHECKLIST.md` §2 — no code change.

## [1.2.7] — 2026-07-22

Follow-up to v1.2.6. The v1.2.6 `Publish images` and `Publish Helm charts`
workflows succeeded (the `dcraft-fusion-web` image was rebuilt for the
first time since v1.2.3), but the `CI` workflow failed at
`npm audit --audit-level=high` on newly-disclosed high-severity
dev-server-only vulnerabilities in transitive dev dependencies. This
release pins the patched versions so the CI gate is green.

### Fixed (CI — npm audit)
- **vite `server.fs.deny` bypass on Windows** (GHSA-fx2h-pf6j-xcff,
  CVE-2026-53571, high): the VitePress docs workspace pulled a nested
  `vite <=6.4.2` that had no auto-fix (VitePress 1.6.4 declares
  `vite: ^5.4.14`, and the advisory is only patched in `vite 6.4.3` /
  `7.3.5` / `8.0.16` — there is no 5.x fix). Added an npm `overrides`
  block in the root `package.json` forcing `vite` to `6.4.3` for both
  the web workspace and VitePress's nested dependency, and regenerated
  `package-lock.json`. VitePress 1.6.4 is compatible with Vite 6.4.3
  (docs build verified). This is a dev-server-only issue (Windows
  `vite dev --host` with sensitive files in `server.fs.allow`); it does
  not affect `vite build` or the shipped static bundle in the Docker
  image.
- `npm audit fix` also pulled in patched versions for the
  `@babel/core` (GHSA-4x5r-pxfx-6jf8), `form-data`
  (GHSA-hmw2-7cc7-3qxx), and `ws` (GHSA-96hv-2xvq-fx4p) advisories
  that had fixes available.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.7`
  / `appVersion: "1.2.7"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.6` → `1.2.7` in
  `infra/helm/*/values.yaml`, `infra/helm/*/examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.7` in `infra/local-dev/k8s/deploy.ps1` and the
  helm install examples in `infra/helm/README.md`.

### Notes
- The v1.2.6 images remain valid in GHCR; v1.2.7 re-publishes identical
  source under a new tag so the CI gate is green for the release.
- The CDC-side version bump (`control-plane/app/main.py` and
  `helm/fusion-cdc/Chart.yaml`) ships in the private `fusion-cdc-engine`
  repo's v1.2.7 release.

## [1.2.6] — 2026-07-22

Follow-up to v1.2.5. The v1.2.4 and v1.2.5 `Publish images` and `CI`
workflows failed, so the `dcraft-fusion-web` image was never rebuilt and
the deployed frontend stayed on the v1.2.3 bundle. This release fixes the
two CI failures and re-publishes the images.

### Fixed (Fusion SPA — CI failures)
- **`auth.ts` UUID polyfill type** (`apps/web/src/ui/auth.ts:5-14`): the
  `crypto.randomUUID` polyfill returned `string` but `Crypto.randomUUID`
  is typed as the template-literal UUID type, so `tsc` failed with
  `TS2322` and the `Publish images` workflow could not build the web
  image. The polyfill now declares the return type as
  `` `${string}-${string}-${string}-${string}-${string}` `` and casts the
  generated string, satisfying the DOM lib type.
- **`SuperadminOverviewScreen` heading race** (`apps/web/src/ui/App.tsx:986-1050`):
  the screen rendered `<AccessPanel title="Superadmin Overview">` (an
  `<h3>`) while the platform overview was loading, then swapped to
  `<PageHeader>` (an `<h2>`) once `/api/v1/platform/overview` resolved.
  Under Node 24 / jsdom 25 the `App.test.tsx` superadmin test's
  `findByRole("heading", { name: "Superadmin Overview" })` resolved
  against the transient `<h3>`, which was then removed by the re-render,
  so `toBeInTheDocument()` failed. The screen now always renders the
  `<PageHeader>` heading and only swaps the body (loading message vs.
  metric grid / panels), so the heading element is stable across the
  loading → loaded transition. The v1.2.5 DEMO-banner and Test-Connection
  chip-flip fixes are unchanged.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.6`
  / `appVersion: "1.2.6"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.5` → `1.2.6` in
  `infra/helm/*/values.yaml`, `infra/helm/*/examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.6` in `infra/local-dev/k8s/deploy.ps1` and the
  helm install examples in `infra/helm/README.md`.

### Notes
- The CDC-side version bump (`control-plane/app/main.py` and
  `helm/fusion-cdc/Chart.yaml`) ships in the private `fusion-cdc-engine`
  repo's v1.2.6 release.

## [1.2.5] — 2026-07-22

Follow-up to v1.2.4. The v1.2.4 live verification + codebase map audit
found that the v1.2.4 SPA fixes never took effect because the deployed
bundle was a stale v1.0.0 build (the image was never rebuilt from source),
and the CDC alert_rules/suppressions migrations were incomplete. This
release fixes the SPA build pipeline, hardens the kernel HTTP API
context scoping, and makes the local-dev Postgres hostPath cross-platform.

### Fixed (Fusion SPA — root cause of v1.2.4 not taking effect)
- **SPA bundle staleness** (`apps/web/Dockerfile`): the production bundle
  is now ALWAYS rebuilt from source inside the image. Added a defensive
  `rm -rf apps/web/dist` before `npm run build` so a cached/committed
  `dist/` can never leak into production, and a build-time check that
  fails if `dist/` is older than the newest `src/` file. Added a root
  `.dockerignore` that excludes `**/dist/`, `**/node_modules/`, and
  secrets from the build context. Fixes #1 (crypto polyfill), #2 (SPA-side
  context headers), #6 (Test Connection chip), #7 (DEMO banner) — all
  were already correct in source but missing from the deployed bundle.
- **CI freshness check** (`.github/workflows/ci.yml`): the `web` job now
  verifies `apps/web/dist` is fresh relative to `src/` and fails if any
  `dist/` artifacts are committed to git.
- **Makefile target** (`Makefile`): added `verify-fresh-bundle` so the
  same check can be run locally.

### Fixed (Fusion kernel — codebase map audit)
- **`listOrganizations` context scoping** (`services/control-plane-kernel/
  internal/httpapi/server.go`): now calls `fusioncontext.FromHeaders` and
  filters the result to the caller's organization when
  `X-Fusion-Organization-Id` is present. A platform superadmin (no org
  header) still sees all organizations. Previously this endpoint was
  effectively public after the auth middleware and bypassed tenant
  scoping.
- **`listAuditEvents` superadmin-only** (`services/control-plane-kernel/
  internal/httpapi/server.go`): now requires a valid request context AND
  platform superadmin access (via `phase1.PlatformOverview`). Audit
  events are platform-wide and may contain cross-tenant data.

### Fixed (local-dev infra)
- **HostPath PV cross-platform** (`infra/local-dev/k8s/00-infra.yaml`):
  the Postgres hostPath `path` is now parameterized via the
  `${DCRAFT_PGDATA_HOSTPATH}` placeholder. `deploy.ps1` substitutes it
  from the `DCRAFT_PGDATA_HOSTPATH` env var (default
  `/var/dcraft-local/postgres-data` for Docker Desktop Linux VM).
  Operators can override to a Windows host bind (e.g.
  `E:/Dcraft/.tmp/pgdata`) without editing the YAML.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.5`
  / `appVersion: "1.2.5"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.4` → `1.2.5` in
  `infra/helm/*/values.yaml`, `infra/helm/*/examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-{cdc,fusion}-local.yaml`.
- Bumped `--version 1.2.5` in `infra/local-dev/k8s/deploy.ps1` and the
  helm install examples in `infra/helm/README.md`.

### Notes
- The CDC-side fixes (alert_rules/suppressions migrations, system_alerts
  model cleanup, JWT/encryption/worker-secret fail-fast, seed health
  surface, dq_policies router dedup) live in the private
  `fusion-cdc-engine` repo and ship in its v1.2.5 release.

## [1.2.4] — 2026-07-22

Follow-up to v1.2.3. The v1.2.3 verification + Fusion UI audit (HTTP-only,
192.168.1.10:8088) found 7 issues. This release fixes the 4 Fusion-SPA-side
issues (1 HIGH, 1 HIGH, 2 LOW). The 3 CDC-side fixes live in the private
`fusion-cdc-engine` repo (alert_rules migration, route ordering, Kafka
widget).

### Fixed
- **Login crashed on HTTP with `crypto.randomUUID is not a function`
  (HIGH).** `apps/web/src/ui/auth.ts:92` called `crypto.randomUUID()` for
  the OIDC `state` nonce, which is `undefined` in non-secure HTTP contexts
  (browsers gate `crypto.randomUUID` behind `window.isSecureContext`). Every
  HTTP deployment — including `http://192.168.1.10:8088` — threw on the
  login redirect. Added a 5-line polyfill at the top of `auth.ts` that
  falls back to a `Math.random`-based RFC-4122 v4 generator when
  `crypto.randomUUID` is missing.
- **Bearer-token API calls returned 400 in `authMode: dev` (HIGH).** The
  local-dev gateway (`infra/local-dev/k8s/values-fusion-local.yaml:12`) runs
  in `authMode: dev`, which makes the middleware a pass-through that does
  NOT inject `X-Fusion-Actor-Id` / `X-Fusion-Organization-Id` /
  `X-Fusion-Tenant-Id` / `X-Fusion-Project-Id` from JWT claims. The SPA's
  `api.ts` only sent `Authorization` + `X-Fusion-Correlation-Id`, so
  tenant-scoped endpoints rejected requests. Added SPA-side defense-in-
  depth: `api.ts` now decodes the JWT payload (base64) and sends the four
  `X-Fusion-*` context headers from claims, so the SPA works regardless of
  gateway auth mode. `values-fusion-local.yaml` is left in `dev` mode (the
  JWKS URL config is fragile; SPA-side injection is the more robust fix).
- **"Test Connection" optimistic UI masked failures (LOW).**
  `apps/web/src/ui/App.tsx:760-767` flipped the connection status chip to
  `healthy` BEFORE calling `/api/v1/connections/{id}/test`, so a failing
  backend still showed a green chip with a small easy-to-miss banner. The
  handler now marks the chip `untested` while the call is in flight, flips
  to `healthy` only on HTTP 200, flips to `blocked` on error, and renders
  the error in a `message-banner error` (red) tone with the failure reason
  prefixed by `Connection test failed:`.
- **Demo data presented as real (LOW).** When `/api/v1/bootstrap` failed,
  `App.tsx` silently fell back to the bundled `createDemoWorkspace()` mock
  (199 lines of hardcoded data in `apps/web/src/ui/demo-workspace.ts`),
  making the UI look live when the backend was down. Added a
  `bootstrapFailed` state that flips on bootstrap error and renders a red
  `DEMO DATA — backend unreachable` banner at the top of the workspace
  shell explaining the situation and that changes will not persist.

### Changed
- Bumped `dcraft-fusion` and `fusion-cdc` Helm charts to `version: 1.2.4` /
  `appVersion: "1.2.4"` (`infra/helm/*/Chart.yaml`).
- Bumped all image tags from `1.2.3` → `1.2.4` in
  `infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`,
  `infra/helm/dcraft-fusion/examples/values-minimal.yaml`,
  `infra/helm/fusion-cdc/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-fusion-local.yaml`,
  `infra/local-dev/k8s/values-cdc-local.yaml`.
- Bumped `--version 1.2.4` in `infra/local-dev/k8s/deploy.ps1` and the
  `infra/helm/README.md` install examples.

## [1.2.3] — 2026-07-21

Follow-up to v1.2.2. The v1.2.2 remote retest (192.168.1.10) confirmed the
self-healing seed worked (6 connectors, 1 source, 2 destinations, 1
connection all present), but found 4 remaining issues. This release fixes
the Fusion-kernel-side and chart-config issues; the CDC frontend / control-
plane fixes live in the private `fusion-cdc-engine` repo.

### Fixed
- **Fusion kernel JWKS misconfig — every protected kernel API route 401s
  with `Get "": unsupported protocol scheme ""` (HIGH).** The kernel runs in
  non-dev auth mode (`FUSION_AUTH_MODE=password`) but
  `FUSION_OIDC_JWKS_URL` was empty in the chart values, so the RS256 JWKS
  fetch failed for every bearer-token request. The Fusion SPA masked this
  by falling back to bundled seed data, so the UI looked alive while the
  kernel API was unusable by any real bearer client.
  - `infra/helm/dcraft-fusion/templates/configmap.yaml`: `FUSION_OIDC_JWKS_URL`
    now auto-defaults to the in-cluster kernel service URL
    (`http://<release>-control-plane-kernel:<port>/oidc/jwks`) when the
    value is empty, so RS256 verification works for any release name without
    operator intervention.
  - `infra/helm/dcraft-fusion/values.yaml`: documented the JWKS URL default
    and the `authMode` switch in `controlPlaneKernel.config`.
  - `infra/local-dev/k8s/values-fusion-local.yaml`: set
    `controlPlaneKernel.config.authMode: dev` for local-dev / community
    deployments, which bypasses the JWKS fetch entirely (per the kernel's
    auth middleware in
    `services/control-plane-kernel/internal/auth/middleware.go`). Prod
    values keep `authMode: password` and rely on the JWKS URL default.

### Changed
- **Helm chart + image tags bumped `1.2.2` → `1.2.3`** so `deploy.ps1`
  pulls the fixed images:
  `infra/helm/dcraft-fusion/Chart.yaml`,
  `infra/helm/fusion-cdc/Chart.yaml`,
  `infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`,
  `infra/helm/dcraft-fusion/examples/values-minimal.yaml`,
  `infra/helm/fusion-cdc/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-fusion-local.yaml`,
  `infra/local-dev/k8s/values-cdc-local.yaml`,
  `infra/helm/README.md`, and the two `--version` flags in
  `infra/local-dev/k8s/deploy.ps1`.
- **`infra/helm/fusion-cdc/templates/control-plane.yaml`:** wire
  `KAFKA_BOOTSTRAP_SERVERS` into the control-plane deployment env (sourced
  from the same `fusion-cdc.kafkaBootstrapServers` helper used by the
  cdc-workers and spark-consumer) so the control-plane's
  `/api/v1/monitoring/health` endpoint can probe Kafka.

## [1.2.2] — 2026-07-21

Hotfix for the v1.2.1 regression: the CDC metadata DB stayed empty
(`connector-definitions`, `sources`, `destinations`, `connections` all
returned `total: 0`) even though `deploy.ps1` reported a successful seed
and the admin user existed. Root cause and self-healing fix live in the
private `fusion-cdc-engine` repo; this public release bumps chart/image
versions so `deploy.ps1` pulls the fixed control-plane image.

### Fixed
- **CDC seed still empty after v1.2.1 (BLOCKER).** The v1.2.1
  `deploy.ps1` fail-fast caught `psql` exit codes but the seed SQL itself
  was broken: `scripts/seed-admin.sql` INSERTed into `destinations` with
  columns (`host`, `port`, `database_name`, `schema_name`, `username`,
  `password_encrypted`, `ssl_enabled`, `ssl_config`, `config`) that do
  NOT exist on the `destinations` table — those fields live inside the
  `connection_config` JSONB column. The `connections` INSERT likewise
  listed `sync_enabled`, `replication_slot`, `publication`,
  `namespace_definition`, `namespace_format`, `stream_prefix`, `config`,
  which do not exist on `connections`. Because the whole seed runs as a
  single atomic `DO $$ ... $$;` block, the first non-existent-column
  error rolled back the entire transaction — including the
  `connector_definitions` INSERTs that ran earlier in the block —
  leaving every CDC table at `total: 0`. The admin user existed only
  because it had been registered manually via `/auth/register` (the
  control-plane has NO startup admin-creation hook). Fixed in the
  private repo by rewriting the `destinations`/`connections` INSERTs
  against the real schema and, more importantly, by adding a
  self-healing seed hook to the control-plane startup
  (`control-plane/app/seed/seed_admin.py`) that re-seeds
  `connector_definitions` whenever it finds an empty DB — so the
  deployment no longer depends on `kubectl cp` succeeding.
- **`infra/local-dev/k8s/deploy.ps1` — seed path is now a fallback.** The
  `kubectl cp + psql -f` seed step is retained as a manual re-seed
  fallback, but the primary seed mechanism is now the control-plane's
  startup hook (the seed SQL is idempotent, so running both is a no-op
  when the DB is already populated). Added a note in `deploy.ps1`
  explaining the new ordering.

### Changed
- **Helm chart + image tags bumped `1.2.1` → `1.2.2`** so `deploy.ps1`
  pulls the fixed control-plane image:
  `infra/helm/dcraft-fusion/Chart.yaml`,
  `infra/helm/fusion-cdc/Chart.yaml`,
  `infra/helm/dcraft-fusion/values.yaml`,
  `infra/helm/fusion-cdc/values.yaml`,
  `infra/helm/dcraft-fusion/examples/values-minimal.yaml`,
  `infra/helm/fusion-cdc/examples/values-minimal.yaml`,
  `infra/local-dev/k8s/values-fusion-local.yaml`,
  `infra/local-dev/k8s/values-cdc-local.yaml`,
  `infra/helm/README.md`, and the two `--version` flags in
  `infra/local-dev/k8s/deploy.ps1`.

## [1.2.1] — 2026-07-21

Hotfix release addressing the three hard blockers found when verifying v1.2.0
against the remote (192.168.1.10) deployment, plus the missing Kafka
dependency for the CDC pipeline.

### Fixed
- **`infra/local-dev/k8s/deploy.ps1` — silent seed failure (BLOCKER).** The
  `kubectl cp` and `kubectl exec ... psql -f` calls were piped to
  `| Out-Null` with no `$LASTEXITCODE` check, so a failed seed left the CDC
  metadata DB unpopulated and the UI showed empty connector lists with no
  error. Removed `| Out-Null`, added `-v ON_ERROR_STOP=1` to the psql call,
  and added `if ($LASTEXITCODE -ne 0) { throw "Seed failed — aborting deploy" }`.
  Audited every other `| Out-Null` / missing exit-code check in the script
  (infra apply, postgres/redis rollout, CREATE DATABASE, helm uninstall,
  kubectl patch, control-plane/frontend/web rollout) and made them all
  fail-fast with descriptive throws.
- **`infra/local-dev/k8s/deploy.ps1` — seed path fallback masked missing
  private repo.** The script silently fell back to the admin-only
  `seed-cdc-admin.sql` stub when the full
  `.tmp/fusion-cdc-engine-private/scripts/seed-admin.sql` was missing,
  producing a CDC DB with an admin user but no connectors. Changed the
  fallback to a hard error with a clear message telling the operator to
  git-fetch / clone the private fusion-cdc-engine repo.
- **`infra/local-dev/k8s/00-infra.yaml` — postgres data wiped on every pod
  restart (BLOCKER).** Postgres ran on an `emptyDir` because the Docker
  Desktop `storage-provisioner` is in CrashLoopBackOff and dynamic PVCs
  stay Pending. Added a `hostPath` PersistentVolume (`/var/dcraft-local/
  postgres-data`) plus a matching PersistentVolumeClaim (storageClassName
  `""`, reclaimPolicy `Retain`) and wired it into the postgres deployment
  with `PGDATA=/var/lib/postgresql/data/pgdata`. Data now survives pod
  restarts. **Local-dev only** — production must use a real storage class
  or external managed Postgres.
- **`infra/local-dev/k8s/seed-connectors-job.yaml` — dead code.** The
  ConfigMap body was a placeholder comment with no SQL and the Job was
  never wired into `deploy.ps1`. Deleted the file; the fixed `deploy.ps1`
  `kubectl cp + psql -f` path is now the single source of truth for
  seeding.

### Added
- **Kafka in the public `fusion-cdc` Helm chart (BLOCKER — user #1
  priority).** The CDC architecture is Debezium-style
  (cdc-workers → Kafka → spark-consumer) but Kafka previously lived only
  in the private repo's kustomize manifests (`kubernetes/base/kafka.yaml`)
  and never reached the published Helm chart. Ported the Apache Kafka 3.7.0
  KRaft StatefulSet (no ZooKeeper) into a new
  `infra/helm/fusion-cdc/templates/kafka.yaml`, gated by a new
  `kafka.enabled` value (default `true`). Added a
  `fusion-cdc.kafkaBootstrapServers` helper and wired
  `KAFKA_BOOTSTRAP_SERVERS` env into the `cdc-workers` and `spark-consumer`
  deployment templates (defaults to the in-cluster service
  `<release>-kafka:9092` when `kafka.enabled=true`, empty otherwise so
  operators can point at external MSK / Confluent Cloud via per-workload
  `env` overrides). Added a `kafka:` block to `values.yaml` (replicas,
  bootstrapServers, image, numPartitions, replicationFactor, retention,
  heap, externalNodePort, storage, resources, security contexts) and a
  matching entry to `values.schema.json`. Local-dev overrides in
  `values-cdc-local.yaml` pin a single small broker that fits the Docker
  Desktop VM.

### Changed
- Bumped chart versions: `dcraft-fusion` `1.1.2 → 1.2.1`,
  `fusion-cdc` `1.1.2 → 1.2.1` (both `version` and `appVersion`).
- Bumped all image tags in chart `values.yaml`, `examples/values-minimal.yaml`,
  and `infra/local-dev/k8s/values-*-local.yaml` from `1.1.1`/`1.1.2` to
  `1.2.1`.
- Bumped the `--version` flags in `deploy.ps1` and the install snippets in
  `infra/helm/README.md` from `1.1.2` to `1.2.1`.

### Notes
- The `dcraft-fusion-stack` umbrella chart (`infra/helm/dcraft-fusion-stack`)
  is a thin file-dependency guide not used by `deploy.ps1`; its dependency
  version pins (`1.0.1`) are pre-existing drift and not updated in this
  release.
