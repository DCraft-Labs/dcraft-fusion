# DCraft Fusion - Technical Architecture V0 (FROZEN)

**Version:** 0.1.0  
**Status:** FROZEN - Implementation Blueprint  
**Date:** 2026-01-18  
**Architect Review Grade:** A-  

---

## 1. FINAL Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CONTROL PLANE (Stateful)                          │
│                                                                             │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                         API Server (HTTP/gRPC)                        │  │
│  │  - Authentication/Authorization (RBAC + ABAC)                        │  │
│  │  - Request validation & admission control                            │  │
│  │  - Rate limiting & quota enforcement                                 │  │
│  │  - Audit logging                                                     │  │
│  └────────────┬─────────────────────────────────────────┬───────────────┘  │
│               │                                         │                   │
│               ▼                                         ▼                   │
│  ┌────────────────────────┐              ┌──────────────────────────────┐  │
│  │   Control Plane DB     │              │    Event Bus (NATS)          │  │
│  │   (PostgreSQL)         │              │  - Resource change events    │  │
│  │                        │              │  - Controller coordination   │  │
│  │  - Resources (tenant)  │              │  - Audit stream              │  │
│  │  - System metadata     │              │  - Telemetry stream          │  │
│  │  - Tenant isolation    │              └──────────┬───────────────────┘  │
│  │    via RLS             │                         │                      │
│  └────────────┬───────────┘                         │                      │
│               │                                     │                      │
│               ▼                                     ▼                      │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │              Controller/Reconciler Runtime (Stateless)              │  │
│  │                                                                     │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │  │
│  │  │  Pipeline    │  │  Dataset     │  │  Quality     │  ... (20+) │  │
│  │  │  Controller  │  │  Controller  │  │  Controller  │            │  │
│  │  └──────────────┘  └──────────────┘  └──────────────┘            │  │
│  │                                                                     │  │
│  │  - Leader election per controller type                             │  │
│  │  - Horizontal sharding by tenant_id hash                           │  │
│  │  - Watch-based + periodic resync (5min default)                    │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                    Provider Registry Service                        │  │
│  │  - Provider manifest storage                                        │  │
│  │  - Capability discovery index                                       │  │
│  │  - Version compatibility matrix                                     │  │
│  │  - Certification tier enforcement                                   │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                    AI Operator Service (Optional)                   │  │
│  │  - LLM routing & prompt management                                  │  │
│  │  - Approval workflow state machine                                  │  │
│  │  - Policy enforcement engine                                        │  │
│  │  - Budget tracking per tenant                                       │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                    Lineage Graph Service                            │  │
│  │  (Neo4j or PostgreSQL with pg_graph extension)                      │  │
│  │  - Node: Dataset, Column, Pipeline, Model                           │  │
│  │  - Edge: produces, consumes, transforms, derives_from               │  │
│  │  - Optimized for upstream/downstream/blast-radius queries           │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │              Telemetry & Audit Store (ClickHouse)                   │  │
│  │  - Execution logs                                                   │  │
│  │  - Audit trail (immutable)                                          │  │
│  │  - Metrics time-series                                              │  │
│  │  - Retention: 90d hot, 2y cold                                      │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ mTLS (outbound-only from edge)
                                    │
┌─────────────────────────────────────────────────────────────────────────────┐
│                          EDGE PLANE (Customer VPC/On-prem)                  │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                      Edge Agent (Single Binary)                     │  │
│  │                                                                     │  │
│  │  ┌──────────────────────┐         ┌──────────────────────┐        │  │
│  │  │   Exec Component     │         │  Observe Component   │        │  │
│  │  │  - Execute commands  │         │  - Metadata scraping │        │  │
│  │  │  - Run providers     │         │  - Log tailing       │        │  │
│  │  │  - Secret injection  │         │  - Metric collection │        │  │
│  │  └──────────────────────┘         └──────────────────────┘        │  │
│  │                                                                     │  │
│  │  - Outbound-only mTLS to control plane                             │  │
│  │  - Local provider runtime (sandboxed)                              │  │
│  │  - Certificate rotation every 30d                                  │  │
│  │  - Offline mode: queue + sync when online                          │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                    Provider Runtime Sandbox                         │  │
│  │  - Isolated execution per provider (Docker/gVisor)                 │  │
│  │  - Resource limits (CPU/Memory/Network)                            │  │
│  │  - Secret access via Unix socket                                   │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │              Execution Engines (BYO or Fusion-default)              │  │
│  │  - dbt, Airflow, Dagster, Prefect, Spark, Flink, etc.             │  │
│  │  - Providers translate Fusion resources → tool-native config       │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Reconciliation Model (FROZEN)

### 2.1 Core Strategy: Hybrid Event-Driven + Periodic Resync

**Decision:** Use event-driven reconciliation as primary mechanism with periodic full resync as safety net.

```
Primary: Watch-based (NATS event stream)
  ├─ Resource CREATE/UPDATE/DELETE → immediate reconciliation
  ├─ Controller subscribes to resource types it owns
  └─ Latency: <500ms for 95th percentile

Safety Net: Periodic resync
  ├─ Default interval: 5 minutes (configurable per controller)
  ├─ Full state comparison: spec vs observed state
  └─ Catches: missed events, external drift, clock skew
```

### 2.2 Drift Detection

**Mechanism:**
1. Controllers maintain `status.observedGeneration` field
2. Compare `metadata.generation` (spec version) vs `status.observedGeneration`
3. If mismatch → reconcile
4. External drift detected via periodic resync comparing actual state vs desired state

**Drift Sources:**
- Manual changes in execution engines (e.g., editing Airflow DAG directly)
- Provider failures that weren't reported
- Network partitions causing missed events
- Time-based state changes (e.g., schedule triggers)

### 2.3 Controller Leases and Sharding Strategy

**Leader Election:**
- Use PostgreSQL advisory locks for leader election per controller type
- Lease duration: 15 seconds
- Renewal interval: 5 seconds
- Lease lost → immediate graceful shutdown, new leader elected within 10s

**Horizontal Sharding:**
```
Sharding Key: hash(tenant_id) % num_shards
  ├─ Each controller instance claims 1+ shards
  ├─ Shard assignment stored in control_plane.controller_shards table
  ├─ Rebalancing on scale-up/down (max 10% shard movement)
  └─ Cross-shard operations use distributed transactions (rare)

Example with 3 controller instances, 12 shards:
  Instance A: shards [0,1,2,3]
  Instance B: shards [4,5,6,7]
  Instance C: shards [8,9,10,11]
```

### 2.4 Idempotency Contract

**MANDATORY for all controller actions:**

1. **Read-Modify-Write with Optimistic Locking:**
   ```
   UPDATE resources 
   SET spec = $1, metadata.generation = metadata.generation + 1
   WHERE id = $2 AND metadata.generation = $3
   RETURNING *;
   
   If affected_rows = 0 → conflict, retry from read
   ```

2. **External Actions Must Be Idempotent:**
   - Use unique request IDs (UUID) for provider calls
   - Providers must deduplicate based on request ID
   - Store action results in `status.actionResults` map keyed by request ID

3. **Status Updates Are Atomic:**
   ```
   PATCH /api/v1/pipelines/{id}/status
   {
     "status": {
       "phase": "Running",
       "observedGeneration": 5,
       "conditions": [...],
       "lastTransitionTime": "2026-01-18T10:00:00Z"
     }
   }
   ```

### 2.5 Conflict Resolution

**Scenario 1: Multiple controllers reconcile same resource (race condition)**
- **Resolution:** Optimistic locking prevents. Loser retries with fresh state.

**Scenario 2: User updates spec while controller is reconciling**
- **Resolution:** Controller's write fails due to generation mismatch. Controller re-reads and reconciles new spec.

**Scenario 3: Two controllers own different aspects of same resource**
- **Resolution:** NOT ALLOWED. Each resource type has exactly one controller. Use separate resources if multiple concerns (e.g., Pipeline + PipelineExecution).

**Scenario 4: External system modifies state (drift)**
- **Resolution:** 
  - If `status.driftDetection.enabled = true` → controller reverts to desired state
  - If `status.driftDetection.enabled = false` → controller updates status to reflect actual state, alerts user

---

## 3. Core Resource Model (FROZEN)

### 3.1 Resource Versioning Strategy

**Decision:** Immutable revisions for spec changes, mutable status.

```
metadata:
  generation: int64        # Increments on spec change (immutable)
  resourceVersion: string  # Increments on any change (spec or status)
  
Revision History:
  - Store last 10 revisions in resource_revisions table
  - Enables rollback and audit trail
  - Old revisions are immutable, compressed after 30d
```

