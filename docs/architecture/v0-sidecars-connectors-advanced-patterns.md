# DCraft Fusion V0: Sidecars, Connectors & Advanced Kubernetes Patterns
## Complete Implementation Guide

## Table of Contents
1. [Sidecar Pattern Strategy](#sidecar-pattern-strategy)
2. [Provider Connector Architecture](#provider-connector-architecture)
3. [Connector Lifecycle Management](#connector-lifecycle-management)
4. [Dynamic Connector Loading](#dynamic-connector-loading)
5. [Kubernetes Operator Pattern](#kubernetes-operator-pattern)
6. [Advanced Kubernetes Patterns](#advanced-kubernetes-patterns)
7. [Complete Implementation Examples](#complete-implementation-examples)

---

# Sidecar Pattern Strategy

## Overview

**Sidecar containers** run alongside the main application container in the same Pod, sharing network namespace, storage volumes, and lifecycle.

**Key Decision:** Use sidecars **sparingly** - they add resource overhead and complexity.

---

## Which Services Get Sidecars?

### ✅ Services WITH Sidecars

#### 1. **All Application Services** → **Linkerd Proxy Sidecar**

**Sidecar:** Linkerd2-proxy (Rust-based micro-proxy)

**Purpose:**
- Automatic mTLS (encrypt all service-to-service traffic)
- Automatic retries (exponential backoff)
- Circuit breaking (prevent cascading failures)
- Load balancing (intelligent request routing)
- Observability (metrics, traces)

**Resource Usage:**
- CPU: 10m (request), 100m (limit)
- Memory: 20Mi (request), 100Mi (limit)

**Configuration:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  namespace: dcraft-fusion
  annotations:
    linkerd.io/inject: enabled  # Auto-inject Linkerd sidecar
spec:
  template:
    metadata:
      annotations:
        linkerd.io/inject: enabled
    spec:
      containers:
      - name: api-gateway
        image: registry.dcraft-fusion.io/api-gateway:v0.1.0
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 50051
          name: grpc
```

**How it works:**
1. Linkerd proxy injector automatically adds sidecar container to Pod
2. Proxy intercepts all inbound/outbound traffic
3. Proxy handles mTLS, retries, circuit breaking, metrics
4. Main container is unaware of proxy (zero code changes)

---

#### 2. **Provider Service** → **Provider Connector Sidecars**

**Sidecars:** One sidecar per active provider (dbt, Snowflake, BigQuery, etc.)

**Purpose:**
- Isolate provider-specific logic
- Enable hot-swapping of providers
- Prevent provider failures from crashing main service
- Allow independent scaling of providers

**Architecture:**
```
┌─────────────────────────────────────────────────────────────┐
│                    Provider Service Pod                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  Main Container: Provider Service (Go)             │     │
│  │  - Provider orchestration                          │     │
│  │  - Provider lifecycle management                   │     │
│  │  - gRPC API server                                 │     │
│  └────────────────┬───────────────────────────────────┘     │
│                   │ (localhost gRPC)                         │
│         ┌─────────┴─────────┬─────────────────┐             │
│         │                   │                 │             │
│         ▼                   ▼                 ▼             │
│  ┌─────────────┐     ┌─────────────┐   ┌─────────────┐     │
│  │  dbt        │     │ Snowflake   │   │ BigQuery    │     │
│  │  Connector  │     │ Connector   │   │ Connector   │     │
│  │  Sidecar    │     │ Sidecar     │   │ Sidecar     │     │
│  │             │     │             │   │             │     │
│  │  Port: 50061│     │ Port: 50062 │   │ Port: 50063 │     │
│  └─────────────┘     └─────────────┘   └─────────────┘     │
│         │                   │                 │             │
│         │                   │                 │             │
│         ▼                   ▼                 ▼             │
│  ┌─────────────┐     ┌─────────────┐   ┌─────────────┐     │
│  │   dbt       │     │  Snowflake  │   │  BigQuery   │     │
│  │   API       │     │  API        │   │  API        │     │
│  └─────────────┘     └─────────────┘   └─────────────┘     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Deployment:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: provider-service
  namespace: dcraft-fusion
spec:
  template:
    spec:
      containers:
      # Main container
      - name: provider-service
        image: registry.dcraft-fusion.io/provider-service:v0.1.0
        ports:
        - containerPort: 50051
          name: grpc
        env:
        - name: DBT_CONNECTOR_URL
          value: "localhost:50061"
        - name: SNOWFLAKE_CONNECTOR_URL
          value: "localhost:50062"
        - name: BIGQUERY_CONNECTOR_URL
          value: "localhost:50063"
      
      # dbt connector sidecar
      - name: dbt-connector
        image: registry.dcraft-fusion.io/connector-dbt:v0.1.0
        ports:
        - containerPort: 50061
          name: grpc
        env:
        - name: DBT_PROFILES_DIR
          value: /etc/dbt
        volumeMounts:
        - name: dbt-profiles
          mountPath: /etc/dbt
      
      # Snowflake connector sidecar
      - name: snowflake-connector
        image: registry.dcraft-fusion.io/connector-snowflake:v0.1.0
        ports:
        - containerPort: 50062
          name: grpc
        env:
        - name: SNOWFLAKE_ACCOUNT
          valueFrom:
            secretKeyRef:
              name: snowflake-credentials
              key: account
      
      # BigQuery connector sidecar
      - name: bigquery-connector
        image: registry.dcraft-fusion.io/connector-bigquery:v0.1.0
        ports:
        - containerPort: 50063
          name: grpc
        env:
        - name: GOOGLE_APPLICATION_CREDENTIALS
          value: /etc/gcp/service-account.json
        volumeMounts:
        - name: gcp-credentials
          mountPath: /etc/gcp
      
      volumes:
      - name: dbt-profiles
        configMap:
          name: dbt-profiles
      - name: gcp-credentials
        secret:
          secretName: gcp-service-account
```

**Why sidecars for connectors?**
1. **Isolation:** Connector failure doesn't crash main service
2. **Hot-swapping:** Add/remove connectors without restarting main service
3. **Independent scaling:** Scale connectors independently
4. **Language flexibility:** Connectors can be in any language (Python for dbt, Go for Snowflake, etc.)
5. **Resource management:** Set different resource limits per connector

---

#### 3. **AI Operator Service** → **Prompt Cache Sidecar**

**Sidecar:** Redis sidecar for local prompt caching

**Purpose:**
- Ultra-low latency cache access (<1ms)
- Reduce load on central Redis cluster
- Improve AI response times

**Configuration:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-operator-service
spec:
  template:
    spec:
      containers:
      # Main container
      - name: ai-operator-service
        image: registry.dcraft-fusion.io/ai-operator-service:v0.1.0
        env:
        - name: REDIS_URL
          value: "localhost:6379"  # Local Redis sidecar
      
      # Redis cache sidecar
      - name: redis-cache
        image: redis:7-alpine
        ports:
        - containerPort: 6379
        resources:
          requests:
            cpu: 50m
            memory: 128Mi
          limits:
            cpu: 200m
            memory: 512Mi
        volumeMounts:
        - name: redis-data
          mountPath: /data
      
      volumes:
      - name: redis-data
        emptyDir: {}  # Ephemeral storage (cache only)
```

**Why local Redis sidecar?**
- **Latency:** <1ms vs 5-10ms for central Redis
- **Availability:** Cache survives if central Redis fails
- **Cost:** Reduces load on central Redis cluster

---

### ❌ Services WITHOUT Sidecars (Except Linkerd)

#### 1. **API Gateway** → No additional sidecars
- Already has Linkerd proxy
- No need for additional sidecars

#### 2. **Resource Service** → No additional sidecars
- Simple CRUD operations
- No need for additional sidecars

#### 3. **Lineage Service** → No additional sidecars
- Direct Neo4j connection
- No need for additional sidecars

#### 4. **Reconciliation Service** → No additional sidecars
- Orchestration logic only
- No need for additional sidecars

#### 5. **Incident Service** → No additional sidecars
- Event handling only
- No need for additional sidecars

#### 6. **Notification Service** → No additional sidecars
- Simple API calls
- No need for additional sidecars

#### 7. **Audit Service** → No additional sidecars
- Simple logging
- No need for additional sidecars

---

## Sidecar Pattern Summary

| Service | Linkerd Proxy | Additional Sidecars | Total Containers |
|---------|---------------|---------------------|------------------|
| API Gateway | ✅ | None | 2 (main + proxy) |
| Resource Service | ✅ | None | 2 (main + proxy) |
| **Provider Service** | ✅ | **3 connectors** | **5 (main + proxy + 3 connectors)** |
| Lineage Service | ✅ | None | 2 (main + proxy) |
| Reconciliation Service | ✅ | None | 2 (main + proxy) |
| Incident Service | ✅ | None | 2 (main + proxy) |
| **AI Operator Service** | ✅ | **Redis cache** | **3 (main + proxy + cache)** |
| Notification Service | ✅ | None | 2 (main + proxy) |
| Audit Service | ✅ | None | 2 (main + proxy) |

**Total sidecars per service:**
- 8 services with 2 containers (main + Linkerd proxy)
- 1 service with 5 containers (Provider Service: main + proxy + 3 connectors)
- 1 service with 3 containers (AI Operator: main + proxy + Redis cache)

**Total containers: 24** (10 main + 10 Linkerd proxies + 3 connector sidecars + 1 Redis sidecar)

---

## Init Containers (Not Sidecars)

**Init containers** run **before** the main container starts, then exit.

### Use Cases in DCraft Fusion:

#### 1. **Database Migration Init Container**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: resource-service
spec:
  template:
    spec:
      initContainers:
      # Run database migrations before starting main container
      - name: db-migrate
        image: registry.dcraft-fusion.io/resource-service:v0.1.0
        command: ["/app/migrate", "up"]
        env:
        - name: POSTGRES_HOST
          value: postgresql
        - name: POSTGRES_PORT
          value: "5432"
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: postgresql-credentials
              key: username
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgresql-credentials
              key: password
      
      containers:
      - name: resource-service
        image: registry.dcraft-fusion.io/resource-service:v0.1.0
```

#### 2. **Wait for Dependencies Init Container**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
spec:
  template:
    spec:
      initContainers:
      # Wait for PostgreSQL to be ready
      - name: wait-for-postgres
        image: busybox:1.36
        command: ['sh', '-c', 'until nc -z postgresql 5432; do echo waiting for postgres; sleep 2; done']
      
      # Wait for Redis to be ready
      - name: wait-for-redis
        image: busybox:1.36
        command: ['sh', '-c', 'until nc -z redis 6379; do echo waiting for redis; sleep 2; done']
      
      # Wait for NATS to be ready
      - name: wait-for-nats
        image: busybox:1.36
        command: ['sh', '-c', 'until nc -z nats 4222; do echo waiting for nats; sleep 2; done']
      
      containers:
      - name: api-gateway
        image: registry.dcraft-fusion.io/api-gateway:v0.1.0
```

#### 3. **Configuration Fetching Init Container**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: provider-service
spec:
  template:
    spec:
      initContainers:
      # Fetch provider configurations from config server
      - name: fetch-config
        image: curlimages/curl:8.5.0
        command:
        - sh
        - -c
        - |
          curl -o /config/providers.yaml https://config-server/providers.yaml
        volumeMounts:
        - name: config
          mountPath: /config
      
      containers:
      - name: provider-service
        image: registry.dcraft-fusion.io/provider-service:v0.1.0
        volumeMounts:
        - name: config
          mountPath: /etc/config
      
      volumes:
      - name: config
        emptyDir: {}
```

---

## DaemonSets (Cluster-Wide Services)

**DaemonSets** ensure a Pod runs on **every node** (or selected nodes).

### Use Cases in DCraft Fusion:

#### 1. **Log Collection DaemonSet (Fluent Bit)**

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluent-bit
  namespace: logging
spec:
  selector:
    matchLabels:
      app: fluent-bit
  template:
    metadata:
      labels:
        app: fluent-bit
    spec:
      containers:
      - name: fluent-bit
        image: fluent/fluent-bit:2.2
        volumeMounts:
        # Mount host's /var/log to collect logs
        - name: varlog
          mountPath: /var/log
          readOnly: true
        # Mount host's /var/lib/docker/containers to collect container logs
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
        # Mount Fluent Bit config
        - name: fluent-bit-config
          mountPath: /fluent-bit/etc/
      
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers
      - name: fluent-bit-config
        configMap:
          name: fluent-bit-config
```

**Why DaemonSet?**
- Collects logs from **all nodes**
- More efficient than sidecar per Pod (1 Pod per node vs 1 sidecar per Pod)

---

#### 2. **Node Monitoring DaemonSet (Prometheus Node Exporter)**

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-exporter
  namespace: monitoring
spec:
  selector:
    matchLabels:
      app: node-exporter
  template:
    metadata:
      labels:
        app: node-exporter
    spec:
      hostNetwork: true
      hostPID: true
      containers:
      - name: node-exporter
        image: prom/node-exporter:v1.7.0
        args:
        - --path.procfs=/host/proc
        - --path.sysfs=/host/sys
        - --path.rootfs=/host/root
        - --collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)
        ports:
        - containerPort: 9100
          name: metrics
        volumeMounts:
        - name: proc
          mountPath: /host/proc
          readOnly: true
        - name: sys
          mountPath: /host/sys
          readOnly: true
        - name: root
          mountPath: /host/root
          readOnly: true
      
      volumes:
      - name: proc
        hostPath:
          path: /proc
      - name: sys
        hostPath:
          path: /sys
      - name: root
        hostPath:
          path: /
```

**Why DaemonSet?**
- Monitors **all nodes**
- More efficient than sidecar per Pod

---

# Provider Connector Architecture

## Overview

**Provider connectors** are plugins that extract metadata and lineage from external data tools (dbt, Snowflake, BigQuery, Airflow, etc.).

**Key Requirements:**
1. **Pluggable:** Add/remove connectors without code changes
2. **Isolated:** Connector failure doesn't crash Provider Service
3. **Scalable:** Scale connectors independently
4. **Dynamic:** Load/unload connectors at runtime (hot reload)
5. **Versioned:** Support multiple versions of same connector

---

## Connector Types

### 1. **Built-in Connectors (V0)**

**Included in Provider Service:**
- dbt (manifest.json parser)
- Snowflake (information_schema queries)
- BigQuery (information_schema queries)
- PostgreSQL (information_schema queries)
- Redshift (information_schema queries)
- Generic SQL (information_schema queries)

**Deployment:** Sidecar containers

---

### 2. **Community Connectors (V1)**

**Contributed by community:**
- Airflow (DAG parser)
- Kafka (topic metadata)
- Fivetran (API integration)
- Airbyte (API integration)
- Spark (job metadata)
- Great Expectations (expectation suites)
- Monte Carlo (monitors)

**Deployment:** Sidecar containers or standalone Pods

---

### 3. **Custom Connectors (V1)**

**Built by users:**
- Proprietary data tools
- Internal data platforms
- Custom data pipelines

**Deployment:** Sidecar containers or standalone Pods

---

## Connector Interface (gRPC)

**All connectors implement the same gRPC interface:**

```protobuf
// proto/provider/v1/connector.proto
syntax = "proto3";

package dcraft.fusion.provider.v1;

// Connector service interface
service Connector {
  // Get connector capabilities
  rpc GetCapabilities(GetCapabilitiesRequest) returns (Capabilities);
  
  // Extract metadata from data source
  rpc ExtractMetadata(ExtractMetadataRequest) returns (ExtractMetadataResponse);
  
  // Extract lineage from data source
  rpc ExtractLineage(ExtractLineageRequest) returns (ExtractLineageResponse);
  
  // Health check
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
}

// Connector capabilities
message Capabilities {
  string name = 1;           // e.g., "dbt"
  string version = 2;        // e.g., "1.0.0"
  repeated string operations = 3;  // e.g., ["extract_metadata", "extract_lineage"]
  map<string, string> config_schema = 4;  // JSON schema for configuration
}

// Extract metadata request
message ExtractMetadataRequest {
  map<string, string> config = 1;  // Connector-specific configuration
  repeated string resource_types = 2;  // e.g., ["dataset", "pipeline"]
}

// Extract metadata response
message ExtractMetadataResponse {
  repeated Dataset datasets = 1;
  repeated Pipeline pipelines = 2;
  map<string, string> metadata = 3;
}

// Extract lineage request
message ExtractLineageRequest {
  map<string, string> config = 1;  // Connector-specific configuration
}

// Extract lineage response
message ExtractLineageResponse {
  repeated LineageEdge edges = 1;
}

message LineageEdge {
  string source_id = 1;
  string target_id = 2;
  string edge_type = 3;  // "dataset_to_dataset", "pipeline_to_dataset"
  map<string, string> metadata = 4;
}

// Health check
message HealthCheckRequest {}

message HealthCheckResponse {
  enum Status {
    UNKNOWN = 0;
    HEALTHY = 1;
    UNHEALTHY = 2;
  }
  Status status = 1;
  string message = 2;
}
```

---

## Connector Implementation Example (dbt Connector in Go)

```go
// connectors/dbt/main.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net"
    "os"
    
    pb "github.com/dcraft-fusion/proto/provider/v1"
    "google.golang.org/grpc"
)

type DbtConnector struct {
    pb.UnimplementedConnectorServer
    manifestPath string
}

func NewDbtConnector() *DbtConnector {
    return &DbtConnector{
        manifestPath: os.Getenv("DBT_MANIFEST_PATH"),
    }
}

func (c *DbtConnector) GetCapabilities(ctx context.Context, req *pb.GetCapabilitiesRequest) (*pb.Capabilities, error) {
    return &pb.Capabilities{
        Name:    "dbt",
        Version: "1.0.0",
        Operations: []string{
            "extract_metadata",
            "extract_lineage",
        },
        ConfigSchema: map[string]string{
            "manifest_path": "string (required): Path to dbt manifest.json",
            "run_results_path": "string (optional): Path to dbt run_results.json",
        },
    }, nil
}

func (c *DbtConnector) ExtractMetadata(ctx context.Context, req *pb.ExtractMetadataRequest) (*pb.ExtractMetadataResponse, error) {
    // 1. Parse dbt manifest.json
    manifestPath := req.Config["manifest_path"]
    if manifestPath == "" {
        manifestPath = c.manifestPath
    }
    
    manifestData, err := os.ReadFile(manifestPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read manifest: %w", err)
    }
    
    var manifest DbtManifest
    if err := json.Unmarshal(manifestData, &manifest); err != nil {
        return nil, fmt.Errorf("failed to parse manifest: %w", err)
    }
    
    // 2. Extract datasets (models, sources, seeds)
    var datasets []*pb.Dataset
    
    // Extract models
    for _, model := range manifest.Nodes {
        if model.ResourceType == "model" {
            datasets = append(datasets, &pb.Dataset{
                Id:   model.UniqueID,
                Name: model.Name,
                Type: "table",
                Owner: model.Config.Meta["owner"],
                Schema: &pb.Schema{
                    Columns: extractColumns(model.Columns),
                },
                Metadata: map[string]string{
                    "database":     model.Database,
                    "schema":       model.Schema,
                    "materialized": model.Config.Materialized,
                    "description":  model.Description,
                },
            })
        }
    }
    
    // Extract sources
    for _, source := range manifest.Sources {
        datasets = append(datasets, &pb.Dataset{
            Id:   source.UniqueID,
            Name: source.Name,
            Type: "table",
            Owner: source.Meta["owner"],
            Schema: &pb.Schema{
                Columns: extractColumns(source.Columns),
            },
            Metadata: map[string]string{
                "database":    source.Database,
                "schema":      source.Schema,
                "source_name": source.SourceName,
                "description": source.Description,
            },
        })
    }
    
    // 3. Extract pipelines (models as transformations)
    var pipelines []*pb.Pipeline
    
    for _, model := range manifest.Nodes {
        if model.ResourceType == "model" {
            pipelines = append(pipelines, &pb.Pipeline{
                Id:   model.UniqueID,
                Name: model.Name,
                Type: "dbt_model",
                Owner: model.Config.Meta["owner"],
                InputDatasets:  extractDependencies(model.DependsOn),
                OutputDatasets: []string{model.UniqueID},
                Metadata: map[string]string{
                    "sql":         model.RawSQL,
                    "description": model.Description,
                },
            })
        }
    }
    
    return &pb.ExtractMetadataResponse{
        Datasets:  datasets,
        Pipelines: pipelines,
    }, nil
}

func (c *DbtConnector) ExtractLineage(ctx context.Context, req *pb.ExtractLineageRequest) (*pb.ExtractLineageResponse, error) {
    // Parse manifest
    manifestPath := req.Config["manifest_path"]
    if manifestPath == "" {
        manifestPath = c.manifestPath
    }
    
    manifestData, err := os.ReadFile(manifestPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read manifest: %w", err)
    }
    
    var manifest DbtManifest
    if err := json.Unmarshal(manifestData, &manifest); err != nil {
        return nil, fmt.Errorf("failed to parse manifest: %w", err)
    }
    
    // Extract lineage edges
    var edges []*pb.LineageEdge
    
    for _, model := range manifest.Nodes {
        if model.ResourceType == "model" {
            // Model depends on upstream datasets
            for _, dep := range model.DependsOn.Nodes {
                edges = append(edges, &pb.LineageEdge{
                    SourceId: dep,
                    TargetId: model.UniqueID,
                    EdgeType: "dataset_to_dataset",
                    Metadata: map[string]string{
                        "type": "dbt_dependency",
                    },
                })
            }
            
            // Model produces output dataset
            edges = append(edges, &pb.LineageEdge{
                SourceId: model.UniqueID,
                TargetId: fmt.Sprintf("%s.%s.%s", model.Database, model.Schema, model.Name),
                EdgeType: "pipeline_to_dataset",
                Metadata: map[string]string{
                    "type": "dbt_model_output",
                },
            })
        }
    }
    
    return &pb.ExtractLineageResponse{
        Edges: edges,
    }, nil
}

func (c *DbtConnector) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
    // Check if manifest file exists
    if _, err := os.Stat(c.manifestPath); os.IsNotExist(err) {
        return &pb.HealthCheckResponse{
            Status:  pb.HealthCheckResponse_UNHEALTHY,
            Message: fmt.Sprintf("manifest file not found: %s", c.manifestPath),
        }, nil
    }
    
    return &pb.HealthCheckResponse{
        Status:  pb.HealthCheckResponse_HEALTHY,
        Message: "dbt connector is healthy",
    }, nil
}

// Helper types
type DbtManifest struct {
    Nodes   map[string]DbtNode   `json:"nodes"`
    Sources map[string]DbtSource `json:"sources"`
}

type DbtNode struct {
    UniqueID     string            `json:"unique_id"`
    Name         string            `json:"name"`
    ResourceType string            `json:"resource_type"`
    Database     string            `json:"database"`
    Schema       string            `json:"schema"`
    Description  string            `json:"description"`
    RawSQL       string            `json:"raw_sql"`
    Columns      map[string]Column `json:"columns"`
    DependsOn    DependsOn         `json:"depends_on"`
    Config       Config            `json:"config"`
}

type DbtSource struct {
    UniqueID   string            `json:"unique_id"`
    Name       string            `json:"name"`
    SourceName string            `json:"source_name"`
    Database   string            `json:"database"`
    Schema     string            `json:"schema"`
    Description string           `json:"description"`
    Columns    map[string]Column `json:"columns"`
    Meta       map[string]string `json:"meta"`
}

type Column struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    Description string `json:"description"`
}

type DependsOn struct {
    Nodes []string `json:"nodes"`
}

type Config struct {
    Materialized string            `json:"materialized"`
    Meta         map[string]string `json:"meta"`
}

func extractColumns(columns map[string]Column) []*pb.Column {
    var result []*pb.Column
    for _, col := range columns {
        result = append(result, &pb.Column{
            Name:        col.Name,
            Type:        col.Type,
            Description: col.Description,
        })
    }
    return result
}

func extractDependencies(deps DependsOn) []string {
    return deps.Nodes
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "50061"
    }
    
    lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    grpcServer := grpc.NewServer()
    pb.RegisterConnectorServer(grpcServer, NewDbtConnector())
    
    log.Printf("dbt connector listening on port %s", port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
```

---

## Connector Dockerfile

```dockerfile
# Dockerfile.connector-dbt
FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY connectors/dbt ./connectors/dbt
COPY proto ./proto

RUN CGO_ENABLED=0 go build -o /app/connector-dbt ./connectors/dbt

FROM gcr.io/distroless/static-debian11:nonroot

COPY --from=builder /app/connector-dbt /app/connector-dbt

USER nonroot:nonroot

EXPOSE 50061

ENTRYPOINT ["/app/connector-dbt"]
```

---

# Connector Lifecycle Management

## Connector Registration

**How connectors are registered with Provider Service:**

### 1. **Static Registration (V0)**

**Configuration file:**
```yaml
# config/providers.yaml
providers:
  - name: dbt
    type: connector
    version: 1.0.0
    enabled: true
    endpoint: localhost:50061
    config:
      manifest_path: /etc/dbt/manifest.json
  
  - name: snowflake
    type: connector
    version: 1.0.0
    enabled: true
    endpoint: localhost:50062
    config:
      account: myaccount
      warehouse: COMPUTE_WH
      database: ANALYTICS
  
  - name: bigquery
    type: connector
    version: 1.0.0
    enabled: true
    endpoint: localhost:50063
    config:
      project_id: my-project
      dataset: analytics
```

**Provider Service loads configuration at startup:**
```go
// pkg/provider/registry.go
package provider

import (
    "context"
    "fmt"
    "os"
    
    pb "github.com/dcraft-fusion/proto/provider/v1"
    "google.golang.org/grpc"
    "gopkg.in/yaml.v3"
)

type Registry struct {
    connectors map[string]*ConnectorClient
}

type ConnectorClient struct {
    Name     string
    Version  string
    Endpoint string
    Config   map[string]string
    Client   pb.ConnectorClient
    Conn     *grpc.ClientConn
}

func NewRegistry(configPath string) (*Registry, error) {
    // Load configuration
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }
    
    var config struct {
        Providers []struct {
            Name     string            `yaml:"name"`
            Type     string            `yaml:"type"`
            Version  string            `yaml:"version"`
            Enabled  bool              `yaml:"enabled"`
            Endpoint string            `yaml:"endpoint"`
            Config   map[string]string `yaml:"config"`
        } `yaml:"providers"`
    }
    
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    // Initialize connectors
    registry := &Registry{
        connectors: make(map[string]*ConnectorClient),
    }
    
    for _, p := range config.Providers {
        if !p.Enabled || p.Type != "connector" {
            continue
        }
        
        // Connect to connector
        conn, err := grpc.Dial(p.Endpoint, grpc.WithInsecure())
        if err != nil {
            return nil, fmt.Errorf("failed to connect to %s: %w", p.Name, err)
        }
        
        client := pb.NewConnectorClient(conn)
        
        // Verify connector capabilities
        caps, err := client.GetCapabilities(context.Background(), &pb.GetCapabilitiesRequest{})
        if err != nil {
            conn.Close()
            return nil, fmt.Errorf("failed to get capabilities for %s: %w", p.Name, err)
        }
        
        registry.connectors[p.Name] = &ConnectorClient{
            Name:     p.Name,
            Version:  p.Version,
            Endpoint: p.Endpoint,
            Config:   p.Config,
            Client:   client,
            Conn:     conn,
        }
        
        log.Printf("Registered connector: %s v%s (capabilities: %v)", 
            caps.Name, caps.Version, caps.Operations)
    }
    
    return registry, nil
}

func (r *Registry) GetConnector(name string) (*ConnectorClient, error) {
    connector, ok := r.connectors[name]
    if !ok {
        return nil, fmt.Errorf("connector not found: %s", name)
    }
    return connector, nil
}

func (r *Registry) ListConnectors() []*ConnectorClient {
    var result []*ConnectorClient
    for _, c := range r.connectors {
        result = append(result, c)
    }
    return result
}

func (r *Registry) Close() {
    for _, c := range r.connectors {
        c.Conn.Close()
    }
}
```

---

### 2. **Dynamic Registration (V1 - Future)**

**API endpoint to register connectors at runtime:**

```go
// cmd/provider-service/main.go
func (s *ProviderService) RegisterConnector(ctx context.Context, req *pb.RegisterConnectorRequest) (*pb.RegisterConnectorResponse, error) {
    // 1. Validate connector image
    if req.Image == "" {
        return nil, fmt.Errorf("connector image is required")
    }
    
    // 2. Create connector Pod
    pod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name:      fmt.Sprintf("connector-%s", req.Name),
            Namespace: "dcraft-fusion",
            Labels: map[string]string{
                "app":       "connector",
                "connector": req.Name,
            },
        },
        Spec: corev1.PodSpec{
            Containers: []corev1.Container{
                {
                    Name:  req.Name,
                    Image: req.Image,
                    Ports: []corev1.ContainerPort{
                        {
                            ContainerPort: 50051,
                            Name:          "grpc",
                        },
                    },
                    Env: convertConfigToEnv(req.Config),
                },
            },
        },
    }
    
    // 3. Create Pod via Kubernetes API
    _, err := s.k8sClient.CoreV1().Pods("dcraft-fusion").Create(ctx, pod, metav1.CreateOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to create connector pod: %w", err)
    }
    
    // 4. Wait for Pod to be ready
    if err := s.waitForPodReady(ctx, pod.Name, 2*time.Minute); err != nil {
        return nil, fmt.Errorf("connector pod not ready: %w", err)
    }
    
    // 5. Connect to connector
    endpoint := fmt.Sprintf("%s.dcraft-fusion.svc.cluster.local:50051", pod.Name)
    conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
    if err != nil {
        return nil, fmt.Errorf("failed to connect to connector: %w", err)
    }
    
    client := pb.NewConnectorClient(conn)
    
    // 6. Verify capabilities
    caps, err := client.GetCapabilities(ctx, &pb.GetCapabilitiesRequest{})
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to get capabilities: %w", err)
    }
    
    // 7. Register in registry
    s.registry.connectors[req.Name] = &ConnectorClient{
        Name:     req.Name,
        Version:  caps.Version,
        Endpoint: endpoint,
        Config:   req.Config,
        Client:   client,
        Conn:     conn,
    }
    
    return &pb.RegisterConnectorResponse{
        ConnectorId: req.Name,
        Status:      "registered",
        Capabilities: caps,
    }, nil
}
```

---

## Connector Execution

**How Provider Service executes connectors:**

```go
// pkg/provider/executor.go
package provider

