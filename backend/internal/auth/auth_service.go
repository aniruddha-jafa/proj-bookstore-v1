package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/config"
	"github.com/aniruddha-jafa/go-auth-v1/internal/logger"
	"github.com/aniruddha-jafa/go-auth-v1/internal/refresh_tokens"
	"github.com/aniruddha-jafa/go-auth-v1/internal/users"
	"github.com/aniruddha-jafa/go-auth-v1/pkg/security"
	"github.com/aniruddha-jafa/go-auth-v1/pkg/util"
)

type AuthService interface {
	SignUp(ctx *context.Context, signupRequest users.UserCreationRequest) (users.UserResponse, error)
	Login(ctx *context.Context, loginRequest LoginRequest) (LoginResponse, error)
	RefreshToken(ctx *context.Context, refreshToken string) (refresh_tokens.RefreshTokenResponse, error)
	Logout(ctx *context.Context, refreshToken string) error
}

type AuthServiceImpl struct {
	UserService      users.UserService
	RefreshTokenRepo refresh_tokens.RefreshTokenRepo
	baseLogger       *slog.Logger
}

func NewAuthServiceImpl(userService users.UserService, refreshTokenRepo refresh_tokens.RefreshTokenRepo) AuthService {
	return &AuthServiceImpl{
		UserService:      userService,
		RefreshTokenRepo: refreshTokenRepo,
		baseLogger:       slog.Default().With(logger.LoggerNameKey, "AuthService"),
	}
}

func (s *AuthServiceImpl) SignUp(ctx *context.Context, signupRequest users.UserCreationRequest) (users.UserResponse, error) {
	userRes, err := s.UserService.Create(ctx, signupRequest)
	if err != nil {
		return users.UserResponse{}, err
	}
	return userRes, nil
}

func (s *AuthServiceImpl) Login(ctx *context.Context, loginRequest LoginRequest) (LoginResponse, error) {
	log := logger.WithContext(s.baseLogger, *ctx)
	log.Info("Logging in user", "email", loginRequest.Email)
	appConfig := config.InitAppConfig()

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
	token, err := security.MakeJwt(user.ID, appConfig.JwtSecret, appConfig.AccessTokenValidity, now)
	if err != nil {
		return LoginResponse{}, errors.New("unable to create jwt: " + err.Error())
	}
	// Create refresh token
	refreshTokenId, err := makeRefreshToken()
	if err != nil {
		return LoginResponse{}, err
	}
	refreshToken, err := s.RefreshTokenRepo.Create(ctx, refresh_tokens.RefreshToken{
		ID:        refreshTokenId,
		UserID:    user.ID,
		ExpiresAt: now.Add(appConfig.RefreshTokenValidity),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		log.Error("Error creating refresh token", "error", err)
		return LoginResponse{}, errors.New("unable to create refresh token: " + err.Error())
	}
	log.Info("Refresh token created", "userId", user.ID, "refreshTokenId", refreshToken.ID)
	csrfTokenValue, err := generateCSRFToken()
	if err != nil {
		return LoginResponse{}, errors.New("unable to generate CSRF token: " + err.Error())
	}
	// Return response with user, token & refresh token info
	return NewLoginResponse(user, token, refreshToken, csrfTokenValue), nil
}

func (s *AuthServiceImpl) RefreshToken(ctx *context.Context, token string) (refresh_tokens.RefreshTokenResponse, error) {
	log := logger.WithContext(s.baseLogger, *ctx)
	log.Info("Trying to refresh token")

	// Check if it exists
	refreshToken, err := s.RefreshTokenRepo.GetById(ctx, token)
	if err != nil {
		log.Error("Error getting refresh token", "error", err)
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
	log.Info("Creating new JWT for refresh token", "refreshTokenId", refreshToken.ID)
	appConfig := config.InitAppConfig()
	newJwt, err := security.MakeJwt(refreshToken.UserID, appConfig.JwtSecret, appConfig.AccessTokenValidity, now)
	if err != nil {
		log.Error("Error creating new JWT", "error", err)
		return refresh_tokens.RefreshTokenResponse{}, errors.New("unable to create jwt: " + err.Error())
	}
	// Return response with new token
	return refresh_tokens.RefreshTokenResponse{Token: newJwt, UserID: refreshToken.UserID}, nil
}

func (s *AuthServiceImpl) Logout(ctx *context.Context, token string) error {
	log := logger.WithContext(s.baseLogger, *ctx)
	log.Info("Logging out user")

	// Get the refresh token
	refreshToken, err := s.RefreshTokenRepo.GetById(ctx, token)
	if err != nil {
		return err
	}
	// Check if already revoked
	// No need to check expiry
	now := util.Now()
	if !refreshToken.RevokedAt.IsZero() && refreshToken.RevokedAt.Before(now) {
		log.Info("Refresh token already revoked", "refreshTokenId", refreshToken.ID, "revokedAt", refreshToken.RevokedAt.Format(time.RFC3339))
		return nil
	}
	revokedToken, err := s.RefreshTokenRepo.Revoke(ctx, token)
	if err != nil {
		return err
	}
	log.Info("Refresh token revoked", "refreshTokenId", revokedToken.ID, "revokedAt", now.Format(time.RFC3339))
	return nil
}
