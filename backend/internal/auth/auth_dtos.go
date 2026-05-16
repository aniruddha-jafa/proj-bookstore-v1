package auth

import (
	"fmt"
	"time"

	"github.com/aniruddha-jafa/go-auth-v1/internal/users"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	// Don't add complex validation here, since we handled it in sign up
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
}

func NewLoginResponse(user users.User, token string, refreshToken string) LoginResponse {
	return LoginResponse{
		ID:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Token:        token,
		RefreshToken: refreshToken,
	}
}

func (r LoginResponse) String() string {
	return fmt.Sprintf("LoginResponse{ID: %s, Email: %s, CreatedAt: %s, UpdatedAt: %s, Token: %s, RefreshToken: %s}", r.ID, r.Email, r.CreatedAt, r.UpdatedAt, r.Token, r.RefreshToken)
}
