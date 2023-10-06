package api

import (
	"github.com/bseleng/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func HadleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Bogdan",
		LastName:  "Seleng",
	}
	return c.JSON(u)
}
func HadleGetUser(c *fiber.Ctx) error {
	return c.JSON("James")
}
