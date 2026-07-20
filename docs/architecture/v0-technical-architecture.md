# DCraft Fusion V0: Complete Technical Architecture

## Table of Contents
1. [Technology Stack Overview](#technology-stack-overview)
2. [Programming Language Choice](#programming-language-choice)
3. [Microservices Architecture](#microservices-architecture)
4. [Communication Patterns](#communication-patterns)
5. [Data Storage Architecture](#data-storage-architecture)
6. [AI Integration Architecture](#ai-integration-architecture)
7. [Deployment Strategy](#deployment-strategy)
8. [Development Workflow](#development-workflow)
9. [Security Architecture](#security-architecture)
10. [Scalability & Performance](#scalability--performance)
11. [Monitoring & Observability](#monitoring--observability)
12. [Cost Optimization](#cost-optimization)

---

# Technology Stack Overview

## Core Technology Decisions

```
┌─────────────────────────────────────────────────────────────┐
│                    DCRAFT FUSION V0 STACK                    │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Backend:           Go 1.21+                                 │
│  API Gateway:       Go (Gin framework)                       │
│  Microservices:     Go                                       │
│  Frontend:          React 18 + TypeScript                    │
│  CLI:               Go (Cobra framework)                     │
│                                                              │
│  Communication:                                              │
│    - Internal:      gRPC (Protocol Buffers)                  │
│    - External:      REST + GraphQL                           │
│                                                              │
│  Storage:                                                    │
│    - Metadata:      PostgreSQL 15                            │
│    - Lineage:       Neo4j 5.x                                │
│    - Cache:         Redis 7.x                                │
│    - Time-series:   TimescaleDB (PostgreSQL extension)       │
│                                                              │
│  Message Queue:     NATS / Apache Kafka                      │
│  Container:         Docker                                   │
│  Orchestration:     Kubernetes                               │
│  Deployment:        Helm + ArgoCD (GitOps)                   │
│  CI/CD:             GitHub Actions                           │
│                                                              │
│  AI/LLM:            OpenAI API / Anthropic Claude API        │
│  Observability:     Prometheus + Grafana + Jaeger            │
│  Logging:           ELK Stack (Elasticsearch, Logstash, Kibana)│
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

# Programming Language Choice

## Primary Language: Go (Golang)

### Why Go for DCraft Fusion Control Plane?

**1. Perfect Fit for Control Plane Systems**
- Kubernetes, Docker, Terraform, and most CNCF projects are written in Go
- Proven track record for building distributed systems
- Designed specifically for cloud-native infrastructure

**2. Performance & Concurrency**
- Excellent concurrency model with goroutines
- Can handle tens of thousands of simultaneous connections
- Minimal runtime footprint
- Fast compilation times (critical for CI/CD)
- Single binary distribution (easy deployment)

**3. Ecosystem & Libraries**
- Rich ecosystem for cloud-native development
- Excellent Kubernetes client libraries (client-go)
- Strong gRPC support
- Mature database drivers (PostgreSQL, Neo4j, Redis)
- Extensive testing frameworks

**4. Developer Productivity**
- Simple, readable syntax
- Fast onboarding for new developers
- Strong standard library
- Built-in tooling (go fmt, go test, go mod)
- Excellent documentation

**5. Production-Ready**
- Battle-tested in production at scale (Google, Uber, Dropbox)
- Stable language with backward compatibility guarantees
- Strong community support
- Excellent debugging and profiling tools

### When NOT to Use Go

**Performance-Critical Hotspots:**
If profiling reveals specific components need 30%+ better performance, consider Rust for those specific microservices. However, start with Go and optimize only when necessary.

### Language Breakdown by Component

```
┌─────────────────────────────────────────────────────────────┐
│  Component                    │  Language  │  Justification  │
├─────────────────────────────────────────────────────────────┤
│  API Gateway                  │  Go        │  Performance    │
│  Reconciliation Engine        │  Go        │  Concurrency    │
│  Provider System              │  Go        │  Ecosystem      │
│  Lineage Engine               │  Go        │  Graph queries  │
│  Incident Manager             │  Go        │  Reliability    │
│  AI Orchestrator              │  Go        │  API integration│
│  CLI                          │  Go        │  Single binary  │
│  Web UI (Frontend)            │  TypeScript│  Type safety    │
│  Infrastructure (IaC)         │  HCL       │  Terraform      │
└─────────────────────────────────────────────────────────────┘
```

---

# Microservices Architecture

## Service Decomposition Strategy

### Core Principle: Domain-Driven Design (DDD)

Decompose by **business capability**, not technical layers.

### V0 Microservices (8 Core Services)

```
┌─────────────────────────────────────────────────────────────┐
│                      API GATEWAY                             │
│  - Authentication/Authorization                              │
│  - Rate Limiting                                             │
│  - Request Routing                                           │
│  - Protocol Translation (REST ↔ gRPC)                        │
└──────────────────┬───────────────────────────────────────────┘
                   │
         ┌─────────┴─────────┬─────────────────┬───────────────┐
         │                   │                 │               │
         ▼                   ▼                 ▼               ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ RESOURCE        │  │ PROVIDER        │  │ LINEAGE         │  │ RECONCILIATION  │
│ SERVICE         │  │ SERVICE         │  │ SERVICE         │  │ SERVICE         │
│                 │  │                 │  │                 │  │                 │
│ - CRUD ops      │  │ - Provider mgmt │  │ - Graph queries │  │ - Drift detect  │
│ - Validation    │  │ - Metadata sync │  │ - Impact analysis│ │ - State compare │
│ - Versioning    │  │ - Provider API  │  │ - Lineage viz   │  │ - Reconcile     │
└─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────────┘
         │                   │                 │               │
         ▼                   ▼                 ▼               ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ INCIDENT        │  │ AI OPERATOR     │  │ NOTIFICATION    │  │ AUDIT           │
│ SERVICE         │  │ SERVICE         │  │ SERVICE         │  │ SERVICE         │
│                 │  │                 │  │                 │  │                 │
│ - Detection     │  │ - LLM integration│ │ - Slack/Email   │  │ - Event logging │
│ - Root cause    │  │ - Prompt mgmt   │  │ - PagerDuty     │  │ - Compliance    │
│ - Resolution    │  │ - Blueprint gen │  │ - Webhooks      │  │ - Audit trail   │
└─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────────┘
         │                   │                 │               │
         └───────────────────┴─────────────────┴───────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │  DATA LAYER     │
                    │                 │
                    │ - PostgreSQL    │
                    │ - Neo4j         │
                    │ - Redis         │
                    │ - TimescaleDB   │
                    └─────────────────┘
```

---

## Detailed Service Specifications

### 1. API Gateway Service

**Responsibility:** Single entry point for all external requests.

**Technology:**
- Go + Gin framework
- Kong or Traefik (alternative: custom Go gateway)

**Features:**
- Authentication (JWT, OAuth2, API Keys)
- Authorization (RBAC)
- Rate limiting (per user, per API key)
- Request routing to internal services
- Protocol translation (REST → gRPC)
- Response caching
- Request/response logging
- Metrics collection

**Endpoints:**
```
REST API:
  GET    /api/v1/datasets
  POST   /api/v1/datasets
  GET    /api/v1/datasets/{id}
  PUT    /api/v1/datasets/{id}
  DELETE /api/v1/datasets/{id}
  
  GET    /api/v1/pipelines
  POST   /api/v1/pipelines
  GET    /api/v1/pipelines/{id}
  
  GET    /api/v1/lineage/downstream/{id}
  GET    /api/v1/lineage/upstream/{id}
  
  GET    /api/v1/incidents
  POST   /api/v1/incidents
  
  POST   /api/v1/ai/ask
  POST   /api/v1/ai/generate-blueprint

GraphQL API:
  POST   /graphql
```

**Deployment:**
- 2-3 replicas (HA)
- Horizontal Pod Autoscaler (HPA): 2-10 replicas
- Resource limits: 1 CPU, 1Gi memory per replica

---

### 2. Resource Service

**Responsibility:** Manage all control-plane resources (Datasets, Pipelines, DataProducts, Policies).

**Technology:**
- Go
- PostgreSQL (primary storage)
- Redis (caching)

**Features:**
- CRUD operations on resources
- Schema validation
- Resource versioning (track changes)
- Resource relationships (foreign keys)
- Search and filtering
- Pagination

**Internal API (gRPC):**
```protobuf
service ResourceService {
  rpc CreateDataset(CreateDatasetRequest) returns (Dataset);
  rpc GetDataset(GetDatasetRequest) returns (Dataset);
  rpc UpdateDataset(UpdateDatasetRequest) returns (Dataset);
  rpc DeleteDataset(DeleteDatasetRequest) returns (Empty);
  rpc ListDatasets(ListDatasetsRequest) returns (ListDatasetsResponse);
  
  rpc CreatePipeline(CreatePipelineRequest) returns (Pipeline);
  rpc GetPipeline(GetPipelineRequest) returns (Pipeline);
  // ... similar for other resources
}
```

**Database Schema:**
```sql
-- Datasets table
CREATE TABLE datasets (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  owner VARCHAR(255),
  schema JSONB,
  sla JSONB,
  contracts JSONB,
  metadata JSONB,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  version INTEGER
);

-- Pipelines table
CREATE TABLE pipelines (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(50),
  cadence VARCHAR(100),
  owner VARCHAR(255),
  input_datasets UUID[],
  output_datasets UUID[],
  execution_engine VARCHAR(100),
  sla JSONB,
  metadata JSONB,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  version INTEGER
);

-- DataProducts table
CREATE TABLE data_products (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  owner VARCHAR(255),
  domain VARCHAR(100),
  datasets UUID[],
  pipelines UUID[],
  consumers JSONB,
  sla JSONB,
  metadata JSONB,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Policies table
CREATE TABLE policies (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(50),
  scope VARCHAR(100),
  rules JSONB,
  enforcement VARCHAR(50),
  metadata JSONB,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Resource versions (audit trail)
CREATE TABLE resource_versions (
  id UUID PRIMARY KEY,
  resource_type VARCHAR(50),
  resource_id UUID,
  version INTEGER,
  changes JSONB,
  changed_by VARCHAR(255),
  changed_at TIMESTAMP
);
```

**Deployment:**
- 2-3 replicas
- HPA: 2-5 replicas
- Resource limits: 500m CPU, 512Mi memory per replica

---

### 3. Provider Service

**Responsibility:** Manage provider lifecycle and metadata extraction.

**Technology:**
- Go
- PostgreSQL (provider configs)
- Redis (provider state cache)

**Features:**
- Provider registration and discovery
- Provider installation and configuration
- Provider execution (scheduled or triggered)
- Provider health monitoring
- Provider error handling and retry logic
- Provider metrics collection

**Internal API (gRPC):**
```protobuf
service ProviderService {
  rpc RegisterProvider(RegisterProviderRequest) returns (Provider);
  rpc ListProviders(ListProvidersRequest) returns (ListProvidersResponse);
  rpc GetProvider(GetProviderRequest) returns (Provider);
  rpc ConfigureProvider(ConfigureProviderRequest) returns (Provider);
  rpc ExecuteProvider(ExecuteProviderRequest) returns (ProviderExecution);
  rpc GetProviderStatus(GetProviderStatusRequest) returns (ProviderStatus);
}
```

**Provider Interface (gRPC Contract):**
```protobuf
service Provider {
  // Capability discovery
  rpc GetCapabilities(Empty) returns (Capabilities);
  
  // Metadata extraction
  rpc ExtractMetadata(ExtractMetadataRequest) returns (ExtractMetadataResponse);
  
  // Lineage extraction
  rpc ExtractLineage(ExtractLineageRequest) returns (ExtractLineageResponse);
  
  // Health check
  rpc HealthCheck(Empty) returns (HealthStatus);
}

message Capabilities {
  string provider_name = 1;
  string provider_version = 2;
  repeated string supported_operations = 3;
  map<string, string> configuration_schema = 4;
}

message ExtractMetadataRequest {
  map<string, string> config = 1;
  repeated string resource_types = 2;
}

message ExtractMetadataResponse {
  repeated Dataset datasets = 1;
  repeated Pipeline pipelines = 2;
  map<string, string> metadata = 3;
}
```

**V0 Providers to Implement:**

**1. dbt Provider**
```go
type DbtProvider struct {
  manifestPath string
  runResultsPath string
}

func (p *DbtProvider) ExtractMetadata() (*Metadata, error) {
  // Parse manifest.json
  // Extract models, sources, tests
  // Return as control-plane resources
}

func (p *DbtProvider) ExtractLineage() (*Lineage, error) {
  // Parse manifest.json
  // Extract model dependencies
  // Return as lineage edges
}
```

**2. Warehouse Provider (Snowflake/BigQuery/Redshift/Postgres)**
```go
type WarehouseProvider struct {
  connectionString string
  warehouse string
}

func (p *WarehouseProvider) ExtractMetadata() (*Metadata, error) {
  // Query information_schema
  // Extract tables, columns, schemas
  // Query query_history (if available)
  // Return as control-plane resources
}

func (p *WarehouseProvider) ExtractLineage() (*Lineage, error) {
  // Parse query logs
  // Extract table dependencies from SQL
  // Return as lineage edges
}
```

**3. Generic SQL Provider**
```go
type GenericSQLProvider struct {
  connectionString string
}

func (p *GenericSQLProvider) ExtractMetadata() (*Metadata, error) {
  // Query information_schema (standard SQL)
  // Extract tables, columns, schemas
  // Return as control-plane resources
}
```

**Deployment:**
- 2 replicas
- HPA: 2-4 replicas
- Resource limits: 500m CPU, 512Mi memory per replica

---

### 4. Lineage Service

**Responsibility:** Store, query, and visualize data lineage.

**Technology:**
- Go
- Neo4j (graph database)
- Redis (query cache)

**Features:**
- Store lineage edges (dataset → dataset, pipeline → dataset)
- Query upstream dependencies
- Query downstream dependencies
- Impact analysis (what breaks if X changes)
- Column-level lineage tracking
- Lineage visualization (graph rendering)
- Lineage search

**Internal API (gRPC):**
```protobuf
service LineageService {
  rpc AddLineageEdge(AddLineageEdgeRequest) returns (Empty);
  rpc GetUpstream(GetUpstreamRequest) returns (GetUpstreamResponse);
  rpc GetDownstream(GetDownstreamRequest) returns (GetDownstreamResponse);
  rpc GetImpact(GetImpactRequest) returns (GetImpactResponse);
  rpc SearchLineage(SearchLineageRequest) returns (SearchLineageResponse);
}

message AddLineageEdgeRequest {
  string source_id = 1;
  string target_id = 2;
  string edge_type = 3; // "dataset_to_dataset", "pipeline_to_dataset"
  map<string, string> metadata = 4;
}

message GetUpstreamRequest {
  string resource_id = 1;
  int32 max_depth = 2; // -1 for unlimited
}

message GetUpstreamResponse {
  repeated LineageNode nodes = 1;
  repeated LineageEdge edges = 2;
}
```

**Neo4j Schema:**
```cypher
// Nodes
CREATE (d:Dataset {
  id: "dataset-123",
  name: "orders",
  type: "table",
  owner: "data-team@company.com"
})

CREATE (p:Pipeline {
  id: "pipeline-456",
  name: "customer_360",
  type: "dbt_model",
  owner: "analytics-team@company.com"
})

// Relationships
CREATE (d1:Dataset)-[:DEPENDS_ON {
  type: "table_dependency",
  created_at: timestamp()
}]->(d2:Dataset)

CREATE (p:Pipeline)-[:PRODUCES {
  type: "transformation",
  created_at: timestamp()
}]->(d:Dataset)

CREATE (p:Pipeline)-[:CONSUMES {
  type: "input",
  created_at: timestamp()
}]->(d:Dataset)
```

**Common Lineage Queries:**

**1. Get all upstream dependencies:**
```cypher
MATCH path = (d:Dataset {id: $dataset_id})<-[:DEPENDS_ON*]-(upstream)
RETURN upstream, path
```

**2. Get all downstream dependencies:**
```cypher
MATCH path = (d:Dataset {id: $dataset_id})-[:DEPENDS_ON*]->(downstream)
RETURN downstream, path
```

**3. Impact analysis (what breaks if dataset is deleted):**
```cypher
MATCH path = (d:Dataset {id: $dataset_id})-[:DEPENDS_ON*]->(affected)
RETURN affected, COUNT(path) as impact_count
ORDER BY impact_count DESC
```

**4. Find critical datasets (most downstream dependencies):**
```cypher
MATCH (d:Dataset)-[:DEPENDS_ON*]->(downstream)
RETURN d, COUNT(downstream) as dependency_count
ORDER BY dependency_count DESC
LIMIT 10
```

**Deployment:**
- Neo4j: 1 primary + 2 read replicas
- Lineage Service: 2 replicas
- HPA: 2-5 replicas
- Resource limits: 1 CPU, 1Gi memory per replica

---

### 5. Reconciliation Service

**Responsibility:** Detect drift between desired state and actual state.

**Technology:**
- Go
- PostgreSQL (desired state)
- Redis (actual state cache)
- NATS/Kafka (event streaming)

**Features:**
- Continuous reconciliation loop (every 1 minute)
- Compare desired state vs actual state
- Detect drift (cadence, SLA, quality, ownership, cost)
- Generate drift reports
- Suggest reconciliation actions
- Track drift history

**Internal API (gRPC):**
```protobuf
service ReconciliationService {
  rpc StartReconciliation(StartReconciliationRequest) returns (ReconciliationJob);
  rpc GetReconciliationStatus(GetReconciliationStatusRequest) returns (ReconciliationStatus);
  rpc ListDrift(ListDriftRequest) returns (ListDriftResponse);
  rpc ReconcileDrift(ReconcileDriftRequest) returns (Empty);
}

message StartReconciliationRequest {
  repeated string resource_ids = 1; // empty = reconcile all
  string reconciliation_type = 2; // "cadence", "sla", "quality", "all"
}

message ReconciliationStatus {
  string job_id = 1;
  string status = 2; // "running", "completed", "failed"
  int32 resources_checked = 3;
  int32 drift_detected = 4;
  repeated Drift drifts = 5;
}

message Drift {
  string resource_id = 1;
  string drift_type = 2; // "cadence", "sla", "quality", "ownership", "cost"
  string severity = 3; // "critical", "high", "medium", "low"
  string desired_state = 4;
  string actual_state = 5;
  string suggested_action = 6;
  int64 detected_at = 7;
}
```

**Reconciliation Algorithm:**
```go
func (r *ReconciliationService) ReconcileLoop() {
  ticker := time.NewTicker(1 * time.Minute)
  defer ticker.Stop()
  
  for {
    select {
    case <-ticker.C:
      // 1. Fetch all resources from Resource Service
      resources, err := r.resourceClient.ListAllResources()
      if err != nil {
        log.Error("Failed to fetch resources", err)
        continue
      }
      
      // 2. For each resource, compare desired vs actual
      for _, resource := range resources {
        // Fetch desired state (from control plane)
        desiredState := resource.DesiredState
        
        // Fetch actual state (from providers)
        actualState, err := r.providerClient.GetActualState(resource.ID)
        if err != nil {
          log.Error("Failed to fetch actual state", err)
          continue
        }
        
        // Compare and detect drift
        drift := r.detectDrift(desiredState, actualState)
        if drift != nil {
          // Store drift
          r.storeDrift(drift)
          
          // Send alert
          r.notificationClient.SendDriftAlert(drift)
          
          // Generate suggested action (via AI)
          action := r.aiClient.GenerateReconciliationAction(drift)
          drift.SuggestedAction = action
        }
      }
    }
  }
}

func (r *ReconciliationService) detectDrift(desired, actual interface{}) *Drift {
  // Cadence drift
  if desired.Cadence != actual.Cadence {
    return &Drift{
      Type: "cadence",
      DesiredState: desired.Cadence,
      ActualState: actual.Cadence,
      Severity: "medium",
    }
  }
  
  // SLA drift
  if actual.Freshness > desired.SLA.Freshness {
    return &Drift{
      Type: "sla",
      DesiredState: fmt.Sprintf("freshness < %s", desired.SLA.Freshness),
      ActualState: fmt.Sprintf("freshness = %s", actual.Freshness),
      Severity: "high",
    }
  }
  
  // Quality drift
  if actual.QualityScore < desired.QualityThreshold {
    return &Drift{
      Type: "quality",
      DesiredState: fmt.Sprintf("quality >= %f", desired.QualityThreshold),
      ActualState: fmt.Sprintf("quality = %f", actual.QualityScore),
      Severity: "critical",
    }
  }
  
  // No drift detected
  return nil
}
```

**Deployment:**
- 2 replicas (for HA)
- HPA: 2-4 replicas
- Resource limits: 500m CPU, 512Mi memory per replica

---

### 6. Incident Service

**Responsibility:** Detect, track, and explain data incidents.

**Technology:**
- Go
- PostgreSQL (incident storage)
- Redis (incident cache)
- NATS/Kafka (event streaming)

**Features:**
- Incident detection (from provider events, drift alerts, quality checks)
- Root cause analysis (lineage-based)
- Impact analysis (affected resources)
- Incident tracking (status, resolution)
- Incident notifications (Slack, email, PagerDuty)
- Incident history and analytics

**Internal API (gRPC):**
```protobuf
service IncidentService {
  rpc CreateIncident(CreateIncidentRequest) returns (Incident);
  rpc GetIncident(GetIncidentRequest) returns (Incident);
  rpc UpdateIncident(UpdateIncidentRequest) returns (Incident);
  rpc ResolveIncident(ResolveIncidentRequest) returns (Incident);
  rpc ListIncidents(ListIncidentsRequest) returns (ListIncidentsResponse);
  rpc AnalyzeRootCause(AnalyzeRootCauseRequest) returns (RootCauseAnalysis);
}

message Incident {
  string id = 1;
  string type = 2; // "pipeline_failure", "quality_breach", "sla_violation"
  string severity = 3; // "critical", "high", "medium", "low"
  string resource_id = 4;
  string description = 5;
  RootCauseAnalysis root_cause = 6;
  repeated string affected_resources = 7;
  string status = 8; // "open", "investigating", "resolved"
  string resolution = 9;
  int64 created_at = 10;
  int64 resolved_at = 11;
}

message RootCauseAnalysis {
  string root_cause = 1;
  string explanation = 2;
  repeated string contributing_factors = 3;
  repeated string suggested_actions = 4;
  float confidence = 5; // 0.0 - 1.0
}
```

**Incident Detection Flow:**
```go
func (i *IncidentService) DetectIncidents() {
  // Listen to events from providers
  i.eventStream.Subscribe("provider.pipeline.failed", func(event Event) {
    // Create incident
    incident := &Incident{
      Type: "pipeline_failure",
      Severity: "high",
      ResourceID: event.ResourceID,
      Description: event.Description,
      Status: "open",
      CreatedAt: time.Now().Unix(),
    }
    
    // Analyze root cause (via AI + lineage)
    rootCause := i.analyzeRootCause(incident)
    incident.RootCause = rootCause
    
    // Analyze impact (via lineage)
    affectedResources := i.lineageClient.GetDownstream(event.ResourceID)
    incident.AffectedResources = affectedResources
    
    // Store incident
    i.storeIncident(incident)
    
    // Send notification
    i.notificationClient.SendIncidentAlert(incident)
  })
  
  // Listen to drift alerts
  i.eventStream.Subscribe("drift.detected", func(event Event) {
    // Create incident if drift is critical
    if event.Severity == "critical" {
      incident := &Incident{
        Type: "sla_violation",
        Severity: event.Severity,
        ResourceID: event.ResourceID,
        Description: event.Description,
        Status: "open",
        CreatedAt: time.Now().Unix(),
      }
      
      // ... same as above
    }
  })
}

func (i *IncidentService) analyzeRootCause(incident *Incident) *RootCauseAnalysis {
  // 1. Trace lineage backwards
  upstream := i.lineageClient.GetUpstream(incident.ResourceID)
  
  // 2. Check upstream for failures/delays
  for _, node := range upstream {
    status := i.providerClient.GetResourceStatus(node.ID)
    if status.Failed || status.Delayed {
      return &RootCauseAnalysis{
        RootCause: fmt.Sprintf("Upstream %s failed/delayed", node.Name),
        Explanation: "...",
        Confidence: 0.9,
      }
    }
  }
  
  // 3. Use AI to analyze logs and generate explanation
  aiAnalysis := i.aiClient.AnalyzeIncident(incident, upstream)
  
  return aiAnalysis
}
```

**Deployment:**
- 2 replicas
- HPA: 2-4 replicas
- Resource limits: 500m CPU, 512Mi memory per replica

---

### 7. AI Operator Service

**Responsibility:** AI capabilities (blueprint generation, incident explanation, NL queries).

**Technology:**
- Go
- OpenAI API / Anthropic Claude API
- Redis (prompt cache, response cache)

**Features:**
- Intent → Blueprint generation
- Blueprint validation
- Incident explanation
- Natural language queries
- Prompt management
- Response caching (reduce LLM API costs)
- Token usage tracking

**Internal API (gRPC):**
```protobuf
service AIOperatorService {
  rpc GenerateBlueprint(GenerateBlueprintRequest) returns (GenerateBlueprintResponse);
  rpc ValidateBlueprint(ValidateBlueprintRequest) returns (ValidateBlueprintResponse);
  rpc ExplainIncident(ExplainIncidentRequest) returns (ExplainIncidentResponse);
  rpc AnswerQuery(AnswerQueryRequest) returns (AnswerQueryResponse);
}

message GenerateBlueprintRequest {
  string intent = 1; // Natural language description
  map<string, string> context = 2;
}

message GenerateBlueprintResponse {
  string blueprint_yaml = 1;
  repeated string warnings = 2;
  float confidence = 3;
}

message ExplainIncidentRequest {
  Incident incident = 1;
  repeated LineageNode upstream = 2;
  repeated string logs = 3;
}

message ExplainIncidentResponse {
  string root_cause = 1;
  string explanation = 2;
  repeated string suggested_actions = 3;
  float confidence = 4;
}
```

**AI Integration Architecture:**

```go
type AIOperator struct {
  llmClient LLMClient // OpenAI or Anthropic
  promptManager *PromptManager
  cache *redis.Client
}

func (ai *AIOperator) GenerateBlueprint(intent string) (string, error) {
  // 1. Check cache
  cacheKey := fmt.Sprintf("blueprint:%s", hash(intent))
  if cached, err := ai.cache.Get(cacheKey).Result(); err == nil {
    return cached, nil
  }
  
  // 2. Build prompt
  prompt := ai.promptManager.GetPrompt("generate_blueprint", map[string]interface{}{
    "intent": intent,
    "examples": ai.promptManager.GetExamples("blueprint"),
  })
  
  // 3. Call LLM API
  response, err := ai.llmClient.Complete(prompt, LLMConfig{
    Model: "gpt-4",
    Temperature: 0.2, // Low temperature for structured output
    MaxTokens: 2000,
  })
  if err != nil {
    return "", err
  }
  
  // 4. Parse response (extract YAML)
  blueprint := extractYAML(response)
  
  // 5. Cache response (24 hours)
  ai.cache.Set(cacheKey, blueprint, 24*time.Hour)
  
  return blueprint, nil
}

func (ai *AIOperator) ExplainIncident(incident *Incident, upstream []LineageNode, logs []string) (string, error) {
  // 1. Build context
  context := map[string]interface{}{
    "incident": incident,
    "upstream": upstream,
    "logs": logs,
  }
  
  // 2. Build prompt
  prompt := ai.promptManager.GetPrompt("explain_incident", context)
  
  // 3. Call LLM API
  response, err := ai.llmClient.Complete(prompt, LLMConfig{
    Model: "claude-2", // Claude is cheaper and has 100k context window
    Temperature: 0.3,
    MaxTokens: 1000,
  })
  if err != nil {
    return "", err
  }
  
  return response, nil
}
```

**Prompt Management:**

Store prompts in files for easy versioning and A/B testing:

```yaml
# prompts/generate_blueprint.yaml
name: generate_blueprint
version: 1.0.0
template: |
  You are a data platform architect. Convert the following natural language intent into a DCraft Fusion blueprint (YAML format).
  
  Intent: {{.intent}}
  
  Output a YAML blueprint with the following structure:
  - dataset: name, owner, sla, contracts, dependencies
  - pipeline: name, type, cadence, owner, inputs, outputs
  - quality: rules, severity
  
  Examples:
  {{range .examples}}
  Intent: {{.intent}}
  Blueprint:
  ```yaml
  {{.blueprint}}
  ```
  {{end}}
  
  Now generate the blueprint for the given intent. Output ONLY the YAML, no explanation.

# prompts/explain_incident.yaml
name: explain_incident
version: 1.0.0
template: |
  You are a data platform expert. Analyze the following incident and explain the root cause.
  
  Incident:
  - Type: {{.incident.type}}
  - Resource: {{.incident.resource_id}}
  - Description: {{.incident.description}}
  
  Upstream Dependencies:
  {{range .upstream}}
  - {{.name}} (status: {{.status}})
  {{end}}
  
  Recent Logs:
  {{range .logs}}
  {{.}}
  {{end}}
  
  Provide:
  1. Root cause (1-2 sentences)
  2. Detailed explanation (3-5 sentences)
  3. Suggested actions (3-5 bullet points)
  
  Be specific and actionable. Focus on facts from the data provided.
```

**Cost Optimization:**

1. **Response Caching:**
   - Cache identical queries for 24 hours
   - Reduces API calls by ~60%

2. **Model Selection:**
   - Use GPT-4 for complex reasoning (blueprint generation)
   - Use Claude-2 for explanations (cheaper, larger context window)
   - Use GPT-3.5 for simple queries (10x cheaper)

3. **Token Optimization:**
   - Remove redundant context
   - Use concise prompts
   - Set appropriate max_tokens limits

4. **Batch Processing:**
   - Batch multiple queries when possible
   - Use async processing for non-urgent queries

**Deployment:**
- 2 replicas
- HPA: 2-6 replicas (based on LLM API latency)
- Resource limits: 500m CPU, 512Mi memory per replica

---

### 8. Notification Service

**Responsibility:** Send notifications to users (Slack, email, PagerDuty, webhooks).

**Technology:**
- Go
- Redis (notification queue)
- SMTP (email)
- Slack API
- PagerDuty API

**Features:**
- Multi-channel notifications (Slack, email, PagerDuty, webhooks)
- Notification templates
- Notification preferences (per user, per team)
- Notification batching (avoid spam)
- Notification history

**Internal API (gRPC):**
```protobuf
service NotificationService {
  rpc SendNotification(SendNotificationRequest) returns (Empty);
  rpc SendBatchNotifications(SendBatchNotificationsRequest) returns (Empty);
  rpc GetNotificationHistory(GetNotificationHistoryRequest) returns (GetNotificationHistoryResponse);
}

message SendNotificationRequest {
  string recipient = 1; // email, Slack channel, PagerDuty user
  string channel = 2; // "email", "slack", "pagerduty", "webhook"
  string template = 3; // "incident_alert", "drift_alert", "sla_violation"
  map<string, string> data = 4;
  string severity = 5; // "critical", "high", "medium", "low"
}
```

**Deployment:**
- 2 replicas
- HPA: 2-4 replicas
- Resource limits: 250m CPU, 256Mi memory per replica

---

# Communication Patterns

## Internal Communication: gRPC

**Why gRPC for internal service-to-service communication?**

1. **Performance:**
   - Binary serialization (Protocol Buffers) is 5-10x faster than JSON
   - HTTP/2 multiplexing reduces latency
   - Streaming support (unary, server streaming, client streaming, bidirectional)

2. **Type Safety:**
   - Strongly typed contracts (.proto files)
   - Auto-generated client/server code
   - Compile-time validation

3. **Ecosystem:**
   - Excellent Go support
   - Built-in load balancing
   - Built-in retries and timeouts
   - Built-in health checking

**gRPC Configuration:**

```go
// Server setup
grpcServer := grpc.NewServer(
  grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
    grpc_prometheus.UnaryServerInterceptor,
    grpc_auth.UnaryServerInterceptor(authFunc),
    grpc_recovery.UnaryServerInterceptor(),
  )),
  grpc.MaxRecvMsgSize(10 * 1024 * 1024), // 10MB
  grpc.MaxSendMsgSize(10 * 1024 * 1024), // 10MB
  grpc.KeepaliveParams(keepalive.ServerParameters{
    MaxConnectionIdle: 15 * time.Second,
    Timeout:           5 * time.Second,
  }),
)

// Client setup
conn, err := grpc.Dial(
  "resource-service:50051",
  grpc.WithInsecure(),
  grpc.WithBlock(),
  grpc.WithTimeout(5*time.Second),
  grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
    grpc_prometheus.UnaryClientInterceptor,
    grpc_retry.UnaryClientInterceptor(
      grpc_retry.WithMax(3),
      grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100*time.Millisecond)),
    ),
  )),
)
```

---

## External Communication: REST + GraphQL

**REST API for simple CRUD operations:**
- Easy to understand and use
- Broad compatibility (curl, Postman, etc.)
- Good for public APIs

**GraphQL for complex queries:**
- Flexible data fetching (clients specify what they need)
- Reduces over-fetching and under-fetching
- Single endpoint for all queries
- Good for frontend applications

**API Gateway handles protocol translation:**
- REST → gRPC (for internal services)
- GraphQL → gRPC (for internal services)

---

## Event Streaming: NATS or Kafka

**For asynchronous communication and event-driven architecture.**

**NATS (Recommended for V0):**
- Lightweight and fast
- Simple pub/sub model
- Built-in persistence (JetStream)
- Lower operational overhead than Kafka
- Good for <100k messages/second

**Kafka (Consider for V1):**
- Higher throughput (>1M messages/second)
- Better for large-scale deployments
- More complex to operate
- Better ecosystem for stream processing

**Use Cases:**
- Provider events (pipeline succeeded/failed)
- Drift alerts
- Incident notifications
- Audit logs

**Event Schema:**
```protobuf
message Event {
  string id = 1;
  string type = 2; // "provider.pipeline.failed", "drift.detected"
  string source = 3; // service that emitted the event
  string resource_id = 4;
  map<string, string> data = 5;
  int64 timestamp = 6;
}
```

---

# Data Storage Architecture

## Storage Strategy: Polyglot Persistence

Use the right database for each use case.

### 1. PostgreSQL (Primary Metadata Store)

**Use for:**
- Structured resources (Datasets, Pipelines, DataProducts, Policies)
- Transactional data (ACID guarantees)
- Complex queries with JOINs
- Full-text search (with pg_trgm extension)

**Configuration:**
```yaml
# PostgreSQL 15
version: 15
replicas: 1 primary + 2 read replicas
storage: 100GB SSD
connection_pool: 100 connections
extensions:
  - pg_trgm (fuzzy text search)
  - uuid-ossp (UUID generation)
  - timescaledb (time-series data)
```

**Schema Design Principles:**
- Normalize to 3NF (avoid data duplication)
- Use UUIDs for primary keys (distributed-friendly)
- Use JSONB for flexible metadata
- Use indexes for frequently queried columns
- Use foreign keys for referential integrity
- Use triggers for audit trails

---

### 2. Neo4j (Lineage Graph Store)

**Use for:**
- Data lineage (dataset → dataset, pipeline → dataset)
- Graph traversals (upstream, downstream, impact analysis)
- Relationship-heavy queries

**Configuration:**
```yaml
# Neo4j 5.x
version: 5.x
deployment: 1 primary + 2 read replicas
storage: 50GB SSD
memory: 8GB heap, 8GB page cache
indexes:
  - Node label indexes (Dataset, Pipeline)
  - Property indexes (id, name)
```

**Why Neo4j over PostgreSQL for lineage?**
- 10-100x faster for graph traversals
- More intuitive query language (Cypher)
- Native graph storage (no expensive JOINs)
- Better for complex lineage queries (multi-hop dependencies)

**Example: PostgreSQL vs Neo4j for lineage query**

PostgreSQL (recursive CTE):
```sql
WITH RECURSIVE lineage AS (
  SELECT id, name, 0 as depth
  FROM datasets
  WHERE id = 'dataset-123'
  
  UNION ALL
  
  SELECT d.id, d.name, l.depth + 1
  FROM datasets d
  JOIN dataset_dependencies dd ON d.id = dd.upstream_id
  JOIN lineage l ON dd.downstream_id = l.id
  WHERE l.depth < 10
)
SELECT * FROM lineage;
```

Neo4j (Cypher):
```cypher
MATCH path = (d:Dataset {id: 'dataset-123'})<-[:DEPENDS_ON*..10]-(upstream)
RETURN upstream, path
```

Neo4j is simpler and faster for this query.

---

### 3. Redis (Caching & Session Store)

**Use for:**
- Response caching (API responses, LLM responses)
- Session storage (user sessions, JWT tokens)
- Rate limiting (token bucket)
- Provider state cache (avoid repeated API calls)
- Real-time metrics (counters, gauges)

**Configuration:**
```yaml
# Redis 7.x
version: 7.x
deployment: 1 primary + 2 read replicas
storage: 10GB memory
persistence: AOF (append-only file)
eviction_policy: allkeys-lru (least recently used)
```

**Caching Strategy:**
```go
// Cache API responses
func (s *ResourceService) GetDataset(id string) (*Dataset, error) {
  // 1. Check cache
  cacheKey := fmt.Sprintf("dataset:%s", id)
  if cached, err := s.cache.Get(cacheKey).Result(); err == nil {
    var dataset Dataset
    json.Unmarshal([]byte(cached), &dataset)
    return &dataset, nil
  }
  
  // 2. Query database
  dataset, err := s.db.GetDataset(id)
  if err != nil {
    return nil, err
  }
  
  // 3. Cache result (1 hour TTL)
  data, _ := json.Marshal(dataset)
  s.cache.Set(cacheKey, data, 1*time.Hour)
  
  return dataset, nil
}
```

---

### 4. TimescaleDB (Time-Series Metrics)

**Use for:**
- Pipeline execution metrics (runtime, success rate)
- SLA metrics (freshness, availability)
- Cost metrics (compute, storage)
- Drift metrics (drift count, severity)

**Why TimescaleDB?**
- PostgreSQL extension (familiar SQL interface)
- Optimized for time-series data (10-100x faster inserts)
- Automatic data retention policies
- Continuous aggregates (pre-computed rollups)

**Configuration:**
```yaml
# TimescaleDB (PostgreSQL extension)
version: 2.x
deployment: Same as PostgreSQL
storage: 50GB SSD
retention: 90 days (raw data), 1 year (aggregates)
```

**Schema:**
```sql
-- Create hypertable (time-series table)
CREATE TABLE pipeline_metrics (
  time TIMESTAMPTZ NOT NULL,
  pipeline_id UUID NOT NULL,
  metric_name VARCHAR(50) NOT NULL,
  metric_value DOUBLE PRECISION NOT NULL,
  tags JSONB
);

SELECT create_hypertable('pipeline_metrics', 'time');

-- Create continuous aggregate (hourly rollup)
CREATE MATERIALIZED VIEW pipeline_metrics_hourly
WITH (timescaledb.continuous) AS
SELECT
  time_bucket('1 hour', time) AS bucket,
  pipeline_id,
  metric_name,
  AVG(metric_value) AS avg_value,
  MAX(metric_value) AS max_value,
  MIN(metric_value) AS min_value
FROM pipeline_metrics
GROUP BY bucket, pipeline_id, metric_name;

-- Add retention policy (delete data older than 90 days)
SELECT add_retention_policy('pipeline_metrics', INTERVAL '90 days');
```

---

## Database Deployment Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      KUBERNETES CLUSTER                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  PostgreSQL (StatefulSet)                          │     │
│  │  - 1 Primary (read/write)                          │     │
│  │  - 2 Read Replicas (read-only)                     │     │
│  │  - Persistent Volume: 100GB SSD                    │     │
│  │  - Backup: Daily to S3                             │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  Neo4j (StatefulSet)                               │     │
│  │  - 1 Primary (read/write)                          │     │
│  │  - 2 Read Replicas (read-only)                     │     │
│  │  - Persistent Volume: 50GB SSD                     │     │
│  │  - Backup: Daily to S3                             │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  Redis (StatefulSet)                               │     │
│  │  - 1 Primary (read/write)                          │     │
│  │  - 2 Read Replicas (read-only)                     │     │
│  │  - Persistent Volume: 10GB SSD                     │     │
│  │  - Backup: Daily AOF snapshots                     │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

# AI Integration Architecture

## LLM Provider Strategy

### Primary: Anthropic Claude

**Why Claude?**
- 100k token context window (vs 8k-32k for GPT)
- 50% cheaper than GPT-4
- Better for long-context tasks (incident analysis, lineage reasoning)
- Strong reasoning capabilities

**Use for:**
- Incident explanation (requires large context: logs, lineage, metrics)
- Blueprint generation (requires examples and context)
- Natural language queries (requires platform metadata)

### Secondary: OpenAI GPT-4

**Why GPT-4?**
- Best-in-class reasoning for complex tasks
- Better for structured output (JSON, YAML)
- Stronger function calling capabilities

**Use for:**
- Blueprint validation (requires structured reasoning)
- Complex impact analysis (requires multi-step reasoning)

### Tertiary: OpenAI GPT-3.5

**Why GPT-3.5?**
- 10x cheaper than GPT-4
- Fast response times
- Good for simple queries

**Use for:**
- Simple natural language queries ("Who owns dataset X?")
- Quick lookups (no complex reasoning required)

---

## AI Service Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    AI OPERATOR SERVICE                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  Prompt Manager                                    │     │
│  │  - Load prompts from files                         │     │
│  │  - Template rendering (Go templates)               │     │
│  │  - Prompt versioning                               │     │
│  │  - A/B testing support                             │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  LLM Router                                        │     │
│  │  - Route requests to appropriate LLM               │     │
│  │  - Model selection based on task complexity        │     │
│  │  - Fallback to cheaper models on error             │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  Response Cache (Redis)                            │     │
│  │  - Cache identical queries (24 hours)              │     │
│  │  - Cache key: hash(prompt + context)               │     │
│  │  - Reduces API calls by ~60%                       │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  Token Usage Tracker                               │     │
│  │  - Track tokens per request                        │     │
│  │  - Track cost per request                          │     │
│  │  - Alert on budget overruns                        │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
└──────────────────┬───────────────────────────────────────────┘
                   │
         ┌─────────┴─────────┬─────────────────┐
         │                   │                 │
         ▼                   ▼                 ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ Anthropic       │  │ OpenAI          │  │ OpenAI          │
│ Claude API      │  │ GPT-4 API       │  │ GPT-3.5 API     │
│                 │  │                 │  │                 │
│ (Primary)       │  │ (Secondary)     │  │ (Tertiary)      │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

---

## Prompt Engineering Best Practices

### 1. Structured Prompts

**Template:**
```
System: You are a [role]. Your goal is [goal].

Context:
[Relevant context from control plane]

Task:
[Specific task description]

Output Format:
[Expected output format]

Examples:
[Few-shot examples]

Now perform the task.
```

### 2. Few-Shot Learning

Provide 2-3 examples of desired output:

```yaml
# prompts/generate_blueprint.yaml
examples:
  - intent: "Sync MySQL orders to Snowflake every hour"
    blueprint: |
      dataset:
        name: orders
        source:
          type: mysql
          table: orders
        destination:
          type: snowflake
          table: raw.orders
        cadence: 1h
        owner: data-team@company.com
  
  - intent: "Transform raw customers with dbt, dedupe by email"
    blueprint: |
      pipeline:
        name: customer_deduplication
        type: dbt_model
        inputs:
          - raw.customers
        outputs:
          - staging.customers
        contract:
          deduplication:
            key: email
            strategy: last_write_wins
        owner: analytics-team@company.com
```

### 3. Chain of Thought (CoT)

For complex reasoning, ask LLM to think step-by-step:

```
Analyze the following incident and explain the root cause.

Think step-by-step:
1. What happened? (describe the incident)
2. When did it happen? (timeline)
3. What are the upstream dependencies? (lineage)
4. Did any upstream dependencies fail or delay? (check status)
5. What is the most likely root cause? (conclusion)
6. What actions should be taken? (recommendations)

Now provide your analysis.
```

### 4. Output Constraints

Constrain output format for structured data:

```
Generate a YAML blueprint for the following intent.

Output ONLY valid YAML. Do not include explanations or markdown code blocks.
The YAML must conform to the following schema:
- dataset: (required)
  - name: string
  - owner: string
  - sla: object
    - freshness: duration
    - severity: string
```

---

## Cost Optimization Strategies

### 1. Response Caching

**Cache identical queries for 24 hours:**

```go
func (ai *AIOperator) AnswerQuery(query string) (string, error) {
  // Generate cache key
  cacheKey := fmt.Sprintf("ai:query:%s", hash(query))
  
  // Check cache
  if cached, err := ai.cache.Get(cacheKey).Result(); err == nil {
    ai.metrics.CacheHits.Inc()
    return cached, nil
  }
  
  // Call LLM API
  response, err := ai.llmClient.Complete(query)
  if err != nil {
    return "", err
  }
  
  // Cache response (24 hours)
  ai.cache.Set(cacheKey, response, 24*time.Hour)
  ai.metrics.CacheMisses.Inc()
  
  return response, nil
}
```

**Expected cache hit rate: 60%**
**Cost reduction: 60%**

---

### 2. Model Selection

**Route requests to appropriate model based on complexity:**

```go
func (ai *AIOperator) SelectModel(task string, contextSize int) string {
  // Simple queries → GPT-3.5 (10x cheaper)
  if task == "simple_query" && contextSize < 1000 {
    return "gpt-3.5-turbo"
  }
  
  // Long context → Claude (100k context window, 50% cheaper than GPT-4)
  if contextSize > 8000 {
    return "claude-2"
  }
  
  // Complex reasoning → GPT-4 (best reasoning)
  if task == "complex_reasoning" {
    return "gpt-4"
  }
  
  // Default → Claude (good balance)
  return "claude-2"
}
```

---

### 3. Token Optimization

**Remove redundant context:**

```go
func (ai *AIOperator) OptimizeContext(context map[string]interface{}) map[string]interface{} {
  // Remove unnecessary fields
  delete(context, "metadata")
  delete(context, "created_at")
  delete(context, "updated_at")
  
  // Truncate long strings
  if desc, ok := context["description"].(string); ok && len(desc) > 500 {
    context["description"] = desc[:500] + "..."
  }
  
  // Remove duplicate information
  // ... (custom logic based on use case)
  
  return context
}
```

---

### 4. Batch Processing

**Batch multiple queries when possible:**

```go
func (ai *AIOperator) AnswerBatchQueries(queries []string) ([]string, error) {
  // Combine queries into single prompt
  prompt := "Answer the following queries:\n\n"
  for i, query := range queries {
    prompt += fmt.Sprintf("%d. %s\n", i+1, query)
  }
  
  // Single LLM API call (cheaper than N separate calls)
  response, err := ai.llmClient.Complete(prompt)
  if err != nil {
    return nil, err
  }
  
  // Parse response (split by numbered answers)
  answers := parseNumberedAnswers(response)
  
  return answers, nil
}
```

---

### 5. Budget Alerts

**Track token usage and alert on overruns:**

```go
func (ai *AIOperator) TrackTokenUsage(model string, promptTokens, completionTokens int) {
  // Calculate cost
  cost := calculateCost(model, promptTokens, completionTokens)
  
  // Update metrics
  ai.metrics.TokensUsed.Add(float64(promptTokens + completionTokens))
  ai.metrics.CostIncurred.Add(cost)
  
  // Check budget
  dailyBudget := 100.0 // $100/day
  dailyCost := ai.metrics.GetDailyCost()
  
  if dailyCost > dailyBudget {
    ai.alertManager.SendAlert("AI budget exceeded", fmt.Sprintf(
      "Daily cost: $%.2f, Budget: $%.2f",
      dailyCost, dailyBudget,
    ))
  }
}

func calculateCost(model string, promptTokens, completionTokens int) float64 {
  costs := map[string]struct{ prompt, completion float64 }{
    "gpt-4":         {0.03, 0.06},   // per 1k tokens
    "gpt-3.5-turbo": {0.0015, 0.002}, // per 1k tokens
    "claude-2":      {0.01102, 0.03268}, // per 1k tokens
  }
  
  c := costs[model]
  return (float64(promptTokens)/1000)*c.prompt + (float64(completionTokens)/1000)*c.completion
}
```

---

## Error Handling & Reliability

### 1. Exponential Backoff Retry

```go
func (ai *AIOperator) CallLLMWithRetry(prompt string) (string, error) {
  maxRetries := 3
  baseDelay := 1 * time.Second
  
  for i := 0; i < maxRetries; i++ {
    response, err := ai.llmClient.Complete(prompt)
    
    if err == nil {
      return response, nil
    }
    
    // Check if error is retryable
    if !isRetryable(err) {
      return "", err
    }
    
    // Exponential backoff
    delay := baseDelay * time.Duration(math.Pow(2, float64(i)))
    time.Sleep(delay)
  }
  
  return "", fmt.Errorf("max retries exceeded")
}

func isRetryable(err error) bool {
  // Retry on rate limit, timeout, server errors
  return strings.Contains(err.Error(), "rate limit") ||
         strings.Contains(err.Error(), "timeout") ||
         strings.Contains(err.Error(), "500") ||
         strings.Contains(err.Error(), "503")
}
```

---

### 2. Fallback Strategy

```go
func (ai *AIOperator) AnswerQueryWithFallback(query string) (string, error) {
  // Try primary model (Claude)
  response, err := ai.callClaude(query)
  if err == nil {
    return response, nil
  }
  
  log.Warn("Claude failed, falling back to GPT-4", err)
  
  // Fallback to GPT-4
  response, err = ai.callGPT4(query)
  if err == nil {
    return response, nil
  }
  
  log.Warn("GPT-4 failed, falling back to GPT-3.5", err)
  
  // Fallback to GPT-3.5
  response, err = ai.callGPT35(query)
  if err == nil {
    return response, nil
  }
  
  // All models failed
  return "", fmt.Errorf("all LLM providers failed")
}
```

---

### 3. Circuit Breaker

```go
type CircuitBreaker struct {
  maxFailures int
  resetTimeout time.Duration
  failures int
  lastFailure time.Time
  state string // "closed", "open", "half-open"
}

func (cb *CircuitBreaker) Call(fn func() (string, error)) (string, error) {
  // Check circuit state
  if cb.state == "open" {
    if time.Since(cb.lastFailure) > cb.resetTimeout {
      cb.state = "half-open"
    } else {
      return "", fmt.Errorf("circuit breaker open")
    }
  }
  
  // Call function
  response, err := fn()
  
  if err != nil {
    cb.failures++
    cb.lastFailure = time.Now()
    
    if cb.failures >= cb.maxFailures {
      cb.state = "open"
    }
    
    return "", err
  }
  
  // Success - reset circuit
  cb.failures = 0
  cb.state = "closed"
  
  return response, nil
}
```

---

# Deployment Strategy

## Container Orchestration: Kubernetes

**Why Kubernetes?**
- Industry standard for container orchestration
- Excellent ecosystem (Helm, ArgoCD, Prometheus, etc.)
- Auto-scaling (HPA, VPA, Cluster Autoscaler)
- Self-healing (automatic restarts, health checks)
- Service discovery and load balancing
- Rolling updates and rollbacks

---

## Deployment Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    KUBERNETES CLUSTER                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  NAMESPACE: dcraft-fusion                          │     │
│  │                                                    │     │
│  │  ┌──────────────────────────────────────────┐     │     │
│  │  │  API Gateway (Deployment)                │     │     │
│  │  │  - Replicas: 2-10 (HPA)                  │     │     │
│  │  │  - Service: LoadBalancer                 │     │     │
│  │  │  - Ingress: HTTPS (cert-manager)         │     │     │
│  │  └──────────────────────────────────────────┘     │     │
│  │                                                    │     │
│  │  ┌──────────────────────────────────────────┐     │     │
│  │  │  Microservices (Deployments)             │     │     │
│  │  │  - Resource Service (2-5 replicas)       │     │     │
│  │  │  - Provider Service (2-4 replicas)       │     │     │
│  │  │  - Lineage Service (2-5 replicas)        │     │     │
│  │  │  - Reconciliation Service (2-4 replicas) │     │     │
│  │  │  - Incident Service (2-4 replicas)       │     │     │
│  │  │  - AI Operator Service (2-6 replicas)    │     │     │
│  │  │  - Notification Service (2-4 replicas)   │     │     │
│  │  │  - Audit Service (2-4 replicas)          │     │     │
│  │  └──────────────────────────────────────────┘     │     │
│  │                                                    │     │
│  │  ┌──────────────────────────────────────────┐     │     │
│  │  │  Databases (StatefulSets)                │     │     │
│  │  │  - PostgreSQL (1 primary + 2 replicas)   │     │     │
│  │  │  - Neo4j (1 primary + 2 replicas)        │     │     │
│  │  │  - Redis (1 primary + 2 replicas)        │     │     │
│  │  └──────────────────────────────────────────┘     │     │
│  │                                                    │     │
│  │  ┌──────────────────────────────────────────┐     │     │
│  │  │  Message Queue                           │     │     │
│  │  │  - NATS (3 replicas)                     │     │     │
│  │  └──────────────────────────────────────────┘     │     │
│  │                                                    │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  NAMESPACE: monitoring                             │     │
│  │  - Prometheus (metrics)                            │     │
│  │  - Grafana (dashboards)                            │     │
│  │  - Jaeger (traces)                                 │     │
│  │  - ELK Stack (logs)                                │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Helm Charts

**Why Helm?**
- Package manager for Kubernetes
- Templating for Kubernetes manifests
- Version control for deployments
- Easy rollbacks
- Dependency management

**Helm Chart Structure:**
```
dcraft-fusion/
├── Chart.yaml              # Chart metadata
├── values.yaml             # Default configuration
├── values-dev.yaml         # Dev environment overrides
├── values-staging.yaml     # Staging environment overrides
├── values-prod.yaml        # Production environment overrides
├── templates/
│   ├── api-gateway/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── ingress.yaml
│   │   └── hpa.yaml
│   ├── resource-service/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── hpa.yaml
│   ├── provider-service/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── hpa.yaml
│   ├── lineage-service/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── hpa.yaml
│   ├── reconciliation-service/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── hpa.yaml
│   ├── incident-service/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── hpa.yaml
│   ├── ai-operator-service/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── hpa.yaml
│   ├── notification-service/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── hpa.yaml
│   ├── audit-service/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── hpa.yaml
│   ├── postgresql/
│   │   ├── statefulset.yaml
│   │   ├── service.yaml
│   │   └── pvc.yaml
│   ├── neo4j/
│   │   ├── statefulset.yaml
│   │   ├── service.yaml
│   │   └── pvc.yaml
│   ├── redis/
│   │   ├── statefulset.yaml
│   │   ├── service.yaml
│   │   └── pvc.yaml
│   └── nats/
│       ├── statefulset.yaml
│       └── service.yaml
└── charts/                 # Subchart dependencies
    ├── postgresql/
    ├── neo4j/
    ├── redis/
    └── nats/
```

**Example: API Gateway Deployment**

```yaml
# templates/api-gateway/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "dcraft-fusion.fullname" . }}-api-gateway
  labels:
    {{- include "dcraft-fusion.labels" . | nindent 4 }}
    app.kubernetes.io/component: api-gateway
spec:
  replicas: {{ .Values.apiGateway.replicaCount }}
  selector:
    matchLabels:
      {{- include "dcraft-fusion.selectorLabels" . | nindent 6 }}
      app.kubernetes.io/component: api-gateway
  template:
    metadata:
      labels:
        {{- include "dcraft-fusion.selectorLabels" . | nindent 8 }}
        app.kubernetes.io/component: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: "{{ .Values.apiGateway.image.repository }}:{{ .Values.apiGateway.image.tag }}"
        imagePullPolicy: {{ .Values.apiGateway.image.pullPolicy }}
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        - name: grpc
          containerPort: 50051
          protocol: TCP
        env:
        - name: POSTGRES_HOST
          value: {{ include "dcraft-fusion.fullname" . }}-postgresql
        - name: POSTGRES_PORT
          value: "5432"
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: {{ include "dcraft-fusion.fullname" . }}-postgresql
              key: username
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: {{ include "dcraft-fusion.fullname" . }}-postgresql
              key: password
        - name: REDIS_HOST
          value: {{ include "dcraft-fusion.fullname" . }}-redis
        - name: REDIS_PORT
          value: "6379"
        livenessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: http
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          {{- toYaml .Values.apiGateway.resources | nindent 10 }}
```

**Example: HPA (Horizontal Pod Autoscaler)**

```yaml
# templates/api-gateway/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: {{ include "dcraft-fusion.fullname" . }}-api-gateway
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ include "dcraft-fusion.fullname" . }}-api-gateway
  minReplicas: {{ .Values.apiGateway.autoscaling.minReplicas }}
  maxReplicas: {{ .Values.apiGateway.autoscaling.maxReplicas }}
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: {{ .Values.apiGateway.autoscaling.targetCPUUtilizationPercentage }}
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: {{ .Values.apiGateway.autoscaling.targetMemoryUtilizationPercentage }}
```

**Example: values.yaml**

```yaml
# Default values for dcraft-fusion

global:
  environment: production
  domain: dcraft-fusion.example.com

apiGateway:
  replicaCount: 2
  image:
    repository: dcraft-fusion/api-gateway
    tag: "v0.1.0"
    pullPolicy: IfNotPresent
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 500m
      memory: 512Mi
  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10
    targetCPUUtilizationPercentage: 70
    targetMemoryUtilizationPercentage: 80

resourceService:
  replicaCount: 2
  image:
    repository: dcraft-fusion/resource-service
    tag: "v0.1.0"
    pullPolicy: IfNotPresent
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 250m
      memory: 256Mi
  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 5
    targetCPUUtilizationPercentage: 70
    targetMemoryUtilizationPercentage: 80

# ... similar for other services

postgresql:
  enabled: true
  primary:
    persistence:
      size: 100Gi
      storageClass: fast-ssd
  readReplicas:
    replicaCount: 2
    persistence:
      size: 100Gi
      storageClass: fast-ssd

neo4j:
  enabled: true
  core:
    numberOfServers: 1
    persistentVolume:
      size: 50Gi
      storageClass: fast-ssd
  readReplica:
    numberOfServers: 2
    persistentVolume:
      size: 50Gi
      storageClass: fast-ssd

redis:
  enabled: true
  master:
    persistence:
      size: 10Gi
      storageClass: fast-ssd
  replica:
    replicaCount: 2
    persistence:
      size: 10Gi
      storageClass: fast-ssd

nats:
  enabled: true
  cluster:
    enabled: true
    replicas: 3
```

---

## GitOps with ArgoCD

**Why GitOps?**
- Git as single source of truth
- Automated deployments (no manual kubectl)
- Drift detection and auto-sync
- Audit trail (Git history)
- Easy rollbacks (Git revert)
- Multi-environment support

**GitOps Workflow:**

```
┌─────────────────────────────────────────────────────────────┐
│                      GITOPS WORKFLOW                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. Developer commits code to Git                            │
│     ↓                                                        │
│  2. CI builds Docker image                                   │
│     ↓                                                        │
│  3. CI pushes image to registry                              │
│     ↓                                                        │
│  4. CI updates Helm values in config repo                    │
│     (e.g., image tag: v0.1.0 → v0.1.1)                       │
│     ↓                                                        │
│  5. ArgoCD detects change in config repo                     │
│     ↓                                                        │
│  6. ArgoCD syncs to Kubernetes cluster                       │
│     ↓                                                        │
│  7. Kubernetes rolls out new version                         │
│     ↓                                                        │
│  8. Health checks pass → deployment complete                 │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Repository Structure:**

```
# Application code repository
dcraft-fusion/
├── cmd/
│   ├── api-gateway/
│   ├── resource-service/
│   ├── provider-service/
│   └── ...
├── pkg/
├── Dockerfile
└── .github/
    └── workflows/
        └── ci.yaml

# Configuration repository (separate)
dcraft-fusion-config/
├── environments/
│   ├── dev/
│   │   ├── values.yaml
│   │   └── secrets.yaml (encrypted)
│   ├── staging/
│   │   ├── values.yaml
│   │   └── secrets.yaml (encrypted)
│   └── prod/
│       ├── values.yaml
│       └── secrets.yaml (encrypted)
└── helm/
    └── dcraft-fusion/
        ├── Chart.yaml
        ├── values.yaml
        └── templates/
```

**ArgoCD Application:**

```yaml
# argocd/application.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: dcraft-fusion-prod
  namespace: argocd
spec:
  project: default
  
  source:
    repoURL: https://github.com/dcraft-fusion/dcraft-fusion-config
    targetRevision: main
    path: helm/dcraft-fusion
    helm:
      valueFiles:
        - ../../environments/prod/values.yaml
  
  destination:
    server: https://kubernetes.default.svc
    namespace: dcraft-fusion
  
  syncPolicy:
    automated:
      prune: true      # Delete resources not in Git
      selfHeal: true   # Sync when drift detected
      allowEmpty: false
    syncOptions:
      - CreateNamespace=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
```

**CI/CD Pipeline (GitHub Actions):**

```yaml
# .github/workflows/ci.yaml
name: CI/CD

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: go test -v ./...
      
      - name: Run linter
        run: golangci-lint run
  
  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      
      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: |
            dcraft-fusion/api-gateway:${{ github.sha }}
            dcraft-fusion/api-gateway:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
  
  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Checkout config repo
        uses: actions/checkout@v3
        with:
          repository: dcraft-fusion/dcraft-fusion-config
          token: ${{ secrets.CONFIG_REPO_TOKEN }}
      
      - name: Update image tag
        run: |
          # Update values.yaml with new image tag
          sed -i 's/tag: .*/tag: "${{ github.sha }}"/' environments/prod/values.yaml
      
      - name: Commit and push
        run: |
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add environments/prod/values.yaml
          git commit -m "Update image tag to ${{ github.sha }}"
          git push
      
      # ArgoCD will automatically detect the change and deploy
```

---

## Deployment Modes

### 1. Single-Cluster Deployment (Recommended for V0)

**Architecture:**
```
┌─────────────────────────────────────────────────────────────┐
│                    KUBERNETES CLUSTER                        │
│                    (Single Region)                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  - All services in one cluster                               │
│  - All databases in one cluster                              │
│  - Simple to manage                                          │
│  - Lower cost                                                │
│  - Good for <100 customers                                   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**When to use:**
- V0 (pilot phase)
- <100 customers
- Single region
- <10,000 datasets

---

### 2. Multi-Cluster Deployment (Consider for V1)

**Architecture:**
```
┌─────────────────────────────────────────────────────────────┐
│                    CONTROL PLANE CLUSTER                     │
│                    (Central)                                 │
├─────────────────────────────────────────────────────────────┤
│  - API Gateway                                               │
│  - AI Operator Service                                       │
│  - Notification Service                                      │
│  - Audit Service                                             │
└──────────────────┬───────────────────────────────────────────┘
                   │
         ┌─────────┴─────────┬─────────────────┐
         │                   │                 │
         ▼                   ▼                 ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ DATA CLUSTER 1  │  │ DATA CLUSTER 2  │  │ DATA CLUSTER 3  │
│ (US-EAST)       │  │ (US-WEST)       │  │ (EU-WEST)       │
├─────────────────┤  ├─────────────────┤  ├─────────────────┤
│ - Resource Svc  │  │ - Resource Svc  │  │ - Resource Svc  │
│ - Provider Svc  │  │ - Provider Svc  │  │ - Provider Svc  │
│ - Lineage Svc   │  │ - Lineage Svc   │  │ - Lineage Svc   │
│ - Databases     │  │ - Databases     │  │ - Databases     │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

**When to use:**
- V1 (scale phase)
- >100 customers
- Multi-region (data residency requirements)
- >10,000 datasets

---

### 3. Federated Deployment (Consider for V1)

**Architecture:**
```
┌─────────────────────────────────────────────────────────────┐
│              CENTRAL CONTROL PLANE (HQ)                      │
│  - Global lineage                                            │
│  - Cross-domain policies                                     │
│  - Federated query                                           │
└──────────────────┬───────────────────────────────────────────┘
                   │
         ┌─────────┴─────────┬─────────────────┐
         │                   │                 │
         ▼                   ▼                 ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ DOMAIN A        │  │ DOMAIN B        │  │ DOMAIN C        │
│ Control Plane   │  │ Control Plane   │  │ Control Plane   │
│                 │  │                 │  │                 │
│ - Local lineage │  │ - Local lineage │  │ - Local lineage │
│ - Local policies│  │ - Local policies│  │ - Local policies│
│ - Local tools   │  │ - Local tools   │  │ - Local tools   │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

**When to use:**
- Large enterprises
- Domain-driven architecture
- Multiple autonomous teams
- >50,000 datasets

---

# Development Workflow

## Local Development Setup

**Prerequisites:**
- Go 1.21+
- Docker Desktop
- kubectl
- Helm
- Tilt (optional, for hot reload)

**Quick Start:**

```bash
# 1. Clone repository
git clone https://github.com/dcraft-fusion/dcraft-fusion
cd dcraft-fusion

# 2. Start local Kubernetes cluster (Docker Desktop or Minikube)
# Docker Desktop: Enable Kubernetes in settings
# OR
minikube start

# 3. Install dependencies
make deps

# 4. Start local development environment
make dev

# This will:
# - Start PostgreSQL, Neo4j, Redis in Docker
# - Run database migrations
# - Start all microservices with hot reload
# - Start frontend dev server
```

**Makefile:**

```makefile
# Makefile

.PHONY: deps test build docker-build dev clean

# Install dependencies
deps:
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@latest

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run

# Build all services
build:
	go build -o bin/api-gateway cmd/api-gateway/main.go
	go build -o bin/resource-service cmd/resource-service/main.go
	go build -o bin/provider-service cmd/provider-service/main.go
	go build -o bin/lineage-service cmd/lineage-service/main.go
	go build -o bin/reconciliation-service cmd/reconciliation-service/main.go
	go build -o bin/incident-service cmd/incident-service/main.go
	go build -o bin/ai-operator-service cmd/ai-operator-service/main.go
	go build -o bin/notification-service cmd/notification-service/main.go
	go build -o bin/audit-service cmd/audit-service/main.go

# Build Docker images
docker-build:
	docker build -t dcraft-fusion/api-gateway:latest -f cmd/api-gateway/Dockerfile .
	docker build -t dcraft-fusion/resource-service:latest -f cmd/resource-service/Dockerfile .
	# ... similar for other services

# Start local development environment
dev:
	docker-compose up -d
	sleep 5
	make migrate
	air # Hot reload for Go services

# Run database migrations
migrate:
	go run cmd/migrate/main.go up

# Clean
clean:
	docker-compose down -v
	rm -rf bin/
	rm -f coverage.out coverage.html
```

**docker-compose.yml (for local development):**

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: dcraft_fusion
      POSTGRES_USER: dcraft
      POSTGRES_PASSWORD: dcraft
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
  
  neo4j:
    image: neo4j:5
    environment:
      NEO4J_AUTH: neo4j/password
    ports:
      - "7474:7474"  # HTTP
      - "7687:7687"  # Bolt
    volumes:
      - neo4j-data:/data
  
  redis:
    image: redis:7
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
  
  nats:
    image: nats:latest
    ports:
      - "4222:4222"  # Client
      - "8222:8222"  # HTTP monitoring

volumes:
  postgres-data:
  neo4j-data:
  redis-data:
```

---

## Hot Reload with Air

**Why Air?**
- Automatic rebuild and restart on code changes
- Fast feedback loop
- No manual restarts

**.air.toml:**

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main cmd/api-gateway/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

---

## Testing Strategy

### 1. Unit Tests

**Test each function/method in isolation:**

```go
// pkg/resource/service_test.go
func TestCreateDataset(t *testing.T) {
  // Setup
  mockDB := &MockDatabase{}
  service := NewResourceService(mockDB)
  
  // Test
  dataset := &Dataset{
    Name: "test_dataset",
    Owner: "test@example.com",
  }
  
  result, err := service.CreateDataset(dataset)
  
  // Assert
  assert.NoError(t, err)
  assert.NotNil(t, result)
  assert.Equal(t, "test_dataset", result.Name)
}
```

**Run unit tests:**
```bash
go test -v ./...
```

---

### 2. Integration Tests

**Test service interactions:**

```go
// tests/integration/resource_service_test.go
func TestResourceServiceIntegration(t *testing.T) {
  // Setup real database (Docker container)
  db := setupTestDatabase(t)
  defer db.Close()
  
  service := NewResourceService(db)
  
  // Test
  dataset := &Dataset{
    Name: "test_dataset",
    Owner: "test@example.com",
  }
  
  created, err := service.CreateDataset(dataset)
  assert.NoError(t, err)
  
  retrieved, err := service.GetDataset(created.ID)
  assert.NoError(t, err)
  assert.Equal(t, created.ID, retrieved.ID)
}
```

**Run integration tests:**
```bash
go test -v -tags=integration ./tests/integration/...
```

---

### 3. End-to-End Tests

**Test complete user flows:**

```go
// tests/e2e/create_pipeline_test.go
func TestCreatePipelineE2E(t *testing.T) {
  // Setup test environment (all services running)
  client := setupE2EClient(t)
  
  // Test: Create dataset → Create pipeline → Verify lineage
  
  // 1. Create dataset
  dataset := &Dataset{Name: "orders", Owner: "test@example.com"}
  createdDataset, err := client.CreateDataset(dataset)
  assert.NoError(t, err)
  
  // 2. Create pipeline
  pipeline := &Pipeline{
    Name: "transform_orders",
    InputDatasets: []string{createdDataset.ID},
    OutputDatasets: []string{"staging.orders"},
  }
  createdPipeline, err := client.CreatePipeline(pipeline)
  assert.NoError(t, err)
  
  // 3. Verify lineage
  lineage, err := client.GetDownstream(createdDataset.ID)
  assert.NoError(t, err)
  assert.Contains(t, lineage, createdPipeline.ID)
}
```

**Run E2E tests:**
```bash
go test -v -tags=e2e ./tests/e2e/...
```

---

### 4. Load Tests

**Test performance and scalability:**

```go
// tests/load/api_load_test.go
func TestAPILoadTest(t *testing.T) {
  // Setup
  client := setupLoadTestClient(t)
  
  // Load test parameters
  concurrency := 100
  requests := 10000
  
  // Run load test
  results := loadtest.Run(loadtest.Config{
    URL: "http://localhost:8080/api/v1/datasets",
    Method: "GET",
    Concurrency: concurrency,
    Requests: requests,
  })
  
  // Assert
  assert.Less(t, results.AvgLatency, 100*time.Millisecond)
  assert.Greater(t, results.SuccessRate, 0.99) // 99% success rate
}
```

**Run load tests:**
```bash
go test -v -tags=load ./tests/load/...
```

---

## Code Quality Tools

### 1. Linter (golangci-lint)

```bash
# Install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run
golangci-lint run

# Config: .golangci.yml
linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - structcheck
    - varcheck
    - ineffassign
    - deadcode
    - typecheck
```

---

### 2. Code Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Coverage threshold: 80%
```

---

### 3. Security Scanner (gosec)

```bash
# Install
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Run
gosec ./...
```

---

# Security Architecture

## Authentication

### 1. JWT (JSON Web Tokens)

**For API authentication:**

```go
// pkg/auth/jwt.go
type JWTManager struct {
  secretKey string
  tokenDuration time.Duration
}

func (m *JWTManager) Generate(userID string, role string) (string, error) {
  claims := jwt.MapClaims{
    "user_id": userID,
    "role": role,
    "exp": time.Now().Add(m.tokenDuration).Unix(),
  }
  
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) Verify(tokenString string) (*jwt.Token, error) {
  return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("unexpected signing method")
    }
    return []byte(m.secretKey), nil
  })
}
```

---

### 2. OAuth2 / SSO

**For enterprise authentication:**

```go
// pkg/auth/oauth.go
type OAuthProvider struct {
  clientID string
  clientSecret string
  redirectURL string
  provider string // "google", "okta", "azure"
}

func (p *OAuthProvider) GetAuthURL() string {
  // Generate OAuth authorization URL
  return fmt.Sprintf(
    "https://%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code",
    p.provider, p.clientID, p.redirectURL,
  )
}

func (p *OAuthProvider) ExchangeCode(code string) (*Token, error) {
  // Exchange authorization code for access token
  // ...
}
```

---

### 3. API Keys

**For service-to-service authentication:**

```go
// pkg/auth/apikey.go
type APIKeyManager struct {
  db *Database
}

func (m *APIKeyManager) Generate(userID string) (string, error) {
  // Generate random API key
  apiKey := generateRandomString(32)
  
  // Hash API key before storing
  hashedKey := hashAPIKey(apiKey)
  
  // Store in database
  err := m.db.StoreAPIKey(userID, hashedKey)
  if err != nil {
    return "", err
  }
  
  return apiKey, nil
}

func (m *APIKeyManager) Verify(apiKey string) (string, error) {
  // Hash API key
  hashedKey := hashAPIKey(apiKey)
  
  // Look up in database
  userID, err := m.db.GetUserIDByAPIKey(hashedKey)
  if err != nil {
    return "", err
  }
  
  return userID, nil
}
```

---

## Authorization (RBAC)

**Role-Based Access Control:**

```go
// pkg/auth/rbac.go
type Role string

const (
  RoleAdmin Role = "admin"
  RoleEditor Role = "editor"
  RoleViewer Role = "viewer"
)

type Permission string

const (
  PermissionReadDatasets Permission = "read:datasets"
  PermissionWriteDatasets Permission = "write:datasets"
  PermissionDeleteDatasets Permission = "delete:datasets"
  PermissionReadPipelines Permission = "read:pipelines"
  PermissionWritePipelines Permission = "write:pipelines"
  PermissionDeletePipelines Permission = "delete:pipelines"
)

var rolePermissions = map[Role][]Permission{
  RoleAdmin: {
    PermissionReadDatasets,
    PermissionWriteDatasets,
    PermissionDeleteDatasets,
    PermissionReadPipelines,
    PermissionWritePipelines,
    PermissionDeletePipelines,
  },
  RoleEditor: {
    PermissionReadDatasets,
    PermissionWriteDatasets,
    PermissionReadPipelines,
    PermissionWritePipelines,
  },
  RoleViewer: {
    PermissionReadDatasets,
    PermissionReadPipelines,
  },
}

func (r Role) HasPermission(permission Permission) bool {
  permissions, ok := rolePermissions[r]
  if !ok {
    return false
  }
  
  for _, p := range permissions {
    if p == permission {
      return true
    }
  }
  
  return false
}
```

**Middleware:**

```go
// pkg/middleware/auth.go
func AuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
  return func(c *gin.Context) {
    // Get token from Authorization header
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
      c.JSON(401, gin.H{"error": "missing authorization header"})
      c.Abort()
      return
    }
    
    // Verify token
    token, err := jwtManager.Verify(strings.TrimPrefix(authHeader, "Bearer "))
    if err != nil {
      c.JSON(401, gin.H{"error": "invalid token"})
      c.Abort()
      return
    }
    
    // Extract claims
    claims := token.Claims.(jwt.MapClaims)
    c.Set("user_id", claims["user_id"])
    c.Set("role", claims["role"])
    
    c.Next()
  }
}

func RequirePermission(permission Permission) gin.HandlerFunc {
  return func(c *gin.Context) {
    // Get role from context
    roleStr, _ := c.Get("role")
    role := Role(roleStr.(string))
    
    // Check permission
    if !role.HasPermission(permission) {
      c.JSON(403, gin.H{"error": "insufficient permissions"})
      c.Abort()
      return
    }
    
    c.Next()
  }
}
```

**Usage:**

```go
// cmd/api-gateway/main.go
router := gin.Default()

// Public endpoints
router.POST("/auth/login", handleLogin)
router.POST("/auth/register", handleRegister)

// Protected endpoints
authorized := router.Group("/api/v1")
authorized.Use(AuthMiddleware(jwtManager))
{
  // Read-only endpoints (viewer, editor, admin)
  authorized.GET("/datasets", RequirePermission(PermissionReadDatasets), handleListDatasets)
  authorized.GET("/datasets/:id", RequirePermission(PermissionReadDatasets), handleGetDataset)
  
  // Write endpoints (editor, admin)
  authorized.POST("/datasets", RequirePermission(PermissionWriteDatasets), handleCreateDataset)
  authorized.PUT("/datasets/:id", RequirePermission(PermissionWriteDatasets), handleUpdateDataset)
  
  // Delete endpoints (admin only)
  authorized.DELETE("/datasets/:id", RequirePermission(PermissionDeleteDatasets), handleDeleteDataset)
}
```

---

## Secrets Management

### 1. Kubernetes Secrets

**Store sensitive data in Kubernetes Secrets:**

```yaml
# secrets.yaml (encrypted with sealed-secrets or SOPS)
apiVersion: v1
kind: Secret
metadata:
  name: dcraft-fusion-secrets
type: Opaque
stringData:
  postgres-password: <encrypted>
  neo4j-password: <encrypted>
  redis-password: <encrypted>
  jwt-secret-key: <encrypted>
  openai-api-key: <encrypted>
  anthropic-api-key: <encrypted>
  slack-webhook-url: <encrypted>
```

---

### 2. External Secrets Operator

**Sync secrets from external secret managers (AWS Secrets Manager, HashiCorp Vault, etc.):**

```yaml
# external-secret.yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: dcraft-fusion-secrets
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: SecretStore
  target:
    name: dcraft-fusion-secrets
    creationPolicy: Owner
  data:
    - secretKey: postgres-password
      remoteRef:
        key: dcraft-fusion/postgres-password
    - secretKey: openai-api-key
      remoteRef:
        key: dcraft-fusion/openai-api-key
```

---

## Network Security

### 1. Network Policies

**Restrict network traffic between pods:**

```yaml
# network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: api-gateway-policy
spec:
  podSelector:
    matchLabels:
      app: api-gateway
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: ingress-controller
      ports:
        - protocol: TCP
          port: 8080
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: resource-service
      ports:
        - protocol: TCP
          port: 50051
    - to:
        - podSelector:
            matchLabels:
              app: provider-service
      ports:
        - protocol: TCP
          port: 50051
```

---

### 2. TLS/SSL

**Encrypt all traffic:**

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: dcraft-fusion-ingress
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
    - hosts:
        - api.dcraft-fusion.example.com
      secretName: dcraft-fusion-tls
  rules:
    - host: api.dcraft-fusion.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-gateway
                port:
                  number: 8080
```

---

## Audit Logging

**Log all important actions:**

```go
// pkg/audit/logger.go
type AuditLogger struct {
  db *Database
}

type AuditLog struct {
  ID string
  Timestamp time.Time
  UserID string
  Action string
  ResourceType string
  ResourceID string
  Changes map[string]interface{}
  IPAddress string
  UserAgent string
}

func (l *AuditLogger) Log(log *AuditLog) error {
  return l.db.InsertAuditLog(log)
}
```

**Middleware:**

```go
// pkg/middleware/audit.go
func AuditMiddleware(auditLogger *AuditLogger) gin.HandlerFunc {
  return func(c *gin.Context) {
    // Capture request
    userID, _ := c.Get("user_id")
    
    // Execute request
    c.Next()
    
    // Log audit event
    if c.Request.Method != "GET" { // Only log write operations
      auditLogger.Log(&AuditLog{
        Timestamp: time.Now(),
        UserID: userID.(string),
        Action: c.Request.Method,
        ResourceType: extractResourceType(c.Request.URL.Path),
        ResourceID: extractResourceID(c.Request.URL.Path),
        IPAddress: c.ClientIP(),
        UserAgent: c.Request.UserAgent(),
      })
    }
  }
}
```

---

# Scalability & Performance

## Horizontal Scaling

**All services are stateless and can scale horizontally:**

```yaml
# HPA (Horizontal Pod Autoscaler)
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-gateway-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api-gateway
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

---

## Database Scaling

### 1. PostgreSQL

**Read replicas for read-heavy workloads:**

```
┌─────────────────────────────────────────────────────────────┐
│  PostgreSQL Cluster                                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Primary (Read/Write)                                        │
│  ↓ (replication)                                             │
│  Read Replica 1 (Read-only)                                  │
│  Read Replica 2 (Read-only)                                  │
│                                                              │
│  Connection Pool:                                            │
│  - Write queries → Primary                                   │
│  - Read queries → Read Replicas (round-robin)                │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

### 2. Neo4j

**Read replicas for lineage queries:**

```
┌─────────────────────────────────────────────────────────────┐
│  Neo4j Cluster                                               │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Core 1 (Read/Write)                                         │
│  ↓ (replication)                                             │
│  Read Replica 1 (Read-only)                                  │
│  Read Replica 2 (Read-only)                                  │
│                                                              │
│  Connection Pool:                                            │
│  - Write queries → Core                                      │
│  - Read queries → Read Replicas (round-robin)                │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

### 3. Redis

**Redis Cluster for high availability:**

```
┌─────────────────────────────────────────────────────────────┐
│  Redis Cluster                                               │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Primary (Read/Write)                                        │
│  ↓ (replication)                                             │
│  Replica 1 (Read-only)                                       │
│  Replica 2 (Read-only)                                       │
│                                                              │
│  Sentinel (monitors health, automatic failover)              │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Caching Strategy

### 1. Multi-Level Caching

```
┌─────────────────────────────────────────────────────────────┐
│  Caching Layers                                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  L1: In-Memory Cache (Go map, 1 minute TTL)                  │
│  ↓ (miss)                                                    │
│  L2: Redis Cache (1 hour TTL)                                │
│  ↓ (miss)                                                    │
│  L3: Database (PostgreSQL, Neo4j)                            │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Implementation:**

```go
// pkg/cache/multilevel.go
type MultiLevelCache struct {
  l1 *sync.Map // In-memory cache
  l2 *redis.Client // Redis cache
  db *Database // Database
}

func (c *MultiLevelCache) Get(key string) (interface{}, error) {
  // L1: In-memory cache
  if value, ok := c.l1.Load(key); ok {
    return value, nil
  }
  
  // L2: Redis cache
  value, err := c.l2.Get(key).Result()
  if err == nil {
    // Store in L1
    c.l1.Store(key, value)
    return value, nil
  }
  
  // L3: Database
  value, err = c.db.Get(key)
  if err != nil {
    return nil, err
  }
  
  // Store in L2 and L1
  c.l2.Set(key, value, 1*time.Hour)
  c.l1.Store(key, value)
  
  return value, nil
}
```

---

## Connection Pooling

**Reuse database connections:**

```go
// pkg/database/pool.go
type ConnectionPool struct {
  db *sql.DB
}

func NewConnectionPool(dsn string) (*ConnectionPool, error) {
  db, err := sql.Open("postgres", dsn)
  if err != nil {
    return nil, err
  }
  
  // Configure connection pool
  db.SetMaxOpenConns(100) // Max open connections
  db.SetMaxIdleConns(10)  // Max idle connections
  db.SetConnMaxLifetime(1 * time.Hour) // Max connection lifetime
  
  return &ConnectionPool{db: db}, nil
}
```

---

## Load Testing Results (Expected)

**Target: 10,000 requests/second**

| Metric | Target | Expected |
|--------|--------|----------|
| Avg Latency | <100ms | 50ms |
| P95 Latency | <200ms | 150ms |
| P99 Latency | <500ms | 300ms |
| Success Rate | >99.9% | 99.95% |
| Throughput | 10k req/s | 15k req/s |

---

# Monitoring & Observability

## Metrics (Prometheus)

**Collect metrics from all services:**

```go
// pkg/metrics/prometheus.go
var (
  httpRequestsTotal = prometheus.NewCounterVec(
    prometheus.CounterOpts{
      Name: "http_requests_total",
      Help: "Total number of HTTP requests",
    },
    []string{"method", "path", "status"},
  )
  
  httpRequestDuration = prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
      Name: "http_request_duration_seconds",
      Help: "HTTP request duration in seconds",
      Buckets: prometheus.DefBuckets,
    },
    []string{"method", "path"},
  )
  
  grpcRequestsTotal = prometheus.NewCounterVec(
    prometheus.CounterOpts{
      Name: "grpc_requests_total",
      Help: "Total number of gRPC requests",
    },
    []string{"method", "status"},
  )
  
  grpcRequestDuration = prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
      Name: "grpc_request_duration_seconds",
      Help: "gRPC request duration in seconds",
      Buckets: prometheus.DefBuckets,
    },
    []string{"method"},
  )
)

