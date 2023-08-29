//model.go
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Type     string             `bson:"type"` // Dean or student
}

type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	SessionID    int32              `bson:"session_id"`
	DeanID       primitive.ObjectID `bson:"dean_id"`
	Dean_Name    string             `bson:"dean_name"`
	StudentID    primitive.ObjectID `bson:"student_id"`
	Student_Name string             `bson:"student_name"`
	Slot         time.Time          `bson:"slot"`
	Status       string             `bson:"status"` // "pending" or "completed" or "available"
}
