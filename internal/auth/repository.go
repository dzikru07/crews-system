package auth

import (
	"context"

	"github.com/dzikru07/crews-system/internal/credential"
	"github.com/dzikru07/crews-system/internal/identity"
	"github.com/dzikru07/crews-system/internal/user_account"
	"github.com/google/uuid"
)

type Repository struct {
	credentialRepo  *credential.Repository
	userAccountRepo *user_account.Repository
	identityRepo    *identity.Repository
}

func NewRepository(
	credentialRepo *credential.Repository,
	userAccountRepo *user_account.Repository,
	identityRepo *identity.Repository,
) *Repository {
	return &Repository{
		credentialRepo:  credentialRepo,
		userAccountRepo: userAccountRepo,
		identityRepo:    identityRepo,
	}
}

func (r *Repository) GetCredential(
	ctx context.Context,
	clientID string,
) (*credential.MachineCredentialWithIdentity, error) {
	return r.credentialRepo.GetWithIdentity(
		ctx,
		clientID,
	)
}

func (r *Repository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*user_account.UserAccount, error) {
	return r.userAccountRepo.GetByEmail(
		ctx,
		email,
	)
}

func (r *Repository) UpdateLastLogin(
	ctx context.Context,
	userID uuid.UUID,
) error {
	return r.userAccountRepo.UpdateLastLogin(
		ctx,
		userID,
	)
}

func (r *Repository) GetIdentityContext(
	ctx context.Context,
	identityID uuid.UUID,
) (*identity.IdentityAuthContext, error) {
	return r.identityRepo.GetAuthContext(
		ctx,
		identityID,
	)
}
