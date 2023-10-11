package main

import (
	"context"
	"flag"
	"log"

	"github.com/bseleng/hotel-reservation/api"
	"github.com/bseleng/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

// ...

func main() {
	listenAddr := flag.String("listenAddr", ":8000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(config)
	apiv1 := app.Group("api/v1")

	// handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	apiv1.Get("/user", userHandler.HadleGetUsers)
	apiv1.Get("/user/:id", userHandler.HadleGetUser)
	app.Listen(*listenAddr)
}
