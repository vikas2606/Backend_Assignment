package routes

import (
	"example/go-rest-api/controller"
	"example/go-rest-api/db"
	"example/go-rest-api/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	tokenTTL = 3600 // Token TTL in seconds (1 hour)
)

func StudentLogin(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var student model.Student
	err := db.StudentCollection.FindOne(db.Context, bson.M{"username": loginData.Username, "password": loginData.Password}).Decode(&student)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a new UUID token
	uuidToken := uuid.New().String()

	// Calculate token expiration time
	expiresAt := time.Now().Add(time.Second * tokenTTL)

	// Store token in the database
	tokenDoc := model.AuthToken{
		StudentID: student.ID, // Assuming student.ID is an ObjectID
		Token:     uuidToken,
		ExpiresAt: expiresAt,
	}

	_, err = db.AuthTokenCollection.InsertOne(db.Context, tokenDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student logged in", "student_uuid": uuidToken})
}

func DeanLogin(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var dean model.Dean
	err := db.DeanCollection.FindOne(db.Context, bson.M{"username": loginData.Username, "password": loginData.Password}).Decode(&dean)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a new UUID token
	uuidToken := uuid.New().String()

	// Calculate token expiration time
	expiresAt := time.Now().Add(time.Second * tokenTTL)

	// Store token in the database
	tokenDoc := model.AuthToken{
		DeanID:    dean.ID, // Assuming student.ID is an ObjectID
		Token:     uuidToken,
		ExpiresAt: expiresAt,
	}

	_, err = db.AuthTokenCollection.InsertOne(db.Context, tokenDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dean logged in", "dean_uuid": uuidToken})
}

func GetAvailableSessions(c *gin.Context) {
	auth, err := controller.ValidateToken(c)
	if err != nil {
		return
	}

	var sessions []model.Session

	filter := bson.M{}

	if auth.StudentID != primitive.NilObjectID {
		filter["status"] = "available"
	} else if auth.DeanID != primitive.NilObjectID {
		filter["status"] = "pending"

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}
	cursor, err := db.SessionCollection.Find(db.Context, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
		return
	}
	defer cursor.Close(db.Context)

	if err := cursor.All(db.Context, &sessions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})

}

func BookSessionSlot(c *gin.Context) {
	auth, err := controller.ValidateToken(c)
	if err != nil {
		return
	}

	sessionIDStr := c.Param("session_id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var session model.Session
	err = db.SessionCollection.FindOne(db.Context, bson.M{"session_id": sessionID, "status": "available"}).Decode(&session)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found or not available"})
		return
	}

	// Fetch the student's username based on their ID
	var student model.Student
	err = db.StudentCollection.FindOne(db.Context, bson.M{"_id": auth.StudentID}).Decode(&student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student details"})
		return
	}

	// Update the session status to "pending" and assign the student ID
	update := bson.M{"$set": bson.M{"status": "pending", "student_id": auth.StudentID, "name": student.Username}}
	_, err = db.SessionCollection.UpdateOne(db.Context, bson.M{"_id": session.ID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Session slot booked successfully",
		"student_id": auth.StudentID.Hex(),
		"username":   student.Username,
	})
}
