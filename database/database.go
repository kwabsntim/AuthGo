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
	// Try to load .env file (optional - for local development)
	godotenv.Load()
	
	//specifying the server api version
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	//applying database pooling options
	opts := options.Client().
		ApplyURI(os.Getenv("MONGODB_URI")).
		SetServerAPIOptions(serverAPI).
		SetMaxPoolSize(20).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(60 * time.Second).
		SetConnectTimeout(5 * time.Second).
		SetSocketTimeout(15 * time.Second).
		SetServerSelectionTimeout(5 * time.Second).
		SetRetryWrites(true).
		SetRetryReads(true)

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
