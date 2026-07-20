package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/dcraft-fusion/control-plane-kernel/internal/audit"
	"github.com/dcraft-fusion/control-plane-kernel/internal/tenancy"
)

type OrganizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (repo *OrganizationRepository) Save(organization tenancy.Organization) tenancy.Organization {
	_, err := repo.db.ExecContext(context.Background(), `INSERT INTO organizations (
		id, name, type, deployment_profile, region, data_plane_location, raw_data_movement_allowed, created_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		type = EXCLUDED.type,
		deployment_profile = EXCLUDED.deployment_profile,
		region = EXCLUDED.region,
		data_plane_location = EXCLUDED.data_plane_location,
		raw_data_movement_allowed = EXCLUDED.raw_data_movement_allowed`,
		organization.ID,
		organization.Name,
		string(organization.Type),
		string(organization.DeploymentProfile),
		organization.Region,
		string(organization.DataPlaneLocation),
		organization.RawDataMovementAllowed,
		organization.CreatedAt,
	)
	if err != nil {
		panic(err)
	}
	return organization
}

func (repo *OrganizationRepository) FindAll() []tenancy.Organization {
	rows, err := repo.db.QueryContext(context.Background(), `SELECT id, name, type, deployment_profile, region, data_plane_location, raw_data_movement_allowed, created_at FROM organizations ORDER BY created_at ASC`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	organizations := []tenancy.Organization{}
	for rows.Next() {
		var organization tenancy.Organization
		if err := rows.Scan(&organization.ID, &organization.Name, &organization.Type, &organization.DeploymentProfile, &organization.Region, &organization.DataPlaneLocation, &organization.RawDataMovementAllowed, &organization.CreatedAt); err != nil {
			panic(err)
		}
		organizations = append(organizations, organization)
	}
	return organizations
}

func (repo *OrganizationRepository) FindByID(id string) (tenancy.Organization, bool) {
	var organization tenancy.Organization
	err := repo.db.QueryRowContext(context.Background(), `SELECT id, name, type, deployment_profile, region, data_plane_location, raw_data_movement_allowed, created_at FROM organizations WHERE id = $1`, id).
		Scan(&organization.ID, &organization.Name, &organization.Type, &organization.DeploymentProfile, &organization.Region, &organization.DataPlaneLocation, &organization.RawDataMovementAllowed, &organization.CreatedAt)
	if err == sql.ErrNoRows {
		return tenancy.Organization{}, false
	}
	if err != nil {
		panic(err)
	}
	return organization, true
}

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (repo *TenantRepository) Save(tenant tenancy.Tenant) tenancy.Tenant {
	_, err := repo.db.ExecContext(context.Background(), `INSERT INTO tenants (
		id, organization_id, name, model, isolation, region, created_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7)
	ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		model = EXCLUDED.model,
		isolation = EXCLUDED.isolation,
		region = EXCLUDED.region`,
		tenant.ID,
		tenant.OrganizationID,
		tenant.Name,
		string(tenant.Model),
		string(tenant.Isolation),
		tenant.Region,
		tenant.CreatedAt,
	)
	if err != nil {
		panic(err)
	}
	return tenant
}

func (repo *TenantRepository) FindByID(id string) (tenancy.Tenant, bool) {
	var tenant tenancy.Tenant
	err := repo.db.QueryRowContext(context.Background(), `SELECT id, organization_id, name, model, isolation, region, created_at FROM tenants WHERE id = $1`, id).
		Scan(&tenant.ID, &tenant.OrganizationID, &tenant.Name, &tenant.Model, &tenant.Isolation, &tenant.Region, &tenant.CreatedAt)
	if err == sql.ErrNoRows {
		return tenancy.Tenant{}, false
	}
	if err != nil {
		panic(err)
	}
	return tenant, true
}

func (repo *TenantRepository) FindByOrganizationID(organizationID string) []tenancy.Tenant {
	rows, err := repo.db.QueryContext(context.Background(), `SELECT id, organization_id, name, model, isolation, region, created_at FROM tenants WHERE organization_id = $1 ORDER BY created_at ASC`, organizationID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	tenants := []tenancy.Tenant{}
	for rows.Next() {
		var tenant tenancy.Tenant
		if err := rows.Scan(&tenant.ID, &tenant.OrganizationID, &tenant.Name, &tenant.Model, &tenant.Isolation, &tenant.Region, &tenant.CreatedAt); err != nil {
			panic(err)
		}
		tenants = append(tenants, tenant)
	}
	return tenants
}

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (repo *ProjectRepository) Save(project tenancy.Project) tenancy.Project {
	_, err := repo.db.ExecContext(context.Background(), `INSERT INTO projects (
		id, organization_id, tenant_id, name, environment, created_at
	) VALUES ($1,$2,$3,$4,$5,$6)
	ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		environment = EXCLUDED.environment`,
		project.ID,
		project.OrganizationID,
		project.TenantID,
		project.Name,
		string(project.Environment),
		project.CreatedAt,
	)
	if err != nil {
		panic(err)
	}
	return project
}

func (repo *ProjectRepository) FindByTenantID(organizationID string, tenantID string) []tenancy.Project {
	rows, err := repo.db.QueryContext(context.Background(), `SELECT id, organization_id, tenant_id, name, environment, created_at FROM projects WHERE organization_id = $1 AND tenant_id = $2 ORDER BY created_at ASC`, organizationID, tenantID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	projects := []tenancy.Project{}
	for rows.Next() {
		var project tenancy.Project
		if err := rows.Scan(&project.ID, &project.OrganizationID, &project.TenantID, &project.Name, &project.Environment, &project.CreatedAt); err != nil {
			panic(err)
		}
		projects = append(projects, project)
	}
	return projects
}

type AuditLog struct {
	db *sql.DB
}

func NewAuditLog(db *sql.DB) *AuditLog {
	return &AuditLog{db: db}
}

func (log *AuditLog) Append(event audit.Event) {
	metadata, err := json.Marshal(event.Metadata)
	if err != nil {
		panic(err)
	}
	_, err = log.db.ExecContext(context.Background(), `INSERT INTO audit_events (
		id, action, actor_id, correlation_id, resource_type, resource_id, occurred_at, metadata
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	ON CONFLICT (id) DO NOTHING`,
		event.ID,
		string(event.Action),
		event.ActorID,
		event.CorrelationID,
		event.ResourceType,
		event.ResourceID,
		event.OccurredAt,
		metadata,
	)
	if err != nil {
		panic(err)
	}
}

func (log *AuditLog) Events() []audit.Event {
	rows, err := log.db.QueryContext(context.Background(), `SELECT id, action, actor_id, correlation_id, resource_type, resource_id, occurred_at, metadata FROM audit_events ORDER BY occurred_at ASC`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	events := []audit.Event{}
	for rows.Next() {
		var event audit.Event
		var metadata []byte
		if err := rows.Scan(&event.ID, &event.Action, &event.ActorID, &event.CorrelationID, &event.ResourceType, &event.ResourceID, &event.OccurredAt, &metadata); err != nil {
			panic(err)
		}
		if err := json.Unmarshal(metadata, &event.Metadata); err != nil {
			panic(err)
		}
		events = append(events, event)
	}
	return events
}
