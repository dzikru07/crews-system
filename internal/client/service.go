package client

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidName = errors.New("client name is required")
	ErrInvalidSlug = errors.New("client slug is required")
	ErrInvalidType = errors.New("invalid client type")
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
	input CreateInput,
) (*Client, error) {
	name := strings.TrimSpace(input.Name)
	slug := strings.TrimSpace(strings.ToLower(input.Slug))
	clientType := strings.TrimSpace(strings.ToLower(input.Type))

	if name == "" {
		return nil, ErrInvalidName
	}

	if slug == "" {
		return nil, ErrInvalidSlug
	}

	switch clientType {
	case TypeMobile:
	case TypeInternal:
	case TypeB2BPartner:
	case TypeService:
	default:
		return nil, ErrInvalidType
	}

	now := time.Now().UTC()

	client := Client{
		ID:            uuid.New(),
		EnvironmentID: input.EnvironmentID,
		Name:          name,
		Slug:          slug,
		Type:          clientType,
		Status:        "active",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repository.Create(ctx, client); err != nil {
		return nil, err
	}

	return &client, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Client, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetByIDAndEnvironment(
	ctx context.Context,
	id uuid.UUID,
	environmentID uuid.UUID,
) (*Client, error) {
	return s.repository.GetByIDAndEnvironment(
		ctx,
		id,
		environmentID,
	)
}
