# Coding Rules

- Prefer clear service boundaries over early distributed complexity.
- Every write path must carry tenant ID, actor ID, and correlation ID.
- Every significant state change must emit an audit event.
- Validate all inputs at API and service boundaries.
- Use DTOs/contracts for APIs.
- Use database migrations for schema changes.
- Keep secrets as references; never store, log, or return raw secret values.
- Enforce authorization on the backend.
- Add tenant-isolation tests for shared storage logic.
- Add contract tests for adapters/providers.
- Make AI-generated changes visible and approval-gated.

