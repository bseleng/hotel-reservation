package main

import (
	"flag"

	"github.com/bseleng/hotel-reservation/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddr := flag.String("listenAddr", ":8000", "The listen address of the API server")
	flag.Parse()
	app := fiber.New()
	apiv1 := app.Group("api/v1")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	apiv1.Get("/user", api.HadleGetUsers)
	apiv1.Get("/user", api.HadleGetUser)
	app.Listen(*listenAddr)
}
