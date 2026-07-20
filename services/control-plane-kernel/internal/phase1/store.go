package phase1

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dcraft-fusion/control-plane-kernel/internal/audit"
	fusioncontext "github.com/dcraft-fusion/control-plane-kernel/internal/context"
	"github.com/dcraft-fusion/control-plane-kernel/internal/ids"
)

type Store struct {
	mu                sync.RWMutex
	memberships       map[string]Membership
	invitations       map[string]Invitation
	serviceAccounts   map[string]ServiceAccount
	connections       map[string]Connection
	datasets          map[string]Dataset
	runs              map[string]WorkflowRun
	policies          map[string]Policy
	aiRecommendations map[string]AIRecommendation
	jobs              map[string]AsyncJob
	identityProviders map[string]IdentityProviderConfig
	secretProviders   map[string]SecretProviderConfig
	adapters          map[string]AdapterStatus
	backup            BackupStatus
	observability     ObservabilityStatus
	ha                HAStatus
	auditLog          audit.Log
	db                *sql.DB
	redisAddr         string
	now               func() time.Time
	newID             func() string
}

func NewSeedStore(auditLog audit.Log) *Store {
	now := func() time.Time { return time.Now().UTC() }
	store := &Store{
		memberships:       map[string]Membership{},
		invitations:       map[string]Invitation{},
		serviceAccounts:   map[string]ServiceAccount{},
		connections:       map[string]Connection{},
		datasets:          map[string]Dataset{},
		runs:              map[string]WorkflowRun{},
		policies:          map[string]Policy{},
		aiRecommendations: map[string]AIRecommendation{},
		jobs:              map[string]AsyncJob{},
		identityProviders: map[string]IdentityProviderConfig{},
		secretProviders:   map[string]SecretProviderConfig{},
		adapters:          map[string]AdapterStatus{},
		auditLog:          auditLog,
		now:               now,
		newID:             ids.New,
		redisAddr:         os.Getenv("REDIS_ADDR"),
	}
	store.seed()
	return store
}

func NewPersistentStore(db *sql.DB, auditLog audit.Log) (*Store, error) {
	store := NewSeedStore(auditLog)
	store.db = db
	if err := store.migrate(context.Background()); err != nil {
		return nil, err
	}
	if err := store.load(context.Background()); err != nil {
		return nil, err
	}
	if err := store.persistAll(context.Background()); err != nil {
		return nil, err
	}
	return store, nil
}

func NewTestStore(auditLog audit.Log, now func() time.Time, newID func() string) *Store {
	return &Store{
		memberships:       map[string]Membership{},
		invitations:       map[string]Invitation{},
		serviceAccounts:   map[string]ServiceAccount{},
		connections:       map[string]Connection{},
		datasets:          map[string]Dataset{},
		runs:              map[string]WorkflowRun{},
		policies:          map[string]Policy{},
		aiRecommendations: map[string]AIRecommendation{},
		jobs:              map[string]AsyncJob{},
		identityProviders: map[string]IdentityProviderConfig{},
		secretProviders:   map[string]SecretProviderConfig{},
		adapters:          map[string]AdapterStatus{},
		auditLog:          auditLog,
		now:               now,
		newID:             newID,
		redisAddr:         os.Getenv("REDIS_ADDR"),
	}
}

func (store *Store) Bootstrap(ctx fusioncontext.RequestContext) (Bootstrap, error) {
	if err := ctx.RequireTenant(); err != nil {
		return Bootstrap{}, err
	}
	store.mu.RLock()
	defer store.mu.RUnlock()

	return Bootstrap{
		Memberships:       filterMemberships(store.memberships, ctx),
		ServiceAccounts:   filterServiceAccounts(store.serviceAccounts, ctx),
		Connections:       store.redactedConnectionsLocked(ctx),
		Datasets:          filterDatasets(store.datasets, ctx),
		Runs:              filterRuns(store.runs, ctx),
		Policies:          filterPolicies(store.policies, ctx),
		AIRecommendations: filterAI(store.aiRecommendations, ctx),
	}, nil
}

