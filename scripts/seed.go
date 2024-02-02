package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/phgh1246/golang_project01/api"
	"github.com/phgh1246/golang_project01/db"
	"github.com/phgh1246/golang_project01/db/fixtures"
)

var store *db.Store

func addRoomsToHotel(hotelID primitive.ObjectID, userIDs []primitive.ObjectID) {
	room := fixtures.AddRoom(store, "small", true, 99.99, hotelID)

	for i, userID := range userIDs {
		booking := fixtures.AddBooking(
			store,
			userID,
			room.ID,
			time.Now().AddDate(0, 0, 4*i),
			time.Now().AddDate(0, 0, 3+(4*i)),
		)
		fmt.Println("booking: ", booking.ID)
	}

	room = fixtures.AddRoom(store, "medium", false, 199.99, hotelID)
	room = fixtures.AddRoom(store, "large", true, 299.99, hotelID)
}

func main() {
	var (
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store = &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	userIDs := []primitive.ObjectID{}

	user := fixtures.AddUser(store, "Alice", "Names", true)
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	userIDs = append(userIDs, user.ID)

	user = fixtures.AddUser(store, "Bob", "Named", false)
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	userIDs = append(userIDs, user.ID)

	user = fixtures.AddUser(store, "Charlie", "Nameless", false)
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	userIDs = append(userIDs, user.ID)

	hotel := fixtures.AddHotel(store, "Fancy Hotel", "Down The Street", 4, nil)
	addRoomsToHotel(hotel.ID, userIDs)

	hotel = fixtures.AddHotel(store, "Classy Hotel", "Down The Lane", 5, nil)
	addRoomsToHotel(hotel.ID, userIDs)

	hotel = fixtures.AddHotel(store, "Shoddy Hotel", "Down The Alley", 2, nil)
	addRoomsToHotel(hotel.ID, userIDs)

	for i := 0; i < 100; i++ {
		fixtures.AddHotel(
			store,
			fmt.Sprintf("Cloned Hotel %d", i),
			fmt.Sprintf("In Parallel Universe %d", i),
			rand.Intn(5)+1,
			nil,
		)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
