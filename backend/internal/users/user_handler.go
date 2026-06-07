package users

import (
	"log/slog"
	"net/http"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler interface {
	// Debug only
	GetAll(ctx *fiber.Ctx) error

	Get(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type UserHandlerImpl struct {
	UserService UserService
	baseLogger  *slog.Logger
}

func NewUserHandlerImpl(userService UserService) UserHandler {
	return &UserHandlerImpl{
		UserService: userService,
		baseLogger:  slog.Default().With(logger.LoggerNameKey, "UserHandler"),
	}
}

func (u *UserHandlerImpl) GetAll(c *fiber.Ctx) error {
	ctx := c.UserContext()
	users, err := u.UserService.GetAll(&ctx)
	if err != nil {
		return err
	}
	c.Status(http.StatusOK).JSON(users)
	return nil
}

func (u *UserHandlerImpl) Get(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.WithContext(u.baseLogger, ctx)
	idParam := c.Params("id")
	if len(idParam) == 0 {
		return apperrors.NewHttpError(http.StatusBadRequest, "id param is required")
	}
	log.Info("Received GET /users with id", "id", idParam)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, "invalid uuid")
	}
	user, err := u.UserService.Get(&ctx, id)
	if err != nil {
		return err
	}
	log.Info("Fetched user", "user", user.String())
	c.Status(http.StatusOK).JSON(user)
	return nil
}

func (u *UserHandlerImpl) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.WithContext(u.baseLogger, ctx)
	log.Info("POST users/:id")

	idParam := c.Params("id")
	if len(idParam) == 0 {
		return apperrors.NewHttpError(http.StatusBadRequest, "id param is required")
	}
	id, err := uuid.Parse(idParam)
	if err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, "invalid uuid")
	}
	log.Info("Trying to update user", "id", id)
	updateRequest := new(UserUpdateRequest)
	if err := c.BodyParser(updateRequest); err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, err.Error())
	}
	userRes, err := u.UserService.Update(&ctx, id, *updateRequest)
	if err != nil {
		return err
	}
	log.Info("User updated", "userId", userRes.ID, "email", userRes.Email)
	c.JSON(userRes)
	return nil
}

func (u *UserHandlerImpl) Delete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.WithContext(u.baseLogger, ctx)
	log.Info("DELETE users/:id")

	idParam := c.Params("id")
	if len(idParam) == 0 {
		return apperrors.NewHttpError(http.StatusBadRequest, "id param is required")
	}
	id, err := uuid.Parse(idParam)
	if err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, "invalid uuid")
	}
	log.Info("Trying to delete user", "id", id)
	err = u.UserService.Delete(&ctx, id)
	if err != nil {
		return err
	}
	log.Info("User deleted", "userId", id)
	c.Status(http.StatusNoContent)
	return nil
}
