# Backup And Restore Runbook

## Scope

This runbook covers the Phase 1 control-plane datastore. It backs up metadata, tenant hierarchy, audit events, platform-admin settings, async job state persisted to Postgres, and Phase 1 resources. It does not move customer raw data.

## Backup

1. Run `pg_dump "$POSTGRES_DSN" --format=custom --file=fusion-control-plane.dump`.
2. Store the dump in the configured tenant-approved backup target.
3. Record the artifact checksum, database version, app image tags, and Kubernetes manifest revision.
4. Keep Redis as recoverable job cache only; Postgres remains the durable source of truth.

## Verification

The Kubernetes `postgres-backup-verify` CronJob performs a daily schema dump and read verification. A production target must extend it to restore the dump into an isolated verification database and run:

```powershell
psql $env:RESTORE_DSN -c "select count(*) from organizations;"
psql $env:RESTORE_DSN -c "select kind, count(*) from phase1_resources group by kind;"
```

## Restore

1. Scale API workers to zero.
2. Restore the verified dump into the target Postgres instance.
3. Run the control-plane process once to apply startup migrations.
4. Scale API workers back to the configured HA replica count.
5. Verify `/healthz`, `/metrics`, `/api/v1/platform/overview`, and one tenant `/api/v1/bootstrap` request.

## RPO/RTO Target

- RPO: 24 hours for Phase 1.
- RTO: 60 minutes for local/private-alpha and customer-controlled deployments.
