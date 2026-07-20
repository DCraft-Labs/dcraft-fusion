# DCraft Fusion — Complete Project Context for Codex

> **Purpose of this file**  
> This Markdown file is designed to be placed inside a repository, ideally as `.ai/PROJECT_CONTEXT.md`, `PROJECT_CONTEXT.md`, or `docs/PROJECT_CONTEXT.md`, so that Codex or any coding agent can understand the full product vision, architectural decisions, implementation direction, constraints, and working assumptions of the project.
>
> This is not a marketing document. It is a detailed engineering/product context document meant to guide code generation, repository structure, architecture decisions, implementation planning, refactoring, and review.

---

## 0. How Codex Should Use This Document

Codex must treat this file as the main product and architecture memory for the repository.

When making code changes, Codex should:

1. Preserve the core product identity and architectural principles.
2. Avoid implementing features that violate the non-data-owning control plane model.
3. Keep execution engines pluggable and loosely coupled.
4. Prefer modular services and adapters instead of hard-wiring every capability into one monolith.
5. Use open-source components where possible, but integrate them carefully through adapters, APIs, metadata contracts, or sidecar services.
6. Keep the AI Operator optional, auditable, approval-driven, and disableable per tenant/project.
7. Build toward a real enterprise-grade data platform control plane, not just a UI wrapper.
8. Prioritise the Control Plane Kernel first before building every 16-layer feature.
9. Treat this as a startup product that must become buildable in phases, with an MVP that demonstrates the core control-plane value.

Codex should not assume that the current repository already has all modules. If missing, Codex should propose and create skeletons incrementally.

---

## 1. Product Identity

### 1.1 Company and Product Naming

The locked naming rules are:

- **Company:** DCraft Labs
- **Product:** DCraft Fusion
- **Short name after first mention:** Fusion

Use **DCraft Fusion** on first mention in documents, UI headers, README files, and marketing pages. After first mention, use **Fusion**.

Avoid inconsistent names such as:

- dcraft fusion
- DcraftFusion
- Data Fusion
- DataCortex, unless referring to earlier product evolution
- Fusion, unless referring to the CDC/reporting platform inspiration or module history

### 1.2 Tagline

Current working tagline:

> **One brain. Many muscles.**

Meaning:

- Fusion is the brain/control plane.
- Existing data tools are the muscles/execution engines.
- Fusion coordinates, governs, observes, and assists.
- Fusion does not replace every engine.
- Fusion does not force customers to move raw data into Fusion by default.

### 1.3 One-Line Product Definition

**DCraft Fusion is an AI-assisted, engine-agnostic control plane for modern data platforms that sits above existing data tools, metadata systems, orchestration engines, transformation engines, BI tools, and governance systems, giving teams one governed operating layer for the full data engineering lifecycle.**

### 1.4 Short Product Description

Fusion helps organisations design, operate, govern, observe, and improve modern data platforms without forcing them into one execution engine. It provides a unified control plane and UX across ingestion, orchestration, transformation, quality, lineage, governance, cost, BI/reporting, reverse ETL, notebooks, CI/CD, ownership, and AI-assisted operations.

The control plane stores metadata, intent, policies, configurations, workflow state, audit logs, approval history, and operational visibility. The actual raw data, table contents, stream payloads, transformations, and compute execution remain inside the customer's own engines and infrastructure unless explicitly configured otherwise.

### 1.5 Central Claim

Fusion sits on top of the customer's existing stack and coordinates it.

It should not be positioned as:

- A replacement for Spark.
- A replacement for dbt.
- A replacement for Airflow.
- A replacement for Kafka.
- A replacement for Superset or BI tools.
- A mandatory data lake or warehouse.
- A SaaS product that must copy all customer data.

It should be positioned as:

- A control plane.
- A metadata and intent layer.
- A governed operational brain.
- A unified UX over fragmented tools.
- A safe AI-assisted operator for data teams.
- A platform that can run in SaaS, customer VPC, or on-prem/air-gapped environments.

---

## 2. Product Evolution and Context

The project has evolved through several related ideas. Codex should understand the history but treat the latest definition as canonical.

### 2.1 Earlier Name: DataCortex

DataCortex was an earlier product idea focused on a self-aware, contract-based data infrastructure platform enhanced with LLMs for automation, semantic understanding, and guided improvements.

Key DataCortex ideas that should survive inside Fusion:

- Data contracts.
- Schema lineage and diff engine.
- Test engine.
- Multi-connector integration.
- Enforcement modes.
- Audit and observability.
- LLM-powered mapping and explanation.
- Contract-aware sync.
- Quarantine for bad records.
- RBAC, field masking, and audit.
- Streaming-first contracts.
- GitOps and CI/CD.
- Data versioning.
- Smart Data Sync Layer.

### 2.2 Fusion CDC Engine

Fusion CDC Engine was a concrete architecture around multi-tenant real-time CDC and transformation.

Important ideas to preserve:

- Control plane versus data plane separation.
- Multi-tenant source configuration.
- CDC workers for MySQL, PostgreSQL, and MongoDB.
- Redis Streams as a lightweight event transport option.
- Checkpointing using local SQLite plus metadata store.
- Destination support for PostgreSQL warehouses and Iceberg/object storage.
- JSON flattening options.
- Sensitive data masking.
- Real-time and scheduled modes.
- Prometheus/Grafana observability.
- Circuit breakers, fallback queues, and retry strategies.
- Kubernetes deployment with HPA.
- Tenant isolation and per-tenant failure containment.

Fusion CDC Engine can become one execution/data-plane module under Fusion, but Fusion itself is broader than CDC.

### 2.3 Superset / BI Platform Research

The Superset fork and BI research produced a strong product direction for the BI/reporting layer.

Important ideas to preserve:

- Keep Apache Superset as the BI shell, charting engine, SQL Lab, dashboard engine, and embedded analytics foundation.
- Do not overload Superset core with every new feature.
- Build surrounding services for catalog, semantic model, quality, AI assistance, guided onboarding, and governance.
- Users should not start from raw tables.
- Users should start from certified business concepts, governed metrics, descriptions, synonyms, lineage, and quality scores.
- BI self-service must include semantic modelling, data understanding, relationship graph, metric registry, guided Explore, chart recommendations, dashboard generator, certification workflow, and AI assistance.

### 2.4 Final Canonical Product: DCraft Fusion

The final frozen product idea is:

> **DCraft Fusion is an AI-assisted, engine-agnostic control plane for modern data platforms with a tightly coupled control plane/UX and loosely coupled execution engines, covering the full 16-layer data engineering lifecycle, with optional and auditable AI Operator support and no mandatory raw data movement.**

This is the definition Codex should follow.

---

## 3. Non-Negotiable Architecture Principle

### 3.1 Core Principle

> **Tightly coupled control plane + UX. Loosely coupled execution engines.**

This is the most important architectural decision in the project.

### 3.2 What It Means

The Fusion control plane and Fusion UX should feel like one product:

- Unified navigation.
- Unified project model.
- Unified tenant model.
- Unified workflow state.
- Unified audit trail.
- Unified access control.
- Unified operational dashboard.
- Unified approval system.
- Unified AI assistant/operator.

But execution engines should remain loosely coupled:

- Airflow can run orchestration.
- dbt can run transformations.
- Spark can run large-scale compute.
- Kafka/Redpanda can run streaming.
- Trino can run federated SQL.
- DuckDB can run local/vectorized transforms.
- Superset can render dashboards.
- OpenMetadata/DataHub can provide metadata and lineage.
- Great Expectations/Soda can provide quality checks.
- Kubernetes can run workloads.
- Customer-managed engines can remain outside Fusion.

Fusion should coordinate these engines through adapters, contracts, metadata events, and APIs.

### 3.3 What Fusion Owns

Fusion owns:

- Metadata.
- Intent.
- Configurations.
- Policies.
- Workflow definitions.
- Workflow versions.
- Environment promotion state.
- Run state.
- Audit logs.
- Approval history.
- Access control mappings.
- Connector definitions.
- Adapter registry.
- Data product definitions.
- Semantic definitions.
- Quality check definitions.
- AI recommendations.
- AI action proposals.
- Observability metadata.
- Cost metadata.
- Lineage metadata.

### 3.4 What Fusion Must Not Own by Default

Fusion must not own by default:

- Raw table data.
- Full stream payloads.
- Customer PII.
- Warehouse copies.
- Production data lake storage.
- Mandatory transformation execution.
- Customer secrets in plaintext.
- AI training datasets from customer data.

Raw data movement must be optional and explicitly configured.

### 3.5 Design Test for Every Feature

Before implementing any feature, ask:

1. Is this control-plane state or execution-plane work?
2. Can Fusion store only metadata/intent instead of raw data?
3. Can the actual work be delegated to an adapter or external engine?
4. Can the user understand and approve what happens?
5. Is the action audited?
6. Can this be disabled per tenant/project?
7. Can this run in SaaS, customer VPC, and on-prem modes?

If the answer breaks the principle, redesign the feature.

---

## 4. Deployment Models

Fusion must support multiple deployment models.

### 4.1 SaaS

In SaaS mode:

- Fusion control plane is hosted by DCraft Labs.
- Customer connects metadata and execution engines through secure connectors.
- Raw data should not be copied unless customer explicitly enables it.
- Secrets must be encrypted and preferably stored in a managed secret backend.
- The SaaS control plane should prioritise metadata sync, orchestration control, observability, and AI recommendations.

### 4.2 Customer VPC / BYOC

In Customer VPC mode:

- Fusion control plane or major components run inside the customer's cloud account.
- Data, secrets, network access, and execution remain under customer control.
- DCraft Labs may provide deployment automation, Helm charts, Terraform, and support.
- This mode is important for regulated customers such as banks, fintechs, healthcare, and enterprises with strict data residency controls.

### 4.3 On-Prem / Air-Gapped

In on-prem or air-gapped mode:

- Fusion must run without mandatory outbound internet access.
- AI Operator must support BYO LLM or be disabled.
- Model calls must be configurable.
- All telemetry must be optional.
- Container images, Helm charts, and offline installation packages may be required.
- License activation should not break if internet is unavailable, unless offline license files are supported.

---

## 5. Target Personas

### 5.1 Data Platform Lead

Primary buyer/user.

Needs:

- One place to see platform health.
- Tenant/project isolation.
- Governance and audit.
- Control over tools and standards.
- Reduced operational sprawl.
- Cost visibility.
- Standardised onboarding.
- Approval gates.
- AI assistance that does not act dangerously.

### 5.2 Data Engineer

Primary builder.

Needs:

- Easy connector setup.
- Pipeline creation.
- dbt/Spark/Airflow integration.
- Schema evolution handling.
- Data quality checks.
- Run debugging.
- Logs and lineage.
- Safe transformation authoring.
- Git integration.

### 5.3 Analytics Engineer

Needs:

- Semantic models.
- Metrics-as-code.
- dbt integration.
- Data contracts.
- BI-ready datasets.
- Quality and freshness checks.
- Certified metrics.
- Impact analysis before changes.

### 5.4 BI / Reporting User

Needs:

- Guided self-service.
- Business-readable datasets.
- No raw table confusion.
- Certified dashboards.
- Prompt-assisted chart/dashboard generation.
- Export and scheduling.
- Permissions and row-level security.

### 5.5 Governance / Compliance User

Needs:

- Audit trail.
- Approval workflows.
- Data classification.
- PII masking policies.
- Lineage.
- Access control reports.
- Policy enforcement.
- Evidence for audits.

### 5.6 SRE / Platform Engineer

Needs:

- Health checks.
- Metrics.
- Logs.
- Runbooks.
- Kubernetes visibility.
- Queue lag.
- Retry state.
- Error grouping.
- Resource usage.
- Cost signals.

---

## 6. 16-Layer Data Engineering Lifecycle

Fusion should eventually cover a 16-layer lifecycle. Codex should not try to build all layers at once, but should understand the target system.

### Layer 1 — Control Plane Kernel

The central platform spine.

Responsibilities:

- Tenant model.
- Organisation model.
- Project/workspace model.
- User and role model.
- Resource registry.
- Workflow registry.
- Connection registry.
- Policy registry.
- Run registry.
- Event gateway.
- Audit log.
- Approval state.
- Environment model.
- Versioning model.

This is the first thing to build.

### Layer 2 — Identity, Access, and Tenant Isolation

Responsibilities:

- Authentication.
- SSO integration.
- Keycloak/OIDC support.
- RBAC.
- ABAC later if needed.
- Tenant-aware authorization.
- Project/workspace-level access.
- Service accounts.
- API tokens.
- Audit of access changes.

### Layer 3 — Source and Destination Connectivity

Responsibilities:

- Source registry.
- Destination registry.
- Credential handling.
- Connection testing.
- Schema discovery.
- Table discovery.
- Column discovery.
- Connector adapter model.
- Per-tenant connections.
- Support for MySQL, PostgreSQL, MongoDB, S3, Iceberg, Trino, warehouses, and APIs.

### Layer 4 — Ingestion and CDC

Responsibilities:

- Batch ingestion.
- Incremental ingestion.
- CDC ingestion.
- Event normalization.
- Schema evolution detection.
- Checkpointing.
- Pause/resume.
- Tenant-level isolation.
- Source-level throttling.
- Destination writes.

This layer can reuse ideas from Fusion CDC Engine.

### Layer 5 — Orchestration

Responsibilities:

- Workflow definitions.
- DAG registration.
- Airflow adapter.
- Temporal/Dagster/Argo adapter options later.
- Trigger runs.
- Get run status.
- Retry workflows.
- Schedule workflows.
- Environment promotion.
- Dependency mapping.

### Layer 6 — Transformation

Responsibilities:

- dbt adapter.
- Spark adapter.
- DuckDB transformation worker.
- SQL transformation definitions.
- Python UDF policies.
- Low-code transformations.
- JSON flattening.
- Masking.
- Type casting.
- Expression validation.
- Transformation lineage.

### Layer 7 — Storage and Lakehouse Control

Responsibilities:

- Iceberg table registry.
- Object storage registry.
- Hot/cold architecture metadata.
- Trino catalog metadata.
- Table lifecycle policies.
- Partition strategy metadata.
- Retention policy metadata.
- Compaction policy triggers.
- Storage cost visibility.

Fusion should not necessarily store data itself. It should manage metadata and coordinate engines.

### Layer 8 — Data Quality

Responsibilities:

