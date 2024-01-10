package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bseleng/hotel-reservation/db"
	"github.com/bseleng/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fname, lname string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		FirstName: fname,
		LastName:  lname,
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(db *db.Store, name, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDs = rooms
	if rooms == nil {
		roomIDs = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIDs,
		Rating:   rating,
	}

	insertedHotel, err := db.Hotel.Insert(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(db *db.Store, size string, seaSide bool, price float64, hotelId primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaSide,
		Price:   price,
		HotelID: hotelId,
	}

	insertedRoom, err := db.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(db *db.Store, userID, roomID primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}

	insertedBooking, err := db.Booking.InsertBooking(context.Background(), &booking)

	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}
