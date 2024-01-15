package api

import (
	"fmt"

	"github.com/bseleng/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})

	if err != nil {
		return ErrNotFound("bookings ")
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrNotFound(fmt.Sprintf("booking %s", id))
	}

	user, err := GetAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}

	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)

	if err != nil {
		return ErrNotFound(fmt.Sprintf("booking %s", id))
	}

	user, err := GetAuthUser(c)
	if err != nil {
		return ErrUnauthorized()

	}

	if booking.UserID != user.ID {
		return ErrUnauthorized()

	}

	if err := h.store.Booking.UpdateBooking(c.Context(), id, bson.M{"canceled": true}); err != nil {
		return err
	}

	return c.JSON(genericResp{Type: "msg", Msg: "updated"})

}
