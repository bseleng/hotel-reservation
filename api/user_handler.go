package api

import (
	"context"

	"github.com/bseleng/hotel-reservation/db"
	"github.com/bseleng/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HadleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)
	user, err := h.userStore.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HadleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Bogdan",
		LastName:  "Seleng",
	}
	return c.JSON(u)
}
