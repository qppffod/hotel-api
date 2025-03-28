package db

import (
	"context"

	"github.com/qppffod/hotel-api/types"
	"github.com/qppffod/hotel-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetAvailableBookings(context.Context, *types.BookRoomParams, string) ([]*types.Booking, error)
	GetBookings(context.Context) ([]*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client, dbname string) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbname).Collection("booking"),
	}
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.ErrInvalidID()
	}

	update := bson.M{"$set": bson.M{"canceled": true}}

	_, err = s.coll.UpdateByID(ctx, oid, update)
	return err
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.ErrInvalidID()
	}

	var booking types.Booking
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var booking []*types.Booking
	if err := cur.All(ctx, &booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *MongoBookingStore) GetAvailableBookings(ctx context.Context, params *types.BookRoomParams, id string) ([]*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.ErrInvalidID()
	}

	filter := bson.M{
		"roomID": oid,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var booking []*types.Booking
	if err := cur.All(ctx, &booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = resp.InsertedID.(primitive.ObjectID)

	return booking, nil
}
