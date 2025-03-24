package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qppffod/hotel-api/api"
	"github.com/qppffod/hotel-api/db"
	"github.com/qppffod/hotel-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin

	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel

}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		NumPersons: 5,
		FromDate:   from,
		TillDate:   till,
	}

	resp, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking:", resp.ID.Hex())
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		HotelID: hotelID,
		Size:    size,
		Seaside: ss,
		Price:   price,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func main() {
	james := seedUser(false, "first", "second", "firstsecond@gmail.com", "supersecurepassword")
	seedUser(true, "admin", "admin", "admin@admin.com", "adminpassword123")
	seedHotel("Bellucia", "France", 3)
	seedHotel("The cozy hotel", "Nederlands", 4)
	hotel := seedHotel("Some hotel", "England", 1)
	seedRoom("small", true, 89.99, hotel.ID)
	seedRoom("medium", true, 189.99, hotel.ID)
	room := seedRoom("large", true, 289.99, hotel.ID)
	seedBooking(james.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoRoomStore(client, hotelStore, db.DBNAME)
	userStore = db.NewMongoUserStore(client, db.DBNAME)
	bookingStore = db.NewMongoBookingStore(client, db.DBNAME)
}
