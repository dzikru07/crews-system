package tenant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

type CreateInput struct {
	Name string
	Slug string
	Type string
}

func (s *Service) Create(
	ctx context.Context,
	input CreateInput,
) (*Tenant, error) {
	name := strings.TrimSpace(input.Name)
	slug := strings.TrimSpace(input.Slug)
	tenantType := strings.TrimSpace(input.Type)

	if name == "" {
		return nil, fmt.Errorf("tenant name is required")
	}

	if slug == "" {
		return nil, fmt.Errorf("tenant slug is required")
	}

	if tenantType == "" {
		return nil, fmt.Errorf("tenant type is required")
	}

	now := time.Now().UTC()

	tenant := Tenant{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		Type:      tenantType,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Create(ctx, tenant); err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Tenant, error) {
	return s.repository.GetByID(ctx, id)
}
