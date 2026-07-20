package tenancy

import "time"

type OrganizationType string

const (
	OrganizationBusiness               OrganizationType = "business"
	OrganizationConsumerApp            OrganizationType = "consumer_app"
	OrganizationManagedServiceProvider OrganizationType = "managed_service_provider"
	OrganizationInternalPlatform       OrganizationType = "internal_platform"
)

type DeploymentProfile string

const (
	DeploymentPooledSaaS        DeploymentProfile = "pooled_saas"
	DeploymentDedicatedDatabase DeploymentProfile = "dedicated_database"
	DeploymentDedicatedStack    DeploymentProfile = "dedicated_stack"
	DeploymentCustomerVPC       DeploymentProfile = "customer_vpc"
	DeploymentOnPrem            DeploymentProfile = "on_prem"
)

type DataPlaneLocation string

const (
	DataPlaneFusionManaged   DataPlaneLocation = "fusion_managed"
	DataPlaneCustomerManaged DataPlaneLocation = "customer_managed"
	DataPlaneHybrid          DataPlaneLocation = "hybrid"
)

type Organization struct {
	ID                     string            `json:"id"`
	Name                   string            `json:"name"`
	Type                   OrganizationType  `json:"type"`
	DeploymentProfile      DeploymentProfile `json:"deploymentProfile"`
	Region                 string            `json:"region"`
	DataPlaneLocation      DataPlaneLocation `json:"dataPlaneLocation"`
	RawDataMovementAllowed bool              `json:"rawDataMovementAllowed"`
	CreatedAt              time.Time         `json:"createdAt"`
}

type TenantModel string

const (
	TenantModelB2B      TenantModel = "b2b"
	TenantModelB2C      TenantModel = "b2c"
	TenantModelB2B2C    TenantModel = "b2b2c"
	TenantModelInternal TenantModel = "internal"
)

type TenantIsolation string

const (
	TenantIsolationSharedInfra    TenantIsolation = "shared_infra"
	TenantIsolationDedicatedInfra TenantIsolation = "dedicated_infra"
	TenantIsolationCustomerInfra  TenantIsolation = "customer_infra"
)

type Tenant struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organizationId"`
	Name           string          `json:"name"`
	Model          TenantModel     `json:"model"`
	Isolation      TenantIsolation `json:"isolation"`
	Region         string          `json:"region"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type CreateTenantRequest struct {
	Name      string          `json:"name"`
	Model     TenantModel     `json:"model"`
	Isolation TenantIsolation `json:"isolation"`
	Region    string          `json:"region"`
}

type EnvironmentType string

const (
	EnvironmentDevelopment EnvironmentType = "development"
	EnvironmentStaging     EnvironmentType = "staging"
	EnvironmentProduction  EnvironmentType = "production"
)

type Project struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organizationId"`
	TenantID       string          `json:"tenantId"`
	Name           string          `json:"name"`
	Environment    EnvironmentType `json:"environment"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type CreateProjectRequest struct {
	Name        string          `json:"name"`
	Environment EnvironmentType `json:"environment"`
}

type CreateOrganizationRequest struct {
	Name                   string            `json:"name"`
	Type                   OrganizationType  `json:"type"`
	DeploymentProfile      DeploymentProfile `json:"deploymentProfile"`
	Region                 string            `json:"region"`
	DataPlaneLocation      DataPlaneLocation `json:"dataPlaneLocation"`
	RawDataMovementAllowed bool              `json:"rawDataMovementAllowed"`
}

func IsCustomerOwned(profile DeploymentProfile) bool {
	return profile == DeploymentCustomerVPC || profile == DeploymentOnPrem
}

func IsDedicated(profile DeploymentProfile) bool {
	return profile == DeploymentDedicatedDatabase || profile == DeploymentDedicatedStack || IsCustomerOwned(profile)
}
