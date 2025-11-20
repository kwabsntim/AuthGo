package main

import (
	db "AuthGo/Database"
	"AuthGo/handlers"
	"AuthGo/middleware"
	"AuthGo/repository"
	"AuthGo/services"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Load environment variables (optional - for local development)
	godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Retry MongoDB connection with exponential backoff
	var client *mongo.Client
	for i := 0; i < 10; i++ {
		client = db.ConnectDB()
		if client != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := client.Ping(ctx, nil); err == nil {
				cancel()
				logger.Info("Connected to MongoDB successfully")
				break
			}
			cancel()
		}
		waitTime := time.Duration(i+1) * time.Second
		logger.Info("MongoDB connection failed, retrying...", "attempt", i+1, "wait", waitTime)
		time.Sleep(waitTime)
	}

	if client == nil {
		logger.Error("Failed to connect to MongoDB after retries")
		return
	}

	// Create dependencies using DI
	userRepo := repository.NewMongoUserRepository(client)

	//service container
	serviceContainer := &handlers.ServiceContainer{
		RegisterService: services.NewRegisterService(userRepo),
		LoginService:    services.NewLoginService(userRepo),
		FetchService:    services.NewFetchUserService(userRepo),
	}
	// Setup routes with container
	mux := handlers.RouteSetup(serviceContainer)
	handlerWithPanicRecovery := middleware.PanicMiddleware(logger)(mux)

	// Setup database indexes
	err := userRepo.SetupIndexes()
	if err != nil {
		logger.Error("Failed to setup indexes", "error", err)
		return
	}

	//disconnecting the mongoDB client when the main function ends
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			logger.Error("Error disconnecting from mongoDb", "error", err)
		}
	}()

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

	quit := make(chan os.Signal, 1)
	logger.Info("Server started on port " + port)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			quit <- syscall.SIGTERM
		}
	}()
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
