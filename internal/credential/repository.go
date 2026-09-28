package credential

import (
	"context"
	"errors"

	"github.com/dzikru07/crews-system/internal/database"
	"github.com/google/uuid"
)

var (
	ErrCredentialNotFound = errors.New("credential not found")
	ErrClientIDExists     = errors.New("client id already exists")
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
	credential *MachineCredential,
) error {
	err := r.db.Pool.QueryRow(
		ctx,
		`
        INSERT INTO machine_credentials (
            id,
            identity_id,
            client_id,
            client_secret_hash,
            status,
            expires_at,
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
            NOW(),
            NOW()
        )
        RETURNING
            created_at,
            updated_at
        `,
		credential.ID,
		credential.IdentityID,
		credential.ClientID,
		credential.ClientSecretHash,
		credential.Status,
		credential.ExpiresAt,
	).Scan(
		&credential.CreatedAt,
		&credential.UpdatedAt,
	)

	return err
}

func (r *Repository) GetByClientID(
	ctx context.Context,
	clientID string,
) (*MachineCredential, error) {
	var credential MachineCredential

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			identity_id,
			client_id,
			client_secret_hash,
			status,
			expires_at,
			last_used_at,
			created_at,
			updated_at
		FROM machine_credentials
		WHERE client_id = $1
		`,
		clientID,
	).Scan(
		&credential.ID,
		&credential.IdentityID,
		&credential.ClientID,
		&credential.ClientSecretHash,
		&credential.Status,
		&credential.ExpiresAt,
		&credential.LastUsedAt,
		&credential.CreatedAt,
		&credential.UpdatedAt,
	)

	if err != nil {
		return nil, ErrCredentialNotFound
	}

	return &credential, nil
}

type MachineCredentialWithIdentity struct {
	MachineCredential

	IdentityType   string
	IdentityStatus string
	TenantID       uuid.UUID
}

func (r *Repository) GetWithIdentity(
	ctx context.Context,
	clientID string,
) (*MachineCredentialWithIdentity, error) {
	var credential MachineCredentialWithIdentity

	err := r.db.Pool.QueryRow(
		ctx,
		`
        SELECT
            mc.id,
            mc.identity_id,
            mc.client_id,
            mc.client_secret_hash,
            mc.status,
            mc.expires_at,
            mc.last_used_at,
            mc.created_at,
            mc.updated_at,

            i.type,
            i.status,
            i.tenant_id

        FROM machine_credentials mc

        INNER JOIN identities i
            ON i.id = mc.identity_id

        WHERE mc.client_id = $1
        `,
		clientID,
	).Scan(
		&credential.ID,
		&credential.IdentityID,
		&credential.ClientID,
		&credential.ClientSecretHash,
		&credential.Status,
		&credential.ExpiresAt,
		&credential.LastUsedAt,
		&credential.CreatedAt,
		&credential.UpdatedAt,

		&credential.IdentityType,
		&credential.IdentityStatus,
		&credential.TenantID,
	)

	if err != nil {
		return nil, ErrCredentialNotFound
	}

	return &credential, nil
}
