package api_registry

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"

	MethodGET    = "GET"
	MethodPOST   = "POST"
	MethodPUT    = "PUT"
	MethodPATCH  = "PATCH"
	MethodDELETE = "DELETE"
)

type APIService struct {
	ID            uuid.UUID `json:"id"`
	ApplicationID uuid.UUID `json:"application_id"`

	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	BasePath    *string `json:"base_path,omitempty"`

	Status string `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type APIEndpoint struct {
	ID        uuid.UUID `json:"id"`
	ServiceID uuid.UUID `json:"service_id"`

	Name        string  `json:"name"`
	Method      string  `json:"method"`
	Path        string  `json:"path"`
	Description *string `json:"description,omitempty"`

	Status       string `json:"status"`
	RequiresAuth bool   `json:"requires_auth"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateServiceRequest struct {
	Name        string  `json:"name" binding:"required,max=150"`
	Slug        string  `json:"slug" binding:"required,max=100"`
	Description *string `json:"description"`
	BasePath    *string `json:"base_path"`
}

type CreateEndpointRequest struct {
	Name         string  `json:"name" binding:"required,max=150"`
	Method       string  `json:"method" binding:"required"`
	Path         string  `json:"path" binding:"required,max=255"`
	Description  *string `json:"description"`
	RequiresAuth *bool   `json:"requires_auth"`
}
