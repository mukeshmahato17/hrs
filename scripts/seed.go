package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mukeshmahato17/hrs/api"
	"github.com/mukeshmahato17/hrs/db"
	"github.com/mukeshmahato17/hrs/db/fixture"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
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
