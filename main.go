package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mukeshmahato17/hrs/api"
	"github.com/mukeshmahato17/hrs/types"
	"go.mongodb.org/mongo-driver/bson"
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
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	user := types.User{
		FirstName: "Alice",
		LastName:  "Gem",
	}
	coll := client.Database(dbname).Collection(users)
	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	var James types.User
	if err := coll.FindOne(ctx, bson.M{}).Decode(&James); err != nil {
		log.Fatal(err)
	}
	fmt.Println(James)

	listenAddr := flag.String("listenAddr", ":3000", "Listen Address")
	flag.Parse()
	app := fiber.New()

	appv1 := app.Group("api/v1")
	appv1.Get("/user", api.HandleUser)

	log.Fatal(app.Listen(*listenAddr))
}
