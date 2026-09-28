package access

import (
	"context"

	"github.com/dzikru07/crews-system/internal/api_registry"
	"github.com/dzikru07/crews-system/internal/database"
	"github.com/google/uuid"
)

type Repository struct {
	db          *database.Database
	apiRegistry *api_registry.Repository
}

func NewRepository(
	db *database.Database,
	apiRegistry *api_registry.Repository,
) *Repository {
	return &Repository{
		db:          db,
		apiRegistry: apiRegistry,
	}
}

func (r *Repository) AssignPermission(
	ctx context.Context,
	endpointID uuid.UUID,
	permissionID uuid.UUID,
) error {

	_, err := r.db.Pool.Exec(
		ctx,
		`
		INSERT INTO endpoint_permissions (
			endpoint_id,
			permission_id
		)
		VALUES ($1, $2)
		ON CONFLICT (
			endpoint_id,
			permission_id
		)
		DO NOTHING
		`,
		endpointID,
		permissionID,
	)

	return err
}

func (r *Repository) RemovePermission(
	ctx context.Context,
	endpointID uuid.UUID,
	permissionID uuid.UUID,
) error {

	_, err := r.db.Pool.Exec(
		ctx,
		`
		DELETE FROM endpoint_permissions
		WHERE endpoint_id = $1
		AND permission_id = $2
		`,
		endpointID,
		permissionID,
	)

	return err
}

func (r *Repository) ListPermissions(
	ctx context.Context,
	endpointID uuid.UUID,
) ([]EndpointPermission, error) {

	rows, err := r.db.Pool.Query(
		ctx,
		`
		SELECT
			endpoint_id,
			permission_id,
			created_at
		FROM endpoint_permissions
		WHERE endpoint_id = $1
		ORDER BY created_at ASC
		`,
		endpointID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []EndpointPermission

	for rows.Next() {
		var item EndpointPermission

		if err := rows.Scan(
			&item.EndpointID,
			&item.PermissionID,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// func (r *Repository) HasPermission(
// 	ctx context.Context,
// 	identityID uuid.UUID,
// 	endpointID uuid.UUID,
// ) (bool, error) {

// 	var exists bool

// 	err := r.db.Pool.QueryRow(
// 		ctx,
// 		`
// 		SELECT EXISTS (
// 			SELECT 1
// 			FROM identities i

// 			INNER JOIN identity_roles ir
// 				ON ir.identity_id = i.id

// 			INNER JOIN roles r
// 				ON r.id = ir.role_id

// 			INNER JOIN role_permissions rp
// 				ON rp.role_id = r.id

// 			INNER JOIN permissions p
// 				ON p.id = rp.permission_id

// 			INNER JOIN endpoint_permissions ep
// 				ON ep.permission_id = p.id

// 			INNER JOIN api_endpoints ae
// 				ON ae.id = ep.endpoint_id

// 			WHERE
// 				i.id = $1
// 				AND ae.id = $2

// 				AND i.status = 'active'
// 				AND r.status = 'active'
// 				AND p.status = 'active'
// 				AND ae.status = 'active'
// 		)
// 		`,
// 		identityID,
// 		endpointID,
// 	).Scan(&exists)

// 	if err != nil {
// 		return false, err
// 	}

//		return exists, nil
//	}
func (r *Repository) HasPermission(
	ctx context.Context,
	identityID uuid.UUID,
	endpointID uuid.UUID,
) (bool, error) {

	var exists bool

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM identities i

			INNER JOIN identity_roles ir
				ON ir.identity_id = i.id

			INNER JOIN roles r
				ON r.id = ir.role_id
				AND r.tenant_id = i.tenant_id

			INNER JOIN role_permissions rp
				ON rp.role_id = r.id

			INNER JOIN permissions p
				ON p.id = rp.permission_id
				AND p.tenant_id = i.tenant_id

			INNER JOIN endpoint_permissions ep
				ON ep.permission_id = p.id

			INNER JOIN api_endpoints ae
				ON ae.id = ep.endpoint_id

			INNER JOIN api_services s
				ON s.id = ae.service_id

			INNER JOIN applications a
				ON a.id = s.application_id

			INNER JOIN environments e
				ON e.application_id = a.id

			WHERE i.id = $1
			  AND ae.id = $2

			  AND i.status = 'active'
			  AND r.status = 'active'
			  AND p.status = 'active'
			  AND ae.status = 'active'
			  AND s.status = 'active'
		)
		`,
		identityID,
		endpointID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) GetEndpointByRoute(
	ctx context.Context,
	method string,
	path string,
) (*api_registry.APIEndpoint, error) {
	return r.apiRegistry.GetEndpointByRoute(
		ctx,
		method,
		path,
	)
}
