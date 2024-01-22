package api

import (
	"github.com/bseleng/hotel-reservation/db"
	"github.com/bseleng/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

type ResourceResp struct {
	Data    []*types.Hotel `json:"data"`
	Results int            `json:"results"`
	Page    int            `json:"page"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func NewResurceResp(data []*types.Hotel, page int) ResourceResp {
	return ResourceResp{
		Results: len(data),
		Data:    data,
		Page:    page,
	}
}
func (h *HotelHandler) HanlderGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	filter := db.Map{
		"rating": params.Rating,
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrNotFound("hotels")
	}
	return c.JSON(NewResurceResp(hotels, int(params.Page)))
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)

	if err != nil {
		return ErrNotFound("hotel")
	}

	return c.JSON(hotel)

}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return ErrInvadidID()
	}

	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err

	}
	return c.JSON(rooms)
}
