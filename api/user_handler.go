package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/types"
)

type UserHandler struct {
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": c.Params("id")})
}

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		ID:        "1",
		FirstName: "James",
		LastName:  "Jackson",
	}
	return c.JSON(u)
}
