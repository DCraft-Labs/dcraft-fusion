package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dcraft-fusion/control-plane-kernel/internal/audit"
	fusioncontext "github.com/dcraft-fusion/control-plane-kernel/internal/context"
	"github.com/dcraft-fusion/control-plane-kernel/internal/phase1"
	"github.com/dcraft-fusion/control-plane-kernel/internal/tenancy"
)

func TestCreateOrganizationAPIRequiresRequestContext(t *testing.T) {
	server := newTestServer(audit.NewMemoryLog())
	request := httptest.NewRequest(http.MethodPost, "/api/v1/organizations", bytes.NewBufferString(`{}`))
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}
	assertBodyContains(t, response.Body.Bytes(), "missing request context header: X-Fusion-Actor-Id")
}

func TestCreateOrganizationAPIEmitsAudit(t *testing.T) {
	auditLog := audit.NewMemoryLog()
	server := newTestServer(auditLog)
	body := []byte(`{
		"name": "Commerce Platform Co.",
		"type": "consumer_app",
		"deploymentProfile": "customer_vpc",
		"region": "us-east-1",
		"dataPlaneLocation": "customer_managed",
		"rawDataMovementAllowed": false
	}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/organizations", bytes.NewReader(body))
	request.Header.Set(fusioncontext.ActorIDHeader, "user-1")
	request.Header.Set(fusioncontext.CorrelationIDHeader, "corr-1")
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d with body %s", response.Code, response.Body.String())
	}
	if len(auditLog.Events()) != 1 {
		t.Fatalf("expected one audit event, got %d", len(auditLog.Events()))
	}
	if auditLog.Events()[0].ActorID != "user-1" {
		t.Fatalf("audit actor mismatch: %+v", auditLog.Events()[0])
	}
}

func TestCreateTenantAPIRequiresOrganizationScope(t *testing.T) {
	server := newTestServer(audit.NewMemoryLog())
	request := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(`{}`))
	request.Header.Set(fusioncontext.ActorIDHeader, "user-1")
	request.Header.Set(fusioncontext.CorrelationIDHeader, "corr-1")
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}
	assertBodyContains(t, response.Body.Bytes(), "missing request context header: X-Fusion-Organization-Id")
}

func TestTenantAndProjectAPIsKeepHierarchyInRequestContext(t *testing.T) {
	auditLog := audit.NewMemoryLog()
	server := newTestServer(auditLog)
	routes := server.Routes()

	orgRequest := httptest.NewRequest(http.MethodPost, "/api/v1/organizations", bytes.NewReader([]byte(`{
		"name": "Acme",
		"type": "business",
		"deploymentProfile": "pooled_saas",
		"region": "us-east-1",
		"dataPlaneLocation": "fusion_managed",
		"rawDataMovementAllowed": false
	}`)))
	orgRequest.Header.Set(fusioncontext.ActorIDHeader, "user-1")
	orgRequest.Header.Set(fusioncontext.CorrelationIDHeader, "corr-1")
	orgResponse := httptest.NewRecorder()
	routes.ServeHTTP(orgResponse, orgRequest)
	if orgResponse.Code != http.StatusCreated {
		t.Fatalf("expected organization 201, got %d with body %s", orgResponse.Code, orgResponse.Body.String())
	}
	var org tenancy.Organization
	if err := json.Unmarshal(orgResponse.Body.Bytes(), &org); err != nil {
		t.Fatalf("organization response was not json: %v", err)
	}

	tenantRequest := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewReader([]byte(`{
		"name": "Acme retail customers",
		"model": "b2b2c",
		"isolation": "shared_infra",
		"region": "us-east-1"
	}`)))
	tenantRequest.Header.Set(fusioncontext.ActorIDHeader, "user-1")
	tenantRequest.Header.Set(fusioncontext.CorrelationIDHeader, "corr-2")
	tenantRequest.Header.Set(fusioncontext.OrganizationIDHeader, org.ID)
	tenantResponse := httptest.NewRecorder()
	routes.ServeHTTP(tenantResponse, tenantRequest)
	if tenantResponse.Code != http.StatusCreated {
		t.Fatalf("expected tenant 201, got %d with body %s", tenantResponse.Code, tenantResponse.Body.String())
	}
	var tenant tenancy.Tenant
	if err := json.Unmarshal(tenantResponse.Body.Bytes(), &tenant); err != nil {
		t.Fatalf("tenant response was not json: %v", err)
	}
	if tenant.OrganizationID != org.ID || tenant.Model != tenancy.TenantModelB2B2C {
		t.Fatalf("unexpected tenant: %+v", tenant)
	}

	projectRequest := httptest.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewReader([]byte(`{
		"name": "Checkout intelligence",
		"environment": "production"
	}`)))
	projectRequest.Header.Set(fusioncontext.ActorIDHeader, "user-1")
	projectRequest.Header.Set(fusioncontext.CorrelationIDHeader, "corr-3")
	projectRequest.Header.Set(fusioncontext.OrganizationIDHeader, org.ID)
	projectRequest.Header.Set(fusioncontext.TenantIDHeader, tenant.ID)
	projectResponse := httptest.NewRecorder()
	routes.ServeHTTP(projectResponse, projectRequest)
	if projectResponse.Code != http.StatusCreated {
		t.Fatalf("expected project 201, got %d with body %s", projectResponse.Code, projectResponse.Body.String())
	}
	var project tenancy.Project
	if err := json.Unmarshal(projectResponse.Body.Bytes(), &project); err != nil {
		t.Fatalf("project response was not json: %v", err)
	}
	if project.OrganizationID != org.ID || project.TenantID != tenant.ID || project.Environment != tenancy.EnvironmentProduction {
		t.Fatalf("unexpected project: %+v", project)
	}
	if len(auditLog.Events()) != 3 {
		t.Fatalf("expected three audit events, got %d", len(auditLog.Events()))
	}
}

func TestPhase1BootstrapRequiresTenantScopeAndReturnsSeededPrivateAlphaData(t *testing.T) {
	server := newTestServer(audit.NewMemoryLog())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/bootstrap", nil)
	request.Header.Set(fusioncontext.ActorIDHeader, "user-founder")
	request.Header.Set(fusioncontext.CorrelationIDHeader, "corr-1")
	request.Header.Set(fusioncontext.OrganizationIDHeader, "org-b2b2c-1")
	request.Header.Set(fusioncontext.TenantIDHeader, "tenant-brand-a")
	request.Header.Set(fusioncontext.ProjectIDHeader, "project-prod")
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", response.Code, response.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("bootstrap response was not json: %v", err)
	}
	if len(body["connections"].([]any)) == 0 || len(body["datasets"].([]any)) == 0 || len(body["runs"].([]any)) == 0 {
		t.Fatalf("bootstrap missed Phase 1 data: %#v", body)
	}
}

func TestCreateConnectionAPIRedactsSecretReference(t *testing.T) {
	auditLog := audit.NewMemoryLog()
	server := newTestServer(auditLog)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/connections", bytes.NewReader([]byte(`{
		"name": "Read-only warehouse",
		"kind": "postgres",
		"secretRef": "vault/brand-a/postgres",
		"readOnly": true
	}`)))
	request.Header.Set(fusioncontext.ActorIDHeader, "user-founder")
	request.Header.Set(fusioncontext.CorrelationIDHeader, "corr-1")
	request.Header.Set(fusioncontext.OrganizationIDHeader, "org-b2b2c-1")
	request.Header.Set(fusioncontext.TenantIDHeader, "tenant-brand-a")
	request.Header.Set(fusioncontext.ProjectIDHeader, "project-prod")
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d with body %s", response.Code, response.Body.String())
	}
	if !bytes.Contains(response.Body.Bytes(), []byte(`"secretRef":"secretref:vault/brand-a/postgres"`)) {
		t.Fatalf("secret reference was not redacted: %s", response.Body.String())
	}
	if len(auditLog.Events()) == 0 {
		t.Fatal("expected connection creation audit event")
	}
}

func TestPlatformOverviewRequiresSuperAdminAndReturnsHardeningState(t *testing.T) {
	server := newTestServer(audit.NewMemoryLog())
	denied := httptest.NewRequest(http.MethodGet, "/api/v1/platform/overview", nil)
	denied.Header.Set(fusioncontext.ActorIDHeader, "user-founder")
	denied.Header.Set(fusioncontext.CorrelationIDHeader, "corr-1")
	deniedResponse := httptest.NewRecorder()

	server.Routes().ServeHTTP(deniedResponse, denied)

	if deniedResponse.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d with body %s", deniedResponse.Code, deniedResponse.Body.String())
	}

	allowed := httptest.NewRequest(http.MethodGet, "/api/v1/platform/overview", nil)
	allowed.Header.Set(fusioncontext.ActorIDHeader, "user-superadmin")
	allowed.Header.Set(fusioncontext.CorrelationIDHeader, "corr-2")
	allowedResponse := httptest.NewRecorder()

	server.Routes().ServeHTTP(allowedResponse, allowed)

	if allowedResponse.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d with body %s", allowedResponse.Code, allowedResponse.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(allowedResponse.Body.Bytes(), &body); err != nil {
		t.Fatalf("overview response was not json: %v", err)
	}
	if len(body["identityProviders"].([]any)) == 0 || len(body["secretProviders"].([]any)) == 0 || len(body["adapters"].([]any)) == 0 {
		t.Fatalf("overview missed enterprise hardening state: %#v", body)
	}
	if len(body["organizations"].([]any)) == 0 {
		t.Fatalf("overview missed organization catalog: %#v", body)
	}
}

func TestPlatformCreateTenantAndProjectWithoutHeaderInjectedScope(t *testing.T) {
	server := newTestServer(audit.NewMemoryLog())
	routes := server.Routes()

	createTenantRequest := httptest.NewRequest(http.MethodPost, "/api/v1/platform/tenants", bytes.NewReader([]byte(`{
		"organizationId": "org-b2b-1",
		"name": "Analytics Workspace",
		"model": "b2b",
		"isolation": "dedicated_infra",
		"region": "us-east-1"
	}`)))
	createTenantRequest.Header.Set(fusioncontext.ActorIDHeader, "user-superadmin")
	createTenantRequest.Header.Set(fusioncontext.CorrelationIDHeader, "corr-super-1")
	createTenantResponse := httptest.NewRecorder()

	routes.ServeHTTP(createTenantResponse, createTenantRequest)

	if createTenantResponse.Code != http.StatusCreated {
		t.Fatalf("expected superadmin tenant creation 201, got %d with body %s", createTenantResponse.Code, createTenantResponse.Body.String())
	}

	var tenant tenancy.Tenant
	if err := json.Unmarshal(createTenantResponse.Body.Bytes(), &tenant); err != nil {
		t.Fatalf("tenant response was not json: %v", err)
	}

	createProjectRequest := httptest.NewRequest(http.MethodPost, "/api/v1/platform/projects", bytes.NewReader([]byte(`{
		"organizationId": "org-b2b-1",
		"tenantId": "`+tenant.ID+`",
		"name": "Production Workspace",
		"environment": "production"
	}`)))
	createProjectRequest.Header.Set(fusioncontext.ActorIDHeader, "user-superadmin")
	createProjectRequest.Header.Set(fusioncontext.CorrelationIDHeader, "corr-super-2")
	createProjectResponse := httptest.NewRecorder()

	routes.ServeHTTP(createProjectResponse, createProjectRequest)

	if createProjectResponse.Code != http.StatusCreated {
		t.Fatalf("expected superadmin project creation 201, got %d with body %s", createProjectResponse.Code, createProjectResponse.Body.String())
	}
}

func assertBodyContains(t *testing.T, body []byte, expected string) {
	t.Helper()
	var parsed map[string]string
	if err := json.Unmarshal(body, &parsed); err != nil {
		t.Fatalf("response was not json: %v", err)
	}
	if parsed["error"] != expected {
		t.Fatalf("expected %q, got %#v", expected, parsed)
	}
}

func newTestServer(auditLog audit.Log) *Server {
	orgRepo := tenancy.NewMemoryOrganizationRepository()
	tenantRepo := tenancy.NewMemoryTenantRepository()
	projectRepo := tenancy.NewMemoryProjectRepository()
	phase1.SeedTenancyCatalog(orgRepo, tenantRepo, projectRepo)
	return NewServer(orgRepo, tenantRepo, projectRepo, auditLog)
}
