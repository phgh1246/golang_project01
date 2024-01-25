package fixtures

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/phgh1246/golang_project01/db"
	"github.com/phgh1246/golang_project01/types"
)

func AddUser(store *db.Store, firstName, lastName string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", strings.ToLower(firstName), strings.ToLower(lastName)),
		FirstName: firstName,
		LastName:  lastName,
		Password: fmt.Sprintf(
			"%s_%spass123",
			strings.ToLower(firstName),
			strings.ToLower(lastName),
		),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddHotel(
	store *db.Store,
	name string,
	location string,
	rating int,
	roomIDs []primitive.ObjectID,
) *types.Hotel {
	rooms := roomIDs
	if roomIDs == nil {
		rooms = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
		Rooms:    rooms,
	}
	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddRoom(
	store *db.Store,
	size string,
	seaside bool,
	price float64,
	hotelID primitive.ObjectID,
) *types.Room {
	room := types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.InsertRoom(context.TODO(), &room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(
	store *db.Store,
	userID, roomID primitive.ObjectID,
	fromDate, tilDate time.Time,
) *types.Booking {
	booking := types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: fromDate,
		TilDate:  tilDate,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.TODO(), &booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