- Quality rule definitions.
- Great Expectations integration.
- Soda Core integration.
- Custom SQL checks.
- Freshness checks.
- Null checks.
- Uniqueness checks.
- Referential checks.
- Schema drift checks.
- Quality score.
- Quality history.
- Failure routing.

### Layer 9 — Lineage and Metadata

Responsibilities:

- OpenMetadata/DataHub integration.
- Column-level lineage where possible.
- Table-level lineage.
- Job-level lineage.
- Dashboard lineage.
- Metric lineage.
- Data product lineage.
- Impact analysis.
- Metadata sync jobs.
- Metadata graph for AI context.

### Layer 10 — Catalog, Glossary, and Data Products

Responsibilities:

- Business glossary.
- Dataset profiles.
- Ownership metadata.
- Data product registry.
- Certification workflow.
- Tags and classifications.
- Data contracts.
- Documentation.
- Search.
- Discovery.
- Trust signals.

### Layer 11 — Semantic and Metrics Layer

Responsibilities:

- Metric registry.
- Dimensions.
- Measures.
- Entities.
- Relationships.
- Grain.
- Cardinality.
- Join rules.
- Synonyms.
- Business descriptions.
- Cube/MetricFlow/dbt semantic integration.
- AI-safe query planning.

This layer is critical for BI and AI.

### Layer 12 — BI, Reporting, and Analytics

Responsibilities:

- Superset fork integration.
- SQL Lab.
- Explore.
- Dashboards.
- Embedded analytics.
- Report generation.
- Scheduled exports.
- SFTP/S3/email delivery.
- Guided chart creation.
- Dashboard generator.
- Semantic dataset onboarding.

### Layer 13 — Observability and Operations

Responsibilities:

- Run status.
- Logs.
- Metrics.
- Alerts.
- Queue lag.
- Worker health.
- Error grouping.
- Retry state.
- SLA monitoring.
- Incident timeline.
- Runbooks.
- Prometheus/Grafana integration.

### Layer 14 — Governance, Security, and Compliance

Responsibilities:

- Access policies.
- PII classification.
- Masking policies.
- Approval workflows.
- Audit logs.
- Policy enforcement.
- Encryption metadata.
- Secret management integration.
- Compliance evidence.
- Tenant isolation.
- PCI DSS/SOX/UK GDPR-aligned controls.

### Layer 15 — Cost, Performance, and FinOps

Responsibilities:

- Job cost metadata.
- Warehouse/query cost metadata.
- Storage cost metadata.
- Idle resource detection.
- Expensive query detection.
- Recommendations.
- Budget alerts.
- Cost by tenant/project/job.
- Performance trend analysis.

### Layer 16 — AI Operator and Automation

Responsibilities:

- AI assistant.
- AI recommendations.
- AI-generated pipeline drafts.
- AI-generated quality rules.
- AI-generated dashboard specs.
- AI incident summaries.
- AI runbook suggestions.
- Human approval before execution.
- Full audit logs.
- BYO LLM support.
- Disableable per tenant/project.

---

## 7. AI Operator Rules

The AI Operator is powerful but must be safe.

### 7.1 Mandatory Constraints

The AI Operator must be:

- Optional.
- Disableable per tenant.
- Disableable per project.
- Disableable per action category.
- Human-approval based for changes.
- Audited.
- Explainable.
- BYO LLM compatible.
- Not trained on customer data by default.
- Not allowed to execute production changes without approval.

### 7.2 AI Can Suggest

AI may suggest:

- Pipeline design.
- Connector settings.
- dbt model drafts.
- SQL transformations.
- Data quality rules.
- Schema drift explanations.
- Incident summaries.
- Root-cause hypotheses.
- Cost optimisation suggestions.
- Dashboard layouts.
- Semantic model descriptions.
- Metric definitions.
- Documentation drafts.

### 7.3 AI Must Not Directly Do Without Approval

AI must not directly:

- Drop tables.
- Modify production pipelines.
- Change access policies.
- Expose secrets.
- Export sensitive data.
- Train on customer data.
- Disable quality checks.
- Promote workflows to production.
- Change tenant isolation rules.
- Execute destructive scripts.

### 7.4 AI Action Lifecycle

Recommended lifecycle:

1. User asks for help or system detects issue.
2. AI reads permitted metadata/context.
3. AI proposes action.
4. Fusion validates action against policy.
5. AI explains expected impact.
6. User reviews diff.
7. User approves or rejects.
8. Fusion executes through adapter.
9. Result is audited.
10. AI summarises outcome.

---

## 8. Control Plane Kernel — First Build Priority

The immediate recommended build is the Control Plane Kernel Skeleton.

This should come before trying to build every module.

### 8.1 Why Kernel First

Without the kernel, all modules become disconnected screens.

The kernel gives:

- Common tenant model.
- Common project model.
- Common registry.
- Common run model.
- Common event model.
- Common policy model.
- Common audit system.
- Common adapter framework.

### 8.2 Core Services to Start With

Recommended initial services/modules:

1. `identity-service`
2. `tenant-service`
3. `project-service`
4. `connection-registry-service`
5. `workflow-registry-service`
6. `run-service`
7. `policy-service`
8. `audit-service`
9. `event-gateway-service`
10. `adapter-registry-service`
11. `projection-service`
12. `ai-assist-service` later, not first

Depending on repo complexity, these can initially be modular packages inside a monorepo instead of separate deployable microservices.

### 8.3 Kernel Domain Objects

Minimum domain objects:

```text
Organization
Tenant
Workspace / Project
Environment
User
Role
Permission
ServiceAccount
Connection
ConnectorDefinition
AdapterDefinition
Workflow
WorkflowVersion
WorkflowRun
TaskRun
Policy
ApprovalRequest
AuditEvent
SystemEvent
DataProduct
Dataset
Metric
QualityRule
LineageNode
LineageEdge
AIRecommendation
```

### 8.4 Event Bus

Use event-driven design early.

Possible technologies:

- Kafka.
- Redpanda.
- Postgres outbox initially.
- Redis Streams for lightweight local/dev or specific CDC workloads.

Recommended pattern:

- Use Postgres transactional outbox for reliable event emission from services.
- Publish to Kafka/Redpanda when available.
- Keep event envelope consistent.

### 8.5 Standard Event Envelope

```json
{
  "event_id": "uuid",
  "event_type": "workflow.run.started",
  "event_version": "1.0",
  "occurred_at": "2026-05-18T00:00:00Z",
  "organization_id": "org_123",
  "tenant_id": "tenant_123",
  "project_id": "project_123",
  "environment": "dev|stage|prod",
  "actor": {
    "type": "user|service|system|ai",
    "id": "user_123"
  },
  "correlation_id": "uuid",
  "causation_id": "uuid",
  "payload": {}
}
```

### 8.6 Workflow Versioning and Promotion

Fusion should treat workflows like governed assets.

States:

```text
DRAFT -> REVIEW_REQUESTED -> APPROVED -> DEPLOYED_TO_DEV -> DEPLOYED_TO_STAGE -> DEPLOYED_TO_PROD -> DEPRECATED -> ARCHIVED
```

Promotion rules:

- Dev can be edited freely by authorised users.
- Stage requires validation.
- Production requires approval.
- AI-generated changes require review.
- All promotions create audit events.

---

## 9. Recommended Repository Structure

A practical monorepo structure:

```text
dcraft-fusion/
  README.md
  PROJECT_CONTEXT.md
  .ai/
    PROJECT_CONTEXT.md
    CODING_RULES.md
    ARCHITECTURE_DECISIONS.md
    AGENT_TASKS.md
    PROMPTS.md
  apps/
    web/
      # React/Next.js frontend
    api-gateway/
      # optional API gateway/BFF
  services/
    identity-service/
    tenant-service/
    project-service/
    connection-registry-service/
    workflow-registry-service/
    run-service/
    policy-service/
    audit-service/
    event-gateway-service/
    adapter-registry-service/
    ai-assist-service/
    metadata-sync-service/
    reporting-service/
    cdc-control-service/
  engines/
    cdc-workers/
    transform-worker/
    duckdb-worker/
    spark-adapter/
    dbt-adapter/
    airflow-adapter/
    superset-adapter/
    openmetadata-adapter/
  packages/
    domain/
    authz/
    event-envelope/
    connector-sdk/
    adapter-sdk/
    policy-sdk/
    ui-kit/
    telemetry/
  infra/
    docker/
    kubernetes/
    helm/
    terraform/
    local-dev/
  docs/
    architecture/
    adr/
    api/
    product/
    runbooks/
  tests/
    integration/
    e2e/
    contract/
```

### 9.1 Suggested Tech Stack

Backend:

- Go 1.26.x for control-plane services, API gateway, reconciliation, provider/adapter registry, audit, workflow/run services, and CLI.
- Python/FastAPI only where it already fits execution workers, AI/metadata jobs, connectors, or CDC/transform workers.
- Rust only for narrowly scoped performance-critical or security-sensitive components after profiling or clear technical need.
- gRPC for internal high-performance service communication where needed.
- REST/GraphQL for frontend-facing APIs.

Frontend:

- React + TypeScript.
- Tailwind/shadcn style system if building from scratch.
- Superset fork for BI module where required.

Data/control storage:

- PostgreSQL for control-plane metadata.
- Redis for cache, queues, lightweight streams, and WebSocket state where useful.
- Kafka/Redpanda for durable event bus.
- Object storage for artifacts/log exports, not raw data by default.

Execution integrations:

- Airflow.
- dbt.
- Spark.
- DuckDB.
- Trino.
- Kafka/Redpanda.
- Superset.
- OpenMetadata/DataHub.
- Great Expectations/Soda.

---

## 10. Fusion CDC Engine Context

The Fusion CDC Engine architecture is an important sub-project/inspiration.

### 10.1 What It Does

It is a multi-tenant CDC and transformation platform that:

- Connects to operational databases.
- Reads changes continuously using binlog/WAL/change streams.
- Streams events through Redis Streams or similar transport.
- Applies transformations and data quality checks.
- Writes cleaned data into destinations such as PostgreSQL warehouses or Iceberg tables on object storage.
- Supports one parent organisation such as a bank with multiple child tenants/customers.

### 10.2 Control Plane vs Data Plane

Control Plane:

- APIs.
- Metadata.
- Scheduling.
- Orchestration.
- Monitoring.
- Authentication.
- Configuration.

Data Plane:

- CDC workers.
- Log readers.
- Event enrichment.
- Transformation workers.
- Spark/DuckDB processors.
- Destination writers.

### 10.3 Major Components

```text
frontend/           React + TypeScript dashboard
control-plane/      APIs and orchestration
cdc-workers/        Python CDC readers
spark-consumer/     PySpark transformation engine, if used
transform-worker/   transformation processing
helm/               Helm charts
kubernetes/         Kubernetes manifests
monitoring/         Grafana + Prometheus
migrations/         DB schema migrations
```

### 10.4 CDC Worker Flow

```text
1. Worker connects to source database host.
2. Worker reads MySQL binlog, PostgreSQL WAL, or MongoDB change stream.
3. Worker normalizes event into standard event structure.
4. Worker adds metadata such as tenant ID, source ID, schema, table, event timestamp.
5. Worker publishes event to stream.
6. Worker maintains checkpoints.
7. Worker handles retries and failure states.
```

### 10.5 Event Routing Pattern

Example stream key pattern:

```text
cdc:<bank>:<tenant>:<source>:<schema>:<table>
```

In Fusion, this should become a configurable naming strategy, not a hardcoded-only pattern.

### 10.6 Checkpointing

Fusion used ideas around:

- SQLite local checkpoint.
- PostgreSQL metadata checkpoint.
- Per-source/per-tenant tracking.
- Resume from last known offset.

Fusion should support checkpoint metadata in the control plane, while allowing individual connectors to manage their native offsets.

### 10.7 Enterprise Features

Required features:

- Schema evolution.
- JSON flattening.
- Sensitive data masking.
- Real-time mode.
- Scheduled mode.
- Prometheus metrics.
- Fallback queues.
- Circuit breakers.
- Retry policies.
- Tenant isolation.
- Error grouping.

### 10.8 Challenges

Known risks:

- Operational complexity.
- Redis bottleneck at high scale if used as central transport.
- Spark operational overhead.
- Tenant balancing across workers.
- Per-tenant pause/resume complexity when source-level logs are shared.
- Backpressure handling.
- Checkpoint correctness.
- Duplicate handling.
- Delete tracking.

---

## 11. CDC Architecture Decision Direction

### 11.1 Previous Pain Point

There was a previous experience where Kafka + Debezium caused operational risk when data surges filled brokers across a 3-node replicated setup and downstream systems were impacted.

This does not mean Kafka/Debezium are bad. It means Fusion must provide:

- Backpressure.
- Quotas.
- Retention policies.
- Tenant isolation.
- Dead-letter queues.
- Circuit breakers.
- Replay controls.
- Per-tenant throttling.
- Monitoring before meltdown.

### 11.2 Redis Streams as an Option

Redis Streams may be useful for:

- Lightweight pipelines.
- Simpler deployments.
- Lower-latency fan-out.
- Internal CDC event transport for smaller installations.
- Local development.

But Redis should not be assumed to be the only durable enterprise backbone for all workloads.

### 11.3 Kafka / Redpanda as an Option

Kafka or Redpanda may be better for:

- High-volume durable streams.
- Replay.
- Consumer groups.
- Long retention.
- Enterprise streaming ecosystems.
- Stronger event bus semantics.

### 11.4 Fusion Design

Fusion should abstract stream transport behind an adapter:

```text
StreamTransportAdapter
  - publish(event)
  - subscribe(topic, consumerGroup)
  - commit(offset)
  - pause(topic/tenant)
  - resume(topic/tenant)
  - getLag()
  - getHealth()
```

Implementations:

- Redis Streams.
- Kafka.
- Redpanda.
- Cloud-native queues later.

---

## 12. Transformation Layer Details

### 12.1 Original 10 Transformation Types

Fusion CDC Engine discussed 10 transformation step types:

1. `cast`
2. `string_op`
3. `math_op`
4. `date_op`
5. `json_extract`
6. `json_flatten_inline`
7. `json_flatten_child`
8. `mask`
9. `expression`
10. `udf`

### 12.2 Transformation Execution Direction

Earlier architecture mentioned Spark, but later thinking moved toward:

- DuckDB + Celery + KEDA for many transformation jobs.
- Queue-based scaling.
- Vectorized single-pass transforms.
- Spark only when workload requires distributed processing.

Codex should not blindly implement Spark for everything.

Recommended rule:

