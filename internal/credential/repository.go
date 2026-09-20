package credential

import (
	"context"
	"errors"

	"github.com/dzikru07/crews-system/internal/database"
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
