package credential

import (
	"context"
	"errors"

	"github.com/dzikru07/crews-system/internal/pkg/crypto"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredential = errors.New("invalid credential")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	identityID uuid.UUID,
	req CreateRequest,
	prefix string,
) (*CreateResponse, error) {

	clientID, err := GenerateClientID(
		"crw_" + prefix,
	)
	if err != nil {
		return nil, err
	}

	clientSecret, err := GenerateClientSecret(
		"crs_" + prefix,
	)
	if err != nil {
		return nil, err
	}

	secretHash, err := crypto.HashSecret(clientSecret)
	if err != nil {
		return nil, err
	}

	credential := &MachineCredential{
		ID:               uuid.New(),
		IdentityID:       identityID,
		ClientID:         clientID,
		ClientSecretHash: secretHash,
		Status:           StatusActive,
		ExpiresAt:        req.ExpiresAt,
	}

	if err := s.repository.Create(ctx, credential); err != nil {
		return nil, err
	}

	return &CreateResponse{
		ID:           credential.ID,
		ClientID:     credential.ClientID,
		ClientSecret: clientSecret,
		Status:       credential.Status,
		ExpiresAt:    credential.ExpiresAt,
		CreatedAt:    credential.CreatedAt,
	}, nil
}
