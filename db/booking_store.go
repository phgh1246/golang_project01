package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/phgh1246/golang_project01/types"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, DataMap) ([]*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, DataMap) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	dbName := os.Getenv(MongoDBNameEnvVar)
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbName).Collection("bookings"),
	}
}

func (s *MongoBookingStore) InsertBooking(
	ctx context.Context,
	booking *types.Booking,
) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)

	return booking, nil
}

func (s *MongoBookingStore) GetBookings(
	ctx context.Context,
	filter DataMap,
) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) GetBookingByID(
	ctx context.Context,
	id string,
) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update DataMap) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	m := bson.M{"$set": update}
	_, err = s.coll.UpdateByID(ctx, oid, m)
	return err
}
