package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/phgh1246/golang_project01/db"
	"github.com/phgh1246/golang_project01/types"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME)
	hotel := types.Hotel{
		Name:     "The Fancy Hotel",
		Location: "Down the Street",
	}

	rooms := []types.Room{
		{
			Type:      types.SingeRoomType,
			BasePrice: 99.99,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 199.99,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 399.99,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 799.99,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(insertedHotel)

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}
}
