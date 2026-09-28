package permission

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRequest struct {
	Name        string  `json:"name" binding:"required,max=150"`
	Description *string `json:"description,omitempty"`
	Resource    string  `json:"resource" binding:"required,max=100"`
	Action      string  `json:"action" binding:"required,max=100"`
}
