package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mukeshmahato17/hrs/db"
	"github.com/mukeshmahato17/hrs/types"
)

func InsertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "email@gmail.com",
		Password:  "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticationFailure(t *testing.T) {
	tbd := setup(t)
	InsertTestUser(t, tbd.UserStore)
	defer tbd.teardown(t)

	app := fiber.New()
	authHandler := NewHandleAuthStore(tbd.UserStore)
	app.Post("/auth", authHandler.HandleAuthentication)

	params := AuthParams{
		Email:    "email@gmail.com",
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
		t.Errorf("expected status 400 but got %d,", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected generic response type to be error but got %s", genResp.Type)
	}

	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected generic res msg to be <invalid credentials> but got %s", genResp.Msg)
	}
}

func TestAuthenticationSucess(t *testing.T) {
	tbd := setup(t)
	insertedUser := InsertTestUser(t, tbd.UserStore)
	defer tbd.teardown(t)

	app := fiber.New()
	authHandler := NewHandleAuthStore(tbd.UserStore)
	app.Post("/auth", authHandler.HandleAuthentication)

	params := AuthParams{
		Email:    "email@gmail.com",
		Password: "password",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		bb, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status 200 but got %d, Response %s,", resp.StatusCode, bb)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatal("Expected Token to be present in the auth response")
	}

	// Set the encrypted password to be nil, because we do not return any
	// JSON response.
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected user to be the inserted user")
	}
}
