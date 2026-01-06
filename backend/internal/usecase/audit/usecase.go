package audit

import (
	"context"
	"time"

	"github.com/google/uuid"

	domainaudit "backend/internal/domain/audit"
)

// Service records audit logs.
type Service struct {
	repo domainaudit.Repository
	now  func() time.Time
}

func NewService(repo domainaudit.Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) Record(ctx context.Context, actorID, action, resource string, metadata map[string]any) error {
	log := domainaudit.Log{
		ID:        uuid.NewString(),
		ActorID:   actorID,
		Action:    action,
		Resource:  resource,
		Metadata:  metadata,
		CreatedAt: s.now(),
	}
	return s.repo.Create(ctx, log)
}

// List returns audit logs optionally filtered by actorID.
func (s *Service) List(ctx context.Context, actorID string) ([]domainaudit.Log, error) {
	return s.repo.List(ctx, actorID)
}
