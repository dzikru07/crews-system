package client

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
	client Client,
) error {
	_, err := r.db.Pool.Exec(
		ctx,
		`
		INSERT INTO clients (
			id,
			environment_id,
			name,
			slug,
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
			$7,
			$8
		)
		`,
		client.ID,
		client.EnvironmentID,
		client.Name,
		client.Slug,
		client.Type,
		client.Status,
		client.CreatedAt,
		client.UpdatedAt,
	)

	return err
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Client, error) {
	var client Client

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			environment_id,
			name,
			slug,
			type,
			status,
			created_at,
			updated_at
		FROM clients
		WHERE id = $1
		`,
		id,
	).Scan(
		&client.ID,
		&client.EnvironmentID,
		&client.Name,
		&client.Slug,
		&client.Type,
		&client.Status,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *Repository) GetByIDAndEnvironment(
	ctx context.Context,
	id uuid.UUID,
	environmentID uuid.UUID,
) (*Client, error) {
	var client Client

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			environment_id,
			name,
			slug,
			type,
			status,
			created_at,
			updated_at
		FROM clients
		WHERE
			id = $1
			AND environment_id = $2
		`,
		id,
		environmentID,
	).Scan(
		&client.ID,
		&client.EnvironmentID,
		&client.Name,
		&client.Slug,
		&client.Type,
		&client.Status,
		&client.CreatedAt,
		&client.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &client, nil
}
