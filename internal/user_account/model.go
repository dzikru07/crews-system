package user_account

import (
	"time"

	"github.com/google/uuid"
)

type UserAccount struct {
	ID              uuid.UUID  `json:"id"`
	IdentityID      uuid.UUID  `json:"identity_id"`
	Email           string     `json:"email"`
	Username        *string    `json:"username,omitempty"`
	PasswordHash    string     `json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CreateRequest struct {
	Email    string  `json:"email"`
	Username *string `json:"username,omitempty"`
	Password string  `json:"password"`
}

type CreateResponse struct {
	ID         uuid.UUID `json:"id"`
	IdentityID uuid.UUID `json:"identity_id"`
	Email      string    `json:"email"`
	Username   *string   `json:"username,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
