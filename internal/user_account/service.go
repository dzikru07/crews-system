package user_account

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/dzikru07/crews-system/internal/identity"
	"github.com/dzikru07/crews-system/internal/pkg/crypto"
	"github.com/google/uuid"
)

var (
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrIdentityNotFound  = errors.New("identity not found")
	ErrIdentityNotHuman  = errors.New("identity must be human")
	ErrIdentityNotActive = errors.New("identity not active")
	ErrUserAlreadyExists = errors.New("user account already exists")
)

type Service struct {
	repository   *Repository
	identityRepo *identity.Repository
}

func NewService(
	repository *Repository,
	identityRepo *identity.Repository,
) *Service {
	return &Service{
		repository:   repository,
		identityRepo: identityRepo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	tenantID uuid.UUID,
	identityID uuid.UUID,
	req CreateRequest,
) (*CreateResponse, error) {

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, ErrInvalidEmail
	}

	req.Email = strings.ToLower(
		strings.TrimSpace(req.Email),
	)

	if len(req.Password) < 8 {
		return nil, ErrInvalidPassword
	}

	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)

		if username == "" {
			return nil, ErrInvalidUsername
		}

		*req.Username = username
	}

	identityRecord, err := s.identityRepo.GetByID(
		ctx,
		identityID,
		tenantID,
	)

	if err != nil {
		return nil, ErrIdentityNotFound
	}

	if identityRecord.Type != identity.TypeHuman {
		return nil, ErrIdentityNotHuman
	}

	if identityRecord.Status != identity.StatusActive {
		return nil, ErrIdentityNotActive
	}

	_, err = s.repository.GetByIdentityID(
		ctx,
		identityID,
	)

	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	passwordHash, err := crypto.HashSecret(req.Password)
	if err != nil {
		return nil, err
	}

	account := &UserAccount{
		ID:           uuid.New(),
		IdentityID:   identityID,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
	}

	if err := s.repository.Create(ctx, account); err != nil {
		return nil, err
	}

	return &CreateResponse{
		ID:         account.ID,
		IdentityID: account.IdentityID,
		Email:      account.Email,
		Username:   account.Username,
		CreatedAt:  account.CreatedAt,
	}, nil
}
