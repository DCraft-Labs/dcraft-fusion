package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/dcraft-fusion/control-plane-kernel/internal/audit"
	fusioncontext "github.com/dcraft-fusion/control-plane-kernel/internal/context"
	"github.com/dcraft-fusion/control-plane-kernel/internal/phase1"
	"github.com/dcraft-fusion/control-plane-kernel/internal/tenancy"
)

type Server struct {
	organizations tenancy.OrganizationRepository
	tenants       tenancy.TenantRepository
	projects      tenancy.ProjectRepository
	auditLog      audit.Log
	tenancy       *tenancy.Service
	phase1        *phase1.Store
}

func NewServer(organizations tenancy.OrganizationRepository, tenants tenancy.TenantRepository, projects tenancy.ProjectRepository, auditLog audit.Log) *Server {
	return NewServerWithPhase1(organizations, tenants, projects, auditLog, phase1.NewSeedStore(auditLog))
}

func NewServerWithPhase1(organizations tenancy.OrganizationRepository, tenants tenancy.TenantRepository, projects tenancy.ProjectRepository, auditLog audit.Log, phase1Store *phase1.Store) *Server {
	return &Server{
		organizations: organizations,
		tenants:       tenants,
		projects:      projects,
		auditLog:      auditLog,
		tenancy:       tenancy.NewService(organizations, tenants, projects, auditLog),
		phase1:        phase1Store,
	}
}

func (server *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", server.health)
	mux.HandleFunc("POST /api/v1/organizations", server.createOrganization)
	mux.HandleFunc("GET /api/v1/organizations", server.listOrganizations)
	mux.HandleFunc("POST /api/v1/platform/organizations", server.platformCreateOrganization)
	mux.HandleFunc("POST /api/v1/platform/tenants", server.platformCreateTenant)
	mux.HandleFunc("POST /api/v1/platform/projects", server.platformCreateProject)
	mux.HandleFunc("POST /api/v1/tenants", server.createTenant)
	mux.HandleFunc("GET /api/v1/tenants", server.listTenants)
	mux.HandleFunc("POST /api/v1/projects", server.createProject)
	mux.HandleFunc("GET /api/v1/projects", server.listProjects)
	mux.HandleFunc("GET /api/v1/bootstrap", server.bootstrap)
	mux.HandleFunc("POST /api/v1/invitations", server.inviteUser)
	mux.HandleFunc("POST /api/v1/service-accounts", server.createServiceAccount)
	mux.HandleFunc("POST /api/v1/connections", server.createConnection)
	mux.HandleFunc("GET /api/v1/connections", server.listConnections)
	mux.HandleFunc("POST /api/v1/connections/{connectionID}/test", server.testConnection)
	mux.HandleFunc("POST /api/v1/connections/{connectionID}/discover", server.discoverMetadata)
	mux.HandleFunc("GET /api/v1/datasets", server.listDatasets)
	mux.HandleFunc("GET /api/v1/runs", server.listRuns)
	mux.HandleFunc("GET /api/v1/policies", server.listPolicies)
	mux.HandleFunc("GET /api/v1/ai/recommendations", server.listAIRecommendations)
	mux.HandleFunc("GET /api/v1/audit-events", server.listAuditEvents)
	mux.HandleFunc("GET /api/v1/platform/overview", server.platformOverview)
	return mux
}

func (server *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "control-plane-kernel"})
}

func (server *Server) createOrganization(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var request tenancy.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	organization, err := server.tenancy.CreateOrganization(ctx, request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, organization)
}

func (server *Server) listOrganizations(w http.ResponseWriter, r *http.Request) {
	// Require a valid request context so this endpoint is not effectively
	// public after the auth middleware. If the caller is scoped to a
	// specific organization (X-Fusion-Organization-Id header present),
	// filter the result to that org only — a scoped user must not see
	// other organizations. A platform superadmin (no org header) sees all.
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	organizations := server.organizations.FindAll()
	if ctx.OrganizationID != "" {
		filtered := make([]tenancy.Organization, 0, len(organizations))
		for _, organization := range organizations {
			if organization.ID == ctx.OrganizationID {
				filtered = append(filtered, organization)
			}
		}
		organizations = filtered
	}
	writeJSON(w, http.StatusOK, organizations)
}

func (server *Server) createTenant(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var request tenancy.CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	tenant, err := server.tenancy.CreateTenant(ctx, request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, tenant)
}

func (server *Server) platformCreateTenant(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if _, err := server.phase1.PlatformOverview(ctx); err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}

	var request struct {
		OrganizationID string                  `json:"organizationId"`
		Name           string                  `json:"name"`
		Model          tenancy.TenantModel     `json:"model"`
		Isolation      tenancy.TenantIsolation `json:"isolation"`
		Region         string                  `json:"region"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	tenant, err := server.tenancy.CreateTenant(fusioncontext.RequestContext{
		ActorID:        ctx.ActorID,
		CorrelationID:  ctx.CorrelationID,
		OrganizationID: request.OrganizationID,
	}, tenancy.CreateTenantRequest{
		Name:      request.Name,
		Model:     request.Model,
		Isolation: request.Isolation,
		Region:    request.Region,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, tenant)
}

func (server *Server) listTenants(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	tenants, err := server.tenancy.ListTenants(ctx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, tenants)
}

func (server *Server) createProject(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var request tenancy.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	project, err := server.tenancy.CreateProject(ctx, request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, project)
}

func (server *Server) platformCreateProject(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if _, err := server.phase1.PlatformOverview(ctx); err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}

	var request struct {
		OrganizationID string                  `json:"organizationId"`
		TenantID       string                  `json:"tenantId"`
		Name           string                  `json:"name"`
		Environment    tenancy.EnvironmentType `json:"environment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	project, err := server.tenancy.CreateProject(fusioncontext.RequestContext{
		ActorID:        ctx.ActorID,
		CorrelationID:  ctx.CorrelationID,
		OrganizationID: request.OrganizationID,
		TenantID:       request.TenantID,
	}, tenancy.CreateProjectRequest{
		Name:        request.Name,
		Environment: request.Environment,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, project)
}

