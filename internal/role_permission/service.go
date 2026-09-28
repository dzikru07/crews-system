package role_permission

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

func (s *Service) AssignPermission(
	ctx context.Context,
	roleID uuid.UUID,
	permissionID uuid.UUID,
) error {

	return s.repository.AssignPermission(
		ctx,
		roleID,
		permissionID,
	)
}

func (s *Service) RemovePermission(
	ctx context.Context,
	roleID uuid.UUID,
	permissionID uuid.UUID,
) error {

	return s.repository.RemovePermission(
		ctx,
		roleID,
		permissionID,
	)
}

func (s *Service) ListPermissions(
	ctx context.Context,
	roleID uuid.UUID,
) ([]RolePermission, error) {

	return s.repository.ListPermissions(
		ctx,
		roleID,
	)
}
