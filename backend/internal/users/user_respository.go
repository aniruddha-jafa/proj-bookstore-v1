package users

import (
	"context"
	"log/slog"

	db "github.com/aniruddha-jafa/go-auth-v1/db/generated"
	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/logger"
	"github.com/aniruddha-jafa/go-auth-v1/pkg/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepo interface {
	// Debug only
	GetAll(ctx *context.Context) ([]User, error)

	EmailIsTaken(ctx *context.Context, email string) (bool, error)
	Get(ctx *context.Context, id uuid.UUID) (User, error)
	GetByEmail(ctx *context.Context, email string) (User, error)
	Create(ctx *context.Context, user User) (User, error)

	Update(ctx *context.Context, id uuid.UUID, userUpdateRequest UserUpdateRequest) (User, error)
	Delete(ctx *context.Context, id uuid.UUID) error
}

type UserRepoImpl struct {
	queries    *db.Queries
	baseLogger *slog.Logger
}

func NewUserRepoImpl(queries *db.Queries) UserRepo {
	return &UserRepoImpl{
		queries:    queries,
		baseLogger: slog.Default().With(logger.LoggerNameKey, "UserRepo"),
	}
}

func (r *UserRepoImpl) Get(ctx *context.Context, id uuid.UUID) (User, error) {
	// Assumption - uuid syntax has been validated at this point
	pgUUID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
	user, err := r.queries.GetUser(*ctx, pgUUID)
	if err != nil {
		return User{}, err
	}
	domainUser, err := ToDomain(&user)
	if err != nil {
		return User{}, err
	}
	return *domainUser, err
}

func (r *UserRepoImpl) GetByEmail(ctx *context.Context, email string) (User, error) {
	user, err := r.queries.GetByEmail(*ctx, email)
	if err != nil {
		return User{}, err
	}
	userDomain, err := ToDomain(&user)
	if err != nil {
		return User{}, err
	}
	return *userDomain, nil
}

func (r *UserRepoImpl) GetAll(ctx *context.Context) ([]User, error) {
	log := logger.WithContext(r.baseLogger, *ctx)
	pgUsers, err := r.queries.GetAllUsers(*ctx)
	if err != nil {
		return nil, err
	}
	log.Info("Found users", "count", len(pgUsers))
	users := make([]User, 0, len(pgUsers))
	for _, u := range pgUsers {
		domainUser, errDomainUser := ToDomain(&u)
		if errDomainUser != nil {
			return nil, errDomainUser
		}
		users = append(users, *domainUser)
	}
	return users, nil
}

func (r *UserRepoImpl) Create(ctx *context.Context, user User) (User, error) {
	pgUUID := util.MakePgUuid(uuid.New())
	timestamp := util.MakePgTimestamp(util.Now())
	params := db.CreateUserParams{
		ID:        pgUUID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
	userCreated, err := r.queries.CreateUser(*ctx, params)
	if err != nil {
		return User{}, err
	}
	userDomain, err := ToDomain(&userCreated)
	if err != nil {
		return User{}, err
	}
	return *userDomain, nil
}

func (r *UserRepoImpl) Update(ctx *context.Context, userId uuid.UUID, userUpdateRequest UserUpdateRequest) (User, error) {
	pgUuid := util.MakePgUuid(userId)
	pgUpdateTimestamp := util.MakePgTimestamp(util.Now())
	params := db.UpdateUserParams{
		ID:        pgUuid,
		Email:     userUpdateRequest.Email,
		UpdatedAt: pgUpdateTimestamp,
	}
	userUpdated, err := r.queries.UpdateUser(*ctx, params)
	userDomain, err := ToDomain(&userUpdated)
	if err != nil {
		return User{}, err
	}
	return *userDomain, nil
}

func (r *UserRepoImpl) EmailIsTaken(ctx *context.Context, email string) (bool, error) {
	isTaken, err := r.queries.EmailIsTaken(*ctx, email)
	if err != nil {
		return false, apperrors.ErrEmailTaken
	}
	return isTaken, nil
}

func (r *UserRepoImpl) Delete(ctx *context.Context, id uuid.UUID) error {
	pgUUID := util.MakePgUuid(id)
	err := r.queries.DeleteUser(*ctx, pgUUID)
	if err != nil {
		return err
	}
	return nil
}
