package tenancy

import (
	"errors"
	"strings"
	"time"

	"github.com/dcraft-fusion/control-plane-kernel/internal/audit"
	fusioncontext "github.com/dcraft-fusion/control-plane-kernel/internal/context"
	"github.com/dcraft-fusion/control-plane-kernel/internal/ids"
)

type Service struct {
	organizations OrganizationRepository
	tenants       TenantRepository
	projects      ProjectRepository
	audit         audit.Log
	now           func() time.Time
	newID         func() string
}

func NewService(organizations OrganizationRepository, tenants TenantRepository, projects ProjectRepository, auditLog audit.Log) *Service {
	return &Service{
		organizations: organizations,
		tenants:       tenants,
		projects:      projects,
		audit:         auditLog,
		now:           time.Now,
		newID:         ids.New,
	}
}

func NewServiceWithClock(organizations OrganizationRepository, tenants TenantRepository, projects ProjectRepository, auditLog audit.Log, now func() time.Time, newID func() string) *Service {
	return &Service{
		organizations: organizations,
		tenants:       tenants,
		projects:      projects,
		audit:         auditLog,
		now:           now,
		newID:         newID,
	}
}

func (service *Service) CreateOrganization(ctx fusioncontext.RequestContext, request CreateOrganizationRequest) (Organization, error) {
	if strings.TrimSpace(request.Name) == "" {
		return Organization{}, errors.New("organization name is required")
	}
	if strings.TrimSpace(request.Region) == "" {
		return Organization{}, errors.New("organization region is required")
	}
	if request.Type == "" {
		return Organization{}, errors.New("organization type is required")
	}
	if request.DeploymentProfile == "" {
		return Organization{}, errors.New("deployment profile is required")
	}
	if request.DataPlaneLocation == "" {
		return Organization{}, errors.New("data plane location is required")
	}

	organization := Organization{
		ID:                     service.newID(),
		Name:                   strings.TrimSpace(request.Name),
		Type:                   request.Type,
		DeploymentProfile:      request.DeploymentProfile,
		Region:                 strings.TrimSpace(request.Region),
		DataPlaneLocation:      request.DataPlaneLocation,
		RawDataMovementAllowed: request.RawDataMovementAllowed,
		CreatedAt:              service.now().UTC(),
	}

	saved := service.organizations.Save(organization)
	service.audit.Append(audit.Event{
		ID:            service.newID(),
		Action:        audit.OrganizationCreated,
		ActorID:       ctx.ActorID,
		CorrelationID: ctx.CorrelationID,
		ResourceType:  "Organization",
		ResourceID:    saved.ID,
		OccurredAt:    service.now().UTC(),
		Metadata: audit.RedactMetadata(map[string]string{
			"deploymentProfile": string(saved.DeploymentProfile),
			"organizationType":  string(saved.Type),
			"region":            saved.Region,
		}),
	})

	return saved, nil
}

func (service *Service) CreateTenant(ctx fusioncontext.RequestContext, request CreateTenantRequest) (Tenant, error) {
	if err := ctx.RequireOrganization(); err != nil {
		return Tenant{}, err
	}
	if _, ok := service.organizations.FindByID(ctx.OrganizationID); !ok {
		return Tenant{}, errors.New("organization not found")
	}
	if strings.TrimSpace(request.Name) == "" {
		return Tenant{}, errors.New("tenant name is required")
	}
	if request.Model == "" {
		return Tenant{}, errors.New("tenant model is required")
	}
	if request.Isolation == "" {
		return Tenant{}, errors.New("tenant isolation is required")
	}
	if strings.TrimSpace(request.Region) == "" {
		return Tenant{}, errors.New("tenant region is required")
	}

	tenant := Tenant{
		ID:             service.newID(),
		OrganizationID: ctx.OrganizationID,
		Name:           strings.TrimSpace(request.Name),
		Model:          request.Model,
		Isolation:      request.Isolation,
		Region:         strings.TrimSpace(request.Region),
		CreatedAt:      service.now().UTC(),
	}

	saved := service.tenants.Save(tenant)
	service.audit.Append(audit.Event{
		ID:            service.newID(),
		Action:        audit.TenantCreated,
		ActorID:       ctx.ActorID,
		CorrelationID: ctx.CorrelationID,
		ResourceType:  "Tenant",
		ResourceID:    saved.ID,
		OccurredAt:    service.now().UTC(),
		Metadata: audit.RedactMetadata(map[string]string{
			"organizationId": saved.OrganizationID,
			"tenantModel":    string(saved.Model),
			"isolation":      string(saved.Isolation),
			"region":         saved.Region,
		}),
	})

	return saved, nil
}

func (service *Service) CreateProject(ctx fusioncontext.RequestContext, request CreateProjectRequest) (Project, error) {
	if err := ctx.RequireTenant(); err != nil {
		return Project{}, err
	}
	tenant, ok := service.tenants.FindByID(ctx.TenantID)
	if !ok {
		return Project{}, errors.New("tenant not found")
	}
	if tenant.OrganizationID != ctx.OrganizationID {
		return Project{}, errors.New("tenant does not belong to organization")
	}
	if strings.TrimSpace(request.Name) == "" {
		return Project{}, errors.New("project name is required")
	}
	if request.Environment == "" {
		return Project{}, errors.New("project environment is required")
	}

	project := Project{
		ID:             service.newID(),
		OrganizationID: ctx.OrganizationID,
		TenantID:       ctx.TenantID,
		Name:           strings.TrimSpace(request.Name),
		Environment:    request.Environment,
		CreatedAt:      service.now().UTC(),
	}

	saved := service.projects.Save(project)
	service.audit.Append(audit.Event{
		ID:            service.newID(),
		Action:        audit.ProjectCreated,
		ActorID:       ctx.ActorID,
		CorrelationID: ctx.CorrelationID,
		ResourceType:  "Project",
		ResourceID:    saved.ID,
		OccurredAt:    service.now().UTC(),
		Metadata: audit.RedactMetadata(map[string]string{
			"organizationId": saved.OrganizationID,
			"tenantId":       saved.TenantID,
			"environment":    string(saved.Environment),
		}),
	})

	return saved, nil
}

func (service *Service) ListTenants(ctx fusioncontext.RequestContext) ([]Tenant, error) {
	if err := ctx.RequireOrganization(); err != nil {
		return nil, err
	}
	return service.tenants.FindByOrganizationID(ctx.OrganizationID), nil
}

func (service *Service) ListProjects(ctx fusioncontext.RequestContext) ([]Project, error) {
	if err := ctx.RequireTenant(); err != nil {
		return nil, err
	}
	return service.projects.FindByTenantID(ctx.OrganizationID, ctx.TenantID), nil
}
