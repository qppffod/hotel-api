package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/types"
	"github.com/qppffod/hotel-api/utils"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return utils.ErrUnAuthorized()
	}
	if !user.IsAdmin {
		return utils.ErrUnAuthorized()
	}

	return c.Next()
}
