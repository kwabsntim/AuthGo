package repository

import (
	"AuthGo/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client

func CreateUser(user *models.User) error {

	user.CreatedAt = time.Now()

	collection := Client.Database("usersdb").Collection("profiles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//checking if the user already exists
	existing := collection.FindOne(ctx, bson.M{"email": user.Email})
	if existing.Err() == nil {
		return errors.New("user already exists")
	}

	//creating user
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return errors.New("failed to create user")
	}

	// Set the generated ID back to the user
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}
	return nil
}
