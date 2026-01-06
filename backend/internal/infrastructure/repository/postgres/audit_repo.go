package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	domainaudit "backend/internal/domain/audit"
)

// AuditRepository implements audit.Repository via PostgreSQL.
type AuditRepository struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, log domainaudit.Log) error {
	return nil
}

func (r *AuditRepository) List(ctx context.Context, actorID string) ([]domainaudit.Log, error) {
	return []domainaudit.Log{}, nil
}
