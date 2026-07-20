package phase1

import (
	"testing"
	"time"

	"github.com/dcraft-fusion/control-plane-kernel/internal/audit"
	fusioncontext "github.com/dcraft-fusion/control-plane-kernel/internal/context"
)

func TestCreateConnectionRejectsRawSecretAndRedactsReference(t *testing.T) {
	auditLog := audit.NewMemoryLog()
	store := newDeterministicStore(auditLog, []string{"conn-000001", "audit-000001"})
	ctx := scopedContext("tenant-a", "project-a")

	_, err := store.CreateConnection(ctx, CreateConnectionRequest{Name: "Raw", Kind: ConnectorPostgres, SecretRef: "postgres://user:password@host/db", ReadOnly: true})
	if err == nil || err.Error() != "connection secretRef must be a reference, not a raw secret" {
		t.Fatalf("expected raw secret rejection, got %v", err)
	}

	connection, err := store.CreateConnection(ctx, CreateConnectionRequest{Name: "Warehouse", Kind: ConnectorPostgres, SecretRef: "vault/project-a/postgres", ReadOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if connection.SecretRef != "secretref:vault/project-a/postgres" {
		t.Fatalf("secret reference was not redacted: %+v", connection)
	}
	if len(auditLog.Events()) != 1 || auditLog.Events()[0].Action != audit.ConnectionCreated {
		t.Fatalf("connection creation audit missing: %+v", auditLog.Events())
	}
}

func TestMetadataDiscoveryIsTenantIsolated(t *testing.T) {
	store := newDeterministicStore(audit.NewMemoryLog(), []string{"conn-000001", "audit-000001", "job-000001", "ds-000001", "ds-000002", "audit-000002"})
	ctxA := scopedContext("tenant-a", "project-a")
	ctxB := scopedContext("tenant-b", "project-b")

	connection, err := store.CreateConnection(ctxA, CreateConnectionRequest{Name: "A", Kind: ConnectorPostgres, SecretRef: "vault/a/postgres", ReadOnly: true})
	if err != nil {
		t.Fatalf("unexpected connection error: %v", err)
	}
	if _, err := store.DiscoverMetadata(ctxB, connection.ID); err == nil || err.Error() != "connection not found" {
		t.Fatalf("expected tenant isolation error, got %v", err)
	}
	datasets, err := store.DiscoverMetadata(ctxA, connection.ID)
	if err != nil {
		t.Fatalf("unexpected discovery error: %v", err)
	}
	if len(datasets) != 2 {
		t.Fatalf("expected two discovered datasets, got %d", len(datasets))
	}
	visibleToB, err := store.Datasets(ctxB)
	if err != nil {
		t.Fatalf("unexpected list error: %v", err)
	}
	if len(visibleToB) != 0 {
		t.Fatalf("tenant B saw tenant A datasets: %+v", visibleToB)
	}
}

func TestAIRecommendationsRemainDraft(t *testing.T) {
	store := NewSeedStore(audit.NewMemoryLog())
	recommendations, err := store.AIRecommendations(fusioncontext.RequestContext{
		ActorID:        "user-founder",
		CorrelationID:  "corr-1",
		OrganizationID: "org-b2b2c-1",
		TenantID:       "tenant-brand-a",
		ProjectID:      "project-prod",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recommendations) == 0 {
		t.Fatal("expected seeded AI recommendation")
	}
	for _, recommendation := range recommendations {
		if recommendation.State != ApprovalDraft {
			t.Fatalf("AI recommendation must stay draft in Phase 1: %+v", recommendation)
		}
	}
}

func newDeterministicStore(auditLog audit.Log, ids []string) *Store {
	return NewTestStore(auditLog, func() time.Time {
		return time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)
	}, func() string {
		id := ids[0]
		ids = ids[1:]
		return id
	})
}

func scopedContext(tenantID string, projectID string) fusioncontext.RequestContext {
	return fusioncontext.RequestContext{
		ActorID:        "user-1",
		CorrelationID:  "corr-1",
		OrganizationID: "org-1",
		TenantID:       tenantID,
		ProjectID:      projectID,
	}
}
