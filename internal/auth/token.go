package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dzikru07/crews-system/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type TokenManager struct {
	config config.JWTConfig
}

func NewTokenManager(cfg config.JWTConfig) *TokenManager {
	return &TokenManager{
		config: cfg,
	}
}

type jwtClaims struct {
	TenantID     string `json:"tenant_id"`
	IdentityID   string `json:"identity_id"`
	ClientID     string `json:"client_id"`
	IdentityType string `json:"identity_type"`
	Environment  string `json:"environment"`
	TokenType    string `json:"token_type"`

	jwt.RegisteredClaims
}

func (m *TokenManager) Generate(
	ctx AuthContext,
) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(m.config.AccessTTL)

	claims := jwtClaims{
		TenantID:     ctx.TenantID.String(),
		IdentityID:   ctx.IdentityID.String(),
		ClientID:     ctx.ClientID,
		IdentityType: ctx.IdentityType,
		Environment:  ctx.Environment,
		TokenType:    AccessTokenType,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			Subject:   ctx.IdentityID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(
		[]byte(m.config.Secret),
	)
	if err != nil {
		return "", 0, fmt.Errorf(
			"sign token: %w",
			err,
		)
	}

	return signedToken, int64(m.config.AccessTTL.Seconds()), nil
}

func (m *TokenManager) Validate(
	tokenString string,
) (*AuthContext, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}

			return []byte(m.config.Secret), nil
		},
		jwt.WithIssuer(m.config.Issuer),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}

		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != AccessTokenType {
		return nil, ErrInvalidToken
	}

	tenantID, err := uuid.Parse(claims.TenantID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	identityID, err := uuid.Parse(claims.IdentityID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &AuthContext{
		TenantID:     tenantID,
		IdentityID:   identityID,
		ClientID:     claims.ClientID,
		IdentityType: claims.IdentityType,
		Environment:  claims.Environment,
	}, nil
}
