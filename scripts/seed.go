package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bseleng/hotel-reservation/api"
	"github.com/bseleng/hotel-reservation/db"
	"github.com/bseleng/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	db := db.Store{
		Room:    db.NewMongoRoomStore(client, hotelStore),
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Hotel:   hotelStore,
	}

	bogdan := fixtures.AddUser(&db, "bogdan", "seleng", false)
	fmt.Printf("---%s--- \n\n%v\n\n", bogdan.FirstName, api.CreateTokenFromUser(bogdan))
	admin := fixtures.AddUser(&db, "admin", "admin", true)
	fmt.Printf("---%s--- \n\n%v\n\n", admin.FirstName, api.CreateTokenFromUser(admin))

	hiltonRu := fixtures.AddHotel(&db, "Hilton", "Russia", 3, nil)
	fixtures.AddRoom(&db, "lagre", true, 199.00, hiltonRu.ID)
	fixtures.AddRoom(&db, "small", true, 99.00, hiltonRu.ID)
	room := fixtures.AddRoom(&db, "medium", true, 149.00, hiltonRu.ID)
	fmt.Printf("--- hilton ru %v\n", hiltonRu)

	booking := fixtures.AddBooking(&db, bogdan.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Printf("--- booking -> %+v\n", booking)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		locaton := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(&db, name, locaton, rand.Intn(5)+1, nil)
	}

}
