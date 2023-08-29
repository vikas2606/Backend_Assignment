// db.go

package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var Context context.Context
var UserCollection *mongo.Collection
var StudentCollection *mongo.Collection
var SessionCollection *mongo.Collection
var DeanCollection *mongo.Collection
var AuthTokenCollection *mongo.Collection // Add this line

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		panic(err) // Stop execution if unable to connect
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		panic(err) // Stop execution if unable to ping
	}

	fmt.Println("Connected to the database!")

	UserCollection = client.Database("university").Collection("users")
	StudentCollection = client.Database("university").Collection("students")
	DeanCollection = client.Database("university").Collection("deans")
	SessionCollection = client.Database("university").Collection("sessions")
	AuthTokenCollection = client.Database("university").Collection("auth_tokens") // Add this line

}
