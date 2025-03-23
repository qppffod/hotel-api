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

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address of the HTTP server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// stores
	userStore := db.NewMongoUserStore(client, db.DBNAME)

	// handlers
	userHandler := api.NewUserHandler(userStore)

	var (
		app   = fiber.New(config)
		apiv1 = app.Group("/api/v1")
	)

	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)

	app.Listen(*listenAddr)
}