func init() {
  prometheus.MustRegister(httpRequestsTotal)
  prometheus.MustRegister(httpRequestDuration)
  prometheus.MustRegister(grpcRequestsTotal)
  prometheus.MustRegister(grpcRequestDuration)
}
```

**Middleware:**

```go
// pkg/middleware/metrics.go
func MetricsMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    start := time.Now()
    
    c.Next()
    
    duration := time.Since(start).Seconds()
    status := strconv.Itoa(c.Writer.Status())
    
    httpRequestsTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, status).Inc()
    httpRequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(duration)
  }
}
```

---

## Dashboards (Grafana)

**Pre-built dashboards:**

1. **Platform Overview**
   - Total requests/second
   - Avg latency
   - Error rate
   - Active users

2. **Service Health**
   - CPU usage per service
   - Memory usage per service
   - Request rate per service
   - Error rate per service

3. **Database Performance**
   - Query latency
   - Connection pool usage
   - Slow queries
   - Cache hit rate

4. **AI Operator**
   - LLM API calls/hour
   - LLM API cost/hour
   - Token usage
   - Cache hit rate

5. **Business Metrics**
   - Datasets tracked
   - Pipelines tracked
   - Incidents detected
   - Drift detected

---

## Distributed Tracing (Jaeger)

**Trace requests across services:**

```go
// pkg/tracing/jaeger.go
func InitTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
  cfg := &config.Configuration{
    ServiceName: serviceName,
    Sampler: &config.SamplerConfig{
      Type: "const",
      Param: 1, // Sample 100% of traces
    },
    Reporter: &config.ReporterConfig{
      LogSpans: true,
      LocalAgentHostPort: "jaeger:6831",
    },
  }
  
  tracer, closer, err := cfg.NewTracer()
  if err != nil {
    return nil, nil, err
  }
  
  opentracing.SetGlobalTracer(tracer)
  
  return tracer, closer, nil
}
```

**Usage:**

```go
// cmd/api-gateway/main.go
func main() {
  tracer, closer, err := InitTracer("api-gateway")
  if err != nil {
    log.Fatal(err)
  }
  defer closer.Close()
  
  // ... rest of main
}

