package api

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/db"
	"github.com/qppffod/hotel-api/types"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams types.AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(authParams.Password))
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}

	fmt.Println("authenticated--->", user)

	return c.JSON(user)
}