- Use DuckDB for local/vectorized, medium-scale, file/table transforms.
- Use Spark for very large distributed transforms.
- Use dbt for warehouse/lakehouse SQL models.
- Use Python UDF only with strict sandboxing and policy.

### 12.3 Transformation Step Model

Example model:

```json
{
  "step_id": "step_001",
  "type": "json_flatten_inline",
  "input_columns": ["payload"],
  "output_columns": ["customer_name", "customer_age"],
  "config": {
    "path": "$.customer",
    "on_missing": "null",
    "on_error": "quarantine"
  },
  "policy": {
    "allowed_in_prod": true,
    "requires_approval": false
  }
}
```

### 12.4 JSON Handling

The platform should support:

- Treat JSON as string.
- Parse JSON.
- Flatten JSON inline into parent table.
- Flatten JSON into child table.
- Preserve original JSON.
- Extract selected paths.
- Infer schema.
- Handle schema drift.
- Quarantine malformed JSON.

### 12.5 Masking

Masking should support:

- Full mask.
- Partial mask.
- Hash.
- Tokenization placeholder.
- Regex-based masking.
- Policy-based masking.
- Role-aware visibility.

---

## 13. Superset Fork / BI Module

### 13.1 Product Direction

The Superset fork should become part of Fusion as the BI and reporting shell.

Do not try to rebuild all BI charts from scratch at first.

Use Superset for:

- SQL Lab.
- Explore.
- Charts.
- Dashboards.
- RBAC base.
- Visualization plugin ecosystem.
- Embedded dashboards.
- Dataset connections.

Build around Superset for:

- Semantic layer.
- Catalog.
- Glossary.
- Guided charting.
- AI chart generation.
- Certification workflow.
- Data quality display.
- Relationship graph.
- Metrics registry.
- Dashboard generator.
- Tenant-specific workflows.
- Scheduling and enterprise delivery.

### 13.2 Core Principle for BI

Users should never start from raw tables.

They should start from certified business concepts such as:

- Payments.
- Customers.
- Transactions.
- Settlement.
- Revenue.
- Churn.
- Active Customers.
- Failed Transactions.
- Successful Transaction Rate.
- Bank.
- Region.
- Date.
- Status.

### 13.3 Why Power BI/Tableau Are Different

The gap is not just charts.

Power BI and Tableau provide:

- Semantic modelling.
- Relationships.
- Measures.
- Business descriptions.
- Data prep.
- Certified sources.
- Governance.
- Collaboration.
- AI readiness.
- Lifecycle management.

Fusion's BI module must solve these gaps around Superset.

### 13.4 Open-Source Components to Reuse or Study

Potential components:

| Component | Role | How to Use |
|---|---|---|
| Apache Superset | BI shell | Keep as dashboard/chart foundation |
| Metabase | UX inspiration | Study visual question builder and simple onboarding |
| Lightdash | Semantic BI inspiration | Study dbt/YAML metrics workflow |
| Cube Core | Headless semantic layer | Strong candidate for runtime semantic API |
| MetricFlow/dbt Semantic Layer | Metrics-as-code | Good for dbt-centric customers |
| OpenMetadata | Catalog/governance/lineage | Strong integration candidate |
| DataHub | Metadata graph/enterprise lineage | Evaluate for larger metadata graph use cases |
| Great Expectations | Data quality | Integrate for expectations and validation docs |
| Soda Core | Data quality | Integrate for SQL-like DQ checks |
| Vega-Lite/ECharts | Chart specs | Use for AI-generated chart specs where needed |

### 13.5 BI Feature Specifications

Required modules:

#### Dataset Profile

- Owner.
- Description.
- Grain.
- Freshness.
- Quality score.
- Row count metadata.
- Column descriptions.
- Tags.
- PII flags.
- Sample disabled by default unless allowed.

#### Business Glossary

- Business terms.
- Synonyms.
- Definitions.
- Owners.
- Related datasets.
- Related metrics.
- Approval state.

#### Metric Registry

- Metric name.
- Business definition.
- Formula.
- Grain.
- Dimensions allowed.
- Filters.
- Owner.
- Certification status.
- Version.

#### Relationship Graph

- Logical relationships.
- Entity keys.
- Cardinality.
- Join rules.
- Fanout warnings.
- Many-to-many warnings.

#### Guided Explore

- Select business question.
- Select metric.
- Select dimension.
- Select time grain.
- Get chart recommendation.
- Generate Superset chart.

#### Natural Language Assistant

- Must use semantic metadata.
- Must not write raw SQL blindly.
- Must validate generated query/chart spec.
- Must show explanation.
- Must ask for approval where needed.

#### Certification Workflow

- Draft dataset/metric.
- Review.
- Approve.
- Certify.
- Deprecate.
- Archive.

#### Dashboard Generator

- Prompt to dashboard draft.
- Uses certified metrics only by default.
- Generates chart specs.
- Creates Superset dashboard through API.
- Requires review before publishing.

---

## 14. Reporting Portal Context

There is a related reporting platform vision from previous discussions.

### 14.1 Modules

The reporting portal has three modules:

1. **Yap Portal** for external users such as banks to download predefined reports with filters.
2. **Internal portal** for report configuration, scheduling, delivery, and notifications.
3. **EMR Reports** that run on headless Spark.

Fusion can absorb these concepts into its BI/reporting layer.

### 14.2 Data Sources

Allowed data sources:

- RDBMS tables.
- JSONs.
- Trino/Iceberg tables.

### 14.3 Output Formats

Required outputs:

- CSV.
- XLS.
- XLSX.
- TXT with custom separator.

### 14.4 Report Configuration

Users should configure:

- Report name.
- Description.
- Query/source.
- Parameters.
- Format.
- Separator.
- Schedule.
- Delivery target.
- Permissions.
- Tenant scope.

### 14.5 Delivery Targets

- Download.
- Email.
- SFTP.
- S3.
- Azure Blob.
- Webhook.
- Kafka notification.

### 14.6 Scale Targets

Previous target mentioned:

- Around 100 reports/sec generation throughput.
- Average report size around 100 MB.
- Strong retry semantics.
- Error-group-based reruns.
- Handling OOM killed jobs and transient failures.

Codex should treat this as a future high-scale target, not necessarily MVP.

### 14.7 Streaming Large Reports

Large report generation should use streaming, not loading entire result sets into memory.

Recommended patterns:

- Cursor-based DB reads.
- Backpressure-aware writer.
- Chunked object-store upload.
- Temporary artifact tracking.
- Resume/retry strategy.
- Memory-safe XLSX generation.

---

## 15. UAT / Real-World Reporting Features

Previous UAT scenarios for a reporting platform included:

### 15.1 Authentication

- SSO via Keycloak.
- Service account auth.
- gRPC interceptor adding headers.
- JWT claims such as `database_user_id`.

### 15.2 Tenancy Examples

Parent bank:

- Visa

Tenants:

- FBN
- Solidaire
- Access
- UBA
- SOFI

### 15.3 Report Scenarios

- Single-tenant report generation.
- Consolidated multi-tenant report generation.
- CSV generation.
- XLSX generation.
- Dataset filters by date range.
- Duplicate parameter handling in SQL.
- Custom CSV separators.
- Double quotes and newline handling.
- Large dataset streaming above 100k rows.
- Multi-tenant output above 150k rows.

### 15.4 Progress Updates

WebSocket progress examples:

```text
10% -> 25% -> 50% -> 75% -> 95% -> 100%
```

### 15.5 Cache

Example cache freshness type:

```text
DAILY_T1 with TTL = 6 hours / 21600 seconds
```

### 15.6 Circuit Breaker

Resilience4j states:

```text
CLOSED -> OPEN -> HALF_OPEN
```

### 15.7 Security

- Parameter binding safety.
- SQL injection guardrails.
- Tenant isolation.
- Per-tenant failure isolation.
- Encrypted fields with KMaaS.
- Credential caching with TTL.

---

## 16. Keycloak and SSO Context

Previous design involved two realms:

- `report-manager`
- `VisapayAgent`

Goal:

- User logs into VisapayAgent portal.
- User clicks a button.
- User lands in Report Manager/Fusion portal without re-login.

### 16.1 Preferred Direction

Use OIDC brokering rather than disabling authentication.

### 16.2 Identity Provider Setup

In `report-manager`, add Identity Provider:

- Type: OpenID Connect v1.0.
- Alias: `visapayagent`.
- Use discovered endpoints from partner realm.

### 16.3 Attribute Mapping Issue

Portal expects attributes such as:

- `bank_id`
- `bank_name`
- `sub_bank_ids`
- `sub_bank_names`

If partner realm does not provide them, use mapper strategy or internal mapping table.

Avoid hardcoding too much into frontend logic.

### 16.4 SSO Cookie Reality

For no-relogin SSO, the user needs an SSO session/cookie in the target realm or brokered flow.

If `prompt=none` is used without a valid session, it can fail with `login_required`.

A lightweight redirect endpoint can initiate the auth code flow properly.

### 16.5 Security Direction

Do not disable auth on the reporting/Fusion side.

Avoid insecure direct token passing unless properly verified.

---

## 17. Data Security and Compliance

Fusion is likely to target enterprise, banking, fintech, healthcare, and regulated data teams. Security must be first-class.

### 17.1 Required Controls

- Tenant isolation.
- RBAC.
- Service accounts.
- Audit logs.
- Secret encryption.
- Field masking.
- PII classification.
- Approval workflows.
- Environment-based permissions.
- Production change controls.
- Policy-as-code where possible.

### 17.2 Secrets

Secrets must not be stored in plaintext.

Integrate with:

- Kubernetes Secrets for local/dev only.
- Vault.
- AWS Secrets Manager.
- Azure Key Vault.
- GCP Secret Manager.
- Customer-managed secret stores.

### 17.3 Audit Event Examples

Audit everything meaningful:

- User login.
- Role change.
- Connection created.
- Connection tested.
- Credential changed.
- Workflow edited.
- Workflow approved.
- Workflow promoted.
- Workflow run started.
- Workflow run failed.
- AI recommendation generated.
- AI recommendation approved.
- Policy changed.
- Report exported.
- Dataset certified.

### 17.4 Data Movement Policy

Default:

- No raw data leaves customer engines.

Optional:

- Metadata sync.
- Profiling statistics.
- Sample data only if explicitly enabled.
- Raw data preview only with explicit permission and masking.

---

## 18. Open-Source Reuse Strategy

The project should use open-source wherever possible, but carefully.

### 18.1 Philosophy

Open-source projects exist so capabilities can be reused, integrated, extended, or studied. Fusion should not build everything from scratch if a mature open-source tool solves part of the problem.

But Fusion must avoid becoming a fragile bundle of unrelated tools.

The value is in:

- Unified UX.
- Unified metadata.
- Unified governance.
- Unified workflow state.
- Unified policy layer.
- Unified AI operator.
- Adapter-based integration.

### 18.2 Reuse Categories

1. **Direct integration** through APIs.
2. **Sidecar service** around an open-source tool.
3. **Adapter** that invokes an engine.
4. **Fork** only when necessary.
5. **UX inspiration** without copying code.
6. **Library reuse** where license allows.

### 18.3 License Watch-Out

Be careful with licenses, especially:

- AGPL.
- SSPL.
- Business Source License.
- Elastic License.
- Non-commercial licenses.

Codex should not copy code from projects without explicit license review.

### 18.4 Candidate Components

| Area | Candidate Tools |
|---|---|
| BI | Apache Superset, Metabase as UX reference, Lightdash as semantic reference |
| Semantic | Cube Core, MetricFlow, dbt Semantic Layer |
| Catalog/Metadata | OpenMetadata, DataHub |
| Quality | Great Expectations, Soda Core |
| Orchestration | Airflow, Dagster, Argo, Temporal |
| Transform | dbt, Spark, DuckDB |
| Streaming | Kafka, Redpanda, Redis Streams |
| Query | Trino, DuckDB, PostgreSQL |
| Lineage | OpenLineage, Marquez, OpenMetadata, DataHub |
| Notebooks | JupyterHub |
| Monitoring | Prometheus, Grafana, OpenTelemetry |
| AI workflows | LangGraph, LlamaIndex, CrewAI where useful |
| Code agents | Codex, Claude, Cursor, Roo, Aider, Continue |

---

## 19. AI Development Workflow for Building Fusion

The user wants to build heavily with AI.

### 19.1 Desired AI Team Concept

The conceptual AI team:

```text
User = final lead / product owner / architect
Ruflo or team manager = AI team manager
Claude Sonnet = technical lead / reviewer
Qwen / Kimi / DeepSeek = developers / testers
Roo / Aider = hands that edit code
Continue + Codebase-Memory MCP = memory/context brain
Playwright MCP = UI tester
OpenRouter = model router
.ai folder = permanent project knowledge
Codex = repo-aware implementation agent
```

### 19.2 Important Working Pattern

The user wants to act like a lead:

- Split tasks.
- Define expected output.
- Assign work.
- Review results.
- Iterate.

AI tools should not randomly edit the whole repo.

### 19.3 `.ai` Folder

Recommended `.ai` folder:

```text
.ai/
  PROJECT_CONTEXT.md
  CODING_RULES.md
  CURRENT_ROADMAP.md
  ARCHITECTURE_DECISIONS.md
  SERVICE_BOUNDARIES.md
  DOMAIN_MODEL.md
  API_CONTRACTS.md
  UI_REQUIREMENTS.md
  TESTING_STRATEGY.md
  AGENT_PROMPTS.md
  KNOWN_RISKS.md
```

### 19.4 Codex Instructions

Codex should:

- Read `.ai/PROJECT_CONTEXT.md` first.
- Read `README.md` second.
- Inspect repo structure before editing.
- Make small commits/changes.
- Explain what changed.
- Add tests where possible.
- Avoid secret leakage.
- Avoid implementing unsafe AI auto-execution.
- Respect service boundaries.

---

## 20. MVP Direction

The full vision is huge. The MVP should prove the control plane concept.

### 20.1 MVP Goal

Build a working Fusion Control Plane Kernel that can:

1. Create tenants/projects.
2. Register connections.
3. Discover metadata from a source.
4. Define a simple workflow.
5. Trigger workflow through adapter.
6. Track run state.
7. Store audit logs.
8. Show UI dashboard.
9. Integrate one execution engine.
10. Demonstrate AI recommendation draft without auto-execution.

### 20.2 MVP Scope Option A — CDC-Centric MVP

Use Fusion CDC as core demo.

Features:

