package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mukeshmahato17/hrs/api"
	"github.com/mukeshmahato17/hrs/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Creates a new client and connects to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	listenAddr := flag.String("listenAddr", ":3000", "Listen Address")
	flag.Parse()

	var (
		app        = fiber.New()
		apiv1      = app.Group("api/v1")
		handleUser = api.NewHandleUserStore(db.NewMongoDBStore(client))
		userStore  = db.NewMongoDBStore(client)
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		store      = &db.Store{
			User:  userStore,
			Hotel: hotelStore,
			Room:  roomStore,
		}
		handleHotel = api.NewHandleHotelStore(store)
	)
	// user handlers
	apiv1.Put("/user/:id", handleUser.HandlePutUser)
	apiv1.Delete("/user/:id", handleUser.HandleDeleteUser)
	apiv1.Get("/user/:id", handleUser.HandleUser)
	apiv1.Get("/users", handleUser.HandleUsers)
	apiv1.Post("/user", handleUser.HandleUserPost)

	// hotel handlers
	apiv1.Get("/hotel", handleHotel.HandleGetHotels)
	apiv1.Get("/hotel/:id/rooms", handleHotel.HandleGetRooms)
	apiv1.Get("/hotel/:id", handleHotel.HandleGetHotelByID)
	log.Fatal(app.Listen(*listenAddr))
}
