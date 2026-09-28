package role

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRequest struct {
	Name        string  `json:"name" binding:"required,max=150"`
	Description *string `json:"description,omitempty"`
}

type UpdateRequest struct {
	Name        *string `json:"name,omitempty,max=150"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty,oneof=active disabled"`
}
