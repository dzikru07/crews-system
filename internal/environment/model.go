package environment

import (
	"time"

	"github.com/google/uuid"
)

const (
	TypeDevelopment = "development"
	TypeStaging     = "staging"
	TypeProduction  = "production"
)

type Environment struct {
	ID            uuid.UUID `json:"id"`
	ApplicationID uuid.UUID `json:"application_id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateInput struct {
	ApplicationID uuid.UUID
	Name          string
	Type          string
}