import (
    "context"
    "fmt"
    "time"
    
    pb "github.com/dcraft-fusion/proto/provider/v1"
)

type Executor struct {
    registry *Registry
}

func NewExecutor(registry *Registry) *Executor {
    return &Executor{registry: registry}
}

func (e *Executor) ExtractMetadata(ctx context.Context, providerName string) (*pb.ExtractMetadataResponse, error) {
    // 1. Get connector
    connector, err := e.registry.GetConnector(providerName)
    if err != nil {
        return nil, err
    }
    
    // 2. Call connector
    ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
    defer cancel()
    
    resp, err := connector.Client.ExtractMetadata(ctx, &pb.ExtractMetadataRequest{
        Config:        connector.Config,
        ResourceTypes: []string{"dataset", "pipeline"},
    })
    if err != nil {
        return nil, fmt.Errorf("connector %s failed: %w", providerName, err)
    }
    
    return resp, nil
}

func (e *Executor) ExtractLineage(ctx context.Context, providerName string) (*pb.ExtractLineageResponse, error) {
    // 1. Get connector
    connector, err := e.registry.GetConnector(providerName)
    if err != nil {
        return nil, err
    }
    
    // 2. Call connector
    ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
    defer cancel()
    
    resp, err := connector.Client.ExtractLineage(ctx, &pb.ExtractLineageRequest{
        Config: connector.Config,
    })
    if err != nil {
        return nil, fmt.Errorf("connector %s failed: %w", providerName, err)
    }
    
    return resp, nil
}

