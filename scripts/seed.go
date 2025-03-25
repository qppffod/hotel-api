package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qppffod/hotel-api/api"
	"github.com/qppffod/hotel-api/db"
	"github.com/qppffod/hotel-api/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)

	store := db.Store{
		Hotel:   db.NewMongoHotelStore(client, db.DBNAME),
		Room:    db.NewMongoRoomStore(client, hotelStore, db.DBNAME),
		User:    db.NewMongoUserStore(client, db.DBNAME),
		Booking: db.NewMongoBookingStore(client, db.DBNAME),
	}

	user := fixtures.AddUser(&store, "first", "second", false)
	fmt.Println("user ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(&store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(&store, "Bellucia", "France", 3, nil)
	room := fixtures.AddRoom(&store, "small", true, 89.99, hotel.ID)
	booking := fixtures.AddBooking(&store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 3))
	fmt.Println("booking ->", booking.ID)
}
