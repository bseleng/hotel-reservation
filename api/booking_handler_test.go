package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/bseleng/hotel-reservation/db/fixtures"
)

func TestGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	user := fixtures.AddUser(db.Store, "bogdan", "seleng", true)
	hotelSpb := fixtures.AddHotel(db.Store, "Plaza", "Saint Petersburg", 3, nil)
	room := fixtures.AddRoom(db.Store, "small", false, 29.90, hotelSpb.ID)

	from := time.Now()
	till := time.Now().AddDate(0, 0, 3)
	booking := fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)

	fmt.Println(booking)
}