func (e *Executor) ExecuteAll(ctx context.Context) error {
    // Execute all connectors in parallel
    connectors := e.registry.ListConnectors()
    
    errChan := make(chan error, len(connectors))
    
    for _, connector := range connectors {
        go func(c *ConnectorClient) {
            // Extract metadata
            metadata, err := e.ExtractMetadata(ctx, c.Name)
            if err != nil {
                errChan <- fmt.Errorf("failed to extract metadata from %s: %w", c.Name, err)
                return
            }
            
            // Store metadata (send to Resource Service)
            if err := e.storeMetadata(ctx, metadata); err != nil {
                errChan <- fmt.Errorf("failed to store metadata from %s: %w", c.Name, err)
                return
            }
            
            // Extract lineage
            lineage, err := e.ExtractLineage(ctx, c.Name)
            if err != nil {
                errChan <- fmt.Errorf("failed to extract lineage from %s: %w", c.Name, err)
                return
            }
            
            // Store lineage (send to Lineage Service)
            if err := e.storeLineage(ctx, lineage); err != nil {
                errChan <- fmt.Errorf("failed to store lineage from %s: %w", c.Name, err)
                return
            }
            
            errChan <- nil
        }(connector)
    }
    
    // Wait for all connectors to finish
    var errors []error
    for i := 0; i < len(connectors); i++ {
        if err := <-errChan; err != nil {
            errors = append(errors, err)
        }
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("connector execution failed: %v", errors)
    }
    
    return nil
}
```

---

## Connector Scheduling

**How connectors are scheduled to run:**

### 1. **Cron-based Scheduling**

```yaml
# CronJob for provider execution
apiVersion: batch/v1
kind: CronJob
metadata:
  name: provider-sync
  namespace: dcraft-fusion
