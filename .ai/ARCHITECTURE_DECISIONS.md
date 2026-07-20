# Architecture Decisions

## ADR-0001: Fusion Is A Control Plane

Status: Accepted

Fusion stores metadata, intent, policy, workflow, run, audit, approval, and operational state. It coordinates external execution engines instead of replacing them.

## ADR-0002: Kernel First

Status: Accepted

The first implementation should focus on tenant/project, connection, workflow, run, policy, audit, adapter registry, and AI assist service boundaries. CDC, BI, governance, and advanced automation should build on that kernel.

## ADR-0003: AI Requires Human Approval For Actions

Status: Accepted

AI may explain, recommend, draft, and prepare changes. It must not directly execute production-impacting actions without explicit approval.

## ADR-0004: Provider Contracts Are First-Class

Status: Accepted

Adapters/providers must expose typed capabilities and support health, config validation, translation, execution, status, lineage extraction, and optional event watching.

