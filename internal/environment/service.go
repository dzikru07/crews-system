package environment

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidName = errors.New("environment name is required")
	ErrInvalidType = errors.New("invalid environment type")
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
) (*Environment, error) {
	name := strings.TrimSpace(input.Name)
	envType := strings.TrimSpace(strings.ToLower(input.Type))

	if name == "" {
		return nil, ErrInvalidName
	}

	switch envType {
	case TypeDevelopment:
	case TypeStaging:
	case TypeProduction:
	default:
		return nil, ErrInvalidType
	}

	now := time.Now().UTC()

	env := Environment{
		ID:            uuid.New(),
		ApplicationID: input.ApplicationID,
		Name:          name,
		Type:          envType,
		Status:        "active",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repository.Create(ctx, env); err != nil {
		return nil, err
	}

	return &env, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Environment, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetByIDAndApplication(
	ctx context.Context,
	id uuid.UUID,
	applicationID uuid.UUID,
) (*Environment, error) {
	return s.repository.GetByIDAndApplication(
		ctx,
		id,
		applicationID,
	)
}