### 3.2 Audit Model

**All mutations logged to audit_log table:**
```sql
CREATE TABLE audit_log (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL,
  actor_id UUID NOT NULL,  -- user or service account
  action VARCHAR(20) NOT NULL,  -- CREATE, UPDATE, DELETE, APPROVE
  resource_type VARCHAR(50) NOT NULL,
  resource_id UUID NOT NULL,
  resource_version VARCHAR(50),
  changes JSONB,  -- diff of spec
  ip_address INET,
  user_agent TEXT
);

-- Partitioned by timestamp (monthly)
-- Immutable (INSERT only)
-- Replicated to ClickHouse for long-term retention
```

### 3.3 Core Resources (20 Resources)

---

#### Resource 1: **Pipeline**

**Purpose:** Declarative definition of a data transformation workflow.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Pipeline
metadata:
  name: customer-analytics-daily
  namespace: production
  labels:
    team: analytics
    criticality: high
spec:
  description: "Daily customer behavior analytics"
  schedule:
    cron: "0 2 * * *"
    timezone: "America/Los_Angeles"
  provider:
    type: dbt  # or airflow, dagster, prefect, fusion-native
    version: "1.7.0"
    config:
      project_dir: "/opt/dbt/customer_analytics"
      profiles_dir: "/opt/dbt/profiles"
  tasks:
    - name: extract-events
      type: sql
      config:
        query: "SELECT * FROM raw.events WHERE date = {{ ds }}"
        destination: staging.events
    - name: transform-sessions
      type: dbt-model
      dependsOn: [extract-events]
      config:
        model: "models/sessions.sql"
    - name: load-to-warehouse
      type: sql
      dependsOn: [transform-sessions]
      config:
        query: "INSERT INTO prod.sessions SELECT * FROM staging.sessions"
  inputs:
    - datasetRef: raw.events
      mode: read
  outputs:
    - datasetRef: prod.sessions
      mode: write
  retryPolicy:
    maxRetries: 3
    backoff: exponential
  timeout: 3600s
  notifications:
    onFailure:
      - type: slack
        channel: "#data-alerts"
  driftDetection:
    enabled: true
    reconcileInterval: 300s
```

**Status:**
```yaml
status:
  phase: Running  # Pending, Running, Succeeded, Failed, Paused
  observedGeneration: 5
  conditions:
    - type: Scheduled
      status: "True"
      lastTransitionTime: "2026-01-18T02:00:00Z"
    - type: Healthy
      status: "True"
      lastTransitionTime: "2026-01-18T02:00:05Z"
  lastExecution:
    runId: exec-abc123
    startTime: "2026-01-18T02:00:00Z"
    endTime: "2026-01-18T02:45:00Z"
    status: Succeeded
  nextScheduledRun: "2026-01-19T02:00:00Z"
  providerStatus:
    syncedAt: "2026-01-18T02:00:10Z"
    externalId: "dbt-project-customer-analytics"
    driftDetected: false
```

**Relationships:**
- `inputs` → Dataset (many-to-many)
- `outputs` → Dataset (many-to-many)
- `owner` → Owner (many-to-one)
- `executions` → PipelineExecution (one-to-many)

---

#### Resource 2: **Dataset**

**Purpose:** Logical representation of a data asset (table, file, stream).

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Dataset
metadata:
  name: prod.sessions
  namespace: production
spec:
  description: "User session aggregates"
  type: table  # table, file, stream, model
  provider:
    type: snowflake
    config:
      database: PROD
      schema: ANALYTICS
      table: SESSIONS
  schema:
    schemaVersionRef: schema-sessions-v3
  quality:
    rules:
      - ruleRef: sessions-freshness-check
      - ruleRef: sessions-row-count-check
  owner:
    ownerRef: team-analytics
  tags:
    - pii
    - customer-data
  retentionPolicy:
    days: 730
  partitioning:
    type: time
    column: event_date
    granularity: day
```

**Status:**
```yaml
status:
  phase: Active  # Active, Deprecated, Archived
  observedGeneration: 3
  lastUpdated: "2026-01-18T02:45:00Z"
  rowCount: 1500000
  sizeBytes: 450000000
  schemaVersion: schema-sessions-v3
  qualityScore: 0.98
  lineage:
    upstreamCount: 3
    downstreamCount: 12
  providerStatus:
    syncedAt: "2026-01-18T02:45:10Z"
    externalId: "PROD.ANALYTICS.SESSIONS"
```

**Relationships:**
- `schemaVersionRef` → SchemaVersion (many-to-one)
- `owner` → Owner (many-to-one)
- `quality.rules` → DataQualityRule (many-to-many)
- `producedBy` → Pipeline (many-to-many, via lineage)
- `consumedBy` → Pipeline (many-to-many, via lineage)

---

#### Resource 3: **DataQualityRule**

**Purpose:** Declarative data quality check.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: DataQualityRule
metadata:
  name: sessions-freshness-check
  namespace: production
spec:
  description: "Ensure sessions table is updated daily"
  type: freshness
  severity: critical  # critical, high, medium, low
  schedule:
    cron: "0 3 * * *"
  target:
    datasetRef: prod.sessions
  config:
    maxAgeDays: 1
    timestampColumn: event_date
  notifications:
    onFailure:
      - type: pagerduty
        severity: critical
```

**Status:**
```yaml
status:
  phase: Active
  observedGeneration: 2
  lastCheck:
    timestamp: "2026-01-18T03:00:00Z"
    result: passed
    value: 0.5  # days old
    threshold: 1.0
  passRate: 0.99  # last 30 days
  consecutiveFailures: 0
```

---

#### Resource 4: **Owner**

**Purpose:** Represents a team or individual responsible for data assets.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Owner
metadata:
  name: team-analytics
  namespace: production
spec:
  displayName: "Analytics Team"
  type: team  # team, individual
  contacts:
    - type: email
      value: analytics@company.com
    - type: slack
      value: "#analytics"
  members:
    - userId: user-alice
      role: lead
    - userId: user-bob
      role: member
```

**Status:**
```yaml
status:
  ownedDatasets: 45
  ownedPipelines: 12
  qualityScore: 0.95  # avg of owned assets
```

---

#### Resource 5: **Policy**

**Purpose:** RBAC + ABAC policy definition.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Policy
metadata:
  name: analytics-read-policy
  namespace: production
spec:
  description: "Analytics team can read all datasets"
  subjects:
    - type: team
      ref: team-analytics
  resources:
    - type: Dataset
      namespaces: [production, staging]
  actions: [read, list]
  conditions:
    - attribute: dataset.tags
      operator: notContains
      value: pii-restricted
```

---

#### Resource 6: **Environment**

**Purpose:** Logical environment (dev, staging, prod) with configuration.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Environment
metadata:
  name: production
spec:
  displayName: "Production"
  type: production
  provider:
    type: snowflake
    config:
      account: xy12345.us-east-1
      warehouse: PROD_WH
      database: PROD
  secrets:
    - name: snowflake-prod-creds
      secretRef: secret-snowflake-prod
  quotas:
    maxPipelines: 100
    maxDatasets: 1000
```

---

#### Resource 7: **LineageEdge**

**Purpose:** Explicit lineage relationship between datasets/pipelines.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: LineageEdge
metadata:
  name: edge-events-to-sessions
  namespace: production
spec:
  source:
    type: Dataset
    ref: raw.events
  target:
    type: Dataset
    ref: prod.sessions
  relationship: transforms  # produces, consumes, transforms, derives_from
  via:
    type: Pipeline
    ref: customer-analytics-daily
  confidence: 1.0  # 0.0-1.0, for inferred lineage
  metadata:
    columns:
      - source: event_id
        target: session_id
        transformation: "HASH(event_id)"
```

---

#### Resource 8: **SchemaVersion**

**Purpose:** Immutable schema snapshot for datasets.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: SchemaVersion
metadata:
  name: schema-sessions-v3
  namespace: production
spec:
  datasetRef: prod.sessions
  version: 3
  schema:
    format: avro  # avro, json-schema, protobuf
    definition: |
      {
        "type": "record",
        "name": "Session",
        "fields": [
          {"name": "session_id", "type": "string"},
          {"name": "user_id", "type": "string"},
          {"name": "event_date", "type": "string"},
          {"name": "duration_seconds", "type": "int"}
        ]
      }
  changeType: breaking  # breaking, compatible, patch
  migratedFrom: schema-sessions-v2
```

---

#### Resource 9: **PipelineExecution**