spec:
  schedule: "*/15 * * * *"  # Every 15 minutes
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: provider-sync
            image: registry.dcraft-fusion.io/provider-service:v0.1.0
            command: ["/app/provider-service", "sync"]
            env:
            - name: PROVIDERS_CONFIG
              value: /etc/config/providers.yaml
            volumeMounts:
            - name: config
              mountPath: /etc/config
          
          volumes:
          - name: config
            configMap:
              name: providers-config
          
          restartPolicy: OnFailure
```

---

### 2. **Event-based Scheduling**

```go
// pkg/provider/scheduler.go
package provider

import (
    "context"
    "time"
    
    "github.com/nats-io/nats.go"
)

type Scheduler struct {
    executor *Executor
    nats     *nats.Conn
}

func NewScheduler(executor *Executor, natsURL string) (*Scheduler, error) {
    nc, err := nats.Connect(natsURL)
    if err != nil {
        return nil, err
    }
    
    return &Scheduler{
        executor: executor,
        nats:     nc,
    }, nil
}

func (s *Scheduler) Start(ctx context.Context) error {
    // Subscribe to provider sync events
    _, err := s.nats.Subscribe("provider.sync.trigger", func(msg *nats.Msg) {
        log.Printf("Received sync trigger: %s", string(msg.Data))
        
        // Execute all connectors
        if err := s.executor.ExecuteAll(ctx); err != nil {
            log.Printf("Connector execution failed: %v", err)
            
            // Publish failure event
            s.nats.Publish("provider.sync.failed", []byte(err.Error()))
        } else {
            log.Printf("Connector execution succeeded")
            
            // Publish success event
            s.nats.Publish("provider.sync.succeeded", []byte("success"))
        }
    })
    
    if err != nil {
        return err
    }
    
    // Keep running
    <-ctx.Done()
    return nil
}
```

---

# Dynamic Connector Loading

## Hot Reload Architecture

**Goal:** Add/remove/update connectors without restarting Provider Service.

### Architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                    Provider Service                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  Connector Registry (In-Memory)                    │     │
│  │  - map[string]*ConnectorClient                     │     │
│  │  - Thread-safe (sync.RWMutex)                      │     │
│  └────────────────────────────────────────────────────┘     │
│                   ▲                                          │
│                   │                                          │
│  ┌────────────────┴───────────────────────────────────┐     │
│  │  Connector Watcher (File System Watcher)           │     │
│  │  - Watches /etc/connectors/*.yaml                  │     │
│  │  - Detects changes (create, update, delete)        │     │
│  │  - Triggers reload                                 │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
│  ┌────────────────────────────────────────────────────┐     │
│  │  Connector Loader                                  │     │
│  │  - Loads connector configuration                   │     │
│  │  - Creates connector Pod (Kubernetes API)          │     │
│  │  - Connects to connector (gRPC)                    │     │
│  │  - Registers in registry                           │     │
│  └────────────────────────────────────────────────────┘     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Implementation:

```go
// pkg/provider/watcher.go
package provider

