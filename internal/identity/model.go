package identity

import (
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	ID                  uuid.UUID  `json:"id"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	Type                string     `json:"type"`
	Name                string     `json:"name"`
	Status              string     `json:"status"`
	Metadata            []byte     `json:"metadata,omitempty"`
	LastAuthenticatedAt *time.Time `json:"last_authenticated_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type CreateRequest struct {
	Name     string         `json:"name" binding:"required,max=150"`
	Type     string         `json:"type" binding:"required,oneof=human machine"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

const (
	TypeHuman   = "human"
	TypeMachine = "machine"

	StatusActive    = "active"
	StatusDisabled  = "disabled"
	StatusSuspended = "suspended"
)