func (server *Server) listProjects(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	projects, err := server.tenancy.ListProjects(ctx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

func (server *Server) bootstrap(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	bootstrap, err := server.phase1.Bootstrap(ctx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, bootstrap)
}

func (server *Server) inviteUser(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var request phase1.InviteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	invitation, err := server.phase1.InviteUser(ctx, request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, invitation)
}

func (server *Server) createServiceAccount(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var request phase1.CreateServiceAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	account, err := server.phase1.CreateServiceAccount(ctx, request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, account)
}

func (server *Server) createConnection(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var request phase1.CreateConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	connection, err := server.phase1.CreateConnection(ctx, request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, connection)
}

func (server *Server) listConnections(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	connections, err := server.phase1.Connections(ctx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, connections)
}

func (server *Server) testConnection(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.phase1.TestConnection(ctx, r.PathValue("connectionID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) discoverMetadata(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	datasets, err := server.phase1.DiscoverMetadata(ctx, r.PathValue("connectionID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, datasets)
}

func (server *Server) listDatasets(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	datasets, err := server.phase1.Datasets(ctx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, datasets)
}

func (server *Server) listRuns(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	runs, err := server.phase1.Runs(ctx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, runs)
}

func (server *Server) listPolicies(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	policies, err := server.phase1.Policies(ctx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, policies)
}

func (server *Server) listAIRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	recommendations, err := server.phase1.AIRecommendations(ctx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, recommendations)
}

func (server *Server) listAuditEvents(w http.ResponseWriter, r *http.Request) {
	// Audit events are platform-wide and may contain cross-tenant data.
	// Require a valid request context AND platform superadmin access so
	// the endpoint is not effectively public after the auth middleware.
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if _, err := server.phase1.PlatformOverview(ctx); err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}
	writeJSON(w, http.StatusOK, server.auditLog.Events())
}

func (server *Server) platformOverview(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	overview, err := server.phase1.PlatformOverview(ctx)
	if err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}
	overview.Organizations = server.platformOrganizations()
	overview.RecentAuditEvents = server.platformAuditEvents()
	writeJSON(w, http.StatusOK, overview)
}

func (server *Server) platformCreateOrganization(w http.ResponseWriter, r *http.Request) {
	ctx, err := fusioncontext.FromHeaders(r.Header)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if _, err := server.phase1.PlatformOverview(ctx); err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}

	var request tenancy.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	organization, err := server.tenancy.CreateOrganization(fusioncontext.RequestContext{
		ActorID:       ctx.ActorID,
		CorrelationID: ctx.CorrelationID,
	}, request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, organization)
}

func (server *Server) platformOrganizations() []phase1.PlatformOrganization {
	organizations := make([]phase1.PlatformOrganization, 0)
	for _, organization := range server.organizations.FindAll() {
		tenants := make([]phase1.PlatformTenant, 0)
		for _, tenant := range server.tenants.FindByOrganizationID(organization.ID) {
			projects := make([]phase1.PlatformProject, 0)
			for _, project := range server.projects.FindByTenantID(organization.ID, tenant.ID) {
				projects = append(projects, phase1.PlatformProject{
					ID:          project.ID,
					Name:        project.Name,
					Environment: string(project.Environment),
					CreatedAt:   project.CreatedAt,
				})
			}
			tenants = append(tenants, phase1.PlatformTenant{
				ID:        tenant.ID,
				Name:      tenant.Name,
				Model:     string(tenant.Model),
				Isolation: string(tenant.Isolation),
				Region:    tenant.Region,
				CreatedAt: tenant.CreatedAt,
				Projects:  projects,
			})
		}
		organizations = append(organizations, phase1.PlatformOrganization{
			ID:                     organization.ID,
			Name:                   organization.Name,
			Type:                   string(organization.Type),
			DeploymentProfile:      string(organization.DeploymentProfile),
			Region:                 organization.Region,
			DataPlaneLocation:      string(organization.DataPlaneLocation),
			RawDataMovementAllowed: organization.RawDataMovementAllowed,
			CreatedAt:              organization.CreatedAt,
			Tenants:                tenants,
		})
	}
	return organizations
}

func (server *Server) platformAuditEvents() []phase1.PlatformAuditEvent {
	events := server.auditLog.Events()
	if len(events) > 12 {
		events = events[len(events)-12:]
	}
	result := make([]phase1.PlatformAuditEvent, 0, len(events))
	for _, event := range events {
		result = append(result, phase1.PlatformAuditEvent{
			ID:            event.ID,
			Action:        string(event.Action),
			ActorID:       event.ActorID,
			CorrelationID: event.CorrelationID,
			ResourceType:  event.ResourceType,
			ResourceID:    event.ResourceID,
			OccurredAt:    event.OccurredAt,
		})
	}
	return result
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
