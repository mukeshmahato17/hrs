package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mukeshmahato17/hrs/api"
	"github.com/mukeshmahato17/hrs/api/middleware"
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
		userStore    = db.NewMongoDBStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingstore = db.NewMongoBookStore(client)
		store        = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingstore,
		}
		bookingHandler = api.NewBookingHandler(*store)
		roomHandler    = api.NewHandleRoomStore(store)
		handleHotel    = api.NewHandleHotelStore(store)
		userHandler    = api.NewHandleUserStore(db.NewMongoDBStore(client))
		authHandler    = api.NewHandleAuthStore(userStore)

		app   = fiber.New()
		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		auth  = app.Group("/api")
		admin = apiv1.Group("/admin", middleware.AdminAuth)
	)

	// auth Handler
	auth.Post("/auth", authHandler.HandleAuthentication)

	// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user/:id", userHandler.HandleUser)
	apiv1.Get("/users", userHandler.HandleUsers)
	apiv1.Post("/user", userHandler.HandleUserPost)

	// hotel handlers
	apiv1.Get("/hotel", handleHotel.HandleGetHotels)
	apiv1.Get("/hotel/:id/rooms", handleHotel.HandleGetRooms)
	apiv1.Get("/hotel/:id", handleHotel.HandleGetHotelByID)

	// room handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// booking handlers
	apiv1.Post("/booking/:id/cancel", bookingHandler.HandleCancelBooking)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)

	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	log.Fatal(app.Listen(*listenAddr))
}
