# ADR-0005: Backend Language

Status: Accepted

## Context

The architecture documents contain some later exploratory backend-language notes, but the V0 technical architecture makes the primary backend decision clearly:

- Backend: Go.
- API gateway: Go.
- Microservices: Go.
- CLI: Go.
- Rust only if profiling proves specific components need materially better performance or stronger low-level guarantees.

That direction also fits Fusion's control-plane identity because Kubernetes, Docker, Terraform, and much of the CNCF ecosystem are Go-native.

## Decision

Core DCraft Fusion backend services will use **Go**.

Runtime strategy:

- Go for control-plane services, API gateway, reconciliation, provider/adapter registry, audit, workflow/run services, and CLI.
- TypeScript/React for UI and frontend-facing contract ergonomics.
- Python remains reserved for CDC workers, metadata/AI workers, and the existing Fusion CDC Engine where it already exists.
- Rust is allowed only for narrowly scoped high-performance or security-sensitive components after profiling or clear technical need.

Current backend toolchain target:

- Go 1.26.x stable line.

## Consequences

- New backend implementation must be Go-first.
- The earlier non-Go backend scaffold was removed.
- Backend tests use `go test`.
- Domain/service code must enforce tenant context, authorization, audit, and secret redaction from the start.