import (
    "context"
    "log"
    "path/filepath"
    
    "github.com/fsnotify/fsnotify"
)

type Watcher struct {
    registry *Registry
    loader   *Loader
    watcher  *fsnotify.Watcher
    configDir string
}

func NewWatcher(registry *Registry, loader *Loader, configDir string) (*Watcher, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, err
    }
    
    return &Watcher{
        registry:  registry,
        loader:    loader,
        watcher:   watcher,
        configDir: configDir,
    }, nil
}

func (w *Watcher) Start(ctx context.Context) error {
    // Add config directory to watcher
    if err := w.watcher.Add(w.configDir); err != nil {
        return err
    }
    
    log.Printf("Watching connector configs in %s", w.configDir)
    
    for {
        select {
        case event := <-w.watcher.Events:
            if filepath.Ext(event.Name) != ".yaml" {
                continue
            }
            
            switch {
            case event.Op&fsnotify.Create == fsnotify.Create:
                log.Printf("New connector config detected: %s", event.Name)
                w.handleCreate(ctx, event.Name)
            
            case event.Op&fsnotify.Write == fsnotify.Write:
                log.Printf("Connector config updated: %s", event.Name)
                w.handleUpdate(ctx, event.Name)
            
            case event.Op&fsnotify.Remove == fsnotify.Remove:
                log.Printf("Connector config removed: %s", event.Name)
                w.handleRemove(ctx, event.Name)
            }
        
        case err := <-w.watcher.Errors:
            log.Printf("Watcher error: %v", err)
        
        case <-ctx.Done():
            return w.watcher.Close()
        }
    }
}

