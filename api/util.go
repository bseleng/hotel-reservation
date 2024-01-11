package api

import (
	"fmt"

	"github.com/bseleng/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func GetAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	fmt.Printf("\n\n GetAuthUser --> %+v", user)


	if !ok {
		return nil, fmt.Errorf("unauthorized")

	}
	return user, nil
}
