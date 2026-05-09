package refresh_tokens

import (
	"errors"
	"fmt"
	"time"

	db "github.com/aniruddha-jafa/go-auth-v1/db/generated"
	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToDomain(dbRefreshToken *db.RefreshToken) (*RefreshToken, error) {
	if dbRefreshToken == nil {
		return nil, errors.New("dbRefreshToken is nil")
	}
	return &RefreshToken{
		ID:        dbRefreshToken.ID,
		UserID:    uuid.UUID(dbRefreshToken.UserID.Bytes),
		ExpiresAt: dbRefreshToken.ExpiresAt.Time,
		RevokedAt: dbRefreshToken.RevokedAt.Time,
		CreatedAt: dbRefreshToken.CreatedAt.Time,
		UpdatedAt: dbRefreshToken.UpdatedAt.Time,
	}, nil
}

func (r RefreshToken) String() string {
	return fmt.Sprintf("RefreshToken{ID: %s, UserID: %s, ExpiresAt: %s, RevokedAt: %s, CreatedAt: %s, UpdatedAt: %s}", r.ID, r.UserID, r.ExpiresAt, r.RevokedAt, r.CreatedAt, r.UpdatedAt)
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}