**Purpose:** Immutable record of a pipeline run.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: PipelineExecution
metadata:
  name: exec-abc123
  namespace: production
spec:
  pipelineRef: customer-analytics-daily
  triggeredBy:
    type: schedule
    userId: system
  parameters:
    ds: "2026-01-18"
```

**Status:**
```yaml
status:
  phase: Succeeded
  startTime: "2026-01-18T02:00:00Z"
  endTime: "2026-01-18T02:45:00Z"
  duration: 2700s
  tasks:
    - name: extract-events
      status: Succeeded
      startTime: "2026-01-18T02:00:00Z"
      endTime: "2026-01-18T02:15:00Z"
    - name: transform-sessions
      status: Succeeded
      startTime: "2026-01-18T02:15:00Z"
      endTime: "2026-01-18T02:40:00Z"
  logs:
    url: "https://logs.fusion.io/exec-abc123"
  metrics:
    rowsProcessed: 5000000
    bytesProcessed: 1500000000
```

---

#### Resource 10: **Alert**

**Purpose:** Active alert triggered by quality rule or anomaly detection.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Alert
metadata:
  name: alert-sessions-freshness-fail
  namespace: production
spec:
  triggeredBy:
    type: DataQualityRule
    ref: sessions-freshness-check
  severity: critical
  message: "Sessions table is 2.5 days old, exceeds threshold of 1 day"
  affectedResources:
    - type: Dataset
      ref: prod.sessions
  createdAt: "2026-01-18T03:00:00Z"
```

**Status:**
```yaml
status:
  state: open  # open, acknowledged, resolved
  acknowledgedBy: user-alice
  acknowledgedAt: "2026-01-18T03:05:00Z"
  resolvedAt: null
  incidentRef: incident-001
```

---

#### Resource 11: **Incident**

**Purpose:** Aggregates related alerts into an incident for tracking.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Incident
metadata:
  name: incident-001
  namespace: production
spec:
  title: "Sessions pipeline failure cascade"
  severity: critical
  alerts:
    - alert-sessions-freshness-fail
    - alert-sessions-row-count-fail
  assignedTo: team-analytics
  createdAt: "2026-01-18T03:00:00Z"
```

**Status:**
```yaml
status:
  state: investigating  # open, investigating, resolved
  updates:
    - timestamp: "2026-01-18T03:10:00Z"
      author: user-alice
      message: "Investigating upstream data source delay"
  resolvedAt: null
  rootCause: null
```

---

#### Resource 12: **Secret**

**Purpose:** Reference to encrypted secret (credentials, API keys).

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Secret
metadata:
  name: secret-snowflake-prod
  namespace: production
spec:
  type: snowflake-credentials
  encryptedData:
    account: <encrypted>
    username: <encrypted>
    password: <encrypted>
  rotationPolicy:
    enabled: true
    intervalDays: 90
  lastRotated: "2025-12-01T00:00:00Z"
```

---

#### Resource 13: **Provider**

**Purpose:** Registered provider plugin with capabilities.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Provider
metadata:
  name: provider-dbt
spec:
  displayName: "dbt Core Provider"
  version: "1.0.0"
  type: orchestration
  capabilities:
    - pipeline-execution
    - lineage-extraction
    - schema-inference
  runtime: edge  # edge or control-plane
  certificationTier: tier1
  manifest:
    image: "fusion/provider-dbt:1.0.0"
    entrypoint: "/app/provider"
  healthCheck:
    endpoint: "/health"
    intervalSeconds: 30
```

---

#### Resource 14: **AIRequest**

**Purpose:** AI-generated change proposal awaiting approval.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: AIRequest
metadata:
  name: ai-req-001
  namespace: production
spec:
  prompt: "Add a data quality check for null values in user_id column"
  targetResource:
    type: Dataset
    ref: prod.sessions
  proposedChanges:
    - action: create
      resource:
        apiVersion: fusion.dcraft.io/v1
        kind: DataQualityRule
        metadata:
          name: sessions-null-check
        spec:
          type: null-check
          target:
            datasetRef: prod.sessions
          config:
            column: user_id
  requestedBy: user-alice
  createdAt: "2026-01-18T10:00:00Z"
```

**Status:**
```yaml
status:
  state: pending-review  # draft, pending-review, approved, rejected, applied
  reviewedBy: user-bob
  reviewedAt: null
  appliedAt: null
  rejectionReason: null
```

---

#### Resource 15: **Tenant**

**Purpose:** Multi-tenant isolation boundary.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Tenant
metadata:
  name: acme-corp
spec:
  displayName: "Acme Corporation"
  plan: enterprise
  quotas:
    maxUsers: 100
    maxPipelines: 500
    maxDatasets: 5000
    aiRequestsPerMonth: 1000
  features:
    aiOperator: true
    advancedLineage: true
  billingContact: billing@acme.com
```

---

#### Resource 16: **User**

**Purpose:** User identity and permissions.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: User
metadata:
  name: user-alice
  namespace: acme-corp
spec:
  email: alice@acme.com
  displayName: "Alice Smith"
  roles:
    - admin
    - data-engineer
  teams:
    - team-analytics
  ssoProvider: okta
  externalId: "okta-user-123"
```

---

#### Resource 17: **EdgeAgent**

**Purpose:** Represents an edge agent instance.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: EdgeAgent
metadata:
  name: edge-agent-prod-us-east
  namespace: production
spec:
  version: "1.0.0"
  environment: production
  location: us-east-1
  capabilities:
    - exec
    - observe
  certificate:
    expiresAt: "2026-02-18T00:00:00Z"
    autoRotate: true
```

**Status:**
```yaml
status:
  phase: Healthy
  lastHeartbeat: "2026-01-18T10:00:00Z"
  connectedSince: "2026-01-15T08:00:00Z"
  providers:
    - name: provider-dbt
      version: "1.0.0"
      status: Running
```

---

#### Resource 18: **Notification**

**Purpose:** Notification configuration for alerts.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Notification
metadata:
  name: slack-data-alerts
  namespace: production
spec:
  type: slack
  config:
    channel: "#data-alerts"
    webhookUrl: <secret-ref>
  filters:
    severity: [critical, high]
    resourceTypes: [Pipeline, DataQualityRule]
```

---

#### Resource 19: **Quota**

**Purpose:** Resource usage limits per tenant/namespace.

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Quota
metadata:
  name: production-quota
  namespace: production
spec:
  limits:
    pipelines: 100
    datasets: 1000
    executions-per-day: 500
    ai-requests-per-month: 100
```

**Status:**
```yaml
status:
  used:
    pipelines: 45
    datasets: 320
    executions-per-day: 234
    ai-requests-per-month: 12
```

---

#### Resource 20: **AuditLog**

**Purpose:** Immutable audit record (read-only resource).

**Spec:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: AuditLog
metadata:
  name: audit-001
spec:
  timestamp: "2026-01-18T10:00:00Z"
  actor:
    type: user
    ref: user-alice
  action: UPDATE
  resource:
    type: Pipeline
    ref: customer-analytics-daily
    version: 5
  changes:
    - field: spec.schedule.cron
      oldValue: "0 1 * * *"
      newValue: "0 2 * * *"
  ipAddress: "203.0.113.42"
  userAgent: "fusion-cli/1.0.0"
```

---

## 4. Provider/Plugin System (FROZEN)

### 4.1 Provider Interface Contract (gRPC)

**See `provider-contract-v0.proto` for full specification.**

**Core Methods:**
- `GetCapabilities()` - Returns provider capabilities manifest
- `ValidateConfig()` - Validates provider-specific configuration
- `TranslateResource()` - Translates Fusion resource → tool-native config
- `ExecuteCommand()` - Executes a command (create/update/delete resource in tool)
- `GetStatus()` - Retrieves current status of managed resource
- `ExtractLineage()` - Extracts lineage metadata from tool
- `HealthCheck()` - Provider health status

### 4.2 Capability Discovery Manifest Schema

```yaml
apiVersion: fusion.dcraft.io/v1
kind: ProviderManifest
metadata:
  name: dbt-provider
  version: "1.0.0"
spec:
  displayName: "dbt Core Provider"
  description: "Orchestrates dbt projects"
  vendor: "Fusion Labs"
  type: orchestration
  capabilities:
    - name: pipeline-execution
      version: "1.0"
      config:
        requiredFields: [project_dir, profiles_dir]
        optionalFields: [target, threads]
    - name: lineage-extraction
      version: "1.0"
      config:
        sources: [manifest.json]
    - name: schema-inference
      version: "1.0"
  resourceTypes:
    - apiVersion: fusion.dcraft.io/v1
      kind: Pipeline
      supportedActions: [create, update, delete, get]
  runtime:
    type: edge  # edge or control-plane
    isolation: container
    resources:
      cpu: "1000m"
      memory: "512Mi"
  dependencies:
    - name: dbt-core
      version: ">=1.7.0"
  secrets:
    - name: database-credentials
      type: snowflake-credentials
      required: true
```

