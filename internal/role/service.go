package role

import (
	"context"

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

func (s *Service) Create(
	ctx context.Context,
	tenantID uuid.UUID,
	req CreateRequest,
) (*Role, error) {

	return s.repository.Create(
		ctx,
		tenantID,
		req,
	)
}

func (s *Service) GetByID(
	ctx context.Context,
	tenantID uuid.UUID,
	roleID uuid.UUID,
) (*Role, error) {

	return s.repository.GetByID(
		ctx,
		tenantID,
		roleID,
	)
}
