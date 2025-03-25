package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/api/middleware"
	"github.com/qppffod/hotel-api/db/fixtures"
	"github.com/qppffod/hotel-api/types"
)

func TestUserBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		nonauthuser    = fixtures.AddUser(db.Store, "John", "Rand", false)
		user           = fixtures.AddUser(db.Store, "first", "second", false)
		hotel          = fixtures.AddHotel(db.Store, "bar hotel", "qaqaq", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 99.99, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(db.Store.User))
		bookingHandler = NewBookingHandler(db.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response: %d", resp.StatusCode)
	}

	var bookingResp types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	if bookingResp.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID.Hex(), bookingResp.ID.Hex())
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID.Hex(), bookingResp.UserID.Hex())
	}

	// non auth user test
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonauthuser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code response but got: %d", resp.StatusCode)
	}
}

func TestGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		adminUser      = fixtures.AddUser(db.Store, "admin", "admin", true)
		user           = fixtures.AddUser(db.Store, "first", "second", false)
		hotel          = fixtures.AddHotel(db.Store, "bar hotel", "qaqaq", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 99.99, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.Store.User), middleware.AdminAuth)
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
		t.Fatalf("non 200 response: %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}

	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID.Hex(), have.ID.Hex())
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID.Hex(), have.UserID.Hex())
	}

	// test for non-admin (user) access
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code response but got: %d", resp.StatusCode)
	}
	fmt.Println(bookings)
}
