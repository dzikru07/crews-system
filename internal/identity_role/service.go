package identity_role

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(
	repository *Repository,
) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AssignRole(
	ctx context.Context,
	identityID uuid.UUID,
	roleID uuid.UUID,
) error {

	return s.repository.AssignRole(
		ctx,
		identityID,
		roleID,
	)
}

func (s *Service) RemoveRole(
	ctx context.Context,
	identityID uuid.UUID,
	roleID uuid.UUID,
) error {

	return s.repository.RemoveRole(
		ctx,
		identityID,
		roleID,
	)
}

func (s *Service) ListRoles(
	ctx context.Context,
	identityID uuid.UUID,
) ([]IdentityRole, error) {

	return s.repository.ListRoles(
		ctx,
		identityID,
	)
}
