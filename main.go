package main

import (
	db "AuthGo/Database"
	"AuthGo/handlers"
	"AuthGo/middleware"
	"AuthGo/repository"
	"AuthGo/services"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: no .env file found")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := db.ConnectDB()
	//checking the database connection for nil pointer
	if client == nil {
		logger.Error("MongoDB client is nil")
		return
	}
	start := time.Now()
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Println("Ping failed:", err)
	}
	fmt.Println("Ping latency:", time.Since(start))
	// Test the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Database ping failed:", err)
	}

	// Create dependencies using DI
	userRepo := repository.NewMongoUserRepository(client)
	userService := services.NewUserService(userRepo)

	// Setup database indexes
	err = userRepo.SetupIndexes()
	if err != nil {
		logger.Error("Failed to setup indexes", "error", err)
	}

	//disconnecting the mongoDB client when the main function ends
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {

			logger.Error("Error disconnecting from mongoDb", "error", err)
		}
	}()
	//using a server mux to map the requests to the handlers
	mux := handlers.RouteSetup(userService)
	handlerWithPanicRecovery := middleware.PanicMiddleware(logger)(mux)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback
	}

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        handlerWithPanicRecovery,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	logger.Info("Server started on port " + port)
	go func() {

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server has been shutdown...")
	}
	logger.Info("Server exited...")
}
