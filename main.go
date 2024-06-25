package main

import (
	"context"
	"errors"
	"fmt"
	"httpserver/config"
	"httpserver/internal/handlers"
	"httpserver/internal/middleware"
	"httpserver/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	r := router.NewRouter()

	r.AddRoute("GET", "/", middleware.Logging(handlers.GetRoot))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}
	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
