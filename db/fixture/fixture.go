package fixture

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mukeshmahato17/hrs/db"
	"github.com/mukeshmahato17/hrs/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		NumPersons: 3,
		FromDate:   from,
		TillDate:   till,
	}
	_, err := store.Booking.InsertBooking(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("booking", booking.ID)
	return booking
}

func AddRoom(store *db.Store, size string, price float64, hid primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Price:   price,
		HotelID: hid,
	}

	instertedRoom, err := store.Room.InsertRoom(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}
	return instertedRoom
}

func AddHotel(store *db.Store, name, loc string, rooms []*primitive.ObjectID, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel

}

func AddUser(store db.Store, fname, lname string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = admin

	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s -> %s\n", user.Email, api.CreateUserToken(user))
	return insertedUser
}
