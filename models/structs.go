package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
	CreatedAt time.Time          `bson:"createdAt" json:"created_at"`
	LastLogin time.Time          `bson:"lastLogin" json:"last_login"`
	UpdatedAt time.Time          `bson:"updated" json:"updated_at"`
}
type JSONresponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
