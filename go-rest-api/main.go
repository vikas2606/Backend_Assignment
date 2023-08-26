package main

import (
	"example/go-rest-api/db"
	"example/go-rest-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize the database connection
	db.ConnectDB()
	router := gin.Default()
	router.POST("/student/login", routes.StudentLogin)
	router.POST("/dean/login", routes.DeanLogin)
	router.GET("/sessions/available", routes.GetAvailableSessions)
	router.POST("/sessions/book/:session_id", routes.BookSessionSlot)

	// Define your API routes and handlers here

	router.Run(":8080")
}
