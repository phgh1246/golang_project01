package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/phgh1246/golang_project01/types"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, DataMap, DataMap) error
	GetHotels(context.Context, DataMap, *Pagination) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	dbName := os.Getenv(MongoDBNameEnvVar)
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbName).Collection("hotels"),
	}
}

func (s *MongoHotelStore) GetHotelByID(
	ctx context.Context,
	id string,
) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var hotel types.Hotel
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (s *MongoHotelStore) GetHotels(
	ctx context.Context,
	filter DataMap,
	pgntn *Pagination,
) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip((pgntn.Page - 1) * pgntn.Limit)
	opts.SetLimit(pgntn.Limit)
	resp, err := s.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filter, update DataMap) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) InsertHotel(
	ctx context.Context,
	hotel *types.Hotel,
) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}