- Register PostgreSQL/MySQL source.
- Discover schemas/tables.
- Enable CDC/incremental sync for selected table.
- Stream changes to destination.
- Track checkpoint.
- Show run/lag/errors.
- Add simple transformation.
- Add quality rule.

Pros:

- Strong data engineering story.
- Shows real platform capability.
- Differentiates from simple catalog tools.

Cons:

- More engineering complexity.
- CDC correctness is hard.

### 20.3 MVP Scope Option B — BI/Semantic Control Plane MVP

Use Superset + semantic/catalog layer as demo.

Features:

- Register warehouse connection.
- Discover datasets.
- Create business dataset profile.
- Define metrics.
- Define relationships.
- Generate Superset chart/dashboard.
- Add certification workflow.
- Add AI chart recommendation.

Pros:

- Easier to demo visually.
- Faster user-facing value.
- Good for first customers.

Cons:

- Less infrastructure-depth than CDC MVP.

### 20.4 MVP Scope Option C — Kernel + One Adapter

Build only kernel and one adapter first.

Features:

- Project/tenant model.
- Connection registry.
- Workflow registry.
- Run registry.
- Audit.
- Airflow or dbt adapter.
- Simple UI.

Pros:

- Clean foundation.
- Avoids overbuilding.
- Makes future layers easier.

Cons:

- Demo may feel less magical unless UI is strong.

### 20.5 Recommended MVP

Recommended path:

1. Build Kernel + UX shell.
2. Add connection registry and metadata discovery.
3. Add dbt or Airflow adapter.
4. Add audit/run tracking.
5. Add Superset/BI integration OR CDC demo as second milestone.
6. Add AI Operator only as recommendation/draft mode.

---

## 21. Four-Month Build Plan

The user wants to use open-source repos, AI coding tools, and complete as much as possible in around four months.

A realistic four-month plan should target a strong MVP, not the entire 16-layer platform.

### Month 1 — Foundation

Goals:

- Repo setup.
- Architecture docs.
- Auth skeleton.
- Tenant/project model.
- Control-plane metadata DB.
- Event envelope.
- Audit service.
- Basic web shell.
- Local Docker Compose.

Deliverables:

- Running web app.
- Running backend.
- PostgreSQL metadata DB.
- Login mock or Keycloak integration.
- Tenant/project CRUD.
- Audit events.

### Month 2 — Connections and Workflow Kernel

Goals:

- Connection registry.
- Connector definitions.
- Source connection test.
- Metadata discovery.
- Workflow registry.
- Workflow versioning.
- Run service.
- Adapter SDK.

Deliverables:

- Register PostgreSQL connection.
- Discover schema/table/columns.
- Create workflow draft.
- Trigger mock adapter.
- Track run status.
- Show logs/events.

### Month 3 — First Real Engine Integration

Choose one:

- dbt adapter.
- Airflow adapter.
- Superset adapter.
- CDC worker adapter.

Deliverables:

- Trigger real job.
- Fetch status.
- Fetch artifacts/logs.
- Emit events into Fusion.
- Show run timeline.
- Add simple policy/approval gate.

### Month 4 — Product Demo Layer

Goals:

- Strong UI/UX.
- AI recommendation draft.
- Semantic/catalog or CDC demo.
- Documentation.
- Deployment manifests.
- E2E tests.

Deliverables:

- Demo-ready MVP.
- Project README.
- Architecture docs.
- Playwright tests.
- Docker/Helm deployment.
- Example seed data.
- Demo script.

---

## 22. Service Boundary Details

### 22.1 Identity Service

Responsibilities:

- Auth integration.
- User profile sync.
- Role assignments.
- Service accounts.
- Token validation.

Do not mix identity with workflow logic.

### 22.2 Tenant Service

Responsibilities:

- Organisations.
- Tenants.
- Parent-child tenant hierarchy.
- Tenant settings.
- Tenant isolation metadata.

### 22.3 Project Service

Responsibilities:

- Workspaces/projects.
- Project environments.
- Project members.
- Project settings.

### 22.4 Connection Registry Service

Responsibilities:

- Connections.
- Connector definitions.
- Credentials references.
- Connection tests.
- Metadata discovery jobs.

### 22.5 Workflow Registry Service

Responsibilities:

- Workflow definitions.
- Versions.
- Promotion states.
- Dependency metadata.

### 22.6 Run Service

Responsibilities:

- Workflow runs.
- Task runs.
- Run status.
- Run timeline.
- Retry state.
- Logs/artifact references.

### 22.7 Policy Service

Responsibilities:

- Approval requirements.
- Environment policies.
- Action policies.
- Data movement policies.
- AI action policies.

### 22.8 Audit Service

Responsibilities:

- Immutable audit events.
- Queryable audit timeline.
- Compliance export.

### 22.9 Adapter Registry Service

Responsibilities:

- Registered adapters.
- Capabilities.
- Health.
- Versions.
- Config schema.

### 22.10 AI Assist Service

Responsibilities:

- Context retrieval.
- Prompt orchestration.
- Recommendation generation.
- Action draft creation.
- Safety validation.
- Approval integration.

---

## 23. Adapter Framework

Adapters are central to the architecture.

### 23.1 Adapter Types

```text
OrchestrationAdapter
TransformationAdapter
CDCAdapter
BIAdapter
CatalogAdapter
QualityAdapter
LineageAdapter
StorageAdapter
QueryAdapter
NotificationAdapter
LLMAdapter
```

### 23.2 Common Adapter Capabilities

Each adapter should expose:

```text
getCapabilities()
validateConfig(config)
testConnection(config)
getHealth()
execute(command)
getStatus(runId)
fetchLogs(runId)
fetchArtifacts(runId)
cancel(runId)
```

Not every adapter needs every function, but capability discovery should state what is supported.

### 23.3 Adapter Command Example

```json
{
  "command_id": "cmd_123",
  "command_type": "dbt.run",
  "tenant_id": "tenant_123",
  "project_id": "project_123",
  "environment": "dev",
  "requested_by": "user_123",
  "config": {
    "project_ref": "dbt_project_1",
    "models": ["customer_mart"],
    "full_refresh": false
  }
}
```

---

## 24. UI / UX Direction

### 24.1 UI Philosophy

The UI should feel like a modern data platform command centre.

It should not look like a collection of admin forms.

### 24.2 Main Navigation

Suggested navigation:

```text
Home / Command Center
Tenants
Projects
Connections
Workflows
Runs
Data Products
Catalog
Quality
Lineage
BI & Reports
Policies
Approvals
AI Operator
Settings
```

### 24.3 Command Center

Show:

- Active runs.
- Failed runs.
- SLA breaches.
- Quality failures.
- Cost warnings.
- AI recommendations.
- Pending approvals.
- Tenant health.
- Engine health.

### 24.4 Connections UI

Features:

- Add source.
- Add destination.
- Test connection.
- Browse schemas.
- Browse tables.
- Select tables.
- Set sync mode.
- Set credentials securely.
- Show last sync metadata.

### 24.5 Workflow UI

Features:

- Workflow list.
- Workflow designer.
- Version history.
- Promotion state.
- Run button.
- Schedule.
- Dependencies.
- Approval panel.
- Run timeline.

### 24.6 AI Operator UI

Must show:

- What AI looked at.
- What AI recommends.
- Why it recommends it.
- Risk level.
- Diff/change preview.
- Required approvals.
- Execute button only after approval.
- Audit trail.

---

## 25. Coding Rules for Codex

### 25.1 General Rules

- Do not hardcode tenant IDs.
- Do not hardcode secrets.
- Do not mix tenant data.
- Every write action must include actor and tenant context.
- Every significant action should emit audit event.
- Use typed contracts.
- Add tests for domain logic.
- Prefer clear service boundaries.
- Avoid premature microservice complexity if repo is still early.

### 25.2 Backend Rules

- Validate all inputs.
- Use DTOs/contracts for APIs.
- Use database migrations.
- Add correlation IDs.
- Add structured logging.
- Add health endpoints.
- Add readiness/liveness endpoints.
- Keep secrets in secret references.
- Enforce tenant filters in repository queries.

### 25.3 Frontend Rules

- Use TypeScript strictly.
- Use typed API clients.
- Show loading/error/empty states.
- Do not expose raw secrets.
- Always include tenant/project context in UI state.
- Make dangerous actions require confirmation.
- Make AI-generated changes visually distinct.

### 25.4 Testing Rules

Required test categories:

- Unit tests.
- Domain tests.
- API tests.
- Integration tests.
- Contract tests for adapters.
- E2E tests for critical flows.
- Authorization tests.
- Tenant isolation tests.

### 25.5 Security Rules

- Never log secrets.
- Never return secrets in API responses.
- Redact credentials.
- Mask sensitive fields.
- Validate authorization on backend, not only frontend.
- Add audit events for access and changes.

---

## 26. Example API Resources

### 26.1 Tenant APIs

```http
POST /api/tenants
GET /api/tenants
GET /api/tenants/{tenantId}
PATCH /api/tenants/{tenantId}
```

### 26.2 Project APIs

```http
POST /api/projects
GET /api/projects?tenantId=...
GET /api/projects/{projectId}
PATCH /api/projects/{projectId}
```

### 26.3 Connection APIs

```http
POST /api/connections
GET /api/connections
GET /api/connections/{connectionId}
POST /api/connections/{connectionId}/test
POST /api/connections/{connectionId}/discover
GET /api/connections/{connectionId}/metadata
```

### 26.4 Workflow APIs

```http
POST /api/workflows
GET /api/workflows
GET /api/workflows/{workflowId}
POST /api/workflows/{workflowId}/versions
POST /api/workflows/{workflowId}/submit-review
POST /api/workflows/{workflowId}/approve
POST /api/workflows/{workflowId}/promote
POST /api/workflows/{workflowId}/runs
```

### 26.5 Run APIs

```http
GET /api/runs
GET /api/runs/{runId}
GET /api/runs/{runId}/timeline
GET /api/runs/{runId}/logs
POST /api/runs/{runId}/cancel
POST /api/runs/{runId}/retry
```

### 26.6 AI APIs

```http
POST /api/ai/recommendations
GET /api/ai/recommendations/{recommendationId}
POST /api/ai/recommendations/{recommendationId}/approve
POST /api/ai/recommendations/{recommendationId}/reject
```

---

## 27. Domain Model Draft

### 27.1 Organization

```json
{
  "id": "org_123",
  "name": "Example Bank",
  "status": "active",
  "created_at": "timestamp"
}
```

### 27.2 Tenant

```json
{
  "id": "tenant_123",
  "organization_id": "org_123",
  "name": "Tenant A",
  "type": "customer|sub_bank|business_unit",
  "status": "active"
}
```

### 27.3 Connection

```json
{
  "id": "conn_123",
  "tenant_id": "tenant_123",
  "project_id": "project_123",
  "name": "Production Postgres",
  "connector_type": "postgresql",
  "direction": "source|destination|both",
  "secret_ref": "vault://...",
  "status": "active",
  "last_test_status": "success"
}
```

### 27.4 Workflow

```json
{
  "id": "wf_123",
  "tenant_id": "tenant_123",
  "project_id": "project_123",
  "name": "Customer CDC Sync",
  "type": "cdc|batch|dbt|report|quality|custom",
  "current_version_id": "wfv_123",
  "status": "draft|active|paused|archived"
}
```

### 27.5 Workflow Run

```json
{
  "id": "run_123",
  "workflow_id": "wf_123",
  "version_id": "wfv_123",
  "status": "queued|running|success|failed|cancelled|retrying",
  "started_at": "timestamp",
  "ended_at": "timestamp",
  "triggered_by": "user_123",
  "engine": "airflow|dbt|spark|duckdb|cdc-worker|superset"
}
```

### 27.6 Audit Event

```json
{
  "id": "audit_123",
  "tenant_id": "tenant_123",
  "project_id": "project_123",
  "actor_type": "user|service|ai|system",
  "actor_id": "user_123",
  "action": "workflow.approved",
  "resource_type": "workflow",
  "resource_id": "wf_123",
  "timestamp": "timestamp",
  "metadata": {}
}
```

---

## 28. Known Risks

### 28.1 Scope Risk

The full 16-layer product is very large. Building all of it in four months is unrealistic unless using many existing components and narrowing MVP scope.

Mitigation:

- Build kernel first.
- Choose one hero workflow.
- Use adapters.
- Avoid building every engine.

### 28.2 Integration Complexity

Integrating Superset, OpenMetadata, dbt, Airflow, Spark, Kafka, and AI tools can become messy.

Mitigation:

- Define clean adapter contracts.
- Use event envelope.
- Avoid direct deep coupling.
- Keep integration-specific code isolated.

### 28.3 Security Risk

Control plane touches sensitive metadata, credentials, policies, and operational actions.

Mitigation:

- Secret references only.
- RBAC from day one.
- Audit logs from day one.
- Tenant isolation tests.

### 28.4 AI Safety Risk

AI may generate unsafe or wrong actions.

Mitigation:

- AI only drafts.
- Human approval required.
- Policy validation.
- No auto-prod execution.
- Full diff and audit.

### 28.5 Open-Source License Risk

Some open-source tools may have restrictive licenses.

Mitigation:

- Prefer API integration.
- Avoid copying code.
- Review license before embedding/forking.

### 28.6 Performance Risk

Control plane can become slow if it tries to process raw data or high-volume streams directly.

Mitigation:

- Keep raw data in execution plane.
- Store metadata and state.
- Use asynchronous jobs.
- Use projections/read models.

---

## 29. What Codex Should Build First

If repository is empty or early-stage, build this order:

### Step 1 — Repo Skeleton

Create:

```text
apps/web
services/api
packages/domain
packages/event-envelope
packages/authz
infra/local-dev
docs/architecture
.ai
```

### Step 2 — Metadata DB

Create PostgreSQL schema for:

- organizations
- tenants
- projects
- connections
- workflows
- workflow_versions
- workflow_runs
- task_runs
- audit_events
- system_events
- approval_requests

### Step 3 — Backend API

Build APIs for:

- tenant CRUD
- project CRUD
- connection CRUD
- audit event write/read
- workflow CRUD
- run state mock

### Step 4 — Frontend Shell

Build pages:

- Command Center
- Tenants
- Projects
- Connections
- Workflows
- Runs
- Audit

### Step 5 — Adapter SDK

Build a simple adapter interface and mock adapter.

### Step 6 — Real Adapter

Pick one:

- dbt CLI adapter.
- Airflow REST adapter.
- Superset API adapter.
- PostgreSQL metadata discovery adapter.

