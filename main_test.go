package main

import (
	"encoding/json"
	"github.com/dan-kirberger/djerk-djym-api/model"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
)

var app App
var ts *httptest.Server

func TestMain(m *testing.M) {
	app = App{}
	app.Initialize()
	ts = httptest.NewServer(app.Handler)
	log.Println("Test server running at " + ts.URL)
	defer ts.Close()

	code := m.Run()

	os.Exit(code)
}

func TestGetAllUsersReturnsEmptyList(t *testing.T) {
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
