package credential

import (
	"context"
	"errors"
	"time"

	"github.com/dzikru07/crews-system/internal/pkg/crypto"
)

var (
	ErrInvalidClientID     = errors.New("invalid client id")
	ErrInvalidClientSecret = errors.New("invalid client secret")
	ErrCredentialRevoked   = errors.New("credential revoked")
	ErrCredentialExpired   = errors.New("credential expired")
	ErrIdentityDisabled    = errors.New("identity disabled")
	ErrIdentitySuspended   = errors.New("identity suspended")
	ErrIdentityNotMachine  = errors.New("identity is not machine")
)

type Validator struct {
	repository *Repository
}

func NewValidator(repository *Repository) *Validator {
	return &Validator{
		repository: repository,
	}
}

func (v *Validator) Validate(
	ctx context.Context,
	clientID string,
	clientSecret string,
) (*MachineCredentialWithIdentity, error) {

	if clientID == "" {
		return nil, ErrInvalidClientID
	}

	if clientSecret == "" {
		return nil, ErrInvalidClientSecret
	}

	credential, err := v.repository.GetWithIdentity(
		ctx,
		clientID,
	)
	if err != nil {
		return nil, ErrInvalidClientID
	}

	if credential.Status == StatusRevoked {
		return nil, ErrCredentialRevoked
	}

	if credential.Status == StatusExpired {
		return nil, ErrCredentialExpired
	}

	if credential.ExpiresAt != nil &&
		time.Now().UTC().After(*credential.ExpiresAt) {
		return nil, ErrCredentialExpired
	}

	if credential.IdentityType != "machine" {
		return nil, ErrIdentityNotMachine
	}

	switch credential.IdentityStatus {
	case "disabled":
		return nil, ErrIdentityDisabled

	case "suspended":
		return nil, ErrIdentitySuspended

	case "active":
		// valid

	default:
		return nil, ErrIdentityDisabled
	}

	if !crypto.VerifySecret(
		clientSecret,
		credential.ClientSecretHash,
	) {
		return nil, ErrInvalidClientSecret
	}

	return credential, nil
}
