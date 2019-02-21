package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dan-kirberger/djerk-djym-api/model"
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
var ts *httptest.Server

func TestMain(m *testing.M) {
	app = App{}
	mongoUri, exists := os.LookupEnv("MONGO_URI")
	if !exists {
		mongoUri = "mongodb://localhost:27017"
	}
	app.Initialize(mongoUri)
	ts = httptest.NewServer(app.Handler)
	log.Println("Test server running at " + ts.URL)
	defer ts.Close()

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

func TestGetAllUsersReturnsEmptyList(t *testing.T) {
	purgeDatabase()
	resp, err := http.Get(ts.URL + "/api/users")
	if err != nil {
		t.Errorf("Failed to fetch users")
	}

	if resp.StatusCode != 200 {
		t.Errorf("Should be 200 yo, not " + strconv.Itoa(resp.StatusCode))
	}
	//responseJson, err := simplejson.NewFromReader(resp.Body)
	//content := responseJson.GetPath("content").MustArray()
	//if len(content) > 0 {
	//	t.Errorf("Expected empty content")
	//}
	decoder := json.NewDecoder(resp.Body)
	var mappedResponse model.UserList
	err = decoder.Decode(&mappedResponse)
	if err != nil {
		t.Errorf("Failed to map response json " + err.Error())
	}

	expectedResponse := model.UserList{Users: []model.User{}}

	if !reflect.DeepEqual(mappedResponse, expectedResponse) {
		t.Errorf("Response received did not match expected")
	}
}

func TestSingleUserCrud(t *testing.T) {
	purgeDatabase()

	userToCreate := model.User{FirstName: "Testy", LastName: "McGee", Weight: 123}
	userCreateJson, _ := json.Marshal(userToCreate)

	resp, _ := http.Post(ts.URL+"/api/users", "application/json", bytes.NewBuffer(userCreateJson))

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

	resp, err = http.Get(ts.URL + "/api/users/" + createdUser.ID)
	if err != nil {
		t.Errorf("Failed to fetch user newly created user")
	}
	decoder = json.NewDecoder(resp.Body)
	var fetchedUser model.User
	err = decoder.Decode(&fetchedUser)
	if err != nil {
		t.Errorf("Failed to map json for fetching user by id " + err.Error())
	}
	if fetchedUser.ID != createdUser.ID {
		t.Fatalf("User ID should be present on the response")
	}
	if createdUser.Weight != userToCreate.Weight ||
		createdUser.FirstName != userToCreate.FirstName ||
		createdUser.LastName != userToCreate.LastName {
		t.Fatalf("Returned user fields do not match input")
	}
	resp, err = http.Get(ts.URL + "/api/users")
	decoder = json.NewDecoder(resp.Body)
	var userListResponse model.UserList
	err = decoder.Decode(&userListResponse)
	if err != nil {
		t.Fatalf("Failed to map response json " + err.Error())
	}
	if len(userListResponse.Users) != 1 {
		t.Fatalf("Expected one user in response list, found " + strconv.Itoa(len(userListResponse.Users)))
	}

	req, _ := http.NewRequest("DELETE", ts.URL+"/api/users/"+createdUser.ID, nil)
	resp, err = http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Errorf("Should receive 200 from delete, got " + strconv.Itoa(resp.StatusCode))
	}
	resp, err = http.Get(ts.URL + "/api/users/" + createdUser.ID)
	if resp.StatusCode != 404 {
		t.Errorf("Should receive 404 when fetching after delete, got " + strconv.Itoa(resp.StatusCode))
	}
}