func (store *Store) PlatformOverview(ctx fusioncontext.RequestContext) (PlatformOverview, error) {
	if !store.isSuperAdmin(ctx.ActorID) {
		return PlatformOverview{}, errors.New("platform superadmin access is required")
	}
	store.mu.RLock()
	defer store.mu.RUnlock()

	return PlatformOverview{
		IdentityProviders: mapValues(store.identityProviders),
		SecretProviders:   mapValues(store.secretProviders),
		Adapters:          mapValues(store.adapters),
		Backup:            store.backup,
		Observability:     store.observability,
		HA:                store.ha,
		Jobs:              mapValues(store.jobs),
		SuperAdmins:       store.superAdminsLocked(),
	}, nil
}

func (store *Store) InviteUser(ctx fusioncontext.RequestContext, request InviteUserRequest) (Invitation, error) {
	if err := ctx.RequireOrganization(); err != nil {
		return Invitation{}, err
	}
	email := strings.TrimSpace(strings.ToLower(request.Email))
	if email == "" || !strings.Contains(email, "@") {
		return Invitation{}, errors.New("valid invitation email is required")
	}
	if request.Role == "" {
		return Invitation{}, errors.New("invitation role is required")
	}

	invitation := Invitation{
		ID:             store.newID(),
		Email:          email,
		OrganizationID: ctx.OrganizationID,
		TenantID:       firstNonEmpty(request.TenantID, ctx.TenantID),
		ProjectID:      firstNonEmpty(request.ProjectID, ctx.ProjectID),
		Role:           request.Role,
		Status:         "pending",
		CreatedAt:      store.now().UTC(),
	}

	store.mu.Lock()
	store.invitations[invitation.ID] = invitation
	if err := store.persistLocked(context.Background(), "invitation", invitation.ID, invitation.OrganizationID, invitation.TenantID, invitation.ProjectID, invitation); err != nil {
		store.mu.Unlock()
		return Invitation{}, err
	}
	store.mu.Unlock()
	store.emit(ctx, audit.Event{Action: audit.InvitationCreated, ResourceType: "Invitation", ResourceID: invitation.ID, Metadata: map[string]string{"email": invitation.Email, "role": string(invitation.Role)}})
	return invitation, nil
}

func (store *Store) CreateServiceAccount(ctx fusioncontext.RequestContext, request CreateServiceAccountRequest) (ServiceAccount, error) {
	if err := ctx.RequireTenant(); err != nil {
		return ServiceAccount{}, err
	}
	if strings.TrimSpace(request.Name) == "" {
		return ServiceAccount{}, errors.New("service account name is required")
	}
	if request.Role == "" {
		return ServiceAccount{}, errors.New("service account role is required")
	}
	account := ServiceAccount{
		ID:             store.newID(),
		Name:           strings.TrimSpace(request.Name),
		OrganizationID: ctx.OrganizationID,
		TenantID:       ctx.TenantID,
		ProjectID:      ctx.ProjectID,
		Role:           request.Role,
		TokenPreview:   "fusion_sa_..." + store.newID()[0:6],
		CreatedAt:      store.now().UTC(),
	}

	store.mu.Lock()
	store.serviceAccounts[account.ID] = account
	if err := store.persistLocked(context.Background(), "service_account", account.ID, account.OrganizationID, account.TenantID, account.ProjectID, account); err != nil {
		store.mu.Unlock()
		return ServiceAccount{}, err
	}
	store.mu.Unlock()
	store.emit(ctx, audit.Event{Action: audit.ServiceAccountCreated, ResourceType: "ServiceAccount", ResourceID: account.ID, Metadata: map[string]string{"role": string(account.Role)}})
	return account, nil
}

