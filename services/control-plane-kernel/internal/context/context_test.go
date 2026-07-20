package context

import (
	"net/http"
	"testing"
)

func TestFromHeadersRequiresActorAndCorrelation(t *testing.T) {
	headers := http.Header{}
	headers.Set(CorrelationIDHeader, "corr-1")

	_, err := FromHeaders(headers)
	if err == nil || err.Error() != "missing request context header: X-Fusion-Actor-Id" {
		t.Fatalf("expected missing actor error, got %v", err)
	}
}

func TestFromHeadersResolvesContext(t *testing.T) {
	headers := http.Header{}
	headers.Set(ActorIDHeader, "user-1")
	headers.Set(CorrelationIDHeader, "corr-1")
	headers.Set(OrganizationIDHeader, "org-1")
	headers.Set(TenantIDHeader, "tenant-1")

	ctx, err := FromHeaders(headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.ActorID != "user-1" || ctx.CorrelationID != "corr-1" || ctx.OrganizationID != "org-1" || ctx.TenantID != "tenant-1" {
		t.Fatalf("unexpected context: %+v", ctx)
	}
}

func TestRequestContextRequiresScopedHeaders(t *testing.T) {
	ctx := RequestContext{ActorID: "user-1", CorrelationID: "corr-1"}
	if err := ctx.RequireOrganization(); err == nil || err.Error() != "missing request context header: X-Fusion-Organization-Id" {
		t.Fatalf("expected missing organization error, got %v", err)
	}

	ctx.OrganizationID = "org-1"
	if err := ctx.RequireTenant(); err == nil || err.Error() != "missing request context header: X-Fusion-Tenant-Id" {
		t.Fatalf("expected missing tenant error, got %v", err)
	}

	ctx.TenantID = "tenant-1"
	if err := ctx.RequireTenant(); err != nil {
		t.Fatalf("unexpected tenant scope error: %v", err)
	}
}
