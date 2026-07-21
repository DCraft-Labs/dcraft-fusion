# Changelog

All notable changes to DCraft Fusion (public repo) are documented here.
This project follows [Keep a Changelog](https://keepachangelog.com/) and
uses [Semantic Versioning](https://semver.org/).

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
