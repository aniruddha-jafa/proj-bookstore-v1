package users

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	db "github.com/aniruddha-jafa/go-auth-v1/db/generated"
	"github.com/aniruddha-jafa/go-auth-v1/internal/config"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u User) String() string {
	return fmt.Sprintf("User{ID: %s, Email: %s, Password: %s, CreatedAt: %s, UpdatedAt: %s}", u.ID, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
}

type UserCreationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserCreationRequest) String() string {
	return fmt.Sprintf("UserCreationRequest{Email: %s, Password: %s}", u.Email, u.Password)
}

func passwordValidator(value interface{}) error {
	minPasswordLength := config.InitAppConfig().MinPasswordLength
	p, ok := value.(string)
	if !ok {
		return errors.New("password must be string")
	}
	if len(p) < minPasswordLength {
		return fmt.Errorf("password must have at least %d chars", minPasswordLength)
	}
	// Example rules: At least one digit, one uppercase, one lowercase
	var (
		hasNumber = regexp.MustCompile(`[0-9]`).MatchString(p)
		hasUpper  = regexp.MustCompile(`[A-Z]`).MatchString(p)
		hasLower  = regexp.MustCompile(`[a-z]`).MatchString(p)
	)
	if !hasNumber || !hasUpper || !hasLower {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}
	return nil
}

func (u UserCreationRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.EmailFormat),
		validation.Field(&u.Password, validation.Required, validation.By(passwordValidator)),
	)
}

type UserUpdateRequest struct {
	Email string `json:"email"`
}

func (u UserUpdateRequest) String() string {
	return fmt.Sprintf("UserUpdateRequest{Email: %s}", u.Email)
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u UserResponse) String() string {
	return fmt.Sprintf("UserResponse{ID: %s, Email: %s, CreatedAt: %s, UpdatedAt: %s}", u.ID, u.Email, u.CreatedAt, u.UpdatedAt)
}

func NewUserResponse(user User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewUserResponseSlice(users []User) []UserResponse {
	res := make([]UserResponse, 0, len(users))
	for _, user := range users {
		res = append(res, NewUserResponse(user))
	}
	return res
}

// ToDomain converts a db.User to a User
func ToDomain(dbUser *db.User) (*User, error) {
	if dbUser == nil {
		return nil, errors.New("dbUser is nil")
	}

	// Assumption - these fields are not nil
	return &User{
		ID:        uuid.UUID(dbUser.ID.Bytes),
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}
