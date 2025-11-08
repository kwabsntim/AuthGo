package models

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"usernme" json:"username"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password"`
	CreatedAt time.Time `bson:"createdAt" json:"created_at"`
	LastLogin time.Time `bson:"lastLogin" json:"last_login"`
	UpdatedAt time.Time `bson:"updated" json:"updated_at"`
}

type Message struct {
	Message string `json:"message"`
}
