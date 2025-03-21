package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address of the HTTP server")
	flag.Parse()
	app := fiber.New()

	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user/:id", api.HandleGetUser)
	apiv1.Get("/users", api.HandleGetUsers)

	app.Listen(*listenAddr)
}
