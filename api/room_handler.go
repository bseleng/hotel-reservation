package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bseleng/hotel-reservation/db"
	"github.com/bseleng/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}
type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(params); err != nil {
		return err
	}
	roomId := c.Params("id")
	roomOID, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return nil
	}

	user, ok := c.Context().Value("user").(*types.User)

	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "internal server error",
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomOID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}
	fmt.Println("%+v\n", booking)
	return nil
}
