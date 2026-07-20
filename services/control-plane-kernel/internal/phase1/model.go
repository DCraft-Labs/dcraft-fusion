package phase1

import "time"

type Role string

const (
	RoleOwner    Role = "owner"
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
	RoleViewer   Role = "viewer"
)

type Membership struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	Email          string    `json:"email"`
	OrganizationID string    `json:"organizationId"`
	TenantID       string    `json:"tenantId,omitempty"`
	ProjectID      string    `json:"projectId,omitempty"`
	Role           Role      `json:"role"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Invitation struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	OrganizationID string    `json:"organizationId"`
	TenantID       string    `json:"tenantId,omitempty"`
	ProjectID      string    `json:"projectId,omitempty"`
	Role           Role      `json:"role"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
}

type InviteUserRequest struct {
	Email     string `json:"email"`
	TenantID  string `json:"tenantId,omitempty"`
	ProjectID string `json:"projectId,omitempty"`
	Role      Role   `json:"role"`
}

type ServiceAccount struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	OrganizationID string    `json:"organizationId"`
	TenantID       string    `json:"tenantId"`
	ProjectID      string    `json:"projectId"`
	Role           Role      `json:"role"`
	TokenPreview   string    `json:"tokenPreview"`
	CreatedAt      time.Time `json:"createdAt"`
}

type CreateServiceAccountRequest struct {
	Name string `json:"name"`
	Role Role   `json:"role"`
}

type ConnectorKind string

const (
	ConnectorPostgres ConnectorKind = "postgres"
	ConnectorAirflow  ConnectorKind = "airflow"
	ConnectorGitHub   ConnectorKind = "github"
	ConnectorSlack    ConnectorKind = "slack"
	ConnectorCDC      ConnectorKind = "fusion_cdc_engine"
)

type ConnectionStatus string

const (
	ConnectionUntested ConnectionStatus = "untested"
	ConnectionHealthy  ConnectionStatus = "healthy"
	ConnectionBlocked  ConnectionStatus = "blocked"
)

type Connection struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Kind             ConnectorKind    `json:"kind"`
	OrganizationID   string           `json:"organizationId"`
	TenantID         string           `json:"tenantId"`
	ProjectID        string           `json:"projectId"`
	Environment      string           `json:"environment"`
	OwnerID          string           `json:"ownerId"`
	SecretRef        string           `json:"secretRef"`
	ReadOnly         bool             `json:"readOnly"`
	Status           ConnectionStatus `json:"status"`
	LastTestedAt     *time.Time       `json:"lastTestedAt,omitempty"`
	LastDiscoveredAt *time.Time       `json:"lastDiscoveredAt,omitempty"`
	CreatedAt        time.Time        `json:"createdAt"`
}

type CreateConnectionRequest struct {
	Name      string        `json:"name"`
	Kind      ConnectorKind `json:"kind"`
	SecretRef string        `json:"secretRef"`
	ReadOnly  bool          `json:"readOnly"`
}

