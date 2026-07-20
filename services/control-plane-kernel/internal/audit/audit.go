package audit

import (
	"regexp"
	"sync"
	"time"
)

type Action string

const (
	OrganizationCreated   Action = "organization.created"
	TenantCreated         Action = "tenant.created"
	ProjectCreated        Action = "project.created"
	InvitationCreated     Action = "invitation.created"
	ServiceAccountCreated Action = "service_account.created"
	ConnectionCreated     Action = "connection.created"
	ConnectionTested      Action = "connection.tested"
	MetadataDiscovered    Action = "metadata.discovered"
)

type Event struct {
	ID            string            `json:"id"`
	Action        Action            `json:"action"`
	ActorID       string            `json:"actorId"`
	CorrelationID string            `json:"correlationId"`
	ResourceType  string            `json:"resourceType"`
	ResourceID    string            `json:"resourceId"`
	OccurredAt    time.Time         `json:"occurredAt"`
	Metadata      map[string]string `json:"metadata"`
}

type Log interface {
	Append(event Event)
	Events() []Event
}

type MemoryLog struct {
	mu     sync.RWMutex
	events []Event
}

func NewMemoryLog() *MemoryLog {
	return &MemoryLog{}
}

func (log *MemoryLog) Append(event Event) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.events = append(log.events, event)
}

func (log *MemoryLog) Events() []Event {
	log.mu.RLock()
	defer log.mu.RUnlock()
	return append([]Event(nil), log.events...)
}

var secretLikeKey = regexp.MustCompile(`(?i)(secret|password|token|credential|private_key|api_key)`)

func RedactMetadata(metadata map[string]string) map[string]string {
	redacted := map[string]string{}
	for key, value := range metadata {
		if secretLikeKey.MatchString(key) {
			continue
		}
		redacted[key] = value
	}
	return redacted
}
