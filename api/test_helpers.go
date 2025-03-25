package api

import (
	"context"
	"log"
	"testing"

	"github.com/qppffod/hotel-api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	Store  *db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.TESTDBNAME).Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.TESTDBNAME)

	return &testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client, db.TESTDBNAME),
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, hotelStore, db.TESTDBNAME),
			Booking: db.NewMongoBookingStore(client, db.TESTDBNAME),
		},
	}
}
