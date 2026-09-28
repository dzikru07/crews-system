package role_permission

import (
	"context"

	"github.com/dzikru07/crews-system/internal/database"
	"github.com/google/uuid"
)

type Repository struct {
	db *database.Database
}

func NewRepository(
	db *database.Database,
) *Repository {
	return &Repository{
		db: db,
	}
}

// func (r *Repository) AssignPermission(
// 	ctx context.Context,
// 	roleID uuid.UUID,
// 	permissionID uuid.UUID,
// ) error {

// 	_, err := r.db.Pool.Exec(
// 		ctx,
// 		`
// 		INSERT INTO role_permissions (
// 			role_id,
// 			permission_id
// 		)
// 		VALUES ($1, $2)
// 		ON CONFLICT (
// 			role_id,
// 			permission_id
// 		)
// 		DO NOTHING
// 		`,
// 		roleID,
// 		permissionID,
// 	)

// 	return err
// }

func (r *Repository) RemovePermission(
	ctx context.Context,
	roleID uuid.UUID,
	permissionID uuid.UUID,
) error {

	_, err := r.db.Pool.Exec(
		ctx,
		`
		DELETE FROM role_permissions
		WHERE role_id = $1
		AND permission_id = $2
		`,
		roleID,
		permissionID,
	)

	return err
}

func (r *Repository) ListPermissions(
	ctx context.Context,
	roleID uuid.UUID,
) ([]RolePermission, error) {

	rows, err := r.db.Pool.Query(
		ctx,
		`
		SELECT
			role_id,
			permission_id,
			created_at
		FROM role_permissions
		WHERE role_id = $1
		ORDER BY created_at ASC
		`,
		roleID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []RolePermission

	for rows.Next() {
		var item RolePermission

		if err := rows.Scan(
			&item.RoleID,
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

func (r *Repository) AssignPermission(
	ctx context.Context,
	roleID uuid.UUID,
	permissionID uuid.UUID,
) error {

	var valid bool

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM roles r
			INNER JOIN permissions p
				ON p.tenant_id = r.tenant_id
			WHERE r.id = $1
			  AND p.id = $2
			  AND r.status = 'active'
			  AND p.status = 'active'
		)
		`,
		roleID,
		permissionID,
	).Scan(&valid)

	if err != nil {
		return err
	}

	if !valid {
		return ErrInvalidAssignment
	}

	_, err = r.db.Pool.Exec(
		ctx,
		`
		INSERT INTO role_permissions (
			role_id,
			permission_id
		)
		VALUES ($1, $2)
		ON CONFLICT (
			role_id,
			permission_id
		)
		DO NOTHING
		`,
		roleID,
		permissionID,
	)

	return err
}