func (w *Watcher) handleCreate(ctx context.Context, configPath string) {
    // Load connector
    connector, err := w.loader.Load(ctx, configPath)
    if err != nil {
        log.Printf("Failed to load connector: %v", err)
        return
    }
    
    // Register connector
    w.registry.Register(connector)
    
    log.Printf("Connector registered: %s", connector.Name)
}

func (w *Watcher) handleUpdate(ctx context.Context, configPath string) {
    // Reload connector
    connector, err := w.loader.Load(ctx, configPath)
    if err != nil {
        log.Printf("Failed to reload connector: %v", err)
        return
    }
    
    // Unregister old connector
    w.registry.Unregister(connector.Name)
    
    // Register new connector
    w.registry.Register(connector)
    
    log.Printf("Connector reloaded: %s", connector.Name)
}

func (w *Watcher) handleRemove(ctx context.Context, configPath string) {
    // Extract connector name from file path
    name := filepath.Base(configPath)
    name = name[:len(name)-len(filepath.Ext(name))]
    
    // Unregister connector
    w.registry.Unregister(name)
    
    log.Printf("Connector unregistered: %s", name)
}
```

---

## Thread-Safe Registry:

```go
// pkg/provider/registry.go
package provider

import (
    "fmt"
    "sync"
)

type Registry struct {
    mu         sync.RWMutex
    connectors map[string]*ConnectorClient
}

func NewRegistry() *Registry {
    return &Registry{
        connectors: make(map[string]*ConnectorClient),
    }
}

func (r *Registry) Register(connector *ConnectorClient) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    r.connectors[connector.Name] = connector
}

func (r *Registry) Unregister(name string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    // Close connection
    if connector, ok := r.connectors[name]; ok {
        connector.Conn.Close()
    }
    
    delete(r.connectors, name)
}

func (r *Registry) GetConnector(name string) (*ConnectorClient, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    connector, ok := r.connectors[name]
    if !ok {
        return nil, fmt.Errorf("connector not found: %s", name)
    }
    return connector, nil
}

func (r *Registry) ListConnectors() []*ConnectorClient {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var result []*ConnectorClient
    for _, c := range r.connectors {
        result = append(result, c)
    }
    return result
}
```

---

# Kubernetes Operator Pattern

## Overview

**Kubernetes Operator** = Custom Resource Definition (CRD) + Custom Controller

**Use Case:** Manage provider connectors as Kubernetes resources.

---

## Custom Resource Definition (CRD)

```yaml
# manifests/crds/provider-crd.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: providers.dcraft.fusion
spec:
  group: dcraft.fusion
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              type:
                type: string
                enum: [dbt, snowflake, bigquery, postgres, airflow]
              version:
                type: string
              enabled:
                type: boolean
              config:
                type: object
                x-kubernetes-preserve-unknown-fields: true
          status:
            type: object
            properties:
              phase:
                type: string
                enum: [Pending, Running, Failed, Succeeded]
              lastSyncTime:
                type: string
                format: date-time
              message:
                type: string
  scope: Namespaced
  names:
    plural: providers
    singular: provider
    kind: Provider
    shortNames:
    - prov
```

---

## Custom Resource (CR)

```yaml
# manifests/providers/dbt-provider.yaml
apiVersion: dcraft.fusion/v1
kind: Provider
metadata:
  name: dbt-provider
  namespace: dcraft-fusion