// pkg/resource/service.go
func (s *ResourceService) GetDataset(ctx context.Context, id string) (*Dataset, error) {
  span, ctx := opentracing.StartSpanFromContext(ctx, "GetDataset")
  defer span.Finish()
  
  span.SetTag("dataset.id", id)
  
  // ... implementation
}
```

---

## Logging (ELK Stack)

**Structured logging:**

```go
// pkg/logging/logger.go
import "go.uber.org/zap"

var logger *zap.Logger

func InitLogger() error {
  var err error
  logger, err = zap.NewProduction()
  if err != nil {
    return err
  }
  return nil
}

func Info(msg string, fields ...zap.Field) {
  logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
  logger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
  logger.Debug(msg, fields...)
}
```

**Usage:**

```go
// pkg/resource/service.go
func (s *ResourceService) CreateDataset(dataset *Dataset) (*Dataset, error) {
  logging.Info("Creating dataset",
    zap.String("name", dataset.Name),
    zap.String("owner", dataset.Owner),
  )
  
  result, err := s.db.InsertDataset(dataset)
  if err != nil {
    logging.Error("Failed to create dataset",
      zap.String("name", dataset.Name),
      zap.Error(err),
    )
    return nil, err
  }
  
  logging.Info("Dataset created successfully",
    zap.String("id", result.ID),
    zap.String("name", result.Name),
  )
  
  return result, nil
}
```

---

## Alerting (Prometheus Alertmanager)

**Alert rules:**

```yaml
# prometheus-alerts.yaml
groups:
  - name: dcraft-fusion
    interval: 30s
    rules:
      # High error rate
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} (threshold: 0.05)"
      
      # High latency
      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 0.5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "P95 latency is {{ $value }}s (threshold: 0.5s)"
      
      # Database connection pool exhausted
      - alert: DatabaseConnectionPoolExhausted
        expr: pg_stat_database_numbackends / pg_settings_max_connections > 0.9
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Database connection pool exhausted"
          description: "Connection pool usage is {{ $value }} (threshold: 0.9)"
      
      # AI API budget exceeded
      - alert: AIAPIBudgetExceeded
        expr: sum(rate(ai_api_cost_dollars[1h])) * 24 > 100
        for: 1h
        labels:
          severity: warning
        annotations:
          summary: "AI API budget exceeded"
          description: "Daily cost is ${{ $value }} (budget: $100)"