Recommended first real adapter:

- PostgreSQL metadata discovery adapter.

Why:

- Easier than full CDC.
- Useful for catalog, BI, and workflow setup.
- Demonstrates connection registry value.

### Step 7 — AI Draft Mode

Add only after kernel is working.

AI should produce recommendations but not execute.

---

## 30. Initial Codex Prompts

### 30.1 Repo Inspection Prompt

```text
Read PROJECT_CONTEXT.md and inspect the repository. Summarize what currently exists, what is missing compared with the Fusion Control Plane Kernel, and propose the next 5 implementation tasks. Do not edit code yet.
```

### 30.2 Skeleton Creation Prompt

```text
Using PROJECT_CONTEXT.md, create the initial monorepo skeleton for DCraft Fusion. Add apps/web, services/api, packages/domain, packages/event-envelope, packages/authz, infra/local-dev, docs/architecture, and .ai. Add README files explaining each directory. Do not implement business logic yet.
```

### 30.3 Domain Model Prompt

```text
Implement the first domain model package for Organization, Tenant, Project, Connection, Workflow, WorkflowVersion, WorkflowRun, ApprovalRequest, AuditEvent, and SystemEvent. Use strong typing, validation, and tenant-aware fields. Add unit tests.
```

### 30.4 Backend Prompt

```text
Build the initial backend API for tenants, projects, connections, workflows, runs, and audit events. Use PostgreSQL migrations, structured logging, correlation IDs, and input validation. Ensure every write action records an audit event.
```

### 30.5 Frontend Prompt

```text
Build the first React/TypeScript Fusion web shell with navigation for Command Center, Tenants, Projects, Connections, Workflows, Runs, Audit, and AI Operator. Use placeholder data if backend is not ready, but keep API client boundaries clean.
```

### 30.6 Adapter Prompt

```text
Create the Adapter SDK with capability discovery, config validation, health check, execute command, get status, fetch logs, fetch artifacts, and cancel methods. Implement a mock adapter and one PostgreSQL metadata discovery adapter.
```

### 30.7 AI Safety Prompt

```text
Implement AI recommendation draft mode only. AI can generate suggested workflow or quality rule changes, but cannot execute them. Add approval state, diff preview, and audit events for approve/reject. Do not add auto-execution.
```

---

## 31. Architecture Decision Records to Create

Create ADR files for:

```text
ADR-001-control-plane-vs-data-plane.md
ADR-002-non-data-owning-control-plane.md
ADR-003-adapter-based-execution.md
ADR-004-event-envelope-and-outbox.md
ADR-005-ai-operator-human-approval.md
ADR-006-superset-as-bi-shell.md
ADR-007-openmetadata-or-datahub-integration.md
ADR-008-duckdb-vs-spark-transform-strategy.md
ADR-009-deployment-modes-saas-byoc-onprem.md
ADR-010-tenant-isolation-model.md
```

---

## 32. Product Positioning

### 32.1 What Fusion Is

Fusion is:

- A control plane.
- A productised operating layer for data platforms.
- A governance and visibility layer.
- An AI-assisted data operations platform.
- A metadata-first coordination plane.
- A platform for fragmented data teams to standardise workflows.

### 32.2 What Fusion Is Not

Fusion is not:

- Just another ETL tool.
- Just a dashboard tool.
- Just a catalog.
- Just an orchestration tool.
- Just an AI chatbot.
- Just a wrapper around Superset.
- Just a wrapper around Airflow.
- A mandatory data warehouse.

### 32.3 Competitor Categories

Fusion overlaps with many categories:

- Databricks workflows/lakehouse.
- Airbyte/Fivetran ingestion.
- dbt Cloud transformations and semantic layer.
- Airflow/Dagster orchestration.
- Monte Carlo/Bigeye observability.
- OpenMetadata/DataHub catalog.
- Power BI/Tableau/Superset BI.
- Atlan/Collibra governance.
- Internal platform portals.

But Fusion's differentiation is the unified control plane across layers, with optional AI Operator and non-data-owning deployment model.

---

## 33. First Customer Profile

Ideal first customers:

- Mid-size fintechs.
- Banks with fragmented data tooling.
- SaaS companies with multi-tenant data platforms.
- Teams using Airflow + dbt + Superset + Kafka but lacking unified control.
- Teams that cannot move data into external SaaS because of compliance.
- Teams struggling with CDC, reporting, lineage, quality, and operational visibility.

Pain points:

- Tool sprawl.
- No unified visibility.
- Manual onboarding.
- Weak governance.
- Repeated pipeline failures.
- No tenant-level operational control.
- BI users confused by raw tables.
- AI cannot be trusted because metadata is poor.

---

## 34. Demo Narrative

A strong demo should show:

1. Create a tenant/project.
2. Register a data source.
3. Discover schemas/tables.
4. Define a workflow or data product.
5. Add a quality rule.
6. Trigger a run through an adapter.
7. View run status and logs.
8. See lineage/metadata update.
9. Ask AI Operator for improvement.
10. AI proposes a safe change.
11. User approves.
12. Fusion executes through adapter.
13. Audit trail proves everything.

The demo should communicate:

> Fusion is the brain. The existing tools are the muscles. The customer keeps control of data and execution.

---

## 35. Current Canonical Summary

DCraft Fusion should be built as an enterprise-grade, AI-assisted, metadata-first control plane for modern data platforms. It should unify the operating experience across ingestion, CDC, orchestration, transformation, quality, lineage, governance, BI/reporting, observability, cost, and AI-assisted operations.

The architecture must preserve a strict separation between:

- Control plane: metadata, intent, policies, approvals, audit, run state, visibility.
- Data/execution plane: actual data movement, transformation, query execution, stream processing, BI rendering, and compute.

The most important build priority is the Control Plane Kernel, not the entire 16-layer vision at once.

The strongest first implementation path is:

1. Build the kernel.
2. Add connection registry.
3. Add metadata discovery.
4. Add workflow/run registry.
5. Add audit and approval.
6. Add one real adapter.
7. Add UI shell.
8. Add AI recommendation draft mode.
9. Expand into CDC, BI, quality, lineage, and governance modules.

---

## 36. Final Instruction for Codex

When working on this repository, do not treat Fusion as a normal CRUD SaaS app. Treat it as a serious data platform control plane.

Every implementation decision should preserve:

- Tenant isolation.
- Non-data-owning control plane.
- Adapter-based execution.
- Auditability.
- Human approval for AI actions.
- Open-source integration strategy.
- Clear service boundaries.
- Future compatibility with SaaS, BYOC, and on-prem deployments.

If uncertain, implement the smallest safe kernel feature and document the decision.



---

# APPENDIX A — Discussion Memory Expansion Pack (Added After User Feedback)

> This appendix exists because the first Codex handoff file was detailed, but it still missed many specific discussions from the wider Startup project.  
> Treat this appendix as a **conversation-to-Codex translation layer**. It captures the important decisions, side discussions, constraints, implementation preferences, and product ideas that were discussed across the project but may not have been fully reflected in the original sections above.

## A0. Important Limitation and How to Use This Appendix

This file is not a literal export of every chat message. It is a consolidated engineering memory created from the project context visible to ChatGPT, uploaded files, and remembered decisions. Codex should use it as the operating memory for implementation.

When there is a conflict:

1. Prefer the latest frozen decision.
2. Prefer explicit constraints over suggestions.
3. Prefer the non-data-owning control-plane model.
4. Prefer building a working Control Plane Kernel before attempting the full 16-layer product.
5. Do not remove earlier features unless the user explicitly says they are deprecated.

The user repeatedly said: **do not leave out features from the initial idea; keep adding on top.** Therefore, Codex should treat old DataCortex, Fusion, Reporting, Superset BI, CDC Engine, and AI-agent workflow discussions as product inspiration unless explicitly marked as separate.

---

## A1. Product Evolution Timeline

The project evolved through multiple names and angles. Codex must understand this history so it does not treat them as unrelated ideas.

### A1.1 DataCortex Phase

Earlier umbrella name: **DataCortex**.

Core idea at that stage:

- A Data Engineering OS.
- Contract-first data infrastructure.
- LLM-native automation.
- Lake-building and observability-first platform.
- Data contracts, validation, schema evolution, CDC/batch sync, lineage, governance, documentation, audit, GitOps, streaming, RBAC, and self-improvement.

Important modules from DataCortex that must not be lost:

- Data Contracts.
- Validation and Testing.
- Smart Sync Engine covering CDC plus batch.
- Lake Construction.
- Schema Evolution and Diffs.
- Lineage and Metadata Graph.
- Observability.
- Refactoring Intelligence.
- Documentation Engine.
- RBAC and Governance.
- LLM-native interface.
- Self-improvement layer.
- Data Packaging and Sharing.
- Audit Trail.
- GitOps / PR workflow.
- Streaming Support.
- LLM-safe data access and metadata-safe prompting.

### A1.2 Fusion Phase

The user also discussed **Fusion** in the context of:

- CDC Engine.
- Reporting Platform.
- Apache Superset fork.
- Data engineering execution infrastructure.
- Multi-tenant banking data platform.

Fusion should be treated as source inspiration and possible execution-plane module history, not necessarily the final product brand.

### A1.3 DCraft Fusion Frozen Phase

Final frozen brand and canonical architecture:

- Company: **DCraft Labs**.
- Product: **DCraft Fusion**.
- Short name: **Fusion**.
- Control plane: **DCraft Fusion Control Plane**.
- AI capability: **AI Operator**.

Frozen product definition:

> DCraft Fusion is an AI-assisted, engine-agnostic control plane for modern data platforms. It sits above existing data tools and centralizes intent, metadata, policy, visibility, governance, operational workflow, and auditable AI assistance while execution remains in customer engines.

Core principle:

> Tightly coupled control plane / UX. Loosely coupled execution engines.

---

## A2. Non-Negotiable Product Boundaries

These are hard constraints.

### A2.1 Fusion Is Not a Data Warehouse

Fusion must not become a warehouse that stores all raw customer data.

Fusion stores:

- Metadata.
- Intent.
- Configuration.
- Policies.
- Ownership.
- Glossary.
- Lineage metadata.
- Quality rule definitions.
- Quality results / summaries.
- Workflow definitions.
- Run state.
- Audit logs.
- AI proposals.
- Approval history.
- Cost and usage summaries.
- Observability signals.

Fusion does not store by default:

- Raw table rows.
- Full customer datasets.
- Stream payload contents.
- Query result dumps.
- Sensitive production data.
- Data files from the customer's lake or warehouse.

Exception: customer explicitly enables a feature that moves or stores data, such as profiling samples, cached previews, or managed execution mode. Even then, it must be governed and configurable.

### A2.2 Execution Remains Pluggable

Fusion coordinates execution through adapters and agents. It does not hard-code one compute engine.

Execution may happen in:

- Customer Airflow.
- Customer dbt Cloud / dbt Core.
- Customer Spark.
- Customer Flink.
- Customer Kubernetes.
- Customer Trino.
- Customer Iceberg catalog.
- Customer warehouse.
- Customer Superset.
- Customer OpenMetadata/DataHub.
- Fusion default engine for greenfield users.

### A2.3 AI Operator Is Optional and Audited

AI Operator must never silently change production pipelines.

AI Operator can:

- Read metadata.
- Explain systems.
- Suggest improvements.
- Draft YAML/SQL/dbt changes.
- Generate quality rules.
- Draft PRs.
- Create remediation plans.
- Summarize incidents.
- Recommend cost optimizations.

AI Operator must not:

- Execute production changes without policy permission.
- Bypass approval workflows.
- Read raw customer data unless explicitly granted.
- Generate unrestricted SQL directly against raw tables for business users.
- Ignore tenant isolation.
- Ignore audit logging.

---

## A3. The 16-Layer Product Scope — Expanded from Discussion

The frozen 16-layer architecture is the product spine. Each layer can support three modes:

1. **BYO mode** — customer already has a tool; Fusion controls and observes through adapter.
2. **Fusion-default mode** — Fusion provides a default open-source-backed implementation.
3. **Hybrid mode** — some components customer-owned, some Fusion-managed.

### A3.1 Layer 1 — Ingestion and Sync

Scope:

- Source registration.
- Destination registration.
- Batch ingestion.
- API ingestion.
- File ingestion.
- Airbyte integration.
- Connector marketplace.
- Incremental syncs.
- Initial full load.
- Retry/backfill management.
- Tenant-aware ingestion configuration.

Relevant discussion points:

- User liked Airbyte-like UX: select tenant → setup source → setup destination → discover schema/tables → configure sync.
- Sources mentioned: MySQL, PostgreSQL, MongoDB.
- Destinations mentioned: PostgreSQL warehouse, Apache Iceberg on object storage, S3/Azure storage, Trino-accessible lake.
- Need test connection, schema browsing, per-table enable/disable.
- Need JSON flattening options at source/destination/connection level.

Possible open-source integrations:

- Airbyte.
- Meltano/Singer-style connectors.
- Debezium for log-based CDC.
- Custom Python CDC for early prototype.

### A3.2 Layer 2 — CDC and Streaming

Scope:

- Real-time change capture.
- Log-based CDC.
- Kafka/Debezium support.
- Redis Streams prototype path from Fusion.
- Schema Registry.
- CDC event normalization.
- Checkpointing.
- Delete handling.
- Deduplication.
- Tenant-level pause/resume.
- Backpressure and throttling.

Important user history:

- User has production experience with Debezium + Kafka Connect across MySQL, PostgreSQL, and MongoDB.
- Previous Kafka/Debezium cluster meltdown was discussed: data surge filled brokers and impacted downstream systems.
- User explored Redis Streams as a lighter central event backbone.
- User discussed a Python CDC engine with Redis Streams and SQLite checkpoints.
- User wanted per-tenant orchestration and throttling, with 2–3 connections per pod configurable.
- Need careful design because one connector may read a full database/server-level binlog/WAL while tenants may be separated by schema/database.

CDC event pattern from Fusion inspiration:

```text
cdc:<bank>:<tenant>:<source>:<schema>:<table>
```

Worker flow:

1. Connect to source database.
2. Read binlog/WAL/change stream.
3. Normalize event.
4. Add tenant metadata and timestamps.
5. Publish to stream.
6. Store checkpoints locally and centrally.
7. Route to transformation/loader.

### A3.3 Layer 3 — Orchestration and Scheduling

Scope:

- Workflow definitions.
- DAG scheduling.
- Run triggering.
- Manual runs.
- Backfills.
- Retries.
- Failure policies.
- Dependency mapping.
- Airflow integration.
- Event-driven orchestration.