### 4.3 Provider Versioning Strategy

**Semantic Versioning (SemVer):**
- Major: Breaking API changes
- Minor: New capabilities, backward-compatible
- Patch: Bug fixes

**Compatibility Matrix:**
```
Fusion Control Plane v1.0.x supports:
  - Provider API v1.x.x (stable)
  - Provider API v2.x.x (beta, opt-in)

Provider Deprecation Policy:
  - v1.x.x supported for 24 months after v2.0.0 release
  - 6-month deprecation notice before removal
```

### 4.4 Provider Execution Location

**Decision Matrix:**

| Provider Type | Execution Location | Justification |
|--------------|-------------------|---------------|
| Orchestration (dbt, Airflow) | **Edge** | Needs access to customer data, credentials, compute |
| Metadata (dbt lineage, Airflow DAG parser) | **Edge** | Reads tool-native metadata files |
| Cloud APIs (Snowflake, BigQuery) | **Control Plane** | Stateless API calls, no customer data access |
| AI/LLM | **Control Plane** | Centralized model management, cost control |
| Observability (log tailing) | **Edge** | Reads local logs, metrics |

**Edge Providers:**
- Run in sandboxed containers on Edge Agent
- Access secrets via Unix socket (no secrets in memory)
- Resource limits enforced (CPU, memory, network)
- Isolated network namespace (can't access other providers)

**Control Plane Providers:**
- Run as Kubernetes pods (if K8s available) or Docker containers
- Stateless, horizontally scalable
- No access to customer data, only metadata

### 4.5 Provider Certification Tiers

**Purpose:** Prevent provider sprawl, ensure quality.

| Tier | Description | Requirements | Support |
|------|-------------|--------------|---------|
| **Tier 1: Fusion-Native** | Built and maintained by Fusion Labs | - Full test coverage<br>- Security audit<br>- SLA commitment | 24/7 support |
| **Tier 2: Verified** | Third-party, verified by Fusion | - Security review<br>- Integration tests<br>- Documentation | Community + vendor |
| **Tier 3: Community** | Community-contributed | - Basic tests<br>- Manifest validation | Community only |
| **Tier 4: Experimental** | Unverified, use at own risk | - Manifest only | None |

**Enforcement:**
- Tenants can restrict allowed tiers via Policy
- Tier 1-2 providers listed in official registry
- Tier 3-4 require explicit opt-in

### 4.6 BYO / Fusion-Default / Hybrid Support

**Per-Layer Configuration:**

```yaml
# Example: Orchestration layer
apiVersion: fusion.dcraft.io/v1
kind: LayerConfig
metadata:
  name: orchestration-config
  namespace: production
spec:
  layer: orchestration
  mode: hybrid  # byo, fusion-default, hybrid
  byo:
    provider: airflow
    version: "2.8.0"
    endpoint: "https://airflow.acme.com"
    credentials:
      secretRef: airflow-creds
  fusionDefault:
    enabled: true
    provider: fusion-orchestrator
    fallbackToBYO: true  # if Fusion provider fails
  routing:
    # Route specific pipelines to specific providers
    rules:
      - match:
          labels:
            team: legacy
        provider: airflow
      - match:
          labels:
            team: modern
        provider: fusion-orchestrator
```

**Hybrid Mode Benefits:**
- Gradual migration from BYO → Fusion
- A/B testing of providers
- Fallback for reliability

---

## 5. Edge Agent Architecture (FROZEN)

### 5.1 Single Binary, Two Logical Components

**Decision:** Use a single binary with two logical components (NOT separate processes).

```
┌─────────────────────────────────────────┐
│         Edge Agent (Single Binary)      │
│                                         │
│  ┌────────────────────────────────────┐ │
│  │      Exec Component                │ │
│  │  - Receives commands from CP       │ │
│  │  - Executes provider operations    │ │
│  │  - Manages provider lifecycle      │ │
│  │  - Injects secrets securely        │ │
│  └────────────────────────────────────┘ │
│                                         │
│  ┌────────────────────────────────────┐ │
│  │      Observe Component             │ │
│  │  - Scrapes metadata (passive)      │ │
│  │  - Tails logs                      │ │
│  │  - Collects metrics                │ │
│  │  - Detects drift                   │ │
│  └────────────────────────────────────┘ │
│                                         │
│  ┌────────────────────────────────────┐ │
│  │      Shared Infrastructure         │ │
│  │  - mTLS client (outbound-only)     │ │
│  │  - Local state DB (SQLite)         │ │
│  │  - Provider runtime manager        │ │
│  │  - Certificate rotation            │ │
│  └────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

**Why Single Binary:**
- Simpler deployment (one process to manage)
- Shared mTLS connection (reduces overhead)
- Easier certificate rotation
- Lower resource footprint

**Component Isolation:**
- Separate goroutines/threads
- Separate gRPC streams to control plane
- Exec component can be disabled via config (observe-only mode)

### 5.2 Outbound-Only mTLS Channel

**Architecture:**
```
Edge Agent                          Control Plane
    │                                     │
    │  1. Establish mTLS connection      │
    │     (client cert + CA validation)  │
    ├────────────────────────────────────>│
    │                                     │
    │  2. Open bidirectional gRPC stream │
    │     (HTTP/2 over TLS)              │
    │<────────────────────────────────────>│
    │                                     │
    │  3. Agent sends heartbeat (30s)    │
    ├────────────────────────────────────>│
    │                                     │
    │  4. CP pushes commands via stream  │
    │<────────────────────────────────────┤
    │                                     │
    │  5. Agent executes, streams result │
    ├────────────────────────────────────>│
```

**Security Properties:**
- No inbound ports on edge
- Agent initiates all connections
- Mutual TLS (both sides authenticated)
- Certificate pinning (agent trusts specific CA)
- Connection multiplexing (single TCP connection)

**Firewall Requirements:**
- Outbound HTTPS (443) to control plane
- No inbound rules needed

### 5.3 Upgrade + Rollback Flow

**In-Place Upgrade (Zero-Downtime):**

```
1. Control plane pushes new agent binary to CDN
2. Agent polls for updates (every 1h, configurable)
3. Agent downloads new binary to /tmp/fusion-agent-new
4. Agent validates signature (Ed25519)
5. Agent spawns new process with new binary
6. New process connects to control plane
7. Control plane routes traffic to new process
8. Old process drains (waits for in-flight ops, max 5min)
9. Old process exits
10. New process renames itself to /usr/local/bin/fusion-agent
```

**Rollback:**
- Agent keeps last 3 versions in `/var/lib/fusion/agent-backups/`
- If new version fails health check 3x → auto-rollback
- Manual rollback via control plane command

**Air-Gapped Upgrade:**
- Admin downloads agent binary manually
- Places in `/var/lib/fusion/agent-updates/`
- Agent detects and upgrades on next poll

### 5.4 Trust Rotation Strategy

**Certificate Lifecycle:**
- Initial certificate issued during agent registration (valid 90 days)
- Auto-rotation at 30 days remaining
- Rotation uses existing mTLS connection (no downtime)

**Rotation Protocol:**
```
1. Agent requests new certificate via gRPC
2. Control plane generates new cert (same CA)
3. Agent receives new cert + private key
4. Agent writes to disk (encrypted at rest)
5. Agent opens new mTLS connection with new cert
6. Agent closes old connection after 5min grace period
```

**CA Rotation (Rare):**
- Control plane pushes new CA cert to agents
- Agents trust both old and new CA for 30 days
- After 30 days, old CA is removed

**Revocation:**
- Control plane maintains CRL (Certificate Revocation List)
- Agents check CRL on connect (cached 1h)
- Revoked agents rejected immediately

### 5.5 Offline / Air-Gapped Support

**Offline Mode:**
- Agent queues commands locally (SQLite)
- Max queue size: 10,000 commands or 1GB
- When connection restored, agent syncs queue
- Control plane reconciles state (may skip stale commands)

**Air-Gapped Mode:**
- Agent runs without control plane connection
- Reads declarative config from local files (`/etc/fusion/resources/`)
- Admin manually places YAML files
- Agent reconciles local state to match files
- No AI, no central lineage, no alerts (local-only)

**Air-Gapped Limitations:**
- No real-time updates
- No centralized monitoring
- No AI assistance
- Manual secret management

---

## 6. Data Model & Storage (FROZEN)

### 6.1 PostgreSQL Tenant Isolation: Row-Level Security (RLS)

**Decision:** Use `tenant_id` column + Row-Level Security (RLS).

**Why RLS over schema-per-tenant or db-per-tenant:**
- ✅ Scales to 10,000+ tenants (schema-per-tenant breaks at ~1000)
- ✅ Simpler backup/restore (single database)
- ✅ Cross-tenant queries possible (for admin/analytics)
- ✅ Lower operational overhead
- ❌ Slightly higher query complexity (tenant_id in every WHERE)

**Schema Structure:**
```sql
-- All resource tables have tenant_id
CREATE TABLE pipelines (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  namespace VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  spec JSONB NOT NULL,
  status JSONB,
  metadata JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,  -- soft delete
  UNIQUE(tenant_id, namespace, name)
);

-- RLS Policy
ALTER TABLE pipelines ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON pipelines
  USING (tenant_id = current_setting('app.current_tenant_id')::UUID);

-- Application sets tenant context per request
SET app.current_tenant_id = '<tenant-uuid>';
```

**Index Strategy:**
```sql
-- Primary lookup: by tenant + namespace + name
CREATE INDEX idx_pipelines_tenant_ns_name 
  ON pipelines(tenant_id, namespace, name);

-- Reconciliation: by tenant + updated_at
CREATE INDEX idx_pipelines_tenant_updated 
  ON pipelines(tenant_id, updated_at) 
  WHERE deleted_at IS NULL;

-- Sharding: by tenant_id hash
CREATE INDEX idx_pipelines_tenant_hash 
  ON pipelines((hashtext(tenant_id::TEXT)));
```

**Partitioning (Future):**
- Partition by `tenant_id` hash when single table exceeds 100M rows
- Use PostgreSQL declarative partitioning (v10+)

### 6.2 Lineage Graph Storage

**Decision:** PostgreSQL with `pg_graph` extension (or Apache AGE).

**Why PostgreSQL over Neo4j:**
- ✅ Simpler operations (one database)
- ✅ ACID transactions with resource updates
- ✅ Sufficient performance for <10M nodes
- ❌ Neo4j better for >10M nodes, complex graph algorithms

**Fallback:** If graph queries become bottleneck, migrate to Neo4j with dual-write.

**Schema:**
```sql
-- Nodes
CREATE TABLE lineage_nodes (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  type VARCHAR(50) NOT NULL,  -- Dataset, Pipeline, Column, Model
  resource_id UUID NOT NULL,  -- FK to datasets/pipelines table
  name VARCHAR(255) NOT NULL,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Edges
CREATE TABLE lineage_edges (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  source_node_id UUID NOT NULL REFERENCES lineage_nodes(id),
  target_node_id UUID NOT NULL REFERENCES lineage_edges(id),
  relationship VARCHAR(50) NOT NULL,  -- produces, consumes, transforms, derives_from
  via_pipeline_id UUID,  -- optional: which pipeline created this edge
  confidence FLOAT NOT NULL DEFAULT 1.0,  -- 0.0-1.0 for inferred lineage
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(tenant_id, source_node_id, target_node_id, relationship)
);

-- Indexes for fast graph queries
CREATE INDEX idx_lineage_edges_source ON lineage_edges(tenant_id, source_node_id);
CREATE INDEX idx_lineage_edges_target ON lineage_edges(tenant_id, target_node_id);
```

**Query Patterns (Optimized):**
```sql
-- Upstream dependencies (recursive CTE)
WITH RECURSIVE upstream AS (
  SELECT source_node_id, target_node_id, 1 as depth
  FROM lineage_edges
  WHERE tenant_id = $1 AND target_node_id = $2
  UNION ALL
  SELECT e.source_node_id, e.target_node_id, u.depth + 1
  FROM lineage_edges e
  JOIN upstream u ON e.target_node_id = u.source_node_id
  WHERE e.tenant_id = $1 AND u.depth < 10  -- max depth
)
SELECT DISTINCT n.* FROM upstream u
JOIN lineage_nodes n ON u.source_node_id = n.id;

-- Blast radius (downstream impact)
-- Similar query, reverse direction
```

### 6.3 Telemetry & Audit Store: ClickHouse

**Decision:** Use ClickHouse for high-volume time-series data.

**Why ClickHouse over PostgreSQL:**
- ✅ 10-100x better compression (columnar storage)
- ✅ 10x faster for analytical queries
- ✅ Built-in TTL and partitioning
- ❌ Not ACID (eventual consistency acceptable for logs)

**Schema:**
```sql
-- Execution logs
CREATE TABLE execution_logs (
  timestamp DateTime64(3),
  tenant_id UUID,
  execution_id UUID,
  pipeline_id UUID,
  level String,  -- INFO, WARN, ERROR
  message String,
  metadata String  -- JSON
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (tenant_id, execution_id, timestamp)
TTL timestamp + INTERVAL 90 DAY;  -- 90 days hot

-- Audit trail
CREATE TABLE audit_log (
  timestamp DateTime64(3),
  tenant_id UUID,
  actor_id UUID,
  action String,
  resource_type String,
  resource_id UUID,
  changes String,  -- JSON diff
  ip_address IPv4,
  user_agent String
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (tenant_id, timestamp)
TTL timestamp + INTERVAL 730 DAY;  -- 2 years

-- Metrics (time-series)
CREATE TABLE metrics (
  timestamp DateTime64(3),
  tenant_id UUID,
  metric_name String,
  value Float64,
  tags Map(String, String)
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (tenant_id, metric_name, timestamp)
TTL timestamp + INTERVAL 90 DAY;
```

**Retention Policy:**
- Hot storage (ClickHouse): 90 days
- Cold storage (S3): 2 years
- After 2 years: delete (compliance permitting)

### 6.4 Secrets Storage

**Decision:** HashiCorp Vault (or AWS Secrets Manager for SaaS).

**Why Vault:**
- ✅ Encryption at rest + in transit
- ✅ Dynamic secrets (auto-rotation)
- ✅ Audit logging
- ✅ Multi-cloud support

**Edge Agent Secret Injection:**
- Agent requests secret via mTLS to control plane
- Control plane fetches from Vault (with tenant-scoped policy)
- Control plane returns encrypted secret (ephemeral key)
- Agent decrypts in memory, passes to provider via Unix socket
- Secret never written to disk on edge

---

## 7. Lineage Subsystem (FROZEN PLAN)

### 7.1 Ingestion Sources (Priority Order)

1. **dbt Manifest (manifest.json)** - Tier 1
   - Explicit lineage from `ref()` and `source()`
   - Confidence: 1.0
   - Ingestion: Parse manifest, create LineageEdge resources

2. **OpenLineage Events** - Tier 1
   - Standard lineage format from Spark, Airflow, etc.
   - Confidence: 0.9-1.0 (depending on emitter quality)
   - Ingestion: Subscribe to OpenLineage event stream, translate to LineageEdge

3. **Runtime Tracing (Fusion-native)** - Tier 1
   - Instrument Fusion-native orchestrator
   - Capture actual reads/writes during execution
   - Confidence: 1.0
   - Ingestion: Real-time via event bus

4. **SQL Parsing (Limited)** - Tier 2
   - Parse SQL from Pipeline tasks
   - Extract table references (FROM, JOIN, INSERT INTO)
   - Confidence: 0.7 (can't handle dynamic SQL, CTEs are hard)
   - Ingestion: Best-effort, mark as inferred

**NOT Supported in V0:**
- Column-level lineage (V1 feature)
- Cross-system lineage (e.g., Airflow → dbt → Snowflake) - requires provider coordination

### 7.2 Conflict Resolution Rules

**Scenario 1: Same edge from multiple sources**
- **Rule:** Highest confidence wins
- **Example:** dbt manifest (1.0) vs SQL parsing (0.7) → use dbt

**Scenario 2: Contradictory edges**
- **Rule:** Most recent source wins, alert user
- **Example:** dbt says A→B, OpenLineage says A→C → use OpenLineage, create Alert

**Scenario 3: Partial lineage**
- **Rule:** Merge, mark as incomplete
- **Example:** dbt knows A→B, runtime tracing knows B→C → create both edges

### 7.3 Confidence Scoring

**Confidence Levels:**
- `1.0` - Explicit declaration (dbt ref, OpenLineage)
- `0.9` - Runtime observed (Fusion tracing)
- `0.7` - Inferred from SQL parsing
- `0.5` - Inferred from naming conventions (e.g., `raw_events` → `stg_events`)
- `0.3` - Guessed (not used in V0)

**UI Display:**
- Solid line: confidence ≥ 0.9
- Dashed line: confidence < 0.9
- Show confidence score on hover

### 7.4 Query Patterns (Must Be Fast)

**Performance SLA:**
- Upstream/downstream (depth ≤ 5): <100ms p95
- Blast radius (depth ≤ 3): <200ms p95
- Full graph export: <5s for 10k nodes

**Optimization:**
- Materialized views for common queries
- Cache frequent queries (Redis, 5min TTL)
- Pre-compute blast radius for critical datasets (nightly job)

---

## 8. AI Operator Architecture (FROZEN)

### 8.1 AI Placement in System

```
┌─────────────────────────────────────────────────────────────┐
│                    Control Plane                            │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              AI Operator Service                     │  │
│  │                                                      │  │
│  │  ┌────────────────┐  ┌────────────────┐            │  │
│  │  │  LLM Router    │  │ Prompt Manager │            │  │
│  │  │  - GPT-4       │  │ - Templates    │            │  │
│  │  │  - Claude      │  │ - Few-shot     │            │  │
│  │  │  - BYOM        │  │ - RAG context  │            │  │
│  │  └────────────────┘  └────────────────┘            │  │
│  │                                                      │  │
│  │  ┌────────────────┐  ┌────────────────┐            │  │
│  │  │ Policy Engine  │  │ Approval FSM   │            │  │
│  │  │ - Data access  │  │ - Draft        │            │  │
│  │  │ - Action scope │  │ - Review       │            │  │
│  │  │ - Budget       │  │ - Approved     │            │  │
│  │  └────────────────┘  └────────────────┘            │  │
│  │                                                      │  │
│  │  ┌────────────────────────────────────────────────┐ │  │
│  │  │         Context Builder                        │ │  │
│  │  │  - Fetches relevant resources (read-only)     │ │  │
│  │  │  - Lineage graph                              │ │  │
│  │  │  - Schema metadata                            │ │  │
│  │  │  - Quality history                            │ │  │
│  │  └────────────────────────────────────────────────┘ │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 8.2 Data Access Policy

**AI can access (read-only):**
- Resource specs (Pipeline, Dataset, etc.)
- Metadata (schema, lineage, quality scores)
- Execution history (success/failure, duration)
- Documentation (descriptions, comments)

**AI CANNOT access:**
- Actual data rows (PII, sensitive data)
- Secrets (credentials, API keys)
- Audit logs (who did what)
- Other tenants' data

**Enforcement:**
- AI service runs with restricted service account
- RLS policies prevent cross-tenant access
- Separate read-only database replica for AI queries

### 8.3 Approval Workflow State Machine

```
┌─────────┐
│  Draft  │  AI generates proposed changes
└────┬────┘
     │ User: "Submit for review"
     ▼
┌──────────────┐
│ PendingReview│  Assigned to reviewer (owner or admin)
└────┬─────────┘
     │
     ├─ Reviewer: "Approve" ──────────────┐
     │                                    ▼
     │                              ┌──────────┐
     │                              │ Approved │
     │                              └────┬─────┘
     │                                   │ Auto-apply or manual trigger
     │                                   ▼
     │                              ┌─────────┐
     │                              │ Applied │  Changes committed
     │                              └─────────┘
     │
     ├─ Reviewer: "Reject" ────────────┐
     │                                  ▼
     │                            ┌──────────┐
     │                            │ Rejected │  End state
     │                            └──────────┘
     │
     └─ User: "Cancel" ────────────────┐
                                       ▼
                                 ┌───────────┐
                                 │ Cancelled │  End state
                                 └───────────┘
```

**State Transitions (Enforced by API):**
- Draft → PendingReview (user action)
- PendingReview → Approved (reviewer action)
- PendingReview → Rejected (reviewer action)
- Draft → Cancelled (user action)
- Approved → Applied (system action or manual trigger)

**Timeout Policy:**
- PendingReview → Auto-reject after 7 days (configurable)
- Approved → Auto-apply after 24 hours (if auto-apply enabled)

### 8.4 Prompt Injection & Hallucination Mitigations

**Prompt Injection:**
1. **Input Sanitization:** Strip special tokens, limit length
2. **Prompt Template Isolation:** User input is clearly marked in prompt
   ```
   System: You are a data engineer assistant.
   Context: <schema>, <lineage>
   User Request: """
   {user_input}  # Clearly delimited
   """
   Task: Generate a data quality rule.
   ```
3. **Output Validation:** Parse AI output, reject if malformed or contains unexpected commands

**Hallucination:**
1. **Grounding in Context:** Provide relevant metadata (schema, existing rules)
2. **Constrained Output Format:** Require structured output (YAML, JSON)
3. **Validation Against Schema:** Reject if generated resource violates API schema
4. **Human Review:** All AI changes require approval (no auto-apply in V0)

**Adversarial Testing:**
- Red team tests for jailbreaks
- Automated tests with known malicious prompts
- Rate limiting per user (10 AI requests/hour)

### 8.5 BYOM (Bring Your Own Model) Routing

**Configuration:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: AIConfig
metadata:
  name: ai-config
  namespace: production
spec:
  defaultModel: gpt-4
  models:
    - name: gpt-4
      provider: openai
      endpoint: https://api.openai.com/v1
      apiKeyRef: secret-openai
    - name: claude-3
      provider: anthropic
      endpoint: https://api.anthropic.com/v1
      apiKeyRef: secret-anthropic
    - name: custom-llm
      provider: custom
      endpoint: https://llm.acme.com/v1
      apiKeyRef: secret-custom-llm
  routing:
    - match:
        promptType: code-generation
      model: gpt-4
    - match:
        promptType: data-quality
      model: claude-3
```

**Routing Logic:**
- Match by prompt type, user preference, or round-robin
- Fallback to default model if preferred model fails
- Track usage per model for billing

### 8.6 AI Budget Controls

**Per-Tenant Budget:**
```yaml
apiVersion: fusion.dcraft.io/v1
kind: Tenant
spec:
  quotas:
    aiRequestsPerMonth: 1000
    aiTokensPerMonth: 1000000  # input + output tokens
```

**Enforcement:**
- API rejects requests if quota exceeded
- Soft limit warning at 80%
- Admin can override (with audit log)

**Cost Tracking:**
- Log tokens used per request
- Aggregate by tenant, user, model
- Monthly billing report

---

## 9. Deployment Modes Matrix (FROZEN)

### 9.1 Deployment Modes

| Component | SaaS | BYOC | On-Prem | Docker Compose | Kubernetes |
|-----------|------|------|---------|----------------|------------|
| **API Server** | ✅ Fusion-hosted | ✅ Customer VPC | ✅ Customer DC | ✅ | ✅ |
| **Control Plane DB** | ✅ RDS | ✅ RDS/self-hosted | ✅ Self-hosted | ✅ PostgreSQL container | ✅ StatefulSet |
| **Controllers** | ✅ Fusion-hosted | ✅ Customer VPC | ✅ Customer DC | ✅ | ✅ |
| **Event Bus (NATS)** | ✅ Fusion-hosted | ✅ Customer VPC | ✅ Customer DC | ✅ | ✅ |
| **Lineage Graph** | ✅ Fusion-hosted | ✅ Customer VPC | ✅ Customer DC | ✅ PostgreSQL | ✅ StatefulSet |
| **Telemetry (ClickHouse)** | ✅ Fusion-hosted | ✅ Customer VPC | ✅ Customer DC | ⚠️ Optional | ✅ StatefulSet |
| **AI Operator** | ✅ Fusion-hosted | ⚠️ Optional | ⚠️ Optional | ⚠️ Optional | ⚠️ Optional |
| **Edge Agent** | ✅ Customer env | ✅ Customer VPC | ✅ Customer DC | ✅ | ✅ DaemonSet |
| **Providers** | ✅ Edge | ✅ Edge | ✅ Edge | ✅ Edge | ✅ Edge |

**Legend:**
- ✅ Fully supported
- ⚠️ Optional or degraded functionality
- ❌ Not supported

### 9.2 What Breaks Without Kubernetes

**Components that require Kubernetes:**
- None (by design)

**Components that benefit from Kubernetes:**
- Controller horizontal scaling (can use Docker Swarm or manual scaling)
- StatefulSets for databases (can use Docker volumes)
- Service discovery (can use DNS or load balancer)

**Fallback for non-K8s deployments:**
- Use Docker Compose with health checks
- Use external load balancer (e.g., HAProxy, Nginx)
- Use external service registry (e.g., Consul)

### 9.3 Reference Deployment Topologies

#### 9.3.1 SaaS (Fusion-Hosted)

```
┌─────────────────────────────────────────────────────────────┐
│                  Fusion Cloud (AWS us-east-1)               │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  EKS Cluster (Control Plane)                        │  │
│  │  - API Server (ALB)                                 │  │
│  │  - Controllers (Deployment, HPA)                    │  │
│  │  - NATS (StatefulSet)                               │  │
│  │  - AI Operator (Deployment)                         │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  RDS PostgreSQL (Multi-AZ)                          │  │
│  │  - Control plane DB                                 │  │
│  │  - Lineage graph                                    │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  ClickHouse (EC2, self-managed)                     │  │
│  │  - Telemetry & audit logs                           │  │
│  └─────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                          │
                          │ mTLS (443)
                          │
┌─────────────────────────────────────────────────────────────┐
│              Customer Environment (Any Cloud/On-Prem)       │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Edge Agent (VM or container)                       │  │
│  │  - Connects to Fusion Cloud                         │  │
│  │  - Runs providers locally                           │  │
│  └─────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

#### 9.3.2 BYOC (Bring Your Own Cloud)

```
┌─────────────────────────────────────────────────────────────┐
│              Customer AWS Account (us-west-2)               │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  EKS Cluster (Control Plane)                        │  │
│  │  - Fusion-provided Helm chart                       │  │
│  │  - Customer manages infrastructure                  │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  RDS PostgreSQL (Customer-managed)                  │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Edge Agent (Same VPC)                              │  │
│  │  - Private networking, no internet egress           │  │
│  └─────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

#### 9.3.3 On-Prem (Air-Gapped)

```
┌─────────────────────────────────────────────────────────────┐
│              Customer Data Center (Air-Gapped)              │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Kubernetes Cluster (or Docker Compose)             │  │
│  │  - All components self-hosted                       │  │
│  │  - No external dependencies                         │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  PostgreSQL (VM or container)                       │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Edge Agent (Same network)                          │  │
│  │  - Reads local config files                         │  │
│  │  - No control plane connection                      │  │
│  └─────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

#### 9.3.4 Docker Compose (Dev/Small Prod)

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15
    volumes:
      - pgdata:/var/lib/postgresql/data
  
  nats:
    image: nats:2.10
  
  api-server:
    image: fusion/api-server:1.0.0
    depends_on: [postgres, nats]
    ports:
      - "8080:8080"
  
  controller:
    image: fusion/controller:1.0.0
    depends_on: [postgres, nats]
    deploy:
      replicas: 2
  
  edge-agent:
    image: fusion/edge-agent:1.0.0
    depends_on: [api-server]
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock  # for provider execution
```

---

## 10. V0 vs V1 Scope Boundaries

### 10.1 V0: Minimal Viable Control Plane

**Goal:** Prove control plane architecture is real, not a dashboard.

**In Scope:**
1. **Core Resources (8):**
   - Pipeline, Dataset, DataQualityRule, Owner, Policy, Environment, PipelineExecution, EdgeAgent

2. **Control Plane:**
   - API Server (CRUD, RBAC)
   - PostgreSQL (RLS-based multi-tenancy)
   - Controllers (Pipeline, Dataset, Quality)
   - Event bus (NATS)

3. **Edge Agent:**
   - Outbound-only mTLS
   - Provider runtime (Docker-based)
   - Exec + Observe components

4. **Providers (3):**
   - dbt (Tier 1)
   - Snowflake (Tier 1)
   - Generic SQL (Tier 2)

5. **Lineage:**
   - dbt manifest ingestion
   - Basic graph queries (upstream/downstream)

6. **AI Operator:**
   - Generate data quality rules
   - Approval workflow (Draft → Review → Approved → Applied)

7. **Deployment:**
   - Docker Compose (dev)
   - Kubernetes (prod)

**Out of Scope (V1):**
- Remaining 12 resources
- Advanced lineage (OpenLineage, SQL parsing)
- Column-level lineage
- AI for pipeline generation
- Incident management
- Advanced RBAC (ABAC)
- Multi-region
- ClickHouse (use PostgreSQL for logs in V0)

### 10.2 V1: Full 16-Layer Platform

**Added in V1 (No Architecture Changes):**

1. **Remaining Resources:**
   - LineageEdge, SchemaVersion, Alert, Incident, Secret, Provider, AIRequest, Tenant, User, Notification, Quota, AuditLog

2. **16 Layers Support:**
   - Ingestion (Airbyte, Fivetran)
   - Storage (S3, GCS, ADLS)
   - Compute (Spark, Flink)
   - Orchestration (Airflow, Dagster, Prefect)
   - Transformation (dbt, Dataform)
   - Quality (Great Expectations, Soda)
   - Catalog (DataHub, Atlan)
   - Governance (Collibra, Alation)
   - Observability (Monte Carlo, Datadog)
   - BI (Tableau, Looker)
   - ML (MLflow, Kubeflow)
   - Reverse ETL (Census, Hightouch)
   - Streaming (Kafka, Pulsar)
   - Query (Trino, Presto)
   - Notebooks (Jupyter, Hex)
   - Metrics (dbt Metrics, Metriql)

3. **Advanced Lineage:**
   - OpenLineage ingestion
   - SQL parsing
   - Column-level lineage
   - Cross-system lineage

4. **Advanced AI:**
   - Pipeline generation
   - Anomaly detection
   - Root cause analysis

5. **Enterprise Features:**
   - ABAC policies
   - Incident management
   - Advanced alerting
   - SLA tracking

**Architecture Unchanged:**
- Same control plane design
- Same reconciliation model
- Same provider interface
- Same edge agent

---

## 11. Doubts / Ambiguities

### 11.1 Doubt: Should providers run in control plane or edge?

**Current Decision:** Edge by default, control plane for stateless API providers.

**Ambiguity:** Some providers (e.g., dbt lineage parser) could run in either location.

**Resolution Needed:**
- Define clear criteria: "If provider needs customer data access → edge, else → control plane"
- Allow override via provider manifest (`runtime: edge` or `runtime: control-plane`)

**Risk if wrong:** Security (control plane accessing customer data) or performance (edge agent overload).

---

### 11.2 Doubt: How to handle schema evolution in resources?

**Current Decision:** Immutable revisions for spec, mutable status.

**Ambiguity:** What if we need to add a required field to spec in V2?

**Resolution Needed:**
- Use API versioning (`v1`, `v2`)
- Support multiple versions simultaneously (conversion webhooks)
- Deprecation policy: v1 supported for 12 months after v2 release

**Risk if wrong:** Breaking changes force all users to migrate simultaneously.

---

### 11.3 Doubt: Should lineage graph be in PostgreSQL or Neo4j?

**Current Decision:** PostgreSQL with pg_graph extension.

**Ambiguity:** At what scale does PostgreSQL become insufficient?

**Resolution Needed:**
- Benchmark at 1M, 10M, 100M nodes
- Define migration path to Neo4j if needed
- Consider dual-write strategy (PostgreSQL + Neo4j) for gradual migration

**Risk if wrong:** Performance bottleneck at scale, expensive migration.

---

### 11.4 Doubt: How to handle multi-region deployments?

**Current Decision:** Not in V0, defer to V1.

**Ambiguity:** Should control plane be multi-region active-active or active-passive?

**Resolution Needed:**
- V1 decision: Active-passive with regional read replicas
- Use PostgreSQL logical replication for cross-region sync
- Edge agents connect to nearest region

**Risk if wrong:** High latency for global customers, complex failover.

---

### 11.5 Doubt: Should AI Operator be mandatory or optional?

**Current Decision:** Optional per tenant.

**Ambiguity:** If optional, how to ensure architecture doesn't bifurcate?

**Resolution Needed:**
- AI Operator is a separate service, can be disabled
- Core control plane has no dependency on AI
- AI uses same resource model (AIRequest resource)

**Risk if wrong:** Two codebases to maintain (AI vs non-AI).

---

### 11.6 Doubt: How to handle provider versioning conflicts?

**Current Decision:** SemVer, compatibility matrix.

**Ambiguity:** What if two pipelines need different versions of same provider?

**Resolution Needed:**
- Allow multiple provider versions to coexist
- Edge agent runs providers in isolated containers
- Pipeline spec includes `provider.version` field

**Risk if wrong:** Dependency hell, edge agent resource exhaustion.

---

### 11.7 Doubt: Should we support custom resources (CRDs)?

**Current Decision:** Not in V0.

**Ambiguity:** Will users need to extend resource model?

**Resolution Needed:**
- V1 feature: Allow custom resources via CRD-like mechanism
- Require schema validation (JSON Schema)
- Custom controllers run in edge or control plane (user choice)

**Risk if wrong:** Platform becomes too rigid, users hack around it.

---

### 11.8 Doubt: How to handle clock skew in distributed system?

**Current Decision:** Use server-side timestamps, periodic resync.

**Ambiguity:** What if edge agent clock is off by hours?

**Resolution Needed:**
- Edge agent syncs time with control plane on connect (NTP-like)
- Use logical clocks (Lamport timestamps) for ordering
- Reject events with timestamp >1h in future or past

**Risk if wrong:** Incorrect event ordering, missed reconciliations.

---

### 11.9 Doubt: Should we support GitOps-style declarative config?

**Current Decision:** Not in V0, but architecture supports it.

**Ambiguity:** How to reconcile Git state vs API state?

**Resolution Needed:**
- V1 feature: GitOps controller watches Git repo
- Converts YAML files → API calls
- Conflict resolution: Git is source of truth (with approval workflow)

**Risk if wrong:** Users expect GitOps, disappointed if not available.

---

### 11.10 Doubt: How to handle provider failures during reconciliation?

**Current Decision:** Retry with exponential backoff, max 3 attempts.

**Ambiguity:** What if provider is down for hours?

**Resolution Needed:**
- Mark resource as `Degraded` after max retries
- Create Alert resource
- Retry on next periodic resync (5min)
- Allow manual retry via API

**Risk if wrong:** Infinite retry loops, resource stuck in bad state.

---

## 12. Final Verdict

### 12.1 Architecture Grade: **A-**

**Strengths:**
- ✅ True control plane architecture (not a dashboard)
- ✅ Clean separation of concerns (control plane vs edge)
- ✅ Declarative resource model with reconciliation
- ✅ Multi-tenancy via RLS (scales to 10k+ tenants)
- ✅ Provider system supports 100+ tools
- ✅ Deployment-agnostic (SaaS, BYOC, on-prem, Docker, K8s)
- ✅ Security-first (mTLS, RLS, secrets management)
- ✅ AI governance (approval workflow, policy enforcement)

**Weaknesses:**
- ⚠️ Lineage graph in PostgreSQL may not scale beyond 10M nodes (mitigation: plan migration to Neo4j)
- ⚠️ Provider versioning conflicts not fully solved (mitigation: allow multiple versions)
- ⚠️ Multi-region not designed yet (acceptable for V0)

**Why not A+:**
- Some ambiguities remain (see section 11)
- Lineage scalability unproven at 100M+ nodes
- No production validation yet (expected for V0)

---

### 12.2 Three Non-Negotiable Risks

#### Risk 1: **Lineage Graph Scalability**

**Problem:** PostgreSQL with pg_graph may not handle 100M+ nodes efficiently.

**Impact:** Slow queries, poor UX, customer churn.

**Mitigation:**
- Benchmark at 10M nodes before V1 launch
- Design migration path to Neo4j (dual-write strategy)
- Implement query result caching (Redis)
- Pre-compute common queries (nightly job)

**Decision Point:** If p95 query latency >1s at 10M nodes → migrate to Neo4j.

---

#### Risk 2: **Provider Sprawl**

**Problem:** Without governance, 100+ providers → maintenance nightmare.

**Impact:** Security vulnerabilities, poor quality, support burden.

**Mitigation:**
- Enforce certification tiers (Tier 1-4)
- Require security review for Tier 2
- Deprecate unused providers (if <10 active users after 6 months)
- Limit Tier 3-4 providers to 10% of total

**Decision Point:** If >50% of providers are Tier 3-4 → tighten certification requirements.

---

#### Risk 3: **AI Hallucination Causes Data Corruption**

**Problem:** AI generates incorrect SQL, deletes production data.

**Impact:** Customer data loss, reputational damage, lawsuits.

**Mitigation:**
- MANDATORY human approval for all AI changes (no auto-apply in V0)
- Dry-run mode (simulate changes, show impact)
- Rollback capability (store last 10 revisions)
- AI cannot access actual data, only metadata
- Rate limiting (10 AI requests/hour per user)

**Decision Point:** If >1% of AI-generated changes cause incidents → disable AI auto-apply permanently.

---

### 12.3 Five Acceptance Tests (Prove This Is a Real Control Plane)

#### Test 1: **Declarative Reconciliation**

**Scenario:** User creates a Pipeline resource via API. Controller reconciles it to dbt project.

**Steps:**
1. POST `/api/v1/pipelines` with Pipeline spec
2. Pipeline controller detects new resource (via event bus)
3. Controller calls dbt provider to create dbt project
4. Provider returns success
5. Controller updates Pipeline status to `Running`
6. User manually edits dbt project (introduces drift)
7. Controller detects drift on next resync (5min)
8. Controller reverts dbt project to match Pipeline spec

**Pass Criteria:**
- Pipeline status reflects actual state within 10s
- Drift detected and corrected within 6min
- No manual intervention required

---

#### Test 2: **Multi-Tenant Isolation**

**Scenario:** Two tenants (Acme, Globex) create pipelines with same name.

**Steps:**
1. Tenant Acme: POST `/api/v1/pipelines` with name `daily-etl`
2. Tenant Globex: POST `/api/v1/pipelines` with name `daily-etl`
3. Both succeed (different tenant_id)
4. Acme user: GET `/api/v1/pipelines` → sees only Acme's pipeline
5. Globex user: GET `/api/v1/pipelines` → sees only Globex's pipeline
6. Attempt to access other tenant's pipeline → 403 Forbidden

**Pass Criteria:**
- No cross-tenant data leakage
- RLS policies enforced at database level
- API returns 403 for unauthorized access

---

#### Test 3: **Provider Pluggability**

**Scenario:** Replace dbt provider with Dagster provider without changing Pipeline resource.

**Steps:**
1. Create Pipeline with `provider.type: dbt`
2. Pipeline runs successfully
3. Update Pipeline: `provider.type: dagster`
4. Controller detects spec change
5. Controller calls dagster provider (not dbt)
6. Dagster provider translates Pipeline → Dagster DAG
7. Pipeline runs successfully with Dagster

**Pass Criteria:**
- No code changes to control plane
- Provider swap is transparent to user
- Pipeline execution succeeds with both providers

---

#### Test 4: **Edge Agent Resilience**

**Scenario:** Edge agent loses connection to control plane, queues commands, reconnects.

**Steps:**
1. Edge agent connected, Pipeline running
2. Simulate network partition (block outbound 443)
3. Control plane sends command to update Pipeline
4. Edge agent queues command locally (SQLite)
5. Restore network after 10 minutes
6. Edge agent reconnects, syncs queue
7. Pipeline update applied

**Pass Criteria:**
- No commands lost during partition
- Queue syncs within 30s of reconnection
- Control plane reconciles state correctly

---

#### Test 5: **AI Approval Workflow**

**Scenario:** AI generates a data quality rule, requires human approval.

**Steps:**
1. User: "Add a null check for user_id column"
2. AI generates AIRequest resource with proposed DataQualityRule
3. AIRequest status: `PendingReview`
4. Reviewer: GET `/api/v1/airequests/{id}` → sees proposed changes
5. Reviewer: PATCH `/api/v1/airequests/{id}` → `Approved`
6. System creates DataQualityRule resource
7. Quality controller reconciles, runs check

**Pass Criteria:**
- AI cannot apply changes without approval
- Approval workflow enforced by API
- Audit log records approval action

---

## Conclusion

This architecture is **production-ready for V0** with the following caveats:
1. Lineage scalability must be validated at 10M+ nodes
2. Provider certification process must be enforced
3. AI approval workflow is mandatory (no auto-apply)

**Next Steps:**
1. Implement V0 scope (8 resources, 3 providers)
2. Run acceptance tests
3. Benchmark lineage at 10M nodes
4. Launch private beta with 10 design partners
5. Iterate based on feedback
6. Plan V1 (16 layers, advanced features)

**Estimated V0 Timeline:**
- Architecture freeze: Week 1 (this document)
- Core control plane: Weeks 2-8
- Edge agent + providers: Weeks 9-12
- AI Operator: Weeks 13-16
- Testing + hardening: Weeks 17-20
- Private beta: Week 21

**Total: 5 months to V0 launch.**

---

**Document Status:** FROZEN  
**Next Review:** After V0 private beta (Week 21)  
**Owner:** Principal Architect  
**Approvers:** CTO, VP Engineering, Head of Product