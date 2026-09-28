package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/dzikru07/crews-system/internal/credential"
	"github.com/dzikru07/crews-system/internal/identity"
	"github.com/dzikru07/crews-system/internal/pkg/crypto"
)

var (
	ErrInvalidGrant         = errors.New("invalid grant")
	ErrInvalidClient        = errors.New("invalid client")
	ErrCredentialInactive   = errors.New("credential inactive")
	ErrIdentityInactive     = errors.New("identity inactive")
	ErrCredentialExpired    = errors.New("credential expired")
	ErrUnsupportedGrantType = errors.New("unsupported grant type")

	ErrInvalidLogin = errors.New("invalid login")
)

type Service struct {
	repository   *Repository
	tokenManager *TokenManager
}

func NewService(
	repository *Repository,
	tokenManager *TokenManager,
) *Service {
	return &Service{
		repository:   repository,
		tokenManager: tokenManager,
	}
}

func (s *Service) Token(
	ctx context.Context,
	req TokenRequest,
	environment string,
) (*TokenResponse, error) {
	if req.GrantType == "" {
		req.GrantType = GrantTypeClientCredentials
	}

	if req.GrantType != GrantTypeClientCredentials {
		return nil, ErrUnsupportedGrantType
	}

	if strings.TrimSpace(req.ClientID) == "" ||
		strings.TrimSpace(req.ClientSecret) == "" {
		return nil, ErrInvalidClient
	}

	credentialRecord, err := s.repository.GetCredential(
		ctx,
		req.ClientID,
	)
	if err != nil {
		return nil, ErrInvalidClient
	}

	if credentialRecord.Status != credential.StatusActive {
		return nil, ErrCredentialInactive
	}

	if credentialRecord.IdentityType != "machine" {
		return nil, ErrInvalidClient
	}

	if credentialRecord.IdentityStatus != "active" {
		return nil, ErrIdentityInactive
	}

	if credentialRecord.ExpiresAt != nil &&
		time.Now().After(*credentialRecord.ExpiresAt) {
		return nil, ErrCredentialExpired
	}

	if !crypto.VerifySecret(
		req.ClientSecret,
		credentialRecord.ClientSecretHash,
	) {
		return nil, ErrInvalidClient
	}

	authContext := AuthContext{
		TenantID:     credentialRecord.TenantID,
		IdentityID:   credentialRecord.IdentityID,
		ClientID:     credentialRecord.ClientID,
		IdentityType: credentialRecord.IdentityType,
		Environment:  environment,
	}

	accessToken, expiresIn, err := s.tokenManager.Generate(
		authContext,
	)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken: accessToken,
		TokenType:   TokenTypeBearer,
		ExpiresIn:   expiresIn,
	}, nil
}

func (s *Service) Login(
	ctx context.Context,
	req LoginRequest,
	environment string,
) (*LoginResponse, error) {

	req.Email = strings.TrimSpace(
		strings.ToLower(req.Email),
	)

	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidLogin
	}

	account, err := s.repository.GetUserByEmail(
		ctx,
		req.Email,
	)
	if err != nil {
		return nil, ErrInvalidLogin
	}

	if !crypto.VerifySecret(
		req.Password,
		account.PasswordHash,
	) {
		return nil, ErrInvalidLogin
	}

	authContext, err := s.repository.GetIdentityContext(
		ctx,
		account.IdentityID,
	)
	if err != nil {
		return nil, ErrInvalidLogin
	}

	if authContext.Type != identity.TypeHuman {
		return nil, ErrInvalidLogin
	}

	if authContext.Status != identity.StatusActive {
		return nil, ErrInvalidLogin
	}

	token, expiresIn, err := s.tokenManager.Generate(
		AuthContext{
			TenantID:     authContext.TenantID,
			IdentityID:   authContext.ID,
			ClientID:     "",
			IdentityType: authContext.Type,
			Environment:  environment,
		},
	)
	if err != nil {
		return nil, err
	}

	if err := s.repository.UpdateLastLogin(
		ctx,
		account.ID,
	); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		TokenType:   TokenTypeBearer,
		ExpiresIn:   expiresIn,
	}, nil
}
