//model.go
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

type Dean struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	SessionID int32              `bson:"session_id"`
	StudentID primitive.ObjectID `bson:"student_id"`
	Name      string             `bson:"name"`
	Slot      time.Time          `bson:"slot"`
	Status    string             `bson:"status"` // "pending" or "completed" or "available"
}

type AuthToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	StudentID primitive.ObjectID `bson:"student_id"`
	DeanID    primitive.ObjectID `bson:"dean_id"`
	Token     string             `bson:"token"`
	ExpiresAt time.Time          `bson:"expires_at"`
}
