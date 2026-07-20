package tenancy

import (
	"testing"
	"time"

	"github.com/dcraft-fusion/control-plane-kernel/internal/audit"
	fusioncontext "github.com/dcraft-fusion/control-plane-kernel/internal/context"
)

func TestCreateOrganizationStoresRecordAndAuditEvent(t *testing.T) {
	repo := NewMemoryOrganizationRepository()
	auditLog := audit.NewMemoryLog()
	ids := []string{"org-1", "audit-1"}
	service := newTestService(repo, NewMemoryTenantRepository(), NewMemoryProjectRepository(), auditLog, ids)

	created, err := service.CreateOrganization(fusioncontext.RequestContext{
		ActorID:       "user-1",
		CorrelationID: "corr-1",
	}, CreateOrganizationRequest{
		Name:                   "Commerce Platform Co.",
		Type:                   OrganizationConsumerApp,
		DeploymentProfile:      DeploymentCustomerVPC,
		Region:                 "us-east-1",
		DataPlaneLocation:      DataPlaneCustomerManaged,
		RawDataMovementAllowed: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID != "org-1" || created.DeploymentProfile != DeploymentCustomerVPC {
		t.Fatalf("unexpected organization: %+v", created)
	}
	events := auditLog.Events()
	if len(events) != 1 {
		t.Fatalf("expected one audit event, got %d", len(events))
	}
	if events[0].ActorID != "user-1" || events[0].CorrelationID != "corr-1" {
		t.Fatalf("audit event missed request context: %+v", events[0])
	}
	if events[0].Metadata["deploymentProfile"] != string(DeploymentCustomerVPC) {
		t.Fatalf("audit event missed deployment metadata: %+v", events[0].Metadata)
	}
}

func TestCreateOrganizationValidatesRequiredFields(t *testing.T) {
	service := NewService(NewMemoryOrganizationRepository(), NewMemoryTenantRepository(), NewMemoryProjectRepository(), audit.NewMemoryLog())
	_, err := service.CreateOrganization(fusioncontext.RequestContext{ActorID: "user-1", CorrelationID: "corr-1"}, CreateOrganizationRequest{})
	if err == nil || err.Error() != "organization name is required" {
		t.Fatalf("expected name validation error, got %v", err)
	}
}

func TestCreateTenantRequiresOrganizationScopeAndExistingOrganization(t *testing.T) {
	organizations := NewMemoryOrganizationRepository()
	service := NewService(organizations, NewMemoryTenantRepository(), NewMemoryProjectRepository(), audit.NewMemoryLog())

	_, err := service.CreateTenant(fusioncontext.RequestContext{ActorID: "user-1", CorrelationID: "corr-1"}, CreateTenantRequest{})
	if err == nil || err.Error() != "missing request context header: X-Fusion-Organization-Id" {
		t.Fatalf("expected organization scope error, got %v", err)
	}

	_, err = service.CreateTenant(fusioncontext.RequestContext{
		ActorID:        "user-1",
		CorrelationID:  "corr-1",
		OrganizationID: "org-missing",
	}, CreateTenantRequest{Name: "Acme tenant", Model: TenantModelB2B, Isolation: TenantIsolationSharedInfra, Region: "us-east-1"})
	if err == nil || err.Error() != "organization not found" {
		t.Fatalf("expected organization existence error, got %v", err)
	}
}

func TestCreateTenantStoresB2B2CIsolationAndAudit(t *testing.T) {
	organizations := NewMemoryOrganizationRepository()
	tenants := NewMemoryTenantRepository()
	auditLog := audit.NewMemoryLog()
	organizations.Save(Organization{ID: "org-1", CreatedAt: time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)})
	service := newTestService(organizations, tenants, NewMemoryProjectRepository(), auditLog, []string{"tenant-1", "audit-1"})

	created, err := service.CreateTenant(fusioncontext.RequestContext{
		ActorID:        "user-1",
		CorrelationID:  "corr-1",
		OrganizationID: "org-1",
	}, CreateTenantRequest{
		Name:      "Marketplace operators",
		Model:     TenantModelB2B2C,
		Isolation: TenantIsolationDedicatedInfra,
		Region:    "ap-south-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID != "tenant-1" || created.OrganizationID != "org-1" || created.Model != TenantModelB2B2C {
		t.Fatalf("unexpected tenant: %+v", created)
	}
	if auditLog.Events()[0].Action != audit.TenantCreated || auditLog.Events()[0].Metadata["isolation"] != string(TenantIsolationDedicatedInfra) {
		t.Fatalf("tenant audit was incomplete: %+v", auditLog.Events()[0])
	}
}

func TestCreateProjectRejectsCrossOrganizationTenant(t *testing.T) {
	tenants := NewMemoryTenantRepository()
	tenants.Save(Tenant{ID: "tenant-1", OrganizationID: "org-2"})
	service := NewService(NewMemoryOrganizationRepository(), tenants, NewMemoryProjectRepository(), audit.NewMemoryLog())

	_, err := service.CreateProject(fusioncontext.RequestContext{
		ActorID:        "user-1",
		CorrelationID:  "corr-1",
		OrganizationID: "org-1",
		TenantID:       "tenant-1",
	}, CreateProjectRequest{Name: "Analytics", Environment: EnvironmentProduction})
	if err == nil || err.Error() != "tenant does not belong to organization" {
		t.Fatalf("expected cross-organization rejection, got %v", err)
	}
}

func TestCreateProjectStoresEnvironmentAndAudit(t *testing.T) {
	projects := NewMemoryProjectRepository()
	auditLog := audit.NewMemoryLog()
	tenants := NewMemoryTenantRepository()
	tenants.Save(Tenant{ID: "tenant-1", OrganizationID: "org-1"})
	service := newTestService(NewMemoryOrganizationRepository(), tenants, projects, auditLog, []string{"project-1", "audit-1"})

	created, err := service.CreateProject(fusioncontext.RequestContext{
		ActorID:        "user-1",
		CorrelationID:  "corr-1",
		OrganizationID: "org-1",
		TenantID:       "tenant-1",
	}, CreateProjectRequest{Name: "Production pipelines", Environment: EnvironmentProduction})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID != "project-1" || created.Environment != EnvironmentProduction {
		t.Fatalf("unexpected project: %+v", created)
	}
	if auditLog.Events()[0].Action != audit.ProjectCreated || auditLog.Events()[0].Metadata["tenantId"] != "tenant-1" {
		t.Fatalf("project audit was incomplete: %+v", auditLog.Events()[0])
	}
}

func newTestService(organizations OrganizationRepository, tenants TenantRepository, projects ProjectRepository, auditLog audit.Log, ids []string) *Service {
	return NewServiceWithClock(organizations, tenants, projects, auditLog, func() time.Time {
		return time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)
	}, func() string {
		id := ids[0]
		ids = ids[1:]
		return id
	})
}
