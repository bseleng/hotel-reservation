package middleware

import (
	"fmt"

	"github.com/bseleng/hotel-reservation/api"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {

	user, err := api.GetAuthUser(c)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return fmt.Errorf("not authorized")
	}

	return c.Next()

}
