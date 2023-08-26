// db.go

package db

import (
	"context"
	"example/go-rest-api/model"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var client *mongo.Client
var Context context.Context
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

	StudentCollection = client.Database("university").Collection("students")
	DeanCollection = client.Database("university").Collection("deans")
	SessionCollection = client.Database("university").Collection("sessions")
	AuthTokenCollection = client.Database("university").Collection("auth_tokens") // Add this line

	printStudents()
}

func printStudents() {
	cursor, err := StudentCollection.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Println("Error fetching students:", err)
		return
	}
	defer cursor.Close(context.Background())

	var students []model.Student
	if err := cursor.All(context.Background(), &students); err != nil {
		fmt.Println("Error decoding students:", err)
		return
	}

	fmt.Println("List of students:")
	for _, student := range students {
		fmt.Printf("ID: %s, Username: %s\n", student.ID, student.Username)
	}
}
