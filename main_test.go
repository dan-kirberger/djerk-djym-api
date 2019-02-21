package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dan-kirberger/djerk-djym-api/model"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var app App
var testServer *httptest.Server

func TestMain(m *testing.M) {
	app = App{}
	mongoUri, exists := os.LookupEnv("MONGO_URI")
	if !exists {
		mongoUri = "mongodb://localhost:27017"
	}
	app.Initialize(mongoUri)
	testServer = httptest.NewServer(app.Handler)
	log.Println("Test server running at " + testServer.URL)
	defer testServer.Close()

	code := m.Run()

	os.Exit(code)
}

func purgeDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := app.MongoClient.Database("testing").Collection("UserProfiles").Drop(ctx)
	if err != nil {
		panic("Failed to purge database of test data")
	}
}

func insertUser(user model.User, t *testing.T) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := app.MongoClient.Database("testing").Collection("UserProfiles")
	res, _ := collection.InsertOne(ctx, bson.M{"firstName": user.FirstName, "lastName": user.LastName, "weight": user.Weight})
	newId := res.InsertedID.(primitive.ObjectID).Hex()

	resp, _ := http.Get(testServer.URL + "/api/users/" + newId)
	if resp.StatusCode != 200 {
		t.Errorf("Precreated user should exist, instead got response " + strconv.Itoa(resp.StatusCode))
	}

	return newId
}

func TestGetAllUsersWithoutResults(t *testing.T) {
	purgeDatabase()
	resp, _ := http.Get(testServer.URL + "/api/users")

	if resp.StatusCode != 200 {
		t.Errorf("Should be 200 yo, not " + strconv.Itoa(resp.StatusCode))
	}

	decoder := json.NewDecoder(resp.Body)
	var mappedResponse model.UserList
	err := decoder.Decode(&mappedResponse)
	if err != nil {
		t.Errorf("Failed to map response json " + err.Error())
	}

	expectedResponse := model.UserList{Users: []model.User{}}

	if !reflect.DeepEqual(mappedResponse, expectedResponse) {
		t.Errorf("Response received did not match expected")
	}
}

func TestGetAllUsersWithResults(t *testing.T) {
	purgeDatabase()
	newId := insertUser(model.User{FirstName: "Testy", LastName: "McGee", Weight: 123}, t)

	resp, err := http.Get(testServer.URL + "/api/users")
	decoder := json.NewDecoder(resp.Body)
	var userListResponse model.UserList
	err = decoder.Decode(&userListResponse)
	if err != nil {
		t.Fatalf("Failed to map response json " + err.Error())
	}
	if len(userListResponse.Users) != 1 {
		t.Fatalf("Expected one user in response list, found " + strconv.Itoa(len(userListResponse.Users)))
	}
	user := userListResponse.Users[0]
	if user.ID != newId ||
		user.Weight != 123 ||
		user.FirstName != "Testy" ||
		user.LastName != "McGee" {
		t.Fatalf("Returned user fields do not match created user")
	}
}

func TestCreateUser(t *testing.T) {
	purgeDatabase()

	userToCreate := model.User{FirstName: "Testy", LastName: "McGee", Weight: 123}
	userCreateJson, _ := json.Marshal(userToCreate)

	resp, _ := http.Post(testServer.URL+"/api/users", "application/json", bytes.NewBuffer(userCreateJson))

	if resp.StatusCode != 200 {
		t.Errorf("Should receive 200 after create")
	}

	decoder := json.NewDecoder(resp.Body)
	var createdUser model.User
	err := decoder.Decode(&createdUser)
	if err != nil {
		t.Errorf("Failed to map json for created user " + err.Error())
	}
	if createdUser.ID == "" {
		t.Fatalf("User ID should be present on the response")
	}
	if createdUser.Weight != userToCreate.Weight ||
		createdUser.FirstName != userToCreate.FirstName ||
		createdUser.LastName != userToCreate.LastName {
		t.Fatalf("Returned user fields do not match input")
	}
}

func TestGetOneExistingUser(t *testing.T) {
	purgeDatabase()
	newId := insertUser(model.User{FirstName: "Testy", LastName: "McGee", Weight: 123}, t)

	resp, err := http.Get(testServer.URL + "/api/users/" + newId)
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200 response for existing user, instead got " + strconv.Itoa(resp.StatusCode))
	}
	decoder := json.NewDecoder(resp.Body)
	var fetchedUser model.User
	err = decoder.Decode(&fetchedUser)
	if err != nil {
		t.Errorf("Failed to map json for fetching user by id " + err.Error())
	}
	if fetchedUser.ID != newId {
		t.Fatalf("User ID should be present on the response")
	}
	if fetchedUser.ID != newId ||
		fetchedUser.Weight != 123 ||
		fetchedUser.FirstName != "Testy" ||
		fetchedUser.LastName != "McGee" {
		t.Fatalf("Returned user fields do not match created user")
	}
}

func TestGetOneMissingUser(t *testing.T) {
	resp, _ := http.Get(testServer.URL + "/api/users/i_dont_exist")
	if resp.StatusCode != 404 {
		t.Errorf("Expected 404 response for non existing user, instead got " + strconv.Itoa(resp.StatusCode))
	}
}

func TestDeleteUser(t *testing.T) {
	purgeDatabase()
	newId := insertUser(model.User{FirstName: "Testy", LastName: "McGee", Weight: 123}, t)

	req, _ := http.NewRequest("DELETE", testServer.URL+"/api/users/"+newId, nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("Should receive 200 from delete, got " + strconv.Itoa(resp.StatusCode))
	}
	resp, _ = http.Get(testServer.URL + "/api/users/" + newId)
	if resp.StatusCode != 404 {
		t.Errorf("Should receive 404 when fetching after delete, got " + strconv.Itoa(resp.StatusCode))
	}
}
