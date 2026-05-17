package refresh_tokens

import (
	"context"
	"log/slog"

	db "github.com/aniruddha-jafa/go-auth-v1/db/generated"
	"github.com/aniruddha-jafa/go-auth-v1/internal/logger"
	"github.com/aniruddha-jafa/go-auth-v1/pkg/util"
)

type RefreshTokenRepo interface {
	Create(ctx *context.Context, refreshToken RefreshToken) (RefreshToken, error)
	GetById(ctx *context.Context, id string) (RefreshToken, error)
	Revoke(ctx *context.Context, refreshTokenId string) (RefreshToken, error)
}

type RefreshTokenRepoImpl struct {
	queries    *db.Queries
	baseLogger *slog.Logger
}

func NewRefreshTokenRepoImpl(queries *db.Queries) RefreshTokenRepo {
	return &RefreshTokenRepoImpl{
		queries:    queries,
		baseLogger: slog.Default().With(logger.LoggerNameKey, "RefreshTokenRepo"),
	}
}

func (r *RefreshTokenRepoImpl) Create(ctx *context.Context, refreshToken RefreshToken) (RefreshToken, error) {
	log := logger.WithContext(r.baseLogger, *ctx)
	log.Info("Creating refresh token", "refreshTokenId", refreshToken.ID)
	refreshTokenCreated, err := r.queries.Create(*ctx, db.CreateParams{
		ID:        refreshToken.ID,
		UserID:    util.MakePgUuid(refreshToken.UserID),
		ExpiresAt: util.MakePgTimestamp(refreshToken.ExpiresAt),
		CreatedAt: util.MakePgTimestamp(refreshToken.CreatedAt),
		UpdatedAt: util.MakePgTimestamp(refreshToken.UpdatedAt),
	})
	if err != nil {
		return RefreshToken{}, err
	}
	refreshTokenDomain, err := ToDomain(&refreshTokenCreated)
	if err != nil {
		return RefreshToken{}, err
	}
	return *refreshTokenDomain, nil
}

func (r *RefreshTokenRepoImpl) GetById(ctx *context.Context, id string) (RefreshToken, error) {
	refreshToken, err := r.queries.GetById(*ctx, id)
	if err != nil {
		return RefreshToken{}, err
	}
	refreshTokenDomain, err := ToDomain(&refreshToken)
	if err != nil {
		return RefreshToken{}, err
	}
	return *refreshTokenDomain, nil
}

func (r *RefreshTokenRepoImpl) Revoke(ctx *context.Context, refreshTokenId string) (RefreshToken, error) {
	log := logger.WithContext(r.baseLogger, *ctx)
	log.Info("Revoking refresh token", "refreshTokenId", refreshTokenId)

	now := util.Now()
	revokedToken, err := r.queries.Revoke(*ctx, db.RevokeParams{
		ID:        refreshTokenId,
		RevokedAt: util.MakePgTimestamp(now),
		UpdatedAt: util.MakePgTimestamp(now),
	})
	if err != nil {
		return RefreshToken{}, err
	}
	revokedTokenDomain, err := ToDomain(&revokedToken)
	if err != nil {
		return RefreshToken{}, err
	}
	return *revokedTokenDomain, nil
}
