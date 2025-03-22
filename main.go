package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/api"
	"github.com/qppffod/hotel-api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi  = "mongodb://localhost:27017"
	dbanme = "hotel-reservation"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address of the HTTP server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	// stores
	userStore := db.NewMongoUserStore(client)

	// handlers
	userHandler := api.NewUserHandler(userStore)

	var (
		app   = fiber.New(config)
		apiv1 = app.Group("/api/v1")
	)

	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)

	app.Listen(*listenAddr)
}
