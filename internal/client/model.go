package client

import (
	"time"

	"github.com/google/uuid"
)

const (
	TypeMobile     = "mobile"
	TypeInternal   = "internal"
	TypeB2BPartner = "b2b_partner"
	TypeService    = "service"
)

type Client struct {
	ID            uuid.UUID `json:"id"`
	EnvironmentID uuid.UUID `json:"environment_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateInput struct {
	EnvironmentID uuid.UUID
	Name          string
	Slug          string
	Type          string
}
