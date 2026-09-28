package identity_role

import (
	"time"

	"github.com/google/uuid"
)

type IdentityRole struct {
	IdentityID uuid.UUID `json:"identity_id"`
	RoleID     uuid.UUID `json:"role_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}