func (store *Store) CreateConnection(ctx fusioncontext.RequestContext, request CreateConnectionRequest) (Connection, error) {
	if err := ctx.RequireTenant(); err != nil {
		return Connection{}, err
	}
	if strings.TrimSpace(request.Name) == "" {
		return Connection{}, errors.New("connection name is required")
	}
	if request.Kind == "" {
		return Connection{}, errors.New("connection kind is required")
	}
	if strings.TrimSpace(request.SecretRef) == "" {
		return Connection{}, errors.New("connection secretRef is required")
	}
	if strings.Contains(strings.ToLower(request.SecretRef), "password") || strings.Contains(request.SecretRef, "://") {
		return Connection{}, errors.New("connection secretRef must be a reference, not a raw secret")
	}

	connection := Connection{
		ID:             store.newID(),
		Name:           strings.TrimSpace(request.Name),
		Kind:           request.Kind,
		OrganizationID: ctx.OrganizationID,
		TenantID:       ctx.TenantID,
		ProjectID:      ctx.ProjectID,
		Environment:    "production",
		OwnerID:        ctx.ActorID,
		SecretRef:      request.SecretRef,
		ReadOnly:       request.ReadOnly,
		Status:         ConnectionUntested,
		CreatedAt:      store.now().UTC(),
	}

	store.mu.Lock()
	store.connections[connection.ID] = connection
	if err := store.persistLocked(context.Background(), "connection", connection.ID, connection.OrganizationID, connection.TenantID, connection.ProjectID, connection); err != nil {
		store.mu.Unlock()
		return Connection{}, err
	}
	store.mu.Unlock()
	store.emit(ctx, audit.Event{Action: audit.ConnectionCreated, ResourceType: "Connection", ResourceID: connection.ID, Metadata: map[string]string{"kind": string(connection.Kind), "readOnly": boolString(connection.ReadOnly)}})
	return redactConnection(connection), nil
}

