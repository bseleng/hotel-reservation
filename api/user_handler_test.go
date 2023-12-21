package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/bseleng/hotel-reservation/db"
	"github.com/bseleng/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testdburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation-test"

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)
	params := types.CreateUserParams{
		Email:     "some@test.com",
		FirstName: "Foo",
		LastName:  "Bar",
		Password:  "1234567",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expected ID length is greater than zero")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected EncryptedPassword not to be included in the JSON response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName %s but got %s", params.FirstName, user.FirstName)
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected FirstName %s but got %s", params.FirstName, user.FirstName)
	}
	if user.Email != params.Email {
		t.Errorf("expected Email %s but got %s", params.Email, user.Email)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected LastName %s but got %s", params.LastName, user.LastName)
	}
	fmt.Println(user)
}
