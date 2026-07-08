package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mukeshmahato17/hrs/api"
	"github.com/mukeshmahato17/hrs/db"
	"github.com/mukeshmahato17/hrs/db/fixture"
	"github.com/mukeshmahato17/hrs/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	userStore    *db.MongoDBStore
	hotelStore   *db.MongoHotelStore
	roomStore    *db.MongoRoomStore
	bookingStore *db.MongoBookStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email string, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin

	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s -> %s\n", user.Email, api.CreateUserToken(user))
	return insertedUser
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		NumPersons: 3,
		FromDate:   from,
		TillDate:   till,
	}
	_, err := bookingStore.InsertBooking(ctx, booking)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("booking", resp.ID)
}

func seedRoom(size string, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Price:   price,
		HotelID: hotelID,
	}

	instertedRoom, err := roomStore.InsertRoom(ctx, room)
	if err != nil {
		log.Fatal(err)
	}
	return instertedRoom
}

func seedHotels(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func main() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	store := &db.Store{
		User:    db.NewMongoDBStore(client),
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: db.NewMongoBookStore(client),
	}

	foo := fixture.AddUser(*store, "foo", "bar", false)
	fmt.Println("foo ->", api.CreateUserToken(foo))
	admin := fixture.AddUser(*store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateUserToken(admin))
	hotel := fixture.AddHotel(store, "HotelAnything", "Anywhere", nil, 5)
	room := fixture.AddRoom(store, "Large", 59.66, hotel.ID)
	booking := fixture.AddBooking(store, foo.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println(booking.ID)
}
