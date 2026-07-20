package tenancy

import (
	"sort"
	"sync"
)

type OrganizationRepository interface {
	Save(organization Organization) Organization
	FindAll() []Organization
	FindByID(id string) (Organization, bool)
}

type TenantRepository interface {
	Save(tenant Tenant) Tenant
	FindByID(id string) (Tenant, bool)
	FindByOrganizationID(organizationID string) []Tenant
}

type ProjectRepository interface {
	Save(project Project) Project
	FindByTenantID(organizationID string, tenantID string) []Project
}

type MemoryOrganizationRepository struct {
	mu            sync.RWMutex
	organizations map[string]Organization
}

func NewMemoryOrganizationRepository() *MemoryOrganizationRepository {
	return &MemoryOrganizationRepository{organizations: map[string]Organization{}}
}

func (repo *MemoryOrganizationRepository) Save(organization Organization) Organization {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.organizations[organization.ID] = organization
	return organization
}

func (repo *MemoryOrganizationRepository) FindAll() []Organization {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	organizations := make([]Organization, 0, len(repo.organizations))
	for _, organization := range repo.organizations {
		organizations = append(organizations, organization)
	}
	sort.Slice(organizations, func(i, j int) bool {
		return organizations[i].CreatedAt.Before(organizations[j].CreatedAt)
	})
	return organizations
}

func (repo *MemoryOrganizationRepository) FindByID(id string) (Organization, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	organization, ok := repo.organizations[id]
	return organization, ok
}

type MemoryTenantRepository struct {
	mu      sync.RWMutex
	tenants map[string]Tenant
}

func NewMemoryTenantRepository() *MemoryTenantRepository {
	return &MemoryTenantRepository{tenants: map[string]Tenant{}}
}

func (repo *MemoryTenantRepository) Save(tenant Tenant) Tenant {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.tenants[tenant.ID] = tenant
	return tenant
}

func (repo *MemoryTenantRepository) FindByID(id string) (Tenant, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	tenant, ok := repo.tenants[id]
	return tenant, ok
}

func (repo *MemoryTenantRepository) FindByOrganizationID(organizationID string) []Tenant {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	tenants := make([]Tenant, 0)
	for _, tenant := range repo.tenants {
		if tenant.OrganizationID == organizationID {
			tenants = append(tenants, tenant)
		}
	}
	sort.Slice(tenants, func(i, j int) bool {
		return tenants[i].CreatedAt.Before(tenants[j].CreatedAt)
	})
	return tenants
}

type MemoryProjectRepository struct {
	mu       sync.RWMutex
	projects map[string]Project
}

func NewMemoryProjectRepository() *MemoryProjectRepository {
	return &MemoryProjectRepository{projects: map[string]Project{}}
}

func (repo *MemoryProjectRepository) Save(project Project) Project {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.projects[project.ID] = project
	return project
}

func (repo *MemoryProjectRepository) FindByTenantID(organizationID string, tenantID string) []Project {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	projects := make([]Project, 0)
	for _, project := range repo.projects {
		if project.OrganizationID == organizationID && project.TenantID == tenantID {
			projects = append(projects, project)
		}
	}
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].CreatedAt.Before(projects[j].CreatedAt)
	})
	return projects
}
