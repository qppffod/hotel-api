package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/qppffod/hotel-api/db/fixtures"
)

func TestGetBookingt(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	user := fixtures.AddUser(db.Store, "first", "second", false)
	hotel := fixtures.AddHotel(db.Store, "bar hotel", "qaqaq", 4, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 99.99, hotel.ID)

	from := time.Now()
	till := from.AddDate(0, 0, 5)
	booking := fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
	fmt.Println(booking)
}
