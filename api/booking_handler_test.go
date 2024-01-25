package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/phgh1246/golang_project01/api/middleware"
	"github.com/phgh1246/golang_project01/db/fixtures"
	"github.com/phgh1246/golang_project01/types"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		user        = fixtures.AddUser(db.Store, "Bob", "Named", false)
		nonAuthUser = fixtures.AddUser(db.Store, "Charlie", "Nameless", false)
		hotel       = fixtures.AddHotel(db.Store, "Fancy Hotel", "Down The Street", 4, nil)
		room        = fixtures.AddRoom(db.Store, "small", true, 99.99, hotel.ID)
		booking     = fixtures.AddBooking(
			db.Store,
			user.ID,
			room.ID,
			time.Now().AddDate(0, 0, 1),
			time.Now().AddDate(0, 0, 4),
		)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected response code 200, got: %d", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(bookingResp)
	if bookingResp.ID != booking.ID {
		t.Fatalf("expected ID %s, got ID %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected UserID %s, got UserID %s", booking.UserID, bookingResp.UserID)
	}

	// testing user trying to access another user's booking
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatal("got status code 200 when expecting failure")
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		adminUser = fixtures.AddUser(db.Store, "Alice", "Names", true)
		user      = fixtures.AddUser(db.Store, "Bob", "Named", false)
		hotel     = fixtures.AddHotel(db.Store, "Fancy Hotel", "Down The Street", 4, nil)
		room      = fixtures.AddRoom(db.Store, "small", true, 99.99, hotel.ID)
		booking   = fixtures.AddBooking(
			db.Store,
			user.ID,
			room.ID,
			time.Now().AddDate(0, 0, 1),
			time.Now().AddDate(0, 0, 4),
		)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)
	_ = booking

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got: %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err = json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}
	retrievedBooking := bookings[0]
	if retrievedBooking.ID != booking.ID {
		t.Fatalf("expected ID %s, got ID %s", booking.ID, retrievedBooking.ID)
	}
	if retrievedBooking.UserID != booking.UserID {
		t.Fatalf("expected UserID %s, got UserID %s", booking.UserID, retrievedBooking.UserID)
	}

	// test for non-admin auth rejection
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatal("got status code 200 when expecting failure")
	}
}
