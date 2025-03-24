package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/db"
	"github.com/qppffod/hotel-api/types"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// admin only
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// user only
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "unauthorized",
		})
	}
	return c.JSON(booking)
}