type ConnectionTestResult struct {
	ConnectionID string    `json:"connectionId"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	TestedAt     time.Time `json:"testedAt"`
}

type Dataset struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Schema           string    `json:"schema"`
	Type             string    `json:"type"`
	SourceSystem     string    `json:"sourceSystem"`
	ConnectionID     string    `json:"connectionId"`
	OrganizationID   string    `json:"organizationId"`
	TenantID         string    `json:"tenantId"`
	ProjectID        string    `json:"projectId"`
	Owner            string    `json:"owner"`
	Tags             []string  `json:"tags"`
	Description      string    `json:"description"`
	FreshnessAt      time.Time `json:"freshnessAt"`
	LastDiscoveredAt time.Time `json:"lastDiscoveredAt"`
}

type RunStatus string

const (
	RunQueued    RunStatus = "queued"
	RunRunning   RunStatus = "running"
	RunSucceeded RunStatus = "succeeded"
	RunFailed    RunStatus = "failed"
	RunCanceled  RunStatus = "canceled"
	RunUnknown   RunStatus = "unknown"
)

type WorkflowRun struct {
	ID             string     `json:"id"`
	WorkflowName   string     `json:"workflowName"`
	Status         RunStatus  `json:"status"`
	OrganizationID string     `json:"organizationId"`
	TenantID       string     `json:"tenantId"`
	ProjectID      string     `json:"projectId"`
	Environment    string     `json:"environment"`
	CorrelationID  string     `json:"correlationId"`
	ArtifactRef    string     `json:"artifactRef"`
	StartedAt      time.Time  `json:"startedAt"`
	FinishedAt     *time.Time `json:"finishedAt,omitempty"`
}

type ApprovalState string

const (
	ApprovalDraft    ApprovalState = "draft"
	ApprovalPending  ApprovalState = "pending_review"
	ApprovalApproved ApprovalState = "approved"
	ApprovalRejected ApprovalState = "rejected"
)

type Policy struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	OrganizationID string `json:"organizationId"`
	TenantID       string `json:"tenantId"`
	ProjectID      string `json:"projectId"`
	Mode           string `json:"mode"`
}

type AIRecommendation struct {
	ID             string        `json:"id"`
	OrganizationID string        `json:"organizationId"`
	TenantID       string        `json:"tenantId"`
	ProjectID      string        `json:"projectId"`
	Title          string        `json:"title"`
	Summary        string        `json:"summary"`
	SourceRefs     []string      `json:"sourceRefs"`
	State          ApprovalState `json:"state"`
	CreatedAt      time.Time     `json:"createdAt"`
}

type Bootstrap struct {
	Memberships       []Membership       `json:"memberships"`
	ServiceAccounts   []ServiceAccount   `json:"serviceAccounts"`
	Connections       []Connection       `json:"connections"`
	Datasets          []Dataset          `json:"datasets"`
	Runs              []WorkflowRun      `json:"runs"`
	Policies          []Policy           `json:"policies"`
	AIRecommendations []AIRecommendation `json:"aiRecommendations"`
}

type JobStatus string

const (
	JobQueued    JobStatus = "queued"
	JobRunning   JobStatus = "running"
	JobSucceeded JobStatus = "succeeded"
	JobFailed    JobStatus = "failed"
)

type AsyncJob struct {
	ID             string            `json:"id"`
	Type           string            `json:"type"`
	Status         JobStatus         `json:"status"`
	OrganizationID string            `json:"organizationId"`
	TenantID       string            `json:"tenantId"`
	ProjectID      string            `json:"projectId"`
	ResourceID     string            `json:"resourceId"`
	Message        string            `json:"message"`
	Metadata       map[string]string `json:"metadata"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

type IdentityProviderConfig struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Protocol       string    `json:"protocol"`
	Issuer         string    `json:"issuer"`
	Audience       string    `json:"audience"`
	JWKSURL        string    `json:"jwksUrl"`
	SSOEnabled     bool      `json:"ssoEnabled"`
	SAMLEnabled    bool      `json:"samlEnabled"`
	Status         string    `json:"status"`
	OrganizationID string    `json:"organizationId,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

type SecretProviderConfig struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	Endpoint  string    `json:"endpoint"`
	Scope     string    `json:"scope"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type AdapterStatus struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Kind        string    `json:"kind"`
	ReadOnly    bool      `json:"readOnly"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type BackupStatus struct {
	ID                  string    `json:"id"`
	Target              string    `json:"target"`
	Schedule            string    `json:"schedule"`
	LastVerifiedAt      time.Time `json:"lastVerifiedAt"`
	LastVerificationJob string    `json:"lastVerificationJob"`
	Status              string    `json:"status"`
	Runbook             string    `json:"runbook"`
}

type ObservabilityStatus struct {
	TracesEnabled      bool     `json:"tracesEnabled"`
	MetricsEndpoint    string   `json:"metricsEndpoint"`
	LogFormat          string   `json:"logFormat"`
	Dashboards         []string `json:"dashboards"`
	Alerts             []string `json:"alerts"`
	OpenTelemetryReady bool     `json:"openTelemetryReady"`
}

type HAStatus struct {
	APIMinReplicas int      `json:"apiMinReplicas"`
	WebMinReplicas int      `json:"webMinReplicas"`
	Autoscaling    bool     `json:"autoscaling"`
	PodDisruption  bool     `json:"podDisruption"`
	LoadTests      []string `json:"loadTests"`
	Status         string   `json:"status"`
}

type PlatformProject struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Environment string    `json:"environment"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PlatformTenant struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Model     string            `json:"model"`
	Isolation string            `json:"isolation"`
	Region    string            `json:"region"`
	CreatedAt time.Time         `json:"createdAt"`
	Projects  []PlatformProject `json:"projects"`
}

type PlatformOrganization struct {
	ID                     string           `json:"id"`
	Name                   string           `json:"name"`
	Type                   string           `json:"type"`
	DeploymentProfile      string           `json:"deploymentProfile"`
	Region                 string           `json:"region"`
	DataPlaneLocation      string           `json:"dataPlaneLocation"`
	RawDataMovementAllowed bool             `json:"rawDataMovementAllowed"`
	CreatedAt              time.Time        `json:"createdAt"`
	Tenants                []PlatformTenant `json:"tenants"`
}

type PlatformAuditEvent struct {
	ID            string    `json:"id"`
	Action        string    `json:"action"`
	ActorID       string    `json:"actorId"`
	CorrelationID string    `json:"correlationId"`
	ResourceType  string    `json:"resourceType"`
	ResourceID    string    `json:"resourceId"`
	OccurredAt    time.Time `json:"occurredAt"`
}

type PlatformOverview struct {
	IdentityProviders []IdentityProviderConfig `json:"identityProviders"`
	SecretProviders   []SecretProviderConfig   `json:"secretProviders"`
	Adapters          []AdapterStatus          `json:"adapters"`
	Backup            BackupStatus             `json:"backup"`
	Observability     ObservabilityStatus      `json:"observability"`
	HA                HAStatus                 `json:"ha"`
	Jobs              []AsyncJob               `json:"jobs"`
	SuperAdmins       []Membership             `json:"superAdmins"`
	Organizations     []PlatformOrganization   `json:"organizations"`
	RecentAuditEvents []PlatformAuditEvent     `json:"recentAuditEvents"`
}
