package application

import (
	"context"
	"errors"

	"github.com/dzikru07/crews-system/internal/database"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("application not found")

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
	app Application,
) error {
	_, err := r.db.Pool.Exec(
		ctx,
		`
		INSERT INTO applications (
			id,
			tenant_id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		)
		`,
		app.ID,
		app.TenantID,
		app.Name,
		app.Slug,
		app.Description,
		app.Status,
		app.CreatedAt,
		app.UpdatedAt,
	)

	return err
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Application, error) {
	var app Application

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			tenant_id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at
		FROM applications
		WHERE id = $1
		`,
		id,
	).Scan(
		&app.ID,
		&app.TenantID,
		&app.Name,
		&app.Slug,
		&app.Description,
		&app.Status,
		&app.CreatedAt,
		&app.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (r *Repository) GetByIDAndTenant(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*Application, error) {
	var app Application

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			tenant_id,
			name,
			slug,
			description,
			status,
			created_at,
			updated_at
		FROM applications
		WHERE
			id = $1
			AND tenant_id = $2
		`,
		id,
		tenantID,
	).Scan(
		&app.ID,
		&app.TenantID,
		&app.Name,
		&app.Slug,
		&app.Description,
		&app.Status,
		&app.CreatedAt,
		&app.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &app, nil
}
