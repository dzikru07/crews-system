package access

import (
	"time"

	"github.com/google/uuid"
)

type EndpointPermission struct {
	EndpointID   uuid.UUID `json:"endpoint_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type AssignPermissionRequest struct {
	PermissionID uuid.UUID `json:"permission_id" binding:"required"`
}

type CheckAccessRequest struct {
	IdentityID uuid.UUID `json:"identity_id" binding:"required"`
	EndpointID uuid.UUID `json:"endpoint_id" binding:"required"`
}
