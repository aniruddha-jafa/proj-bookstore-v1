package users

import (
	"log"
	"net/http"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler interface {
	// Debug only
	GetAll(ctx *fiber.Ctx) error

	Get(ctx *fiber.Ctx, d uuid.UUID) error
	Update(ctx fiber.Ctx, id uuid.UUID) error
}

type UserHandlerImpl struct {
	UserService UserService
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
	idParam := c.Params("id")
	if len(idParam) == 0 {
		return apperrors.NewHttpError(http.StatusBadRequest, "id param is required")
	}
	log.Printf("Received GET /users with id %s", idParam)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, "invalid uuid")
	}
	user, err := u.UserService.Get(&ctx, id)
	if err != nil {
		return err
	}
	log.Printf("fetched user %s", user)
	c.Status(http.StatusOK).JSON(user)
	return nil
}

func (u *UserHandlerImpl) Update(c *fiber.Ctx) error {
	log.Printf("POST users/:id")
	ctx := c.UserContext()
	idParam := c.Params("id")
	if len(idParam) == 0 {
		return apperrors.NewHttpError(http.StatusBadRequest, "id param is required")
	}
	id, err := uuid.Parse(idParam)
	if err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, "invalid uuid")
	}
	log.Printf("Trying to update user with id: %s", id)
	updateRequest := new(UserUpdateRequest)
	if err := c.BodyParser(updateRequest); err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, err.Error())
	}
	userRes, err := u.UserService.Update(&ctx, id, *updateRequest)
	if err != nil {
		return err
	}
	c.JSON(userRes)
	return nil
}
