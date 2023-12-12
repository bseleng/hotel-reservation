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

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

// ...

func main() {
	listenAddr := flag.String("listenAddr", ":8000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	var (
		userHandler  = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		hotelHandler = api.NewHotelHandler(hotelStore, roomStore)
		app          = fiber.New(config)
		apiv1        = app.Group("api/v1")
	)
	// user handles
	apiv1.Get("/user", userHandler.HadleGetUsers)
	apiv1.Get("/user/:id", userHandler.HadleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("user/:id", userHandler.HandlePutUser)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HanlderGetHotels)

	app.Listen(*listenAddr)
}
