package main

import (
	"context"
	"flag"
	"log"

	"github.com/bseleng/hotel-reservation/api"
	"github.com/bseleng/hotel-reservation/api/middleware"
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
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		hotelHandler   = api.NewHotelHandler(store)
		userHandler    = api.NewUserHandler(store.User)
		authHandler    = api.NewAuthHandler(store.User)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", middleware.AdminAuth)
	)

	//auth handles
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// versioned api routes
	// user handles
	apiv1.Get("/user", userHandler.HadleGetUsers)
	apiv1.Get("/user/:id", userHandler.HadleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HanlderGetHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)

	// room handlers
	apiv1.Get("room/", roomHandler.HandleGetRooms)
	apiv1.Post("room/:id/book", roomHandler.HandleBookRoom)

	//booking handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)

	//admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	app.Listen(*listenAddr)
}
