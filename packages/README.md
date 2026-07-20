# Packages

Shared implementation packages live here.

- `domain`: Core domain models and invariants.
- `authz`: Authorization checks and policy helpers.
- `event-envelope`: Standard event envelope definitions.
- `connector-sdk`: Connector-facing SDK code.
- `adapter-sdk`: Adapter/provider SDK code.
- `policy-sdk`: Policy evaluation helpers and contracts.
- `ui-kit`: Shared UI components.
- `telemetry`: Logging, metrics, tracing, and correlation helpers.

Keep shared packages small and stable. Do not use this folder as a dumping ground for service-specific logic.

