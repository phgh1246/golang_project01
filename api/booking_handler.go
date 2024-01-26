package api

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/phgh1246/golang_project01/db"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrorResourceNotFound("booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrorUnauthorized()
	}
	if booking.UserID != user.ID {
		return ErrorUnauthorized()
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResp{
		Type: "msg",
		Msg:  "updated",
	})
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return ErrorResourceNotFound("bookings")
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrorResourceNotFound("booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrorUnauthorized()
	}
	if booking.UserID != user.ID {
		return ErrorUnauthorized()
	}
	return c.JSON(booking)
}
