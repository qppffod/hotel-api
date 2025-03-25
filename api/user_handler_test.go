package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/types"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.Store.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "James",
		LastName:  "Smith",
		Email:     "jamessmith@gmail.com",
		Password:  "ll12l3lflefl",
	}
	b, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Error(err)
	}

	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EcnryptedPassword not to be included in the JSON response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName to be %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastName to be %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email to be %s but got %s", params.Email, user.Email)
	}
}
