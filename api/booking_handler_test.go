package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/phgh1246/golang_project01/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	user := fixtures.AddUser(db.Store, "Alice", "Names", false)
	hotel := fixtures.AddHotel(db.Store, "Fancy Hotel", "Down The Street", 4, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 99.99, hotel.ID)
	booking := fixtures.AddBooking(
		db.Store,
		user.ID,
		room.ID,
		time.Now().AddDate(0, 0, 1),
		time.Now().AddDate(0, 0, 4),
	)
	fmt.Println(booking)
}
