package credential

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusActive  = "active"
	StatusRevoked = "revoked"
	StatusExpired = "expired"
)

type MachineCredential struct {
	ID               uuid.UUID  `json:"id"`
	IdentityID       uuid.UUID  `json:"identity_id"`
	ClientID         string     `json:"client_id"`
	ClientSecretHash string     `json:"-"`
	Status           string     `json:"status"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	LastUsedAt       *time.Time `json:"last_used_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type CreateRequest struct {
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type CreateResponse struct {
	ID           uuid.UUID  `json:"id"`
	ClientID     string     `json:"client_id"`
	ClientSecret string     `json:"client_secret"`
	Status       string     `json:"status"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