spec:
  type: dbt
  version: "1.0.0"
  enabled: true
  config:
    manifestPath: /etc/dbt/manifest.json
    runResultsPath: /etc/dbt/run_results.json
```

---

## Operator Controller

```go
// cmd/provider-operator/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/watch"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/cache"
)

type ProviderController struct {
    k8sClient *kubernetes.Clientset
    informer  cache.SharedIndexInformer
}

func NewProviderController() (*ProviderController, error) {
    // Create Kubernetes client
    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, err
    }
    
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }
    
    // Create informer (watches Provider resources)
    informer := cache.NewSharedIndexInformer(
        &cache.ListWatch{
            ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
                return listProviders(clientset)
            },
            WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
                return watchProviders(clientset)
            },
        },
        &Provider{},
        time.Minute,
        cache.Indexers{},
    )
    
    return &ProviderController{
        k8sClient: clientset,
        informer:  informer,
    }, nil
}

func (c *ProviderController) Run(ctx context.Context) error {
    // Add event handlers
    c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            provider := obj.(*Provider)
            log.Printf("Provider added: %s", provider.Name)
            c.reconcile(ctx, provider)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            provider := newObj.(*Provider)
            log.Printf("Provider updated: %s", provider.Name)
            c.reconcile(ctx, provider)
        },
        DeleteFunc: func(obj interface{}) {
            provider := obj.(*Provider)
            log.Printf("Provider deleted: %s", provider.Name)
            c.cleanup(ctx, provider)
        },
    })
    
    // Start informer
    go c.informer.Run(ctx.Done())
    
    // Wait for cache sync
    if !cache.WaitForCacheSync(ctx.Done(), c.informer.HasSynced) {
        return fmt.Errorf("failed to sync cache")
    }
    
    log.Println("Provider controller started")
    
    <-ctx.Done()
    return nil
}

func (c *ProviderController) reconcile(ctx context.Context, provider *Provider) {
    // Reconciliation loop: ensure actual state matches desired state
    
    // 1. Check if connector Pod exists
    podName := fmt.Sprintf("connector-%s", provider.Name)
    pod, err := c.k8sClient.CoreV1().Pods(provider.Namespace).Get(ctx, podName, metav1.GetOptions{})
    
    if err != nil {
        // Pod doesn't exist, create it
        log.Printf("Creating connector Pod for %s", provider.Name)
        
        pod = &corev1.Pod{
            ObjectMeta: metav1.ObjectMeta{
                Name:      podName,
                Namespace: provider.Namespace,
                Labels: map[string]string{
                    "app":       "connector",
                    "connector": provider.Name,
                    "type":      provider.Spec.Type,
                },
                OwnerReferences: []metav1.OwnerReference{
                    {
                        APIVersion: "dcraft.fusion/v1",
                        Kind:       "Provider",
                        Name:       provider.Name,
                        UID:        provider.UID,
                    },
                },
            },
            Spec: corev1.PodSpec{
                Containers: []corev1.Container{
                    {
                        Name:  provider.Spec.Type,
                        Image: fmt.Sprintf("registry.dcraft-fusion.io/connector-%s:%s", provider.Spec.Type, provider.Spec.Version),
                        Ports: []corev1.ContainerPort{
                            {
                                ContainerPort: 50051,
                                Name:          "grpc",
                            },
                        },
                        Env: convertConfigToEnv(provider.Spec.Config),
                    },
                },
            },
        }
        
        _, err = c.k8sClient.CoreV1().Pods(provider.Namespace).Create(ctx, pod, metav1.CreateOptions{})
        if err != nil {
            log.Printf("Failed to create Pod: %v", err)
            c.updateStatus(ctx, provider, "Failed", err.Error())
            return
        }
        
        c.updateStatus(ctx, provider, "Pending", "Pod created")
        return
    }
    
    // 2. Check Pod status
    if pod.Status.Phase == corev1.PodRunning {
        c.updateStatus(ctx, provider, "Running", "Connector is running")
    } else if pod.Status.Phase == corev1.PodFailed {
        c.updateStatus(ctx, provider, "Failed", "Pod failed")
    }
}

func (c *ProviderController) cleanup(ctx context.Context, provider *Provider) {
    // Delete connector Pod
    podName := fmt.Sprintf("connector-%s", provider.Name)
    err := c.k8sClient.CoreV1().Pods(provider.Namespace).Delete(ctx, podName, metav1.DeleteOptions{})
    if err != nil {
        log.Printf("Failed to delete Pod: %v", err)
    }
}

func (c *ProviderController) updateStatus(ctx context.Context, provider *Provider, phase, message string) {
    provider.Status.Phase = phase
    provider.Status.Message = message
    provider.Status.LastSyncTime = metav1.Now()
    
    // Update Provider status (via Kubernetes API)
    // ... (implementation depends on custom client)
}

type Provider struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    
    Spec   ProviderSpec   `json:"spec"`
    Status ProviderStatus `json:"status,omitempty"`
}

type ProviderSpec struct {
    Type    string            `json:"type"`
    Version string            `json:"version"`
    Enabled bool              `json:"enabled"`
    Config  map[string]string `json:"config"`
}

type ProviderStatus struct {
    Phase        string      `json:"phase"`
    LastSyncTime metav1.Time `json:"lastSyncTime"`
    Message      string      `json:"message"`
}

func main() {
    controller, err := NewProviderController()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    if err := controller.Run(ctx); err != nil {
        log.Fatal(err)
    }
}
```

---

## Benefits of Operator Pattern:

1. **Declarative:** Define desired state in YAML, operator ensures it
2. **Automated:** Operator handles creation, updates, deletion
3. **Self-healing:** Operator recreates failed Pods
4. **Kubernetes-native:** Integrates seamlessly with kubectl, GitOps, etc.

**Usage:**
```bash
# Create provider
kubectl apply -f manifests/providers/dbt-provider.yaml

# List providers
kubectl get providers -n dcraft-fusion

# Describe provider
kubectl describe provider dbt-provider -n dcraft-fusion

# Delete provider
kubectl delete provider dbt-provider -n dcraft-fusion
```

---

# Advanced Kubernetes Patterns

## 1. StatefulSets (for Databases)

**Use Case:** PostgreSQL, Neo4j, Redis (stateful workloads)

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgresql
  namespace: dcraft-fusion
spec:
  serviceName: postgresql
  replicas: 3
  selector:
    matchLabels:
      app: postgresql
  template:
    metadata:
      labels:
        app: postgresql
    spec:
      containers:
      - name: postgresql
        image: postgres:15
        ports:
        - containerPort: 5432
          name: postgres
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
  
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      storageClassName: fast-ssd
      resources:
        requests:
          storage: 100Gi
```

**Why StatefulSet?**
- Stable network identity (postgresql-0, postgresql-1, postgresql-2)
- Stable persistent storage (each Pod gets its own PVC)
- Ordered deployment and scaling
- Ordered rolling updates

---

## 2. Jobs (for One-Time Tasks)

