package db

import (
	"context"

	"github.com/qppffod/hotel-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRoomsByHotelID(context.Context, string) ([]*types.Room, error)
	GetRooms(context.Context) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore, dbname string) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbname).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) GetRooms(ctx context.Context) ([]*types.Room, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) GetRoomsByHotelID(ctx context.Context, id string) ([]*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"hotelID": oid}
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = res.InsertedID.(primitive.ObjectID)

	if err := s.HotelStore.UpdateHotelRooms(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}
