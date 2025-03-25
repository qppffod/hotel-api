package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/hotel-api/db/fixtures"
	"github.com/qppffod/hotel-api/types"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	// insertedUser := insertTestUser(t, tdb.Store.User)
	insertedUser := fixtures.AddUser(tdb.Store, "first", "second", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.Store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := types.AuthParams{
		Email:    "first@second.com",
		Password: "first_second",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var authResp types.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	insertedUser.EncryptedPassword = ""

	if authResp.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user to be the inserted user")
	}

	fmt.Println(authResp)
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	fixtures.AddUser(tdb.Store, "first", "second", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.Store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := types.AuthParams{
		Email:    "first@second.com",
		Password: "wrongpassword",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s", genResp.Type)
	}
	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected gen response msg to be <invalid credentials> but got %s", genResp.Msg)
	}

}
