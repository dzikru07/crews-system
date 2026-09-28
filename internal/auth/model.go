package auth

import (
	"time"

	"github.com/google/uuid"
)

const (
	GrantTypeClientCredentials = "client_credentials"
	TokenTypeBearer            = "Bearer"
	AccessTokenType            = "access"
)

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type AuthContext struct {
	TenantID     uuid.UUID
	IdentityID   uuid.UUID
	ClientID     string
	ClientUUID   uuid.UUID
	IdentityType string
	Environment  string
}

type AccessTokenClaims struct {
	TenantID     string `json:"tenant_id"`
	IdentityID   string `json:"identity_id"`
	ClientID     string `json:"client_id"`
	IdentityType string `json:"identity_type"`
	Environment  string `json:"environment"`
	TokenType    string `json:"token_type"`

	IssuedAt  time.Time `json:"-"`
	ExpiresAt time.Time `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}
