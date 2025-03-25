package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qppffod/hotel-api/db"
	"github.com/qppffod/hotel-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		NumPersons: 5,
		FromDate:   from,
		TillDate:   till,
	}

	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		HotelID: hotelID,
		Size:    size,
		Seaside: ss,
		Price:   price,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddHotel(store *db.Store, name, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin

	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}
