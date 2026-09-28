package api_registry

import (
	"context"
	"errors"

	"github.com/dzikru07/crews-system/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrServiceNotFound       = errors.New("api service not found")
	ErrEndpointNotFound      = errors.New("api endpoint not found")
	ErrServiceExists         = errors.New("api service already exists")
	ErrEndpointExists        = errors.New("api endpoint already exists")
	ErrEndpointRouteNotFound = errors.New("api endpoint route not found")
)

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateService(
	ctx context.Context,
	applicationID uuid.UUID,
	req CreateServiceRequest,
) (*APIService, error) {

	id := uuid.New()

	var service APIService

	err := r.db.Pool.QueryRow(
		ctx,
		`
		INSERT INTO api_services (
			id,
			application_id,
			name,
			slug,
			description,
			base_path,
			status
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			'active'
		)
		RETURNING
			id,
			application_id,
			name,
			slug,
			description,
			base_path,
			status,
			created_at,
			updated_at
		`,
		id,
		applicationID,
		req.Name,
		req.Slug,
		req.Description,
		req.BasePath,
	).Scan(
		&service.ID,
		&service.ApplicationID,
		&service.Name,
		&service.Slug,
		&service.Description,
		&service.BasePath,
		&service.Status,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *Repository) GetServiceByID(
	ctx context.Context,
	applicationID uuid.UUID,
	serviceID uuid.UUID,
) (*APIService, error) {

	var service APIService

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			application_id,
			name,
			slug,
			description,
			base_path,
			status,
			created_at,
			updated_at
		FROM api_services
		WHERE
			id = $1
			AND application_id = $2
		`,
		serviceID,
		applicationID,
	).Scan(
		&service.ID,
		&service.ApplicationID,
		&service.Name,
		&service.Slug,
		&service.Description,
		&service.BasePath,
		&service.Status,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrServiceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *Repository) CreateEndpoint(
	ctx context.Context,
	serviceID uuid.UUID,
	req CreateEndpointRequest,
) (*APIEndpoint, error) {

	id := uuid.New()

	requiresAuth := true

	if req.RequiresAuth != nil {
		requiresAuth = *req.RequiresAuth
	}

	var endpoint APIEndpoint

	err := r.db.Pool.QueryRow(
		ctx,
		`
		INSERT INTO api_endpoints (
			id,
			service_id,
			name,
			method,
			path,
			description,
			status,
			requires_auth
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			'active',
			$7
		)
		RETURNING
			id,
			service_id,
			name,
			method,
			path,
			description,
			status,
			requires_auth,
			created_at,
			updated_at
		`,
		id,
		serviceID,
		req.Name,
		req.Method,
		req.Path,
		req.Description,
		requiresAuth,
	).Scan(
		&endpoint.ID,
		&endpoint.ServiceID,
		&endpoint.Name,
		&endpoint.Method,
		&endpoint.Path,
		&endpoint.Description,
		&endpoint.Status,
		&endpoint.RequiresAuth,
		&endpoint.CreatedAt,
		&endpoint.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}

func (r *Repository) GetEndpointByID(
	ctx context.Context,
	serviceID uuid.UUID,
	endpointID uuid.UUID,
) (*APIEndpoint, error) {

	var endpoint APIEndpoint

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			service_id,
			name,
			method,
			path,
			description,
			status,
			requires_auth,
			created_at,
			updated_at
		FROM api_endpoints
		WHERE
			id = $1
			AND service_id = $2
		`,
		endpointID,
		serviceID,
	).Scan(
		&endpoint.ID,
		&endpoint.ServiceID,
		&endpoint.Name,
		&endpoint.Method,
		&endpoint.Path,
		&endpoint.Description,
		&endpoint.Status,
		&endpoint.RequiresAuth,
		&endpoint.CreatedAt,
		&endpoint.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEndpointNotFound
	}

	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}

func (r *Repository) GetEndpointByRoute(
	ctx context.Context,
	method string,
	path string,
) (*APIEndpoint, error) {

	var endpoint APIEndpoint

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			ae.id,
			ae.service_id,
			ae.name,
			ae.method,
			ae.path,
			ae.description,
			ae.status,
			ae.requires_auth,
			ae.created_at,
			ae.updated_at
		FROM api_endpoints ae
		INNER JOIN api_services s
			ON s.id = ae.service_id
		WHERE
			ae.method = $1
			AND ae.path = $2
			AND ae.status = 'active'
			AND s.status = 'active'
		`,
		method,
		path,
	).Scan(
		&endpoint.ID,
		&endpoint.ServiceID,
		&endpoint.Name,
		&endpoint.Method,
		&endpoint.Path,
		&endpoint.Description,
		&endpoint.Status,
		&endpoint.RequiresAuth,
		&endpoint.CreatedAt,
		&endpoint.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEndpointRouteNotFound
	}

	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}