Discussion points:

- User has used Airflow with S3-synced DAGs and Spark cluster in client mode.
- Fusion should not replace Airflow immediately; it should control, observe, and generate workflows through adapter.
- Need event-first control-plane architecture.
- Run Service must track run state independent of execution engine.

### A3.4 Layer 4 — Transformation

Scope:

- dbt integration.
- SQL transformations.
- Spark/PySpark transformations.
- Python transformations.
- Low-code transformation builder.
- Transformation rules at column, table, connection, and pipeline level.
- Versioned transformation models.
- Transformation lineage.

Discussion points:

- User discussed a 10-step transformation engine in Fusion CDC context.
- Some reviewed project claimed no Spark; user questioned whether that was realistic.
- Fusion should support both Spark and non-Spark engines depending on scale and customer setup.
- Transformation layer should be engine-agnostic: dbt for SQL models, Spark for heavy workloads, Python for custom logic, Flink for streaming where needed.

### A3.5 Layer 5 — Data Quality and Contracts

Scope:

- Data contracts.
- Schema expectations.
- Freshness checks.
- Null/unique/range checks.
- Referential integrity checks.
- Business rules.
- Contract versioning.
- Contract enforcement modes: warn, block, quarantine.
- Quality score.
- Quality history.
- Quality gates in CI/CD.

Discussion points:

- DataCortex originally had contract-first design.
- Superset BI blueprint requires data quality/freshness scores shown to users.
- Data quality should feed AI Operator context.
- Great Expectations or Soda Core can be integrated, but custom profiling should come first.

### A3.6 Layer 6 — Schema Evolution

Scope:

- Schema drift detection.
- Backward/forward compatibility.
- Breaking change detection.
- Schema Registry integration.
- Avro/JSON schema compatibility.
- dbt model impact analysis.
- BI dashboard impact analysis.
- Approval workflow for breaking changes.

Discussion points:

- User has production CDC with Confluent Schema Registry and Avro.
- Need schema governance to protect downstream consumers and dbt models.
- User wants versioned dbt models and schema evolution.

### A3.7 Layer 7 — Lineage

Scope:

- Table-level lineage.
- Column-level lineage.
- Pipeline lineage.
- Dashboard lineage.
- Metric lineage.
- CDC-to-dbt-to-BI lineage.
- Impact analysis.
- OpenLineage support.
- Integration with OpenMetadata or DataHub.

Discussion points:

- User repeatedly values lineage for auditability and incident resolution.
- Superset fork should show lineage and trust context.
- Fusion should build a metadata graph rather than isolated configuration tables.

### A3.8 Layer 8 — Governance and Ownership

Scope:

- Data owners.
- Stewards.
- Domain ownership.
- Glossary.
- Certification workflow.
- Approvals.
- Review process.
- Auditability.
- Compliance metadata.

Discussion points:

- Superset BI blueprint requires certified datasets/metrics and owner approval.
- Business users should start from certified concepts, not raw tables.
- Governance must be product-native, not a PDF process outside the tool.

### A3.9 Layer 9 — Access and Policy Enforcement

Scope:

- RBAC.
- ABAC.
- Tenant isolation.
- Row-level policies.
- Column masking.
- PII sensitivity.
- Secrets management.
- Policy-as-code.
- Keycloak integration.

Discussion points:

- User has forked Superset with Keycloak SSO and RBAC.
- Reporting platform requires bank/sub-bank/tenant isolation.
- Real-world Keycloak SSO constraints were discussed: report-manager realm, VisapayAgent realm, identity brokering, OIDC, hardcoded mappers for attributes like bank_id and sub_bank_ids where needed.
- Must avoid disabling auth just to make SSO easier.

### A3.10 Layer 10 — Observability: Health and SLAs

Scope:

- Pipeline health.
- Job status.
- Lag.
- Failure rate.
- SLA miss.
- Freshness.
- Volume anomalies.
- Data quality trend.
- Runtime trend.
- Prometheus/Grafana integration.

Discussion points:

- Fusion included Prometheus + Grafana.
- Reporting platform requires WebSocket progress updates.
- User has built reusable Python alerting layer.
- AI Operator should use observability context for diagnosis.

### A3.11 Layer 11 — Cost and Usage Attribution

Scope:

- Per-tenant usage.
- Per-pipeline cost.
- Compute attribution.
- Storage cost.
- Query cost.
- Idle resources.
- Cost anomalies.
- FinOps recommendations.

Discussion points:

- User achieved 25% infrastructure cost reduction by replacing or optimizing platforms.
- JupyterHub was used as Databricks replacement to reduce cost.
- Fusion should show cost by tenant, workload, engine, and data product.

### A3.12 Layer 12 — Incident Detection

Scope:

- Failure grouping.
- Alert deduplication.
- Error classification.
- Incident timeline.
- Business impact mapping.
- Auto-ticket creation.
- Notification routing.

Discussion points:

- Reporting platform needs error group-based reruns.
- User wants robust retry semantics beyond what Kafka alone provides.
- Need to distinguish OOM, transient DB issue, bad query, bad data, schema drift, permission issue, downstream unavailable.

### A3.13 Layer 13 — Root Cause Analysis

Scope:

- Dependency impact tree.
- Recent changes.
- Run logs.
- Data drift.
- Schema drift.
- Downstream dashboard breakage.
- AI-generated RCA summary.
- Suggested fix.

Discussion points:

- AI Operator should explain failures and propose remediation.
- Lineage, observability, schema history, and audit logs must feed RCA.

### A3.14 Layer 14 — Reverse ETL and Activation

Scope:

- Publish curated data to business tools.
- CRM activation.
- Webhook delivery.
- API delivery.
- SFTP/S3/email delivery.
- Governed export.
- Audit of outbound movement.

Discussion points:

- Reporting/delivery platform includes email, SFTP, S3, Azure/blob delivery.
- PDF generation service and scheduled report delivery are relevant execution modules.
- Reverse ETL must enforce policies and approvals.

### A3.15 Layer 15 — CI/CD and Environments

Scope:

- GitOps.
- PR workflow.
- Environment promotion.
- Dev/stage/prod.
- Contract validation in CI.
- dbt checks in CI.
- Policy checks.
- AI-generated PRs.
- Rollback.

Discussion points:

- DataCortex included GitOps/PR workflow.
- User wants AI agents to edit code, but human remains final lead/reviewer.
- Codex should create repo structures that support staged development, not random files.

### A3.16 Layer 16 — AI Operator

Scope:

- Metadata-aware assistant.
- Safe code generation.
- Pipeline design suggestions.
- Contract generation.
- Query explanation.
- RCA.
- Cost optimization.
- Data onboarding assistant.
- BI assistant.
- PR drafting.
- Policy-aware actions.

Discussion points:

- AI Operator is a product feature, not just an internal coding helper.
- It should use RAG over metadata, glossary, lineage, schema, run state, audit logs, and docs.
- It must be optional per tenant/project.
- Every recommendation should be explainable and auditable.

---

## A4. Control Plane Kernel Skeleton — Immediate Build Focus

The user explicitly decided that the next step is **not more documents**, but the implementation of the Control Plane Kernel Skeleton.

Codex should prioritize this build.

### A4.1 Core Services to Create

These service names were explicitly defined and should be surfaced in the repo:

```text
identity-service
tenant-service
project-service
registry-service
run-service
connection-service
policy-service
audit-service
event-gateway
projection-service
```

### A4.2 Service Responsibilities

#### identity-service

Purpose:

- Users.
- Groups.
- Roles.
- Service accounts.
- Auth provider integration.
- Keycloak/OIDC support.

Should not:

- Store raw customer data.
- Execute pipelines.

#### tenant-service

Purpose:

- Tenant hierarchy.
- Organization/bank/sub-bank/customer mapping.
- Tenant settings.
- Deployment mode.
- Isolation boundaries.

Important examples from discussions:

- Bank: Visa.
- Tenants/Sub-banks: FBN, Solidaire, Access, UBA, SOFI.
- Need bank_id, bank_name, sub_bank_ids, sub_bank_names claims/attributes in some portal flows.

#### project-service

Purpose:

- Project/workspace model.
- Data product grouping.
- Environment mapping.
- Ownership.
- Domain mapping.

#### registry-service

Purpose:

- Central resource registry.
- Resource definitions.
- Resource versions.
- Metadata object models.
- Adapter registrations.
- Provider capabilities.

Resource examples:

- Connection.
- Dataset.
- Metric.
- Workflow.
- Pipeline.
- Contract.
- QualityRule.
- Run.
- Policy.
- Dashboard.
- Report.
- Notebook.
- Incident.

#### run-service

Purpose:

- Trigger runs.
- Track execution state.
- Store run attempts.
- Store status transitions.
- Handle retries/backfills.
- Connect to adapter APIs.

Run states should include at least:

```text
CREATED
QUEUED
DISPATCHED
RUNNING
SUCCEEDED
FAILED
CANCELLED
TIMED_OUT
RETRYING
SKIPPED
```

#### connection-service

Purpose:

- Source/destination connection definitions.
- Secrets references.
- Test connection.
- Schema discovery orchestration.
- Adapter selection.

Important:

- Store secret references, not raw secrets where possible.
- Support KMaaS/Vault/Kubernetes secret integration later.

#### policy-service

Purpose:

- RBAC/ABAC policies.
- Row/column restrictions.
- AI permissions.
- Approval policies.
- Data movement restrictions.
- Environment-level controls.

#### audit-service

Purpose:

- Immutable audit trail.
- Who did what, when, where, why.
- AI proposal audit.
- Approval audit.
- Execution audit.
- Policy evaluation audit.

#### event-gateway

Purpose:

- Central event ingest/egress.
- Kafka/Redpanda publisher.
- Webhook receiver.
- Standard event envelope enforcement.

#### projection-service

Purpose:

- Read models.
- Materialized views.
- Dashboard-friendly projections.
- Search indexes.
- Denormalized query models for UI.

### A4.3 Event-First Architecture

User explicitly wants event-first architecture using:

- Kafka or Redpanda.
- Schema Registry.
- Outbox pattern.
- Standard event envelope.

Standard event envelope should include:

```json
{
  "event_id": "uuid",
  "event_type": "fusion.resource.created",
  "event_version": "1.0",
  "occurred_at": "2026-05-18T00:00:00Z",
  "producer": "registry-service",
  "tenant_id": "tenant_123",
  "project_id": "project_123",
  "actor": {
    "actor_type": "USER|SERVICE|AI_OPERATOR",
    "actor_id": "user_123"
  },
  "resource": {
    "resource_type": "Connection",
    "resource_id": "conn_123",
    "resource_version": 1
  },
  "correlation_id": "uuid",
  "causation_id": "uuid",
  "payload": {}
}
```

Event topics can start simple:

```text
fusion.identity.events
fusion.tenant.events
fusion.project.events
fusion.registry.events
fusion.connection.events
fusion.run.events
fusion.policy.events
fusion.audit.events
fusion.ai.events
fusion.integration.events
```

Outbox pattern:

- Each service writes transaction + outbox event in same DB transaction.
- Event publisher reads outbox and publishes to Kafka/Redpanda.
- Consumers are idempotent.
- Projection service builds read models.

---

## A5. Reporting Portal Discussion — Must Not Be Lost

The user has a detailed reporting platform vision that overlaps with Fusion's Reporting/BI/Delivery layers.

### A5.1 Reporting Modules

Three main modules:

1. **Yap Portal / External portal**
   - External users such as banks download predefined reports with filters.

2. **Internal report configuration portal**
   - Internal users configure reports.
   - Scheduling.
   - Delivery.
   - Notifications.
   - Report templates.

3. **EMR/Spark Reports**
   - Headless Spark execution for heavy reports.
   - Jobs can run separately from UI.

### A5.2 Report Data Sources

Allowed data sources:

- RDBMS tables.
- JSONs.
- Trino/Iceberg tables.
- PostgreSQL.
- Object storage-backed data.

### A5.3 Output Formats

Required outputs:

- CSV.
- XLS.
- XLSX.
- TXT.
- Custom separator support.
- Quote/newline handling.

### A5.4 Delivery Channels

Discussed delivery channels:

- Portal download.
- Email.
- SFTP.
- S3.
- Azure blob.
- Webhook/Kafka notifications.

### A5.5 Performance Targets

User discussed ambitious performance:

- Around 100 reports/sec generation throughput target.
- Average report size around 100 MB.
- Need HPA/Kubernetes autoscaling.
- Need CPU/memory optimization.
- Need robust error/retry semantics.

### A5.6 Error and Retry Requirements

Need classify and rerun by error group:

- OOMKilled.
- Transient DB issue.
- Query timeout.
- Permission denied.
- File write failure.
- Network/SFTP failure.
- Object store unavailable.
- Bad parameter.
- Data unavailable.

Retry design must support:

- Automatic retry for transient failures.
- Manual rerun by group.
- Tenant-isolated rerun.
- Idempotent file output.
- Status tracking.
- Audit logging.

### A5.7 Go/gRPC Architecture Preferences

User prefers:

- Go backend services for the control plane and reporting portal APIs.
- gRPC service for query execution.
- Kafka communication between services.
- Common authentication.
- Multi-tenant context propagation.

User noted:

- Report service should not directly connect to DB.
- Query execution should happen via separate gRPC service.
- This keeps data execution separate from frontend-facing APIs and improves deployment isolation.

### A5.8 Session/Auth Requirements

Session logic discussed:

- Initial session: 30 minutes.
- If user active, extend by 10–15 minutes.
- Maximum 3 extensions.
- After session ends, JWT must be revoked.
- Keycloak should be used for authentication and identity management.

### A5.9 PCI DSS 4.0 Reporting Portal

The reporting platform must be PCI DSS 4.0-aware:

- Tenant-specific schema/database context.
- Audit trail.
- Access control.
- Secure report delivery.
- Sensitive data masking/encryption.
- No accidental cross-tenant access.

---

## A6. Fusion CDC Engine — Detailed Memory

Fusion CDC Engine is a separate but highly relevant architecture reference.

### A6.1 Core Analogy

- Control Plane = Brain.
- Data Plane = Muscles.
- Redis Streams = Nervous System.
- Spark = Processing Engine.
- PostgreSQL = Metadata Memory.

### A6.2 Core Components

