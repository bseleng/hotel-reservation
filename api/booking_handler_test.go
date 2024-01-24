package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin          = app.Group("/", JWTAuthentication(db.User), AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

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

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}

	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %q got %q", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %q got %q", booking.UserID, have.UserID)
	}

	// test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized but got %d", resp.StatusCode)
	}

	fmt.Println(bookings)
}

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		user           = fixtures.AddUser(db.Store, "bogdan", "seleng", false)
		nonAuthUser    = fixtures.AddUser(db.Store, "bogdan", "___", false)
		hotelSpb       = fixtures.AddHotel(db.Store, "Plaza", "Saint Petersburg", 3, nil)
		room           = fixtures.AddRoom(db.Store, "small", false, 29.90, hotelSpb.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 3)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		api            = app.Group("/", JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)

	fmt.Printf("\n\n %q %[1]T\n", booking.ID.Hex())

	api.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 got %d", resp.StatusCode)
	}

	var bookingResp *types.Booking

	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	if bookingResp.ID != booking.ID {
		t.Fatalf("expectes %s got %s", bookingResp.ID, booking.ID)
	}

	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expectes %s got %s", bookingResp.UserID, booking.UserID)
	}

	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Errorf("expected a non 200 got %d", resp.StatusCode)
	}
}
