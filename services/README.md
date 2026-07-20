# Services

Control-plane domain services live here.

Initial kernel services:

- `identity-service`
- `tenant-service`
- `project-service`
- `connection-registry-service`
- `workflow-registry-service`
- `run-service`
- `policy-service`
- `audit-service`
- `event-gateway-service`
- `adapter-registry-service`
- `ai-assist-service`
- `metadata-sync-service`
- `reporting-service`
- `cdc-control-service`
- `control-plane-kernel`: first Go Phase 1 kernel implementation for organization/deployment context and audit foundation.

Avoid premature distributed complexity while the implementation is early. These folders mark ownership boundaries; services can still be implemented together until there is a real operational need to split them.
