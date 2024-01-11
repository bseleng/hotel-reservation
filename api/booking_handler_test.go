package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bseleng/hotel-reservation/api/middleware"
	"github.com/bseleng/hotel-reservation/db/fixtures"
	"github.com/bseleng/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		adminUser      = fixtures.AddUser(db.Store, "admin", "admin", true)
		user           = fixtures.AddUser(db.Store, "bogdan", "seleng", false)
		hotelSpb       = fixtures.AddHotel(db.Store, "Plaza", "Saint Petersburg", 3, nil)
		room           = fixtures.AddRoom(db.Store, "small", false, 29.90, hotelSpb.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

	_ = booking
	fmt.Printf("\n\n adminUser --> %+v\n", adminUser)

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("none 200 response %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	fmt.Println(bookings)
}