- Frontend: React + TypeScript dashboard.
- Control Plane: FastAPI backend.
- CDC workers: Python.
- Spark consumer: PySpark.
- Transform worker.
- Redis Streams.
- PostgreSQL metadata DB.
- Airflow.
- Prometheus + Grafana.
- Helm/Kubernetes deployment.

### A6.3 Main Flow

1. User configures source and destination in UI.
2. Control plane stores config and orchestrates workers.
3. CDC worker reads database changes.
4. Worker publishes normalized events to Redis Streams.
5. Spark consumer reads events.
6. Transformations and DQ validations run.
7. Schema evolution is handled.
8. Data is written to Postgres warehouse or Iceberg.
9. Metrics are emitted.
10. Checkpoints are updated.

### A6.4 Features Discussed

- Multi-tenant CDC.
- MySQL/PostgreSQL/MongoDB source support.
- Postgres/Iceberg destination support.
- JSON flattening.
- Sensitive data masking.
- Real-time mode.
- Scheduled mode.
- Fallback queues when Redis fails.
- Circuit breakers.
- Retries.
- Prometheus metrics.
- HPA.

### A6.5 Challenges to Address in Fusion

- Redis can become bottleneck at extreme scale.
- Spark adds operational complexity.
- Large tenant count needs careful worker balancing.
- CDC connector per tenant vs per physical source is tricky.
- Need deduplication and checkpoint correctness.
- Need delete handling.
- Need schema evolution governance.

### A6.6 How This Maps to Fusion

Fusion should not blindly copy Fusion as-is. Instead:

- Treat Fusion as one possible execution-plane module.
- Extract control-plane ideas into Fusion kernel.
- Keep CDC adapter pluggable.
- Support Debezium/Kafka and Redis/Python CDC options.
- Let customers BYO CDC engine where possible.

---

## A7. Superset / BI Platform Discussion — Expanded

The user researched how a forked Superset can compete with Power BI/Tableau.

### A7.1 Core Conclusion

Do not replace Superset. Use Superset as the BI shell and surround it with governed semantic, metadata, quality, and AI layers.

Superset remains responsible for:

- SQL Lab.
- Explore.
- Charts.
- Dashboards.
- Embedded analytics.
- RBAC foundation.
- Visualization plugins.

Fusion or sidecar platform layers provide:

- Semantic layer.
- Business glossary.
- Dataset profiling.
- Data quality score.
- Certified datasets/metrics.
- Relationship graph.
- AI assistant.
- Governed chart generation.
- Lifecycle management.

### A7.2 Product Principle

Users should never start from raw tables. They should start from certified business concepts.

Examples:

- Payments.
- Customers.
- Transactions.
- Settlement.
- Revenue.
- Churn.
- Active Customers.
- Failed Transactions.
- Successful Transaction Rate.
- Bank.
- Region.
- Date.
- Status.

### A7.3 Why Superset Alone Is Not Enough

Superset is strong for technical users but weaker than Power BI/Tableau in:

- Semantic modelling.
- Business-readable metadata.
- Relationship modelling.
- Data prep.
- Dataset certification.
- Metric lifecycle.
- Non-technical guided exploration.
- AI-readiness.
- Data quality and freshness visibility.
- Governance workflows.

### A7.4 Recommended Components

- Apache Superset: BI shell.
- Cube Core: headless semantic layer candidate.
- MetricFlow/dbt Semantic Layer: metrics-as-code candidate.
- OpenMetadata: catalog/governance/glossary/lineage.
- DataHub: graph-heavy metadata/context option.
- Great Expectations: formal DQ checks.
- Soda Core: lightweight SQL-based checks.
- Apache ECharts: custom visualizations.
- Vega-Lite: AI-generated chart specs.
- Metabase: UX inspiration only; do not copy code without legal review.
- Lightdash: semantic BI/dbt-first inspiration.
- Evidence.dev: BI-as-code/narrative dashboard inspiration.

### A7.5 Data Onboarding Wizard

When a dataset is added, automatically scan:

- Schema.
- Column types.
- Null percentages.
- Distinct counts.
- Min/max values.
- Top values.
- Date ranges.
- Key candidates.
- PII candidates.
- Status fields.
- Currency fields.
- JSON fields.
- Grain.
- Likely dimensions.
- Likely measures.
- Default time column.
- Suggested metrics.
- Data quality warnings.
- Join candidates.

### A7.6 Semantic Model Requirements

Every field should include:

- Technical name.
- Business name.
- Description.
- Owner.
- Domain.
- Sensitivity.
- Allowed aggregations.
- Synonyms.
- Example questions.

Metrics should include:

- Formula.
- Grain.
- Valid dimensions.
- Time dimension.
- Filters.
- Owner.
- Certification status.
- Tests.
- Version.

### A7.7 Relationship Graph Requirements

Discover/store:

- Primary keys.
- Foreign keys.
- Natural keys.
- Cardinality.
- Join type.
- Tenant keys.
- Date keys.
- Fanout risk.
- Many-to-many ambiguity.

Warn user if selected chart may duplicate measures.

### A7.8 AI Chart Generation Guardrails

AI must:

- Retrieve semantic metadata first.
- Use certified metrics/dimensions where possible.
- Validate permissions.
- Generate semantic query or validated SQL.
- Dry-run query.
- Validate chart spec.
- Explain metric definition.
- Warn about uncertified data.

AI must not:

- Write arbitrary raw SQL against unknown tables.
- Invent metrics.
- Ignore row-level security.
- Join tables without relationship validation.

### A7.9 BI Roadmap From Research

- Phase 0: Discovery, fork audit, requirement mapping.
- Phase 1: Data understanding MVP — profiler, dataset summary, column classification, quality/freshness score.
- Phase 2: Semantic foundation — metric registry, relationship graph, certified dataset/metric workflow.
- Phase 3: Guided Explore.
- Phase 4: AI assistant.
- Phase 5: Enterprise operations — scheduling, lifecycle, governance, monitoring.

---

## A8. Keycloak / SSO Real-World Discussion

The user discussed integrating a Report Manager portal with an existing VisapayAgent portal.

### A8.1 Realms

- Own realm: `report-manager`.
- Client: `report-manager-client`.
- Partner/client realm: `VisapayAgent`.
- Client portal: `visapay-customer-portal`.

### A8.2 Desired Flow

User logs into VisapayAgent portal, clicks a button, lands in Report Manager portal without re-login.

### A8.3 Discussed Solution

OIDC brokering:

- Add Identity Provider in `report-manager` realm.
- Use OpenID Connect v1.0 provider.
- Alias: `visapayagent`.
- Discover endpoints from partner realm.
- Authorization code flow creates SSO session in own realm.

### A8.4 Attribute Mapping Problem

Report Manager frontend expects attributes like:

- bank_id.
- bank_name.
- sub_bank_ids.
- sub_bank_names.

Partner realm may not provide them.

Possible workaround discussed:

- Add hardcoded or mapped attributes at identity-provider level in own realm.
- Keep frontend logic compatible.

### A8.5 Important Warning

Do not disable auth in Report Manager. That is not acceptable for regulated multi-tenant reporting.

---

## A9. Kubernetes / Autoscaling / JVM Concerns

The user repeatedly discussed production Kubernetes concerns.

### A9.1 HPA

Current pattern:

- Scale based on CPU.
- Scale based on memory.
- If either crosses threshold, pods scale.
- Threshold around 80% was discussed.

Need clarify in implementation docs:

- Requests drive HPA percentages.
- Limits protect nodes but can cause throttling/OOM.
- JVM Xmx should be aligned with container memory.
- For Java services, leave headroom for metaspace, direct memory, thread stacks, native memory.

### A9.2 Node Scheduling

User wants:

- Node selectors/tolerations where possible.
- Noise isolation in shared infra.
- But also a strategy when infra constraints do not allow selectors/tolerations.

### A9.3 Multi-Tenant Datasource Concern

Current pattern discussed:

- App startup reads tenant table.
- Builds map: tenant → HikariDataSource.
- Concern: this does not scale if tenant count grows heavily.

Fusion design guidance:

- Avoid one permanent pool per tenant for large tenant counts.
- Prefer connection broker/query service/pooled dynamic connection management.
- Centralize tenant context.
- Consider per-tenant pool only for limited high-value tenants.

---

## A10. JupyterHub / Databricks Replacement Discussion

User has real experience replacing Databricks-like workloads with JupyterHub + Spark.

Important product inspiration:

- Notebook layer should support multi-tenant isolated notebooks.
- Spark in headless mode can be used for cost-effective workloads.
- JupyterHub on Kubernetes can be a notebook execution plane.
- Fusion should not necessarily build notebook execution first, but should support notebook registry/visibility later.

Fusion notebook layer should eventually provide:

- Notebook registration.
- Owner/domain metadata.
- Environment mapping.
- Execution history.
- Dataset lineage from notebooks where possible.
- Cost attribution.
- Policy checks.

---

## A11. Iceberg / Postgres ETL Discrepancy Discussion

User investigated a case where Postgres ETL behaved correctly but Iceberg destination had extra/mismatched records.

Relevant architecture patterns from that discussion:

- PySpark JDBC pipeline.
- MySQL/Postgres source reads.
- Mongo via PySpark or relational staging logic.
- Cursor-field windows for incremental loads.
- Staging tables for relational joins.
- Postgres JDBC writes.
- SQLAlchemy final inserts.
- Delete matching primary keys in target before insert for UPSERT-like behavior.
- max(cursor_field) calculation per target table.
- Timezone checks including Spark session/JVM/default and IST conversions.
- Quality checks.
- Teams webhook operational alerts.

Fusion lesson:

- Incremental logic must be explicit, testable, and auditable.
- Timezone handling must be standardized.
- Iceberg write semantics differ from Postgres upsert semantics.
- Data reconciliation should be first-class.
- Cursor windows, deletes, and idempotency need metadata tracking.

---

## A12. AI Coding / Agent Workflow Discussion

The user wants to build the product heavily using AI tools, with the user acting as final lead.

### A12.1 Desired Human/AI Team Model

User analogy:

- User = final lead / architect / reviewer.
- Sonnet or strongest model = tech lead / reviewer.
- Qwen/Kimi/DeepSeek = developers/testers.
- Roo/Aider/OpenCode = hands that edit code.
- Continue + Codebase-Memory MCP = memory/context brain.
- Playwright MCP = UI tester.
- OpenRouter = model router.
- `.ai` folder = permanent project knowledge.

### A12.2 Target Setup

Suggested pattern:

```text
Ruflo / AI team manager
  -> Sonnet as tech lead/reviewer
  -> Qwen/Kimi/DeepSeek as developer/tester agents
  -> Roo/Aider/OpenCode edits files
  -> Continue + memory MCP provide repo context
  -> Playwright MCP tests UI
  -> OpenRouter routes models
  -> GitHub repo + CI validates
  -> User gives final approval
```

### A12.3 Repo Memory Requirement

Codex or any AI agent should read files in `.ai/`:

```text
.ai/
  PROJECT_CONTEXT.md
  ARCHITECTURE.md
  DECISIONS.md
  ROADMAP.md
  CODING_STANDARDS.md
  OPEN_SOURCE_STRATEGY.md
  AI_OPERATOR_RULES.md
  DOMAIN_MODEL.md
  SERVICE_BOUNDARIES.md
```

### A12.4 AI Agent Rules

Agents must:

- Read `.ai/PROJECT_CONTEXT.md` before changing code.
- Make small PR-sized changes.
- Write tests for core logic.
- Never make architecture changes silently.
- Preserve service boundaries.
- Update `DECISIONS.md` when making irreversible decisions.
- Use Playwright for UI flows.
- Use API tests for backend services.
- Avoid hardcoding vendor-specific choices unless adapter pattern exists.

### A12.5 Tooling Preferences

The user compared:

- Cursor.
- Windsurf.
- GitHub Copilot.
- Claude Code.
- Cline.
- Roo Code.
- Aider.
- OpenRouter.
- Continue.
- Codebase memory MCP.

For this project, preferred characteristics:

- Large repo context.
- Agentic multi-file editing.
- Java + React support.
- Architecture-aware changes.
- Testing/debugging loop.
- Cost-effective multi-model routing.
- Persistent memory.

---

## A13. Startup / Market / Positioning Discussions

The user asked for VC/growth/operator-style analysis of the startup.

### A13.1 Ideal First Customers

Likely first customers:

- Mid-size fintechs.
- Regulated SaaS companies with multi-tenant data platforms.
- Banks or banking infrastructure providers.
- Data teams struggling with Airbyte/dbt/Superset/DataHub fragmentation.
- Companies with many pipelines but weak governance.
- Teams replacing expensive Databricks/managed BI workflows with open-source stack.
- Companies needing audit-ready CDC/reporting.

### A13.2 Pain Points Fusion Solves

- Too many tools with no single control plane.
- Data teams cannot trace failures quickly.
- Business users cannot trust raw tables.
- CDC pipelines break due to schema changes.
- Metrics differ across reports.
- Governance is disconnected from execution.
- BI dashboards exist without quality/freshness context.
- AI over raw SQL is unsafe.
- Multi-tenant reporting is hard.
- Cost attribution is weak.
- Open-source stacks need enterprise glue.

### A13.3 Positioning

Strong positioning:

> DCraft Fusion is the control plane for open and hybrid data platforms — giving data teams Databricks-like operating visibility and governance without forcing them into one proprietary engine.

Alternative:

> The AI-assisted operating layer for Airbyte, dbt, Spark, Superset, Iceberg, Trino, OpenMetadata, and your existing data stack.

### A13.4 What Fusion Is Not

- Not only a dashboard tool.
- Not only a CDC tool.
- Not only a data catalog.
- Not only an AI chatbot.
- Not a replacement for every data tool.
- Not a warehouse.
- Not a Databricks clone that must own compute/data.

---

## A14. Open-Source Reuse Strategy

The user strongly prefers using open-source repos wherever possible, modifying and integrating instead of building everything from scratch.

### A14.1 Principle

Use open-source tools as providers, engines, or inspiration, but keep Fusion core native.

Fusion core must own:

- Resource API.
- Control plane UX.
- Service boundaries.
- Provider SDK.
- Adapter contracts.
- Edge agent model.
- AI governance.
- Policy engine.
- Audit model.
- Tenant model.
- Run registry.
- Projection/read model.

### A14.2 Open-Source Tools by Layer

Ingestion:

- Airbyte.
- Meltano/Singer.

CDC:

- Debezium.
- Kafka Connect.
- Redpanda/Kafka.
- Redis Streams for lightweight prototype.

Orchestration:

- Airflow.
- Dagster.
- Prefect.

