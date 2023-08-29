package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ActiveTokens []string

type TokenDetails struct {
	ID       primitive.ObjectID
	Username string
}

var ActiveTokenDetails map[string]TokenDetails

func init() {
	ActiveTokenDetails = make(map[string]TokenDetails)
}

func ValidateToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return "", errors.New("Authorization header missing")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return "", errors.New("Invalid authorization header format")
	}

	uuidToken := authHeaderParts[1]

	_, err := uuid.Parse(uuidToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
		return "", errors.New("Invalid token format")
	}

	found := false
	for _, token := range ActiveTokens {
		if token == uuidToken {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return "", errors.New("Invalid token")
	}
	fmt.Println(uuidToken)
	return uuidToken, nil
}
