package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidName = errors.New("application name is required")
	ErrInvalidSlug = errors.New("application slug is required")
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
) (*Application, error) {
	name := strings.TrimSpace(input.Name)
	slug := strings.TrimSpace(strings.ToLower(input.Slug))

	if name == "" {
		return nil, ErrInvalidName
	}

	if slug == "" {
		return nil, ErrInvalidSlug
	}

	now := time.Now().UTC()

	app := Application{
		ID:          uuid.New(),
		TenantID:    input.TenantID,
		Name:        name,
		Slug:        slug,
		Description: input.Description,
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repository.Create(ctx, app); err != nil {
		return nil, err
	}

	return &app, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Application, error) {

	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetByIDAndTenant(
	ctx context.Context,
	id uuid.UUID,
	tenantID uuid.UUID,
) (*Application, error) {
	return s.repository.GetByIDAndTenant(
		ctx,
		id,
		tenantID,
	)
}
