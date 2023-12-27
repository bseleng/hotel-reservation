package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bseleng/hotel-reservation/api"
	"github.com/bseleng/hotel-reservation/db"
	"github.com/bseleng/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel

}

func seedUser(isadmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isadmin
	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%s --> %s\n\n", user.Email, api.CreateTokenFromUser(user))

	return insertedUser
}

func seedRoom(size string, seaSide bool, price float64, hotelId primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaSide,
		Price:   price,
		HotelID: hotelId,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}

	if _, err := bookingStore.InsertBooking(context.Background(), &booking); err != nil {
		log.Fatal(err)
	}
}

func main() {

	bogdan := seedUser(false, "bogdan", "seleng", "bseleng@test.com", "superPassword")
	seedUser(true, "admin", "admin", "admin@test.com", "adminPassword")
	seedHotel("Hilton", "Russia", 2)
	seedHotel("Ritz Carlton", "Switzerland", 5)
	fsHotel := seedHotel("Four Seasons", "Italy", 3)
	seedRoom("small", true, 79.99, fsHotel.ID)
	seedRoom("medium", false, 99.99, fsHotel.ID)
	roomToBook := seedRoom("large", true, 179.99, fsHotel.ID)
	seedBooking(bogdan.ID, roomToBook.ID, time.Now(), time.Now().AddDate(0, 0, 2))

}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}
