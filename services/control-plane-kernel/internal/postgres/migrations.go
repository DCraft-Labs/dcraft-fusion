package postgres

import (
	"context"
	"database/sql"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS organizations (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			deployment_profile TEXT NOT NULL,
			region TEXT NOT NULL,
			data_plane_location TEXT NOT NULL,
			raw_data_movement_allowed BOOLEAN NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tenants (
			id TEXT PRIMARY KEY,
			organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			model TEXT NOT NULL,
			isolation TEXT NOT NULL,
			region TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			organization_id TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			tenant_id TEXT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			environment TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS audit_events (
			id TEXT PRIMARY KEY,
			action TEXT NOT NULL,
			actor_id TEXT NOT NULL,
			correlation_id TEXT NOT NULL,
			resource_type TEXT NOT NULL,
			resource_id TEXT NOT NULL,
			occurred_at TIMESTAMPTZ NOT NULL,
			metadata JSONB NOT NULL DEFAULT '{}'::jsonb
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tenants_org ON tenants(organization_id)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_scope ON projects(organization_id, tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_events_occurred_at ON audit_events(occurred_at DESC)`,
	}
	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}
