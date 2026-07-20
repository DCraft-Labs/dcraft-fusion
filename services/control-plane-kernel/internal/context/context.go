package context

import (
	"errors"
	"net/http"
	"strings"
)

const (
	ActorIDHeader        = "X-Fusion-Actor-Id"
	CorrelationIDHeader  = "X-Fusion-Correlation-Id"
	OrganizationIDHeader = "X-Fusion-Organization-Id"
	TenantIDHeader       = "X-Fusion-Tenant-Id"
	ProjectIDHeader      = "X-Fusion-Project-Id"
)

type RequestContext struct {
	ActorID        string
	CorrelationID  string
	OrganizationID string
	TenantID       string
	ProjectID      string
}

func FromHeaders(headers http.Header) (RequestContext, error) {
	actorID := strings.TrimSpace(headers.Get(ActorIDHeader))
	if actorID == "" {
		return RequestContext{}, errors.New("missing request context header: " + ActorIDHeader)
	}

	correlationID := strings.TrimSpace(headers.Get(CorrelationIDHeader))
	if correlationID == "" {
		return RequestContext{}, errors.New("missing request context header: " + CorrelationIDHeader)
	}

	return RequestContext{
		ActorID:        actorID,
		CorrelationID:  correlationID,
		OrganizationID: strings.TrimSpace(headers.Get(OrganizationIDHeader)),
		TenantID:       strings.TrimSpace(headers.Get(TenantIDHeader)),
		ProjectID:      strings.TrimSpace(headers.Get(ProjectIDHeader)),
	}, nil
}

func (ctx RequestContext) RequireOrganization() error {
	if ctx.OrganizationID == "" {
		return errors.New("missing request context header: " + OrganizationIDHeader)
	}
	return nil
}

func (ctx RequestContext) RequireTenant() error {
	if err := ctx.RequireOrganization(); err != nil {
		return err
	}
	if ctx.TenantID == "" {
		return errors.New("missing request context header: " + TenantIDHeader)
	}
	return nil
}
