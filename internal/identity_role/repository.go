package identity_role

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

// func (r *Repository) AssignRole(
// 	ctx context.Context,
// 	identityID uuid.UUID,
// 	roleID uuid.UUID,
// ) error {

// 	_, err := r.db.Pool.Exec(
// 		ctx,
// 		`
// 		INSERT INTO identity_roles (
// 			identity_id,
// 			role_id
// 		)
// 		VALUES ($1, $2)
// 		ON CONFLICT (
// 			identity_id,
// 			role_id
// 		)
// 		DO NOTHING
// 		`,
// 		identityID,
// 		roleID,
// 	)

// 	return err
// }

func (r *Repository) RemoveRole(
	ctx context.Context,
	identityID uuid.UUID,
	roleID uuid.UUID,
) error {

	_, err := r.db.Pool.Exec(
		ctx,
		`
		DELETE FROM identity_roles
		WHERE identity_id = $1
		AND role_id = $2
		`,
		identityID,
		roleID,
	)

	return err
}

func (r *Repository) ListRoles(
	ctx context.Context,
	identityID uuid.UUID,
) ([]IdentityRole, error) {

	rows, err := r.db.Pool.Query(
		ctx,
		`
		SELECT
			identity_id,
			role_id,
			created_at
		FROM identity_roles
		WHERE identity_id = $1
		ORDER BY created_at ASC
		`,
		identityID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []IdentityRole

	for rows.Next() {
		var item IdentityRole

		if err := rows.Scan(
			&item.IdentityID,
			&item.RoleID,
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

func (r *Repository) AssignRole(
	ctx context.Context,
	identityID uuid.UUID,
	roleID uuid.UUID,
) error {

	var valid bool

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM identities i
			INNER JOIN roles r
				ON r.tenant_id = i.tenant_id
			WHERE i.id = $1
			  AND r.id = $2
			  AND i.status = 'active'
			  AND r.status = 'active'
		)
		`,
		identityID,
		roleID,
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
		INSERT INTO identity_roles (
			identity_id,
			role_id
		)
		VALUES ($1, $2)
		ON CONFLICT (
			identity_id,
			role_id
		)
		DO NOTHING
		`,
		identityID,
		roleID,
	)

	return err
}
