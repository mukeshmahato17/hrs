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
		app         = fiber.New()
		apiv1       = app.Group("api/v1")
		handleUser  = api.NewHandleUserStore(db.NewMongoDBStore(client, db.DBNAME))
		hotelStore  = db.NewMongoHotelStore(client)
		roomStore   = db.NewMongoRoomStore(client, hotelStore)
		handleHotel = api.NewHandleHotelStore(hotelStore, roomStore)
	)
	// user handlers
	apiv1.Put("/user/:id", handleUser.HandlePutUser)
	apiv1.Delete("/user/:id", handleUser.HandleDeleteUser)
	apiv1.Get("/user/:id", handleUser.HandleUser)
	apiv1.Get("/users", handleUser.HandleUsers)
	apiv1.Post("/user", handleUser.HandleUserPost)

	// hotel handlers
	apiv1.Get("/hotels", handleHotel.HandleGetHotels)
	log.Fatal(app.Listen(*listenAddr))
}
