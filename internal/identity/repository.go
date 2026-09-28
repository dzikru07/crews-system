package identity

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dzikru07/crews-system/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound = errors.New("identity not found")
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
) (*Identity, error) {
	id := uuid.New()

	var metadata []byte

	if req.Metadata != nil {
		data, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, err
		}

		metadata = data
	}

	var identity Identity

	err := r.db.Pool.QueryRow(
		ctx,
		`
		INSERT INTO identities (
			id,
			tenant_id,
			type,
			name,
			status,
			metadata
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			'active',
			$5
		)
		RETURNING
			id,
			tenant_id,
			type,
			name,
			status,
			metadata,
			last_authenticated_at,
			created_at,
			updated_at
		`,
		id,
		tenantID,
		req.Type,
		req.Name,
		metadata,
	).Scan(
		&identity.ID,
		&identity.TenantID,
		&identity.Type,
		&identity.Name,
		&identity.Status,
		&identity.Metadata,
		&identity.LastAuthenticatedAt,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &identity, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*Identity, error) {
	var identity Identity

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			tenant_id,
			type,
			name,
			status,
			metadata,
			last_authenticated_at,
			created_at,
			updated_at
		FROM identities
		WHERE
			id = $1
			AND tenant_id = $2
		`,
		id,
		tenantID,
	).Scan(
		&identity.ID,
		&identity.TenantID,
		&identity.Type,
		&identity.Name,
		&identity.Status,
		&identity.Metadata,
		&identity.LastAuthenticatedAt,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &identity, nil
}

type IdentityAuthContext struct {
	ID       uuid.UUID
	TenantID uuid.UUID
	Type     string
	Status   string
}

func (r *Repository) GetAuthContext(
	ctx context.Context,
	id uuid.UUID,
) (*IdentityAuthContext, error) {

	var result IdentityAuthContext

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			tenant_id,
			type,
			status
		FROM identities
		WHERE id = $1
		`,
		id,
	).Scan(
		&result.ID,
		&result.TenantID,
		&result.Type,
		&result.Status,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
