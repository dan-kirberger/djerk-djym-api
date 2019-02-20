package main

import (
	"context"
	"encoding/json"
	"github.com/dan-kirberger/djerk-djym-api/model"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"log"
	"net/http"
	"regexp"
	"time"
)

type App struct {
	Handler *http.ServeMux
}

func handler(w http.ResponseWriter, r *http.Request) {
	//log.Printf("Received request for %s", r.URL.Path[1:])
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	user := model.User{ID: "asdf", FirstName: "dan", LastName: "iel", Weight: 999}
	response, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func writeHandler(writer http.ResponseWriter, request *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, _ := mongo.Connect(ctx, &options.ClientOptions{Hosts: []string{"localhost"}})
	_ = client.Ping(ctx, readpref.Primary())

	collection := client.Database("testing").Collection("UserProfiles")
	res, _ := collection.InsertOne(ctx, bson.M{"firstName": "Dan", "lastName": "K"})

	user := model.User{ID: res.InsertedID.(primitive.ObjectID).Hex(), FirstName: "dan", LastName: "iel", Weight: 999}
	response, _ := json.Marshal(user)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(response)
}

func users(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		getAllUsers(writer, request)
	case http.MethodPost:
		addUser(writer, request)
	default:
		notFound(writer, "i dunno how to do that yet")
	}
}

func getAllUsers(writer http.ResponseWriter, request *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, _ := mongo.Connect(ctx, &options.ClientOptions{Hosts: []string{"localhost"}})
	//_ = client.Ping(ctx, readpref.Primary())

	mongoUsers, _ := client.Database("testing").Collection("UserProfiles").Find(ctx, bson.D{})

	userModels := make([]model.User, 0)

	for mongoUsers.Next(ctx) {
		mongoDoc := &bson.D{}
		_ = mongoUsers.Decode(mongoDoc)
		m := mongoDoc.Map()
		user := model.User{
			ID:        m["_id"].(primitive.ObjectID).Hex(),
			FirstName: m["firstName"].(string),
			LastName:  m["lastName"].(string),
			Weight:    m["weight"].(int32),
		}
		userModels = append(userModels, user)
	}
	response, _ := json.Marshal(model.UserList{Users: userModels})
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(response)
}

func addUser(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var newUser model.User
	err := decoder.Decode(&newUser)
	if err != nil {
		kaboom(writer, err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, _ := mongo.Connect(ctx, &options.ClientOptions{Hosts: []string{"localhost"}})
	_ = client.Ping(ctx, readpref.Primary())

	collection := client.Database("testing").Collection("UserProfiles")
	res, _ := collection.InsertOne(ctx, bson.M{"firstName": newUser.FirstName, "lastName": newUser.LastName, "weight": newUser.Weight})

	filter := bson.D{{"_id", res.InsertedID}}
	insertedUser := client.Database("testing").Collection("UserProfiles").FindOne(ctx, filter)
	mongoDoc := &bson.D{}
	err = insertedUser.Decode(mongoDoc)
	if err != nil {
		kaboom(writer, err)
		return
	}
	m := mongoDoc.Map()
	user := model.User{
		ID:        m["_id"].(primitive.ObjectID).Hex(),
		FirstName: m["firstName"].(string),
		LastName:  m["lastName"].(string),
		Weight:    m["weight"].(int32),
	}

	response, _ := json.Marshal(user)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(response)
}

func deleteUser(writer http.ResponseWriter, userId string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, _ := mongo.Connect(ctx, &options.ClientOptions{Hosts: []string{"localhost"}})
	_ = client.Ping(ctx, readpref.Primary())

	collection := client.Database("testing").Collection("UserProfiles")
	objectId, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.D{{"_id", objectId}}
	res, _ := collection.DeleteOne(ctx, filter)
	if res.DeletedCount == 0 {
		notFound(writer, "User ID not found:"+userId)
	}
}

func getOneUser(writer http.ResponseWriter, userId string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, _ := mongo.Connect(ctx, &options.ClientOptions{Hosts: []string{"localhost"}})
	_ = client.Ping(ctx, readpref.Primary())

	objectId, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.D{{"_id", objectId}}
	theUser := client.Database("testing").Collection("UserProfiles").FindOne(ctx, filter)
	mongoDoc := &bson.D{}
	err := theUser.Decode(mongoDoc)
	if err != nil {
		kaboom(writer, err)
		return
	}
	m := mongoDoc.Map()
	user := model.User{
		ID:        m["_id"].(primitive.ObjectID).Hex(),
		FirstName: m["firstName"].(string),
		LastName:  m["lastName"].(string),
		Weight:    m["weight"].(int32),
	}

	response, _ := json.Marshal(user)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(response)
}

func oneUser(writer http.ResponseWriter, request *http.Request) {
	r, _ := regexp.Compile("^/api/users/([a-f0-9]+)$")

	//possibleId := strings.TrimPrefix(request.URL.Path, "/api/users/")
	//r.FindString(request.URL.Path)
	//if strings.Contains(possibleId, "/") {
	if !r.MatchString(request.URL.Path) {
		notFound(writer, "Resource not found: "+request.URL.Path)
		return
	}
	switch request.Method {
	case http.MethodDelete:
		deleteUser(writer, r.FindStringSubmatch(request.URL.Path)[1])
	case http.MethodGet:
		getOneUser(writer, r.FindStringSubmatch(request.URL.Path)[1])
	default:
		notFound(writer, "i dunno how to do that yet")
	}
}

func kaboom(writer http.ResponseWriter, e error) {
	log.Println("Some shit blew up yo: " + e.Error())
	errorResponse := model.ErrorResponse{Status: 500, Message: "Something blew up yo: " + e.Error()}
	response, _ := json.Marshal(errorResponse)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(500)
	writer.Write(response)
}

func notFound(writer http.ResponseWriter, msg string) {
	errorResponse := model.ErrorResponse{Status: 404, Message: msg}
	response, _ := json.Marshal(errorResponse)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(404)
	writer.Write(response)
}

func (a *App) Initialize() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/api/users", users)
	serveMux.HandleFunc("/api/users/", oneUser)
	a.Handler = serveMux
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Handler))
}
