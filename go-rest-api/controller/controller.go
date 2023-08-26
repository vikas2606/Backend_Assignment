package controller

import (
	"errors"
	"example/go-rest-api/db"
	"example/go-rest-api/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateToken(c *gin.Context) (*model.AuthToken, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return nil, errors.New("Authorization header missing")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return nil, errors.New("Invalid authorization header format")
	}

	uuidToken := authHeaderParts[1]

	token, err := uuid.Parse(uuidToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
		return nil, errors.New("Invalid token format")
	}

	var auth model.AuthToken
	err = db.AuthTokenCollection.FindOne(db.Context, bson.M{"token": token.String()}).Decode(&auth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return nil, errors.New("Invalid token")
	}

	return &auth, nil
}