```

---

# Cost Optimization

## Infrastructure Cost Breakdown (Estimated)

**For 100 customers, 10,000 datasets, 1,000 pipelines:**

| Component | Configuration | Monthly Cost |
|-----------|---------------|--------------|
| **Kubernetes Cluster** | 3 nodes (8 CPU, 32GB RAM each) | $500 |
| **PostgreSQL** | 1 primary + 2 replicas (4 CPU, 16GB RAM each) | $400 |
| **Neo4j** | 1 primary + 2 replicas (4 CPU, 16GB RAM each) | $400 |
| **Redis** | 1 primary + 2 replicas (2 CPU, 8GB RAM each) | $200 |
| **Load Balancer** | 1 load balancer | $50 |
| **Storage** | 500GB SSD | $100 |
| **Backup** | 1TB S3 storage | $25 |
| **Monitoring** | Prometheus + Grafana + Jaeger | $100 |
| **AI API** | OpenAI + Anthropic (10k queries/month) | $500 |
| **Total** | | **$2,275/month** |

**Cost per customer:** $22.75/month

**Target pricing:** $50-100/customer/month

**Gross margin:** 55-78%

---

## Cost Optimization Strategies

### 1. Right-Size Resources

**Start small, scale up as needed:**

```yaml
# Development environment
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi

# Production environment
resources:
  requests:
    cpu: 500m
    memory: 512Mi
  limits:
    cpu: 1000m
    memory: 1Gi
```

---

### 2. Use Spot Instances (AWS)

**Save 70% on compute costs:**

```yaml
# nodepool-spot.yaml
apiVersion: v1
kind: NodePool
metadata:
  name: spot-pool
spec:
  instanceTypes:
    - t3.large
    - t3a.large
  capacityType: SPOT
  minSize: 1
  maxSize: 10
  labels:
    workload: stateless
```

**Taint spot nodes:**

```yaml
taints:
  - key: spot
    value: "true"
    effect: NoSchedule
```

**Tolerate spot nodes for stateless services:**

```yaml
tolerations:
  - key: spot
    operator: Equal
    value: "true"
    effect: NoSchedule
```

---

### 3. Auto-Scaling

**Scale down during off-hours:**

```yaml
# HPA with schedule
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-gateway-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api-gateway
  minReplicas: 2  # Business hours
  maxReplicas: 10
  # Scale down to 1 replica during off-hours (via CronJob)
```

---

### 4. Database Cost Optimization

**Use managed databases with auto-pause:**

- AWS RDS: Auto-pause after 5 minutes of inactivity
- Azure Database: Serverless tier with auto-pause
- GCP Cloud SQL: Automatic storage increase

---

### 5. AI API Cost Optimization

**Reduce LLM API costs by 60%:**

1. **Response caching:** Cache identical queries (60% cache hit rate)
2. **Model selection:** Use cheaper models when possible (GPT-3.5 vs GPT-4)
3. **Token optimization:** Remove redundant context
4. **Batch processing:** Batch multiple queries

**Expected savings:** $500/month → $200/month

---

# Summary

## Technology Stack (Final)

```
Backend:           Go 1.21+
Frontend:          React 18 + TypeScript
CLI:               Go (Cobra)
API:               REST + GraphQL + gRPC
Databases:         PostgreSQL, Neo4j, Redis, TimescaleDB
Message Queue:     NATS
Container:         Docker
Orchestration:     Kubernetes
Deployment:        Helm + ArgoCD (GitOps)
CI/CD:             GitHub Actions
AI/LLM:            OpenAI + Anthropic Claude
Monitoring:        Prometheus + Grafana + Jaeger
Logging:           ELK Stack
```

---

## Key Decisions

1. **Go for backend** - Performance, concurrency, ecosystem
2. **gRPC for internal communication** - Performance, type safety
3. **REST + GraphQL for external API** - Compatibility, flexibility
4. **PostgreSQL + Neo4j** - Polyglot persistence (right tool for right job)
5. **Kubernetes + Helm + ArgoCD** - Industry standard, GitOps
6. **Multi-level caching** - Performance, cost optimization
7. **Horizontal scaling** - Stateless services, auto-scaling
8. **Comprehensive observability** - Metrics, traces, logs

---

## Next Steps

1. **Review and approve this technical architecture**
2. **Set up development environment**
3. **Create project structure and boilerplate**
4. **Implement core services (Resource, Provider, Lineage)**
5. **Set up CI/CD pipeline**
6. **Deploy to development environment**
7. **Start V0 development**

---

**This technical architecture is production-ready and scalable to 100+ customers and 10,000+ datasets.**