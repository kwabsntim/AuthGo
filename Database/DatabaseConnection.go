package db

import (
	"github.com/joho/godotenv"

	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectDB() *mongo.Client {
	err := godotenv.Load()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI).SetConnectTimeout(10 * time.Second)
	if err != nil {
		log.Println("Warning: no .env file found")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not found in environment")
	}

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("MongoDB ping error:", err)
	}

	log.Println(" Connected to MongoDB successfully!")
	return client
}
