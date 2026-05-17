package users

import (
	"context"
	"log/slog"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/logger"
	"github.com/aniruddha-jafa/go-auth-v1/pkg/security"
	"github.com/google/uuid"
)

type UserService interface {
	// Debug only
	GetAll(ctx *context.Context) ([]UserResponse, error)
	GetByEmail(ctx *context.Context, email string) (User, error)
	Get(ctx *context.Context, uuid uuid.UUID) (UserResponse, error)
	Create(ctx *context.Context, userCreationRequest UserCreationRequest) (UserResponse, error)
	Update(ctx *context.Context, id uuid.UUID, userUpdateRequest UserUpdateRequest) (UserResponse, error)
	// Delete(ctx context.Context, id uuid.UUID) error
}

type UserServiceImpl struct {
	UserRepo   UserRepo
	baseLogger *slog.Logger
}

func NewUserServiceImpl(userRepo UserRepo) UserService {
	return &UserServiceImpl{
		UserRepo:   userRepo,
		baseLogger: slog.Default().With(logger.LoggerNameKey, "UserService"),
	}
}

func (u *UserServiceImpl) GetAll(ctx *context.Context) ([]UserResponse, error) {
	users, err := u.UserRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return NewUserResponseSlice(users), nil
}

func (u *UserServiceImpl) Get(ctx *context.Context, uuid uuid.UUID) (UserResponse, error) {
	log := logger.WithContext(u.baseLogger, *ctx)
	log.Info("Trying to get user with id", "userId", uuid)

	user, err := u.UserRepo.Get(ctx, uuid)
	if err != nil {
		return UserResponse{}, err
	}
	return NewUserResponse(user), nil
}

func (u *UserServiceImpl) GetByEmail(ctx *context.Context, email string) (User, error) {
	user, err := u.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *UserServiceImpl) Create(ctx *context.Context, userRequest UserCreationRequest) (UserResponse, error) {
	log := logger.WithContext(u.baseLogger, *ctx)

	log.Info("Trying to create a user")
	hashedPassword, err := security.HashPassword(userRequest.Password)
	if err != nil {
		return UserResponse{}, err
	}
	user := User{
		Email:    userRequest.Email,
		Password: hashedPassword,
	}
	emailIsTaken, err := u.UserRepo.EmailIsTaken(ctx, user.Email)
	if err != nil {
		return UserResponse{}, err
	}
	if emailIsTaken {
		return UserResponse{}, apperrors.ErrEmailTaken
	}
	user, err = u.UserRepo.Create(ctx, user)
	if err != nil {
		return UserResponse{}, err
	}
	return NewUserResponse(user), nil
}

func (u *UserServiceImpl) Update(ctx *context.Context, userId uuid.UUID, userUpdateRequest UserUpdateRequest) (UserResponse, error) {
	log := logger.WithContext(u.baseLogger, *ctx)
	log.Info("Trying to update a user", "userId", userId)

	currentUser, err := u.UserRepo.Get(ctx, userId)
	if err != nil {
		return UserResponse{}, err
	}
	// If email is the same - no change needed
	if currentUser.Email == userUpdateRequest.Email {
		return NewUserResponse(currentUser), nil
	}
	// Else check if the email is available
	emailIsTaken, err := u.UserRepo.EmailIsTaken(ctx, userUpdateRequest.Email)
	if err != nil {
		return UserResponse{}, err
	}
	if emailIsTaken {
		return UserResponse{}, apperrors.ErrEmailTaken
	}
	// Perform the update
	user, err := u.UserRepo.Update(ctx, userId, userUpdateRequest)
	if err != nil {
		return UserResponse{}, nil
	}
	return NewUserResponse(user), nil
}
