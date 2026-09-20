package tenant

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzikru07/crews-system/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound = errors.New("tenant not found")
	ErrConflict = errors.New("tenant already exists")
)

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	tenant Tenant,
) error {
	const query = `
		INSERT INTO tenants (
			id,
			name,
			slug,
			type,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Pool.Exec(
		ctx,
		query,
		tenant.ID,
		tenant.Name,
		tenant.Slug,
		tenant.Type,
		tenant.Status,
		tenant.CreatedAt,
		tenant.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create tenant: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Tenant, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			type,
			status,
			created_at,
			updated_at
		FROM tenants
		WHERE id = $1
	`

	var tenant Tenant

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Slug,
		&tenant.Type,
		&tenant.Status,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get tenant: %w", err)
	}

	return &tenant, nil
}
