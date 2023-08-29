package routes

import (
	"example/go-rest-api/controller"
	"example/go-rest-api/db"
	"example/go-rest-api/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

	var student model.User
	err := db.UserCollection.FindOne(db.Context, bson.M{"username": loginData.Username, "password": loginData.Password, "type": "student"}).Decode(&student)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a new UUID token
	uuidToken := uuid.New().String()

	// Calculate token expiration time
	expiresAt := time.Now().Add(time.Second * tokenTTL)

	cookie := http.Cookie{
		Name:     "student_token",
		Value:    uuidToken,
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	controller.ActiveTokens = append(controller.ActiveTokens, uuidToken)
	controller.ActiveTokenDetails[uuidToken] = controller.TokenDetails{
		ID:       student.ID,
		Username: student.Username,
	}

	http.SetCookie(c.Writer, &cookie)
	fmt.Println(controller.ActiveTokenDetails)
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

	var dean model.User
	err := db.UserCollection.FindOne(db.Context, bson.M{"username": loginData.Username, "password": loginData.Password, "type": "dean"}).Decode(&dean)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a new UUID token
	uuidToken := uuid.New().String()

	// Calculate token expiration time
	expiresAt := time.Now().Add(time.Second * tokenTTL)

	// Store token in the database
	cookie := http.Cookie{
		Name:     "dean_token",
		Value:    uuidToken,
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	controller.ActiveTokens = append(controller.ActiveTokens, uuidToken)
	controller.ActiveTokenDetails[uuidToken] = controller.TokenDetails{
		ID:       dean.ID,
		Username: dean.Username,
	}

	http.SetCookie(c.Writer, &cookie)
	c.JSON(http.StatusOK, gin.H{"message": "Dean logged in"})
}

func GetAvailableSessions(c *gin.Context) {
	_, err := controller.ValidateToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var sessions []model.Session

	filter := bson.M{
		"status": "available",
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

func GetPendingSessions(c *gin.Context) {
	uuidToken, err := controller.ValidateToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenDetails, found := controller.ActiveTokenDetails[uuidToken]
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var sessions []model.Session

	filter := bson.M{
		"status":    "pending",
		"dean_id":   tokenDetails.ID,
		"dean_name": tokenDetails.Username,
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
	uuidToken, err := controller.ValidateToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenDetails, found := controller.ActiveTokenDetails[uuidToken]
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Update the session fields with token details
	update := bson.M{
		"$set": bson.M{
			"status":       "pending",
			"student_id":   tokenDetails.ID,
			"student_name": tokenDetails.Username,
		},
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

	_, err = db.SessionCollection.UpdateOne(db.Context, bson.M{"_id": session.ID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session slot booked successfully",
	})
}
