package permission

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
) (*Permission, error) {

	id := uuid.New()

	var result Permission

	err := r.db.Pool.QueryRow(
		ctx,
		`
		INSERT INTO permissions (
			id,
			tenant_id,
			name,
			description,
			resource,
			action,
			status
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			'active'
		)
		RETURNING
			id,
			tenant_id,
			name,
			description,
			resource,
			action,
			status,
			created_at,
			updated_at
		`,
		id,
		tenantID,
		req.Name,
		req.Description,
		req.Resource,
		req.Action,
	).Scan(
		&result.ID,
		&result.TenantID,
		&result.Name,
		&result.Description,
		&result.Resource,
		&result.Action,
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
	permissionID uuid.UUID,
) (*Permission, error) {

	var result Permission

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			tenant_id,
			name,
			description,
			resource,
			action,
			status,
			created_at,
			updated_at
		FROM permissions
		WHERE id = $1
		AND tenant_id = $2
		`,
		permissionID,
		tenantID,
	).Scan(
		&result.ID,
		&result.TenantID,
		&result.Name,
		&result.Description,
		&result.Resource,
		&result.Action,
		&result.Status,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPermissionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
