package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/phgh1246/golang_project01/db"
	"github.com/phgh1246/golang_project01/types"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TilDate    time.Time `json:"tilDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) Validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.FromDate) {
		return fmt.Errorf("cannot book a past date")
	}

	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), db.DataMap{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return err
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "interal server error",
		})
	}

	ok, err = h.isRoomAvailableForBooking(c.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked", c.Params("id")),
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TilDate:    params.TilDate,
		NumPersons: params.NumPersons,
	}

	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(
	ctx context.Context,
	roomID primitive.ObjectID,
	params BookRoomParams,
) (bool, error) {
	filter := db.DataMap{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tilDate": bson.M{
			"$lte": params.TilDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, filter)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}
