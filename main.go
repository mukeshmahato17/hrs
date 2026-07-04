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

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const users = "users"

func main() {
	// Creates a new client and connects to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := flag.String("listenAddr", ":3000", "Listen Address")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("api/v1")

	handleUser := api.NewHandleUserStore(db.NewMongoDBStore(client))
	apiv1.Get("/user/:id", handleUser.HandleUser)
	apiv1.Get("/users", handleUser.HandleUsers)
	apiv1.Post("/user", handleUser.HandleUserPost)

	log.Fatal(app.Listen(*listenAddr))
}
