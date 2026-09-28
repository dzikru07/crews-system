package role

import (
	"context"
	"errors"

	"github.com/dzikru07/crews-system/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	tenantID uuid.UUID,
	req CreateRequest,
) (*Role, error) {

	id := uuid.New()

	var result Role

	err := r.db.Pool.QueryRow(
		ctx,
		`
		INSERT INTO roles (
			id,
			tenant_id,
			name,
			description,
			status
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			'active'
		)
		RETURNING
			id,
			tenant_id,
			name,
			description,
			status,
			created_at,
			updated_at
		`,
		id,
		tenantID,
		req.Name,
		req.Description,
	).Scan(
		&result.ID,
		&result.TenantID,
		&result.Name,
		&result.Description,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	tenantID uuid.UUID,
	roleID uuid.UUID,
) (*Role, error) {

	var result Role

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			tenant_id,
			name,
			description,
			status,
			created_at,
			updated_at
		FROM roles
		WHERE id = $1
		AND tenant_id = $2
		`,
		roleID,
		tenantID,
	).Scan(
		&result.ID,
		&result.TenantID,
		&result.Name,
		&result.Description,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRoleNotFound
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