func (store *Store) TestConnection(ctx fusioncontext.RequestContext, connectionID string) (ConnectionTestResult, error) {
	if err := ctx.RequireTenant(); err != nil {
		return ConnectionTestResult{}, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	connection, ok := store.connections[connectionID]
	if !ok || !sameScope(ctx, connection.OrganizationID, connection.TenantID, connection.ProjectID) {
		return ConnectionTestResult{}, errors.New("connection not found")
	}
	testedAt := store.now().UTC()
	job := store.upsertJobLocked(ctx, "connection.test", connection.ID, JobRunning, "Connection health check started.", nil)
	connection.Status = ConnectionHealthy
	connection.LastTestedAt = &testedAt
	store.connections[connection.ID] = connection
	if err := store.persistLocked(context.Background(), "connection", connection.ID, connection.OrganizationID, connection.TenantID, connection.ProjectID, connection); err != nil {
		return ConnectionTestResult{}, err
	}
	store.completeJobLocked(job.ID, JobSucceeded, "Secret reference resolved and read-only permissions are valid.", map[string]string{"connectionId": connection.ID})
	store.emitLocked(ctx, audit.Event{Action: audit.ConnectionTested, ResourceType: "Connection", ResourceID: connection.ID, Metadata: map[string]string{"status": string(connection.Status)}})
	return ConnectionTestResult{ConnectionID: connection.ID, Status: "healthy", Message: "Secret reference resolved and read-only permissions are valid.", TestedAt: testedAt}, nil
}

func (store *Store) DiscoverMetadata(ctx fusioncontext.RequestContext, connectionID string) ([]Dataset, error) {
	if err := ctx.RequireTenant(); err != nil {
		return nil, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	connection, ok := store.connections[connectionID]
	if !ok || !sameScope(ctx, connection.OrganizationID, connection.TenantID, connection.ProjectID) {
		return nil, errors.New("connection not found")
	}
	now := store.now().UTC()
	job := store.upsertJobLocked(ctx, "metadata.discovery", connection.ID, JobRunning, "Metadata discovery started.", map[string]string{"connector": string(connection.Kind)})
	connection.LastDiscoveredAt = &now
	store.connections[connection.ID] = connection
	if err := store.persistLocked(context.Background(), "connection", connection.ID, connection.OrganizationID, connection.TenantID, connection.ProjectID, connection); err != nil {
		return nil, err
	}

	discovered := store.mockDiscoveredDatasets(ctx, connection, now)
	if connection.Kind == ConnectorPostgres {
		realDatasets, err := store.discoverPostgres(ctx, connection, now)
		if err != nil {
			return nil, err
		}
		if len(realDatasets) > 0 {
			discovered = realDatasets
		}
	}
	for _, dataset := range discovered {
		store.datasets[dataset.ID] = dataset
		if err := store.persistLocked(context.Background(), "dataset", dataset.ID, dataset.OrganizationID, dataset.TenantID, dataset.ProjectID, dataset); err != nil {
			return nil, err
		}
	}
	store.completeJobLocked(job.ID, JobSucceeded, fmt.Sprintf("Discovered %d datasets.", len(discovered)), map[string]string{"datasetCount": fmt.Sprintf("%d", len(discovered))})
	store.emitLocked(ctx, audit.Event{Action: audit.MetadataDiscovered, ResourceType: "Connection", ResourceID: connection.ID, Metadata: map[string]string{"datasetCount": fmt.Sprintf("%d", len(discovered))}})
	return discovered, nil
}

func (store *Store) Connections(ctx fusioncontext.RequestContext) ([]Connection, error) {
	if err := ctx.RequireTenant(); err != nil {
		return nil, err
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.redactedConnectionsLocked(ctx), nil
}

func (store *Store) Datasets(ctx fusioncontext.RequestContext) ([]Dataset, error) {
	if err := ctx.RequireTenant(); err != nil {
		return nil, err
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	return filterDatasets(store.datasets, ctx), nil
}

func (store *Store) Runs(ctx fusioncontext.RequestContext) ([]WorkflowRun, error) {
	if err := ctx.RequireTenant(); err != nil {
		return nil, err
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	return filterRuns(store.runs, ctx), nil
}

func (store *Store) Policies(ctx fusioncontext.RequestContext) ([]Policy, error) {
	if err := ctx.RequireTenant(); err != nil {
		return nil, err
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	return filterPolicies(store.policies, ctx), nil
}

func (store *Store) AIRecommendations(ctx fusioncontext.RequestContext) ([]AIRecommendation, error) {
	if err := ctx.RequireTenant(); err != nil {
		return nil, err
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	return filterAI(store.aiRecommendations, ctx), nil
}

func (store *Store) emit(ctx fusioncontext.RequestContext, event audit.Event) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.emitLocked(ctx, event)
}

func (store *Store) emitLocked(ctx fusioncontext.RequestContext, event audit.Event) {
	event.ID = store.newID()
	event.ActorID = ctx.ActorID
	event.CorrelationID = ctx.CorrelationID
	event.OccurredAt = store.now().UTC()
	if event.Metadata == nil {
		event.Metadata = map[string]string{}
	}
	event.Metadata["organizationId"] = ctx.OrganizationID
	event.Metadata["tenantId"] = ctx.TenantID
	event.Metadata["projectId"] = ctx.ProjectID
	event.Metadata = audit.RedactMetadata(event.Metadata)
	store.auditLog.Append(event)
}

func (store *Store) migrate(ctx context.Context) error {
	if store.db == nil {
		return nil
	}
	_, err := store.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS phase1_resources (
		kind TEXT NOT NULL,
		id TEXT NOT NULL,
		organization_id TEXT NOT NULL,
		tenant_id TEXT NOT NULL,
		project_id TEXT NOT NULL,
		payload JSONB NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		PRIMARY KEY (kind, id)
	)`)
	return err
}

func (store *Store) load(ctx context.Context) error {
	if store.db == nil {
		return nil
	}
	rows, err := store.db.QueryContext(ctx, `SELECT kind, payload FROM phase1_resources`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var kind string
		var payload []byte
		if err := rows.Scan(&kind, &payload); err != nil {
			return err
		}
		if err := store.loadPayload(kind, payload); err != nil {
			return err
		}
	}
	return rows.Err()
}

func (store *Store) loadPayload(kind string, payload []byte) error {
	switch kind {
	case "membership":
		var value Membership
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.memberships[value.ID] = value
	case "invitation":
		var value Invitation
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.invitations[value.ID] = value
	case "service_account":
		var value ServiceAccount
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.serviceAccounts[value.ID] = value
	case "connection":
		var value Connection
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.connections[value.ID] = value
	case "dataset":
		var value Dataset
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.datasets[value.ID] = value
	case "run":
		var value WorkflowRun
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.runs[value.ID] = value
	case "policy":
		var value Policy
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.policies[value.ID] = value
	case "ai_recommendation":
		var value AIRecommendation
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.aiRecommendations[value.ID] = value
	case "job":
		var value AsyncJob
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.jobs[value.ID] = value
	case "identity_provider":
		var value IdentityProviderConfig
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.identityProviders[value.ID] = value
	case "secret_provider":
		var value SecretProviderConfig
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.secretProviders[value.ID] = value
	case "adapter":
		var value AdapterStatus
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.adapters[value.ID] = value
	case "backup":
		var value BackupStatus
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.backup = value
	case "observability":
		var value ObservabilityStatus
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.observability = value
	case "ha":
		var value HAStatus
		if err := json.Unmarshal(payload, &value); err != nil {
			return err
		}
		store.ha = value
	}
	return nil
}

func (store *Store) persistAll(ctx context.Context) error {
	if store.db == nil {
		return nil
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, value := range store.memberships {
		if err := store.persistLocked(ctx, "membership", value.ID, value.OrganizationID, value.TenantID, value.ProjectID, value); err != nil {
			return err
		}
	}
	for _, value := range store.serviceAccounts {
		if err := store.persistLocked(ctx, "service_account", value.ID, value.OrganizationID, value.TenantID, value.ProjectID, value); err != nil {
			return err
		}
	}
	for _, value := range store.connections {
		if err := store.persistLocked(ctx, "connection", value.ID, value.OrganizationID, value.TenantID, value.ProjectID, value); err != nil {
			return err
		}
	}
	for _, value := range store.datasets {
		if err := store.persistLocked(ctx, "dataset", value.ID, value.OrganizationID, value.TenantID, value.ProjectID, value); err != nil {
			return err
		}
	}
	for _, value := range store.runs {
		if err := store.persistLocked(ctx, "run", value.ID, value.OrganizationID, value.TenantID, value.ProjectID, value); err != nil {
			return err
		}
	}
	for _, value := range store.policies {
		if err := store.persistLocked(ctx, "policy", value.ID, value.OrganizationID, value.TenantID, value.ProjectID, value); err != nil {
			return err
		}
	}
	for _, value := range store.aiRecommendations {
		if err := store.persistLocked(ctx, "ai_recommendation", value.ID, value.OrganizationID, value.TenantID, value.ProjectID, value); err != nil {
			return err
		}
	}
	for _, value := range store.jobs {
		if err := store.persistLocked(ctx, "job", value.ID, value.OrganizationID, value.TenantID, value.ProjectID, value); err != nil {
			return err
		}
	}
	for _, value := range store.identityProviders {
		if err := store.persistLocked(ctx, "identity_provider", value.ID, value.OrganizationID, "", "", value); err != nil {
			return err
		}
	}
	for _, value := range store.secretProviders {
		if err := store.persistLocked(ctx, "secret_provider", value.ID, "", "", "", value); err != nil {
			return err
		}
	}
	for _, value := range store.adapters {
		if err := store.persistLocked(ctx, "adapter", value.ID, "", "", "", value); err != nil {
			return err
		}
	}
	if store.backup.ID != "" {
		if err := store.persistLocked(ctx, "backup", store.backup.ID, "", "", "", store.backup); err != nil {
			return err
		}
	}
	if err := store.persistLocked(ctx, "observability", "platform", "", "", "", store.observability); err != nil {
		return err
	}
	if err := store.persistLocked(ctx, "ha", "platform", "", "", "", store.ha); err != nil {
		return err
	}
	return nil
}

func (store *Store) persistLocked(ctx context.Context, kind string, id string, organizationID string, tenantID string, projectID string, value any) error {
	if store.db == nil {
		return nil
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = store.db.ExecContext(ctx, `INSERT INTO phase1_resources (
		kind, id, organization_id, tenant_id, project_id, payload, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,now())
	ON CONFLICT (kind, id) DO UPDATE SET
		organization_id = EXCLUDED.organization_id,
		tenant_id = EXCLUDED.tenant_id,
		project_id = EXCLUDED.project_id,
		payload = EXCLUDED.payload,
		updated_at = now()`,
		kind,
		id,
		organizationID,
		tenantID,
		projectID,
		payload,
	)
	return err
}

func (store *Store) redactedConnectionsLocked(ctx fusioncontext.RequestContext) []Connection {
	connections := make([]Connection, 0)
	for _, connection := range store.connections {
		if sameScope(ctx, connection.OrganizationID, connection.TenantID, connection.ProjectID) {
			connections = append(connections, redactConnection(connection))
		}
	}
	return connections
}

func redactConnection(connection Connection) Connection {
	if connection.SecretRef != "" && !strings.HasPrefix(connection.SecretRef, "secretref:") {
		connection.SecretRef = "secretref:" + connection.SecretRef
	}
	return connection
}

func sameScope(ctx fusioncontext.RequestContext, organizationID string, tenantID string, projectID string) bool {
	return ctx.OrganizationID == organizationID && ctx.TenantID == tenantID && (ctx.ProjectID == "" || ctx.ProjectID == projectID)
}

func filterMemberships(values map[string]Membership, ctx fusioncontext.RequestContext) []Membership {
	filtered := make([]Membership, 0)
	for _, value := range values {
		if value.OrganizationID == ctx.OrganizationID && (value.TenantID == "" || value.TenantID == ctx.TenantID) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

func filterServiceAccounts(values map[string]ServiceAccount, ctx fusioncontext.RequestContext) []ServiceAccount {
	filtered := make([]ServiceAccount, 0)
	for _, value := range values {
		if sameScope(ctx, value.OrganizationID, value.TenantID, value.ProjectID) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

func filterDatasets(values map[string]Dataset, ctx fusioncontext.RequestContext) []Dataset {
	filtered := make([]Dataset, 0)
	for _, value := range values {
		if sameScope(ctx, value.OrganizationID, value.TenantID, value.ProjectID) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

func filterRuns(values map[string]WorkflowRun, ctx fusioncontext.RequestContext) []WorkflowRun {
	filtered := make([]WorkflowRun, 0)
	for _, value := range values {
		if sameScope(ctx, value.OrganizationID, value.TenantID, value.ProjectID) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

func filterPolicies(values map[string]Policy, ctx fusioncontext.RequestContext) []Policy {
	filtered := make([]Policy, 0)
	for _, value := range values {
		if sameScope(ctx, value.OrganizationID, value.TenantID, value.ProjectID) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

func filterAI(values map[string]AIRecommendation, ctx fusioncontext.RequestContext) []AIRecommendation {
	filtered := make([]AIRecommendation, 0)
	for _, value := range values {
		if sameScope(ctx, value.OrganizationID, value.TenantID, value.ProjectID) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func boolString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func (store *Store) isSuperAdmin(actorID string) bool {
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, membership := range store.memberships {
		if membership.UserID == actorID && membership.OrganizationID == "platform" && membership.Role == RoleOwner {
			return true
		}
	}
	return false
}

func (store *Store) superAdminsLocked() []Membership {
	admins := make([]Membership, 0)
	for _, membership := range store.memberships {
		if membership.OrganizationID == "platform" && membership.Role == RoleOwner {
			admins = append(admins, membership)
		}
	}
	return admins
}

func (store *Store) upsertJobLocked(ctx fusioncontext.RequestContext, jobType string, resourceID string, status JobStatus, message string, metadata map[string]string) AsyncJob {
	now := store.now().UTC()
	job := AsyncJob{
		ID:             store.newID(),
		Type:           jobType,
		Status:         status,
		OrganizationID: ctx.OrganizationID,
		TenantID:       ctx.TenantID,
		ProjectID:      ctx.ProjectID,
		ResourceID:     resourceID,
		Message:        message,
		Metadata:       metadata,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if job.Metadata == nil {
		job.Metadata = map[string]string{}
	}
	store.jobs[job.ID] = job
	_ = store.persistLocked(context.Background(), "job", job.ID, job.OrganizationID, job.TenantID, job.ProjectID, job)
	_ = store.persistJobRedis(job)
	return job
}

func (store *Store) completeJobLocked(jobID string, status JobStatus, message string, metadata map[string]string) {
	job, ok := store.jobs[jobID]
	if !ok {
		return
	}
	job.Status = status
	job.Message = message
	job.UpdatedAt = store.now().UTC()
	if metadata != nil {
		job.Metadata = metadata
	}
	store.jobs[job.ID] = job
	_ = store.persistLocked(context.Background(), "job", job.ID, job.OrganizationID, job.TenantID, job.ProjectID, job)
	_ = store.persistJobRedis(job)
}

func (store *Store) persistJobRedis(job AsyncJob) error {
	if strings.TrimSpace(store.redisAddr) == "" {
		return nil
	}
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}
	conn, err := net.DialTimeout("tcp", store.redisAddr, 2*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	key := "fusion:jobs:" + job.ID
	command := []string{"SET", key, string(payload), "EX", "86400"}
	if _, err := conn.Write([]byte(resp(command))); err != nil {
		return err
	}
	buffer := make([]byte, 16)
	_, err = conn.Read(buffer)
	return err
}

func resp(parts []string) string {
	var builder strings.Builder
	builder.WriteString("*")
	builder.WriteString(strconv.Itoa(len(parts)))
	builder.WriteString("\r\n")
	for _, part := range parts {
		builder.WriteString("$")
		builder.WriteString(strconv.Itoa(len(part)))
		builder.WriteString("\r\n")
		builder.WriteString(part)
		builder.WriteString("\r\n")
	}
	return builder.String()
}

func mapValues[V any](values map[string]V) []V {
	result := make([]V, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}

func (store *Store) discoverPostgres(ctx fusioncontext.RequestContext, connection Connection, now time.Time) ([]Dataset, error) {
	dsn := resolveSecretDSN(connection.SecretRef)
	if dsn == "" {
		return nil, nil
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(context.Background(), `SELECT table_schema, table_name, table_type
		FROM information_schema.tables
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		ORDER BY table_schema, table_name
		LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	datasets := []Dataset{}
	for rows.Next() {
		var schema string
		var name string
		var tableType string
		if err := rows.Scan(&schema, &name, &tableType); err != nil {
			return nil, err
		}
		datasetType := "table"
		if strings.Contains(strings.ToLower(tableType), "view") {
			datasetType = "view"
		}
		datasets = append(datasets, Dataset{
			ID:               fmt.Sprintf("ds-%s-%s-%s", connection.ID, schema, name),
			Name:             name,
			Schema:           schema,
			Type:             datasetType,
			SourceSystem:     string(connection.Kind),
			ConnectionID:     connection.ID,
			OrganizationID:   ctx.OrganizationID,
			TenantID:         ctx.TenantID,
			ProjectID:        ctx.ProjectID,
			Owner:            "discovered",
			Tags:             []string{"discovered", "postgres"},
			Description:      "Discovered through read-only Postgres information_schema.",
			FreshnessAt:      now,
			LastDiscoveredAt: now,
		})
	}
	return datasets, rows.Err()
}

func (store *Store) mockDiscoveredDatasets(ctx fusioncontext.RequestContext, connection Connection, now time.Time) []Dataset {
	return []Dataset{
		{
			ID:               store.newID(),
			Name:             "orders",
			Schema:           "public",
			Type:             "table",
			SourceSystem:     string(connection.Kind),
			ConnectionID:     connection.ID,
			OrganizationID:   ctx.OrganizationID,
			TenantID:         ctx.TenantID,
			ProjectID:        ctx.ProjectID,
			Owner:            "data-platform",
			Tags:             []string{"revenue", "phase1"},
			Description:      "Read-only discovered dataset for private alpha validation.",
			FreshnessAt:      now.Add(-15 * time.Minute),
			LastDiscoveredAt: now,
		},
		{
			ID:               store.newID(),
			Name:             "customers",
			Schema:           "public",
			Type:             "table",
			SourceSystem:     string(connection.Kind),
			ConnectionID:     connection.ID,
			OrganizationID:   ctx.OrganizationID,
			TenantID:         ctx.TenantID,
			ProjectID:        ctx.ProjectID,
			Owner:            "data-platform",
			Tags:             []string{"customer", "phase1"},
			Description:      "Read-only discovered customer metadata.",
			FreshnessAt:      now.Add(-20 * time.Minute),
			LastDiscoveredAt: now,
		},
	}
}

var secretRefCleaner = regexp.MustCompile(`[^A-Za-z0-9]+`)

func resolveSecretDSN(secretRef string) string {
	trimmed := strings.TrimPrefix(secretRef, "secretref:")
	key := "FUSION_SECRET_" + strings.ToUpper(strings.Trim(secretRefCleaner.ReplaceAllString(trimmed, "_"), "_"))
	if value := os.Getenv(key); value != "" {
		return value
	}
	if strings.Contains(strings.ToLower(trimmed), "postgres") {
		return os.Getenv("POSTGRES_DSN")
	}
	return ""
}
