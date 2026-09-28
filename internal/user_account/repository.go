package user_account

import (
	"context"
	"errors"

	"github.com/dzikru07/crews-system/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound       = errors.New("user account not found")
	ErrEmailExists    = errors.New("email already exists")
	ErrUsernameExists = errors.New("username already exists")
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
	account *UserAccount,
) error {
	err := r.db.Pool.QueryRow(
		ctx,
		`
		INSERT INTO user_accounts (
			id,
			identity_id,
			email,
			username,
			password_hash,
			email_verified_at,
			last_login_at,
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
			NOW(),
			NOW()
		)
		RETURNING
			created_at,
			updated_at
		`,
		account.ID,
		account.IdentityID,
		account.Email,
		account.Username,
		account.PasswordHash,
		account.EmailVerifiedAt,
		account.LastLoginAt,
	).Scan(
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	return err
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*UserAccount, error) {
	var account UserAccount

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			identity_id,
			email,
			username,
			password_hash,
			email_verified_at,
			last_login_at,
			created_at,
			updated_at
		FROM user_accounts
		WHERE id = $1
		`,
		id,
	).Scan(
		&account.ID,
		&account.IdentityID,
		&account.Email,
		&account.Username,
		&account.PasswordHash,
		&account.EmailVerifiedAt,
		&account.LastLoginAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *Repository) GetByIdentityID(
	ctx context.Context,
	identityID uuid.UUID,
) (*UserAccount, error) {
	var account UserAccount

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			identity_id,
			email,
			username,
			password_hash,
			email_verified_at,
			last_login_at,
			created_at,
			updated_at
		FROM user_accounts
		WHERE identity_id = $1
		`,
		identityID,
	).Scan(
		&account.ID,
		&account.IdentityID,
		&account.Email,
		&account.Username,
		&account.PasswordHash,
		&account.EmailVerifiedAt,
		&account.LastLoginAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *Repository) GetByEmail(
	ctx context.Context,
	email string,
) (*UserAccount, error) {
	var account UserAccount

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			identity_id,
			email,
			username,
			password_hash,
			email_verified_at,
			last_login_at,
			created_at,
			updated_at
		FROM user_accounts
		WHERE email = $1
		`,
		email,
	).Scan(
		&account.ID,
		&account.IdentityID,
		&account.Email,
		&account.Username,
		&account.PasswordHash,
		&account.EmailVerifiedAt,
		&account.LastLoginAt,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *Repository) UpdateLastLogin(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := r.db.Pool.Exec(
		ctx,
		`
		UPDATE user_accounts
		SET
			last_login_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		`,
		id,
	)

	return err
}
