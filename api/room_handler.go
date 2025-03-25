package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/db"
	"github.com/qppffod/hotel-api/types"
	"github.com/qppffod/hotel-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params types.BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.ErrInvalidID()
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "internal server error",
		})
	}

	ok, err = h.isRoomAvailableForBooking(c.Context(), &params, c.Params("id"))
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked", roomID.Hex()),
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		NumPersons: params.NumPersons,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
	}

	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(insertedBooking)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, params *types.BookRoomParams, id string) (bool, error) {
	bookings, err := h.store.Booking.GetAvailableBookings(ctx, params, id)
	if err != nil {
		return false, utils.ErrResourceNotFound("available bookings")
	}

	ok := len(bookings) == 0
	return ok, nil
}
