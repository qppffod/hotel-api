package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/api"
	"github.com/qppffod/hotel-api/db"
	"github.com/qppffod/hotel-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if apiError, ok := err.(utils.Error); ok {
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := utils.NewError(http.StatusInternalServerError, err.Error())
		return c.Status(apiError.Code).JSON(apiError)
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address of the HTTP server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		// stores
		userStore    = db.NewMongoUserStore(client, db.DBNAME)
		hotelStore   = db.NewMongoHotelStore(client, db.DBNAME)
		roomStore    = db.NewMongoRoomStore(client, hotelStore, db.DBNAME)
		bookingStore = db.NewMongoBookingStore(client, db.DBNAME)
		store        = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}

		// handlers
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiv1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin = apiv1.Group("/admin", api.AdminAuth)
	)

	// auth handlers
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Versioned API routes
	// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/room", hotelHandler.HandleGetHotelRooms)

	// room handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// boking handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin only
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	app.Listen(*listenAddr)
}

// "fromDate": "2025-04-01T00:00:00.0Z",
// "tillDate": "2025-04-05T00:00:00.0Z",
// "numPersons": 5
