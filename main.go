package main

import (
	db "AuthGo/Database"
	"AuthGo/handlers"
	"AuthGo/middleware"
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
	//

	//logger using slog to log in json format
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {

			logger.Error("Error disconnecting from mongoDb", "error", err)
		}
	}()
	//using a server mux to map the requests to the handlers
	mux := handlers.RouteSetup()
	//adding middleware for panic recovery
	handlerforPanicRecovery := middleware.PanicMiddleware(logger)(mux)
	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerforPanicRecovery,
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
