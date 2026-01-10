package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	domainaudit "backend/internal/domain/audit"
	"backend/internal/infrastructure/repository/postgres/sqlc"
)

// AuditRepository implements audit.Repository via PostgreSQL.
type AuditRepository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{db: db, q: sqlc.New(db)}
}

func (r *AuditRepository) Create(ctx context.Context, log domainaudit.Log) error {
	var actor *uuid.UUID
	if log.ActorID != "" {
		uid, err := uuid.Parse(log.ActorID)
		if err != nil {
			return err
		}
		actor = &uid
	}
	id, err := uuid.Parse(log.ID)
	if err != nil {
		return err
	}
	return r.q.CreateAuditLog(ctx, sqlc.CreateAuditLogParams{
		ID:        id,
		ActorID:   actor,
		Action:    log.Action,
		Resource:  log.Resource,
		Metadata:  log.Metadata,
		CreatedAt: log.CreatedAt,
	})
}

func (r *AuditRepository) List(ctx context.Context, actorID string) ([]domainaudit.Log, error) {
	var actor *uuid.UUID
	if actorID != "" {
		uid, err := uuid.Parse(actorID)
		if err != nil {
			return nil, err
		}
		actor = &uid
	}
	rows, err := r.q.ListAuditLogs(ctx, actor)
	if err != nil {
		return nil, err
	}
	var logs []domainaudit.Log
	for _, row := range rows {
		entry := domainaudit.Log{
			ID:        row.ID.String(),
			Action:    row.Action,
			Resource:  row.Resource,
			CreatedAt: row.CreatedAt,
		}
		if row.ActorID.Valid {
			entry.ActorID = row.ActorID.UUID.String()
		}
		if row.Metadata != nil {
			entry.Metadata = row.Metadata
		}
		logs = append(logs, entry)
	}
	return logs, nil
}
