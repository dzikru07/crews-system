package role_permission

import (
	"time"

	"github.com/google/uuid"
)

type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type AssignPermissionRequest struct {
	PermissionID uuid.UUID `json:"permission_id" binding:"required"`
}
