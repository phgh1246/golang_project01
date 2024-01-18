package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/phgh1246/golang_project01/db"
	"github.com/phgh1246/golang_project01/types"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(firstName, lastName, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  "defaultseedpassword",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
		Rooms:    []primitive.ObjectID{},
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.99,
		},
		{
			Size:  "normal",
			Price: 199.99,
		},
		{
			Size:  "large",
			Price: 399.99,
		},
		{
			Size:  "jumbo",
			Price: 799.99,
		},
	}
	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("Fancy Hotel", "Down The Street", 4)
	seedHotel("Classy Hotel", "Down The Lane", 5)
	seedHotel("Shoddy Hotel", "Down The Alley", 2)

	seedUser("Alice", "Names", "alice@emails.com")
	seedUser("Bob", "Named", "bob@emails.co")
	seedUser("Charlie", "Nameless", "charlie@emails.org")
}

func init() {
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
