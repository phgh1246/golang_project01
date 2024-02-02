package api

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/phgh1246/golang_project01/db"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}
	dbName := os.Getenv(db.MongoDBNameEnvVar)
	if err := tdb.client.Database(dbName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}
	dbURI := os.Getenv("TEST_MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	tdb := &testdb{
		client: client,
		Store: &db.Store{
			Hotel:   hotelStore,
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
		},
	}

	tdb.teardown(t)

	return tdb
}