Transformation:

- dbt Core.
- Spark.
- Flink.
- SQLMesh.

Quality:

- Great Expectations.
- Soda Core.

Lineage/Metadata:

- OpenLineage.
- Marquez.
- OpenMetadata.
- DataHub.

BI:

- Superset.
- Cube.
- MetricFlow.
- Lightdash reference.
- Evidence.dev reference.

Visualization:

- ECharts.
- Vega-Lite.

Policy/Governance:

- Open Policy Agent.
- Keycloak.

Secrets:

- Vault.
- Cloud KMS.
- Kubernetes Secrets for local/dev only.

Monitoring:

- Prometheus.
- Grafana.
- OpenTelemetry.

---

## A15. Canonical MVP Build Plan — More Concrete

The MVP must prove the control plane, not every layer.

### A15.1 MVP Goal

Show that Fusion can:

1. Register a tenant.
2. Register a project.
3. Register a connection.
4. Discover metadata.
5. Define a workflow/resource.
6. Trigger a run through an adapter.
7. Track run state.
8. Emit events.
9. Store audit logs.
10. Show UI visibility.
11. Let AI Operator propose a safe improvement without auto-executing.

### A15.2 MVP Stack Suggestion

Backend:

- Go 1.26.x.
- Standard Go HTTP services first, gRPC where service-to-service streaming/execution is needed.
- Go workspaces/modules for backend packages.
- PostgreSQL.
- Kafka or Redpanda.
- Flyway/Liquibase.
- OpenAPI.

Frontend:

- React.
- TypeScript.
- Tailwind or a structured component library.
- Auth-aware UI.
- Data platform admin UX.

Infra:

- Docker Compose for local.
- Kubernetes manifests/Helm later.
- Local Redpanda/Kafka.
- Local Postgres.
- Optional Keycloak dev realm.

AI:

- Start with AI proposal service stub.
- No production execution.
- Store proposals and approvals.

### A15.3 MVP Services to Implement First

Recommended order:

1. `registry-service`.
2. `tenant-service`.
3. `connection-service`.
4. `run-service`.
5. `audit-service`.
6. `event-gateway`.
7. `projection-service`.
8. UI shell.
9. One adapter: dbt or Airflow or mock execution adapter.
10. AI proposal stub.

### A15.4 First Demo Flow

```text
Login
 -> choose tenant
 -> create project
 -> create Postgres/Trino connection metadata
 -> test connection through adapter/mock
 -> discover schemas/tables
 -> register dataset
 -> create data quality rule
 -> trigger profile/validation run
 -> observe run state
 -> show audit trail
 -> AI Operator recommends adding freshness check
 -> user approves draft only
 -> decision/event stored
```

---

## A16. Repo Structure — Expanded Codex Target

Codex should organize repo cleanly.

```text
dcraft-fusion/
  .ai/
    PROJECT_CONTEXT.md
    ARCHITECTURE.md
    DECISIONS.md
    DOMAIN_MODEL.md
    SERVICE_BOUNDARIES.md
    ROADMAP.md
    AI_OPERATOR_RULES.md
    OPEN_SOURCE_STRATEGY.md
    CODING_STANDARDS.md

  apps/
    web-console/
      package.json
      src/
        app/
        components/
        features/
          tenants/
          projects/
          registry/
          connections/
          runs/
          audit/
          ai-operator/
        routes/
        api/
        auth/

  services/
    identity-service/
    tenant-service/
    project-service/
    registry-service/
    connection-service/
    run-service/
    policy-service/
    audit-service/
    event-gateway/
    projection-service/
    ai-operator-service/

  packages/
    fusion-domain/
    fusion-events/
    fusion-adapter-sdk/
    fusion-policy-sdk/
    fusion-testkit/

  adapters/
    airflow-adapter/
    dbt-adapter/
    superset-adapter/
    openmetadata-adapter/
    datahub-adapter/
    airbyte-adapter/
    debezium-adapter/
    trino-adapter/
    iceberg-adapter/
    spark-adapter/
    great-expectations-adapter/
    soda-adapter/

  infra/
    docker-compose/
    helm/
    kubernetes/
    keycloak/
    redpanda/
    postgres/
    observability/

  docs/
    architecture/
    api/
    adr/
    product/
    diagrams/

  examples/
    demo-bank/
    demo-fintech/
    demo-retail/
    sample-dbt-project/
    sample-airflow-dags/
```

---

## A17. Domain Model — Expanded Draft

### A17.1 Tenant Model

```text
Organization
  id
  name
  type: CUSTOMER | INTERNAL | PARTNER
  deployment_mode: SAAS | BYOC | ON_PREM

Tenant
  id
  organization_id
  name
  external_ref
  parent_tenant_id
  isolation_level
  status

TenantContext
  tenant_id
  organization_id
  project_id
  environment
  actor_id
  roles
  policies
```

### A17.2 Project Model

```text
Project
  id
  tenant_id
  name
  domain
  owner_user_id
  lifecycle_state
  default_environment
  created_at
```

### A17.3 Resource Model

```text
Resource
  id
  tenant_id
  project_id
  resource_type
  name
  version
  status
  spec_json
  metadata_json
  created_by
  created_at
  updated_at
```

Resource types:

```text
CONNECTION
DATASET
PIPELINE
WORKFLOW
JOB
DASHBOARD
REPORT
METRIC
CONTRACT
QUALITY_RULE
POLICY
INCIDENT
AI_PROPOSAL
```

### A17.4 Connection Model

```text
Connection
  id
  tenant_id
  project_id
  name
  connection_type
  provider
  host_ref
  database_ref
  secret_ref
  metadata
  status
```

Connection types:

- POSTGRES.
- MYSQL.
- MONGODB.
- TRINO.
- ICEBERG.
- S3.
- AZURE_BLOB.
- KAFKA.
- REDIS_STREAMS.
- SUPERSET.
- OPENMETADATA.
- AIRFLOW.
- DBT.

### A17.5 Run Model

```text
Run
  id
  tenant_id
  project_id
  resource_id
  resource_type
  run_type
  status
  trigger_type
  triggered_by
  started_at
  ended_at
  attempt
  engine
  adapter
  external_run_id
  error_code
  error_group
  error_message
```

### A17.6 AI Proposal Model

```text
AIProposal
  id
  tenant_id
  project_id
  proposal_type
  target_resource_id
  target_resource_type
  summary
  rationale
  risk_level
  diff_payload
  required_approvals
  status: DRAFT | PENDING_APPROVAL | APPROVED | REJECTED | APPLIED
  created_by_model
  created_at
```

---

## A18. Adapter SDK Contract

Every adapter should implement a standard shape.

```text
Adapter
  getCapabilities()
  validateConnection(connectionSpec)
  discoverMetadata(connectionSpec)
  planAction(actionSpec)
  executeAction(actionSpec)
  getRunStatus(externalRunId)
  cancelRun(externalRunId)
  fetchLogs(externalRunId)
  normalizeError(error)
```

Adapter result should include:

```json
{
  "status": "SUCCESS|FAILED|RUNNING|CANCELLED",
  "external_id": "engine-specific-id",
  "logs_ref": "...",
  "metrics": {},
  "error": {
    "code": "QUERY_TIMEOUT",
    "group": "TRANSIENT_ENGINE_FAILURE",
    "message": "...",
    "retryable": true
  }
}
```

Adapters must not leak engine-specific complexity into UI.

---

## A19. Security and Compliance Requirements

### A19.1 Multi-Tenant Isolation

Every table should include tenant/project context where applicable.

Every request should carry tenant context.

Every event should include tenant_id.

Every adapter call should be tenant-scoped.

No cross-tenant queries without explicit consolidated permission.

### A19.2 Secrets

Do not store plain secrets in service DB.

Store references to:

- Vault.
- Cloud secret manager.
- Kubernetes secret in dev.
- KMaaS.

### A19.3 Audit

Audit everything important:

- Login.
- Connection creation.
- Secret reference change.
- Policy change.
- Workflow change.
- Run trigger.
- Run cancel.
- AI proposal creation.
- AI proposal approval/rejection.
- Export/download.
- Admin access.

### A19.4 Compliance Inspiration

User context includes regulated banking and PCI DSS/SOX-style environments. Fusion should be designed audit-first.

---

## A20. Detailed Product Backlog

### A20.1 Kernel Backlog

- Create mono-repo structure.
- Add `.ai` memory files.
- Add common domain package.
- Add event envelope package.
- Add tenant-service skeleton.
- Add registry-service skeleton.
- Add connection-service skeleton.
- Add run-service skeleton.
- Add audit-service skeleton.
- Add event-gateway skeleton.
- Add projection-service skeleton.
- Add Postgres migrations.
- Add Redpanda/Kafka local compose.
- Add OpenAPI specs.
- Add health endpoints.
- Add basic auth placeholder.
- Add tenant context middleware.
- Add outbox table.
- Add outbox publisher.
- Add idempotent consumer base.

### A20.2 UI Backlog

- Login shell.
- Tenant switcher.
- Project switcher.
- Resource registry page.
- Connection list page.
- Connection creation wizard.
- Dataset discovery page.
- Runs page.
- Run detail page.
- Audit timeline page.
- AI Operator proposal page.
- Settings/policies page.

### A20.3 CDC Backlog

- Source connector registry.
- Debezium adapter draft.
- Redis Streams adapter draft.
- CDC event schema.
- Checkpoint model.
- Dedup model.
- Delete event handling.
- Schema evolution alerts.
- Tenant pause/resume.
- Backfill support.

### A20.4 BI Backlog

- Superset adapter.
- OpenMetadata adapter.
- Dataset profiler.
- Column classifier.
- Quality/freshness score.
- Metric registry.
- Relationship graph.
- Certified dataset workflow.
- AI chart draft.
- Chart spec validator.
- Dashboard lifecycle.

### A20.5 Reporting Backlog

- Report definition model.
- Report parameter model.
- Report schedule model.
- Export format model.
- Delivery target model.
- Report run model.
- Streaming CSV writer.
- XLSX writer.
- TXT writer with custom separator.
- Delivery service.
- Retry/error grouping.
- WebSocket progress.
- MinIO/S3 path strategy.

---

## A21. Concrete Coding Standards for Codex

### A21.1 General

- Keep code modular.
- Avoid premature mega-abstractions.
- But preserve adapter boundaries.
- Write clear package names.
- Do not hardcode tenant IDs.
- Do not hardcode secrets.
- Add tests for domain logic.
- Add migrations for DB changes.
- Update OpenAPI/contracts when adding endpoints.
- Add ADR for major decisions.

### A21.2 Go Backend

- Go 1.26.x.
- Prefer small domain-focused packages.
- Use explicit DTO structs for API boundaries.
- Separate domain, application, infrastructure, API layers.
- Keep validation in service/domain boundaries.
- Use Flyway/Liquibase.
- Use Testcontainers for integration tests where possible.
- Use structured logging.
- Include correlation_id and tenant_id in logs.

### A21.3 React Frontend

- TypeScript.
- Feature-based folders.
- API client layer.
- Strong typing for backend DTOs.
- Tenant/project context provider.
- Do not embed backend URLs all over components.
- Use table/detail/wizard UX for platform flows.

### A21.4 Eventing

- Standard event envelope.
- Use schema versioning.
- Include tenant/project context.
- Ensure idempotency.
- Use outbox.
- Avoid direct cross-service DB reads.

---

## A22. What Codex Should Build Next — Exact Instruction

When Codex receives this file in the repo, the first implementation request should be something like:

```text
Read .ai/PROJECT_CONTEXT.md completely.
Do not summarize only. Treat it as the source of truth.
Create the initial DCraft Fusion Control Plane Kernel skeleton.

Use Go 1.26.x for backend services and React + TypeScript for web-console.
Start with:
- tenant-service
- registry-service
- connection-service
- run-service
- audit-service
- event-gateway
- projection-service
- shared fusion-domain package
- shared fusion-events package
- docker-compose with Postgres + Redpanda
- basic migrations
- standard event envelope
- outbox table design
- health endpoints
- README explaining how to run locally

Do not implement full CDC/Superset/AI yet. Create clean extension points and TODOs.
```

---

## A23. Missing / Uncertain Items to Fill Later

The following items were discussed conceptually but need final decisions before production implementation:

- Exact frontend framework/design system.
- Whether backend modules stay as separate Go modules or move to a Go workspace as service count grows.
- Whether MVP uses Kafka or Redpanda locally.
- Whether first adapter should be dbt, Airflow, Superset, or mock adapter.
- Whether identity-service is custom or Keycloak-first wrapper.
- Whether control-plane APIs are REST-first, GraphQL-first, or mixed REST/gRPC.
- Whether OpenMetadata or DataHub is first metadata integration.
- Whether Cube or MetricFlow is first semantic layer integration.
- Whether CDC MVP uses Debezium/Kafka or Redis/Python prototype.
- Whether deployment target is SaaS-first, BYOC-first, or local-demo-first.

Codex should not block on these. It should implement safe extension points and document assumptions.

---

## A24. One-Page Mental Model for Every AI Agent

Fusion is not trying to own everything. Fusion is trying to control, connect, govern, observe, and improve everything.

```text
User / Data Team / Business Team
        |
        v
DCraft Fusion UX + Control Plane
        |
        +--> Metadata / Registry / Policy / Audit / AI Proposals
        |
        +--> Adapters / Edge Agents
                |
                +--> Airbyte / Debezium / Kafka / Redis
                +--> Airflow / Dagster / Prefect
                +--> dbt / Spark / Flink / SQLMesh
                +--> Iceberg / Trino / Warehouses
                +--> OpenMetadata / DataHub / OpenLineage
                +--> Superset / Cube / MetricFlow
                +--> Great Expectations / Soda
                +--> S3 / SFTP / Email / Webhooks
```

The customer keeps their data and engines. Fusion gives them one brain.

---

## A25. Final Reminder After Expansion

The user was correct that a high-level Markdown file is not enough. For Codex, the file must include:

- Product identity.
- Historical context.
- Frozen decisions.
- Non-negotiable boundaries.
- 16-layer architecture.
- Kernel build plan.
- Service names.
- Event architecture.
- Reporting portal requirements.
- CDC engine requirements.
- Superset BI roadmap.
- Keycloak SSO lessons.
- Kubernetes/HPA lessons.
- AI-agent workflow.
- Open-source reuse strategy.
- Domain model.
- Adapter SDK.
- Concrete backlog.
- Exact next prompt for Codex.

This expanded file now contains those missing project discussions in one place.