**Use Case:** Database migrations, data imports

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: db-migrate
  namespace: dcraft-fusion
spec:
  template:
    spec:
      containers:
      - name: migrate
        image: registry.dcraft-fusion.io/resource-service:v0.1.0
        command: ["/app/migrate", "up"]
        env:
        - name: POSTGRES_HOST
          value: postgresql
      restartPolicy: OnFailure
  backoffLimit: 3
```

---

## 3. CronJobs (for Scheduled Tasks)

**Use Case:** Provider sync, backup, cleanup

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: provider-sync
  namespace: dcraft-fusion
spec:
  schedule: "*/15 * * * *"  # Every 15 minutes
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: sync
            image: registry.dcraft-fusion.io/provider-service:v0.1.0
            command: ["/app/provider-service", "sync"]
          restartPolicy: OnFailure
```

---

## 4. ConfigMaps & Secrets

**ConfigMaps** for non-sensitive configuration:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: providers-config
  namespace: dcraft-fusion
data:
  providers.yaml: |
    providers:
      - name: dbt
        type: connector
        enabled: true
```

**Secrets** for sensitive data:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: postgresql-credentials
  namespace: dcraft-fusion
type: Opaque
stringData:
  username: dcraft
  password: supersecret
```

---

## 5. Service Accounts & RBAC

**Service Account** for Provider Operator:
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: provider-operator
  namespace: dcraft-fusion
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: provider-operator
  namespace: dcraft-fusion
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch", "create", "update", "delete"]
- apiGroups: ["dcraft.fusion"]
  resources: ["providers"]
  verbs: ["get", "list", "watch", "update"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: provider-operator
  namespace: dcraft-fusion
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: provider-operator
subjects:
- kind: ServiceAccount
  name: provider-operator
  namespace: dcraft-fusion
```

---

# Complete Implementation Examples

## Example 1: Adding a New Connector (Airflow)

### Step 1: Implement Connector

```go
// connectors/airflow/main.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    
    pb "github.com/dcraft-fusion/proto/provider/v1"
)

type AirflowConnector struct {
    pb.UnimplementedConnectorServer
    baseURL string
    apiKey  string
}

func (c *AirflowConnector) ExtractMetadata(ctx context.Context, req *pb.ExtractMetadataRequest) (*pb.ExtractMetadataResponse, error) {
    // 1. Call Airflow API to get DAGs
    resp, err := http.Get(fmt.Sprintf("%s/api/v1/dags", c.baseURL))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var dags struct {
        DAGs []struct {
            DAGID       string `json:"dag_id"`
            Description string `json:"description"`
            Owners      []string `json:"owners"`
            Tags        []struct {
                Name string `json:"name"`
            } `json:"tags"`
        } `json:"dags"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&dags); err != nil {
        return nil, err
    }
    
    // 2. Convert DAGs to pipelines
    var pipelines []*pb.Pipeline
    for _, dag := range dags.DAGs {
        pipelines = append(pipelines, &pb.Pipeline{
            Id:    dag.DAGID,
            Name:  dag.DAGID,
            Type:  "airflow_dag",
            Owner: dag.Owners[0],
            Metadata: map[string]string{
                "description": dag.Description,
            },
        })
    }
    
    return &pb.ExtractMetadataResponse{
        Pipelines: pipelines,
    }, nil
}

// ... (implement other methods)
```

### Step 2: Build Docker Image

```dockerfile
# Dockerfile.connector-airflow
FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY connectors/airflow ./connectors/airflow
COPY proto ./proto
RUN CGO_ENABLED=0 go build -o /app/connector-airflow ./connectors/airflow

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=builder /app/connector-airflow /app/connector-airflow
USER nonroot:nonroot
EXPOSE 50064
ENTRYPOINT ["/app/connector-airflow"]
```

```bash
docker build -t registry.dcraft-fusion.io/connector-airflow:v0.1.0 -f Dockerfile.connector-airflow .
docker push registry.dcraft-fusion.io/connector-airflow:v0.1.0
```

### Step 3: Create Provider CR

```yaml
# manifests/providers/airflow-provider.yaml
apiVersion: dcraft.fusion/v1
kind: Provider
metadata:
  name: airflow-provider
  namespace: dcraft-fusion
spec:
  type: airflow
  version: "1.0.0"
  enabled: true
  config:
    baseURL: "https://airflow.example.com"
    apiKey: "secret"
```

### Step 4: Apply

```bash
kubectl apply -f manifests/providers/airflow-provider.yaml
```

**Operator automatically:**
1. Creates connector Pod
2. Waits for Pod to be ready
3. Registers connector in Provider Service
4. Starts extracting metadata

---

## Example 2: Hot Reload Connector Configuration

### Step 1: Update Configuration

```bash
# Edit connector config
kubectl edit provider dbt-provider -n dcraft-fusion
```

Change:
```yaml
spec:
  config:
    manifestPath: /etc/dbt/manifest.json
    # Add new config
    runResultsPath: /etc/dbt/run_results.json
```

### Step 2: Operator Detects Change

**Operator automatically:**
1. Detects Provider resource update
2. Recreates connector Pod with new config
3. Waits for new Pod to be ready
4. Re-registers connector in Provider Service

**No manual intervention required!**

---

# Summary

## Sidecar Strategy

| Service | Sidecars | Purpose |
|---------|----------|---------|
| All Services | Linkerd Proxy | mTLS, retries, circuit breakers, observability |
| Provider Service | 3 Connector Sidecars | Isolate connector logic, enable hot-swapping |
| AI Operator Service | Redis Cache Sidecar | Ultra-low latency cache access |

**Total containers: 24** (10 main + 10 Linkerd + 3 connectors + 1 Redis)

---

## Connector Architecture

1. **Interface:** gRPC (GetCapabilities, ExtractMetadata, ExtractLineage, HealthCheck)
2. **Deployment:** Sidecar containers (V0) or standalone Pods (V1)
3. **Registration:** Static (config file) or Dynamic (Kubernetes API)
4. **Execution:** Scheduled (CronJob) or Event-driven (NATS)
5. **Hot Reload:** File system watcher + thread-safe registry

---

## Advanced Kubernetes Patterns

1. **Init Containers:** Database migrations, dependency waiting
2. **DaemonSets:** Log collection, node monitoring
3. **StatefulSets:** Databases (PostgreSQL, Neo4j, Redis)
4. **Jobs:** One-time tasks (migrations, imports)
5. **CronJobs:** Scheduled tasks (sync, backup)
6. **Operator Pattern:** Custom resources + custom controller

---

## Key Benefits

✅ **Isolation:** Connector failures don't crash main service
✅ **Hot-swapping:** Add/remove connectors without restarts
✅ **Scalability:** Scale connectors independently
✅ **Flexibility:** Connectors can be in any language
✅ **Kubernetes-native:** Integrates with kubectl, GitOps, etc.
✅ **Self-healing:** Operator recreates failed Pods
✅ **Declarative:** Define desired state, operator ensures it

---

**This architecture provides a production-ready, extensible, and maintainable connector system for DCraft Fusion V0.**