package phase1

import (
	"time"

	"github.com/dcraft-fusion/control-plane-kernel/internal/tenancy"
)

func SeedTenancyCatalog(
	organizations tenancy.OrganizationRepository,
	tenants tenancy.TenantRepository,
	projects tenancy.ProjectRepository,
) {
	createdAt := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)

	organizations.Save(tenancy.Organization{
		ID:                     "org-b2b2c-1",
		Name:                   "Commerce Platform Co.",
		Type:                   tenancy.OrganizationConsumerApp,
		DeploymentProfile:      tenancy.DeploymentCustomerVPC,
		Region:                 "us-east-1",
		DataPlaneLocation:      tenancy.DataPlaneCustomerManaged,
		RawDataMovementAllowed: false,
		CreatedAt:              createdAt,
	})
	organizations.Save(tenancy.Organization{
		ID:                     "org-b2b-1",
		Name:                   "Acme Analytics",
		Type:                   tenancy.OrganizationBusiness,
		DeploymentProfile:      tenancy.DeploymentPooledSaaS,
		Region:                 "us-east-1",
		DataPlaneLocation:      tenancy.DataPlaneFusionManaged,
		RawDataMovementAllowed: false,
		CreatedAt:              createdAt.Add(1 * time.Minute),
	})
	organizations.Save(tenancy.Organization{
		ID:                     "org-internal-1",
		Name:                   "Global Enterprise Platform",
		Type:                   tenancy.OrganizationInternalPlatform,
		DeploymentProfile:      tenancy.DeploymentDedicatedStack,
		Region:                 "eu-west-1",
		DataPlaneLocation:      tenancy.DataPlaneHybrid,
		RawDataMovementAllowed: false,
		CreatedAt:              createdAt.Add(2 * time.Minute),
	})

	tenants.Save(tenancy.Tenant{
		ID:             "tenant-brand-a",
		OrganizationID: "org-b2b2c-1",
		Name:           "Brand A",
		Model:          tenancy.TenantModelB2B2C,
		Isolation:      tenancy.TenantIsolationCustomerInfra,
		Region:         "us-east-1",
		CreatedAt:      createdAt,
	})
	tenants.Save(tenancy.Tenant{
		ID:             "tenant-brand-b",
		OrganizationID: "org-b2b2c-1",
		Name:           "Brand B",
		Model:          tenancy.TenantModelB2B2C,
		Isolation:      tenancy.TenantIsolationCustomerInfra,
		Region:         "us-east-1",
		CreatedAt:      createdAt.Add(1 * time.Minute),
	})
	tenants.Save(tenancy.Tenant{
		ID:             "tenant-finance",
		OrganizationID: "org-b2b-1",
		Name:           "Finance",
		Model:          tenancy.TenantModelB2B,
		Isolation:      tenancy.TenantIsolationSharedInfra,
		Region:         "us-east-1",
		CreatedAt:      createdAt.Add(2 * time.Minute),
	})

	projects.Save(tenancy.Project{
		ID:             "project-prod",
		OrganizationID: "org-b2b2c-1",
		TenantID:       "tenant-brand-a",
		Name:           "Production Data Platform",
		Environment:    tenancy.EnvironmentProduction,
		CreatedAt:      createdAt,
	})
	projects.Save(tenancy.Project{
		ID:             "project-sandbox",
		OrganizationID: "org-b2b2c-1",
		TenantID:       "tenant-brand-a",
		Name:           "Sandbox",
		Environment:    tenancy.EnvironmentStaging,
		CreatedAt:      createdAt.Add(1 * time.Minute),
	})
	projects.Save(tenancy.Project{
		ID:             "project-finance-prod",
		OrganizationID: "org-b2b-1",
		TenantID:       "tenant-finance",
		Name:           "Finance Operations",
		Environment:    tenancy.EnvironmentProduction,
		CreatedAt:      createdAt.Add(2 * time.Minute),
	})
}
