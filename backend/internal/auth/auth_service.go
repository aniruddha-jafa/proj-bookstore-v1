package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/config"
	"github.com/aniruddha-jafa/go-auth-v1/internal/refresh_tokens"
	"github.com/aniruddha-jafa/go-auth-v1/internal/users"
	"github.com/aniruddha-jafa/go-auth-v1/pkg/security"
	"github.com/aniruddha-jafa/go-auth-v1/pkg/util"
)

// Short validity for testing
const DEFAULT_TOKEN_VALIDITY = time.Second * 30
const DEFAULT_REFRESH_TOKEN_VALIDITY = time.Minute * 2

type AuthService interface {
	SignUp(ctx *context.Context, signupRequest users.UserCreationRequest) (users.UserResponse, error)
	Login(ctx *context.Context, loginRequest LoginRequest) (LoginResponse, error)
	RefreshToken(ctx *context.Context, refreshToken string) (refresh_tokens.RefreshTokenResponse, error)
	Logout(ctx *context.Context, refreshToken string) error
}

type AuthServiceImpl struct {
	UserService      users.UserService
	RefreshTokenRepo refresh_tokens.RefreshTokenRepo
}

func (s *AuthServiceImpl) SignUp(ctx *context.Context, signupRequest users.UserCreationRequest) (users.UserResponse, error) {
	userRes, err := s.UserService.Create(ctx, signupRequest)
	if err != nil {
		return users.UserResponse{}, err
	}
	return userRes, nil
}

func (s *AuthServiceImpl) Login(ctx *context.Context, loginRequest LoginRequest) (LoginResponse, error) {
	log.Printf("Logging in user: %v", loginRequest)
	// Lookup email
	// Don't return a 404 to avoid leaking info that the user doesn't exist
	user, err := s.UserService.GetByEmail(ctx, loginRequest.Email)
	if err != nil {
		return LoginResponse{}, apperrors.ErrInvalidCredentials
	}
	// Compare password
	ok, err := security.VerifyPassword(loginRequest.Password, user.Password)
	if err != nil {
		return LoginResponse{}, apperrors.ErrInvalidCredentials
	}
	if !ok {
		return LoginResponse{}, apperrors.ErrInvalidCredentials
	}
	// Create token
	now := util.Now()
	token, err := security.MakeJwt(user.ID, config.InitAppConfig().JwtSecret, DEFAULT_TOKEN_VALIDITY, now)
	if err != nil {
		return LoginResponse{}, errors.New("unable to create jwt: " + err.Error())
	}
	// Create refresh token
	refreshTokenId, err := MakeRefreshToken()
	if err != nil {
		return LoginResponse{}, err
	}
	refreshToken, err := s.RefreshTokenRepo.Create(ctx, refresh_tokens.RefreshToken{
		ID:        refreshTokenId,
		UserID:    user.ID,
		ExpiresAt: now.Add(DEFAULT_REFRESH_TOKEN_VALIDITY),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		return LoginResponse{}, errors.New("unable to create refresh token: " + err.Error())
	}
	log.Printf("Refresh token created: %v", refreshToken)
	// Return response with user, token & refresh token info
	return NewLoginResponse(user, token, refreshToken.ID), nil
}

func (s *AuthServiceImpl) RefreshToken(ctx *context.Context, token string) (refresh_tokens.RefreshTokenResponse, error) {
	log.Printf("Trying to refresh token: %s", token)
	// Check if it exists
	refreshToken, err := s.RefreshTokenRepo.GetById(ctx, token)
	if err != nil {
		return refresh_tokens.RefreshTokenResponse{}, err
	}
	// Check if revoked
	now := util.Now()
	if !refreshToken.RevokedAt.IsZero() && refreshToken.RevokedAt.Before(now) {
		return refresh_tokens.RefreshTokenResponse{}, apperrors.ErrTokenRevoked
	}
	// Check if expired
	if !refreshToken.ExpiresAt.IsZero() && refreshToken.ExpiresAt.Before(now) {
		return refresh_tokens.RefreshTokenResponse{}, apperrors.ErrTokenExpired
	}
	// Token is valid - create a new JWT
	log.Printf("Creating new JWT for refresh token: %s", refreshToken)
	newJwt, err := security.MakeJwt(refreshToken.UserID, config.InitAppConfig().JwtSecret, DEFAULT_TOKEN_VALIDITY, now)
	if err != nil {
		return refresh_tokens.RefreshTokenResponse{}, errors.New("unable to create jwt: " + err.Error())
	}
	// Return response with new token
	return refresh_tokens.RefreshTokenResponse{Token: newJwt}, nil
}

func (s *AuthServiceImpl) Logout(ctx *context.Context, token string) error {
	// Get the refresh token
	refreshToken, err := s.RefreshTokenRepo.GetById(ctx, token)
	if err != nil {
		return err
	}
	// Check if already revoked
	// No need to check expiry
	now := util.Now()
	if !refreshToken.RevokedAt.IsZero() && refreshToken.RevokedAt.Before(now) {
		log.Printf("Refresh token already revoked: %v", refreshToken)
		return nil
	}
	revokedToken, err := s.RefreshTokenRepo.Revoke(ctx, token)
	if err != nil {
		return err
	}
	log.Printf("Refresh token revoked at %s: %v", now.Format(time.RFC3339), revokedToken)
	return nil
}
