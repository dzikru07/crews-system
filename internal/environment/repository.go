package environment

import (
	"context"

	"github.com/dzikru07/crews-system/internal/database"
	"github.com/google/uuid"
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
	env Environment,
) error {
	_, err := r.db.Pool.Exec(
		ctx,
		`
		INSERT INTO environments (
			id,
			application_id,
			name,
			type,
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
			$7
		)
		`,
		env.ID,
		env.ApplicationID,
		env.Name,
		env.Type,
		env.Status,
		env.CreatedAt,
		env.UpdatedAt,
	)

	return err
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Environment, error) {
	var env Environment

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			application_id,
			name,
			type,
			status,
			created_at,
			updated_at
		FROM environments
		WHERE id = $1
		`,
		id,
	).Scan(
		&env.ID,
		&env.ApplicationID,
		&env.Name,
		&env.Type,
		&env.Status,
		&env.CreatedAt,
		&env.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (r *Repository) GetByIDAndApplication(
	ctx context.Context,
	id uuid.UUID,
	applicationID uuid.UUID,
) (*Environment, error) {
	var env Environment

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			application_id,
			name,
			type,
			status,
			created_at,
			updated_at
		FROM environments
		WHERE
			id = $1
			AND application_id = $2
		`,
		id,
		applicationID,
	).Scan(
		&env.ID,
		&env.ApplicationID,
		&env.Name,
		&env.Type,
		&env.Status,
		&env.CreatedAt,
		&env.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &env, nil
}
