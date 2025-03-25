package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/db"
	"github.com/qppffod/hotel-api/utils"
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
		return utils.ErrResourceNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return utils.ErrUnAuthorized()
	}

	if booking.UserID != user.ID {
		return utils.ErrUnAuthorized()
	}

	err = h.store.Booking.UpdateBooking(c.Context(), booking.ID.Hex())
	if err != nil {
		return err
	}

	return c.JSON(genericResp{
		Type: "msg",
		Msg:  "updated",
	})
}

// admin only
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context())
	if err != nil {
		return utils.ErrResourceNotFound("bookings")
	}

	return c.JSON(bookings)
}

// user only
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return utils.ErrResourceNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return utils.ErrUnAuthorized()
	}
	if booking.UserID != user.ID {
		return utils.ErrUnAuthorized()
	}
	return c.JSON(booking)
}
