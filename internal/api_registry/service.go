package api_registry

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidMethod = errors.New("invalid http method")
	ErrInvalidPath   = errors.New("invalid endpoint path")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateService(
	ctx context.Context,
	applicationID uuid.UUID,
	req CreateServiceRequest,
) (*APIService, error) {

	req.Name = strings.TrimSpace(req.Name)
	req.Slug = strings.TrimSpace(req.Slug)

	if req.Name == "" || req.Slug == "" {
		return nil, errors.New("name and slug are required")
	}

	return s.repository.CreateService(
		ctx,
		applicationID,
		req,
	)
}

func (s *Service) CreateEndpoint(
	ctx context.Context,
	serviceID uuid.UUID,
	req CreateEndpointRequest,
) (*APIEndpoint, error) {

	req.Name = strings.TrimSpace(req.Name)
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	req.Path = strings.TrimSpace(req.Path)

	if req.Name == "" {
		return nil, errors.New("endpoint name is required")
	}

	if req.Path == "" || !strings.HasPrefix(req.Path, "/") {
		return nil, ErrInvalidPath
	}

	switch req.Method {
	case MethodGET,
		MethodPOST,
		MethodPUT,
		MethodPATCH,
		MethodDELETE:

	default:
		return nil, ErrInvalidMethod
	}

	return s.repository.CreateEndpoint(
		ctx,
		serviceID,
		req,
	)
}
