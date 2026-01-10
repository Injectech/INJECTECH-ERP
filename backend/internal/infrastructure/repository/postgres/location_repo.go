package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainlocation "backend/internal/domain/location"
)

// LocationRepository implements location.Repository via PostgreSQL.
type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Create(ctx context.Context, l domainlocation.Location) (domainlocation.Location, error) {
	id, err := uuid.Parse(l.ID)
	if err != nil {
		return domainlocation.Location{}, err
	}
	if l.IsDefault {
		_, err := r.db.Exec(ctx, `UPDATE locations SET is_default = false WHERE is_default = true AND deleted_at IS NULL`)
		if err != nil {
			return domainlocation.Location{}, err
		}
	}

	row := r.db.QueryRow(ctx, `
		INSERT INTO locations (id, name, description, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, is_default, created_at, updated_at, deleted_at
	`, id, l.Name, l.Description, l.IsDefault, l.CreatedAt, l.UpdatedAt)

	var created domainlocation.Location
	var deletedAt *time.Time
	if err := row.Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.IsDefault,
		&created.CreatedAt,
		&created.UpdatedAt,
		&deletedAt,
	); err != nil {
		return domainlocation.Location{}, err
	}
	created.ID = id.String()
	created.DeletedAt = deletedAt
	return created, nil
}

func (r *LocationRepository) List(ctx context.Context) ([]domainlocation.Location, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, description, is_default, created_at, updated_at, deleted_at
		FROM locations
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []domainlocation.Location
	for rows.Next() {
		var loc domainlocation.Location
		var id uuid.UUID
		if err := rows.Scan(
			&id,
			&loc.Name,
			&loc.Description,
			&loc.IsDefault,
			&loc.CreatedAt,
			&loc.UpdatedAt,
			&loc.DeletedAt,
		); err != nil {
			return nil, err
		}
		loc.ID = id.String()
		locations = append(locations, loc)
	}
	return locations, rows.Err()
}

func (r *LocationRepository) GetDefault(ctx context.Context) (domainlocation.Location, error) {
	var loc domainlocation.Location
	var id uuid.UUID
	row := r.db.QueryRow(ctx, `
		SELECT id, name, description, is_default, created_at, updated_at, deleted_at
		FROM locations
		WHERE is_default = true AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT 1
	`)
	if err := row.Scan(
		&id,
		&loc.Name,
		&loc.Description,
		&loc.IsDefault,
		&loc.CreatedAt,
		&loc.UpdatedAt,
		&loc.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainlocation.Location{}, domainlocation.ErrNotFound
		}
		return domainlocation.Location{}, err
	}
	loc.ID = id.String()
	return loc, nil
}
