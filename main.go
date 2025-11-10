package main

import (
	db "AuthGo/Database"
	"AuthGo/handlers"
	"AuthGo/middleware"
	"AuthGo/repository"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := db.ConnectDB()
	//checking the database connection for nil pointer
	if client == nil {
		logger.Error("MongoDB client is nil")
		return
	}
	repository.Client = client

	// Test the connection
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Database ping failed:", err)
	}

	//disconnecting the mongoDB client when the main function ends
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {

			logger.Error("Error disconnecting from mongoDb", "error", err)
		}
	}()
	//using a server mux to map the requests to the handlers
	mux := handlers.RouteSetup()
	handlerWithPanicRecovery := middleware.PanicMiddleware(logger)(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithPanicRecovery, // Use middleware
	}

	logger.Info("Server started on port 8080")
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
