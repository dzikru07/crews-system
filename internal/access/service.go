package access

import (
	"context"

	"github.com/dzikru07/crews-system/internal/api_registry"
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

func (s *Service) AssignPermission(
	ctx context.Context,
	endpointID uuid.UUID,
	permissionID uuid.UUID,
) error {

	return s.repository.AssignPermission(
		ctx,
		endpointID,
		permissionID,
	)
}

func (s *Service) RemovePermission(
	ctx context.Context,
	endpointID uuid.UUID,
	permissionID uuid.UUID,
) error {

	return s.repository.RemovePermission(
		ctx,
		endpointID,
		permissionID,
	)
}

func (s *Service) ListPermissions(
	ctx context.Context,
	endpointID uuid.UUID,
) ([]EndpointPermission, error) {

	return s.repository.ListPermissions(
		ctx,
		endpointID,
	)
}

func (s *Service) HasPermission(
	ctx context.Context,
	identityID uuid.UUID,
	endpointID uuid.UUID,
) (bool, error) {

	return s.repository.HasPermission(
		ctx,
		identityID,
		endpointID,
	)
}

func (s *Service) ResolveEndpoint(
	ctx context.Context,
	method string,
	path string,
) (*api_registry.APIEndpoint, error) {

	return s.repository.GetEndpointByRoute(
		ctx,
		method,
		path,
	)
}

func (s *Service) GetEndpointByRoute(
	ctx context.Context,
	method string,
	path string,
) (*api_registry.APIEndpoint, error) {
	return s.repository.GetEndpointByRoute(
		ctx,
		method,
		path,
	)
}

func (s *Service) CheckAccess(
	ctx context.Context,
	identityID uuid.UUID,
	endpointID uuid.UUID,
) (bool, error) {

	return s.repository.HasPermission(
		ctx,
		identityID,
		endpointID,
	)
}
